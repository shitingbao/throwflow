package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	ctos "github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"io"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"weixin/internal/conf"
	"weixin/internal/domain"
	"weixin/internal/pkg/mini/oauth2"
	"weixin/internal/pkg/mini/wxa"
	"weixin/internal/pkg/tool"
)

var (
	WeixinUserCouponNotFound      = errors.NotFound("WEIXIN_USER_COUPON_NOT_FOUND", "微信用户券码不存在")
	WeixinUserCouponListError     = errors.InternalServer("WEIXIN_USER_COUPON_LIST_ERROR", "微信用户券码列表获取失败")
	WeixinUserCouponBindError     = errors.InternalServer("WEIXIN_USER_COUPON_BIND_ERROR", "微信用户券码绑定失败")
	WeixinUserCouponActivateError = errors.InternalServer("WEIXIN_USER_COUPON_ACTIVATE_ERROR", "微信用户券码激活失败")
	WeixinUserCouponCreateError   = errors.InternalServer("WEIXIN_USER_COUPON_CREATE_ERROR", "微信用户券码创建失败")
)

type UserCouponRepo interface {
	Get(context.Context, uint64, string, string) (*domain.UserCoupon, error)
	GetByPhone(context.Context, uint64, string, string, string) (*domain.UserCoupon, error)
	List(context.Context, int, int, uint64, uint64, string) ([]*domain.UserCoupon, error)
	Count(context.Context, uint64, uint64, string) (int64, error)
	Save(context.Context, *domain.UserCoupon) (*domain.UserCoupon, error)
	Update(context.Context, *domain.UserCoupon) (*domain.UserCoupon, error)

	SaveCacheString(context.Context, string, string, time.Duration) (bool, error)
	DeleteCache(context.Context, string) error

	PutContent(context.Context, string, io.Reader) (*ctos.PutObjectV2Output, error)
}

type UserCouponUsecase struct {
	repo     UserCouponRepo
	ucclrepo UserCouponCreateLogRepo
	uorrepo  UserOrganizationRelationRepo
	uirrepo  UserIntegralRelationRepo
	ucrepo   UserCommissionRepo
	urepo    UserRepo
	usrrepo  UserScanRecordRepo
	corepo   CompanyOrganizationRepo
	screpo   ShortCodeRepo
	tm       Transaction
	conf     *conf.Data
	oconf    *conf.Organization
	wconf    *conf.Weixin
	vconf    *conf.Volcengine
	log      *log.Helper
}

func NewUserCouponUsecase(repo UserCouponRepo, ucclrepo UserCouponCreateLogRepo, uorrepo UserOrganizationRelationRepo, uirrepo UserIntegralRelationRepo, ucrepo UserCommissionRepo, urepo UserRepo, usrrepo UserScanRecordRepo, corepo CompanyOrganizationRepo, screpo ShortCodeRepo, tm Transaction, conf *conf.Data, oconf *conf.Organization, wconf *conf.Weixin, vconf *conf.Volcengine, logger log.Logger) *UserCouponUsecase {
	return &UserCouponUsecase{repo: repo, ucclrepo: ucclrepo, uorrepo: uorrepo, uirrepo: uirrepo, ucrepo: ucrepo, urepo: urepo, usrrepo: usrrepo, corepo: corepo, screpo: screpo, tm: tm, conf: conf, oconf: oconf, wconf: wconf, vconf: vconf, log: log.NewHelper(logger)}
}

func (ucuc *UserCouponUsecase) GetUserCoupons(ctx context.Context, userId, organizationId uint64) (*domain.UserCoupon, error) {
	user, err := ucuc.urepo.Get(ctx, userId)

	if err != nil {
		return nil, WeixinLoginError
	}

	companyOrganization, err := ucuc.corepo.Get(ctx, organizationId)

	if err != nil {
		return nil, WeixinCompanyOrganizationNotFound
	}

	if _, err := ucuc.uorrepo.GetByUserId(ctx, user.Id, companyOrganization.Data.OrganizationId, 0, "0"); err == nil {
		return nil, WeixinUserOrganizationRelationExist
	}

	userCoupon, err := ucuc.repo.GetByPhone(ctx, organizationId, user.Phone, "2", "2")

	if err != nil {
		return nil, WeixinUserCouponNotFound
	}

	return userCoupon, nil
}

func (ucuc *UserCouponUsecase) ListUserCoupons(ctx context.Context, pageNum, pageSize, userId, organizationId uint64) (*domain.UserCouponList, error) {
	if _, err := ucuc.urepo.Get(ctx, userId); err != nil {
		return nil, WeixinLoginError
	}

	if _, err := ucuc.corepo.Get(ctx, organizationId); err != nil {
		return nil, WeixinCompanyOrganizationNotFound
	}

	list := make([]*domain.UserCoupon, 0)

	userCoupons, err := ucuc.repo.List(ctx, int(pageNum), int(pageSize), userId, organizationId, "")

	if err != nil {
		return nil, WeixinUserCouponListError
	}

	for _, userCoupon := range userCoupons {
		userCoupon.SetContent(ctx)

		list = append(list, &domain.UserCoupon{
			Id:               userCoupon.Id,
			UserId:           userCoupon.UserId,
			CouponCode:       userCoupon.CouponCode,
			Level:            userCoupon.Level,
			Phone:            userCoupon.Phone,
			UserCouponStatus: userCoupon.UserCouponStatus,
			OrganizationId:   userCoupon.OrganizationId,
			Content:          userCoupon.Content,
			CreateTime:       userCoupon.CreateTime,
			UpdateTime:       userCoupon.UpdateTime,
		})
	}

	total, err := ucuc.repo.Count(ctx, userId, organizationId, "")

	if err != nil {
		return nil, WeixinUserCouponListError
	}

	notUsedTotal, err := ucuc.repo.Count(ctx, userId, organizationId, "1")

	if err != nil {
		return nil, WeixinUserCouponListError
	}

	return &domain.UserCouponList{
		PageNum:   pageNum,
		PageSize:  pageSize,
		Total:     uint64(total),
		TotalUsed: uint64(total - notUsedTotal),
		List:      list,
	}, nil
}

func (ucuc *UserCouponUsecase) BindUserCoupons(ctx context.Context, userId uint64, phone string) error {
	user, err := ucuc.urepo.Get(ctx, userId)

	if err != nil {
		return WeixinLoginError
	}

	isFail := true

	for num := 0; num <= 1; num++ {
		result, err := ucuc.repo.SaveCacheString(ctx, "weixin:user:coupon:"+strconv.FormatUint(user.Id, 10), "go", ucuc.conf.Redis.UserCouponLockTimeout.AsDuration())

		if err != nil {
			return WeixinUserCouponBindError
		}

		if result {
			err := ucuc.tm.InTx(ctx, func(ctx context.Context) error {
				if _, err := ucuc.repo.GetByPhone(ctx, 0, phone, "2", ""); err == nil {
					return WeixinUserCouponBindError
				}

				inUserCoupon, err := ucuc.repo.Get(ctx, user.Id, "2", "1")

				if err != nil {
					return err
				}

				inUserCoupon.SetPhone(ctx, phone)
				inUserCoupon.SetUserCouponStatus(ctx, 2)
				inUserCoupon.SetUpdateTime(ctx)

				if _, err := ucuc.repo.Update(ctx, inUserCoupon); err != nil {
					return err
				}

				return nil
			})

			ucuc.repo.DeleteCache(ctx, "weixin:user:coupon:"+strconv.FormatUint(user.Id, 10))

			if err != nil {
				return WeixinUserCouponBindError
			}

			isFail = false

			break
		} else {
			time.Sleep(200 * time.Millisecond)

			continue
		}
	}

	if isFail {
		return WeixinUserCouponBindError
	}

	return nil
}

func (ucuc *UserCouponUsecase) ActivateUserCoupons(ctx context.Context, userId, parentUserId, organizationId uint64) error {
	user, err := ucuc.urepo.Get(ctx, userId)

	if err != nil {
		return WeixinLoginError
	}

	companyOrganization, err := ucuc.corepo.Get(ctx, organizationId)

	if err != nil {
		return WeixinCompanyOrganizationNotFound
	}

	inUserCoupon, err := ucuc.repo.GetByPhone(ctx, organizationId, user.Phone, "2", "2")

	if err != nil {
		return WeixinUserCouponNotFound
	}

	organizationCourses := make([]domain.OrganizationCourse, 0)

	for _, organizationCourse := range companyOrganization.Data.OrganizationCourses {
		courseModules := make([]domain.OrganizationCourseModule, 0)

		for _, courseModule := range organizationCourse.CourseModules {
			courseModules = append(courseModules, domain.OrganizationCourseModule{
				CourseModuleName:    courseModule.CourseModuleName,
				CourseModuleContent: courseModule.CourseModuleContent,
			})
		}

		organizationCourses = append(organizationCourses, domain.OrganizationCourse{
			CourseName:     organizationCourse.CourseName,
			CourseSubName:  organizationCourse.CourseSubName,
			CoursePrice:    organizationCourse.CoursePrice,
			CourseDuration: organizationCourse.CourseDuration,
			CourseLevel:    uint8(organizationCourse.CourseLevel),
			CourseModules:  courseModules,
		})
	}

	sort.Sort(domain.OrganizationCourses(organizationCourses))

	if _, err := ucuc.uorrepo.GetByUserId(ctx, user.Id, organizationId, 0, "0"); err == nil {
		return WeixinUserOrganizationRelationExist
	}

	var organizationUserId uint64 = 0
	var organizationTutorId uint64 = 0
	var parentUserOrganizationRelation *domain.UserOrganizationRelation
	var tutorUserIntegralRelation *domain.UserIntegralRelation

	if parentUserId > 0 {
		parentUser, err := ucuc.urepo.Get(ctx, parentUserId)

		if err != nil {
			return WeixinUserNotFound
		}

		parentUserOrganizationRelation, err = ucuc.uorrepo.GetByUserId(ctx, parentUser.Id, organizationId, 0, "0")

		if err != nil {
			return WeixinUserOrganizationRelationNotFound
		}

		organizationUserId = parentUserOrganizationRelation.UserId

		if parentUserOrganizationRelation.Level == 4 {
			organizationTutorId = parentUserOrganizationRelation.UserId

			tutorUserIntegralRelation = &domain.UserIntegralRelation{
				UserId:             parentUserOrganizationRelation.UserId,
				OrganizationId:     parentUserOrganizationRelation.OrganizationId,
				OrganizationUserId: parentUserOrganizationRelation.OrganizationUserId,
				Level:              parentUserOrganizationRelation.Level,
			}
		} else {
			userIntegralRelations, err := ucuc.uirrepo.List(ctx, organizationId)

			if err != nil {
				return WeixinUserOrganizationRelationListError
			}

			tutorUserIntegralRelation = ucuc.uirrepo.GetSuperior(ctx, parentUserOrganizationRelation.UserId, 4, userIntegralRelations)

			if tutorUserIntegralRelation != nil {
				organizationTutorId = tutorUserIntegralRelation.UserId
			}
		}
	}

	err = ucuc.tm.InTx(ctx, func(ctx context.Context) error {
		inUserCoupon.SetUserCouponStatus(ctx, 3)
		inUserCoupon.SetUpdateTime(ctx)

		if _, err = ucuc.repo.Update(ctx, inUserCoupon); err != nil {
			return err
		}

		var wcontent *wxa.GetUnlimitedQRCodeResponse

		if inUserCoupon.OrganizationId == ucuc.oconf.DjOrganizationId {
			accessToken, err := oauth2.GetAccessToken(ucuc.wconf.DjMini.Appid, ucuc.wconf.DjMini.Secret)

			if err != nil {
				return err
			}

			wcontent, err = wxa.GetUnlimitedQRCode(accessToken.AccessToken, "oId="+strconv.FormatUint(inUserCoupon.OrganizationId, 10)+"&uId="+strconv.FormatUint(user.Id, 10), ucuc.wconf.DjMini.QrCodeEnvVersion)

			if err != nil {
				return err
			}
		} else {
			accessToken, err := oauth2.GetAccessToken(ucuc.wconf.Mini.Appid, ucuc.wconf.Mini.Secret)

			if err != nil {
				return err
			}

			wcontent, err = wxa.GetUnlimitedQRCode(accessToken.AccessToken, "oId="+strconv.FormatUint(inUserCoupon.OrganizationId, 10)+"&uId="+strconv.FormatUint(user.Id, 10), ucuc.wconf.Mini.QrCodeEnvVersion)

			if err != nil {
				return err
			}
		}

		objectKey := tool.GetRandCode(time.Now().String())

		if _, err := ucuc.repo.PutContent(ctx, ucuc.vconf.Tos.Company.SubFolder+"/"+objectKey+".png", strings.NewReader(wcontent.Buffer)); err != nil {
			return err
		}

		inUserOrganizationRelation := domain.NewUserOrganizationRelation(ctx, user.Id, inUserCoupon.OrganizationId, organizationUserId, organizationTutorId, inUserCoupon.Level, 0, ucuc.vconf.Tos.Company.Url+"/"+ucuc.vconf.Tos.Company.SubFolder+"/"+objectKey+".png")
		inUserOrganizationRelation.SetCreateTime(ctx)
		inUserOrganizationRelation.SetUpdateTime(ctx)

		if _, err := ucuc.uorrepo.Save(ctx, inUserOrganizationRelation); err != nil {
			return err
		}

		inUserIntegralRelation := domain.NewUserIntegralRelation(ctx, user.Id, inUserCoupon.OrganizationId, organizationUserId, inUserCoupon.Level)
		inUserIntegralRelation.SetCreateTime(ctx)
		inUserIntegralRelation.SetUpdateTime(ctx)

		if _, err := ucuc.uirrepo.Save(ctx, inUserIntegralRelation); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return WeixinUserCouponActivateError
	}

	return nil
}

func (ucuc *UserCouponUsecase) CreateUserCoupons(ctx context.Context, userId, organizationId uint64, num uint32, level uint8) error {
	if _, err := ucuc.urepo.Get(ctx, userId); err != nil {
		return WeixinUserNotFound
	}

	companyOrganization, err := ucuc.corepo.Get(ctx, organizationId)

	if err != nil {
		return WeixinCompanyOrganizationNotFound
	}

	inUserCouponCreateLog := domain.NewUserCouponCreateLog(ctx, userId, companyOrganization.Data.OrganizationId, num, level, 0)
	inUserCouponCreateLog.SetCreateTime(ctx)
	inUserCouponCreateLog.SetUpdateTime(ctx)

	if _, err := ucuc.ucclrepo.Save(ctx, inUserCouponCreateLog); err != nil {
		return WeixinUserCouponCreateError
	}

	return nil
}

func (ucuc *UserCouponUsecase) SyncUserCoupons(ctx context.Context) error {
	userCouponCreateLogs, err := ucuc.ucclrepo.List(ctx, "0")

	if err != nil {
		return WeixinUserCouponCreateLogNotFound
	}

	var wg sync.WaitGroup

	for _, inUserCouponCreateLog := range userCouponCreateLogs {
		inUserCouponCreateLog.SetIsHandle(ctx, 1)
		inUserCouponCreateLog.SetUpdateTime(ctx)

		if _, err := ucuc.ucclrepo.Update(ctx, inUserCouponCreateLog); err != nil {
			return WeixinUserCouponCreateLogUpdateError
		}

		wg.Add(1)

		go ucuc.SyncUserCoupon(ctx, &wg, inUserCouponCreateLog)
	}

	wg.Wait()

	userCoupons, err := ucuc.repo.List(ctx, 0, 0, 0, 0, "2")

	if err != nil {
		return WeixinUserCouponListError
	}

	for _, inUserCoupon := range userCoupons {
		if time.Now().After(inUserCoupon.UpdateTime.AddDate(0, 0, 1)) {
			inUserCoupon.SetPhone(ctx, "")
			inUserCoupon.SetUserCouponStatus(ctx, 1)
			inUserCoupon.SetUpdateTime(ctx)

			ucuc.repo.Update(ctx, inUserCoupon)
		}
	}

	return nil
}

func (ucuc *UserCouponUsecase) SyncUserCoupon(ctx context.Context, wg *sync.WaitGroup, userCouponCreateLog *domain.UserCouponCreateLog) {
	defer func() {
		ucuc.ucclrepo.Delete(ctx, userCouponCreateLog)

		wg.Done()
	}()

	num := userCouponCreateLog.Num

	for num > 0 {
		if shortCode, err := ucuc.screpo.Create(ctx); err == nil {
			inUserCoupon := domain.NewUserCoupon(ctx, userCouponCreateLog.UserId, userCouponCreateLog.OrganizationId, userCouponCreateLog.Level, 1, shortCode.Data.ShortCode, "")
			inUserCoupon.SetCreateTime(ctx)
			inUserCoupon.SetUpdateTime(ctx)

			if _, err := ucuc.repo.Save(ctx, inUserCoupon); err == nil {
				num--
			}
		}
	}
}

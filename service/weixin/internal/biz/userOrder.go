package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	ctos "github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
	"weixin/internal/conf"
	"weixin/internal/domain"
	"weixin/internal/pkg/mini/oauth2"
	"weixin/internal/pkg/mini/wxa"
	"weixin/internal/pkg/tool"
)

var (
	WeixinUserOrderNotFound       = errors.NotFound("WEIXIN_USER_ORDER_NOT_FOUND", "微信用户订单不存在")
	WeixinUserOrderListError      = errors.InternalServer("WEIXIN_USER_ORDER_LIST_ERROR", "微信用户订单列表获取失败")
	WeixinUserOrderCreateError    = errors.InternalServer("WEIXIN_USER_ORDER_CREATE_ERROR", "微信用户订单创建失败")
	WeixinUserOrderCourseNotExist = errors.InternalServer("WEIXIN_USER_ORDER_COURSE_NOT_EXIST", "微信用户订单套餐不存在")
	WeixinUserOrderParentNotExist = errors.InternalServer("WEIXIN_USER_ORDER_PARENT_NOT_EXIST", "暂无推荐人，无法激活会员请点击老会员的邀请二维码")
)

type UserOrderRepo interface {
	NextId(context.Context) (uint64, error)
	GetByOutTradeNo(context.Context, string) (*domain.UserOrder, error)
	GetByUserId(context.Context, uint64, uint64, string) (*domain.UserOrder, error)
	List(context.Context, int, int) ([]*domain.UserOrder, error)
	ListOperation(context.Context) ([]*domain.UserOrder, error)
	Count(context.Context) (int64, error)
	StatisticsPayAmount(context.Context, uint64, uint64, string) (*domain.UserOrder, error)
	Update(context.Context, *domain.UserOrder) (*domain.UserOrder, error)
	Save(context.Context, *domain.UserOrder) (*domain.UserOrder, error)

	PutContent(context.Context, string, io.Reader) (*ctos.PutObjectV2Output, error)
}

type UserOrderUsecase struct {
	repo     UserOrderRepo
	urepo    UserRepo
	uorrepo  UserOrganizationRelationRepo
	uirrepo  UserIntegralRelationRepo
	usrrepo  UserScanRecordRepo
	ucrepo   UserCommissionRepo
	uoirepo  UserOpenIdRepo
	ublrepo  UserBalanceLogRepo
	ucclrepo UserCouponCreateLogRepo
	crepo    CompanyRepo
	corepo   CompanyOrganizationRepo
	prepo    PayRepo
	qcrepo   QrCodeRepo
	surepo   ShortUrlRepo
	screpo   ShortCodeRepo
	tm       Transaction
	conf     *conf.Data
	oconf    *conf.Organization
	wconf    *conf.Weixin
	vconf    *conf.Volcengine
	log      *log.Helper
}

func NewUserOrderUsecase(repo UserOrderRepo, urepo UserRepo, uorrepo UserOrganizationRelationRepo, uirrepo UserIntegralRelationRepo, usrrepo UserScanRecordRepo, ucrepo UserCommissionRepo, uoirepo UserOpenIdRepo, ublrepo UserBalanceLogRepo, ucclrepo UserCouponCreateLogRepo, crepo CompanyRepo, corepo CompanyOrganizationRepo, prepo PayRepo, qcrepo QrCodeRepo, surepo ShortUrlRepo, screpo ShortCodeRepo, tm Transaction, conf *conf.Data, oconf *conf.Organization, wconf *conf.Weixin, vconf *conf.Volcengine, logger log.Logger) *UserOrderUsecase {
	return &UserOrderUsecase{repo: repo, urepo: urepo, uorrepo: uorrepo, uirrepo: uirrepo, usrrepo: usrrepo, ucrepo: ucrepo, uoirepo: uoirepo, ublrepo: ublrepo, ucclrepo: ucclrepo, crepo: crepo, corepo: corepo, prepo: prepo, qcrepo: qcrepo, surepo: surepo, screpo: screpo, tm: tm, conf: conf, oconf: oconf, wconf: wconf, vconf: vconf, log: log.NewHelper(logger)}
}

func (uouc *UserOrderUsecase) ListUserOrders(ctx context.Context, pageNum, pageSize uint64) (*domain.UserOrderList, error) {
	list, err := uouc.repo.List(ctx, int(pageNum), int(pageSize))

	if err != nil {
		return nil, WeixinUserOrderListError
	}

	total, err := uouc.repo.Count(ctx)

	if err != nil {
		return nil, WeixinUserOrderListError
	}

	for _, l := range list {
		if len(l.NickName) == 0 {
			l.NickName = tool.FormatPhone(l.Phone)
		}

		if len(l.AvatarUrl) == 0 {
			l.AvatarUrl = uouc.vconf.Tos.Avatar.Url + "/" + uouc.vconf.Tos.Avatar.SubFolder + "/" + uouc.vconf.Tos.Avatar.DefaultAvatar
		}
	}

	return &domain.UserOrderList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (uouc *UserOrderUsecase) CreateUserOrders(ctx context.Context, userId, parentUserId, organizationId uint64, payAmount float64, clientIp string) (*domain.UserOrderRequestPayment, error) {
	user, err := uouc.urepo.Get(ctx, userId)

	if err != nil {
		return nil, WeixinLoginError
	}

	companyOrganization, err := uouc.corepo.Get(ctx, organizationId)

	if err != nil {
		return nil, WeixinCompanyOrganizationNotFound
	}

	if _, err := uouc.uorrepo.GetByUserId(ctx, user.Id, companyOrganization.Data.OrganizationId, 0, "0"); err == nil {
		return nil, WeixinUserOrganizationRelationExist
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

	var organizationUserId uint64 = 0
	var organizationTutorId uint64 = 0

	if parentUserId > 0 {
		parentUser, err := uouc.urepo.Get(ctx, parentUserId)

		if err != nil {
			return nil, WeixinUserNotFound
		}

		parentUserOrganizationRelation, err := uouc.uorrepo.GetByUserId(ctx, parentUser.Id, companyOrganization.Data.OrganizationId, 0, "0")

		if err != nil {
			return nil, WeixinUserOrganizationRelationNotFound
		}

		if parentUserOrganizationRelation.OrganizationId != organizationId {
			return nil, WeixinUserOrderCreateError
		}

		organizationUserId = parentUserOrganizationRelation.UserId

		if parentUserOrganizationRelation.Level == 4 {
			organizationTutorId = parentUserOrganizationRelation.UserId
		} else {
			userIntegralRelations, err := uouc.uirrepo.List(ctx, companyOrganization.Data.OrganizationId)

			if err != nil {
				return nil, WeixinUserOrganizationRelationListError
			}

			tutorUserIntegralRelation := uouc.uirrepo.GetSuperior(ctx, parentUserOrganizationRelation.UserId, 4, userIntegralRelations)

			if tutorUserIntegralRelation != nil {
				organizationTutorId = tutorUserIntegralRelation.UserId
			}
		}
	}

	if organizationUserId == 0 {
		return nil, WeixinUserOrderParentNotExist
	}

	appid := ""

	if uouc.oconf.DjOrganizationId == organizationId {
		appid = uouc.wconf.DjMini.Appid
	} else if uouc.oconf.DefaultOrganizationId == organizationId {
		appid = uouc.wconf.Mini.Appid
	} else if uouc.oconf.LbOrganizationId == organizationId {
		appid = uouc.wconf.Mini.Appid
	} else {
		return nil, WeixinUserOpenidNotFound
	}

	userOpenId, err := uouc.uoirepo.Get(ctx, user.Id, appid, "")

	if err != nil {
		return nil, WeixinUserOpenidNotFound
	}

	outTradeNo, err := uouc.repo.NextId(ctx)

	if err != nil {
		return nil, WeixinUserOrderCreateError
	}

	if len(organizationCourses) > 0 {
		var level uint8 = 0
		isNotExist := true

		for _, organizationCourse := range organizationCourses {
			if organizationCourse.CoursePrice == payAmount {
				level = organizationCourse.CourseLevel
				isNotExist = false

				break
			}
		}

		if isNotExist {
			return nil, WeixinUserOrderCourseNotExist
		}

		inUserOrder := domain.NewUserOrder(ctx, user.Id, companyOrganization.Data.OrganizationId, organizationUserId, organizationTutorId, strconv.FormatUint(outTradeNo, 10), "", "", float32(organizationCourses[level-1].CoursePrice), float32(organizationCourses[level-1].CoursePrice), nil, level, 0, 1)
		inUserOrder.SetCreateTime(ctx)
		inUserOrder.SetUpdateTime(ctx)

		if _, err := uouc.repo.Save(ctx, inUserOrder); err != nil {
			return nil, WeixinUserOrderCreateError
		}

		payData, err := uouc.prepo.Pay(ctx, companyOrganization.Data.OrganizationId, organizationCourses[level-1].CoursePrice, strconv.FormatUint(outTradeNo, 10), "课程购买", tool.GetRandCode(time.Now().String())[0:30], userOpenId.OpenId, clientIp)

		if err != nil {
			return nil, WeixinUserOrderCreateError
		}

		return &domain.UserOrderRequestPayment{
			TimeStamp:  payData.Data.TimeStamp,
			NonceStr:   payData.Data.NonceStr,
			Package:    payData.Data.Package,
			SignType:   payData.Data.SignType,
			PaySign:    payData.Data.PaySign,
			OutTradeNo: strconv.FormatUint(outTradeNo, 10),
			PayAmount:  fmt.Sprintf("%.2f", tool.Decimal(payAmount, 2)),
			LevelName:  WeixinUserOrganizationRelationLevel[level-1],
		}, nil
	} else {
		return nil, WeixinUserOrderCreateError
	}
}

func (uouc *UserOrderUsecase) UpgradeUserOrders(ctx context.Context, userId, organizationId uint64, payAmount float64, clientIp string) (*domain.UserOrderRequestPayment, error) {
	user, err := uouc.urepo.Get(ctx, userId)

	if err != nil {
		return nil, WeixinLoginError
	}

	companyOrganization, err := uouc.corepo.Get(ctx, organizationId)

	if err != nil {
		return nil, WeixinCompanyOrganizationNotFound
	}

	userOrganizationRelation, err := uouc.uorrepo.GetByUserId(ctx, user.Id, companyOrganization.Data.OrganizationId, 0, "0")

	if err != nil {
		return nil, WeixinUserOrganizationRelationNotFound
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

	var level uint8 = 0
	isNotExist := true

	for _, organizationCourse := range organizationCourses {
		if organizationCourse.CoursePrice == payAmount {
			level = organizationCourse.CourseLevel
			isNotExist = false

			break
		}
	}

	if isNotExist {
		return nil, WeixinUserOrderCourseNotExist
	}

	if level < userOrganizationRelation.Level {
		return nil, WeixinUserOrderCreateError
	}

	var rulePayAmount float64 = payAmount - organizationCourses[userOrganizationRelation.Level-1].CoursePrice

	if rulePayAmount < 0 {
		return nil, WeixinUserOrderCreateError
	}

	appid := ""

	if uouc.oconf.DjOrganizationId == userOrganizationRelation.OrganizationId {
		appid = uouc.wconf.DjMini.Appid
	} else if uouc.oconf.DefaultOrganizationId == userOrganizationRelation.OrganizationId {
		appid = uouc.wconf.Mini.Appid
	} else if uouc.oconf.LbOrganizationId == userOrganizationRelation.OrganizationId {
		appid = uouc.wconf.Mini.Appid
	} else {
		return nil, WeixinUserOpenidNotFound
	}

	userOpenId, err := uouc.uoirepo.Get(ctx, user.Id, appid, "")

	if err != nil {
		return nil, WeixinUserOpenidNotFound
	}

	outTradeNo, err := uouc.repo.NextId(ctx)

	if err != nil {
		return nil, WeixinUserOrderCreateError
	}

	inUserOrder := domain.NewUserOrder(ctx, user.Id, userOrganizationRelation.OrganizationId, userOrganizationRelation.OrganizationUserId, userOrganizationRelation.OrganizationTutorId, strconv.FormatUint(outTradeNo, 10), "", "", float32(payAmount), float32(rulePayAmount), nil, level, 0, 2)
	inUserOrder.SetCreateTime(ctx)
	inUserOrder.SetUpdateTime(ctx)

	if _, err := uouc.repo.Save(ctx, inUserOrder); err != nil {
		return nil, WeixinUserOrderCreateError
	}

	payData, err := uouc.prepo.Pay(ctx, companyOrganization.Data.OrganizationId, rulePayAmount, strconv.FormatUint(outTradeNo, 10), "课程购买", tool.GetRandCode(time.Now().String())[0:30], userOpenId.OpenId, clientIp)

	if err != nil {
		return nil, WeixinUserOrderCreateError
	}

	return &domain.UserOrderRequestPayment{
		TimeStamp:  payData.Data.TimeStamp,
		NonceStr:   payData.Data.NonceStr,
		Package:    payData.Data.Package,
		SignType:   payData.Data.SignType,
		PaySign:    payData.Data.PaySign,
		OutTradeNo: strconv.FormatUint(outTradeNo, 10),
		PayAmount:  fmt.Sprintf("%.2f", tool.Decimal(rulePayAmount, 2)),
		LevelName:  WeixinUserOrganizationRelationLevel[level-1],
	}, nil
}

func (uouc *UserOrderUsecase) AsyncNotificationUserOrders(ctx context.Context, content string) error {
	asyncNotification, err := uouc.prepo.PayAsyncNotification(ctx, content)

	if err != nil {
		return err
	}

	payTime, err := tool.StringToTime("2006-01-02T15:04:05-07:00", asyncNotification.Data.PayTime)

	if err != nil {
		return err
	}

	inUserOrder, err := uouc.repo.GetByOutTradeNo(ctx, asyncNotification.Data.OutTradeNo)

	if err != nil {
		return WeixinUserOrderNotFound
	}

	if inUserOrder.PayStatus == 1 {
		return nil
	}

	var inUserOrganizationRelation *domain.UserOrganizationRelation

	err = uouc.tm.InTx(ctx, func(ctx context.Context) error {
		inUserOrder.SetOutTransactionId(ctx, asyncNotification.Data.OutTransactionId)
		inUserOrder.SetTransactionId(ctx, asyncNotification.Data.TransactionId)
		inUserOrder.SetPayTime(ctx, &payTime)
		inUserOrder.SetPayStatus(ctx, 1)
		inUserOrder.SetUpdateTime(ctx)

		if _, err = uouc.repo.Update(ctx, inUserOrder); err != nil {
			return err
		}

		if err := uouc.uorrepo.DeleteByUserId(ctx, inUserOrder.UserId, "1"); err != nil {
			return err
		}

		if inUserOrganizationRelation, err = uouc.uorrepo.GetByUserId(ctx, inUserOrder.UserId, inUserOrder.OrganizationId, 0, "0"); err == nil {
			if inUserOrder.Level > inUserOrganizationRelation.Level {
				if inUserOrder.Level == 4 {
					inUserOrganizationRelation.SetOrganizationTutorId(ctx, 0)
				}

				inUserOrganizationRelation.SetLevel(ctx, inUserOrder.Level)
				inUserOrganizationRelation.SetUpdateTime(ctx)

				if _, err := uouc.uorrepo.Update(ctx, inUserOrganizationRelation); err != nil {
					return err
				}

				if inUserIntegralRelation, err := uouc.uirrepo.GetByUserId(ctx, inUserOrder.UserId, inUserOrder.OrganizationId, 0); err == nil {
					inUserIntegralRelation.SetLevel(ctx, inUserOrder.Level)
					inUserIntegralRelation.SetUpdateTime(ctx)

					if _, err := uouc.uirrepo.Update(ctx, inUserIntegralRelation); err != nil {
						return err
					}
				} else {
					inUserIntegralRelation = domain.NewUserIntegralRelation(ctx, inUserOrder.UserId, inUserOrder.OrganizationId, inUserOrder.OrganizationUserId, inUserOrder.Level)
					inUserIntegralRelation.SetCreateTime(ctx)
					inUserIntegralRelation.SetUpdateTime(ctx)

					if _, err := uouc.uirrepo.Save(ctx, inUserIntegralRelation); err != nil {
						return err
					}
				}

				if inUserOrder.Level == 4 {
					tmpChildIds := make([]uint64, 0)
					childIds := make([]uint64, 0)

					userIntegralRelations, err := uouc.uirrepo.List(ctx, inUserOrder.OrganizationId)

					if err != nil {
						return WeixinUserOrganizationRelationListError
					}

					uouc.uirrepo.ListChildId(ctx, inUserOrder.UserId, &tmpChildIds, userIntegralRelations)

					if len(tmpChildIds) > 0 {
						userOrganizationRelations, err := uouc.uorrepo.List(ctx, inUserOrder.OrganizationId)

						if err != nil {
							return WeixinUserOrganizationRelationListError
						}

						for _, tmpChildId := range tmpChildIds {
							for _, userOrganizationRelation := range userOrganizationRelations {
								if userOrganizationRelation.UserId == tmpChildId {
									if (userOrganizationRelation.OrganizationTutorId == 0 && userOrganizationRelation.Level != 4) || userOrganizationRelation.OrganizationTutorId != 0 {
										isNotExist := true

										for _, ltmpChildId := range tmpChildIds {
											if ltmpChildId == userOrganizationRelation.OrganizationTutorId {
												isNotExist = false

												break
											}
										}

										if isNotExist {
											childIds = append(childIds, tmpChildId)
										}
									}

									break
								}
							}
						}

						if len(childIds) > 0 {
							if err := uouc.uorrepo.UpdateSuperior(ctx, inUserOrder.UserId, childIds); err != nil {
								return err
							}
						}
					}
				}
			}
		} else {
			var wcontent *wxa.GetUnlimitedQRCodeResponse

			if inUserOrder.OrganizationId == uouc.oconf.DjOrganizationId {
				accessToken, err := oauth2.GetAccessToken(uouc.wconf.DjMini.Appid, uouc.wconf.DjMini.Secret)

				if err != nil {
					return err
				}

				wcontent, err = wxa.GetUnlimitedQRCode(accessToken.AccessToken, "oId="+strconv.FormatUint(inUserOrder.OrganizationId, 10)+"&uId="+strconv.FormatUint(inUserOrder.UserId, 10), uouc.wconf.DjMini.QrCodeEnvVersion)

				if err != nil {
					return err
				}
			} else {
				accessToken, err := oauth2.GetAccessToken(uouc.wconf.Mini.Appid, uouc.wconf.Mini.Secret)

				if err != nil {
					return err
				}

				wcontent, err = wxa.GetUnlimitedQRCode(accessToken.AccessToken, "oId="+strconv.FormatUint(inUserOrder.OrganizationId, 10)+"&uId="+strconv.FormatUint(inUserOrder.UserId, 10), uouc.wconf.Mini.QrCodeEnvVersion)

				if err != nil {
					return err
				}
			}

			objectKey := tool.GetRandCode(time.Now().String())

			if _, err := uouc.repo.PutContent(ctx, uouc.vconf.Tos.Company.SubFolder+"/"+objectKey+".png", strings.NewReader(wcontent.Buffer)); err != nil {
				return err
			}

			organizationTutorId := inUserOrder.OrganizationTutorId

			if inUserOrder.Level == 4 {
				organizationTutorId = 0
			}

			inUserOrganizationRelation = domain.NewUserOrganizationRelation(ctx, inUserOrder.UserId, inUserOrder.OrganizationId, inUserOrder.OrganizationUserId, organizationTutorId, inUserOrder.Level, 0, uouc.vconf.Tos.Company.Url+"/"+uouc.vconf.Tos.Company.SubFolder+"/"+objectKey+".png")
			inUserOrganizationRelation.SetCreateTime(ctx)
			inUserOrganizationRelation.SetUpdateTime(ctx)

			if _, err := uouc.uorrepo.Save(ctx, inUserOrganizationRelation); err != nil {
				return err
			}

			inUserIntegralRelation := domain.NewUserIntegralRelation(ctx, inUserOrder.UserId, inUserOrder.OrganizationId, inUserOrder.OrganizationUserId, inUserOrder.Level)
			inUserIntegralRelation.SetCreateTime(ctx)
			inUserIntegralRelation.SetUpdateTime(ctx)

			if _, err := uouc.uirrepo.Save(ctx, inUserIntegralRelation); err != nil {
				return err
			}
		}

		if inUserOrder.Level == 3 || inUserOrder.Level == 4 {
			var num uint32

			if inUserOrder.Level == 3 {
				num = 15
			} else if inUserOrder.Level == 4 {
				num = 48
			}

			inUserCouponCreateLog := domain.NewUserCouponCreateLog(ctx, inUserOrder.UserId, inUserOrder.OrganizationId, num, 2, 0)
			inUserCouponCreateLog.SetCreateTime(ctx)
			inUserCouponCreateLog.SetUpdateTime(ctx)

			if _, err := uouc.ucclrepo.Save(ctx, inUserCouponCreateLog); err != nil {
				return err
			}
		}

		uiday, _ := strconv.ParseUint(time.Now().Format("20060102"), 10, 64)

		err = uouc.getUserComission(ctx, uiday, inUserOrder, inUserOrganizationRelation)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (uouc *UserOrderUsecase) getUserComission(ctx context.Context, uiday uint64, userOrder *domain.UserOrder, userOrganizationRelation *domain.UserOrganizationRelation) error {
	companyOrganization, err := uouc.corepo.Get(ctx, userOrder.OrganizationId)

	if err != nil {
		return err
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

	var coursePrice float64
	isNotExist := true

	for _, organizationCourse := range organizationCourses {
		if organizationCourse.CourseLevel == userOrder.Level {
			coursePrice = organizationCourse.CoursePrice
			isNotExist = false

			break
		}
	}

	if isNotExist {
		return WeixinUserOrderCourseNotExist
	}

	var courseCommissionPool float64

	if userOrder.Level == 1 {
		courseCommissionPool = float64(userOrder.PayAmount) * companyOrganization.Data.OrganizationColonelCommission.ZeroCourseRatio / 100
	} else if userOrder.Level == 2 {
		courseCommissionPool = float64(userOrder.PayAmount) * companyOrganization.Data.OrganizationColonelCommission.PrimaryCourseRatio / 100
	} else if userOrder.Level == 3 {
		courseCommissionPool = float64(userOrder.PayAmount) * companyOrganization.Data.OrganizationColonelCommission.IntermediateCourseRatio / 100
	} else if userOrder.Level == 4 {
		courseCommissionPool = float64(userOrder.PayAmount) * companyOrganization.Data.OrganizationColonelCommission.AdvancedCourseRatio / 100
	}

	var organizationParentUser *domain.UserOrganizationRelation
	var organizationTutorUser *domain.UserOrganizationRelation

	if userOrganizationRelation.OrganizationUserId > 0 {
		organizationParentUser, err = uouc.uorrepo.GetByUserId(ctx, userOrganizationRelation.OrganizationUserId, userOrganizationRelation.OrganizationId, 0, "0")

		if err != nil {
			return WeixinUserOrganizationRelationNotFound
		}

		if userOrganizationRelation.OrganizationTutorId > 0 {
			organizationTutorUser, err = uouc.uorrepo.GetByUserId(ctx, userOrganizationRelation.OrganizationTutorId, userOrganizationRelation.OrganizationId, 0, "0")

			if err != nil {
				return WeixinUserOrganizationRelationNotFound
			}
		}

		var realCommission float64

		if userOrder.Level == 1 {
			if organizationParentUser.Level == 1 {
				realCommission = companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedPresenterZeroCourseCommissionRule * float64(userOrder.PayAmount) / coursePrice
			} else if organizationParentUser.Level == 2 {
				realCommission = companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterZeroCourseCommissionRule * float64(userOrder.PayAmount) / coursePrice
			} else if organizationParentUser.Level == 3 {
				realCommission = companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterZeroCourseCommissionRule * float64(userOrder.PayAmount) / coursePrice
			} else if organizationParentUser.Level == 4 {
				realCommission = companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterZeroCourseCommissionRule * float64(userOrder.PayAmount) / coursePrice
			}

			inUserCommission := domain.NewUserCommission(ctx, organizationParentUser.UserId, organizationParentUser.OrganizationId, userOrganizationRelation.UserId, userOrder.Id, userOrganizationRelation.Level, organizationParentUser.Level, 1, 2, 1, 1, 1, userOrder.PayAmount, float32(courseCommissionPool), float32(tool.Decimal(realCommission, 0)))
			inUserCommission.SetCreateTime(ctx, *userOrder.PayTime)
			inUserCommission.SetUpdateTime(ctx)

			if _, err := uouc.ucrepo.Save(ctx, inUserCommission); err != nil {
				return err
			}

			if organizationTutorUser != nil && organizationParentUser.Level != 4 {
				if organizationTutorUser.Level == 4 {
					if organizationParentUser.Level == 1 {
						realCommission = companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedTutorZeroCourseCommissionRule * float64(userOrder.PayAmount) / coursePrice
					} else if organizationParentUser.Level == 2 {
						realCommission = companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorZeroCourseCommissionRule * float64(userOrder.PayAmount) / coursePrice
					} else if organizationParentUser.Level == 3 {
						realCommission = companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorZeroCourseCommissionRule * float64(userOrder.PayAmount) / coursePrice
					}

					inUserCommission = domain.NewUserCommission(ctx, organizationTutorUser.UserId, organizationTutorUser.OrganizationId, userOrganizationRelation.UserId, userOrder.Id, userOrganizationRelation.Level, organizationTutorUser.Level, 2, 2, 1, 1, 1, userOrder.PayAmount, float32(courseCommissionPool), float32(tool.Decimal(realCommission, 0)))
					inUserCommission.SetCreateTime(ctx, *userOrder.PayTime)
					inUserCommission.SetUpdateTime(ctx)

					if _, err = uouc.ucrepo.Save(ctx, inUserCommission); err != nil {
						return err
					}
				}
			}
		} else if userOrder.Level == 2 {
			if organizationParentUser.Level == 2 {
				realCommission = companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterPrimaryCourseCommissionRule * float64(userOrder.PayAmount) / coursePrice
			} else if organizationParentUser.Level == 3 {
				realCommission = companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterPrimaryCourseCommissionRule * float64(userOrder.PayAmount) / coursePrice
			} else if organizationParentUser.Level == 4 {
				realCommission = companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterPrimaryCourseCommissionRule * float64(userOrder.PayAmount) / coursePrice
			}

			inUserCommission := domain.NewUserCommission(ctx, organizationParentUser.UserId, organizationParentUser.OrganizationId, userOrganizationRelation.UserId, userOrder.Id, userOrganizationRelation.Level, organizationParentUser.Level, 1, 2, 1, 1, 1, userOrder.PayAmount, float32(courseCommissionPool), float32(tool.Decimal(realCommission, 0)))
			inUserCommission.SetCreateTime(ctx, *userOrder.PayTime)
			inUserCommission.SetUpdateTime(ctx)

			if _, err := uouc.ucrepo.Save(ctx, inUserCommission); err != nil {
				return err
			}

			if organizationTutorUser != nil && organizationParentUser.Level != 4 {
				if organizationTutorUser.Level == 4 {
					if organizationParentUser.Level == 2 {
						realCommission = companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorPrimaryCourseCommissionRule * float64(userOrder.PayAmount) / coursePrice
					} else if organizationParentUser.Level == 3 {
						realCommission = companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorPrimaryCourseCommissionRule * float64(userOrder.PayAmount) / coursePrice
					}

					inUserCommission = domain.NewUserCommission(ctx, organizationTutorUser.UserId, organizationTutorUser.OrganizationId, userOrganizationRelation.UserId, userOrder.Id, userOrganizationRelation.Level, organizationTutorUser.Level, 2, 2, 1, 1, 1, userOrder.PayAmount, float32(courseCommissionPool), float32(tool.Decimal(realCommission, 0)))
					inUserCommission.SetCreateTime(ctx, *userOrder.PayTime)
					inUserCommission.SetUpdateTime(ctx)

					if _, err = uouc.ucrepo.Save(ctx, inUserCommission); err != nil {
						return err
					}
				}
			}
		} else if userOrder.Level == 3 {
			if organizationParentUser.Level == 2 {
				realCommission = companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterIntermediateCourseCommissionRule * float64(userOrder.PayAmount) / coursePrice
			} else if organizationParentUser.Level == 3 {
				realCommission = companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterIntermediateCourseCommissionRule * float64(userOrder.PayAmount) / coursePrice
			} else if organizationParentUser.Level == 4 {
				realCommission = companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterIntermediateCourseCommissionRule * float64(userOrder.PayAmount) / coursePrice
			}

			inUserCommission := domain.NewUserCommission(ctx, organizationParentUser.UserId, organizationParentUser.OrganizationId, userOrganizationRelation.UserId, userOrder.Id, userOrganizationRelation.Level, organizationParentUser.Level, 1, 2, 1, 1, 1, userOrder.PayAmount, float32(courseCommissionPool), float32(tool.Decimal(realCommission, 0)))
			inUserCommission.SetCreateTime(ctx, *userOrder.PayTime)
			inUserCommission.SetUpdateTime(ctx)

			if _, err := uouc.ucrepo.Save(ctx, inUserCommission); err != nil {
				return err
			}

			if organizationTutorUser != nil && organizationParentUser.Level != 4 {
				if organizationTutorUser.Level == 4 {
					if organizationParentUser.Level == 2 {
						realCommission = companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorIntermediateCourseCommissionRule * float64(userOrder.PayAmount) / coursePrice
					} else if organizationParentUser.Level == 3 {
						realCommission = companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorIntermediateCourseCommissionRule * float64(userOrder.PayAmount) / coursePrice
					}

					inUserCommission = domain.NewUserCommission(ctx, organizationTutorUser.UserId, organizationTutorUser.OrganizationId, userOrganizationRelation.UserId, userOrder.Id, userOrganizationRelation.Level, organizationTutorUser.Level, 2, 2, 1, 1, 1, userOrder.PayAmount, float32(courseCommissionPool), float32(tool.Decimal(realCommission, 0)))
					inUserCommission.SetCreateTime(ctx, *userOrder.PayTime)
					inUserCommission.SetUpdateTime(ctx)

					if _, err = uouc.ucrepo.Save(ctx, inUserCommission); err != nil {
						return err
					}
				}
			}
		} else if userOrder.Level == 4 {
			if organizationParentUser.Level == 2 {
				realCommission = companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterAdvancedCourseCommissionRule * float64(userOrder.PayAmount) / coursePrice
			} else if organizationParentUser.Level == 3 {
				realCommission = companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterAdvancedCourseCommissionRule * float64(userOrder.PayAmount) / coursePrice
			} else if organizationParentUser.Level == 4 {
				realCommission = companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterAdvancedCourseCommissionRule * float64(userOrder.PayAmount) / coursePrice
			}

			inUserCommission := domain.NewUserCommission(ctx, organizationParentUser.UserId, organizationParentUser.OrganizationId, userOrganizationRelation.UserId, userOrder.Id, userOrganizationRelation.Level, organizationParentUser.Level, 1, 2, 1, 1, 1, userOrder.PayAmount, float32(courseCommissionPool), float32(tool.Decimal(realCommission, 0)))
			inUserCommission.SetCreateTime(ctx, *userOrder.PayTime)
			inUserCommission.SetUpdateTime(ctx)

			if _, err := uouc.ucrepo.Save(ctx, inUserCommission); err != nil {
				return err
			}

			if organizationTutorUser != nil && organizationParentUser.Level != 4 {
				if organizationTutorUser.Level == 4 {
					if organizationParentUser.Level == 2 {
						realCommission = companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorAdvancedCourseCommissionRule * float64(userOrder.PayAmount) / coursePrice
					} else if organizationParentUser.Level == 3 {
						realCommission = companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorAdvancedCourseCommissionRule * float64(userOrder.PayAmount) / coursePrice
					}

					inUserCommission = domain.NewUserCommission(ctx, organizationTutorUser.UserId, organizationTutorUser.OrganizationId, userOrganizationRelation.UserId, userOrder.Id, userOrganizationRelation.Level, organizationTutorUser.Level, 2, 2, 1, 1, 1, userOrder.PayAmount, float32(courseCommissionPool), float32(tool.Decimal(realCommission, 0)))
					inUserCommission.SetCreateTime(ctx, *userOrder.PayTime)
					inUserCommission.SetUpdateTime(ctx)

					if _, err = uouc.ucrepo.Save(ctx, inUserCommission); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (uouc *UserOrderUsecase) SyncQrCodeUserOrganizationRelations(ctx context.Context) error {
	if userOrganizationRelations, err := uouc.uorrepo.List(ctx, uouc.oconf.DjOrganizationId); err == nil {
		var wcontent *wxa.GetUnlimitedQRCodeResponse

		if accessToken, err := oauth2.GetAccessToken(uouc.wconf.DjMini.Appid, uouc.wconf.DjMini.Secret); err == nil {
			for _, inUserOrganizationRelation := range userOrganizationRelations {
				if inUserOrganizationRelation.IsOrderRelation == 0 {
					if wcontent, err = wxa.GetUnlimitedQRCode(accessToken.AccessToken, "oId="+strconv.FormatUint(inUserOrganizationRelation.OrganizationId, 10)+"&uId="+strconv.FormatUint(inUserOrganizationRelation.UserId, 10), uouc.wconf.DjMini.QrCodeEnvVersion); err == nil {
						objectKey := tool.GetRandCode(time.Now().String())

						if _, err := uouc.repo.PutContent(ctx, uouc.vconf.Tos.Company.SubFolder+"/"+objectKey+".png", strings.NewReader(wcontent.Buffer)); err == nil {
							inUserOrganizationRelation.SetOrganizationUserQrCodeUrl(ctx, uouc.vconf.Tos.Company.Url+"/"+uouc.vconf.Tos.Company.SubFolder+"/"+objectKey+".png")
							inUserOrganizationRelation.SetUpdateTime(ctx)

							uouc.uorrepo.Update(ctx, inUserOrganizationRelation)
						}
					}
				}
			}
		}
	}

	if userOrganizationRelations, err := uouc.uorrepo.List(ctx, uouc.oconf.DefaultOrganizationId); err == nil {
		var wcontent *wxa.GetUnlimitedQRCodeResponse

		if accessToken, err := oauth2.GetAccessToken(uouc.wconf.Mini.Appid, uouc.wconf.Mini.Secret); err == nil {
			for _, inUserOrganizationRelation := range userOrganizationRelations {
				if inUserOrganizationRelation.IsOrderRelation == 0 {
					if wcontent, err = wxa.GetUnlimitedQRCode(accessToken.AccessToken, "oId="+strconv.FormatUint(inUserOrganizationRelation.OrganizationId, 10)+"&uId="+strconv.FormatUint(inUserOrganizationRelation.UserId, 10), uouc.wconf.Mini.QrCodeEnvVersion); err == nil {
						objectKey := tool.GetRandCode(time.Now().String())

						if _, err := uouc.repo.PutContent(ctx, uouc.vconf.Tos.Company.SubFolder+"/"+objectKey+".png", strings.NewReader(wcontent.Buffer)); err != nil {
							inUserOrganizationRelation.SetOrganizationUserQrCodeUrl(ctx, uouc.vconf.Tos.Company.Url+"/"+uouc.vconf.Tos.Company.SubFolder+"/"+objectKey+".png")
							inUserOrganizationRelation.SetUpdateTime(ctx)

							uouc.uorrepo.Update(ctx, inUserOrganizationRelation)
						}

					}
				}
			}
		}
	}

	return nil
}

func (uouc *UserOrderUsecase) OperationUserOrders(ctx context.Context) error {
	userOrders, _ := uouc.repo.ListOperation(ctx)

	uiday, _ := strconv.ParseUint(time.Now().Format("20060102"), 10, 64)

	for _, userOrder := range userOrders {
		if userOrganizationRelation, err := uouc.uorrepo.GetByUserId(ctx, userOrder.UserId, userOrder.OrganizationId, 0, "0"); err == nil {
			if _, err := uouc.ucrepo.GetByRelevanceId(ctx, userOrder.Id, 1); err != nil {
				uouc.getUserComission(ctx, uiday, userOrder, userOrganizationRelation)
			}
		}
	}

	return nil
}

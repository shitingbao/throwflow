package biz

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	v1 "weixin/api/service/company/v1"
	"weixin/internal/conf"
	"weixin/internal/domain"
	"weixin/internal/pkg/tool"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	WeixinUserOrganizationRelationNotFound       = errors.NotFound("WEIXIN_USER_ORGANIZATION_RELATION_NOT_FOUND", "微信用户绑定机构不存在")
	WeixinUserOrganizationRelationExist          = errors.InternalServer("WEIXIN_USER_ORGANIZATION_RELATION_EXIST", "微信用户已绑定机构")
	WeixinUserOrganizationRelationOtherExist     = errors.InternalServer("WEIXIN_USER_ORGANIZATION_RELATION_OTHER_EXIST", "微信用户已绑定其他机构")
	WeixinUserOrganizationRelationNotNeedUpgrade = errors.InternalServer("WEIXIN_USER_ORGANIZATION_RELATION_NOT_NEED_UPGRADE", "微信用户绑定机构不需要升级")
	WeixinUserOrganizationRelationUnbindError    = errors.InternalServer("WEIXIN_USER_ORGANIZATION_RELATION_UNBIND_ERROR", "微信用户解绑机构失败")
	WeixinUserOrganizationRelationbindError      = errors.InternalServer("WEIXIN_USER_ORGANIZATION_RELATION_BIND_ERROR", "微信用户绑定机构失败")
	WeixinUserOrganizationRelationListError      = errors.InternalServer("WEIXIN_USER_ORGANIZATION_RELATION_LIST_ERROR", "微信用户机构关系列表获取失败")
	WeixinUserOrganizationRelationUpdateError    = errors.InternalServer("WEIXIN_USER_ORGANIZATION_RELATION_UPDATE_ERROR", "微信用户机构关系更新失败")

	WeixinUserOrganizationRelationLevel = [4]string{"零级", "初级", "中级", "高级"}
)

type UserOrganizationRelationRepo interface {
	Get(context.Context, uint64, uint64, string) (*domain.UserOrganizationRelation, error)
	GetByUserId(context.Context, uint64, uint64, uint64, string) (*domain.UserOrganizationRelation, error)
	List(context.Context, uint64) ([]*domain.UserOrganizationRelation, error)
	ListDirectChild(context.Context, uint64, uint64) ([]*domain.UserOrganizationRelation, error)
	Count(context.Context, uint64, uint64, string) (int64, error)
	Save(context.Context, *domain.UserOrganizationRelation) (*domain.UserOrganizationRelation, error)
	Update(context.Context, *domain.UserOrganizationRelation) (*domain.UserOrganizationRelation, error)
	UpdateSuperior(context.Context, uint64, []uint64) error
	DeleteByUserId(context.Context, uint64, string) error
}

type UserOrganizationRelationUsecase struct {
	repo    UserOrganizationRelationRepo
	uirrepo UserIntegralRelationRepo
	urepo   UserRepo
	uodrepo UserOpenDouyinRepo
	usrrepo UserScanRecordRepo
	uorepo  UserOrderRepo
	ucrepo  UserCommissionRepo
	ucorepo UserCouponRepo
	crepo   CompanyRepo
	corepo  CompanyOrganizationRepo
	jorepo  JinritemaiOrderRepo
	dorepo  DoukeOrderRepo
	darepo  DjAwemeRepo
	tlrepo  TaskLogRepo
	tm      Transaction
	conf    *conf.Data
	oconf   *conf.Organization
	wconf   *conf.Weixin
	log     *log.Helper
}

func NewUserOrganizationRelationUsecase(repo UserOrganizationRelationRepo, uirrepo UserIntegralRelationRepo, urepo UserRepo, uodrepo UserOpenDouyinRepo, usrrepo UserScanRecordRepo, uorepo UserOrderRepo, ucrepo UserCommissionRepo, ucorepo UserCouponRepo, crepo CompanyRepo, corepo CompanyOrganizationRepo, jorepo JinritemaiOrderRepo, dorepo DoukeOrderRepo, darepo DjAwemeRepo, tlrepo TaskLogRepo, tm Transaction, conf *conf.Data, oconf *conf.Organization, wconf *conf.Weixin, logger log.Logger) *UserOrganizationRelationUsecase {
	return &UserOrganizationRelationUsecase{repo: repo, uirrepo: uirrepo, urepo: urepo, uodrepo: uodrepo, usrrepo: usrrepo, uorepo: uorepo, ucrepo: ucrepo, ucorepo: ucorepo, crepo: crepo, corepo: corepo, jorepo: jorepo, dorepo: dorepo, darepo: darepo, tlrepo: tlrepo, tm: tm, conf: conf, oconf: oconf, wconf: wconf, log: log.NewHelper(logger)}
}

func (uoruc *UserOrganizationRelationUsecase) GetUserOrganizationRelations(ctx context.Context, userId uint64) (*domain.UserOrganizationRelationInfo, error) {
	user, err := uoruc.urepo.Get(ctx, userId)

	if err != nil {
		return nil, WeixinLoginError
	}

	total, _ := uoruc.urepo.Count(ctx)

	if userOrganizationRelation, err := uoruc.repo.GetByUserId(ctx, user.Id, 0, 0, "0"); err == nil {
		companyOrganization, err := uoruc.corepo.Get(ctx, userOrganizationRelation.OrganizationId)

		if err != nil {
			return nil, WeixinCompanyOrganizationNotFound
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
				CourseName:          organizationCourse.CourseName,
				CourseSubName:       organizationCourse.CourseSubName,
				CoursePrice:         organizationCourse.CoursePrice,
				CourseDuration:      organizationCourse.CourseDuration,
				CourseLevel:         uint8(organizationCourse.CourseLevel),
				CourseOriginalPrice: organizationCourse.CourseOriginalPrice,
				CourseModules:       courseModules,
			})
		}

		sort.Sort(domain.OrganizationCourses(organizationCourses))

		levelName := WeixinUserOrganizationRelationLevel[userOrganizationRelation.Level-1]

		return &domain.UserOrganizationRelationInfo{
			OrganizationId:            companyOrganization.Data.OrganizationId,
			OrganizationName:          companyOrganization.Data.OrganizationName,
			OrganizationLogoUrl:       companyOrganization.Data.OrganizationLogoUrl,
			OrganizationCourses:       organizationCourses,
			CompanyName:               companyOrganization.Data.CompanyName,
			BankCode:                  companyOrganization.Data.BankCode,
			BankDeposit:               companyOrganization.Data.BankDeposit,
			ActivationTime:            userOrganizationRelation.CreateTime,
			LevelName:                 levelName,
			Level:                     userOrganizationRelation.Level,
			OrganizationUserQrCodeUrl: userOrganizationRelation.OrganizationUserQrCodeUrl,
			Total:                     uint64(total),
		}, nil
	} else {
		userScanRecord, err := uoruc.usrrepo.Get(ctx, user.Id, 0, 1)

		if err != nil {
			return nil, WeixinUserScanRecordNotFound
		}

		companyOrganization, err := uoruc.corepo.Get(ctx, userScanRecord.OrganizationId)

		if err != nil {
			return nil, WeixinCompanyOrganizationNotFound
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
				CourseName:          organizationCourse.CourseName,
				CourseSubName:       organizationCourse.CourseSubName,
				CoursePrice:         organizationCourse.CoursePrice,
				CourseDuration:      organizationCourse.CourseDuration,
				CourseOriginalPrice: organizationCourse.CourseOriginalPrice,
				CourseLevel:         uint8(organizationCourse.CourseLevel),
				CourseModules:       courseModules,
			})
		}

		sort.Sort(domain.OrganizationCourses(organizationCourses))

		var parentUserId uint64 = 0
		parentNickName, parentAvatarUrl := "", ""

		if userScanRecord.OrganizationUserId > 0 {
			if parentUser, err := uoruc.urepo.Get(ctx, userScanRecord.OrganizationUserId); err == nil {
				parentUserId = parentUser.Id
				parentNickName = parentUser.NickName
				parentAvatarUrl = parentUser.AvatarUrl
			}
		}

		return &domain.UserOrganizationRelationInfo{
			OrganizationId:      companyOrganization.Data.OrganizationId,
			OrganizationName:    companyOrganization.Data.OrganizationName,
			OrganizationLogoUrl: companyOrganization.Data.OrganizationLogoUrl,
			OrganizationCourses: organizationCourses,
			CompanyName:         companyOrganization.Data.CompanyName,
			BankCode:            companyOrganization.Data.BankCode,
			BankDeposit:         companyOrganization.Data.BankDeposit,
			ParentUserId:        parentUserId,
			ParentNickName:      parentNickName,
			ParentAvatarUrl:     parentAvatarUrl,
			Total:               uint64(total),
		}, nil
	}
}

func (uoruc *UserOrganizationRelationUsecase) GetBindUserOrganizationRelations(ctx context.Context, userId, organizationId uint64) (*domain.BindUserOrganizationRelationInfo, error) {
	user, err := uoruc.urepo.Get(ctx, userId)

	if err != nil {
		return nil, WeixinLoginError
	}

	userOrganizationRelation, err := uoruc.repo.GetByUserId(ctx, user.Id, organizationId, 0, "0")

	if err != nil {
		return nil, WeixinUserOrganizationRelationNotFound
	}

	bindUserOrganizationRelationInfo := &domain.BindUserOrganizationRelationInfo{
		OrganizationId: userOrganizationRelation.OrganizationId,
		CreateTime:     userOrganizationRelation.CreateTime,
		TutorId:        userOrganizationRelation.OrganizationTutorId,
	}

	companyOrganization, err := uoruc.corepo.Get(ctx, userOrganizationRelation.OrganizationId)

	if err != nil {
		return nil, WeixinCompanyOrganizationNotFound
	}

	if userOrganizationRelation.OrganizationUserId > 0 {
		if parentUser, err := uoruc.urepo.Get(ctx, userOrganizationRelation.OrganizationUserId); err == nil {
			nickName := parentUser.NickName

			if len(nickName) == 0 {
				nickName = tool.FormatPhone(parentUser.Phone)
			}

			bindUserOrganizationRelationInfo.ParentNickName = nickName
		}
	} else {
		bindUserOrganizationRelationInfo.ParentNickName = companyOrganization.Data.OrganizationName
	}

	if userOrganizationRelation.OrganizationTutorId > 0 {
		if tutorUser, err := uoruc.urepo.Get(ctx, userOrganizationRelation.OrganizationTutorId); err == nil {
			nickName := tutorUser.NickName

			if len(nickName) == 0 {
				nickName = tool.FormatPhone(tutorUser.Phone)
			}

			bindUserOrganizationRelationInfo.TutorNickName = nickName
		}
	}

	mcns := make([]*domain.Mcn, 0)

	if len(companyOrganization.Data.OrganizationMcns) > 0 {
		organizationMcns := make([]string, 0)

		for _, organizationMcn := range companyOrganization.Data.OrganizationMcns {
			organizationMcns = append(organizationMcns, organizationMcn.OrganizationMcn)
		}

		if userOpenDouyins, err := uoruc.uodrepo.List(ctx, 0, 40, user.Id, ""); err == nil {
			for _, userOpenDouyin := range userOpenDouyins {
				if djAweme, err := uoruc.darepo.Get(ctx, userOpenDouyin.AccountId, strings.Join(organizationMcns, ",")); err == nil {
					isNotExist := true

					for _, lmcn := range mcns {
						if lmcn.Name == djAweme.Account {
							isNotExist = false

							break
						}
					}

					if isNotExist {
						bindStartTime := ""
						bindEndTime := ""

						if fbindStartTime, err := tool.StringToTime("2006-01-02", djAweme.BindStartTime); err == nil {
							bindStartTime = tool.TimeToString("2006/01/02", fbindStartTime)
						}

						if fbindEndTime, err := tool.StringToTime("2006-01-02", djAweme.BindEndTime); err == nil {
							bindEndTime = tool.TimeToString("2006/01/02", fbindEndTime)
						}

						mcns = append(mcns, &domain.Mcn{
							Name:          djAweme.Account,
							BindStartTime: bindStartTime,
							BindEndTime:   bindEndTime,
						})
					}
				}
			}
		}
	}

	bindUserOrganizationRelationInfo.Mcn = mcns

	return bindUserOrganizationRelationInfo, nil
}

func (uoruc *UserOrganizationRelationUsecase) ListParentUserOrganizationRelations(ctx context.Context, userId, organizationId uint64, relationType string) ([]*domain.ParentUserOrganizationRelation, error) {
	user, err := uoruc.urepo.Get(ctx, userId)

	if err != nil {
		return nil, WeixinLoginError
	}

	if _, err := uoruc.repo.GetByUserId(ctx, user.Id, 0, 0, "0"); err == nil {
		return nil, WeixinUserOrganizationRelationExist
	}

	companyOrganization, err := uoruc.corepo.Get(ctx, organizationId)

	if err != nil {
		return nil, WeixinCompanyOrganizationNotFound
	}

	list := make([]*domain.ParentUserOrganizationRelation, 0)
	userIds := make([]uint64, 0)

	userScanRecord, err := uoruc.usrrepo.Get(ctx, user.Id, organizationId, 1)

	if err != nil {
		list = append(list, &domain.ParentUserOrganizationRelation{
			ParentUserId:    0,
			ParentNickName:  companyOrganization.Data.OrganizationName,
			ParentAvatarUrl: companyOrganization.Data.OrganizationLogoUrl,
			ParentUserType:  "companyOrganization",
		})
	} else {
		if userScanRecord.OrganizationUserId > 0 {
			if parentUser, err := uoruc.urepo.Get(ctx, userScanRecord.OrganizationUserId); err == nil {
				isNotExist := true

				for _, luserId := range userIds {
					if luserId == parentUser.Id {
						isNotExist = false

						break
					}
				}

				if isNotExist {
					userIds = append(userIds, parentUser.Id)

					list = append(list, &domain.ParentUserOrganizationRelation{
						ParentUserId:    parentUser.Id,
						ParentNickName:  parentUser.NickName,
						ParentAvatarUrl: parentUser.AvatarUrl,
						ParentUserType:  "user",
					})
				}
			}
		} else {
			list = append(list, &domain.ParentUserOrganizationRelation{
				ParentUserId:    0,
				ParentNickName:  companyOrganization.Data.OrganizationName,
				ParentAvatarUrl: companyOrganization.Data.OrganizationLogoUrl,
				ParentUserType:  "organization",
			})
		}
	}

	if relationType == "coupon" {
		if userCoupon, err := uoruc.ucorepo.GetByPhone(ctx, organizationId, user.Phone, "", "2"); err == nil {
			if parentUser, err := uoruc.urepo.Get(ctx, userCoupon.UserId); err == nil {
				isNotExist := true

				for _, luserId := range userIds {
					if luserId == parentUser.Id {
						isNotExist = false

						break
					}
				}

				if isNotExist {
					userIds = append(userIds, parentUser.Id)

					list = append(list, &domain.ParentUserOrganizationRelation{
						ParentUserId:    parentUser.Id,
						ParentNickName:  parentUser.NickName,
						ParentAvatarUrl: parentUser.AvatarUrl,
						ParentUserType:  "user",
					})
				}
			}
		}
	}

	return list, nil
}

func (uoruc *UserOrganizationRelationUsecase) UpdateLevelUserOrganizationRelations(ctx context.Context, userId, organizationId uint64, level uint8) error {
	user, err := uoruc.urepo.Get(ctx, userId)

	if err != nil {
		return WeixinLoginError
	}

	if _, err := uoruc.corepo.Get(ctx, organizationId); err != nil {
		return WeixinCompanyOrganizationNotFound
	}

	inUserOrganizationRelation, err := uoruc.repo.GetByUserId(ctx, user.Id, organizationId, 0, "0")

	if err != nil {
		return WeixinUserOrganizationRelationNotFound
	}

	if level == inUserOrganizationRelation.Level {
		return nil
	}

	inUserIntegralRelation, err := uoruc.uirrepo.GetByUserId(ctx, user.Id, organizationId, 0)

	if err != nil {
		return WeixinUserOrganizationRelationNotFound
	}

	oldLevel := inUserOrganizationRelation.Level

	userIntegralRelations, err := uoruc.uirrepo.List(ctx, organizationId)

	if err != nil {
		return WeixinUserOrganizationRelationListError
	}

	err = uoruc.tm.InTx(ctx, func(ctx context.Context) error {
		inUserOrganizationRelation.SetLevel(ctx, level)
		inUserOrganizationRelation.SetUpdateTime(ctx)

		if level == 4 {
			inUserOrganizationRelation.SetOrganizationTutorId(ctx, 0)
		}

		inUserOrganizationRelation, err = uoruc.repo.Update(ctx, inUserOrganizationRelation)

		if err != nil {
			return err
		}

		inUserIntegralRelation.SetLevel(ctx, level)
		inUserIntegralRelation.SetUpdateTime(ctx)

		if _, err = uoruc.uirrepo.Update(ctx, inUserIntegralRelation); err != nil {
			return err
		}

		if level == 4 {
			tmpChildIds := make([]uint64, 0)
			childIds := make([]uint64, 0)

			uoruc.uirrepo.ListChildId(ctx, user.Id, &tmpChildIds, userIntegralRelations)

			if len(tmpChildIds) > 0 {
				userOrganizationRelations, err := uoruc.repo.List(ctx, organizationId)

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
					if err := uoruc.repo.UpdateSuperior(ctx, user.Id, childIds); err != nil {
						return err
					}
				}
			}
		} else {
			if oldLevel == 4 {
				childIds := make([]uint64, 0)

				uoruc.uirrepo.ListChildId(ctx, user.Id, &childIds, userIntegralRelations)

				childIds = append(childIds, user.Id)

				var tutorUserId uint64 = 0

				if inUserOrganizationRelation.OrganizationUserId > 0 {
					tutorUserIntegralRelation := uoruc.uirrepo.GetSuperior(ctx, inUserOrganizationRelation.OrganizationUserId, 4, userIntegralRelations)

					if tutorUserIntegralRelation != nil {
						tutorUserId = tutorUserIntegralRelation.UserId
					}
				}

				if len(childIds) > 0 {
					if err := uoruc.repo.UpdateSuperior(ctx, tutorUserId, childIds); err != nil {
						return err
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return WeixinUserOrganizationRelationUpdateError
	}

	return nil
}

func (uoruc *UserOrganizationRelationUsecase) UpdateTutorUserOrganizationRelations(ctx context.Context) error {
	var wg sync.WaitGroup

	companyOrganizations, err := uoruc.corepo.List(ctx)

	if err != nil {
		return err
	}

	for _, companyOrganization := range companyOrganizations.Data.List {
		wg.Add(1)

		go uoruc.UpdateTutorUserOrganizationRelation(ctx, &wg, companyOrganization)
	}

	wg.Wait()

	return nil
}

func (uoruc *UserOrganizationRelationUsecase) UpdateTutorUserOrganizationRelation(ctx context.Context, wg *sync.WaitGroup, companyOrganization *v1.ListCompanyOrganizationsReply_CompanyOrganization) {
	defer wg.Done()

	userIntegralRelations, uirerr := uoruc.uirrepo.List(ctx, companyOrganization.OrganizationId)
	userOrganizationRelations, uorerr := uoruc.repo.List(ctx, companyOrganization.OrganizationId)

	if uirerr == nil && uorerr == nil {
		for _, inUserOrganizationRelation := range userOrganizationRelations {
			if inUserOrganizationRelation.Level != 4 {
				childIds := make([]uint64, 0)
				childIds = append(childIds, inUserOrganizationRelation.UserId)

				var tutorUserId uint64 = 0

				if inUserOrganizationRelation.OrganizationUserId > 0 {
					tutorUserIntegralRelation := uoruc.uirrepo.GetSuperior(ctx, inUserOrganizationRelation.OrganizationUserId, 4, userIntegralRelations)

					if tutorUserIntegralRelation != nil {
						tutorUserId = tutorUserIntegralRelation.UserId
					}
				}

				uoruc.repo.UpdateSuperior(ctx, tutorUserId, childIds)
			}
		}
	}
}

func (uoruc *UserOrganizationRelationUsecase) SyncUserOrganizationRelations(ctx context.Context) error {
	doukeOrders, err := uoruc.dorepo.ListUserId(ctx)

	if err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "SyncOrganizationUsers", fmt.Sprintf("[SyncOrganizationUsersError] Description=%s", "获取抖客订单列表失败"))
		inTaskLog.SetCreateTime(ctx)

		uoruc.tlrepo.Save(ctx, inTaskLog)

		return err
	}

	for _, doukeOrder := range doukeOrders.Data.List {
		if _, err := uoruc.repo.GetByUserId(ctx, doukeOrder.UserId, 0, 0, ""); err != nil {
			if userScanRecord, err := uoruc.usrrepo.Get(ctx, doukeOrder.UserId, 0, 1); err == nil {
				var organizationTutorId uint64 = 0

				if userScanRecord.OrganizationUserId > 0 {
					if parentUserOrganizationRelation, err := uoruc.repo.GetByUserId(ctx, userScanRecord.OrganizationUserId, userScanRecord.OrganizationId, 0, "0"); err == nil {
						if parentUserOrganizationRelation.Level == 4 {
							organizationTutorId = parentUserOrganizationRelation.UserId
						} else {
							if userIntegralRelations, err := uoruc.uirrepo.List(ctx, userScanRecord.OrganizationId); err == nil {
								tutorUserIntegralRelation := uoruc.uirrepo.GetSuperior(ctx, parentUserOrganizationRelation.UserId, 4, userIntegralRelations)

								if tutorUserIntegralRelation != nil {
									organizationTutorId = tutorUserIntegralRelation.UserId
								}
							}
						}
					}
				}

				inUserOrganizationRelation := domain.NewUserOrganizationRelation(ctx, userScanRecord.UserId, userScanRecord.OrganizationId, userScanRecord.OrganizationUserId, organizationTutorId, 0, 1, "")
				inUserOrganizationRelation.SetCreateTime(ctx)
				inUserOrganizationRelation.SetUpdateTime(ctx)

				uoruc.repo.Save(ctx, inUserOrganizationRelation)
			}
		}
	}

	jinritemaiOrders, err := uoruc.jorepo.ListByPickExtra(ctx)

	if err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "SyncOrganizationUsers", fmt.Sprintf("[SyncOrganizationUsersError] Description=%s", "获取达人订单列表失败"))
		inTaskLog.SetCreateTime(ctx)

		uoruc.tlrepo.Save(ctx, inTaskLog)

		return err
	}

	clientKeyAndOpenIds := make([]*domain.UserOpenDouyin, 0)

	for _, jinritemaiOrder := range jinritemaiOrders.Data.List {
		clientKeyAndOpenIds = append(clientKeyAndOpenIds, &domain.UserOpenDouyin{
			ClientKey: jinritemaiOrder.ClientKey,
			OpenId:    jinritemaiOrder.OpenId,
		})
	}

	if len(clientKeyAndOpenIds) > 0 {
		list, err := uoruc.uodrepo.ListByClientKeyAndOpenId(ctx, 0, 40, clientKeyAndOpenIds, "")

		if err != nil {
			inTaskLog := domain.NewTaskLog(ctx, "SyncOrganizationUsers", fmt.Sprintf("[SyncOrganizationUsersError] Description=%s", "获取微信用户关联抖音用户列表失败"))
			inTaskLog.SetCreateTime(ctx)

			uoruc.tlrepo.Save(ctx, inTaskLog)

			return err
		}

		var wg sync.WaitGroup

		for _, l := range list {
			wg.Add(1)

			go func(l *domain.UserOpenDouyin) {
				defer wg.Done()

				if _, err := uoruc.repo.GetByUserId(ctx, l.UserId, 0, 0, ""); err != nil {
					if userScanRecord, err := uoruc.usrrepo.Get(ctx, l.UserId, 0, 1); err == nil {
						var organizationTutorId uint64 = 0

						if userScanRecord.OrganizationUserId > 0 {
							if parentUserOrganizationRelation, err := uoruc.repo.GetByUserId(ctx, userScanRecord.OrganizationUserId, userScanRecord.OrganizationId, 0, "0"); err == nil {
								if parentUserOrganizationRelation.Level == 4 {
									organizationTutorId = parentUserOrganizationRelation.UserId
								} else {
									if userIntegralRelations, err := uoruc.uirrepo.List(ctx, userScanRecord.OrganizationId); err == nil {
										tutorUserIntegralRelation := uoruc.uirrepo.GetSuperior(ctx, parentUserOrganizationRelation.UserId, 4, userIntegralRelations)

										if tutorUserIntegralRelation != nil {
											organizationTutorId = tutorUserIntegralRelation.UserId
										}
									}
								}
							}
						}

						inUserOrganizationRelation := domain.NewUserOrganizationRelation(ctx, userScanRecord.UserId, userScanRecord.OrganizationId, userScanRecord.OrganizationUserId, organizationTutorId, 0, 1, "")
						inUserOrganizationRelation.SetCreateTime(ctx)
						inUserOrganizationRelation.SetUpdateTime(ctx)

						uoruc.repo.Save(ctx, inUserOrganizationRelation)
					}
				}
			}(l)
		}

		wg.Wait()
	}

	return nil
}

package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"strings"
	"sync"
	"time"
	v1 "weixin/api/service/company/v1"
	"weixin/internal/conf"
	"weixin/internal/domain"
	"weixin/internal/pkg/tool"
)

var (
	WeixinUserCommissionListError            = errors.InternalServer("WEIXIN_USER_COMMISSION_LIST_ERROR", "微信用户分佣明细列表获取失败")
	WeixinUserCommissionOrderCreateError     = errors.InternalServer("WEIXIN_USER_COMMISSION_ORDER_CREATE_ERROR", "微信用户电商分佣创建失败")
	WeixinUserCommissionCostOrderCreateError = errors.InternalServer("WEIXIN_USER_COMMISSION_COST_ORDER_CREATE_ERROR", "微信用户成本购分佣创建失败")
	WeixinUserCommissionTaskCreateError      = errors.InternalServer("WEIXIN_USER_COMMISSION_TASK_CREATE_ERROR", "微信用户任务分佣创建失败")
)

type UserCommissionRepo interface {
	GetByOutTradeNo(context.Context, string) (*domain.UserCommission, error)
	GetByRelevanceId(context.Context, uint64, uint8) (*domain.UserCommission, error)
	List(context.Context, int, int, uint64, uint64, []uint64, uint8, string, string, string) ([]*domain.UserCommission, error)
	ListByRelevanceId(context.Context, uint64, uint64, []string) ([]*domain.UserCommission, error)
	ListTask(context.Context, string, string) ([]*domain.UserCommission, error)
	ListBalance(context.Context, int, int, uint64, uint8, string) ([]*domain.UserCommission, error)
	ListCashable(context.Context) ([]*domain.UserCommission, error)
	ListOperation(context.Context) ([]*domain.UserCommission, error)
	Count(context.Context, uint64, uint64, []uint64, uint8, string, string, string) (int64, error)
	CountBalance(context.Context, uint64, uint8, string) (int64, error)
	Statistics(context.Context, uint64, uint8, uint8, []uint8, []uint8) (*domain.UserCommission, error)
	StatisticsDetail(context.Context, uint64, uint64, []uint64, uint8, string, string, string) (*domain.UserCommission, error)
	Update(context.Context, *domain.UserCommission) (*domain.UserCommission, error)
	Save(context.Context, *domain.UserCommission) (*domain.UserCommission, error)
	DeleteByRelevanceId(context.Context, uint64, uint8) error
}

type UserCommissionUsecase struct {
	repo     UserCommissionRepo
	urepo    UserRepo
	uorerepo UserOrganizationRelationRepo
	uirrepo  UserIntegralRelationRepo
	uodrepo  UserOpenDouyinRepo
	uorrepo  UserOrderRepo
	ucrepo   UserCouponRepo
	usrrepo  UserScanRecordRepo
	ublrepo  UserBalanceLogRepo
	corepo   CompanyOrganizationRepo
	ctrepo   CompanyTaskRepo
	cprepo   CompanyProductRepo
	jorepo   JinritemaiOrderRepo
	dorepo   DoukeOrderRepo
	darepo   DjAwemeRepo
	tlrepo   TaskLogRepo
	tm       Transaction
	conf     *conf.Data
	cconf    *conf.Company
	oconf    *conf.Organization
	log      *log.Helper
}

func NewUserCommissionUsecase(repo UserCommissionRepo, urepo UserRepo, uorerepo UserOrganizationRelationRepo, uirrepo UserIntegralRelationRepo, uodrepo UserOpenDouyinRepo, uorrepo UserOrderRepo, ucrepo UserCouponRepo, usrrepo UserScanRecordRepo, ublrepo UserBalanceLogRepo, corepo CompanyOrganizationRepo, ctrepo CompanyTaskRepo, cprepo CompanyProductRepo, jorepo JinritemaiOrderRepo, dorepo DoukeOrderRepo, darepo DjAwemeRepo, tlrepo TaskLogRepo, tm Transaction, conf *conf.Data, cconf *conf.Company, oconf *conf.Organization, logger log.Logger) *UserCommissionUsecase {
	return &UserCommissionUsecase{repo: repo, urepo: urepo, uorerepo: uorerepo, uirrepo: uirrepo, uodrepo: uodrepo, uorrepo: uorrepo, ucrepo: ucrepo, usrrepo: usrrepo, ublrepo: ublrepo, corepo: corepo, ctrepo: ctrepo, cprepo: cprepo, jorepo: jorepo, dorepo: dorepo, darepo: darepo, tlrepo: tlrepo, tm: tm, conf: conf, cconf: cconf, oconf: oconf, log: log.NewHelper(logger)}
}

func (ucuc *UserCommissionUsecase) ListUserCommissions(ctx context.Context, pageNum, pageSize, userId, organizationId uint64, commissionType uint8, month, keyword string) (*domain.UserCommissionList, error) {
	user, err := ucuc.urepo.Get(ctx, userId)

	if err != nil {
		return nil, WeixinLoginError
	}

	if _, err := ucuc.uorerepo.GetByUserId(ctx, user.Id, organizationId, 0, "0"); err != nil {
		return nil, WeixinUserOrganizationRelationNotFound
	}

	startDay, endDay := "", ""

	if len(month) > 0 {
		startTime, endTime := tool.GetMonthStartTimeAndEndTime(month)

		startDay = startTime.Format("2006-01-02")
		endDay = endTime.Format("2006-01-02")
	}

	userOrganizationRelations, err := ucuc.uorerepo.List(ctx, organizationId)

	if err != nil {
		return nil, WeixinUserOrganizationRelationListError
	}

	childIds := make([]uint64, 0)

	ucuc.uorerepo.ListChildId(ctx, userId, &childIds, userOrganizationRelations)

	userCommissions, err := ucuc.repo.List(ctx, int(pageNum), int(pageSize), userId, organizationId, childIds, commissionType, startDay, endDay, keyword)

	if err != nil {
		return nil, WeixinUserCommissionListError
	}

	total, err := ucuc.repo.Count(ctx, userId, organizationId, childIds, commissionType, startDay, endDay, keyword)

	if err != nil {
		return nil, WeixinUserCommissionListError
	}

	list := make([]*domain.UserCommission, 0)

	for _, userCommission := range userCommissions {
		activationTime, relationName := "", ""

		for _, userOrganizationRelation := range userOrganizationRelations {
			if userOrganizationRelation.UserId == userCommission.ChildUserId {
				activationTime = tool.TimeToString("2006/01/02 15:04", userOrganizationRelation.CreateTime)

				if userOrganizationRelation.OrganizationUserId == user.Id {
					relationName = "直接"
				} else {
					relationName = "间接"
				}

				break
			}
		}

		var commissionRatio float32 = 0.00

		if userCommission.CommissionPool > 0 {
			commissionRatio = userCommission.CommissionAmount / userCommission.CommissionPool
		}

		list = append(list, &domain.UserCommission{
			ChildUserId:        userCommission.ChildUserId,
			ChildNickName:      userCommission.ChildNickName,
			ChildAvatarUrl:     userCommission.ChildAvatarUrl,
			ChildPhone:         userCommission.ChildPhone,
			RelationName:       relationName,
			TotalPayAmount:     userCommission.TotalPayAmount,
			CommissionPool:     userCommission.CommissionPool,
			CommissionAmount:   userCommission.CommissionAmount,
			CommissionRatio:    commissionRatio * 100,
			CommissionType:     userCommission.CommissionType,
			CommissionTypeName: userCommission.GetCommissionTypeName(ctx),
			ActivationTime:     activationTime,
		})
	}

	return &domain.UserCommissionList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (ucuc *UserCommissionUsecase) StatisticsUserCommissions(ctx context.Context, userId, organizationId uint64) (*domain.StatisticsUserOrganizationRelations, error) {
	statistics := make([]*domain.StatisticsUserOrganizationRelation, 0)

	user, err := ucuc.urepo.Get(ctx, userId)

	if err != nil {
		return nil, WeixinLoginError
	}

	userOrganizationRelation, err := ucuc.uorerepo.GetByUserId(ctx, user.Id, organizationId, 0, "0")

	if err != nil {
		return nil, WeixinUserOrganizationRelationNotFound
	}

	levelName := WeixinUserOrganizationRelationLevel[userOrganizationRelation.Level-1]

	statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
		Key:   "基本身份",
		Value: levelName,
	})

	var wg sync.WaitGroup
	var presenterNum int64 = 0
	var childNum uint64 = 0
	var userOrganizationRelations []*domain.UserOrganizationRelation
	var statisticsCourse, statisticsOrder, statisticsCostOrder *domain.UserCommission
	var uerr, scerr, soerr, scoerr error

	wg.Add(5)

	go func() {
		defer wg.Done()

		presenterNum, _ = ucuc.uorerepo.Count(ctx, userId, organizationId, "")
	}()

	go func() {
		defer wg.Done()

		userOrganizationRelations, uerr = ucuc.uorerepo.List(ctx, organizationId)
	}()

	go func() {
		defer wg.Done()

		statisticsCourse, scerr = ucuc.repo.Statistics(ctx, user.Id, 0, 1, []uint8{1}, []uint8{1})
	}()

	go func() {
		defer wg.Done()

		statisticsOrder, soerr = ucuc.repo.Statistics(ctx, user.Id, 0, 1, []uint8{2}, []uint8{1})
	}()

	go func() {
		defer wg.Done()

		statisticsCostOrder, scoerr = ucuc.repo.Statistics(ctx, user.Id, 0, 1, []uint8{4}, []uint8{1})
	}()

	wg.Wait()

	if uerr == nil {
		ucuc.uorerepo.GetChildNum(ctx, userId, &childNum, userOrganizationRelations)
	}

	statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
		Key:   "基本直推人数",
		Value: strconv.FormatUint(uint64(presenterNum), 10),
	})

	statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
		Key:   "基本总人数",
		Value: strconv.FormatUint(childNum, 10),
	})

	if scerr == nil {
		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "会员销额",
			Value: fmt.Sprintf("%.2f", tool.Decimal(float64(statisticsCourse.TotalPayAmount), 2)),
		})

		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "会员分佣池",
			Value: fmt.Sprintf("%.2f", tool.Decimal(float64(statisticsCourse.CommissionPool), 2)),
		})

		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "会员总分佣",
			Value: fmt.Sprintf("%.2f", tool.Decimal(float64(statisticsCourse.CommissionAmount), 2)),
		})
	} else {
		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "会员销额",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "会员分佣池",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "会员总分佣",
			Value: "0.00",
		})
	}

	if scoerr == nil {
		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "成本购销额",
			Value: fmt.Sprintf("%.2f", tool.Decimal(float64(statisticsCostOrder.TotalPayAmount), 2)),
		})

		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "成本购分佣池",
			Value: fmt.Sprintf("%.2f", tool.Decimal(float64(statisticsCostOrder.CommissionPool), 2)),
		})

		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "成本购总分佣",
			Value: fmt.Sprintf("%.2f", tool.Decimal(float64(statisticsCostOrder.CommissionAmount), 2)),
		})
	} else {
		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "成本购销额",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "成本购分佣池",
			Value: "0.00",
		})
		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "成本购总分佣",
			Value: "0.00",
		})

	}

	if soerr == nil {
		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "带货销额",
			Value: fmt.Sprintf("%.2f", tool.Decimal(float64(statisticsOrder.TotalPayAmount), 2)),
		})

		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "带货分佣池",
			Value: fmt.Sprintf("%.2f", tool.Decimal(float64(statisticsOrder.CommissionPool), 2)),
		})

		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "带货总分佣",
			Value: fmt.Sprintf("%.2f", tool.Decimal(float64(statisticsOrder.CommissionAmount), 2)),
		})
	} else {
		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "带货销额",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "带货分佣池",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "带货总分佣",
			Value: "0.00",
		})
	}

	return &domain.StatisticsUserOrganizationRelations{
		Statistics: statistics,
	}, nil
}

func (ucuc *UserCommissionUsecase) StatisticsDetailUserCommissions(ctx context.Context, userId, organizationId uint64, commissionType uint8, month, keyword string) (*domain.StatisticsUserOrganizationRelations, error) {
	statistics := make([]*domain.StatisticsUserOrganizationRelation, 0)

	user, err := ucuc.urepo.Get(ctx, userId)

	if err != nil {
		return nil, WeixinLoginError
	}

	if _, err := ucuc.uorerepo.GetByUserId(ctx, user.Id, organizationId, 0, "0"); err != nil {
		return nil, WeixinUserOrganizationRelationNotFound
	}

	startDay, endDay := "", ""

	if len(month) > 0 {
		startTime, endTime := tool.GetMonthStartTimeAndEndTime(month)

		startDay = startTime.Format("2006-01-02")
		endDay = endTime.Format("2006-01-02")
	}

	userOrganizationRelations, err := ucuc.uorerepo.List(ctx, organizationId)

	if err != nil {
		return nil, WeixinUserOrganizationRelationListError
	}

	childIds := make([]uint64, 0)

	ucuc.uorerepo.ListChildId(ctx, userId, &childIds, userOrganizationRelations)

	if userCommission, err := ucuc.repo.StatisticsDetail(ctx, user.Id, organizationId, childIds, commissionType, startDay, endDay, keyword); err == nil {
		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "销售金额",
			Value: fmt.Sprintf("%.2f", tool.Decimal(float64(userCommission.TotalPayAmount), 2)),
		})

		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "分佣池",
			Value: fmt.Sprintf("%.2f", tool.Decimal(float64(userCommission.CommissionPool), 2)),
		})

		var commissionRatio float32 = 0.00

		if userCommission.CommissionPool > 0 {
			commissionRatio = userCommission.CommissionAmount / userCommission.CommissionPool
		}

		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "分佣比例",
			Value: fmt.Sprintf("%.2f", tool.Decimal(float64(commissionRatio*100), 2)),
		})

		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "总分佣",
			Value: fmt.Sprintf("%.2f", tool.Decimal(float64(userCommission.CommissionAmount), 2)),
		})
	} else {
		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "销售金额",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "分佣池",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "分佣比例",
			Value: "0.00",
		})

		statistics = append(statistics, &domain.StatisticsUserOrganizationRelation{
			Key:   "总分佣",
			Value: "0.00",
		})
	}

	return &domain.StatisticsUserOrganizationRelations{
		Statistics: statistics,
	}, nil
}

func (ucuc *UserCommissionUsecase) CreateOrderUserCommissions(ctx context.Context, totalPayAmount, commission float64, clientKey, openId, orderId, flowPoint, paySuccessTime string) error {
	relevanceId, err := strconv.ParseUint(orderId, 10, 64)

	if err != nil {
		return WeixinUserCommissionOrderCreateError
	}

	tpaySuccessTime, err := tool.StringToTime("2006-01-02 15:04:05", paySuccessTime)

	if err != nil {
		return WeixinUserCommissionOrderCreateError
	}

	tendTime, err := tool.StringToTime("2006-01-02 15:04:05", "2024-04-01 00:00:00")

	if err != nil {
		return WeixinUserCommissionOrderCreateError
	}

	if err != nil {
		return WeixinUserCommissionOrderCreateError
	}

	if flowPoint == "REFUND" {
		err := ucuc.repo.DeleteByRelevanceId(ctx, relevanceId, 2)

		if err != nil {
			return WeixinUserCommissionOrderCreateError
		}

		return nil
	}

	userOpenDouyin, err := ucuc.uodrepo.GetByClientKeyAndOpenId(ctx, clientKey, openId)

	if err != nil {
		return WeixinUserOpenDouyinNotFound
	}

	userCommissions, err := ucuc.repo.ListByRelevanceId(ctx, 0, relevanceId, []string{"2"})

	if err != nil {
		return WeixinUserCommissionOrderCreateError
	}

	if len(userCommissions) == 0 {
		var commissionStatus uint8

		if flowPoint == "SETTLE" {
			commissionStatus = 2
		} else {
			commissionStatus = 1
		}

		userOrganizationRelation, err := ucuc.uorerepo.GetByUserId(ctx, userOpenDouyin.UserId, 0, 0, "0")

		if err != nil {
			return WeixinUserOrganizationRelationNotFound
		}

		companyOrganization, err := ucuc.corepo.Get(ctx, userOrganizationRelation.OrganizationId)

		if err != nil {
			return WeixinCompanyOrganizationNotFound
		}

		organizationMcns := make([]string, 0)

		for _, organizationMcn := range companyOrganization.Data.OrganizationMcns {
			organizationMcns = append(organizationMcns, organizationMcn.OrganizationMcn)
		}

		if len(organizationMcns) > 0 && len(userOpenDouyin.AccountId) > 0 {
			djAweme, err := ucuc.darepo.Get(ctx, userOpenDouyin.AccountId, strings.Join(organizationMcns, ","))

			if err != nil {
				return WeixinDjAwemeNotFound
			}

			if ratio, err := strconv.ParseFloat(djAweme.Ratio, 64); err == nil {
				if bindStartTime, err := tool.StringToTime("2006-01-02", djAweme.BindStartTime); err == nil {
					if ratio > 0.00 && bindStartTime.Add(24*time.Hour).Before(tpaySuccessTime) {
						err = ucuc.tm.InTx(ctx, func(ctx context.Context) error {
							if tpaySuccessTime.Before(tendTime) {
								if ratio == 70 {
									ratio = 90
								}
							}

							return ucuc.getOrderComission(ctx, relevanceId, commissionStatus, totalPayAmount, commission, ratio, tpaySuccessTime, userOrganizationRelation, companyOrganization)
						})

						if err != nil {
							return WeixinUserCommissionOrderCreateError
						}
					}
				}
			}
		}
	} else {
		var commissionStatus uint8
		var organizationId uint64
		var parentLevel uint8

		for _, userCommission := range userCommissions {
			commissionStatus = userCommission.CommissionStatus
			organizationId = userCommission.OrganizationId

			if userCommission.Relation == 1 {
				parentLevel = userCommission.Level

				break
			}
		}

		if commissionStatus == 1 && flowPoint == "SETTLE" {
			companyOrganization, err := ucuc.corepo.Get(ctx, organizationId)

			if err != nil {
				return WeixinCompanyOrganizationNotFound
			}

			err = ucuc.tm.InTx(ctx, func(ctx context.Context) error {
				for _, inUserCommission := range userCommissions {
					commissionPool := commission * 100 / float64(inUserCommission.CommissionMcnRatio) * companyOrganization.Data.OrganizationColonelCommission.OrderRatio / 100

					var realCommission float64

					if inUserCommission.Relation == 1 {
						if inUserCommission.Level == 2 {
							realCommission = commissionPool * companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterOrderCommissionRule / 100
						} else if inUserCommission.Level == 3 {
							realCommission = commissionPool * companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterOrderCommissionRule / 100
						} else if inUserCommission.Level == 4 {
							realCommission = commissionPool * companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterOrderCommissionRule / 100
						}
					} else if inUserCommission.Relation == 2 {
						if inUserCommission.Level == 4 {
							if parentLevel == 2 {
								realCommission = commissionPool * companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorOrderCommissionRule / 100
							} else if parentLevel == 3 {
								realCommission = commissionPool * companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorOrderCommissionRule / 100
							}
						}
					}

					inUserCommission.SetCommissionStatus(ctx, 2)
					inUserCommission.SetCommissionPool(ctx, float32(commissionPool))
					inUserCommission.SetCommissionAmount(ctx, float32(tool.Decimal(realCommission, 2)))

					if _, err := ucuc.repo.Update(ctx, inUserCommission); err != nil {
						return err
					}
				}

				return nil
			})

			if err != nil {
				return WeixinUserCommissionOrderCreateError
			}
		}
	}

	return nil
}

func (ucuc *UserCommissionUsecase) getOrderComission(ctx context.Context, relevanceId uint64, commissionStatus uint8, totalPayAmount, commission, ratio float64, createTime time.Time, userOrganizationRelation *domain.UserOrganizationRelation, companyOrganization *v1.GetCompanyOrganizationsReply) error {
	var organizationParentUser *domain.UserOrganizationRelation
	var organizationTutorUser *domain.UserOrganizationRelation
	var err error

	if userOrganizationRelation.OrganizationUserId > 0 {
		organizationParentUser, err = ucuc.uorerepo.GetByUserId(ctx, userOrganizationRelation.OrganizationUserId, userOrganizationRelation.OrganizationId, 0, "0")

		if err != nil {
			return WeixinUserOrganizationRelationNotFound
		}

		if userOrganizationRelation.OrganizationTutorId > 0 {
			organizationTutorUser, err = ucuc.uorerepo.GetByUserId(ctx, userOrganizationRelation.OrganizationTutorId, userOrganizationRelation.OrganizationId, 0, "0")

			if err != nil {
				return WeixinUserOrganizationRelationNotFound
			}
		}

		commissionPool := commission * 100 / ratio * companyOrganization.Data.OrganizationColonelCommission.OrderRatio / 100

		var realCommission float64

		if organizationParentUser.Level == 2 {
			realCommission = commissionPool * companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterOrderCommissionRule / 100
		} else if organizationParentUser.Level == 3 {
			realCommission = commissionPool * companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterOrderCommissionRule / 100
		} else if organizationParentUser.Level == 4 {
			realCommission = commissionPool * companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterOrderCommissionRule / 100
		}

		inUserCommission := domain.NewUserCommission(ctx, organizationParentUser.UserId, organizationParentUser.OrganizationId, userOrganizationRelation.UserId, relevanceId, userOrganizationRelation.Level, organizationParentUser.Level, 1, commissionStatus, 2, 1, 1, float32(tool.Decimal(totalPayAmount, 2)), float32(commissionPool), float32(tool.Decimal(realCommission, 2)))
		inUserCommission.SetCommissionMcnRatio(ctx, float32(ratio))
		inUserCommission.SetCreateTime(ctx, createTime)
		inUserCommission.SetUpdateTime(ctx)

		if _, err := ucuc.repo.Save(ctx, inUserCommission); err != nil {
			return err
		}

		if organizationTutorUser != nil {
			if organizationParentUser.Level == 2 {
				if organizationTutorUser.Level == 4 {
					realCommission = commissionPool * companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorOrderCommissionRule / 100

					inUserCommission = domain.NewUserCommission(ctx, organizationTutorUser.UserId, organizationTutorUser.OrganizationId, userOrganizationRelation.UserId, relevanceId, userOrganizationRelation.Level, organizationTutorUser.Level, 2, commissionStatus, 2, 1, 1, float32(tool.Decimal(totalPayAmount, 2)), float32(commissionPool), float32(tool.Decimal(realCommission, 2)))
					inUserCommission.SetCommissionMcnRatio(ctx, float32(ratio))
					inUserCommission.SetCreateTime(ctx, createTime)
					inUserCommission.SetUpdateTime(ctx)

					if _, err := ucuc.repo.Save(ctx, inUserCommission); err != nil {
						return err
					}
				}
			} else if organizationParentUser.Level == 3 {
				if organizationTutorUser.Level == 4 {
					realCommission = commissionPool * companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorOrderCommissionRule / 100

					inUserCommission = domain.NewUserCommission(ctx, organizationTutorUser.UserId, organizationTutorUser.OrganizationId, userOrganizationRelation.UserId, relevanceId, userOrganizationRelation.Level, organizationTutorUser.Level, 2, commissionStatus, 2, 1, 1, float32(tool.Decimal(totalPayAmount, 2)), float32(commissionPool), float32(tool.Decimal(realCommission, 2)))
					inUserCommission.SetCommissionMcnRatio(ctx, float32(ratio))
					inUserCommission.SetCreateTime(ctx, createTime)
					inUserCommission.SetUpdateTime(ctx)

					if _, err := ucuc.repo.Save(ctx, inUserCommission); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (ucuc *UserCommissionUsecase) CreateCostOrderUserCommissions(ctx context.Context, userId uint64, totalPayAmount, commission float64, orderId, productId, flowPoint, paySuccessTime string) error {
	relevanceId, err := strconv.ParseUint(orderId, 10, 64)

	if err != nil {
		return WeixinUserCommissionCostOrderCreateError
	}

	if flowPoint == "REFUND" {
		err = ucuc.tm.InTx(ctx, func(ctx context.Context) error {
			err = ucuc.repo.DeleteByRelevanceId(ctx, relevanceId, 3)

			if err != nil {
				return err
			}

			err = ucuc.repo.DeleteByRelevanceId(ctx, relevanceId, 4)

			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return WeixinUserCommissionCostOrderCreateError
		}

		return nil
	}

	user, err := ucuc.urepo.Get(ctx, userId)

	if err != nil {
		return WeixinUserNotFound
	}

	userCommissions, err := ucuc.repo.ListByRelevanceId(ctx, user.Id, relevanceId, []string{"3", "4"})

	if err != nil {
		return WeixinUserCommissionOrderCreateError
	}

	createTime, _ := tool.StringToTime("2006-01-02 15:04:05", paySuccessTime)

	if len(userCommissions) == 0 {
		var commissionStatus uint8

		if flowPoint == "SETTLE" {
			commissionStatus = 2
		} else {
			commissionStatus = 1
		}

		userOrganizationRelation, err := ucuc.uorerepo.GetByUserId(ctx, user.Id, 0, 0, "")

		if err != nil {
			if userScanRecord, err := ucuc.usrrepo.Get(ctx, user.Id, 0, 1); err == nil {
				var organizationTutorId uint64 = 0

				if userScanRecord.OrganizationUserId > 0 {
					if parentUserOrganizationRelation, err := ucuc.uorerepo.GetByUserId(ctx, userScanRecord.OrganizationUserId, userScanRecord.OrganizationId, 0, "0"); err == nil {
						if parentUserOrganizationRelation.Level == 4 {
							organizationTutorId = parentUserOrganizationRelation.UserId
						} else {
							if userIntegralRelations, err := ucuc.uirrepo.List(ctx, userScanRecord.OrganizationId); err == nil {
								tutorUserIntegralRelation := ucuc.uirrepo.GetSuperior(ctx, parentUserOrganizationRelation.UserId, 4, userIntegralRelations)

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

				userOrganizationRelation, err = ucuc.uorerepo.Save(ctx, inUserOrganizationRelation)

				if err != nil {
					return WeixinCompanyOrganizationNotFound
				}
			} else {
				return WeixinUserScanRecordNotFound
			}
		}

		companyOrganization, err := ucuc.corepo.Get(ctx, userOrganizationRelation.OrganizationId)

		if err != nil {
			return WeixinCompanyOrganizationNotFound
		}

		err = ucuc.tm.InTx(ctx, func(ctx context.Context) error {
			return ucuc.getCostOrderComission(ctx, relevanceId, commissionStatus, totalPayAmount, commission, createTime, userOrganizationRelation, companyOrganization)
		})

		if err != nil {
			return WeixinUserCommissionCostOrderCreateError
		}
	} else {
		var commissionStatus uint8
		var organizationId uint64
		var parentLevel uint8

		for _, userCommission := range userCommissions {
			commissionStatus = userCommission.CommissionStatus
			organizationId = userCommission.OrganizationId

			if userCommission.CommissionType == 4 && userCommission.Relation == 1 {
				parentLevel = userCommission.Level

				break
			}
		}

		if commissionStatus == 1 && flowPoint == "SETTLE" {
			companyOrganization, err := ucuc.corepo.Get(ctx, organizationId)

			if err != nil {
				return WeixinCompanyOrganizationNotFound
			}

			err = ucuc.tm.InTx(ctx, func(ctx context.Context) error {
				for _, inUserCommission := range userCommissions {
					if inUserCommission.CommissionType == 3 {
						commissionPool := commission * 0.6
						realCommission := commission * 0.6

						inUserCommission.SetCommissionStatus(ctx, 2)
						inUserCommission.SetCommissionPool(ctx, float32(commissionPool))
						inUserCommission.SetCommissionAmount(ctx, float32(tool.Decimal(realCommission, 2)))

						if _, err := ucuc.repo.Update(ctx, inUserCommission); err != nil {
							return err
						}
					} else {
						commissionPool := commission * companyOrganization.Data.OrganizationColonelCommission.CostOrderRatio / 100

						var realCommission float64

						if inUserCommission.Relation == 1 {
							if inUserCommission.Level == 1 {
								realCommission = commissionPool * companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedPresenterCostOrderCommissionRule / 100
							} else if inUserCommission.Level == 2 {
								realCommission = commissionPool * companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterCostOrderCommissionRule / 100
							} else if inUserCommission.Level == 3 {
								realCommission = commissionPool * companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterCostOrderCommissionRule / 100
							} else if inUserCommission.Level == 4 {
								realCommission = commissionPool * companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterCostOrderCommissionRule / 100
							}
						} else if inUserCommission.Relation == 2 {
							if parentLevel == 1 {
								realCommission = commissionPool * companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedTutorCostOrderCommissionRule / 100
							} else if parentLevel == 2 {
								realCommission = commissionPool * companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorCostOrderCommissionRule / 100
							} else if parentLevel == 3 {
								realCommission = commissionPool * companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorCostOrderCommissionRule / 100
							}
						}

						inUserCommission.SetCommissionStatus(ctx, 2)
						inUserCommission.SetCommissionPool(ctx, float32(commissionPool))
						inUserCommission.SetCommissionAmount(ctx, float32(tool.Decimal(realCommission, 2)))

						if _, err := ucuc.repo.Update(ctx, inUserCommission); err != nil {
							return err
						}
					}
				}

				return nil
			})

			if err != nil {
				return WeixinUserCommissionCostOrderCreateError
			}
		}
	}

	return nil
}

func (ucuc *UserCommissionUsecase) getCostOrderComission(ctx context.Context, relevanceId uint64, commissionStatus uint8, totalPayAmount, commission float64, createTime time.Time, userOrganizationRelation *domain.UserOrganizationRelation, companyOrganization *v1.GetCompanyOrganizationsReply) error {
	var organizationParentUser *domain.UserOrganizationRelation
	var organizationTutorUser *domain.UserOrganizationRelation
	var err error

	commissionPool := commission * 0.6
	realCommission := commission * 0.6

	inUserCommission := domain.NewUserCommission(ctx, userOrganizationRelation.UserId, userOrganizationRelation.OrganizationId, userOrganizationRelation.UserId, relevanceId, userOrganizationRelation.Level, userOrganizationRelation.Level, 1, commissionStatus, 3, 1, 1, float32(tool.Decimal(totalPayAmount, 2)), float32(commissionPool), float32(tool.Decimal(realCommission, 2)))
	inUserCommission.SetCreateTime(ctx, createTime)
	inUserCommission.SetUpdateTime(ctx)

	if _, err := ucuc.repo.Save(ctx, inUserCommission); err != nil {
		return err
	}

	if userOrganizationRelation.OrganizationUserId > 0 {
		organizationParentUser, err = ucuc.uorerepo.GetByUserId(ctx, userOrganizationRelation.OrganizationUserId, userOrganizationRelation.OrganizationId, 0, "0")

		if err != nil {
			return WeixinUserOrganizationRelationNotFound
		}

		if userOrganizationRelation.IsOrderRelation == 1 {
			if organizationParentUser.OrganizationTutorId > 0 {
				organizationTutorUser, err = ucuc.uorerepo.GetByUserId(ctx, organizationParentUser.OrganizationTutorId, userOrganizationRelation.OrganizationId, 0, "0")

				if err != nil {
					return WeixinUserOrganizationRelationNotFound
				}
			}
		} else {
			if userOrganizationRelation.OrganizationTutorId > 0 {
				organizationTutorUser, err = ucuc.uorerepo.GetByUserId(ctx, userOrganizationRelation.OrganizationTutorId, userOrganizationRelation.OrganizationId, 0, "0")

				if err != nil {
					return WeixinUserOrganizationRelationNotFound
				}
			}
		}

		commissionPool = commission * companyOrganization.Data.OrganizationColonelCommission.CostOrderRatio / 100

		if organizationParentUser.Level == 1 {
			realCommission = commissionPool * companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedPresenterCostOrderCommissionRule / 100
		} else if organizationParentUser.Level == 2 {
			realCommission = commissionPool * companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedPresenterCostOrderCommissionRule / 100
		} else if organizationParentUser.Level == 3 {
			realCommission = commissionPool * companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedPresenterCostOrderCommissionRule / 100
		} else if organizationParentUser.Level == 4 {
			realCommission = commissionPool * companyOrganization.Data.OrganizationColonelCommission.AdvancedPresenterCostOrderCommissionRule / 100
		}

		inUserCommission = domain.NewUserCommission(ctx, organizationParentUser.UserId, organizationParentUser.OrganizationId, userOrganizationRelation.UserId, relevanceId, userOrganizationRelation.Level, organizationParentUser.Level, 1, commissionStatus, 4, 1, 1, float32(tool.Decimal(totalPayAmount, 2)), float32(commissionPool), float32(tool.Decimal(realCommission, 2)))
		inUserCommission.SetCreateTime(ctx, createTime)
		inUserCommission.SetUpdateTime(ctx)

		if _, err := ucuc.repo.Save(ctx, inUserCommission); err != nil {
			return err
		}

		if organizationTutorUser != nil {
			if organizationTutorUser.Level == 4 && organizationParentUser.Level != 4 {
				if organizationParentUser.Level == 1 {
					realCommission = commissionPool * companyOrganization.Data.OrganizationColonelCommission.ZeroAdvancedTutorCostOrderCommissionRule / 100
				} else if organizationParentUser.Level == 2 {
					realCommission = commissionPool * companyOrganization.Data.OrganizationColonelCommission.PrimaryAdvancedTutorCostOrderCommissionRule / 100
				} else if organizationParentUser.Level == 3 {
					realCommission = commissionPool * companyOrganization.Data.OrganizationColonelCommission.IntermediateAdvancedTutorCostOrderCommissionRule / 100
				}

				inUserCommission = domain.NewUserCommission(ctx, organizationTutorUser.UserId, organizationTutorUser.OrganizationId, userOrganizationRelation.UserId, relevanceId, userOrganizationRelation.Level, organizationTutorUser.Level, 2, commissionStatus, 4, 1, 1, float32(tool.Decimal(totalPayAmount, 2)), float32(commissionPool), float32(tool.Decimal(realCommission, 2)))
				inUserCommission.SetCreateTime(ctx, createTime)
				inUserCommission.SetUpdateTime(ctx)

				if _, err := ucuc.repo.Save(ctx, inUserCommission); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (ucuc *UserCommissionUsecase) CreateTaskUserCommissions(ctx context.Context, userId, taskRelationId uint64, commission float64, flowPoint, successTime string) error {
	user, err := ucuc.urepo.Get(ctx, userId)

	if err != nil {
		return WeixinUserNotFound
	}

	if flowPoint == "REFUND" {
		err := ucuc.repo.DeleteByRelevanceId(ctx, taskRelationId, 5)

		if err != nil {
			return WeixinUserCommissionTaskCreateError
		}

		return nil
	}

	createTime, _ := tool.StringToTime("2006-01-02 15:04:05", successTime)

	userCommissions, err := ucuc.repo.ListByRelevanceId(ctx, user.Id, taskRelationId, []string{"5"})

	if err != nil {
		return WeixinUserCommissionTaskCreateError
	}

	if len(userCommissions) == 0 {
		inUserCommission := domain.NewUserCommission(ctx, user.Id, 0, user.Id, taskRelationId, 0, 0, 1, 1, 5, 1, 1, float32(tool.Decimal(commission, 2)), float32(commission), float32(tool.Decimal(commission, 2)))
		inUserCommission.SetCreateTime(ctx, createTime)
		inUserCommission.SetUpdateTime(ctx)

		if _, err := ucuc.repo.Save(ctx, inUserCommission); err != nil {
			return WeixinUserCommissionTaskCreateError
		}
	}

	return nil
}

func (ucuc *UserCommissionUsecase) SyncTaskUserCommissions(ctx context.Context) error {
	userCommissions, err := ucuc.repo.ListTask(ctx, "1", "5")

	if err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "SyncTaskUserCommissions", fmt.Sprintf("[SyncTaskUserCommissionsError], Description=%s", "获取用户任务佣金列表失败"))
		inTaskLog.SetCreateTime(ctx)

		ucuc.tlrepo.Save(ctx, inTaskLog)

		return err
	}

	for _, inUserCommission := range userCommissions {
		if comapnyTaskAccountRelation, err := ucuc.ctrepo.GetCompanyTaskAccountRelation(ctx, inUserCommission.RelevanceId); err == nil {
			if statisticses, err := ucuc.dorepo.StatisticsByPaySuccessTime(ctx, comapnyTaskAccountRelation.Data.UserId, comapnyTaskAccountRelation.Data.ProductOutId, "SETTLE", comapnyTaskAccountRelation.Data.ClaimTime, comapnyTaskAccountRelation.Data.ExpireTime); err == nil {
				for _, statistics := range statisticses.Data.Statistics {
					if statistics.Key == "orderNum" {
						if orderNum, err := strconv.ParseUint(statistics.Value, 10, 64); err == nil {
							if orderNum > 0 {
								inUserCommission.SetCommissionStatus(ctx, 2)
								inUserCommission.SetUpdateTime(ctx)

								ucuc.repo.Update(ctx, inUserCommission)
							}
						}
					}
				}
			}
		}
	}

	return nil
}

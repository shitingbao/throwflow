package biz

import (
	"context"
	"encoding/json"
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
	WeixinUserCommissionListError = errors.InternalServer("WEIXIN_USER_COMMISSION_LIST_ERROR", "微信用户分佣明细列表获取失败")
)

type UserCommissionRepo interface {
	GetByChildUserId(context.Context, uint64, uint64, uint64, uint8, string, string) (*domain.UserCommission, error)
	List(context.Context, int, int, uint64, uint64, uint8, string, string, string) ([]*domain.UserCommission, error)
	ListByDay(context.Context, uint32, []string) ([]*domain.UserCommission, error)
	Count(context.Context, uint64, uint64, uint8, string, string, string) (int64, error)
	Statistics(context.Context, uint64, uint32, uint8) (*domain.UserCommission, error)
	StatisticsReal(context.Context, uint64, uint32, uint8) (*domain.UserCommission, error)
	Update(context.Context, *domain.UserCommission) (*domain.UserCommission, error)
	Save(context.Context, *domain.UserCommission) (*domain.UserCommission, error)
	DeleteByDay(context.Context, uint32, []string) error
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

func NewUserCommissionUsecase(repo UserCommissionRepo, urepo UserRepo, uorerepo UserOrganizationRelationRepo, uirrepo UserIntegralRelationRepo, uodrepo UserOpenDouyinRepo, uorrepo UserOrderRepo, ucrepo UserCouponRepo, usrrepo UserScanRecordRepo, ublrepo UserBalanceLogRepo, corepo CompanyOrganizationRepo, jorepo JinritemaiOrderRepo, dorepo DoukeOrderRepo, darepo DjAwemeRepo, tlrepo TaskLogRepo, tm Transaction, conf *conf.Data, cconf *conf.Company, oconf *conf.Organization, logger log.Logger) *UserCommissionUsecase {
	return &UserCommissionUsecase{repo: repo, urepo: urepo, uorerepo: uorerepo, uirrepo: uirrepo, uodrepo: uodrepo, uorrepo: uorrepo, ucrepo: ucrepo, usrrepo: usrrepo, ublrepo: ublrepo, corepo: corepo, jorepo: jorepo, dorepo: dorepo, darepo: darepo, tlrepo: tlrepo, tm: tm, conf: conf, cconf: cconf, oconf: oconf, log: log.NewHelper(logger)}
}

func (ucuc *UserCommissionUsecase) ListUserCommissions(ctx context.Context, pageNum, pageSize, userId, organizationId uint64, isDirect uint8, month, keyword string) (*domain.UserCommissionList, error) {
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

		startDay = startTime.Format("20060102")
		endDay = endTime.Format("20060102")
	}

	userCommissions, err := ucuc.repo.List(ctx, int(pageNum), int(pageSize), userId, organizationId, isDirect, startDay, endDay, keyword)

	if err != nil {
		return nil, WeixinUserCommissionListError
	}

	total, err := ucuc.repo.Count(ctx, userId, organizationId, isDirect, startDay, endDay, keyword)

	if err != nil {
		return nil, WeixinUserCommissionListError
	}

	list := make([]*domain.UserCommission, 0)
	activationTimeMap := make(map[uint64]string)

	for _, userCommission := range userCommissions {
		activationTime := ""

		if activationTimeVal, ok := activationTimeMap[userCommission.ChildUserId]; !ok {
			if userOrganizationRelation, err := ucuc.uorerepo.GetByUserId(ctx, userCommission.ChildUserId, organizationId, userCommission.UserId, "0"); err == nil {
				activationTime = tool.TimeToString("2006/01/02 15:04", userOrganizationRelation.CreateTime)

				activationTimeMap[userCommission.ChildUserId] = activationTime
			}
		} else {
			activationTime = activationTimeVal
		}

		var totalPayAmount, commissionPool, estimatedUserCommission, commissionRatio float32 = 0.00, 0.00, 0.00, 0.00

		if userCommission.UserCommissionType == 1 {
			totalPayAmount = userCommission.TotalPayAmount
			commissionPool = userCommission.CommissionPool
			estimatedUserCommission = userCommission.EstimatedUserCommission

			if userCommission.CommissionPool > 0 {
				commissionRatio = userCommission.EstimatedUserCommission / userCommission.CommissionPool
			}
		} else if userCommission.UserCommissionType == 2 || userCommission.UserCommissionType == 3 {
			if userCommissionInfo, err := ucuc.repo.GetByChildUserId(ctx, userId, userCommission.ChildUserId, organizationId, userCommission.UserCommissionType, startDay, endDay); err == nil {
				totalPayAmount = userCommissionInfo.TotalPayAmount
				commissionPool = userCommissionInfo.CommissionPool
				estimatedUserCommission = userCommissionInfo.EstimatedUserCommission

				if userCommission.CommissionPool > 0 {
					commissionRatio = userCommissionInfo.EstimatedUserCommission / userCommissionInfo.CommissionPool
				}
			}
		}

		list = append(list, &domain.UserCommission{
			ChildUserId:             userCommission.ChildUserId,
			ChildNickName:           userCommission.ChildNickName,
			ChildAvatarUrl:          userCommission.ChildAvatarUrl,
			ChildPhone:              userCommission.ChildPhone,
			RelationName:            userCommission.GetRelationName(ctx),
			TotalPayAmount:          totalPayAmount,
			CommissionPool:          commissionPool,
			EstimatedUserCommission: estimatedUserCommission,
			RealUserCommission:      userCommission.RealUserCommission,
			CommissionRatio:         commissionRatio * 100,
			UserCommissionType:      userCommission.UserCommissionType,
			UserCommissionTypeName:  userCommission.GetUserCommissionTypeName(ctx),
			ActivationTime:          activationTime,
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
		Key:   "基本等级",
		Value: levelName,
	})

	uiday, _ := strconv.ParseUint(time.Now().Format("20060102"), 10, 64)

	var wg sync.WaitGroup
	var presenterNum int64 = 0
	var childNum uint64 = 0
	var userIntegralRelations []*domain.UserIntegralRelation
	var statisticsCourse, statisticsOrder, statisticsCostOrder *domain.UserCommission
	var uerr, scerr, soerr, scoerr error

	wg.Add(5)

	go func() {
		defer wg.Done()

		presenterNum, _ = ucuc.uorerepo.Count(ctx, userId, organizationId, "")
	}()

	go func() {
		defer wg.Done()

		userIntegralRelations, uerr = ucuc.uirrepo.List(ctx, organizationId)
	}()

	go func() {
		defer wg.Done()

		statisticsCourse, scerr = ucuc.repo.Statistics(ctx, user.Id, uint32(uiday), 1)
	}()

	go func() {
		defer wg.Done()

		statisticsOrder, soerr = ucuc.repo.Statistics(ctx, user.Id, uint32(uiday), 2)
	}()

	go func() {
		defer wg.Done()

		statisticsCostOrder, scoerr = ucuc.repo.Statistics(ctx, user.Id, uint32(uiday), 3)
	}()

	wg.Wait()

	if uerr == nil {
		ucuc.uirrepo.GetChildNum(ctx, userId, &childNum, userIntegralRelations)
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
			Key:   "会员预估分佣",
			Value: fmt.Sprintf("%.2f", tool.Decimal(float64(statisticsCourse.EstimatedUserCommission), 2)),
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
			Key:   "会员预估分佣",
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
			Key:   "成本购预估分佣",
			Value: fmt.Sprintf("%.2f", tool.Decimal(float64(statisticsCostOrder.EstimatedUserCommission), 2)),
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
			Key:   "成本购预估分佣",
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
			Key:   "带货预估分佣",
			Value: fmt.Sprintf("%.2f", tool.Decimal(float64(statisticsOrder.EstimatedUserCommission), 2)),
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
			Key:   "带货预估分佣",
			Value: "0.00",
		})
	}

	return &domain.StatisticsUserOrganizationRelations{
		Statistics: statistics,
	}, nil
}

func (ucuc *UserCommissionUsecase) SyncOrderUserCommissions(ctx context.Context, day string) error {
	var wg sync.WaitGroup

	if len(day) == 0 {
		day = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}

	tday, _ := tool.StringToTime("2006-01-02", day)
	uiday, _ := strconv.ParseUint(tday.Format("20060102"), 10, 64)

	companyOrganizations, err := ucuc.corepo.List(ctx)

	if err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "SyncOrderUserCommissions", fmt.Sprintf("[SyncOrderUserCommissionsError], Description=%s", "获取企业机构列表失败"))
		inTaskLog.SetCreateTime(ctx)

		ucuc.tlrepo.Save(ctx, inTaskLog)

		return err
	}

	err = ucuc.tm.InTx(ctx, func(ctx context.Context) error {
		err = ucuc.repo.DeleteByDay(ctx, uint32(uiday), []string{"2"})

		if err != nil {
			return err
		}

		err = ucuc.ublrepo.DeleteByDay(ctx, 1, uint32(uiday), []string{"2"})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "SyncOrderUserCommissions", fmt.Sprintf("[SyncOrderUserCommissionsError], Description=%s", "根据时间用户电商分佣删除失败"))
		inTaskLog.SetCreateTime(ctx)

		ucuc.tlrepo.Save(ctx, inTaskLog)

		return err
	}

	for _, companyOrganization := range companyOrganizations.Data.List {
		if len(companyOrganization.OrganizationMcns) > 0 {
			if companyOrganization.OrganizationId == 5 {
				wg.Add(1)

				go ucuc.SyncOrderUserCommission(ctx, &wg, day, uiday, companyOrganization)
			}
		}
	}

	wg.Wait()

	return nil
}

func (ucuc *UserCommissionUsecase) SyncOrderUserCommission(ctx context.Context, wg *sync.WaitGroup, day string, uiday uint64, companyOrganization *v1.ListCompanyOrganizationsReply_CompanyOrganization) {
	defer wg.Done()

	userOrganizationRelations, err := ucuc.uorerepo.List(ctx, companyOrganization.OrganizationId)

	if err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "SyncOrderUserCommission", fmt.Sprintf("[SyncOrderUserCommissionError], OrganizationId=%d, Description=%s", companyOrganization.OrganizationId, "获取企业机构用户关系列表失败"))
		inTaskLog.SetCreateTime(ctx)

		ucuc.tlrepo.Save(ctx, inTaskLog)

		return
	}

	userOrganizationCommissions := make([]*domain.UserOrganizationCommission, 0)

	for _, userOrganizationRelation := range userOrganizationRelations {
		userOrganizationCommissions = append(userOrganizationCommissions, &domain.UserOrganizationCommission{
			UserId:              userOrganizationRelation.UserId,
			OrganizationUserId:  userOrganizationRelation.OrganizationUserId,
			OrganizationTutorId: userOrganizationRelation.OrganizationTutorId,
			Level:               userOrganizationRelation.Level,
			IsOrderRelation:     userOrganizationRelation.IsOrderRelation,
		})
		
		var orderNum, orderRefundNum uint64
		var totalPayAmount, estimatedCommission, realCommission float64

		organizationMcns := make([]string, 0)

		for _, organizationMcn := range companyOrganization.OrganizationMcns {
			organizationMcns = append(organizationMcns, organizationMcn.OrganizationMcn)
		}

		djAwemes, _ := ucuc.darepo.List(ctx, "90", strings.Join(organizationMcns, ","))

		if userOpenDouyins, err := ucuc.uodrepo.List(ctx, 0, 40, userOrganizationRelation.UserId, ""); err == nil {
			userCommissionOpenDouyins := make([]*domain.UserCommissionOpenDouyin, 0)

			for _, userOpenDouyin := range userOpenDouyins {
				for _, djAweme := range djAwemes {
					if djAweme.AwemeId == userOpenDouyin.AccountId {
						userCommissionOpenDouyins = append(userCommissionOpenDouyins, &domain.UserCommissionOpenDouyin{
							ClientKey: userOpenDouyin.ClientKey,
							OpenId:    userOpenDouyin.OpenId,
						})

						if len(userCommissionOpenDouyins) > 0 {
							if content, err := json.Marshal(userCommissionOpenDouyins); err == nil {
								if userOpenDouyin.OpenId == "_000Iem5lJ6nbYAD1r4v-Xi7OuSw-iBEnOoz" {
									fmt.Println("#########################")
									fmt.Println(djAweme.BindStartTime)
									fmt.Println("#########################")
								}

								if len(djAweme.BindStartTime) > 0 {
									if statistics, err := ucuc.jorepo.StatisticsByPayTimeDay(ctx, djAweme.BindStartTime, day, string(content), ""); err == nil {
										if userOpenDouyin.OpenId == "_000Iem5lJ6nbYAD1r4v-Xi7OuSw-iBEnOoz" {
											fmt.Println("#########################")
											fmt.Println(statistics)
											fmt.Println("#########################")
										}

										for _, statistic := range statistics.Data.Statistics {
											if statistic.Key == "orderRefundNum" {
												orderRefundNum, _ = strconv.ParseUint(statistic.Value, 10, 64)
											} else if statistic.Key == "orderNum" {
												orderNum, _ = strconv.ParseUint(statistic.Value, 10, 64)
											} else if statistic.Key == "totalPayAmount" {
												totalPayAmount, _ = strconv.ParseFloat(statistic.Value, 10)
											} else if statistic.Key == "estimatedCommission" {
												estimatedCommission, _ = strconv.ParseFloat(statistic.Value, 10)
											} else if statistic.Key == "realCommission" {
												realCommission, _ = strconv.ParseFloat(statistic.Value, 10)
											}
										}

										isNotExist := true

										for _, userOrganizationCommission := range userOrganizationCommissions {
											if userOrganizationCommission.UserId == userOrganizationRelation.UserId {
												userOrganizationCommission.OrderNum += orderNum
												userOrganizationCommission.OrderRefundNum += orderRefundNum
												userOrganizationCommission.TotalPayAmount += totalPayAmount
												userOrganizationCommission.EstimatedCommission += estimatedCommission
												userOrganizationCommission.RealCommission += realCommission

												isNotExist = false

												break
											}
										}

										if isNotExist {
											userOrganizationCommissions = append(userOrganizationCommissions, &domain.UserOrganizationCommission{
												UserId:              userOrganizationRelation.UserId,
												OrganizationUserId:  userOrganizationRelation.OrganizationUserId,
												OrganizationTutorId: userOrganizationRelation.OrganizationTutorId,
												Level:               userOrganizationRelation.Level,
												IsOrderRelation:     userOrganizationRelation.IsOrderRelation,
												OrderNum:            orderNum,
												OrderRefundNum:      orderRefundNum,
												TotalPayAmount:      totalPayAmount,
												EstimatedCommission: estimatedCommission,
												RealCommission:      realCommission,
											})
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	for _, userOrganizationCommission := range userOrganizationCommissions {
		if userOrganizationCommission.UserId == 26044 {
			fmt.Println("#########################3")
			fmt.Println(userOrganizationCommission.RealCommission)
			fmt.Println("#########################3")
		}
	}

	ucuc.getOrderComission(ctx, day, uiday, userOrganizationCommissions, companyOrganization)
}

func (ucuc *UserCommissionUsecase) getOrderComission(ctx context.Context, day string, uiday uint64, userOrganizationCommissions []*domain.UserOrganizationCommission, companyOrganization *v1.ListCompanyOrganizationsReply_CompanyOrganization) {
	for _, userOrganizationCommission := range userOrganizationCommissions {
		if userOrganizationCommission.UserId == 26044 {
			fmt.Println("#########################4")
			fmt.Println(userOrganizationCommission)
			fmt.Println("#########################4")
		}

		if userOrganizationCommission.RealCommission > 0 && userOrganizationCommission.OrganizationUserId > 0 {
			baseServiceRealCommission := userOrganizationCommission.RealCommission / 0.9 * companyOrganization.OrganizationColonelCommission.OrderRatio / 100
			baseServiceEstimatedCommission := userOrganizationCommission.EstimatedCommission / 0.9 * companyOrganization.OrganizationColonelCommission.OrderRatio / 100

			for _, parentUserOrganizationCommission := range userOrganizationCommissions {
				var estimatedUserCommission, realUserCommission float64 = 0.00, 0.00

				if parentUserOrganizationCommission.UserId == userOrganizationCommission.OrganizationUserId {
					if userOrganizationCommission.UserId == 26044 {
						fmt.Println("#########################3")
						fmt.Println(parentUserOrganizationCommission)
						fmt.Println("#########################3")
					}

					if parentUserOrganizationCommission.Level == 3 {
						estimatedUserCommission = baseServiceEstimatedCommission * companyOrganization.OrganizationColonelCommission.IntermediateAdvancedPresenterOrderCommissionRule / 100
						realUserCommission = baseServiceRealCommission * companyOrganization.OrganizationColonelCommission.IntermediateAdvancedPresenterOrderCommissionRule / 100

						inUserCommission := domain.NewUserCommission(ctx, parentUserOrganizationCommission.UserId, companyOrganization.OrganizationId, userOrganizationCommission.UserId, userOrganizationCommission.OrderNum, userOrganizationCommission.OrderRefundNum, uint32(uiday), userOrganizationCommission.Level, parentUserOrganizationCommission.Level, 1, 2, float32(userOrganizationCommission.TotalPayAmount), float32(baseServiceEstimatedCommission), float32(userOrganizationCommission.EstimatedCommission), float32(userOrganizationCommission.RealCommission), float32(estimatedUserCommission), float32(realUserCommission))
						inUserCommission.SetCreateTime(ctx)
						inUserCommission.SetUpdateTime(ctx)

						if userCommission, err := ucuc.repo.Save(ctx, inUserCommission); err == nil {
							inUserBalanceLog := domain.NewUserBalanceLog(ctx, parentUserOrganizationCommission.UserId, userCommission.Id, 2, 1, 1, float32(realUserCommission), "电商佣金")
							inUserBalanceLog.SetDay(ctx, uint32(uiday))
							inUserBalanceLog.SetCreateTime(ctx)
							inUserBalanceLog.SetUpdateTime(ctx)

							ucuc.ublrepo.Save(ctx, inUserBalanceLog)
						}

						for _, tutorUserOrganizationCommission := range userOrganizationCommissions {
							if tutorUserOrganizationCommission.UserId == userOrganizationCommission.OrganizationTutorId {
								if tutorUserOrganizationCommission.Level == 4 {
									estimatedUserCommission = baseServiceEstimatedCommission * companyOrganization.OrganizationColonelCommission.IntermediateAdvancedTutorOrderCommissionRule / 100
									realUserCommission = baseServiceRealCommission * companyOrganization.OrganizationColonelCommission.IntermediateAdvancedTutorOrderCommissionRule / 100

									inUserCommission = domain.NewUserCommission(ctx, tutorUserOrganizationCommission.UserId, companyOrganization.OrganizationId, userOrganizationCommission.UserId, userOrganizationCommission.OrderNum, userOrganizationCommission.OrderRefundNum, uint32(uiday), userOrganizationCommission.Level, tutorUserOrganizationCommission.Level, 2, 2, float32(userOrganizationCommission.TotalPayAmount), float32(baseServiceEstimatedCommission), float32(userOrganizationCommission.EstimatedCommission), float32(userOrganizationCommission.RealCommission), float32(estimatedUserCommission), float32(realUserCommission))
									inUserCommission.SetCreateTime(ctx)
									inUserCommission.SetUpdateTime(ctx)

									if userCommission, err := ucuc.repo.Save(ctx, inUserCommission); err == nil {
										inUserBalanceLog := domain.NewUserBalanceLog(ctx, tutorUserOrganizationCommission.UserId, userCommission.Id, 2, 1, 1, float32(realUserCommission), "电商佣金")
										inUserBalanceLog.SetDay(ctx, uint32(uiday))
										inUserBalanceLog.SetCreateTime(ctx)
										inUserBalanceLog.SetUpdateTime(ctx)

										ucuc.ublrepo.Save(ctx, inUserBalanceLog)
									}
								}

								break
							}
						}
					} else if parentUserOrganizationCommission.Level == 4 {
						estimatedUserCommission = baseServiceEstimatedCommission * companyOrganization.OrganizationColonelCommission.AdvancedPresenterOrderCommissionRule / 100
						realUserCommission = baseServiceRealCommission * companyOrganization.OrganizationColonelCommission.AdvancedPresenterOrderCommissionRule / 100

						inUserCommission := domain.NewUserCommission(ctx, parentUserOrganizationCommission.UserId, companyOrganization.OrganizationId, userOrganizationCommission.UserId, userOrganizationCommission.OrderNum, userOrganizationCommission.OrderRefundNum, uint32(uiday), userOrganizationCommission.Level, parentUserOrganizationCommission.Level, 1, 2, float32(userOrganizationCommission.TotalPayAmount), float32(baseServiceEstimatedCommission), float32(userOrganizationCommission.EstimatedCommission), float32(userOrganizationCommission.RealCommission), float32(estimatedUserCommission), float32(realUserCommission))
						inUserCommission.SetCreateTime(ctx)
						inUserCommission.SetUpdateTime(ctx)

						if userCommission, err := ucuc.repo.Save(ctx, inUserCommission); err == nil {
							inUserBalanceLog := domain.NewUserBalanceLog(ctx, parentUserOrganizationCommission.UserId, userCommission.Id, 2, 1, 1, float32(realUserCommission), "电商佣金")
							inUserBalanceLog.SetDay(ctx, uint32(uiday))
							inUserBalanceLog.SetCreateTime(ctx)
							inUserBalanceLog.SetUpdateTime(ctx)

							ucuc.ublrepo.Save(ctx, inUserBalanceLog)
						}
					}

					break
				}
			}
		}
	}
}

func (ucuc *UserCommissionUsecase) SyncCostOrderUserCommissions(ctx context.Context, day string) error {
	var wg sync.WaitGroup

	if len(day) == 0 {
		day = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}

	tday, _ := tool.StringToTime("2006-01-02", day)
	uiday, _ := strconv.ParseUint(tday.Format("20060102"), 10, 64)

	companyOrganizations, err := ucuc.corepo.List(ctx)

	if err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "SyncCostOrderUserCommissions", fmt.Sprintf("[SyncCostOrderUserCommissionsError], Description=%s", "获取企业机构列表失败"))
		inTaskLog.SetCreateTime(ctx)

		ucuc.tlrepo.Save(ctx, inTaskLog)

		return err
	}

	err = ucuc.tm.InTx(ctx, func(ctx context.Context) error {
		err = ucuc.repo.DeleteByDay(ctx, uint32(uiday), []string{"3"})

		if err != nil {
			return err
		}

		err = ucuc.ublrepo.DeleteByDay(ctx, 1, uint32(uiday), []string{"3"})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "SyncCostOrderUserCommissions", fmt.Sprintf("[SyncCostOrderUserCommissionsError], Description=%s", "根据时间用户成本购分佣删除失败"))
		inTaskLog.SetCreateTime(ctx)

		ucuc.tlrepo.Save(ctx, inTaskLog)

		return err
	}

	for _, companyOrganization := range companyOrganizations.Data.List {
		wg.Add(1)

		go ucuc.SyncCostOrderUserCommission(ctx, &wg, day, uiday, companyOrganization)
	}

	wg.Wait()

	return nil
}

func (ucuc *UserCommissionUsecase) SyncCostOrderUserCommission(ctx context.Context, wg *sync.WaitGroup, day string, uiday uint64, companyOrganization *v1.ListCompanyOrganizationsReply_CompanyOrganization) {
	defer wg.Done()

	userOrganizationRelations, err := ucuc.uorerepo.List(ctx, companyOrganization.OrganizationId)

	if err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "SyncCostOrderUserCommission", fmt.Sprintf("[SyncCostOrderUserCommissionError], OrganizationId=%d, Description=%s", companyOrganization.OrganizationId, "获取企业机构用户关系列表失败"))
		inTaskLog.SetCreateTime(ctx)

		ucuc.tlrepo.Save(ctx, inTaskLog)

		return
	}

	userOrganizationCommissions := make([]*domain.UserOrganizationCommission, 0)

	for _, userOrganizationRelation := range userOrganizationRelations {
		var costOrderTotalPayAmount, costOrderEstimatedCommission, costOrderRealCommission float64

		if companyOrganization.OrganizationId == ucuc.oconf.DefaultOrganizationId {
			if statistics, err := ucuc.dorepo.StatisticsByDay(ctx, userOrganizationRelation.UserId, day); err == nil {
				for _, statistic := range statistics.Data.Statistics {
					if statistic.Key == "totalPayAmount" {
						costOrderTotalPayAmount, _ = strconv.ParseFloat(statistic.Value, 10)
					} else if statistic.Key == "estimatedCommission" {
						costOrderEstimatedCommission, _ = strconv.ParseFloat(statistic.Value, 10)
					} else if statistic.Key == "realCommission" {
						costOrderRealCommission, _ = strconv.ParseFloat(statistic.Value, 10)
					}
				}
			}
		} else if companyOrganization.OrganizationId == ucuc.oconf.DjOrganizationId {
			if statistics, err := ucuc.dorepo.StatisticsByDay(ctx, userOrganizationRelation.UserId, day); err == nil {
				for _, statistic := range statistics.Data.Statistics {
					if statistic.Key == "totalPayAmount" {
						costOrderTotalPayAmount, _ = strconv.ParseFloat(statistic.Value, 10)
					} else if statistic.Key == "estimatedCommission" {
						costOrderEstimatedCommission, _ = strconv.ParseFloat(statistic.Value, 10)
					} else if statistic.Key == "realCommission" {
						costOrderRealCommission, _ = strconv.ParseFloat(statistic.Value, 10)
					}
				}
			}
		}

		userOrganizationCommissions = append(userOrganizationCommissions, &domain.UserOrganizationCommission{
			UserId:                       userOrganizationRelation.UserId,
			OrganizationUserId:           userOrganizationRelation.OrganizationUserId,
			OrganizationTutorId:          userOrganizationRelation.OrganizationTutorId,
			Level:                        userOrganizationRelation.Level,
			IsOrderRelation:              userOrganizationRelation.IsOrderRelation,
			CostOrderTotalPayAmount:      costOrderTotalPayAmount,
			CostOrderEstimatedCommission: costOrderEstimatedCommission,
			CostOrderRealCommission:      costOrderRealCommission,
		})
	}

	ucuc.getCostOrderComission(ctx, day, uiday, userOrganizationCommissions, companyOrganization)
}

func (ucuc *UserCommissionUsecase) getCostOrderComission(ctx context.Context, day string, uiday uint64, userOrganizationCommissions []*domain.UserOrganizationCommission, companyOrganization *v1.ListCompanyOrganizationsReply_CompanyOrganization) {
	for _, userOrganizationCommission := range userOrganizationCommissions {
		if userOrganizationCommission.CostOrderRealCommission > 0 {
			baseServiceRealCommission := userOrganizationCommission.CostOrderRealCommission * companyOrganization.OrganizationColonelCommission.CostOrderRatio / 100
			baseServiceEstimatedCommission := userOrganizationCommission.CostOrderEstimatedCommission * companyOrganization.OrganizationColonelCommission.CostOrderRatio / 100

			inUserBalanceLog := domain.NewUserBalanceLog(ctx, userOrganizationCommission.UserId, 0, 3, 1, 1, float32(userOrganizationCommission.CostOrderRealCommission*(1-companyOrganization.OrganizationColonelCommission.CostOrderRatio/100)), "成本购佣金")
			inUserBalanceLog.SetDay(ctx, uint32(uiday))
			inUserBalanceLog.SetCreateTime(ctx)
			inUserBalanceLog.SetUpdateTime(ctx)

			ucuc.ublrepo.Save(ctx, inUserBalanceLog)

			for _, parentUserOrganizationCommission := range userOrganizationCommissions {
				var estimatedUserCommission, realUserCommission float64 = 0.00, 0.00

				if parentUserOrganizationCommission.UserId == userOrganizationCommission.OrganizationUserId {
					if parentUserOrganizationCommission.Level == 2 {
						estimatedUserCommission = baseServiceEstimatedCommission * companyOrganization.OrganizationColonelCommission.PrimaryAdvancedPresenterCostOrderCommissionRule / 100
						realUserCommission = baseServiceRealCommission * companyOrganization.OrganizationColonelCommission.PrimaryAdvancedPresenterCostOrderCommissionRule / 100

						inUserCommission := domain.NewUserCommission(ctx, parentUserOrganizationCommission.UserId, companyOrganization.OrganizationId, userOrganizationCommission.UserId, 0, 0, uint32(uiday), userOrganizationCommission.Level, parentUserOrganizationCommission.Level, 1, 3, float32(userOrganizationCommission.CostOrderTotalPayAmount), float32(baseServiceEstimatedCommission), float32(userOrganizationCommission.CostOrderEstimatedCommission), float32(userOrganizationCommission.CostOrderRealCommission), float32(estimatedUserCommission), float32(realUserCommission))
						inUserCommission.SetCreateTime(ctx)
						inUserCommission.SetUpdateTime(ctx)

						if userCommission, err := ucuc.repo.Save(ctx, inUserCommission); err == nil {
							inUserBalanceLog = domain.NewUserBalanceLog(ctx, parentUserOrganizationCommission.UserId, userCommission.Id, 3, 1, 1, float32(realUserCommission), "成本购佣金")
							inUserBalanceLog.SetDay(ctx, uint32(uiday))
							inUserBalanceLog.SetCreateTime(ctx)
							inUserBalanceLog.SetUpdateTime(ctx)

							ucuc.ublrepo.Save(ctx, inUserBalanceLog)
						}

						for _, tutorUserOrganizationCommission := range userOrganizationCommissions {
							if tutorUserOrganizationCommission.UserId == userOrganizationCommission.OrganizationTutorId {
								if tutorUserOrganizationCommission.Level == 4 {
									estimatedUserCommission = baseServiceEstimatedCommission * companyOrganization.OrganizationColonelCommission.PrimaryAdvancedTutorCostOrderCommissionRule / 100
									realUserCommission = baseServiceRealCommission * companyOrganization.OrganizationColonelCommission.PrimaryAdvancedTutorCostOrderCommissionRule / 100

									inUserCommission = domain.NewUserCommission(ctx, tutorUserOrganizationCommission.UserId, companyOrganization.OrganizationId, userOrganizationCommission.UserId, 0, 0, uint32(uiday), userOrganizationCommission.Level, tutorUserOrganizationCommission.Level, 2, 3, float32(userOrganizationCommission.CostOrderTotalPayAmount), float32(baseServiceEstimatedCommission), float32(userOrganizationCommission.CostOrderEstimatedCommission), float32(userOrganizationCommission.CostOrderRealCommission), float32(estimatedUserCommission), float32(realUserCommission))
									inUserCommission.SetCreateTime(ctx)
									inUserCommission.SetUpdateTime(ctx)

									if userCommission, err := ucuc.repo.Save(ctx, inUserCommission); err == nil {
										inUserBalanceLog = domain.NewUserBalanceLog(ctx, tutorUserOrganizationCommission.UserId, userCommission.Id, 3, 1, 1, float32(realUserCommission), "成本购佣金")
										inUserBalanceLog.SetDay(ctx, uint32(uiday))
										inUserBalanceLog.SetCreateTime(ctx)
										inUserBalanceLog.SetUpdateTime(ctx)

										ucuc.ublrepo.Save(ctx, inUserBalanceLog)
									}
								}

								break
							}
						}
					} else if parentUserOrganizationCommission.Level == 3 {
						estimatedUserCommission = baseServiceEstimatedCommission * companyOrganization.OrganizationColonelCommission.IntermediateAdvancedPresenterCostOrderCommissionRule / 100
						realUserCommission = baseServiceRealCommission * companyOrganization.OrganizationColonelCommission.IntermediateAdvancedPresenterCostOrderCommissionRule / 100

						inUserCommission := domain.NewUserCommission(ctx, parentUserOrganizationCommission.UserId, companyOrganization.OrganizationId, userOrganizationCommission.UserId, 0, 0, uint32(uiday), userOrganizationCommission.Level, parentUserOrganizationCommission.Level, 1, 3, float32(userOrganizationCommission.CostOrderTotalPayAmount), float32(baseServiceEstimatedCommission), float32(userOrganizationCommission.CostOrderEstimatedCommission), float32(userOrganizationCommission.CostOrderRealCommission), float32(estimatedUserCommission), float32(realUserCommission))
						inUserCommission.SetCreateTime(ctx)
						inUserCommission.SetUpdateTime(ctx)

						if userCommission, err := ucuc.repo.Save(ctx, inUserCommission); err == nil {
							inUserBalanceLog = domain.NewUserBalanceLog(ctx, parentUserOrganizationCommission.UserId, userCommission.Id, 3, 1, 1, float32(realUserCommission), "成本购佣金")
							inUserBalanceLog.SetDay(ctx, uint32(uiday))
							inUserBalanceLog.SetCreateTime(ctx)
							inUserBalanceLog.SetUpdateTime(ctx)

							ucuc.ublrepo.Save(ctx, inUserBalanceLog)
						}

						for _, tutorUserOrganizationCommission := range userOrganizationCommissions {
							if tutorUserOrganizationCommission.UserId == userOrganizationCommission.OrganizationTutorId {
								if tutorUserOrganizationCommission.Level == 4 {
									estimatedUserCommission = baseServiceEstimatedCommission * companyOrganization.OrganizationColonelCommission.IntermediateAdvancedTutorCostOrderCommissionRule / 100
									realUserCommission = baseServiceRealCommission * companyOrganization.OrganizationColonelCommission.IntermediateAdvancedTutorCostOrderCommissionRule / 100

									inUserCommission = domain.NewUserCommission(ctx, tutorUserOrganizationCommission.UserId, companyOrganization.OrganizationId, userOrganizationCommission.UserId, 0, 0, uint32(uiday), userOrganizationCommission.Level, tutorUserOrganizationCommission.Level, 2, 3, float32(userOrganizationCommission.CostOrderTotalPayAmount), float32(baseServiceEstimatedCommission), float32(userOrganizationCommission.CostOrderEstimatedCommission), float32(userOrganizationCommission.CostOrderRealCommission), float32(estimatedUserCommission), float32(realUserCommission))
									inUserCommission.SetCreateTime(ctx)
									inUserCommission.SetUpdateTime(ctx)

									if userCommission, err := ucuc.repo.Save(ctx, inUserCommission); err == nil {
										inUserBalanceLog = domain.NewUserBalanceLog(ctx, tutorUserOrganizationCommission.UserId, userCommission.Id, 3, 1, 1, float32(realUserCommission), "成本购佣金")
										inUserBalanceLog.SetDay(ctx, uint32(uiday))
										inUserBalanceLog.SetCreateTime(ctx)
										inUserBalanceLog.SetUpdateTime(ctx)

										ucuc.ublrepo.Save(ctx, inUserBalanceLog)
									}
								}

								break
							}
						}
					} else if parentUserOrganizationCommission.Level == 4 {
						estimatedUserCommission = baseServiceEstimatedCommission * companyOrganization.OrganizationColonelCommission.AdvancedPresenterCostOrderCommissionRule / 100
						realUserCommission = baseServiceRealCommission * companyOrganization.OrganizationColonelCommission.AdvancedPresenterCostOrderCommissionRule / 100

						inUserCommission := domain.NewUserCommission(ctx, parentUserOrganizationCommission.UserId, companyOrganization.OrganizationId, userOrganizationCommission.UserId, 0, 0, uint32(uiday), userOrganizationCommission.Level, parentUserOrganizationCommission.Level, 1, 3, float32(userOrganizationCommission.CostOrderTotalPayAmount), float32(baseServiceEstimatedCommission), float32(userOrganizationCommission.CostOrderEstimatedCommission), float32(userOrganizationCommission.CostOrderRealCommission), float32(estimatedUserCommission), float32(realUserCommission))
						inUserCommission.SetCreateTime(ctx)
						inUserCommission.SetUpdateTime(ctx)

						if userCommission, err := ucuc.repo.Save(ctx, inUserCommission); err == nil {
							inUserBalanceLog = domain.NewUserBalanceLog(ctx, parentUserOrganizationCommission.UserId, userCommission.Id, 3, 1, 1, float32(realUserCommission), "成本购佣金")
							inUserBalanceLog.SetDay(ctx, uint32(uiday))
							inUserBalanceLog.SetCreateTime(ctx)
							inUserBalanceLog.SetUpdateTime(ctx)

							ucuc.ublrepo.Save(ctx, inUserBalanceLog)
						}
					}

					break
				}
			}
		}
	}
}

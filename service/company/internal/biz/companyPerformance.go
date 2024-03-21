package biz

import (
	"company/internal/conf"
	"company/internal/domain"
	"company/internal/pkg/tool"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	CompanyCompanyPerformanceNotFound = errors.NotFound("COMPANY_COMAPNY_PERFORMANCE_NOT_FOUND", "团队绩效不存在")
)

type QianchuanAdAdvertiserRepo interface {
	List(context.Context, uint64, string) ([]*domain.QianchuanAdAdvertiser, error)
}

type CompanyPerformanceUsecase struct {
	cpmrepo  CompanyPerformanceMonthlyRepo
	cpdrepo  CompanyPerformanceDailyRepo
	cprrepo  CompanyPerformanceRebalanceRepo
	crepo    CompanyRepo
	curepo   CompanyUserRepo
	currepo  CompanyUserRoleRepo
	cprurepo CompanyPerformanceRuleRepo
	qadrepo  QianchuanAdvertiserRepo
	qarepo   QianchuanAdAdvertiserRepo
	conf     *conf.Data
	log      *log.Helper
}

func NewCompanyPerformanceUsecase(cpmrepo CompanyPerformanceMonthlyRepo, cpdrepo CompanyPerformanceDailyRepo, cprrepo CompanyPerformanceRebalanceRepo, crepo CompanyRepo, curepo CompanyUserRepo, currepo CompanyUserRoleRepo, cprurepo CompanyPerformanceRuleRepo, qadrepo QianchuanAdvertiserRepo, qarepo QianchuanAdAdvertiserRepo, conf *conf.Data, logger log.Logger) *CompanyPerformanceUsecase {
	return &CompanyPerformanceUsecase{cpmrepo: cpmrepo, cpdrepo: cpdrepo, cprrepo: cprrepo, crepo: crepo, curepo: curepo, currepo: currepo, cprurepo: cprurepo, qadrepo: qadrepo, qarepo: qarepo, conf: conf, log: log.NewHelper(logger)}
}

func (cpuc *CompanyPerformanceUsecase) ListCompanyPerformances(ctx context.Context, pageNum, pageSize, companyId uint64, updateDay uint32) (*domain.CompanyPerformanceMonthlyList, error) {
	if _, err := cpuc.crepo.GetById(ctx, companyId); err != nil {
		return nil, CompanyCompanyNotFound
	}

	sumStatCost, err := cpuc.cpmrepo.Sum(ctx, companyId, updateDay)

	if err != nil {
		return nil, CompanyDataError
	}

	companyPerformances, err := cpuc.cpmrepo.List(ctx, int(pageNum), int(pageSize), companyId, updateDay)

	if err != nil {
		return nil, CompanyCompanyPerformanceNotFound
	}

	list := make([]*domain.CompanyPerformanceMonthly, 0)

	for _, companyPerformance := range companyPerformances {
		if sumStatCost == 0.00 {
			companyPerformance.StatCostProportion = 0.00
		} else {
			companyPerformance.StatCostProportion = companyPerformance.StatCost / sumStatCost
		}

		companyPerformance.GetAdvertisers(ctx)
		companyPerformance.GetQianchuanAdvertisers(ctx)

		list = append(list, companyPerformance)
	}

	total, err := cpuc.cpmrepo.Count(ctx, companyId, updateDay)

	if err != nil {
		return nil, CompanyDataError
	}

	return &domain.CompanyPerformanceMonthlyList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (cpuc *CompanyPerformanceUsecase) GetCompanyPerformances(ctx context.Context, userId, companyId uint64, updateDay string) (*domain.CompanyPerformanceDailyList, error) {
	if _, err := cpuc.crepo.GetById(ctx, companyId); err != nil {
		return nil, CompanyCompanyNotFound
	}

	companyPerformances, err := cpuc.cpdrepo.List(ctx, userId, companyId, updateDay)

	if err != nil {
		return nil, CompanyCompanyPerformanceNotFound
	}

	companyPerformanceRebalances, err := cpuc.cprrepo.List(ctx, userId, companyId, updateDay)

	list := make([]*domain.CompanyPerformanceDaily, 0)

	updateTime, _ := tool.StringToTime("2006-01", updateDay)

	var endDay int

	if updateDay == fmt.Sprintf("%d", time.Now().Year())+"-"+time.Now().Format("01") {
		endDay = time.Now().Day()
	} else {
		month, _ := strconv.Atoi(updateTime.Format("01"))

		endDay = tool.GetDays(updateTime.Year(), month)
	}

	for day := 1; day <= endDay; day++ {
		sUpdateDay := fmt.Sprintf("%d", updateTime.Year()) + "-" + updateTime.Format("01") + "-" + fmt.Sprintf("%02d", day)

		isNotExist := true
		companyPerformanceDaily := &domain.CompanyPerformanceDaily{}
		rebalanceCosts := make([]*domain.CompanyPerformanceRebalance, 0)

		for _, companyPerformance := range companyPerformances {
			if sUpdateDay == tool.TimeToString("2006-01-02", companyPerformance.UpdateDay) {
				isNotExist = false

				companyPerformanceDaily.Id = companyPerformance.Id
				companyPerformanceDaily.CompanyId = companyPerformance.CompanyId
				companyPerformanceDaily.UserId = companyPerformance.UserId
				companyPerformanceDaily.Advertisers = companyPerformance.Advertisers
				companyPerformanceDaily.StatCost = companyPerformance.StatCost
				companyPerformanceDaily.PayOrderAmount = companyPerformance.PayOrderAmount
				companyPerformanceDaily.Roi = companyPerformance.Roi
				companyPerformanceDaily.Cost = companyPerformance.Cost
				companyPerformanceDaily.RebalanceCosts = companyPerformance.RebalanceCosts
				companyPerformanceDaily.UpdateDay = companyPerformance.UpdateDay
				companyPerformanceDaily.CreateTime = companyPerformance.CreateTime
				companyPerformanceDaily.UpdateTime = companyPerformance.UpdateTime

				break
			}
		}

		if isNotExist {
			tUpdateDay, _ := tool.StringToTime("2006-01-02", sUpdateDay)

			companyPerformanceDaily.UpdateDay = tUpdateDay
		}

		for _, companyPerformanceRebalance := range companyPerformanceRebalances {
			if tool.TimeToString("2006-01-02", companyPerformanceRebalance.UpdateDay) == sUpdateDay {
				rebalanceCosts = append(rebalanceCosts, &domain.CompanyPerformanceRebalance{
					Id:            companyPerformanceRebalance.Id,
					CompanyId:     companyPerformanceRebalance.CompanyId,
					UserId:        companyPerformanceRebalance.UserId,
					Cost:          companyPerformanceRebalance.Cost,
					RebalanceType: companyPerformanceRebalance.RebalanceType,
					Reason:        companyPerformanceRebalance.Reason,
					UpdateDay:     companyPerformanceRebalance.UpdateDay,
					CreateTime:    companyPerformanceRebalance.CreateTime,
					UpdateTime:    companyPerformanceRebalance.UpdateTime,
				})
			}
		}

		companyPerformanceDaily.RebalanceCosts = rebalanceCosts

		list = append(list, companyPerformanceDaily)

	}

	return &domain.CompanyPerformanceDailyList{
		List: list,
	}, nil
}

func (cpuc *CompanyPerformanceUsecase) syncCompanyPerformance(ctx context.Context, companyId uint64, updateDay string, wg *sync.WaitGroup) {
	defer wg.Done()

	qianchuanAdvertisers, err := cpuc.qadrepo.List(ctx, companyId, 0, 0, "", "", "1")

	if err != nil {
		return
	}

	adAdvertisers := make(map[string]map[string]map[string]float32)
	performanceRules := make(map[string]map[string][]*domain.PRules)
	performances := make(map[string]map[string]float32)

	if qianchuanAdAdvertisers, err := cpuc.qarepo.List(ctx, companyId, updateDay); err == nil {
		for _, qianchuanAdAdvertiser := range qianchuanAdAdvertisers {
			key := strconv.FormatUint(qianchuanAdAdvertiser.AdvertiserId, 10)

			if _, ok := adAdvertisers[key]; !ok {
				adAdvertisers[key] = make(map[string]map[string]float32)
			}

			if _, ok := adAdvertisers[key][qianchuanAdAdvertiser.MarketingGoal]; !ok {
				adAdvertisers[key][qianchuanAdAdvertiser.MarketingGoal] = make(map[string]float32)
			}

			adAdvertisers[key][qianchuanAdAdvertiser.MarketingGoal]["statCost"] = qianchuanAdAdvertiser.StatCost
			adAdvertisers[key][qianchuanAdAdvertiser.MarketingGoal]["payOrderAmount"] = qianchuanAdAdvertiser.PayOrderAmount

			if qianchuanAdAdvertiser.StatCost == 0.00 {
				adAdvertisers[key][qianchuanAdAdvertiser.MarketingGoal]["roi"] = 0.00
			} else {
				adAdvertisers[key][qianchuanAdAdvertiser.MarketingGoal]["roi"] = qianchuanAdAdvertiser.PayOrderAmount / qianchuanAdAdvertiser.StatCost
			}
		}

		if companyPerformanceRules, err := cpuc.cprurepo.ListAdvertiserIdsAndPerformanceRuleNotNull(ctx, companyId); err == nil {
			for _, companyPerformanceRule := range companyPerformanceRules {
				advertiserIds := tool.RemoveEmptyString(strings.Split(companyPerformanceRule.AdvertiserIds, ","))

				for _, advertiserId := range advertiserIds {
					if _, ok := performanceRules[advertiserId]; !ok {
						performanceRules[advertiserId] = make(map[string][]*domain.PRules)
					}

					companyPerformanceRule.GetPerformanceRules(ctx)

					for _, performanceRule := range companyPerformanceRule.PerformanceRulest {
						performanceRules[advertiserId][performanceRule.Ptype] = performanceRule.Rules
					}
				}
			}
		}

		for index, performanceRule := range performanceRules {
			if _, ok := performances[index]; !ok {
				performances[index] = make(map[string]float32)
			}

			if _, ok := adAdvertisers[index]; ok {
				for iadAdvertiser, adAdvertiser := range adAdvertisers[index] {
					rules := make([]*domain.PRules, 0)

					if iadAdvertiser == "LIVE_PROM_GOODS" {
						rules = performanceRule["live"]
					} else if iadAdvertiser == "VIDEO_PROM_GOODS" {
						rules = performanceRule["video"]
					}

					for _, rule := range rules {
						cMin, _ := strconv.ParseFloat(rule.Condition.Min, 10)
						cMax, _ := strconv.ParseFloat(rule.Condition.Max, 10)
						cPercentage, _ := strconv.ParseFloat(rule.Commission.Percentage, 10)

						if float32(cMin) <= adAdvertiser["roi"] && float32(cMax) > adAdvertiser["roi"] {
							performances[index]["cost"] += adAdvertiser["statCost"] * float32(cPercentage) / 100
						}
					}
				}
			}
		}
	}

	if companyUsers, err := cpuc.curepo.ListByCompanyId(ctx, companyId); err == nil {
		for _, companyUser := range companyUsers {
			var statCost float32 = 0.00
			var payOrderAmount float32 = 0.00
			var roi float32 = 0.00
			var cost float32 = 0.00
			var advertisers []*domain.Advertiser = make([]*domain.Advertiser, 0)

			if companyUser.Status == 1 && companyUser.Role != 1 {
				if companyUserRoles, err := cpuc.currepo.ListByUserIdAndCompanyId(ctx, companyUser.Id, companyUser.CompanyId); err == nil {
					if len(companyUserRoles) == 0 && companyUser.Role == 2 {
						for _, qianchuanAdvertiser := range qianchuanAdvertisers.Data.List {
							advertisers = append(advertisers, &domain.Advertiser{
								AdvertiserId:   qianchuanAdvertiser.AdvertiserId,
								AdvertiserName: qianchuanAdvertiser.AdvertiserName,
							})

							key := strconv.FormatUint(qianchuanAdvertiser.AdvertiserId, 10)

							for _, adAdvertiser := range adAdvertisers[key] {
								statCost += adAdvertiser["statCost"]
								payOrderAmount += adAdvertiser["payOrderAmount"]
							}

							if statCost > 0.00 {
								roi = payOrderAmount / statCost
							}

							cost += performances[key]["cost"]
						}
					} else {
						for _, companyUserRole := range companyUserRoles {
							for _, qianchuanAdvertiser := range qianchuanAdvertisers.Data.List {
								if companyUserRole.AdvertiserId == qianchuanAdvertiser.AdvertiserId {
									advertisers = append(advertisers, &domain.Advertiser{
										AdvertiserId:   qianchuanAdvertiser.AdvertiserId,
										AdvertiserName: qianchuanAdvertiser.AdvertiserName,
									})

									break
								}
							}

							key := strconv.FormatUint(companyUserRole.AdvertiserId, 10)

							for _, adAdvertiser := range adAdvertisers[key] {
								statCost += adAdvertiser["statCost"]
								roi += adAdvertiser["roi"]
							}

							cost += performances[key]["cost"]
						}
					}
				}

				var inCompanyPerformanceDaily *domain.CompanyPerformanceDaily

				badvertisers, _ := json.Marshal(advertisers)

				if inCompanyPerformanceDaily, err = cpuc.cpdrepo.GetByUserIdAndCompanyIdAndUpdateDay(ctx, companyUser.Id, companyUser.CompanyId, updateDay); err == nil {
					inCompanyPerformanceDaily.SetCompanyId(ctx, companyId)
					inCompanyPerformanceDaily.SetUserId(ctx, companyUser.Id)
					inCompanyPerformanceDaily.SetAdvertisers(ctx, string(badvertisers))
					inCompanyPerformanceDaily.SetStatCost(ctx, statCost)
					inCompanyPerformanceDaily.SetPayOrderAmount(ctx, payOrderAmount)
					inCompanyPerformanceDaily.SetRoi(ctx, roi)
					inCompanyPerformanceDaily.SetCost(ctx, cost)
					inCompanyPerformanceDaily.SetUpdateDay(ctx, updateDay)
					inCompanyPerformanceDaily.SetUpdateTime(ctx)

					cpuc.cpdrepo.Update(ctx, inCompanyPerformanceDaily)
				} else {
					inCompanyPerformanceDaily = domain.NewCompanyPerformanceDaily(ctx, companyUser.Id, companyId, statCost, payOrderAmount, roi, cost, string(badvertisers))
					inCompanyPerformanceDaily.SetUpdateDay(ctx, updateDay)
					inCompanyPerformanceDaily.SetCreateTime(ctx)
					inCompanyPerformanceDaily.SetUpdateTime(ctx)

					cpuc.cpdrepo.Save(ctx, inCompanyPerformanceDaily)
				}

				if updateTime, err := tool.StringToTime("2006-01-02", updateDay); err == nil {
					if iUpdateDay, err := strconv.ParseUint(tool.TimeToString("200601", updateTime), 10, 64); err == nil {
						var inCompanyPerformanceMonthly *domain.CompanyPerformanceMonthly
						var mstatCost float32 = 0.00
						var mpayOrderAmount float32 = 0.00
						var mroi float32 = 0.00
						var mcost float32 = 0.00
						var mrebalanceCost float32 = 0.00
						var madvertisers string

						if companyUserPerformanceDailies, err := cpuc.cpdrepo.List(ctx, companyUser.Id, companyId, tool.TimeToString("2006-01", updateTime)); err == nil {
							for index, companyUserPerformanceDaily := range companyUserPerformanceDailies {
								if index == 0 {
									madvertisers = companyUserPerformanceDaily.Advertisers
								}

								mstatCost += companyUserPerformanceDaily.StatCost
								mpayOrderAmount += companyUserPerformanceDaily.PayOrderAmount
								mcost += companyUserPerformanceDaily.Cost
							}

							if mstatCost > 0.00 {
								mroi = mpayOrderAmount / mstatCost
							}
						}

						if companyUserPerformanceRebalances, err := cpuc.cprrepo.List(ctx, companyUser.Id, companyId, tool.TimeToString("2006-01", updateTime)); err == nil {
							for _, companyUserPerformanceRebalance := range companyUserPerformanceRebalances {
								if companyUserPerformanceRebalance.RebalanceType == 1 {
									mrebalanceCost += companyUserPerformanceRebalance.Cost
								} else {
									mrebalanceCost -= companyUserPerformanceRebalance.Cost
								}
							}
						}

						if inCompanyPerformanceMonthly, err = cpuc.cpmrepo.GetByUserIdAndCompanyIdAndUpdateDay(ctx, companyUser.Id, companyUser.CompanyId, uint32(iUpdateDay)); err == nil {
							inCompanyPerformanceMonthly.SetCompanyId(ctx, companyId)
							inCompanyPerformanceMonthly.SetUserId(ctx, companyUser.Id)
							inCompanyPerformanceMonthly.SetUsername(ctx, companyUser.Username)
							inCompanyPerformanceMonthly.SetJob(ctx, companyUser.Job)
							inCompanyPerformanceMonthly.SetAdvertisers(ctx, madvertisers)
							inCompanyPerformanceMonthly.SetStatCost(ctx, mstatCost)
							inCompanyPerformanceMonthly.SetRoi(ctx, mroi)
							inCompanyPerformanceMonthly.SetCost(ctx, mcost)
							inCompanyPerformanceMonthly.SetRebalanceCost(ctx, mrebalanceCost)
							inCompanyPerformanceMonthly.SetTotalCost(ctx, mcost+mrebalanceCost)
							inCompanyPerformanceMonthly.SetUpdateDay(ctx, uint32(iUpdateDay))
							inCompanyPerformanceMonthly.SetUpdateTime(ctx)

							cpuc.cpmrepo.Update(ctx, inCompanyPerformanceMonthly)
						} else {
							inCompanyPerformanceMonthly = domain.NewCompanyPerformanceMonthly(ctx, companyUser.Id, companyId, mstatCost, mroi, mcost, mrebalanceCost, (mcost + mrebalanceCost), companyUser.Username, companyUser.Job, madvertisers)
							inCompanyPerformanceMonthly.SetUpdateDay(ctx, uint32(iUpdateDay))
							inCompanyPerformanceMonthly.SetCreateTime(ctx)
							inCompanyPerformanceMonthly.SetUpdateTime(ctx)

							cpuc.cpmrepo.Save(ctx, inCompanyPerformanceMonthly)
						}
					}
				}
			}
		}
	}
}

func (cpuc *CompanyPerformanceUsecase) RobotCompanyPerformances(ctx context.Context, updateDay string) error {
	/*var wg sync.WaitGroup
	wg.Add(1)

	go c.syncCompanyPerformance(ctx, 7, updateDay, &wg)

	wg.Wait()
	spew.Dump("--------------------------------------------")*/

	companys, err := cpuc.crepo.List(ctx, 0, 0, 0, "", "1", 0)

	if err != nil {
		return CompanyCompanyNotFound
	}

	var wg sync.WaitGroup

	for _, company := range companys {
		wg.Add(1)

		go cpuc.syncCompanyPerformance(ctx, company.Id, updateDay, &wg)
	}

	wg.Wait()

	return nil
}

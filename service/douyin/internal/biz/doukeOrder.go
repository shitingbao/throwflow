package biz

import (
	"context"
	"douyin/internal/conf"
	"douyin/internal/domain"
	"douyin/internal/pkg/csj"
	order1 "douyin/internal/pkg/csj/order"
	"douyin/internal/pkg/tool"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
	"math"
	"strconv"
	"sync"
	"time"
)

var (
	DouyinDoukeOrderConfigError            = errors.InternalServer("DOUYIN_DOUKE_ORDER_CONFIG_ERROR", "抖客订单配置信息错误")
	DouyinDoukeOrderMessageVerifySignError = errors.InternalServer("DOUYIN_DOUKE_ORDER_MESSAGE_VERIFY_SIGN_ERROR", "抖客订单消息验签失败")
	DouyinDoukeOrderMessageUnmarshalError  = errors.InternalServer("DOUYIN_DOUKE_ORDER_MESSAGE_UNMARSHAL_ERROR", "抖客订单消息解析失败")
	DouyinDoukeOrderListError              = errors.InternalServer("DOUYIN_DOUKE_ORDER_LIST_ERROR", "抖客订单列表获取失败")
	DouyinDoukeCountByUserIdError          = errors.InternalServer("DOUYIN_DOUKE_COUNT_BY_USER_ID_ERROR", "抖客订单数量获取失败")
)

type DoukeOrderRepo interface {
	SaveIndex(context.Context)
	Upsert(context.Context, *domain.DoukeOrder) error
}

type DoukeOrderUsecase struct {
	repo    DoukeOrderRepo
	doirepo DoukeOrderInfoRepo
	cprepo  CompanyProductRepo
	calrepo CsjApiLogRepo
	wucrepo WeixinUserCommissionRepo
	tlrepo  TaskLogRepo
	conf    *conf.Data
	cconf   *conf.Csj
	log     *log.Helper
}

func NewDoukeOrderUsecase(repo DoukeOrderRepo, doirepo DoukeOrderInfoRepo, cprepo CompanyProductRepo, calrepo CsjApiLogRepo, tlrepo TaskLogRepo, wucrepo WeixinUserCommissionRepo, conf *conf.Data, cconf *conf.Csj, logger log.Logger) *DoukeOrderUsecase {
	return &DoukeOrderUsecase{repo: repo, doirepo: doirepo, cprepo: cprepo, calrepo: calrepo, tlrepo: tlrepo, wucrepo: wucrepo, conf: conf, cconf: cconf, log: log.NewHelper(logger)}
}

func (douc *DoukeOrderUsecase) GetDoukeOrders(ctx context.Context, userId uint64, productId, flowPoint, createTime string) (*domain.DoukeOrderInfo, error) {
	doukeOrder, err := douc.doirepo.GetByUserIdAndProductId(ctx, userId, productId, flowPoint, createTime)

	if err != nil {
		return nil, DouyinDoukeCountByUserIdError
	}

	return doukeOrder, nil
}

func (douc *DoukeOrderUsecase) ListUserIdDoukeOrders(ctx context.Context) ([]*domain.DoukeOrderInfo, error) {
	list, err := douc.doirepo.List(ctx)

	if err != nil {
		return nil, DouyinDoukeOrderListError
	}

	return list, nil
}

func (douc *DoukeOrderUsecase) StatisticsDoukeOrders(ctx context.Context, userId uint64) (*domain.StatisticsDoukeOrders, error) {
	statistics := make([]*domain.StatisticsDoukeOrder, 0)
	statistics = append(statistics, &domain.StatisticsDoukeOrder{
		Key:   "orderNum",
		Value: "0",
	})

	doukeOrder, _ := douc.doirepo.Count(ctx, userId, "", "")

	for _, statistic := range statistics {
		if statistic.Key == "orderNum" {
			statistic.Value = strconv.FormatInt(doukeOrder, 10)
		}
	}

	return &domain.StatisticsDoukeOrders{
		Statistics: statistics,
	}, nil
}

func (douc *DoukeOrderUsecase) StatisticsDoukeOrderByDays(ctx context.Context, userId uint64, day string) (*domain.StatisticsDoukeOrders, error) {
	statistics := make([]*domain.StatisticsDoukeOrder, 0)

	statistics = append(statistics, &domain.StatisticsDoukeOrder{
		Key:   "totalPayAmount",
		Value: "0.00",
	})

	statistics = append(statistics, &domain.StatisticsDoukeOrder{
		Key:   "estimatedCommission",
		Value: "0.00",
	})

	statistics = append(statistics, &domain.StatisticsDoukeOrder{
		Key:   "realCommission",
		Value: "0.00",
	})

	var wg sync.WaitGroup

	var doukeNotRefundOrder *domain.DoukeOrderInfo
	var doukeOrderRealcommission *domain.DoukeOrderInfo

	wg.Add(2)

	go func() {
		defer wg.Done()

		doukeNotRefundOrder, _ = douc.doirepo.Statistics(ctx, userId, "", day, "")
	}()

	go func() {
		defer wg.Done()

		doukeOrderRealcommission, _ = douc.doirepo.StatisticsRealcommission(ctx, userId, day, "")
	}()

	wg.Wait()

	for _, statistic := range statistics {
		if statistic.Key == "totalPayAmount" {
			statistic.Value = fmt.Sprintf("%.2f", tool.Decimal(float64(doukeNotRefundOrder.TotalPayAmount), 2))
		} else if statistic.Key == "estimatedCommission" {
			statistic.Value = fmt.Sprintf("%.2f", tool.Decimal(float64(doukeNotRefundOrder.EstimatedCommission), 2))
		} else if statistic.Key == "realCommission" {
			statistic.Value = fmt.Sprintf("%.2f", tool.Decimal(float64(doukeOrderRealcommission.RealCommission), 2))
		}
	}

	return &domain.StatisticsDoukeOrders{
		Statistics: statistics,
	}, nil
}

func (douc *DoukeOrderUsecase) StatisticsDoukeOrderByPaySuccessTimes(ctx context.Context, userId, productId uint64, flowPoint, startTime, endTime string) (*domain.StatisticsDoukeOrders, error) {
	statistics := make([]*domain.StatisticsDoukeOrder, 0)

	statistics = append(statistics, &domain.StatisticsDoukeOrder{
		Key:   "orderNum",
		Value: "0.00",
	})

	if doukeOrder, err := douc.doirepo.StatisticsByProductId(ctx, userId, productId, startTime, endTime, flowPoint); err == nil {
		for _, statistic := range statistics {
			if statistic.Key == "orderNum" {
				statistic.Value = strconv.FormatUint(doukeOrder.ItemNum, 10)
			}
		}

	}

	return &domain.StatisticsDoukeOrders{
		Statistics: statistics,
	}, nil
}

/*func (douc *DoukeOrderUsecase) AsyncNotificationDoukeOrders(ctx context.Context, eventSign, appId, content string) error {
	if douc.dconf.AppKey != appId {
		return DouyinDoukeOrderConfigError
	}

	if ok := douke.VerifyMessageSign(douc.dconf.AppKey, douc.dconf.AppSecret, content, eventSign); !ok {
		return DouyinDoukeOrderMessageVerifySignError
	}

	var doudianAllianceDistributorOrderResponse []orderMessage.DoudianAllianceDistributorOrderResponse

	if err := json.Unmarshal([]byte(content), &doudianAllianceDistributorOrderResponse); err != nil {
		return DouyinDoukeOrderMessageUnmarshalError
	}

	for _, doukeOrder := range doudianAllianceDistributorOrderResponse {
		var doudianAllianceDistributorOrderDataResponse orderMessage.DoudianAllianceDistributorOrderDataResponse

		if err := json.Unmarshal([]byte(doukeOrder.Data), &doudianAllianceDistributorOrderDataResponse); err == nil {
			pidInfo := domain.PidInfo{
				Pid:           doudianAllianceDistributorOrderDataResponse.PidInfo.Pid,
				ExternalInfo:  doudianAllianceDistributorOrderDataResponse.PidInfo.ExternalInfo,
				MediaTypeName: doudianAllianceDistributorOrderDataResponse.PidInfo.MediaTypeName,
			}

			productTags := domain.ProductTags{
				HasSubsidyTag:     doudianAllianceDistributorOrderDataResponse.ProductTags.HasSubsidyTag,
				HasSupermarketTag: doudianAllianceDistributorOrderDataResponse.ProductTags.HasSupermarketTag,
			}

			inDoukeOrder := domain.NewDoukeOrder(ctx, doudianAllianceDistributorOrderDataResponse.AdsActivityId, doudianAllianceDistributorOrderDataResponse.AdsRealCommission, doudianAllianceDistributorOrderDataResponse.AdsEstimatedCommission, doudianAllianceDistributorOrderDataResponse.TotalPayAmount, doudianAllianceDistributorOrderDataResponse.SettledGoodsAmount, doudianAllianceDistributorOrderDataResponse.ItemNum, doudianAllianceDistributorOrderDataResponse.ShopId, doudianAllianceDistributorOrderDataResponse.PayGoodsAmount, doudianAllianceDistributorOrderDataResponse.AdsDistributorId, doudianAllianceDistributorOrderDataResponse.AdsPromotionTate, doudianAllianceDistributorOrderDataResponse.AuthorUid, doudianAllianceDistributorOrderDataResponse.AuthorAccount, doudianAllianceDistributorOrderDataResponse.MediaType, doudianAllianceDistributorOrderDataResponse.ProductImg, doudianAllianceDistributorOrderDataResponse.UpdateTime, doudianAllianceDistributorOrderDataResponse.PaySuccessTime, doudianAllianceDistributorOrderDataResponse.ProductId, doudianAllianceDistributorOrderDataResponse.FlowPoint, doudianAllianceDistributorOrderDataResponse.SettleTime, doudianAllianceDistributorOrderDataResponse.AuthorBuyinId, doudianAllianceDistributorOrderDataResponse.OrderId, doudianAllianceDistributorOrderDataResponse.ProductName, doudianAllianceDistributorOrderDataResponse.DistributionType, doudianAllianceDistributorOrderDataResponse.ShopName, doudianAllianceDistributorOrderDataResponse.ProductActivityId, doudianAllianceDistributorOrderDataResponse.MaterialId, doudianAllianceDistributorOrderDataResponse.RefundTime, doudianAllianceDistributorOrderDataResponse.ConfirmTime, doudianAllianceDistributorOrderDataResponse.BuyerAppId, doudianAllianceDistributorOrderDataResponse.DistributorRightType, pidInfo, productTags)

			if err := douc.repo.Upsert(ctx, inDoukeOrder); err != nil {
				sinDoukeOrder, _ := json.Marshal(inDoukeOrder)

				inTaskLog := domain.NewTaskLog(ctx, "AsyncNotificationDoukeOrders", fmt.Sprintf("[AsyncNotificationDoukeOrdersError AsyncNotificationDoukeOrderError], Data=%s, Description=%s", sinDoukeOrder, "同步抖客订单数据，插入数据库失败"))
				inTaskLog.SetCreateTime(ctx)

				douc.tlrepo.Save(ctx, inTaskLog)
			}

			paySuccessTime, _ := tool.StringToTime("2006-01-02 15:04:05", doudianAllianceDistributorOrderDataResponse.PaySuccessTime)

			var inDoukeOrderInfo *domain.DoukeOrderInfo

			userId, _ := strconv.ParseUint(pidInfo.ExternalInfo, 10, 64)

			if len(doudianAllianceDistributorOrderDataResponse.SettleTime) > 0 {
				settleTime, _ := tool.StringToTime("2006-01-02 15:04:05", doudianAllianceDistributorOrderDataResponse.SettleTime)

				inDoukeOrderInfo = domain.NewDoukeOrderInfo(ctx, userId, uint64(doudianAllianceDistributorOrderDataResponse.ItemNum), float32(tool.Decimal(float64(doudianAllianceDistributorOrderDataResponse.TotalPayAmount)/100, 2)), float32(tool.Decimal(float64(doudianAllianceDistributorOrderDataResponse.PayGoodsAmount)/100, 2)), float32(tool.Decimal(float64(doudianAllianceDistributorOrderDataResponse.AdsEstimatedCommission)/100, 2)), float32(tool.Decimal(float64(doudianAllianceDistributorOrderDataResponse.AdsRealCommission)/100, 2)), doudianAllianceDistributorOrderDataResponse.OrderId, doudianAllianceDistributorOrderDataResponse.ProductId, doudianAllianceDistributorOrderDataResponse.ProductName, doudianAllianceDistributorOrderDataResponse.ProductImg, doudianAllianceDistributorOrderDataResponse.FlowPoint, paySuccessTime, &settleTime)
			} else {
				inDoukeOrderInfo = domain.NewDoukeOrderInfo(ctx, userId, uint64(doudianAllianceDistributorOrderDataResponse.ItemNum), float32(tool.Decimal(float64(doudianAllianceDistributorOrderDataResponse.TotalPayAmount)/100, 2)), float32(tool.Decimal(float64(doudianAllianceDistributorOrderDataResponse.PayGoodsAmount)/100, 2)), float32(tool.Decimal(float64(doudianAllianceDistributorOrderDataResponse.AdsEstimatedCommission)/100, 2)), float32(tool.Decimal(float64(doudianAllianceDistributorOrderDataResponse.AdsRealCommission)/100, 2)), doudianAllianceDistributorOrderDataResponse.OrderId, doudianAllianceDistributorOrderDataResponse.ProductId, doudianAllianceDistributorOrderDataResponse.ProductName, doudianAllianceDistributorOrderDataResponse.ProductImg, doudianAllianceDistributorOrderDataResponse.FlowPoint, paySuccessTime, nil)
			}

			inDoukeOrderInfo.SetCreateTime(ctx)
			inDoukeOrderInfo.SetUpdateTime(ctx)

			if err := douc.doirepo.Upsert(ctx, inDoukeOrderInfo); err != nil {
				sinDoukeOrderInfo, _ := json.Marshal(inDoukeOrderInfo)

				inTaskLog := domain.NewTaskLog(ctx, "SyncJinritemaiOrders", fmt.Sprintf("[SyncJinritemaiOrdersError SyncJinritemaiOrderError], Data=%s, Description=%s", sinDoukeOrderInfo, "同步抖客订单详情数据，插入数据库失败"))
				inTaskLog.SetCreateTime(ctx)

				douc.tlrepo.Save(ctx, inTaskLog)
			}
		}
	}

	return nil
}
*/

func (douc *DoukeOrderUsecase) SyncDoukeOrders(ctx context.Context, day string) error {
	var wg sync.WaitGroup

	limiter := rate.NewLimiter(0, 60)
	limiter.SetLimit(rate.Limit(60))

	//startTime := day + " " + time.Now().Add(-time.Minute*10).Format("15:04") + ":00"
	startTime := day + " " + time.Now().Add(-time.Hour*5).Format("15:04") + ":00"
	endTime := day + " " + time.Now().Format("15:04") + ":00"

	tstartTime, _ := tool.StringToTime("2006-01-02 15:04:05", startTime)
	tendTime, _ := tool.StringToTime("2006-01-02 15:04:05", endTime)

	wg.Add(1)

	go douc.SyncDoukeOrder(ctx, &wg, limiter, int(tstartTime.Unix()), int(tendTime.Unix()))

	wg.Wait()

	return nil
}

func (douc *DoukeOrderUsecase) SyncDoukeOrder(ctx context.Context, wg *sync.WaitGroup, limiter *rate.Limiter, startTime, endTime int) {
	defer func() {
		wg.Done()
	}()

	orders, err := douc.listOrders(ctx, limiter, startTime, endTime, "0")

	if err == nil {
		for len(orders.Data.Orders) == csj.PageSize50 {
			orders, err = douc.listOrders(ctx, limiter, startTime, endTime, orders.Data.Cursor)
		}
	}
}

func (douc *DoukeOrderUsecase) listOrders(ctx context.Context, limiter *rate.Limiter, startTime, endTime int, cursor string) (*order1.ListResponse, error) {
	var orders *order1.ListResponse
	var err error

	for retryNum := 0; retryNum < 3; retryNum++ {
		limiter.Wait(ctx)

		orders, err = order1.List(startTime, endTime, douc.cconf.AppId, douc.cconf.AppSecret, cursor, uuid.NewString())

		if err != nil {
			if retryNum == 2 {
				inCsjApiLog := domain.NewCsjApiLog(ctx, err.Error())
				inCsjApiLog.SetCreateTime(ctx)

				douc.calrepo.Save(ctx, inCsjApiLog)
			}
		} else {
			for _, order := range orders.Data.Orders {
				inDoukeOrder := domain.NewDoukeOrder(ctx, order.TotalPayAmount, order.PayGoodsAmount, order.AfterSalesStatus, order.EstimatedTechServiceFee, order.EstimatedCommission, order.AdsRealCommission, order.SplitRate, order.OrderId, order.AppId, order.ProductId, order.ProductName, order.AuthorAccount, order.AdsAttribution, order.ProductImg, order.PaySuccessTime, order.RefundTime, order.FlowPoint, order.ExternalInfo, order.SettleTime, order.ConfirmTime, order.MediaTypeName, order.UpdateTime)

				if err := douc.repo.Upsert(ctx, inDoukeOrder); err != nil {
					sinDoukeOrder, _ := json.Marshal(inDoukeOrder)

					inTaskLog := domain.NewTaskLog(ctx, "SyncDoukeOrders", fmt.Sprintf("[SyncDoukeOrdersError SyncDoukeOrderError], Data=%s, Description=%s", sinDoukeOrder, "同步抖客订单数据，插入数据库失败"))
					inTaskLog.SetCreateTime(ctx)

					douc.tlrepo.Save(ctx, inTaskLog)
				}

				userId, _ := strconv.ParseUint(order.ExternalInfo, 10, 64)
				paySuccessTime, _ := tool.StringToTime("2006-01-02 15:04:05", order.PaySuccessTime)

				inDoukeOrderInfo := domain.NewDoukeOrderInfo(ctx, userId, 0, int64(order.AfterSalesStatus), float32(tool.Decimal(float64(order.TotalPayAmount)/float64(100), 2)), float32(tool.Decimal(float64(order.PayGoodsAmount)/float64(100), 2)), float32(tool.Decimal(float64(order.EstimatedCommission)/float64(100), 2)), float32(tool.Decimal(float64(order.AdsRealCommission)/float64(100), 2)), order.OrderId, order.ProductId, order.ProductName, order.ProductImg, order.FlowPoint, paySuccessTime)

				if len(order.SettleTime) > 0 {
					if settleTime, err := tool.StringToTime("2006-01-02 15:04:05", order.SettleTime); err == nil {
						inDoukeOrderInfo.SetSettleTime(ctx, &settleTime)
					}
				}

				if len(order.RefundTime) > 0 {
					if refundTime, err := tool.StringToTime("2006-01-02 15:04:05", order.RefundTime); err == nil {
						inDoukeOrderInfo.SetRefundTime(ctx, &refundTime)
					}
				}

				if len(order.ConfirmTime) > 0 {
					if confirmTime, err := tool.StringToTime("2006-01-02 15:04:05", order.ConfirmTime); err == nil {
						inDoukeOrderInfo.SetConfirmTime(ctx, &confirmTime)
					}
				}

				inDoukeOrderInfo.SetCreateTime(ctx)
				inDoukeOrderInfo.SetUpdateTime(ctx)

				if err := douc.doirepo.Upsert(ctx, inDoukeOrderInfo); err != nil {
					sinDoukeOrderInfo, _ := json.Marshal(inDoukeOrderInfo)

					inTaskLog := domain.NewTaskLog(ctx, "SyncDoukeOrders", fmt.Sprintf("[SyncDoukeOrdersError SyncDoukeOrderError], Data=%s, Description=%s", sinDoukeOrderInfo, "同步抖客订单详情数据，插入数据库失败"))
					inTaskLog.SetCreateTime(ctx)

					douc.tlrepo.Save(ctx, inTaskLog)
				}

				if order.FlowPoint == "SETTLE" {
					if order.AdsRealCommission > 0 {
						douc.wucrepo.CreateCostOrder(ctx, userId, tool.Decimal(float64(order.TotalPayAmount)/float64(100), 2), tool.Decimal(float64(order.AdsRealCommission)/float64(100), 2), order.OrderId, order.ProductId, order.FlowPoint, order.PaySuccessTime)
					}
				} else {
					if order.EstimatedCommission > 0 {
						douc.wucrepo.CreateCostOrder(ctx, userId, tool.Decimal(float64(order.TotalPayAmount)/float64(100), 2), tool.Decimal(float64(order.EstimatedCommission)/float64(100), 2), order.OrderId, order.ProductId, order.FlowPoint, order.PaySuccessTime)
					}
				}
			}

			break
		}
	}

	return orders, err
}

func (douc *DoukeOrderUsecase) OperationDoukeOrders(ctx context.Context) error {
	total, err := douc.doirepo.CountOperation(ctx)

	if err != nil {
		return DouyinDoukeOrderListError
	}

	totalPage := uint64(math.Ceil(float64(total) / float64(40000)))

	var i uint64 = 0

	for ; i < totalPage; i++ {
		if doukeOrders, err := douc.doirepo.ListOperation(ctx, int(i), 40000); err == nil {
			for _, doukeOrder := range doukeOrders {
				if doukeOrder.FlowPoint == "SETTLE" {
					douc.wucrepo.CreateCostOrder(ctx, doukeOrder.UserId, float64(doukeOrder.TotalPayAmount), float64(doukeOrder.RealCommission), doukeOrder.OrderId, doukeOrder.ProductId, doukeOrder.FlowPoint, tool.TimeToString("2006-01-02 15:04:05", doukeOrder.PaySuccessTime))
				} else {
					douc.wucrepo.CreateCostOrder(ctx, doukeOrder.UserId, float64(doukeOrder.TotalPayAmount), float64(doukeOrder.EstimatedCommission), doukeOrder.OrderId, doukeOrder.ProductId, doukeOrder.FlowPoint, tool.TimeToString("2006-01-02 15:04:05", doukeOrder.PaySuccessTime))
				}
			}
		}
	}

	return nil
}

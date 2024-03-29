package biz

import (
	"context"
	"douyin/internal/conf"
	"douyin/internal/domain"
	"douyin/internal/pkg/douke"
	"douyin/internal/pkg/douke/orderMessage"
	"douyin/internal/pkg/tool"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
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
	tlrepo  TaskLogRepo
	conf    *conf.Data
	dconf   *conf.Douke
	log     *log.Helper
}

func NewDoukeOrderUsecase(repo DoukeOrderRepo, doirepo DoukeOrderInfoRepo, cprepo CompanyProductRepo, tlrepo TaskLogRepo, conf *conf.Data, dconf *conf.Douke, logger log.Logger) *DoukeOrderUsecase {
	return &DoukeOrderUsecase{repo: repo, doirepo: doirepo, cprepo: cprepo, tlrepo: tlrepo, conf: conf, dconf: dconf, log: log.NewHelper(logger)}
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

	doukeOrder, _ := douc.doirepo.Statistics(ctx, userId, "", "", "pay_succ")

	for _, statistic := range statistics {
		if statistic.Key == "orderNum" {
			statistic.Value = strconv.FormatUint(doukeOrder.ItemNum, 10)
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

func (douc *DoukeOrderUsecase) AsyncNotificationDoukeOrders(ctx context.Context, eventSign, appId, content string) error {
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

func (douc *DoukeOrderUsecase) GetCompanyTaskUserOrderStatus(ctx context.Context, userId uint64, productId, flowPoint, createTime string) (*domain.DoukeOrderInfo, error) {
	doukeOrder, err := douc.doirepo.GetByUserIdAndProductId(ctx, userId, productId, flowPoint, createTime)

	if err != nil {
		return nil, DouyinDoukeCountByUserIdError
	}

	return doukeOrder, nil
}

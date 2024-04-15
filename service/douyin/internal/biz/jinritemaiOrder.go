package biz

import (
	"context"
	"douyin/internal/conf"
	"douyin/internal/domain"
	"douyin/internal/pkg/event/event"
	"douyin/internal/pkg/jinritemai"
	"douyin/internal/pkg/jinritemai/order"
	"douyin/internal/pkg/jinritemai/orderMessage"
	"douyin/internal/pkg/tool"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/time/rate"
	"math"
	"strconv"
	"sync"
	"time"
)

var (
	DouyinJinritemaiOrderNotFound               = errors.NotFound("DOUYIN_JINRITEMAI_ORDER_NOT_FOUND", "达人订单不存在")
	DouyinJinritemaiOrderListError              = errors.InternalServer("DOUYIN_JINRITEMAI_ORDER_LIST_ERROR", "达人订单列表获取失败")
	DouyinJinritemaiOrderItemNumGetError        = errors.InternalServer("DOUYIN_JINRITEMAI_ORDER_ITEM_NUM_GET_ERROR", "达人订单总销量获取失败")
	DouyinJinritemaiOrderConfigError            = errors.InternalServer("DOUYIN_JINRITEMAI_ORDER_CONFIG_ERROR", "达人订单配置信息错误")
	DouyinJinritemaiOrderMessageVerifySignError = errors.InternalServer("DOUYIN_JINRITEMAI_ORDER_MESSAGE_VERIFY_SIGN_ERROR", "达人订单消息验签失败")
	DouyinJinritemaiOrderMessageUnmarshalError  = errors.InternalServer("DOUYIN_JINRITEMAI_ORDER_MESSAGE_UNMARSHAL_ERROR", "达人订单消息解析失败")
)

type JinritemaiOrderRepo interface {
	ListByPickExtra(context.Context) ([]*domain.JinritemaiOrder, error)
	SaveIndex(context.Context)
	Upsert(context.Context, *domain.JinritemaiOrder) error

	Send(context.Context, event.Event) error
}

type JinritemaiOrderUsecase struct {
	repo       JinritemaiOrderRepo
	joirepo    JinritemaiOrderInfoRepo
	doirepo    DoukeOrderInfoRepo
	wurepo     WeixinUserRepo
	wuodrepo   WeixinUserOpenDouyinRepo
	wucrepo    WeixinUserCommissionRepo
	odtrepo    OpenDouyinTokenRepo
	oduirepo   OpenDouyinUserInfoRepo
	oduiclrepo OpenDouyinUserInfoCreateLogRepo
	odvrepo    OpenDouyinVideoRepo
	tlrepo     TaskLogRepo
	jalrepo    JinritemaiApiLogRepo
	cprepo     CompanyProductRepo
	qasrepo    QianchuanAdvertiserStatusRepo
	qaoirepo   QianchuanAwemeOrderInfoRepo
	conf       *conf.Data
	econf      *conf.Event
	dconf      *conf.Developer
	log        *log.Helper
}

func NewJinritemaiOrderUsecase(repo JinritemaiOrderRepo, joirepo JinritemaiOrderInfoRepo, doirepo DoukeOrderInfoRepo, wurepo WeixinUserRepo, wuodrepo WeixinUserOpenDouyinRepo, wucrepo WeixinUserCommissionRepo, odtrepo OpenDouyinTokenRepo, oduirepo OpenDouyinUserInfoRepo, oduiclrepo OpenDouyinUserInfoCreateLogRepo, odvrepo OpenDouyinVideoRepo, tlrepo TaskLogRepo, jalrepo JinritemaiApiLogRepo, cprepo CompanyProductRepo, qasrepo QianchuanAdvertiserStatusRepo, qaoirepo QianchuanAwemeOrderInfoRepo, conf *conf.Data, econf *conf.Event, dconf *conf.Developer, logger log.Logger) *JinritemaiOrderUsecase {
	return &JinritemaiOrderUsecase{repo: repo, joirepo: joirepo, doirepo: doirepo, wurepo: wurepo, wuodrepo: wuodrepo, wucrepo: wucrepo, odtrepo: odtrepo, oduirepo: oduirepo, oduiclrepo: oduiclrepo, odvrepo: odvrepo, tlrepo: tlrepo, jalrepo: jalrepo, cprepo: cprepo, qasrepo: qasrepo, qaoirepo: qaoirepo, conf: conf, econf: econf, dconf: dconf, log: log.NewHelper(logger)}
}

func (jouc *JinritemaiOrderUsecase) GetStorePreferenceJinritemaiOrders(ctx context.Context, userId uint64) ([]*domain.StorePreference, error) {
	list := make([]*domain.StorePreference, 0)

	weixinUser, err := jouc.wurepo.GetById(ctx, userId)

	if err != nil {
		return nil, DouyinWeixinUserNotFound
	}

	weixinUserOpenDouyins, err := jouc.wuodrepo.List(ctx, weixinUser.Data.UserId)

	if err != nil {
		return nil, DouyinWeixinUserOpenDouyinListError
	}

	day := time.Now().Format("2006-01-02")

	if len(weixinUserOpenDouyins.Data.List) > 0 {
		openDouyinUserInfos := make([]*domain.OpenDouyinUserInfo, 0)

		for _, weixinUserOpenDouyin := range weixinUserOpenDouyins.Data.List {
			if weixinUserOpenDouyin.CreateDate != day {
				openDouyinUserInfos = append(openDouyinUserInfos, &domain.OpenDouyinUserInfo{
					ClientKey: weixinUserOpenDouyin.ClientKey,
					OpenId:    weixinUserOpenDouyin.OpenId,
				})
			}
		}

		statisticsAwemeIndustries, _ := jouc.joirepo.StatisticsAwemeIndustry(ctx, 0, "", "", openDouyinUserInfos)

		var itemNum uint64 = 0

		for _, statisticsAwemeIndustry := range statisticsAwemeIndustries {
			isNotExist := true

			itemNum += statisticsAwemeIndustry.ItemNum

			for _, l := range list {
				if statisticsAwemeIndustry.IndustryId == l.IndustryId {
					l.ItemNum += statisticsAwemeIndustry.ItemNum

					isNotExist = false

					break
				}
			}

			if isNotExist {
				list = append(list, &domain.StorePreference{
					IndustryId:   statisticsAwemeIndustry.IndustryId,
					IndustryName: statisticsAwemeIndustry.IndustryName,
					ItemNum:      statisticsAwemeIndustry.ItemNum,
				})
			}
		}

		for _, l := range list {
			l.SetIndustryRatio(ctx, itemNum)
		}
	}

	return list, nil
}

func (jouc *JinritemaiOrderUsecase) GetIsTopJinritemaiOrders(ctx context.Context, productId uint64) (*domain.JinritemaiOrderInfo, error) {
	jinritemaiOrder, err := jouc.joirepo.GetIsTopByProductId(ctx, productId)

	if err != nil {
		return nil, DouyinJinritemaiOrderNotFound
	}

	return jinritemaiOrder, nil
}

func (jouc *JinritemaiOrderUsecase) ListJinritemaiOrders(ctx context.Context, pageNum, pageSize, userId uint64, startDay, endDay string) (*domain.JinritemaiOrderList, error) {
	weixinUser, err := jouc.wurepo.GetById(ctx, userId)

	if err != nil {
		return nil, DouyinWeixinUserNotFound
	}

	weixinUserOpenDouyins, err := jouc.wuodrepo.List(ctx, weixinUser.Data.UserId)

	if err != nil {
		return nil, DouyinWeixinUserOpenDouyinListError
	}

	var total int64
	list := make([]*domain.JinritemaiOrderInfo, 0)

	if len(weixinUserOpenDouyins.Data.List) > 0 {
		openDouyinTokens := make([]*domain.OpenDouyinToken, 0)

		for _, weixinUserOpenDouyin := range weixinUserOpenDouyins.Data.List {
			openDouyinTokens = append(openDouyinTokens, &domain.OpenDouyinToken{
				ClientKey: weixinUserOpenDouyin.ClientKey,
				OpenId:    weixinUserOpenDouyin.OpenId,
			})
		}

		jinritemaiOrders, err := jouc.joirepo.List(ctx, int(pageNum), int(pageSize), openDouyinTokens, startDay, endDay)

		if err != nil {
			return nil, DouyinJinritemaiOrderListError
		}

		total, err = jouc.joirepo.Count(ctx, openDouyinTokens, startDay, endDay)

		if err != nil {
			return nil, DouyinJinritemaiOrderListError
		}

		for _, jinritemaiOrder := range jinritemaiOrders {
			avatar := ""
			var isShow uint8 = 0

			if jinritemaiOrder.MediaType != "video" {
				for _, weixinUserOpenDouyin := range weixinUserOpenDouyins.Data.List {
					if weixinUserOpenDouyin.OpenId == jinritemaiOrder.OpenId && weixinUserOpenDouyin.ClientKey == jinritemaiOrder.ClientKey {
						avatar = weixinUserOpenDouyin.Avatar

						if len(weixinUserOpenDouyin.CooperativeCode) == 0 {
							isShow = 1
						}

						break
					}
				}
			}

			jinritemaiOrder.SetMediaTypeName(ctx)

			list = append(list, &domain.JinritemaiOrderInfo{
				ClientKey:          jinritemaiOrder.ClientKey,
				OpenId:             jinritemaiOrder.OpenId,
				ProductId:          jinritemaiOrder.ProductId,
				ProductName:        jinritemaiOrder.ProductName,
				ProductImg:         jinritemaiOrder.ProductImg,
				TotalPayAmount:     jinritemaiOrder.TotalPayAmount,
				RealCommission:     jinritemaiOrder.EstimatedCommission,
				RealCommissionRate: jinritemaiOrder.GetRealCommissionRate(ctx),
				ItemNum:            jinritemaiOrder.ItemNum,
				MediaType:          jinritemaiOrder.MediaType,
				MediaTypeName:      jinritemaiOrder.MediaTypeName,
				MediaId:            jinritemaiOrder.MediaId,
				MediaCover:         jinritemaiOrder.MediaCover,
				Avatar:             avatar,
				IsShow:             isShow,
			})
		}
	}

	return &domain.JinritemaiOrderList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (jouc *JinritemaiOrderUsecase) ListProductIdAndVideoIdJinritemaiOrders(ctx context.Context, pageNum, pageSize uint64) (*domain.JinritemaiOrderInfoList, error) {
	list := make([]*domain.JinritemaiOrderInfo, 0)

	list, err := jouc.joirepo.ListProductIdAndMediaIds(ctx, int(pageNum), int(pageSize))

	if err != nil {
		return nil, DouyinJinritemaiOrderListError
	}

	total, err := jouc.joirepo.CountProductIdAndMediaIds(ctx)

	if err != nil {
		return nil, DouyinJinritemaiOrderListError
	}

	return &domain.JinritemaiOrderInfoList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (jouc *JinritemaiOrderUsecase) ListJinritemaiOrderByPickExtras(ctx context.Context) ([]*domain.JinritemaiOrder, error) {
	list, err := jouc.repo.ListByPickExtra(ctx)

	if err != nil {
		return nil, DouyinOpenDouyinUserInfoListError
	}

	return list, nil
}

func (jouc *JinritemaiOrderUsecase) ListCommissionRateJinritemaiOrders(ctx context.Context, content string) ([]*domain.JinritemaiOrderInfo, error) {
	list := make([]*domain.JinritemaiOrderInfo, 0)

	var commissionRateJinritemaiOrders []*domain.CommissionRateJinritemaiOrder

	if err := json.Unmarshal([]byte(content), &commissionRateJinritemaiOrders); err != nil {
		return nil, DouyinJinritemaiOrderListError
	}

	if len(commissionRateJinritemaiOrders) == 0 {
		return list, nil
	}

	list, err := jouc.joirepo.ListByProductIdAndMediaIds(ctx, commissionRateJinritemaiOrders)

	if err != nil {
		return nil, DouyinJinritemaiOrderListError
	}

	return list, nil
}

func (jouc *JinritemaiOrderUsecase) StatisticsJinritemaiOrders(ctx context.Context, userId uint64, startDay, endDay string) (*domain.StatisticsJinritemaiOrders, error) {
	weixinUser, err := jouc.wurepo.GetById(ctx, userId)

	if err != nil {
		return nil, DouyinWeixinUserNotFound
	}

	weixinUserOpenDouyins, err := jouc.wuodrepo.List(ctx, weixinUser.Data.UserId)

	if err != nil {
		return nil, DouyinWeixinUserOpenDouyinListError
	}

	statistics := make([]*domain.StatisticsJinritemaiOrder, 0)
	statistics = append(statistics, &domain.StatisticsJinritemaiOrder{
		Key:   "带货销量",
		Value: "0",
	})

	statistics = append(statistics, &domain.StatisticsJinritemaiOrder{
		Key:   "带货销额",
		Value: "0.00",
	})

	statistics = append(statistics, &domain.StatisticsJinritemaiOrder{
		Key:   "带货佣金",
		Value: "0.00",
	})

	statistics = append(statistics, &domain.StatisticsJinritemaiOrder{
		Key:   "成本购订单",
		Value: "0",
	})

	statistics = append(statistics, &domain.StatisticsJinritemaiOrder{
		Key:   "成本购金额",
		Value: "0.00",
	})

	statistics = append(statistics, &domain.StatisticsJinritemaiOrder{
		Key:   "成本购返佣",
		Value: "0.00",
	})

	statistics = append(statistics, &domain.StatisticsJinritemaiOrder{
		Key:   "有效销量",
		Value: "0",
	})

	statistics = append(statistics, &domain.StatisticsJinritemaiOrder{
		Key:   "有效销额",
		Value: "0.00",
	})

	statistics = append(statistics, &domain.StatisticsJinritemaiOrder{
		Key:   "有效佣金",
		Value: "0",
	})

	statistics = append(statistics, &domain.StatisticsJinritemaiOrder{
		Key:   "佣金率",
		Value: "0%",
	})

	statistics = append(statistics, &domain.StatisticsJinritemaiOrder{
		Key:   "退款率",
		Value: "0%",
	})

	statistics = append(statistics, &domain.StatisticsJinritemaiOrder{
		Key:   "投流",
		Value: "0.00",
	})

	var doukeOrder *domain.DoukeOrderInfo
	var doukeOrderNum int64

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()

		doukeOrder, _ = jouc.doirepo.Statistics(ctx, userId, startDay, endDay, "pay_succ")
	}()

	go func() {
		defer wg.Done()

		doukeOrderNum, _ = jouc.doirepo.Count(ctx, userId, startDay, endDay)
	}()

	if len(weixinUserOpenDouyins.Data.List) > 0 {
		openDouyinTokens := make([]*domain.OpenDouyinToken, 0)

		for _, weixinUserOpenDouyin := range weixinUserOpenDouyins.Data.List {
			openDouyinTokens = append(openDouyinTokens, &domain.OpenDouyinToken{
				ClientKey: weixinUserOpenDouyin.ClientKey,
				OpenId:    weixinUserOpenDouyin.OpenId,
			})
		}

		wg.Add(3)

		var jinritemaiOrder, jinritemaiRefundOrder, jinritemaiOrderRealcommission *domain.JinritemaiOrderInfo

		go func() {
			defer wg.Done()

			jinritemaiOrder, _ = jouc.joirepo.Statistics(ctx, openDouyinTokens, startDay, endDay, "pay_succ", "")
		}()

		go func() {
			defer wg.Done()

			jinritemaiRefundOrder, _ = jouc.joirepo.Statistics(ctx, openDouyinTokens, startDay, endDay, "refund", "")
		}()

		go func() {
			defer wg.Done()

			jinritemaiOrderRealcommission, _ = jouc.joirepo.StatisticsRealcommission(ctx, openDouyinTokens, startDay, endDay, "")
		}()

		wg.Wait()

		for _, statistic := range statistics {
			if statistic.Key == "带货销量" {
				statistic.Value = strconv.FormatUint(jinritemaiOrder.ItemNum+jinritemaiRefundOrder.ItemNum, 10)
			} else if statistic.Key == "带货销额" {
				statistic.Value = strconv.FormatFloat(tool.Decimal(float64(jinritemaiOrder.TotalPayAmount+jinritemaiRefundOrder.TotalPayAmount), 2), 'f', 2, 64)
			} else if statistic.Key == "带货佣金" {
				statistic.Value = jinritemaiOrder.GetEstimatedCommission(ctx)
			} else if statistic.Key == "成本购订单" {
				statistic.Value = strconv.FormatInt(doukeOrderNum, 10)
			} else if statistic.Key == "成本购金额" {
				statistic.Value = doukeOrder.GetTotalPayAmount(ctx)
			} else if statistic.Key == "成本购返佣" {
				statistic.Value = doukeOrder.GetEstimatedCommission(ctx)
			} else if statistic.Key == "有效销量" {
				statistic.Value = jinritemaiOrder.GetItemNum(ctx)
			} else if statistic.Key == "有效销额" {
				statistic.Value = jinritemaiOrder.GetTotalPayAmount(ctx)
			} else if statistic.Key == "有效佣金" {
				statistic.Value = jinritemaiOrderRealcommission.GetRealCommission(ctx)
			} else if statistic.Key == "佣金率" {
				statistic.Value = jinritemaiOrder.GetRealCommissionRate(ctx)
			} else if statistic.Key == "退款率" {
				var refundRate float64

				if (jinritemaiOrder.ItemNum + jinritemaiRefundOrder.ItemNum) > 0 {
					refundRate = float64(jinritemaiRefundOrder.ItemNum) / float64(jinritemaiOrder.ItemNum+jinritemaiRefundOrder.ItemNum)
				}

				statistic.Value = strconv.FormatFloat(tool.Decimal(refundRate*100, 2), 'f', 2, 64) + "%"
			}
		}
	} else {
		wg.Wait()

		for _, statistic := range statistics {
			if statistic.Key == "成本购订单" {
				statistic.Value = strconv.FormatInt(doukeOrderNum, 10)
			} else if statistic.Key == "成本购金额" {
				statistic.Value = doukeOrder.GetTotalPayAmount(ctx)
			} else if statistic.Key == "成本购返佣" {
				statistic.Value = doukeOrder.GetEstimatedCommission(ctx)
			}
		}
	}

	return &domain.StatisticsJinritemaiOrders{
		Statistics: statistics,
	}, nil
}

func (jouc *JinritemaiOrderUsecase) StatisticsJinritemaiOrderByClientKeyAndOpenIds(ctx context.Context, clientKey, openId, startDay, endDay string) (*domain.StatisticsJinritemaiOrders, error) {
	statistics := make([]*domain.StatisticsJinritemaiOrder, 0)
	statistics = append(statistics, &domain.StatisticsJinritemaiOrder{
		Key:   "销量",
		Value: "0",
	})

	statistics = append(statistics, &domain.StatisticsJinritemaiOrder{
		Key:   "销额",
		Value: "0.00",
	})

	openDouyinTokens := make([]*domain.OpenDouyinToken, 0)
	openDouyinTokens = append(openDouyinTokens, &domain.OpenDouyinToken{
		ClientKey: clientKey,
		OpenId:    openId,
	})

	jinritemaiOrder, _ := jouc.joirepo.Statistics(ctx, openDouyinTokens, startDay, endDay, "pay_succ", "")

	for _, statistic := range statistics {
		if statistic.Key == "销量" {
			statistic.Value = jinritemaiOrder.GetItemNum(ctx)
		} else if statistic.Key == "销额" {
			statistic.Value = jinritemaiOrder.GetTotalPayAmount(ctx)
		}
	}

	return &domain.StatisticsJinritemaiOrders{
		Statistics: statistics,
	}, nil
}

func (jouc *JinritemaiOrderUsecase) StatisticsJinritemaiOrderByDays(ctx context.Context, day, content, pickExtra string) (*domain.StatisticsJinritemaiOrders, error) {
	statistics := make([]*domain.StatisticsJinritemaiOrder, 0)
	statistics = append(statistics, &domain.StatisticsJinritemaiOrder{
		Key:   "orderNum",
		Value: "0",
	})

	statistics = append(statistics, &domain.StatisticsJinritemaiOrder{
		Key:   "orderRefundNum",
		Value: "0",
	})

	statistics = append(statistics, &domain.StatisticsJinritemaiOrder{
		Key:   "totalPayAmount",
		Value: "0.00",
	})

	statistics = append(statistics, &domain.StatisticsJinritemaiOrder{
		Key:   "estimatedCommission",
		Value: "0.00",
	})

	statistics = append(statistics, &domain.StatisticsJinritemaiOrder{
		Key:   "realCommission",
		Value: "0.00",
	})

	var openDouyinTokens []*domain.OpenDouyinToken

	if err := json.Unmarshal([]byte(content), &openDouyinTokens); err == nil {
		var wg sync.WaitGroup

		var jinritemaiRefundOrder *domain.JinritemaiOrderInfo
		var jinritemaiNotRefundOrder *domain.JinritemaiOrderInfo
		var jinritemaiOrderRealcommission *domain.JinritemaiOrderInfo

		wg.Add(3)

		go func() {
			defer wg.Done()

			jinritemaiRefundOrder, _ = jouc.joirepo.Statistics(ctx, openDouyinTokens, "", day, "refund", pickExtra)
		}()

		go func() {
			defer wg.Done()

			jinritemaiNotRefundOrder, _ = jouc.joirepo.Statistics(ctx, openDouyinTokens, "", day, "", pickExtra)
		}()

		go func() {
			defer wg.Done()

			jinritemaiOrderRealcommission, _ = jouc.joirepo.StatisticsRealcommission(ctx, openDouyinTokens, day, "", pickExtra)
		}()

		wg.Wait()

		for _, statistic := range statistics {
			if statistic.Key == "orderRefundNum" {
				statistic.Value = strconv.FormatUint(jinritemaiRefundOrder.ItemNum, 10)
			} else if statistic.Key == "orderNum" {
				statistic.Value = strconv.FormatUint(jinritemaiRefundOrder.ItemNum+jinritemaiNotRefundOrder.ItemNum, 10)
			} else if statistic.Key == "totalPayAmount" {
				statistic.Value = fmt.Sprintf("%.2f", tool.Decimal(float64(jinritemaiNotRefundOrder.TotalPayAmount), 2))
			} else if statistic.Key == "estimatedCommission" {
				statistic.Value = fmt.Sprintf("%.2f", tool.Decimal(float64(jinritemaiNotRefundOrder.EstimatedCommission), 2))
			} else if statistic.Key == "realCommission" {
				statistic.Value = fmt.Sprintf("%.2f", tool.Decimal(float64(jinritemaiOrderRealcommission.RealCommission), 2))
			}
		}
	}

	return &domain.StatisticsJinritemaiOrders{
		Statistics: statistics,
	}, nil
}

func (jouc *JinritemaiOrderUsecase) StatisticsJinritemaiOrderByPayTimeAndDays(ctx context.Context, day, payTime, content, pickExtra string) (*domain.StatisticsJinritemaiOrders, error) {
	statistics := make([]*domain.StatisticsJinritemaiOrder, 0)
	statistics = append(statistics, &domain.StatisticsJinritemaiOrder{
		Key:   "orderNum",
		Value: "0",
	})

	statistics = append(statistics, &domain.StatisticsJinritemaiOrder{
		Key:   "orderRefundNum",
		Value: "0",
	})

	statistics = append(statistics, &domain.StatisticsJinritemaiOrder{
		Key:   "totalPayAmount",
		Value: "0.00",
	})

	statistics = append(statistics, &domain.StatisticsJinritemaiOrder{
		Key:   "estimatedCommission",
		Value: "0.00",
	})

	statistics = append(statistics, &domain.StatisticsJinritemaiOrder{
		Key:   "realCommission",
		Value: "0.00",
	})

	var openDouyinTokens []*domain.OpenDouyinToken

	if err := json.Unmarshal([]byte(content), &openDouyinTokens); err == nil {
		var wg sync.WaitGroup

		var jinritemaiRefundOrder *domain.JinritemaiOrderInfo
		var jinritemaiNotRefundOrder *domain.JinritemaiOrderInfo
		var jinritemaiOrderRealcommission *domain.JinritemaiOrderInfo

		wg.Add(3)

		go func() {
			defer wg.Done()

			jinritemaiRefundOrder, _ = jouc.joirepo.Statistics(ctx, openDouyinTokens, "", day, "refund", pickExtra)
		}()

		go func() {
			defer wg.Done()

			jinritemaiNotRefundOrder, _ = jouc.joirepo.Statistics(ctx, openDouyinTokens, "", day, "", pickExtra)
		}()

		go func() {
			defer wg.Done()

			jinritemaiOrderRealcommission, _ = jouc.joirepo.StatisticsRealcommissionPayTime(ctx, openDouyinTokens, payTime, day, day, pickExtra)
		}()

		wg.Wait()

		for _, statistic := range statistics {
			if statistic.Key == "orderRefundNum" {
				statistic.Value = strconv.FormatUint(jinritemaiRefundOrder.ItemNum, 10)
			} else if statistic.Key == "orderNum" {
				statistic.Value = strconv.FormatUint(jinritemaiRefundOrder.ItemNum+jinritemaiNotRefundOrder.ItemNum, 10)
			} else if statistic.Key == "totalPayAmount" {
				statistic.Value = fmt.Sprintf("%.2f", tool.Decimal(float64(jinritemaiNotRefundOrder.TotalPayAmount), 2))
			} else if statistic.Key == "estimatedCommission" {
				statistic.Value = fmt.Sprintf("%.2f", tool.Decimal(float64(jinritemaiNotRefundOrder.EstimatedCommission), 2))
			} else if statistic.Key == "realCommission" {
				statistic.Value = fmt.Sprintf("%.2f", tool.Decimal(float64(jinritemaiOrderRealcommission.RealCommission), 2))
			}
		}
	}

	return &domain.StatisticsJinritemaiOrders{
		Statistics: statistics,
	}, nil
}

func (jouc *JinritemaiOrderUsecase) AsyncNotificationJinritemaiOrders(ctx context.Context, msgId, sign, content string) error {
	if ok := jinritemai.VerifyMessageSign(jouc.dconf.Aweme.OpenDouyin.ClientSecret, content, sign); !ok {
		return DouyinJinritemaiOrderMessageVerifySignError
	}

	var allianceDarenOrderResponse orderMessage.AllianceDarenOrderResponse

	if err := json.Unmarshal([]byte(content), &allianceDarenOrderResponse); err != nil {
		return DouyinJinritemaiOrderMessageUnmarshalError
	}

	if jouc.dconf.Aweme.OpenDouyin.ClientKey != allianceDarenOrderResponse.ClientKey {
		return DouyinJinritemaiOrderConfigError
	}

	var allianceDarenOrderDataResponse orderMessage.AllianceDarenOrderDataResponse

	if err := json.Unmarshal([]byte(allianceDarenOrderResponse.Content), &allianceDarenOrderDataResponse); err != nil {
		return DouyinJinritemaiOrderMessageUnmarshalError
	}

	pidInfo := domain.PidInfo{
		Pid:           allianceDarenOrderDataResponse.PidInfo.Pid,
		ExternalInfo:  allianceDarenOrderDataResponse.PidInfo.ExternalInfo,
		MediaTypeName: allianceDarenOrderDataResponse.PidInfo.MediaTypeName,
	}

	inJinritemaiOrder := domain.NewJinritemaiOrder(ctx, allianceDarenOrderDataResponse.PayGoodsAmount, allianceDarenOrderDataResponse.SettledGoodsAmount, allianceDarenOrderDataResponse.EstimatedCommission, allianceDarenOrderDataResponse.RealCommission, allianceDarenOrderDataResponse.ItemNum, allianceDarenOrderDataResponse.ShopId, allianceDarenOrderDataResponse.EstimatedTotalCommission, allianceDarenOrderDataResponse.EstimatedTechServiceFee, allianceDarenOrderDataResponse.PlatformSubsidy, allianceDarenOrderDataResponse.AuthorSubsidy, allianceDarenOrderDataResponse.AppId, allianceDarenOrderDataResponse.SettleUserSteppedCommission, allianceDarenOrderDataResponse.SettleInstSteppedCommission, allianceDarenOrderDataResponse.PaySubsidy, allianceDarenOrderDataResponse.MediaId, allianceDarenOrderDataResponse.EstimatedInstSteppedCommission, allianceDarenOrderDataResponse.EstimatedUserSteppedCommission, allianceDarenOrderDataResponse.TotalPayAmount, allianceDarenOrderDataResponse.CommissionRate, allianceDarenOrderDataResponse.IsSteppedPlan, allianceDarenOrderDataResponse.OrderId, allianceDarenOrderDataResponse.ProductId, allianceDarenOrderDataResponse.ProductName, allianceDarenOrderDataResponse.ProductImg, allianceDarenOrderDataResponse.AuthorAccount, allianceDarenOrderResponse.ClientKey, allianceDarenOrderDataResponse.AuthorOpenId, allianceDarenOrderDataResponse.ShopName, allianceDarenOrderDataResponse.FlowPoint, allianceDarenOrderDataResponse.App, allianceDarenOrderDataResponse.UpdateTime, allianceDarenOrderDataResponse.PaySuccessTime, allianceDarenOrderDataResponse.SettleTime, allianceDarenOrderDataResponse.Extra, allianceDarenOrderDataResponse.RefundTime, allianceDarenOrderDataResponse.PickSourceClientKey, allianceDarenOrderDataResponse.PickExtra, allianceDarenOrderDataResponse.AuthorShortId, allianceDarenOrderDataResponse.MediaType, allianceDarenOrderDataResponse.AuthorBuyinId, allianceDarenOrderDataResponse.ConfirmTime, allianceDarenOrderDataResponse.ProductActivityId, pidInfo)

	if err := jouc.repo.Upsert(ctx, inJinritemaiOrder); err != nil {
		sinJinritemaiOrder, _ := json.Marshal(inJinritemaiOrder)

		inTaskLog := domain.NewTaskLog(ctx, "AsyncNotificationJinritemaiOrders", fmt.Sprintf("[AsyncNotificationJinritemaiOrdersError] Data=%s, msg=%s, Description=%s", sinJinritemaiOrder, content, "同步精选联盟达人订单数据，插入数据库失败"))
		inTaskLog.SetCreateTime(ctx)

		jouc.tlrepo.Save(ctx, inTaskLog)
	}

	paySuccessTime, _ := tool.StringToTime("2006-01-02 15:04:05", allianceDarenOrderDataResponse.PaySuccessTime)

	var commissionRate uint8

	if allianceDarenOrderDataResponse.CommissionRate > 0 {
		commissionRate = uint8(allianceDarenOrderDataResponse.CommissionRate / 100)
	} else {
		commissionRate = uint8(allianceDarenOrderDataResponse.CommissionAte / 100)
	}

	openId := ""

	if len(allianceDarenOrderDataResponse.AuthorOpenId) > 0 {
		openId = allianceDarenOrderDataResponse.AuthorOpenId
	} else {
		openId = allianceDarenOrderResponse.FromUserId
	}

	var inJinritemaiOrderInfo *domain.JinritemaiOrderInfo

	if len(allianceDarenOrderDataResponse.SettleTime) > 0 {
		settleTime, _ := tool.StringToTime("2006-01-02 15:04:05", allianceDarenOrderDataResponse.SettleTime)

		inJinritemaiOrderInfo = domain.NewJinritemaiOrderInfo(ctx, uint64(allianceDarenOrderDataResponse.ItemNum), uint64(allianceDarenOrderDataResponse.MediaId), commissionRate, float32(tool.Decimal(float64(allianceDarenOrderDataResponse.TotalPayAmount)/float64(100), 2)), float32(tool.Decimal(float64(allianceDarenOrderDataResponse.PayGoodsAmount)/float64(100), 2)), float32(tool.Decimal(float64(allianceDarenOrderDataResponse.EstimatedCommission)/float64(100), 2)), float32(tool.Decimal(float64(allianceDarenOrderDataResponse.RealCommission)/float64(100), 2)), allianceDarenOrderResponse.ClientKey, openId, "", allianceDarenOrderDataResponse.OrderId, allianceDarenOrderDataResponse.ProductId, allianceDarenOrderDataResponse.ProductName, allianceDarenOrderDataResponse.ProductImg, allianceDarenOrderDataResponse.FlowPoint, allianceDarenOrderDataResponse.PickExtra, allianceDarenOrderDataResponse.MediaType, paySuccessTime, &settleTime)
	} else {
		inJinritemaiOrderInfo = domain.NewJinritemaiOrderInfo(ctx, uint64(allianceDarenOrderDataResponse.ItemNum), uint64(allianceDarenOrderDataResponse.MediaId), commissionRate, float32(tool.Decimal(float64(allianceDarenOrderDataResponse.TotalPayAmount)/float64(100), 2)), float32(tool.Decimal(float64(allianceDarenOrderDataResponse.PayGoodsAmount)/float64(100), 2)), float32(tool.Decimal(float64(allianceDarenOrderDataResponse.EstimatedCommission)/float64(100), 2)), float32(tool.Decimal(float64(allianceDarenOrderDataResponse.RealCommission)/float64(100), 2)), allianceDarenOrderResponse.ClientKey, openId, "", allianceDarenOrderDataResponse.OrderId, allianceDarenOrderDataResponse.ProductId, allianceDarenOrderDataResponse.ProductName, allianceDarenOrderDataResponse.ProductImg, allianceDarenOrderDataResponse.FlowPoint, allianceDarenOrderDataResponse.PickExtra, allianceDarenOrderDataResponse.MediaType, paySuccessTime, nil)
	}

	inJinritemaiOrderInfo.SetCreateTime(ctx)
	inJinritemaiOrderInfo.SetUpdateTime(ctx)

	if err := jouc.joirepo.Upsert(ctx, inJinritemaiOrderInfo); err != nil {
		sinJinritemaiOrderInfo, _ := json.Marshal(inJinritemaiOrderInfo)

		inTaskLog := domain.NewTaskLog(ctx, "AsyncNotificationJinritemaiOrders", fmt.Sprintf("[AsyncNotificationJinritemaiOrdersError] Data=%s, msg=%s, Description=%s", sinJinritemaiOrderInfo, content, "同步精选联盟达人订单详情数据，插入数据库失败"))
		inTaskLog.SetCreateTime(ctx)

		jouc.tlrepo.Save(ctx, inTaskLog)
	}

	if allianceDarenOrderDataResponse.FlowPoint == "SETTLE" {
		if allianceDarenOrderDataResponse.RealCommission > 0 {
			jouc.wucrepo.CreateOrder(ctx, tool.Decimal(float64(allianceDarenOrderDataResponse.TotalPayAmount)/float64(100), 2), tool.Decimal(float64(allianceDarenOrderDataResponse.RealCommission)/float64(100), 2), allianceDarenOrderResponse.ClientKey, openId, allianceDarenOrderDataResponse.OrderId, allianceDarenOrderDataResponse.FlowPoint, allianceDarenOrderDataResponse.PaySuccessTime)
		}
	} else {
		if allianceDarenOrderDataResponse.EstimatedCommission > 0 {
			jouc.wucrepo.CreateOrder(ctx, tool.Decimal(float64(allianceDarenOrderDataResponse.TotalPayAmount)/float64(100), 2), tool.Decimal(float64(allianceDarenOrderDataResponse.EstimatedCommission)/float64(100), 2), allianceDarenOrderResponse.ClientKey, openId, allianceDarenOrderDataResponse.OrderId, allianceDarenOrderDataResponse.FlowPoint, allianceDarenOrderDataResponse.PaySuccessTime)
		}
	}

	return nil
}

func (jouc *JinritemaiOrderUsecase) SyncJinritemaiOrders(ctx context.Context, day string) error {
	openDouyinTokens, err := jouc.odtrepo.List(ctx)

	if err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "SyncJinritemaiOrders", fmt.Sprintf("[SyncJinritemaiOrdersError] Description=%s", "获取抖音开放平台token列表失败"))
		inTaskLog.SetCreateTime(ctx)

		jouc.tlrepo.Save(ctx, inTaskLog)

		return err
	}

	if len(day) == 0 {
		day = tool.TimeToString("2006-01-02", time.Now())
	}

	jouc.repo.SaveIndex(ctx)

	var wg sync.WaitGroup

	limiter := rate.NewLimiter(0, 60)
	limiter.SetLimit(rate.Limit(60))

	for _, openDouyinToken := range openDouyinTokens {
		wg.Add(1)

		startTime := day + " " + time.Now().Add(-time.Minute*10).Format("15:04") + ":00"
		endTime := day + " " + time.Now().Format("15:04") + ":00"

		openDouyinUserInfo := &domain.OpenDouyinUserInfo{}

		openDouyinUserInfo, _ = jouc.oduirepo.GetByClientKeyAndOpenId(ctx, openDouyinToken.ClientKey, openDouyinToken.OpenId)

		go jouc.SyncJinritemaiOrder(ctx, &wg, limiter, startTime, endTime, openDouyinToken, openDouyinUserInfo)
	}

	wg.Wait()

	return nil
}

func (jouc *JinritemaiOrderUsecase) SyncJinritemaiOrder(ctx context.Context, wg *sync.WaitGroup, limiter *rate.Limiter, startTime, endTime string, openDouyinToken *domain.OpenDouyinToken, openDouyinUserInfo *domain.OpenDouyinUserInfo) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("获取抖音开放平台达人订单列表异常，", err)
		}

		wg.Done()
	}()

	orders, err := jouc.listOrders(ctx, limiter, startTime, endTime, "", openDouyinToken, openDouyinUserInfo)

	if err == nil {
		for len(orders.Data.Orders) == jinritemai.PageSize20 {
			orders, err = jouc.listOrders(ctx, limiter, startTime, endTime, orders.Data.Cursor, openDouyinToken, openDouyinUserInfo)
		}
	}
}

func (jouc *JinritemaiOrderUsecase) listOrders(ctx context.Context, limiter *rate.Limiter, startTime, endTime, cursor string, openDouyinToken *domain.OpenDouyinToken, openDouyinUserInfo *domain.OpenDouyinUserInfo) (*order.ListOrderResponse, error) {
	var orders *order.ListOrderResponse
	var err error

	for retryNum := 0; retryNum < 3; retryNum++ {
		limiter.Wait(ctx)

		orders, err = order.ListOrder(openDouyinToken.AccessToken, openDouyinToken.OpenId, startTime, endTime, cursor)

		if err != nil {
			if retryNum == 2 {
				inJinritemaiApiLog := domain.NewJinritemaiApiLog(ctx, openDouyinToken.ClientKey, openDouyinToken.OpenId, openDouyinToken.AccessToken, err.Error())
				inJinritemaiApiLog.SetCreateTime(ctx)

				jouc.jalrepo.Save(ctx, inJinritemaiApiLog)
			}
		} else {
			for _, order := range orders.Data.Orders {
				pidInfo := domain.PidInfo{
					Pid:           order.PidInfo.Pid,
					ExternalInfo:  order.PidInfo.ExternalInfo,
					MediaTypeName: order.PidInfo.MediaTypeName,
				}

				inJinritemaiOrder := domain.NewJinritemaiOrder(ctx, order.PayGoodsAmount, order.SettledGoodsAmount, order.EstimatedCommission, order.RealCommission, order.ItemNum, order.ShopId, order.EstimatedTotalCommission, order.EstimatedTechServiceFee, order.PlatformSubsidy, order.AuthorSubsidy, order.AppId, order.SettleUserSteppedCommission, order.SettleInstSteppedCommission, order.PaySubsidy, order.MediaId, order.EstimatedInstSteppedCommission, order.EstimatedUserSteppedCommission, order.TotalPayAmount, order.CommissionRate, order.IsSteppedPlan, order.OrderId, order.ProductId, order.ProductName, order.ProductImg, order.AuthorAccount, openDouyinToken.ClientKey, order.AuthorOpenId, order.ShopName, order.FlowPoint, order.App, order.UpdateTime, order.PaySuccessTime, order.SettleTime, order.Extra, order.RefundTime, order.PickSourceClientKey, order.PickExtra, order.AuthorShortId, order.MediaType, order.AuthorBuyinId, order.ConfirmTime, order.ProductActivityId, pidInfo)

				if err := jouc.repo.Upsert(ctx, inJinritemaiOrder); err != nil {
					sinJinritemaiOrder, _ := json.Marshal(inJinritemaiOrder)

					inTaskLog := domain.NewTaskLog(ctx, "SyncJinritemaiOrders", fmt.Sprintf("[SyncJinritemaiOrdersError SyncJinritemaiOrderError] ClientKey=%d, OpenId=%d, Data=%s, Description=%s", openDouyinToken.ClientKey, openDouyinToken.OpenId, sinJinritemaiOrder, "同步精选联盟达人订单数据，插入数据库失败"))
					inTaskLog.SetCreateTime(ctx)

					jouc.tlrepo.Save(ctx, inTaskLog)
				}

				paySuccessTime, _ := tool.StringToTime("2006-01-02 15:04:05", order.PaySuccessTime)
				commissionRate := uint8(order.CommissionRate / 100)

				var inJinritemaiOrderInfo *domain.JinritemaiOrderInfo

				if len(order.SettleTime) > 0 {
					settleTime, _ := tool.StringToTime("2006-01-02 15:04:05", order.SettleTime)

					inJinritemaiOrderInfo = domain.NewJinritemaiOrderInfo(ctx, uint64(order.ItemNum), uint64(order.MediaId), commissionRate, float32(tool.Decimal(float64(order.TotalPayAmount)/float64(100), 2)), float32(tool.Decimal(float64(order.PayGoodsAmount)/float64(100), 2)), float32(tool.Decimal(float64(order.EstimatedCommission)/float64(100), 2)), float32(tool.Decimal(float64(order.RealCommission)/float64(100), 2)), openDouyinToken.ClientKey, order.AuthorOpenId, openDouyinUserInfo.BuyinId, order.OrderId, order.ProductId, order.ProductName, order.ProductImg, order.FlowPoint, order.PickExtra, order.MediaType, paySuccessTime, &settleTime)
				} else {
					inJinritemaiOrderInfo = domain.NewJinritemaiOrderInfo(ctx, uint64(order.ItemNum), uint64(order.MediaId), commissionRate, float32(tool.Decimal(float64(order.TotalPayAmount)/float64(100), 2)), float32(tool.Decimal(float64(order.PayGoodsAmount)/float64(100), 2)), float32(tool.Decimal(float64(order.EstimatedCommission)/float64(100), 2)), float32(tool.Decimal(float64(order.RealCommission)/float64(100), 2)), openDouyinToken.ClientKey, order.AuthorOpenId, openDouyinUserInfo.BuyinId, order.OrderId, order.ProductId, order.ProductName, order.ProductImg, order.FlowPoint, order.PickExtra, order.MediaType, paySuccessTime, nil)
				}

				inJinritemaiOrderInfo.SetCreateTime(ctx)
				inJinritemaiOrderInfo.SetUpdateTime(ctx)

				if err := jouc.joirepo.Upsert(ctx, inJinritemaiOrderInfo); err != nil {
					sinJinritemaiOrderInfo, _ := json.Marshal(inJinritemaiOrderInfo)

					inTaskLog := domain.NewTaskLog(ctx, "SyncJinritemaiOrders", fmt.Sprintf("[SyncJinritemaiOrdersError SyncJinritemaiOrderError] ClientKey=%d, OpenId=%d, Data=%s, Description=%s", openDouyinToken.ClientKey, openDouyinToken.OpenId, sinJinritemaiOrderInfo, "同步精选联盟达人订单详情数据，插入数据库失败"))
					inTaskLog.SetCreateTime(ctx)

					jouc.tlrepo.Save(ctx, inTaskLog)
				}
			}

			break
		}
	}

	return orders, err
}

func (jouc *JinritemaiOrderUsecase) Sync90DayJinritemaiOrders(ctx context.Context) error {
	openDouyinUserInfoCreateLogs, err := jouc.oduiclrepo.List(ctx, "0")

	if err != nil {
		return DouyinOpenDouyinUserInfoCreateLogNotFound
	}

	var wg sync.WaitGroup
	wg.Add(1)

	limiter := rate.NewLimiter(0, 60)
	limiter.SetLimit(rate.Limit(60))

	for _, inOpenDouyinUserInfoCreateLog := range openDouyinUserInfoCreateLogs {
		inOpenDouyinUserInfoCreateLog.SetIsHandle(ctx, 1)
		inOpenDouyinUserInfoCreateLog.SetUpdateTime(ctx)

		if _, err := jouc.oduiclrepo.Update(ctx, inOpenDouyinUserInfoCreateLog); err != nil {
			return DouyinOpenDouyinUserInfoCreateLogUpdateError
		}

		openDouyinToken, err := jouc.odtrepo.GetByClientKeyAndOpenId(ctx, inOpenDouyinUserInfoCreateLog.ClientKey, inOpenDouyinUserInfoCreateLog.OpenId)

		if err != nil {
			return DouyinOpenDouyinTokenNotFound
		}

		openDouyinUserInfo, err := jouc.oduirepo.GetByClientKeyAndOpenId(ctx, inOpenDouyinUserInfoCreateLog.ClientKey, inOpenDouyinUserInfoCreateLog.OpenId)

		if err != nil {
			return DouyinOpenDouyinUserInfoNotFound
		}

		wg.Add(1)

		go jouc.Sync90DayJinritemaiOrder(ctx, &wg, limiter, inOpenDouyinUserInfoCreateLog, openDouyinToken, openDouyinUserInfo)
	}

	wg.Wait()

	return nil
}

func (jouc *JinritemaiOrderUsecase) Sync90DayJinritemaiOrder(ctx context.Context, wg *sync.WaitGroup, limiter *rate.Limiter, openDouyinUserInfoCreateLog *domain.OpenDouyinUserInfoCreateLog, openDouyinToken *domain.OpenDouyinToken, openDouyinUserInfo *domain.OpenDouyinUserInfo) {
	defer func() {
		jouc.oduiclrepo.Delete(ctx, openDouyinUserInfoCreateLog)

		wg.Done()
	}()

	var sdjowg sync.WaitGroup

	for i := 90; i >= 0; i-- {
		sdjowg.Add(1)

		day := time.Now().AddDate(0, 0, -i).Format("2006-01-02")

		jouc.repo.SaveIndex(ctx)
		jouc.SyncJinritemaiOrder(ctx, &sdjowg, limiter, day+" 00:00:00", day+" 23:59:59", openDouyinToken, openDouyinUserInfo)
	}

	sdjowg.Wait()
}

func (jouc *JinritemaiOrderUsecase) CompensateJinritemaiOrders(ctx context.Context) error {
	openDouyinUserInfos, err := jouc.oduirepo.List(ctx, 0, 40)

	if err != nil {
		return DouyinOpenDouyinUserInfoListError
	}

	var wg sync.WaitGroup

	limiter := rate.NewLimiter(0, 60)
	limiter.SetLimit(rate.Limit(60))

	for _, openDouyinUserInfo := range openDouyinUserInfos {
		openDouyinToken, err := jouc.odtrepo.GetByClientKeyAndOpenId(ctx, openDouyinUserInfo.ClientKey, openDouyinUserInfo.OpenId)

		if err != nil {
			return DouyinOpenDouyinTokenNotFound
		}

		wg.Add(1)

		go jouc.CompensateJinritemaiOrder(ctx, &wg, limiter, openDouyinToken, openDouyinUserInfo)
	}

	wg.Wait()

	return nil
}

func (jouc *JinritemaiOrderUsecase) CompensateJinritemaiOrder(ctx context.Context, wg *sync.WaitGroup, limiter *rate.Limiter, openDouyinToken *domain.OpenDouyinToken, openDouyinUserInfo *domain.OpenDouyinUserInfo) {
	day := time.Now().AddDate(0, 0, -2).Format("2006-01-02")

	jouc.repo.SaveIndex(ctx)
	jouc.SyncJinritemaiOrder(ctx, wg, limiter, day+" 00:00:00", day+" 23:59:59", openDouyinToken, openDouyinUserInfo)
}

func (jouc *JinritemaiOrderUsecase) OperationJinritemaiOrders(ctx context.Context) error {
	total, err := jouc.joirepo.CountOperation(ctx)

	if err != nil {
		return DouyinJinritemaiOrderListError
	}

	totalPage := uint64(math.Ceil(float64(total) / float64(40000)))

	var i uint64 = 0

	for ; i < totalPage; i++ {
		if jinritemaiOrders, err := jouc.joirepo.ListOperation(ctx, int(i), 40000); err == nil {
			for _, jinritemaiOrder := range jinritemaiOrders {
				if jinritemaiOrder.FlowPoint == "SETTLE" {
					jouc.wucrepo.CreateOrder(ctx, tool.Decimal(float64(jinritemaiOrder.TotalPayAmount), 2), tool.Decimal(float64(jinritemaiOrder.RealCommission), 2), jinritemaiOrder.ClientKey, jinritemaiOrder.OpenId, jinritemaiOrder.OrderId, jinritemaiOrder.FlowPoint, tool.TimeToString("2006-01-02 15:04:05", jinritemaiOrder.PaySuccessTime))
				} else {
					jouc.wucrepo.CreateOrder(ctx, tool.Decimal(float64(jinritemaiOrder.TotalPayAmount), 2), tool.Decimal(float64(jinritemaiOrder.EstimatedCommission), 2), jinritemaiOrder.ClientKey, jinritemaiOrder.OpenId, jinritemaiOrder.OrderId, jinritemaiOrder.FlowPoint, tool.TimeToString("2006-01-02 15:04:05", jinritemaiOrder.PaySuccessTime))
				}
			}
		}
	}

	return nil
}

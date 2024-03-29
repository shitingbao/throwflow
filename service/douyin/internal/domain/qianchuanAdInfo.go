package domain

import (
	"context"
	"douyin/internal/pkg/tool"
	"time"
)

type QianchuanAdInfo struct {
	AdId                    uint64
	AdvertiserId            uint64
	AdvertiserName          string
	CampaignId              uint64
	CampaignName            string
	CampaignBudget          float64
	CampaignBudgetMode      string
	CampaignStatus          string
	CampaignCreateDate      string
	PromotionWay            string
	MarketingGoal           string
	MarketingScene          string
	Name                    string
	Status                  string
	OptStatus               string
	AdCreateTime            string
	AdModifyTime            string
	LabAdType               string
	StatCost                float64
	Roi                     float64
	PayOrderCount           int64
	PayOrderAmount          float64
	CreateOrderAmount       float64
	CreateOrderCount        int64
	ClickCnt                int64
	ShowCnt                 int64
	ConvertCnt              int64
	ClickRate               float64
	CpmPlatform             float64
	DyFollow                int64
	PayConvertRate          float64
	ConvertCost             float64
	ConvertRate             float64
	AveragePayOrderStatCost float64
	PayOrderAveragePrice    float64
	ProductInfo             []*ProductInfo
	AwemeInfo               []*AwemeInfo
	DeliverySetting         DeliverySetting
	CreateTime              time.Time
	UpdateTime              time.Time
}

func NewQianchuanAdInfoQianchuanAd(ctx context.Context, adId, advertiserId uint64, advertiserName, promotionWay, marketingGoal, marketingScene, name, status, optStatus, adCreateTime, adModifyTime, labAdType string, productInfo []*ProductInfo, awemeInfo []*AwemeInfo, deliverySetting DeliverySetting) *QianchuanAdInfo {
	return &QianchuanAdInfo{
		AdId:            adId,
		AdvertiserId:    advertiserId,
		AdvertiserName:  advertiserName,
		PromotionWay:    promotionWay,
		MarketingGoal:   marketingGoal,
		MarketingScene:  marketingScene,
		Name:            name,
		Status:          status,
		OptStatus:       optStatus,
		AdCreateTime:    adCreateTime,
		AdModifyTime:    adModifyTime,
		LabAdType:       labAdType,
		ProductInfo:     productInfo,
		AwemeInfo:       awemeInfo,
		DeliverySetting: deliverySetting,
	}
}

func NewQianchuanAdInfoQianchuanReportAd(ctx context.Context, adId, advertiserId uint64, statCost, payOrderAmount, createOrderAmount float64, payOrderCount, createOrderCount, clickCnt, showCnt, convertCnt, dyFollow int64) *QianchuanAdInfo {
	return &QianchuanAdInfo{
		AdId:              adId,
		AdvertiserId:      advertiserId,
		StatCost:          statCost,
		PayOrderCount:     payOrderCount,
		PayOrderAmount:    payOrderAmount,
		CreateOrderCount:  createOrderCount,
		CreateOrderAmount: createOrderAmount,
		ClickCnt:          clickCnt,
		ShowCnt:           showCnt,
		ConvertCnt:        convertCnt,
		DyFollow:          dyFollow,
	}
}

func NewQianchuanAdQianchuanCampaign(ctx context.Context, advertiserId, campaignId uint64, campaignBudget float64, campaignName, campaignBudgetMode, campaignStatus, campaignCreateDate string) *QianchuanAdInfo {
	return &QianchuanAdInfo{
		AdvertiserId:       advertiserId,
		CampaignId:         campaignId,
		CampaignName:       campaignName,
		CampaignBudget:     campaignBudget,
		CampaignBudgetMode: campaignBudgetMode,
		CampaignStatus:     campaignStatus,
		CampaignCreateDate: campaignCreateDate,
	}
}

func (qai *QianchuanAdInfo) SetAdId(ctx context.Context, adId uint64) {
	qai.AdId = adId
}

func (qai *QianchuanAdInfo) SetAdvertiserId(ctx context.Context, advertiserId uint64) {
	qai.AdvertiserId = advertiserId
}

func (qai *QianchuanAdInfo) SetAdvertiserName(ctx context.Context, advertiserName string) {
	qai.AdvertiserName = advertiserName
}

func (qai *QianchuanAdInfo) SetCampaignId(ctx context.Context, campaignId uint64) {
	qai.CampaignId = campaignId
}

func (qai *QianchuanAdInfo) SetCampaignName(ctx context.Context, campaignName string) {
	qai.CampaignName = campaignName
}

func (qai *QianchuanAdInfo) SetCampaignBudget(ctx context.Context, campaignBudget float64) {
	qai.CampaignBudget = campaignBudget
}

func (qai *QianchuanAdInfo) SetCampaignBudgetMode(ctx context.Context, campaignBudgetMode string) {
	qai.CampaignBudgetMode = campaignBudgetMode
}

func (qai *QianchuanAdInfo) SetCampaignStatus(ctx context.Context, campaignStatus string) {
	qai.CampaignStatus = campaignStatus
}

func (qai *QianchuanAdInfo) SetCampaignCreateDate(ctx context.Context, campaignCreateDate string) {
	qai.CampaignCreateDate = campaignCreateDate
}

func (qai *QianchuanAdInfo) SetPromotionWay(ctx context.Context, promotionWay string) {
	qai.PromotionWay = promotionWay
}

func (qai *QianchuanAdInfo) SetMarketingGoal(ctx context.Context, marketingGoal string) {
	qai.MarketingGoal = marketingGoal
}

func (qai *QianchuanAdInfo) SetMarketingScene(ctx context.Context, marketingScene string) {
	qai.MarketingScene = marketingScene
}

func (qai *QianchuanAdInfo) SetName(ctx context.Context, name string) {
	qai.Name = name
}

func (qai *QianchuanAdInfo) SetStatus(ctx context.Context, status string) {
	qai.Status = status
}

func (qai *QianchuanAdInfo) SetOptStatus(ctx context.Context, optStatus string) {
	qai.OptStatus = optStatus
}

func (qai *QianchuanAdInfo) SetAdCreateTime(ctx context.Context, adCreateTime string) {
	qai.AdCreateTime = adCreateTime
}

func (qai *QianchuanAdInfo) SetAdModifyTime(ctx context.Context, adModifyTime string) {
	qai.AdModifyTime = adModifyTime
}

func (qai *QianchuanAdInfo) SetLabAdType(ctx context.Context, labAdType string) {
	qai.LabAdType = labAdType
}

func (qai *QianchuanAdInfo) SetStatCost(ctx context.Context, statCost float64) {
	qai.StatCost = statCost
}

func (qai *QianchuanAdInfo) SetRoi(ctx context.Context) {
	if qai.StatCost > 0 {
		qai.Roi = tool.Decimal(qai.PayOrderAmount/qai.StatCost, 2)
	}
}

func (qai *QianchuanAdInfo) SetPayOrderCount(ctx context.Context, payOrderCount int64) {
	qai.PayOrderCount = payOrderCount
}

func (qai *QianchuanAdInfo) SetPayOrderAmount(ctx context.Context, payOrderAmount float64) {
	qai.PayOrderAmount = payOrderAmount
}

func (qai *QianchuanAdInfo) SetCreateOrderCount(ctx context.Context, createOrderCount int64) {
	qai.CreateOrderCount = createOrderCount
}

func (qai *QianchuanAdInfo) SetCreateOrderAmount(ctx context.Context, createOrderAmount float64) {
	qai.CreateOrderAmount = createOrderAmount
}

func (qai *QianchuanAdInfo) SetClickCnt(ctx context.Context, clickCnt int64) {
	qai.ClickCnt = clickCnt
}

func (qai *QianchuanAdInfo) SetShowCnt(ctx context.Context, showCnt int64) {
	qai.ShowCnt = showCnt
}

func (qai *QianchuanAdInfo) SetConvertCnt(ctx context.Context, convertCnt int64) {
	qai.ConvertCnt = convertCnt
}

func (qai *QianchuanAdInfo) SetClickRate(ctx context.Context) {
	if qai.ShowCnt > 0 {
		qai.ClickRate = tool.Decimal(float64(qai.ClickCnt)/float64(qai.ShowCnt), 2)
	}
}

func (qai *QianchuanAdInfo) SetCpmPlatform(ctx context.Context) {
	if qai.ShowCnt > 0 {
		qai.CpmPlatform = tool.Decimal(qai.StatCost/float64(qai.ShowCnt)*1000, 2)
	}
}

func (qai *QianchuanAdInfo) SetDyFollow(ctx context.Context, dyFollow int64) {
	qai.DyFollow = dyFollow
}

func (qai *QianchuanAdInfo) SetPayConvertRate(ctx context.Context) {
	if qai.ClickCnt > 0 {
		qai.PayConvertRate = tool.Decimal(float64(qai.PayOrderCount)/float64(qai.ClickCnt), 2)
	}
}

func (qai *QianchuanAdInfo) SetConvertCost(ctx context.Context) {
	if qai.ConvertCnt > 0 {
		qai.ConvertCost = tool.Decimal(float64(qai.StatCost)/float64(qai.ConvertCnt), 2)
	}
}

func (qai *QianchuanAdInfo) SetConvertRate(ctx context.Context) {
	if qai.ClickCnt > 0 {
		qai.ConvertRate = tool.Decimal(float64(qai.ConvertCnt)/float64(qai.ClickCnt), 2)
	}
}

func (qai *QianchuanAdInfo) SetAveragePayOrderStatCost(ctx context.Context) {
	if qai.PayOrderCount > 0 {
		qai.AveragePayOrderStatCost = tool.Decimal(qai.StatCost/float64(qai.PayOrderCount), 2)
	}
}

func (qai *QianchuanAdInfo) SetPayOrderAveragePrice(ctx context.Context) {
	if qai.PayOrderCount > 0 {
		qai.PayOrderAveragePrice = tool.Decimal(qai.PayOrderAmount/float64(qai.PayOrderCount), 2)
	}
}

func (qai *QianchuanAdInfo) SetProductInfo(ctx context.Context, productInfo []*ProductInfo) {
	qai.ProductInfo = productInfo
}

func (qai *QianchuanAdInfo) SetAwemeInfo(ctx context.Context, awemeInfo []*AwemeInfo) {
	qai.AwemeInfo = awemeInfo
}

func (qai *QianchuanAdInfo) SetDeliverySetting(ctx context.Context, deliverySetting DeliverySetting) {
	qai.DeliverySetting = deliverySetting
}

func (qai *QianchuanAdInfo) SetUpdateTime(ctx context.Context) {
	qai.UpdateTime = time.Now()
}

func (qai *QianchuanAdInfo) SetCreateTime(ctx context.Context) {
	qai.CreateTime = time.Now()
}

func (qai *QianchuanAdInfo) GetLabAdTypeName(ctx context.Context) (labAdTypeName string) {
	if qai.LabAdType == "NOT_LAB_AD" {
		labAdTypeName = "非托管计划"
	} else if qai.LabAdType == "LAB_AD" {
		labAdTypeName = "托管计划"
	}

	return
}

func (qai *QianchuanAdInfo) GetMarketingGoalName(ctx context.Context) (marketingGoalName string) {
	if qai.MarketingGoal == "VIDEO_PROM_GOODS" {
		marketingGoalName = "短视频带货"
	} else if qai.MarketingGoal == "LIVE_PROM_GOODS" {
		marketingGoalName = "直播带货"
	}

	return
}

func (qai *QianchuanAdInfo) GetStatusName(ctx context.Context) (statusName string) {
	switch qai.Status {
	case "DELETED":
		statusName = "已删除"
	case "AUDIT":
		statusName = "新建审核中"
	case "TIME_DONE":
		statusName = "已完成"
	case "DISABLE":
		statusName = "已暂停"
	case "TIME_NO_REACH":
		statusName = "未到达投放时间"
	case "OFFLINE_BALANCE":
		statusName = "账户余额不足"
	case "OFFLINE_BUDGET":
		statusName = "广告预算不足（已超出预算）"
	case "DELIVERY_OK":
		statusName = "投放中"
	case "NO_SCHEDULE":
		statusName = "不在投放时段"
	case "REAUDIT":
		statusName = "修改审核中"
	case "OFFLINE_AUDIT":
		statusName = "审核不通过"
	case "EXTERNAL_URL_DISABLE":
		statusName = "落地页暂不可用"
	case "LIVE_ROOM_OFF":
		statusName = "关联直播间未开播"
	case "FROZEN":
		statusName = "已终止"
	case "SYSTEM_DISABLE":
		statusName = "系统暂停"
	case "ALL_INCLUDE_DELETED":
		statusName = "全部（包含已删除）"
	case "QUOTA_DISABLE":
		statusName = "在投计划配额超限"
	case "ROI2_DISABLE":
		statusName = "全域推广暂停"
	case "DELETE":
		statusName = "已删除"
	case "DRAFT":
		statusName = "草稿"
	case "CREATE":
		statusName = "计划新建"
	case "PRE_OFFLINE_BUDGET":
		statusName = "广告预算不足（即将超出预算）"
	case "PRE_ONLINE":
		statusName = "预上线"
	case "ERROR":
		statusName = "数据错误"
	case "AUDIT_STATUS_ERROR":
		statusName = "异常，请联系审核人员"
	case "ADVERTISER_OFFLINE_BUDGET":
		statusName = "账户超出预算"
	case "ADVERTISER_PRE_OFFLINE_BUDGET":
		statusName = "账户接近预算"
	case "CAMPAIGN_DISABLE":
		statusName = "已被广告组暂停"
	case "CAMPAIGN_OFFLINE_BUDGET":
		statusName = "广告组超出预算"
	case "CAMPAIGN_PREOFFLINE_BUDGET":
		statusName = "广告组接近预算"
	}

	return
}

func (qai *QianchuanAdInfo) GetBudgetModeName(ctx context.Context) (budgetModeName string) {
	if qai.DeliverySetting.BudgetMode == "BUDGET_MODE_INFINITE" {
		budgetModeName = "不限"
	} else if qai.DeliverySetting.BudgetMode == "BUDGET_MODE_DAY" {
		budgetModeName = "日预算"
	} else if qai.DeliverySetting.BudgetMode == "BUDGET_MODE_TOTAL" {
		budgetModeName = "总预算"
	}

	return
}

func (qai *QianchuanAdInfo) GetOptStatusName(ctx context.Context) (optStatusName string) {
	switch qai.OptStatus {
	case "ENABLE":
		optStatusName = "启用"
	case "DISABLE":
		optStatusName = "暂停"
	case "DELETE":
		optStatusName = "删除"
	case "QUOTA_DISABLE":
		optStatusName = "因在投计划配额超限而暂停"
	case "ROI2_DISABLE":
		optStatusName = "因该计划关联的抖音号开启全域推广，因此本计划被系统暂停"
	case "SYSTEM_DISABLE":
		optStatusName = "系统暂停，因低效计划被系统自动暂停"
	}

	return
}

func (qai *QianchuanAdInfo) GetExternalActionName(ctx context.Context) (externalActionName string) {
	switch qai.DeliverySetting.ExternalAction {
	case "AD_CONVERT_TYPE_SHOPPING":
		externalActionName = "商品购买"
	case "AD_CONVERT_TYPE_QC_FOLLOW_ACTION":
		externalActionName = "粉丝提升"
	case "AD_CONVERT_TYPE_QC_MUST_BUY":
		externalActionName = "点赞评论"
	case "AD_CONVERT_TYPE_LIVE_ENTER_ACTION":
		externalActionName = "进入直播间"
	case "AD_CONVERT_TYPE_LIVE_CLICK_PRODUCT_ACTION":
		externalActionName = "直播间商品点击"
	case "AD_CONVERT_TYPE_LIVE_SUCCESSORDER_ACTION":
		externalActionName = "直播间下单"
	case "AD_CONVERT_TYPE_NEW_FOLLOW_ACTION":
		externalActionName = "直播间粉丝提升"
	case "AD_CONVERT_TYPE_LIVE_COMMENT_ACTION":
		externalActionName = "直播间评论"
	case "AD_CONVERT_TYPE_LIVE_SUCCESSORDER_PAY":
		if qai.DeliverySetting.DeepExternalAction == "AD_CONVERT_TYPE_LIVE_PAY_ROI" {
			externalActionName = "支付ROI"
		} else {
			externalActionName = "直播间成交"
		}
	case "AD_CONVERT_TYPE_LIVE_PAY_ROI":
		externalActionName = "ROI相关"
	}

	return
}

func (qai *QianchuanAdInfo) GetPromotionId(ctx context.Context) (promotionId uint64) {
	if qai.MarketingGoal == "VIDEO_PROM_GOODS" {
		if len(qai.ProductInfo) > 0 {
			promotionId = qai.ProductInfo[0].Id
		}
	} else if qai.MarketingGoal == "LIVE_PROM_GOODS" {
		if len(qai.AwemeInfo) > 0 {
			promotionId = qai.AwemeInfo[0].AwemeId
		}
	}

	return
}

func (qai *QianchuanAdInfo) GetPromotionShowId(ctx context.Context) (promotionShowId string) {
	if qai.MarketingGoal == "LIVE_PROM_GOODS" {
		if len(qai.AwemeInfo) > 0 {
			promotionShowId = qai.AwemeInfo[0].AwemeShowId
		}
	}

	return
}

func (qai *QianchuanAdInfo) GetPromotionName(ctx context.Context) (promotionName string) {
	if qai.MarketingGoal == "VIDEO_PROM_GOODS" {
		if len(qai.ProductInfo) > 0 {
			promotionName = qai.ProductInfo[0].Name
		}
	} else if qai.MarketingGoal == "LIVE_PROM_GOODS" {
		if len(qai.AwemeInfo) > 0 {
			promotionName = qai.AwemeInfo[0].AwemeName
		}
	}

	return
}

func (qai *QianchuanAdInfo) GetPromotionImg(ctx context.Context) (promotionImg string) {
	if qai.MarketingGoal == "VIDEO_PROM_GOODS" {
		if len(qai.ProductInfo) > 0 {
			promotionImg = qai.ProductInfo[0].Img
		}
	}

	return
}

func (qai *QianchuanAdInfo) GetPromotionAvatar(ctx context.Context) (promotionAvatar string) {
	if qai.MarketingGoal == "LIVE_PROM_GOODS" {
		if len(qai.AwemeInfo) > 0 {
			promotionAvatar = qai.AwemeInfo[0].AwemeAvatar
		}
	}

	return
}

func (qai *QianchuanAdInfo) GetPromotionType(ctx context.Context) (promotionType string) {
	if qai.MarketingGoal == "VIDEO_PROM_GOODS" {
		promotionType = "product"
	} else if qai.MarketingGoal == "LIVE_PROM_GOODS" {
		promotionType = "aweme"
	}

	return
}

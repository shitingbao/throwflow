package domain

import (
	"context"
	"time"
)

type ProductInfo struct {
	Id                  uint64  `json:"id" bson:"id"`
	Name                string  `json:"name" bson:"name"`
	DiscountPrice       float64 `json:"discount_price" bson:"discount_price"`
	Img                 string  `json:"img" bson:"img"`
	MarketPrice         float64 `json:"market_price" bson:"market_price"`
	DiscountLowerPrice  float64 `json:"discount_lower_price" bson:"discount_lower_price"`
	DiscountHigherPrice float64 `json:"discount_higher_price" bson:"discount_higher_price"`
}

type MessageAd struct {
	Type    string `json:"type"`
	Message struct {
		Name     string `json:"name"`
		SyncDate string `json:"sync_date"`
		SyncTime string `json:"sync_time"`
		SendTime string `json:"send_time"`
		Content  string `json:"content"`
	} `json:"message"`
}

type AwemeInfo struct {
	AwemeId     uint64 `json:"aweme_id" bson:"aweme_id"`
	AwemeName   string `json:"aweme_name" bson:"aweme_name"`
	AwemeShowId string `json:"aweme_show_id" bson:"aweme_show_id"`
	AwemeAvatar string `json:"aweme_avatar" bson:"aweme_avatar"`
}

type DeliverySetting struct {
	DeepExternalAction string  `json:"deep_external_action" bson:"deep_external_action"`
	DeepBidType        string  `json:"deep_bid_type" bson:"deep_bid_type"`
	RoiGoal            float64 `json:"roi_goal" bson:"roi_goal"`
	SmartBidType       string  `json:"smart_bid_type" bson:"smart_bid_type"`
	ExternalAction     string  `json:"external_action" bson:"external_action"`
	Budget             float64 `json:"budget" bson:"budget"`
	ReviveBudget       float64 `json:"revive_budget" bson:"revive_budget"`
	BudgetMode         string  `json:"budget_mode" bson:"budget_mode"`
	CpaBid             float64 `json:"cpa_bid" bson:"cpa_bid"`
	StartTime          string  `json:"start_time" bson:"start_time"`
	EndTime            string  `json:"end_time" bson:"end_time"`
}

type QianchuanAd struct {
	AdId            uint64
	AdvertiserId    uint64
	CampaignId      uint64
	PromotionWay    string
	MarketingGoal   string
	MarketingScene  string
	Name            string
	Status          string
	OptStatus       string
	AdCreateTime    string
	AdModifyTime    string
	LabAdType       string
	ProductInfo     []*ProductInfo
	AwemeInfo       []*AwemeInfo
	DeliverySetting DeliverySetting
	CreateTime      time.Time
	UpdateTime      time.Time
}

type QianchuanAdList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*QianchuanAd
}

type ExternalQianchuanAd struct {
	AdId                    uint64
	AdvertiserId            uint64
	AdvertiserName          string
	CampaignId              uint64
	CampaignName            string
	AdName                  string
	LabAdType               string
	LabAdTypeName           string
	MarketingGoal           string
	MarketingGoalName       string
	Status                  string
	StatusName              string
	OptStatus               string
	OptStatusName           string
	ExternalAction          string
	ExternalActionName      string
	DeepExternalAction      string
	DeepBidType             string
	PromotionId             uint64
	PromotionShowId         string
	PromotionName           string
	PromotionImg            string
	PromotionAvatar         string
	PromotionType           string
	StatCost                float64
	Roi                     float64
	YesterdayStatCost       float64
	YesterdayRoi            float64
	YesterdayPayOrderAmount float64
	CpaBid                  float64
	RoiGoal                 float64
	Budget                  float64
	BudgetMode              string
	BudgetModeName          string
	PayOrderCount           int64
	PayOrderAmount          float64
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
	AdCreateTime            time.Time
	AdModifyTime            time.Time
}

type ExternalQianchuanAdList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*ExternalQianchuanAd
}

type StatisticsExternalQianchuanAd struct {
	Key   string
	Value string
}

type StatisticsExternalQianchuanAds struct {
	Statistics []*StatisticsExternalQianchuanAd
}

type Filter struct {
	Key   string
	Value string
}

type SelectExternalQianchuanAds struct {
	Filter []*Filter
}

func NewSelectExternalQianchuanAds() *SelectExternalQianchuanAds {
	filter := make([]*Filter, 0)

	filter = append(filter, &Filter{Key: "labAd", Value: "托管"})
	filter = append(filter, &Filter{Key: "notLabAd", Value: "非托管"})
	filter = append(filter, &Filter{Key: "videoPromGoods", Value: "商品"})
	filter = append(filter, &Filter{Key: "livePromGoods", Value: "直播"})

	return &SelectExternalQianchuanAds{
		Filter: filter,
	}
}

func NewQianchuanAd(ctx context.Context, adId, advertiserId, campaignId uint64, promotionWay, marketingGoal, marketingScene, name, status, optStatus, adCreateTime, adModifyTime, labAdType string, productInfo []*ProductInfo, awemeInfo []*AwemeInfo, deliverySetting DeliverySetting) *QianchuanAd {
	return &QianchuanAd{
		AdId:            adId,
		AdvertiserId:    advertiserId,
		CampaignId:      campaignId,
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

func (qa *QianchuanAd) SetAdId(ctx context.Context, adId uint64) {
	qa.AdId = adId
}

func (qa *QianchuanAd) SetAdvertiserId(ctx context.Context, advertiserId uint64) {
	qa.AdvertiserId = advertiserId
}

func (qa *QianchuanAd) SetCampaignId(ctx context.Context, campaignId uint64) {
	qa.CampaignId = campaignId
}

func (qa *QianchuanAd) SetPromotionWay(ctx context.Context, promotionWay string) {
	qa.PromotionWay = promotionWay
}

func (qa *QianchuanAd) SetMarketingGoal(ctx context.Context, marketingGoal string) {
	qa.MarketingGoal = marketingGoal
}

func (qa *QianchuanAd) SetMarketingScene(ctx context.Context, marketingScene string) {
	qa.MarketingScene = marketingScene
}

func (qa *QianchuanAd) SetName(ctx context.Context, name string) {
	qa.Name = name
}

func (qa *QianchuanAd) SetStatus(ctx context.Context, status string) {
	qa.Status = status
}

func (qa *QianchuanAd) SetOptStatus(ctx context.Context, optStatus string) {
	qa.OptStatus = optStatus
}

func (qa *QianchuanAd) SetAdCreateTime(ctx context.Context, adCreateTime string) {
	qa.AdCreateTime = adCreateTime
}

func (qa *QianchuanAd) SetAdModifyTime(ctx context.Context, adModifyTime string) {
	qa.AdModifyTime = adModifyTime
}

func (qa *QianchuanAd) SetLabAdType(ctx context.Context, labAdType string) {
	qa.LabAdType = labAdType
}

func (qa *QianchuanAd) SetProductInfo(ctx context.Context, productInfo []*ProductInfo) {
	qa.ProductInfo = productInfo
}

func (qa *QianchuanAd) SetAwemeInfo(ctx context.Context, awemeInfo []*AwemeInfo) {
	qa.AwemeInfo = awemeInfo
}

func (qa *QianchuanAd) SetDeliverySetting(ctx context.Context, deliverySetting DeliverySetting) {
	qa.DeliverySetting = deliverySetting
}

func (qa *QianchuanAd) SetUpdateTime(ctx context.Context) {
	qa.UpdateTime = time.Now()
}

func (qa *QianchuanAd) SetCreateTime(ctx context.Context) {
	qa.CreateTime = time.Now()
}

func (eqa *ExternalQianchuanAd) GetRoi(ctx context.Context) (roi float64) {
	if eqa.StatCost > 0 {
		roi = eqa.PayOrderAmount / eqa.StatCost
	}

	return
}

func (eqa *ExternalQianchuanAd) GetLabAdTypeName(ctx context.Context) (labAdTypeName string) {
	if eqa.LabAdType == "NOT_LAB_AD" {
		labAdTypeName = "非托管计划"
	} else if eqa.LabAdType == "LAB_AD" {
		labAdTypeName = "托管计划"
	}

	return
}

func (eqa *ExternalQianchuanAd) GetMarketingGoalName(ctx context.Context) (marketingGoalName string) {
	if eqa.MarketingGoal == "VIDEO_PROM_GOODS" {
		marketingGoalName = "短视频带货"
	} else if eqa.MarketingGoal == "LIVE_PROM_GOODS" {
		marketingGoalName = "直播带货"
	}

	return
}

func (eqa *ExternalQianchuanAd) GetStatusName(ctx context.Context) (statusName string) {
	if eqa.Status == "DELETED" {
		statusName = "已删除"
	} else if eqa.Status == "AUDIT" {
		statusName = "新建审核中"
	} else if eqa.Status == "TIME_DONE" {
		statusName = "已完成"
	} else if eqa.Status == "DISABLE" {
		statusName = "已暂停"
	} else if eqa.Status == "TIME_NO_REACH" {
		statusName = "未到达投放时间"
	} else if eqa.Status == "OFFLINE_BALANCE" {
		statusName = "账户余额不足"
	} else if eqa.Status == "OFFLINE_BUDGET" {
		statusName = "广告预算不足（已超出预算）"
	} else if eqa.Status == "DELIVERY_OK" {
		statusName = "投放中"
	} else if eqa.Status == "NO_SCHEDULE" {
		statusName = "不在投放时段"
	} else if eqa.Status == "REAUDIT" {
		statusName = "修改审核中"
	} else if eqa.Status == "OFFLINE_AUDIT" {
		statusName = "审核不通过"
	} else if eqa.Status == "EXTERNAL_URL_DISABLE" {
		statusName = "落地页暂不可用"
	} else if eqa.Status == "LIVE_ROOM_OFF" {
		statusName = "关联直播间未开播"
	} else if eqa.Status == "FROZEN" {
		statusName = "已终止"
	} else if eqa.Status == "SYSTEM_DISABLE" {
		statusName = "系统暂停"
	} else if eqa.Status == "ALL_INCLUDE_DELETED" {
		statusName = "全部（包含已删除）"
	} else if eqa.Status == "QUOTA_DISABLE" {
		statusName = "在投计划配额超限"
	} else if eqa.Status == "ROI2_DISABLE" {
		statusName = "全域推广暂停"
	} else if eqa.Status == "DELETE" {
		statusName = "已删除"
	} else if eqa.Status == "DRAFT" {
		statusName = "草稿"
	} else if eqa.Status == "CREATE" {
		statusName = "计划新建"
	} else if eqa.Status == "PRE_OFFLINE_BUDGET" {
		statusName = "广告预算不足（即将超出预算）"
	} else if eqa.Status == "PRE_ONLINE" {
		statusName = "预上线"
	} else if eqa.Status == "ERROR" {
		statusName = "数据错误"
	} else if eqa.Status == "AUDIT_STATUS_ERROR" {
		statusName = "异常，请联系审核人员"
	} else if eqa.Status == "ADVERTISER_OFFLINE_BUDGET" {
		statusName = "账户超出预算"
	} else if eqa.Status == "ADVERTISER_PRE_OFFLINE_BUDGET" {
		statusName = "账户接近预算"
	} else if eqa.Status == "CAMPAIGN_DISABLE" {
		statusName = "已被广告组暂停"
	} else if eqa.Status == "CAMPAIGN_OFFLINE_BUDGET" {
		statusName = "广告组超出预算"
	} else if eqa.Status == "CAMPAIGN_PREOFFLINE_BUDGET" {
		statusName = "广告组接近预算"
	}

	return
}

func (eqa *ExternalQianchuanAd) GetOptStatusName(ctx context.Context) (optStatusName string) {
	if eqa.OptStatus == "ENABLE" {
		optStatusName = "启用"
	} else if eqa.OptStatus == "DISABLE" {
		optStatusName = "暂停"
	} else if eqa.OptStatus == "DELETE" {
		optStatusName = "删除"
	} else if eqa.OptStatus == "QUOTA_DISABLE" {
		optStatusName = "因在投计划配额超限而暂停"
	} else if eqa.OptStatus == "ROI2_DISABLE" {
		optStatusName = "因该计划关联的抖音号开启全域推广，因此本计划被系统暂停"
	} else if eqa.OptStatus == "SYSTEM_DISABLE" {
		optStatusName = "系统暂停，因低效计划被系统自动暂停"
	}

	return
}

func (eqa *ExternalQianchuanAd) GetClickRate(ctx context.Context) (clickRate float64) {
	if eqa.ShowCnt > 0 {
		clickRate = float64(eqa.ClickCnt) / float64(eqa.ShowCnt)
	}

	return
}

func (eqa *ExternalQianchuanAd) GetPayConvertRate(ctx context.Context) (payConvertRate float64) {
	if eqa.ClickCnt > 0 {
		payConvertRate = float64(eqa.PayOrderCount) / float64(eqa.ClickCnt)
	}

	return
}

func (eqa *ExternalQianchuanAd) GetConvertCost(ctx context.Context) (convertCost float64) {
	if eqa.ConvertCnt > 0 {
		convertCost = float64(eqa.StatCost) / float64(eqa.ConvertCnt)
	}

	return
}

func (eqa *ExternalQianchuanAd) GetConvertRate(ctx context.Context) (convertRate float64) {
	if eqa.ClickCnt > 0 {
		convertRate = float64(eqa.ConvertCnt) / float64(eqa.ClickCnt)
	}

	return
}

func (eqa *ExternalQianchuanAd) GetAveragePayOrderStatCost(ctx context.Context) (averagePayOrderStatCost float64) {
	if eqa.PayOrderCount > 0 {
		averagePayOrderStatCost = eqa.StatCost / float64(eqa.PayOrderCount)
	}

	return
}

func (eqa *ExternalQianchuanAd) GetPayOrderAveragePrice(ctx context.Context) (payOrderAveragePrice float64) {
	if eqa.PayOrderCount > 0 {
		payOrderAveragePrice = eqa.PayOrderAmount / float64(eqa.PayOrderCount)
	}

	return
}

func (eqa *ExternalQianchuanAd) GetCpmPlatform(ctx context.Context) (cpmPlatform float64) {
	if eqa.ShowCnt > 0 {
		cpmPlatform = eqa.StatCost / float64(eqa.ShowCnt) * 1000
	}

	return
}

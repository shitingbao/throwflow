package domain

import (
	"context"
)

type PidInfo struct {
	Pid           string `json:"pid" bson:"pid"`
	ExternalInfo  string `json:"external_info" bson:"external_info"`
	MediaTypeName string `json:"media_type_name" bson:"media_type_name"`
}

type JinritemaiOrder struct {
	OrderId                        string
	ProductId                      string
	ProductName                    string
	ProductImg                     string
	AuthorAccount                  string
	AuthorClientKey                string
	AuthorOpenId                   string
	ShopName                       string
	TotalPayAmount                 float64
	CommissionRate                 float64
	FlowPoint                      string
	App                            string
	UpdateTime                     string
	PaySuccessTime                 string
	SettleTime                     string
	PayGoodsAmount                 int64
	SettledGoodsAmount             int64
	EstimatedCommission            int64
	RealCommission                 int64
	Extra                          string
	ItemNum                        int64
	ShopId                         int64
	RefundTime                     string
	PidInfo                        PidInfo
	EstimatedTotalCommission       int64
	EstimatedTechServiceFee        int64
	PickSourceClientKey            string
	PickExtra                      string
	AuthorShortId                  string
	MediaType                      string
	IsSteppedPlan                  bool
	PlatformSubsidy                int64
	AuthorSubsidy                  int64
	ProductActivityId              string
	AppId                          int64
	SettleUserSteppedCommission    int64
	SettleInstSteppedCommission    int64
	PaySubsidy                     int64
	MediaId                        int64
	AuthorBuyinId                  string
	ConfirmTime                    string
	EstimatedInstSteppedCommission int64
	EstimatedUserSteppedCommission int64
}

type JinritemaiOrderList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*JinritemaiOrderInfo
}

type CommissionRateJinritemaiOrder struct {
	ProductId uint64 `json:"productId"`
	VideoId   uint64 `json:"videoId"`
}

type StatisticsJinritemaiOrder struct {
	Key   string
	Value string
}

type StatisticsJinritemaiOrders struct {
	Statistics []*StatisticsJinritemaiOrder
}

type MessageJinritemaiOrder struct {
	ClientKey string `json:"client_key"`
	OpenId    string `json:"open_id"`
}

func NewJinritemaiOrder(ctx context.Context, payGoodsAmount, settledGoodsAmount, estimatedCommission, realCommission, itemNum, shopId, estimatedTotalCommission, estimatedTechServiceFee, platformSubsidy, authorSubsidy, appId, settleUserSteppedCommission, settleInstSteppedCommission, paySubsidy, mediaId, estimatedInstSteppedCommission, estimatedUserSteppedCommission int64, totalPayAmount, commissionRate float64, isSteppedPlan bool, orderId, productId, productName, productImg, authorAccount, authorClientKey, authorOpenId, shopName, flowPoint, app, updateTime, paySuccessTime, settleTime, extra, refundTime, pickSourceClientKey, pickExtra, authorShortId, mediaType, authorBuyinId, confirmTime, productActivityId string, pidInfo PidInfo) *JinritemaiOrder {
	return &JinritemaiOrder{
		OrderId:                        orderId,
		ProductId:                      productId,
		ProductName:                    productName,
		ProductImg:                     productImg,
		AuthorAccount:                  authorAccount,
		AuthorClientKey:                authorClientKey,
		AuthorOpenId:                   authorOpenId,
		ShopName:                       shopName,
		TotalPayAmount:                 totalPayAmount,
		CommissionRate:                 commissionRate,
		FlowPoint:                      flowPoint,
		App:                            app,
		UpdateTime:                     updateTime,
		PaySuccessTime:                 paySuccessTime,
		SettleTime:                     settleTime,
		PayGoodsAmount:                 payGoodsAmount,
		SettledGoodsAmount:             settledGoodsAmount,
		EstimatedCommission:            estimatedCommission,
		RealCommission:                 realCommission,
		Extra:                          extra,
		ItemNum:                        itemNum,
		ShopId:                         shopId,
		RefundTime:                     refundTime,
		PidInfo:                        pidInfo,
		EstimatedTotalCommission:       estimatedTotalCommission,
		EstimatedTechServiceFee:        estimatedTechServiceFee,
		PickSourceClientKey:            pickSourceClientKey,
		PickExtra:                      pickExtra,
		AuthorShortId:                  authorShortId,
		MediaType:                      mediaType,
		IsSteppedPlan:                  isSteppedPlan,
		PlatformSubsidy:                platformSubsidy,
		AuthorSubsidy:                  authorSubsidy,
		ProductActivityId:              productActivityId,
		AppId:                          appId,
		SettleUserSteppedCommission:    settleUserSteppedCommission,
		SettleInstSteppedCommission:    settleInstSteppedCommission,
		PaySubsidy:                     paySubsidy,
		MediaId:                        mediaId,
		AuthorBuyinId:                  authorBuyinId,
		ConfirmTime:                    confirmTime,
		EstimatedInstSteppedCommission: estimatedInstSteppedCommission,
		EstimatedUserSteppedCommission: estimatedUserSteppedCommission,
	}
}

func (jo *JinritemaiOrder) SetOrderId(ctx context.Context, orderId string) {
	jo.OrderId = orderId
}

func (jo *JinritemaiOrder) SetProductId(ctx context.Context, productId string) {
	jo.ProductId = productId
}

func (jo *JinritemaiOrder) SetProductName(ctx context.Context, productName string) {
	jo.ProductName = productName
}

func (jo *JinritemaiOrder) SetProductImg(ctx context.Context, productImg string) {
	jo.ProductImg = productImg
}

func (jo *JinritemaiOrder) SetAuthorAccount(ctx context.Context, authorAccount string) {
	jo.AuthorAccount = authorAccount
}

func (jo *JinritemaiOrder) SetAuthorClientKey(ctx context.Context, authorClientKey string) {
	jo.AuthorClientKey = authorClientKey
}

func (jo *JinritemaiOrder) SetAuthorOpenId(ctx context.Context, authorOpenId string) {
	jo.AuthorOpenId = authorOpenId
}

func (jo *JinritemaiOrder) SetShopName(ctx context.Context, shopName string) {
	jo.ShopName = shopName
}

func (jo *JinritemaiOrder) SetTotalPayAmount(ctx context.Context, totalPayAmount float64) {
	jo.TotalPayAmount = totalPayAmount
}

func (jo *JinritemaiOrder) SetCommissionRate(ctx context.Context, commissionRate float64) {
	jo.CommissionRate = commissionRate
}

func (jo *JinritemaiOrder) SetFlowPoint(ctx context.Context, flowPoint string) {
	jo.FlowPoint = flowPoint
}

func (jo *JinritemaiOrder) SetApp(ctx context.Context, app string) {
	jo.App = app
}

func (jo *JinritemaiOrder) SetUpdateTime(ctx context.Context, updateTime string) {
	jo.UpdateTime = updateTime
}

func (jo *JinritemaiOrder) SetPaySuccessTime(ctx context.Context, paySuccessTime string) {
	jo.PaySuccessTime = paySuccessTime
}

func (jo *JinritemaiOrder) SetSettleTime(ctx context.Context, settleTime string) {
	jo.SettleTime = settleTime
}

func (jo *JinritemaiOrder) SetPayGoodsAmount(ctx context.Context, payGoodsAmount int64) {
	jo.PayGoodsAmount = payGoodsAmount
}

func (jo *JinritemaiOrder) SetSettledGoodsAmount(ctx context.Context, settledGoodsAmount int64) {
	jo.SettledGoodsAmount = settledGoodsAmount
}

func (jo *JinritemaiOrder) SetEstimatedCommission(ctx context.Context, estimatedCommission int64) {
	jo.EstimatedCommission = estimatedCommission
}

func (jo *JinritemaiOrder) SetRealCommission(ctx context.Context, realCommission int64) {
	jo.RealCommission = realCommission
}

func (jo *JinritemaiOrder) SetExtra(ctx context.Context, extra string) {
	jo.Extra = extra
}

func (jo *JinritemaiOrder) SetItemNum(ctx context.Context, itemNum int64) {
	jo.ItemNum = itemNum
}

func (jo *JinritemaiOrder) SetShopId(ctx context.Context, shopId int64) {
	jo.ShopId = shopId
}

func (jo *JinritemaiOrder) SetRefundTime(ctx context.Context, refundTime string) {
	jo.RefundTime = refundTime
}

func (jo *JinritemaiOrder) SetPidInfo(ctx context.Context, pidInfo PidInfo) {
	jo.PidInfo = pidInfo
}

func (jo *JinritemaiOrder) SetEstimatedTotalCommission(ctx context.Context, estimatedTotalCommission int64) {
	jo.EstimatedTotalCommission = estimatedTotalCommission
}

func (jo *JinritemaiOrder) SetEstimatedTechServiceFee(ctx context.Context, estimatedTechServiceFee int64) {
	jo.EstimatedTechServiceFee = estimatedTechServiceFee
}

func (jo *JinritemaiOrder) SetPickSourceClientKey(ctx context.Context, pickSourceClientKey string) {
	jo.PickSourceClientKey = pickSourceClientKey
}

func (jo *JinritemaiOrder) SetPickExtra(ctx context.Context, pickExtra string) {
	jo.PickExtra = pickExtra
}

func (jo *JinritemaiOrder) SetAuthorShortId(ctx context.Context, authorShortId string) {
	jo.AuthorShortId = authorShortId
}

func (jo *JinritemaiOrder) SetMediaType(ctx context.Context, mediaType string) {
	jo.MediaType = mediaType
}

func (jo *JinritemaiOrder) SetIsSteppedPlan(ctx context.Context, isSteppedPlan bool) {
	jo.IsSteppedPlan = isSteppedPlan
}

func (jo *JinritemaiOrder) SetPlatformSubsidy(ctx context.Context, platformSubsidy int64) {
	jo.PlatformSubsidy = platformSubsidy
}

func (jo *JinritemaiOrder) SetAuthorSubsidy(ctx context.Context, authorSubsidy int64) {
	jo.AuthorSubsidy = authorSubsidy
}

func (jo *JinritemaiOrder) SetProductActivityId(ctx context.Context, productActivityId string) {
	jo.ProductActivityId = productActivityId
}

func (jo *JinritemaiOrder) SetAppId(ctx context.Context, appId int64) {
	jo.AppId = appId
}

func (jo *JinritemaiOrder) SetSettleUserSteppedCommission(ctx context.Context, settleUserSteppedCommission int64) {
	jo.SettleUserSteppedCommission = settleUserSteppedCommission
}

func (jo *JinritemaiOrder) SetPaySubsidy(ctx context.Context, paySubsidy int64) {
	jo.PaySubsidy = paySubsidy
}

func (jo *JinritemaiOrder) SetMediaId(ctx context.Context, mediaId int64) {
	jo.MediaId = mediaId
}

func (jo *JinritemaiOrder) SetAuthorBuyinId(ctx context.Context, authorBuyinId string) {
	jo.AuthorBuyinId = authorBuyinId
}

func (jo *JinritemaiOrder) SetConfirmTime(ctx context.Context, confirmTime string) {
	jo.ConfirmTime = confirmTime
}

func (jo *JinritemaiOrder) SetEstimatedInstSteppedCommission(ctx context.Context, estimatedInstSteppedCommission int64) {
	jo.EstimatedInstSteppedCommission = estimatedInstSteppedCommission
}

func (jo *JinritemaiOrder) SetEstimatedUserSteppedCommission(ctx context.Context, estimatedUserSteppedCommission int64) {
	jo.EstimatedUserSteppedCommission = estimatedUserSteppedCommission
}

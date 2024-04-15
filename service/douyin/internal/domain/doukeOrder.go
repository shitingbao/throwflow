package domain

import (
	"context"
)

type StatisticsDoukeOrder struct {
	Key   string
	Value string
}

type StatisticsDoukeOrders struct {
	Statistics []*StatisticsDoukeOrder
}

type DoukeOrder struct {
	OrderId                 string  `json:"order_id"`
	AppId                   string  `json:"app_id"`
	ProductId               string  `json:"product_id"`
	ProductName             string  `json:"product_name"`
	AuthorAccount           string  `json:"author_account"`
	AdsAttribution          string  `json:"ads_attribution"`
	ProductImg              string  `json:"product_img"`
	TotalPayAmount          int     `json:"total_pay_amount"`
	PaySuccessTime          string  `json:"pay_success_time"`
	RefundTime              string  `json:"refund_time"`
	PayGoodsAmount          int     `json:"pay_goods_amount"`
	EstimatedCommission     float32 `json:"estimated_commission"`
	AdsRealCommission       float32 `json:"ads_real_commission"`
	SplitRate               float32 `json:"split_rate"`
	AfterSalesStatus        int     `json:"after_sales_status"`
	FlowPoint               string  `json:"flow_point"`
	ExternalInfo            string  `json:"external_info"`
	SettleTime              string  `json:"settle_time"`
	ConfirmTime             string  `json:"confirm_time"`
	MediaTypeName           string  `json:"media_type_name"`
	UpdateTime              string  `json:"update_time"`
	EstimatedTechServiceFee int     `json:"estimated_tech_service_fee"`
}

func NewDoukeOrder(ctx context.Context, totalPayAmount, payGoodsAmount, afterSalesStatus, estimatedTechServiceFee int, estimatedCommission, adsRealCommission, splitRate float32, orderId, appId, productId, productName, authorAccount, adsAttribution, productImg, paySuccessTime, refundTime, flowPoint, externalInfo, settleTime, confirmTime, mediaTypeName, updateTime string) *DoukeOrder {
	return &DoukeOrder{
		OrderId:                 orderId,
		AppId:                   appId,
		ProductId:               productId,
		ProductName:             productName,
		AuthorAccount:           authorAccount,
		AdsAttribution:          adsAttribution,
		ProductImg:              productImg,
		TotalPayAmount:          totalPayAmount,
		PaySuccessTime:          paySuccessTime,
		RefundTime:              refundTime,
		PayGoodsAmount:          payGoodsAmount,
		EstimatedCommission:     estimatedCommission,
		AdsRealCommission:       adsRealCommission,
		SplitRate:               splitRate,
		AfterSalesStatus:        afterSalesStatus,
		FlowPoint:               flowPoint,
		ExternalInfo:            externalInfo,
		SettleTime:              settleTime,
		ConfirmTime:             confirmTime,
		MediaTypeName:           mediaTypeName,
		UpdateTime:              updateTime,
		EstimatedTechServiceFee: estimatedTechServiceFee,
	}
}

func (do *DoukeOrder) SetOrderId(ctx context.Context, orderId string) {
	do.OrderId = orderId
}

func (do *DoukeOrder) SetAppId(ctx context.Context, appId string) {
	do.AppId = appId
}

func (do *DoukeOrder) SetProductId(ctx context.Context, productId string) {
	do.ProductId = productId
}

func (do *DoukeOrder) SetProductName(ctx context.Context, productName string) {
	do.ProductName = productName
}

func (do *DoukeOrder) SetAuthorAccount(ctx context.Context, authorAccount string) {
	do.AuthorAccount = authorAccount
}

func (do *DoukeOrder) SetAdsAttribution(ctx context.Context, adsAttribution string) {
	do.AdsAttribution = adsAttribution
}

func (do *DoukeOrder) SetProductImg(ctx context.Context, productImg string) {
	do.ProductImg = productImg
}

func (do *DoukeOrder) SetTotalPayAmount(ctx context.Context, totalPayAmount int) {
	do.TotalPayAmount = totalPayAmount
}

func (do *DoukeOrder) SetPaySuccessTime(ctx context.Context, paySuccessTime string) {
	do.PaySuccessTime = paySuccessTime
}

func (do *DoukeOrder) SetRefundTime(ctx context.Context, refundTime string) {
	do.RefundTime = refundTime
}

func (do *DoukeOrder) SetPayGoodsAmount(ctx context.Context, payGoodsAmount int) {
	do.PayGoodsAmount = payGoodsAmount
}

func (do *DoukeOrder) SetEstimatedCommission(ctx context.Context, estimatedCommission float32) {
	do.EstimatedCommission = estimatedCommission
}

func (do *DoukeOrder) SetAdsRealCommission(ctx context.Context, adsRealCommission float32) {
	do.AdsRealCommission = adsRealCommission
}

func (do *DoukeOrder) SetSplitRate(ctx context.Context, splitRate float32) {
	do.SplitRate = splitRate
}

func (do *DoukeOrder) SetAfterSalesStatus(ctx context.Context, afterSalesStatus int) {
	do.AfterSalesStatus = afterSalesStatus
}

func (do *DoukeOrder) SetFlowPoint(ctx context.Context, flowPoint string) {
	do.FlowPoint = flowPoint
}

func (do *DoukeOrder) SetExternalInfo(ctx context.Context, externalInfo string) {
	do.ExternalInfo = externalInfo
}

func (do *DoukeOrder) SetSettleTime(ctx context.Context, settleTime string) {
	do.SettleTime = settleTime
}

func (do *DoukeOrder) SetConfirmTime(ctx context.Context, confirmTime string) {
	do.ConfirmTime = confirmTime
}

func (do *DoukeOrder) SetMediaTypeName(ctx context.Context, mediaTypeName string) {
	do.MediaTypeName = mediaTypeName
}

func (do *DoukeOrder) SetUpdateTime(ctx context.Context, updateTime string) {
	do.UpdateTime = updateTime
}

func (do *DoukeOrder) SetEstimatedTechServiceFee(ctx context.Context, estimatedTechServiceFee int) {
	do.EstimatedTechServiceFee = estimatedTechServiceFee
}

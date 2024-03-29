package domain

import (
	"context"
)

type ProductTags struct {
	HasSubsidyTag     bool `json:"has_subsidy_tag" bson:"has_subsidy_tag"`
	HasSupermarketTag bool `json:"has_supermarket_tag" bson:"has_supermarket_tag"`
}

type StatisticsDoukeOrder struct {
	Key   string
	Value string
}

type StatisticsDoukeOrders struct {
	Statistics []*StatisticsDoukeOrder
}

type DoukeOrder struct {
	AuthorAccount          string
	MediaType              string
	AdsActivityId          int64
	ProductImg             string
	UpdateTime             string
	PaySuccessTime         string
	AdsRealCommission      int64
	AdsEstimatedCommission int64
	ProductId              string
	TotalPayAmount         int64
	FlowPoint              string
	SettleTime             string
	SettledGoodsAmount     int64
	PidInfo                PidInfo
	ItemNum                int64
	AuthorBuyinId          string
	ShopId                 int64
	PayGoodsAmount         int64
	AdsDistributorId       int64
	AdsPromotionTate       int64
	OrderId                string
	ProductName            string
	DistributionType       string
	AuthorUid              int64
	ShopName               string
	ProductActivityId      string
	MaterialId             string
	RefundTime             string
	ConfirmTime            string
	ProductTags            ProductTags
	BuyerAppId             string
	DistributorRightType   string
}

func NewDoukeOrder(ctx context.Context, adsActivityId, adsRealCommission, adsEstimatedCommission, totalPayAmount, settledGoodsAmount, itemNum, shopId, payGoodsAmount, adsDistributorId, adsPromotionTate, authorUid int64, authorAccount, mediaType, productImg, updateTime, paySuccessTime, productId, flowPoint, settleTime, authorBuyinId, orderId, productName, distributionType, shopName, productActivityId, materialId, refundTime, confirmTime, buyerAppId, distributorRightType string, pidInfo PidInfo, productTags ProductTags) *DoukeOrder {
	return &DoukeOrder{
		AuthorAccount:          authorAccount,
		MediaType:              mediaType,
		AdsActivityId:          adsActivityId,
		ProductImg:             productImg,
		UpdateTime:             updateTime,
		PaySuccessTime:         paySuccessTime,
		AdsRealCommission:      adsRealCommission,
		AdsEstimatedCommission: adsEstimatedCommission,
		ProductId:              productId,
		TotalPayAmount:         totalPayAmount,
		FlowPoint:              flowPoint,
		SettleTime:             settleTime,
		SettledGoodsAmount:     settledGoodsAmount,
		PidInfo:                pidInfo,
		ItemNum:                itemNum,
		AuthorBuyinId:          authorBuyinId,
		ShopId:                 shopId,
		PayGoodsAmount:         payGoodsAmount,
		AdsDistributorId:       adsDistributorId,
		AdsPromotionTate:       adsPromotionTate,
		OrderId:                orderId,
		ProductName:            productName,
		DistributionType:       distributionType,
		AuthorUid:              authorUid,
		ShopName:               shopName,
		ProductActivityId:      productActivityId,
		MaterialId:             materialId,
		RefundTime:             refundTime,
		ConfirmTime:            confirmTime,
		ProductTags:            productTags,
		BuyerAppId:             buyerAppId,
		DistributorRightType:   distributorRightType,
	}
}

func (do *DoukeOrder) SetAuthorAccount(ctx context.Context, authorAccount string) {
	do.AuthorAccount = authorAccount
}

func (do *DoukeOrder) SetMediaType(ctx context.Context, mediaType string) {
	do.MediaType = mediaType
}

func (do *DoukeOrder) SetAdsActivityId(ctx context.Context, adsActivityId int64) {
	do.AdsActivityId = adsActivityId
}

func (do *DoukeOrder) SetProductImg(ctx context.Context, productImg string) {
	do.ProductImg = productImg
}

func (do *DoukeOrder) SetUpdateTime(ctx context.Context, updateTime string) {
	do.UpdateTime = updateTime
}

func (do *DoukeOrder) SetPaySuccessTime(ctx context.Context, paySuccessTime string) {
	do.PaySuccessTime = paySuccessTime
}

func (do *DoukeOrder) SetAdsRealCommission(ctx context.Context, adsRealCommission int64) {
	do.AdsRealCommission = adsRealCommission
}

func (do *DoukeOrder) SetAdsEstimatedCommission(ctx context.Context, adsEstimatedCommission int64) {
	do.AdsEstimatedCommission = adsEstimatedCommission
}

func (do *DoukeOrder) SetProductId(ctx context.Context, productId string) {
	do.ProductId = productId
}

func (do *DoukeOrder) SetTotalPayAmount(ctx context.Context, totalPayAmount int64) {
	do.TotalPayAmount = totalPayAmount
}

func (do *DoukeOrder) SetFlowPoint(ctx context.Context, flowPoint string) {
	do.FlowPoint = flowPoint
}

func (do *DoukeOrder) SetSettleTime(ctx context.Context, settleTime string) {
	do.SettleTime = settleTime
}

func (do *DoukeOrder) SetSettledGoodsAmount(ctx context.Context, settledGoodsAmount int64) {
	do.SettledGoodsAmount = settledGoodsAmount
}

func (do *DoukeOrder) SetPidInfo(ctx context.Context, pidInfo PidInfo) {
	do.PidInfo = pidInfo
}

func (do *DoukeOrder) SetItemNum(ctx context.Context, itemNum int64) {
	do.ItemNum = itemNum
}

func (do *DoukeOrder) SetAuthorBuyinId(ctx context.Context, authorBuyinId string) {
	do.AuthorBuyinId = authorBuyinId
}

func (do *DoukeOrder) SetShopId(ctx context.Context, shopId int64) {
	do.ShopId = shopId
}

func (do *DoukeOrder) SetPayGoodsAmount(ctx context.Context, payGoodsAmount int64) {
	do.PayGoodsAmount = payGoodsAmount
}

func (do *DoukeOrder) SetAdsDistributorId(ctx context.Context, adsDistributorId int64) {
	do.AdsDistributorId = adsDistributorId
}

func (do *DoukeOrder) SetAdsPromotionTate(ctx context.Context, adsPromotionTate int64) {
	do.AdsPromotionTate = adsPromotionTate
}

func (do *DoukeOrder) SetOrderId(ctx context.Context, orderId string) {
	do.OrderId = orderId
}

func (do *DoukeOrder) SetProductName(ctx context.Context, productName string) {
	do.ProductName = productName
}

func (do *DoukeOrder) SetDistributionType(ctx context.Context, distributionType string) {
	do.DistributionType = distributionType
}

func (do *DoukeOrder) SetAuthorUid(ctx context.Context, authorUid int64) {
	do.AuthorUid = authorUid
}

func (do *DoukeOrder) SetShopName(ctx context.Context, shopName string) {
	do.ShopName = shopName
}

func (do *DoukeOrder) SetProductActivityId(ctx context.Context, productActivityId string) {
	do.ProductActivityId = productActivityId
}

func (do *DoukeOrder) SetMaterialId(ctx context.Context, materialId string) {
	do.MaterialId = materialId
}

func (do *DoukeOrder) SetRefundTime(ctx context.Context, refundTime string) {
	do.RefundTime = refundTime
}

func (do *DoukeOrder) SetConfirmTime(ctx context.Context, confirmTime string) {
	do.ConfirmTime = confirmTime
}

func (do *DoukeOrder) SetProductTags(ctx context.Context, productTags ProductTags) {
	do.ProductTags = productTags
}

func (do *DoukeOrder) SetBuyerAppId(ctx context.Context, buyerAppId string) {
	do.BuyerAppId = buyerAppId
}

func (do *DoukeOrder) SetDistributorRightType(ctx context.Context, distributorRightType string) {
	do.DistributorRightType = distributorRightType
}

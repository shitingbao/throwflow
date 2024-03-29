package domain

import (
	"context"
	"douyin/internal/pkg/tool"
	"strconv"
	"time"
)

type JinritemaiOrderInfo struct {
	Id                  uint64
	ClientKey           string
	OpenId              string
	BuyinId             string
	OrderId             string
	ProductId           string
	ProductName         string
	ProductImg          string
	CommissionRate      uint8
	PaySuccessTime      time.Time
	SettleTime          *time.Time
	TotalPayAmount      float32
	PayGoodsAmount      float32
	FlowPoint           string
	EstimatedCommission float32
	RealCommission      float32
	RealCommissionRate  string
	ItemNum             uint64
	PickExtra           string
	MediaType           string
	MediaTypeName       string
	MediaId             uint64
	MediaCover          string
	Avatar              string
	IsShow              uint8
	CreateTime          time.Time
	UpdateTime          time.Time
}

type JinritemaiOrderInfoList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*JinritemaiOrderInfo
}

type JinritemaiOrderInfoStatisticsAwemeIndustry struct {
	ClientKey    string
	OpenId       string
	IndustryId   uint64
	IndustryName string
	ItemNum      uint64
}

type JinritemaiOrderInfoStatisticsAwemeIndustryList struct {
	ClientKey string
	OpenId    string
	ItemNum   uint64
	List      []*JinritemaiOrderInfoStatisticsAwemeIndustry
}

type JinritemaiOrderInfoReportArea struct {
	Area string
}

type StorePreference struct {
	IndustryId    uint64
	IndustryName  string
	IndustryRatio string
	ItemNum       uint64
}

func NewJinritemaiOrderInfo(ctx context.Context, itemNum, mediaId uint64, commissionRate uint8, totalPayAmount, payGoodsAmount, estimatedCommission, realCommission float32, clientKey, openId, buyinId, orderId, productId, productName, productImg, flowPoint, pickExtra, mediaType string, paySuccessTime time.Time, settleTime *time.Time) *JinritemaiOrderInfo {
	return &JinritemaiOrderInfo{
		ClientKey:           clientKey,
		OpenId:              openId,
		BuyinId:             buyinId,
		OrderId:             orderId,
		ProductId:           productId,
		ProductName:         productName,
		ProductImg:          productImg,
		CommissionRate:      commissionRate,
		PaySuccessTime:      paySuccessTime,
		SettleTime:          settleTime,
		TotalPayAmount:      totalPayAmount,
		PayGoodsAmount:      payGoodsAmount,
		FlowPoint:           flowPoint,
		EstimatedCommission: estimatedCommission,
		RealCommission:      realCommission,
		ItemNum:             itemNum,
		PickExtra:           pickExtra,
		MediaType:           mediaType,
		MediaId:             mediaId,
	}
}

func (joi *JinritemaiOrderInfo) SetClientKey(ctx context.Context, clientKey string) {
	joi.ClientKey = clientKey
}

func (joi *JinritemaiOrderInfo) SetOpenId(ctx context.Context, openId string) {
	joi.OpenId = openId
}

func (joi *JinritemaiOrderInfo) SetBuyinId(ctx context.Context, buyinId string) {
	joi.BuyinId = buyinId
}

func (joi *JinritemaiOrderInfo) SetOrderId(ctx context.Context, orderId string) {
	joi.OrderId = orderId
}

func (joi *JinritemaiOrderInfo) SetProductId(ctx context.Context, productId string) {
	joi.ProductId = productId
}

func (joi *JinritemaiOrderInfo) SetProductName(ctx context.Context, productName string) {
	joi.ProductName = productName
}

func (joi *JinritemaiOrderInfo) SetProductImg(ctx context.Context, productImg string) {
	joi.ProductImg = productImg
}

func (joi *JinritemaiOrderInfo) SetCommissionRate(ctx context.Context, commissionRate uint8) {
	joi.CommissionRate = commissionRate
}

func (joi *JinritemaiOrderInfo) SetPaySuccessTime(ctx context.Context, paySuccessTime time.Time) {
	joi.PaySuccessTime = paySuccessTime
}

func (joi *JinritemaiOrderInfo) SetSettleTime(ctx context.Context, settleTime *time.Time) {
	joi.SettleTime = settleTime
}

func (joi *JinritemaiOrderInfo) SetTotalPayAmount(ctx context.Context, totalPayAmount float32) {
	joi.TotalPayAmount = totalPayAmount
}

func (joi *JinritemaiOrderInfo) SetPayGoodsAmount(ctx context.Context, payGoodsAmount float32) {
	joi.PayGoodsAmount = payGoodsAmount
}

func (joi *JinritemaiOrderInfo) SetFlowPoint(ctx context.Context, flowPoint string) {
	joi.FlowPoint = flowPoint
}

func (joi *JinritemaiOrderInfo) SetEstimatedCommission(ctx context.Context, estimatedCommission float32) {
	joi.EstimatedCommission = estimatedCommission
}

func (joi *JinritemaiOrderInfo) SetRealCommission(ctx context.Context, realCommission float32) {
	joi.RealCommission = realCommission
}

func (joi *JinritemaiOrderInfo) SetItemNum(ctx context.Context, itemNum uint64) {
	joi.ItemNum = itemNum
}

func (joi *JinritemaiOrderInfo) SetPickExtra(ctx context.Context, pickExtra string) {
	joi.PickExtra = pickExtra
}

func (joi *JinritemaiOrderInfo) SetMediaType(ctx context.Context, mediaType string) {
	joi.MediaType = mediaType
}

func (joi *JinritemaiOrderInfo) SetMediaTypeName(ctx context.Context) {
	if joi.MediaType == "shop_list" {
		joi.MediaTypeName = "橱窗"
	} else if joi.MediaType == "video" {
		joi.MediaTypeName = "视频"
	} else if joi.MediaType == "live" {
		joi.MediaTypeName = "直播"
	} else if joi.MediaType == "others" {
		joi.MediaTypeName = "其他"
	}
}

func (joi *JinritemaiOrderInfo) SetMediaId(ctx context.Context, mediaId uint64) {
	joi.MediaId = mediaId
}

func (joi *JinritemaiOrderInfo) SetCreateTime(ctx context.Context) {
	joi.CreateTime = time.Now()
}

func (joi *JinritemaiOrderInfo) SetUpdateTime(ctx context.Context) {
	joi.UpdateTime = time.Now()
}

func (joi *JinritemaiOrderInfo) GetItemNum(ctx context.Context) string {
	return strconv.FormatUint(joi.ItemNum, 10)
}

func (joi *JinritemaiOrderInfo) GetTotalPayAmount(ctx context.Context) string {
	return strconv.FormatFloat(tool.Decimal(float64(joi.TotalPayAmount), 2), 'f', 2, 64)
}

func (joi *JinritemaiOrderInfo) GetPayGoodsAmount(ctx context.Context) string {
	return strconv.FormatFloat(tool.Decimal(float64(joi.PayGoodsAmount), 2), 'f', 2, 64)
}

func (joi *JinritemaiOrderInfo) GetEstimatedCommission(ctx context.Context) string {
	return strconv.FormatFloat(tool.Decimal(float64(joi.EstimatedCommission), 2), 'f', 2, 64)
}

func (joi *JinritemaiOrderInfo) GetRealCommission(ctx context.Context) string {
	return strconv.FormatFloat(tool.Decimal(float64(joi.RealCommission), 2), 'f', 2, 64)
}

func (joi *JinritemaiOrderInfo) GetRealCommissionRate(ctx context.Context) string {
	var realCommissionRate float64

	if joi.TotalPayAmount > 0 {
		realCommissionRate = float64(joi.EstimatedCommission) / float64(joi.PayGoodsAmount)
	}

	return strconv.FormatFloat(tool.Decimal(realCommissionRate*100, 2), 'f', 2, 64) + "%"
}

func (joi *JinritemaiOrderInfo) GetRefundRate(ctx context.Context, refundItemNum uint64) (refundRate float32) {
	if joi.ItemNum > 0 {
		refundRate = float32(refundItemNum) / float32(joi.ItemNum)
	}

	return
}

func (joisai *JinritemaiOrderInfoStatisticsAwemeIndustry) GetIndustryRatio(ctx context.Context, itemNum uint64) string {
	var industryRatio float64

	if itemNum > 0 {
		industryRatio = float64(joisai.ItemNum) / float64(itemNum)
	}

	return strconv.FormatFloat(tool.Decimal(industryRatio*100, 2), 'f', 2, 64) + "%"
}

func (sp *StorePreference) SetIndustryRatio(ctx context.Context, itemNum uint64) {
	var industryRatio float64

	if itemNum > 0 {
		industryRatio = float64(sp.ItemNum) / float64(itemNum)
	}

	sp.IndustryRatio = strconv.FormatFloat(tool.Decimal(industryRatio*100, 2), 'f', 2, 64)
}

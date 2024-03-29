package domain

import (
	"context"
	"douyin/internal/pkg/tool"
	"time"
)

type QianchuanReportProduct struct {
	AdvertiserId            uint64
	AdvertiserName          string
	AdIds                   []uint64
	ProductId               uint64
	DiscountPrice           float64
	ProductName             string
	ProductImg              string
	StatCost                float64
	PayOrderCount           int64
	PayOrderAmount          float64
	ShowCnt                 int64
	ClickCnt                int64
	ConvertCnt              int64
	DyFollow                int64
	Roi                     float64
	PayOrderAveragePrice    float64
	ClickRate               float64
	ConvertRate             float64
	AveragePayOrderStatCost float64
	CreateTime              time.Time
	UpdateTime              time.Time
}

type QianchuanReportProductList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*QianchuanReportProduct
}

type StatisticsQianchuanReportProduct struct {
	Key   string
	Value string
}

type StatisticsQianchuanReportProducts struct {
	Statistics []*StatisticsQianchuanReportProduct
}

func NewQianchuanReportProduct(ctx context.Context, advertiserId, productId uint64, payOrderCount, showCnt, clickCnt, convertCnt, dyFollow int64, discountPrice, statCost, payOrderAmount float64, advertiserName, productName, productImg string) *QianchuanReportProduct {
	return &QianchuanReportProduct{
		AdvertiserId:   advertiserId,
		AdvertiserName: advertiserName,
		ProductId:      productId,
		DiscountPrice:  discountPrice,
		ProductName:    productName,
		ProductImg:     productImg,
		StatCost:       statCost,
		PayOrderCount:  payOrderCount,
		PayOrderAmount: payOrderAmount,
		ShowCnt:        showCnt,
		ClickCnt:       clickCnt,
		ConvertCnt:     convertCnt,
		DyFollow:       dyFollow,
	}
}

func (qrp *QianchuanReportProduct) SetAdvertiserId(ctx context.Context, advertiserId uint64) {
	qrp.AdvertiserId = advertiserId
}

func (qrp *QianchuanReportProduct) SetAdvertiserName(ctx context.Context, advertiserName string) {
	qrp.AdvertiserName = advertiserName
}

func (qrp *QianchuanReportProduct) SetProductId(ctx context.Context, productId uint64) {
	qrp.ProductId = productId
}

func (qrp *QianchuanReportProduct) SetDiscountPrice(ctx context.Context, discountPrice float64) {
	qrp.DiscountPrice = discountPrice
}

func (qrp *QianchuanReportProduct) SetProductName(ctx context.Context, productName string) {
	qrp.ProductName = productName
}

func (qrp *QianchuanReportProduct) SetProductImg(ctx context.Context, productImg string) {
	qrp.ProductImg = productImg
}

func (qrp *QianchuanReportProduct) SetStatCost(ctx context.Context, statCost float64) {
	qrp.StatCost = statCost
}

func (qrp *QianchuanReportProduct) SetPayOrderCount(ctx context.Context, payOrderCount int64) {
	qrp.PayOrderCount = payOrderCount
}

func (qrp *QianchuanReportProduct) SetPayOrderAmount(ctx context.Context, payOrderAmount float64) {
	qrp.PayOrderAmount = payOrderAmount
}

func (qrp *QianchuanReportProduct) SetShowCnt(ctx context.Context, showCnt int64) {
	qrp.ShowCnt = showCnt
}

func (qrp *QianchuanReportProduct) SetClickCnt(ctx context.Context, clickCnt int64) {
	qrp.ClickCnt = clickCnt
}

func (qrp *QianchuanReportProduct) SetConvertCnt(ctx context.Context, convertCnt int64) {
	qrp.ConvertCnt = convertCnt
}

func (qrp *QianchuanReportProduct) SetDyFollow(ctx context.Context, dyFollow int64) {
	qrp.DyFollow = dyFollow
}

func (qrp *QianchuanReportProduct) SetRoi(ctx context.Context) {
	if qrp.StatCost > 0 {
		qrp.Roi = tool.Decimal(qrp.PayOrderAmount/qrp.StatCost, 2)
	}
}

func (qrp *QianchuanReportProduct) SetPayOrderAveragePrice(ctx context.Context) {
	if qrp.PayOrderCount > 0 {
		qrp.PayOrderAveragePrice = tool.Decimal(qrp.PayOrderAmount/float64(qrp.PayOrderCount), 2)
	}

}

func (qrp *QianchuanReportProduct) SetClickRate(ctx context.Context) {
	if qrp.ShowCnt > 0 {
		qrp.ClickRate = tool.Decimal(float64(qrp.ClickCnt)/float64(qrp.ShowCnt), 2)
	}

}

func (qrp *QianchuanReportProduct) SetConvertRate(ctx context.Context) {
	if qrp.ClickCnt > 0 {
		qrp.ConvertRate = tool.Decimal(float64(qrp.ConvertCnt)/float64(qrp.ClickCnt), 2)
	}
}

func (qrp *QianchuanReportProduct) SetAveragePayOrderStatCost(ctx context.Context) {
	if qrp.PayOrderCount > 0 {
		qrp.AveragePayOrderStatCost = tool.Decimal(qrp.StatCost/float64(qrp.PayOrderCount), 2)
	}
}

func (qrp *QianchuanReportProduct) SetUpdateTime(ctx context.Context) {
	qrp.UpdateTime = time.Now()
}

func (qrp *QianchuanReportProduct) SetCreateTime(ctx context.Context) {
	qrp.CreateTime = time.Now()
}

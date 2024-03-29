package domain

import (
	"context"
	"time"
)

type QianchuanReportAd struct {
	AdId              uint64
	AdvertiserId      uint64
	AwemeId           uint64
	MarketingGoal     int64
	StatCost          float64
	ShowCnt           int64
	ClickCnt          int64
	PayOrderCount     int64
	CreateOrderAmount float64
	CreateOrderCount  int64
	PayOrderAmount    float64
	DyFollow          int64
	ConvertCnt        int64
	CreateTime        time.Time
	UpdateTime        time.Time
}

func NewQianchuanReportAd(ctx context.Context, adId, advertiserId, awemeId uint64, showCnt, clickCnt, payOrderCount, createOrderCount, dyFollow, convertCnt, marketingGoal int64, statCost, createOrderAmount, payOrderAmount float64) *QianchuanReportAd {
	return &QianchuanReportAd{
		AdId:              adId,
		AdvertiserId:      advertiserId,
		AwemeId:           awemeId,
		MarketingGoal:     marketingGoal,
		StatCost:          statCost,
		ShowCnt:           showCnt,
		ClickCnt:          clickCnt,
		PayOrderCount:     payOrderCount,
		CreateOrderAmount: createOrderAmount,
		CreateOrderCount:  createOrderCount,
		PayOrderAmount:    payOrderAmount,
		DyFollow:          dyFollow,
		ConvertCnt:        convertCnt,
	}
}

func (qra *QianchuanReportAd) SetAdId(ctx context.Context, adId uint64) {
	qra.AdId = adId
}

func (qra *QianchuanReportAd) SetAdvertiserId(ctx context.Context, advertiserId uint64) {
	qra.AdvertiserId = advertiserId
}

func (qra *QianchuanReportAd) SetAwemeId(ctx context.Context, awemeId uint64) {
	qra.AwemeId = awemeId
}

func (qra *QianchuanReportAd) SetMarketingGoal(ctx context.Context, marketingGoal int64) {
	qra.MarketingGoal = marketingGoal
}

func (qra *QianchuanReportAd) SetStatCost(ctx context.Context, statCost float64) {
	qra.StatCost = statCost
}

func (qra *QianchuanReportAd) SetShowCnt(ctx context.Context, showCnt int64) {
	qra.ShowCnt = showCnt
}

func (qra *QianchuanReportAd) SetClickCnt(ctx context.Context, clickCnt int64) {
	qra.ClickCnt = clickCnt
}

func (qra *QianchuanReportAd) SetPayOrderCount(ctx context.Context, payOrderCount int64) {
	qra.PayOrderCount = payOrderCount
}

func (qra *QianchuanReportAd) SetCreateOrderAmount(ctx context.Context, createOrderAmount float64) {
	qra.CreateOrderAmount = createOrderAmount
}

func (qra *QianchuanReportAd) SetCreateOrderCount(ctx context.Context, createOrderCount int64) {
	qra.CreateOrderCount = createOrderCount
}

func (qra *QianchuanReportAd) SetPayOrderAmount(ctx context.Context, payOrderAmount float64) {
	qra.PayOrderAmount = payOrderAmount
}

func (qra *QianchuanReportAd) SetDyFollow(ctx context.Context, dyFollow int64) {
	qra.DyFollow = dyFollow
}

func (qra *QianchuanReportAd) SetConvertCnt(ctx context.Context, convertCnt int64) {
	qra.ConvertCnt = convertCnt
}

func (qra *QianchuanReportAd) SetUpdateTime(ctx context.Context) {
	qra.UpdateTime = time.Now()
}

func (qra *QianchuanReportAd) SetCreateTime(ctx context.Context) {
	qra.CreateTime = time.Now()
}

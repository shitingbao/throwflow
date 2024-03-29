package domain

import (
	"context"
	"time"
)

type QianchuanReportAdRealtime struct {
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
	Time              int64
	CreateTime        time.Time
	UpdateTime        time.Time
}

type ExternalQianchuanReportAdRealtime struct {
	AdId                    uint64
	AdvertiserId            uint64
	AwemeId                 uint64
	MarketingGoal           int64
	StatCost                float64
	Roi                     float64
	ShowCnt                 int64
	ClickCnt                int64
	PayOrderCount           int64
	CreateOrderAmount       float64
	CreateOrderCount        int64
	PayOrderAmount          float64
	DyFollow                int64
	ConvertCnt              int64
	ClickRate               float64
	CpmPlatform             float64
	PayConvertRate          float64
	ConvertCost             float64
	ConvertRate             float64
	AveragePayOrderStatCost float64
	PayOrderAveragePrice    float64
	Time                    int64
}

func NewQianchuanReportAdRealtime(ctx context.Context, adId, advertiserId, awemeId uint64, marketingGoal, showCnt, clickCnt, payOrderCount, createOrderCount, dyFollow, convertCnt, time int64, statCost, createOrderAmount, payOrderAmount float64) *QianchuanReportAdRealtime {
	return &QianchuanReportAdRealtime{
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
		Time:              time,
	}
}

func (qrar *QianchuanReportAdRealtime) SetAdId(ctx context.Context, adId uint64) {
	qrar.AdId = adId
}

func (qrar *QianchuanReportAdRealtime) SetAdvertiserId(ctx context.Context, advertiserId uint64) {
	qrar.AdvertiserId = advertiserId
}

func (qrar *QianchuanReportAdRealtime) SetAwemeId(ctx context.Context, awemeId uint64) {
	qrar.AwemeId = awemeId
}

func (qrar *QianchuanReportAdRealtime) SetMarketingGoal(ctx context.Context, marketingGoal int64) {
	qrar.MarketingGoal = marketingGoal
}

func (qrar *QianchuanReportAdRealtime) SetStatCost(ctx context.Context, statCost float64) {
	qrar.StatCost = statCost
}

func (qrar *QianchuanReportAdRealtime) SetShowCnt(ctx context.Context, showCnt int64) {
	qrar.ShowCnt = showCnt
}

func (qrar *QianchuanReportAdRealtime) SetClickCnt(ctx context.Context, clickCnt int64) {
	qrar.ClickCnt = clickCnt
}

func (qrar *QianchuanReportAdRealtime) SetPayOrderCount(ctx context.Context, payOrderCount int64) {
	qrar.PayOrderCount = payOrderCount
}

func (qrar *QianchuanReportAdRealtime) SetCreateOrderAmount(ctx context.Context, createOrderAmount float64) {
	qrar.CreateOrderAmount = createOrderAmount
}

func (qrar *QianchuanReportAdRealtime) SetCreateOrderCount(ctx context.Context, createOrderCount int64) {
	qrar.CreateOrderCount = createOrderCount
}

func (qrar *QianchuanReportAdRealtime) SetPayOrderAmount(ctx context.Context, payOrderAmount float64) {
	qrar.PayOrderAmount = payOrderAmount
}

func (qrar *QianchuanReportAdRealtime) SetDyFollow(ctx context.Context, dyFollow int64) {
	qrar.DyFollow = dyFollow
}

func (qrar *QianchuanReportAdRealtime) SetConvertCnt(ctx context.Context, convertCnt int64) {
	qrar.ConvertCnt = convertCnt
}

func (qrar *QianchuanReportAdRealtime) SetTime(ctx context.Context, time int64) {
	qrar.Time = time
}

func (qrar *QianchuanReportAdRealtime) SetUpdateTime(ctx context.Context) {
	qrar.UpdateTime = time.Now()
}

func (qrar *QianchuanReportAdRealtime) SetCreateTime(ctx context.Context) {
	qrar.CreateTime = time.Now()
}

func (eqrar *ExternalQianchuanReportAdRealtime) SetRoi(ctx context.Context) {
	if eqrar.StatCost > 0 {
		eqrar.Roi = eqrar.PayOrderAmount / eqrar.StatCost
	}
}

func (eqrar *ExternalQianchuanReportAdRealtime) SetClickRate(ctx context.Context) {
	if eqrar.ShowCnt > 0 {
		eqrar.ClickRate = float64(eqrar.ClickCnt) / float64(eqrar.ShowCnt)
	}

	return
}

func (eqrar *ExternalQianchuanReportAdRealtime) SetPayConvertRate(ctx context.Context) {
	if eqrar.ClickCnt > 0 {
		eqrar.PayConvertRate = float64(eqrar.PayOrderCount) / float64(eqrar.ClickCnt)
	}

	return
}

func (eqrar *ExternalQianchuanReportAdRealtime) SetConvertCost(ctx context.Context) {
	if eqrar.ConvertCnt > 0 {
		eqrar.ConvertCost = float64(eqrar.StatCost) / float64(eqrar.ConvertCnt)
	}

	return
}

func (eqrar *ExternalQianchuanReportAdRealtime) SetConvertRate(ctx context.Context) {
	if eqrar.ClickCnt > 0 {
		eqrar.ConvertRate = float64(eqrar.ConvertCnt) / float64(eqrar.ClickCnt)
	}

	return
}

func (eqrar *ExternalQianchuanReportAdRealtime) SetAveragePayOrderStatCost(ctx context.Context) {
	if eqrar.PayOrderCount > 0 {
		eqrar.AveragePayOrderStatCost = eqrar.StatCost / float64(eqrar.PayOrderCount)
	}

	return
}

func (eqrar *ExternalQianchuanReportAdRealtime) SetPayOrderAveragePrice(ctx context.Context) {
	if eqrar.PayOrderCount > 0 {
		eqrar.PayOrderAveragePrice = eqrar.PayOrderAmount / float64(eqrar.PayOrderCount)
	}

	return
}

func (eqrar *ExternalQianchuanReportAdRealtime) SetCpmPlatform(ctx context.Context) {
	if eqrar.ShowCnt > 0 {
		eqrar.CpmPlatform = eqrar.StatCost / float64(eqrar.ShowCnt) * 1000
	}

	return
}

package domain

import (
	"context"
	"douyin/internal/pkg/tool"
	"time"
)

type QianchuanAdvertiserInfo struct {
	Id                      uint64
	Name                    string
	GeneralTotalBalance     float64
	Campaigns               uint64
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
	CreateTime              time.Time
	UpdateTime              time.Time
}

func NewQianchuanAdvertiserInfo(ctx context.Context, id, campaigns uint64, generalTotalBalance, statCost, payOrderAmount, createOrderAmount float64, payOrderCount, createOrderCount, clickCnt, showCnt, convertCnt, dyFollow int64, name string) *QianchuanAdvertiserInfo {
	return &QianchuanAdvertiserInfo{
		Id:                  id,
		Name:                name,
		GeneralTotalBalance: generalTotalBalance,
		Campaigns:           campaigns,
		StatCost:            statCost,
		PayOrderCount:       payOrderCount,
		PayOrderAmount:      payOrderAmount,
		CreateOrderCount:    createOrderCount,
		CreateOrderAmount:   createOrderAmount,
		ClickCnt:            clickCnt,
		ShowCnt:             showCnt,
		ConvertCnt:          convertCnt,
		DyFollow:            dyFollow,
	}
}

func (qai *QianchuanAdvertiserInfo) SetId(ctx context.Context, id uint64) {
	qai.Id = id
}

func (qai *QianchuanAdvertiserInfo) SetName(ctx context.Context, name string) {
	qai.Name = name
}

func (qai *QianchuanAdvertiserInfo) SetGeneralTotalBalance(ctx context.Context, generalTotalBalance float64) {
	qai.GeneralTotalBalance = generalTotalBalance
}

func (qai *QianchuanAdvertiserInfo) SetCampaigns(ctx context.Context, campaigns uint64) {
	qai.Campaigns = campaigns
}

func (qai *QianchuanAdvertiserInfo) SetStatCost(ctx context.Context, statCost float64) {
	qai.StatCost = statCost
}

func (qai *QianchuanAdvertiserInfo) SetRoi(ctx context.Context) {
	if qai.StatCost > 0 {
		qai.Roi = tool.Decimal(qai.PayOrderAmount/qai.StatCost, 2)
	}
}

func (qai *QianchuanAdvertiserInfo) SetPayOrderCount(ctx context.Context, payOrderCount int64) {
	qai.PayOrderCount = payOrderCount
}

func (qai *QianchuanAdvertiserInfo) SetPayOrderAmount(ctx context.Context, payOrderAmount float64) {
	qai.PayOrderAmount = payOrderAmount
}

func (qai *QianchuanAdvertiserInfo) SetCreateOrderAmount(ctx context.Context, createOrderAmount float64) {
	qai.CreateOrderAmount = createOrderAmount
}

func (qai *QianchuanAdvertiserInfo) SetCreateOrderCount(ctx context.Context, createOrderCount int64) {
	qai.CreateOrderCount = createOrderCount
}

func (qai *QianchuanAdvertiserInfo) SetClickCnt(ctx context.Context, clickCnt int64) {
	qai.ClickCnt = clickCnt
}

func (qai *QianchuanAdvertiserInfo) SetShowCnt(ctx context.Context, showCnt int64) {
	qai.ShowCnt = showCnt
}

func (qai *QianchuanAdvertiserInfo) SetConvertCnt(ctx context.Context, convertCnt int64) {
	qai.ConvertCnt = convertCnt
}

func (qai *QianchuanAdvertiserInfo) SetClickRate(ctx context.Context) {
	if qai.ShowCnt > 0 {
		qai.ClickRate = tool.Decimal(float64(qai.ClickCnt)/float64(qai.ShowCnt), 2)
	}
}

func (qai *QianchuanAdvertiserInfo) SetCpmPlatform(ctx context.Context) {
	if qai.ShowCnt > 0 {
		qai.CpmPlatform = tool.Decimal(qai.StatCost/float64(qai.ShowCnt)*1000, 2)
	}
}

func (qai *QianchuanAdvertiserInfo) SetDyFollow(ctx context.Context, dyFollow int64) {
	qai.DyFollow = dyFollow
}

func (qai *QianchuanAdvertiserInfo) SetPayConvertRate(ctx context.Context) {
	if qai.ClickCnt > 0 {
		qai.PayConvertRate = tool.Decimal(float64(qai.PayOrderCount)/float64(qai.ClickCnt), 2)
	}
}

func (qai *QianchuanAdvertiserInfo) SetConvertCost(ctx context.Context) {
	if qai.ConvertCnt > 0 {
		qai.ConvertCost = tool.Decimal(float64(qai.StatCost)/float64(qai.ConvertCnt), 2)
	}
}

func (qai *QianchuanAdvertiserInfo) SetConvertRate(ctx context.Context) {
	if qai.ClickCnt > 0 {
		qai.ConvertRate = tool.Decimal(float64(qai.ConvertCnt)/float64(qai.ClickCnt), 2)
	}
}

func (qai *QianchuanAdvertiserInfo) SetAveragePayOrderStatCost(ctx context.Context) {
	if qai.PayOrderCount > 0 {
		qai.AveragePayOrderStatCost = tool.Decimal(qai.StatCost/float64(qai.PayOrderCount), 2)
	}
}

func (qai *QianchuanAdvertiserInfo) SetPayOrderAveragePrice(ctx context.Context) {
	if qai.PayOrderCount > 0 {
		qai.PayOrderAveragePrice = tool.Decimal(qai.PayOrderAmount/float64(qai.PayOrderCount), 2)
	}
}

func (qai *QianchuanAdvertiserInfo) SetUpdateTime(ctx context.Context) {
	qai.UpdateTime = time.Now()
}

func (qai *QianchuanAdvertiserInfo) SetCreateTime(ctx context.Context) {
	qai.CreateTime = time.Now()
}

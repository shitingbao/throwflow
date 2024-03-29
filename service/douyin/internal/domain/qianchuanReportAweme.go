package domain

import (
	"context"
	"douyin/internal/pkg/tool"
	"time"
)

type QianchuanReportAweme struct {
	AdvertiserId            uint64
	AdvertiserName          string
	AdIds                   []uint64
	AwemeId                 uint64
	AwemeName               string
	AwemeShowId             string
	AwemeAvatar             string
	DyFollow                int64
	StatCost                float64
	PayOrderCount           int64
	PayOrderAmount          float64
	ShowCnt                 int64
	ClickCnt                int64
	ConvertCnt              int64
	Roi                     float64
	PayOrderAveragePrice    float64
	ClickRate               float64
	ConvertRate             float64
	AveragePayOrderStatCost float64
	CreateTime              time.Time
	UpdateTime              time.Time
}

type QianchuanReportAwemeList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*QianchuanReportAweme
}

type StatisticsQianchuanReportAweme struct {
	Key   string
	Value string
}

type StatisticsQianchuanReportAwemes struct {
	Statistics []*StatisticsQianchuanReportAweme
}

func NewQianchuanReportAweme(ctx context.Context, advertiserId, awemeId uint64, dyFollow, payOrderCount, showCnt, clickCnt, convertCnt int64, statCost, payOrderAmount float64, advertiserName, awemeName, awemeShowId, awemeAvatar string) *QianchuanReportAweme {
	return &QianchuanReportAweme{
		AdvertiserId:   advertiserId,
		AdvertiserName: advertiserName,
		AwemeId:        awemeId,
		AwemeName:      awemeName,
		AwemeShowId:    awemeShowId,
		AwemeAvatar:    awemeAvatar,
		DyFollow:       dyFollow,
		StatCost:       statCost,
		PayOrderCount:  payOrderCount,
		PayOrderAmount: payOrderAmount,
		ShowCnt:        showCnt,
		ClickCnt:       clickCnt,
		ConvertCnt:     convertCnt,
	}
}

func (qra *QianchuanReportAweme) SetAdvertiserId(ctx context.Context, advertiserId uint64) {
	qra.AdvertiserId = advertiserId
}

func (qra *QianchuanReportAweme) SetAdvertiserName(ctx context.Context, advertiserName string) {
	qra.AdvertiserName = advertiserName
}

func (qra *QianchuanReportAweme) SetAwemeId(ctx context.Context, awemeId uint64) {
	qra.AwemeId = awemeId
}

func (qra *QianchuanReportAweme) SetAwemeName(ctx context.Context, awemeName string) {
	qra.AwemeName = awemeName
}

func (qra *QianchuanReportAweme) SetAwemeShowId(ctx context.Context, awemeShowId string) {
	qra.AwemeShowId = awemeShowId
}

func (qra *QianchuanReportAweme) SetAwemeAvatar(ctx context.Context, awemeAvatar string) {
	qra.AwemeAvatar = awemeAvatar
}

func (qra *QianchuanReportAweme) SetDyFollow(ctx context.Context, dyFollow int64) {
	qra.DyFollow = dyFollow
}

func (qra *QianchuanReportAweme) SetStatCost(ctx context.Context, statCost float64) {
	qra.StatCost = statCost
}

func (qra *QianchuanReportAweme) SetPayOrderCount(ctx context.Context, payOrderCount int64) {
	qra.PayOrderCount = payOrderCount
}

func (qra *QianchuanReportAweme) SetPayOrderAmount(ctx context.Context, payOrderAmount float64) {
	qra.PayOrderAmount = payOrderAmount
}

func (qra *QianchuanReportAweme) SetShowCnt(ctx context.Context, showCnt int64) {
	qra.ShowCnt = showCnt
}

func (qra *QianchuanReportAweme) SetClickCnt(ctx context.Context, clickCnt int64) {
	qra.ClickCnt = clickCnt
}

func (qra *QianchuanReportAweme) SetConvertCnt(ctx context.Context, convertCnt int64) {
	qra.ConvertCnt = convertCnt
}

func (qra *QianchuanReportAweme) SetRoi(ctx context.Context) {
	if qra.StatCost > 0 {
		qra.Roi = tool.Decimal(qra.PayOrderAmount/qra.StatCost, 2)
	}
}

func (qra *QianchuanReportAweme) SetPayOrderAveragePrice(ctx context.Context) {
	if qra.PayOrderCount > 0 {
		qra.PayOrderAveragePrice = tool.Decimal(qra.PayOrderAmount/float64(qra.PayOrderCount), 2)
	}

}

func (qra *QianchuanReportAweme) SetClickRate(ctx context.Context) {
	if qra.ShowCnt > 0 {
		qra.ClickRate = tool.Decimal(float64(qra.ClickCnt)/float64(qra.ShowCnt), 2)
	}

}

func (qra *QianchuanReportAweme) SetConvertRate(ctx context.Context) {
	if qra.ClickCnt > 0 {
		qra.ConvertRate = tool.Decimal(float64(qra.ConvertCnt)/float64(qra.ClickCnt), 2)
	}
}

func (qra *QianchuanReportAweme) SetAveragePayOrderStatCost(ctx context.Context) {
	if qra.PayOrderCount > 0 {
		qra.AveragePayOrderStatCost = tool.Decimal(qra.StatCost/float64(qra.PayOrderCount), 2)
	}
}

func (qra *QianchuanReportAweme) SetUpdateTime(ctx context.Context) {
	qra.UpdateTime = time.Now()
}

func (qra *QianchuanReportAweme) SetCreateTime(ctx context.Context) {
	qra.CreateTime = time.Now()
}

package domain

import (
	"company/internal/pkg/tool"
	"context"
	"time"
)

type CompanyPerformanceDaily struct {
	Id             uint64
	CompanyId      uint64
	UserId         uint64
	Advertisers    string
	StatCost       float32
	PayOrderAmount float32
	Roi            float32
	Cost           float32
	RebalanceCosts []*CompanyPerformanceRebalance
	UpdateDay      time.Time
	CreateTime     time.Time
	UpdateTime     time.Time
}

type CompanyPerformanceDailyList struct {
	List []*CompanyPerformanceDaily
}

func NewCompanyPerformanceDaily(ctx context.Context, userId, companyId uint64, statCost, payOrderAmount, roi, cost float32, advertisers string) *CompanyPerformanceDaily {
	return &CompanyPerformanceDaily{
		CompanyId:      companyId,
		UserId:         userId,
		Advertisers:    advertisers,
		StatCost:       statCost,
		PayOrderAmount: payOrderAmount,
		Roi:            roi,
		Cost:           cost,
	}
}

func (cpd *CompanyPerformanceDaily) SetCompanyId(ctx context.Context, companyId uint64) {
	cpd.CompanyId = companyId
}

func (cpd *CompanyPerformanceDaily) SetUserId(ctx context.Context, userId uint64) {
	cpd.UserId = userId
}

func (cpd *CompanyPerformanceDaily) SetAdvertisers(ctx context.Context, advertisers string) {
	cpd.Advertisers = advertisers
}

func (cpd *CompanyPerformanceDaily) SetStatCost(ctx context.Context, statCost float32) {
	cpd.StatCost = statCost
}

func (cpd *CompanyPerformanceDaily) SetPayOrderAmount(ctx context.Context, payOrderAmount float32) {
	cpd.PayOrderAmount = payOrderAmount
}

func (cpd *CompanyPerformanceDaily) SetRoi(ctx context.Context, roi float32) {
	cpd.Roi = roi
}

func (cpd *CompanyPerformanceDaily) SetCost(ctx context.Context, cost float32) {
	cpd.Cost = cost
}

func (cpd *CompanyPerformanceDaily) SetUpdateDay(ctx context.Context, updateDay string) {
	updateTime, _ := tool.StringToTime("2006-01-02", updateDay)

	cpd.UpdateDay = updateTime
}

func (cpd *CompanyPerformanceDaily) SetUpdateTime(ctx context.Context) {
	cpd.UpdateTime = time.Now()
}

func (cpd *CompanyPerformanceDaily) SetCreateTime(ctx context.Context) {
	cpd.CreateTime = time.Now()
}

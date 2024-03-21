package domain

import (
	"context"
	"time"
)

type CompanyPerformanceRebalance struct {
	Id            uint64
	CompanyId     uint64
	UserId        uint64
	Cost          float32
	RebalanceType uint8
	Reason        string
	UpdateDay     time.Time
	CreateTime    time.Time
	UpdateTime    time.Time
}

func NewCompanyPerformanceRebalance(ctx context.Context, userId, companyId uint64, cost float32, ctype uint8, reason string) *CompanyPerformanceRebalance {
	return &CompanyPerformanceRebalance{
		CompanyId:     companyId,
		UserId:        userId,
		Cost:          cost,
		RebalanceType: ctype,
		Reason:        reason,
	}
}

func (cpr *CompanyPerformanceRebalance) SetUpdateDay(ctx context.Context) {
	cpr.UpdateDay = time.Now()
}

func (cpr *CompanyPerformanceRebalance) SetUpdateTime(ctx context.Context) {
	cpr.UpdateTime = time.Now()
}

func (cpr *CompanyPerformanceRebalance) SetCreateTime(ctx context.Context) {
	cpr.CreateTime = time.Now()
}

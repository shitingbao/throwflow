package domain

import (
	"context"
	"encoding/json"
	"time"
)

type CompanyPerformanceMonthly struct {
	Id                   uint64
	CompanyId            uint64
	UserId               uint64
	Username             string
	Job                  string
	QianchuanAdvertisers uint8
	Advertisers          string
	Advertiserst         []*Advertiser
	StatCost             float32
	StatCostProportion   float32
	Roi                  float32
	Cost                 float32
	RebalanceCost        float32
	TotalCost            float32
	UpdateDay            uint32
	CreateTime           time.Time
	UpdateTime           time.Time
}

type Advertiser struct {
	AdvertiserId   uint64 `json:"advertiserId"`
	AdvertiserName string `json:"advertiserName"`
}

type Advertisers []*Advertiser

type CompanyPerformanceMonthlyList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*CompanyPerformanceMonthly
}

func NewCompanyPerformanceMonthly(ctx context.Context, userId, companyId uint64, statCost, roi, cost, rebalanceCost, totalCost float32, username, job, advertisers string) *CompanyPerformanceMonthly {
	return &CompanyPerformanceMonthly{
		CompanyId:     companyId,
		UserId:        userId,
		Username:      username,
		Job:           job,
		Advertisers:   advertisers,
		StatCost:      statCost,
		Roi:           roi,
		Cost:          cost,
		RebalanceCost: rebalanceCost,
		TotalCost:     totalCost,
	}
}

func (cpm *CompanyPerformanceMonthly) SetCompanyId(ctx context.Context, companyId uint64) {
	cpm.CompanyId = companyId
}

func (cpm *CompanyPerformanceMonthly) SetUserId(ctx context.Context, userId uint64) {
	cpm.UserId = userId
}

func (cpm *CompanyPerformanceMonthly) SetUsername(ctx context.Context, username string) {
	cpm.Username = username
}

func (cpm *CompanyPerformanceMonthly) SetJob(ctx context.Context, job string) {
	cpm.Job = job
}

func (cpm *CompanyPerformanceMonthly) SetAdvertisers(ctx context.Context, advertisers string) {
	cpm.Advertisers = advertisers
}

func (cpm *CompanyPerformanceMonthly) SetStatCost(ctx context.Context, statCost float32) {
	cpm.StatCost = statCost
}

func (cpm *CompanyPerformanceMonthly) SetRoi(ctx context.Context, roi float32) {
	cpm.Roi = roi
}

func (cpm *CompanyPerformanceMonthly) SetCost(ctx context.Context, cost float32) {
	cpm.Cost = cost
}

func (cpm *CompanyPerformanceMonthly) SetRebalanceCost(ctx context.Context, rebalanceCost float32) {
	cpm.RebalanceCost = rebalanceCost
}

func (cpm *CompanyPerformanceMonthly) SetTotalCost(ctx context.Context, totalCost float32) {
	cpm.TotalCost = totalCost
}

func (cpm *CompanyPerformanceMonthly) SetUpdateDay(ctx context.Context, updateDay uint32) {
	cpm.UpdateDay = updateDay
}

func (cpm *CompanyPerformanceMonthly) SetUpdateTime(ctx context.Context) {
	cpm.UpdateTime = time.Now()
}

func (cpm *CompanyPerformanceMonthly) SetCreateTime(ctx context.Context) {
	cpm.CreateTime = time.Now()
}

func (cpm *CompanyPerformanceMonthly) GetAdvertisers(ctx context.Context) {
	var advertisers Advertisers

	if err := json.Unmarshal([]byte(cpm.Advertisers), &advertisers); err == nil {
		cpm.Advertiserst = advertisers
	}
}

func (cpm *CompanyPerformanceMonthly) GetQianchuanAdvertisers(ctx context.Context) {
	cpm.QianchuanAdvertisers = uint8(len(cpm.Advertiserst))
}

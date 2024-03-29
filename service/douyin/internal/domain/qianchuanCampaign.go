package domain

import (
	"context"
	"time"
)

type QianchuanCampaign struct {
	Id             uint64
	AdvertiserId   uint64
	Name           string
	Budget         float64
	BudgetMode     string
	MarketingGoal  string
	MarketingScene string
	Status         string
	CreateDate     string
	CreateTime     time.Time
	UpdateTime     time.Time
}

type QianchuanCampaignList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*QianchuanCampaign
}

func NewQianchuanCampaign(ctx context.Context, id, advertiserId uint64, budget float64, name, budgetMode, marketingGoal, marketingScene, status, createDate string) *QianchuanCampaign {
	return &QianchuanCampaign{
		Id:             id,
		AdvertiserId:   advertiserId,
		Name:           name,
		Budget:         budget,
		BudgetMode:     budgetMode,
		MarketingGoal:  marketingGoal,
		MarketingScene: marketingScene,
		Status:         status,
		CreateDate:     createDate,
	}
}

func (qc *QianchuanCampaign) SetId(ctx context.Context, id uint64) {
	qc.Id = id
}

func (qc *QianchuanCampaign) SetAdvertiserId(ctx context.Context, advertiserId uint64) {
	qc.AdvertiserId = advertiserId
}

func (qc *QianchuanCampaign) SetBudget(ctx context.Context, budget float64) {
	qc.Budget = budget
}

func (qc *QianchuanCampaign) SetName(ctx context.Context, name string) {
	qc.Name = name
}

func (qc *QianchuanCampaign) SetBudgetMode(ctx context.Context, budgetMode string) {
	qc.BudgetMode = budgetMode
}

func (qc *QianchuanCampaign) SetMarketingGoal(ctx context.Context, marketingGoal string) {
	qc.MarketingGoal = marketingGoal
}

func (qc *QianchuanCampaign) SetMarketingScene(ctx context.Context, marketingScene string) {
	qc.MarketingScene = marketingScene
}

func (qc *QianchuanCampaign) SetStatus(ctx context.Context, status string) {
	qc.Status = status
}

func (qc *QianchuanCampaign) SetCreateDate(ctx context.Context, createDate string) {
	qc.CreateDate = createDate
}

func (qc *QianchuanCampaign) SetUpdateTime(ctx context.Context) {
	qc.UpdateTime = time.Now()
}

func (qc *QianchuanCampaign) SetCreateTime(ctx context.Context) {
	qc.CreateTime = time.Now()
}

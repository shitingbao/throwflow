package domain

import (
	"context"
	"time"
)

type CompanyTask struct {
	Id             uint64
	ProductOutId   uint64
	ExpireTime     uint64
	PlayNum        uint64
	Price          uint64
	Quota          uint64
	ClaimQuota     uint64
	SuccessQuota   uint64
	IsTop          uint8
	IsDel          uint8
	IsGoodReviews  uint8
	CreateTime     time.Time
	UpdateTime     time.Time
	CompanyProduct CompanyProduct
}

type CompanyTaskList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*CompanyTask
}

func NewCompanyTask(ctx context.Context, productOutId, expireTime, playNum, price, quota uint64, isGoodReviews uint8) *CompanyTask {
	return &CompanyTask{
		ProductOutId:  productOutId,
		ExpireTime:    expireTime,
		PlayNum:       playNum,
		Price:         price,
		Quota:         quota,
		IsGoodReviews: isGoodReviews,
	}
}

func (c *CompanyTask) SetQuota(ctx context.Context, quota uint64) {
	c.Quota = quota
}

func (c *CompanyTask) SetClaimQuota(ctx context.Context, ClaimQuota uint64) {
	c.ClaimQuota = ClaimQuota
}

func (c *CompanyTask) SetSuccessQuota(ctx context.Context, SuccessQuota uint64) {
	c.SuccessQuota = SuccessQuota
}

func (c *CompanyTask) SetIsTop(ctx context.Context, isTop uint8) {
	c.IsTop = isTop
}

func (c *CompanyTask) SetUpdateTime(ctx context.Context) {
	c.UpdateTime = time.Now()
}

func (c *CompanyTask) SetCreateTime(ctx context.Context) {
	c.CreateTime = time.Now()
}

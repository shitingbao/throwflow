package domain

import (
	"context"
	"time"
)

type Collect struct {
	CompanyId  uint64
	Phone      string
	VideoId    uint64
	CreateTime time.Time
	UpdateTime time.Time
}

func NewCollect(ctx context.Context, companyId, videoId uint64, phone string) *Collect {
	return &Collect{
		CompanyId: companyId,
		VideoId:   videoId,
		Phone:     phone,
	}
}

func (c *Collect) SetCompanyId(ctx context.Context, companyId uint64) {
	c.CompanyId = companyId
}

func (c *Collect) SetPhone(ctx context.Context, phone string) {
	c.Phone = phone
}

func (c *Collect) SetVideoId(ctx context.Context, videoId uint64) {
	c.VideoId = videoId
}

func (c *Collect) SetUpdateTime(ctx context.Context) {
	c.UpdateTime = time.Now()
}

func (c *Collect) SetCreateTime(ctx context.Context) {
	c.CreateTime = time.Now()
}

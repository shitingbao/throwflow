package domain

import (
	"context"
	"time"
)

type QianchuanAdvertiserStatus struct {
	CompanyId    uint64
	AdvertiserId uint64
	Status       uint8
	Day          uint32
	CreateTime   time.Time
	UpdateTime   time.Time
}

func NewQianchuanAdvertiserStatus(ctx context.Context, companyId, advertiserId uint64, day uint32, status uint8) *QianchuanAdvertiserStatus {
	return &QianchuanAdvertiserStatus{
		CompanyId:    companyId,
		AdvertiserId: advertiserId,
		Status:       status,
		Day:          day,
	}
}

func (qas *QianchuanAdvertiserStatus) SetCompanyId(ctx context.Context, companyId uint64) {
	qas.CompanyId = companyId
}

func (qas *QianchuanAdvertiserStatus) SetAdvertiserId(ctx context.Context, advertiserId uint64) {
	qas.AdvertiserId = advertiserId
}

func (qas *QianchuanAdvertiserStatus) SetStatus(ctx context.Context, status uint8) {
	qas.Status = status
}

func (qas *QianchuanAdvertiserStatus) SetDay(ctx context.Context, day uint32) {
	qas.Day = day
}

func (qas *QianchuanAdvertiserStatus) SetUpdateTime(ctx context.Context) {
	qas.UpdateTime = time.Now()
}

func (qas *QianchuanAdvertiserStatus) SetCreateTime(ctx context.Context) {
	qas.CreateTime = time.Now()
}

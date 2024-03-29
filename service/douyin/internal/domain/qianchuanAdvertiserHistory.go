package domain

import (
	"context"
	"time"
)

type QianchuanAdvertiserHistory struct {
	Id           uint64
	AdvertiserId uint64
	Day          uint32
	CreateTime   time.Time
	UpdateTime   time.Time
}

func NewQianchuanAdvertiserHistory(ctx context.Context, advertiserId uint64, day uint32) *QianchuanAdvertiserHistory {
	return &QianchuanAdvertiserHistory{
		AdvertiserId: advertiserId,
		Day:          day,
	}
}

func (qah *QianchuanAdvertiserHistory) SetAdvertiserId(ctx context.Context, advertiserId uint64) {
	qah.AdvertiserId = advertiserId
}

func (qah *QianchuanAdvertiserHistory) SetDay(ctx context.Context, day uint32) {
	qah.Day = day
}

func (qah *QianchuanAdvertiserHistory) SetUpdateTime(ctx context.Context) {
	qah.UpdateTime = time.Now()
}

func (qah *QianchuanAdvertiserHistory) SetCreateTime(ctx context.Context) {
	qah.CreateTime = time.Now()
}

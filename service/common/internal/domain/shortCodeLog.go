package domain

import (
	"context"
	"time"
)

type ShortCodeLog struct {
	Id         uint64
	ShortCode  string
	CreateTime time.Time
	UpdateTime time.Time
}

func NewShortCodeLog(ctx context.Context, shortCode string) *ShortCodeLog {
	return &ShortCodeLog{
		ShortCode: shortCode,
	}
}

func (scl *ShortCodeLog) SetShortCode(ctx context.Context, shortCode string) {
	scl.ShortCode = shortCode
}

func (scl *ShortCodeLog) SetUpdateTime(ctx context.Context) {
	scl.UpdateTime = time.Now()
}

func (scl *ShortCodeLog) SetCreateTime(ctx context.Context) {
	scl.CreateTime = time.Now()
}

package biz

import (
	"context"
	"douyin/internal/domain"
)

type QianchuanReportAdRealtimeRepo interface {
	Get(context.Context, uint64, string) (*domain.QianchuanReportAdRealtime, error)
	GetByTime(context.Context, uint64, int64, string) (*domain.QianchuanReportAdRealtime, error)
	ListByAdIds(context.Context, []uint64, int64, string) ([]*domain.QianchuanReportAdRealtime, error)
	ListAdvertisers(context.Context, string, string) ([]*domain.QianchuanReportAdRealtime, error)
	StatisticsAdvertisers(context.Context, string, string) (int64, error)
	SaveIndex(context.Context, string)
	Upsert(context.Context, string, *domain.QianchuanReportAdRealtime) error
}

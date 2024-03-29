package biz

import (
	"context"
	"douyin/internal/domain"
)

type QianchuanReportAdRepo interface {
	Get(context.Context, uint64, uint64, string) (*domain.QianchuanReportAd, error)
	ListByMarketingGoal(context.Context, string, string) ([]*domain.QianchuanReportAd, error)
	SaveIndex(context.Context, string)
	Upsert(context.Context, string, *domain.QianchuanReportAd) error
}

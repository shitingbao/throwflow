package biz

import (
	"context"
	"douyin/internal/domain"
)

type DoukeOrderInfoRepo interface {
	List(context.Context) ([]*domain.DoukeOrderInfo, error)
	Statistics(context.Context, uint64, string, string, string) (*domain.DoukeOrderInfo, error)
	StatisticsRealcommission(context.Context, uint64, string, string) (*domain.DoukeOrderInfo, error)
	Upsert(context.Context, *domain.DoukeOrderInfo) error
	GetByUserIdAndProductId(context.Context, uint64, string, string, string) (*domain.DoukeOrderInfo, error)
}

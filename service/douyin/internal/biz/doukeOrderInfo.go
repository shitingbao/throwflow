package biz

import (
	"context"
	"douyin/internal/domain"
)

type DoukeOrderInfoRepo interface {
	List(context.Context) ([]*domain.DoukeOrderInfo, error)
	ListOperation(context.Context, int, int) ([]*domain.DoukeOrderInfo, error)
	Count(context.Context, uint64, string, string) (int64, error)
	CountOperation(context.Context) (int64, error)
	Statistics(context.Context, uint64, string, string, string) (*domain.DoukeOrderInfo, error)
	StatisticsRealcommission(context.Context, uint64, string, string) (*domain.DoukeOrderInfo, error)
	StatisticsByProductId(context.Context, uint64, uint64, string, string, string) (*domain.DoukeOrderInfo, error)
	Upsert(context.Context, *domain.DoukeOrderInfo) error
	GetByUserIdAndProductId(context.Context, uint64, string, string, string) (*domain.DoukeOrderInfo, error)
}

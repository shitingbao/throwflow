package biz

import (
	"context"
	"douyin/internal/domain"
)

type JinritemaiOrderInfoRepo interface {
	GetByClientKeyAndOpenId(context.Context, string, string) (*domain.JinritemaiOrderInfo, error)
	GetByClientKeyAndOpenIdAndMediaTypeAndMediaId(context.Context, string, string, string, string) (*domain.JinritemaiOrderInfo, error)
	GetIsTopByProductId(context.Context, uint64) (*domain.JinritemaiOrderInfo, error)
	List(context.Context, int, int, []*domain.OpenDouyinToken, string, string) ([]*domain.JinritemaiOrderInfo, error)
	ListOperation(context.Context, int, int) ([]*domain.JinritemaiOrderInfo, error)
	ListByProductIds(context.Context, []*domain.OpenDouyinToken, []string) ([]*domain.JinritemaiOrderInfo, error)
	ListProductByClientKeyAndOpenIdAndMediaType(context.Context, string, string, string) ([]*domain.JinritemaiOrderInfo, error)
	ListNotRefund(context.Context) ([]*domain.JinritemaiOrderInfo, error)
	ListRefund(context.Context) ([]*domain.JinritemaiOrderInfo, error)
	ListMediaId(context.Context) ([]*domain.JinritemaiOrderInfo, error)
	ListProductId(context.Context) ([]*domain.JinritemaiOrderInfo, error)
	ListByPickExtra(context.Context, uint8) ([]*domain.JinritemaiOrderInfo, error)
	ListByProductIdAndMediaIds(context.Context, []*domain.CommissionRateJinritemaiOrder) ([]*domain.JinritemaiOrderInfo, error)
	ListProductIdAndMediaIds(context.Context, int, int) ([]*domain.JinritemaiOrderInfo, error)
	Count(context.Context, []*domain.OpenDouyinToken, string, string) (int64, error)
	CountOperation(context.Context) (int64, error)
	CountProductIdAndMediaIds(context.Context) (int64, error)
	Statistics(context.Context, []*domain.OpenDouyinToken, string, string, string, string) (*domain.JinritemaiOrderInfo, error)
	StatisticsRealcommission(context.Context, []*domain.OpenDouyinToken, string, string, string) (*domain.JinritemaiOrderInfo, error)
	StatisticsRealcommissionPayTime(context.Context, []*domain.OpenDouyinToken, string, string, string, string) (*domain.JinritemaiOrderInfo, error)
	StatisticsAwemeIndustry(context.Context, uint64, string, string, []*domain.OpenDouyinUserInfo) ([]*domain.JinritemaiOrderInfoStatisticsAwemeIndustry, error)
	Upsert(context.Context, *domain.JinritemaiOrderInfo) error
}

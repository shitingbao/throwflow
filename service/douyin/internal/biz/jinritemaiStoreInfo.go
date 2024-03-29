package biz

import (
	"context"
	"douyin/internal/domain"
)

type JinritemaiStoreInfoRepo interface {
	List(context.Context, int, int, []*domain.OpenDouyinToken) ([]*domain.JinritemaiStoreInfo, error)
	ListByIds(context.Context, []uint64) ([]*domain.JinritemaiStoreInfo, error)
	ListByProductId(context.Context, string) ([]*domain.OpenDouyinUserInfo, error)
	ListByClientKeyAndOpenIdAndProductIds(context.Context, string, string, []string) ([]*domain.JinritemaiStoreInfo, error)
	ListProductId(context.Context, int, int) ([]*domain.JinritemaiStoreInfo, error)
	Count(context.Context, []*domain.OpenDouyinToken) (int64, error)
	CountProductId(context.Context) (int64, error)
	DeleteByDayAndClientKeyAndOpenIdAndProductIds(context.Context, string, string, []string) error
}

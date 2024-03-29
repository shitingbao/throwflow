package biz

import (
	"context"
	"douyin/internal/domain"
)

type QianchuanAdvertiserInfoRepo interface {
	Get(context.Context, uint64, string) (*domain.QianchuanAdvertiserInfo, error)
	GetByDay(context.Context, uint64, string, string) (*domain.QianchuanAdvertiserInfo, error)
	List(context.Context, string, string, string, uint64, uint64) ([]*domain.QianchuanAdvertiserInfo, error)
	Count(context.Context, string, string, string) (uint64, error)
	Statistics(context.Context, string, string, string) (*domain.QianchuanAdvertiserInfo, error)
	SaveIndex(context.Context, string)
	Upsert(context.Context, string, *domain.QianchuanAdvertiserInfo) error
	Delete(context.Context, string, []uint64) error
}

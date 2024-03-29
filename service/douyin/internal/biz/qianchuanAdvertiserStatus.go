package biz

import (
	"context"
	"douyin/internal/conf"
	"douyin/internal/domain"
	"github.com/go-kratos/kratos/v2/log"
)

type QianchuanAdvertiserStatusRepo interface {
	GetByCompanyIdAndAdvertiserIdAndDay(context.Context, uint64, uint64, uint32) (*domain.QianchuanAdvertiserStatus, error)
	List(context.Context, uint64, uint32) ([]*domain.QianchuanAdvertiserStatus, error)
	Count(context.Context, uint64, uint64, uint32) (int64, error)
	DeleteByCompanyIdAndAdvertiserIdAndDay(context.Context, uint64, uint64, uint32) error
	Save(context.Context, *domain.QianchuanAdvertiserStatus) (*domain.QianchuanAdvertiserStatus, error)
	Update(context.Context, *domain.QianchuanAdvertiserStatus) (*domain.QianchuanAdvertiserStatus, error)
}

type QianchuanAdvertiserStatusUsecase struct {
	repo QianchuanAdvertiserStatusRepo
	conf *conf.Data
	log  *log.Helper
}

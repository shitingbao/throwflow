package biz

import (
	"context"
	"douyin/internal/conf"
	"douyin/internal/domain"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	DouyinOpenDouyinUserInfoCreateLogNotFound    = errors.InternalServer("DOUYIN_OPEN_DOUYIN_USER_INFO_CREATE_LOG_NOT_FOUND", "抖音开放平台用户信息新增记录不存在")
	DouyinOpenDouyinUserInfoCreateLogUpdateError = errors.InternalServer("DOUYIN_OPEN_DOUYIN_USER_INFO_CREATE_LOG_UPDATE_ERROR", "抖音开放平台用户信息新增记录更新失败")
)

type OpenDouyinUserInfoCreateLogRepo interface {
	List(context.Context, string) ([]*domain.OpenDouyinUserInfoCreateLog, error)
	Save(context.Context, *domain.OpenDouyinUserInfoCreateLog) (*domain.OpenDouyinUserInfoCreateLog, error)
	Update(context.Context, *domain.OpenDouyinUserInfoCreateLog) (*domain.OpenDouyinUserInfoCreateLog, error)
	Delete(context.Context, *domain.OpenDouyinUserInfoCreateLog) error
}

type OpenDouyinUserInfoCreateLogUsecase struct {
	repo OpenDouyinUserInfoCreateLogRepo
	tm   Transaction
	conf *conf.Data
	log  *log.Helper
}

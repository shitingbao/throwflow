package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "task/api/service/douyin/v1"
	"task/internal/pkg/tool"
)

type JinritemaiStoreRepo interface {
	Sync(context.Context) (*v1.SyncJinritemaiStoresReply, error)
}

type JinritemaiStoreUsecase struct {
	repo JinritemaiStoreRepo
	log  *log.Helper
}

func NewJinritemaiStoreUsecase(repo JinritemaiStoreRepo, logger log.Logger) *JinritemaiStoreUsecase {
	return &JinritemaiStoreUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (jsuc *JinritemaiStoreUsecase) SyncJinritemaiStores(ctx context.Context) (*v1.SyncJinritemaiStoresReply, error) {
	jinritemaiStore, err := jsuc.repo.Sync(ctx)

	if err != nil {
		return nil, errors.InternalServer("TASK_SYNC_JINRITEMAI_STORE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return jinritemaiStore, nil
}

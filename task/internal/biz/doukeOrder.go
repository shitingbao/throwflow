package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "task/api/service/douyin/v1"
	"task/internal/pkg/tool"
	"time"
)

type DoukeOrderRepo interface {
	Sync(context.Context, string) (*v1.SyncDoukeOrdersReply, error)
}

type DoukeOrderUsecase struct {
	repo DoukeOrderRepo
	log  *log.Helper
}

func NewDoukeOrderUsecase(repo DoukeOrderRepo, logger log.Logger) *DoukeOrderUsecase {
	return &DoukeOrderUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (douc *DoukeOrderUsecase) SyncDoukeOrders(ctx context.Context) (*v1.SyncDoukeOrdersReply, error) {
	day := time.Now().Format("2006-01-02")

	doukeOrder, err := douc.repo.Sync(ctx, day)

	if err != nil {
		return nil, errors.InternalServer("TASK_SYNC_DOUKE_ORDER_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return doukeOrder, nil
}

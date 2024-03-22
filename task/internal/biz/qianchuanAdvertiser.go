package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "task/api/service/douyin/v1"
	"task/internal/pkg/tool"
)

type QianchuanAdvertiserRepo interface {
	Sync(context.Context) (*v1.SyncQianchuanDatasReply, error)
	SyncRDS(context.Context) (*v1.SyncRdsDatasReply, error)
}

type QianchuanAdvertiserUsecase struct {
	repo QianchuanAdvertiserRepo
	log  *log.Helper
}

func NewQianchuanAdvertiserUsecase(repo QianchuanAdvertiserRepo, logger log.Logger) *QianchuanAdvertiserUsecase {
	return &QianchuanAdvertiserUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (qauc *QianchuanAdvertiserUsecase) SyncQianchuanDatas(ctx context.Context) (*v1.SyncQianchuanDatasReply, error) {
	qianchuanAdvertiser, err := qauc.repo.Sync(ctx)

	if err != nil {
		return nil, errors.InternalServer("TASK_SYNC_QIANCHUAN_DATA_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return qianchuanAdvertiser, nil
}

func (qauc *QianchuanAdvertiserUsecase) SyncRdsDatas(ctx context.Context) (*v1.SyncRdsDatasReply, error) {
	rds, err := qauc.repo.SyncRDS(ctx)

	if err != nil {
		return nil, errors.InternalServer("TASK_SYNC_RDS_DATA_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return rds, nil
}

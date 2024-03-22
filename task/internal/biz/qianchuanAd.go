package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "task/api/service/douyin/v1"
	"task/internal/pkg/tool"
)

type QianchuanAdRepo interface {
	Sync(context.Context) (*v1.SyncQianchuanAdsReply, error)
}

type QianchuanAdUsecase struct {
	repo QianchuanAdRepo
	log  *log.Helper
}

func NewQianchuanAdUsecase(repo QianchuanAdRepo, logger log.Logger) *QianchuanAdUsecase {
	return &QianchuanAdUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (qauc *QianchuanAdUsecase) SyncQianchuanAds(ctx context.Context) (*v1.SyncQianchuanAdsReply, error) {
	qianchuanAd, err := qauc.repo.Sync(ctx)

	if err != nil {
		return nil, errors.InternalServer("TASK_SYNC_QIANCHUAN_AD_DATA_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return qianchuanAd, nil
}

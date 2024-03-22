package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "task/api/service/douyin/v1"
	"task/internal/pkg/tool"
)

type OpenDouyinVideoRepo interface {
	Sync(context.Context) (*v1.SyncOpenDouyinVideosReply, error)
}

type OpenDouyinVideoUsecase struct {
	repo OpenDouyinVideoRepo
	log  *log.Helper
}

func NewOpenDouyinVideoUsecase(repo OpenDouyinVideoRepo, logger log.Logger) *OpenDouyinVideoUsecase {
	return &OpenDouyinVideoUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (odvuc *OpenDouyinVideoUsecase) SyncOpenDouyinVideos(ctx context.Context) (*v1.SyncOpenDouyinVideosReply, error) {
	openDouyinVideo, err := odvuc.repo.Sync(ctx)

	if err != nil {
		return nil, errors.InternalServer("TASK_SYNC_OPEN_DOUYIN_VIDEO_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return openDouyinVideo, nil
}

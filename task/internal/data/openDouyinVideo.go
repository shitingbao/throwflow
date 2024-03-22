package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "task/api/service/douyin/v1"
	"task/internal/biz"
)

type openDouyinVideoRepo struct {
	data *Data
	log  *log.Helper
}

func NewOpenDouyinVideoRepo(data *Data, logger log.Logger) biz.OpenDouyinVideoRepo {
	return &openDouyinVideoRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (odvr *openDouyinVideoRepo) Sync(ctx context.Context) (*v1.SyncOpenDouyinVideosReply, error) {
	openDouyinVideo, err := odvr.data.douyinuc.SyncOpenDouyinVideos(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return openDouyinVideo, err
}

package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "task/api/service/douyin/v1"
	"task/internal/biz"
)

type qianchuanAdvertiserRepo struct {
	data *Data
	log  *log.Helper
}

func NewQianchuanAdvertiserRepo(data *Data, logger log.Logger) biz.QianchuanAdvertiserRepo {
	return &qianchuanAdvertiserRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (qar *qianchuanAdvertiserRepo) Sync(ctx context.Context) (*v1.SyncQianchuanDatasReply, error) {
	qianchuanAdvertiser, err := qar.data.douyinuc.SyncQianchuanDatas(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return qianchuanAdvertiser, err
}

func (qar *qianchuanAdvertiserRepo) SyncRDS(ctx context.Context) (*v1.SyncRdsDatasReply, error) {
	rds, err := qar.data.douyinuc.SyncRdsDatas(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return rds, err
}

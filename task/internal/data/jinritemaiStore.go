package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "task/api/service/douyin/v1"
	"task/internal/biz"
)

type jinritemaiStoreRepo struct {
	data *Data
	log  *log.Helper
}

func NewJinritemaiStoreRepo(data *Data, logger log.Logger) biz.JinritemaiStoreRepo {
	return &jinritemaiStoreRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (jsr *jinritemaiStoreRepo) Sync(ctx context.Context) (*v1.SyncJinritemaiStoresReply, error) {
	jinritemaiStore, err := jsr.data.douyinuc.SyncJinritemaiStores(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return jinritemaiStore, err
}

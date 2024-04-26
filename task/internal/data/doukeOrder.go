package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "task/api/service/douyin/v1"
	"task/internal/biz"
)

type doukeOrderRepo struct {
	data *Data
	log  *log.Helper
}

func NewDoukeOrderRepo(data *Data, logger log.Logger) biz.DoukeOrderRepo {
	return &doukeOrderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (dor *doukeOrderRepo) Sync(ctx context.Context, day string) (*v1.SyncDoukeOrdersReply, error) {
	doukeOrder, err := dor.data.douyinuc.SyncDoukeOrders(ctx, &v1.SyncDoukeOrdersRequest{
		Day: day,
	})

	if err != nil {
		return nil, err
	}

	return doukeOrder, err
}

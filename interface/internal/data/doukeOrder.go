package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/douyin/v1"
	"interface/internal/biz"
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

func (dor *doukeOrderRepo) Statistics(ctx context.Context, userId uint64) (*v1.StatisticsDoukeOrdersReply, error) {
	list, err := dor.data.douyinuc.StatisticsDoukeOrders(ctx, &v1.StatisticsDoukeOrdersRequest{
		UserId: userId,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

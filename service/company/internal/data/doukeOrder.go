package data

import (
	v1 "company/api/service/douyin/v1"
	"company/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
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

func (dor *doukeOrderRepo) Get(ctx context.Context, userId uint64, productOutId, flowPoint, createTime string) (*v1.GetDoukeOrdersReply, error) {
	doukeOrder, err := dor.data.douyinuc.GetDoukeOrders(ctx, &v1.GetDoukeOrdersRequest{
		UserId:     userId,
		FlowPoint:  flowPoint,
		ProductId:  productOutId,
		CreateTime: createTime,
	})

	if err != nil {
		return nil, err
	}

	return doukeOrder, nil
}

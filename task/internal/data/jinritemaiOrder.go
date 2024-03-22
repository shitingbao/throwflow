package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "task/api/service/douyin/v1"
	"task/internal/biz"
)

type jinritemaiOrderRepo struct {
	data *Data
	log  *log.Helper
}

func NewJinritemaiOrderRepo(data *Data, logger log.Logger) biz.JinritemaiOrderRepo {
	return &jinritemaiOrderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (jor *jinritemaiOrderRepo) Sync90Day(ctx context.Context) (*v1.Sync90DayJinritemaiOrdersReply, error) {
	jinritemaiOrder, err := jor.data.douyinuc.Sync90DayJinritemaiOrders(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return jinritemaiOrder, err
}

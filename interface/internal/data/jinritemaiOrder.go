package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/douyin/v1"
	"interface/internal/biz"
)

type jinritemailOrderRepo struct {
	data *Data
	log  *log.Helper
}

func NewJinritemailOrderRepo(data *Data, logger log.Logger) biz.JinritemailOrderRepo {
	return &jinritemailOrderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (jor *jinritemailOrderRepo) GetStorePreference(ctx context.Context, userId uint64) (*v1.GetStorePreferenceJinritemaiOrdersReply, error) {
	list, err := jor.data.douyinuc.GetStorePreferenceJinritemaiOrders(ctx, &v1.GetStorePreferenceJinritemaiOrdersRequest{
		UserId: userId,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (jor *jinritemailOrderRepo) List(ctx context.Context, userId, pageNum, pageSize uint64, startDay, endDay string) (*v1.ListJinritemaiOrdersReply, error) {
	list, err := jor.data.douyinuc.ListJinritemaiOrders(ctx, &v1.ListJinritemaiOrdersRequest{
		PageNum:  pageNum,
		PageSize: pageSize,
		UserId:   userId,
		StartDay: startDay,
		EndDay:   endDay,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (jor *jinritemailOrderRepo) Statistics(ctx context.Context, userId uint64, startDay, endDay string) (*v1.StatisticsJinritemaiOrdersReply, error) {
	list, err := jor.data.douyinuc.StatisticsJinritemaiOrders(ctx, &v1.StatisticsJinritemaiOrdersRequest{
		UserId:   userId,
		StartDay: startDay,
		EndDay:   endDay,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

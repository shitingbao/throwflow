package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "weixin/api/service/douyin/v1"
	"weixin/internal/biz"
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

func (jor *jinritemaiOrderRepo) ListByPickExtra(ctx context.Context) (*v1.ListJinritemaiOrderByPickExtrasReply, error) {
	list, err := jor.data.douyinuc.ListJinritemaiOrderByPickExtras(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (jor *jinritemaiOrderRepo) Statistics(ctx context.Context, userId uint64, startDay, endDay string) (*v1.StatisticsJinritemaiOrdersReply, error) {
	statistics, err := jor.data.douyinuc.StatisticsJinritemaiOrders(ctx, &v1.StatisticsJinritemaiOrdersRequest{
		UserId:   userId,
		StartDay: startDay,
		EndDay:   endDay,
	})

	if err != nil {
		return nil, err
	}

	return statistics, err
}

func (jor *jinritemaiOrderRepo) StatisticsByDay(ctx context.Context, day, content, pickExtra string) (*v1.StatisticsJinritemaiOrderByDaysReply, error) {
	statistics, err := jor.data.douyinuc.StatisticsJinritemaiOrderByDays(ctx, &v1.StatisticsJinritemaiOrderByDaysRequest{
		Day:       day,
		Content:   content,
		PickExtra: pickExtra,
	})

	if err != nil {
		return nil, err
	}

	return statistics, err
}

func (jor *jinritemaiOrderRepo) StatisticsByPayTimeDay(ctx context.Context, payTime, day, content, pickExtra string) (*v1.StatisticsJinritemaiOrderByPayTimeAndDaysReply, error) {
	statistics, err := jor.data.douyinuc.StatisticsJinritemaiOrderByPayTimeAndDays(ctx, &v1.StatisticsJinritemaiOrderByPayTimeAndDaysRequest{
		Day:       day,
		PayTime:   payTime,
		Content:   content,
		PickExtra: pickExtra,
	})

	if err != nil {
		return nil, err
	}

	return statistics, err
}

func (jor *jinritemaiOrderRepo) StatisticsByClientKeyAndOpenId(ctx context.Context, clientKey, openId, startDay, endDay string) (*v1.StatisticsJinritemaiOrderByClientKeyAndOpenIdsReply, error) {
	statistics, err := jor.data.douyinuc.StatisticsJinritemaiOrderByClientKeyAndOpenIds(ctx, &v1.StatisticsJinritemaiOrderByClientKeyAndOpenIdsRequest{
		ClientKey: clientKey,
		OpenId:    openId,
		StartDay:  startDay,
		EndDay:    endDay,
	})

	if err != nil {
		return nil, err
	}

	return statistics, err
}

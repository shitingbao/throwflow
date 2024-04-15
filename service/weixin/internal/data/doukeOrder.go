package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "weixin/api/service/douyin/v1"
	"weixin/internal/biz"
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

func (dor *doukeOrderRepo) ListUserId(ctx context.Context) (*v1.ListUserIdDoukeOrdersReply, error) {
	list, err := dor.data.douyinuc.ListUserIdDoukeOrders(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (dor *doukeOrderRepo) StatisticsByDay(ctx context.Context, userId uint64, day string) (*v1.StatisticsDoukeOrderByDaysReply, error) {
	statistics, err := dor.data.douyinuc.StatisticsDoukeOrderByDays(ctx, &v1.StatisticsDoukeOrderByDaysRequest{
		UserId: userId,
		Day:    day,
	})

	if err != nil {
		return nil, err
	}

	return statistics, err
}

func (dor *doukeOrderRepo) StatisticsByPaySuccessTime(ctx context.Context, userId, productId uint64, flowPoint, startTime, endTime string) (*v1.StatisticsDoukeOrderByPaySuccessTimesReply, error) {
	statistics, err := dor.data.douyinuc.StatisticsDoukeOrderByPaySuccessTimes(ctx, &v1.StatisticsDoukeOrderByPaySuccessTimesRequest{
		UserId:    userId,
		ProductId: productId,
		FlowPoint: flowPoint,
		StartTime: startTime,
		EndTime:   endTime,
	})

	if err != nil {
		return nil, err
	}

	return statistics, err
}

package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	v1 "interface/api/interface/v1"
)

func (is *InterfaceService) StatisticsUserDoukeOrders(ctx context.Context, in *empty.Empty) (*v1.StatisticsUserDoukeOrdersReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	statistics, err := is.douc.StatisticsUserDoukeOrders(ctx, userInfo.Data.UserId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsUserDoukeOrdersReply_Statistics, 0)

	for _, statistic := range statistics.Data.Statistics {
		list = append(list, &v1.StatisticsUserDoukeOrdersReply_Statistics{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsUserDoukeOrdersReply{
		Code: 200,
		Data: &v1.StatisticsUserDoukeOrdersReply_Data{
			Statistics: list,
		},
	}, nil
}

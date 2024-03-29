package service

import (
	"context"
	v1 "douyin/api/douyin/v1"
	"douyin/internal/pkg/tool"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (ds *DouyinService) ListUserIdDoukeOrders(ctx context.Context, in *emptypb.Empty) (*v1.ListUserIdDoukeOrdersReply, error) {
	doukeOrders, err := ds.douc.ListUserIdDoukeOrders(ctx)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListUserIdDoukeOrdersReply_UserId, 0)

	for _, doukeOrder := range doukeOrders {
		list = append(list, &v1.ListUserIdDoukeOrdersReply_UserId{
			UserId: doukeOrder.UserId,
		})
	}

	return &v1.ListUserIdDoukeOrdersReply{
		Code: 200,
		Data: &v1.ListUserIdDoukeOrdersReply_Data{
			List: list,
		},
	}, nil
}

func (ds *DouyinService) StatisticsDoukeOrders(ctx context.Context, in *v1.StatisticsDoukeOrdersRequest) (*v1.StatisticsDoukeOrdersReply, error) {
	statistics, err := ds.douc.StatisticsDoukeOrders(ctx, in.UserId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsDoukeOrdersReply_Statistics, 0)

	for _, statistic := range statistics.Statistics {
		list = append(list, &v1.StatisticsDoukeOrdersReply_Statistics{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsDoukeOrdersReply{
		Code: 200,
		Data: &v1.StatisticsDoukeOrdersReply_Data{
			Statistics: list,
		},
	}, nil
}

func (ds *DouyinService) StatisticsDoukeOrderByDays(ctx context.Context, in *v1.StatisticsDoukeOrderByDaysRequest) (*v1.StatisticsDoukeOrderByDaysReply, error) {
	statistics, err := ds.douc.StatisticsDoukeOrderByDays(ctx, in.UserId, in.Day)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsDoukeOrderByDaysReply_Statistics, 0)

	for _, statistic := range statistics.Statistics {
		list = append(list, &v1.StatisticsDoukeOrderByDaysReply_Statistics{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsDoukeOrderByDaysReply{
		Code: 200,
		Data: &v1.StatisticsDoukeOrderByDaysReply_Data{
			Statistics: list,
		},
	}, nil
}

func (ds *DouyinService) AsyncNotificationDoukeOrders(ctx context.Context, in *v1.AsyncNotificationDoukeOrdersRequest) (*v1.AsyncNotificationDoukeOrdersReply, error) {
	ds.log.Infof("同步抖客订单详情数据, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ds.douc.AsyncNotificationDoukeOrders(ctx, in.EventSign, in.AppId, in.Content); err != nil {
		return nil, err
	}

	ds.log.Infof("同步抖客订单详情数据, 结束时间 %s \n", time.Now())

	return &v1.AsyncNotificationDoukeOrdersReply{
		Code: 200,
		Data: &v1.AsyncNotificationDoukeOrdersReply_Data{},
	}, nil
}

func (ds *DouyinService) GetCompanyTaskUserOrderStatus(ctx context.Context, in *v1.GetCompanyTaskUserOrderStatusRequest) (*v1.GetCompanyTaskUserOrderStatusReply, error) {
	order, err := ds.douc.GetCompanyTaskUserOrderStatus(ctx, in.UserId, in.ProductId, in.FlowPoint, in.CreateTime)

	if err != nil {
		return nil, err
	}

	return &v1.GetCompanyTaskUserOrderStatusReply{
		Code: 200,
		Data: &v1.GetCompanyTaskUserOrderStatusReply_Data{
			Id:             order.Id,
			UserId:         order.UserId,
			ProductId:      order.ProductId,
			PaySuccessTime: tool.TimeToString("2006-01-02 15:04:05", order.PaySuccessTime),
			FlowPoint:      order.FlowPoint,
		},
	}, nil
}

package service

import (
	"context"
	v1 "douyin/api/douyin/v1"
	"douyin/internal/pkg/tool"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

func (ds *DouyinService) GetDoukeOrders(ctx context.Context, in *v1.GetDoukeOrdersRequest) (*v1.GetDoukeOrdersReply, error) {
	doukeOrder, err := ds.douc.GetDoukeOrders(ctx, in.UserId, in.ProductId, in.FlowPoint, in.CreateTime)

	if err != nil {
		return nil, err
	}

	settleTime := ""

	if doukeOrder.SettleTime != nil {
		settleTime = tool.TimeToString("2006-01-02 15:04:05", *doukeOrder.SettleTime)
	}

	return &v1.GetDoukeOrdersReply{
		Code: 200,
		Data: &v1.GetDoukeOrdersReply_Data{
			UserId:              doukeOrder.UserId,
			OrderId:             doukeOrder.OrderId,
			ProductId:           doukeOrder.ProductId,
			ProductName:         doukeOrder.ProductName,
			ProductImg:          doukeOrder.ProductImg,
			PaySuccessTime:      tool.TimeToString("2006-01-02 15:04:05", doukeOrder.PaySuccessTime),
			SettleTime:          settleTime,
			TotalPayAmount:      tool.Decimal(float64(doukeOrder.TotalPayAmount), 2),
			PayGoodsAmount:      tool.Decimal(float64(doukeOrder.PayGoodsAmount), 2),
			FlowPoint:           doukeOrder.FlowPoint,
			EstimatedCommission: tool.Decimal(float64(doukeOrder.EstimatedCommission), 2),
			RealCommission:      tool.Decimal(float64(doukeOrder.RealCommission), 2),
			ItemNum:             doukeOrder.ItemNum,
			CreateTime:          tool.TimeToString("2006-01-02 15:04:05", doukeOrder.CreateTime),
			UpdateTime:          tool.TimeToString("2006-01-02 15:04:05", doukeOrder.UpdateTime),
		},
	}, nil
}

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

func (ds *DouyinService) StatisticsDoukeOrderByPaySuccessTimes(ctx context.Context, in *v1.StatisticsDoukeOrderByPaySuccessTimesRequest) (*v1.StatisticsDoukeOrderByPaySuccessTimesReply, error) {
	statistics, err := ds.douc.StatisticsDoukeOrderByPaySuccessTimes(ctx, in.UserId, in.ProductId, in.FlowPoint, in.StartTime, in.EndTime)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsDoukeOrderByPaySuccessTimesReply_Statistics, 0)

	for _, statistic := range statistics.Statistics {
		list = append(list, &v1.StatisticsDoukeOrderByPaySuccessTimesReply_Statistics{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsDoukeOrderByPaySuccessTimesReply{
		Code: 200,
		Data: &v1.StatisticsDoukeOrderByPaySuccessTimesReply_Data{
			Statistics: list,
		},
	}, nil
}

/*func (ds *DouyinService) AsyncNotificationDoukeOrders(ctx context.Context, in *v1.AsyncNotificationDoukeOrdersRequest) (*v1.AsyncNotificationDoukeOrdersReply, error) {
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
*/

func (ds *DouyinService) SyncDoukeOrders(ctx context.Context, in *v1.SyncDoukeOrdersRequest) (*v1.SyncDoukeOrdersReply, error) {
	ds.log.Infof("同步抖客订单数据, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ds.douc.SyncDoukeOrders(ctx, in.Day); err != nil {
		return nil, err
	}

	ds.log.Infof("同步抖客订单数据, 结束时间 %s \n", time.Now())

	return &v1.SyncDoukeOrdersReply{
		Code: 200,
		Data: &v1.SyncDoukeOrdersReply_Data{},
	}, nil
}

func (ds *DouyinService) OperationDoukeOrders(ctx context.Context, in *empty.Empty) (*v1.OperationDoukeOrdersReply, error) {
	ds.log.Infof("抖客订单佣金数据, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ds.douc.OperationDoukeOrders(ctx); err != nil {
		return nil, err
	}

	ds.log.Infof("抖客订单佣金数据, 结束时间 %s \n", time.Now())

	return &v1.OperationDoukeOrdersReply{
		Code: 200,
		Data: &v1.OperationDoukeOrdersReply_Data{},
	}, nil
}

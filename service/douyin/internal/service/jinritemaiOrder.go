package service

import (
	"context"
	v1 "douyin/api/douyin/v1"
	"douyin/internal/pkg/tool"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"
	"math"
	"strconv"
	"time"
)

func (ds *DouyinService) GetStorePreferenceJinritemaiOrders(ctx context.Context, in *v1.GetStorePreferenceJinritemaiOrdersRequest) (*v1.GetStorePreferenceJinritemaiOrdersReply, error) {
	storePreferences, err := ds.jouc.GetStorePreferenceJinritemaiOrders(ctx, in.UserId)

	if err != nil {
		return nil, err
	}

	storePreference := make([]*v1.GetStorePreferenceJinritemaiOrdersReply_StorePreference, 0)

	for _, lstorePreference := range storePreferences {
		storePreference = append(storePreference, &v1.GetStorePreferenceJinritemaiOrdersReply_StorePreference{
			IndustryId:    lstorePreference.IndustryId,
			IndustryName:  lstorePreference.IndustryName,
			IndustryRatio: lstorePreference.IndustryRatio,
		})
	}

	return &v1.GetStorePreferenceJinritemaiOrdersReply{
		Code: 200,
		Data: &v1.GetStorePreferenceJinritemaiOrdersReply_Data{
			StorePreference: storePreference,
		},
	}, nil
}

func (ds *DouyinService) GetIsTopJinritemaiOrders(ctx context.Context, in *v1.GetIsTopJinritemaiOrdersRequest) (*v1.GetIsTopJinritemaiOrdersReply, error) {
	var isTop uint32 = 0

	if _, err := ds.jouc.GetIsTopJinritemaiOrders(ctx, in.ProductId); err == nil {
		isTop = 1
	}

	return &v1.GetIsTopJinritemaiOrdersReply{
		Code: 200,
		Data: &v1.GetIsTopJinritemaiOrdersReply_Data{
			IsTop: isTop,
		},
	}, nil
}

func (ds *DouyinService) ListJinritemaiOrders(ctx context.Context, in *v1.ListJinritemaiOrdersRequest) (*v1.ListJinritemaiOrdersReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jinritemaiOrders, err := ds.jouc.ListJinritemaiOrders(ctx, in.PageNum, in.PageSize, in.UserId, in.StartDay, in.EndDay)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListJinritemaiOrdersReply_JinritemaiOrders, 0)

	for _, jinritemaiOrder := range jinritemaiOrders.List {
		list = append(list, &v1.ListJinritemaiOrdersReply_JinritemaiOrders{
			ProductId:          jinritemaiOrder.ProductId,
			ProductName:        jinritemaiOrder.ProductName,
			ProductImg:         jinritemaiOrder.ProductImg,
			TotalPayAmount:     tool.Decimal(float64(jinritemaiOrder.TotalPayAmount), 2),
			RealCommission:     tool.Decimal(float64(jinritemaiOrder.RealCommission), 2),
			RealCommissionRate: jinritemaiOrder.RealCommissionRate,
			ItemNum:            jinritemaiOrder.ItemNum,
			MediaType:          jinritemaiOrder.MediaType,
			MediaTypeName:      jinritemaiOrder.MediaTypeName,
			MediaId:            jinritemaiOrder.MediaId,
			MediaCover:         jinritemaiOrder.MediaCover,
			Avatar:             jinritemaiOrder.Avatar,
			IsShow:             uint32(jinritemaiOrder.IsShow),
		})
	}

	totalPage := uint64(math.Ceil(float64(jinritemaiOrders.Total) / float64(jinritemaiOrders.PageSize)))

	return &v1.ListJinritemaiOrdersReply{
		Code: 200,
		Data: &v1.ListJinritemaiOrdersReply_Data{
			PageNum:   jinritemaiOrders.PageNum,
			PageSize:  jinritemaiOrders.PageSize,
			Total:     jinritemaiOrders.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ds *DouyinService) ListProductIdAndVideoIdJinritemaiOrders(ctx context.Context, in *v1.ListProductIdAndVideoIdJinritemaiOrdersRequest) (*v1.ListProductIdAndVideoIdJinritemaiOrdersReply, error) {
	jinritemaiOrders, err := ds.jouc.ListProductIdAndVideoIdJinritemaiOrders(ctx, in.PageNum, in.PageSize)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListProductIdAndVideoIdJinritemaiOrdersReply_JinritemaiOrder, 0)

	for _, jinritemaiOrder := range jinritemaiOrders.List {
		productId, _ := strconv.ParseUint(jinritemaiOrder.ProductId, 10, 64)

		list = append(list, &v1.ListProductIdAndVideoIdJinritemaiOrdersReply_JinritemaiOrder{
			VideoId:      jinritemaiOrder.MediaId,
			ProductOutId: productId,
		})
	}

	totalPage := uint64(math.Ceil(float64(jinritemaiOrders.Total) / float64(jinritemaiOrders.PageSize)))

	return &v1.ListProductIdAndVideoIdJinritemaiOrdersReply{
		Code: 200,
		Data: &v1.ListProductIdAndVideoIdJinritemaiOrdersReply_Data{
			PageNum:   jinritemaiOrders.PageNum,
			PageSize:  jinritemaiOrders.PageSize,
			Total:     jinritemaiOrders.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ds *DouyinService) ListJinritemaiOrderByPickExtras(ctx context.Context, in *emptypb.Empty) (*v1.ListJinritemaiOrderByPickExtrasReply, error) {
	jinritemaiOrders, err := ds.jouc.ListJinritemaiOrderByPickExtras(ctx)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListJinritemaiOrderByPickExtrasReply_OpenDouyinUserInfo, 0)

	for _, jinritemaiOrder := range jinritemaiOrders {
		list = append(list, &v1.ListJinritemaiOrderByPickExtrasReply_OpenDouyinUserInfo{
			ClientKey: jinritemaiOrder.AuthorClientKey,
			OpenId:    jinritemaiOrder.AuthorOpenId,
		})
	}

	return &v1.ListJinritemaiOrderByPickExtrasReply{
		Code: 200,
		Data: &v1.ListJinritemaiOrderByPickExtrasReply_Data{
			List: list,
		},
	}, nil
}

func (ds *DouyinService) ListCommissionRateJinritemaiOrders(ctx context.Context, in *v1.ListCommissionRateJinritemaiOrdersRequest) (*v1.ListCommissionRateJinritemaiOrdersReply, error) {
	jinritemaiOrders, err := ds.jouc.ListCommissionRateJinritemaiOrders(ctx, in.Content)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCommissionRateJinritemaiOrdersReply_CommissionRate, 0)

	for _, jinritemaiOrder := range jinritemaiOrders {
		list = append(list, &v1.ListCommissionRateJinritemaiOrdersReply_CommissionRate{
			ProductId:       jinritemaiOrder.ProductId,
			MediaId:         jinritemaiOrder.MediaId,
			CommissionRatio: uint32(jinritemaiOrder.CommissionRate),
		})
	}

	return &v1.ListCommissionRateJinritemaiOrdersReply{
		Code: 200,
		Data: &v1.ListCommissionRateJinritemaiOrdersReply_Data{
			List: list,
		},
	}, nil
}

func (ds *DouyinService) StatisticsJinritemaiOrders(ctx context.Context, in *v1.StatisticsJinritemaiOrdersRequest) (*v1.StatisticsJinritemaiOrdersReply, error) {
	statistics, err := ds.jouc.StatisticsJinritemaiOrders(ctx, in.UserId, in.StartDay, in.EndDay)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsJinritemaiOrdersReply_Statistics, 0)

	for _, statistic := range statistics.Statistics {
		list = append(list, &v1.StatisticsJinritemaiOrdersReply_Statistics{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsJinritemaiOrdersReply{
		Code: 200,
		Data: &v1.StatisticsJinritemaiOrdersReply_Data{
			Statistics: list,
		},
	}, nil
}

func (ds *DouyinService) StatisticsJinritemaiOrderByClientKeyAndOpenIds(ctx context.Context, in *v1.StatisticsJinritemaiOrderByClientKeyAndOpenIdsRequest) (*v1.StatisticsJinritemaiOrderByClientKeyAndOpenIdsReply, error) {
	statistics, err := ds.jouc.StatisticsJinritemaiOrderByClientKeyAndOpenIds(ctx, in.ClientKey, in.OpenId, in.StartDay, in.EndDay)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsJinritemaiOrderByClientKeyAndOpenIdsReply_Statistics, 0)

	for _, statistic := range statistics.Statistics {
		list = append(list, &v1.StatisticsJinritemaiOrderByClientKeyAndOpenIdsReply_Statistics{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsJinritemaiOrderByClientKeyAndOpenIdsReply{
		Code: 200,
		Data: &v1.StatisticsJinritemaiOrderByClientKeyAndOpenIdsReply_Data{
			Statistics: list,
		},
	}, nil
}

func (ds *DouyinService) StatisticsJinritemaiOrderByDays(ctx context.Context, in *v1.StatisticsJinritemaiOrderByDaysRequest) (*v1.StatisticsJinritemaiOrderByDaysReply, error) {
	statistics, err := ds.jouc.StatisticsJinritemaiOrderByDays(ctx, in.Day, in.Content, in.PickExtra)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsJinritemaiOrderByDaysReply_Statistics, 0)

	for _, statistic := range statistics.Statistics {
		list = append(list, &v1.StatisticsJinritemaiOrderByDaysReply_Statistics{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsJinritemaiOrderByDaysReply{
		Code: 200,
		Data: &v1.StatisticsJinritemaiOrderByDaysReply_Data{
			Statistics: list,
		},
	}, nil
}

func (ds *DouyinService) AsyncNotificationJinritemaiOrders(ctx context.Context, in *v1.AsyncNotificationJinritemaiOrdersRequest) (*v1.AsyncNotificationJinritemaiOrdersReply, error) {
	ds.log.Infof("同步精选联盟达人详情数据, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ds.jouc.AsyncNotificationJinritemaiOrders(ctx, in.MsgId, in.Sign, in.Content); err != nil {
		return nil, err
	}

	ds.log.Infof("同步精选联盟达人详情数据, 结束时间 %s \n", time.Now())

	return &v1.AsyncNotificationJinritemaiOrdersReply{
		Code: 200,
		Data: &v1.AsyncNotificationJinritemaiOrdersReply_Data{},
	}, nil
}

func (ds *DouyinService) SyncJinritemaiOrders(ctx context.Context, in *v1.SyncJinritemaiOrdersRequest) (*v1.SyncJinritemaiOrdersReply, error) {
	ds.log.Infof("同步精选联盟达人订单数据, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ds.jouc.SyncJinritemaiOrders(ctx, in.Day); err != nil {
		return nil, err
	}

	ds.log.Infof("同步精选联盟达人订单数据, 结束时间 %s \n", time.Now())

	return &v1.SyncJinritemaiOrdersReply{
		Code: 200,
		Data: &v1.SyncJinritemaiOrdersReply_Data{},
	}, nil
}

func (ds *DouyinService) Sync90DayJinritemaiOrders(ctx context.Context, in *empty.Empty) (*v1.Sync90DayJinritemaiOrdersReply, error) {
	ds.log.Infof("同步新授权的精选联盟达人90天订单数据, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ds.jouc.Sync90DayJinritemaiOrders(ctx); err != nil {
		return nil, err
	}

	ds.log.Infof("同步新授权的精选联盟达人90天订单数据, 结束时间 %s \n", time.Now())

	return &v1.Sync90DayJinritemaiOrdersReply{
		Code: 200,
		Data: &v1.Sync90DayJinritemaiOrdersReply_Data{},
	}, nil
}

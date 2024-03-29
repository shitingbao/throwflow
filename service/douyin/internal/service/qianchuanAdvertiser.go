package service

import (
	"context"
	v1 "douyin/api/douyin/v1"
	"douyin/internal/pkg/tool"
	"google.golang.org/protobuf/types/known/emptypb"
	"math"
	"time"
)

func (ds *DouyinService) GetQianchuanAdvertiserByCompanyIdAndAdvertiserIds(ctx context.Context, in *v1.GetQianchuanAdvertiserByCompanyIdAndAdvertiserIdsRequest) (*v1.GetQianchuanAdvertiserByCompanyIdAndAdvertiserIdsReply, error) {
	qianchuanAdvertiser, err := ds.qauc.GetQianchuanAdvertiserByCompanyIdAndAdvertiserIds(ctx, in.CompanyId, in.AdvertiserId)

	if err != nil {
		return nil, err
	}

	return &v1.GetQianchuanAdvertiserByCompanyIdAndAdvertiserIdsReply{
		Code: 200,
		Data: &v1.GetQianchuanAdvertiserByCompanyIdAndAdvertiserIdsReply_Data{
			AdvertiserId:   qianchuanAdvertiser.AdvertiserId,
			CompanyId:      qianchuanAdvertiser.CompanyId,
			AccountId:      qianchuanAdvertiser.AccountId,
			AdvertiserName: qianchuanAdvertiser.AdvertiserName,
			CompanyName:    qianchuanAdvertiser.CompanyName,
			Status:         uint32(qianchuanAdvertiser.Status),
			CreateTime:     tool.TimeToString("2006-01-02 15:04:05", qianchuanAdvertiser.CreateTime),
			UpdateTime:     tool.TimeToString("2006-01-02 15:04:05", qianchuanAdvertiser.UpdateTime),
		},
	}, nil
}

func (ds *DouyinService) ListQianchuanAdvertisers(ctx context.Context, in *v1.ListQianchuanAdvertisersRequest) (*v1.ListQianchuanAdvertisersReply, error) {
	qianchuanAdvertisers, err := ds.qauc.ListQianchuanAdvertisers(ctx, in.PageNum, in.PageSize, in.CompanyId, in.Keyword, in.AdvertiserIds, in.Status)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListQianchuanAdvertisersReply_Advertisers, 0)

	for _, qianchuanAdvertiser := range qianchuanAdvertisers.List {
		list = append(list, &v1.ListQianchuanAdvertisersReply_Advertisers{
			AdvertiserId:     qianchuanAdvertiser.AdvertiserId,
			CompanyId:        qianchuanAdvertiser.CompanyId,
			AccountId:        qianchuanAdvertiser.AccountId,
			AdvertiserName:   qianchuanAdvertiser.AdvertiserName,
			CompanyName:      qianchuanAdvertiser.CompanyName,
			Status:           uint32(qianchuanAdvertiser.Status),
			OtherCompanyName: qianchuanAdvertiser.OtherCompanyName,
			OtherCompanyId:   qianchuanAdvertiser.OtherCompanyId,
			CreateTime:       tool.TimeToString("2006-01-02 15:04:05", qianchuanAdvertiser.CreateTime),
			UpdateTime:       tool.TimeToString("2006-01-02 15:04:05", qianchuanAdvertiser.UpdateTime),
		})
	}

	totalPage := uint64(math.Ceil(float64(qianchuanAdvertisers.Total) / float64(qianchuanAdvertisers.PageSize)))

	return &v1.ListQianchuanAdvertisersReply{
		Code: 200,
		Data: &v1.ListQianchuanAdvertisersReply_Data{
			PageNum:   qianchuanAdvertisers.PageNum,
			PageSize:  qianchuanAdvertisers.PageSize,
			Total:     qianchuanAdvertisers.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ds *DouyinService) ListQianchuanAdvertiserByDays(ctx context.Context, in *v1.ListQianchuanAdvertiserByDaysRequest) (*v1.ListQianchuanAdvertiserByDaysReply, error) {
	qianchuanAdvertiserstatuses, err := ds.qauc.ListQianchuanAdvertiserByDays(ctx, in.CompanyId, in.Day)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListQianchuanAdvertiserByDaysReply_Advertisers, 0)

	for _, qianchuanAdvertiserstatus := range qianchuanAdvertiserstatuses {
		list = append(list, &v1.ListQianchuanAdvertiserByDaysReply_Advertisers{
			AdvertiserId: qianchuanAdvertiserstatus.AdvertiserId,
			Status:       uint32(qianchuanAdvertiserstatus.Status),
		})
	}

	return &v1.ListQianchuanAdvertiserByDaysReply{
		Code: 200,
		Data: &v1.ListQianchuanAdvertiserByDaysReply_Data{
			List: list,
		},
	}, nil
}

func (ds *DouyinService) ListExternalQianchuanAdvertisers(ctx context.Context, in *v1.ListExternalQianchuanAdvertisersRequest) (*v1.ListExternalQianchuanAdvertisersReply, error) {
	qianchuanAdvertisers, err := ds.qauc.ListExternalQianchuanAdvertisers(ctx, in.PageNum, in.PageSize, in.StartDay, in.EndDay, in.AdvertiserIds)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListExternalQianchuanAdvertisersReply_Advertisers, 0)

	for _, qianchuanAdvertiser := range qianchuanAdvertisers.List {
		list = append(list, &v1.ListExternalQianchuanAdvertisersReply_Advertisers{
			AdvertiserId:            qianchuanAdvertiser.AdvertiserId,
			AdvertiserName:          qianchuanAdvertiser.AdvertiserName,
			GeneralTotalBalance:     qianchuanAdvertiser.GeneralTotalBalance,
			StatCost:                qianchuanAdvertiser.StatCost,
			Roi:                     qianchuanAdvertiser.Roi,
			Campaigns:               int64(qianchuanAdvertiser.Campaigns),
			PayOrderCount:           qianchuanAdvertiser.PayOrderCount,
			PayOrderAmount:          qianchuanAdvertiser.PayOrderAmount,
			ClickCnt:                qianchuanAdvertiser.ClickCnt,
			ClickRate:               qianchuanAdvertiser.ClickRate,
			PayConvertRate:          qianchuanAdvertiser.PayConvertRate,
			AveragePayOrderStatCost: qianchuanAdvertiser.AveragePayOrderStatCost,
			PayOrderAveragePrice:    qianchuanAdvertiser.PayOrderAveragePrice,
			ShowCnt:                 qianchuanAdvertiser.ShowCnt,
			DyFollow:                qianchuanAdvertiser.DyFollow,
		})
	}

	totalPage := uint64(math.Ceil(float64(qianchuanAdvertisers.Total) / float64(qianchuanAdvertisers.PageSize)))

	return &v1.ListExternalQianchuanAdvertisersReply{
		Code: 200,
		Data: &v1.ListExternalQianchuanAdvertisersReply_Data{
			PageNum:   qianchuanAdvertisers.PageNum,
			PageSize:  qianchuanAdvertisers.PageSize,
			Total:     qianchuanAdvertisers.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ds *DouyinService) StatisticsQianchuanAdvertisers(ctx context.Context, in *v1.StatisticsQianchuanAdvertisersRequest) (*v1.StatisticsQianchuanAdvertisersReply, error) {
	statistics, err := ds.qauc.StatisticsQianchuanAdvertisers(ctx, in.CompanyId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsQianchuanAdvertisersReply_Statistics, 0)

	for _, statistic := range statistics.Statistics {
		list = append(list, &v1.StatisticsQianchuanAdvertisersReply_Statistics{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsQianchuanAdvertisersReply{
		Code: 200,
		Data: &v1.StatisticsQianchuanAdvertisersReply_Data{
			Statistics: list,
		},
	}, nil
}

func (ds *DouyinService) StatisticsDashboardQianchuanAdvertisers(ctx context.Context, in *v1.StatisticsDashboardQianchuanAdvertisersRequest) (*v1.StatisticsDashboardQianchuanAdvertisersReply, error) {
	statistics, err := ds.qauc.StatisticsDashboardQianchuanAdvertisers(ctx, in.CompanyId, in.Day, in.AdvertiserIds)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsDashboardQianchuanAdvertisersReply_Statistics, 0)

	for _, statistic := range statistics.Statistics {
		list = append(list, &v1.StatisticsDashboardQianchuanAdvertisersReply_Statistics{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsDashboardQianchuanAdvertisersReply{
		Code: 200,
		Data: &v1.StatisticsDashboardQianchuanAdvertisersReply_Data{
			Statistics: list,
		},
	}, nil
}

func (ds *DouyinService) StatisticsExternalQianchuanAdvertisers(ctx context.Context, in *v1.StatisticsExternalQianchuanAdvertisersRequest) (*v1.StatisticsExternalQianchuanAdvertisersReply, error) {
	statistics, err := ds.qauc.StatisticsExternalQianchuanAdvertisers(ctx, in.StartDay, in.EndDay, in.AdvertiserIds)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsExternalQianchuanAdvertisersReply_Statistics, 0)

	for _, statistic := range statistics.Statistics {
		list = append(list, &v1.StatisticsExternalQianchuanAdvertisersReply_Statistics{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsExternalQianchuanAdvertisersReply{
		Code: 200,
		Data: &v1.StatisticsExternalQianchuanAdvertisersReply_Data{
			Statistics: list,
		},
	}, nil
}

func (ds *DouyinService) UpdateStatusQianchuanAdvertisers(ctx context.Context, in *v1.UpdateStatusQianchuanAdvertisersRequest) (*v1.UpdateStatusQianchuanAdvertisersReply, error) {
	qianchuanAdvertiser, err := ds.qauc.UpdateStatusQianchuanAdvertisers(ctx, in.CompanyId, in.AdvertiserId, uint8(in.Status))

	if err != nil {
		return nil, err
	}

	return &v1.UpdateStatusQianchuanAdvertisersReply{
		Code: 200,
		Data: &v1.UpdateStatusQianchuanAdvertisersReply_Data{
			CompanyId:      qianchuanAdvertiser.CompanyId,
			AdvertiserId:   qianchuanAdvertiser.AdvertiserId,
			AdvertiserName: qianchuanAdvertiser.AdvertiserName,
			CompanyName:    qianchuanAdvertiser.CompanyName,
			Status:         uint32(qianchuanAdvertiser.Status),
			CreateTime:     tool.TimeToString("2006-01-02 15:04:05", qianchuanAdvertiser.CreateTime),
			UpdateTime:     tool.TimeToString("2006-01-02 15:04:05", qianchuanAdvertiser.UpdateTime),
		},
	}, nil
}

func (ds *DouyinService) UpdateStatusQianchuanAdvertisersByCompanyId(ctx context.Context, in *v1.UpdateStatusQianchuanAdvertisersByCompanyIdRequest) (*v1.UpdateStatusQianchuanAdvertisersByCompanyIdReply, error) {
	err := ds.qauc.UpdateStatusQianchuanAdvertisersByCompanyId(ctx, in.CompanyId, uint8(in.Status))

	if err != nil {
		return nil, err
	}

	return &v1.UpdateStatusQianchuanAdvertisersByCompanyIdReply{
		Code: 200,
		Data: &v1.UpdateStatusQianchuanAdvertisersByCompanyIdReply_Data{},
	}, nil
}

func (ds *DouyinService) SyncQianchuanDatas(ctx context.Context, in *emptypb.Empty) (*v1.SyncQianchuanDatasReply, error) {
	ds.log.Infof("同步千川数据, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ds.qauc.SyncQianchuanDatas(ctx); err != nil {
		return nil, err
	}

	ds.log.Infof("同步千川数据, 结束时间 %s \n", time.Now())

	return &v1.SyncQianchuanDatasReply{
		Code: 200,
		Data: &v1.SyncQianchuanDatasReply_Data{},
	}, nil
}

func (ds *DouyinService) SyncRdsDatas(ctx context.Context, in *emptypb.Empty) (*v1.SyncRdsDatasReply, error) {
	ds.log.Infof("同步RDS数据, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ds.qauc.SyncRdsDatas(ctx); err != nil {
		return nil, err
	}

	ds.log.Infof("同步RDS数据, 结束时间 %s \n", time.Now())

	return &v1.SyncRdsDatasReply{
		Code: 200,
		Data: &v1.SyncRdsDatasReply_Data{},
	}, nil
}

package service

import (
	"context"
	v1 "douyin/api/douyin/v1"
	"douyin/internal/pkg/tool"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"math"
	"time"
)

func (ds *DouyinService) GetExternalQianchuanAds(ctx context.Context, in *v1.GetExternalQianchuanAdsRequest) (*v1.GetExternalQianchuanAdsReply, error) {
	qianchuanAds, err := ds.qaduc.GetExternalQianchuanAds(ctx, in.AdId, in.Day)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.GetExternalQianchuanAdsReply_Ads, 0)

	for _, qianchuanAd := range qianchuanAds {
		list = append(list, &v1.GetExternalQianchuanAdsReply_Ads{
			Time:                    tool.TimeToString("15:04", time.Unix(qianchuanAd.Time, 0)),
			StatCost:                qianchuanAd.StatCost,
			Roi:                     qianchuanAd.Roi,
			PayOrderCount:           qianchuanAd.PayOrderCount,
			PayOrderAmount:          qianchuanAd.PayOrderAmount,
			ClickCnt:                qianchuanAd.ClickCnt,
			ShowCnt:                 qianchuanAd.ShowCnt,
			ConvertCnt:              qianchuanAd.ConvertCnt,
			ClickRate:               qianchuanAd.ClickRate,
			CpmPlatform:             qianchuanAd.CpmPlatform,
			DyFollow:                qianchuanAd.DyFollow,
			PayConvertRate:          qianchuanAd.PayConvertRate,
			ConvertCost:             qianchuanAd.ConvertCost,
			AveragePayOrderStatCost: qianchuanAd.AveragePayOrderStatCost,
			PayOrderAveragePrice:    qianchuanAd.PayOrderAveragePrice,
		})
	}

	return &v1.GetExternalQianchuanAdsReply{
		Code: 200,
		Data: &v1.GetExternalQianchuanAdsReply_Data{
			List: list,
		},
	}, nil
}

func (ds *DouyinService) GetExternalHistoryQianchuanAds(ctx context.Context, in *v1.GetExternalHistoryQianchuanAdsRequest) (*v1.GetExternalHistoryQianchuanAdsReply, error) {
	qianchuanAds, err := ds.qaduc.GetExternalHistoryQianchuanAds(ctx, in.AdId, in.StartDay, in.EndDay)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.GetExternalHistoryQianchuanAdsReply_Ads, 0)

	for _, qianchuanAd := range qianchuanAds {
		list = append(list, &v1.GetExternalHistoryQianchuanAdsReply_Ads{
			Time:                    tool.TimeToString("2006-01-02", time.Unix(qianchuanAd.Time, 0)),
			StatCost:                qianchuanAd.StatCost,
			Roi:                     qianchuanAd.Roi,
			PayOrderCount:           qianchuanAd.PayOrderCount,
			PayOrderAmount:          qianchuanAd.PayOrderAmount,
			ClickCnt:                qianchuanAd.ClickCnt,
			ShowCnt:                 qianchuanAd.ShowCnt,
			ConvertCnt:              qianchuanAd.ConvertCnt,
			ClickRate:               qianchuanAd.ClickRate,
			CpmPlatform:             qianchuanAd.CpmPlatform,
			DyFollow:                qianchuanAd.DyFollow,
			PayConvertRate:          qianchuanAd.PayConvertRate,
			ConvertCost:             qianchuanAd.ConvertCost,
			AveragePayOrderStatCost: qianchuanAd.AveragePayOrderStatCost,
			PayOrderAveragePrice:    qianchuanAd.PayOrderAveragePrice,
		})
	}

	return &v1.GetExternalHistoryQianchuanAdsReply{
		Code: 200,
		Data: &v1.GetExternalHistoryQianchuanAdsReply_Data{
			List: list,
		},
	}, nil
}

func (ds *DouyinService) ListQianchuanAds(ctx context.Context, in *v1.ListQianchuanAdsRequest) (*v1.ListQianchuanAdsReply, error) {
	qianchuanAds, err := ds.qaduc.ListQianchuanAds(ctx, in.PageNum, in.PageSize, in.Day, in.Keyword, in.AdvertiserIds)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListQianchuanAdsReply_QianchuanAds, 0)

	for _, qianchuanAd := range qianchuanAds.List {
		list = append(list, &v1.ListQianchuanAdsReply_QianchuanAds{
			AdId:           qianchuanAd.AdId,
			AdvertiserId:   qianchuanAd.AdvertiserId,
			CampaignId:     qianchuanAd.CampaignId,
			PromotionWay:   qianchuanAd.PromotionWay,
			MarketingGoal:  qianchuanAd.MarketingGoal,
			MarketingScene: qianchuanAd.MarketingScene,
			Name:           qianchuanAd.Name,
			Status:         qianchuanAd.Status,
			OptStatus:      qianchuanAd.OptStatus,
			AdCreateTime:   qianchuanAd.AdCreateTime,
			AdModifyTime:   qianchuanAd.AdModifyTime,
			LabAdType:      qianchuanAd.LabAdType,
		})
	}

	totalPage := uint64(math.Ceil(float64(qianchuanAds.Total) / float64(qianchuanAds.PageSize)))

	return &v1.ListQianchuanAdsReply{
		Code: 200,
		Data: &v1.ListQianchuanAdsReply_Data{
			PageNum:   qianchuanAds.PageNum,
			PageSize:  qianchuanAds.PageSize,
			Total:     qianchuanAds.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ds *DouyinService) ListQianchuanReportAdvertisers(ctx context.Context, in *v1.ListQianchuanReportAdvertisersRequest) (*v1.ListQianchuanReportAdvertisersReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	qianchuanReportAdvertisers, err := ds.qaduc.ListQianchuanReportAdvertisers(ctx, in.Day, in.AdvertiserIds)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListQianchuanReportAdvertisersReply_QianchuanReportAdvertisers, 0)

	for _, qianchuanReportAdvertiser := range qianchuanReportAdvertisers {
		list = append(list, &v1.ListQianchuanReportAdvertisersReply_QianchuanReportAdvertisers{
			Time:           tool.TimeToString("15:04", time.Unix(qianchuanReportAdvertiser.Time, 0)),
			StatCost:       fmt.Sprintf("%.2f", qianchuanReportAdvertiser.StatCost),
			PayOrderAmount: fmt.Sprintf("%.2f", qianchuanReportAdvertiser.PayOrderAmount),
		})
	}

	return &v1.ListQianchuanReportAdvertisersReply{
		Code: 200,
		Data: &v1.ListQianchuanReportAdvertisersReply_Data{
			List: list,
		},
	}, nil
}

func (ds *DouyinService) ListExternalQianchuanAds(ctx context.Context, in *v1.ListExternalQianchuanAdsRequest) (*v1.ListExternalQianchuanAdsReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	qianchuanAds, err := ds.qaduc.ListExternalQianchuanAds(ctx, in.PageNum, in.PageSize, in.StartDay, in.EndDay, in.Keyword, in.AdvertiserIds, in.Filter, in.OrderName, in.OrderType)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListExternalQianchuanAdsReply_Ads, 0)

	for _, qianchuanAd := range qianchuanAds.List {
		var campaignId uint64
		var campaignName string

		if qianchuanAd.LabAdType == "NOT_LAB_AD" {
			campaignId = qianchuanAd.CampaignId
			campaignName = qianchuanAd.CampaignName
		}

		list = append(list, &v1.ListExternalQianchuanAdsReply_Ads{
			AdId:                    qianchuanAd.AdId,
			AdName:                  qianchuanAd.AdName,
			AdvertiserId:            qianchuanAd.AdvertiserId,
			AdvertiserName:          qianchuanAd.AdvertiserName,
			CampaignId:              campaignId,
			CampaignName:            campaignName,
			LabAdType:               qianchuanAd.LabAdType,
			LabAdTypeName:           qianchuanAd.LabAdTypeName,
			MarketingGoal:           qianchuanAd.MarketingGoal,
			MarketingGoalName:       qianchuanAd.MarketingGoalName,
			Status:                  qianchuanAd.Status,
			StatusName:              qianchuanAd.StatusName,
			OptStatus:               qianchuanAd.OptStatus,
			OptStatusName:           qianchuanAd.OptStatusName,
			ExternalAction:          qianchuanAd.ExternalAction,
			ExternalActionName:      qianchuanAd.ExternalActionName,
			DeepExternalAction:      qianchuanAd.DeepExternalAction,
			DeepBidType:             qianchuanAd.DeepBidType,
			PromotionId:             qianchuanAd.PromotionId,
			PromotionShowId:         qianchuanAd.PromotionShowId,
			PromotionName:           qianchuanAd.PromotionName,
			PromotionImg:            qianchuanAd.PromotionImg,
			PromotionAvatar:         qianchuanAd.PromotionAvatar,
			PromotionType:           qianchuanAd.PromotionType,
			StatCost:                qianchuanAd.StatCost,
			Roi:                     qianchuanAd.Roi,
			CpaBid:                  qianchuanAd.CpaBid,
			RoiGoal:                 qianchuanAd.RoiGoal,
			Budget:                  qianchuanAd.Budget,
			BudgetMode:              qianchuanAd.BudgetMode,
			BudgetModeName:          qianchuanAd.BudgetModeName,
			PayOrderCount:           qianchuanAd.PayOrderCount,
			PayOrderAmount:          qianchuanAd.PayOrderAmount,
			ClickCnt:                qianchuanAd.ClickCnt,
			ShowCnt:                 qianchuanAd.ShowCnt,
			ConvertCnt:              qianchuanAd.ConvertCnt,
			ClickRate:               qianchuanAd.ClickRate,
			CpmPlatform:             qianchuanAd.CpmPlatform,
			DyFollow:                qianchuanAd.DyFollow,
			PayConvertRate:          qianchuanAd.PayConvertRate,
			ConvertCost:             qianchuanAd.ConvertCost,
			AveragePayOrderStatCost: qianchuanAd.AveragePayOrderStatCost,
			PayOrderAveragePrice:    qianchuanAd.PayOrderAveragePrice,
			ConvertRate:             qianchuanAd.ConvertRate,
			AdCreateTime:            tool.TimeToString("2006-01-02 15:04:05", qianchuanAd.AdCreateTime),
			AdModifyTime:            tool.TimeToString("2006-01-02 15:04:05", qianchuanAd.AdModifyTime),
		})
	}

	totalPage := uint64(math.Ceil(float64(qianchuanAds.Total) / float64(qianchuanAds.PageSize)))

	return &v1.ListExternalQianchuanAdsReply{
		Code: 200,
		Data: &v1.ListExternalQianchuanAdsReply_Data{
			PageNum:   qianchuanAds.PageNum,
			PageSize:  qianchuanAds.PageSize,
			Total:     qianchuanAds.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ds *DouyinService) StatisticsExternalQianchuanAds(ctx context.Context, in *v1.StatisticsExternalQianchuanAdsRequest) (*v1.StatisticsExternalQianchuanAdsReply, error) {
	statistics, err := ds.qaduc.StatisticsExternalQianchuanAds(ctx, in.StartDay, in.EndDay, in.Keyword, in.AdvertiserIds, in.Filter)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsExternalQianchuanAdsReply_Statistics, 0)

	for _, statistic := range statistics.Statistics {
		list = append(list, &v1.StatisticsExternalQianchuanAdsReply_Statistics{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsExternalQianchuanAdsReply{
		Code: 200,
		Data: &v1.StatisticsExternalQianchuanAdsReply_Data{
			Statistics: list,
		},
	}, nil
}

func (ds *DouyinService) ListSelectExternalQianchuanAds(ctx context.Context, in *empty.Empty) (*v1.ListSelectExternalQianchuanAdsReply, error) {
	selects, err := ds.qaduc.ListSelectExternalQianchuanAds(ctx)

	if err != nil {
		return nil, err
	}

	filter := make([]*v1.ListSelectExternalQianchuanAdsReply_Filter, 0)

	for _, lfilter := range selects.Filter {
		filter = append(filter, &v1.ListSelectExternalQianchuanAdsReply_Filter{
			Key:   lfilter.Key,
			Value: lfilter.Value,
		})
	}

	return &v1.ListSelectExternalQianchuanAdsReply{
		Code: 200,
		Data: &v1.ListSelectExternalQianchuanAdsReply_Data{
			Filter: filter,
		},
	}, nil
}

func (ds *DouyinService) SyncQianchuanAds(ctx context.Context, in *v1.SyncQianchuanAdsRequest) (*v1.SyncQianchuanAdsReply, error) {
	ds.log.Infof("同步千川计划数据, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ds.qaduc.SyncQianchuanAds(ctx, in.Day); err != nil {
		return nil, err
	}

	ds.log.Infof("同步千川计划数据, 结束时间 %s \n", time.Now())

	return &v1.SyncQianchuanAdsReply{
		Code: 200,
		Data: &v1.SyncQianchuanAdsReply_Data{},
	}, nil
}

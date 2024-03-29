package service

import (
	"context"
	v1 "douyin/api/douyin/v1"
	"fmt"
	"math"
	"time"
)

func (ds *DouyinService) ListQianchuanReportAwemes(ctx context.Context, in *v1.ListQianchuanReportAwemesRequest) (*v1.ListQianchuanReportAwemesReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	qianchuanReportAwemes, err := ds.qrauc.ListQianchuanReportAwemes(ctx, in.PageNum, in.PageSize, uint8(in.IsDistinction), in.Day, in.Keyword, in.AdvertiserIds)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListQianchuanReportAwemesReply_Awemes, 0)

	for _, qianchuanReportAweme := range qianchuanReportAwemes.List {
		list = append(list, &v1.ListQianchuanReportAwemesReply_Awemes{
			AdvertiserId:            qianchuanReportAweme.AdvertiserId,
			AdvertiserName:          qianchuanReportAweme.AdvertiserName,
			AwemeId:                 qianchuanReportAweme.AwemeId,
			AwemeName:               qianchuanReportAweme.AwemeName,
			AwemeShowId:             qianchuanReportAweme.AwemeShowId,
			AwemeAvatar:             qianchuanReportAweme.AwemeAvatar,
			AwemeUrl:                fmt.Sprintf("https://live.douyin.com/%d", qianchuanReportAweme.AwemeId),
			DyFollow:                qianchuanReportAweme.DyFollow,
			StatCost:                fmt.Sprintf("%.2f", qianchuanReportAweme.StatCost),
			PayOrderCount:           qianchuanReportAweme.PayOrderCount,
			PayOrderAmount:          fmt.Sprintf("%.2f", qianchuanReportAweme.PayOrderAmount),
			PayOrderAveragePrice:    fmt.Sprintf("%.2f¥", qianchuanReportAweme.PayOrderAveragePrice),
			Roi:                     fmt.Sprintf("%.2f", qianchuanReportAweme.Roi),
			ShowCnt:                 qianchuanReportAweme.ShowCnt,
			ClickCnt:                qianchuanReportAweme.ClickCnt,
			ClickRate:               fmt.Sprintf("%.2f%%", qianchuanReportAweme.ClickRate*100),
			ConvertCnt:              qianchuanReportAweme.ConvertCnt,
			ConvertRate:             fmt.Sprintf("%.2f%%", qianchuanReportAweme.ConvertRate*100),
			AveragePayOrderStatCost: fmt.Sprintf("%.2f¥", qianchuanReportAweme.AveragePayOrderStatCost),
		})
	}

	totalPage := uint64(math.Ceil(float64(qianchuanReportAwemes.Total) / float64(qianchuanReportAwemes.PageSize)))

	return &v1.ListQianchuanReportAwemesReply{
		Code: 200,
		Data: &v1.ListQianchuanReportAwemesReply_Data{
			PageNum:   qianchuanReportAwemes.PageNum,
			PageSize:  qianchuanReportAwemes.PageSize,
			Total:     qianchuanReportAwemes.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ds *DouyinService) StatisticsQianchuanReportAwemes(ctx context.Context, in *v1.StatisticsQianchuanReportAwemesRequest) (*v1.StatisticsQianchuanReportAwemesReply, error) {
	statistics, err := ds.qrauc.StatisticsQianchuanReportAwemes(ctx, in.Day, in.AdvertiserIds)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsQianchuanReportAwemesReply_Statistics, 0)

	for _, statistic := range statistics.Statistics {
		list = append(list, &v1.StatisticsQianchuanReportAwemesReply_Statistics{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsQianchuanReportAwemesReply{
		Code: 200,
		Data: &v1.StatisticsQianchuanReportAwemesReply_Data{
			Statistics: list,
		},
	}, nil
}

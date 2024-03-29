package service

import (
	"context"
	v1 "douyin/api/douyin/v1"
	"fmt"
	"math"
	"time"
)

func (ds *DouyinService) ListQianchuanReportProducts(ctx context.Context, in *v1.ListQianchuanReportProductsRequest) (*v1.ListQianchuanReportProductsReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	qianchuanReportProducts, err := ds.qrpuc.ListQianchuanReportProducts(ctx, in.PageNum, in.PageSize, uint8(in.IsDistinction), in.Day, in.Keyword, in.AdvertiserIds)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListQianchuanReportProductsReply_Products, 0)

	for _, qianchuanReportProduct := range qianchuanReportProducts.List {
		list = append(list, &v1.ListQianchuanReportProductsReply_Products{
			AdvertiserId:            qianchuanReportProduct.AdvertiserId,
			AdvertiserName:          qianchuanReportProduct.AdvertiserName,
			ProductId:               qianchuanReportProduct.ProductId,
			DiscountPrice:           fmt.Sprintf("%.2f¥", qianchuanReportProduct.DiscountPrice),
			ProductName:             qianchuanReportProduct.ProductName,
			ProductImg:              qianchuanReportProduct.ProductImg,
			ProductUrl:              fmt.Sprintf("http://haohuo.jinritemai.com/views/product/item2?id=%d", qianchuanReportProduct.ProductId),
			StatCost:                fmt.Sprintf("%.2f", qianchuanReportProduct.StatCost),
			PayOrderCount:           qianchuanReportProduct.PayOrderCount,
			PayOrderAmount:          fmt.Sprintf("%.2f", qianchuanReportProduct.PayOrderAmount),
			PayOrderAveragePrice:    fmt.Sprintf("%.2f¥", qianchuanReportProduct.PayOrderAveragePrice),
			Roi:                     fmt.Sprintf("%.2f", qianchuanReportProduct.Roi),
			ShowCnt:                 qianchuanReportProduct.ShowCnt,
			ClickCnt:                qianchuanReportProduct.ClickCnt,
			ClickRate:               fmt.Sprintf("%.2f%%", qianchuanReportProduct.ClickRate*100),
			ConvertCnt:              qianchuanReportProduct.ConvertCnt,
			ConvertRate:             fmt.Sprintf("%.2f%%", qianchuanReportProduct.ConvertRate*100),
			AveragePayOrderStatCost: fmt.Sprintf("%.2f¥", qianchuanReportProduct.AveragePayOrderStatCost),
		})
	}

	totalPage := uint64(math.Ceil(float64(qianchuanReportProducts.Total) / float64(qianchuanReportProducts.PageSize)))

	return &v1.ListQianchuanReportProductsReply{
		Code: 200,
		Data: &v1.ListQianchuanReportProductsReply_Data{
			PageNum:   qianchuanReportProducts.PageNum,
			PageSize:  qianchuanReportProducts.PageSize,
			Total:     qianchuanReportProducts.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ds *DouyinService) StatisticsQianchuanReportProducts(ctx context.Context, in *v1.StatisticsQianchuanReportProductsRequest) (*v1.StatisticsQianchuanReportProductsReply, error) {
	statistics, err := ds.qrpuc.StatisticsQianchuanReportProducts(ctx, in.Day, in.AdvertiserIds)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsQianchuanReportProductsReply_Statistics, 0)

	for _, statistic := range statistics.Statistics {
		list = append(list, &v1.StatisticsQianchuanReportProductsReply_Statistics{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsQianchuanReportProductsReply{
		Code: 200,
		Data: &v1.StatisticsQianchuanReportProductsReply_Data{
			Statistics: list,
		},
	}, nil
}

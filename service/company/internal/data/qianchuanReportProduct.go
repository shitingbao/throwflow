package data

import (
	v1 "company/api/service/douyin/v1"
	"company/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type qianchuanReportProductRepo struct {
	data *Data
	log  *log.Helper
}

func NewQianchuanReportProductRepo(data *Data, logger log.Logger) biz.QianchuanReportProductRepo {
	return &qianchuanReportProductRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (qrpr *qianchuanReportProductRepo) List(ctx context.Context, pageNum, pageSize uint64, isDistinction uint32, day, keyword, advertiserIds string) (*v1.ListQianchuanReportProductsReply, error) {
	list, err := qrpr.data.douyinuc.ListQianchuanReportProducts(ctx, &v1.ListQianchuanReportProductsRequest{
		PageNum:       pageNum,
		PageSize:      pageSize,
		Day:           day,
		Keyword:       keyword,
		AdvertiserIds: advertiserIds,
		IsDistinction: isDistinction,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (qrpr *qianchuanReportProductRepo) Statistics(ctx context.Context, day, advertiserIds string) (*v1.StatisticsQianchuanReportProductsReply, error) {
	list, err := qrpr.data.douyinuc.StatisticsQianchuanReportProducts(ctx, &v1.StatisticsQianchuanReportProductsRequest{
		Day:           day,
		AdvertiserIds: advertiserIds,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

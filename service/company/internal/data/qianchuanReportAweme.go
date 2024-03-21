package data

import (
	v1 "company/api/service/douyin/v1"
	"company/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type qianchuanReportAwemeRepo struct {
	data *Data
	log  *log.Helper
}

func NewQianchuanReportAwemeRepo(data *Data, logger log.Logger) biz.QianchuanReportAwemeRepo {
	return &qianchuanReportAwemeRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (qrar *qianchuanReportAwemeRepo) List(ctx context.Context, pageNum, pageSize uint64, isDistinction uint32, day, keyword, advertiserIds string) (*v1.ListQianchuanReportAwemesReply, error) {
	list, err := qrar.data.douyinuc.ListQianchuanReportAwemes(ctx, &v1.ListQianchuanReportAwemesRequest{
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

func (qrar *qianchuanReportAwemeRepo) Statistics(ctx context.Context, day, advertiserIds string) (*v1.StatisticsQianchuanReportAwemesReply, error) {
	list, err := qrar.data.douyinuc.StatisticsQianchuanReportAwemes(ctx, &v1.StatisticsQianchuanReportAwemesRequest{
		Day:           day,
		AdvertiserIds: advertiserIds,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

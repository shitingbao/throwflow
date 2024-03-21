package data

import (
	v1 "company/api/service/douyin/v1"
	"company/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type qianchuanReportAdvertiserRepo struct {
	data *Data
	log  *log.Helper
}

func NewQianchuanReportAdvertiserRepo(data *Data, logger log.Logger) biz.QianchuanReportAdvertiserRepo {
	return &qianchuanReportAdvertiserRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (qrar *qianchuanReportAdvertiserRepo) List(ctx context.Context, advertiserIds, day string) (*v1.ListQianchuanReportAdvertisersReply, error) {
	list, err := qrar.data.douyinuc.ListQianchuanReportAdvertisers(ctx, &v1.ListQianchuanReportAdvertisersRequest{
		Day:           day,
		AdvertiserIds: advertiserIds,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

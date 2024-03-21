package data

import (
	v1 "company/api/service/douyin/v1"
	"company/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type qianchuanAdvertiserHistoryRepo struct {
	data *Data
	log  *log.Helper
}

func NewQianchuanAdvertiserHistoryRepo(data *Data, logger log.Logger) biz.QianchuanAdvertiserHistoryRepo {
	return &qianchuanAdvertiserHistoryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (qahr *qianchuanAdvertiserHistoryRepo) List(ctx context.Context, day, advertiserIds string) (*v1.ListQianchuanAdvertiserHistorysReply, error) {
	list, err := qahr.data.douyinuc.ListQianchuanAdvertiserHistorys(ctx, &v1.ListQianchuanAdvertiserHistorysRequest{
		Day:           day,
		AdvertiserIds: advertiserIds,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

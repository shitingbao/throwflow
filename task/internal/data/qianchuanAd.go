package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "task/api/service/douyin/v1"
	"task/internal/biz"
)

type qianchuanAdRepo struct {
	data *Data
	log  *log.Helper
}

func NewQianchuanAdRepo(data *Data, logger log.Logger) biz.QianchuanAdRepo {
	return &qianchuanAdRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (qar *qianchuanAdRepo) Sync(ctx context.Context) (*v1.SyncQianchuanAdsReply, error) {
	qianchuanAd, err := qar.data.douyinuc.SyncQianchuanAds(ctx, &v1.SyncQianchuanAdsRequest{
		Day: "",
	})

	if err != nil {
		return nil, err
	}

	return qianchuanAd, err
}

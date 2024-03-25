package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/douyin/v1"
	"interface/internal/biz"
)

type qianchuanAdvertiserRepo struct {
	data *Data
	log  *log.Helper
}

func NewQianchuanAdvertiserRepo(data *Data, logger log.Logger) biz.QianchuanAdvertiserRepo {
	return &qianchuanAdvertiserRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (qar *qianchuanAdvertiserRepo) List(ctx context.Context, companyId, pageNum, pageSize uint64, keyword string) (*v1.ListQianchuanAdvertisersReply, error) {
	list, err := qar.data.douyinuc.ListQianchuanAdvertisers(ctx, &v1.ListQianchuanAdvertisersRequest{
		PageNum:   pageNum,
		PageSize:  pageSize,
		CompanyId: companyId,
		Keyword:   keyword,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (qar *qianchuanAdvertiserRepo) Statistics(ctx context.Context, companyId uint64) (*v1.StatisticsQianchuanAdvertisersReply, error) {
	list, err := qar.data.douyinuc.StatisticsQianchuanAdvertisers(ctx, &v1.StatisticsQianchuanAdvertisersRequest{
		CompanyId: companyId,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (qar *qianchuanAdvertiserRepo) Update(ctx context.Context, companyId, advertiserId uint64, status uint32) (*v1.UpdateStatusQianchuanAdvertisersReply, error) {
	qianchuanAdvertiser, err := qar.data.douyinuc.UpdateStatusQianchuanAdvertisers(ctx, &v1.UpdateStatusQianchuanAdvertisersRequest{
		CompanyId:    companyId,
		AdvertiserId: advertiserId,
		Status:       status,
	})

	if err != nil {
		return nil, err
	}

	return qianchuanAdvertiser, err
}

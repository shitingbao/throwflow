package data

import (
	v1 "company/api/service/douyin/v1"
	"company/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
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

func (qar *qianchuanAdvertiserRepo) GetByCompanyIdAndAdvertiserIds(ctx context.Context, companyId, advertiserId uint64) (*v1.GetQianchuanAdvertiserByCompanyIdAndAdvertiserIdsReply, error) {
	qianchuanAdvertiser, err := qar.data.douyinuc.GetQianchuanAdvertiserByCompanyIdAndAdvertiserIds(ctx, &v1.GetQianchuanAdvertiserByCompanyIdAndAdvertiserIdsRequest{
		CompanyId:    companyId,
		AdvertiserId: advertiserId,
	})

	if err != nil {
		return nil, err
	}

	return qianchuanAdvertiser, err
}

func (qar *qianchuanAdvertiserRepo) List(ctx context.Context, companyId, pageNum, pageSize uint64, keyword, advertiserIds, status string) (*v1.ListQianchuanAdvertisersReply, error) {
	list, err := qar.data.douyinuc.ListQianchuanAdvertisers(ctx, &v1.ListQianchuanAdvertisersRequest{
		PageNum:       pageNum,
		PageSize:      pageSize,
		CompanyId:     companyId,
		Keyword:       keyword,
		AdvertiserIds: advertiserIds,
		Status:        status,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (qar *qianchuanAdvertiserRepo) ListByDays(ctx context.Context, companyId uint64, day string) (*v1.ListQianchuanAdvertiserByDaysReply, error) {
	list, err := qar.data.douyinuc.ListQianchuanAdvertiserByDays(ctx, &v1.ListQianchuanAdvertiserByDaysRequest{
		CompanyId: companyId,
		Day:       day,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (qar *qianchuanAdvertiserRepo) ListExternal(ctx context.Context, pageNum, pageSize uint64, startDay, endDay, advertiserIds string) (*v1.ListExternalQianchuanAdvertisersReply, error) {
	list, err := qar.data.douyinuc.ListExternalQianchuanAdvertisers(ctx, &v1.ListExternalQianchuanAdvertisersRequest{
		PageNum:       pageNum,
		PageSize:      pageSize,
		StartDay:      startDay,
		EndDay:        endDay,
		AdvertiserIds: advertiserIds,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (qar *qianchuanAdvertiserRepo) StatisticsDashboard(ctx context.Context, companyId uint64, day, advertiserIds string) (*v1.StatisticsDashboardQianchuanAdvertisersReply, error) {
	list, err := qar.data.douyinuc.StatisticsDashboardQianchuanAdvertisers(ctx, &v1.StatisticsDashboardQianchuanAdvertisersRequest{
		CompanyId:     companyId,
		Day:           day,
		AdvertiserIds: advertiserIds,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (qar *qianchuanAdvertiserRepo) StatisticsExternal(ctx context.Context, startDay, endDay, advertiserIds string) (*v1.StatisticsExternalQianchuanAdvertisersReply, error) {
	list, err := qar.data.douyinuc.StatisticsExternalQianchuanAdvertisers(ctx, &v1.StatisticsExternalQianchuanAdvertisersRequest{
		StartDay:      startDay,
		EndDay:        endDay,
		AdvertiserIds: advertiserIds,
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

func (qar *qianchuanAdvertiserRepo) UpdateStatusByCompanyId(ctx context.Context, companyId uint64, status uint32) (*v1.UpdateStatusQianchuanAdvertisersByCompanyIdReply, error) {
	qianchuanAdvertiser, err := qar.data.douyinuc.UpdateStatusQianchuanAdvertisersByCompanyId(ctx, &v1.UpdateStatusQianchuanAdvertisersByCompanyIdRequest{
		CompanyId: companyId,
		Status:    status,
	})

	if err != nil {
		return nil, err
	}

	return qianchuanAdvertiser, err
}

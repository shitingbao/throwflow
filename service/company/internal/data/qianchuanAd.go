package data

import (
	v1 "company/api/service/douyin/v1"
	"company/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (qar *qianchuanAdRepo) GetExternal(ctx context.Context, adId uint64, day string) (*v1.GetExternalQianchuanAdsReply, error) {
	qianchuanAd, err := qar.data.douyinuc.GetExternalQianchuanAds(ctx, &v1.GetExternalQianchuanAdsRequest{
		AdId: adId,
		Day:  day,
	})

	if err != nil {
		return nil, err
	}

	return qianchuanAd, err
}

func (qar *qianchuanAdRepo) GetExternalHistory(ctx context.Context, adId uint64, startDay, endDay string) (*v1.GetExternalHistoryQianchuanAdsReply, error) {
	qianchuanAd, err := qar.data.douyinuc.GetExternalHistoryQianchuanAds(ctx, &v1.GetExternalHistoryQianchuanAdsRequest{
		AdId:     adId,
		StartDay: startDay,
		EndDay:   endDay,
	})

	if err != nil {
		return nil, err
	}

	return qianchuanAd, err
}

func (qar *qianchuanAdRepo) List(ctx context.Context, pageNum, pageSize uint64, day, keyword, advertiserIds string) (*v1.ListQianchuanAdsReply, error) {
	list, err := qar.data.douyinuc.ListQianchuanAds(ctx, &v1.ListQianchuanAdsRequest{
		PageNum:       pageNum,
		PageSize:      pageSize,
		Day:           day,
		Keyword:       keyword,
		AdvertiserIds: advertiserIds,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (qar *qianchuanAdRepo) ListExternal(ctx context.Context, pageNum, pageSize uint64, startDay, endDay, keyword, advertiserIds, filter, orderName, orderType string) (*v1.ListExternalQianchuanAdsReply, error) {
	list, err := qar.data.douyinuc.ListExternalQianchuanAds(ctx, &v1.ListExternalQianchuanAdsRequest{
		PageNum:       pageNum,
		PageSize:      pageSize,
		StartDay:      startDay,
		EndDay:        endDay,
		Keyword:       keyword,
		AdvertiserIds: advertiserIds,
		Filter:        filter,
		OrderName:     orderName,
		OrderType:     orderType,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (qar *qianchuanAdRepo) ListSelectExternal(ctx context.Context) (*v1.ListSelectExternalQianchuanAdsReply, error) {
	list, err := qar.data.douyinuc.ListSelectExternalQianchuanAds(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (qar *qianchuanAdRepo) StatisticsExternal(ctx context.Context, startDay, endDay, keyword, advertiserIds, filter string) (*v1.StatisticsExternalQianchuanAdsReply, error) {
	list, err := qar.data.douyinuc.StatisticsExternalQianchuanAds(ctx, &v1.StatisticsExternalQianchuanAdsRequest{
		StartDay:      startDay,
		EndDay:        endDay,
		Keyword:       keyword,
		AdvertiserIds: advertiserIds,
		Filter:        filter,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

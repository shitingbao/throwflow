package service

import (
	"context"
	v1 "douyin/api/douyin/v1"
)

func (ds *DouyinService) ListQianchuanAdvertiserHistorys(ctx context.Context, in *v1.ListQianchuanAdvertiserHistorysRequest) (*v1.ListQianchuanAdvertiserHistorysReply, error) {
	qianchuanAdvertiserHistorys, err := ds.qahuc.ListQianchuanAdvertiserHistorys(ctx, in.Day, in.AdvertiserIds)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListQianchuanAdvertiserHistorysReply_QianchuanAdvertisers, 0)

	for _, qianchuanAdvertiserHistory := range qianchuanAdvertiserHistorys {
		list = append(list, &v1.ListQianchuanAdvertiserHistorysReply_QianchuanAdvertisers{
			AdvertiserId: qianchuanAdvertiserHistory.AdvertiserId,
		})
	}

	return &v1.ListQianchuanAdvertiserHistorysReply{
		Code: 200,
		Data: &v1.ListQianchuanAdvertiserHistorysReply_Data{
			List: list,
		},
	}, nil
}

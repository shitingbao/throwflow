package service

import (
	"context"
	v1 "douyin/api/douyin/v1"
	"math"
)

func (ds *DouyinService) ListOpenDouyinUserInfos(ctx context.Context, in *v1.ListOpenDouyinUserInfosRequest) (*v1.ListOpenDouyinUserInfosReply, error) {
	openDouyinUserInfos, err := ds.oduiuc.ListOpenDouyinUserInfos(ctx, in.PageNum, in.PageSize)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListOpenDouyinUserInfosReply_OpenDouyinUserInfos, 0)

	for _, openDouyinUserInfo := range openDouyinUserInfos.List {
		list = append(list, &v1.ListOpenDouyinUserInfosReply_OpenDouyinUserInfos{
			ClientKey: openDouyinUserInfo.ClientKey,
			OpenId:    openDouyinUserInfo.OpenId,
			VideoId:   openDouyinUserInfo.VideoId,
		})
	}

	totalPage := uint64(math.Ceil(float64(openDouyinUserInfos.Total) / float64(openDouyinUserInfos.PageSize)))

	return &v1.ListOpenDouyinUserInfosReply{
		Code: 200,
		Data: &v1.ListOpenDouyinUserInfosReply_Data{
			PageNum:   openDouyinUserInfos.PageNum,
			PageSize:  openDouyinUserInfos.PageSize,
			Total:     openDouyinUserInfos.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ds *DouyinService) ListOpenDouyinUserInfosByProductId(ctx context.Context, in *v1.ListOpenDouyinUserInfosByProductIdRequest) (*v1.ListOpenDouyinUserInfosByProductIdReply, error) {
	openDouyinUserInfos, err := ds.oduiuc.ListOpenDouyinUserInfosByProductId(ctx, in.ProductId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListOpenDouyinUserInfosByProductIdReply_OpenDouyinUserInfos, 0)

	for _, openDouyinUserInfo := range openDouyinUserInfos {
		list = append(list, &v1.ListOpenDouyinUserInfosByProductIdReply_OpenDouyinUserInfos{
			AccountId:    openDouyinUserInfo.AccountId,
			Nickname:     openDouyinUserInfo.Nickname,
			Avatar:       openDouyinUserInfo.Avatar,
			AvatarLarger: openDouyinUserInfo.AvatarLarger,
		})
	}

	return &v1.ListOpenDouyinUserInfosByProductIdReply{
		Code: 200,
		Data: &v1.ListOpenDouyinUserInfosByProductIdReply_Data{
			List: list,
		},
	}, nil
}

func (ds *DouyinService) UpdateOpenDouyinUserInfos(ctx context.Context, in *v1.UpdateOpenDouyinUserInfosRequest) (*v1.UpdateOpenDouyinUserInfosReply, error) {
	if err := ds.oduiuc.UpdateOpenDouyinUserInfos(ctx, in.AwemeId, in.ClientKey, in.OpenId, in.AccountId, in.Nickname, in.Avatar, in.AvatarLarger, in.Area); err != nil {
		return nil, err
	}

	return &v1.UpdateOpenDouyinUserInfosReply{
		Code: 200,
		Data: &v1.UpdateOpenDouyinUserInfosReply_Data{},
	}, nil
}

package service

import (
	"context"
	"math"
	v1 "weixin/api/weixin/v1"
	"weixin/internal/pkg/tool"
)

func (ws *WeixinService) GetOpenDouyinUsers(ctx context.Context, in *v1.GetOpenDouyinUsersRequest) (*v1.GetOpenDouyinUsersReply, error) {
	openDouyinUser, err := ws.uoduc.GetOpenDouyinUsers(ctx, in.ClientKey, in.OpenId)

	if err != nil {
		return nil, err
	}

	return &v1.GetOpenDouyinUsersReply{
		Code: 200,
		Data: &v1.GetOpenDouyinUsersReply_Data{
			OpenDouyinUserId: openDouyinUser.Id,
			ClientKey:        openDouyinUser.ClientKey,
			OpenId:           openDouyinUser.OpenId,
			AwemeId:          openDouyinUser.AwemeId,
			AccountId:        openDouyinUser.AccountId,
			Nickname:         openDouyinUser.Nickname,
			Avatar:           openDouyinUser.Avatar,
			AvatarLarger:     openDouyinUser.AvatarLarger,
			CooperativeCode:  openDouyinUser.CooperativeCode,
			Fans:             openDouyinUser.Fans,
			FansShow:         openDouyinUser.FansShow,
			Area:             openDouyinUser.Area,
			CreateDate:       tool.TimeToString("2006-01-02", openDouyinUser.CreateTime),
		},
	}, nil
}

func (ws *WeixinService) ListOpenDouyinUsers(ctx context.Context, in *v1.ListOpenDouyinUsersRequest) (*v1.ListOpenDouyinUsersReply, error) {
	openDouyinUsers, err := ws.uoduc.ListOpenDouyinUsers(ctx, in.PageNum, in.PageSize, in.UserId, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListOpenDouyinUsersReply_OpenDouyinUser, 0)

	for _, openDouyinUser := range openDouyinUsers.List {
		list = append(list, &v1.ListOpenDouyinUsersReply_OpenDouyinUser{
			OpenDouyinUserId: openDouyinUser.Id,
			ClientKey:        openDouyinUser.ClientKey,
			OpenId:           openDouyinUser.OpenId,
			AwemeId:          openDouyinUser.AwemeId,
			AccountId:        openDouyinUser.AccountId,
			Nickname:         openDouyinUser.Nickname,
			Avatar:           openDouyinUser.Avatar,
			AvatarLarger:     openDouyinUser.AvatarLarger,
			CooperativeCode:  openDouyinUser.CooperativeCode,
			Fans:             openDouyinUser.Fans,
			FansShow:         openDouyinUser.FansShow,
			Area:             openDouyinUser.Area,
			CreateDate:       tool.TimeToString("2006-01-02", openDouyinUser.CreateTime),
			Level:            openDouyinUser.Level,
		})
	}

	totalPage := uint64(math.Ceil(float64(openDouyinUsers.Total) / float64(openDouyinUsers.PageSize)))

	return &v1.ListOpenDouyinUsersReply{
		Code: 200,
		Data: &v1.ListOpenDouyinUsersReply_Data{
			PageNum:   openDouyinUsers.PageNum,
			PageSize:  openDouyinUsers.PageSize,
			Total:     openDouyinUsers.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ws *WeixinService) ListByClientKeyAndOpenIds(ctx context.Context, in *v1.ListByClientKeyAndOpenIdsRequest) (*v1.ListByClientKeyAndOpenIdsReply, error) {
	openDouyinUsers, err := ws.uoduc.ListByClientKeyAndOpenIds(ctx, in.PageNum, in.PageSize, in.ClientKeyAndOpenIds, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListByClientKeyAndOpenIdsReply_OpenDouyinUser, 0)

	for _, openDouyinUser := range openDouyinUsers.List {
		list = append(list, &v1.ListByClientKeyAndOpenIdsReply_OpenDouyinUser{
			OpenDouyinUserId: openDouyinUser.Id,
			ClientKey:        openDouyinUser.ClientKey,
			OpenId:           openDouyinUser.OpenId,
			AwemeId:          openDouyinUser.AwemeId,
			AccountId:        openDouyinUser.AccountId,
			Nickname:         openDouyinUser.Nickname,
			Avatar:           openDouyinUser.Avatar,
			AvatarLarger:     openDouyinUser.AvatarLarger,
			CooperativeCode:  openDouyinUser.CooperativeCode,
			Fans:             openDouyinUser.Fans,
			FansShow:         openDouyinUser.FansShow,
			Area:             openDouyinUser.Area,
			CreateDate:       tool.TimeToString("2006-01-02", openDouyinUser.CreateTime),
		})
	}

	totalPage := uint64(math.Ceil(float64(openDouyinUsers.Total) / float64(openDouyinUsers.PageSize)))

	return &v1.ListByClientKeyAndOpenIdsReply{
		Code: 200,
		Data: &v1.ListByClientKeyAndOpenIdsReply_Data{
			PageNum:   openDouyinUsers.PageNum,
			PageSize:  openDouyinUsers.PageSize,
			Total:     openDouyinUsers.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ws *WeixinService) UpdateOpenDouyinUsers(ctx context.Context, in *v1.UpdateOpenDouyinUsersRequest) (*v1.UpdateOpenDouyinUsersReply, error) {
	err := ws.uoduc.UpdateOpenDouyinUsers(ctx, in.UserId, in.AwemeId, in.ClientKey, in.OpenId, in.AccountId, in.Nickname, in.Avatar, in.AvatarLarger, in.Area)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateOpenDouyinUsersReply{
		Code: 200,
		Data: &v1.UpdateOpenDouyinUsersReply_Data{},
	}, nil
}

func (ws *WeixinService) UpdateUserInfoOpenDouyinUsers(ctx context.Context, in *v1.UpdateUserInfoOpenDouyinUsersRequest) (*v1.UpdateUserInfoOpenDouyinUsersReply, error) {
	err := ws.uoduc.UpdateUserInfoOpenDouyinUsers(ctx, in.AwemeId, in.ClientKey, in.OpenId, in.AccountId, in.Nickname, in.Avatar, in.AvatarLarger, in.Area)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateUserInfoOpenDouyinUsersReply{
		Code: 200,
		Data: &v1.UpdateUserInfoOpenDouyinUsersReply_Data{},
	}, nil
}

func (ws *WeixinService) UpdateCooperativeCodeOpenDouyinUsers(ctx context.Context, in *v1.UpdateCooperativeCodeOpenDouyinUsersRequest) (*v1.UpdateCooperativeCodeOpenDouyinUsersReply, error) {
	openDouyinUsers, err := ws.uoduc.UpdateCooperativeCodeOpenDouyinUsers(ctx, in.UserId, in.OpenDouyinUserId, in.CooperativeCode)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.UpdateCooperativeCodeOpenDouyinUsersReply_OpenDouyinUser, 0)

	for _, openDouyinUser := range openDouyinUsers {
		list = append(list, &v1.UpdateCooperativeCodeOpenDouyinUsersReply_OpenDouyinUser{
			OpenDouyinUserId: openDouyinUser.Id,
			Nickname:         openDouyinUser.Nickname,
			AccountId:        openDouyinUser.AccountId,
			Avatar:           openDouyinUser.Avatar,
			AvatarLarger:     openDouyinUser.AvatarLarger,
			CooperativeCode:  openDouyinUser.CooperativeCode,
		})
	}

	return &v1.UpdateCooperativeCodeOpenDouyinUsersReply{
		Code: 200,
		Data: &v1.UpdateCooperativeCodeOpenDouyinUsersReply_Data{
			List: list,
		},
	}, nil
}

func (ws *WeixinService) UpdateUserFansOpenDouyinUsers(ctx context.Context, in *v1.UpdateUserFansOpenDouyinUsersRequest) (*v1.UpdateUserFansOpenDouyinUsersReply, error) {
	err := ws.uoduc.UpdateUserFansOpenDouyinUsers(ctx, in.ClientKey, in.OpenId, in.Fans)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateUserFansOpenDouyinUsersReply{
		Code: 200,
		Data: &v1.UpdateUserFansOpenDouyinUsersReply_Data{},
	}, nil
}

func (ws *WeixinService) DeleteOpenDouyinUsers(ctx context.Context, in *v1.DeleteOpenDouyinUsersRequest) (*v1.DeleteOpenDouyinUsersReply, error) {
	err := ws.uoduc.DeleteOpenDouyinUsers(ctx, in.UserId, in.OpenDouyinUserId)

	if err != nil {
		return nil, err
	}

	return &v1.DeleteOpenDouyinUsersReply{
		Code: 200,
		Data: &v1.DeleteOpenDouyinUsersReply_Data{},
	}, nil
}

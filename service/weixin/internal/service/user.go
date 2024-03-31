package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
	v1 "weixin/api/weixin/v1"
)

func (ws *WeixinService) GetUsers(ctx context.Context, in *v1.GetUsersRequest) (*v1.GetUsersReply, error) {
	user, err := ws.uuc.GetUsers(ctx, in.Token)

	if err != nil {
		return nil, err
	}

	openIds := make([]*v1.GetUsersReply_OpenId, 0)

	for _, openId := range user.OpenIds {
		openIds = append(openIds, &v1.GetUsersReply_OpenId{
			Appid:  openId.Appid,
			OpenId: openId.OpenId,
		})
	}

	return &v1.GetUsersReply{
		Code: 200,
		Data: &v1.GetUsersReply_Data{
			UserId:            user.Id,
			Phone:             user.Phone,
			CountryCode:       user.CountryCode,
			NickName:          user.NickName,
			AvatarUrl:         user.AvatarUrl,
			Balance:           user.Balance,
			Integral:          user.Integral,
			IntegralLevelName: user.IntegralLevelName,
			Ranking:           user.Ranking,
			Total:             user.Total,
			TotalRanking:      user.TotalRanking,
			OpenIds:           openIds,
		},
	}, nil
}

func (ws *WeixinService) GetByIdUsers(ctx context.Context, in *v1.GetByIdUsersRequest) (*v1.GetByIdUsersReply, error) {
	user, err := ws.uuc.GetByIdUsers(ctx, in.UserId)

	if err != nil {
		return nil, err
	}

	openIds := make([]*v1.GetByIdUsersReply_OpenId, 0)

	for _, openId := range user.OpenIds {
		openIds = append(openIds, &v1.GetByIdUsersReply_OpenId{
			Appid:  openId.Appid,
			OpenId: openId.OpenId,
		})
	}

	return &v1.GetByIdUsersReply{
		Code: 200,
		Data: &v1.GetByIdUsersReply_Data{
			UserId:            user.Id,
			Phone:             user.Phone,
			CountryCode:       user.CountryCode,
			NickName:          user.NickName,
			AvatarUrl:         user.AvatarUrl,
			Balance:           user.Balance,
			Integral:          user.Integral,
			IntegralLevelName: user.IntegralLevelName,
			Ranking:           user.Ranking,
			Total:             user.Total,
			TotalRanking:      user.TotalRanking,
			OpenIds:           openIds,
		},
	}, nil
}

func (ws *WeixinService) GetFollowUsers(ctx context.Context, in *v1.GetFollowUsersRequest) (*v1.GetFollowUsersReply, error) {
	followData, err := ws.uuc.GetFollowUsers(ctx, in.OrganizationId, in.ParentUserId)

	if err != nil {
		return nil, err
	}

	return &v1.GetFollowUsersReply{
		Code: 200,
		Data: &v1.GetFollowUsersReply_Data{
			FollowType:    followData.FollowType,
			FollowName:    followData.FollowName,
			FollowLogoUrl: followData.FollowLogoUrl,
			QrCodeUrl:     followData.QrCodeUrl,
			TotalNum:      followData.TotalNum,
		},
	}, nil
}

func (ws *WeixinService) CreateUsers(ctx context.Context, in *v1.CreateUsersRequest) (*v1.CreateUsersReply, error) {
	user, err := ws.uuc.CreateUsers(ctx, in.OrganizationId, in.LoginCode, in.PhoneCode)

	if err != nil {
		return nil, err
	}

	openIds := make([]*v1.CreateUsersReply_OpenId, 0)

	for _, openId := range user.OpenIds {
		openIds = append(openIds, &v1.CreateUsersReply_OpenId{
			Appid:  openId.Appid,
			OpenId: openId.OpenId,
		})
	}

	return &v1.CreateUsersReply{
		Code: 200,
		Data: &v1.CreateUsersReply_Data{
			UserId:            user.Id,
			Phone:             user.Phone,
			CountryCode:       user.CountryCode,
			NickName:          user.NickName,
			AvatarUrl:         user.AvatarUrl,
			Balance:           user.Balance,
			Integral:          user.Integral,
			IntegralLevelName: user.IntegralLevelName,
			Ranking:           user.Ranking,
			Total:             user.Total,
			TotalRanking:      user.TotalRanking,
			Token:             user.Token,
			OpenIds:           openIds,
		},
	}, nil
}

func (ws *WeixinService) UpdateUsers(ctx context.Context, in *v1.UpdateUsersRequest) (*v1.UpdateUsersReply, error) {
	user, err := ws.uuc.UpdateUsers(ctx, in.UserId, in.NickName, in.Avatar)

	if err != nil {
		return nil, err
	}

	openIds := make([]*v1.UpdateUsersReply_OpenId, 0)

	for _, openId := range user.OpenIds {
		openIds = append(openIds, &v1.UpdateUsersReply_OpenId{
			Appid:  openId.Appid,
			OpenId: openId.OpenId,
		})
	}

	return &v1.UpdateUsersReply{
		Code: 200,
		Data: &v1.UpdateUsersReply_Data{
			UserId:            user.Id,
			Phone:             user.Phone,
			CountryCode:       user.CountryCode,
			NickName:          user.NickName,
			AvatarUrl:         user.AvatarUrl,
			Balance:           user.Balance,
			Integral:          user.Integral,
			IntegralLevelName: user.IntegralLevelName,
			Ranking:           user.Ranking,
			Total:             user.Total,
			TotalRanking:      user.TotalRanking,
			OpenIds:           openIds,
		},
	}, nil
}

func (ws *WeixinService) SyncIntegralUsers(ctx context.Context, in *emptypb.Empty) (*v1.SyncIntegralUsersReply, error) {
	ws.log.Infof("同步微信用户积分数据, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ws.uuc.SyncIntegralUsers(ctx); err != nil {
		return nil, err
	}

	ws.log.Infof("同步微信用户积分数据, 结束时间 %s \n", time.Now())

	return &v1.SyncIntegralUsersReply{
		Code: 200,
		Data: &v1.SyncIntegralUsersReply_Data{},
	}, nil
}

func (ws *WeixinService) ImportDatas(ctx context.Context, in *emptypb.Empty) (*v1.ImportDatasReply, error) {
	ws.log.Infof("导入用户数据, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ws.uuc.ImportDatas(ctx); err != nil {
		return nil, err
	}

	ws.log.Infof("导入用户数据, 结束时间 %s \n", time.Now())

	return &v1.ImportDatasReply{
		Code: 200,
		Data: &v1.ImportDatasReply_Data{},
	}, nil
}

func (ws *WeixinService) ParentUserDatas(ctx context.Context, in *emptypb.Empty) (*v1.ParentUserDatasReply, error) {
	ws.log.Infof("导入用户数据, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ws.uuc.ParentUserDatas(ctx); err != nil {
		return nil, err
	}

	ws.log.Infof("导入用户数据, 结束时间 %s \n", time.Now())

	return &v1.ParentUserDatasReply{
		Code: 200,
		Data: &v1.ParentUserDatasReply_Data{},
	}, nil
}

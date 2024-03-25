package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "interface/api/interface/v1"
	"interface/internal/biz"
	"interface/internal/pkg/tool"
)

func (is *InterfaceService) GetUsers(ctx context.Context, in *emptypb.Empty) (*v1.GetUsersReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	user, err := is.uuc.GetUsers(ctx, userInfo.Data.UserId)

	if err != nil {
		return nil, err
	}

	return &v1.GetUsersReply{
		Code: 200,
		Data: &v1.GetUsersReply_Data{
			UserId:            user.Data.UserId,
			Phone:             tool.FormatPhone(user.Data.Phone),
			NickName:          user.Data.NickName,
			AvatarUrl:         user.Data.AvatarUrl,
			Balance:           user.Data.Balance,
			Integral:          user.Data.Integral,
			IntegralLevelName: user.Data.IntegralLevelName,
			Ranking:           user.Data.Ranking,
			Total:             10000 + user.Data.Total,
			TotalRanking:      10000 + user.Data.TotalRanking,
		},
	}, nil
}

func (is *InterfaceService) GetStorePreferenceUserOpenDouyins(ctx context.Context, in *emptypb.Empty) (*v1.GetStorePreferenceUserOpenDouyinsReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	storePreferences, err := is.jouc.GetStorePreferenceUserOpenDouyins(ctx, userInfo.Data.UserId)

	if err != nil {
		return nil, err
	}

	storePreference := make([]*v1.GetStorePreferenceUserOpenDouyinsReply_StorePreference, 0)

	for _, lstorePreference := range storePreferences.Data.StorePreference {
		storePreference = append(storePreference, &v1.GetStorePreferenceUserOpenDouyinsReply_StorePreference{
			IndustryId:    lstorePreference.IndustryId,
			IndustryName:  lstorePreference.IndustryName,
			IndustryRatio: lstorePreference.IndustryRatio,
		})
	}

	return &v1.GetStorePreferenceUserOpenDouyinsReply{
		Code: 200,
		Data: &v1.GetStorePreferenceUserOpenDouyinsReply_Data{
			StorePreference: storePreference,
		},
	}, nil
}

func (is *InterfaceService) GetFollows(ctx context.Context, in *v1.GetFollowsRequest) (*v1.GetFollowsReply, error) {
	follow, err := is.uuc.GetFollows(ctx, in.OrganizationId, in.ParentUserId)

	if err != nil {
		return nil, err
	}

	return &v1.GetFollowsReply{
		Code: 200,
		Data: &v1.GetFollowsReply_Data{
			FollowType:    follow.Data.FollowType,
			FollowName:    follow.Data.FollowName,
			FollowLogoUrl: follow.Data.FollowLogoUrl,
			QrCodeUrl:     follow.Data.QrCodeUrl,
			TotalNum:      follow.Data.TotalNum,
		},
	}, nil
}

func (is *InterfaceService) CreateUsers(ctx context.Context, in *v1.CreateUsersRequest) (*v1.CreateUsersReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	user, err := is.uuc.CreateUsers(ctx, in.OrganizationId, in.LoginCode, in.PhoneCode)

	if err != nil {
		return nil, err
	}

	return &v1.CreateUsersReply{
		Code: 200,
		Data: &v1.CreateUsersReply_Data{
			UserId:            user.Data.UserId,
			Phone:             tool.FormatPhone(user.Data.Phone),
			NickName:          user.Data.NickName,
			AvatarUrl:         user.Data.AvatarUrl,
			Balance:           user.Data.Balance,
			Integral:          user.Data.Integral,
			IntegralLevelName: user.Data.IntegralLevelName,
			Ranking:           user.Data.Ranking,
			Total:             user.Data.Total,
			TotalRanking:      user.Data.TotalRanking,
			Token:             user.Data.Token,
		},
	}, nil
}

func (is *InterfaceService) UpdateUsers(ctx context.Context, in *v1.UpdateUsersRequest) (*v1.UpdateUsersReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	user, err := is.uuc.UpdateUsers(ctx, userInfo.Data.UserId, in.NickName, in.Avatar)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateUsersReply{
		Code: 200,
		Data: &v1.UpdateUsersReply_Data{
			UserId:            user.Data.UserId,
			Phone:             tool.FormatPhone(user.Data.Phone),
			NickName:          user.Data.NickName,
			AvatarUrl:         user.Data.AvatarUrl,
			Balance:           user.Data.Balance,
			Integral:          user.Data.Integral,
			IntegralLevelName: user.Data.IntegralLevelName,
			Ranking:           user.Data.Ranking,
			Total:             user.Data.Total,
			TotalRanking:      user.Data.TotalRanking,
		},
	}, nil
}

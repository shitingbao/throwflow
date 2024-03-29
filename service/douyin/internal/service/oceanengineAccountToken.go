package service

import (
	"context"
	v1 "douyin/api/douyin/v1"
	"douyin/internal/pkg/tool"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

func (ds *DouyinService) GetOceanengineAccountTokens(ctx context.Context, in *v1.GetOceanengineAccountTokensRequest) (*v1.GetOceanengineAccountTokensReply, error) {
	oceanengineAccountToken, err := ds.oatuc.GetOceanengineAccountTokens(ctx, in.CompanyId, in.AdvertiserId)

	if err != nil {
		return nil, err
	}

	return &v1.GetOceanengineAccountTokensReply{
		Code: 200,
		Data: &v1.GetOceanengineAccountTokensReply_Data{
			AccountId:             oceanengineAccountToken.AccountId,
			CompanyId:             oceanengineAccountToken.CompanyId,
			AppId:                 oceanengineAccountToken.AppId,
			AccessToken:           oceanengineAccountToken.AccessToken,
			ExpiresIn:             oceanengineAccountToken.ExpiresIn,
			RefreshToken:          oceanengineAccountToken.RefreshToken,
			RefreshTokenExpiresIn: oceanengineAccountToken.RefreshTokenExpiresIn,
			AccountRole:           oceanengineAccountToken.AccountRole,
			UpdateTime:            tool.TimeToString("2006-01-02 15:04:05", oceanengineAccountToken.UpdateTime),
		},
	}, nil
}

func (ds *DouyinService) RefreshOceanengineAccountTokens(ctx context.Context, in *emptypb.Empty) (*v1.RefreshOceanengineAccountTokensReply, error) {
	ds.log.Infof("刷新千川授权账户Token, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ds.oatuc.RefreshOceanengineAccountTokens(ctx); err != nil {
		return nil, err
	}

	ds.log.Infof("刷新千川授权账户Token, 结束时间 %s \n", time.Now())

	return &v1.RefreshOceanengineAccountTokensReply{
		Code: 200,
		Data: &v1.RefreshOceanengineAccountTokensReply_Data{},
	}, nil
}

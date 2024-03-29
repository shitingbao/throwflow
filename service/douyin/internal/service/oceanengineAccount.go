package service

import (
	"context"
	v1 "douyin/api/douyin/v1"
)

func (ds *DouyinService) CreateOceanengineAccounts(ctx context.Context, in *v1.CreateOceanengineAccountsRequest) (*v1.CreateOceanengineAccountsReply, error) {
	ctx = context.Background()

	if err := ds.oauc.CreateOceanengineAccounts(ctx, in.CompanyId, in.AppId, in.AuthCode); err != nil {
		return nil, err
	}

	return &v1.CreateOceanengineAccountsReply{
		Code: 200,
		Data: &v1.CreateOceanengineAccountsReply_Data{},
	}, nil
}

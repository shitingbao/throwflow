package service

import (
	v1 "admin/api/admin/v1"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (as *AdminService) GetToken(ctx context.Context, in *emptypb.Empty) (*v1.GetTokenReply, error) {
	if _, err := as.verifyPermission(ctx, "all"); err != nil {
		return nil, err
	}

	token, err := as.tuc.CreateToken(ctx)

	if err != nil {
		return nil, err
	}

	return &v1.GetTokenReply{
		Code: 200,
		Data: &v1.GetTokenReply_Data{
			Token: token,
		},
	}, err
}

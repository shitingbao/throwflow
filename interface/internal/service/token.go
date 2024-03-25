package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "interface/api/interface/v1"
)

func (is *InterfaceService) GetToken(ctx context.Context, in *emptypb.Empty) (*v1.GetTokenReply, error) {
	token, err := is.tuc.CreateToken(ctx)

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

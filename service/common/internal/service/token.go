package service

import (
	v1 "common/api/common/v1"
	"context"
)

func (cs *CommonService) GetToken(ctx context.Context, in *v1.GetTokenRequest) (*v1.GetTokenReply, error) {
	token, err := cs.tuc.CreateToken(ctx, in.Key)

	if err != nil {
		return nil, err
	}

	return &v1.GetTokenReply{
		Code: 200,
		Data: &v1.GetTokenReply_Data{
			Token: token,
		},
	}, nil
}

func (cs *CommonService) VerifyToken(ctx context.Context, in *v1.VerifyTokenRequest) (*v1.VerifyTokenReply, error) {
	if err := cs.tuc.VerifyToken(ctx, in.Key); err != nil {
		return nil, err
	}

	return &v1.VerifyTokenReply{
		Code: 200,
		Data: &v1.VerifyTokenReply_Data{},
	}, nil
}

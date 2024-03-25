package service

import (
	"context"
	v1 "interface/api/interface/v1"
	"interface/internal/biz"
)

func (is *InterfaceService) CreateUserScanRecords(ctx context.Context, in *v1.CreateUserScanRecordsRequest) (*v1.CreateUserScanRecordsReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	if _, err := is.usruc.CreateUserScanRecords(ctx, userInfo.Data.UserId, in.OrganizationId, in.ParentUserId); err != nil {
		return nil, err
	}

	return &v1.CreateUserScanRecordsReply{
		Code: 200,
		Data: &v1.CreateUserScanRecordsReply_Data{},
	}, nil
}

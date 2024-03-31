package service

import (
	"context"
	v1 "weixin/api/weixin/v1"
)

func (ws *WeixinService) CreateUserScanRecords(ctx context.Context, in *v1.CreateUserScanRecordsRequest) (*v1.CreateUserScanRecordsReply, error) {
	if err := ws.usruc.CreateUserScanRecords(ctx, in.UserId, in.OrganizationId, in.ParentUserId); err != nil {
		return nil, err
	}

	return &v1.CreateUserScanRecordsReply{
		Code: 200,
		Data: &v1.CreateUserScanRecordsReply_Data{},
	}, nil
}

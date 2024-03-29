package service

import (
	"context"
	v1 "douyin/api/douyin/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

func (ds *DouyinService) CreateDoukeTokens(ctx context.Context, in *v1.CreateDoukeTokensRequest) (*v1.CreateDoukeTokensReply, error) {
	ctx = context.Background()

	if err := ds.dtuc.CreateDoukeTokens(ctx, in.Code); err != nil {
		return nil, err
	}

	return &v1.CreateDoukeTokensReply{
		Code: 200,
		Data: &v1.CreateDoukeTokensReply_Data{},
	}, nil
}

func (ds *DouyinService) RefreshDoukeTokens(ctx context.Context, in *emptypb.Empty) (*v1.RefreshDoukeTokensReply, error) {
	ds.log.Infof("刷新抖客平台授权账户Token, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ds.dtuc.RefreshDoukeTokens(ctx); err != nil {
		return nil, err
	}

	ds.log.Infof("刷新抖客平台授权账户Token, 结束时间 %s \n", time.Now())

	return &v1.RefreshDoukeTokensReply{
		Code: 200,
		Data: &v1.RefreshDoukeTokensReply_Data{},
	}, nil
}

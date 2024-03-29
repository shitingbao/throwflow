package service

import (
	"context"
	v1 "douyin/api/douyin/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

func (ds *DouyinService) GetUrlOpenDouyinTokens(ctx context.Context, in *v1.GetUrlOpenDouyinTokensRequest) (*v1.GetUrlOpenDouyinTokensReply, error) {
	redirectUrl, err := ds.odtuc.GetUrlOpenDouyinTokens(ctx, in.UserId)

	if err != nil {
		return nil, err
	}

	return &v1.GetUrlOpenDouyinTokensReply{
		Code: 200,
		Data: &v1.GetUrlOpenDouyinTokensReply_Data{
			RedirectUrl: redirectUrl,
		},
	}, nil
}

func (ds *DouyinService) GetTicketOpenDouyinTokens(ctx context.Context, in *emptypb.Empty) (*v1.GetTicketOpenDouyinTokensReply, error) {
	ticket, err := ds.odtuc.GetTicketOpenDouyinTokens(ctx)

	if err != nil {
		return nil, err
	}

	return &v1.GetTicketOpenDouyinTokensReply{
		Code: 200,
		Data: &v1.GetTicketOpenDouyinTokensReply_Data{
			Ticket: ticket.Data.Ticket,
		},
	}, nil
}

func (ds *DouyinService) GetQrCodeOpenDouyinTokens(ctx context.Context, in *v1.GetQrCodeOpenDouyinTokensRequest) (*v1.GetQrCodeOpenDouyinTokensReply, error) {
	qrCode, err := ds.odtuc.GetQrCodeOpenDouyinTokens(ctx, in.UserId)

	if err != nil {
		return nil, err
	}

	return &v1.GetQrCodeOpenDouyinTokensReply{
		Code: 200,
		Data: &v1.GetQrCodeOpenDouyinTokensReply_Data{
			QrCode: qrCode.QrCode,
			State:  qrCode.State,
			Token:  qrCode.Token,
		},
	}, nil
}

func (ds *DouyinService) GetStatusOpenDouyinTokens(ctx context.Context, in *v1.GetStatusOpenDouyinTokensRequest) (*v1.GetStatusOpenDouyinTokensReply, error) {
	qrCodeStatus, err := ds.odtuc.GetStatusOpenDouyinTokens(ctx, in.Token, in.State, in.Timestamp)

	if err != nil {
		return nil, err
	}

	return &v1.GetStatusOpenDouyinTokensReply{
		Code: 200,
		Data: &v1.GetStatusOpenDouyinTokensReply_Data{
			Code:   qrCodeStatus.Code,
			Status: qrCodeStatus.Status,
		},
	}, nil
}

func (ds *DouyinService) CreateOpenDouyinTokens(ctx context.Context, in *v1.CreateOpenDouyinTokensRequest) (*v1.CreateOpenDouyinTokensReply, error) {
	ctx = context.Background()

	if err := ds.odtuc.CreateOpenDouyinTokens(ctx, in.State, in.Code); err != nil {
		return nil, err
	}

	return &v1.CreateOpenDouyinTokensReply{
		Code: 200,
		Data: &v1.CreateOpenDouyinTokensReply_Data{},
	}, nil
}

func (ds *DouyinService) UpdateCooperativeCodeDouyinTokens(ctx context.Context, in *v1.UpdateCooperativeCodeDouyinTokensRequest) (*v1.UpdateCooperativeCodeDouyinTokensReply, error) {
	if err := ds.odtuc.UpdateCooperativeCodeDouyinTokens(ctx, in.ClientKey, in.OpenId, in.CooperativeCode); err != nil {
		return nil, err
	}

	return &v1.UpdateCooperativeCodeDouyinTokensReply{
		Code: 200,
		Data: &v1.UpdateCooperativeCodeDouyinTokensReply_Data{},
	}, nil
}

func (ds *DouyinService) RefreshOpenDouyinTokens(ctx context.Context, in *emptypb.Empty) (*v1.RefreshOpenDouyinTokensReply, error) {
	ds.log.Infof("刷新抖音开放平台授权账户Token, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ds.odtuc.RefreshOpenDouyinTokens(ctx); err != nil {
		return nil, err
	}

	ds.log.Infof("刷新抖音开放平台授权账户Token, 结束时间 %s \n", time.Now())

	return &v1.RefreshOpenDouyinTokensReply{
		Code: 200,
		Data: &v1.RefreshOpenDouyinTokensReply_Data{},
	}, nil
}

func (ds *DouyinService) RenewRefreshTokensOpenDouyinTokens(ctx context.Context, in *emptypb.Empty) (*v1.RenewRefreshTokensOpenDouyinTokensReply, error) {
	ds.log.Infof("刷新抖音开放平台授权账户Refresh Token, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ds.odtuc.RenewRefreshTokensOpenDouyinTokens(ctx); err != nil {
		return nil, err
	}

	ds.log.Infof("刷新抖音开放平台授权账户Refresh Token, 结束时间 %s \n", time.Now())

	return &v1.RenewRefreshTokensOpenDouyinTokensReply{
		Code: 200,
		Data: &v1.RenewRefreshTokensOpenDouyinTokensReply_Data{},
	}, nil
}

func (ds *DouyinService) SyncUserFansOpenDouyinTokens(ctx context.Context, in *emptypb.Empty) (*v1.SyncUserFansOpenDouyinTokensReply, error) {
	ds.log.Infof("同步精选联盟达人粉丝数据, 开始时间 %s \n", time.Now())

	if err := ds.odtuc.SyncUserFansOpenDouyinTokens(ctx); err != nil {
		return nil, err
	}

	ds.log.Infof("同步精选联盟达人粉丝数据, 结束时间 %s \n", time.Now())

	return &v1.SyncUserFansOpenDouyinTokensReply{
		Code: 200,
		Data: &v1.SyncUserFansOpenDouyinTokensReply_Data{},
	}, nil
}

package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	dv1 "interface/api/service/douyin/v1"
	v1 "interface/api/service/weixin/v1"
	"interface/internal/biz"
)

type userOpenDouyinRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserOpenDouyinRepo(data *Data, logger log.Logger) biz.UserOpenDouyinRepo {
	return &userOpenDouyinRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (uodr *userOpenDouyinRepo) GetUrl(ctx context.Context, userId uint64) (*dv1.GetUrlOpenDouyinTokensReply, error) {
	url, err := uodr.data.douyinuc.GetUrlOpenDouyinTokens(ctx, &dv1.GetUrlOpenDouyinTokensRequest{
		UserId: userId,
	})

	if err != nil {
		return nil, err
	}

	return url, err
}

func (uodr *userOpenDouyinRepo) GetTicket(ctx context.Context) (*dv1.GetTicketOpenDouyinTokensReply, error) {
	ticket, err := uodr.data.douyinuc.GetTicketOpenDouyinTokens(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return ticket, err
}

func (uodr *userOpenDouyinRepo) GetQrCode(ctx context.Context, userId uint64) (*dv1.GetQrCodeOpenDouyinTokensReply, error) {
	qrCode, err := uodr.data.douyinuc.GetQrCodeOpenDouyinTokens(ctx, &dv1.GetQrCodeOpenDouyinTokensRequest{
		UserId: userId,
	})

	if err != nil {
		return nil, err
	}

	return qrCode, err
}

func (uodr *userOpenDouyinRepo) GetStatusQrCode(ctx context.Context, token, state, timestamp string) (*dv1.GetStatusOpenDouyinTokensReply, error) {
	qrCodeStatus, err := uodr.data.douyinuc.GetStatusOpenDouyinTokens(ctx, &dv1.GetStatusOpenDouyinTokensRequest{
		Token:     token,
		State:     state,
		Timestamp: timestamp,
	})

	if err != nil {
		return nil, err
	}

	return qrCodeStatus, err
}

func (uodr *userOpenDouyinRepo) List(ctx context.Context, userId, pageNum, pageSize uint64, keyword string) (*v1.ListOpenDouyinUsersReply, error) {
	list, err := uodr.data.weixinuc.ListOpenDouyinUsers(ctx, &v1.ListOpenDouyinUsersRequest{
		PageNum:  pageNum,
		PageSize: pageSize,
		UserId:   userId,
		Keyword:  keyword,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (uodr *userOpenDouyinRepo) Save(ctx context.Context, state, code string) (*dv1.CreateOpenDouyinTokensReply, error) {
	openDouyinToken, err := uodr.data.douyinuc.CreateOpenDouyinTokens(ctx, &dv1.CreateOpenDouyinTokensRequest{
		State: state,
		Code:  code,
	})

	if err != nil {
		return nil, err
	}

	return openDouyinToken, err
}

func (uodr *userOpenDouyinRepo) UpdateCooperativeCode(ctx context.Context, userId, openDouyinUserId uint64, cooperativeCode string) (*v1.UpdateCooperativeCodeOpenDouyinUsersReply, error) {
	list, err := uodr.data.weixinuc.UpdateCooperativeCodeOpenDouyinUsers(ctx, &v1.UpdateCooperativeCodeOpenDouyinUsersRequest{
		UserId:           userId,
		OpenDouyinUserId: openDouyinUserId,
		CooperativeCode:  cooperativeCode,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (uodr *userOpenDouyinRepo) Delete(ctx context.Context, userId, openDouyinUserId uint64) (*v1.DeleteOpenDouyinUsersReply, error) {
	openDouyin, err := uodr.data.weixinuc.DeleteOpenDouyinUsers(ctx, &v1.DeleteOpenDouyinUsersRequest{
		UserId:           userId,
		OpenDouyinUserId: openDouyinUserId,
	})

	if err != nil {
		return nil, err
	}

	return openDouyin, err
}

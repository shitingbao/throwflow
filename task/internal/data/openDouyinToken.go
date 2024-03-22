package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "task/api/service/douyin/v1"
	"task/internal/biz"
)

type openDouyinTokenRepo struct {
	data *Data
	log  *log.Helper
}

func NewOpenDouyinTokenRepo(data *Data, logger log.Logger) biz.OpenDouyinTokenRepo {
	return &openDouyinTokenRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (odtr *openDouyinTokenRepo) Refresh(ctx context.Context) (*v1.RefreshOpenDouyinTokensReply, error) {
	token, err := odtr.data.douyinuc.RefreshOpenDouyinTokens(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return token, err
}

func (odtr *openDouyinTokenRepo) RenewRefresh(ctx context.Context) (*v1.RenewRefreshTokensOpenDouyinTokensReply, error) {
	token, err := odtr.data.douyinuc.RenewRefreshTokensOpenDouyinTokens(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return token, err
}

func (odtr *openDouyinTokenRepo) SyncUserFans(ctx context.Context) (*v1.SyncUserFansOpenDouyinTokensReply, error) {
	userFans, err := odtr.data.douyinuc.SyncUserFansOpenDouyinTokens(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return userFans, err
}

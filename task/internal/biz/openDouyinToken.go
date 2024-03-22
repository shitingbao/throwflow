package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "task/api/service/douyin/v1"
	"task/internal/pkg/tool"
)

type OpenDouyinTokenRepo interface {
	Refresh(context.Context) (*v1.RefreshOpenDouyinTokensReply, error)
	RenewRefresh(context.Context) (*v1.RenewRefreshTokensOpenDouyinTokensReply, error)
	SyncUserFans(context.Context) (*v1.SyncUserFansOpenDouyinTokensReply, error)
}

type OpenDouyinTokenUsecase struct {
	repo OpenDouyinTokenRepo
	log  *log.Helper
}

func NewOpenDouyinTokenUsecase(repo OpenDouyinTokenRepo, logger log.Logger) *OpenDouyinTokenUsecase {
	return &OpenDouyinTokenUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (odtuc *OpenDouyinTokenUsecase) RefreshOpenDouyinTokens(ctx context.Context) (*v1.RefreshOpenDouyinTokensReply, error) {
	token, err := odtuc.repo.Refresh(ctx)

	if err != nil {
		return nil, errors.InternalServer("TASK_REFRESH_OPEN_DOUYIN_TOKEN_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return token, nil
}

func (odtuc *OpenDouyinTokenUsecase) RenewRefreshTokensOpenDouyinTokens(ctx context.Context) (*v1.RenewRefreshTokensOpenDouyinTokensReply, error) {
	token, err := odtuc.repo.RenewRefresh(ctx)

	if err != nil {
		return nil, errors.InternalServer("TASK_RENEW_REFRESH_OPEN_DOUYIN_TOKEN_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return token, nil
}

func (odtuc *OpenDouyinTokenUsecase) SyncUserFansOpenDouyinTokens(ctx context.Context) (*v1.SyncUserFansOpenDouyinTokensReply, error) {
	userFans, err := odtuc.repo.SyncUserFans(ctx)

	if err != nil {
		return nil, errors.InternalServer("TASK_SYNC_USER_FANS_OPEN_DOUYIN_TOKEN_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return userFans, nil
}

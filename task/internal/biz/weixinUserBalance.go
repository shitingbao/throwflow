package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "task/api/service/weixin/v1"
	"task/internal/pkg/tool"
)

type WeixinUserBalanceRepo interface {
	Sync(ctx context.Context) (*v1.SyncUserBalancesReply, error)
}

type WeixinUserBalanceUsecase struct {
	repo WeixinUserBalanceRepo
	log  *log.Helper
}

func NewWeixinUserBalanceUsecase(repo WeixinUserBalanceRepo, logger log.Logger) *WeixinUserBalanceUsecase {
	return &WeixinUserBalanceUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (wubuc *WeixinUserBalanceUsecase) SyncUserBalances(ctx context.Context) (*v1.SyncUserBalancesReply, error) {
	userBalance, err := wubuc.repo.Sync(ctx)

	if err != nil {
		return nil, errors.InternalServer("TASK_SYNC_USER_BALANCE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return userBalance, nil
}

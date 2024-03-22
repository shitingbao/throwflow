package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "task/api/service/weixin/v1"
	"task/internal/pkg/tool"
)

type WeixinUserCommissionRepo interface {
	SyncOrder(ctx context.Context) (*v1.SyncOrderUserCommissionsReply, error)
	SyncCostOrder(ctx context.Context) (*v1.SyncCostOrderUserCommissionsReply, error)
}

type WeixinUserCommissionUsecase struct {
	repo WeixinUserCommissionRepo
	log  *log.Helper
}

func NewWeixinUserCommissionUsecase(repo WeixinUserCommissionRepo, logger log.Logger) *WeixinUserCommissionUsecase {
	return &WeixinUserCommissionUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (wucuc *WeixinUserCommissionUsecase) SyncOrderUserCommissions(ctx context.Context) (*v1.SyncOrderUserCommissionsReply, error) {
	userCommission, err := wucuc.repo.SyncOrder(ctx)

	if err != nil {
		return nil, errors.InternalServer("TASK_SYNC_ORDER_USER_COMMISSION_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return userCommission, nil
}

func (wucuc *WeixinUserCommissionUsecase) SyncCostOrderUserCommissions(ctx context.Context) (*v1.SyncCostOrderUserCommissionsReply, error) {
	userCommission, err := wucuc.repo.SyncCostOrder(ctx)

	if err != nil {
		return nil, errors.InternalServer("TASK_SYNC_COST_ORDER_USER_COMMISSION_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return userCommission, nil
}

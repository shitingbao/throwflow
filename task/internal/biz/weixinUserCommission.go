package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "task/api/service/weixin/v1"
	"task/internal/pkg/tool"
)

type WeixinUserCommissionRepo interface {
	SyncTask(ctx context.Context) (*v1.SyncTaskUserCommissionsReply, error)
}

type WeixinUserCommissionUsecase struct {
	repo WeixinUserCommissionRepo
	log  *log.Helper
}

func NewWeixinUserCommissionUsecase(repo WeixinUserCommissionRepo, logger log.Logger) *WeixinUserCommissionUsecase {
	return &WeixinUserCommissionUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (wucuc *WeixinUserCommissionUsecase) SyncTaskUserCommissions(ctx context.Context) (*v1.SyncTaskUserCommissionsReply, error) {
	userCommission, err := wucuc.repo.SyncTask(ctx)

	if err != nil {
		return nil, errors.InternalServer("TASK_SYNC_TASK_USER_COMMISSION_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return userCommission, nil
}

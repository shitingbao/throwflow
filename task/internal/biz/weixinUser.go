package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "task/api/service/weixin/v1"
	"task/internal/pkg/tool"
)

type WeixinUserRepo interface {
	Sync(context.Context) (*v1.SyncIntegralUsersReply, error)
}

type WeixinUserUsecase struct {
	repo WeixinUserRepo
	log  *log.Helper
}

func NewWeixinUserUsecase(repo WeixinUserRepo, logger log.Logger) *WeixinUserUsecase {
	return &WeixinUserUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (wuuc *WeixinUserUsecase) SyncIntegralUsers(ctx context.Context) (*v1.SyncIntegralUsersReply, error) {
	user, err := wuuc.repo.Sync(ctx)

	if err != nil {
		return nil, errors.InternalServer("TASK_SYNC_INTEGRA_USER_DATA_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return user, nil
}

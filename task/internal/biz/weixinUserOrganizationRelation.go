package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "task/api/service/weixin/v1"
	"task/internal/pkg/tool"
)

type WeixinUserOrganizationRepo interface {
	Sync(ctx context.Context) (*v1.SyncUserOrganizationRelationsReply, error)
}

type WeixinUserOrganizationUsecase struct {
	repo WeixinUserOrganizationRepo
	log  *log.Helper
}

func NewWeixinUserOrganizationUsecase(repo WeixinUserOrganizationRepo, logger log.Logger) *WeixinUserOrganizationUsecase {
	return &WeixinUserOrganizationUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (wuouc *WeixinUserOrganizationUsecase) SyncUserOrganizationRelations(ctx context.Context) (*v1.SyncUserOrganizationRelationsReply, error) {
	userOrganization, err := wuouc.repo.Sync(ctx)

	if err != nil {
		return nil, errors.InternalServer("TASK_SYNC_USER_ORGANIZATION_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return userOrganization, nil
}

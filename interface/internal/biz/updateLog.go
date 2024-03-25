package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/common/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"
)

type UpdateLogRepo interface {
	List(context.Context) (*v1.ListUpdateLogsReply, error)
}

type UpdateLogUsecase struct {
	repo UpdateLogRepo
	conf *conf.Data
	log  *log.Helper
}

func NewUpdateLogUsecase(repo UpdateLogRepo, conf *conf.Data, logger log.Logger) *UpdateLogUsecase {
	return &UpdateLogUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (uluc *UpdateLogUsecase) ListUpdateLogs(ctx context.Context) (*v1.ListUpdateLogsReply, error) {
	materials, err := uluc.repo.List(ctx)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_UPDATE_LOG_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return materials, nil
}

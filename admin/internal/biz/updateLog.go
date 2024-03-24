package biz

import (
	v1 "admin/api/service/common/v1"
	"admin/internal/pkg/tool"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type UpdateLogRepo interface {
	List(context.Context) (*v1.ListUpdateLogsReply, error)
	Save(context.Context, string, string) (*v1.CreateUpdateLogsReply, error)
	Update(context.Context, uint64, string, string) (*v1.UpdateUpdateLogsReply, error)
}

type UpdateLogUsecase struct {
	repo UpdateLogRepo
	log  *log.Helper
}

func NewUpdateLogUsecase(repo UpdateLogRepo, logger log.Logger) *UpdateLogUsecase {
	return &UpdateLogUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uluc *UpdateLogUsecase) ListUpdateLogs(ctx context.Context) (*v1.ListUpdateLogsReply, error) {
	list, err := uluc.repo.List(ctx)

	if err != nil {
		return nil, AdminDataError
	}

	return list, nil
}

func (uluc *UpdateLogUsecase) CreateUpdateLogs(ctx context.Context, name, content string) (*v1.CreateUpdateLogsReply, error) {
	updateLog, err := uluc.repo.Save(ctx, name, content)

	if err != nil {
		return nil, errors.InternalServer("ADMIN_UPDATE_LOG_CREATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return updateLog, nil
}

func (uluc *UpdateLogUsecase) UpdateUpdateLogs(ctx context.Context, id uint64, name, content string) (*v1.UpdateUpdateLogsReply, error) {
	updateLog, err := uluc.repo.Update(ctx, id, name, content)

	if err != nil {
		return nil, errors.InternalServer("ADMIN_UPDATE_LOG_UPDATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return updateLog, nil
}

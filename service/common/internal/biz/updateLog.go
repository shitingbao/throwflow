package biz

import (
	"common/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	CommonUpdateLogNotFound    = errors.InternalServer("COMMON_UPDATE_LOG_NOT_FOUND", "更新日志不存在")
	CommonUpdateLogCreateError = errors.InternalServer("COMMON_UPDATE_LOG_CREATE_ERROR", "更新日志创建失败")
	CommonUpdateLogUpdateError = errors.InternalServer("COMMON_UPDATE_LOG_UPDATE_ERROR", "更新日志更新失败")
)

type UpdateLogRepo interface {
	Get(context.Context, uint64) (*domain.UpdateLog, error)
	List(context.Context) ([]*domain.UpdateLog, error)
	Save(context.Context, *domain.UpdateLog) (*domain.UpdateLog, error)
	Update(context.Context, *domain.UpdateLog) (*domain.UpdateLog, error)
}

type UpdateLogUsecase struct {
	repo UpdateLogRepo
	log  *log.Helper
}

func NewUpdateLogUsecase(repo UpdateLogRepo, logger log.Logger) *UpdateLogUsecase {
	return &UpdateLogUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uluc *UpdateLogUsecase) ListUpdateLogs(ctx context.Context) ([]*domain.UpdateLog, error) {
	list, err := uluc.repo.List(ctx)

	if err != nil {
		return nil, CommonUpdateLogNotFound
	}

	return list, err
}

func (uluc *UpdateLogUsecase) CreateUpdateLogs(ctx context.Context, name, content string) (*domain.UpdateLog, error) {
	inUpdateLog := domain.NewUpdateLog(ctx, name, content)
	inUpdateLog.SetCreateTime(ctx)
	inUpdateLog.SetUpdateTime(ctx)

	updateLog, err := uluc.repo.Save(ctx, inUpdateLog)

	if err != nil {
		return nil, CommonUpdateLogCreateError
	}

	return updateLog, nil
}

func (uluc *UpdateLogUsecase) UpdateUpdateLogs(ctx context.Context, id uint64, name, content string) (*domain.UpdateLog, error) {
	inUpdateLog, err := uluc.repo.Get(ctx, id)

	if err != nil {
		return nil, CommonUpdateLogNotFound
	}

	inUpdateLog.SetName(ctx, name)
	inUpdateLog.SetContent(ctx, content)
	inUpdateLog.SetUpdateTime(ctx)

	updateLog, err := uluc.repo.Update(ctx, inUpdateLog)

	if err != nil {
		return nil, CommonUpdateLogUpdateError
	}

	return updateLog, nil
}

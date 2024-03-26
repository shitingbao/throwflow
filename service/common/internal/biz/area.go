package biz

import (
	"common/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	CommonAreaNotFound = errors.InternalServer("COMMON_AREA_NOT_FOUND", "行政区划不存在")
)

type AreaRepo interface {
	GetByAreaCode(context.Context, uint64) (*domain.Area, error)
	ListByParentAreaCode(ctx context.Context, parentAreaCode uint64) ([]*domain.Area, error)
}

type AreaUsecase struct {
	repo AreaRepo
	log  *log.Helper
}

func NewAreaUsecase(repo AreaRepo, logger log.Logger) *AreaUsecase {
	return &AreaUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (auc *AreaUsecase) ListAreas(ctx context.Context, parentAreaCode uint64) ([]*domain.Area, error) {
	list, err := auc.repo.ListByParentAreaCode(ctx, parentAreaCode)

	if err != nil {
		return nil, CommonDataError
	}

	return list, nil
}

func (auc *AreaUsecase) GetAreas(ctx context.Context, areaCode uint64) (*domain.Area, error) {
	area, err := auc.repo.GetByAreaCode(ctx, areaCode)

	if err != nil {
		return nil, CommonAreaNotFound
	}

	return area, nil
}

package biz

import (
	v1 "admin/api/service/common/v1"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type AreaRepo interface {
	List(context.Context, uint64) (*v1.ListAreasReply, error)
}

type AreaUsecase struct {
	repo AreaRepo
	log  *log.Helper
}

func NewAreaUsecase(repo AreaRepo, logger log.Logger) *AreaUsecase {
	return &AreaUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (auc *AreaUsecase) ListAreas(ctx context.Context, parentAreaCode uint64) (*v1.ListAreasReply, error) {
	list, err := auc.repo.List(ctx, parentAreaCode)

	if err != nil {
		return nil, AdminDataError
	}

	return list, nil
}

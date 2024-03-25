package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/common/v1"
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
		return nil, InterfaceDataError
	}

	return list, nil
}

package biz

import (
	v1 "company/api/service/common/v1"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	CompanyAreaNotFound = errors.NotFound("COMPANY_AREA_NOT_FOUND", "行政区划不存在")
)

type AreaRepo interface {
	GetByAreaCode(context.Context, uint64) (*v1.GetAreasReply, error)
}

type AreaUsecase struct {
	repo AreaRepo
	log  *log.Helper
}

func NewAreaUsecase(repo AreaRepo, logger log.Logger) *AreaUsecase {
	return &AreaUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (auc *AreaUsecase) GetByAreaCode(ctx context.Context, areaCode uint64) (*v1.GetAreasReply, error) {
	area, err := auc.repo.GetByAreaCode(ctx, areaCode)

	if err != nil {
		return nil, CompanyAreaNotFound
	}

	return area, nil
}

package biz

import (
	"company/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	CompanyIndustryNotFound = errors.NotFound("COMPANY_INDUSTRY_NOT_FOUND", "行业分类不存在")
)

type IndustryRepo interface {
	GetById(context.Context, uint64) (*domain.Industry, error)
	List(context.Context, uint8) ([]*domain.Industry, error)
}

type IndustryUsecase struct {
	repo IndustryRepo
	log  *log.Helper
}

func NewIndustryUsecase(repo IndustryRepo, logger log.Logger) *IndustryUsecase {
	return &IndustryUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (iuc *IndustryUsecase) ListIndustries(ctx context.Context) ([]*domain.Industry, error) {
	list, err := iuc.repo.List(ctx, 1)

	if err != nil {
		return nil, CompanyDataError
	}

	return list, nil
}

func (iuc *IndustryUsecase) GetIndustryById(ctx context.Context, id uint64) (*domain.Industry, error) {
	industry, err := iuc.repo.GetById(ctx, id)

	if err != nil {
		return nil, CompanyIndustryNotFound
	}

	return industry, nil
}

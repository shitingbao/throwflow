package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/company/v1"
)

type IndustryRepo interface {
	List(context.Context) (*v1.ListIndustriesReply, error)
}

type IndustryUsecase struct {
	repo IndustryRepo
	log  *log.Helper
}

func NewIndustryUsecase(repo IndustryRepo, logger log.Logger) *IndustryUsecase {
	return &IndustryUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (iuc *IndustryUsecase) ListIndustries(ctx context.Context) (*v1.ListIndustriesReply, error) {
	list, err := iuc.repo.List(ctx)

	if err != nil {
		return nil, InterfaceDataError
	}

	return list, nil
}

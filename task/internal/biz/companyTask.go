package biz

import (
	"context"
	v1 "task/api/service/company/v1"

	"github.com/go-kratos/kratos/v2/log"
)

type CompanyTaskRepo interface {
	Sync(ctx context.Context) (*v1.SyncUpdateCompanyTaskDetailReply, error)
}

type CompanyTaskUsecase struct {
	repo CompanyTaskRepo
	log  *log.Helper
}

func NewCompanyTaskUsecase(repo CompanyTaskRepo, logger log.Logger) *CompanyTaskUsecase {
	return &CompanyTaskUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (cuc *CompanyTaskUsecase) SyncCompanyTasks(ctx context.Context) (*v1.SyncUpdateCompanyTaskDetailReply, error) {
	res, err := cuc.repo.Sync(ctx)

	if err != nil {
		return nil, err
	}

	return res, nil
}

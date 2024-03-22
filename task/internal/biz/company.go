package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "task/api/service/company/v1"
	"task/internal/pkg/tool"
)

type CompanyRepo interface {
	SyncUpdateStatus(context.Context) (*v1.SyncUpdateStatusCompanysReply, error)
}

type CompanyUsecase struct {
	repo CompanyRepo
	log  *log.Helper
}

func NewCompanyUsecase(repo CompanyRepo, logger log.Logger) *CompanyUsecase {
	return &CompanyUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (cuc *CompanyUsecase) SyncUpdateStatusCompanys(ctx context.Context) (*v1.SyncUpdateStatusCompanysReply, error) {
	company, err := cuc.repo.SyncUpdateStatus(ctx)

	if err != nil {
		return nil, errors.InternalServer("TASK_SYNC_UPDATE_STATUS_COMPANY_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return company, nil
}

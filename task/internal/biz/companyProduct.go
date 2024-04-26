package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "task/api/service/company/v1"
	"task/internal/pkg/tool"
)

type CompanyProductRepo interface {
	Sync(ctx context.Context) (*v1.SyncCompanyProductsReply, error)
}

type CompanyProductUsecase struct {
	repo CompanyProductRepo
	log  *log.Helper
}

func NewCompanyProductUsecase(repo CompanyProductRepo, logger log.Logger) *CompanyProductUsecase {
	return &CompanyProductUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (cpuc *CompanyProductUsecase) SyncCompanyProducts(ctx context.Context) (*v1.SyncCompanyProductsReply, error) {
	companyProduct, err := cpuc.repo.Sync(ctx)

	if err != nil {
		return nil, errors.InternalServer("TASK_SYNC_COMPANY_PRODUCT_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return companyProduct, nil
}

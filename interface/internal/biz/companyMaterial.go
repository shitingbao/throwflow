package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/company/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"
)

type CompanyMaterialRepo interface {
	List(context.Context, uint64) (*v1.ListCompanyMaterialsReply, error)
}

type CompanyMaterialUsecase struct {
	repo CompanyMaterialRepo
	conf *conf.Data
	log  *log.Helper
}

func NewCompanyMaterialUsecase(repo CompanyMaterialRepo, conf *conf.Data, logger log.Logger) *CompanyMaterialUsecase {
	return &CompanyMaterialUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (cmuc *CompanyMaterialUsecase) ListCompanyMaterials(ctx context.Context, companyId uint64) (*v1.ListCompanyMaterialsReply, error) {
	companyMaterials, err := cmuc.repo.List(ctx, companyId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_COMPANY_MATERIAL_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return companyMaterials, nil
}

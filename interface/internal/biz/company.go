package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/company/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"
)

type CompanyRepo interface {
	Get(context.Context, uint64) (*v1.GetCompanysReply, error)
}

type CompanyUsecase struct {
	repo CompanyRepo
	conf *conf.Data
	log  *log.Helper
}

func NewCompanyUsecase(repo CompanyRepo, conf *conf.Data, logger log.Logger) *CompanyUsecase {
	return &CompanyUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (cuc *CompanyUsecase) GetMiniQrCodeCompanys(ctx context.Context, companyId uint64) (*v1.GetCompanysReply, error) {
	company, err := cuc.repo.Get(ctx, companyId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_MINI_QR_CODE_COMPANY_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return company, nil
}

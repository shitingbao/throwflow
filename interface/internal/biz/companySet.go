package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/company/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"
	"time"
)

type CompanySetRepo interface {
	Get(context.Context, uint64, string, string) (*v1.GetCompanySetsReply, error)
	Update(context.Context, uint64, string, string) (*v1.UpdateCompanySetsReply, error)
}

type CompanySetUsecase struct {
	repo CompanySetRepo
	conf *conf.Data
	log  *log.Helper
}

func NewCompanySetUsecase(repo CompanySetRepo, conf *conf.Data, logger log.Logger) *CompanySetUsecase {
	return &CompanySetUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (csuc *CompanySetUsecase) GetCompanySets(ctx context.Context, companyId uint64, setKey string) (*v1.GetCompanySetsReply, error) {
	companySet, err := csuc.repo.Get(ctx, companyId, time.Now().Format("2006-01-02"), setKey)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_COMPANY_SET_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return companySet, nil
}

func (csuc *CompanySetUsecase) UpdateCompanySets(ctx context.Context, companyId uint64, setKey, setValue string) (*v1.UpdateCompanySetsReply, error) {
	companySet, err := csuc.repo.Update(ctx, companyId, setKey, setValue)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_UPDATE_COMPANY_SET_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return companySet, nil
}

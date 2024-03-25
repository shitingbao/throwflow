package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/company/v1"
	"interface/internal/biz"
)

type companySetRepo struct {
	data *Data
	log  *log.Helper
}

func NewCompanySetRepo(data *Data, logger log.Logger) biz.CompanySetRepo {
	return &companySetRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (csr *companySetRepo) Get(ctx context.Context, companyId uint64, day, setKey string) (*v1.GetCompanySetsReply, error) {
	companySet, err := csr.data.companyuc.GetCompanySets(ctx, &v1.GetCompanySetsRequest{
		CompanyId: companyId,
		SetKey:    setKey,
		Day:       day,
	})

	if err != nil {
		return nil, err
	}

	return companySet, err
}

func (csr *companySetRepo) Update(ctx context.Context, companyId uint64, setKey, setValue string) (*v1.UpdateCompanySetsReply, error) {
	companySet, err := csr.data.companyuc.UpdateCompanySets(ctx, &v1.UpdateCompanySetsRequest{
		CompanyId: companyId,
		SetKey:    setKey,
		SetValue:  setValue,
	})

	if err != nil {
		return nil, err
	}

	return companySet, err
}

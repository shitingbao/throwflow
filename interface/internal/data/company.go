package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/company/v1"
	"interface/internal/biz"
)

type companyRepo struct {
	data *Data
	log  *log.Helper
}

func NewCompanyRepo(data *Data, logger log.Logger) biz.CompanyRepo {
	return &companyRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cr *companyRepo) Get(ctx context.Context, companyId uint64) (*v1.GetCompanysReply, error) {
	company, err := cr.data.companyuc.GetCompanys(ctx, &v1.GetCompanysRequest{
		Id: companyId,
	})

	if err != nil {
		return nil, err
	}

	return company, err
}

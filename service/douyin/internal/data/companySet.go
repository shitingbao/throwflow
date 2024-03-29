package data

import (
	"context"
	v1 "douyin/api/service/company/v1"
	"douyin/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
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

func (csr *companySetRepo) Get(ctx context.Context, companyId uint64, setkey, day string) (*v1.GetCompanySetsReply, error) {
	companySet, err := csr.data.companyuc.GetCompanySets(ctx, &v1.GetCompanySetsRequest{
		CompanyId: companyId,
		SetKey:    setkey,
		Day:       day,
	})

	if err != nil {
		return nil, err
	}

	return companySet, err
}

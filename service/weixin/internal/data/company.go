package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "weixin/api/service/company/v1"
	"weixin/internal/biz"
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

func (cr *companyRepo) List(ctx context.Context, status string) (*v1.ListCompanysReply, error) {
	list, err := cr.data.companyuc.ListCompanys(ctx, &v1.ListCompanysRequest{
		PageSize: 40,
		Status:   status,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

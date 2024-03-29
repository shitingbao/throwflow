package data

import (
	"context"
	v1 "douyin/api/service/company/v1"
	"douyin/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
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

func (cr *companyRepo) GetById(ctx context.Context, companyId uint64) (*v1.GetCompanysReply, error) {
	list, err := cr.data.companyuc.GetCompanys(ctx, &v1.GetCompanysRequest{
		Id: companyId,
	})

	if err != nil {
		return nil, err
	}

	return list, err
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

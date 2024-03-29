package data

import (
	"context"
	v1 "douyin/api/service/company/v1"
	"douyin/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type companyProductRepo struct {
	data *Data
	log  *log.Helper
}

func NewCompanyProductRepo(data *Data, logger log.Logger) biz.CompanyProductRepo {
	return &companyProductRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cpr *companyProductRepo) GetByProductOutId(ctx context.Context, productOutId uint64, productStatus string) (*v1.GetCompanyProductByProductOutIdsReply, error) {
	prodcut, err := cpr.data.companyuc.GetCompanyProductByProductOutIds(ctx, &v1.GetCompanyProductByProductOutIdsRequest{
		ProductOutId:  productOutId,
		ProductStatus: productStatus,
	})

	if err != nil {
		return nil, err
	}

	return prodcut, err
}

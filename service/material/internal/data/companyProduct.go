package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "material/api/service/company/v1"
	"material/internal/biz"
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

func (cpr *companyProductRepo) GetByProductOutId(ctx context.Context, productOutId uint64) (*v1.GetCompanyProductByProductOutIdsReply, error) {
	companyProduct, err := cpr.data.companyuc.GetCompanyProductByProductOutIds(ctx, &v1.GetCompanyProductByProductOutIdsRequest{
		ProductOutId: productOutId,
	})

	if err != nil {
		return nil, err
	}

	return companyProduct, err
}

func (cpr *companyProductRepo) GetExternal(ctx context.Context, productId uint64) (*v1.GetExternalCompanyProductsReply, error) {
	companyProduct, err := cpr.data.companyuc.GetExternalCompanyProducts(ctx, &v1.GetExternalCompanyProductsRequest{
		ProductId: productId,
	})

	if err != nil {
		return nil, err
	}

	return companyProduct, err
}

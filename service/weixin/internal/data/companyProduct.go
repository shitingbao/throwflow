package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "weixin/api/service/company/v1"
	"weixin/internal/biz"
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

func (cpr *companyProductRepo) GetExternal(ctx context.Context, productOutId uint64) (*v1.GetExternalCompanyProductsReply, error) {
	companyProduct, err := cpr.data.companyuc.GetExternalCompanyProducts(ctx, &v1.GetExternalCompanyProductsRequest{
		ProductId: productOutId,
	})

	if err != nil {
		return nil, err
	}

	return companyProduct, err
}

func (cpr *companyProductRepo) Statistics(ctx context.Context, industryId, categoryId, subCategoryId uint64, keyword string) (*v1.StatisticsCompanyProductsReply, error) {
	list, err := cpr.data.companyuc.StatisticsCompanyProducts(ctx, &v1.StatisticsCompanyProductsRequest{
		IndustryId:    industryId,
		CategoryId:    categoryId,
		SubCategoryId: subCategoryId,
		Keyword:       keyword,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

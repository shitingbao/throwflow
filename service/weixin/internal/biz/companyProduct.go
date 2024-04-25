package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	v1 "weixin/api/service/company/v1"
)

var (
	WeixinCompanyProductNotFound = errors.NotFound("WEIXIN_COMPANY_PRODUCT_NOT_FOUND", "企业商品不存在")
)

type CompanyProductRepo interface {
	GetByProductOutId(context.Context, uint64) (*v1.GetCompanyProductByProductOutIdsReply, error)
	GetExternal(context.Context, uint64) (*v1.GetExternalCompanyProductsReply, error)
	Statistics(context.Context, uint64, uint64, uint64, string) (*v1.StatisticsCompanyProductsReply, error)
}

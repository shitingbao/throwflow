package biz

import (
	v1 "company/api/service/douyin/v1"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	CompanyDoukeProductListError         = errors.InternalServer("COMPANY_DOUKE_PRODUCT_LIST_ERROR", "抖客商品列表获取失败")
	CompanyDoukeProductCategoryListError = errors.InternalServer("COMPANY_DOUKE_PRODUCT_CATEGORY_LIST_ERROR", "抖客商品分类列表获取失败")
	CompanyDoukeProductShareCreateError  = errors.InternalServer("COMPANY_DOUKE_PRODUCT_SHARE_CREATE_ERROR", "该商品已下架，或商家店铺设置了不可推广，暂时无法成本购")
)

type DoukeProductRepo interface {
	List(context.Context, uint64, uint64, uint64, uint64, uint64, uint64) (*v1.ListDoukeProductsReply, error)
	ListByProductId(context.Context, string) (*v1.ListDoukeProductByProductIdsReply, error)
	ListCategory(context.Context, uint64) (*v1.ListCategoryDoukeProductsReply, error)
	Save(context.Context, string, string) (*v1.CreateShareDoukeProductsReply, error)
	Parse(context.Context, string) (*v1.ParseDoukeProductsReply, error)
}

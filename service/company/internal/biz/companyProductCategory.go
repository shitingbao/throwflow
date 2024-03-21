package biz

import (
	"company/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	CompanyProductCategoryNotFound = errors.NotFound("COMPANY_PRODUCT_CATEGORY_NOT_FOUND", "行业分类不存在")
)

type CompanyProductCategoryRepo interface {
	List(context.Context) ([]*domain.CompanyProductCategory, error)
}

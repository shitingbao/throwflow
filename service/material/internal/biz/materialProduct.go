package biz

import (
	"context"
	"material/internal/domain"
)

type MaterialProductRepo interface {
	List(context.Context, int, int, string, string, string, string, []domain.CompanyProductCategory) ([]*domain.MaterialProduct, error)
	Count(context.Context, string, string, string, []domain.CompanyProductCategory) (int64, error)
}

package biz

import (
	"company/internal/domain"
	"context"
)

type CompanyPerformanceMonthlyRepo interface {
	GetByUserIdAndCompanyIdAndUpdateDay(context.Context, uint64, uint64, uint32) (*domain.CompanyPerformanceMonthly, error)
	List(context.Context, int, int, uint64, uint32) ([]*domain.CompanyPerformanceMonthly, error)
	Count(context.Context, uint64, uint32) (int64, error)
	Sum(context.Context, uint64, uint32) (float32, error)
	Save(context.Context, *domain.CompanyPerformanceMonthly) (*domain.CompanyPerformanceMonthly, error)
	Update(context.Context, *domain.CompanyPerformanceMonthly) (*domain.CompanyPerformanceMonthly, error)
	DeleteByCompanyId(context.Context, uint64) error
}

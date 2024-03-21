package biz

import (
	"company/internal/domain"
	"context"
)

type CompanyPerformanceDailyRepo interface {
	GetByUserIdAndCompanyIdAndUpdateDay(context.Context, uint64, uint64, string) (*domain.CompanyPerformanceDaily, error)
	List(context.Context, uint64, uint64, string) ([]*domain.CompanyPerformanceDaily, error)
	Save(context.Context, *domain.CompanyPerformanceDaily) (*domain.CompanyPerformanceDaily, error)
	Update(context.Context, *domain.CompanyPerformanceDaily) (*domain.CompanyPerformanceDaily, error)
	DeleteByCompanyId(context.Context, uint64) error
}

package biz

import (
	"company/internal/domain"
	"context"
)

type CompanyUserCompanyRepo interface {
	GetByPhone(context.Context, string) (*domain.CompanyUserCompany, error)
	Save(context.Context, *domain.CompanyUserCompany) (*domain.CompanyUserCompany, error)
	Update(context.Context, *domain.CompanyUserCompany) (*domain.CompanyUserCompany, error)
	DeleteByCompanyId(context.Context, uint64) error
	Delete(context.Context, *domain.CompanyUserCompany) error
}

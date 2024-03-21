package biz

import (
	"company/internal/domain"
	"context"
)

type CompanyUserWhiteRepo interface {
	GetByPhone(context.Context, string) (*domain.CompanyUserWhite, error)
	Save(context.Context, *domain.CompanyUserWhite) (*domain.CompanyUserWhite, error)
	Delete(context.Context, string) error
}

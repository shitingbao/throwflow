package biz

import (
	"company/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	CompanyMaterialLibraryNotFound = errors.NotFound("COMPANY_MATERIAL_LIBRARY_NOT_FOUND", "企业素材云库不存在")
)

type CompanyMaterialLibraryRepo interface {
	GetByCompanyId(context.Context, uint64, uint64) (*domain.CompanyMaterialLibrary, error)
	GetByCompanyIdAndLibraryName(context.Context, uint64, string) (*domain.CompanyMaterialLibrary, error)
	GetByCompanyIdAndParentIdAndLibraryName(context.Context, uint64, uint64, string) (*domain.CompanyMaterialLibrary, error)
	GetByCompanyIdAndParentIdAndProductId(context.Context, uint64, uint64, uint64) (*domain.CompanyMaterialLibrary, error)
	ListByParentIdAndLibraryType(context.Context, uint64, uint64, uint8) ([]*domain.CompanyMaterialLibrary, error)
	ListByParentId(context.Context, uint64, uint64) ([]*domain.CompanyMaterialLibrary, error)
	CountByParentIdAndLibraryType(context.Context, uint64, uint8) (int64, error)
	Save(context.Context, *domain.CompanyMaterialLibrary) (*domain.CompanyMaterialLibrary, error)
	Update(context.Context, *domain.CompanyMaterialLibrary) (*domain.CompanyMaterialLibrary, error)
}

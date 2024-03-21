package biz

import (
	"company/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	CompanyCompanyUserRoleListError = errors.NotFound("COMPANY_COMPANY_USER_ROLE_LIST_ERROR", "公司用户权限列表获取失败")
)

type CompanyUserRoleRepo interface {
	GetByUserIdAndCompanyIdAndAdvertiserIdAndDay(context.Context, uint64, uint64, uint64, uint32) (*domain.CompanyUserRole, error)
	ListByCompanyIdAndDay(context.Context, uint64, uint32) ([]*domain.CompanyUserRole, error)
	ListByUserIdAndCompanyIdAndDay(context.Context, uint64, uint64, uint32) ([]*domain.CompanyUserRole, error)
	ListByCompanyIdAndAdvertiserIdAndDay(context.Context, uint64, uint64, uint32) ([]*domain.CompanyUserRole, error)
	Save(context.Context, *domain.CompanyUserRole) (*domain.CompanyUserRole, error)
	Update(context.Context, *domain.CompanyUserRole) (*domain.CompanyUserRole, error)
	DeleteByUserIdAndCompanyIdAndAdvertiserIdAndDay(context.Context, uint64, uint64, uint64, uint32) error

	////////////////////和团队绩效相关需要再处理////////////////////////////
	ListByUserIdAndCompanyId(context.Context, uint64, uint64) ([]*domain.CompanyUserRole, error)
	////////////////////////////////////////////
}

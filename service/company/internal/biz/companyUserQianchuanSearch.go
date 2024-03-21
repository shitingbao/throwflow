package biz

import (
	"company/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	CompanyUserQianchuanSearchListError = errors.InternalServer("COMPANY_QIANCHUAN_AD_LIST_ERROR", "千川搜索列表获取失败")
)

type CompanyUserQianchuanSearchRepo interface {
	Get(context.Context, uint64, string, string) (*domain.CompanyUserQianchuanSearch, error)
	List(context.Context, uint64, uint32, uint32, string) ([]*domain.CompanyUserQianchuanSearch, error)
	Save(context.Context, *domain.CompanyUserQianchuanSearch) (*domain.CompanyUserQianchuanSearch, error)
	Update(context.Context, *domain.CompanyUserQianchuanSearch) (*domain.CompanyUserQianchuanSearch, error)
}

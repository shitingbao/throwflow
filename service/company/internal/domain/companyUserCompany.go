package domain

import (
	"context"
	"time"
)

type CompanyUserCompany struct {
	Id         uint64
	Phone      string
	CompanyId  uint64
	CreateTime time.Time
	UpdateTime time.Time
}

func NewCompanyUserCompany(ctx context.Context, companyId uint64, phone string) *CompanyUserCompany {
	return &CompanyUserCompany{
		CompanyId: companyId,
		Phone:     phone,
	}
}

func (cuc *CompanyUserCompany) SetPhone(ctx context.Context, phone string) {
	cuc.Phone = phone
}

func (cuc *CompanyUserCompany) SetCompanyId(ctx context.Context, companyId uint64) {
	cuc.CompanyId = companyId
}

func (cuc *CompanyUserCompany) SetUpdateTime(ctx context.Context) {
	cuc.UpdateTime = time.Now()
}

func (cuc *CompanyUserCompany) SetCreateTime(ctx context.Context) {
	cuc.CreateTime = time.Now()
}

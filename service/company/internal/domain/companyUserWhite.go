package domain

import (
	"context"
	"time"
)

type CompanyUserWhite struct {
	Phone      string
	CreateTime time.Time
	UpdateTime time.Time
}

func NewCompanyUserWhite(ctx context.Context, phone string) *CompanyUserWhite {
	return &CompanyUserWhite{
		Phone: phone,
	}
}

func (cuw *CompanyUserWhite) SetUpdateTime(ctx context.Context) {
	cuw.UpdateTime = time.Now()
}

func (cuw *CompanyUserWhite) SetCreateTime(ctx context.Context) {
	cuw.CreateTime = time.Now()
}

package biz

import (
	"context"
	v1 "douyin/api/service/company/v1"
)

type CompanyProductRepo interface {
	GetByProductOutId(context.Context, uint64, string) (*v1.GetCompanyProductByProductOutIdsReply, error)
}

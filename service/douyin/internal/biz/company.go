package biz

import (
	"context"
	v1 "douyin/api/service/company/v1"
)

type CompanyRepo interface {
	GetById(context.Context, uint64) (*v1.GetCompanysReply, error)
	List(context.Context, string) (*v1.ListCompanysReply, error)
}

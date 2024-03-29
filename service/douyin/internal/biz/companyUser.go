package biz

import (
	"context"
	v1 "douyin/api/service/company/v1"
)

type CompanyUserRepo interface {
	List(context.Context, uint64) (*v1.ListCompanyUsersReply, error)
	DeleteByCompanyIdAndAdvertiserId(context.Context, uint64, uint64) (*v1.DeleteRoleCompanyUsersReply, error)
}

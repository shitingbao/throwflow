package data

import (
	"context"
	v1 "douyin/api/service/company/v1"
	"douyin/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type companyUserRepo struct {
	data *Data
	log  *log.Helper
}

func NewCompanyUserRepo(data *Data, logger log.Logger) biz.CompanyUserRepo {
	return &companyUserRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cur *companyUserRepo) List(ctx context.Context, companyId uint64) (*v1.ListCompanyUsersReply, error) {
	list, err := cur.data.companyuc.ListCompanyUsers(ctx, &v1.ListCompanyUsersRequest{
		CompanyId: companyId,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (cur *companyUserRepo) DeleteByCompanyIdAndAdvertiserId(ctx context.Context, companyId, advertiserId uint64) (*v1.DeleteRoleCompanyUsersReply, error) {
	companyUserRole, err := cur.data.companyuc.DeleteRoleCompanyUsers(ctx, &v1.DeleteRoleCompanyUsersRequest{
		CompanyId:    companyId,
		AdvertiserId: advertiserId,
	})

	if err != nil {
		return nil, err
	}

	return companyUserRole, err
}

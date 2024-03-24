package data

import (
	v1 "admin/api/service/company/v1"
	"admin/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (cur *companyUserRepo) List(ctx context.Context, companyId, pageNum uint64, keyword string) (*v1.ListCompanyUsersReply, error) {
	list, err := cur.data.companyuc.ListCompanyUsers(ctx, &v1.ListCompanyUsersRequest{
		CompanyId: companyId,
		PageNum:   pageNum,
		PageSize:  uint64(cur.data.conf.Database.PageSize),
		Keyword:   keyword,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (cur *companyUserRepo) ListSelect(ctx context.Context) (*v1.ListSelectCompanyUsersReply, error) {
	list, err := cur.data.companyuc.ListSelectCompanyUsers(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (cur *companyUserRepo) ListQianchuanAdvertisers(ctx context.Context, id, companyId uint64) (*v1.ListQianchuanAdvertisersCompanyUsersReply, error) {
	list, err := cur.data.companyuc.ListQianchuanAdvertisersCompanyUsers(ctx, &v1.ListQianchuanAdvertisersCompanyUsersRequest{
		Id:        id,
		CompanyId: companyId,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (cur *companyUserRepo) Save(ctx context.Context, companyId uint64, username, job, phone string, role uint32) (*v1.CreateCompanyUsersReply, error) {
	companyUser, err := cur.data.companyuc.CreateCompanyUsers(ctx, &v1.CreateCompanyUsersRequest{
		CompanyId: companyId,
		Username:  username,
		Job:       job,
		Phone:     phone,
		Role:      role,
	})

	if err != nil {
		return nil, err
	}

	return companyUser, err
}

func (cur *companyUserRepo) Update(ctx context.Context, id, companyId uint64, username, job, phone string, role uint32) (*v1.UpdateCompanyUsersReply, error) {
	companyUser, err := cur.data.companyuc.UpdateCompanyUsers(ctx, &v1.UpdateCompanyUsersRequest{
		Id:        id,
		CompanyId: companyId,
		Username:  username,
		Job:       job,
		Phone:     phone,
		Role:      role,
	})

	if err != nil {
		return nil, err
	}

	return companyUser, err
}

func (cur *companyUserRepo) UpdateStatus(ctx context.Context, id, companyId uint64, status uint32) (*v1.UpdateStatusCompanyUsersReply, error) {
	companyUser, err := cur.data.companyuc.UpdateStatusCompanyUsers(ctx, &v1.UpdateStatusCompanyUsersRequest{
		Id:        id,
		CompanyId: companyId,
		Status:    status,
	})

	if err != nil {
		return nil, err
	}

	return companyUser, err
}

func (cur *companyUserRepo) UpdateWhite(ctx context.Context, id, companyId uint64, isWhite uint32) (*v1.UpdateWhiteCompanyUsersReply, error) {
	companyUser, err := cur.data.companyuc.UpdateWhiteCompanyUsers(ctx, &v1.UpdateWhiteCompanyUsersRequest{
		Id:        id,
		CompanyId: companyId,
		IsWhite:   isWhite,
	})

	if err != nil {
		return nil, err
	}

	return companyUser, err
}

func (cur *companyUserRepo) UpdateRole(ctx context.Context, id, companyId uint64, roleIds string) (*v1.UpdateRoleCompanyUsersReply, error) {
	companyUser, err := cur.data.companyuc.UpdateRoleCompanyUsers(ctx, &v1.UpdateRoleCompanyUsersRequest{
		Id:        id,
		CompanyId: companyId,
		RoleIds:   roleIds,
	})

	if err != nil {
		return nil, err
	}

	return companyUser, err
}

func (cur *companyUserRepo) Delete(ctx context.Context, id, companyId uint64) (*v1.DeleteCompanyUsersReply, error) {
	companyUser, err := cur.data.companyuc.DeleteCompanyUsers(ctx, &v1.DeleteCompanyUsersRequest{
		Id:        id,
		CompanyId: companyId,
	})

	if err != nil {
		return nil, err
	}

	return companyUser, err
}

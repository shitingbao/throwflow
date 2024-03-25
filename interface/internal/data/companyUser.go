package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/company/v1"
	"interface/internal/biz"
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

func (cur *companyUserRepo) ListByPhone(ctx context.Context, phone string) (*v1.ListCompanyUsersByPhoneReply, error) {
	companyUsers, err := cur.data.companyuc.ListCompanyUsersByPhone(ctx, &v1.ListCompanyUsersByPhoneRequest{
		Phone: phone,
	})

	if err != nil {
		return nil, err
	}

	return companyUsers, err
}

func (cur *companyUserRepo) GetCompanyUser(ctx context.Context, token string) (*v1.GetCompanyUserReply, error) {
	companyUser, err := cur.data.companyuc.GetCompanyUser(ctx, &v1.GetCompanyUserRequest{
		Token: token,
	})

	if err != nil {
		return nil, err
	}

	return companyUser, err
}

func (cur *companyUserRepo) ListCompanyUserMenu(ctx context.Context, companyId uint64) (*v1.ListCompanyUserMenuReply, error) {
	companyUserMenus, err := cur.data.companyuc.ListCompanyUserMenu(ctx, &v1.ListCompanyUserMenuRequest{
		CompanyId: companyId,
	})

	if err != nil {
		return nil, err
	}

	return companyUserMenus, err
}

func (cur *companyUserRepo) List(ctx context.Context, companyId, pageNum, pageSize uint64) (*v1.ListCompanyUsersReply, error) {
	companyUsers, err := cur.data.companyuc.ListCompanyUsers(ctx, &v1.ListCompanyUsersRequest{
		CompanyId: companyId,
		PageNum:   pageNum,
		PageSize:  pageSize,
		Keyword:   "",
	})

	if err != nil {
		return nil, err
	}

	return companyUsers, err
}

func (cur *companyUserRepo) Statistics(ctx context.Context, companyId uint64) (*v1.StatisticsCompanyUsersReply, error) {
	companyUsers, err := cur.data.companyuc.StatisticsCompanyUsers(ctx, &v1.StatisticsCompanyUsersRequest{
		CompanyId: companyId,
	})

	if err != nil {
		return nil, err
	}

	return companyUsers, err
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

func (cur *companyUserRepo) Get(ctx context.Context, id, companyId uint64) (*v1.GetCompanyUsersReply, error) {
	companyUser, err := cur.data.companyuc.GetCompanyUsers(ctx, &v1.GetCompanyUsersRequest{
		Id:        id,
		CompanyId: companyId,
	})

	if err != nil {
		return nil, err
	}

	return companyUser, err
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

func (cur *companyUserRepo) ChangeCompanyUserCompany(ctx context.Context, token string, companyId uint64) (*v1.ChangeCompanyUserCompanyReply, error) {
	companyUser, err := cur.data.companyuc.ChangeCompanyUserCompany(ctx, &v1.ChangeCompanyUserCompanyRequest{
		Token:     token,
		CompanyId: companyId,
	})

	if err != nil {
		return nil, err
	}

	return companyUser, err
}

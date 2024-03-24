package service

import (
	v1 "admin/api/admin/v1"
	"admin/internal/biz"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"math"
)

func (as *AdminService) ListCompanyUsers(ctx context.Context, in *v1.ListCompanyUsersRequest) (*v1.ListCompanyUsersReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:companyUser:list"); err != nil {
		return nil, err
	}

	companyUsers, err := as.cuuc.ListCompanyUsers(ctx, in.CompanyId, in.PageNum, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyUsersReply_CompanyUsers, 0)

	for _, companyUser := range companyUsers.Data.List {
		list = append(list, &v1.ListCompanyUsersReply_CompanyUsers{
			Id:        companyUser.Id,
			CompanyId: companyUser.CompanyId,
			Username:  companyUser.Username,
			Job:       companyUser.Job,
			Phone:     companyUser.Phone,
			Role:      companyUser.Role,
			Status:    companyUser.Status,
			IsWhite:   companyUser.IsWhite,
		})
	}

	totalPage := uint64(math.Ceil(float64(companyUsers.Data.Total) / float64(companyUsers.Data.PageSize)))

	return &v1.ListCompanyUsersReply{
		Code: 200,
		Data: &v1.ListCompanyUsersReply_Data{
			PageNum:   companyUsers.Data.PageNum,
			PageSize:  companyUsers.Data.PageSize,
			Total:     companyUsers.Data.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (as *AdminService) ListSelectCompanyUsers(ctx context.Context, in *emptypb.Empty) (*v1.ListSelectCompanyUsersReply, error) {
	if _, err := as.verifyPermission(ctx, "all"); err != nil {
		return nil, err
	}

	selects, err := as.cuuc.ListSelectCompanyUsers(ctx)

	if err != nil {
		return nil, err
	}

	role := make([]*v1.ListSelectCompanyUsersReply_Role, 0)

	for _, lrole := range selects.Data.Role {
		role = append(role, &v1.ListSelectCompanyUsersReply_Role{
			Key:   lrole.Key,
			Value: lrole.Value,
		})
	}

	return &v1.ListSelectCompanyUsersReply{
		Code: 200,
		Data: &v1.ListSelectCompanyUsersReply_Data{
			Role: role,
		},
	}, nil
}

func (as *AdminService) ListQianchuanAdvertisersCompanyUsers(ctx context.Context, in *v1.ListQianchuanAdvertisersCompanyUsersRequest) (*v1.ListQianchuanAdvertisersCompanyUsersReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:companyUser:list"); err != nil {
		return nil, err
	}

	qianchuanAdvertisers, err := as.cuuc.ListQianchuanAdvertisersCompanyUsersRequest(ctx, in.Id, in.CompanyId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListQianchuanAdvertisersCompanyUsersReply_Advertisers, 0)

	for _, qianchuanAdvertiser := range qianchuanAdvertisers.Data.List {
		list = append(list, &v1.ListQianchuanAdvertisersCompanyUsersReply_Advertisers{
			AdvertiserId:   qianchuanAdvertiser.AdvertiserId,
			AdvertiserName: qianchuanAdvertiser.AdvertiserName,
			CompanyName:    qianchuanAdvertiser.CompanyName,
			Status:         qianchuanAdvertiser.Status,
			IsSelect:       qianchuanAdvertiser.IsSelect,
		})
	}

	return &v1.ListQianchuanAdvertisersCompanyUsersReply{
		Code: 200,
		Data: &v1.ListQianchuanAdvertisersCompanyUsersReply_Data{
			List: list,
		},
	}, nil
}

func (as *AdminService) CreateCompanyUsers(ctx context.Context, in *v1.CreateCompanyUsersRequest) (*v1.CreateCompanyUsersReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:companyUser:create"); err != nil {
		return nil, err
	}

	if ok := as.verifyToken(ctx, in.Token); !ok {
		return nil, biz.AdminTokenVerifyError
	}

	companyUser, err := as.cuuc.CreateCompanyUsers(ctx, in.CompanyId, in.Username, in.Job, in.Phone, in.Role)

	if err != nil {
		return nil, err
	}

	return &v1.CreateCompanyUsersReply{
		Code: 200,
		Data: &v1.CreateCompanyUsersReply_Data{
			Id:        companyUser.Data.Id,
			CompanyId: companyUser.Data.CompanyId,
			Username:  companyUser.Data.Username,
			Job:       companyUser.Data.Job,
			Phone:     companyUser.Data.Phone,
			Role:      companyUser.Data.Role,
			Status:    companyUser.Data.Status,
			IsWhite:   companyUser.Data.IsWhite,
		},
	}, nil
}

func (as *AdminService) UpdateCompanyUsers(ctx context.Context, in *v1.UpdateCompanyUsersRequest) (*v1.UpdateCompanyUsersReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:companyUser:update"); err != nil {
		return nil, err
	}

	companyUser, err := as.cuuc.UpdateCompanyUsers(ctx, in.Id, in.CompanyId, in.Username, in.Job, in.Phone, in.Role)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateCompanyUsersReply{
		Code: 200,
		Data: &v1.UpdateCompanyUsersReply_Data{
			Id:        companyUser.Data.Id,
			CompanyId: companyUser.Data.CompanyId,
			Username:  companyUser.Data.Username,
			Job:       companyUser.Data.Job,
			Phone:     companyUser.Data.Phone,
			Role:      companyUser.Data.Role,
			Status:    companyUser.Data.Status,
			IsWhite:   companyUser.Data.IsWhite,
		},
	}, nil
}

func (as *AdminService) UpdateStatusCompanyUsers(ctx context.Context, in *v1.UpdateStatusCompanyUsersRequest) (*v1.UpdateStatusCompanyUsersReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:companyUser:updateStatus"); err != nil {
		return nil, err
	}

	companyUser, err := as.cuuc.UpdateStatusCompanyUsers(ctx, in.Id, in.CompanyId, in.Status)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateStatusCompanyUsersReply{
		Code: 200,
		Data: &v1.UpdateStatusCompanyUsersReply_Data{
			Id:        companyUser.Data.Id,
			CompanyId: companyUser.Data.CompanyId,
			Username:  companyUser.Data.Username,
			Job:       companyUser.Data.Job,
			Phone:     companyUser.Data.Phone,
			Role:      companyUser.Data.Role,
			Status:    companyUser.Data.Status,
			IsWhite:   companyUser.Data.IsWhite,
		},
	}, nil
}

func (as *AdminService) UpdateWhiteCompanyUsers(ctx context.Context, in *v1.UpdateWhiteCompanyUsersRequest) (*v1.UpdateWhiteCompanyUsersReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:companyUser:updateWhite"); err != nil {
		return nil, err
	}

	companyUser, err := as.cuuc.UpdateWhiteCompanyUsers(ctx, in.Id, in.CompanyId, in.IsWhite)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateWhiteCompanyUsersReply{
		Code: 200,
		Data: &v1.UpdateWhiteCompanyUsersReply_Data{
			Id:        companyUser.Data.Id,
			CompanyId: companyUser.Data.CompanyId,
			Username:  companyUser.Data.Username,
			Job:       companyUser.Data.Job,
			Phone:     companyUser.Data.Phone,
			Role:      companyUser.Data.Role,
			Status:    companyUser.Data.Status,
			IsWhite:   companyUser.Data.IsWhite,
		},
	}, nil
}

func (as *AdminService) UpdateRoleCompanyUsers(ctx context.Context, in *v1.UpdateRoleCompanyUsersRequest) (*v1.UpdateRoleCompanyUsersReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:companyUser:updateRole"); err != nil {
		return nil, err
	}

	companyUser, err := as.cuuc.UpdateRoleCompanyUsers(ctx, in.Id, in.CompanyId, in.RoleIds)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateRoleCompanyUsersReply{
		Code: 200,
		Data: &v1.UpdateRoleCompanyUsersReply_Data{
			Id:        companyUser.Data.Id,
			CompanyId: companyUser.Data.CompanyId,
			Username:  companyUser.Data.Username,
			Job:       companyUser.Data.Job,
			Phone:     companyUser.Data.Phone,
			Role:      companyUser.Data.Role,
			Status:    companyUser.Data.Status,
			IsWhite:   companyUser.Data.IsWhite,
		},
	}, nil
}

func (as *AdminService) DeleteCompanyUsers(ctx context.Context, in *v1.DeleteCompanyUsersRequest) (*v1.DeleteCompanyUsersReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:companyUser:delete"); err != nil {
		return nil, err
	}

	if _, err := as.cuuc.DeleteCompanyUsers(ctx, in.Id, in.CompanyId); err != nil {
		return nil, err
	}

	return &v1.DeleteCompanyUsersReply{
		Code: 200,
		Data: &v1.DeleteCompanyUsersReply_Data{},
	}, nil
}

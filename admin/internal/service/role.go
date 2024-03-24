package service

import (
	v1 "admin/api/admin/v1"
	"admin/internal/biz"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (as *AdminService) ListRoles(ctx context.Context, in *emptypb.Empty) (*v1.ListRolesReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:role:list"); err != nil {
		return nil, err
	}

	roles, err := as.ruc.ListRoles(ctx)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListRolesReply_Roles, 0)

	for _, role := range roles {
		list = append(list, &v1.ListRolesReply_Roles{
			Id:          role.Id,
			RoleName:    role.RoleName,
			RoleExplain: role.RoleExplain,
			MenuIds:     role.MenuIds,
			Status:      uint32(role.Status),
		})
	}

	return &v1.ListRolesReply{
		Code: 200,
		Data: &v1.ListRolesReply_Data{
			List: list,
		},
	}, nil
}

func (as *AdminService) CreateRoles(ctx context.Context, in *v1.CreateRolesRequest) (*v1.CreateRolesReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:role:create"); err != nil {
		return nil, err
	}

	if ok := as.verifyToken(ctx, in.Token); !ok {
		return nil, biz.AdminTokenVerifyError
	}

	ids, err := as.verifyMenu(ctx, in.MenuIds)

	if err != nil {
		return nil, err
	}

	role, err := as.ruc.CreateRoles(ctx, in.RoleName, in.RoleExplain, ids, uint8(in.Status))

	if err != nil {
		return nil, err
	}

	return &v1.CreateRolesReply{
		Code: 200,
		Data: &v1.CreateRolesReply_Data{
			Id:          role.Id,
			RoleName:    role.RoleName,
			RoleExplain: role.RoleExplain,
			MenuIds:     role.MenuIds,
			Status:      uint32(role.Status),
		},
	}, nil
}

func (as *AdminService) UpdateRoles(ctx context.Context, in *v1.UpdateRolesRequest) (*v1.UpdateRolesReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:role:update"); err != nil {
		return nil, err
	}

	ids, err := as.verifyMenu(ctx, in.MenuIds)

	if err != nil {
		return nil, err
	}

	role, err := as.ruc.UpdateRoles(ctx, in.Id, in.RoleName, in.RoleExplain, ids, uint8(in.Status))

	if err != nil {
		return nil, err
	}

	return &v1.UpdateRolesReply{
		Code: 200,
		Data: &v1.UpdateRolesReply_Data{
			Id:          role.Id,
			RoleName:    role.RoleName,
			RoleExplain: role.RoleExplain,
			MenuIds:     role.MenuIds,
			Status:      uint32(role.Status),
		},
	}, nil
}

func (as *AdminService) UpdateStatusRoles(ctx context.Context, in *v1.UpdateStatusRolesRequest) (*v1.UpdateStatusRolesReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:role:updateStatus"); err != nil {
		return nil, err
	}

	role, err := as.ruc.UpdateStatusRoles(ctx, in.Id, uint8(in.Status))

	if err != nil {
		return nil, err
	}

	return &v1.UpdateStatusRolesReply{
		Code: 200,
		Data: &v1.UpdateStatusRolesReply_Data{
			Id:          role.Id,
			RoleName:    role.RoleName,
			RoleExplain: role.RoleExplain,
			MenuIds:     role.MenuIds,
			Status:      uint32(role.Status),
		},
	}, nil
}

func (as *AdminService) DeleteRoles(ctx context.Context, in *v1.DeleteRolesRequest) (*v1.DeleteRolesReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:role:delete"); err != nil {
		return nil, err
	}

	err := as.ruc.DeleteRoles(ctx, in.Id)

	if err != nil {
		return nil, err
	}

	return &v1.DeleteRolesReply{
		Code: 200,
		Data: &v1.DeleteRolesReply_Data{},
	}, nil
}

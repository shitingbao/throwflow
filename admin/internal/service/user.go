package service

import (
	v1 "admin/api/admin/v1"
	"admin/internal/biz"
	"admin/internal/pkg/tool"
	"context"
	"encoding/json"
	"google.golang.org/protobuf/types/known/emptypb"
	"math"
	"unicode/utf8"
)

type ListUserMenusReply_Menus struct {
	Id        uint64                      `json:"id"`
	MenuName  string                      `json:"menuName"`
	FileName  string                      `json:"fileName"`
	IconName  string                      `json:"iconName"`
	ChildList []*ListUserMenusReply_Menus `json:"childList"`
}

func (as *AdminService) GetUserInfo(ctx context.Context, in *emptypb.Empty) (*v1.GetUserInfoReply, error) {
	user, err := as.verifyPermission(ctx, "all")

	if err != nil {
		return nil, err
	}

	permissionCodes := make([]*v1.GetUserInfoReply_PermissionCodes, 0)

	for _, menu := range user.Menus {
		if menu.Status == 1 {
			if l := utf8.RuneCountInString(menu.PermissionCode); l > 0 {
				permissionCodes = append(permissionCodes, &v1.GetUserInfoReply_PermissionCodes{
					PermissionCode: menu.PermissionCode,
				})
			}
		}
	}

	return &v1.GetUserInfoReply{
		Code: 200,
		Data: &v1.GetUserInfoReply_Data{
			Username:        user.Username,
			Nickname:        user.Nickname,
			Email:           user.Email,
			LastLoginTime:   tool.TimeToString("2006-01-02 15:04:05", user.LastLoginTime),
			RoleName:        user.Role.RoleName,
			PermissionCodes: permissionCodes,
		},
	}, nil
}

func (as *AdminService) ListUsers(ctx context.Context, in *v1.ListUsersRequest) (*v1.ListUsersReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:user:list"); err != nil {
		return nil, err
	}

	users, err := as.uuc.ListUsers(ctx, in.PageNum)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListUsersReply_Users, 0)

	for _, user := range users.List {
		list = append(list, &v1.ListUsersReply_Users{
			Id:            user.Id,
			Username:      user.Username,
			Nickname:      user.Nickname,
			Email:         user.Email,
			RoleId:        user.RoleId,
			Status:        uint32(user.Status),
			LastLoginIp:   user.LastLoginIp,
			LastLoginTime: tool.TimeToString("2006-01-02 15:04:05", user.LastLoginTime),
		})
	}

	totalPage := uint64(math.Ceil(float64(users.Total) / float64(users.PageSize)))

	return &v1.ListUsersReply{
		Code: 200,
		Data: &v1.ListUsersReply_Data{
			PageNum:   users.PageNum,
			PageSize:  users.PageSize,
			Total:     users.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (as *AdminService) ListUserMenus(ctx context.Context, in *emptypb.Empty) (*v1.ListUserMenusReply, error) {
	user, err := as.verifyPermission(ctx, "all")

	if err != nil {
		return nil, err
	}

	treeMenus := as.treeMenus(ctx, user.Menus, 0)

	list, _ := json.Marshal(treeMenus)

	return &v1.ListUserMenusReply{
		Code: 200,
		Data: &v1.ListUserMenusReply_Data{
			List: string(list),
		},
	}, nil
}

func (as *AdminService) CreateUsers(ctx context.Context, in *v1.CreateUsersRequest) (*v1.CreateUsersReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:user:create"); err != nil {
		return nil, err
	}

	if ok := as.verifyToken(ctx, in.Token); !ok {
		return nil, biz.AdminTokenVerifyError
	}

	if !as.verifyPassword(ctx, in.Password) {
		return nil, biz.AdminValidatorError
	}

	if err := as.verifyRole(ctx, in.RoleId); err != nil {
		return nil, err
	}

	if len := utf8.RuneCountInString(in.Email); len > 0 {
		if !as.verifyEmail(ctx, in.Email) {
			return nil, biz.AdminValidatorError
		}
	}

	user, err := as.uuc.CreateUsers(ctx, in.Username, in.Nickname, in.Email, in.Password, in.RoleId, uint8(in.Status))

	if err != nil {
		return nil, err
	}

	return &v1.CreateUsersReply{
		Code: 200,
		Data: &v1.CreateUsersReply_Data{
			Id:            user.Id,
			Username:      user.Username,
			Nickname:      user.Nickname,
			Email:         user.Email,
			RoleId:        user.RoleId,
			Status:        uint32(user.Status),
			LastLoginIp:   user.LastLoginIp,
			LastLoginTime: tool.TimeToString("2006-01-02 15:04:05", user.LastLoginTime),
		},
	}, nil
}

func (as *AdminService) UpdateUsers(ctx context.Context, in *v1.UpdateUsersRequest) (*v1.UpdateUsersReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:user:update"); err != nil {
		return nil, err
	}

	if len := utf8.RuneCountInString(in.Password); len > 0 {
		if !as.verifyPassword(ctx, in.Password) {
			return nil, biz.AdminValidatorError
		}
	}

	if err := as.verifyRole(ctx, in.RoleId); err != nil {
		return nil, err
	}

	if len := utf8.RuneCountInString(in.Email); len > 0 {
		if !as.verifyEmail(ctx, in.Email) {
			return nil, biz.AdminValidatorError
		}
	}

	user, err := as.uuc.UpdateUsers(ctx, in.Id, in.Username, in.Nickname, in.Email, in.Password, in.RoleId, uint8(in.Status))

	if err != nil {
		return nil, err
	}

	return &v1.UpdateUsersReply{
		Code: 200,
		Data: &v1.UpdateUsersReply_Data{
			Id:            user.Id,
			Username:      user.Username,
			Nickname:      user.Nickname,
			Email:         user.Email,
			RoleId:        user.RoleId,
			Status:        uint32(user.Status),
			LastLoginIp:   user.LastLoginIp,
			LastLoginTime: tool.TimeToString("2006-01-02 15:04:05", user.LastLoginTime),
		},
	}, nil
}

func (as *AdminService) UpdateStatusUsers(ctx context.Context, in *v1.UpdateStatusUsersRequest) (*v1.UpdateStatusUsersReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:user:updateStatus"); err != nil {
		return nil, err
	}

	user, err := as.uuc.UpdateStatusUsers(ctx, in.Id, uint8(in.Status))

	if err != nil {
		return nil, err
	}

	return &v1.UpdateStatusUsersReply{
		Code: 200,
		Data: &v1.UpdateStatusUsersReply_Data{
			Id:            user.Id,
			Username:      user.Username,
			Nickname:      user.Nickname,
			Email:         user.Email,
			RoleId:        user.RoleId,
			Status:        uint32(user.Status),
			LastLoginIp:   user.LastLoginIp,
			LastLoginTime: tool.TimeToString("2006-01-02 15:04:05", user.LastLoginTime),
		},
	}, nil
}

func (as *AdminService) DeleteUsers(ctx context.Context, in *v1.DeleteUsersRequest) (*v1.DeleteUsersReply, error) {
	user, err := as.verifyPermission(ctx, "admin:user:delete")

	if err != nil {
		return nil, err
	}

	err = as.uuc.DeleteUsers(ctx, user.Id, in.Id)

	if err != nil {
		return nil, err
	}

	return &v1.DeleteUsersReply{
		Code: 200,
		Data: &v1.DeleteUsersReply_Data{},
	}, nil
}

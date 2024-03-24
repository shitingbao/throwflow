package service

import (
	v1 "admin/api/admin/v1"
	"admin/internal/biz"
	"admin/internal/pkg/tool"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"unicode/utf8"
)

func (as *AdminService) Login(ctx context.Context, in *v1.LoginRequest) (*v1.LoginReply, error) {
	if !as.verifyPassword(ctx, in.Password) {
		return nil, biz.AdminValidatorError
	}

	loginReply, err := as.uuc.Login(ctx, in.Username, in.Password)

	if err != nil {
		return nil, err
	}

	permissionCodes := make([]*v1.LoginReply_PermissionCodes, 0)

	for _, menu := range loginReply.Menus {
		if menu.Status == 1 {
			if l := utf8.RuneCountInString(menu.PermissionCode); l > 0 {
				permissionCodes = append(permissionCodes, &v1.LoginReply_PermissionCodes{
					PermissionCode: menu.PermissionCode,
				})
			}
		}
	}

	return &v1.LoginReply{
		Code: 200,
		Data: &v1.LoginReply_Data{
			Username:        loginReply.Username,
			Nickname:        loginReply.Nickname,
			Email:           loginReply.Email,
			LastLoginTime:   tool.TimeToString("2006-01-02 15:04:05", loginReply.LastLoginTime),
			Token:           loginReply.Token,
			RoleName:        loginReply.RoleName,
			PermissionCodes: permissionCodes,
		},
	}, nil
}

func (as *AdminService) Logout(ctx context.Context, in *emptypb.Empty) (*v1.LogoutReply, error) {
	if err := as.uuc.Logout(ctx); err != nil {
		return nil, err
	}

	return &v1.LogoutReply{
		Code: 200,
		Data: &v1.LogoutReply_Data{},
	}, nil
}

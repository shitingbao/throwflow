package service

import (
	v1 "admin/api/admin/v1"
	"admin/internal/biz"
	"context"
	"encoding/json"
	"google.golang.org/protobuf/types/known/emptypb"
	"unicode/utf8"
)

type ListMenusReply_Menus struct {
	Id             uint64                  `json:"id"`
	MenuName       string                  `json:"menuName"`
	ParentId       uint64                  `json:"parentId"`
	MenuType       uint32                  `json:"menuType"`
	MenuTypeName   string                  `json:"menuTypeName"`
	Sort           uint32                  `json:"sort"`
	FileName       string                  `json:"fileName"`
	IconName       string                  `json:"iconName"`
	PermissionCode string                  `json:"permissionCode"`
	Status         uint32                  `json:"status"`
	ChildList      []*ListMenusReply_Menus `json:"childList"`
}

func (as *AdminService) ListMenus(ctx context.Context, in *emptypb.Empty) (*v1.ListMenusReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:menu:list"); err != nil {
		return nil, err
	}

	menus, err := as.muc.ListMenus(ctx)

	if err != nil {
		return nil, err
	}

	listMenus := as.listMenus(ctx, menus)

	list, _ := json.Marshal(listMenus)

	return &v1.ListMenusReply{
		Code: 200,
		Data: &v1.ListMenusReply_Data{
			List: string(list),
		},
	}, nil
}

func (as *AdminService) ListPermissionCodes(ctx context.Context, in *emptypb.Empty) (*v1.ListPermissionCodesReply, error) {
	if _, err := as.verifyPermission(ctx, "all"); err != nil {
		return nil, err
	}

	permissionCodes, _ := as.muc.ListPermissionCodes(ctx)

	list := make([]*v1.ListPermissionCodesReply_PermissionCodes, 0)

	for _, permissionCode := range permissionCodes {
		childList := make([]*v1.ListPermissionCodesReply_ChildPermissionCodes, 0)

		for _, childPermissionCode := range permissionCode.ChildPermissionCodes {
			childList = append(childList, &v1.ListPermissionCodesReply_ChildPermissionCodes{
				Name: childPermissionCode.Name,
				Code: childPermissionCode.Code,
			})
		}

		list = append(list, &v1.ListPermissionCodesReply_PermissionCodes{
			Name:      permissionCode.Name,
			Code:      permissionCode.Code,
			ChildList: childList,
		})
	}

	return &v1.ListPermissionCodesReply{
		Code: 200,
		Data: &v1.ListPermissionCodesReply_Data{
			List: list,
		},
	}, nil
}

func (as *AdminService) CreateMenus(ctx context.Context, in *v1.CreateMenusRequest) (*v1.CreateMenusReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:menu:create"); err != nil {
		return nil, err
	}

	if ok := as.verifyToken(ctx, in.Token); !ok {
		return nil, biz.AdminTokenVerifyError
	}

	if in.MenuType == 2 {
		if l := utf8.RuneCountInString(in.FileName); l == 0 {
			return nil, biz.AdminValidatorError
		}

		if l := utf8.RuneCountInString(in.PermissionCode); l > 0 {
			if ok := as.verifyPermissionCode(ctx, in.PermissionCode); !ok {
				return nil, biz.AdminMenuPermissionCodeNotFound
			}
		} else {
			return nil, biz.AdminMenuPermissionCodeNotFound
		}
	} else if in.MenuType == 3 {
		if l := utf8.RuneCountInString(in.PermissionCode); l > 0 {
			if ok := as.verifyPermissionCode(ctx, in.PermissionCode); !ok {
				return nil, biz.AdminMenuPermissionCodeNotFound
			}
		} else {
			return nil, biz.AdminMenuPermissionCodeNotFound
		}
	}

	menu, err := as.muc.CreateMenus(ctx, in.MenuName, in.ParentId, uint8(in.MenuType), uint8(in.Sort), in.FileName, in.IconName, in.PermissionCode, uint8(in.Status))

	if err != nil {
		return nil, err
	}

	return &v1.CreateMenusReply{
		Code: 200,
		Data: &v1.CreateMenusReply_Data{
			Id:             menu.Id,
			MenuName:       menu.MenuName,
			ParentId:       menu.ParentId,
			MenuType:       uint32(menu.MenuType),
			Sort:           uint32(menu.Sort),
			FileName:       menu.FileName,
			IconName:       menu.IconName,
			PermissionCode: menu.PermissionCode,
			Status:         uint32(menu.Status),
		},
	}, nil
}

func (as *AdminService) UpdateMenus(ctx context.Context, in *v1.UpdateMenusRequest) (*v1.UpdateMenusReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:menu:update"); err != nil {
		return nil, err
	}

	if l := utf8.RuneCountInString(in.PermissionCode); l > 0 {
		if ok := as.verifyPermissionCode(ctx, in.PermissionCode); !ok {
			return nil, biz.AdminMenuPermissionCodeNotFound
		}
	}

	menu, err := as.muc.UpdateMenus(ctx, in.Id, in.MenuName, uint8(in.Sort), in.FileName, in.IconName, in.PermissionCode)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateMenusReply{
		Code: 200,
		Data: &v1.UpdateMenusReply_Data{
			Id:             menu.Id,
			MenuName:       menu.MenuName,
			ParentId:       menu.ParentId,
			MenuType:       uint32(menu.MenuType),
			Sort:           uint32(menu.Sort),
			FileName:       menu.FileName,
			IconName:       menu.IconName,
			PermissionCode: menu.PermissionCode,
			Status:         uint32(menu.Status),
		},
	}, nil
}

func (as *AdminService) UpdateStatusMenus(ctx context.Context, in *v1.UpdateStatusMenusRequest) (*v1.UpdateStatusMenusReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:menu:updateStatus"); err != nil {
		return nil, err
	}

	menu, err := as.muc.UpdateStatusMenus(ctx, in.Id, uint8(in.Status))

	if err != nil {
		return nil, err
	}

	return &v1.UpdateStatusMenusReply{
		Code: 200,
		Data: &v1.UpdateStatusMenusReply_Data{
			Id:             menu.Id,
			MenuName:       menu.MenuName,
			ParentId:       menu.ParentId,
			MenuType:       uint32(menu.MenuType),
			Sort:           uint32(menu.Sort),
			FileName:       menu.FileName,
			IconName:       menu.IconName,
			PermissionCode: menu.PermissionCode,
			Status:         uint32(menu.Status),
		},
	}, nil
}

func (as *AdminService) DeleteMenus(ctx context.Context, in *v1.DeleteMenusRequest) (*v1.DeleteMenusReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:menu:delete"); err != nil {
		return nil, err
	}

	err := as.muc.DeleteMenus(ctx, in.Id)

	if err != nil {
		return nil, err
	}

	return &v1.DeleteMenusReply{
		Code: 200,
		Data: &v1.DeleteMenusReply_Data{},
	}, nil
}

package service

import (
	v1 "company/api/company/v1"
	"company/internal/biz"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"unicode/utf8"
)

func (cs *CompanyService) ListMenus(ctx context.Context, in *emptypb.Empty) (*v1.ListMenusReply, error) {
	menus, err := cs.muc.ListMenus(ctx)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListMenusReply_Menu, 0)

	for _, menu := range menus {
		childList := make([]*v1.ListMenusReply_ChildMenu, 0)

		for _, childMenu := range menu.ChildList {
			childList = append(childList, &v1.ListMenusReply_ChildMenu{
				Id:             childMenu.Id,
				MenuName:       childMenu.MenuName,
				ParentId:       childMenu.ParentId,
				Sort:           uint32(childMenu.Sort),
				MenuType:       childMenu.MenuType,
				Filename:       childMenu.Filename,
				Filepath:       childMenu.Filepath,
				PermissionCode: childMenu.PermissionCode,
				Status:         uint32(childMenu.Status),
			})
		}

		list = append(list, &v1.ListMenusReply_Menu{
			Id:             menu.Id,
			MenuName:       menu.MenuName,
			ChildList:      childList,
			ParentId:       menu.ParentId,
			Sort:           uint32(menu.Sort),
			MenuType:       menu.MenuType,
			Filename:       menu.Filename,
			Filepath:       menu.Filepath,
			PermissionCode: menu.PermissionCode,
			Status:         uint32(menu.Status),
		})
	}

	return &v1.ListMenusReply{
		Code: 200,
		Data: &v1.ListMenusReply_Data{
			List: list,
		},
	}, nil
}

func (cs *CompanyService) ListPermissionCodes(ctx context.Context, in *emptypb.Empty) (*v1.ListPermissionCodesReply, error) {
	permissionCodes, _ := cs.muc.ListPermissionCodes(ctx)

	list := make([]*v1.ListPermissionCodesReply_PermissionCode, 0)

	for _, permissionCode := range permissionCodes {
		list = append(list, &v1.ListPermissionCodesReply_PermissionCode{
			Name: permissionCode.Name,
			Code: permissionCode.Code,
		})
	}

	return &v1.ListPermissionCodesReply{
		Code: 200,
		Data: &v1.ListPermissionCodesReply_Data{
			List: list,
		},
	}, nil
}

func (cs *CompanyService) CreateMenus(ctx context.Context, in *v1.CreateMenusRequest) (*v1.CreateMenusReply, error) {
	if l := utf8.RuneCountInString(in.PermissionCode); l > 0 {
		if ok := cs.verifyPermissionCode(ctx, in.PermissionCode); !ok {
			return nil, biz.CompanyMenuPermissionCodeNotFound
		}
	}

	menu, err := cs.muc.CreateMenus(ctx, in.MenuName, in.MenuType, in.Filename, in.Filepath, in.PermissionCode, in.ParentId, uint8(in.Sort), uint8(in.Status))

	if err != nil {
		return nil, err
	}

	return &v1.CreateMenusReply{
		Code: 200,
		Data: &v1.CreateMenusReply_Data{
			Id:             menu.Id,
			MenuName:       menu.MenuName,
			ParentId:       menu.ParentId,
			Sort:           uint32(menu.Sort),
			MenuType:       menu.MenuType,
			Filename:       menu.Filename,
			Filepath:       menu.Filepath,
			PermissionCode: menu.PermissionCode,
			Status:         uint32(menu.Status),
		},
	}, nil
}

func (cs *CompanyService) UpdateMenus(ctx context.Context, in *v1.UpdateMenusRequest) (*v1.UpdateMenusReply, error) {
	if l := utf8.RuneCountInString(in.PermissionCode); l > 0 {
		if ok := cs.verifyPermissionCode(ctx, in.PermissionCode); !ok {
			return nil, biz.CompanyMenuPermissionCodeNotFound
		}
	}

	menu, err := cs.muc.UpdateMenus(ctx, in.Id, in.MenuName, in.MenuType, in.Filename, in.Filepath, in.PermissionCode, uint8(in.Sort))

	if err != nil {
		return nil, err
	}

	return &v1.UpdateMenusReply{
		Code: 200,
		Data: &v1.UpdateMenusReply_Data{
			Id:             menu.Id,
			MenuName:       menu.MenuName,
			ParentId:       menu.ParentId,
			Sort:           uint32(menu.Sort),
			MenuType:       menu.MenuType,
			Filename:       menu.Filename,
			Filepath:       menu.Filepath,
			PermissionCode: menu.PermissionCode,
			Status:         uint32(menu.Status),
		},
	}, nil
}

func (cs *CompanyService) UpdateStatusMenus(ctx context.Context, in *v1.UpdateStatusMenusRequest) (*v1.UpdateStatusMenusReply, error) {
	menu, err := cs.muc.UpdateStatusMenus(ctx, in.Id, uint8(in.Status))

	if err != nil {
		return nil, err
	}

	return &v1.UpdateStatusMenusReply{
		Code: 200,
		Data: &v1.UpdateStatusMenusReply_Data{
			Id:             menu.Id,
			MenuName:       menu.MenuName,
			ParentId:       menu.ParentId,
			Sort:           uint32(menu.Sort),
			MenuType:       menu.MenuType,
			Filename:       menu.Filename,
			Filepath:       menu.Filepath,
			PermissionCode: menu.PermissionCode,
			Status:         uint32(menu.Status),
		},
	}, nil
}

func (cs *CompanyService) DeleteMenus(ctx context.Context, in *v1.DeleteMenusRequest) (*v1.DeleteMenusReply, error) {
	err := cs.muc.DeleteMenus(ctx, in.Id)

	if err != nil {
		return nil, err
	}

	return &v1.DeleteMenusReply{
		Code: 200,
		Data: &v1.DeleteMenusReply_Data{},
	}, nil
}

package service

import (
	v1 "admin/api/admin/v1"
	"admin/internal/biz"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (as *AdminService) ListCompanyMenus(ctx context.Context, in *emptypb.Empty) (*v1.ListCompanyMenusReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:companyMenu:list"); err != nil {
		return nil, err
	}

	companyMenus, err := as.cmuc.ListCompanyMenus(ctx)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyMenusReply_Menus, 0)

	for _, companyMenu := range companyMenus.Data.List {
		childList := make([]*v1.ListCompanyMenusReply_ChildMenus, 0)

		for _, childMenu := range companyMenu.ChildList {
			childList = append(childList, &v1.ListCompanyMenusReply_ChildMenus{
				Id:             childMenu.Id,
				MenuName:       childMenu.MenuName,
				ParentId:       childMenu.ParentId,
				Sort:           childMenu.Sort,
				MenuType:       childMenu.MenuType,
				Filename:       childMenu.Filename,
				Filepath:       childMenu.Filepath,
				PermissionCode: childMenu.PermissionCode,
				Status:         childMenu.Status,
			})
		}

		list = append(list, &v1.ListCompanyMenusReply_Menus{
			Id:             companyMenu.Id,
			MenuName:       companyMenu.MenuName,
			ChildList:      childList,
			ParentId:       companyMenu.ParentId,
			Sort:           companyMenu.Sort,
			MenuType:       companyMenu.MenuType,
			Filename:       companyMenu.Filename,
			Filepath:       companyMenu.Filepath,
			PermissionCode: companyMenu.PermissionCode,
			Status:         companyMenu.Status,
		})
	}

	return &v1.ListCompanyMenusReply{
		Code: 200,
		Data: &v1.ListCompanyMenusReply_Data{
			List: list,
		},
	}, nil
}

func (as *AdminService) ListCompanyPermissionCodes(ctx context.Context, in *emptypb.Empty) (*v1.ListCompanyPermissionCodesReply, error) {
	if _, err := as.verifyPermission(ctx, "all"); err != nil {
		return nil, err
	}

	permissionCodes, _ := as.cmuc.ListCompanyPermissionCodes(ctx)

	list := make([]*v1.ListCompanyPermissionCodesReply_PermissionCodes, 0)

	for _, permissionCode := range permissionCodes.Data.List {
		list = append(list, &v1.ListCompanyPermissionCodesReply_PermissionCodes{
			Name: permissionCode.Name,
			Code: permissionCode.Code,
		})
	}

	return &v1.ListCompanyPermissionCodesReply{
		Code: 200,
		Data: &v1.ListCompanyPermissionCodesReply_Data{
			List: list,
		},
	}, nil
}

func (as *AdminService) CreateCompanyMenus(ctx context.Context, in *v1.CreateCompanyMenusRequest) (*v1.CreateCompanyMenusReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:companyMenu:create"); err != nil {
		return nil, err
	}

	if ok := as.verifyToken(ctx, in.Token); !ok {
		return nil, biz.AdminTokenVerifyError
	}

	companyMenu, err := as.cmuc.CreateCompanyMenus(ctx, in.MenuName, in.MenuType, in.Filename, in.Filepath, in.PermissionCode, in.ParentId, in.Sort, in.Status)

	if err != nil {
		return nil, err
	}

	return &v1.CreateCompanyMenusReply{
		Code: 200,
		Data: &v1.CreateCompanyMenusReply_Data{
			Id:             companyMenu.Data.Id,
			MenuName:       companyMenu.Data.MenuName,
			ParentId:       companyMenu.Data.ParentId,
			Sort:           companyMenu.Data.Sort,
			MenuType:       companyMenu.Data.MenuType,
			Filename:       companyMenu.Data.Filename,
			Filepath:       companyMenu.Data.Filepath,
			PermissionCode: companyMenu.Data.PermissionCode,
			Status:         companyMenu.Data.Status,
		},
	}, nil
}

func (as *AdminService) UpdateCompanyMenus(ctx context.Context, in *v1.UpdateCompanyMenusRequest) (*v1.UpdateCompanyMenusReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:companyMenu:update"); err != nil {
		return nil, err
	}

	companyMenu, err := as.cmuc.UpdateCompanyMenus(ctx, in.Id, in.MenuName, in.MenuType, in.Filename, in.Filepath, in.PermissionCode, in.Sort)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateCompanyMenusReply{
		Code: 200,
		Data: &v1.UpdateCompanyMenusReply_Data{
			Id:             companyMenu.Data.Id,
			MenuName:       companyMenu.Data.MenuName,
			ParentId:       companyMenu.Data.ParentId,
			Sort:           companyMenu.Data.Sort,
			MenuType:       companyMenu.Data.MenuType,
			Filename:       companyMenu.Data.Filename,
			Filepath:       companyMenu.Data.Filepath,
			PermissionCode: companyMenu.Data.PermissionCode,
			Status:         companyMenu.Data.Status,
		},
	}, nil
}

func (as *AdminService) UpdateStatusCompanyMenus(ctx context.Context, in *v1.UpdateStatusCompanyMenusRequest) (*v1.UpdateStatusCompanyMenusReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:companyMenu:updateStatus"); err != nil {
		return nil, err
	}

	companyMenu, err := as.cmuc.UpdateStatusCompanyMenus(ctx, in.Id, in.Status)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateStatusCompanyMenusReply{
		Code: 200,
		Data: &v1.UpdateStatusCompanyMenusReply_Data{
			Id:             companyMenu.Data.Id,
			MenuName:       companyMenu.Data.MenuName,
			ParentId:       companyMenu.Data.ParentId,
			Sort:           companyMenu.Data.Sort,
			MenuType:       companyMenu.Data.MenuType,
			Filename:       companyMenu.Data.Filename,
			Filepath:       companyMenu.Data.Filepath,
			PermissionCode: companyMenu.Data.PermissionCode,
			Status:         companyMenu.Data.Status,
		},
	}, nil
}

func (as *AdminService) DeleteCompanyMenus(ctx context.Context, in *v1.DeleteCompanyMenusRequest) (*v1.DeleteCompanyMenusReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:companyMenu:delete"); err != nil {
		return nil, err
	}

	if _, err := as.cmuc.DeleteCompanyMenus(ctx, in.Id); err != nil {
		return nil, err
	}

	return &v1.DeleteCompanyMenusReply{
		Code: 200,
		Data: &v1.DeleteCompanyMenusReply_Data{},
	}, nil
}

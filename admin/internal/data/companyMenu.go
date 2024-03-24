package data

import (
	v1 "admin/api/service/company/v1"
	"admin/internal/biz"
	"google.golang.org/protobuf/types/known/emptypb"

	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type companyMenuRepo struct {
	data *Data
	log  *log.Helper
}

func NewCompanyMenuRepo(data *Data, logger log.Logger) biz.CompanyMenuRepo {
	return &companyMenuRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cmr *companyMenuRepo) List(ctx context.Context) (*v1.ListMenusReply, error) {
	list, err := cmr.data.companyuc.ListMenus(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (cmr *companyMenuRepo) ListPermissionCodes(ctx context.Context) (*v1.ListPermissionCodesReply, error) {
	list, err := cmr.data.companyuc.ListPermissionCodes(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (cmr *companyMenuRepo) Save(ctx context.Context, menuName, menuType, filename, filepath, permissionCode string, parentId uint64, sort, status uint32) (*v1.CreateMenusReply, error) {
	companyMenu, err := cmr.data.companyuc.CreateMenus(ctx, &v1.CreateMenusRequest{
		MenuName:       menuName,
		ParentId:       parentId,
		Sort:           sort,
		MenuType:       menuType,
		Filename:       filename,
		Filepath:       filepath,
		PermissionCode: permissionCode,
		Status:         status,
	})

	if err != nil {
		return nil, err
	}

	return companyMenu, err
}

func (cmr *companyMenuRepo) Update(ctx context.Context, id uint64, menuName, menuType, filename, filepath, permissionCode string, sort uint32) (*v1.UpdateMenusReply, error) {
	companyMenu, err := cmr.data.companyuc.UpdateMenus(ctx, &v1.UpdateMenusRequest{
		Id:             id,
		MenuName:       menuName,
		Sort:           sort,
		MenuType:       menuType,
		Filename:       filename,
		Filepath:       filepath,
		PermissionCode: permissionCode,
	})

	if err != nil {
		return nil, err
	}

	return companyMenu, err
}

func (cmr *companyMenuRepo) UpdateStatus(ctx context.Context, id uint64, status uint32) (*v1.UpdateStatusMenusReply, error) {
	companyMenu, err := cmr.data.companyuc.UpdateStatusMenus(ctx, &v1.UpdateStatusMenusRequest{
		Id:     id,
		Status: status,
	})

	if err != nil {
		return nil, err
	}

	return companyMenu, err
}

func (cmr *companyMenuRepo) Delete(ctx context.Context, id uint64) (*v1.DeleteMenusReply, error) {
	companyMenu, err := cmr.data.companyuc.DeleteMenus(ctx, &v1.DeleteMenusRequest{
		Id: id,
	})

	if err != nil {
		return nil, err
	}

	return companyMenu, err
}

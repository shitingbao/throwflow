package biz

import (
	v1 "admin/api/service/company/v1"
	"admin/internal/pkg/tool"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	AdminCompanyMenuNotFound       = errors.NotFound("ADMIN_COMPANY_MENU_NOT_FOUND", "saas菜单不存在")
	AdminCompanyMenuParentNotFound = errors.NotFound("ADMIN_COMPANY_MENU_PARENT_NOT_FOUND", "saas父级菜单不存在")
	AdminCompanyMenuCreateError    = errors.InternalServer("ADMIN_COMPANY_MENU_CREATE_ERROR", "saas菜单创建失败")
	AdminCompanyMenuUpdateError    = errors.InternalServer("ADMIN_COMPANY_MENU_UPDATE_ERROR", "saas菜单更新失败")
	AdminCompanyMenuDeleteError    = errors.InternalServer("ADMIN_COMPANY_MENU_DELETE_ERROR", "saas菜单删除失败")
)

type CompanyMenuRepo interface {
	List(context.Context) (*v1.ListMenusReply, error)
	ListPermissionCodes(context.Context) (*v1.ListPermissionCodesReply, error)
	Save(context.Context, string, string, string, string, string, uint64, uint32, uint32) (*v1.CreateMenusReply, error)
	Update(context.Context, uint64, string, string, string, string, string, uint32) (*v1.UpdateMenusReply, error)
	UpdateStatus(context.Context, uint64, uint32) (*v1.UpdateStatusMenusReply, error)
	Delete(context.Context, uint64) (*v1.DeleteMenusReply, error)
}

type CompanyMenuUsecase struct {
	repo CompanyMenuRepo
	log  *log.Helper
}

func NewCompanyMenuUsecase(repo CompanyMenuRepo, logger log.Logger) *CompanyMenuUsecase {
	return &CompanyMenuUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (cmuc *CompanyMenuUsecase) ListCompanyMenus(ctx context.Context) (*v1.ListMenusReply, error) {
	list, err := cmuc.repo.List(ctx)

	if err != nil {
		return nil, AdminDataError
	}

	return list, nil
}

func (cmuc *CompanyMenuUsecase) ListCompanyPermissionCodes(ctx context.Context) (*v1.ListPermissionCodesReply, error) {
	list, err := cmuc.repo.ListPermissionCodes(ctx)

	if err != nil {
		return nil, AdminDataError
	}

	return list, nil
}

func (cmuc *CompanyMenuUsecase) CreateCompanyMenus(ctx context.Context, menuName, menuType, filename, filepath, permissionCode string, parentId uint64, sort, status uint32) (*v1.CreateMenusReply, error) {
	companyMenu, err := cmuc.repo.Save(ctx, menuName, menuType, filename, filepath, permissionCode, parentId, sort, status)

	if err != nil {
		return nil, errors.InternalServer("ADMIN_COMPANY_MENU_CREATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return companyMenu, nil
}

func (cmuc *CompanyMenuUsecase) UpdateCompanyMenus(ctx context.Context, id uint64, menuName, menuType, filename, filepath, permissionCode string, sort uint32) (*v1.UpdateMenusReply, error) {
	companyMenu, err := cmuc.repo.Update(ctx, id, menuName, menuType, filename, filepath, permissionCode, sort)

	if err != nil {
		return nil, errors.InternalServer("ADMIN_COMPANY_MENU_UPDATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return companyMenu, nil
}

func (cmuc *CompanyMenuUsecase) UpdateStatusCompanyMenus(ctx context.Context, id uint64, status uint32) (*v1.UpdateStatusMenusReply, error) {
	companyMenu, err := cmuc.repo.UpdateStatus(ctx, id, status)

	if err != nil {
		return nil, errors.InternalServer("ADMIN_COMPANY_MENU_UPDATE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return companyMenu, nil
}

func (cmuc *CompanyMenuUsecase) DeleteCompanyMenus(ctx context.Context, id uint64) (*v1.DeleteMenusReply, error) {
	companyMenu, err := cmuc.repo.Delete(ctx, id)

	if err != nil {
		return nil, errors.InternalServer("ADMIN_COMPANY_MENU_DELETE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	return companyMenu, nil
}

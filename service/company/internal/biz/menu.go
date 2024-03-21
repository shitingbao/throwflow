package biz

import (
	"company/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	CompanyMenuNotFound               = errors.NotFound("COMPANY_MENU_NOT_FOUND", "菜单不存在")
	CompanyMenuParentNotFound         = errors.NotFound("COMPANY_MENU_PARENT_NOT_FOUND", "父级菜单不存在")
	CompanyMenuPermissionCodeNotFound = errors.NotFound("COMPANY_MENU_PERMISSION_CODE_NOT_FOUND", "权限标识不存在")
	CompanyMenuCreateError            = errors.InternalServer("COMPANY_MENU_CREATE_ERROR", "菜单创建失败")
	CompanyMenuUpdateError            = errors.InternalServer("COMPANY_MENU_UPDATE_ERROR", "菜单更新失败")
	CompanyMenuDeleteError            = errors.InternalServer("COMPANY_MENU_DELETE_ERROR", "菜单删除失败")
)

type MenuRepo interface {
	GetById(context.Context, uint64) (*domain.Menu, error)
	ListByIds(context.Context, []string) ([]*domain.Menu, error)
	ListByParentId(context.Context, uint64) ([]*domain.Menu, error)
	Save(context.Context, *domain.Menu) (*domain.Menu, error)
	Update(context.Context, *domain.Menu) (*domain.Menu, error)
	Delete(context.Context, *domain.Menu) error
}

type MenuUsecase struct {
	repo MenuRepo
	tm   Transaction
	log  *log.Helper
}

func NewMenuUsecase(repo MenuRepo, tm Transaction, logger log.Logger) *MenuUsecase {
	return &MenuUsecase{repo: repo, tm: tm, log: log.NewHelper(logger)}
}

func (muc *MenuUsecase) GetMenuById(ctx context.Context, id uint64) (*domain.Menu, error) {
	menu, err := muc.getMenuById(ctx, id)

	if err != nil {
		return nil, CompanyMenuNotFound
	}

	return menu, nil
}

func (muc *MenuUsecase) ListMenus(ctx context.Context) ([]*domain.Menu, error) {
	list, err := muc.repo.ListByParentId(ctx, 0)

	if err != nil {
		return nil, CompanyDataError
	}

	for _, menu := range list {
		menu.ChildList, _ = muc.repo.ListByParentId(ctx, menu.Id)
	}

	return list, nil
}

func (muc *MenuUsecase) ListPermissionCodes(ctx context.Context) ([]*domain.PermissionCodes, error) {
	return domain.NewPermissionCodes(), nil
}

func (muc *MenuUsecase) CreateMenus(ctx context.Context, menuName, menuType, filename, filepath, permissionCode string, parentId uint64, sort, status uint8) (*domain.Menu, error) {
	if parentId > 0 {
		if _, err := muc.getMenuById(ctx, parentId); err != nil {
			return nil, CompanyMenuParentNotFound
		}
	}

	inMenu := domain.NewMenu(ctx, menuName, menuType, filename, filepath, permissionCode, parentId, sort, status)
	inMenu.SetCreateTime(ctx)
	inMenu.SetUpdateTime(ctx)

	menu, err := muc.repo.Save(ctx, inMenu)

	if err != nil {
		return nil, CompanyMenuCreateError
	}

	return menu, nil
}

func (muc *MenuUsecase) UpdateMenus(ctx context.Context, id uint64, menuName, menuType, filename, filepath, permissionCode string, sort uint8) (*domain.Menu, error) {
	inMenu, err := muc.getMenuById(ctx, id)

	if err != nil {
		return nil, CompanyMenuNotFound
	}

	inMenu.SetMenuName(ctx, menuName)
	inMenu.SetPermissionCode(ctx, permissionCode)
	inMenu.SetSort(ctx, sort)
	inMenu.SetMenuType(ctx, menuType)
	inMenu.SetFilename(ctx, filename)
	inMenu.SetFilepath(ctx, filepath)
	inMenu.SetUpdateTime(ctx)

	menu, err := muc.repo.Update(ctx, inMenu)

	if err != nil {
		return nil, CompanyMenuUpdateError
	}

	return menu, nil
}

func (muc *MenuUsecase) UpdateStatusMenus(ctx context.Context, id uint64, status uint8) (*domain.Menu, error) {
	inMenu, err := muc.getMenuById(ctx, id)

	if err != nil {
		return nil, CompanyMenuNotFound
	}

	inMenu.SetStatus(ctx, status)
	inMenu.SetUpdateTime(ctx)

	var menu *domain.Menu

	err = muc.tm.InTx(ctx, func(ctx context.Context) error {
		menu, err = muc.repo.Update(ctx, inMenu)

		if err != nil {
			return err
		}

		for _, childMenu := range inMenu.ChildList {
			childMenu.SetStatus(ctx, status)
			childMenu.SetUpdateTime(ctx)

			_, err = muc.repo.Update(ctx, childMenu)

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, CompanyMenuUpdateError
	}

	return menu, nil
}

func (muc *MenuUsecase) DeleteMenus(ctx context.Context, id uint64) error {
	inMenu, err := muc.getMenuById(ctx, id)

	if err != nil {
		return CompanyMenuNotFound
	}

	err = muc.tm.InTx(ctx, func(ctx context.Context) error {
		if err := muc.repo.Delete(ctx, inMenu); err != nil {
			return err
		}

		for _, childMenu := range inMenu.ChildList {
			if err := muc.repo.Delete(ctx, childMenu); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return CompanyMenuDeleteError
	}

	return nil
}

func (muc *MenuUsecase) getMenuById(ctx context.Context, id uint64) (*domain.Menu, error) {
	menu, err := muc.repo.GetById(ctx, id)

	if err != nil {
		return nil, err
	}

	menu.ChildList, _ = muc.repo.ListByParentId(ctx, menu.Id)

	return menu, nil
}

package biz

import (
	"admin/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	AdminMenuNotFound               = errors.NotFound("ADMIN_MENU_NOT_FOUND", "菜单不存在")
	AdminMenuParentNotFound         = errors.NotFound("ADMIN_MENU_PARENT_NOT_FOUND", "父级菜单不存在")
	AdminMenuPermissionCodeNotFound = errors.NotFound("ADMIN_MENU_PERMISSION_CODE_NOT_FOUND", "权限标识不存在")
	AdminMenuCreateError            = errors.InternalServer("ADMIN_MENU_CREATE_ERROR", "菜单创建失败")
	AdminMenuUpdateError            = errors.InternalServer("ADMIN_MENU_UPDATE_ERROR", "菜单更新失败")
	AdminMenuDeleteError            = errors.InternalServer("ADMIN_MENU_DELETE_ERROR", "菜单删除失败")
	AdminMenuRoleNotDelete          = errors.InternalServer("ADMIN_MENU_ROLE_NOT_DISABLED", "该菜单下有角色不能被删除")
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
	repo  MenuRepo
	rrepo RoleRepo
	tm    Transaction
	log   *log.Helper
}

func NewMenuUsecase(repo MenuRepo, rrepo RoleRepo, tm Transaction, logger log.Logger) *MenuUsecase {
	return &MenuUsecase{repo: repo, rrepo: rrepo, tm: tm, log: log.NewHelper(logger)}
}

func (muc *MenuUsecase) GetMenuById(ctx context.Context, id uint64) (*domain.Menu, error) {
	menu, err := muc.getMenuById(ctx, id)

	if err != nil {
		return nil, AdminMenuNotFound
	}

	return menu, nil
}

func (muc *MenuUsecase) ListMenus(ctx context.Context) ([]*domain.Menu, error) {
	list, err := muc.repo.ListByParentId(ctx, 0)

	if err != nil {
		return nil, AdminDataError
	}

	for _, menu := range list {
		muc.getChildMenu(ctx, menu)
	}

	return list, nil
}

func (muc *MenuUsecase) ListPermissionCodes(ctx context.Context) ([]*domain.PermissionCodes, error) {
	return domain.NewPermissionCodes(), nil
}

func (muc *MenuUsecase) CreateMenus(ctx context.Context, menuName string, parentId uint64, menuType, sort uint8, fileName, iconName, permissionCode string, status uint8) (*domain.Menu, error) {
	if parentId > 0 {
		parentMenu, err := muc.getMenuById(ctx, parentId)

		if err != nil {
			return nil, AdminMenuParentNotFound
		}

		if (parentMenu.MenuType == 2 || parentMenu.MenuType == 3) && menuType != 3 {
			return nil, AdminValidatorError
		} else if parentMenu.MenuType == 1 && menuType == 3 {
			return nil, AdminValidatorError
		}
	}

	inMenu := domain.NewMenu(ctx, menuName, parentId, menuType, sort, fileName, iconName, permissionCode, status)
	inMenu.SetCreateTime(ctx)
	inMenu.SetUpdateTime(ctx)

	menu, err := muc.repo.Save(ctx, inMenu)

	if err != nil {
		return nil, AdminMenuCreateError
	}

	return menu, nil
}

func (muc *MenuUsecase) UpdateMenus(ctx context.Context, id uint64, menuName string, sort uint8, fileName, iconName, permissionCode string) (*domain.Menu, error) {
	inMenu, err := muc.getMenuById(ctx, id)

	if err != nil {
		return nil, AdminMenuNotFound
	}

	inMenu.SetMenuName(ctx, menuName)
	inMenu.SetSort(ctx, sort)
	inMenu.SetFileName(ctx, fileName)
	inMenu.SetIconName(ctx, iconName)
	inMenu.SetPermissionCode(ctx, permissionCode)
	inMenu.SetUpdateTime(ctx)

	menu, err := muc.repo.Update(ctx, inMenu)

	if err != nil {
		return nil, AdminMenuUpdateError
	}

	return menu, nil
}

func (muc *MenuUsecase) UpdateStatusMenus(ctx context.Context, id uint64, status uint8) (*domain.Menu, error) {
	inMenu, err := muc.getMenuById(ctx, id)

	if err != nil {
		return nil, AdminMenuNotFound
	}

	inMenu.SetStatus(ctx, status)
	inMenu.SetUpdateTime(ctx)

	var menu *domain.Menu

	err = muc.tm.InTx(ctx, func(ctx context.Context) error {
		menu, err = muc.repo.Update(ctx, inMenu)

		if err != nil {
			return err
		}

		err = muc.updateStatusChildMenu(ctx, inMenu, status)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, AdminMenuUpdateError
	}

	return menu, nil
}

func (muc *MenuUsecase) DeleteMenus(ctx context.Context, id uint64) error {
	inMenu, err := muc.getMenuById(ctx, id)

	if err != nil {
		return AdminMenuNotFound
	}

	if len(inMenu.Roles) > 0 {
		return AdminMenuRoleNotDelete
	}

	err = muc.tm.InTx(ctx, func(ctx context.Context) error {
		if err := muc.repo.Delete(ctx, inMenu); err != nil {
			return err
		}

		err = muc.DeleteChildMenu(ctx, inMenu)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return AdminMenuDeleteError
	}

	return nil
}

func (muc *MenuUsecase) getMenuById(ctx context.Context, id uint64) (*domain.Menu, error) {
	menu, err := muc.repo.GetById(ctx, id)

	if err != nil {
		return nil, err
	}

	roles := make([]*domain.Role, 0)

	muc.getChildMenu(ctx, menu)
	muc.getRoles(ctx, menu, &roles)

	menu.Roles = roles

	return menu, nil
}

func (muc *MenuUsecase) getRoles(ctx context.Context, menu *domain.Menu, roles *[]*domain.Role) {
	if menu == nil {
		return
	}

	if lroles, err := muc.rrepo.ListByMenuId(ctx, menu.Id); err == nil {
		for _, lrole := range lroles {
			isNotExist := true

			for _, role := range *roles {
				if lrole.Id == role.Id {
					isNotExist = false

					break
				}
			}

			if isNotExist {
				*roles = append(*roles, lrole)
			}
		}
	}

	for _, childMenu := range menu.ChildList {
		muc.getRoles(ctx, childMenu, roles)
	}
}

func (muc *MenuUsecase) getChildMenu(ctx context.Context, menu *domain.Menu) {
	if menus, err := muc.repo.ListByParentId(ctx, menu.Id); err == nil {
		if len(menus) == 0 {
			return
		}

		childList := make([]*domain.Menu, 0)

		for _, lmenu := range menus {
			muc.getChildMenu(ctx, lmenu)

			childList = append(childList, lmenu)
		}

		menu.ChildList = childList
	}
}

func (muc *MenuUsecase) updateStatusChildMenu(ctx context.Context, menu *domain.Menu, status uint8) error {
	if menu == nil {
		return nil
	}

	for _, childMenu := range menu.ChildList {
		childMenu.SetStatus(ctx, status)
		childMenu.SetUpdateTime(ctx)

		_, err := muc.repo.Update(ctx, childMenu)

		if err != nil {
			return err
		}

		muc.updateStatusChildMenu(ctx, childMenu, status)
	}

	return nil
}

func (muc *MenuUsecase) DeleteChildMenu(ctx context.Context, menu *domain.Menu) error {
	if menu == nil {
		return nil
	}

	for _, childMenu := range menu.ChildList {
		if err := muc.repo.Delete(ctx, childMenu); err != nil {
			return err
		}

		muc.DeleteChildMenu(ctx, childMenu)
	}

	return nil
}

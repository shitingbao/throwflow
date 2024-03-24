package biz

import (
	"admin/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
)

var (
	AdminRoleNotFound        = errors.NotFound("ADMIN_ROLE_NOT_FOUND", "角色不存在")
	AdminRoleDisabled        = errors.NotFound("ADMIN_ROLE_DISABLED", "角色被禁用")
	AdminRoleMenuNotFound    = errors.InternalServer("ADMIN_ROLE_MENU_NOT_FOUND", "角色关联菜单不存在")
	AdminRoleCreateError     = errors.InternalServer("ADMIN_ROLE_CREATE_ERROR", "角色创建失败")
	AdminRoleUpdateError     = errors.InternalServer("ADMIN_ROLE_UPDATE_ERROR", "角色更新失败")
	AdminRoleDeleteError     = errors.InternalServer("ADMIN_ROLE_DELETE_ERROR", "角色删除失败")
	AdminRoleUserNotDisbaled = errors.InternalServer("ADMIN_ROLE_USER_NOT_DISABLED", "该角色下有用户不能被禁用")
	AdminRoleUserNotDelete   = errors.InternalServer("ADMIN_ROLE_USER_NOT_DELETE", "该角色下有用户不能被删除")
)

type RoleRepo interface {
	GetById(context.Context, uint64) (*domain.Role, error)
	List(context.Context) ([]*domain.Role, error)
	ListByMenuId(context.Context, uint64) ([]*domain.Role, error)
	Save(context.Context, *domain.Role) (*domain.Role, error)
	Update(context.Context, *domain.Role) (*domain.Role, error)
	Delete(context.Context, *domain.Role) error
}

type RoleUsecase struct {
	repo  RoleRepo
	mrepo MenuRepo
	log   *log.Helper
}

func NewRoleUsecase(repo RoleRepo, mrepo MenuRepo, logger log.Logger) *RoleUsecase {
	return &RoleUsecase{repo: repo, mrepo: mrepo, log: log.NewHelper(logger)}
}

func (ruc *RoleUsecase) GetRoleById(ctx context.Context, id uint64) (*domain.Role, error) {
	role, err := ruc.getRoleById(ctx, id)

	if err != nil {
		return nil, AdminRoleNotFound
	}

	return role, nil
}

func (ruc *RoleUsecase) ListRoles(ctx context.Context) ([]*domain.Role, error) {
	list, err := ruc.repo.List(ctx)

	if err != nil {
		return nil, AdminDataError
	}

	return list, nil
}

func (ruc *RoleUsecase) CreateRoles(ctx context.Context, RoleName, RoleExplain string, ids []string, status uint8) (*domain.Role, error) {
	inRole := domain.NewRole(ctx, RoleName, RoleExplain, strings.Join(ids, ","), status)
	inRole.SetCreateTime(ctx)
	inRole.SetUpdateTime(ctx)

	role, err := ruc.repo.Save(ctx, inRole)

	if err != nil {
		return nil, AdminRoleCreateError
	}

	return role, nil
}

func (ruc *RoleUsecase) UpdateRoles(ctx context.Context, id uint64, roleName, roleExplain string, ids []string, status uint8) (*domain.Role, error) {
	inRole, err := ruc.getRoleById(ctx, id)

	if err != nil {
		return nil, AdminRoleNotFound
	}

	inRole.SetRoleName(ctx, roleName)
	inRole.SetRoleExplain(ctx, roleExplain)
	inRole.SetMenuIds(ctx, strings.Join(ids, ","))
	inRole.SetUpdateTime(ctx)

	if status == 0 && len(inRole.Users) > 0 {
		return nil, AdminRoleUserNotDisbaled
	}

	inRole.SetStatus(ctx, status)

	role, err := ruc.repo.Update(ctx, inRole)

	if err != nil {
		return nil, AdminRoleUpdateError
	}

	return role, nil
}

func (ruc *RoleUsecase) UpdateStatusRoles(ctx context.Context, id uint64, status uint8) (*domain.Role, error) {
	inRole, err := ruc.getRoleById(ctx, id)

	if err != nil {
		return nil, AdminRoleNotFound
	}

	if status == 0 && len(inRole.Users) > 0 {
		return nil, AdminRoleUserNotDisbaled
	}

	inRole.SetStatus(ctx, status)
	inRole.SetUpdateTime(ctx)

	role, err := ruc.repo.Update(ctx, inRole)

	if err != nil {
		return nil, AdminRoleUpdateError
	}

	return role, nil
}

func (ruc *RoleUsecase) DeleteRoles(ctx context.Context, id uint64) error {
	inRole, err := ruc.getRoleById(ctx, id)

	if err != nil {
		return AdminRoleNotFound
	}

	if len(inRole.Users) > 0 {
		return AdminRoleUserNotDelete
	}

	if err := ruc.repo.Delete(ctx, inRole); err != nil {
		return AdminRoleDeleteError
	}

	return nil
}

func (ruc *RoleUsecase) getRoleById(ctx context.Context, id uint64) (*domain.Role, error) {
	role, err := ruc.repo.GetById(ctx, id)

	if err != nil {
		return nil, err
	}

	role.Menus, _ = ruc.mrepo.ListByIds(ctx, strings.Split(role.MenuIds, ","))

	return role, nil
}

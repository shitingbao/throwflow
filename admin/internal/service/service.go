package service

import (
	v1 "admin/api/admin/v1"
	"admin/internal/biz"
	"admin/internal/conf"
	"admin/internal/domain"
	"context"
	"github.com/dlclark/regexp2"
	"github.com/google/wire"
	"net/mail"
	"strconv"
	"strings"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewAdminService)

type AdminService struct {
	v1.UnimplementedAdminServer

	uuc  *biz.UserUsecase
	ruc  *biz.RoleUsecase
	tuc  *biz.TokenUsecase
	muc  *biz.MenuUsecase
	sluc *biz.SmsLogUsecase
	ocuc *biz.OceanengineConfigUsecase
	cmuc *biz.CompanyMenuUsecase
	cluc *biz.ClueUsecase
	iuc  *biz.IndustryUsecase
	couc *biz.CompanyUsecase
	cuuc *biz.CompanyUserUsecase
	auc  *biz.AreaUsecase
	uluc *biz.UpdateLogUsecase

	conf *conf.Data
}

func NewAdminService(uuc *biz.UserUsecase, ruc *biz.RoleUsecase, tuc *biz.TokenUsecase, muc *biz.MenuUsecase, sluc *biz.SmsLogUsecase, ocuc *biz.OceanengineConfigUsecase, cmuc *biz.CompanyMenuUsecase, cluc *biz.ClueUsecase, iuc *biz.IndustryUsecase, couc *biz.CompanyUsecase, cuuc *biz.CompanyUserUsecase, auc *biz.AreaUsecase, uluc *biz.UpdateLogUsecase, conf *conf.Data) *AdminService {
	return &AdminService{uuc: uuc, ruc: ruc, tuc: tuc, muc: muc, sluc: sluc, ocuc: ocuc, cmuc: cmuc, cluc: cluc, iuc: iuc, couc: couc, cuuc: cuuc, auc: auc, uluc: uluc, conf: conf}
}

func (as *AdminService) verifyPermission(ctx context.Context, permissionCode string) (*domain.User, error) {
	user, err := as.uuc.GetUserByToken(ctx)

	if err != nil {
		return nil, err
	}

	if !user.VerifyPermission(ctx, permissionCode) {
		return nil, biz.AdminLoginPermissionError
	}

	return user, nil
}

func (as *AdminService) verifyPassword(ctx context.Context, password string) bool {
	passwordRule := regexp2.MustCompile("^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)[a-zA-Z\\d]{8,}$", 0)

	if ok, _ := passwordRule.MatchString(password); !ok {
		return false
	}

	return true
}

func (as *AdminService) verifyRole(ctx context.Context, roleId uint64) error {
	role, err := as.ruc.GetRoleById(ctx, roleId)

	if err != nil {
		return biz.AdminRoleNotFound
	}

	if role.Status == 0 {
		return biz.AdminRoleDisabled
	}

	return nil
}

func (as *AdminService) verifyToken(ctx context.Context, token string) bool {
	return as.tuc.VerifyToken(ctx, token)
}

func (as *AdminService) verifyEmail(ctx context.Context, addr string) bool {
	a, err := mail.ParseAddress(addr)

	if err != nil {
		return false
	}

	addr = a.Address

	if len(addr) > 254 {
		return false
	}

	parts := strings.SplitN(addr, "@", 2)

	if len(parts[0]) > 64 {
		return false
	}

	host := parts[1]
	hostLower := strings.ToLower(strings.TrimSuffix(host, "."))

	if len(host) > 253 {
		return false
	}

	for _, part := range strings.Split(hostLower, ".") {
		if l := len(part); l == 0 || l > 63 {
			return false
		}

		if part[0] == '-' {
			return false
		}

		if part[len(part)-1] == '-' {
			return false
		}

		for _, r := range part {
			if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
				return false
			}
		}
	}

	return true
}

func (as *AdminService) verifyMenu(ctx context.Context, ids string) ([]string, error) {
	inIds := strings.Split(ids, ",")
	list := make([]string, 0)

	for _, id := range inIds {
		idUint, _ := strconv.ParseUint(id, 10, 64)

		if _, err := as.muc.GetMenuById(ctx, idUint); err != nil {
			continue
		}

		isExist := true

		for _, lid := range list {
			if lid == id {
				isExist = false
				break
			}
		}

		if isExist {
			list = append(list, id)
		}
	}

	if len(list) == 0 {
		return nil, biz.AdminRoleMenuNotFound
	}

	return list, nil
}

func (as *AdminService) verifyPermissionCode(ctx context.Context, permissionCode string) bool {
	permissionCodes := make([]*domain.PermissionCodes, 0)
	permissionCodes = domain.NewPermissionCodes()

	isExist := false

	for _, pc := range permissionCodes {
		if pc.Code == permissionCode {
			isExist = true
			break
		}

		for _, cpc := range pc.ChildPermissionCodes {
			if cpc.Code == permissionCode {
				isExist = true
				break
			}
		}
	}

	return isExist
}

func (as *AdminService) listMenus(ctx context.Context, menus []*domain.Menu) []*ListMenusReply_Menus {
	if menus == nil {
		return nil
	}

	list := make([]*ListMenusReply_Menus, 0)

	for _, lmenu := range menus {
		var menuTypeName string

		if lmenu.MenuType == 1 {
			menuTypeName = "目录"
		} else if lmenu.MenuType == 2 {
			menuTypeName = "菜单"
		} else if lmenu.MenuType == 3 {
			menuTypeName = "按钮"
		}

		menu := &ListMenusReply_Menus{
			Id:             lmenu.Id,
			MenuName:       lmenu.MenuName,
			ParentId:       lmenu.ParentId,
			MenuType:       uint32(lmenu.MenuType),
			MenuTypeName:   menuTypeName,
			Sort:           uint32(lmenu.Sort),
			FileName:       lmenu.FileName,
			IconName:       lmenu.IconName,
			PermissionCode: lmenu.PermissionCode,
			Status:         uint32(lmenu.Status),
		}

		menu.ChildList = as.listMenus(ctx, lmenu.ChildList)

		list = append(list, menu)
	}

	return list
}

func (as *AdminService) treeMenus(ctx context.Context, menus []*domain.Menu, parentId uint64) []*ListUserMenusReply_Menus {
	list := make([]*ListUserMenusReply_Menus, 0)

	for _, lmenu := range menus {
		if lmenu.ParentId == parentId && lmenu.Status == 1 && lmenu.MenuType != 3 {
			menu := &ListUserMenusReply_Menus{
				Id:       lmenu.Id,
				MenuName: lmenu.MenuName,
				FileName: lmenu.FileName,
				IconName: lmenu.IconName,
			}

			menu.ChildList = as.treeMenus(ctx, menus, lmenu.Id)

			list = append(list, menu)
		}
	}

	return list
}

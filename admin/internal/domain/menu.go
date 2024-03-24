package domain

import (
	"context"
	"time"
)

type Menu struct {
	Id             uint64
	MenuName       string
	ParentId       uint64
	MenuType       uint8
	Sort           uint8
	FileName       string
	IconName       string
	PermissionCode string
	Status         uint8
	CreateTime     time.Time
	UpdateTime     time.Time
	Roles          []*Role
	ChildList      []*Menu
}

type ChildPermissionCodes struct {
	Name string
	Code string
}

type PermissionCodes struct {
	Name                 string
	Code                 string
	ChildPermissionCodes []*ChildPermissionCodes
}

func NewMenu(ctx context.Context, menuName string, parentId uint64, menuType, sort uint8, fileName, iconName, permissionCode string, status uint8) *Menu {
	return &Menu{
		MenuName:       menuName,
		ParentId:       parentId,
		MenuType:       menuType,
		Sort:           sort,
		FileName:       fileName,
		IconName:       iconName,
		PermissionCode: permissionCode,
		Status:         status,
	}
}

func NewPermissionCodes() []*PermissionCodes {
	permissionCodes := make([]*PermissionCodes, 0)

	childPermissionCodes := make([]*ChildPermissionCodes, 0)
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "菜单查询", Code: "admin:menu:list"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "菜单新增", Code: "admin:menu:create"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "菜单修改", Code: "admin:menu:update"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "菜单状态修改", Code: "admin:menu:updateStatus"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "菜单删除", Code: "admin:menu:delete"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "菜单管理", Code: "admin:menu:view", ChildPermissionCodes: childPermissionCodes})

	childPermissionCodes = make([]*ChildPermissionCodes, 0)
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "角色查询", Code: "admin:role:list"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "角色新增", Code: "admin:role:create"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "角色修改", Code: "admin:role:update"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "角色状态修改", Code: "admin:role:updateStatus"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "角色删除", Code: "admin:role:delete"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "角色管理", Code: "admin:role:view", ChildPermissionCodes: childPermissionCodes})

	childPermissionCodes = make([]*ChildPermissionCodes, 0)
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "账号查询", Code: "admin:user:list"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "账号新增", Code: "admin:user:create"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "账号修改", Code: "admin:user:update"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "账号状态修改", Code: "admin:user:updateStatus"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "账号删除", Code: "admin:user:delete"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "账号管理", Code: "admin:user:view", ChildPermissionCodes: childPermissionCodes})

	childPermissionCodes = make([]*ChildPermissionCodes, 0)
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "盟码菜单查询", Code: "admin:companyMenu:list"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "盟码菜单新增", Code: "admin:companyMenu:create"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "盟码菜单修改", Code: "admin:companyMenu:update"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "盟码菜单状态修改", Code: "admin:companyMenu:updateStatus"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "盟码菜单删除", Code: "admin:companyMenu:delete"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "盟码菜单管理", Code: "admin:companyMenu:view", ChildPermissionCodes: childPermissionCodes})

	childPermissionCodes = make([]*ChildPermissionCodes, 0)
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "巨量引擎账号查询", Code: "admin:oceanengineConfig:list"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "巨量引擎账号新增", Code: "admin:oceanengineConfig:create"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "巨量引擎账号修改", Code: "admin:oceanengineConfig:update"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "巨量引擎账号状态修改", Code: "admin:oceanengineConfig:updateStatus"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "巨量引擎账号删除", Code: "admin:oceanengineConfig:delete"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "巨量引擎账号管理", Code: "admin:oceanengineConfig:view", ChildPermissionCodes: childPermissionCodes})

	childPermissionCodes = make([]*ChildPermissionCodes, 0)
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "短信日志查询", Code: "admin:smslog:list"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "短信日志管理", Code: "admin:smslog:view", ChildPermissionCodes: childPermissionCodes})

	childPermissionCodes = make([]*ChildPermissionCodes, 0)
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "线索库查询", Code: "admin:clue:list"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "线索库新增", Code: "admin:clue:create"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "线索库修改", Code: "admin:clue:update"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "线索库操作记录修改", Code: "admin:clue:updateOperationLog"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "线索库删除", Code: "admin:clue:delete"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "线索库管理", Code: "admin:clue:view", ChildPermissionCodes: childPermissionCodes})

	childPermissionCodes = make([]*ChildPermissionCodes, 0)
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "企业库查询", Code: "admin:company:list"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "企业库新增", Code: "admin:company:create"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "企业库修改", Code: "admin:company:update"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "企业库状态修改", Code: "admin:company:updateStatus"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "企业库权限修改", Code: "admin:company:updateRole"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "企业库删除", Code: "admin:company:delete"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "企业库管理", Code: "admin:company:view", ChildPermissionCodes: childPermissionCodes})

	childPermissionCodes = make([]*ChildPermissionCodes, 0)
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "企业用户查询", Code: "admin:companyUser:list"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "企业用户新增", Code: "admin:companyUser:create"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "企业用户修改", Code: "admin:companyUser:update"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "企业用户状态修改", Code: "admin:companyUser:updateStatus"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "企业用户权限修改", Code: "admin:companyUser:updateRole"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "企业用户白名单修改", Code: "admin:companyUser:updateWhite"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "企业用户删除", Code: "admin:companyUser:delete"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "企业用户管理", ChildPermissionCodes: childPermissionCodes})

	childPermissionCodes = make([]*ChildPermissionCodes, 0)
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "更新日志查询", Code: "admin:updateLog:list"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "更新日志新增", Code: "admin:updateLog:create"})
	childPermissionCodes = append(childPermissionCodes, &ChildPermissionCodes{Name: "更新日志修改", Code: "admin:updateLog:update"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "更新日志管理", Code: "admin:updateLog:view", ChildPermissionCodes: childPermissionCodes})

	return permissionCodes
}

func (m *Menu) SetMenuName(ctx context.Context, menuName string) {
	m.MenuName = menuName
}

func (m *Menu) SetSort(ctx context.Context, sort uint8) {
	m.Sort = sort
}

func (m *Menu) SetFileName(ctx context.Context, fileName string) {
	m.FileName = fileName
}

func (m *Menu) SetIconName(ctx context.Context, iconName string) {
	m.IconName = iconName
}

func (m *Menu) SetPermissionCode(ctx context.Context, permissionCode string) {
	m.PermissionCode = permissionCode
}

func (m *Menu) SetStatus(ctx context.Context, status uint8) {
	m.Status = status
}

func (m *Menu) SetUpdateTime(ctx context.Context) {
	m.UpdateTime = time.Now()
}

func (m *Menu) SetCreateTime(ctx context.Context) {
	m.CreateTime = time.Now()
}

func (m *Menu) VerifyStatus(ctx context.Context) bool {
	if m.Status == 1 {
		return true
	}

	return false
}

package domain

import (
	"context"
	"time"
)

type Menu struct {
	Id             uint64
	MenuName       string
	ParentId       uint64
	Sort           uint8
	MenuType       string
	Filename       string
	Filepath       string
	PermissionCode string
	Status         uint8
	CreateTime     time.Time
	UpdateTime     time.Time
	ChildList      []*Menu
}

func NewMenu(ctx context.Context, menuName, menuType, filename, filepath, permissionCode string, parentId uint64, sort, status uint8) *Menu {
	return &Menu{
		MenuName:       menuName,
		ParentId:       parentId,
		Sort:           sort,
		MenuType:       menuType,
		Filename:       filename,
		Filepath:       filepath,
		PermissionCode: permissionCode,
		Status:         status,
	}
}

func (m *Menu) SetMenuName(ctx context.Context, menuName string) {
	m.MenuName = menuName
}

func (m *Menu) SetSort(ctx context.Context, sort uint8) {
	m.Sort = sort
}

func (m *Menu) SetPermissionCode(ctx context.Context, permissionCode string) {
	m.PermissionCode = permissionCode
}

func (m *Menu) SetStatus(ctx context.Context, status uint8) {
	m.Status = status
}

func (m *Menu) SetMenuType(ctx context.Context, menuType string) {
	m.MenuType = menuType
}

func (m *Menu) SetFilename(ctx context.Context, filename string) {
	m.Filename = filename
}

func (m *Menu) SetFilepath(ctx context.Context, filepath string) {
	m.Filepath = filepath
}

func (m *Menu) SetUpdateTime(ctx context.Context) {
	m.UpdateTime = time.Now()
}

func (m *Menu) SetCreateTime(ctx context.Context) {
	m.CreateTime = time.Now()
}

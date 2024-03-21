package data

import (
	"company/internal/biz"
	"company/internal/domain"
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// 菜单表
type Menu struct {
	Id             uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	MenuName       string    `gorm:"column:menu_name;type:varchar(20);not null;comment:菜单名称"`
	ParentId       uint64    `gorm:"column:parent_id;type:bigint(20) UNSIGNED;not null;default:0;comment:父级ID"`
	Sort           uint8     `gorm:"column:sort;type:tinyint(3) UNSIGNED;not null;default:0;comment:菜单顺序 数值越小排位越靠前"`
	MenuType       string    `gorm:"column:type;type:enum('menu','button');not null;default:'menu';comment:类型，menu：菜单 button：按钮"`
	Filename       string    `gorm:"column:filename;type:varchar(50);not null;comment:文件名"`
	Filepath       string    `gorm:"column:filepath;type:varchar(50);not null;comment:路径"`
	PermissionCode string    `gorm:"column:permission_code;type:varchar(50);not null;comment:权限标识码"`
	Status         uint8     `gorm:"column:status;type:tinyint(3) UNSIGNED;not null;default:0;comment:状态：1：启用，0：禁用"`
	CreateTime     time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime     time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (Menu) TableName() string {
	return "company_menu"
}

type menuRepo struct {
	data *Data
	log  *log.Helper
}

func (m *Menu) ToDomain() *domain.Menu {
	return &domain.Menu{
		Id:             m.Id,
		MenuName:       m.MenuName,
		ParentId:       m.ParentId,
		Sort:           m.Sort,
		MenuType:       m.MenuType,
		Filename:       m.Filename,
		Filepath:       m.Filepath,
		PermissionCode: m.PermissionCode,
		Status:         m.Status,
		CreateTime:     m.CreateTime,
		UpdateTime:     m.UpdateTime,
	}
}

func NewMenuRepo(data *Data, logger log.Logger) biz.MenuRepo {
	return &menuRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (mr *menuRepo) GetById(ctx context.Context, id uint64) (*domain.Menu, error) {
	menu := &Menu{}

	if result := mr.data.db.WithContext(ctx).First(menu, id); result.Error != nil {
		return nil, result.Error
	}

	return menu.ToDomain(), nil
}

func (mr *menuRepo) ListByIds(ctx context.Context, ids []string) ([]*domain.Menu, error) {
	var menus []Menu
	list := make([]*domain.Menu, 0)

	if result := mr.data.db.WithContext(ctx).Order("sort ASC,id DESC").Find(&menus, ids); result.Error != nil {
		return nil, result.Error
	}

	for _, menu := range menus {
		list = append(list, menu.ToDomain())
	}

	return list, nil
}

func (mr *menuRepo) ListByParentId(ctx context.Context, id uint64) ([]*domain.Menu, error) {
	var menus []Menu
	list := make([]*domain.Menu, 0)

	if result := mr.data.db.WithContext(ctx).
		Where("parent_id = ?", id).
		Order("sort ASC").
		Find(&menus); result.Error != nil {
		return nil, result.Error
	}

	for _, menu := range menus {
		list = append(list, menu.ToDomain())
	}

	return list, nil
}

func (mr *menuRepo) Save(ctx context.Context, in *domain.Menu) (*domain.Menu, error) {
	menu := &Menu{
		MenuName:       in.MenuName,
		ParentId:       in.ParentId,
		Sort:           in.Sort,
		MenuType:       in.MenuType,
		Filename:       in.Filename,
		Filepath:       in.Filepath,
		PermissionCode: in.PermissionCode,
		Status:         in.Status,
		CreateTime:     in.CreateTime,
		UpdateTime:     in.UpdateTime,
	}

	if result := mr.data.db.WithContext(ctx).Create(menu); result.Error != nil {
		return nil, result.Error
	}

	return menu.ToDomain(), nil
}

func (mr *menuRepo) Update(ctx context.Context, in *domain.Menu) (*domain.Menu, error) {
	menu := &Menu{
		Id:             in.Id,
		MenuName:       in.MenuName,
		ParentId:       in.ParentId,
		Sort:           in.Sort,
		MenuType:       in.MenuType,
		Filename:       in.Filename,
		Filepath:       in.Filepath,
		PermissionCode: in.PermissionCode,
		Status:         in.Status,
		CreateTime:     in.CreateTime,
		UpdateTime:     in.UpdateTime,
	}

	if result := mr.data.DB(ctx).Save(menu); result.Error != nil {
		return nil, result.Error
	}

	return menu.ToDomain(), nil
}

func (mr *menuRepo) Delete(ctx context.Context, in *domain.Menu) error {
	menu := &Menu{
		Id:             in.Id,
		MenuName:       in.MenuName,
		ParentId:       in.ParentId,
		Sort:           in.Sort,
		MenuType:       in.MenuType,
		Filename:       in.Filename,
		Filepath:       in.Filepath,
		PermissionCode: in.PermissionCode,
		Status:         in.Status,
		CreateTime:     in.CreateTime,
		UpdateTime:     in.UpdateTime,
	}

	if result := mr.data.DB(ctx).Delete(menu); result.Error != nil {
		return result.Error
	}

	return nil
}

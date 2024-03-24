package data

import (
	"admin/internal/biz"
	"admin/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// 角色表
type Role struct {
	Id          uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	RoleName    string    `gorm:"column:role_name;type:varchar(20);not null;comment:权限名称"`
	RoleExplain string    `gorm:"column:role_explain;type:varchar(255);not null;comment:权限描述"`
	MenuIds     string    `gorm:"column:menu_ids;type:varchar(255);not null;comment:菜单IDs"`
	Status      uint8     `gorm:"column:status;type:tinyint(3) UNSIGNED;not null;default:0;comment:状态：1：启用，0：禁用"`
	CreateTime  time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime  time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
	User        []User
}

func (Role) TableName() string {
	return "admin_role"
}

type roleRepo struct {
	data *Data
	log  *log.Helper
}

func (r *Role) ToDomain() *domain.Role {
	role := &domain.Role{
		Id:          r.Id,
		RoleName:    r.RoleName,
		RoleExplain: r.RoleExplain,
		MenuIds:     r.MenuIds,
		Status:      r.Status,
		CreateTime:  r.CreateTime,
		UpdateTime:  r.UpdateTime,
	}

	role.Users = make([]*domain.User, 0)

	for _, user := range r.User {
		role.Users = append(role.Users, &domain.User{
			Id:            user.Id,
			Username:      user.Username,
			Nickname:      user.Nickname,
			Password:      user.Password,
			Email:         user.Email,
			RoleId:        user.RoleId,
			Status:        user.Status,
			LastLoginIp:   user.LastLoginIp,
			LastLoginTime: user.LastLoginTime,
			CreateTime:    user.CreateTime,
			UpdateTime:    user.UpdateTime,
		})
	}

	return role
}

func NewRoleRepo(data *Data, logger log.Logger) biz.RoleRepo {
	return &roleRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (rr *roleRepo) GetById(ctx context.Context, id uint64) (*domain.Role, error) {
	role := &Role{}

	if result := rr.data.db.WithContext(ctx).Preload("User").First(role, id); result.Error != nil {
		return nil, result.Error
	}

	return role.ToDomain(), nil
}

func (rr *roleRepo) List(ctx context.Context) ([]*domain.Role, error) {
	var roles []Role
	list := make([]*domain.Role, 0)

	if result := rr.data.db.WithContext(ctx).
		Order("update_time DESC").
		Find(&roles); result.Error != nil {
		return nil, result.Error
	}

	for _, role := range roles {
		list = append(list, role.ToDomain())
	}

	return list, nil
}

func (rr *roleRepo) ListByMenuId(ctx context.Context, id uint64) ([]*domain.Role, error) {
	var roles []Role
	list := make([]*domain.Role, 0)

	if result := rr.data.db.WithContext(ctx).
		Where("find_in_set(?, `menu_ids`)", id).
		Order("update_time DESC").
		Find(&roles); result.Error != nil {
		return nil, result.Error
	}

	for _, role := range roles {
		list = append(list, role.ToDomain())
	}

	return list, nil
}

func (rr *roleRepo) Save(ctx context.Context, in *domain.Role) (*domain.Role, error) {
	role := &Role{
		RoleName:    in.RoleName,
		RoleExplain: in.RoleExplain,
		MenuIds:     in.MenuIds,
		Status:      in.Status,
		CreateTime:  in.CreateTime,
		UpdateTime:  in.UpdateTime,
	}

	if result := rr.data.db.WithContext(ctx).Create(role); result.Error != nil {
		return nil, result.Error
	}

	return role.ToDomain(), nil
}

func (rr *roleRepo) Update(ctx context.Context, in *domain.Role) (*domain.Role, error) {
	role := &Role{
		Id:          in.Id,
		RoleName:    in.RoleName,
		RoleExplain: in.RoleExplain,
		MenuIds:     in.MenuIds,
		Status:      in.Status,
		CreateTime:  in.CreateTime,
		UpdateTime:  in.UpdateTime,
	}

	if result := rr.data.db.WithContext(ctx).Save(role); result.Error != nil {
		return nil, result.Error
	}

	return role.ToDomain(), nil
}

func (rr *roleRepo) Delete(ctx context.Context, in *domain.Role) error {
	role := &Role{
		Id:          in.Id,
		RoleName:    in.RoleName,
		RoleExplain: in.RoleExplain,
		MenuIds:     in.MenuIds,
		Status:      in.Status,
		CreateTime:  in.CreateTime,
		UpdateTime:  in.UpdateTime,
	}

	if result := rr.data.db.WithContext(ctx).Delete(role); result.Error != nil {
		return result.Error
	}

	return nil
}

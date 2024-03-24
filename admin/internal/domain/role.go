package domain

import (
	"context"
	"time"
)

type Role struct {
	Id          uint64
	RoleName    string
	RoleExplain string
	MenuIds     string
	Status      uint8
	CreateTime  time.Time
	UpdateTime  time.Time
	Menus       []*Menu
	Users       []*User
}

func NewRole(ctx context.Context, roleName, roleExplain, ids string, status uint8) *Role {
	return &Role{
		RoleName:    roleName,
		RoleExplain: roleExplain,
		MenuIds:     ids,
		Status:      status,
	}
}

func (r *Role) SetRoleName(ctx context.Context, roleName string) {
	r.RoleName = roleName
}

func (r *Role) SetRoleExplain(ctx context.Context, roleExplain string) {
	r.RoleExplain = roleExplain
}

func (r *Role) SetMenuIds(ctx context.Context, ids string) {
	r.MenuIds = ids
}

func (r *Role) SetStatus(ctx context.Context, status uint8) {
	r.Status = status
}

func (r *Role) SetUpdateTime(ctx context.Context) {
	r.UpdateTime = time.Now()
}

func (r *Role) SetCreateTime(ctx context.Context) {
	r.CreateTime = time.Now()
}

func (r *Role) VerifyStatus(ctx context.Context) bool {
	if r.Status == 1 {
		return true
	}

	return false
}

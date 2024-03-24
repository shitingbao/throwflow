package domain

import (
	"admin/internal/pkg/tool"
	"context"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Login struct {
	Username      string
	Nickname      string
	Email         string
	LastLoginTime time.Time
	LastLoginIp   string
	Token         string
	RoleName      string
	Menus         []*Menu
}

type User struct {
	Id            uint64
	Username      string
	Nickname      string
	Password      string
	Email         string
	RoleId        uint64
	Status        uint8
	LastLoginIp   string
	LastLoginTime time.Time
	CreateTime    time.Time
	UpdateTime    time.Time
	Role          *Role
	Menus         []*Menu
}

type UserList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*User
}

type AllUserList struct {
	List []*User
}

func NewUser(ctx context.Context, username, nickname, email, password string, roleId uint64, status uint8) *User {
	return &User{
		Username: username,
		Nickname: nickname,
		Password: password,
		Email:    email,
		RoleId:   roleId,
		Status:   status,
	}
}

func (u *User) SetUsername(ctx context.Context, username string) {
	u.Username = username
}

func (u *User) SetNickname(ctx context.Context, nickname string) {
	u.Nickname = nickname
}

func (u *User) SetPassword(ctx context.Context, password string) {
	passwd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	u.Password = string(passwd)
}

func (u *User) SetEmail(ctx context.Context, email string) {
	u.Email = email
}

func (u *User) SetRoleId(ctx context.Context, roleId uint64) {
	u.RoleId = roleId
}

func (u *User) SetStatus(ctx context.Context, status uint8) {
	u.Status = status
}

func (u *User) SetUpdateTime(ctx context.Context) {
	u.UpdateTime = time.Now()
}

func (u *User) SetCreateTime(ctx context.Context) {
	u.CreateTime = time.Now()
}

func (u *User) SetLastLoginIp(ctx context.Context) {
	u.LastLoginIp = tool.GetClientIp(ctx)
}

func (u *User) SetLastLoginTime(ctx context.Context) {
	u.LastLoginTime = time.Now()
}

func (u *User) VerifyPassword(ctx context.Context, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return false
	}

	return true
}

func (u *User) VerifyStatus(ctx context.Context) bool {
	if u.Status == 1 {
		return true
	}

	return false
}

func (u *User) VerifyPermission(ctx context.Context, permissionCode string) bool {
	if ok := u.VerifyStatus(ctx); !ok {
		return false
	}

	if permissionCode == "all" {
		return true
	}

	if ok := u.Role.VerifyStatus(ctx); !ok {
		return false
	}

	for _, menu := range u.Menus {
		if menu.PermissionCode == permissionCode && menu.VerifyStatus(ctx) {
			return true
		}
	}

	return false
}

package data

import (
	"admin/internal/domain"
	"strconv"
	"time"

	"admin/internal/biz"

	"context"
	"github.com/go-kratos/kratos/v2/log"
)

// 管理员表
type User struct {
	Id            uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	Username      string    `gorm:"column:username;uniqueIndex;type:varchar(12);not null;comment:登录账号"`
	Nickname      string    `gorm:"column:nickname;type:varchar(50);not null;comment:昵称"`
	Password      string    `gorm:"column:password;type:varchar(60);not null;comment:密码"`
	Email         string    `gorm:"column:email;type:varchar(80);not null;comment:邮箱"`
	RoleId        uint64    `gorm:"column:role_id;type:bigint(20) UNSIGNED;not null;comment:角色ID"`
	Status        uint8     `gorm:"column:status;type:tinyint(3) UNSIGNED;not null;default:0;comment:状态：1：启用，0：禁用"`
	LastLoginIp   string    `gorm:"column:last_login_ip;type:varchar(50);not null;comment:最后一次登录IP"`
	LastLoginTime time.Time `gorm:"column:last_login_time;type:datetime;not null;comment:最后一次登录时间"`
	CreateTime    time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime    time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
	Role          Role      `gorm:"foreignKey:RoleId"`
}

func (User) TableName() string {
	return "admin_user"
}

type userRepo struct {
	data *Data
	log  *log.Helper
}

func (u *User) ToDomain() *domain.User {
	return &domain.User{
		Id:            u.Id,
		Username:      u.Username,
		Nickname:      u.Nickname,
		Password:      u.Password,
		Email:         u.Email,
		RoleId:        u.RoleId,
		Status:        u.Status,
		LastLoginIp:   u.LastLoginIp,
		LastLoginTime: u.LastLoginTime,
		CreateTime:    u.CreateTime,
		UpdateTime:    u.UpdateTime,
		Role: &domain.Role{
			Id:          u.Role.Id,
			RoleName:    u.Role.RoleName,
			RoleExplain: u.Role.RoleExplain,
			MenuIds:     u.Role.MenuIds,
			Status:      u.Role.Status,
			CreateTime:  u.Role.CreateTime,
			UpdateTime:  u.Role.UpdateTime,
		},
	}
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ur *userRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	user := &User{}

	if result := ur.data.db.WithContext(ctx).Preload("Role").Where("username = ?", username).First(user); result.Error != nil {
		return nil, result.Error
	}

	return user.ToDomain(), nil
}

func (ur *userRepo) GetById(ctx context.Context, id uint64) (*domain.User, error) {
	user := &User{}

	if result := ur.data.db.WithContext(ctx).Preload("Role").First(user, id); result.Error != nil {
		return nil, result.Error
	}

	return user.ToDomain(), nil
}

func (ur *userRepo) List(ctx context.Context, pageNum int) ([]*domain.User, error) {
	var users []User
	list := make([]*domain.User, 0)
	pageSize := int(ur.data.conf.Database.PageSize)

	if result := ur.data.db.WithContext(ctx).
		Order("update_time DESC,id DESC").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&users); result.Error != nil {
		return nil, result.Error
	}

	for _, user := range users {
		list = append(list, user.ToDomain())
	}

	return list, nil
}

func (ur *userRepo) Count(ctx context.Context) (int64, error) {
	var count int64

	if result := ur.data.db.Model(&User{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (ur *userRepo) Save(ctx context.Context, in *domain.User) (*domain.User, error) {
	user := &User{
		Username:      in.Username,
		Nickname:      in.Nickname,
		Password:      in.Password,
		Email:         in.Email,
		RoleId:        in.RoleId,
		Status:        in.Status,
		LastLoginIp:   in.LastLoginIp,
		LastLoginTime: in.LastLoginTime,
		CreateTime:    in.CreateTime,
		UpdateTime:    in.UpdateTime,
	}

	if result := ur.data.db.WithContext(ctx).Create(user); result.Error != nil {
		return nil, result.Error
	}

	return user.ToDomain(), nil
}

func (ur *userRepo) Update(ctx context.Context, in *domain.User) (*domain.User, error) {
	user := &User{
		Id:            in.Id,
		Username:      in.Username,
		Nickname:      in.Nickname,
		Password:      in.Password,
		Email:         in.Email,
		RoleId:        in.RoleId,
		Status:        in.Status,
		LastLoginIp:   in.LastLoginIp,
		LastLoginTime: in.LastLoginTime,
		CreateTime:    in.CreateTime,
		UpdateTime:    in.UpdateTime,
	}

	if result := ur.data.db.WithContext(ctx).Save(user); result.Error != nil {
		return nil, result.Error
	}

	return user.ToDomain(), nil
}

func (ur *userRepo) Delete(ctx context.Context, in *domain.User) error {
	user := &User{
		Id:            in.Id,
		Username:      in.Username,
		Nickname:      in.Nickname,
		Password:      in.Password,
		Email:         in.Email,
		RoleId:        in.RoleId,
		Status:        in.Status,
		LastLoginIp:   in.LastLoginIp,
		LastLoginTime: in.LastLoginTime,
		CreateTime:    in.CreateTime,
		UpdateTime:    in.UpdateTime,
	}

	if result := ur.data.db.WithContext(ctx).Delete(user); result.Error != nil {
		return result.Error
	}

	return nil
}

func (ur *userRepo) SaveCacheHash(ctx context.Context, key string, val map[string]string, timeout time.Duration) error {
	_, err := ur.data.rdb.HMSet(ctx, key, val).Result()

	if err != nil {
		return err
	}

	_, err = ur.data.rdb.Expire(ctx, key, timeout).Result()

	if err != nil {
		return err
	}

	return nil
}

func (ur *userRepo) GetCacheHash(ctx context.Context, key string, field string) (uint64, error) {
	val, err := ur.data.rdb.HGet(ctx, key, field).Result()

	if err != nil {
		return 0, err
	}

	uintId, err := strconv.Atoi(val)

	if err != nil {
		return 0, err
	}

	return uint64(uintId), nil
}

func (ur *userRepo) DeleteCache(ctx context.Context, key string) error {
	if _, err := ur.data.rdb.Del(ctx, key).Result(); err != nil {
		return err
	}

	return nil
}

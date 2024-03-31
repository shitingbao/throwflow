package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"weixin/internal/biz"
	"weixin/internal/domain"
)

// 微信小程序用户收获地址表
type UserAddress struct {
	Id               uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	UserId           uint64    `gorm:"column:user_id;type:bigint(20) UNSIGNED;not null;comment:微信小程序用户ID"`
	Name             string    `gorm:"column:name;type:varchar(100);not null;comment:收货人"`
	Phone            string    `gorm:"column:phone;type:varchar(20);not null;comment:手机号码"`
	ProvinceAreaCode uint64    `gorm:"column:province_area_code;type:bigint(20) UNSIGNED;not null;default:0;comment:省区划代码"`
	CityAreaCode     uint64    `gorm:"column:city_area_code;type:bigint(20) UNSIGNED;not null;default:0;comment:市区划代码"`
	AreaAreaCode     uint64    `gorm:"column:area_area_code;type:bigint(20) UNSIGNED;not null;default:0;comment:区区划代码"`
	AddressInfo      string    `gorm:"column:address_info;type:text;not null;comment:详细地址"`
	IsDefault        uint8     `gorm:"column:is_default;type:tinyint(3) UNSIGNED;not null;default:0;comment:是否默认地址，1：是，0：否"`
	CreateTime       time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime       time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (UserAddress) TableName() string {
	return "weixin_user_address"
}

type userAddressRepo struct {
	data *Data
	log  *log.Helper
}

func (ua *UserAddress) ToDomain(ctx context.Context) *domain.UserAddress {
	userAddress := &domain.UserAddress{
		Id:               ua.Id,
		UserId:           ua.UserId,
		Name:             ua.Name,
		Phone:            ua.Phone,
		ProvinceAreaCode: ua.ProvinceAreaCode,
		CityAreaCode:     ua.CityAreaCode,
		AreaAreaCode:     ua.AreaAreaCode,
		AddressInfo:      ua.AddressInfo,
		IsDefault:        ua.IsDefault,
		CreateTime:       ua.CreateTime,
		UpdateTime:       ua.UpdateTime,
	}

	return userAddress
}

func NewUserAddressRepo(data *Data, logger log.Logger) biz.UserAddressRepo {
	return &userAddressRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (uar *userAddressRepo) GetById(ctx context.Context, userId, userAddressId uint64) (*domain.UserAddress, error) {
	userAddress := &UserAddress{}

	if result := uar.data.db.WithContext(ctx).Where("id = ?", userAddressId).Where("user_id = ?", userId).First(userAddress); result.Error != nil {
		return nil, result.Error
	}

	return userAddress.ToDomain(ctx), nil
}

func (uar *userAddressRepo) List(ctx context.Context, pageNum, pageSize int, userId uint64) ([]*domain.UserAddress, error) {
	var userAddresses []UserAddress
	list := make([]*domain.UserAddress, 0)

	if result := uar.data.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Order("create_time DESC,id DESC").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&userAddresses); result.Error != nil {
		return nil, result.Error
	}

	for _, userAddress := range userAddresses {
		list = append(list, userAddress.ToDomain(ctx))
	}

	return list, nil
}

func (uar *userAddressRepo) Count(ctx context.Context, userId uint64) (int64, error) {
	var count int64

	if result := uar.data.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Model(&UserAddress{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (uar *userAddressRepo) Save(ctx context.Context, in *domain.UserAddress) (*domain.UserAddress, error) {
	userAddress := &UserAddress{
		UserId:           in.UserId,
		Name:             in.Name,
		Phone:            in.Phone,
		ProvinceAreaCode: in.ProvinceAreaCode,
		CityAreaCode:     in.CityAreaCode,
		AreaAreaCode:     in.AreaAreaCode,
		AddressInfo:      in.AddressInfo,
		IsDefault:        in.IsDefault,
		CreateTime:       in.CreateTime,
		UpdateTime:       in.UpdateTime,
	}

	if result := uar.data.DB(ctx).Create(userAddress); result.Error != nil {
		return nil, result.Error
	}

	return userAddress.ToDomain(ctx), nil
}

func (uar *userAddressRepo) Update(ctx context.Context, in *domain.UserAddress) (*domain.UserAddress, error) {
	userAddress := &UserAddress{
		Id:               in.Id,
		UserId:           in.UserId,
		Name:             in.Name,
		Phone:            in.Phone,
		ProvinceAreaCode: in.ProvinceAreaCode,
		CityAreaCode:     in.CityAreaCode,
		AreaAreaCode:     in.AreaAreaCode,
		AddressInfo:      in.AddressInfo,
		IsDefault:        in.IsDefault,
		CreateTime:       in.CreateTime,
		UpdateTime:       in.UpdateTime,
	}

	if result := uar.data.DB(ctx).Save(userAddress); result.Error != nil {
		return nil, result.Error
	}

	return userAddress.ToDomain(ctx), nil
}

func (uar *userAddressRepo) UpdateIsDefaults(ctx context.Context, userId uint64) error {
	if result := uar.data.DB(ctx).Model(UserAddress{}).
		Where("user_id", userId).
		Updates(map[string]interface{}{"is_default": 0}); result.Error != nil {
		return result.Error
	}

	return nil
}

func (uar *userAddressRepo) Delete(ctx context.Context, in *domain.UserAddress) error {
	userAddress := &UserAddress{
		Id:               in.Id,
		UserId:           in.UserId,
		Name:             in.Name,
		Phone:            in.Phone,
		ProvinceAreaCode: in.ProvinceAreaCode,
		CityAreaCode:     in.CityAreaCode,
		AreaAreaCode:     in.AreaAreaCode,
		AddressInfo:      in.AddressInfo,
		IsDefault:        in.IsDefault,
		CreateTime:       in.CreateTime,
		UpdateTime:       in.UpdateTime,
	}

	if result := uar.data.db.WithContext(ctx).Delete(userAddress); result.Error != nil {
		return result.Error
	}

	return nil
}

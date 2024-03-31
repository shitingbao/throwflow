package domain

import (
	"context"
	"time"
)

type UserAddress struct {
	Id               uint64
	UserId           uint64
	Name             string
	Phone            string
	ProvinceAreaCode uint64
	ProvinceAreaName string
	CityAreaCode     uint64
	CityAreaName     string
	AreaAreaCode     uint64
	AreaAreaName     string
	AddressInfo      string
	IsDefault        uint8
	CreateTime       time.Time
	UpdateTime       time.Time
}

type UserAddressList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*UserAddress
}

func NewUserAddress(ctx context.Context, userId, provinceAreaCode, cityAreaCode, areaAreaCode uint64, isDefault uint8, name, phone, addressInfo string) *UserAddress {
	return &UserAddress{
		UserId:           userId,
		Name:             name,
		Phone:            phone,
		ProvinceAreaCode: provinceAreaCode,
		CityAreaCode:     cityAreaCode,
		AreaAreaCode:     areaAreaCode,
		AddressInfo:      addressInfo,
		IsDefault:        isDefault,
	}
}

func (ua *UserAddress) SetUserId(ctx context.Context, userId uint64) {
	ua.UserId = userId
}

func (ua *UserAddress) SetName(ctx context.Context, name string) {
	ua.Name = name
}

func (ua *UserAddress) SetPhone(ctx context.Context, phone string) {
	ua.Phone = phone
}

func (ua *UserAddress) SetProvinceAreaCode(ctx context.Context, provinceAreaCode uint64) {
	ua.ProvinceAreaCode = provinceAreaCode
}

func (ua *UserAddress) SetProvinceAreaName(ctx context.Context, provinceAreaName string) {
	ua.ProvinceAreaName = provinceAreaName
}

func (ua *UserAddress) SetCityAreaCode(ctx context.Context, cityAreaCode uint64) {
	ua.CityAreaCode = cityAreaCode
}

func (ua *UserAddress) SetCityAreaName(ctx context.Context, cityAreaName string) {
	ua.CityAreaName = cityAreaName
}

func (ua *UserAddress) SetAreaAreaCode(ctx context.Context, areaAreaCode uint64) {
	ua.AreaAreaCode = areaAreaCode
}

func (ua *UserAddress) SetAreaAreaName(ctx context.Context, areaAreaName string) {
	ua.AreaAreaName = areaAreaName
}

func (ua *UserAddress) SetAddressInfo(ctx context.Context, addressInfo string) {
	ua.AddressInfo = addressInfo
}

func (ua *UserAddress) SetIsDefault(ctx context.Context, isDefault uint8) {
	ua.IsDefault = isDefault
}

func (ua *UserAddress) SetUpdateTime(ctx context.Context) {
	ua.UpdateTime = time.Now()
}

func (ua *UserAddress) SetCreateTime(ctx context.Context) {
	ua.CreateTime = time.Now()
}

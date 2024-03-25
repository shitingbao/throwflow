package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/weixin/v1"
	"interface/internal/biz"
)

type userAddressRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserAddressRepo(data *Data, logger log.Logger) biz.UserAddressRepo {
	return &userAddressRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (uar *userAddressRepo) List(ctx context.Context, userId, pageNum, pageSize uint64) (*v1.ListUserAddressesReply, error) {
	list, err := uar.data.weixinuc.ListUserAddresses(ctx, &v1.ListUserAddressesRequest{
		PageNum:  pageNum,
		PageSize: pageSize,
		UserId:   userId,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (uar *userAddressRepo) Save(ctx context.Context, userId, provinceAreaCode, cityAreaCode, areaAreaCode uint64, isDefault uint32, name, phone, addressInfo string) (*v1.CreateUserAddressesReply, error) {
	userAddress, err := uar.data.weixinuc.CreateUserAddresses(ctx, &v1.CreateUserAddressesRequest{
		UserId:           userId,
		Name:             name,
		Phone:            phone,
		ProvinceAreaCode: provinceAreaCode,
		CityAreaCode:     cityAreaCode,
		AreaAreaCode:     areaAreaCode,
		AddressInfo:      addressInfo,
		IsDefault:        isDefault,
	})

	if err != nil {
		return nil, err
	}

	return userAddress, err
}

func (uar *userAddressRepo) Update(ctx context.Context, userId, userAddressId, provinceAreaCode, cityAreaCode, areaAreaCode uint64, isDefault uint32, name, phone, addressInfo string) (*v1.UpdateUserAddressesReply, error) {
	userAddress, err := uar.data.weixinuc.UpdateUserAddresses(ctx, &v1.UpdateUserAddressesRequest{
		UserId:           userId,
		UserAddressId:    userAddressId,
		Name:             name,
		Phone:            phone,
		ProvinceAreaCode: provinceAreaCode,
		CityAreaCode:     cityAreaCode,
		AreaAreaCode:     areaAreaCode,
		AddressInfo:      addressInfo,
		IsDefault:        isDefault,
	})

	if err != nil {
		return nil, err
	}

	return userAddress, err
}

func (uar *userAddressRepo) UpdateDefault(ctx context.Context, userId, userAddressId uint64, isDefault uint32) (*v1.UpdateDefaultUserAddressesReply, error) {
	userAddress, err := uar.data.weixinuc.UpdateDefaultUserAddresses(ctx, &v1.UpdateDefaultUserAddressesRequest{
		UserId:        userId,
		UserAddressId: userAddressId,
		IsDefault:     isDefault,
	})

	if err != nil {
		return nil, err
	}

	return userAddress, err
}

func (uar *userAddressRepo) Delete(ctx context.Context, userId, userAddressId uint64) (*v1.DeleteUserAddressesReply, error) {
	userAddress, err := uar.data.weixinuc.DeleteUserAddresses(ctx, &v1.DeleteUserAddressesRequest{
		UserId:        userId,
		UserAddressId: userAddressId,
	})

	if err != nil {
		return nil, err
	}

	return userAddress, err
}

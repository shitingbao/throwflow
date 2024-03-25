package service

import (
	"context"
	v1 "interface/api/interface/v1"
	"interface/internal/biz"
)

func (is *InterfaceService) ListUserAddresses(ctx context.Context, in *v1.ListUserAddressesRequest) (*v1.ListUserAddressesReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	userAddresses, err := is.uauc.ListUserAddresses(ctx, in.PageNum, in.PageSize, userInfo.Data.UserId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListUserAddressesReply_UserAddress, 0)

	for _, userAddress := range userAddresses.Data.List {
		list = append(list, &v1.ListUserAddressesReply_UserAddress{
			UserAddressId:    userAddress.UserAddressId,
			Name:             userAddress.Name,
			Phone:            userAddress.Phone,
			ProvinceAreaCode: userAddress.ProvinceAreaCode,
			ProvinceAreaName: userAddress.ProvinceAreaName,
			CityAreaCode:     userAddress.CityAreaCode,
			CityAreaName:     userAddress.CityAreaName,
			AreaAreaCode:     userAddress.AreaAreaCode,
			AreaAreaName:     userAddress.AreaAreaName,
			AddressInfo:      userAddress.AddressInfo,
			IsDefault:        userAddress.IsDefault,
		})
	}

	return &v1.ListUserAddressesReply{
		Code: 200,
		Data: &v1.ListUserAddressesReply_Data{
			PageNum:   userAddresses.Data.PageNum,
			PageSize:  userAddresses.Data.PageSize,
			Total:     userAddresses.Data.Total,
			TotalPage: userAddresses.Data.TotalPage,
			List:      list,
		},
	}, nil
}

func (is *InterfaceService) CreateUserAddresses(ctx context.Context, in *v1.CreateUserAddressesRequest) (*v1.CreateUserAddressesReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	userAddress, err := is.uauc.CreateUserAddresses(ctx, userInfo.Data.UserId, in.ProvinceAreaCode, in.CityAreaCode, in.AreaAreaCode, in.IsDefault, in.Name, in.Phone, in.AddressInfo)

	if err != nil {
		return nil, err
	}

	return &v1.CreateUserAddressesReply{
		Code: 200,
		Data: &v1.CreateUserAddressesReply_Data{
			UserAddressId:    userAddress.Data.UserAddressId,
			Name:             userAddress.Data.Name,
			Phone:            userAddress.Data.Phone,
			ProvinceAreaCode: userAddress.Data.ProvinceAreaCode,
			ProvinceAreaName: userAddress.Data.ProvinceAreaName,
			CityAreaCode:     userAddress.Data.CityAreaCode,
			CityAreaName:     userAddress.Data.CityAreaName,
			AreaAreaCode:     userAddress.Data.AreaAreaCode,
			AreaAreaName:     userAddress.Data.AreaAreaName,
			AddressInfo:      userAddress.Data.AddressInfo,
			IsDefault:        userAddress.Data.IsDefault,
		},
	}, nil
}

func (is *InterfaceService) UpdateUserAddresses(ctx context.Context, in *v1.UpdateUserAddressesRequest) (*v1.UpdateUserAddressesReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	userAddress, err := is.uauc.UpdateUserAddresses(ctx, userInfo.Data.UserId, in.UserAddressId, in.ProvinceAreaCode, in.CityAreaCode, in.AreaAreaCode, in.IsDefault, in.Name, in.Phone, in.AddressInfo)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateUserAddressesReply{
		Code: 200,
		Data: &v1.UpdateUserAddressesReply_Data{
			UserAddressId:    userAddress.Data.UserAddressId,
			Name:             userAddress.Data.Name,
			Phone:            userAddress.Data.Phone,
			ProvinceAreaCode: userAddress.Data.ProvinceAreaCode,
			ProvinceAreaName: userAddress.Data.ProvinceAreaName,
			CityAreaCode:     userAddress.Data.CityAreaCode,
			CityAreaName:     userAddress.Data.CityAreaName,
			AreaAreaCode:     userAddress.Data.AreaAreaCode,
			AreaAreaName:     userAddress.Data.AreaAreaName,
			AddressInfo:      userAddress.Data.AddressInfo,
			IsDefault:        userAddress.Data.IsDefault,
		},
	}, nil
}

func (is *InterfaceService) UpdateDefaultUserAddresses(ctx context.Context, in *v1.UpdateDefaultUserAddressesRequest) (*v1.UpdateDefaultUserAddressesReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	userAddress, err := is.uauc.UpdateDefaultUserAddresses(ctx, userInfo.Data.UserId, in.UserAddressId, in.IsDefault)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateDefaultUserAddressesReply{
		Code: 200,
		Data: &v1.UpdateDefaultUserAddressesReply_Data{
			UserAddressId:    userAddress.Data.UserAddressId,
			Name:             userAddress.Data.Name,
			Phone:            userAddress.Data.Phone,
			ProvinceAreaCode: userAddress.Data.ProvinceAreaCode,
			ProvinceAreaName: userAddress.Data.ProvinceAreaName,
			CityAreaCode:     userAddress.Data.CityAreaCode,
			CityAreaName:     userAddress.Data.CityAreaName,
			AreaAreaCode:     userAddress.Data.AreaAreaCode,
			AreaAreaName:     userAddress.Data.AreaAreaName,
			AddressInfo:      userAddress.Data.AddressInfo,
			IsDefault:        userAddress.Data.IsDefault,
		},
	}, nil
}

func (is *InterfaceService) DeleteUserAddresses(ctx context.Context, in *v1.DeleteUserAddressesRequest) (*v1.DeleteUserAddressesReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	if _, err := is.uauc.DeleteUserAddresses(ctx, userInfo.Data.UserId, in.UserAddressId); err != nil {
		return nil, err
	}

	return &v1.DeleteUserAddressesReply{
		Code: 200,
		Data: &v1.DeleteUserAddressesReply_Data{},
	}, nil
}

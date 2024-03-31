package service

import (
	"context"
	"math"
	v1 "weixin/api/weixin/v1"
)

func (ws *WeixinService) ListUserAddresses(ctx context.Context, in *v1.ListUserAddressesRequest) (*v1.ListUserAddressesReply, error) {
	userAddresses, err := ws.uauc.ListUserAddresses(ctx, in.PageNum, in.PageSize, in.UserId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListUserAddressesReply_UserAddress, 0)

	for _, userAddress := range userAddresses.List {
		list = append(list, &v1.ListUserAddressesReply_UserAddress{
			UserAddressId:    userAddress.Id,
			Name:             userAddress.Name,
			Phone:            userAddress.Phone,
			ProvinceAreaCode: userAddress.ProvinceAreaCode,
			ProvinceAreaName: userAddress.ProvinceAreaName,
			CityAreaCode:     userAddress.CityAreaCode,
			CityAreaName:     userAddress.CityAreaName,
			AreaAreaCode:     userAddress.AreaAreaCode,
			AreaAreaName:     userAddress.AreaAreaName,
			AddressInfo:      userAddress.AddressInfo,
			IsDefault:        uint32(userAddress.IsDefault),
		})
	}

	totalPage := uint64(math.Ceil(float64(userAddresses.Total) / float64(userAddresses.PageSize)))

	return &v1.ListUserAddressesReply{
		Code: 200,
		Data: &v1.ListUserAddressesReply_Data{
			PageNum:   userAddresses.PageNum,
			PageSize:  userAddresses.PageSize,
			Total:     userAddresses.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ws *WeixinService) CreateUserAddresses(ctx context.Context, in *v1.CreateUserAddressesRequest) (*v1.CreateUserAddressesReply, error) {
	userAddress, err := ws.uauc.CreateUserAddresses(ctx, in.UserId, in.ProvinceAreaCode, in.CityAreaCode, in.AreaAreaCode, uint8(in.IsDefault), in.Name, in.Phone, in.AddressInfo)

	if err != nil {
		return nil, err
	}

	return &v1.CreateUserAddressesReply{
		Code: 200,
		Data: &v1.CreateUserAddressesReply_Data{
			UserAddressId:    userAddress.Id,
			Name:             userAddress.Name,
			Phone:            userAddress.Phone,
			ProvinceAreaCode: userAddress.ProvinceAreaCode,
			ProvinceAreaName: userAddress.ProvinceAreaName,
			CityAreaCode:     userAddress.CityAreaCode,
			CityAreaName:     userAddress.CityAreaName,
			AreaAreaCode:     userAddress.AreaAreaCode,
			AreaAreaName:     userAddress.AreaAreaName,
			AddressInfo:      userAddress.AddressInfo,
			IsDefault:        uint32(userAddress.IsDefault),
		},
	}, nil
}

func (ws *WeixinService) UpdateUserAddresses(ctx context.Context, in *v1.UpdateUserAddressesRequest) (*v1.UpdateUserAddressesReply, error) {
	userAddress, err := ws.uauc.UpdateUserAddresses(ctx, in.UserId, in.UserAddressId, in.ProvinceAreaCode, in.CityAreaCode, in.AreaAreaCode, uint8(in.IsDefault), in.Name, in.Phone, in.AddressInfo)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateUserAddressesReply{
		Code: 200,
		Data: &v1.UpdateUserAddressesReply_Data{
			UserAddressId:    userAddress.Id,
			Name:             userAddress.Name,
			Phone:            userAddress.Phone,
			ProvinceAreaCode: userAddress.ProvinceAreaCode,
			ProvinceAreaName: userAddress.ProvinceAreaName,
			CityAreaCode:     userAddress.CityAreaCode,
			CityAreaName:     userAddress.CityAreaName,
			AreaAreaCode:     userAddress.AreaAreaCode,
			AreaAreaName:     userAddress.AreaAreaName,
			AddressInfo:      userAddress.AddressInfo,
			IsDefault:        uint32(userAddress.IsDefault),
		},
	}, nil
}

func (ws *WeixinService) UpdateDefaultUserAddresses(ctx context.Context, in *v1.UpdateDefaultUserAddressesRequest) (*v1.UpdateDefaultUserAddressesReply, error) {
	userAddress, err := ws.uauc.UpdateDefaultUserAddresses(ctx, in.UserId, in.UserAddressId, uint8(in.IsDefault))

	if err != nil {
		return nil, err
	}

	return &v1.UpdateDefaultUserAddressesReply{
		Code: 200,
		Data: &v1.UpdateDefaultUserAddressesReply_Data{
			UserAddressId:    userAddress.Id,
			Name:             userAddress.Name,
			Phone:            userAddress.Phone,
			ProvinceAreaCode: userAddress.ProvinceAreaCode,
			ProvinceAreaName: userAddress.ProvinceAreaName,
			CityAreaCode:     userAddress.CityAreaCode,
			CityAreaName:     userAddress.CityAreaName,
			AreaAreaCode:     userAddress.AreaAreaCode,
			AreaAreaName:     userAddress.AreaAreaName,
			AddressInfo:      userAddress.AddressInfo,
			IsDefault:        uint32(userAddress.IsDefault),
		},
	}, nil
}

func (ws *WeixinService) DeleteUserAddresses(ctx context.Context, in *v1.DeleteUserAddressesRequest) (*v1.DeleteUserAddressesReply, error) {
	err := ws.uauc.DeleteUserAddresses(ctx, in.UserId, in.UserAddressId)

	if err != nil {
		return nil, err
	}

	return &v1.DeleteUserAddressesReply{
		Code: 200,
		Data: &v1.DeleteUserAddressesReply_Data{},
	}, nil
}

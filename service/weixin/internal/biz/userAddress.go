package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"weixin/internal/conf"
	"weixin/internal/domain"
)

var (
	WeixinUserAddressNotFound    = errors.NotFound("WEIXIN_USER_ADDRESS_NOT_FOUND", "微信用户收货地址不存在")
	WeixinUserAddressListError   = errors.InternalServer("WEIXIN_USER_ADDRESS_LIST_ERROR", "微信用户收货地址列表获取失败")
	WeixinUserAddressCreateError = errors.InternalServer("WEIXIN_USER_ADDRESS_CREATE_ERROR", "微信用户收货地址创建失败")
	WeixinUserAddressUpdateError = errors.InternalServer("WEIXIN_USER_ADDRESS_UPDATE_ERROR", "微信用户收货地址更新失败")
	WeixinUserAddressDeleteError = errors.InternalServer("WEIXIN_USER_ADDRESS_DELETE_ERROR", "微信用户收货地址删除失败")
)

type UserAddressRepo interface {
	GetById(context.Context, uint64, uint64) (*domain.UserAddress, error)
	List(context.Context, int, int, uint64) ([]*domain.UserAddress, error)
	Count(context.Context, uint64) (int64, error)
	Save(context.Context, *domain.UserAddress) (*domain.UserAddress, error)
	Update(context.Context, *domain.UserAddress) (*domain.UserAddress, error)
	UpdateIsDefaults(context.Context, uint64) error
	Delete(context.Context, *domain.UserAddress) error
}

type UserAddressUsecase struct {
	repo  UserAddressRepo
	arepo AreaRepo
	tm    Transaction
	conf  *conf.Data
	log   *log.Helper
}

func NewUserAddressUsecase(repo UserAddressRepo, arepo AreaRepo, tm Transaction, conf *conf.Data, logger log.Logger) *UserAddressUsecase {
	return &UserAddressUsecase{repo: repo, arepo: arepo, tm: tm, conf: conf, log: log.NewHelper(logger)}
}

func (uauc *UserAddressUsecase) ListUserAddresses(ctx context.Context, pageNum, pageSize, userId uint64) (*domain.UserAddressList, error) {
	userAddresses, err := uauc.repo.List(ctx, int(pageNum), int(pageSize), userId)

	if err != nil {
		return nil, WeixinUserAddressListError
	}

	total, err := uauc.repo.Count(ctx, userId)

	if err != nil {
		return nil, WeixinUserAddressListError
	}

	list := make([]*domain.UserAddress, 0)

	for _, userAddress := range userAddresses {
		userAddress, _ = uauc.getUserAddress(ctx, userAddress)

		list = append(list, userAddress)
	}

	return &domain.UserAddressList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (uauc *UserAddressUsecase) CreateUserAddresses(ctx context.Context, userId, provinceAreaCode, cityAreaCode, areaAreaCode uint64, isDefault uint8, name, phone, addressInfo string) (*domain.UserAddress, error) {
	province, err := uauc.arepo.GetByAreaCode(ctx, provinceAreaCode)

	if err != nil {
		return nil, WeixinAreaNotFound
	}

	city, err := uauc.arepo.GetByAreaCode(ctx, cityAreaCode)

	if err != nil {
		return nil, WeixinAreaNotFound
	}

	if province.Data.AreaCode != city.Data.ParentAreaCode {
		return nil, WeixinAreaNotFound
	}

	area, err := uauc.arepo.GetByAreaCode(ctx, areaAreaCode)

	if err != nil {
		return nil, WeixinAreaNotFound
	}

	if city.Data.AreaCode != area.Data.ParentAreaCode {
		return nil, WeixinAreaNotFound
	}

	var userAddress *domain.UserAddress

	err = uauc.tm.InTx(ctx, func(ctx context.Context) error {
		if isDefault == 1 {
			if err = uauc.repo.UpdateIsDefaults(ctx, userId); err != nil {
				return err
			}
		}

		inUserAddress := domain.NewUserAddress(ctx, userId, province.Data.AreaCode, city.Data.AreaCode, area.Data.AreaCode, isDefault, name, phone, addressInfo)
		inUserAddress.SetCreateTime(ctx)
		inUserAddress.SetUpdateTime(ctx)

		userAddress, err = uauc.repo.Save(ctx, inUserAddress)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, WeixinUserAddressCreateError
	}

	userAddress.SetProvinceAreaName(ctx, province.Data.AreaName)
	userAddress.SetCityAreaName(ctx, city.Data.AreaName)
	userAddress.SetAreaAreaName(ctx, area.Data.AreaName)

	return userAddress, nil
}

func (uauc *UserAddressUsecase) UpdateUserAddresses(ctx context.Context, userId, userAddressId, provinceAreaCode, cityAreaCode, areaAreaCode uint64, isDefault uint8, name, phone, addressInfo string) (*domain.UserAddress, error) {
	inUserAddress, err := uauc.repo.GetById(ctx, userId, userAddressId)

	if err != nil {
		return nil, WeixinUserAddressNotFound
	}

	province, err := uauc.arepo.GetByAreaCode(ctx, provinceAreaCode)

	if err != nil {
		return nil, WeixinAreaNotFound
	}

	city, err := uauc.arepo.GetByAreaCode(ctx, cityAreaCode)

	if err != nil {
		return nil, WeixinAreaNotFound
	}

	if province.Data.AreaCode != city.Data.ParentAreaCode {
		return nil, WeixinAreaNotFound
	}

	area, err := uauc.arepo.GetByAreaCode(ctx, areaAreaCode)

	if err != nil {
		return nil, WeixinAreaNotFound
	}

	if city.Data.AreaCode != area.Data.ParentAreaCode {
		return nil, WeixinAreaNotFound
	}

	var userAddress *domain.UserAddress

	err = uauc.tm.InTx(ctx, func(ctx context.Context) error {
		if isDefault == 1 {
			if err = uauc.repo.UpdateIsDefaults(ctx, userId); err != nil {
				return err
			}
		}

		inUserAddress.SetName(ctx, name)
		inUserAddress.SetPhone(ctx, phone)
		inUserAddress.SetProvinceAreaCode(ctx, province.Data.AreaCode)
		inUserAddress.SetCityAreaCode(ctx, city.Data.AreaCode)
		inUserAddress.SetAreaAreaCode(ctx, area.Data.AreaCode)
		inUserAddress.SetAddressInfo(ctx, addressInfo)
		inUserAddress.SetIsDefault(ctx, isDefault)
		inUserAddress.SetUpdateTime(ctx)

		userAddress, err = uauc.repo.Update(ctx, inUserAddress)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, WeixinUserAddressUpdateError
	}

	userAddress.SetProvinceAreaName(ctx, province.Data.AreaName)
	userAddress.SetCityAreaName(ctx, city.Data.AreaName)
	userAddress.SetAreaAreaName(ctx, area.Data.AreaName)

	return userAddress, nil
}

func (uauc *UserAddressUsecase) UpdateDefaultUserAddresses(ctx context.Context, userId, userAddressId uint64, isDefault uint8) (*domain.UserAddress, error) {
	inUserAddress, err := uauc.repo.GetById(ctx, userId, userAddressId)

	if err != nil {
		return nil, WeixinUserAddressNotFound
	}

	var userAddress *domain.UserAddress

	err = uauc.tm.InTx(ctx, func(ctx context.Context) error {
		if isDefault == 1 {
			if err = uauc.repo.UpdateIsDefaults(ctx, userId); err != nil {
				return err
			}
		}

		inUserAddress.SetIsDefault(ctx, isDefault)
		inUserAddress.SetUpdateTime(ctx)

		userAddress, err = uauc.repo.Update(ctx, inUserAddress)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, WeixinUserAddressUpdateError
	}

	userAddress, _ = uauc.getUserAddress(ctx, userAddress)

	return userAddress, nil
}

func (uauc *UserAddressUsecase) DeleteUserAddresses(ctx context.Context, userId, userAddressId uint64) error {
	inUserAddress, err := uauc.repo.GetById(ctx, userId, userAddressId)

	if err != nil {
		return WeixinUserAddressNotFound
	}

	if err := uauc.repo.Delete(ctx, inUserAddress); err != nil {
		return WeixinUserAddressDeleteError
	}

	return nil
}

func (uauc *UserAddressUsecase) getUserAddress(ctx context.Context, userAddress *domain.UserAddress) (*domain.UserAddress, error) {
	if province, err := uauc.arepo.GetByAreaCode(ctx, userAddress.ProvinceAreaCode); err == nil {
		userAddress.ProvinceAreaName = province.Data.AreaName

		if city, err := uauc.arepo.GetByAreaCode(ctx, userAddress.CityAreaCode); err == nil {
			if city.Data.ParentAreaCode == province.Data.AreaCode {
				userAddress.CityAreaName = city.Data.AreaName

				if area, err := uauc.arepo.GetByAreaCode(ctx, userAddress.AreaAreaCode); err == nil {
					if area.Data.ParentAreaCode == city.Data.AreaCode {
						userAddress.AreaAreaName = area.Data.AreaName
					}
				}
			}
		}
	}

	return userAddress, nil
}

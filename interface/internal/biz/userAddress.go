package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/weixin/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"
)

type UserAddressRepo interface {
	List(context.Context, uint64, uint64, uint64) (*v1.ListUserAddressesReply, error)
	Save(context.Context, uint64, uint64, uint64, uint64, uint32, string, string, string) (*v1.CreateUserAddressesReply, error)
	Update(context.Context, uint64, uint64, uint64, uint64, uint64, uint32, string, string, string) (*v1.UpdateUserAddressesReply, error)
	UpdateDefault(context.Context, uint64, uint64, uint32) (*v1.UpdateDefaultUserAddressesReply, error)
	Delete(context.Context, uint64, uint64) (*v1.DeleteUserAddressesReply, error)
}

type UserAddressUsecase struct {
	repo UserAddressRepo
	conf *conf.Data
	log  *log.Helper
}

func NewUserAddressUsecase(repo UserAddressRepo, conf *conf.Data, logger log.Logger) *UserAddressUsecase {
	return &UserAddressUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (uauc *UserAddressUsecase) ListUserAddresses(ctx context.Context, pageNum, pageSize, userId uint64) (*v1.ListUserAddressesReply, error) {
	if pageSize == 0 {
		pageSize = uint64(uauc.conf.Database.PageSize)
	}

	list, err := uauc.repo.List(ctx, userId, pageNum, pageSize)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_USER_ADDRESS_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (uauc *UserAddressUsecase) CreateUserAddresses(ctx context.Context, userId, provinceAreaCode, cityAreaCode, areaAreaCode uint64, isDefault uint32, name, phone, addressInfo string) (*v1.CreateUserAddressesReply, error) {
	userAddress, err := uauc.repo.Save(ctx, userId, provinceAreaCode, cityAreaCode, areaAreaCode, isDefault, name, phone, addressInfo)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_CREATE_USER_ADDRESS_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userAddress, nil
}

func (uauc *UserAddressUsecase) UpdateUserAddresses(ctx context.Context, userId, userAddressId, provinceAreaCode, cityAreaCode, areaAreaCode uint64, isDefault uint32, name, phone, addressInfo string) (*v1.UpdateUserAddressesReply, error) {
	userAddress, err := uauc.repo.Update(ctx, userId, userAddressId, provinceAreaCode, cityAreaCode, areaAreaCode, isDefault, name, phone, addressInfo)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_UPDATE_USER_ADDRESS_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userAddress, nil
}

func (uauc *UserAddressUsecase) UpdateDefaultUserAddresses(ctx context.Context, userId, userAddressId uint64, isDefault uint32) (*v1.UpdateDefaultUserAddressesReply, error) {
	userAddress, err := uauc.repo.UpdateDefault(ctx, userId, userAddressId, isDefault)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_UPDATE_DEFAULT_USER_ADDRESS_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userAddress, nil
}

func (uauc *UserAddressUsecase) DeleteUserAddresses(ctx context.Context, userId, userAddressId uint64) (*v1.DeleteUserAddressesReply, error) {
	userAddress, err := uauc.repo.Delete(ctx, userId, userAddressId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_DELETE_USER_ADDRESS_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userAddress, nil
}

package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	cv1 "interface/api/service/company/v1"
	v1 "interface/api/service/douyin/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"
)

type UserStoreRepo interface {
	List(context.Context, uint64, uint64, uint64) (*v1.ListJinritemaiStoresReply, error)
	Save(context.Context, uint64, uint64, string) (*cv1.CreateJinritemaiStoreCompanyProductsReply, error)
	Update(context.Context, uint64, string) (*v1.UpdateJinritemaiStoresReply, error)
	Delete(context.Context, uint64, string) (*v1.DeleteJinritemaiStoresReply, error)
}

type UserStoreUsecase struct {
	repo UserStoreRepo
	conf *conf.Data
	log  *log.Helper
}

func NewUserStoreUsecase(repo UserStoreRepo, conf *conf.Data, logger log.Logger) *UserStoreUsecase {
	return &UserStoreUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (usuc *UserStoreUsecase) ListUserStores(ctx context.Context, pageNum, pageSize, userId uint64) (*v1.ListJinritemaiStoresReply, error) {
	if pageSize == 0 {
		pageSize = uint64(usuc.conf.Database.PageSize)
	}

	list, err := usuc.repo.List(ctx, userId, pageNum, pageSize)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_USER_STORE_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (usuc *UserStoreUsecase) CreateUserStores(ctx context.Context, userId, productId uint64, openDouyinUserIds string) (*cv1.CreateJinritemaiStoreCompanyProductsReply, error) {
	userStore, err := usuc.repo.Save(ctx, userId, productId, openDouyinUserIds)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_CREATE_USER_STORE_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userStore, nil
}

func (usuc *UserStoreUsecase) UpdateUserStores(ctx context.Context, userId uint64, stores string) (*v1.UpdateJinritemaiStoresReply, error) {
	userStore, err := usuc.repo.Update(ctx, userId, stores)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_UPDATE_USER_STORE_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userStore, nil
}

func (usuc *UserStoreUsecase) DeleteUserStores(ctx context.Context, userId uint64, stores string) (*v1.DeleteJinritemaiStoresReply, error) {
	userStore, err := usuc.repo.Delete(ctx, userId, stores)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_DELETE_USER_STORE_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userStore, nil
}

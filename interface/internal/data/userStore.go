package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	cv1 "interface/api/service/company/v1"
	v1 "interface/api/service/douyin/v1"
	"interface/internal/biz"
)

type userStoreRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserStoreRepo(data *Data, logger log.Logger) biz.UserStoreRepo {
	return &userStoreRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (usr *userStoreRepo) List(ctx context.Context, userId, pageNum, pageSize uint64) (*v1.ListJinritemaiStoresReply, error) {
	list, err := usr.data.douyinuc.ListJinritemaiStores(ctx, &v1.ListJinritemaiStoresRequest{
		PageNum:  pageNum,
		PageSize: pageSize,
		UserId:   userId,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (usr *userStoreRepo) Save(ctx context.Context, userId, productId uint64, openDouyinUserIds string) (*cv1.CreateJinritemaiStoreCompanyProductsReply, error) {
	store, err := usr.data.companyuc.CreateJinritemaiStoreCompanyProducts(ctx, &cv1.CreateJinritemaiStoreCompanyProductsRequest{
		UserId:            userId,
		OpenDouyinUserIds: openDouyinUserIds,
		ProductId:         productId,
	})

	if err != nil {
		return nil, err
	}

	return store, err
}

func (usr *userStoreRepo) Update(ctx context.Context, userId uint64, stores string) (*v1.UpdateJinritemaiStoresReply, error) {
	store, err := usr.data.douyinuc.UpdateJinritemaiStores(ctx, &v1.UpdateJinritemaiStoresRequest{
		UserId: userId,
		Stores: stores,
	})

	if err != nil {
		return nil, err
	}

	return store, err
}

func (usr *userStoreRepo) Delete(ctx context.Context, userId uint64, stores string) (*v1.DeleteJinritemaiStoresReply, error) {
	store, err := usr.data.douyinuc.DeleteJinritemaiStores(ctx, &v1.DeleteJinritemaiStoresRequest{
		UserId: userId,
		Stores: stores,
	})

	if err != nil {
		return nil, err
	}

	return store, err
}

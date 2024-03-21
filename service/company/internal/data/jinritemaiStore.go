package data

import (
	v1 "company/api/service/douyin/v1"
	"company/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type jinritemaiStoreRepo struct {
	data *Data
	log  *log.Helper
}

func NewJinritemaiStoreRepo(data *Data, logger log.Logger) biz.JinritemaiStoreRepo {
	return &jinritemaiStoreRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (jsr *jinritemaiStoreRepo) List(ctx context.Context, userId uint64) (*v1.ListJinritemaiStoresReply, error) {
	list, err := jsr.data.douyinuc.ListJinritemaiStores(ctx, &v1.ListJinritemaiStoresRequest{
		PageSize: 40,
		UserId:   userId,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (jsr *jinritemaiStoreRepo) Save(ctx context.Context, userId, companyId, productId uint64, openDouyinUserIds, activityUrl string) (*v1.CreateJinritemaiStoresReply, error) {
	list, err := jsr.data.douyinuc.CreateJinritemaiStores(ctx, &v1.CreateJinritemaiStoresRequest{
		UserId:            userId,
		OpenDouyinUserIds: openDouyinUserIds,
		CompanyId:         companyId,
		ProductId:         productId,
		ActivityUrl:       activityUrl,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

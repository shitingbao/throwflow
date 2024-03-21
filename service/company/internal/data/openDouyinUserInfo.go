package data

import (
	v1 "company/api/service/douyin/v1"
	"company/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type openDouyinUserInfoRepo struct {
	data *Data
	log  *log.Helper
}

func NewOpenDouyinUserInfoRepo(data *Data, logger log.Logger) biz.OpenDouyinUserInfoRepo {
	return &openDouyinUserInfoRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (oduir *openDouyinUserInfoRepo) ListByProductId(ctx context.Context, productId string) (*v1.ListOpenDouyinUserInfosByProductIdReply, error) {
	list, err := oduir.data.douyinuc.ListOpenDouyinUserInfosByProductId(ctx, &v1.ListOpenDouyinUserInfosByProductIdRequest{
		ProductId: productId,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

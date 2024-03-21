package data

import (
	v1 "company/api/service/douyin/v1"
	"company/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type doukeProductShareRepo struct {
	data *Data
	log  *log.Helper
}

func NewDoukeProductShareRepo(data *Data, logger log.Logger) biz.DoukeProductShareRepo {
	return &doukeProductShareRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (dpsr *doukeProductShareRepo) Save(ctx context.Context, productUrl, externalInfo string) (*v1.CreateDoukeProductSharesReply, error) {
	productShare, err := dpsr.data.douyinuc.CreateDoukeProductShares(ctx, &v1.CreateDoukeProductSharesRequest{
		ProductUrl:   productUrl,
		ExternalInfo: externalInfo,
	})

	if err != nil {
		return nil, err
	}

	return productShare, err
}

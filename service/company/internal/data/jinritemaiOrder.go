package data

import (
	v1 "company/api/service/douyin/v1"
	"company/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type jinritemaiOrderRepo struct {
	data *Data
	log  *log.Helper
}

func NewJinritemaiOrderRepo(data *Data, logger log.Logger) biz.JinritemaiOrderRepo {
	return &jinritemaiOrderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (jor *jinritemaiOrderRepo) ListCommissionRate(ctx context.Context, content string) (*v1.ListCommissionRateJinritemaiOrdersReply, error) {
	list, err := jor.data.douyinuc.ListCommissionRateJinritemaiOrders(ctx, &v1.ListCommissionRateJinritemaiOrdersRequest{
		Content: content,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (jor *jinritemaiOrderRepo) GetIsTop(ctx context.Context, productId uint64) (*v1.GetIsTopJinritemaiOrdersReply, error) {
	jinritemaiOrder, err := jor.data.douyinuc.GetIsTopJinritemaiOrders(ctx, &v1.GetIsTopJinritemaiOrdersRequest{
		ProductId: productId,
	})

	if err != nil {
		return nil, err
	}

	return jinritemaiOrder, err
}

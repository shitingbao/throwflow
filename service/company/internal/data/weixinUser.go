package data

import (
	v1 "company/api/service/weixin/v1"
	"company/internal/biz"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type weixinUserRepo struct {
	data *Data
	log  *log.Helper
}

func NewWeixinUserRepo(data *Data, logger log.Logger) biz.WeixinUserRepo {
	return &weixinUserRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (wur *weixinUserRepo) ListByIds(ctx context.Context, phone, keyword, userIds string) (*v1.ListByIdsReply, error) {
	list, err := wur.data.weixinuc.ListByIds(ctx, &v1.ListByIdsRequest{
		Phone:   phone,
		Keyword: keyword,
		UserIds: userIds,
	})

	if err != nil {
		return nil, err
	}

	return list, nil
}

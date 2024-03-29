package data

import (
	"context"
	v1 "douyin/api/service/weixin/v1"
	"douyin/internal/biz"
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

func (wur *weixinUserRepo) GetById(ctx context.Context, userId uint64) (*v1.GetByIdUsersReply, error) {
	weixinUser, err := wur.data.weixinuc.GetByIdUsers(ctx, &v1.GetByIdUsersRequest{
		UserId: userId,
	})

	if err != nil {
		return nil, err
	}

	return weixinUser, err
}

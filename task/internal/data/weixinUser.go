package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "task/api/service/weixin/v1"
	"task/internal/biz"
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

func (wur *weixinUserRepo) Sync(ctx context.Context) (*v1.SyncIntegralUsersReply, error) {
	user, err := wur.data.weixinuc.SyncIntegralUsers(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return user, err
}

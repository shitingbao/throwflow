package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "task/api/service/weixin/v1"
	"task/internal/biz"
)

type weixinUserBalanceRepo struct {
	data *Data
	log  *log.Helper
}

func NewWeixinUserBalanceRepo(data *Data, logger log.Logger) biz.WeixinUserBalanceRepo {
	return &weixinUserBalanceRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (wur *weixinUserBalanceRepo) Sync(ctx context.Context) (*v1.SyncUserBalancesReply, error) {
	userBalance, err := wur.data.weixinuc.SyncUserBalances(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return userBalance, err
}

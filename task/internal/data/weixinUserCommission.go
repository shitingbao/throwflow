package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "task/api/service/weixin/v1"
	"task/internal/biz"
)

type weixinUserCommissionRepo struct {
	data *Data
	log  *log.Helper
}

func NewWeixinUserCommissionRepo(data *Data, logger log.Logger) biz.WeixinUserCommissionRepo {
	return &weixinUserCommissionRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (wucr *weixinUserCommissionRepo) SyncOrder(ctx context.Context) (*v1.SyncOrderUserCommissionsReply, error) {
	userCommission, err := wucr.data.weixinuc.SyncOrderUserCommissions(ctx, &v1.SyncOrderUserCommissionsRequest{})

	if err != nil {
		return nil, err
	}

	return userCommission, err
}

func (wucr *weixinUserCommissionRepo) SyncCostOrder(ctx context.Context) (*v1.SyncCostOrderUserCommissionsReply, error) {
	userCommission, err := wucr.data.weixinuc.SyncCostOrderUserCommissions(ctx, &v1.SyncCostOrderUserCommissionsRequest{})

	if err != nil {
		return nil, err
	}

	return userCommission, err
}

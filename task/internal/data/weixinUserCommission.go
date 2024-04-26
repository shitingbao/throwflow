package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (wucr *weixinUserCommissionRepo) SyncTask(ctx context.Context) (*v1.SyncTaskUserCommissionsReply, error) {
	userCommission, err := wucr.data.weixinuc.SyncTaskUserCommissions(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return userCommission, err
}

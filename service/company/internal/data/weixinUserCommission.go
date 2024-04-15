package data

import (
	v1 "company/api/service/weixin/v1"
	"company/internal/biz"
	"context"

	"github.com/go-kratos/kratos/v2/log"
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

func (wur *weixinUserCommissionRepo) CreateTaskUserCommissions(ctx context.Context, UserId, TaskRelationId uint64, FlowPoint string, Commission float64, SuccessTime string) (*v1.CreateTaskUserCommissionsReply, error) {
	commission, err := wur.data.weixinuc.CreateTaskUserCommissions(ctx, &v1.CreateTaskUserCommissionsRequest{
		UserId:         UserId,
		TaskRelationId: TaskRelationId,
		FlowPoint:      FlowPoint,
		Commission:     Commission,
		SuccessTime:    SuccessTime,
	})

	if err != nil {
		return nil, err
	}

	return commission, nil
}

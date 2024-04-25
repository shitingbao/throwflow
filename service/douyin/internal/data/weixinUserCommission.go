package data

import (
	"context"
	v1 "douyin/api/service/weixin/v1"
	"douyin/internal/biz"
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

func (wucr *weixinUserCommissionRepo) CreateOrder(ctx context.Context, totalPayAmount, commission float64, clientKey, openId, orderId, flowPoint, paySuccessTime string) (*v1.CreateOrderUserCommissionsReply, error) {
	userCommission, err := wucr.data.weixinuc.CreateOrderUserCommissions(ctx, &v1.CreateOrderUserCommissionsRequest{
		ClientKey:      clientKey,
		OpenId:         openId,
		OrderId:        orderId,
		FlowPoint:      flowPoint,
		PaySuccessTime: paySuccessTime,
		TotalPayAmount: totalPayAmount,
		Commission:     commission,
	})

	if err != nil {
		return nil, err
	}

	return userCommission, err
}

func (wucr *weixinUserCommissionRepo) CreateCostOrder(ctx context.Context, userId uint64, totalPayAmount, commission float64, orderId, productId, flowPoint, paySuccessTime string) (*v1.CreateCostOrderUserCommissionsReply, error) {
	userCommission, err := wucr.data.weixinuc.CreateCostOrderUserCommissions(ctx, &v1.CreateCostOrderUserCommissionsRequest{
		UserId:         userId,
		OrderId:        orderId,
		ProductId:      productId,
		FlowPoint:      flowPoint,
		PaySuccessTime: paySuccessTime,
		TotalPayAmount: totalPayAmount,
		Commission:     commission,
	})

	if err != nil {
		return nil, err
	}

	return userCommission, err
}

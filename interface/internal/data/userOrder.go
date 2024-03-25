package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/weixin/v1"
	"interface/internal/biz"
)

type userOrderRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserOrderRepo(data *Data, logger log.Logger) biz.UserOrderRepo {
	return &userOrderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (uor *userOrderRepo) Create(ctx context.Context, userId, parentUserId, organizationId uint64, payAmount float64, clientIp string) (*v1.CreateUserOrdersReply, error) {
	userOrder, err := uor.data.weixinuc.CreateUserOrders(ctx, &v1.CreateUserOrdersRequest{
		UserId:         userId,
		ParentUserId:   parentUserId,
		OrganizationId: organizationId,
		PayAmount:      payAmount,
		ClientIp:       clientIp,
	})

	if err != nil {
		return nil, err
	}

	return userOrder, err
}

func (uor *userOrderRepo) Upgrade(ctx context.Context, userId, organizationId uint64, payAmount float64, clientIp string) (*v1.UpgradeUserOrdersReply, error) {
	userOrder, err := uor.data.weixinuc.UpgradeUserOrders(ctx, &v1.UpgradeUserOrdersRequest{
		UserId:         userId,
		OrganizationId: organizationId,
		PayAmount:      payAmount,
		ClientIp:       clientIp,
	})

	if err != nil {
		return nil, err
	}

	return userOrder, err
}

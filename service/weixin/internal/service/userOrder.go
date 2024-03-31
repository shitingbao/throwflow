package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"math"
	"time"
	v1 "weixin/api/weixin/v1"
	"weixin/internal/pkg/tool"
)

func (ws *WeixinService) ListUserOrders(ctx context.Context, in *v1.ListUserOrdersRequest) (*v1.ListUserOrdersReply, error) {
	userOrders, err := ws.uouc.ListUserOrders(ctx, in.PageNum, in.PageSize)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListUserOrdersReply_UserOrder, 0)

	for _, userOrder := range userOrders.List {
		list = append(list, &v1.ListUserOrdersReply_UserOrder{
			NickName:  userOrder.NickName,
			AvatarUrl: userOrder.AvatarUrl,
			PayTime:   tool.TimeToString("2006/01/02 15:04", *userOrder.PayTime),
		})
	}

	totalPage := uint64(math.Ceil(float64(userOrders.Total) / float64(userOrders.PageSize)))

	return &v1.ListUserOrdersReply{
		Code: 200,
		Data: &v1.ListUserOrdersReply_Data{
			PageNum:   userOrders.PageNum,
			PageSize:  userOrders.PageSize,
			Total:     userOrders.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ws *WeixinService) CreateUserOrders(ctx context.Context, in *v1.CreateUserOrdersRequest) (*v1.CreateUserOrdersReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	order, err := ws.uouc.CreateUserOrders(ctx, in.UserId, in.ParentUserId, in.OrganizationId, in.PayAmount, in.ClientIp)

	if err != nil {
		return nil, err
	}

	return &v1.CreateUserOrdersReply{
		Code: 200,
		Data: &v1.CreateUserOrdersReply_Data{
			TimeStamp:  order.TimeStamp,
			NonceStr:   order.NonceStr,
			Package:    order.Package,
			SignType:   order.SignType,
			PaySign:    order.PaySign,
			OutTradeNo: order.OutTradeNo,
			PayAmount:  order.PayAmount,
			LevelName:  order.LevelName,
		},
	}, nil
}

func (ws *WeixinService) UpgradeUserOrders(ctx context.Context, in *v1.UpgradeUserOrdersRequest) (*v1.UpgradeUserOrdersReply, error) {
	order, err := ws.uouc.UpgradeUserOrders(ctx, in.UserId, in.OrganizationId, in.PayAmount, in.ClientIp)

	if err != nil {
		return nil, err
	}

	return &v1.UpgradeUserOrdersReply{
		Code: 200,
		Data: &v1.UpgradeUserOrdersReply_Data{
			TimeStamp:  order.TimeStamp,
			NonceStr:   order.NonceStr,
			Package:    order.Package,
			SignType:   order.SignType,
			PaySign:    order.PaySign,
			OutTradeNo: order.OutTradeNo,
			PayAmount:  order.PayAmount,
			LevelName:  order.LevelName,
		},
	}, nil
}

func (ws *WeixinService) AsyncNotificationUserOrders(ctx context.Context, in *v1.AsyncNotificationUserOrdersRequest) (*v1.AsyncNotificationUserOrdersReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := ws.uouc.AsyncNotificationUserOrders(ctx, in.Content); err != nil {
		return nil, err
	}

	return &v1.AsyncNotificationUserOrdersReply{
		Code: 200,
		Data: &v1.AsyncNotificationUserOrdersReply_Data{},
	}, nil
}

func (ws *WeixinService) SyncQrCodeUserOrganizationRelations(ctx context.Context, in *emptypb.Empty) (*v1.SyncQrCodeUserOrganizationRelationsReply, error) {
	ws.log.Infof("同步微信用户账单机构关系二维码数据, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := ws.uouc.SyncQrCodeUserOrganizationRelations(ctx); err != nil {
		return nil, err
	}

	ws.log.Infof("同步微信用户账单机构关系二维码数据, 结束时间 %s \n", time.Now())

	return &v1.SyncQrCodeUserOrganizationRelationsReply{
		Code: 200,
		Data: &v1.SyncQrCodeUserOrganizationRelationsReply_Data{},
	}, nil
}

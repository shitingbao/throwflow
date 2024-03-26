package service

import (
	v1 "common/api/common/v1"
	"context"
)

func (cs *CommonService) Pay(ctx context.Context, in *v1.PayRequest) (*v1.PayReply, error) {
	data, err := cs.puc.Pay(ctx, in.OrganizationId, in.TotalFee, in.OutTradeNo, in.Content, in.NonceStr, in.OpenId, in.ClientIp)

	if err != nil {
		return nil, err
	}

	return &v1.PayReply{
		Code: 200,
		Data: &v1.PayReply_Data{
			TimeStamp: data.TimeStamp,
			NonceStr:  data.NonceStr,
			Package:   data.Package,
			SignType:  data.SignType,
			PaySign:   data.PaySign,
		},
	}, nil
}

func (cs *CommonService) PayAsyncNotification(ctx context.Context, in *v1.PayAsyncNotificationRequest) (*v1.PayAsyncNotificationReply, error) {
	data, err := cs.puc.PayAsyncNotification(ctx, in.Content)

	if err != nil {
		return nil, err
	}

	return &v1.PayAsyncNotificationReply{
		Code: 200,
		Data: &v1.PayAsyncNotificationReply_Data{
			OutTradeNo:       data.OutTradeNo,
			OutTransactionId: data.OutTransactionId,
			TransactionId:    data.TransactionId,
			PayTime:          data.TimeEnd,
			PayAmount:        float64(data.TotalFee) / 100,
		},
	}, nil
}

func (cs *CommonService) Divide(ctx context.Context, in *v1.DivideRequest) (*v1.DivideReply, error) {
	if _, err := cs.puc.Divide(ctx, in.OutTradeNo, in.TransactionNo); err != nil {
		return nil, err
	}

	return &v1.DivideReply{
		Code: 200,
		Data: &v1.DivideReply_Data{},
	}, nil
}

package service

import (
	"context"
	v1 "interface/api/interface/v1"
)

func (is *InterfaceService) SendSms(ctx context.Context, in *v1.SendSmsRequest) (*v1.SendSmsReply, error) {
	_, err := is.suc.SendSms(ctx, in.Phone, in.Types)

	if err != nil {
		return nil, err
	}

	return &v1.SendSmsReply{
		Code: 200,
		Data: &v1.SendSmsReply_Data{},
	}, nil
}

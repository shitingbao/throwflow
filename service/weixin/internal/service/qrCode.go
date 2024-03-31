package service

import (
	"context"
	"time"
	v1 "weixin/api/weixin/v1"
)

func (ws *WeixinService) GetQrCodes(ctx context.Context, in *v1.GetQrCodesRequest) (*v1.GetQrCodesReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	staticUrl, err := ws.qcuc.GetQrCodes(ctx, in.OrganizationId, in.Scene)

	if err != nil {
		return nil, err
	}

	return &v1.GetQrCodesReply{
		Code: 200,
		Data: &v1.GetQrCodesReply_Data{
			StaticUrl: staticUrl,
		},
	}, nil
}

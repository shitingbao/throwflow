package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "weixin/api/service/common/v1"
	"weixin/internal/biz"
)

type payRepo struct {
	data *Data
	log  *log.Helper
}

func NewPayRepo(data *Data, logger log.Logger) biz.PayRepo {
	return &payRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (pr *payRepo) Pay(ctx context.Context, organizationId uint64, totalFee float64, outTradeNo, content, nonceStr, openId, clientIp string) (*v1.PayReply, error) {
	pay, err := pr.data.commonuc.Pay(ctx, &v1.PayRequest{
		OrganizationId: organizationId,
		OutTradeNo:     outTradeNo,
		Content:        content,
		NonceStr:       nonceStr,
		OpenId:         openId,
		ClientIp:       clientIp,
		TotalFee:       totalFee,
	})
	
	if err != nil {
		return nil, err
	}

	return pay, err
}

func (pr *payRepo) PayAsyncNotification(ctx context.Context, content string) (*v1.PayAsyncNotificationReply, error) {
	pay, err := pr.data.commonuc.PayAsyncNotification(ctx, &v1.PayAsyncNotificationRequest{
		Content: content,
	})

	if err != nil {
		return nil, err
	}

	return pay, err
}

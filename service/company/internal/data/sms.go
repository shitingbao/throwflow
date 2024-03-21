package data

import (
	v1 "company/api/service/common/v1"
	"company/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type smsRepo struct {
	data *Data
	log  *log.Helper
}

func NewSmsRepo(data *Data, logger log.Logger) biz.SmsRepo {
	return &smsRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

type SmsRepo interface {
	Send(context.Context, string, string, string, string) (*v1.SendSmsReply, error)
}

func (sr *smsRepo) Send(ctx context.Context, phone, content, types, ip string) (*v1.SendSmsReply, error) {
	sms, err := sr.data.commonuc.SendSms(ctx, &v1.SendSmsRequest{
		Phone:   phone,
		Content: content,
		Types:   types,
		Ip:      ip,
	})

	if err != nil {
		return nil, err
	}

	return sms, err
}

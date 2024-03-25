package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/common/v1"
	"interface/internal/biz"
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

func (sr *smsRepo) Verify(ctx context.Context, phone, types, code string) (*v1.VerifySmsReply, error) {
	sms, err := sr.data.commonuc.VerifySms(ctx, &v1.VerifySmsRequest{
		Phone: phone,
		Code:  code,
		Types: types,
	})

	if err != nil {
		return nil, err
	}

	return sms, err
}

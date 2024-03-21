package biz

import (
	v1 "company/api/service/common/v1"
	"context"
)

type SmsRepo interface {
	Send(context.Context, string, string, string, string) (*v1.SendSmsReply, error)
}

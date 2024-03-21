package biz

import (
	v1 "company/api/service/weixin/v1"
	"context"
)

type WeixinUserScanRecordRepo interface {
	Create(ctx context.Context, userId uint64) (*v1.CreateUserScanRecordsReply, error)
}

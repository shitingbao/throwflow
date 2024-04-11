package biz

import (
	v1 "company/api/service/douyin/v1"
	"context"
)

type DoukeOrderRepo interface {
	Get(context.Context, uint64, string, string, string) (*v1.GetDoukeOrdersReply, error)
}

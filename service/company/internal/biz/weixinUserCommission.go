package biz

import (
	v1 "company/api/service/weixin/v1"
	"context"
)

type WeixinUserCommissionRepo interface {
	CreateTaskUserCommissions(ctx context.Context, UserId, TaskRelationId uint64, FlowPoint string, Commission float64, SuccessTime string) (*v1.CreateTaskUserCommissionsReply, error)
}

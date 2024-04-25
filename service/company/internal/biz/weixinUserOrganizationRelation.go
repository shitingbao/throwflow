package biz

import (
	v1 "company/api/service/weixin/v1"
	"context"
)

type WeixinUserOrganizationRelationRepo interface {
	GetUserOrganizationRelations(context.Context, uint64) (*v1.GetUserOrganizationRelationsReply, error)
}

package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	v1 "weixin/api/service/company/v1"
)

var (
	WeixinCompanyTaskAccountRelationNotFound = errors.NotFound("WEIXIN_COMPANY_TASK_ACCOUNT_RELATION_NOT_FOUND", "微信用户领取任务不存在")
)

type CompanyTaskRepo interface {
	GetCompanyTaskAccountRelation(context.Context, uint64) (*v1.GetCompanyTaskAccountRelationsReply, error)
}

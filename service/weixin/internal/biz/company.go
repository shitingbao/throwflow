package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	v1 "weixin/api/service/company/v1"
)

var (
	WeixinCompanyListError = errors.NotFound("WEIXIN_COMAPNY_LIST_ERROR", "企业列表获取失败")
	WeixinCompanyNotFound  = errors.NotFound("WEIXIN_COMAPNY_NOT_FOUND", "企业不存在")
)

type CompanyRepo interface {
	List(context.Context, string) (*v1.ListCompanysReply, error)
}

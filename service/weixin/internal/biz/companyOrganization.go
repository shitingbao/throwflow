package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	v1 "weixin/api/service/company/v1"
)

var (
	WeixinCompanyOrganizationNotFound = errors.NotFound("WEIXIN_COMPANY_ORGANIZATION_NOT_FOUND", "企业机构不存在")
)

type CompanyOrganizationRepo interface {
	Get(context.Context, uint64) (*v1.GetCompanyOrganizationsReply, error)
	GetByOrganizationCode(context.Context, string) (*v1.GetCompanyOrganizationByOrganizationCodesReply, error)
	List(context.Context) (*v1.ListCompanyOrganizationsReply, error)
}

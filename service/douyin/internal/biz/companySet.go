package biz

import (
	"context"
	v1 "douyin/api/service/company/v1"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	DouyinCompanySetGetError = errors.InternalServer("DOUYIN_COMPANY_SET_GET_ERROR", "公司设置获取失败")
)

type CompanySetRepo interface {
	Get(context.Context, uint64, string, string) (*v1.GetCompanySetsReply, error)
}

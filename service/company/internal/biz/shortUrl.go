package biz

import (
	v1 "company/api/service/common/v1"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	CompanyShortUrlCreateError = errors.InternalServer("COMPANY_SHORT_URL_CREATE_ERROR", "短链接创建失败")
)

type ShortUrlRepo interface {
	Create(context.Context, string) (*v1.CreateShortUrlReply, error)
}

package biz

import (
	v1 "company/api/service/common/v1"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	CompanyShortCodeCreateError = errors.InternalServer("COMPANY_SHORT_CODE_CREATE_ERROR", "短码创建失败")
)

type ShortCodeRepo interface {
	Create(ctx context.Context) (*v1.CreateShortCodeReply, error)
}

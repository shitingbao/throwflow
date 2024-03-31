package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	v1 "weixin/api/service/common/v1"
)

var (
	WeixinShortCodeCreateError = errors.InternalServer("WEIXIN_SHORT_CODE_CREATE_ERROR", "短码创建失败")
)

type ShortCodeRepo interface {
	Create(context.Context) (*v1.CreateShortCodeReply, error)
}

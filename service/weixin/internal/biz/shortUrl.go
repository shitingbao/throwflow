package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	v1 "weixin/api/service/common/v1"
)

var (
	WeixinShortUrlCreateError = errors.InternalServer("WEIXIN_SHORT_URL_CREATE_ERROR", "短链接创建失败")
)

type ShortUrlRepo interface {
	Create(context.Context, string) (*v1.CreateShortUrlReply, error)
}

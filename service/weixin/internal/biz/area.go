package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	v1 "weixin/api/service/common/v1"
)

var (
	WeixinAreaNotFound = errors.NotFound("WEIXIN_AREA_NOT_FOUND", "行政区划不存在")
)

type AreaRepo interface {
	GetByAreaCode(context.Context, uint64) (*v1.GetAreasReply, error)
}

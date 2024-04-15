package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"weixin/internal/domain"
)

var (
	WeixinDjAwemeNotFound  = errors.NotFound("WEIXIN_DJ_AWEME_NOT_FOUND", "MCN达人不存在")
	WeixinDjAwemeListError = errors.InternalServer("WEIXIN_DJ_AWEME_LIST_ERROR", "MCN达人列表获取失败")
)

type DjAwemeRepo interface {
	Get(context.Context, string, string) (*domain.DjAweme, error)
	List(context.Context, string, string) ([]*domain.DjAweme, error)
	SaveIndex(context.Context)
}

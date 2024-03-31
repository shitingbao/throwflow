package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	v1 "weixin/api/service/common/v1"
)

var (
	WeixinKuaidiInfoGetError = errors.InternalServer("WEIXIN_KUAIDI_INFO_GET_ERROR", "运单号信息获取失败")
)

type KuaidiInfoRepo interface {
	Get(context.Context, string, string, string) (*v1.GetKuaidiInfosReply, error)
}

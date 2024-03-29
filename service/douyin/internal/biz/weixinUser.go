package biz

import (
	"context"
	v1 "douyin/api/service/weixin/v1"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	DouyinWeixinUserNotFound = errors.NotFound("DOUYIN_WEIXIN_USER_NOT_FOUND", "微信用户不存在")
)

type WeixinUserRepo interface {
	GetById(context.Context, uint64) (*v1.GetByIdUsersReply, error)
}

package biz

import (
	"context"
	v1 "douyin/api/service/weixin/v1"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	DouyinWeixinUserOpenDouyinNotFound             = errors.NotFound("DOUYIN_WEIXIN_USER_OPEN_DOUYIN_NOT_FOUND", "微信用户关联抖音用户不存在")
	DouyinWeixinUserOpenDouyinListError            = errors.NotFound("DOUYIN_WEIXIN_USER_OPEN_DOUYIN_LIST_ERROR", "微信用户关联抖音用户列表获取失败")
	DouyinWeixinUserOpenDouyinListByCompanyIdError = errors.NotFound("DOUYIN_WEIXIN_USER_OPEN_DOUYIN_LIST_BY_COMPANY_ID_ERROR", "企业关联抖音用户列表获取失败")
)

type WeixinUserOpenDouyinRepo interface {
	Get(context.Context, string, string) (*v1.GetOpenDouyinUsersReply, error)
	List(context.Context, uint64) (*v1.ListOpenDouyinUsersReply, error)
	Update(context.Context, uint64, uint64, string, string, string, string, string, string, string) (*v1.UpdateOpenDouyinUsersReply, error)
	UpdateUserInfos(context.Context, uint64, string, string, string, string, string, string, string) (*v1.UpdateUserInfoOpenDouyinUsersReply, error)
	UpdateUserFans(context.Context, string, string, uint64) (*v1.UpdateUserFansOpenDouyinUsersReply, error)
}

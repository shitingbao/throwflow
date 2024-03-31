package biz

import (
	"context"
	v1 "weixin/api/service/douyin/v1"
)

type OpenDouyinUserInfoRepo interface {
	UpdateCooperativeCodes(context.Context, string, string, string) (*v1.UpdateCooperativeCodeDouyinTokensReply, error)
}

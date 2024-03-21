package biz

import (
	v1 "company/api/service/douyin/v1"
	"context"
)

type QianchuanAdvertiserHistoryRepo interface {
	List(context.Context, string, string) (*v1.ListQianchuanAdvertiserHistorysReply, error)
}

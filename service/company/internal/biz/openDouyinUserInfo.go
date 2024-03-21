package biz

import (
	v1 "company/api/service/douyin/v1"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	CompanyOpenDouyinUserInfoListError = errors.InternalServer("COMPANY_OPEN_DOUYIN_USER_INFO_LIST_ERROR", "达人列表获取失败")
)

type OpenDouyinUserInfoRepo interface {
	ListByProductId(context.Context, string) (*v1.ListOpenDouyinUserInfosByProductIdReply, error)
}

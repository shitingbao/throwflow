package biz

import (
	v1 "company/api/service/douyin/v1"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	CompanyDoukeProductShareCreateError = errors.InternalServer("COMPANY_DOUKE_PRODUCT_SHARE_CREATE_ERROR", "抖客商品分销转链创建失败")
)

type DoukeProductShareRepo interface {
	Save(context.Context, string, string) (*v1.CreateDoukeProductSharesReply, error)
}

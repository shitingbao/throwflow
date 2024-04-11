package biz

import (
	v1 "company/api/service/douyin/v1"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	CompanyDoukeProductShareCreateError = errors.InternalServer("COMPANY_DOUKE_PRODUCT_SHARE_CREATE_ERROR", "该商品已下架，或商家店铺设置了不可推广，暂时无法成本购")
)

type DoukeProductShareRepo interface {
	Save(context.Context, string, string) (*v1.CreateDoukeProductSharesReply, error)
}

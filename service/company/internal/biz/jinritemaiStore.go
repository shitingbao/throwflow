package biz

import (
	v1 "company/api/service/douyin/v1"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	CompanyJinritemaiStoreListError   = errors.InternalServer("COMPANY_JINRITEMAI_STORE_LIST_ERROR", "达人橱窗列表获取失败")
	CompanyJinritemaiStoreCreateError = errors.InternalServer("COMPANY_JINRITEMAI_STORE_CREATE_ERROR", "达人橱窗创建失败")
)

type JinritemaiStoreRepo interface {
	List(context.Context, uint64) (*v1.ListJinritemaiStoresReply, error)
	Save(context.Context, uint64, uint64, uint64, string, string) (*v1.CreateJinritemaiStoresReply, error)
}

package biz

import (
	v1 "company/api/service/douyin/v1"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	CompanyQianchuanReportProductListError = errors.InternalServer("COMPANY_QIANCHUAN_REPORT_PRODUCT_LIST_ERROR", "短视频带货获取失败")
)

type QianchuanReportProductRepo interface {
	List(context.Context, uint64, uint64, uint32, string, string, string) (*v1.ListQianchuanReportProductsReply, error)
	Statistics(context.Context, string, string) (*v1.StatisticsQianchuanReportProductsReply, error)
}

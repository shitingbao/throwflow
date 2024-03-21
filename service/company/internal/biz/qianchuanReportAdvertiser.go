package biz

import (
	v1 "company/api/service/douyin/v1"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	CompanyQianchuanReportAdvertiserListError = errors.InternalServer("COMPANY_QIANCHUAN_REPORT_ADVERTISER_LIST_ERROR", "千川广告账户数据获取失败")
)

type QianchuanReportAdvertiserRepo interface {
	List(context.Context, string, string) (*v1.ListQianchuanReportAdvertisersReply, error)
}

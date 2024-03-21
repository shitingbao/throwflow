package biz

import (
	v1 "company/api/service/douyin/v1"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	CompanyQianchuanReportAwemeListError = errors.InternalServer("COMPANY_QIANCHUAN_REPORT_AWEME_LIST_ERROR", "直播带货获取失败")
)

type QianchuanReportAwemeRepo interface {
	List(context.Context, uint64, uint64, uint32, string, string, string) (*v1.ListQianchuanReportAwemesReply, error)
	Statistics(context.Context, string, string) (*v1.StatisticsQianchuanReportAwemesReply, error)
}

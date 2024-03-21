package biz

import (
	v1 "company/api/service/douyin/v1"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	CompanyQianchuanAdListError = errors.InternalServer("COMPANY_QIANCHUAN_AD_LIST_ERROR", "千川广告计划列表获取失败")
	CompanyQianchuanAdNotFound  = errors.NotFound("COMPANY_QIANCHUAN_AD_NOT_FOUND", "千川广告计划不存在")
)

type QianchuanAdRepo interface {
	GetExternal(context.Context, uint64, string) (*v1.GetExternalQianchuanAdsReply, error)
	GetExternalHistory(context.Context, uint64, string, string) (*v1.GetExternalHistoryQianchuanAdsReply, error)
	List(context.Context, uint64, uint64, string, string, string) (*v1.ListQianchuanAdsReply, error)
	ListExternal(context.Context, uint64, uint64, string, string, string, string, string, string, string) (*v1.ListExternalQianchuanAdsReply, error)
	ListSelectExternal(context.Context) (*v1.ListSelectExternalQianchuanAdsReply, error)
	StatisticsExternal(context.Context, string, string, string, string, string) (*v1.StatisticsExternalQianchuanAdsReply, error)
}

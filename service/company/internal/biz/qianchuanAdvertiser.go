package biz

import (
	v1 "company/api/service/douyin/v1"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	CompanyQianchuanAdvertiserNotFound  = errors.NotFound("COMPANY_QIANCHUAN_ADVERTISER_NOT_FOUND", "千川广告账户不存在")
	CompanyQianchuanAdvertiserListError = errors.InternalServer("COMPANY_QIANCHUAN_ADVERTISER_LIST_ERROR", "千川账户列表获取失败")
)

type QianchuanAdvertiserRepo interface {
	GetByCompanyIdAndAdvertiserIds(context.Context, uint64, uint64) (*v1.GetQianchuanAdvertiserByCompanyIdAndAdvertiserIdsReply, error)
	List(context.Context, uint64, uint64, uint64, string, string, string) (*v1.ListQianchuanAdvertisersReply, error)
	ListByDays(context.Context, uint64, string) (*v1.ListQianchuanAdvertiserByDaysReply, error)
	ListExternal(context.Context, uint64, uint64, string, string, string) (*v1.ListExternalQianchuanAdvertisersReply, error)
	StatisticsDashboard(context.Context, uint64, string, string) (*v1.StatisticsDashboardQianchuanAdvertisersReply, error)
	StatisticsExternal(context.Context, string, string, string) (*v1.StatisticsExternalQianchuanAdvertisersReply, error)
	Update(context.Context, uint64, uint64, uint32) (*v1.UpdateStatusQianchuanAdvertisersReply, error)
	UpdateStatusByCompanyId(context.Context, uint64, uint32) (*v1.UpdateStatusQianchuanAdvertisersByCompanyIdReply, error)
}

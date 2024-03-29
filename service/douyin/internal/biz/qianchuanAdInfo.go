package biz

import (
	"context"
	"douyin/internal/domain"
)

type QianchuanAdInfoRepo interface {
	Get(context.Context, uint64, string) (*domain.QianchuanAdInfo, error)
	GetByDay(context.Context, uint64, string, string) (*domain.QianchuanAdInfo, error)
	List(context.Context, string, string, string, string, string, string, string, uint64, uint64) ([]*domain.QianchuanAdInfo, error)
	ListNotLabAd(context.Context, string, string) ([]*domain.QianchuanAdInfo, error)
	ListAdvertiser(context.Context, string) ([]*domain.QianchuanAdvertiserInfo, error)
	ListAdvertiserCampaigns(context.Context, string) ([]*domain.QianchuanAdvertiserInfo, error)
	Count(context.Context, string, string, string, string, string) (uint64, error)
	Statistics(context.Context, string, string, string, string, string) (*domain.QianchuanAdInfo, error)
	SaveIndex(context.Context, string)
	UpsertQianchuanReportAd(context.Context, string, *domain.QianchuanAdInfo) error
	UpsertQianchuanAd(context.Context, string, *domain.QianchuanAdInfo) error
	UpsertQianchuanCampaign(context.Context, string, *domain.QianchuanAdInfo) error
}

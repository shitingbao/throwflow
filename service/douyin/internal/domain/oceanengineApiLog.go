package domain

import (
	"context"
	"time"
)

type OceanengineApiLog struct {
	CompanyId    uint64
	AccountId    uint64
	AppId        string
	AdvertiserId uint64
	CampaignId   uint64
	AdId         uint64
	AccessToken  string
	Content      string
	CreateTime   time.Time
}

func NewOceanengineApiLog(ctx context.Context, companyId, accountId, advertiserId, campaignId, adId uint64, appId, accessToken, content string) *OceanengineApiLog {
	return &OceanengineApiLog{
		CompanyId:    companyId,
		AccountId:    accountId,
		AppId:        appId,
		AdvertiserId: advertiserId,
		CampaignId:   campaignId,
		AdId:         adId,
		AccessToken:  accessToken,
		Content:      content,
	}
}

func (oal *OceanengineApiLog) SetCompanyId(ctx context.Context, companyId uint64) {
	oal.CompanyId = companyId
}

func (oal *OceanengineApiLog) SetAccountId(ctx context.Context, accountId uint64) {
	oal.AccountId = accountId
}

func (oal *OceanengineApiLog) SetAdvertiserId(ctx context.Context, advertiserId uint64) {
	oal.AdvertiserId = advertiserId
}

func (oal *OceanengineApiLog) SetCampaignId(ctx context.Context, campaignId uint64) {
	oal.CampaignId = campaignId
}

func (oal *OceanengineApiLog) SetAdId(ctx context.Context, adId uint64) {
	oal.AdId = adId
}

func (oal *OceanengineApiLog) SetAppId(ctx context.Context, appId string) {
	oal.AppId = appId
}

func (oal *OceanengineApiLog) SetAccessToken(ctx context.Context, accessToken string) {
	oal.AccessToken = accessToken
}

func (oal *OceanengineApiLog) SetContent(ctx context.Context, content string) {
	oal.Content = content
}

func (oal *OceanengineApiLog) SetCreateTime(ctx context.Context) {
	oal.CreateTime = time.Now()
}

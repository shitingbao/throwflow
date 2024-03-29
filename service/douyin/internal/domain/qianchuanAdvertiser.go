package domain

import (
	"context"
	"time"
)

type QianchuanAdvertiser struct {
	Id               uint64
	AppId            string
	CompanyId        uint64
	AccountId        uint64
	AdvertiserId     uint64
	AdvertiserName   string
	CompanyName      string
	Status           uint8
	OtherCompanyName string
	OtherCompanyId   string
	CreateTime       time.Time
	UpdateTime       time.Time
}

type QianchuanAdvertiserList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*QianchuanAdvertiser
}

type ExternalQianchuanAdvertiser struct {
	AdvertiserId            uint64
	AdvertiserName          string
	CompanyName             string
	StatCost                float64
	Roi                     float64
	YesterdayStatCost       float64
	YesterdayRoi            float64
	YesterdayPayOrderAmount float64
	GeneralTotalBalance     float64
	Campaigns               uint64
	PayOrderCount           int64
	PayOrderAmount          float64
	ClickCnt                int64
	ClickRate               float64
	PayConvertRate          float64
	AveragePayOrderStatCost float64
	PayOrderAveragePrice    float64
	DyFollow                int64
	ShowCnt                 int64
}

type ExternalQianchuanAdvertiserList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*ExternalQianchuanAdvertiser
}

type RobbotQianchuanAdvertiser struct {
	CompanyId      uint64
	AccountId      uint64
	AdvertiserId   uint64
	AdvertiserName string
	AccessToken    string
}

type Status struct {
	Key   string
	Value string
}

type SelectQianchuanAdvertisers struct {
	Status []*Status
}

type StatisticsQianchuanAdvertiser struct {
	Key   string
	Value string
}

type StatisticsQianchuanAdvertisers struct {
	Statistics []*StatisticsQianchuanAdvertiser
}

func NewSelectQianchuanAdvertisers() *SelectQianchuanAdvertisers {
	status := make([]*Status, 0)

	status = append(status, &Status{Key: "0", Value: "待授权"})
	status = append(status, &Status{Key: "1", Value: "已授权"})

	return &SelectQianchuanAdvertisers{
		Status: status,
	}
}

func NewQianchuanAdvertiser(ctx context.Context, companyId, accountId, advertiserId uint64, appId, advertiserName, companyName string) *QianchuanAdvertiser {
	return &QianchuanAdvertiser{
		AppId:          appId,
		CompanyId:      companyId,
		AccountId:      accountId,
		AdvertiserId:   advertiserId,
		AdvertiserName: advertiserName,
		CompanyName:    companyName,
	}
}

func (qa *QianchuanAdvertiser) SetAppId(ctx context.Context, appId string) {
	qa.AppId = appId
}

func (qa *QianchuanAdvertiser) SetAdvertiserName(ctx context.Context, advertiserName string) {
	qa.AdvertiserName = advertiserName
}

func (qa *QianchuanAdvertiser) SetCompanyName(ctx context.Context, companyName string) {
	qa.CompanyName = companyName
}

func (qa *QianchuanAdvertiser) SetAccountId(ctx context.Context, accountId uint64) {
	qa.AccountId = accountId
}

func (qa *QianchuanAdvertiser) SetStatus(ctx context.Context, status uint8) {
	qa.Status = status
}

func (qa *QianchuanAdvertiser) SetUpdateTime(ctx context.Context) {
	qa.UpdateTime = time.Now()
}

func (qa *QianchuanAdvertiser) SetCreateTime(ctx context.Context) {
	qa.CreateTime = time.Now()
}

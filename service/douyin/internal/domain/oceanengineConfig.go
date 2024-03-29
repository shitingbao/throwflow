package domain

import (
	"context"
	"time"
)

type OceanengineConfig struct {
	Id              uint64
	OceanengineType uint8
	AppId           string
	AppName         string
	AppSecret       string
	RedirectUrl     string
	Concurrents     uint8
	Status          uint8
	CreateTime      time.Time
	UpdateTime      time.Time
}

type OceanengineConfigList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*OceanengineConfig
}

type OceanengineType struct {
	Key   string
	Value string
}

type SelectOceanengineConfigs struct {
	OceanengineType []*OceanengineType
}

func NewSelectOceanengineConfigs() *SelectOceanengineConfigs {
	oceanengineType := make([]*OceanengineType, 0)

	oceanengineType = append(oceanengineType, &OceanengineType{Key: "1", Value: "巨量千川"})
	oceanengineType = append(oceanengineType, &OceanengineType{Key: "2", Value: "巨量广告"})

	return &SelectOceanengineConfigs{
		OceanengineType: oceanengineType,
	}
}

func NewOceanengineConfig(ctx context.Context, appId, appName, appSecret, redirectUrl string, oceanengineType, concurrents, status uint8) *OceanengineConfig {
	return &OceanengineConfig{
		OceanengineType: oceanengineType,
		AppId:           appId,
		AppName:         appName,
		AppSecret:       appSecret,
		RedirectUrl:     redirectUrl,
		Concurrents:     concurrents,
		Status:          status,
	}
}

func (oc *OceanengineConfig) SetOceanengineType(ctx context.Context, oceanengineType uint8) {
	oc.OceanengineType = oceanengineType
}

func (oc *OceanengineConfig) SetAppId(ctx context.Context, appId string) {
	oc.AppId = appId
}

func (oc *OceanengineConfig) SetAppName(ctx context.Context, appName string) {
	oc.AppName = appName
}

func (oc *OceanengineConfig) SetAppSecret(ctx context.Context, appSecret string) {
	oc.AppSecret = appSecret
}

func (oc *OceanengineConfig) SetRedirectUrl(ctx context.Context, redirectUrl string) {
	oc.RedirectUrl = redirectUrl
}

func (oc *OceanengineConfig) SetConcurrents(ctx context.Context, concurrents uint8) {
	oc.Concurrents = concurrents
}

func (oc *OceanengineConfig) SetStatus(ctx context.Context, status uint8) {
	oc.Status = status
}

func (oc *OceanengineConfig) SetUpdateTime(ctx context.Context) {
	oc.UpdateTime = time.Now()
}

func (oc *OceanengineConfig) SetCreateTime(ctx context.Context) {
	oc.CreateTime = time.Now()
}

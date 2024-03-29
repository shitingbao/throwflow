package domain

import (
	"context"
	"time"
)

type OpenDouyinUserInfoCreateLog struct {
	Id         uint64
	ClientKey  string
	OpenId     string
	IsHandle   uint8
	CreateTime time.Time
	UpdateTime time.Time
}

func NewOpenDouyinUserInfoCreateLog(ctx context.Context, clientKey, openId string) *OpenDouyinUserInfoCreateLog {
	return &OpenDouyinUserInfoCreateLog{
		ClientKey: clientKey,
		OpenId:    openId,
	}
}

func (oduicl *OpenDouyinUserInfoCreateLog) SetClientKey(ctx context.Context, clientKey string) {
	oduicl.ClientKey = clientKey
}

func (oduicl *OpenDouyinUserInfoCreateLog) SetOpenId(ctx context.Context, openId string) {
	oduicl.OpenId = openId
}

func (oduicl *OpenDouyinUserInfoCreateLog) SetIsHandle(ctx context.Context, isHandle uint8) {
	oduicl.IsHandle = isHandle
}

func (oduicl *OpenDouyinUserInfoCreateLog) SetUpdateTime(ctx context.Context) {
	oduicl.UpdateTime = time.Now()
}

func (oduicl *OpenDouyinUserInfoCreateLog) SetCreateTime(ctx context.Context) {
	oduicl.CreateTime = time.Now()
}

package domain

import (
	"context"
	"time"
)

type OpenDouyinApiLog struct {
	ClientKey   string
	OpenId      string
	AccessToken string
	Content     string
	CreateTime  time.Time
}

func NewOpenDouyinApiLog(ctx context.Context, clientKey, openId, accessToken, content string) *OpenDouyinApiLog {
	return &OpenDouyinApiLog{
		ClientKey:   clientKey,
		OpenId:      openId,
		AccessToken: accessToken,
		Content:     content,
	}
}

func (odal *OpenDouyinApiLog) SetClientKey(ctx context.Context, clientKey string) {
	odal.ClientKey = clientKey
}

func (odal *OpenDouyinApiLog) SetOpenId(ctx context.Context, openId string) {
	odal.OpenId = openId
}

func (odal *OpenDouyinApiLog) SetAccessToken(ctx context.Context, accessToken string) {
	odal.AccessToken = accessToken
}

func (odal *OpenDouyinApiLog) SetContent(ctx context.Context, content string) {
	odal.Content = content
}

func (odal *OpenDouyinApiLog) SetCreateTime(ctx context.Context) {
	odal.CreateTime = time.Now()
}

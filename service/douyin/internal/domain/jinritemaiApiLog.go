package domain

import (
	"context"
	"time"
)

type JinritemaiApiLog struct {
	ClientKey   string
	OpenId      string
	AccessToken string
	Content     string
	CreateTime  time.Time
}

func NewJinritemaiApiLog(ctx context.Context, clientKey, openId, accessToken, content string) *JinritemaiApiLog {
	return &JinritemaiApiLog{
		ClientKey:   clientKey,
		OpenId:      openId,
		AccessToken: accessToken,
		Content:     content,
	}
}

func (jal *JinritemaiApiLog) SetClientKey(ctx context.Context, clientKey string) {
	jal.ClientKey = clientKey
}

func (jal *JinritemaiApiLog) SetOpenId(ctx context.Context, openId string) {
	jal.OpenId = openId
}

func (jal *JinritemaiApiLog) SetAccessToken(ctx context.Context, accessToken string) {
	jal.AccessToken = accessToken
}

func (jal *JinritemaiApiLog) SetContent(ctx context.Context, content string) {
	jal.Content = content
}

func (jal *JinritemaiApiLog) SetCreateTime(ctx context.Context) {
	jal.CreateTime = time.Now()
}

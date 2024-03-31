package domain

import (
	"context"
)

type AwemesAdvertiserWeixinAuth struct {
	ClientKey       string
	OpenId          string
	CooperativeCode string
	AuthStatus      int32
}

func NewAwemesAdvertiserWeixinAuth(ctx context.Context, authStatus int32, clientKey, openId, cooperativeCode string) *AwemesAdvertiserWeixinAuth {
	return &AwemesAdvertiserWeixinAuth{
		ClientKey:       clientKey,
		OpenId:          openId,
		CooperativeCode: cooperativeCode,
		AuthStatus:      authStatus,
	}
}

func (aawa *AwemesAdvertiserWeixinAuth) SetClientKey(ctx context.Context, clientKey string) {
	aawa.ClientKey = clientKey
}

func (aawa *AwemesAdvertiserWeixinAuth) SetOpenId(ctx context.Context, openId string) {
	aawa.OpenId = openId
}

func (aawa *AwemesAdvertiserWeixinAuth) SetCooperativeCode(ctx context.Context, cooperativeCode string) {
	aawa.CooperativeCode = cooperativeCode
}

func (aawa *AwemesAdvertiserWeixinAuth) SetAuthStatus(ctx context.Context, authStatus int32) {
	aawa.AuthStatus = authStatus
}

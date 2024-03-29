package domain

import (
	"context"
	"time"
)

type DoukeToken struct {
	Id              uint64
	AuthorityId     string
	AuthSubjectType string
	AccessToken     string
	ExpiresIn       uint64
	RefreshToken    string
	CreateTime      time.Time
	UpdateTime      time.Time
}

func NewDoukeToken(ctx context.Context, authorityId, authSubjectType, accessToken, refreshToken string, expiresIn uint64) *DoukeToken {
	return &DoukeToken{
		AuthorityId:     authorityId,
		AuthSubjectType: authSubjectType,
		AccessToken:     accessToken,
		ExpiresIn:       expiresIn,
		RefreshToken:    refreshToken,
	}
}

func (dt *DoukeToken) SetAuthorityId(ctx context.Context, authorityId string) {
	dt.AuthorityId = authorityId
}

func (dt *DoukeToken) SetAuthSubjectType(ctx context.Context, authSubjectType string) {
	dt.AuthSubjectType = authSubjectType
}

func (dt *DoukeToken) SetAccessToken(ctx context.Context, accessToken string) {
	dt.AccessToken = accessToken
}

func (dt *DoukeToken) SetExpiresIn(ctx context.Context, expiresIn uint64) {
	dt.ExpiresIn = expiresIn
}

func (dt *DoukeToken) SetRefreshToken(ctx context.Context, refreshToken string) {
	dt.RefreshToken = refreshToken
}

func (dt *DoukeToken) SetCreateTime(ctx context.Context) {
	dt.CreateTime = time.Now()
}

func (dt *DoukeToken) SetUpdateTime(ctx context.Context) {
	dt.UpdateTime = time.Now()
}

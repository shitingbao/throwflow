package domain

import (
	"context"
	"time"
)

type OceanengineAccountToken struct {
	Id                    uint64
	AppId                 string
	CompanyId             uint64
	AccountId             uint64
	AccessToken           string
	ExpiresIn             uint32
	RefreshToken          string
	RefreshTokenExpiresIn uint32
	AccountRole           string
	CreateTime            time.Time
	UpdateTime            time.Time
}

func NewOceanengineAccountToken(ctx context.Context, appId, accessToken, refreshToken string, companyId, accountId uint64, expiresIn, refreshTokenExpiresIn uint32) *OceanengineAccountToken {
	return &OceanengineAccountToken{
		AppId:                 appId,
		CompanyId:             companyId,
		AccountId:             accountId,
		AccessToken:           accessToken,
		ExpiresIn:             expiresIn,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresIn: refreshTokenExpiresIn,
	}
}

func (oat *OceanengineAccountToken) SetAppId(ctx context.Context, appId string) {
	oat.AppId = appId
}

func (oat *OceanengineAccountToken) SetAccessToken(ctx context.Context, accessToken string) {
	oat.AccessToken = accessToken
}

func (oat *OceanengineAccountToken) SetRefreshToken(ctx context.Context, refreshToken string) {
	oat.RefreshToken = refreshToken
}

func (oat *OceanengineAccountToken) SetCompanyId(ctx context.Context, companyId uint64) {
	oat.CompanyId = companyId
}

func (oat *OceanengineAccountToken) SetExpiresIn(ctx context.Context, expiresIn uint32) {
	oat.ExpiresIn = expiresIn
}

func (oat *OceanengineAccountToken) SetRefreshTokenExpiresIn(ctx context.Context, refreshTokenExpiresIn uint32) {
	oat.RefreshTokenExpiresIn = refreshTokenExpiresIn
}

func (oat *OceanengineAccountToken) SetUpdateTime(ctx context.Context) {
	oat.UpdateTime = time.Now()
}

func (oat *OceanengineAccountToken) SetCreateTime(ctx context.Context) {
	oat.CreateTime = time.Now()
}

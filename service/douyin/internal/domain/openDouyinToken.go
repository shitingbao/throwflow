package domain

import (
	"context"
	"time"
)

type OpenDouyinToken struct {
	Id               uint64
	ClientKey        string `json:"clientKey"`
	OpenId           string `json:"openId"`
	AccessToken      string
	ExpiresIn        uint64
	RefreshToken     string
	RefreshExpiresIn uint64
	CreateTime       time.Time
	UpdateTime       time.Time
}

type ConfigOpenDouyinToken struct {
	ClientKey   string
	CallBackUrl string
	Scope       string
}

type QrCodeOpenDouyinTokenData struct {
	Captcha        string `json:"captcha"`
	DescUrl        string `json:"desc_url"`
	Description    string `json:"description"`
	ErrorCode      uint64 `json:"error_code"`
	IsFrontier     bool   `json:"is_frontier"`
	Qrcode         string `json:"qrcode"`
	QrcodeIndexUrl string `json:"qrcode_index_url"`
	Token          string `json:"token"`
}

type QrCodeOpenDouyinToken struct {
	Message string                    `json:"message"`
	Data    QrCodeOpenDouyinTokenData `json:"data"`
}

type QrCodeOpenDouyinStatusData struct {
	Captcha         string `json:"captcha"`
	Code            string `json:"code"`
	ConfirmedScopes string `json:"confirmed_scopes"`
	DescUrl         string `json:"desc_url"`
	Description     string `json:"description"`
	ErrorCode       uint64 `json:"error_code"`
	RedirectUrl     string `json:"redirect_url"`
	State           string `json:"state"`
	Status          string `json:"status"`
}

type QrCodeOpenDouyinStatus struct {
	Message string                     `json:"message"`
	Data    QrCodeOpenDouyinStatusData `json:"data"`
}

type QrCode struct {
	QrCode string
	State  string
	Token  string
}

type QrCodeStatus struct {
	Code   string
	Status string
}

type StateParam struct {
	UserId    string `json:"userId"`
	ClientKey string `json:"clientKey"`
}

func NewOpenDouyinToken(ctx context.Context, clientKey, openId, accessToken, refreshToken string, expiresIn, refreshExpiresIn uint64) *OpenDouyinToken {
	return &OpenDouyinToken{
		ClientKey:        clientKey,
		OpenId:           openId,
		AccessToken:      accessToken,
		ExpiresIn:        expiresIn,
		RefreshToken:     refreshToken,
		RefreshExpiresIn: refreshExpiresIn,
	}
}

func (odt *OpenDouyinToken) SetClientKey(ctx context.Context, clientKey string) {
	odt.ClientKey = clientKey
}

func (odt *OpenDouyinToken) SetOpenId(ctx context.Context, openId string) {
	odt.OpenId = openId
}

func (odt *OpenDouyinToken) SetAccessToken(ctx context.Context, accessToken string) {
	odt.AccessToken = accessToken
}

func (odt *OpenDouyinToken) SetExpiresIn(ctx context.Context, expiresIn uint64) {
	odt.ExpiresIn = expiresIn
}

func (odt *OpenDouyinToken) SetRefreshToken(ctx context.Context, refreshToken string) {
	odt.RefreshToken = refreshToken
}

func (odt *OpenDouyinToken) SetRefreshExpiresIn(ctx context.Context, refreshExpiresIn uint64) {
	odt.RefreshExpiresIn = refreshExpiresIn
}

func (odt *OpenDouyinToken) SetUpdateTime(ctx context.Context) {
	odt.UpdateTime = time.Now()
}

func (odt *OpenDouyinToken) SetCreateTime(ctx context.Context) {
	odt.CreateTime = time.Now()
}

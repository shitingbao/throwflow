package token

import (
	"douyin/internal/pkg/douke"
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
)

func CreateToken(appKey, appSecret, code string) (*CreateTokenResponse, error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	accessTokenRequest := CreateTokenRequest{
		Code:            code,
		GrantType:       douke.CreateTokenGrantType,
		TestShop:        "",
		ShopId:          "",
		AuthId:          "",
		AuthSubjectType: "",
	}

	param := accessTokenRequest.String()

	commonRequest := douke.CommonRequest{
		Method:     "token.create",
		AppKey:     appKey,
		ParamJson:  param,
		Timestamp:  timestamp,
		V:          douke.V,
		Sign:       douke.Sign(appKey, appSecret, "token.create", param, timestamp, douke.V),
		SignMethod: douke.SignMethod,
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(commonRequest.String()).
		SetHeader("Content-Type", douke.ApplicationJson).
		SetBody(accessTokenRequest.String()).
		Post(douke.BaseDomain + "/token/create")

	if err != nil {
		return nil, douke.NewDoukeError(90101, douke.BaseDomain+"/token/create", "请求失败："+err.Error(), "")
	}

	if resp.StatusCode() != 200 {
		return nil, douke.NewDoukeError(90101, douke.BaseDomain+"/token/create", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var createTokenResponse CreateTokenResponse

	if err := createTokenResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &createTokenResponse, nil
}

func RefreshToken(appKey, appSecret, refreshToken string) (*RefreshTokenResponse, error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	refreshTokenRequest := RefreshTokenRequest{
		RefreshToken: refreshToken,
		GrantType:    douke.RefreshTokenGrantType,
	}

	param := refreshTokenRequest.String()

	commonRequest := douke.CommonRequest{
		Method:     "token.refresh",
		AppKey:     appKey,
		ParamJson:  param,
		Timestamp:  timestamp,
		V:          douke.V,
		Sign:       douke.Sign(appKey, appSecret, "token.refresh", param, timestamp, douke.V),
		SignMethod: douke.SignMethod,
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(commonRequest.String()).
		SetHeader("Content-Type", douke.ApplicationJson).
		SetBody(refreshTokenRequest.String()).
		Post(douke.BaseDomain + "/token/refresh")

	if err != nil {
		return nil, douke.NewDoukeError(90101, douke.BaseDomain+"/token/refresh", "请求失败："+err.Error(), resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, douke.NewDoukeError(90101, douke.BaseDomain+"/token/refresh", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var refreshTokenResponse RefreshTokenResponse

	if err := refreshTokenResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &refreshTokenResponse, nil
}

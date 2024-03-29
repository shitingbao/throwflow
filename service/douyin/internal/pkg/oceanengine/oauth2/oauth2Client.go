package oauth2

import (
	"douyin/internal/pkg/oceanengine"
	"github.com/go-resty/resty/v2"
	"strconv"
)

func AccessToken(appId, appSecret, authCode string) (*AccessTokenResponse, error) {
	accessTokenRequest := AccessTokenRequest{
		AppId:     appId,
		Secret:    appSecret,
		GrantType: oceanengine.AccessTokenGrantType,
		AuthCode:  authCode,
	}

	resp, err := resty.
		New().
		R().
		SetHeader("Content-Type", oceanengine.ApplicationJson).
		SetBody(accessTokenRequest.String()).
		Post(oceanengine.BaseDomain + "/open_api/oauth2/access_token/")

	if err != nil {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/oauth2/access_token/", "REQUEST_ERROR", "请求失败："+err.Error(), "", resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/oauth2/access_token/", "REQUEST_ERROR", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), "", resp.String())
	}

	var accessTokenResponse AccessTokenResponse

	if err := accessTokenResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &accessTokenResponse, nil
}

func RefreshToken(appId, appSecret, refreshToken string) (*RefreshTokenResponse, error) {
	refreshTokenRequest := RefreshTokenRequest{
		AppId:        appId,
		Secret:       appSecret,
		GrantType:    oceanengine.RefreshTokenGrantType,
		RefreshToken: refreshToken,
	}

	resp, err := resty.
		New().
		R().
		SetHeader("Content-Type", oceanengine.ApplicationJson).
		SetBody(refreshTokenRequest.String()).
		Post(oceanengine.BaseDomain + "/open_api/oauth2/refresh_token/")

	if err != nil {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/oauth2/refresh_token/", "REQUEST_ERROR", "请求失败："+err.Error(), "", resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/oauth2/refresh_token/", "REQUEST_ERROR", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), "", resp.String())
	}

	var refreshTokenResponse RefreshTokenResponse

	if err := refreshTokenResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &refreshTokenResponse, nil
}

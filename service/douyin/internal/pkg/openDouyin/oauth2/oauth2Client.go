package oauth2

import (
	"douyin/internal/pkg/openDouyin"
	"github.com/go-resty/resty/v2"
	"strconv"
)

func AccessToken(clientKey, clientSecret, code string) (*AccessTokenResponse, error) {
	accessTokenRequest := AccessTokenRequest{
		ClientSecret: clientSecret,
		Code:         code,
		GrantType:    openDouyin.AccessTokenGrantType,
		ClientKey:    clientKey,
	}

	resp, err := resty.
		New().
		R().
		SetHeader("Content-Type", openDouyin.ApplicationJson).
		SetBody(accessTokenRequest.String()).
		Post(openDouyin.BaseDomain + "/oauth/access_token/")

	if err != nil {
		return nil, openDouyin.NewOpenDouyinError(90001, openDouyin.BaseDomain+"/oauth/access_token/", "请求失败："+err.Error(), resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, openDouyin.NewOpenDouyinError(90001, openDouyin.BaseDomain+"/oauth/access_token/", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var accessTokenResponse AccessTokenResponse

	if err := accessTokenResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &accessTokenResponse, nil
}

func RefreshToken(clientKey, refreshToken string) (*RefreshTokenResponse, error) {
	refreshTokenRequest := make(map[string]string)
	refreshTokenRequest["client_key"] = clientKey
	refreshTokenRequest["grant_type"] = openDouyin.RefreshTokenGrantType
	refreshTokenRequest["refresh_token"] = refreshToken

	resp, err := resty.
		New().
		R().
		SetHeader("Content-Type", openDouyin.ApplicationForm).
		SetMultipartFormData(refreshTokenRequest).
		Post(openDouyin.BaseDomain + "/oauth/refresh_token/")

	if err != nil {
		return nil, openDouyin.NewOpenDouyinError(90001, openDouyin.BaseDomain+"/oauth/refresh_token/", "请求失败："+err.Error(), resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, openDouyin.NewOpenDouyinError(90001, openDouyin.BaseDomain+"/oauth/refresh_token/", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var refreshTokenResponse RefreshTokenResponse

	if err := refreshTokenResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &refreshTokenResponse, nil
}

func RenewRefreshToken(clientKey, refreshToken string) (*RenewRefreshTokenResponse, error) {
	renewRefreshTokenRequest := make(map[string]string)
	renewRefreshTokenRequest["client_key"] = clientKey
	renewRefreshTokenRequest["refresh_token"] = refreshToken

	resp, err := resty.
		New().
		R().
		SetHeader("Content-Type", openDouyin.ApplicationForm).
		SetMultipartFormData(renewRefreshTokenRequest).
		Post(openDouyin.BaseDomain + "/oauth/renew_refresh_token/")

	if err != nil {
		return nil, openDouyin.NewOpenDouyinError(90001, openDouyin.BaseDomain+"/oauth/renew_refresh_token/", "请求失败："+err.Error(), resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, openDouyin.NewOpenDouyinError(90001, openDouyin.BaseDomain+"/oauth/renew_refresh_token/", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var renewRefreshTokenResponse RenewRefreshTokenResponse

	if err := renewRefreshTokenResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &renewRefreshTokenResponse, nil
}

func ClientToken(clientKey, clientSecret string) (*ClientTokenResponse, error) {
	clientTokenRequest := ClientTokenRequest{
		ClientKey:    clientKey,
		ClientSecret: clientSecret,
		GrantType:    openDouyin.ClientTokenGrantType,
	}

	resp, err := resty.
		New().
		R().
		SetHeader("Content-Type", openDouyin.ApplicationJson).
		SetBody(clientTokenRequest.String()).
		Post(openDouyin.BaseDomain + "/oauth/client_token/")

	if err != nil {
		return nil, openDouyin.NewOpenDouyinError(90001, openDouyin.BaseDomain+"/oauth/client_token/", "请求失败："+err.Error(), resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, openDouyin.NewOpenDouyinError(90001, openDouyin.BaseDomain+"/oauth/client_token/", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var clientTokenResponse ClientTokenResponse

	if err := clientTokenResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &clientTokenResponse, nil
}

func GetUserInfo(accessToken, openId string) (*GetUserInfoResponse, error) {
	getUserInfoRequest := make(map[string]string)
	getUserInfoRequest["access_token"] = accessToken
	getUserInfoRequest["open_id"] = openId

	resp, err := resty.
		New().
		R().
		SetHeader("Content-Type", openDouyin.ApplicationForm).
		SetMultipartFormData(getUserInfoRequest).
		Post(openDouyin.BaseDomain + "/oauth/userinfo/")

	if err != nil {
		return nil, openDouyin.NewOpenDouyinError(90001, openDouyin.BaseDomain+"/oauth/userinfo/", "请求失败："+err.Error(), resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, openDouyin.NewOpenDouyinError(90001, openDouyin.BaseDomain+"/oauth/userinfo/", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var getUserInfoResponse GetUserInfoResponse

	if err := getUserInfoResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &getUserInfoResponse, nil
}

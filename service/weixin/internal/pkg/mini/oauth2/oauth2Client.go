package oauth2

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"weixin/internal/pkg/mini"
)

func GetAccessToken(appid, secret string) (*GetAccessTokenResponse, error) {
	getAccessTokenRequest := GetAccessTokenRequest{
		GrantType: mini.GetAccessTokenGrantType,
		Appid:     appid,
		Secret:    secret,
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(getAccessTokenRequest.String()).
		Get(mini.BaseDomain + "/cgi-bin/token")

	if err != nil {
		return nil, mini.NewMiniError(err.Error())
	}
	fmt.Println(resp.String())
	var getAccessTokenResponse GetAccessTokenResponse

	if err := getAccessTokenResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &getAccessTokenResponse, nil
}

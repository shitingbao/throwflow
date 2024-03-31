package wxa

import (
	"github.com/go-resty/resty/v2"
	"weixin/internal/pkg/mini"
)

func GetUnlimitedQRCode(accessToken, scene, envVersion string) (*GetUnlimitedQRCodeResponse, error) {
	getUnlimitedQRCodeQueryRequest := GetUnlimitedQRCodeQueryRequest{
		AccessToken: accessToken,
	}

	getUnlimitedQRCodeRequest := GetUnlimitedQRCodeRequest{
		Scene:      scene,
		Page:       "pages/user/user",
		EnvVersion: envVersion,
		CheckPath:  false,
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(getUnlimitedQRCodeQueryRequest.String()).
		SetBody(getUnlimitedQRCodeRequest.String()).
		Post(mini.BaseDomain + "/wxa/getwxacodeunlimit")

	if err != nil {
		return nil, mini.NewMiniError(err.Error())
	}

	var getUnlimitedQRCodeResponse GetUnlimitedQRCodeResponse

	if err := getUnlimitedQRCodeResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &getUnlimitedQRCodeResponse, nil
}

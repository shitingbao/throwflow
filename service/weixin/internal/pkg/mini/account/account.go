package account

import (
	"github.com/go-resty/resty/v2"
	"weixin/internal/pkg/mini"
)

func Code2Session(appid, secret, code string) (*Code2SessionResponse, error) {
	code2SessionRequest := Code2SessionRequest{
		Appid:     appid,
		Secret:    secret,
		JsCode:    code,
		GrantType: mini.Code2SessionTokenGrantType,
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(code2SessionRequest.String()).
		Get(mini.BaseDomain + "/sns/jscode2session")

	if err != nil {
		return nil, mini.NewMiniError(err.Error())
	}

	var code2SessionResponse Code2SessionResponse

	if err := code2SessionResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &code2SessionResponse, nil
}

func GetPhoneNumber(accessToken, code string) (*GetPhoneNumberResponse, error) {
	getPhoneNumberQueryRequest := GetPhoneNumberQueryRequest{
		AccessToken: accessToken,
	}

	getPhoneNumberRequest := GetPhoneNumberRequest{
		Code: code,
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(getPhoneNumberQueryRequest.String()).
		SetBody(getPhoneNumberRequest.String()).
		Post(mini.BaseDomain + "/wxa/business/getuserphonenumber")

	if err != nil {
		return nil, mini.NewMiniError(err.Error())
	}

	var getPhoneNumberResponse GetPhoneNumberResponse

	if err := getPhoneNumberResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &getPhoneNumberResponse, nil
}

package account

import (
	"encoding/json"
	"weixin/internal/pkg/mini"
)

type Watermark struct {
	Timestamp uint64 `json:"timestamp"`
	Appid     string `json:"appid"`
}

type PhoneInfo struct {
	PhoneNumber     string    `json:"phoneNumber"`
	PurePhoneNumber string    `json:"purePhoneNumber"`
	CountryCode     string    `json:"countryCode"`
	Watermark       Watermark `json:"watermark"`
}

type GetPhoneNumberResponse struct {
	Errmsg    string    `json:"errmsg"`
	Errcode   int32     `json:"errcode"`
	PhoneInfo PhoneInfo `json:"phone_info"`
}

func (gpnr *GetPhoneNumberResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), gpnr); err != nil {
		return mini.NewMiniError("解析json失败：" + err.Error())
	} else {
		if gpnr.Errcode != 0 {
			return mini.NewMiniError(mini.ResponseDescription[gpnr.Errcode])
		}
	}

	return nil
}

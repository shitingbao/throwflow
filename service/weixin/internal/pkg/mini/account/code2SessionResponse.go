package account

import (
	"encoding/json"
	"weixin/internal/pkg/mini"
)

type Code2SessionResponse struct {
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	Errmsg     string `json:"errmsg"`
	Openid     string `json:"openid"`
	Errcode    int32  `json:"errcode"`
}

func (csr *Code2SessionResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), csr); err != nil {
		return mini.NewMiniError("解析json失败：" + err.Error())
	} else {
		if csr.Errcode != 0 {
			return mini.NewMiniError(mini.ResponseDescription[csr.Errcode])
		}
	}

	return nil
}

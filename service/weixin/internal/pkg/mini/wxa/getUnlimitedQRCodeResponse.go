package wxa

import (
	"encoding/json"
	"strings"
	"weixin/internal/pkg/mini"
)

type GetUnlimitedQRCodeResponse struct {
	Buffer  string `json:"buffer"`
	Errmsg  string `json:"errmsg"`
	Errcode int32  `json:"errcode"`
}

func (guqrcr *GetUnlimitedQRCodeResponse) Struct(response string) error {
	if strings.HasPrefix(response, "{") {
		if err := json.Unmarshal([]byte(response), guqrcr); err != nil {
			return mini.NewMiniError("解析json失败：" + err.Error())
		} else {
			if guqrcr.Errcode != 0 {
				return mini.NewMiniError(mini.ResponseDescription[guqrcr.Errcode])
			}
		}
	} else {
		guqrcr.Buffer = response
	}

	return nil
}

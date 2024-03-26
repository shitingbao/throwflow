package kuaidi

import (
	"encoding/json"
)

type GetKuaidiParamRequest struct {
	Com      string `json:"com"`
	Num      string `json:"num"`
	Phone    string `json:"phone"`
	From     string `json:"from"`
	To       string `json:"to"`
	Resultv2 string `json:"resultv2"`
	Show     string `json:"show"`
	Order    string `json:"order"`
}

func (gkpr GetKuaidiParamRequest) String() string {
	data, _ := json.Marshal(gkpr)

	return string(data)
}

package wxa

import (
	"github.com/google/go-querystring/query"
	"weixin/internal/pkg/tool"
)

type GetUnlimitedQRCodeQueryRequest struct {
	AccessToken string `url:"access_token"`
}

type GetUnlimitedQRCodeRequest struct {
	Scene      string `json:"scene"`
	Page       string `json:"page"`
	EnvVersion string `json:"env_version"`
	CheckPath  bool   `json:"check_path"`
}

func (guqrcqr GetUnlimitedQRCodeQueryRequest) String() string {
	v, _ := query.Values(guqrcqr)

	return v.Encode()
}

func (guqrcr GetUnlimitedQRCodeRequest) String() string {
	data, _ := tool.Marshal(guqrcr)

	return string(data)
}

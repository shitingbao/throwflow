package product

import (
	"douyin/internal/pkg/csj"
	"douyin/internal/pkg/douke"
	"encoding/json"
)

type LinkDataResponse struct {
	DyPassword string `json:"dy_password"`
}

type LinkResponse struct {
	csj.CommonResponse
	Data LinkDataResponse `json:"data"`
}

func (lr *LinkResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), lr); err != nil {
		return csj.NewCsjError(90100, douke.BaseDomain+"/product/link", "解析json失败："+err.Error(), response)
	} else {
		if lr.CommonResponse.Code != 0 {
			return csj.NewCsjError(lr.CommonResponse.Code, douke.BaseDomain+"/product/link", lr.CommonResponse.Desc, response)
		}
	}

	return nil
}

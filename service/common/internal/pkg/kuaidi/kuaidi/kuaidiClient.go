package kuaidi

import (
	"common/internal/pkg/kuaidi"
	"common/internal/pkg/tool"
	"github.com/go-resty/resty/v2"
	"strconv"
	"strings"
)

func GetKuaidi(key, customer, com, num, phone string) (*GetKuaidiResponse, error) {
	getKuaidiParamRequest := GetKuaidiParamRequest{
		Com:      com,
		Num:      num,
		Phone:    phone,
		From:     "",
		To:       "",
		Resultv2: "1",
		Show:     "0",
		Order:    "desc",
	}

	getKuaidiRequest := make(map[string]string)
	getKuaidiRequest["customer"] = customer
	getKuaidiRequest["sign"] = strings.ToUpper(tool.GetMd5(getKuaidiParamRequest.String() + key + customer))
	getKuaidiRequest["param"] = getKuaidiParamRequest.String()

	resp, err := resty.
		New().
		R().
		SetHeader("Content-Type", kuaidi.ApplicationForm).
		SetFormData(getKuaidiRequest).
		Post(kuaidi.BaseDomain + "/poll/query.do")

	if err != nil {
		return nil, kuaidi.NewKuaidiError(kuaidi.BaseDomain+"/poll/query.do", "请求失败："+err.Error(), resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, kuaidi.NewKuaidiError(kuaidi.BaseDomain+"/poll/query.do", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var getKuaidiResponse GetKuaidiResponse

	if err := getKuaidiResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &getKuaidiResponse, nil
}

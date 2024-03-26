package suolink

import (
	"common/internal/pkg/suolink"
	"github.com/go-resty/resty/v2"
	"strconv"
)

func GetSuoLink(key, url string) (*GetSuoLinkResponse, error) {
	getSuoLinkRequest := GetSuoLinkRequest{
		Url:        url,
		Format:     suolink.Format,
		Key:        key,
		ExpireDate: suolink.ExpireDate,
		Domain:     suolink.Domain,
		Protocol:   "0",
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(getSuoLinkRequest.String()).
		Get(suolink.BaseDomain)

	if err != nil {
		return nil, suolink.NewSuoLinkError(suolink.BaseDomain, "请求失败："+err.Error(), resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, suolink.NewSuoLinkError(suolink.BaseDomain, "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var getSuoLinkResponse GetSuoLinkResponse

	if err := getSuoLinkResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &getSuoLinkResponse, nil
}

package tools

import (
	"douyin/internal/pkg/jinritemai"
	"douyin/internal/pkg/oceanengine"
	"github.com/go-resty/resty/v2"
	"strconv"
)

func AddAwemeAuth(advertiserId uint64, accessToken, awemeId, code string) (*AddAwemeAuthResponse, error) {
	addAwemeAuthRequest := AddAwemeAuthRequest{
		AdvertiserId: advertiserId,
		AwemeId:      awemeId,
		Code:         code,
		AuthType:     oceanengine.AwemeAuthAuthType,
		EndTime:      oceanengine.AwemeAuthEndTime,
	}

	resp, err := resty.
		New().
		R().
		SetHeader("Access-Token", accessToken).
		SetHeader("Content-Type", jinritemai.ApplicationJson).
		SetBody(addAwemeAuthRequest.String()).
		Post(oceanengine.AwemeAuthDomain + "/open_api/v1.0/qianchuan/tools/aweme_auth/")

	if err != nil {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.AwemeAuthDomain+"/open_api/v1.0/qianchuan/tools/aweme_auth/", "REQUEST_ERROR", "请求失败："+err.Error(), "", resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.AwemeAuthDomain+"/open_api/v1.0/qianchuan/tools/aweme_auth/", "REQUEST_ERROR", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), "", resp.String())
	}

	var addAwemeAuthResponse AddAwemeAuthResponse

	if err := addAwemeAuthResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &addAwemeAuthResponse, nil
}

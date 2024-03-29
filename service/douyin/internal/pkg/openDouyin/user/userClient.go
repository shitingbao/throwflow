package user

import (
	"douyin/internal/pkg/openDouyin"
	"github.com/go-resty/resty/v2"
	"strconv"
)

func ListUserFans(accessToken, openId string) (*ListUserFansResponse, error) {
	listUserFansRequest := ListUserFansRequest{
		OpenId:   openId,
		DateType: openDouyin.DataType,
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(listUserFansRequest.String()).
		SetHeader("Access-Token", accessToken).
		SetHeader("Content-Type", openDouyin.ApplicationJson).
		Get(openDouyin.BaseDomain + "/data/external/user/fans/")

	if err != nil {
		return nil, openDouyin.NewOpenDouyinError(90001, openDouyin.BaseDomain+"/data/external/user/fans/", "请求失败："+err.Error(), resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, openDouyin.NewOpenDouyinError(90001, openDouyin.BaseDomain+"/data/external/user/fans/", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var listUserFansResponse ListUserFansResponse

	if err := listUserFansResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &listUserFansResponse, nil
}

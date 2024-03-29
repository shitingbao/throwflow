package js

import (
	"douyin/internal/pkg/openDouyin"
	"github.com/go-resty/resty/v2"
	"strconv"
)

func GetTicket(accessToken string) (*GetTicketResponse, error) {
	resp, err := resty.
		New().
		R().
		SetHeader("Content-Type", openDouyin.ApplicationJson).
		SetHeader("access-token", accessToken).
		Post(openDouyin.BaseDomain + "/js/getticket/")

	if err != nil {
		return nil, openDouyin.NewOpenDouyinError(90001, openDouyin.BaseDomain+"/js/getticket/", "请求失败："+err.Error(), resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, openDouyin.NewOpenDouyinError(90001, openDouyin.BaseDomain+"/js/getticket/", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var getTicketResponse GetTicketResponse

	if err := getTicketResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &getTicketResponse, nil
}

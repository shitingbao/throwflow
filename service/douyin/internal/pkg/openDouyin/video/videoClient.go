package video

import (
	"douyin/internal/pkg/openDouyin"
	"github.com/go-resty/resty/v2"
	"strconv"
)

func ListVideo(accessToken, openId string, cursor int64) (*ListVideoResponse, error) {
	listVideoRequest := ListVideoRequest{
		OpenId: openId,
		Cursor: cursor,
		Count:  openDouyin.PageSize20,
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(listVideoRequest.String()).
		SetHeader("Access-Token", accessToken).
		SetHeader("Content-Type", openDouyin.ApplicationJson).
		Get(openDouyin.BaseDomain + "/api/douyin/v1/video/video_list/")

	if err != nil {
		return nil, openDouyin.NewOpenDouyinError(90001, openDouyin.BaseDomain+"/api/douyin/v1/video/video_list/", "请求失败："+err.Error(), resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, openDouyin.NewOpenDouyinError(90001, openDouyin.BaseDomain+"/api/douyin/v1/video/video_list/", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var listVideoResponse ListVideoResponse

	if err := listVideoResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &listVideoResponse, nil
}

package video

import (
	"douyin/internal/pkg/openDouyin"
	"encoding/json"
)

type VideoStatistics struct {
	CommentCount  int32 `json:"comment_count"`  // 评论数
	DiggCount     int32 `json:"digg_count"`     // 点赞数
	DownloadCount int32 `json:"download_count"` // 下载数
	ForwardCount  int32 `json:"forward_count"`  // 转发数
	PlayCount     int32 `json:"play_count"`     // 播放数，只有作者本人可见。公开视频设为私密后，播放数也会返回0。
	ShareCount    int32 `json:"share_count"`    // 分享数
}

type Video struct {
	Title       string          `json:"title"`        // 视频标题
	Cover       string          `json:"cover"`        // 视频封面
	CreateTime  int64           `json:"create_time"`  // 视频创建时间戳
	IsReviewed  bool            `json:"is_reviewed"`  // 表示是否审核结束。审核通过或者失败都会返回true，审核中返回false。
	ItemId      string          `json:"item_id"`      // 视频id
	Statistics  VideoStatistics `json:"statistics"`   // 统计数据
	IsTop       bool            `json:"is_top"`       // 是否置顶
	MediaType   int             `json:"media_type"`   // 媒体类型。2:图集;4:视频
	ShareUrl    string          `json:"share_url"`    // 视频播放页面。视频播放页可能会失效，请在观看视频前调用/video/data/获取最新的播放页。
	VideoId     string          `json:"video_id"`     // 视频id
	VideoStatus int32           `json:"video_status"` // 表示视频状态。1:细化为5、6、7三种状态;2:不适宜公开;4:审核中;5:公开视频;6:好友可见;7:私密视频
}

type ListVideoDataResponse struct {
	ErrorCode   uint64  `json:"error_code"`
	Description string  `json:"description"`
	Cursor      int64   `json:"cursor"`
	HasMore     bool    `json:"has_more"`
	List        []Video `json:"list"`
}

type ListVideoResponse struct {
	openDouyin.CommonResponse
	Data ListVideoDataResponse `json:"data"`
}

func (lvr *ListVideoResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), lvr); err != nil {
		return openDouyin.NewOpenDouyinError(90000, openDouyin.BaseDomain+"/api/douyin/v1/video/video_list/", "解析json失败："+err.Error(), response)
	} else {
		if lvr.CommonResponse.Extra.ErrorCode != 0 {
			return openDouyin.NewOpenDouyinError(lvr.CommonResponse.Extra.ErrorCode, openDouyin.BaseDomain+"/api/douyin/v1/video/video_list/", lvr.CommonResponse.Extra.Description, response)
		}
	}

	return nil
}

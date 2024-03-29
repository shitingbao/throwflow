package tools

import (
	"douyin/internal/pkg/tool"
)

type AddAwemeAuthRequest struct {
	AdvertiserId uint64 `json:"advertiser_id"`
	AwemeId      string `json:"aweme_id"`
	Code         string `json:"code"`
	AuthType     string `json:"auth_type"`
	EndTime      string `json:"end_time"`
}

func (aaar AddAwemeAuthRequest) String() string {
	data, _ := tool.Marshal(aaar)

	return string(data)
}

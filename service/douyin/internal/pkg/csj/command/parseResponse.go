package command

import (
	"douyin/internal/pkg/csj"
	"douyin/internal/pkg/douke"
	"encoding/json"
	"strings"
)

type Qrcode struct {
	Url    string `json:"url"`
	Width  uint64 `json:"width"`  // 图片宽度
	Height uint64 `json:"height"` // 图片高度
}

type ShareInfo struct {
	DeepLink     string `json:"deep_link"`     // deeplink
	ShareCommand string `json:"share_command"` // 抖口令
	Qrcode       Qrcode `json:"qrcode"`        // 二维码
	Zlink        string `json:"zlink"`         // zlink
}

type ProductInfo struct {
	ProductId        uint64    `json:"product_id"`         // 商品id
	InsActivityParam string    `json:"ins_activity_param"` // 团长参数
	ShareInfo        ShareInfo `json:"share_info"`         // 分享物料
}

type LiveInfo struct {
	AuthorBuyinId string    `json:"author_buyin_id"` // 直播间buyin_id
	ProductId     uint64    `json:"product_id"`      // 商品id
	ShareInfo     ShareInfo `json:"share_info"`      // 分享物料
}

type ActivityInfo struct {
	MaterialId  string    `json:"material_id"`  // 聚合页物料id
	ExtraParams string    `json:"extra_params"` // 活动页自定义参数
	ShareInfo   ShareInfo `json:"share_info"`   // 分享物料
}

type ParseDataResponse struct {
	CommandType  uint64       `json:"command_type"`   // 抖口令解析出的原始物料类型。1-商品，2-直播间/直播间商品，3-活动页, 可解析聚合页接口生成的抖口令
	ProductInfo  ProductInfo  `json:"product_info"`   // 抖口令解析出的商品信息
	LiveInfo     LiveInfo     `json:"live_info"`      // 抖口令解析出的直播间/直播间商品信息
	ActivityInfo ActivityInfo `json:"activity_info"`  // 抖口令解析出的活动页信息
	IsOwnCommand bool         `json:"is_own_command"` // 是否是自己的抖口令吗，true:属于自己，false:属于其他抖客，空:没有抖客归属的抖口令，解析有抖客归属的抖口令时一定返回
}

type ParseResponse struct {
	csj.CommonResponse
	Data ParseDataResponse `json:"data"`
}

func (pr *ParseResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), pr); err != nil {
		return csj.NewCsjError(90100, douke.BaseDomain+"/command_parse", "解析json失败："+err.Error(), response)
	} else {
		if pr.CommonResponse.Code != 0 {
			return csj.NewCsjError(pr.CommonResponse.Code, douke.BaseDomain+"/command_parse", strings.Replace(pr.CommonResponse.Desc, "command parse error:", "", -1), response)
		}
	}

	return nil
}

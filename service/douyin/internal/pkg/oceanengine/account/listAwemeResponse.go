package account

import (
	"douyin/internal/pkg/oceanengine"
	"encoding/json"
)

type ListAwemes struct {
	AwemeAvatar             string   `json:"aweme_avatar"`
	AwemeId                 uint64   `json:"aweme_id"`
	AwemeShowId             string   `json:"aweme_show_id"`
	AwemeName               string   `json:"aweme_name"`
	AwemeStatus             string   `json:"aweme_status"`
	BindType                []string `json:"bind_type"`
	AwemeHasVideoPermission bool     `json:"aweme_has_video_permission"`
	AwemeHasLivePermission  bool     `json:"aweme_has_live_permission"`
	AwemeHasUniProm         bool     `json:"aweme_has_uni_prom"`
}

type ListAwemeData struct {
	AwemeIdList []*ListAwemes            `json:"aweme_id_list"`
	PageInfo    oceanengine.PageResponse `json:"page_info"`
}

type ListAwemeResponse struct {
	oceanengine.CommonResponse
	Data ListAwemeData `json:"data"`
}

func (lar *ListAwemeResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), lar); err != nil {
		return oceanengine.NewOceanengineError(90000, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/aweme/authorized/get/", "JSON_UNMARSHAL_ERROR", "解析json失败："+err.Error(), "", response)
	} else {
		if lar.Code != 0 {
			return oceanengine.NewOceanengineError(lar.Code, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/aweme/authorized/get/", lar.Message, oceanengine.ResponseDescription[lar.Code], lar.RequestId, response)
		}
	}

	return nil
}

package account

import (
	"douyin/internal/pkg/oceanengine"
	"encoding/json"
)

type ListAgentAdvertiserData struct {
	List     []*uint64                `json:"list"`
	PageInfo oceanengine.PageResponse `json:"page_info"`
}

type ListAgentAdvertiserResponse struct {
	oceanengine.CommonResponse
	Data ListAgentAdvertiserData `json:"data"`
}

func (laar *ListAgentAdvertiserResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), laar); err != nil {
		return oceanengine.NewOceanengineError(90000, oceanengine.BaseDomain+"/open_api/2/agent/advertiser/select/", "JSON_UNMARSHAL_ERROR", "解析json失败："+err.Error(), "", response)
	} else {
		if laar.Code != 0 {
			return oceanengine.NewOceanengineError(laar.Code, oceanengine.BaseDomain+"/open_api/2/agent/advertiser/select/", laar.Message, oceanengine.ResponseDescription[laar.Code], laar.RequestId, response)
		}
	}

	return nil
}

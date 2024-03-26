package suolink

import (
	"common/internal/pkg/suolink"
	"encoding/json"
)

type GetSuoLinkResponse struct {
	Url string `json:"url"`
	Err string `json:"err"`
}

func (gslr *GetSuoLinkResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), gslr); err != nil {
		return suolink.NewSuoLinkError(suolink.BaseDomain, "解析json失败："+err.Error(), response)
	} else {
		if len(gslr.Err) > 0 {
			return suolink.NewSuoLinkError(suolink.BaseDomain, gslr.Err, response)
		}
	}

	return nil
}

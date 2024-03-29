package user

import (
	"douyin/internal/pkg/openDouyin"
	"encoding/json"
)

type UserFan struct {
	Date      string `json:"date"`       // 日期	yyyy-MM-dd
	TotalFans int64  `json:"total_fans"` // 每日总粉丝数
	NewFans   int64  `json:"new_fans"`   // 每天新粉丝数
}

type ListUserFansDataResponse struct {
	ErrorCode   uint64    `json:"error_code"`
	Description string    `json:"description"`
	ResultList  []UserFan `json:"result_list"`
}

type ListUserFansResponse struct {
	openDouyin.CommonResponse
	Data ListUserFansDataResponse `json:"data"`
}

func (lufr *ListUserFansResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), lufr); err != nil {
		return openDouyin.NewOpenDouyinError(90000, openDouyin.BaseDomain+"/data/external/user/fans/", "解析json失败："+err.Error(), response)
	} else {
		if lufr.CommonResponse.Extra.ErrorCode != 0 {
			return openDouyin.NewOpenDouyinError(lufr.CommonResponse.Extra.ErrorCode, openDouyin.BaseDomain+"/data/external/user/fans/", lufr.CommonResponse.Extra.Description, response)
		}
	}

	return nil
}

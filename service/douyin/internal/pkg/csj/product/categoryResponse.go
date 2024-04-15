package product

import (
	"douyin/internal/pkg/csj"
	"douyin/internal/pkg/douke"
	"encoding/json"
)

type ListDataCategoryResponse struct {
	Id    uint64 `json:"id"`    // 类⽬id
	Name  string `json:"name"`  // 类⽬名称
	Level uint64 `json:"level"` // 类⽬层级， 最多三级
}

type CategoryDataResponse struct {
	Total        uint64                     `json:"total"`
	CategoryList []ListDataCategoryResponse `json:"category_list"`
}

type CategoryResponse struct {
	csj.CommonResponse
	Data CategoryDataResponse `json:"data"`
}

func (cr *CategoryResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), cr); err != nil {
		return csj.NewCsjError(90100, douke.BaseDomain+"/product/category", "解析json失败："+err.Error(), response)
	} else {
		if cr.CommonResponse.Code != 0 {
			return csj.NewCsjError(cr.CommonResponse.Code, douke.BaseDomain+"/product/category", cr.CommonResponse.Desc, response)
		}
	}

	return nil
}

package product

import (
	"douyin/internal/pkg/oceanengine"
	"encoding/json"
)

type ProductImgList struct {
	ImgUrl string `json:"img_url"`
}

type ListAvailableProducts struct {
	CategoryName          string           `json:"category_name"`
	DiscountPrice         float64          `json:"discount_price"`
	MarketPrice           float64          `json:"market_price"`
	DiscountLowerPrice    float64          `json:"discount_lower_price"`
	DiscountHigherPrice   float64          `json:"discount_higher_price"`
	Id                    uint64           `json:"id"`
	Img                   string           `json:"img"`
	Inventory             uint64           `json:"inventory"`
	Name                  string           `json:"name"`
	ProductRate           float64          `json:"product_rate"`
	SaleTime              string           `json:"sale_time"`
	Tags                  string           `json:"tags"`
	SupportProductNewOpen bool             `json:"support_product_new_open"`
	ImgList               []ProductImgList `json:"img_list"`
	ChannelId             uint64           `json:"channel_id"`
	ChannelType           string           `json:"channel_type"`
}

type ListAvailableProductData struct {
	ProductList []*ListAvailableProducts `json:"product_list"`
	PageInfo    oceanengine.PageResponse `json:"page_info"`
}

type ListAvailableProductResponse struct {
	oceanengine.CommonResponse
	Data ListAvailableProductData `json:"data"`
}

func (lapr *ListAvailableProductResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), lapr); err != nil {
		return oceanengine.NewOceanengineError(90000, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/product/available/get/", "JSON_UNMARSHAL_ERROR", "解析json失败："+err.Error(), "", response)
	} else {
		if lapr.Code != 0 {
			return oceanengine.NewOceanengineError(lapr.Code, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/product/available/get/", lapr.Message, oceanengine.ResponseDescription[lapr.Code], lapr.RequestId, response)
		}
	}

	return nil
}

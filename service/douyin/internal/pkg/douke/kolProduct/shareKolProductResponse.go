package kolProduct

import (
	"douyin/internal/pkg/douke"
	"encoding/json"
)

type ShareKolProductQrcodeResponse struct {
	Url    string `json:"url"`
	Width  uint32 `json:"width"`
	Height uint32 `json:"height"`
}

type ShareKolProductCouponLinkResponse struct {
	CouponStatus int64                         `json:"coupon_status"`
	Qrcode       ShareKolProductQrcodeResponse `json:"qrcode"`
	ShareCommand string                        `json:"share_command"`
	Deeplink     string                        `json:"deeplink"`
	ShareLink    string                        `json:"share_link"`
}

type ShareKolProductDataResponse struct {
	UseInsActivity string                            `json:"use_ins_activity"`
	ShareLink      string                            `json:"share_link"`
	DyPassword     string                            `json:"dy_password"`
	QrCode         ShareKolProductQrcodeResponse     `json:"qr_code"`
	DyDeeplink     string                            `json:"dy_deeplink"`
	DyZlink        string                            `json:"dy_zlink"`
	CouponLink     ShareKolProductCouponLinkResponse `json:"coupon_link"`
}

type ShareKolProductResponse struct {
	douke.CommonResponse
	Data ShareKolProductDataResponse `json:"data"`
}

func (skpr *ShareKolProductResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), skpr); err != nil {
		return douke.NewDoukeError(90100, douke.BaseDomain+"/buyin/kolProductShare", "解析json失败："+err.Error(), response)
	} else {
		if skpr.CommonResponse.Code != 10000 {
			return douke.NewDoukeError(skpr.CommonResponse.Code, douke.BaseDomain+"/buyin/kolProductShare", skpr.CommonResponse.SubMsg, response)
		}
	}

	return nil
}

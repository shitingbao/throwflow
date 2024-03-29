package kolProduct

import (
	"douyin/internal/pkg/douke"
)

type ShareKolProductRequest struct {
	ProductUrl       string `json:"product_url"`
	Pid              string `json:"pid"`
	ExternalInfo     string `json:"external_info"`
	NeedQrCode       bool   `json:"need_qr_code"`
	Platform         int32  `json:"platform"`
	UseCoupon        bool   `json:"use_coupon"`
	NeedShareLink    bool   `json:"need_share_link"`
	InsActivityParam string `json:"ins_activity_param"`
	NeedZlink        bool   `json:"need_zlink"`
}

func (skpr ShareKolProductRequest) String() string {
	return douke.Marshal(skpr)
}

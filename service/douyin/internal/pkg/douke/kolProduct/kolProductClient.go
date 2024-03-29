package kolProduct

import (
	"douyin/internal/pkg/douke"
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
)

func ShareKolProduct(appKey, appSecret, accessToken, productUrl, pid, externalInfo string) (*ShareKolProductResponse, error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	shareKolProductRequest := ShareKolProductRequest{
		ProductUrl:       productUrl,
		Pid:              pid,
		ExternalInfo:     externalInfo,
		NeedQrCode:       false,
		Platform:         0,
		UseCoupon:        false,
		NeedShareLink:    false,
		InsActivityParam: "",
		NeedZlink:        false,
	}

	param := shareKolProductRequest.String()

	commonRequest := douke.CommonRequest{
		Method:      "buyin.kolProductShare",
		AppKey:      appKey,
		AccessToken: accessToken,
		ParamJson:   param,
		Timestamp:   timestamp,
		V:           douke.V,
		Sign:        douke.Sign(appKey, appSecret, "buyin.kolProductShare", param, timestamp, douke.V),
		SignMethod:  douke.SignMethod,
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(commonRequest.String()).
		SetHeader("Content-Type", douke.ApplicationJson).
		SetBody(shareKolProductRequest.String()).
		Post(douke.BaseDomain + "/buyin/kolProductShare")

	if err != nil {
		return nil, douke.NewDoukeError(90101, douke.BaseDomain+"/buyin/kolProductShare", "请求失败："+err.Error(), resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, douke.NewDoukeError(90101, douke.BaseDomain+"/buyin/kolProductShare", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), resp.String())
	}

	var shareKolProductResponse ShareKolProductResponse

	if err := shareKolProductResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &shareKolProductResponse, nil
}

package finance

import (
	"douyin/internal/pkg/oceanengine"
	"github.com/go-resty/resty/v2"
	"strconv"
)

func GetWallet(advertiserId uint64, accessToken string) (*GetWalletResponse, error) {
	getWalletRequest := GetWalletRequest{
		AdvertiserId: advertiserId,
	}

	resp, err := resty.
		New().
		R().
		SetQueryString(getWalletRequest.String()).
		SetHeader("Access-Token", accessToken).
		SetHeader("Content-Type", oceanengine.ApplicationJson).
		Get(oceanengine.BaseDomain + "/open_api/v1.0/qianchuan/finance/wallet/get/")

	if err != nil {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/finance/wallet/get/", "REQUEST_ERROR", "请求失败："+err.Error(), "", resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/finance/wallet/get/", "REQUEST_ERROR", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), "", resp.String())
	}

	var getWalletResponse GetWalletResponse

	if err := getWalletResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &getWalletResponse, nil
}

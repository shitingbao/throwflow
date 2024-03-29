package report

import (
	"douyin/internal/pkg/oceanengine"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"strconv"
)

func GetReportLive(advertiserId, awemeId uint64, accessToken, startTime, endTime string) (*GetReportLiveResponse, error) {
	fields := []string{}
	fields = append(fields, "stat_cost")
	fields = append(fields, "cpm_platform")
	fields = append(fields, "click_cnt")
	fields = append(fields, "ctr")
	fields = append(fields, "total_live_pay_order_gpm")
	fields = append(fields, "luban_live_pay_order_gpm")
	fields = append(fields, "cpc_platform")
	fields = append(fields, "convert_cnt")
	fields = append(fields, "convert_rate")
	fields = append(fields, "cpa_platform")
	fields = append(fields, "live_pay_order_gmv_alias")
	fields = append(fields, "luban_live_pay_order_gmv")
	fields = append(fields, "live_pay_order_gmv_roi")
	fields = append(fields, "ad_live_prepay_and_pay_order_gmv_roi")
	fields = append(fields, "live_create_order_count_alias")
	fields = append(fields, "live_create_order_rate")
	fields = append(fields, "luban_live_order_count")
	fields = append(fields, "ad_live_create_order_rate")
	fields = append(fields, "live_pay_order_count_alias")
	fields = append(fields, "live_pay_order_rate")
	fields = append(fields, "luban_live_pay_order_count")
	fields = append(fields, "ad_live_pay_order_rate")
	fields = append(fields, "live_pay_order_gmv_avg")
	fields = append(fields, "ad_live_pay_order_gmv_avg")
	fields = append(fields, "luban_live_prepay_order_count")
	fields = append(fields, "luban_live_prepay_order_gmv")
	fields = append(fields, "live_prepay_order_count_alias")
	fields = append(fields, "live_prepay_order_gmv_alias")
	fields = append(fields, "live_order_pay_coupon_amount")
	fields = append(fields, "total_live_watch_cnt")
	fields = append(fields, "total_live_follow_cnt")
	fields = append(fields, "live_watch_one_minute_count")
	fields = append(fields, "total_live_fans_club_join_cnt")
	fields = append(fields, "live_click_cart_count_alias")
	fields = append(fields, "live_click_product_count_alias")
	fields = append(fields, "total_live_comment_cnt")
	fields = append(fields, "total_live_share_cnt")
	fields = append(fields, "total_live_gift_cnt")
	fields = append(fields, "total_live_gift_amount")

	sfields, _ := json.Marshal(fields)

	getReportLiveRequest := GetReportLiveRequest{
		AdvertiserId: advertiserId,
		AwemeId:      awemeId,
		StartTime:    startTime,
		EndTime:      endTime,
		Fields:       string(sfields),
	}
	
	resp, err := resty.
		New().
		R().
		SetQueryString(getReportLiveRequest.String()).
		SetHeader("Access-Token", accessToken).
		SetHeader("Content-Type", oceanengine.ApplicationJson).
		Get(oceanengine.BaseDomain + "/open_api/v1.0/qianchuan/report/live/get/")

	if err != nil {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/report/live/get/", "REQUEST_ERROR", "请求失败："+err.Error(), "", resp.String())
	}

	if resp.StatusCode() != 200 {
		return nil, oceanengine.NewOceanengineError(90001, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/report/live/get/", "REQUEST_ERROR", "请求失败，状态码："+strconv.Itoa(resp.StatusCode()), "", resp.String())
	}

	var getReportLiveResponse GetReportLiveResponse

	if err := getReportLiveResponse.Struct(resp.String()); err != nil {
		return nil, err
	}

	return &getReportLiveResponse, nil
}

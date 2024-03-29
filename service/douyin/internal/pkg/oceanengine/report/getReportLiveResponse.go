package report

import (
	"douyin/internal/pkg/oceanengine"
	"encoding/json"
)

type GetReportLiveData struct {
	StatCost                      float64 `json:"stat_cost"`
	CpmPlatform                   float64 `json:"cpm_platform"`
	ClickCnt                      int64   `json:"click_cnt"`
	Ctr                           float64 `json:"ctr"`
	TotalLivePayOrderGpm          float64 `json:"total_live_pay_order_gpm"`
	LubanLivePayOrderGpm          float64 `json:"luban_live_pay_order_gpm"`
	CpcPlatform                   float64 `json:"cpc_platform"`
	ConvertCnt                    int64   `json:"convert_cnt"`
	ConvertRate                   float64 `json:"convert_rate"`
	CpaPlatform                   float64 `json:"cpa_platform"`
	LivePayOrderGmvAlias          float64 `json:"live_pay_order_gmv_alias"`
	LubanLivePayOrderGmv          float64 `json:"luban_live_pay_order_gmv"`
	LivePayOrderGmvRoi            float64 `json:"live_pay_order_gmv_roi"`
	AdLivePrepayAndPayOrderGmvRoi float64 `json:"ad_live_prepay_and_pay_order_gmv_roi"`
	LiveCreateOrderCountAlias     int64   `json:"live_create_order_count_alias"`
	LiveCreateOrderRate           float64 `json:"live_create_order_rate"`
	LubanLiveOrderCount           int64   `json:"luban_live_order_count"`
	AdLiveCreateOrderRate         float64 `json:"ad_live_create_order_rate"`
	LivePayOrderCountAlias        int64   `json:"live_pay_order_count_alias"`
	LivePayOrderRate              float64 `json:"live_pay_order_rate"`
	LubanLivePayOrderCount        int64   `json:"luban_live_pay_order_count"`
	AdLivePayOrderRate            float64 `json:"ad_live_pay_order_rate"`
	LivePayOrderGmvAvg            float64 `json:"live_pay_order_gmv_avg"`
	AdLivePayOrderGmvAvg          float64 `json:"ad_live_pay_order_gmv_avg"`
	LubanLivePrepayOrderCount     int64   `json:"luban_live_prepay_order_count"`
	LubanLivePrepayOrderGmv       float64 `json:"luban_live_prepay_order_gmv"`
	LivePrepayOrderCountAlias     float64 `json:"live_prepay_order_count_alias"`
	LivePrepayOrderGmvAlias       float64 `json:"live_prepay_order_gmv_alias"`
	LiveOrderPayCouponAmount      float64 `json:"live_order_pay_coupon_amount"`
	TotalLiveWatchCnt             int64   `json:"total_live_watch_cnt"`
	TotalLiveFollowCnt            int64   `json:"total_live_follow_cnt"`
	LiveWatchOneMinuteCount       int64   `json:"live_watch_one_minute_count"`
	TotalLiveFansClubJoinCnt      int64   `json:"total_live_fans_club_join_cnt"`
	LiveClickCartCountAlias       int64   `json:"live_click_cart_count_alias"`
	LiveClickProductCountAlias    int64   `json:"live_click_product_count_alias"`
	TotalLiveCommentCnt           int64   `json:"total_live_comment_cnt"`
	TotalLiveShareCnt             int64   `json:"total_live_share_cnt"`
	TotalLiveGiftCnt              int64   `json:"total_live_gift_cnt"`
	TotalLiveGiftAmount           float64 `json:"total_live_gift_amount"`
}

type GetReportLiveResponse struct {
	oceanengine.CommonResponse
	Data GetReportLiveData `json:"data"`
}

func (grlr *GetReportLiveResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), grlr); err != nil {
		return oceanengine.NewOceanengineError(90000, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/report/live/get/", "JSON_UNMARSHAL_ERROR", "解析json失败："+err.Error(), "", response)
	} else {
		if grlr.Code != 0 {
			return oceanengine.NewOceanengineError(grlr.Code, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/report/live/get/", grlr.Message, oceanengine.ResponseDescription[grlr.Code], grlr.RequestId, response)
		}
	}

	return nil
}

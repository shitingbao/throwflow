package finance

import (
	"douyin/internal/pkg/oceanengine"
	"encoding/json"
)

type ShareExpiringDetailList struct {
	Category   string  `json:"category"`
	Amount     float64 `json:"amount"`
	ExpireTime uint64  `json:"expire_time"`
}

type GetWalletData struct {
	TotalBalanceAbs                 float64                    `json:"total_balance_abs"`
	GrantBalance                    float64                    `json:"grant_balance"`
	UnionValidGrantBalance          float64                    `json:"union_valid_grant_balance"`
	SearchValidGrantBalance         float64                    `json:"search_valid_grant_balance"`
	CommonValidGrantBalance         float64                    `json:"common_valid_grant_balance"`
	DefaultValidGrantBalance        float64                    `json:"default_valid_grant_balance"`
	GeneralTotalBalance             float64                    `json:"general_total_balance"`
	GeneralBalanceValid             float64                    `json:"general_balance_valid"`
	GeneralBalanceValidNonGrant     float64                    `json:"general_balance_valid_non_grant"`
	GeneralBalanceValidGrantUnion   float64                    `json:"general_balance_valid_grant_union"`
	GeneralBalanceValidGrantSearch  float64                    `json:"general_balance_valid_grant_search"`
	GeneralBalanceValidGrantCommon  float64                    `json:"general_balance_valid_grant_common"`
	GeneralBalanceValidGrantDefault float64                    `json:"general_balance_valid_grant_default"`
	GeneralBalanceInvalid           float64                    `json:"general_balance_invalid"`
	GeneralBalanceInvalidOrder      float64                    `json:"general_balance_invalid_order"`
	GeneralBalanceInvalidFrozen     float64                    `json:"general_balance_invalid_frozen"`
	BrandBalance                    float64                    `json:"brand_balance"`
	BrandBalanceValid               float64                    `json:"brand_balance_valid"`
	BrandBalanceValidNonGrant       float64                    `json:"brand_balance_valid_non_grant"`
	BrandBalanceValidGrant          float64                    `json:"brand_balance_valid_grant"`
	BrandBalanceInvalid             float64                    `json:"brand_balance_invalid"`
	BrandBalanceInvalidFrozen       float64                    `json:"brand_balance_invalid_frozen"`
	DeductionCouponBalance          float64                    `json:"deduction_coupon_balance"`
	DeductionCouponBalanceAll       float64                    `json:"deduction_coupon_balance_all"`
	DeductionCouponBalanceOther     float64                    `json:"deduction_coupon_balance_other"`
	DeductionCouponBalanceSelf      float64                    `json:"deduction_coupon_balance_self"`
	GrantExpiring                   float64                    `json:"grant_expiring"`
	ShareBalance                    float64                    `json:"share_balance"`
	ShareBalanceValidGrantUnion     float64                    `json:"share_balance_valid_grant_union"`
	ShareBalanceValidGrantSearch    float64                    `json:"share_balance_valid_grant_search"`
	ShareBalanceValidGrantCommon    float64                    `json:"share_balance_valid_grant_common"`
	ShareBalanceValidGrantDefault   float64                    `json:"share_balance_valid_grant_default"`
	ShareBalanceValid               float64                    `json:"share_balance_valid"`
	ShareBalanceExpiring            float64                    `json:"share_balance_expiring"`
	ShareExpiringDetailList         []*ShareExpiringDetailList `json:"share_expiring_detail_list"`
}

type GetWalletResponse struct {
	oceanengine.CommonResponse
	Data GetWalletData `json:"data"`
}

func (gwr *GetWalletResponse) Struct(response string) error {
	if err := json.Unmarshal([]byte(response), gwr); err != nil {
		return oceanengine.NewOceanengineError(90000, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/finance/wallet/get/", "JSON_UNMARSHAL_ERROR", "解析json失败："+err.Error(), "", response)
	} else {
		if gwr.Code != 0 {
			return oceanengine.NewOceanengineError(gwr.Code, oceanengine.BaseDomain+"/open_api/v1.0/qianchuan/finance/wallet/get/", gwr.Message, oceanengine.ResponseDescription[gwr.Code], gwr.RequestId, response)
		}
	}

	return nil
}

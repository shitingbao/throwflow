package domain

import (
	"context"
	"time"
)

type ShareExpiringDetailList struct {
	Category   string
	Amount     float64
	ExpireTime uint64
}

type QianchuanWallet struct {
	AdvertiserId                    uint64
	TotalBalanceAbs                 float64
	GrantBalance                    float64
	UnionValidGrantBalance          float64
	SearchValidGrantBalance         float64
	CommonValidGrantBalance         float64
	DefaultValidGrantBalance        float64
	GeneralTotalBalance             float64
	GeneralBalanceValid             float64
	GeneralBalanceValidNonGrant     float64
	GeneralBalanceValidGrantUnion   float64
	GeneralBalanceValidGrantSearch  float64
	GeneralBalanceValidGrantCommon  float64
	GeneralBalanceValidGrantDefault float64
	GeneralBalanceInvalid           float64
	GeneralBalanceInvalidOrder      float64
	GeneralBalanceInvalidFrozen     float64
	BrandBalance                    float64
	BrandBalanceValid               float64
	BrandBalanceValidNonGrant       float64
	BrandBalanceValidGrant          float64
	BrandBalanceInvalid             float64
	BrandBalanceInvalidFrozen       float64
	DeductionCouponBalance          float64
	DeductionCouponBalanceAll       float64
	DeductionCouponBalanceOther     float64
	DeductionCouponBalanceSelf      float64
	GrantExpiring                   float64
	ShareBalance                    float64
	ShareBalanceValidGrantUnion     float64
	ShareBalanceValidGrantSearch    float64
	ShareBalanceValidGrantCommon    float64
	ShareBalanceValidGrantDefault   float64
	ShareBalanceValid               float64
	ShareBalanceExpiring            float64
	ShareExpiringDetailList         []*ShareExpiringDetailList
	CreateTime                      time.Time
	UpdateTime                      time.Time
}

func NewQianchuanWallet(ctx context.Context, advertiserId uint64,
	totalBalanceAbs,
	grantBalance,
	unionValidGrantBalance,
	searchValidGrantBalance,
	commonValidGrantBalance,
	defaultValidGrantBalance,
	generalTotalBalance,
	generalBalanceValid,
	generalBalanceValidNonGrant,
	generalBalanceValidGrantUnion,
	generalBalanceValidGrantSearch,
	generalBalanceValidGrantCommon,
	generalBalanceValidGrantDefault,
	generalBalanceInvalid,
	generalBalanceInvalidOrder,
	generalBalanceInvalidFrozen,
	brandBalance,
	brandBalanceValid,
	brandBalanceValidNonGrant,
	brandBalanceValidGrant,
	brandBalanceInvalid,
	brandBalanceInvalidFrozen,
	deductionCouponBalance,
	deductionCouponBalanceAll,
	deductionCouponBalanceOther,
	deductionCouponBalanceSelf,
	grantExpiring,
	shareBalance,
	shareBalanceValidGrantUnion,
	shareBalanceValidGrantSearch,
	shareBalanceValidGrantCommon,
	shareBalanceValidGrantDefault,
	shareBalanceValid,
	shareBalanceExpiring float64,
	shareExpiringDetailList []*ShareExpiringDetailList) *QianchuanWallet {
	return &QianchuanWallet{
		AdvertiserId:                    advertiserId,
		TotalBalanceAbs:                 totalBalanceAbs,
		GrantBalance:                    grantBalance,
		UnionValidGrantBalance:          unionValidGrantBalance,
		SearchValidGrantBalance:         searchValidGrantBalance,
		CommonValidGrantBalance:         commonValidGrantBalance,
		DefaultValidGrantBalance:        defaultValidGrantBalance,
		GeneralTotalBalance:             generalTotalBalance,
		GeneralBalanceValid:             generalBalanceValid,
		GeneralBalanceValidNonGrant:     generalBalanceValidNonGrant,
		GeneralBalanceValidGrantUnion:   generalBalanceValidGrantUnion,
		GeneralBalanceValidGrantSearch:  generalBalanceValidGrantSearch,
		GeneralBalanceValidGrantCommon:  generalBalanceValidGrantCommon,
		GeneralBalanceValidGrantDefault: generalBalanceValidGrantDefault,
		GeneralBalanceInvalid:           generalBalanceInvalid,
		GeneralBalanceInvalidOrder:      generalBalanceInvalidOrder,
		GeneralBalanceInvalidFrozen:     generalBalanceInvalidFrozen,
		BrandBalance:                    brandBalance,
		BrandBalanceValid:               brandBalanceValid,
		BrandBalanceValidNonGrant:       brandBalanceValidNonGrant,
		BrandBalanceValidGrant:          brandBalanceValidGrant,
		BrandBalanceInvalid:             brandBalanceInvalid,
		BrandBalanceInvalidFrozen:       brandBalanceInvalidFrozen,
		DeductionCouponBalance:          deductionCouponBalance,
		DeductionCouponBalanceAll:       deductionCouponBalanceAll,
		DeductionCouponBalanceOther:     deductionCouponBalanceOther,
		DeductionCouponBalanceSelf:      deductionCouponBalanceSelf,
		GrantExpiring:                   grantExpiring,
		ShareBalance:                    shareBalance,
		ShareBalanceValidGrantUnion:     shareBalanceValidGrantUnion,
		ShareBalanceValidGrantSearch:    shareBalanceValidGrantSearch,
		ShareBalanceValidGrantCommon:    shareBalanceValidGrantCommon,
		ShareBalanceValidGrantDefault:   shareBalanceValidGrantDefault,
		ShareBalanceValid:               shareBalanceValid,
		ShareBalanceExpiring:            shareBalanceExpiring,
		ShareExpiringDetailList:         shareExpiringDetailList,
	}
}

func (qw *QianchuanWallet) SetAdvertiserId(ctx context.Context, advertiserId uint64) {
	qw.AdvertiserId = advertiserId
}

func (qw *QianchuanWallet) SetTotalBalanceAbs(ctx context.Context, totalBalanceAbs float64) {
	qw.TotalBalanceAbs = totalBalanceAbs
}

func (qw *QianchuanWallet) SetGrantBalance(ctx context.Context, grantBalance float64) {
	qw.GrantBalance = grantBalance
}

func (qw *QianchuanWallet) SetUnionValidGrantBalance(ctx context.Context, unionValidGrantBalance float64) {
	qw.UnionValidGrantBalance = unionValidGrantBalance
}

func (qw *QianchuanWallet) SetSearchValidGrantBalance(ctx context.Context, searchValidGrantBalance float64) {
	qw.SearchValidGrantBalance = searchValidGrantBalance
}

func (qw *QianchuanWallet) SetCommonValidGrantBalance(ctx context.Context, commonValidGrantBalance float64) {
	qw.CommonValidGrantBalance = commonValidGrantBalance
}

func (qw *QianchuanWallet) SetDefaultValidGrantBalance(ctx context.Context, defaultValidGrantBalance float64) {
	qw.DefaultValidGrantBalance = defaultValidGrantBalance
}

func (qw *QianchuanWallet) SetGeneralTotalBalance(ctx context.Context, generalTotalBalance float64) {
	qw.GeneralTotalBalance = generalTotalBalance
}

func (qw *QianchuanWallet) SetGeneralBalanceValid(ctx context.Context, generalBalanceValid float64) {
	qw.GeneralBalanceValid = generalBalanceValid
}

func (qw *QianchuanWallet) SetGeneralBalanceValidNonGrant(ctx context.Context, generalBalanceValidNonGrant float64) {
	qw.GeneralBalanceValidNonGrant = generalBalanceValidNonGrant
}

func (qw *QianchuanWallet) SetGeneralBalanceValidGrantUnion(ctx context.Context, generalBalanceValidGrantUnion float64) {
	qw.GeneralBalanceValidGrantUnion = generalBalanceValidGrantUnion
}

func (qw *QianchuanWallet) SetGeneralBalanceValidGrantSearch(ctx context.Context, generalBalanceValidGrantSearch float64) {
	qw.GeneralBalanceValidGrantSearch = generalBalanceValidGrantSearch
}

func (qw *QianchuanWallet) SetGeneralBalanceValidGrantCommon(ctx context.Context, generalBalanceValidGrantCommon float64) {
	qw.GeneralBalanceValidGrantCommon = generalBalanceValidGrantCommon
}

func (qw *QianchuanWallet) SetGeneralBalanceValidGrantDefault(ctx context.Context, generalBalanceValidGrantDefault float64) {
	qw.GeneralBalanceValidGrantDefault = generalBalanceValidGrantDefault
}

func (qw *QianchuanWallet) SetGeneralBalanceInvalid(ctx context.Context, generalBalanceInvalid float64) {
	qw.GeneralBalanceInvalid = generalBalanceInvalid
}

func (qw *QianchuanWallet) SetGeneralBalanceInvalidOrder(ctx context.Context, generalBalanceInvalidOrder float64) {
	qw.GeneralBalanceInvalidOrder = generalBalanceInvalidOrder
}

func (qw *QianchuanWallet) SetGeneralBalanceInvalidFrozen(ctx context.Context, generalBalanceInvalidFrozen float64) {
	qw.GeneralBalanceInvalidFrozen = generalBalanceInvalidFrozen
}

func (qw *QianchuanWallet) SetBrandBalance(ctx context.Context, brandBalance float64) {
	qw.BrandBalance = brandBalance
}

func (qw *QianchuanWallet) SetBrandBalanceValid(ctx context.Context, brandBalanceValid float64) {
	qw.BrandBalanceValid = brandBalanceValid
}

func (qw *QianchuanWallet) SetBrandBalanceValidNonGrant(ctx context.Context, brandBalanceValidNonGrant float64) {
	qw.BrandBalanceValidNonGrant = brandBalanceValidNonGrant
}

func (qw *QianchuanWallet) SetBrandBalanceValidGrant(ctx context.Context, brandBalanceValidGrant float64) {
	qw.BrandBalanceValidGrant = brandBalanceValidGrant
}

func (qw *QianchuanWallet) SetBrandBalanceInvalid(ctx context.Context, brandBalanceInvalid float64) {
	qw.BrandBalanceInvalid = brandBalanceInvalid
}

func (qw *QianchuanWallet) SetBrandBalanceInvalidFrozen(ctx context.Context, brandBalanceInvalidFrozen float64) {
	qw.BrandBalanceInvalidFrozen = brandBalanceInvalidFrozen
}

func (qw *QianchuanWallet) SetDeductionCouponBalance(ctx context.Context, deductionCouponBalance float64) {
	qw.DeductionCouponBalance = deductionCouponBalance
}

func (qw *QianchuanWallet) SetDeductionCouponBalanceAll(ctx context.Context, deductionCouponBalanceAll float64) {
	qw.DeductionCouponBalanceAll = deductionCouponBalanceAll
}

func (qw *QianchuanWallet) SetDeductionCouponBalanceOther(ctx context.Context, deductionCouponBalanceOther float64) {
	qw.DeductionCouponBalanceOther = deductionCouponBalanceOther
}

func (qw *QianchuanWallet) SetDeductionCouponBalanceSelf(ctx context.Context, deductionCouponBalanceSelf float64) {
	qw.DeductionCouponBalanceSelf = deductionCouponBalanceSelf
}

func (qw *QianchuanWallet) SetGrantExpiring(ctx context.Context, grantExpiring float64) {
	qw.GrantExpiring = grantExpiring
}

func (qw *QianchuanWallet) SetShareBalance(ctx context.Context, shareBalance float64) {
	qw.ShareBalance = shareBalance
}

func (qw *QianchuanWallet) SetShareBalanceValidGrantUnion(ctx context.Context, shareBalanceValidGrantUnion float64) {
	qw.ShareBalanceValidGrantUnion = shareBalanceValidGrantUnion
}

func (qw *QianchuanWallet) SetShareBalanceValidGrantSearch(ctx context.Context, shareBalanceValidGrantSearch float64) {
	qw.ShareBalanceValidGrantSearch = shareBalanceValidGrantSearch
}

func (qw *QianchuanWallet) SetShareBalanceValidGrantCommon(ctx context.Context, shareBalanceValidGrantCommon float64) {
	qw.ShareBalanceValidGrantCommon = shareBalanceValidGrantCommon
}

func (qw *QianchuanWallet) SetShareBalanceValidGrantDefault(ctx context.Context, shareBalanceValidGrantDefault float64) {
	qw.ShareBalanceValidGrantDefault = shareBalanceValidGrantDefault
}

func (qw *QianchuanWallet) SetShareBalanceValid(ctx context.Context, shareBalanceValid float64) {
	qw.ShareBalanceValid = shareBalanceValid
}

func (qw *QianchuanWallet) SetShareBalanceExpiring(ctx context.Context, shareBalanceExpiring float64) {
	qw.ShareBalanceExpiring = shareBalanceExpiring
}

func (qw *QianchuanWallet) SetShareExpiringDetailList(ctx context.Context, shareExpiringDetailList []*ShareExpiringDetailList) {
	qw.ShareExpiringDetailList = shareExpiringDetailList
}

func (qw *QianchuanWallet) SetUpdateTime(ctx context.Context) {
	qw.UpdateTime = time.Now()
}

func (qw *QianchuanWallet) SetCreateTime(ctx context.Context) {
	qw.CreateTime = time.Now()
}

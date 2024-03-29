package domain

import (
	"context"
	"time"
)

type JinritemaiStoreInfo struct {
	Id                 uint64
	ClientKey          string
	OpenId             string
	ProductId          string
	ProductName        string
	ProductImg         string
	ProductPrice       float32
	CommissionType     uint8
	CommissionTypeName string
	CommissionRatio    uint8
	PromotionId        uint64
	PromotionType      uint8
	PromotionTypeName  string
	ColonelActivityId  uint64
	CreateTime         time.Time
	UpdateTime         time.Time
}

func NewJinritemaiStoreInfo(ctx context.Context, colonelActivityId, promotionId uint64, commissionType, commissionRatio, promotionType uint8, productPrice float32, clientKey, openId, productId, productName, productImg string) *JinritemaiStoreInfo {
	return &JinritemaiStoreInfo{
		ClientKey:         clientKey,
		OpenId:            openId,
		ProductId:         productId,
		ProductName:       productName,
		ProductImg:        productImg,
		ProductPrice:      productPrice,
		CommissionType:    commissionType,
		CommissionRatio:   commissionRatio,
		PromotionId:       promotionId,
		PromotionType:     promotionType,
		ColonelActivityId: colonelActivityId,
	}
}

func (jsi *JinritemaiStoreInfo) SetClientKey(ctx context.Context, clientKey string) {
	jsi.ClientKey = clientKey
}

func (jsi *JinritemaiStoreInfo) SetOpenId(ctx context.Context, openId string) {
	jsi.OpenId = openId
}

func (jsi *JinritemaiStoreInfo) SetProductId(ctx context.Context, productId string) {
	jsi.ProductId = productId
}

func (jsi *JinritemaiStoreInfo) SetProductName(ctx context.Context, productName string) {
	jsi.ProductName = productName
}

func (jsi *JinritemaiStoreInfo) SetProductImg(ctx context.Context, productImg string) {
	jsi.ProductImg = productImg
}

func (jsi *JinritemaiStoreInfo) SetProductPrice(ctx context.Context, productPrice float32) {
	jsi.ProductPrice = productPrice
}

func (jsi *JinritemaiStoreInfo) SetCommissionType(ctx context.Context, commissionType uint8) {
	jsi.CommissionType = commissionType
}

func (jsi *JinritemaiStoreInfo) SetCommissionTypeName(ctx context.Context) {
	switch jsi.CommissionType {
	case 0:
		jsi.CommissionTypeName = "未定义（异常）"
	case 1:
		jsi.CommissionTypeName = "专属团长"
	case 2:
		jsi.CommissionTypeName = "普通佣金"
	case 3:
		jsi.CommissionTypeName = "定向佣金"
	case 4:
		jsi.CommissionTypeName = "提报活动"
	case 5:
		jsi.CommissionTypeName = "招募佣金"
	}
}

func (jsi *JinritemaiStoreInfo) SetCommissionRatio(ctx context.Context, commissionRatio uint8) {
	jsi.CommissionRatio = commissionRatio
}

func (jsi *JinritemaiStoreInfo) SetPromotionId(ctx context.Context, promotionId uint64) {
	jsi.PromotionId = promotionId
}

func (jsi *JinritemaiStoreInfo) SetPromotionType(ctx context.Context, promotionType uint8) {
	jsi.PromotionType = promotionType
}

func (jsi *JinritemaiStoreInfo) SetPromotionTypeName(ctx context.Context) {
	switch jsi.PromotionType {
	case 0:
		jsi.PromotionTypeName = "非团长"
	case 1:
		jsi.CommissionTypeName = "团长"
	}
}

func (jsi *JinritemaiStoreInfo) SetColonelActivityId(ctx context.Context, colonelActivityId uint64) {
	jsi.ColonelActivityId = colonelActivityId
}

func (jsi *JinritemaiStoreInfo) SetCreateTime(ctx context.Context, createTime time.Time) {
	jsi.CreateTime = createTime
}

func (jsi *JinritemaiStoreInfo) SetUpdateTime(ctx context.Context, updateTime time.Time) {
	jsi.UpdateTime = updateTime
}

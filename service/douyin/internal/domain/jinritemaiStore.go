package domain

import (
	"context"
	"encoding/json"
	"strings"
)

type Store struct {
	ClientKey string `json:"clientKey"`
	OpenId    string `json:"openId"`
	ProductId string `json:"productId"`
}

type JinritemaiStore struct {
	ClientKey         string
	OpenId            string
	ProductId         int64
	PromotionId       int64
	Title             string
	Cover             string
	PromotionType     int64
	Price             int64
	CosType           int64
	CosRatio          float64
	ColonelActivityId int64
	HideStatus        bool
	Store             string
	Stores            []*Store
}

type JinritemaiStoreList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*JinritemaiStoreInfo
}

type AddStoreMessage struct {
	ProductName string
	AwemeName   string
	Content     string
}

func NewJinritemaiStore(ctx context.Context, productId, promotionId, promotionType, price, cosType, colonelActivityId int64, cosRatio float64, hideStatus bool, clientKey, openId, title, cover string) *JinritemaiStore {
	return &JinritemaiStore{
		ClientKey:         clientKey,
		OpenId:            openId,
		ProductId:         productId,
		PromotionId:       promotionId,
		Title:             title,
		Cover:             cover,
		PromotionType:     promotionType,
		Price:             price,
		CosType:           cosType,
		CosRatio:          cosRatio,
		ColonelActivityId: colonelActivityId,
		HideStatus:        hideStatus,
	}
}

func (js *JinritemaiStore) SetClientKey(ctx context.Context, clientKey string) {
	js.ClientKey = clientKey
}

func (js *JinritemaiStore) SetOpenId(ctx context.Context, openId string) {
	js.OpenId = openId
}

func (js *JinritemaiStore) SetProductId(ctx context.Context, productId int64) {
	js.ProductId = productId
}

func (js *JinritemaiStore) SetPromotionId(ctx context.Context, promotionId int64) {
	js.PromotionId = promotionId
}

func (js *JinritemaiStore) SetTitle(ctx context.Context, title string) {
	js.Title = title
}

func (js *JinritemaiStore) SetCover(ctx context.Context, cover string) {
	js.Cover = cover
}

func (js *JinritemaiStore) SetPromotionType(ctx context.Context, promotionType int64) {
	js.PromotionType = promotionType
}

func (js *JinritemaiStore) SetPrice(ctx context.Context, price int64) {
	js.Price = price
}

func (js *JinritemaiStore) SetCosType(ctx context.Context, cosType int64) {
	js.CosType = cosType
}

func (js *JinritemaiStore) SetCosRatio(ctx context.Context, cosRatio float64) {
	js.CosRatio = cosRatio
}

func (js *JinritemaiStore) SetColonelActivityId(ctx context.Context, colonelActivityId int64) {
	js.ColonelActivityId = colonelActivityId
}

func (js *JinritemaiStore) SetHideStatus(ctx context.Context, hideStatus bool) {
	js.HideStatus = hideStatus
}

func (js *JinritemaiStore) SetStore(ctx context.Context, store string) {
	js.Store = store
}

func (js *JinritemaiStore) VerifyStores(ctx context.Context) bool {
	var stores []*Store

	if err := json.Unmarshal([]byte(js.Store), &stores); err != nil {
		return false
	}

	for _, store := range stores {
		if len(store.ClientKey) == 0 {
			return false
		}

		if len(store.OpenId) == 0 {
			return false
		}

		if len(store.ProductId) == 0 {
			return false
		}

		productIds := strings.Split(store.ProductId, ",")

		if len(productIds) == 0 {
			return false
		}
	}

	js.Stores = stores

	return true
}

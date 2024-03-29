package domain

import (
	"context"
	"time"
)

type ImgUrl struct {
	ImgUrl string `json:"img_url" bson:"img_url"`
}

type QianchuanProduct struct {
	Id                  uint64
	AdvertiserId        uint64
	Name                string
	Img                 string
	CategoryName        string
	DiscountPrice       float64
	DiscountLowerPrice  float64
	DiscountHigherPrice float64
	ImgList             []*ImgUrl
	Inventory           uint64
	MarketPrice         float64
	ProductRate         float64
	SaleTime            string
	Tags                string
	CreateTime          time.Time
	UpdateTime          time.Time
}

func NewQianchuanProduct(ctx context.Context, id, advertiserId, inventory uint64, discountPrice, discountLowerPrice, discountHigherPrice, marketPrice, productRate float64, name, img, categoryName, saleTime, tags string, imgList []*ImgUrl) *QianchuanProduct {
	return &QianchuanProduct{
		Id:                  id,
		AdvertiserId:        advertiserId,
		Name:                name,
		Img:                 img,
		CategoryName:        categoryName,
		DiscountPrice:       discountPrice,
		DiscountLowerPrice:  discountLowerPrice,
		DiscountHigherPrice: discountHigherPrice,
		ImgList:             imgList,
		Inventory:           inventory,
		MarketPrice:         marketPrice,
		ProductRate:         productRate,
		SaleTime:            saleTime,
		Tags:                tags,
	}
}

func (qp *QianchuanProduct) SetId(ctx context.Context, id uint64) {
	qp.Id = id
}

func (qp *QianchuanProduct) SetAdvertiserId(ctx context.Context, advertiserId uint64) {
	qp.AdvertiserId = advertiserId
}

func (qp *QianchuanProduct) SetName(ctx context.Context, name string) {
	qp.Name = name
}

func (qp *QianchuanProduct) SetImg(ctx context.Context, img string) {
	qp.Img = img
}

func (qp *QianchuanProduct) SetCategoryName(ctx context.Context, categoryName string) {
	qp.CategoryName = categoryName
}

func (qp *QianchuanProduct) SetDiscountPrice(ctx context.Context, discountPrice float64) {
	qp.DiscountPrice = discountPrice
}

func (qp *QianchuanProduct) SetDiscountLowerPrice(ctx context.Context, discountLowerPrice float64) {
	qp.DiscountLowerPrice = discountLowerPrice
}

func (qp *QianchuanProduct) SetDiscountHigherPrice(ctx context.Context, discountHigherPrice float64) {
	qp.DiscountHigherPrice = discountHigherPrice
}

func (qp *QianchuanProduct) SetImgList(ctx context.Context, imgList []*ImgUrl) {
	qp.ImgList = imgList
}

func (qp *QianchuanProduct) SetInventory(ctx context.Context, inventory uint64) {
	qp.Inventory = inventory
}

func (qp *QianchuanProduct) SetMarketPrice(ctx context.Context, marketPrice float64) {
	qp.MarketPrice = marketPrice
}

func (qp *QianchuanProduct) SetProductRate(ctx context.Context, productRate float64) {
	qp.ProductRate = productRate
}

func (qp *QianchuanProduct) SetSaleTime(ctx context.Context, saleTime string) {
	qp.SaleTime = saleTime
}

func (qp *QianchuanProduct) SetTags(ctx context.Context, tags string) {
	qp.Tags = tags
}

func (qp *QianchuanProduct) SetUpdateTime(ctx context.Context) {
	qp.UpdateTime = time.Now()
}

func (qp *QianchuanProduct) SetCreateTime(ctx context.Context) {
	qp.CreateTime = time.Now()
}

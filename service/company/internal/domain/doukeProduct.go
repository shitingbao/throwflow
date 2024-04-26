package domain

import (
	"context"
	"time"
)

type DoukeProductGorm struct {
	ProductOutId    uint64
	ProductType     uint8
	ProductStatus   uint8
	ProductName     string
	ProductImg      string
	ProductPrice    string
	IndustryId      uint64
	CategoryId      uint64
	SubCategoryId   uint64
	ShopName        string
	ShopScore       float64
	TotalSale       uint64
	CommissionRatio float32
	CreateTime      time.Time
	UpdateTime      time.Time
}

func NewDoukeProductGorm(ctx context.Context, productOutId, industryId, categoryId, subCategoryId, totalSale uint64, productType, productStatus uint8, commissionRatio float32, productName, productImg, productPrice, shopName string) *DoukeProductGorm {
	return &DoukeProductGorm{
		ProductOutId:    productOutId,
		ProductType:     productType,
		ProductStatus:   productStatus,
		ProductName:     productName,
		ProductImg:      productImg,
		ProductPrice:    productPrice,
		IndustryId:      industryId,
		CategoryId:      categoryId,
		SubCategoryId:   subCategoryId,
		ShopName:        shopName,
		TotalSale:       totalSale,
		CommissionRatio: commissionRatio,
	}
}

func (dpg *DoukeProductGorm) SetProductOutId(ctx context.Context, productOutId uint64) {
	dpg.ProductOutId = productOutId
}

func (dpg *DoukeProductGorm) SetProductType(ctx context.Context, productType uint8) {
	dpg.ProductType = productType
}

func (dpg *DoukeProductGorm) SetProductStatus(ctx context.Context, productStatus uint8) {
	dpg.ProductStatus = productStatus
}

func (dpg *DoukeProductGorm) SetProductName(ctx context.Context, productName string) {
	dpg.ProductName = productName
}

func (dpg *DoukeProductGorm) SetProductImg(ctx context.Context, productImg string) {
	dpg.ProductImg = productImg
}

func (dpg *DoukeProductGorm) SetProductPrice(ctx context.Context, productPrice string) {
	dpg.ProductPrice = productPrice
}

func (dpg *DoukeProductGorm) SetIndustryId(ctx context.Context, industryId uint64) {
	dpg.IndustryId = industryId
}

func (dpg *DoukeProductGorm) SetCategoryId(ctx context.Context, categoryId uint64) {
	dpg.CategoryId = categoryId
}

func (dpg *DoukeProductGorm) SetSubCategoryId(ctx context.Context, subCategoryId uint64) {
	dpg.SubCategoryId = subCategoryId
}

func (dpg *DoukeProductGorm) SetShopName(ctx context.Context, shopName string) {
	dpg.ShopName = shopName
}

func (dpg *DoukeProductGorm) SetShopScore(ctx context.Context, shopScore float64) {
	dpg.ShopScore = shopScore
}

func (dpg *DoukeProductGorm) SetTotalSale(ctx context.Context, totalSale uint64) {
	dpg.TotalSale = totalSale
}

func (dpg *DoukeProductGorm) SetCommissionRatio(ctx context.Context, commissionRatio float32) {
	dpg.CommissionRatio = commissionRatio
}

func (dpg *DoukeProductGorm) SetCreateTime(ctx context.Context) {
	dpg.CreateTime = time.Now()
}

func (dpg *DoukeProductGorm) SetUpdateTime(ctx context.Context) {
	dpg.UpdateTime = time.Now()
}

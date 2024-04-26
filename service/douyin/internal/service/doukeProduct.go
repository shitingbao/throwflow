package service

import (
	"context"
	v1 "douyin/api/douyin/v1"
	"douyin/internal/biz"
	"douyin/internal/pkg/csj"
	"douyin/internal/pkg/tool"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func (ds *DouyinService) ListDoukeProducts(ctx context.Context, in *v1.ListDoukeProductsRequest) (*v1.ListDoukeProductsReply, error) {
	if in.PageSize > csj.PageSize20 {
		return nil, biz.DouyinValidatorError
	}

	products, err := ds.dpuc.ListDoukeProducts(ctx, in.PageNum, in.PageSize, in.CosRatioMin, in.IndustryId, in.CategoryId, in.SubCategoryId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListDoukeProductsReply_Product, 0)

	for _, product := range products.Data.Products {
		list = append(list, &v1.ListDoukeProductsReply_Product{
			ProductOutId:    product.ProductId,
			ProductName:     product.Title,
			ProductImg:      product.Cover,
			ProductPrice:    fmt.Sprintf("%.2f", tool.Decimal(float64(product.Price)/float64(100), 2)),
			IndustryId:      product.FirstCid,
			CategoryId:      product.SecondCid,
			SubCategoryId:   product.ThirdCid,
			ShopName:        product.ShopName,
			TotalSale:       product.Sales,
			CommissionRatio: product.CosRatio / 100,
		})
	}

	totalPage := uint64(math.Ceil(float64(products.Data.Total) / float64(csj.PageSize20)))

	return &v1.ListDoukeProductsReply{
		Code: 200,
		Data: &v1.ListDoukeProductsReply_Data{
			PageNum:   in.PageNum,
			PageSize:  in.PageSize,
			Total:     products.Data.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ds *DouyinService) ListDoukeProductByProductIds(ctx context.Context, in *v1.ListDoukeProductByProductIdsRequest) (*v1.ListDoukeProductByProductIdsReply, error) {
	products, err := ds.dpuc.ListDoukeProductByProductIds(ctx, in.ProductIds)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListDoukeProductByProductIdsReply_Product, 0)

	for _, product := range products.Data.Products {
		shopScore, _ := strconv.ParseFloat(product.ShopTotalScore.ShopScore.Score, 10)

		list = append(list, &v1.ListDoukeProductByProductIdsReply_Product{
			ProductOutId:    product.ProductId,
			ProductName:     product.Title,
			ProductImg:      strings.Join(product.Imgs, ","),
			ProductPrice:    fmt.Sprintf("%.2f", tool.Decimal(float64(product.Price)/float64(100), 2)),
			IndustryId:      product.FirstCid,
			CategoryId:      product.SecondCid,
			SubCategoryId:   product.ThirdCid,
			ShopName:        product.ShopName,
			ShopScore:       shopScore,
			TotalSale:       product.Sales,
			CommissionRatio: product.CosRatio / 100,
		})
	}

	return &v1.ListDoukeProductByProductIdsReply{
		Code: 200,
		Data: &v1.ListDoukeProductByProductIdsReply_Data{
			List: list,
		},
	}, nil
}

func (ds *DouyinService) ListCategoryDoukeProducts(ctx context.Context, in *v1.ListCategoryDoukeProductsRequest) (*v1.ListCategoryDoukeProductsReply, error) {
	categories, err := ds.dpuc.ListCategoryDoukeProducts(ctx, in.ParentId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCategoryDoukeProductsReply_Category, 0)

	for _, category := range categories.Data.CategoryList {
		if category.Id != in.ParentId {
			list = append(list, &v1.ListCategoryDoukeProductsReply_Category{
				CategoryId:   category.Id,
				CategoryName: category.Name,
			})
		}
	}

	return &v1.ListCategoryDoukeProductsReply{
		Code: 200,
		Data: &v1.ListCategoryDoukeProductsReply_Data{
			List: list,
		},
	}, nil
}

func (ds *DouyinService) CreateShareDoukeProducts(ctx context.Context, in *v1.CreateShareDoukeProductsRequest) (*v1.CreateShareDoukeProductsReply, error) {
	productShare, err := ds.dpuc.CreateShareDoukeProducts(ctx, in.ProductUrl, in.ExternalInfo)

	if err != nil {
		return nil, err
	}

	return &v1.CreateShareDoukeProductsReply{
		Code: 200,
		Data: &v1.CreateShareDoukeProductsReply_Data{
			DyPassword: productShare.Data.DyPassword,
		},
	}, nil
}

func (ds *DouyinService) ParseDoukeProducts(ctx context.Context, in *v1.ParseDoukeProductsRequest) (*v1.ParseDoukeProductsReply, error) {
	productShare, err := ds.dpuc.ParseDoukeProducts(ctx, in.Content)

	if err != nil {
		return nil, err
	}

	return &v1.ParseDoukeProductsReply{
		Code: 200,
		Data: &v1.ParseDoukeProductsReply_Data{
			ProductId: productShare.Data.ProductInfo.ProductId,
		},
	}, nil
}

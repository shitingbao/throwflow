package service

import (
	"context"
	v1 "douyin/api/douyin/v1"
	"douyin/internal/pkg/tool"
	"fmt"
	"strconv"
	"strings"
)

func (ds *DouyinService) GetDoukeProducts(ctx context.Context, in *v1.GetDoukeProductsRequest) (*v1.GetDoukeProductsReply, error) {
	product, err := ds.dpuc.GetDoukeProducts(ctx, in.ProductId)

	if err != nil {
		return nil, err
	}

	shopScore, _ := strconv.ParseFloat(product.Data.Products[0].ShopTotalScore.ShopScore.Score, 10)

	return &v1.GetDoukeProductsReply{
		Code: 200,
		Data: &v1.GetDoukeProductsReply_Data{
			ProductOutId:    product.Data.Products[0].ProductId,
			ProductName:     product.Data.Products[0].Title,
			ProductImg:      strings.Join(product.Data.Products[0].Imgs, ","),
			ProductPrice:    fmt.Sprintf("%.2f", tool.Decimal(float64(product.Data.Products[0].Price)/float64(100), 2)),
			IndustryId:      product.Data.Products[0].FirstCid,
			CategoryId:      product.Data.Products[0].SecondCid,
			SubCategoryId:   product.Data.Products[0].ThirdCid,
			ShopName:        product.Data.Products[0].ShopName,
			ShopScore:       shopScore,
			TotalSale:       product.Data.Products[0].Sales,
			CommissionRatio: product.Data.Products[0].CosRatio / 100,
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

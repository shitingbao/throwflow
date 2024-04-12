package service

import (
	"context"
	v1 "interface/api/interface/v1"
	"interface/internal/biz"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (is *InterfaceService) ListProducts(ctx context.Context, in *v1.ListProductsRequest) (*v1.ListProductsReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "product"); err != nil {
		return nil, err
	}

	products, err := is.prouc.ListProducts(ctx, in.PageNum, in.PageSize, in.IndustryId, in.CategoryId, in.SubCategoryId, in.ProductStatus, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListProductsReply_CompanyProduct, 0)

	for _, product := range products.Data.List {
		productImgs := make([]*v1.ListProductsReply_ProductImg, 0)
		materialOutUrls := make([]*v1.ListProductsReply_MaterialOutUrl, 0)
		commissions := make([]*v1.ListProductsReply_Commission, 0)

		for _, productImg := range product.ProductImgs {
			productImgs = append(productImgs, &v1.ListProductsReply_ProductImg{
				ProductImg: productImg.ProductImg,
			})
		}

		for _, materialOutUrl := range product.MaterialOutUrls {
			materialOutUrls = append(materialOutUrls, &v1.ListProductsReply_MaterialOutUrl{
				MaterialOutUrl: materialOutUrl.MaterialOutUrl,
			})
		}

		for _, commission := range product.Commissions {
			commissions = append(commissions, &v1.ListProductsReply_Commission{
				CommissionRatio:  float64(commission.CommissionRatio),
				ServiceRatio:     float64(commission.ServiceRatio),
				CommissionOutUrl: commission.CommissionOutUrl,
			})
		}

		list = append(list, &v1.ListProductsReply_CompanyProduct{
			ProductId:            product.ProductId,
			ProductOutId:         product.ProductOutId,
			ProductType:          product.ProductType,
			ProductStatus:        product.ProductStatus,
			ProductName:          product.ProductName,
			ProductImgs:          productImgs,
			ProductPrice:         product.ProductPrice,
			ProductUrl:           product.ProductUrl,
			IsTop:                product.IsTop,
			MaterialOutUrls:      materialOutUrls,
			SampleThresholdType:  product.SampleThresholdType,
			SampleThresholdValue: product.SampleThresholdValue,
			Commissions:          commissions,
			InvestmentRatio:      product.InvestmentRatio,
			ForbidReason:         product.ForbidReason,
			IsTask:               product.IsTask,
		})
	}

	return &v1.ListProductsReply{
		Code: 200,
		Data: &v1.ListProductsReply_Data{
			PageNum:   products.Data.PageNum,
			PageSize:  products.Data.PageSize,
			Total:     products.Data.Total,
			TotalPage: products.Data.TotalPage,
			List:      list,
		},
	}, nil
}

func (is *InterfaceService) ListMiniProducts(ctx context.Context, in *v1.ListMiniProductsRequest) (*v1.ListMiniProductsReply, error) {
	products, err := is.prouc.ListMiniProducts(ctx, in.PageNum, in.PageSize, in.IndustryId, in.CategoryId, in.SubCategoryId, in.IsInvestment, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListMiniProductsReply_CompanyProduct, 0)

	for _, product := range products.Data.List {
		list = append(list, &v1.ListMiniProductsReply_CompanyProduct{
			ProductId:             product.ProductId,
			ProductName:           product.ProductName,
			ProductImg:            product.ProductImg,
			ProductPrice:          product.ProductPrice,
			IsTop:                 product.IsTop,
			PureCommission:        product.PureCommission,
			PureServiceCommission: product.PureServiceCommission,
			CommonCommission:      product.CommonCommission,
			InvestmentRatio:       product.InvestmentRatio,
			IsHot:                 product.IsHot,
			TotalSale:             product.TotalSale,
			IsTask:                product.IsTask,
		})
	}

	return &v1.ListMiniProductsReply{
		Code: 200,
		Data: &v1.ListMiniProductsReply_Data{
			PageNum:   products.Data.PageNum,
			PageSize:  products.Data.PageSize,
			Total:     products.Data.Total,
			TotalPage: products.Data.TotalPage,
			List:      list,
		},
	}, nil
}

func (is *InterfaceService) ListCategoryProducts(ctx context.Context, in *emptypb.Empty) (*v1.ListCategoryProductsReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "product"); err != nil {
		return nil, err
	}

	productCategories, err := is.prouc.ListCategorys(ctx)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCategoryProductsReply_Category, 0)

	for _, productCategory := range productCategories.Data.Category {
		childList := make([]*v1.ListCategoryProductsReply_ChildCategory, 0)

		for _, lchildList := range productCategory.ChildList {
			subChildList := make([]*v1.ListCategoryProductsReply_SubChildCategory, 0)

			for _, llchildList := range lchildList.ChildList {
				subChildList = append(subChildList, &v1.ListCategoryProductsReply_SubChildCategory{
					Key:   llchildList.Key,
					Value: llchildList.Value,
				})
			}

			childList = append(childList, &v1.ListCategoryProductsReply_ChildCategory{
				Key:       lchildList.Key,
				Value:     lchildList.Value,
				ChildList: subChildList,
			})
		}

		list = append(list, &v1.ListCategoryProductsReply_Category{
			Key:       productCategory.Key,
			Value:     productCategory.Value,
			ChildList: childList,
		})
	}

	return &v1.ListCategoryProductsReply{
		Code: 200,
		Data: &v1.ListCategoryProductsReply_Data{
			Category: list,
		},
	}, nil
}

func (is *InterfaceService) ListMiniCategorys(ctx context.Context, in *emptypb.Empty) (*v1.ListMiniCategorysReply, error) {
	productCategories, err := is.prouc.ListMiniCategorys(ctx)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListMiniCategorysReply_Category, 0)

	for _, productCategory := range productCategories.Data.Category {
		childList := make([]*v1.ListMiniCategorysReply_ChildCategory, 0)

		for _, lchildList := range productCategory.ChildList {
			subChildList := make([]*v1.ListMiniCategorysReply_SubChildCategory, 0)

			for _, llchildList := range lchildList.ChildList {
				subChildList = append(subChildList, &v1.ListMiniCategorysReply_SubChildCategory{
					Key:   llchildList.Key,
					Value: llchildList.Value,
				})
			}

			childList = append(childList, &v1.ListMiniCategorysReply_ChildCategory{
				Key:       lchildList.Key,
				Value:     lchildList.Value,
				ChildList: subChildList,
			})
		}

		list = append(list, &v1.ListMiniCategorysReply_Category{
			Key:       productCategory.Key,
			Value:     productCategory.Value,
			ChildList: childList,
		})
	}

	return &v1.ListMiniCategorysReply{
		Code: 200,
		Data: &v1.ListMiniCategorysReply_Data{
			Category: list,
		},
	}, nil
}

func (is *InterfaceService) ListCompanyTaskProducts(ctx context.Context, in *v1.ListCompanyTaskProductsRequest) (*v1.ListCompanyTaskProductsReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "product"); err != nil {
		return nil, err
	}

	products, err := is.prouc.ListCompanyTaskProducts(ctx, in.PageNum, in.PageSize, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyTaskProductsReply_CompanyProduct, 0)

	for _, product := range products.Data.List {
		productImgs := make([]*v1.ListCompanyTaskProductsReply_ProductImg, 0)
		materialOutUrls := make([]*v1.ListCompanyTaskProductsReply_MaterialOutUrl, 0)
		commissions := make([]*v1.ListCompanyTaskProductsReply_Commission, 0)

		for _, productImg := range product.ProductImgs {
			productImgs = append(productImgs, &v1.ListCompanyTaskProductsReply_ProductImg{
				ProductImg: productImg.ProductImg,
			})
		}

		for _, materialOutUrl := range product.MaterialOutUrls {
			materialOutUrls = append(materialOutUrls, &v1.ListCompanyTaskProductsReply_MaterialOutUrl{
				MaterialOutUrl: materialOutUrl.MaterialOutUrl,
			})
		}

		for _, commission := range product.Commissions {
			commissions = append(commissions, &v1.ListCompanyTaskProductsReply_Commission{
				CommissionRatio:  float64(commission.CommissionRatio),
				ServiceRatio:     float64(commission.ServiceRatio),
				CommissionOutUrl: commission.CommissionOutUrl,
			})
		}

		list = append(list, &v1.ListCompanyTaskProductsReply_CompanyProduct{
			ProductId:            product.ProductId,
			ProductOutId:         product.ProductOutId,
			ProductType:          product.ProductType,
			ProductStatus:        product.ProductStatus,
			ProductName:          product.ProductName,
			ProductImgs:          productImgs,
			ProductPrice:         product.ProductPrice,
			ProductUrl:           product.ProductUrl,
			IsTop:                product.IsTop,
			MaterialOutUrls:      materialOutUrls,
			SampleThresholdType:  product.SampleThresholdType,
			SampleThresholdValue: product.SampleThresholdValue,
			Commissions:          commissions,
			InvestmentRatio:      product.InvestmentRatio,
			ForbidReason:         product.ForbidReason,
			IsTask:               product.IsTask,
		})
	}

	return &v1.ListCompanyTaskProductsReply{
		Code: 200,
		Data: &v1.ListCompanyTaskProductsReply_Data{
			PageNum:   products.Data.PageNum,
			PageSize:  products.Data.PageSize,
			Total:     products.Data.Total,
			TotalPage: products.Data.TotalPage,
			List:      list,
		},
	}, nil
}

func (is *InterfaceService) StatisticsMiniProducts(ctx context.Context, in *v1.StatisticsMiniProductsRequest) (*v1.StatisticsMiniProductsReply, error) {
	statistics, err := is.prouc.StatisticsMiniProducts(ctx, in.IndustryId, in.CategoryId, in.SubCategoryId, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsMiniProductsReply_Statistic, 0)

	for _, statistic := range statistics.Data.Statistics {
		list = append(list, &v1.StatisticsMiniProductsReply_Statistic{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsMiniProductsReply{
		Code: 200,
		Data: &v1.StatisticsMiniProductsReply_Data{
			Statistics: list,
		},
	}, nil
}

func (is *InterfaceService) GetUploadIdProducts(ctx context.Context, in *v1.GetUploadIdProductsRequest) (*v1.GetUploadIdProductsReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "product"); err != nil {
		return nil, err
	}

	uploadId, err := is.prouc.GetUploadIdProducts(ctx, in.Suffix)

	if err != nil {
		return nil, err
	}

	return &v1.GetUploadIdProductsReply{
		Code: 200,
		Data: &v1.GetUploadIdProductsReply_Data{
			UploadId: uploadId.Data.UploadId,
		},
	}, nil
}

func (is *InterfaceService) GetMiniProducts(ctx context.Context, in *v1.GetMiniProductsRequest) (*v1.GetMiniProductsReply, error) {
	product, err := is.prouc.GetMiniProducts(ctx, in.ProductId)

	if err != nil {
		return nil, err
	}

	productImgs := make([]*v1.GetMiniProductsReply_ProductImg, 0)
	productDetailImgs := make([]*v1.GetMiniProductsReply_ProductDetailImg, 0)
	awemes := make([]*v1.GetMiniProductsReply_Aweme, 0)

	for _, productImg := range product.Data.ProductImgs {
		productImgs = append(productImgs, &v1.GetMiniProductsReply_ProductImg{
			ProductImg: productImg.ProductImg,
		})
	}

	for _, productDetailImg := range product.Data.ProductDetailImgs {
		productDetailImgs = append(productDetailImgs, &v1.GetMiniProductsReply_ProductDetailImg{
			ProductDetailImg: productDetailImg.ProductDetailImg,
		})
	}

	for _, aweme := range product.Data.Awemes {
		awemes = append(awemes, &v1.GetMiniProductsReply_Aweme{
			AccountId:    aweme.AccountId,
			Nickname:     aweme.Nickname,
			Avatar:       aweme.Avatar,
			AvatarLarger: aweme.AvatarLarger,
		})
	}

	return &v1.GetMiniProductsReply{
		Code: 200,
		Data: &v1.GetMiniProductsReply_Data{
			ProductId:             product.Data.ProductId,
			ProductName:           product.Data.ProductName,
			ProductImgs:           productImgs,
			ProductDetailImgs:     productDetailImgs,
			ProductPrice:          product.Data.ProductPrice,
			ShopName:              product.Data.ShopName,
			ShopLogo:              product.Data.ShopLogo,
			ShopScore:             product.Data.ShopScore,
			IsTop:                 product.Data.IsTop,
			SampleThresholdType:   product.Data.SampleThresholdType,
			SampleThresholdValue:  product.Data.SampleThresholdValue,
			PureCommission:        product.Data.PureCommission,
			PureServiceCommission: product.Data.PureServiceCommission,
			CommonCommission:      product.Data.CommonCommission,
			IsHot:                 product.Data.IsHot,
			ProductUrl:            product.Data.ProductUrl,
			TotalSale:             product.Data.TotalSale,
			InvestmentRatio:       product.Data.InvestmentRatio,
			Awemes:                awemes,
		},
	}, nil
}

func (is *InterfaceService) GetMiniProductShareProducts(ctx context.Context, in *v1.GetMiniProductShareProductsRequest) (*v1.GetMiniProductShareProductsReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	productShare, err := is.prouc.GetMiniProductShareProducts(ctx, userInfo.Data.UserId, in.ProductId)

	if err != nil {
		return nil, err
	}

	return &v1.GetMiniProductShareProductsReply{
		Code: 200,
		Data: &v1.GetMiniProductShareProductsReply_Data{
			DyPassword: productShare.Data.DyPassword,
		},
	}, nil
}

func (is *InterfaceService) CreateProducts(ctx context.Context, in *v1.CreateProductsRequest) (*v1.CreateProductsReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	if _, err := is.verifyLogin(ctx, true, false, "product"); err != nil {
		return nil, err
	}

	product, err := is.prouc.CreateProducts(ctx, in.Commission)

	if err != nil {
		return nil, err
	}

	productImgs := make([]*v1.CreateProductsReply_ProductImg, 0)
	materialOutUrls := make([]*v1.CreateProductsReply_MaterialOutUrl, 0)
	commissions := make([]*v1.CreateProductsReply_Commission, 0)

	for _, productImg := range product.Data.ProductImgs {
		productImgs = append(productImgs, &v1.CreateProductsReply_ProductImg{
			ProductImg: productImg.ProductImg,
		})
	}

	for _, materialOutUrl := range product.Data.MaterialOutUrls {
		materialOutUrls = append(materialOutUrls, &v1.CreateProductsReply_MaterialOutUrl{
			MaterialOutUrl: materialOutUrl.MaterialOutUrl,
		})
	}

	for _, commission := range product.Data.Commissions {
		commissions = append(commissions, &v1.CreateProductsReply_Commission{
			CommissionRatio:  float64(commission.CommissionRatio),
			ServiceRatio:     float64(commission.ServiceRatio),
			CommissionOutUrl: commission.CommissionOutUrl,
		})
	}

	return &v1.CreateProductsReply{
		Code: 200,
		Data: &v1.CreateProductsReply_Data{
			ProductId:            product.Data.ProductId,
			ProductOutId:         product.Data.ProductOutId,
			ProductType:          product.Data.ProductType,
			ProductStatus:        product.Data.ProductStatus,
			ProductName:          product.Data.ProductName,
			ProductImgs:          productImgs,
			ProductPrice:         product.Data.ProductPrice,
			ProductUrl:           product.Data.ProductUrl,
			IsTop:                product.Data.IsTop,
			MaterialOutUrls:      materialOutUrls,
			SampleThresholdType:  product.Data.SampleThresholdType,
			SampleThresholdValue: product.Data.SampleThresholdValue,
			Commissions:          commissions,
			InvestmentRatio:      product.Data.InvestmentRatio,
			ForbidReason:         product.Data.ForbidReason,
			IsTask:               product.Data.IsTask,
		},
	}, nil
}

func (is *InterfaceService) UpdateCommissionProducts(ctx context.Context, in *v1.UpdateCommissionProductsRequest) (*v1.UpdateCommissionProductsReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	if _, err := is.verifyLogin(ctx, true, false, "product"); err != nil {
		return nil, err
	}

	product, err := is.prouc.UpdateCommissionProducts(ctx, in.ProductId, in.Commission)

	if err != nil {
		return nil, err
	}

	productImgs := make([]*v1.UpdateCommissionProductsReply_ProductImg, 0)
	materialOutUrls := make([]*v1.UpdateCommissionProductsReply_MaterialOutUrl, 0)
	commissions := make([]*v1.UpdateCommissionProductsReply_Commission, 0)

	for _, productImg := range product.Data.ProductImgs {
		productImgs = append(productImgs, &v1.UpdateCommissionProductsReply_ProductImg{
			ProductImg: productImg.ProductImg,
		})
	}

	for _, materialOutUrl := range product.Data.MaterialOutUrls {
		materialOutUrls = append(materialOutUrls, &v1.UpdateCommissionProductsReply_MaterialOutUrl{
			MaterialOutUrl: materialOutUrl.MaterialOutUrl,
		})
	}

	for _, commission := range product.Data.Commissions {
		commissions = append(commissions, &v1.UpdateCommissionProductsReply_Commission{
			CommissionRatio:  float64(commission.CommissionRatio),
			ServiceRatio:     float64(commission.ServiceRatio),
			CommissionOutUrl: commission.CommissionOutUrl,
		})
	}

	return &v1.UpdateCommissionProductsReply{
		Code: 200,
		Data: &v1.UpdateCommissionProductsReply_Data{
			ProductId:            product.Data.ProductId,
			ProductOutId:         product.Data.ProductOutId,
			ProductType:          product.Data.ProductType,
			ProductStatus:        product.Data.ProductStatus,
			ProductName:          product.Data.ProductName,
			ProductImgs:          productImgs,
			ProductPrice:         product.Data.ProductPrice,
			ProductUrl:           product.Data.ProductUrl,
			IsTop:                product.Data.IsTop,
			MaterialOutUrls:      materialOutUrls,
			SampleThresholdType:  product.Data.SampleThresholdType,
			SampleThresholdValue: product.Data.SampleThresholdValue,
			Commissions:          commissions,
			InvestmentRatio:      product.Data.InvestmentRatio,
			ForbidReason:         product.Data.ForbidReason,
			IsTask:               product.Data.IsTask,
		},
	}, nil
}

func (is *InterfaceService) UpdateStatusProducts(ctx context.Context, in *v1.UpdateStatusProductsRequest) (*v1.UpdateStatusProductsReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "product"); err != nil {
		return nil, err
	}

	product, err := is.prouc.UpdateStatusProducts(ctx, in.ProductId, in.Status)

	if err != nil {
		return nil, err
	}

	productImgs := make([]*v1.UpdateStatusProductsReply_ProductImg, 0)
	materialOutUrls := make([]*v1.UpdateStatusProductsReply_MaterialOutUrl, 0)
	commissions := make([]*v1.UpdateStatusProductsReply_Commission, 0)

	for _, productImg := range product.Data.ProductImgs {
		productImgs = append(productImgs, &v1.UpdateStatusProductsReply_ProductImg{
			ProductImg: productImg.ProductImg,
		})
	}

	for _, materialOutUrl := range product.Data.MaterialOutUrls {
		materialOutUrls = append(materialOutUrls, &v1.UpdateStatusProductsReply_MaterialOutUrl{
			MaterialOutUrl: materialOutUrl.MaterialOutUrl,
		})
	}

	for _, commission := range product.Data.Commissions {
		commissions = append(commissions, &v1.UpdateStatusProductsReply_Commission{
			CommissionRatio:  float64(commission.CommissionRatio),
			ServiceRatio:     float64(commission.ServiceRatio),
			CommissionOutUrl: commission.CommissionOutUrl,
		})
	}

	return &v1.UpdateStatusProductsReply{
		Code: 200,
		Data: &v1.UpdateStatusProductsReply_Data{
			ProductId:            product.Data.ProductId,
			ProductOutId:         product.Data.ProductOutId,
			ProductType:          product.Data.ProductType,
			ProductStatus:        product.Data.ProductStatus,
			ProductName:          product.Data.ProductName,
			ProductImgs:          productImgs,
			ProductPrice:         product.Data.ProductPrice,
			ProductUrl:           product.Data.ProductUrl,
			IsTop:                product.Data.IsTop,
			MaterialOutUrls:      materialOutUrls,
			SampleThresholdType:  product.Data.SampleThresholdType,
			SampleThresholdValue: product.Data.SampleThresholdValue,
			Commissions:          commissions,
			InvestmentRatio:      product.Data.InvestmentRatio,
			ForbidReason:         product.Data.ForbidReason,
			IsTask:               product.Data.IsTask,
		},
	}, nil
}

func (is *InterfaceService) UpdateIsTopProducts(ctx context.Context, in *v1.UpdateIsTopProductsRequest) (*v1.UpdateIsTopProductsReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "product"); err != nil {
		return nil, err
	}

	product, err := is.prouc.UpdateIsTopProducts(ctx, in.ProductId, in.IsTop)

	if err != nil {
		return nil, err
	}

	productImgs := make([]*v1.UpdateIsTopProductsReply_ProductImg, 0)
	materialOutUrls := make([]*v1.UpdateIsTopProductsReply_MaterialOutUrl, 0)
	commissions := make([]*v1.UpdateIsTopProductsReply_Commission, 0)

	for _, productImg := range product.Data.ProductImgs {
		productImgs = append(productImgs, &v1.UpdateIsTopProductsReply_ProductImg{
			ProductImg: productImg.ProductImg,
		})
	}

	for _, materialOutUrl := range product.Data.MaterialOutUrls {
		materialOutUrls = append(materialOutUrls, &v1.UpdateIsTopProductsReply_MaterialOutUrl{
			MaterialOutUrl: materialOutUrl.MaterialOutUrl,
		})
	}

	for _, commission := range product.Data.Commissions {
		commissions = append(commissions, &v1.UpdateIsTopProductsReply_Commission{
			CommissionRatio:  float64(commission.CommissionRatio),
			ServiceRatio:     float64(commission.ServiceRatio),
			CommissionOutUrl: commission.CommissionOutUrl,
		})
	}

	return &v1.UpdateIsTopProductsReply{
		Code: 200,
		Data: &v1.UpdateIsTopProductsReply_Data{
			ProductId:            product.Data.ProductId,
			ProductOutId:         product.Data.ProductOutId,
			ProductType:          product.Data.ProductType,
			ProductStatus:        product.Data.ProductStatus,
			ProductName:          product.Data.ProductName,
			ProductImgs:          productImgs,
			ProductPrice:         product.Data.ProductPrice,
			ProductUrl:           product.Data.ProductUrl,
			IsTop:                product.Data.IsTop,
			MaterialOutUrls:      materialOutUrls,
			SampleThresholdType:  product.Data.SampleThresholdType,
			SampleThresholdValue: product.Data.SampleThresholdValue,
			Commissions:          commissions,
			InvestmentRatio:      product.Data.InvestmentRatio,
			ForbidReason:         product.Data.ForbidReason,
			IsTask:               product.Data.IsTask,
		},
	}, nil
}

func (is *InterfaceService) UpdateSampleThresholdProducts(ctx context.Context, in *v1.UpdateSampleThresholdProductsRequest) (*v1.UpdateSampleThresholdProductsReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "product"); err != nil {
		return nil, err
	}

	product, err := is.prouc.UpdateSampleThresholdProducts(ctx, in.ProductId, in.SampleThresholdValue, in.SampleThresholdType)

	if err != nil {
		return nil, err
	}

	productImgs := make([]*v1.UpdateSampleThresholdProductsReply_ProductImg, 0)
	materialOutUrls := make([]*v1.UpdateSampleThresholdProductsReply_MaterialOutUrl, 0)
	commissions := make([]*v1.UpdateSampleThresholdProductsReply_Commission, 0)

	for _, productImg := range product.Data.ProductImgs {
		productImgs = append(productImgs, &v1.UpdateSampleThresholdProductsReply_ProductImg{
			ProductImg: productImg.ProductImg,
		})
	}

	for _, materialOutUrl := range product.Data.MaterialOutUrls {
		materialOutUrls = append(materialOutUrls, &v1.UpdateSampleThresholdProductsReply_MaterialOutUrl{
			MaterialOutUrl: materialOutUrl.MaterialOutUrl,
		})
	}

	for _, commission := range product.Data.Commissions {
		commissions = append(commissions, &v1.UpdateSampleThresholdProductsReply_Commission{
			CommissionRatio:  float64(commission.CommissionRatio),
			ServiceRatio:     float64(commission.ServiceRatio),
			CommissionOutUrl: commission.CommissionOutUrl,
		})
	}

	return &v1.UpdateSampleThresholdProductsReply{
		Code: 200,
		Data: &v1.UpdateSampleThresholdProductsReply_Data{
			ProductId:            product.Data.ProductId,
			ProductOutId:         product.Data.ProductOutId,
			ProductType:          product.Data.ProductType,
			ProductStatus:        product.Data.ProductStatus,
			ProductName:          product.Data.ProductName,
			ProductImgs:          productImgs,
			ProductPrice:         product.Data.ProductPrice,
			ProductUrl:           product.Data.ProductUrl,
			IsTop:                product.Data.IsTop,
			MaterialOutUrls:      materialOutUrls,
			SampleThresholdType:  product.Data.SampleThresholdType,
			SampleThresholdValue: product.Data.SampleThresholdValue,
			Commissions:          commissions,
			InvestmentRatio:      product.Data.InvestmentRatio,
			ForbidReason:         product.Data.ForbidReason,
			IsTask:               product.Data.IsTask,
		},
	}, nil
}

func (is *InterfaceService) UpdateMaterialProducts(ctx context.Context, in *v1.UpdateMaterialProductsRequest) (*v1.UpdateMaterialProductsReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	if _, err := is.verifyLogin(ctx, true, false, "product"); err != nil {
		return nil, err
	}

	if _, err := is.prouc.UpdateMaterialProducts(ctx, in.ProductId, in.ProductMaterial); err != nil {
		return nil, err
	}

	return &v1.UpdateMaterialProductsReply{
		Code: 200,
		Data: &v1.UpdateMaterialProductsReply_Data{},
	}, nil
}

func (is *InterfaceService) UpdateInvestmentRatioProducts(ctx context.Context, in *v1.UpdateInvestmentRatioProductsRequest) (*v1.UpdateInvestmentRatioProductsReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "product"); err != nil {
		return nil, err
	}

	product, err := is.prouc.UpdateInvestmentRatioProducts(ctx, in.ProductId, in.InvestmentRatio)

	if err != nil {
		return nil, err
	}

	productImgs := make([]*v1.UpdateInvestmentRatioProductsReply_ProductImg, 0)
	materialOutUrls := make([]*v1.UpdateInvestmentRatioProductsReply_MaterialOutUrl, 0)
	commissions := make([]*v1.UpdateInvestmentRatioProductsReply_Commission, 0)

	for _, productImg := range product.Data.ProductImgs {
		productImgs = append(productImgs, &v1.UpdateInvestmentRatioProductsReply_ProductImg{
			ProductImg: productImg.ProductImg,
		})
	}

	for _, materialOutUrl := range product.Data.MaterialOutUrls {
		materialOutUrls = append(materialOutUrls, &v1.UpdateInvestmentRatioProductsReply_MaterialOutUrl{
			MaterialOutUrl: materialOutUrl.MaterialOutUrl,
		})
	}

	for _, commission := range product.Data.Commissions {
		commissions = append(commissions, &v1.UpdateInvestmentRatioProductsReply_Commission{
			CommissionRatio:  float64(commission.CommissionRatio),
			ServiceRatio:     float64(commission.ServiceRatio),
			CommissionOutUrl: commission.CommissionOutUrl,
		})
	}

	return &v1.UpdateInvestmentRatioProductsReply{
		Code: 200,
		Data: &v1.UpdateInvestmentRatioProductsReply_Data{
			ProductId:            product.Data.ProductId,
			ProductOutId:         product.Data.ProductOutId,
			ProductType:          product.Data.ProductType,
			ProductStatus:        product.Data.ProductStatus,
			ProductName:          product.Data.ProductName,
			ProductImgs:          productImgs,
			ProductPrice:         product.Data.ProductPrice,
			ProductUrl:           product.Data.ProductUrl,
			IsTop:                product.Data.IsTop,
			MaterialOutUrls:      materialOutUrls,
			SampleThresholdType:  product.Data.SampleThresholdType,
			SampleThresholdValue: product.Data.SampleThresholdValue,
			Commissions:          commissions,
			InvestmentRatio:      product.Data.InvestmentRatio,
			ForbidReason:         product.Data.ForbidReason,
			IsTask:               product.Data.IsTask,
		},
	}, nil
}

func (is *InterfaceService) ParseMiniProductProducts(ctx context.Context, in *v1.ParseMiniProductProductsRequest) (*v1.ParseMiniProductProductsReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	if _, err := is.verifyMiniUserLogin(ctx); err != nil {
		return nil, err
	}

	product, err := is.prouc.ParseMiniProductProducts(ctx, in.Content)

	if err != nil {
		return nil, err
	}

	return &v1.ParseMiniProductProductsReply{
		Code: 200,
		Data: &v1.ParseMiniProductProductsReply_Data{
			ProductId:             product.Data.ProductId,
			ProductName:           product.Data.ProductName,
			ProductImg:            product.Data.ProductImg,
			ProductPrice:          product.Data.ProductPrice,
			IsTop:                 product.Data.IsTop,
			PureCommission:        product.Data.PureCommission,
			PureServiceCommission: product.Data.PureServiceCommission,
			CommonCommission:      product.Data.CommonCommission,
			InvestmentRatio:       product.Data.InvestmentRatio,
			IsHot:                 product.Data.IsHot,
			TotalSale:             product.Data.TotalSale,
			IsTask:                product.Data.IsTask,
		},
	}, nil
}

func (is *InterfaceService) UploadPartProducts(ctx context.Context, in *v1.UploadPartProductsRequest) (*v1.UploadPartProductsReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "product"); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := is.prouc.UploadPartProducts(ctx, in.PartNumber, in.TotalPart, in.ContentLength, in.UploadId, in.Content); err != nil {
		return nil, err
	}

	return &v1.UploadPartProductsReply{
		Code: 200,
		Data: &v1.UploadPartProductsReply_Data{},
	}, nil
}

func (is *InterfaceService) CompleteUploadProducts(ctx context.Context, in *v1.CompleteUploadRequest) (*v1.CompleteUploadReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	if _, err := is.verifyLogin(ctx, true, false, "product"); err != nil {
		return nil, err
	}

	staticUrl, err := is.prouc.CompleteUploadProducts(ctx, in.UploadId)

	if err != nil {
		return nil, err
	}

	return &v1.CompleteUploadReply{
		Code: 200,
		Data: &v1.CompleteUploadReply_Data{
			StaticUrl: staticUrl.Data.StaticUrl,
		},
	}, nil
}

func (is *InterfaceService) AbortUploadProducts(ctx context.Context, in *v1.AbortUploadProductsRequest) (*v1.AbortUploadProductsReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	if _, err := is.verifyLogin(ctx, true, false, "product"); err != nil {
		return nil, err
	}

	if err := is.prouc.AbortUploadProducts(ctx, in.UploadId); err != nil {
		return nil, err
	}

	return &v1.AbortUploadProductsReply{
		Code: 200,
		Data: &v1.AbortUploadProductsReply_Data{},
	}, nil
}

func (is *InterfaceService) DeleteProducts(ctx context.Context, in *v1.DeleteProductsRequest) (*v1.DeleteProductsReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "product"); err != nil {
		return nil, err
	}

	if err := is.prouc.DeleteProducts(ctx, in.ProductId); err != nil {
		return nil, err
	}

	return &v1.DeleteProductsReply{
		Code: 200,
		Data: &v1.DeleteProductsReply_Data{},
	}, nil
}

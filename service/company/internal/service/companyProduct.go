package service

import (
	v1 "company/api/company/v1"
	"company/internal/pkg/tool"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"math"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (cs *CompanyService) GetCompanyProducts(ctx context.Context, in *v1.GetCompanyProductsRequest) (*v1.GetCompanyProductsReply, error) {
	companyProduct, err := cs.cprouc.GetCompanyProducts(ctx, in.ProductId, in.ProductStatus)

	if err != nil {
		return nil, err
	}

	productImgs := make([]*v1.GetCompanyProductsReply_ProductImg, 0)
	materialOutUrls := make([]*v1.GetCompanyProductsReply_MaterialOutUrl, 0)
	commissions := make([]*v1.GetCompanyProductsReply_Commission, 0)

	for _, productImg := range companyProduct.ProductImgs {
		productImgs = append(productImgs, &v1.GetCompanyProductsReply_ProductImg{
			ProductImg: productImg,
		})
	}

	for _, materialOutUrl := range companyProduct.MaterialOutUrls {
		materialOutUrls = append(materialOutUrls, &v1.GetCompanyProductsReply_MaterialOutUrl{
			MaterialOutUrl: materialOutUrl,
		})
	}

	for _, commission := range companyProduct.Commissions {
		commissions = append(commissions, &v1.GetCompanyProductsReply_Commission{
			CommissionRatio:  float64(commission.CommissionRatio),
			ServiceRatio:     float64(commission.ServiceRatio),
			CommissionOutUrl: commission.CommissionOutUrl,
		})
	}

	return &v1.GetCompanyProductsReply{
		Code: 200,
		Data: &v1.GetCompanyProductsReply_Data{
			ProductId:       companyProduct.Id,
			ProductOutId:    companyProduct.ProductOutId,
			ProductType:     uint32(companyProduct.ProductType),
			ProductStatus:   uint32(companyProduct.ProductStatus),
			ProductName:     companyProduct.ProductName,
			ProductImgs:     productImgs,
			ProductPrice:    companyProduct.ProductPrice,
			ProductUrl:      companyProduct.ProductUrl,
			IsTop:           uint32(companyProduct.IsTop),
			MaterialOutUrls: materialOutUrls,
			Commissions:     commissions,
			InvestmentRatio: float64(companyProduct.InvestmentRatio),
			ForbidReason:    companyProduct.ForbidReason,
			IsTask:          uint32(companyProduct.IsTask),
		},
	}, nil
}

func (cs *CompanyService) GetCompanyProductByProductOutIds(ctx context.Context, in *v1.GetCompanyProductByProductOutIdsRequest) (*v1.GetCompanyProductByProductOutIdsReply, error) {
	companyProduct, err := cs.cprouc.GetCompanyProductByProductOutIds(ctx, in.ProductOutId, in.ProductStatus)

	if err != nil {
		return nil, err
	}

	productImgs := make([]*v1.GetCompanyProductByProductOutIdsReply_ProductImg, 0)
	materialOutUrls := make([]*v1.GetCompanyProductByProductOutIdsReply_MaterialOutUrl, 0)
	commissions := make([]*v1.GetCompanyProductByProductOutIdsReply_Commission, 0)

	for _, productImg := range companyProduct.ProductImgs {
		productImgs = append(productImgs, &v1.GetCompanyProductByProductOutIdsReply_ProductImg{
			ProductImg: productImg,
		})
	}

	for _, materialOutUrl := range companyProduct.MaterialOutUrls {
		materialOutUrls = append(materialOutUrls, &v1.GetCompanyProductByProductOutIdsReply_MaterialOutUrl{
			MaterialOutUrl: materialOutUrl,
		})
	}

	for _, commission := range companyProduct.Commissions {
		commissions = append(commissions, &v1.GetCompanyProductByProductOutIdsReply_Commission{
			CommissionRatio:  float64(commission.CommissionRatio),
			ServiceRatio:     float64(commission.ServiceRatio),
			CommissionOutUrl: commission.CommissionOutUrl,
		})
	}

	return &v1.GetCompanyProductByProductOutIdsReply{
		Code: 200,
		Data: &v1.GetCompanyProductByProductOutIdsReply_Data{
			ProductId:       companyProduct.Id,
			ProductOutId:    companyProduct.ProductOutId,
			ProductType:     uint32(companyProduct.ProductType),
			ProductStatus:   uint32(companyProduct.ProductStatus),
			ProductName:     companyProduct.ProductName,
			ProductImgs:     productImgs,
			ProductPrice:    companyProduct.ProductPrice,
			ProductUrl:      companyProduct.ProductUrl,
			IsTop:           uint32(companyProduct.IsTop),
			MaterialOutUrls: materialOutUrls,
			Commissions:     commissions,
			InvestmentRatio: float64(companyProduct.InvestmentRatio),
			ForbidReason:    companyProduct.ForbidReason,
			IsTask:          uint32(companyProduct.IsTask),
		},
	}, nil
}

func (cs *CompanyService) GetExternalCompanyProducts(ctx context.Context, in *v1.GetExternalCompanyProductsRequest) (*v1.GetExternalCompanyProductsReply, error) {
	companyProduct, err := cs.cprouc.GetExternalCompanyProducts(ctx, in.ProductId)

	if err != nil {
		return nil, err
	}

	productImgs := make([]*v1.GetExternalCompanyProductsReply_ProductImg, 0)
	awemes := make([]*v1.GetExternalCompanyProductsReply_Aweme, 0)

	for _, productImg := range companyProduct.ProductImgs {
		productImgs = append(productImgs, &v1.GetExternalCompanyProductsReply_ProductImg{
			ProductImg: productImg,
		})
	}

	for _, aweme := range companyProduct.Awemes {
		awemes = append(awemes, &v1.GetExternalCompanyProductsReply_Aweme{
			AccountId:    aweme.AccountId,
			Nickname:     aweme.Nickname,
			Avatar:       aweme.Avatar,
			AvatarLarger: aweme.AvatarLarger,
		})
	}

	return &v1.GetExternalCompanyProductsReply{
		Code: 200,
		Data: &v1.GetExternalCompanyProductsReply_Data{
			ProductId:             companyProduct.ProductOutId,
			ProductName:           companyProduct.ProductName,
			ProductImgs:           productImgs,
			ProductPrice:          companyProduct.ProductPrice,
			ShopName:              companyProduct.ShopName,
			ShopScore:             fmt.Sprintf("%.2f", tool.Decimal(float64(companyProduct.ShopScore), 2)),
			IsTop:                 uint32(companyProduct.IsTop),
			PureCommission:        companyProduct.PureCommission,
			PureServiceCommission: companyProduct.PureServiceCommission,
			CommonCommission:      fmt.Sprintf("%.f", tool.Decimal(float64(companyProduct.CommissionRatio), 2)),
			IsHot:                 uint32(companyProduct.IsHot),
			ProductUrl:            companyProduct.ProductUrl,
			TotalSale:             companyProduct.TotalSale,
			InvestmentRatio:       fmt.Sprintf("%.f", tool.Decimal(float64(companyProduct.InvestmentRatio), 2)),
			Awemes:                awemes,
		},
	}, nil
}

func (cs *CompanyService) GetExternalProductShareCompanyProducts(ctx context.Context, in *v1.GetExternalProductShareCompanyProductsRequest) (*v1.GetExternalProductShareCompanyProductsReply, error) {
	productShare, err := cs.cprouc.GetExternalProductShareCompanyProducts(ctx, in.ProductId, in.UserId)

	if err != nil {
		return nil, err
	}

	return &v1.GetExternalProductShareCompanyProductsReply{
		Code: 200,
		Data: &v1.GetExternalProductShareCompanyProductsReply_Data{
			DyPassword: productShare.Data.DyPassword,
		},
	}, nil
}

func (cs *CompanyService) GetUploadIdCompanyProducts(ctx context.Context, in *v1.GetUploadIdCompanyProductsRequest) (*v1.GetUploadIdCompanyProductsReply, error) {
	uploadId, err := cs.cprouc.GetUploadIdCompanyProducts(ctx, in.Suffix)

	if err != nil {
		return nil, err
	}

	return &v1.GetUploadIdCompanyProductsReply{
		Code: 200,
		Data: &v1.GetUploadIdCompanyProductsReply_Data{
			UploadId: uploadId,
		},
	}, nil
}

func (cs *CompanyService) ListCompanyProducts(ctx context.Context, in *v1.ListCompanyProductsRequest) (*v1.ListCompanyProductsReply, error) {
	if in.IndustryId == 0 {
		in.CategoryId = 0
		in.SubCategoryId = 0
	} else if in.CategoryId == 0 {
		in.SubCategoryId = 0
	}

	companyProducts, err := cs.cprouc.ListCompanyProducts(ctx, in.PageNum, in.PageSize, in.IndustryId, in.CategoryId, in.SubCategoryId, in.ProductStatus, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyProductsReply_CompanyProduct, 0)

	for _, companyProduct := range companyProducts.List {
		productImgs := make([]*v1.ListCompanyProductsReply_ProductImg, 0)
		materialOutUrls := make([]*v1.ListCompanyProductsReply_MaterialOutUrl, 0)
		commissions := make([]*v1.ListCompanyProductsReply_Commission, 0)

		for _, productImg := range companyProduct.ProductImgs {
			productImgs = append(productImgs, &v1.ListCompanyProductsReply_ProductImg{
				ProductImg: productImg,
			})
		}

		for _, materialOutUrl := range companyProduct.MaterialOutUrls {
			materialOutUrls = append(materialOutUrls, &v1.ListCompanyProductsReply_MaterialOutUrl{
				MaterialOutUrl: materialOutUrl,
			})
		}

		for _, commission := range companyProduct.Commissions {
			commissions = append(commissions, &v1.ListCompanyProductsReply_Commission{
				CommissionRatio:  float64(commission.CommissionRatio),
				ServiceRatio:     float64(commission.ServiceRatio),
				CommissionOutUrl: commission.CommissionOutUrl,
			})
		}

		list = append(list, &v1.ListCompanyProductsReply_CompanyProduct{
			ProductId:       companyProduct.Id,
			ProductOutId:    companyProduct.ProductOutId,
			ProductType:     uint32(companyProduct.ProductType),
			ProductStatus:   uint32(companyProduct.ProductStatus),
			ProductName:     companyProduct.ProductName,
			ProductImgs:     productImgs,
			ProductPrice:    companyProduct.ProductPrice,
			ProductUrl:      companyProduct.ProductUrl,
			IsTop:           uint32(companyProduct.IsTop),
			MaterialOutUrls: materialOutUrls,
			Commissions:     commissions,
			InvestmentRatio: float64(companyProduct.InvestmentRatio),
			ForbidReason:    companyProduct.ForbidReason,
			IsTask:          uint32(companyProduct.IsTask),
		})
	}

	totalPage := uint64(math.Ceil(float64(companyProducts.Total) / float64(companyProducts.PageSize)))

	return &v1.ListCompanyProductsReply{
		Code: 200,
		Data: &v1.ListCompanyProductsReply_Data{
			PageNum:   companyProducts.PageNum,
			PageSize:  companyProducts.PageSize,
			Total:     companyProducts.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (cs *CompanyService) ListSimpleCompanyProducts(ctx context.Context, in *v1.ListSimpleCompanyProductsRequest) (*v1.ListSimpleCompanyProductsReply, error) {
	if in.IndustryId == 0 {
		in.CategoryId = 0
		in.SubCategoryId = 0
	} else if in.CategoryId == 0 {
		in.SubCategoryId = 0
	}

	companyProducts, err := cs.cprouc.ListCompanyProducts(ctx, in.PageNum, in.PageSize, in.IndustryId, in.CategoryId, in.SubCategoryId, in.ProductStatus, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListSimpleCompanyProductsReply_CompanyProduct, 0)

	for _, companyProduct := range companyProducts.List {
		productImg := ""

		if len(companyProduct.ProductImgs) > 0 {
			productImg = companyProduct.ProductImgs[0]
		}

		list = append(list, &v1.ListSimpleCompanyProductsReply_CompanyProduct{
			ProductId:    companyProduct.Id,
			ProductName:  companyProduct.ProductName,
			ProductImg:   productImg,
			ProductPrice: companyProduct.ProductPrice + "¥",
		})
	}

	totalPage := uint64(math.Ceil(float64(companyProducts.Total) / float64(companyProducts.PageSize)))

	return &v1.ListSimpleCompanyProductsReply{
		Code: 200,
		Data: &v1.ListSimpleCompanyProductsReply_Data{
			PageNum:   companyProducts.PageNum,
			PageSize:  companyProducts.PageSize,
			Total:     companyProducts.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (cs *CompanyService) ListExternalCompanyProducts(ctx context.Context, in *v1.ListExternalCompanyProductsRequest) (*v1.ListExternalCompanyProductsReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if in.IndustryId == 0 {
		in.CategoryId = 0
		in.SubCategoryId = 0
	} else if in.CategoryId == 0 {
		in.SubCategoryId = 0
	}

	companyProducts, err := cs.cprouc.ListExternalCompanyProducts(ctx, in.PageNum, in.PageSize, in.IndustryId, in.CategoryId, in.SubCategoryId, uint8(in.IsInvestment), in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListExternalCompanyProductsReply_CompanyProduct, 0)

	for _, companyProduct := range companyProducts.List {
		productImg := ""

		if len(companyProduct.ProductImgs) > 0 {
			productImg = companyProduct.ProductImgs[0]
		}

		list = append(list, &v1.ListExternalCompanyProductsReply_CompanyProduct{
			ProductId:             companyProduct.ProductOutId,
			ProductName:           companyProduct.ProductName,
			ProductImg:            productImg,
			ProductPrice:          companyProduct.ProductPrice,
			IsTop:                 uint32(companyProduct.IsTop),
			PureCommission:        companyProduct.PureCommission,
			PureServiceCommission: companyProduct.PureServiceCommission,
			CommonCommission:      fmt.Sprintf("%.f", tool.Decimal(float64(companyProduct.CommissionRatio), 2)),
			InvestmentRatio:       fmt.Sprintf("%.f", tool.Decimal(float64(companyProduct.InvestmentRatio), 2)),
			IsHot:                 uint32(companyProduct.IsHot),
			TotalSale:             companyProduct.TotalSale,
			IsTask:                uint32(companyProduct.IsTask),
		})
	}

	totalPage := uint64(math.Ceil(float64(companyProducts.Total) / float64(companyProducts.PageSize)))

	return &v1.ListExternalCompanyProductsReply{
		Code: 200,
		Data: &v1.ListExternalCompanyProductsReply_Data{
			PageNum:   companyProducts.PageNum,
			PageSize:  companyProducts.PageSize,
			Total:     companyProducts.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (cs *CompanyService) ListCompanyProductCategorys(ctx context.Context, in *emptypb.Empty) (*v1.ListCompanyProductCategorysReply, error) {
	companyProductCategories, err := cs.cprouc.ListCompanyProductCategorys(ctx)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyProductCategorysReply_Category, 0)

	for _, companyProductCategory := range companyProductCategories {
		childList := make([]*v1.ListCompanyProductCategorysReply_ChildCategory, 0)

		for _, lchildList := range companyProductCategory.ChildList {
			subChildList := make([]*v1.ListCompanyProductCategorysReply_SubChildCategory, 0)

			for _, llchildList := range lchildList.ChildList {
				subChildList = append(subChildList, &v1.ListCompanyProductCategorysReply_SubChildCategory{
					Key:   strconv.FormatUint(llchildList.CategoryId, 10),
					Value: llchildList.CategoryName,
				})
			}

			childList = append(childList, &v1.ListCompanyProductCategorysReply_ChildCategory{
				Key:       strconv.FormatUint(lchildList.CategoryId, 10),
				Value:     lchildList.CategoryName,
				ChildList: subChildList,
			})
		}

		list = append(list, &v1.ListCompanyProductCategorysReply_Category{
			Key:       strconv.FormatUint(companyProductCategory.CategoryId, 10),
			Value:     companyProductCategory.CategoryName,
			ChildList: childList,
		})
	}

	return &v1.ListCompanyProductCategorysReply{
		Code: 200,
		Data: &v1.ListCompanyProductCategorysReply_Data{
			Category: list,
		},
	}, nil
}

func (cs *CompanyService) ListCompanyTaskProducts(ctx context.Context, in *v1.ListCompanyTaskProductsRequest) (*v1.ListCompanyTaskProductsReply, error) {
	companyProducts, err := cs.cprouc.ListCompanyTaskProducts(ctx, in.PageNum, in.PageSize, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCompanyTaskProductsReply_CompanyProduct, 0)

	for _, companyProduct := range companyProducts.List {
		productImgs := make([]*v1.ListCompanyTaskProductsReply_ProductImg, 0)
		materialOutUrls := make([]*v1.ListCompanyTaskProductsReply_MaterialOutUrl, 0)
		commissions := make([]*v1.ListCompanyTaskProductsReply_Commission, 0)

		for _, productImg := range companyProduct.ProductImgs {
			productImgs = append(productImgs, &v1.ListCompanyTaskProductsReply_ProductImg{
				ProductImg: productImg,
			})
		}

		for _, materialOutUrl := range companyProduct.MaterialOutUrls {
			materialOutUrls = append(materialOutUrls, &v1.ListCompanyTaskProductsReply_MaterialOutUrl{
				MaterialOutUrl: materialOutUrl,
			})
		}

		for _, commission := range companyProduct.Commissions {
			commissions = append(commissions, &v1.ListCompanyTaskProductsReply_Commission{
				CommissionRatio:  float64(commission.CommissionRatio),
				ServiceRatio:     float64(commission.ServiceRatio),
				CommissionOutUrl: commission.CommissionOutUrl,
			})
		}

		list = append(list, &v1.ListCompanyTaskProductsReply_CompanyProduct{
			ProductId:       companyProduct.Id,
			ProductOutId:    companyProduct.ProductOutId,
			ProductType:     uint32(companyProduct.ProductType),
			ProductStatus:   uint32(companyProduct.ProductStatus),
			ProductName:     companyProduct.ProductName,
			ProductImgs:     productImgs,
			ProductPrice:    companyProduct.ProductPrice,
			ProductUrl:      companyProduct.ProductUrl,
			IsTop:           uint32(companyProduct.IsTop),
			MaterialOutUrls: materialOutUrls,
			Commissions:     commissions,
			InvestmentRatio: float64(companyProduct.InvestmentRatio),
			ForbidReason:    companyProduct.ForbidReason,
			IsTask:          uint32(companyProduct.IsTask),
		})
	}

	totalPage := uint64(math.Ceil(float64(companyProducts.Total) / float64(companyProducts.PageSize)))

	return &v1.ListCompanyTaskProductsReply{
		Code: 200,
		Data: &v1.ListCompanyTaskProductsReply_Data{
			PageNum:   companyProducts.PageNum,
			PageSize:  companyProducts.PageSize,
			Total:     companyProducts.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (cs *CompanyService) StatisticsCompanyProducts(ctx context.Context, in *v1.StatisticsCompanyProductsRequest) (*v1.StatisticsCompanyProductsReply, error) {
	if in.IndustryId == 0 {
		in.CategoryId = 0
		in.SubCategoryId = 0
	} else if in.CategoryId == 0 {
		in.SubCategoryId = 0
	}

	statistics, err := cs.cprouc.StatisticsCompanyProducts(ctx, in.IndustryId, in.CategoryId, in.SubCategoryId, in.ProductStatus, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsCompanyProductsReply_Statistic, 0)

	for _, statistic := range statistics {
		list = append(list, &v1.StatisticsCompanyProductsReply_Statistic{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsCompanyProductsReply{
		Code: 200,
		Data: &v1.StatisticsCompanyProductsReply_Data{
			Statistics: list,
		},
	}, nil
}

func (cs *CompanyService) CreateCompanyProducts(ctx context.Context, in *v1.CreateCompanyProductsRequest) (*v1.CreateCompanyProductsReply, error) {
	companyProduct, err := cs.cprouc.CreateCompanyProducts(ctx, in.Commission)

	if err != nil {
		return nil, err
	}

	productImgs := make([]*v1.CreateCompanyProductsReply_ProductImg, 0)
	materialOutUrls := make([]*v1.CreateCompanyProductsReply_MaterialOutUrl, 0)
	commissions := make([]*v1.CreateCompanyProductsReply_Commission, 0)

	for _, productImg := range companyProduct.ProductImgs {
		productImgs = append(productImgs, &v1.CreateCompanyProductsReply_ProductImg{
			ProductImg: productImg,
		})
	}

	for _, materialOutUrl := range companyProduct.MaterialOutUrls {
		materialOutUrls = append(materialOutUrls, &v1.CreateCompanyProductsReply_MaterialOutUrl{
			MaterialOutUrl: materialOutUrl,
		})
	}

	for _, commission := range companyProduct.Commissions {
		commissions = append(commissions, &v1.CreateCompanyProductsReply_Commission{
			CommissionRatio:  float64(commission.CommissionRatio),
			ServiceRatio:     float64(commission.ServiceRatio),
			CommissionOutUrl: commission.CommissionOutUrl,
		})
	}

	return &v1.CreateCompanyProductsReply{
		Code: 200,
		Data: &v1.CreateCompanyProductsReply_Data{
			ProductId:       companyProduct.Id,
			ProductOutId:    companyProduct.ProductOutId,
			ProductType:     uint32(companyProduct.ProductType),
			ProductStatus:   uint32(companyProduct.ProductStatus),
			ProductName:     companyProduct.ProductName,
			ProductImgs:     productImgs,
			ProductPrice:    companyProduct.ProductPrice,
			ProductUrl:      companyProduct.ProductUrl,
			IsTop:           uint32(companyProduct.IsTop),
			MaterialOutUrls: materialOutUrls,
			Commissions:     commissions,
			InvestmentRatio: float64(companyProduct.InvestmentRatio),
			ForbidReason:    companyProduct.ForbidReason,
			IsTask:          uint32(companyProduct.IsTask),
		},
	}, nil
}

func (cs *CompanyService) CreateJinritemaiStoreCompanyProducts(ctx context.Context, in *v1.CreateJinritemaiStoreCompanyProductsRequest) (*v1.CreateJinritemaiStoreCompanyProductsReply, error) {
	companyProductJinritemaiStore, err := cs.cprouc.CreateJinritemaiStoreCompanyProducts(ctx, in.UserId, in.ProductId, in.OpenDouyinUserIds)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.CreateJinritemaiStoreCompanyProductsReply_Message, 0)

	for _, message := range companyProductJinritemaiStore.Messages {
		list = append(list, &v1.CreateJinritemaiStoreCompanyProductsReply_Message{
			ProductName: message.ProductName,
			AwemeName:   message.AwemeName,
			Content:     message.Content,
		})
	}

	return &v1.CreateJinritemaiStoreCompanyProductsReply{
		Code: 200,
		Data: &v1.CreateJinritemaiStoreCompanyProductsReply_Data{
			Content: companyProductJinritemaiStore.Content,
			List:    list,
		},
	}, nil
}

func (cs *CompanyService) UpdateStatusCompanyProducts(ctx context.Context, in *v1.UpdateStatusCompanyProductsRequest) (*v1.UpdateStatusCompanyProductsReply, error) {
	companyProduct, err := cs.cprouc.UpdateStatusCompanyProducts(ctx, in.ProductId, uint8(in.Status))

	if err != nil {
		return nil, err
	}

	productImgs := make([]*v1.UpdateStatusCompanyProductsReply_ProductImg, 0)
	materialOutUrls := make([]*v1.UpdateStatusCompanyProductsReply_MaterialOutUrl, 0)
	commissions := make([]*v1.UpdateStatusCompanyProductsReply_Commission, 0)

	for _, productImg := range companyProduct.ProductImgs {
		productImgs = append(productImgs, &v1.UpdateStatusCompanyProductsReply_ProductImg{
			ProductImg: productImg,
		})
	}

	for _, materialOutUrl := range companyProduct.MaterialOutUrls {
		materialOutUrls = append(materialOutUrls, &v1.UpdateStatusCompanyProductsReply_MaterialOutUrl{
			MaterialOutUrl: materialOutUrl,
		})
	}

	for _, commission := range companyProduct.Commissions {
		commissions = append(commissions, &v1.UpdateStatusCompanyProductsReply_Commission{
			CommissionRatio:  float64(commission.CommissionRatio),
			ServiceRatio:     float64(commission.ServiceRatio),
			CommissionOutUrl: commission.CommissionOutUrl,
		})
	}

	return &v1.UpdateStatusCompanyProductsReply{
		Code: 200,
		Data: &v1.UpdateStatusCompanyProductsReply_Data{
			ProductId:       companyProduct.Id,
			ProductOutId:    companyProduct.ProductOutId,
			ProductType:     uint32(companyProduct.ProductType),
			ProductStatus:   uint32(companyProduct.ProductStatus),
			ProductName:     companyProduct.ProductName,
			ProductImgs:     productImgs,
			ProductPrice:    companyProduct.ProductPrice,
			ProductUrl:      companyProduct.ProductUrl,
			IsTop:           uint32(companyProduct.IsTop),
			MaterialOutUrls: materialOutUrls,
			Commissions:     commissions,
			InvestmentRatio: float64(companyProduct.InvestmentRatio),
			ForbidReason:    companyProduct.ForbidReason,
			IsTask:          uint32(companyProduct.IsTask),
		},
	}, nil
}

func (cs *CompanyService) UpdateIsTopCompanyProducts(ctx context.Context, in *v1.UpdateIsTopCompanyProductsRequest) (*v1.UpdateIsTopCompanyProductsReply, error) {
	companyProduct, err := cs.cprouc.UpdateIsTopCompanyProducts(ctx, in.ProductId, uint8(in.IsTop))

	if err != nil {
		return nil, err
	}

	productImgs := make([]*v1.UpdateIsTopCompanyProductsReply_ProductImg, 0)
	materialOutUrls := make([]*v1.UpdateIsTopCompanyProductsReply_MaterialOutUrl, 0)
	commissions := make([]*v1.UpdateIsTopCompanyProductsReply_Commission, 0)

	for _, productImg := range companyProduct.ProductImgs {
		productImgs = append(productImgs, &v1.UpdateIsTopCompanyProductsReply_ProductImg{
			ProductImg: productImg,
		})
	}

	for _, materialOutUrl := range companyProduct.MaterialOutUrls {
		materialOutUrls = append(materialOutUrls, &v1.UpdateIsTopCompanyProductsReply_MaterialOutUrl{
			MaterialOutUrl: materialOutUrl,
		})
	}

	for _, commission := range companyProduct.Commissions {
		commissions = append(commissions, &v1.UpdateIsTopCompanyProductsReply_Commission{
			CommissionRatio:  float64(commission.CommissionRatio),
			ServiceRatio:     float64(commission.ServiceRatio),
			CommissionOutUrl: commission.CommissionOutUrl,
		})
	}

	return &v1.UpdateIsTopCompanyProductsReply{
		Code: 200,
		Data: &v1.UpdateIsTopCompanyProductsReply_Data{
			ProductId:       companyProduct.Id,
			ProductOutId:    companyProduct.ProductOutId,
			ProductType:     uint32(companyProduct.ProductType),
			ProductStatus:   uint32(companyProduct.ProductStatus),
			ProductName:     companyProduct.ProductName,
			ProductImgs:     productImgs,
			ProductPrice:    companyProduct.ProductPrice,
			ProductUrl:      companyProduct.ProductUrl,
			IsTop:           uint32(companyProduct.IsTop),
			MaterialOutUrls: materialOutUrls,
			Commissions:     commissions,
			InvestmentRatio: float64(companyProduct.InvestmentRatio),
			ForbidReason:    companyProduct.ForbidReason,
			IsTask:          uint32(companyProduct.IsTask),
		},
	}, nil
}

func (cs *CompanyService) UpdateCommissionCompanyProducts(ctx context.Context, in *v1.UpdateCommissionCompanyProductsRequest) (*v1.UpdateCommissionCompanyProductsReply, error) {
	companyProduct, err := cs.cprouc.UpdateCommissionCompanyProducts(ctx, in.ProductId, in.Commission)

	if err != nil {
		return nil, err
	}

	if companyProduct == nil {
		return &v1.UpdateCommissionCompanyProductsReply{
			Code: 200,
			Data: &v1.UpdateCommissionCompanyProductsReply_Data{},
		}, nil
	} else {
		productImgs := make([]*v1.UpdateCommissionCompanyProductsReply_ProductImg, 0)
		materialOutUrls := make([]*v1.UpdateCommissionCompanyProductsReply_MaterialOutUrl, 0)
		commissions := make([]*v1.UpdateCommissionCompanyProductsReply_Commission, 0)

		for _, productImg := range companyProduct.ProductImgs {
			productImgs = append(productImgs, &v1.UpdateCommissionCompanyProductsReply_ProductImg{
				ProductImg: productImg,
			})
		}

		for _, materialOutUrl := range companyProduct.MaterialOutUrls {
			materialOutUrls = append(materialOutUrls, &v1.UpdateCommissionCompanyProductsReply_MaterialOutUrl{
				MaterialOutUrl: materialOutUrl,
			})
		}

		for _, commission := range companyProduct.Commissions {
			commissions = append(commissions, &v1.UpdateCommissionCompanyProductsReply_Commission{
				CommissionRatio:  float64(commission.CommissionRatio),
				ServiceRatio:     float64(commission.ServiceRatio),
				CommissionOutUrl: commission.CommissionOutUrl,
			})
		}

		return &v1.UpdateCommissionCompanyProductsReply{
			Code: 200,
			Data: &v1.UpdateCommissionCompanyProductsReply_Data{
				ProductId:       companyProduct.Id,
				ProductOutId:    companyProduct.ProductOutId,
				ProductType:     uint32(companyProduct.ProductType),
				ProductStatus:   uint32(companyProduct.ProductStatus),
				ProductName:     companyProduct.ProductName,
				ProductImgs:     productImgs,
				ProductPrice:    companyProduct.ProductPrice,
				ProductUrl:      companyProduct.ProductUrl,
				IsTop:           uint32(companyProduct.IsTop),
				MaterialOutUrls: materialOutUrls,
				Commissions:     commissions,
				InvestmentRatio: float64(companyProduct.InvestmentRatio),
				ForbidReason:    companyProduct.ForbidReason,
				IsTask:          uint32(companyProduct.IsTask),
			},
		}, nil
	}
}

func (cs *CompanyService) UpdateMaterialCompanyProducts(ctx context.Context, in *v1.UpdateMaterialCompanyProductsRequest) (*v1.UpdateMaterialCompanyProductsReply, error) {
	if err := cs.cprouc.UpdateMaterialCompanyProducts(ctx, in.ProductId, in.ProductMaterial); err != nil {
		return nil, err
	}

	return &v1.UpdateMaterialCompanyProductsReply{
		Code: 200,
		Data: &v1.UpdateMaterialCompanyProductsReply_Data{},
	}, nil
}

func (cs *CompanyService) UpdateInvestmentRatioCompanyProducts(ctx context.Context, in *v1.UpdateInvestmentRatioCompanyProductsRequest) (*v1.UpdateInvestmentRatioCompanyProductsReply, error) {
	companyProduct, err := cs.cprouc.UpdateInvestmentRatioCompanyProducts(ctx, in.ProductId, in.InvestmentRatio)

	if err != nil {
		return nil, err
	}

	productImgs := make([]*v1.UpdateInvestmentRatioCompanyProductsReply_ProductImg, 0)
	materialOutUrls := make([]*v1.UpdateInvestmentRatioCompanyProductsReply_MaterialOutUrl, 0)
	commissions := make([]*v1.UpdateInvestmentRatioCompanyProductsReply_Commission, 0)

	for _, productImg := range companyProduct.ProductImgs {
		productImgs = append(productImgs, &v1.UpdateInvestmentRatioCompanyProductsReply_ProductImg{
			ProductImg: productImg,
		})
	}

	for _, materialOutUrl := range companyProduct.MaterialOutUrls {
		materialOutUrls = append(materialOutUrls, &v1.UpdateInvestmentRatioCompanyProductsReply_MaterialOutUrl{
			MaterialOutUrl: materialOutUrl,
		})
	}

	for _, commission := range companyProduct.Commissions {
		commissions = append(commissions, &v1.UpdateInvestmentRatioCompanyProductsReply_Commission{
			CommissionRatio:  float64(commission.CommissionRatio),
			ServiceRatio:     float64(commission.ServiceRatio),
			CommissionOutUrl: commission.CommissionOutUrl,
		})
	}

	return &v1.UpdateInvestmentRatioCompanyProductsReply{
		Code: 200,
		Data: &v1.UpdateInvestmentRatioCompanyProductsReply_Data{
			ProductId:       companyProduct.Id,
			ProductOutId:    companyProduct.ProductOutId,
			ProductType:     uint32(companyProduct.ProductType),
			ProductStatus:   uint32(companyProduct.ProductStatus),
			ProductName:     companyProduct.ProductName,
			ProductImgs:     productImgs,
			ProductPrice:    companyProduct.ProductPrice,
			ProductUrl:      companyProduct.ProductUrl,
			IsTop:           uint32(companyProduct.IsTop),
			MaterialOutUrls: materialOutUrls,
			Commissions:     commissions,
			InvestmentRatio: float64(companyProduct.InvestmentRatio),
			ForbidReason:    companyProduct.ForbidReason,
			IsTask:          uint32(companyProduct.IsTask),
		},
	}, nil
}

func (cs *CompanyService) UpdateOutProductCompanyProducts(ctx context.Context, in *v1.UpdateOutProductCompanyProductsRequest) (*v1.UpdateOutProductCompanyProductsReply, error) {
	err := cs.cprouc.UpdateOutProductCompanyProducts(ctx, in.ProductOutId, in.IndustryId, in.CategoryId, in.SubCategoryId, in.TotalSale, uint8(in.ProductStatus), in.ShopScore, in.CommissionRatio, in.ProductName, in.ProductImg, in.ProductDetailImg, in.ProductPrice, in.IndustryName, in.CategoryName, in.SubCategoryName, in.ShopName, in.ShopLogo, in.ForbidReason)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateOutProductCompanyProductsReply{
		Code: 200,
		Data: &v1.UpdateOutProductCompanyProductsReply_Data{},
	}, nil
}

func (cs *CompanyService) UpdateCompanyProductByProductIds(ctx context.Context, in *v1.UpdateCompanyProductByProductIdsRequest) (*v1.UpdateCompanyProductByProductIdsReply, error) {
	err := cs.cprouc.UpdateCompanyProductByProductIds(ctx, in.ProductIds)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateCompanyProductByProductIdsReply{
		Code: 200,
		Data: &v1.UpdateCompanyProductByProductIdsReply_Data{},
	}, nil
}

func (cs *CompanyService) ParseCompanyProducts(ctx context.Context, in *v1.ParseCompanyProductsRequest) (*v1.ParseCompanyProductsReply, error) {
	companyProduct, err := cs.cprouc.ParseCompanyProducts(ctx, in.Content)

	if err != nil {
		return nil, err
	}

	productImg := ""

	if len(companyProduct.ProductImgs) > 0 {
		productImg = companyProduct.ProductImgs[0]
	}

	return &v1.ParseCompanyProductsReply{
		Code: 200,
		Data: &v1.ParseCompanyProductsReply_Data{
			ProductId:             companyProduct.ProductOutId,
			ProductName:           companyProduct.ProductName,
			ProductImg:            productImg,
			ProductPrice:          companyProduct.ProductPrice,
			IsTop:                 uint32(companyProduct.IsTop),
			PureCommission:        companyProduct.PureCommission,
			PureServiceCommission: companyProduct.PureServiceCommission,
			CommonCommission:      fmt.Sprintf("%.f", tool.Decimal(float64(companyProduct.CommissionRatio), 2)),
			InvestmentRatio:       fmt.Sprintf("%.f", tool.Decimal(float64(companyProduct.InvestmentRatio), 2)),
			IsHot:                 uint32(companyProduct.IsHot),
			TotalSale:             companyProduct.TotalSale,
			IsTask:                uint32(companyProduct.IsTask),
		},
	}, nil
}

func (cs *CompanyService) VerificationCompanyProducts(ctx context.Context, in *v1.VerificationCompanyProductsRequest) (*v1.VerificationCompanyProductsReply, error) {
	err := cs.cprouc.VerificationCompanyProducts(ctx, in.ProductId)

	if err != nil {
		return nil, err
	}

	return &v1.VerificationCompanyProductsReply{
		Code: 200,
		Data: &v1.VerificationCompanyProductsReply_Data{},
	}, nil
}

func (cs *CompanyService) UploadPartCompanyProducts(ctx context.Context, in *v1.UploadPartCompanyProductsRequest) (*v1.UploadPartCompanyProductsReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := cs.cprouc.UploadPartCompanyProducts(ctx, in.PartNumber, in.TotalPart, in.ContentLength, in.UploadId, in.Content)

	if err != nil {
		return nil, err
	}

	return &v1.UploadPartCompanyProductsReply{
		Code: 200,
		Data: &v1.UploadPartCompanyProductsReply_Data{},
	}, nil
}

func (cs *CompanyService) CompleteUploadCompanyProducts(ctx context.Context, in *v1.CompleteUploadCompanyProductsRequest) (*v1.CompleteUploadCompanyProductsReply, error) {
	staticUrl, err := cs.cprouc.CompleteUploadCompanyProducts(ctx, in.UploadId)

	if err != nil {
		return nil, err
	}

	return &v1.CompleteUploadCompanyProductsReply{
		Code: 200,
		Data: &v1.CompleteUploadCompanyProductsReply_Data{
			StaticUrl: staticUrl,
		},
	}, nil
}

func (cs *CompanyService) AbortUploadCompanyProducts(ctx context.Context, in *v1.AbortUploadCompanyProductsRequest) (*v1.AbortUploadCompanyProductsReply, error) {
	err := cs.cprouc.AbortUploadCompanyProducts(ctx, in.UploadId)

	if err != nil {
		return nil, err
	}

	return &v1.AbortUploadCompanyProductsReply{
		Code: 200,
		Data: &v1.AbortUploadCompanyProductsReply_Data{},
	}, nil
}

func (cs *CompanyService) DeleteCompanyProducts(ctx context.Context, in *v1.DeleteCompanyProductsRequest) (*v1.DeleteCompanyProductsReply, error) {
	err := cs.cprouc.DeleteCompanyProducts(ctx, in.ProductId)

	if err != nil {
		return nil, err
	}

	return &v1.DeleteCompanyProductsReply{
		Code: 200,
		Data: &v1.DeleteCompanyProductsReply_Data{},
	}, nil
}

func (cs *CompanyService) SyncCompanyProducts(ctx context.Context, in *empty.Empty) (*v1.SyncCompanyProductsReply, error) {
	cs.log.Infof("同步商品数据, 开始时间 %s \n", time.Now())

	ctx = context.Background()

	if err := cs.cprouc.SyncCompanyProducts(ctx); err != nil {
		return nil, err
	}

	cs.log.Infof("同步商品数据, 结束时间 %s \n", time.Now())

	return &v1.SyncCompanyProductsReply{
		Code: 200,
		Data: &v1.SyncCompanyProductsReply_Data{},
	}, nil
}

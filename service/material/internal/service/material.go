package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "material/api/material/v1"
	"material/internal/biz"
	"material/internal/pkg/tool"
	"math"
	"time"
)

func (ms *MaterialService) ListMaterials(ctx context.Context, in *v1.ListMaterialsRequest) (*v1.ListMaterialsReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if in.IsShowCollect == 1 {
		if in.CompanyId == 0 {
			return nil, biz.MaterialValidatorError
		}

		if len(in.Phone) == 0 {
			return nil, biz.MaterialValidatorError
		}
	}

	materials, err := ms.muc.ListMaterials(ctx, in.PageNum, in.PageSize, in.CompanyId, in.ProductId, uint8(in.IsShowCollect), in.Phone, in.Category, in.Keyword, in.Search, in.Msort, in.Mplatform)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListMaterialsReply_Materials, 0)

	for _, material := range materials.List {
		list = append(list, &v1.ListMaterialsReply_Materials{
			VideoId:            material.VideoId,
			VideoName:          material.VideoName,
			VideoUrl:           material.VideoUrl,
			VideoCover:         material.VideoCover,
			VideoLike:          material.VideoLike,
			VideoLikeShowA:     material.VideoLikeShowA,
			VideoLikeShowB:     material.VideoLikeShowB,
			AwemeId:            material.AwemeId,
			AwemeName:          material.AwemeName,
			AwemeAccount:       material.AwemeAccount,
			AwemeFollowers:     material.AwemeFollowers,
			AwemeFollowersShow: material.AwemeFollowersShow,
			AwemeImg:           material.AwemeImg,
			AwemeLandingPage:   material.AwemeLandingPage,
			ProductId:          material.ProductId,
			ProductName:        material.ProductName,
			ProductImg:         material.ProductImg,
			ProductLandingPage: material.ProductLandingPage,
			ProductPrice:       material.ProductPrice,
			IsHot:              uint32(material.IsHot),
			TotalItemNum:       material.TotalItemNum,
			PlatformName:       material.PlatformName,
			IsCollect:          uint32(material.IsCollect),
			UpdateTime:         tool.TimeToString("2006/01/02 15:04", material.UpdateDay),
			UpdateTimeF:        tool.TimeToString("2006-01-02 15:04:05", material.UpdateDay),
		})
	}

	totalPage := uint64(math.Ceil(float64(materials.Total) / float64(materials.PageSize)))

	return &v1.ListMaterialsReply{
		Code: 200,
		Data: &v1.ListMaterialsReply_Data{
			PageNum:   materials.PageNum,
			PageSize:  materials.PageSize,
			Total:     materials.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ms *MaterialService) ListProducts(ctx context.Context, in *v1.ListProductsRequest) (*v1.ListProductsReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	products, err := ms.muc.ListProducts(ctx, in.PageNum, in.PageSize, in.Category, in.Keyword, in.Search, in.Msort, in.Mplatform)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListProductsReply_Products, 0)

	for _, product := range products.List {
		list = append(list, &v1.ListProductsReply_Products{
			ProductId:          product.ProductId,
			ProductName:        product.ProductName,
			ProductImg:         product.ProductImg,
			ProductLandingPage: product.ProductLandingPage,
			ProductPrice:       product.ProductPrice,
			IsHot:              uint32(product.IsHot),
			VideoLike:          product.VideoLike,
			VideoLikeShowA:     product.VideoLikeShowA,
			VideoLikeShowB:     product.VideoLikeShowB,
			Awemes:             product.Awemes,
			Videos:             product.Videos,
			PlatformName:       product.PlatformName,
		})
	}

	totalPage := uint64(math.Ceil(float64(products.Total) / float64(products.PageSize)))

	return &v1.ListProductsReply{
		Code: 200,
		Data: &v1.ListProductsReply_Data{
			PageNum:   products.PageNum,
			PageSize:  products.PageSize,
			Total:     products.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (ms *MaterialService) ListAwemesByProductId(ctx context.Context, in *v1.ListAwemesByProductIdRequest) (*v1.ListAwemesByProductIdReply, error) {
	materials, err := ms.muc.ListAwemesByProductId(ctx, in.ProductId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListAwemesByProductIdReply_Aweme, 0)

	for _, material := range materials {
		list = append(list, &v1.ListAwemesByProductIdReply_Aweme{
			AwemeName:    material.AwemeName,
			AwemeAccount: material.AwemeAccount,
			AwemeImg:     material.AwemeImg,
		})
	}

	return &v1.ListAwemesByProductIdReply{
		Code: 200,
		Data: &v1.ListAwemesByProductIdReply_Data{
			List: list,
		},
	}, nil
}

func (ms *MaterialService) ListSelectMaterials(ctx context.Context, in *emptypb.Empty) (*v1.ListSelectMaterialsReply, error) {
	selects, err := ms.muc.ListSelectMaterials(ctx)

	if err != nil {
		return nil, err
	}

	category := make([]*v1.ListSelectMaterialsReply_Category, 0)
	msort := make([]*v1.ListSelectMaterialsReply_Msort, 0)
	mplatform := make([]*v1.ListSelectMaterialsReply_Mplatform, 0)
	search := make([]*v1.ListSelectMaterialsReply_Search, 0)

	for _, lcategory := range selects.Category {
		childList := make([]*v1.ListSelectMaterialsReply_ChildCategory, 0)

		for _, clcategory := range lcategory.ChildList {
			childList = append(childList, &v1.ListSelectMaterialsReply_ChildCategory{
				Key:   clcategory.Key,
				Value: clcategory.Value,
			})
		}

		category = append(category, &v1.ListSelectMaterialsReply_Category{
			Key:       lcategory.Key,
			Value:     lcategory.Value,
			ChildList: childList,
		})
	}

	for _, lsort := range selects.Msort {
		msort = append(msort, &v1.ListSelectMaterialsReply_Msort{
			Key:   lsort.Key,
			Value: lsort.Value,
		})
	}

	for _, lplatform := range selects.Mplatform {
		mplatform = append(mplatform, &v1.ListSelectMaterialsReply_Mplatform{
			Key:   lplatform.Key,
			Value: lplatform.Value,
		})
	}

	for _, lsearch := range selects.Search {
		search = append(search, &v1.ListSelectMaterialsReply_Search{
			Key:   lsearch.Key,
			Value: lsearch.Value,
		})
	}

	return &v1.ListSelectMaterialsReply{
		Code: 200,
		Data: &v1.ListSelectMaterialsReply_Data{
			Category:  category,
			Msort:     msort,
			Mplatform: mplatform,
			Search:    search,
		},
	}, nil
}

func (ms *MaterialService) GetDownUrlVideoUrls(ctx context.Context, in *v1.GetDownUrlVideoUrlsRequest) (*v1.GetDownUrlVideoUrlsReply, error) {
	videoUrl, err := ms.muc.GetDownUrlVideoUrls(ctx, in.VideoId)

	if err != nil {
		return nil, err
	}

	return &v1.GetDownUrlVideoUrlsReply{
		Code: 200,
		Data: &v1.GetDownUrlVideoUrlsReply_Data{
			VideoUrl: videoUrl,
		},
	}, nil
}

func (ms *MaterialService) GetVideoUrls(ctx context.Context, in *v1.GetVideoUrlsRequest) (*v1.GetVideoUrlsReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	videoUrl, err := ms.muc.GetVideoUrls(ctx, in.VideoId)

	if err != nil {
		return nil, err
	}

	imageUrls := make([]*v1.GetVideoUrlsReply_ImageUrls, 0)

	for _, imageUrl := range videoUrl.Data.ImageUrls {
		imageUrls = append(imageUrls, &v1.GetVideoUrlsReply_ImageUrls{
			ImageUrl: imageUrl.ImageUrl,
		})
	}

	return &v1.GetVideoUrlsReply{
		Code: 200,
		Data: &v1.GetVideoUrlsReply_Data{
			VideoJumpUrl: videoUrl.Data.VideoUrl,
			ImageUrls:    imageUrls,
		},
	}, nil
}

func (ms *MaterialService) GetPromotions(ctx context.Context, in *v1.GetPromotionsRequest) (*v1.GetPromotionsReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	promotion, err := ms.muc.GetPromotions(ctx, in.PageNum, in.PageSize, in.PromotionId, in.Ptype)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.GetPromotionsReply_Materials, 0)

	for _, material := range promotion.List {
		list = append(list, &v1.GetPromotionsReply_Materials{
			VideoId:            material.VideoId,
			VideoName:          material.VideoName,
			VideoUrl:           material.VideoUrl,
			VideoCover:         material.VideoCover,
			AwemeId:            material.AwemeId,
			AwemeName:          material.AwemeName,
			AwemeAccount:       material.AwemeAccount,
			AwemeFollowers:     material.AwemeFollowers,
			AwemeFollowersShow: material.AwemeFollowersShow,
			AwemeImg:           material.AwemeImg,
			AwemeLandingPage:   material.AwemeLandingPage,
			ProductId:          material.ProductId,
			ProductName:        material.ProductName,
			ProductImg:         material.ProductImg,
			ProductLandingPage: material.ProductLandingPage,
			ProductPrice:       material.ProductPrice,
			VideoLike:          material.VideoLike,
			VideoLikeShowA:     material.VideoLikeShowA,
			VideoLikeShowB:     material.VideoLikeShowB,
			TotalItemNum:       material.TotalItemNum,
			UpdateTime:         tool.TimeToString("2006/01/02 15:04", material.UpdateDay),
		})
	}

	totalPage := uint64(math.Ceil(float64(promotion.Total) / float64(promotion.PageSize)))

	industry := make([]*v1.GetPromotionsReply_Industry, 0)

	for _, lindustry := range promotion.Industry {
		industry = append(industry, &v1.GetPromotionsReply_Industry{
			IndustryName:  lindustry.IndustryName,
			IndustryRatio: lindustry.IndustryRatio,
		})
	}

	return &v1.GetPromotionsReply{
		Code: 200,
		Data: &v1.GetPromotionsReply_Data{
			PromotionId:            promotion.PromotionId,
			PromotionName:          promotion.PromotionName,
			PromotionType:          promotion.PromotionType,
			PromotionAccount:       promotion.PromotionAccount,
			PromotionImg:           promotion.PromotionImg,
			PromotionLandingPage:   promotion.PromotionLandingPage,
			PromotionFollowers:     promotion.PromotionFollowers,
			PromotionFollowersShow: promotion.PromotionFollowersShow,
			PromotionPrice:         promotion.PromotionPrice,
			PromotionPlatformName:  promotion.PromotionPlatformName,
			IndustryName:           promotion.IndustryName,
			ShopName:               promotion.ShopName,
			ShopLogo:               promotion.ShopLogo,
			Industry:               industry,
			PageNum:                promotion.PageNum,
			PageSize:               promotion.PageSize,
			Total:                  promotion.Total,
			TotalPage:              totalPage,
			List:                   list,
		},
	}, nil
}

func (ms *MaterialService) GetMaterials(ctx context.Context, in *v1.GetMaterialsRequest) (*v1.GetMaterialsReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	material, companyProduct, err := ms.muc.GetMaterials(ctx, in.VideoId)

	if err != nil {
		return nil, err
	}

	return &v1.GetMaterialsReply{
		Code: 200,
		Data: &v1.GetMaterialsReply_Data{
			VideoId:               material.VideoId,
			VideoName:             material.VideoName,
			VideoUrl:              material.VideoUrl,
			VideoCover:            material.VideoCover,
			VideoLike:             material.VideoLike,
			VideoLikeShowA:        material.VideoLikeShowA,
			VideoLikeShowB:        material.VideoLikeShowB,
			AwemeId:               material.AwemeId,
			AwemeName:             material.AwemeName,
			AwemeAccount:          material.AwemeAccount,
			AwemeFollowers:        material.AwemeFollowers,
			AwemeFollowersShow:    material.AwemeFollowersShow,
			AwemeImg:              material.AwemeImg,
			AwemeLandingPage:      material.AwemeLandingPage,
			ProductId:             material.ProductId,
			ProductName:           material.ProductName,
			ProductImg:            material.ProductImg,
			ProductLandingPage:    material.ProductLandingPage,
			ProductPrice:          material.ProductPrice,
			PureCommission:        companyProduct.Data.PureCommission,
			PureServiceCommission: companyProduct.Data.PureServiceCommission,
			CommonCommission:      companyProduct.Data.CommonCommission,
			IsHot:                 uint32(material.IsHot),
			TotalItemNum:          material.TotalItemNum,
			PlatformName:          material.PlatformName,
			IsCollect:             uint32(material.IsCollect),
			UpdateTime:            tool.TimeToString("2006/01/02 15:04", material.UpdateDay),
			UpdateTimeF:           tool.TimeToString("2006/01/02 15:04", material.UpdateDay),
		},
	}, nil
}

func (ms *MaterialService) GetUploadIdMaterials(ctx context.Context, in *v1.GetUploadIdMaterialsRequest) (*v1.GetUploadIdMaterialsReply, error) {
	uploadId, err := ms.muc.GetUploadIdMaterials(ctx, in.Suffix)

	if err != nil {
		return nil, err
	}

	return &v1.GetUploadIdMaterialsReply{
		Code: 200,
		Data: &v1.GetUploadIdMaterialsReply_Data{
			UploadId: uploadId,
		},
	}, nil
}

func (ms *MaterialService) GetFileSizeMaterials(ctx context.Context, in *v1.GetFileSizeMaterialsRequest) (*v1.GetFileSizeMaterialsReply, error) {
	fileSize, err := ms.muc.GetFileSizeMaterials(ctx, in.MaterialUrl)

	if err != nil {
		return nil, err
	}

	return &v1.GetFileSizeMaterialsReply{
		Code: 200,
		Data: &v1.GetFileSizeMaterialsReply_Data{
			FileSize: fileSize,
		},
	}, nil
}

func (ms *MaterialService) GetIsTopMaterials(ctx context.Context, in *v1.GetIsTopMaterialsRequest) (*v1.GetIsTopMaterialsReply, error) {
	var isTop uint32 = 0

	if _, err := ms.muc.GetIsTopMaterials(ctx, in.ProductId); err == nil {
		isTop = 1
	}

	return &v1.GetIsTopMaterialsReply{
		Code: 200,
		Data: &v1.GetIsTopMaterialsReply_Data{
			IsTop: isTop,
		},
	}, nil
}

func (ms *MaterialService) StatisticsMaterials(ctx context.Context, in *emptypb.Empty) (*v1.StatisticsMaterialsReply, error) {
	statistics, err := ms.muc.StatisticsMaterials(ctx)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsMaterialsReply_Statistics, 0)

	for _, statistic := range statistics.Statistics {
		list = append(list, &v1.StatisticsMaterialsReply_Statistics{
			Key:   statistic.Key,
			Value: statistic.Value,
		})
	}

	return &v1.StatisticsMaterialsReply{
		Code: 200,
		Data: &v1.StatisticsMaterialsReply_Data{
			Statistics: list,
		},
	}, nil
}

func (ms *MaterialService) CreateMaterials(ctx context.Context, in *v1.CreateMaterialsRequest) (*v1.CreateMaterialsReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := ms.muc.CreateMaterials(ctx, in.ProductId, in.VideoId, in.MaterialUrl, in.FileName, in.MaterialType, in.FileType)

	if err != nil {
		return nil, err
	}

	return &v1.CreateMaterialsReply{
		Code: 200,
		Data: &v1.CreateMaterialsReply_Data{},
	}, nil
}

func (ms *MaterialService) UploadMaterials(ctx context.Context, in *v1.UploadMaterialsRequest) (*v1.UploadMaterialsReply, error) {
	err := ms.muc.UploadMaterials(ctx, in.VideoId, in.VideoUrl)

	if err != nil {
		return nil, err
	}

	return &v1.UploadMaterialsReply{
		Code: 200,
		Data: &v1.UploadMaterialsReply_Data{},
	}, nil
}

func (ms *MaterialService) UploadPartMaterials(ctx context.Context, in *v1.UploadPartMaterialsRequest) (*v1.UploadPartMaterialsReply, error) {
	/*ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()*/

	ctx = context.Background()

	err := ms.muc.UploadPartMaterials(ctx, in.PartNumber, in.TotalPart, in.ContentLength, in.UploadId, in.Content)

	if err != nil {
		return nil, err
	}

	return &v1.UploadPartMaterialsReply{
		Code: 200,
		Data: &v1.UploadPartMaterialsReply_Data{},
	}, nil
}

func (ms *MaterialService) CompleteUploadMaterials(ctx context.Context, in *v1.CompleteUploadMaterialsRequest) (*v1.CompleteUploadMaterialsReply, error) {
	staticUrl, err := ms.muc.CompleteUploadMaterials(ctx, in.UploadId)

	if err != nil {
		return nil, err
	}

	return &v1.CompleteUploadMaterialsReply{
		Code: 200,
		Data: &v1.CompleteUploadMaterialsReply_Data{
			StaticUrl: staticUrl,
		},
	}, nil
}

func (ms *MaterialService) AbortUploadMaterials(ctx context.Context, in *v1.AbortUploadMaterialsRequest) (*v1.AbortUploadMaterialsReply, error) {
	err := ms.muc.AbortUploadMaterials(ctx, in.UploadId)

	if err != nil {
		return nil, err
	}

	return &v1.AbortUploadMaterialsReply{
		Code: 200,
		Data: &v1.AbortUploadMaterialsReply_Data{},
	}, nil
}

func (ms *MaterialService) DownMaterials(ctx context.Context, in *v1.DownMaterialsRequest) (*v1.DownMaterialsReply, error) {
	if in.DownType == "cloud" {
		if in.CompanyMaterialId == 0 {
			return nil, biz.MaterialValidatorError
		}
	}

	err := ms.muc.DownMaterials(ctx, in.CompanyId, in.VideoId, in.CompanyMaterialId, in.DownType)

	if err != nil {
		return nil, err
	}

	return &v1.DownMaterialsReply{
		Code: 200,
		Data: &v1.DownMaterialsReply_Data{},
	}, nil
}

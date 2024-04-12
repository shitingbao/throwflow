package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "interface/api/interface/v1"
	"math"
	"time"
)

func (is *InterfaceService) ListMaterials(ctx context.Context, in *v1.ListMaterialsRequest) (*v1.ListMaterialsReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "material")

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	materials, err := is.muc.ListMaterials(ctx, in.PageNum, in.PageSize, companyUser.Data.CurrentCompanyId, companyUser.Data.Phone, in.Category, in.Keyword, in.Search, in.Msort, in.Mplatform)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListMaterialsReply_Materials, 0)

	for _, material := range materials.Data.List {
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
			IsHot:              material.IsHot,
			PlatformName:       material.PlatformName,
			IsCollect:          material.IsCollect,
		})
	}

	totalPage := uint64(math.Ceil(float64(materials.Data.Total) / float64(materials.Data.PageSize)))

	if totalPage > 20 {
		totalPage = 20
	}

	return &v1.ListMaterialsReply{
		Code: 200,
		Data: &v1.ListMaterialsReply_Data{
			PageNum:   materials.Data.PageNum,
			PageSize:  materials.Data.PageSize,
			Total:     materials.Data.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (is *InterfaceService) ListMiniMaterials(ctx context.Context, in *v1.ListMiniMaterialsRequest) (*v1.ListMiniMaterialsReply, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	materials, err := is.muc.ListMiniMaterials(ctx, in.PageNum, in.PageSize)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListMiniMaterialsReply_Materials, 0)

	for _, material := range materials.Data.List {
		list = append(list, &v1.ListMiniMaterialsReply_Materials{
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
			IsHot:              material.IsHot,
			TotalItemNum:       material.TotalItemNum,
			PlatformName:       material.PlatformName,
			UpdateTime:         material.UpdateTime,
		})
	}

	totalPage := uint64(math.Ceil(float64(materials.Data.Total) / float64(materials.Data.PageSize)))

	if totalPage > 20 {
		totalPage = 20
	}

	return &v1.ListMiniMaterialsReply{
		Code: 200,
		Data: &v1.ListMiniMaterialsReply_Data{
			PageNum:   materials.Data.PageNum,
			PageSize:  materials.Data.PageSize,
			Total:     materials.Data.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (is *InterfaceService) ListMiniMaterialProducts(ctx context.Context, in *v1.ListMiniMaterialProductsRequest) (*v1.ListMiniMaterialProductsReply, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	materials, err := is.muc.ListMiniMaterialProducts(ctx, in.PageNum, in.PageSize, in.ProductId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListMiniMaterialProductsReply_Materials, 0)

	for _, material := range materials.Data.List {
		list = append(list, &v1.ListMiniMaterialProductsReply_Materials{
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
			IsHot:              material.IsHot,
			TotalItemNum:       material.TotalItemNum,
			PlatformName:       material.PlatformName,
			UpdateTime:         material.UpdateTime,
		})
	}

	totalPage := uint64(math.Ceil(float64(materials.Data.Total) / float64(materials.Data.PageSize)))

	if totalPage > 20 {
		totalPage = 20
	}

	return &v1.ListMiniMaterialProductsReply{
		Code: 200,
		Data: &v1.ListMiniMaterialProductsReply_Data{
			PageNum:   materials.Data.PageNum,
			PageSize:  materials.Data.PageSize,
			Total:     materials.Data.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (is *InterfaceService) ListMaterialProducts(ctx context.Context, in *v1.ListMaterialProductsRequest) (*v1.ListMaterialProductsReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "material"); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	materialProducts, err := is.muc.ListMaterialProducts(ctx, in.PageNum, in.PageSize, in.Category, in.Keyword, "product", in.Msort, in.Mplatform)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListMaterialProductsReply_Products, 0)

	for _, materialProduct := range materialProducts.Data.List {
		list = append(list, &v1.ListMaterialProductsReply_Products{
			ProductId:          materialProduct.ProductId,
			ProductName:        materialProduct.ProductName,
			ProductImg:         materialProduct.ProductImg,
			ProductLandingPage: materialProduct.ProductLandingPage,
			ProductPrice:       materialProduct.ProductPrice,
			IsHot:              materialProduct.IsHot,
			VideoLike:          materialProduct.VideoLike,
			VideoLikeShowA:     materialProduct.VideoLikeShowA,
			VideoLikeShowB:     materialProduct.VideoLikeShowB,
			Awemes:             materialProduct.Awemes,
			Videos:             materialProduct.Videos,
			PlatformName:       materialProduct.PlatformName,
		})
	}

	totalPage := uint64(math.Ceil(float64(materialProducts.Data.Total) / float64(materialProducts.Data.PageSize)))

	if totalPage > 20 {
		totalPage = 20
	}

	return &v1.ListMaterialProductsReply{
		Code: 200,
		Data: &v1.ListMaterialProductsReply_Data{
			PageNum:   materialProducts.Data.PageNum,
			PageSize:  materialProducts.Data.PageSize,
			Total:     materialProducts.Data.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (is *InterfaceService) ListCollectMaterials(ctx context.Context, in *v1.ListCollectMaterialsRequest) (*v1.ListCollectMaterialsReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "material")

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	materials, err := is.muc.ListCollectMaterials(ctx, in.PageNum, in.PageSize, companyUser.Data.CurrentCompanyId, companyUser.Data.Phone, in.Category, in.Keyword, in.Search, in.Msort, in.Mplatform)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListCollectMaterialsReply_Materials, 0)

	for _, material := range materials.Data.List {
		list = append(list, &v1.ListCollectMaterialsReply_Materials{
			VideoId:            material.VideoId,
			VideoName:          material.VideoName,
			VideoUrl:           material.VideoUrl,
			VideoCover:         material.VideoCover,
			VideoCategory:      material.VideoCategory,
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
			PlatformName:       material.PlatformName,
			UpdateDay:          material.UpdateDay,
			IsCollect:          material.IsCollect,
		})
	}

	totalPage := uint64(math.Ceil(float64(materials.Data.Total) / float64(materials.Data.PageSize)))

	if totalPage > 20 {
		totalPage = 20
	}

	return &v1.ListCollectMaterialsReply{
		Code: 200,
		Data: &v1.ListCollectMaterialsReply_Data{
			PageNum:   materials.Data.PageNum,
			PageSize:  materials.Data.PageSize,
			Total:     materials.Data.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (is *InterfaceService) ListSelectMaterials(ctx context.Context, in *emptypb.Empty) (*v1.ListSelectMaterialsReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "material"); err != nil {
		return nil, err
	}

	selects, err := is.muc.ListSelectMaterials(ctx)

	if err != nil {
		return nil, err
	}

	category := make([]*v1.ListSelectMaterialsReply_Category, 0)
	msort := make([]*v1.ListSelectMaterialsReply_Msort, 0)
	mplatform := make([]*v1.ListSelectMaterialsReply_Mplatform, 0)
	search := make([]*v1.ListSelectMaterialsReply_Search, 0)

	for _, lcategory := range selects.Data.Category {
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

	for _, lmsort := range selects.Data.Msort {
		msort = append(msort, &v1.ListSelectMaterialsReply_Msort{
			Key:   lmsort.Key,
			Value: lmsort.Value,
		})
	}

	for _, lplatform := range selects.Data.Mplatform {
		mplatform = append(mplatform, &v1.ListSelectMaterialsReply_Mplatform{
			Key:   lplatform.Key,
			Value: lplatform.Value,
		})
	}

	for _, lsearch := range selects.Data.Search {
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

func (is *InterfaceService) StatisticsMaterials(ctx context.Context, in *emptypb.Empty) (*v1.StatisticsMaterialsReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "material"); err != nil {
		return nil, err
	}

	statistics, err := is.muc.StatisticsMaterials(ctx)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsMaterialsReply_Statistics, 0)

	for _, statistic := range statistics.Data.Statistics {
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

func (is *InterfaceService) StatisticsMiniMaterials(ctx context.Context, in *emptypb.Empty) (*v1.StatisticsMaterialsReply, error) {
	statistics, err := is.muc.StatisticsMaterials(ctx)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.StatisticsMaterialsReply_Statistics, 0)

	for _, statistic := range statistics.Data.Statistics {
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

func (is *InterfaceService) GetDownVideoUrls(ctx context.Context, in *v1.GetDownVideoUrlsRequest) (*v1.GetDownVideoUrlsReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "material"); err != nil {
		return nil, err
	}

	videoUrl, err := is.muc.GetDownVideoUrls(ctx, in.VideoId)

	if err != nil {
		return nil, err
	}

	return &v1.GetDownVideoUrlsReply{
		Code: 200,
		Data: &v1.GetDownVideoUrlsReply_Data{
			VideoUrl: videoUrl.Data.VideoUrl,
		},
	}, nil
}

func (is *InterfaceService) GetMiniDownVideoUrls(ctx context.Context, in *v1.GetMiniDownVideoUrlsRequest) (*v1.GetMiniDownVideoUrlsReply, error) {
	if _, err := is.verifyMiniUserLogin(ctx); err != nil {
		return nil, err
	}

	videoUrl, err := is.muc.GetMiniDownVideoUrls(ctx, in.VideoId)

	if err != nil {
		return nil, err
	}

	return &v1.GetMiniDownVideoUrlsReply{
		Code: 200,
		Data: &v1.GetMiniDownVideoUrlsReply_Data{
			VideoUrl: videoUrl.Data.VideoUrl,
		},
	}, nil
}

func (is *InterfaceService) GetVideoUrls(ctx context.Context, in *v1.GetVideoUrlsRequest) (*v1.GetVideoUrlsReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "material"); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	videoUrl, err := is.muc.GetVideoUrls(ctx, in.VideoId)

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
			VideoJumpUrl: videoUrl.Data.VideoJumpUrl,
			ImageUrls:    imageUrls,
		},
	}, nil
}

func (is *InterfaceService) GetPromotions(ctx context.Context, in *v1.GetPromotionsRequest) (*v1.GetPromotionsReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "material"); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	promotion, err := is.muc.GetPromotions(ctx, in.PageNum, in.PageSize, in.PromotionId, in.Ptype)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.GetPromotionsReply_Materials, 0)

	for _, material := range promotion.Data.List {
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
			UpdateTime:         material.UpdateTime,
		})
	}

	industry := make([]*v1.GetPromotionsReply_Industry, 0)

	for _, lindustry := range promotion.Data.Industry {
		industry = append(industry, &v1.GetPromotionsReply_Industry{
			IndustryName:  lindustry.IndustryName,
			IndustryRatio: lindustry.IndustryRatio,
		})
	}

	return &v1.GetPromotionsReply{
		Code: 200,
		Data: &v1.GetPromotionsReply_Data{
			PromotionId:            promotion.Data.PromotionId,
			PromotionName:          promotion.Data.PromotionName,
			PromotionType:          promotion.Data.PromotionType,
			PromotionAccount:       promotion.Data.PromotionAccount,
			PromotionImg:           promotion.Data.PromotionImg,
			PromotionLandingPage:   promotion.Data.PromotionLandingPage,
			PromotionFollowers:     promotion.Data.PromotionFollowers,
			PromotionFollowersShow: promotion.Data.PromotionFollowersShow,
			PromotionPrice:         promotion.Data.PromotionPrice,
			PromotionPlatformName:  promotion.Data.PromotionPlatformName,
			IndustryName:           promotion.Data.IndustryName,
			ShopName:               promotion.Data.ShopName,
			ShopLogo:               promotion.Data.ShopLogo,
			Industry:               industry,
			PageNum:                promotion.Data.PageNum,
			PageSize:               promotion.Data.PageSize,
			Total:                  promotion.Data.Total,
			TotalPage:              promotion.Data.TotalPage,
			List:                   list,
		},
	}, nil
}

func (is *InterfaceService) GetMiniMaterials(ctx context.Context, in *v1.GetMiniMaterialsRequest) (*v1.GetMiniMaterialsReply, error) {
	material, err := is.muc.GetMiniMaterials(ctx, in.VideoId)

	if err != nil {
		return nil, err
	}

	return &v1.GetMiniMaterialsReply{
		Code: 200,
		Data: &v1.GetMiniMaterialsReply_Data{
			VideoId:               material.Data.VideoId,
			VideoName:             material.Data.VideoName,
			VideoUrl:              material.Data.VideoUrl,
			VideoCover:            material.Data.VideoCover,
			VideoLike:             material.Data.VideoLike,
			VideoLikeShowA:        material.Data.VideoLikeShowA,
			VideoLikeShowB:        material.Data.VideoLikeShowB,
			AwemeId:               material.Data.AwemeId,
			AwemeName:             material.Data.AwemeName,
			AwemeAccount:          material.Data.AwemeAccount,
			AwemeFollowers:        material.Data.AwemeFollowers,
			AwemeFollowersShow:    material.Data.AwemeFollowersShow,
			AwemeImg:              material.Data.AwemeImg,
			AwemeLandingPage:      material.Data.AwemeLandingPage,
			ProductId:             material.Data.ProductId,
			ProductName:           material.Data.ProductName,
			ProductImg:            material.Data.ProductImg,
			ProductLandingPage:    material.Data.ProductLandingPage,
			ProductPrice:          material.Data.ProductPrice,
			PureCommission:        material.Data.PureCommission,
			PureServiceCommission: material.Data.PureServiceCommission,
			CommonCommission:      material.Data.CommonCommission,
			IsHot:                 material.Data.IsHot,
			TotalItemNum:          material.Data.TotalItemNum,
			PlatformName:          material.Data.PlatformName,
			IsCollect:             material.Data.IsCollect,
			UpdateTime:            material.Data.UpdateTime,
			UpdateTimeF:           material.Data.UpdateTimeF,
		},
	}, nil
}

func (is *InterfaceService) UpdateCollects(ctx context.Context, in *v1.UpdateCollectsRequest) (*v1.UpdateCollectsReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "material")

	if err != nil {
		return nil, err
	}

	if err := is.muc.UpdateCollects(ctx, companyUser.Data.CurrentCompanyId, in.VideoId, companyUser.Data.Phone); err != nil {
		return nil, err
	}

	return &v1.UpdateCollectsReply{
		Code: 200,
		Data: &v1.UpdateCollectsReply_Data{},
	}, nil
}

func (is *InterfaceService) DownMaterials(ctx context.Context, in *v1.DownMaterialsRequest) (*v1.DownMaterialsReply, error) {
	companyUser, err := is.verifyLogin(ctx, true, false, "material")

	if err != nil {
		return nil, err
	}

	if err := is.muc.DownMaterials(ctx, companyUser.Data.CurrentCompanyId, in.VideoId, in.CompanyMaterialId, in.DownType); err != nil {
		return nil, err
	}

	return &v1.DownMaterialsReply{
		Code: 200,
		Data: &v1.DownMaterialsReply_Data{},
	}, nil
}

func (is *InterfaceService) DownMiniMaterials(ctx context.Context, in *v1.DownMiniMaterialsRequest) (*v1.DownMiniMaterialsReply, error) {
	if _, err := is.verifyMiniUserLogin(ctx); err != nil {
		return nil, err
	}

	if err := is.muc.DownMiniMaterials(ctx, in.VideoId, in.DownType); err != nil {
		return nil, err
	}

	return &v1.DownMiniMaterialsReply{
		Code: 200,
		Data: &v1.DownMiniMaterialsReply_Data{},
	}, nil
}

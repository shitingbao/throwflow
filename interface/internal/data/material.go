package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "interface/api/service/material/v1"
	"interface/internal/biz"
)

type materialRepo struct {
	data *Data
	log  *log.Helper
}

func NewMaterialRepo(data *Data, logger log.Logger) biz.MaterialRepo {
	return &materialRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (mr *materialRepo) Get(ctx context.Context, videoId uint64) (*v1.GetMaterialsReply, error) {
	material, err := mr.data.materialuc.GetMaterials(ctx, &v1.GetMaterialsRequest{
		VideoId: videoId,
	})

	if err != nil {
		return nil, err
	}

	return material, err
}

func (mr *materialRepo) GetDownUrl(ctx context.Context, videoId uint64) (*v1.GetDownUrlVideoUrlsReply, error) {
	videoUrl, err := mr.data.materialuc.GetDownUrlVideoUrls(ctx, &v1.GetDownUrlVideoUrlsRequest{
		VideoId: videoId,
	})

	if err != nil {
		return nil, err
	}

	return videoUrl, err
}

func (mr *materialRepo) GetVideoUrl(ctx context.Context, videoId uint64) (*v1.GetVideoUrlsReply, error) {
	videoUrl, err := mr.data.materialuc.GetVideoUrls(ctx, &v1.GetVideoUrlsRequest{
		VideoId: videoId,
	})

	if err != nil {
		return nil, err
	}

	return videoUrl, err
}

func (mr *materialRepo) GetPromotion(ctx context.Context, pageNum, pageSize, promotionId uint64, ptype string) (*v1.GetPromotionsReply, error) {
	promotion, err := mr.data.materialuc.GetPromotions(ctx, &v1.GetPromotionsRequest{
		PageNum:     pageNum,
		PageSize:    pageSize,
		PromotionId: promotionId,
		Ptype:       ptype,
	})

	if err != nil {
		return nil, err
	}

	return promotion, err
}

func (mr *materialRepo) List(ctx context.Context, pageNum, pageSize, companyId, productId uint64, isShowCollect uint32, phone, category, keyword, search, msort, mplatform string) (*v1.ListMaterialsReply, error) {
	list, err := mr.data.materialuc.ListMaterials(ctx, &v1.ListMaterialsRequest{
		PageNum:       pageNum,
		PageSize:      pageSize,
		CompanyId:     companyId,
		ProductId:     productId,
		Keyword:       keyword,
		Search:        search,
		Category:      category,
		Msort:         msort,
		Mplatform:     mplatform,
		IsShowCollect: isShowCollect,
		Phone:         phone,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (mr *materialRepo) ListProduct(ctx context.Context, pageNum, pageSize uint64, category, keyword, search, msort, mplatform string) (*v1.ListProductsReply, error) {
	list, err := mr.data.materialuc.ListProducts(ctx, &v1.ListProductsRequest{
		PageNum:   pageNum,
		PageSize:  pageSize,
		Keyword:   keyword,
		Search:    search,
		Category:  category,
		Msort:     msort,
		Mplatform: mplatform,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (mr *materialRepo) ListCollects(ctx context.Context, pageNum, pageSize, companyId uint64, phone, category, keyword, search, msort, mplatform string) (*v1.ListCollectMaterialsReply, error) {
	list, err := mr.data.materialuc.ListCollectMaterials(ctx, &v1.ListCollectMaterialsRequest{
		PageNum:   pageNum,
		PageSize:  pageSize,
		CompanyId: companyId,
		Keyword:   keyword,
		Search:    search,
		Category:  category,
		Msort:     msort,
		Mplatform: mplatform,
		Phone:     phone,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (mr *materialRepo) Select(ctx context.Context) (*v1.ListSelectMaterialsReply, error) {
	list, err := mr.data.materialuc.ListSelectMaterials(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (mr *materialRepo) Statistics(ctx context.Context) (*v1.StatisticsMaterialsReply, error) {
	list, err := mr.data.materialuc.StatisticsMaterials(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (mr *materialRepo) UpdateCollect(ctx context.Context, companyId, videoId uint64, phone string) (*v1.UpdateCollectsReply, error) {
	collect, err := mr.data.materialuc.UpdateCollects(ctx, &v1.UpdateCollectsRequest{
		CompanyId: companyId,
		VideoId:   videoId,
		Phone:     phone,
	})

	if err != nil {
		return nil, err
	}

	return collect, err
}

func (mr *materialRepo) Down(ctx context.Context, companyId, videoId, companyMaterialId uint64, downType string) (*v1.DownMaterialsReply, error) {
	down, err := mr.data.materialuc.DownMaterials(ctx, &v1.DownMaterialsRequest{
		CompanyId:         companyId,
		VideoId:           videoId,
		CompanyMaterialId: companyMaterialId,
		DownType:          downType,
	})

	if err != nil {
		return nil, err
	}

	return down, err
}

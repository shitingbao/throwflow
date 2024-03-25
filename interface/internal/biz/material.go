package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/material/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"
)

type MaterialRepo interface {
	Get(context.Context, uint64) (*v1.GetMaterialsReply, error)
	GetDownUrl(context.Context, uint64) (*v1.GetDownUrlVideoUrlsReply, error)
	GetVideoUrl(context.Context, uint64) (*v1.GetVideoUrlsReply, error)
	GetPromotion(context.Context, uint64, uint64, uint64, string) (*v1.GetPromotionsReply, error)
	List(context.Context, uint64, uint64, uint64, uint64, uint32, string, string, string, string, string, string) (*v1.ListMaterialsReply, error)
	ListProduct(context.Context, uint64, uint64, string, string, string, string, string) (*v1.ListProductsReply, error)
	ListCollects(context.Context, uint64, uint64, uint64, string, string, string, string, string, string) (*v1.ListCollectMaterialsReply, error)
	Select(context.Context) (*v1.ListSelectMaterialsReply, error)
	Statistics(context.Context) (*v1.StatisticsMaterialsReply, error)
	UpdateCollect(context.Context, uint64, uint64, string) (*v1.UpdateCollectsReply, error)
	Down(context.Context, uint64, uint64, uint64, string) (*v1.DownMaterialsReply, error)
}

type MaterialUsecase struct {
	repo MaterialRepo
	conf *conf.Data
	log  *log.Helper
}

func NewMaterialUsecase(repo MaterialRepo, conf *conf.Data, logger log.Logger) *MaterialUsecase {
	return &MaterialUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (muc *MaterialUsecase) ListMaterials(ctx context.Context, pageNum, pageSize, companyId uint64, phone, category, keyword, search, msort, mplatform string) (*v1.ListMaterialsReply, error) {
	if pageSize == 0 {
		pageSize = uint64(muc.conf.Database.PageSize)
	}

	list, err := muc.repo.List(ctx, pageNum, pageSize, companyId, 0, 1, phone, category, keyword, search, msort, mplatform)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_MATERIAL_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (muc *MaterialUsecase) ListMiniMaterials(ctx context.Context, pageNum, pageSize uint64) (*v1.ListMaterialsReply, error) {
	if pageSize == 0 {
		pageSize = uint64(muc.conf.Database.PageSize)
	}

	list, err := muc.repo.List(ctx, pageNum, pageSize, 0, 0, 0, "", "", "", "", "", "")

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_MATERIAL_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (muc *MaterialUsecase) ListMiniMaterialProducts(ctx context.Context, pageNum, pageSize, productId uint64) (*v1.ListMaterialsReply, error) {
	if pageSize == 0 {
		pageSize = uint64(muc.conf.Database.PageSize)
	}

	list, err := muc.repo.List(ctx, pageNum, pageSize, 0, productId, 0, "", "", "", "", "", "")

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_MATERIAL_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (muc *MaterialUsecase) ListMaterialProducts(ctx context.Context, pageNum, pageSize uint64, category, keyword, search, msort, mplatform string) (*v1.ListProductsReply, error) {
	if pageSize == 0 {
		pageSize = uint64(muc.conf.Database.PageSize)
	}

	list, err := muc.repo.ListProduct(ctx, pageNum, pageSize, category, keyword, search, msort, mplatform)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_PRODUCT_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (muc *MaterialUsecase) ListCollectMaterials(ctx context.Context, pageNum, pageSize, companyId uint64, phone, category, keyword, search, msort, mplatform string) (*v1.ListCollectMaterialsReply, error) {
	if pageSize == 0 {
		pageSize = uint64(muc.conf.Database.PageSize)
	}

	list, err := muc.repo.ListCollects(ctx, pageNum, pageSize, companyId, phone, category, keyword, search, msort, mplatform)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_COLLECT_MATERIAL_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (muc *MaterialUsecase) ListSelectMaterials(ctx context.Context) (*v1.ListSelectMaterialsReply, error) {
	list, err := muc.repo.Select(ctx)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_SELECT_MATERIAL_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (muc *MaterialUsecase) StatisticsMaterials(ctx context.Context) (*v1.StatisticsMaterialsReply, error) {
	list, err := muc.repo.Statistics(ctx)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_STATISTICS_MATERIAL_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (muc *MaterialUsecase) GetDownVideoUrls(ctx context.Context, videoId uint64) (*v1.GetDownUrlVideoUrlsReply, error) {
	videoUrl, err := muc.repo.GetDownUrl(ctx, videoId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_DOWN_VIDEO_URL_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return videoUrl, nil
}

func (muc *MaterialUsecase) GetMiniDownVideoUrls(ctx context.Context, videoId uint64) (*v1.GetDownUrlVideoUrlsReply, error) {
	videoUrl, err := muc.repo.GetDownUrl(ctx, videoId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_DOWN_VIDEO_URL_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return videoUrl, nil
}

func (muc *MaterialUsecase) GetVideoUrls(ctx context.Context, videoId uint64) (*v1.GetVideoUrlsReply, error) {
	videoUrl, err := muc.repo.GetVideoUrl(ctx, videoId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_VIDEO_URL_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return videoUrl, nil
}

func (muc *MaterialUsecase) GetPromotions(ctx context.Context, pageNum, pageSize, promotionId uint64, ptype string) (*v1.GetPromotionsReply, error) {
	if pageSize == 0 {
		pageSize = uint64(muc.conf.Database.PageSize)
	}

	promotion, err := muc.repo.GetPromotion(ctx, pageNum, pageSize, promotionId, ptype)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_PROMOTION_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return promotion, nil
}

func (muc *MaterialUsecase) GetMiniMaterials(ctx context.Context, videoId uint64) (*v1.GetMaterialsReply, error) {
	material, err := muc.repo.Get(ctx, videoId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_MATERIAL_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return material, nil
}

func (muc *MaterialUsecase) UpdateCollects(ctx context.Context, companyId, videoId uint64, phone string) error {
	if _, err := muc.repo.UpdateCollect(ctx, companyId, videoId, phone); err != nil {
		return errors.InternalServer("INTERFACE_UPDATE_COLLECT_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return nil
}

func (muc *MaterialUsecase) DownMaterials(ctx context.Context, companyId, videoId, companyMaterialId uint64, downType string) error {
	if _, err := muc.repo.Down(ctx, companyId, videoId, companyMaterialId, downType); err != nil {
		return errors.InternalServer("INTERFACE_DOWN_MATERIAL_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return nil
}

func (muc *MaterialUsecase) DownMiniMaterials(ctx context.Context, videoId uint64, downType string) error {
	if _, err := muc.repo.Down(ctx, 0, videoId, 0, downType); err != nil {
		return errors.InternalServer("INTERFACE_DOWN_MATERIAL_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return nil
}

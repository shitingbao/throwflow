package biz

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"material/internal/domain"
)

var (
	MaterialCollectListError   = errors.InternalServer("MATERIAL_COLLECT_LIST_ERROR", "素材参谋收藏获取失败")
	MaterialCollectCreateError = errors.InternalServer("MATERIAL_COLLECT_CREATE_ERROR", "素材参谋收藏创建失败")
	MaterialCollectUpdateError = errors.InternalServer("MATERIAL_COLLECT_UPDATE_ERROR", "素材参谋收藏更新失败")
	MaterialCollectDeleteError = errors.InternalServer("MATERIAL_COLLECT_DELETE_ERROR", "素材参谋收藏删除失败")
)

type CollectRepo interface {
	Get(context.Context, uint64, uint64, string) (*domain.Collect, error)
	ListByVideoIds(context.Context, uint64, string, []string) ([]*domain.Collect, error)
	List(context.Context, int, int, uint64, string, string, string, string, string, []domain.CompanyProductCategory) ([]*domain.Material, error)
	Count(context.Context, uint64, string, string, string, string, []domain.CompanyProductCategory) (int64, error)
	Save(context.Context, *domain.Collect) (*domain.Collect, error)
	Delete(context.Context, *domain.Collect) error
}

type CollectUsecase struct {
	repo  CollectRepo
	mrepo MaterialRepo
	log   *log.Helper
}

func NewCollectUsecase(repo CollectRepo, mrepo MaterialRepo, logger log.Logger) *CollectUsecase {
	return &CollectUsecase{repo: repo, mrepo: mrepo, log: log.NewHelper(logger)}
}

func (cuc *CollectUsecase) ListCollectMaterials(ctx context.Context, pageNum, pageSize, companyId uint64, phone, category, keyword, search, msort, mplatform string) (*domain.MaterialList, error) {
	sortBy := "update_day"

	if msort == "like" {
		sortBy = "video_like"
	} else if msort == "isHot" {
		sortBy = "is_hot"
	}

	var categories []domain.CompanyProductCategory

	if len(category) > 0 {
		json.Unmarshal([]byte(category), &categories)
	}

	materials, err := cuc.repo.List(ctx, int(pageNum), int(pageSize), companyId, phone, keyword, search, sortBy, mplatform, categories)

	if err != nil {
		return nil, MaterialCollectListError
	}

	total, err := cuc.repo.Count(ctx, companyId, phone, keyword, search, mplatform, categories)

	if err != nil {
		return nil, MaterialCollectListError
	}

	list := make([]*domain.Material, 0)

	for _, material := range materials {
		material.IsCollect = 1
		material.SetPlatformName(ctx)
		material.SetAwemeFollowersShow(ctx)
		material.SetVideoLikeShowA(ctx)
		material.SetVideoLikeShowB(ctx)

		list = append(list, material)
	}

	return &domain.MaterialList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (cuc *CollectUsecase) UpdateCollects(ctx context.Context, companyId, videoId uint64, phone string) error {
	if _, err := cuc.mrepo.GetByVideoId(ctx, videoId); err != nil {
		return MaterialMaterialNotFound
	}

	if collect, err := cuc.repo.Get(ctx, companyId, videoId, phone); err != nil {
		inCollect := domain.NewCollect(ctx, companyId, videoId, phone)
		inCollect.SetCreateTime(ctx)
		inCollect.SetUpdateTime(ctx)

		if _, err := cuc.repo.Save(ctx, inCollect); err != nil {
			return MaterialCollectCreateError
		} else {
			return nil
		}
	} else {
		if err := cuc.repo.Delete(ctx, collect); err != nil {
			return MaterialCollectDeleteError
		} else {
			return nil
		}
	}
}

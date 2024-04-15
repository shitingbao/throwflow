package biz

import (
	"context"
	"encoding/json"
	"io"
	"material/internal/conf"
	"material/internal/domain"
	"net/http"
	"strconv"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

var (
	MaterialContentListError   = errors.InternalServer("MATERIAL_CONTENT_LIST_ERROR", "AI文案列表出错")
	MaterialContentGetError    = errors.InternalServer("MATERIAL_CONTENT_GET_ERROR", "AI文案获取出错")
	MaterialContentCreateError = errors.InternalServer("MATERIAL_CONTENT_CREATE_ERROR", "AI文案创建出错")
	MaterialContentUpdateError = errors.InternalServer("MATERIAL_CONTENT_UPDATE_ERROR", "AI文案更新出错")
)

type MaterialContentRepo interface {
	Get(context.Context, uint64, uint64) (*domain.MaterialContent, error)
	List(context.Context, int, int, uint64) ([]*domain.MaterialContent, error)
	Count(context.Context, uint64) (int64, error)
	Save(context.Context, *domain.MaterialContent) (*domain.MaterialContent, error)
	Update(context.Context, *domain.MaterialContent) (*domain.MaterialContent, error)
}

type MaterialContentUsecase struct {
	repo  MaterialContentRepo
	mrepo MaterialRepo
	conf  *conf.Data
	log   *log.Helper
}

func NewMaterialContentUsecase(repo MaterialContentRepo, mrepo MaterialRepo, conf *conf.Data, logger log.Logger) *MaterialContentUsecase {
	return &MaterialContentUsecase{repo: repo, mrepo: mrepo, conf: conf, log: log.NewHelper(logger)}
}

func (mc *MaterialContentUsecase) ListMaterialContent(ctx context.Context, pageNum, pageSize, userId uint64) (*domain.MaterialContentList, error) {
	list, err := mc.repo.List(ctx, int(pageNum), int(pageSize), userId)

	if err != nil {
		return nil, MaterialContentListError
	}

	total, err := mc.repo.Count(ctx, userId)

	if err != nil {
		return nil, MaterialContentListError
	}

	videoIds := []uint64{}

	for _, material := range list {
		videoIds = append(videoIds, material.VideoId)
	}

	materials, err := mc.mrepo.ListByVideoIds(ctx, videoIds)

	if err != nil {
		return nil, MaterialMaterialListError
	}

	materialMap := make(map[uint64]*domain.Material)

	for _, material := range materials {
		materialMap[material.VideoId] = material
	}

	for _, materialContent := range list {
		material := materialMap[materialContent.VideoId]

		if material != nil {
			materialContent.SetVideoName(ctx, material.VideoName)
			materialContent.SetVideoUrl(ctx, material.VideoUrl)
			materialContent.SetVideoCover(ctx, material.VideoCover)
		}
	}

	return &domain.MaterialContentList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, err
}

func (mc *MaterialContentUsecase) CreateMaterialContent(ctx context.Context, userId, videoId uint64) (*domain.MaterialContent, error) {
	content, err := mc.repo.Get(ctx, userId, videoId)

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, MaterialContentGetError
	}

	if content != nil {
		return content, nil
	}

	res, err := mc.getMaterialContent(ctx, videoId)

	if err != nil || res.Data.OcrContent == "" {
		return nil, MaterialContentGetError
	}

	material, err := mc.mrepo.GetByVideoId(ctx, videoId)

	if err != nil {
		return nil, MaterialMaterialListError
	}

	materialContent := domain.NewMaterialContent(ctx, material.ProductId, userId, videoId, res.Data.OcrContent)
	materialContent.SetCreateTime(ctx)
	materialContent.SetUpdateTime(ctx)

	materialContent, err = mc.repo.Save(ctx, materialContent)

	if err != nil {
		return nil, MaterialContentCreateError
	}

	return materialContent, nil
}

func (mc *MaterialContentUsecase) UpdateMaterialContent(ctx context.Context, userId, videoId uint64, content string) (*domain.MaterialContent, error) {
	materialContent, err := mc.repo.Get(ctx, userId, videoId)

	if err != nil {
		return nil, MaterialContentGetError
	}

	materialContent.SetContent(ctx, content)
	materialContent.SetUpdateTime(ctx)

	materialContent, err = mc.repo.Update(ctx, materialContent)

	if err != nil {
		return nil, MaterialContentUpdateError
	}

	return materialContent, nil
}

func (mc *MaterialContentUsecase) RecoveMaterialContent(ctx context.Context, userId, videoId uint64) (*domain.MaterialContent, error) {
	materialContent, err := mc.repo.Get(ctx, userId, videoId)

	if err != nil {
		return nil, MaterialContentGetError
	}

	res, err := mc.getMaterialContent(ctx, videoId)

	if err != nil || res.Data.OcrContent == "" {
		return nil, MaterialContentGetError
	}

	materialContent.SetContent(ctx, res.Data.OcrContent)
	materialContent.SetUpdateTime(ctx)

	materialContent, err = mc.repo.Update(ctx, materialContent)

	if err != nil {
		return nil, MaterialContentCreateError
	}

	return materialContent, nil
}

func (mc *MaterialContentUsecase) getMaterialContent(ctx context.Context, videoId uint64) (*domain.MaterialContentResult, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://saas.mengma.cn/v1/sucai/video/content", nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("video_id", strconv.FormatUint(videoId, 10))

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	res := &domain.MaterialContentResult{}

	if err := json.Unmarshal(b, res); err != nil {
		return nil, err
	}

	return res, nil
}

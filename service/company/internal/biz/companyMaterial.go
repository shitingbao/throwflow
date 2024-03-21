package biz

import (
	"company/internal/conf"
	"company/internal/domain"
	"company/internal/pkg/tool"
	"context"
	"encoding/base64"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	ctos "github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	CompanyMaterialUploadInitializationError = errors.InternalServer("COMPANY_MATERIAL_UPLOAD_INITIALIZATION_ERROR", "上传组件初始化失败")
	CompanyMaterialUploadCreateError         = errors.InternalServer("COMPANY_MATERIAL_UPLOAD_CREATE_ERROR", "上传任务创建失败")
	CompanyMaterialUploadError               = errors.InternalServer("COMPANY_MATERIAL_UPLOAD_ERROR", "上传失败")
	CompanyMaterialUploadAbortError          = errors.InternalServer("COMPANY_MATERIAL_UPLOAD_ABORT_ERROR", "上传任务取消失败")
	CompanyMaterialDeleteError               = errors.InternalServer("COMPANY_MATERIAL_DELETE_ERROR", "删除失败")
	CompanyMaterialListError                 = errors.InternalServer("COMPANY_MATERIAL_LIST_ERROR", "企业素材云库列表获取失败")
	CompanyMaterialCreateError               = errors.InternalServer("COMPANY_MATERIAL_CREATE_ERROR", "企业素材云库创建失败")
	CompanyMaterialUpdateError               = errors.InternalServer("COMPANY_MATERIAL_UPDATE_ERROR", "企业素材云库更新失败")

	AllowMaterialImageType = map[string]struct{}{
		"jpg":  struct{}{},
		"png":  struct{}{},
		"gif":  struct{}{},
		"svg":  struct{}{},
		"bmp":  struct{}{},
		"webp": struct{}{},
	}

	AllowMaterialVideoType = map[string]struct{}{
		"mpeg": struct{}{},
		"mp4":  struct{}{},
		"mov":  struct{}{},
		"avi":  struct{}{},
		"flv":  struct{}{},
		"mkv":  struct{}{},
		"mp3":  struct{}{},
		"wav":  struct{}{},
		"aac":  struct{}{},
		"ogg":  struct{}{},
		"flac": struct{}{},
		"m4a":  struct{}{},
		"wma":  struct{}{},
	}
)

type CompanyMaterialRepo interface {
	SaveCacheHash(context.Context, string, map[string]string, time.Duration) error
	UpdateCacheHash(context.Context, string, map[string]string) error
	GetCacheHash(context.Context, string, string) (string, error)

	DeleteCache(context.Context, string) error

	CreateMultipartUpload(context.Context, string) (*ctos.CreateMultipartUploadV2Output, error)
	UploadPart(context.Context, int, int64, string, string, io.Reader) (*ctos.UploadPartV2Output, error)
	CompleteMultipartUpload(context.Context, string, string, []ctos.UploadedPartV2) (*ctos.CompleteMultipartUploadV2Output, error)
	AbortMultipartUpload(context.Context, string, string) (*ctos.AbortMultipartUploadOutput, error)
	DeleteObjectV2(context.Context, string) (*ctos.DeleteObjectV2Output, error)
}

type CompanyMaterialUsecase struct {
	repo    CompanyMaterialRepo
	crepo   CompanyRepo
	cmlrepo CompanyMaterialLibraryRepo
	tm      Transaction
	conf    *conf.Data
	vconf   *conf.Volcengine
	econf   *conf.Event
	log     *log.Helper
}

func NewCompanyMaterialUsecase(repo CompanyMaterialRepo, crepo CompanyRepo, cmlrepo CompanyMaterialLibraryRepo, tm Transaction, conf *conf.Data, vconf *conf.Volcengine, econf *conf.Event, logger log.Logger) *CompanyMaterialUsecase {
	return &CompanyMaterialUsecase{repo: repo, crepo: crepo, cmlrepo: cmlrepo, tm: tm, conf: conf, vconf: vconf, econf: econf, log: log.NewHelper(logger)}
}

func (cmuc *CompanyMaterialUsecase) ListCompanyMaterials(ctx context.Context, companyId uint64) ([]*domain.CompanyMaterialLibrary, error) {
	list, err := cmuc.cmlrepo.ListByParentIdAndLibraryType(ctx, companyId, 0, 1)

	if err != nil {
		return nil, CompanyMaterialListError
	}

	for _, companyMaterial := range list {
		cmuc.getChildCompanyMaterial(ctx, companyMaterial)
	}

	return list, nil
}

func (cmuc *CompanyMaterialUsecase) GetCompanyMaterials(ctx context.Context, companyId, companyMaterialId uint64) (*domain.CompanyMaterialLibrary, error) {
	companyMaterialLibrary, err := cmuc.cmlrepo.GetByCompanyId(ctx, companyId, companyMaterialId)

	if err != nil {
		return nil, CompanyMaterialLibraryNotFound
	}

	return companyMaterialLibrary, nil
}

func (cmuc *CompanyMaterialUsecase) GetCompanyMaterialByLibraryNames(ctx context.Context, companyId, parentId uint64, materialLibraryName string) (*domain.CompanyMaterialLibrary, error) {
	companyMaterialLibrary, err := cmuc.cmlrepo.GetByCompanyIdAndParentIdAndLibraryName(ctx, companyId, parentId, materialLibraryName)

	if err != nil {
		return nil, CompanyMaterialLibraryNotFound
	}

	return companyMaterialLibrary, nil
}

func (cmuc *CompanyMaterialUsecase) GetCompanyMaterialByProductIds(ctx context.Context, companyId, parentId, productId uint64) (*domain.CompanyMaterialLibrary, error) {
	companyMaterialLibrary, err := cmuc.cmlrepo.GetByCompanyIdAndParentIdAndProductId(ctx, companyId, parentId, productId)

	if err != nil {
		return nil, CompanyMaterialLibraryNotFound
	}

	return companyMaterialLibrary, nil
}

func (cmuc *CompanyMaterialUsecase) GetUploadIdCompanyMaterials(ctx context.Context, companyId uint64, suffix string) (string, error) {
	objectKey := tool.GetRandCode(time.Now().String())

	createMultipartOutput, err := cmuc.repo.CreateMultipartUpload(ctx, cmuc.vconf.Tos.Material.SubFolder+"/"+objectKey+"."+suffix)

	if err != nil {
		return "", CompanyMaterialUploadCreateError
	}

	cacheData := make(map[string]string)
	cacheData["fileName"] = cmuc.vconf.Tos.Material.SubFolder + "/" + objectKey + "." + suffix
	cacheData["companyId"] = strconv.FormatUint(companyId, 10)

	if err := cmuc.repo.SaveCacheHash(ctx, "company:material:"+createMultipartOutput.UploadID, cacheData, cmuc.conf.Redis.MaterialTokenTimeout.AsDuration()); err != nil {
		return "", CompanyMaterialUploadCreateError
	}

	return createMultipartOutput.UploadID, nil
}

func (cmuc *CompanyMaterialUsecase) CreateCompanyMaterials(ctx context.Context, companyId, parentId, productId, viedoId, materialLiberaryFileSize uint64, materialLibraryType uint8, materialLibraryName, materialLibraryUrl, materialLiberaryFileType string) (*domain.CompanyMaterialLibrary, error) {
	var parentCompanyMaterialLibrary *domain.CompanyMaterialLibrary

	if parentId > 0 {
		var err error

		parentCompanyMaterialLibrary, err = cmuc.cmlrepo.GetByCompanyId(ctx, companyId, parentId)

		if err != nil {
			return nil, CompanyMaterialCreateError
		}
	} else {
		if materialLibraryType == 2 {
			return nil, CompanyMaterialCreateError
		}
	}

	var inCompanyMaterialLibrary *domain.CompanyMaterialLibrary

	if parentCompanyMaterialLibrary == nil {
		inCompanyMaterialLibrary = domain.NewCompanyMaterialLibrary(ctx, companyId, 0, materialLiberaryFileSize, materialLibraryType, materialLibraryName, materialLibraryUrl, materialLiberaryFileType)
	} else {
		inCompanyMaterialLibrary = domain.NewCompanyMaterialLibrary(ctx, companyId, parentCompanyMaterialLibrary.Id, materialLiberaryFileSize, materialLibraryType, materialLibraryName, materialLibraryUrl, materialLiberaryFileType)
	}

	inCompanyMaterialLibrary.SetProductId(ctx, productId)
	inCompanyMaterialLibrary.SetVideoId(ctx, viedoId)
	inCompanyMaterialLibrary.SetCreateTime(ctx)
	inCompanyMaterialLibrary.SetUpdateTime(ctx)

	companyMaterialLibrary, err := cmuc.cmlrepo.Save(ctx, inCompanyMaterialLibrary)

	if err != nil {
		return nil, CompanyMaterialCreateError
	}

	return companyMaterialLibrary, nil
}

func (cmuc *CompanyMaterialUsecase) UpdateFileSizeCompanyMaterials(ctx context.Context, companyId, companyMaterialId, materialLiberaryFileSize uint64) error {
	if _, err := cmuc.crepo.GetById(ctx, companyId); err != nil {
		return CompanyCompanyNotFound
	}

	inCompanyMaterialLibrary, err := cmuc.cmlrepo.GetByCompanyId(ctx, companyId, companyMaterialId)

	if err != nil {
		return CompanyMaterialLibraryNotFound

	}

	inCompanyMaterialLibrary.SetLiberaryFileSize(ctx, inCompanyMaterialLibrary.LiberaryFileSize+materialLiberaryFileSize)
	inCompanyMaterialLibrary.SetUpdateTime(ctx)

	if _, err := cmuc.cmlrepo.Update(ctx, inCompanyMaterialLibrary); err != nil {
		return CompanyMaterialUpdateError

	}

	return nil
}

func (cmuc *CompanyMaterialUsecase) UploadPartCompanyMaterials(ctx context.Context, companyId, partNumber, totalPart, contentLength uint64, uploadId, content string) error {
	if scompanyId, err := cmuc.repo.GetCacheHash(ctx, "company:material:"+uploadId, "companyId"); err != nil {
		return CompanyMaterialUploadError
	} else {
		if scompanyId != strconv.FormatUint(companyId, 10) {
			return CompanyMaterialUploadError
		}
	}

	if partNumberVal, err := cmuc.repo.GetCacheHash(ctx, "company:material:"+uploadId, "partNumber:"+strconv.FormatUint(partNumber, 10)); err == nil {
		if len(partNumberVal) > 0 {
			return nil
		}
	}

	fileName, err := cmuc.repo.GetCacheHash(ctx, "company:material:"+uploadId, "fileName")

	if err != nil {
		return CompanyMaterialUploadError
	}

	scontent := strings.Split(content, ",")

	if len(scontent) != 2 {
		return CompanyMaterialUploadError
	}

	bcontent, err := base64.StdEncoding.DecodeString(scontent[1])

	if err != nil {
		return CompanyMaterialUploadError
	}

	if uint64(len(bcontent)) != contentLength {
		return CompanyMaterialUploadError
	}

	partOutput, err := cmuc.repo.UploadPart(ctx, int(partNumber), int64(contentLength), fileName, uploadId, strings.NewReader(string(bcontent)))

	if err != nil {
		return CompanyMaterialUploadError
	}

	cacheData := make(map[string]string)
	cacheData["partNumber:"+strconv.FormatUint(partNumber, 10)] = partOutput.ETag
	cacheData["totalPart"] = strconv.FormatUint(totalPart, 10)

	cmuc.repo.UpdateCacheHash(ctx, "company:material:"+uploadId, cacheData)

	return nil
}

func (cmuc *CompanyMaterialUsecase) CompleteUploadCompanyMaterials(ctx context.Context, companyId uint64, uploadId string) (string, error) {
	if scompanyId, err := cmuc.repo.GetCacheHash(ctx, "company:material:"+uploadId, "companyId"); err != nil {
		return "", CompanyMaterialUploadError
	} else {
		if scompanyId != strconv.FormatUint(companyId, 10) {
			return "", CompanyMaterialUploadError
		}
	}

	fileName, err := cmuc.repo.GetCacheHash(ctx, "company:material:"+uploadId, "fileName")

	if err != nil {
		return "", CompanyMaterialUploadError
	}

	totalPart, err := cmuc.repo.GetCacheHash(ctx, "company:material:"+uploadId, "totalPart")

	if err != nil {
		return "", CompanyMaterialUploadError
	}

	itotalPart, _ := strconv.ParseUint(totalPart, 10, 64)
	var index uint64

	parts := make([]tos.UploadedPartV2, 0)

	for index = 0; index < itotalPart; index++ {
		eTag, err := cmuc.repo.GetCacheHash(ctx, "company:material:"+uploadId, "partNumber:"+strconv.FormatUint((index+1), 10))

		if err != nil {
			return "", CompanyMaterialUploadError
		}

		parts = append(parts, tos.UploadedPartV2{
			PartNumber: int(index + 1),
			ETag:       eTag,
		})
	}

	completeOutput, err := cmuc.repo.CompleteMultipartUpload(ctx, fileName, uploadId, parts)

	if err != nil {
		return "", CompanyMaterialUploadError
	}

	cmuc.repo.DeleteCache(ctx, "company:material:"+uploadId)

	return cmuc.vconf.Tos.Material.Url + "/" + completeOutput.Key, nil
}

func (cmuc *CompanyMaterialUsecase) AbortUploadCompanyMaterials(ctx context.Context, companyId uint64, uploadId string) error {
	if scompanyId, err := cmuc.repo.GetCacheHash(ctx, "company:material:"+uploadId, "companyId"); err != nil {
		return CompanyMaterialUploadAbortError
	} else {
		if scompanyId != strconv.FormatUint(companyId, 10) {
			return CompanyMaterialUploadAbortError
		}
	}

	fileName, err := cmuc.repo.GetCacheHash(ctx, "company:material:"+uploadId, "fileName")

	if err != nil {
		return CompanyMaterialUploadAbortError
	}

	if _, err := cmuc.repo.AbortMultipartUpload(ctx, fileName, uploadId); err != nil {
		return CompanyMaterialUploadAbortError
	}

	cmuc.repo.DeleteCache(ctx, "company:material:"+uploadId)

	return nil
}

func (cmuc *CompanyMaterialUsecase) DeleteCompanyMaterials(ctx context.Context, staticUrl string) error {
	if !regexp.MustCompile("^" + cmuc.vconf.Tos.Material.Url).MatchString(staticUrl) {
		return CompanyMaterialDeleteError
	}

	if _, err := cmuc.repo.DeleteObjectV2(ctx, strings.Replace(staticUrl, cmuc.vconf.Tos.Material.Url+"/", "", -1)); err != nil {
		return CompanyMaterialDeleteError
	}

	return nil
}

func (cmuc *CompanyMaterialUsecase) getChildCompanyMaterial(ctx context.Context, companyMaterial *domain.CompanyMaterialLibrary) {
	if companyMaterials, err := cmuc.cmlrepo.ListByParentIdAndLibraryType(ctx, companyMaterial.CompanyId, companyMaterial.Id, 1); err == nil {
		if len(companyMaterials) == 0 {
			return
		}

		childList := make([]*domain.CompanyMaterialLibrary, 0)

		for _, lcompanyMaterial := range companyMaterials {
			cmuc.getChildCompanyMaterial(ctx, lcompanyMaterial)

			childList = append(childList, lcompanyMaterial)
		}

		companyMaterial.ChildList = childList
	}
}

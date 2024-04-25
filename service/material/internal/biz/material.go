package biz

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"io"
	"io/ioutil"
	v1 "material/api/service/company/v1"
	"material/internal/conf"
	"material/internal/domain"
	"material/internal/pkg/event/event"
	"material/internal/pkg/event/kafka"
	"material/internal/pkg/tool"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

var (
	MaterialMaterialListError         = errors.InternalServer("MATERIAL_MATERIAL_LIST_ERROR", "素材参谋获取失败")
	MaterialProductListError          = errors.InternalServer("MATERIAL_PRODUCT_LIST_ERROR", "爆品参谋获取失败")
	MaterialMaterialNotFound          = errors.NotFound("MATERIAL_MATERIAL_NOT_FOUND", "素材参谋不存在")
	MaterialVideoUrlGetError          = errors.InternalServer("MATERIAL_VIDEO_URL_GET_ERROR", "视频播放地址获取失败")
	MaterialFileSizeGetError          = errors.InternalServer("MATERIAL_FILE_SIZE_GET_ERROR", "素材参谋文件大小获取失败")
	MaterialUploadInitializationError = errors.InternalServer("MATERIAL_UPLOAD_INITIALIZATION_ERROR", "上传组件初始化失败")
	MaterialUploadCreateError         = errors.InternalServer("MATERIAL_UPLOAD_CREATE_ERROR", "上传任务创建失败")
	MaterialUploadError               = errors.InternalServer("MATERIAL_UPLOAD_ERROR", "上传失败")
	MaterialUploadAbortError          = errors.InternalServer("MATERIAL_UPLOAD_ABORT_ERROR", "上传任务取消失败")
)

type MaterialRepo interface {
	GetByVideoId(context.Context, uint64) (*domain.Material, error)
	GetByProductId(context.Context, uint64) (*domain.Material, error)
	GetIsTopByProductId(context.Context, uint64) (*domain.Material, error)
	GetByAwemeId(context.Context, uint64) (*domain.Material, error)
	List(context.Context, int, int, uint64, string, string, string, string, []domain.CompanyProductCategory) ([]*domain.Material, error)
	ListAwemeByProductId(context.Context, uint64) ([]*domain.Material, error)
	ListByPromotionId(context.Context, int, int, uint64, string) ([]*domain.Material, error)
	Count(context.Context, uint64, string, string, string, []domain.CompanyProductCategory) (int64, error)
	CountByPromotionId(context.Context, uint64, string) (int64, error)
	Statistics(context.Context, string) (int64, error)
	StatisticsAwemeIndustry(context.Context, uint64) ([]*domain.MaterialAwemeIndustry, error)
	UpdateVideoUrl(context.Context, uint64, string) error

	SaveCacheHash(context.Context, string, map[string]string, time.Duration) error
	UpdateCacheHash(context.Context, string, map[string]string) error
	GetCacheHash(context.Context, string, string) (string, error)

	DeleteCache(context.Context, string) error

	GetObject(context.Context, string) (*tos.GetObjectV2Output, error)
	CreateMultipartUpload(context.Context, string) (*tos.CreateMultipartUploadV2Output, error)
	UploadPart(context.Context, int, int64, string, string, io.Reader) (*tos.UploadPartV2Output, error)
	CompleteMultipartUpload(context.Context, string, string, []tos.UploadedPartV2) (*tos.CompleteMultipartUploadV2Output, error)
	AbortMultipartUpload(context.Context, string, string) (*tos.AbortMultipartUploadOutput, error)

	Send(context.Context, event.Event) error
}

type MaterialUsecase struct {
	repo   MaterialRepo
	mcrepo MaterialCategoryRepo
	mprepo MaterialProductRepo
	cprepo CompanyProductRepo
	cmrepo CompanyMaterialRepo
	crepo  CollectRepo
	conf   *conf.Data
	vconf  *conf.Volcengine
	cconf  *conf.Company
	log    *log.Helper
}

func NewMaterialUsecase(repo MaterialRepo, mcrepo MaterialCategoryRepo, mprepo MaterialProductRepo, cprepo CompanyProductRepo, cmrepo CompanyMaterialRepo, crepo CollectRepo, conf *conf.Data, vconf *conf.Volcengine, cconf *conf.Company, logger log.Logger) *MaterialUsecase {
	return &MaterialUsecase{repo: repo, mcrepo: mcrepo, mprepo: mprepo, cprepo: cprepo, cmrepo: cmrepo, crepo: crepo, conf: conf, vconf: vconf, cconf: cconf, log: log.NewHelper(logger)}
}

func (muc *MaterialUsecase) ListMaterials(ctx context.Context, pageNum, pageSize, companyId, productId uint64, isShowCollect uint8, phone, category, keyword, search, msort, mplatform string) (*domain.MaterialList, error) {
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

	materials, err := muc.repo.List(ctx, int(pageNum), int(pageSize), productId, keyword, search, sortBy, mplatform, categories)

	if err != nil {
		return nil, MaterialMaterialListError
	}

	total, err := muc.repo.Count(ctx, productId, keyword, search, mplatform, categories)

	if err != nil {
		return nil, MaterialMaterialListError
	}

	var collects []*domain.Collect

	if isShowCollect == 1 {
		videoIds := make([]string, 0)

		for _, material := range materials {
			videoIds = append(videoIds, strconv.FormatUint(material.VideoId, 10))
		}

		collects, _ = muc.crepo.ListByVideoIds(ctx, companyId, phone, videoIds)
	}

	list := make([]*domain.Material, 0)

	for _, material := range materials {
		if isShowCollect == 1 {
			for _, dcollect := range collects {
				if dcollect.VideoId == material.VideoId {
					material.IsCollect = 1
					break
				}
			}
		} else {
			material.IsCollect = 0
		}

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

func (muc *MaterialUsecase) ListProducts(ctx context.Context, pageNum, pageSize uint64, category, keyword, search, msort, mplatform string) (*domain.ProductList, error) {
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

	products, err := muc.mprepo.List(ctx, int(pageNum), int(pageSize), keyword, search, sortBy, mplatform, categories)

	if err != nil {
		return nil, MaterialProductListError
	}

	total, err := muc.mprepo.Count(ctx, keyword, search, mplatform, categories)

	if err != nil {
		return nil, MaterialProductListError
	}

	list := make([]*domain.Product, 0)

	var wg sync.WaitGroup

	for _, lproduct := range products {
		wg.Add(1)

		product := &domain.Product{
			ProductId: lproduct.ProductId,
			IsHot:     lproduct.IsHot,
			VideoLike: lproduct.VideoLike,
			Awemes:    lproduct.Awemes,
			Videos:    lproduct.Videos,
			Platform:  lproduct.Platform,
		}

		product.SetVideoLikeShowA(ctx)
		product.SetVideoLikeShowB(ctx)

		go func(product *domain.Product) {
			defer wg.Done()

			if material, err := muc.repo.GetByProductId(ctx, product.ProductId); err == nil {
				material.SetPlatformName(ctx)

				product.ProductName = material.ProductName
				product.ProductImg = material.ProductImg
				product.ProductLandingPage = material.ProductLandingPage
				product.ProductPrice = material.ProductPrice
				product.Platform = material.Platform
				product.PlatformName = material.PlatformName
			}
		}(product)

		list = append(list, product)
	}

	wg.Wait()

	return &domain.ProductList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (muc *MaterialUsecase) ListAwemesByProductId(ctx context.Context, productId uint64) ([]*domain.Material, error) {
	list, err := muc.repo.ListAwemeByProductId(ctx, productId)

	if err != nil {
		return nil, MaterialMaterialListError
	}

	return list, nil
}

func (muc *MaterialUsecase) ListSelectMaterials(ctx context.Context) (*domain.SelectMaterials, error) {
	selectMaterials := domain.NewSelectMaterials()

	materialCategories, err := muc.mcrepo.List(ctx)

	if err != nil {
		return nil, MaterialMaterialCategoryNotFound
	}

	categories := make([]*domain.Category, 0)

	for _, materialCategory := range materialCategories {
		if materialCategory.ParentId == 0 {
			childList := make([]*domain.ChildListCategory, 0)

			categories = append(categories, &domain.Category{
				Key:       strconv.FormatUint(materialCategory.CategoryId, 10),
				Value:     materialCategory.CategoryName,
				ChildList: childList,
			})
		} else {
			for _, l := range categories {
				if l.Key == strconv.FormatUint(materialCategory.ParentId, 10) {
					l.ChildList = append(l.ChildList, &domain.ChildListCategory{
						Key:   strconv.FormatUint(materialCategory.CategoryId, 10),
						Value: materialCategory.CategoryName,
					})
				}
			}
		}
	}

	selectMaterials.SetCategory(ctx, categories)

	return selectMaterials, nil
}

func (muc *MaterialUsecase) GetDownUrlVideoUrls(ctx context.Context, videoId uint64) (string, error) {
	material, err := muc.repo.GetByVideoId(ctx, videoId)

	if err != nil {
		return "", MaterialMaterialNotFound
	}

	return material.VideoUrl, nil
}

func (muc *MaterialUsecase) GetVideoUrls(ctx context.Context, videoId uint64) (*domain.VideoUrlBody, error) {
	/*material, err := muc.repo.GetByVideoId(ctx, videoId)

	if err != nil {
		return nil, MaterialMaterialNotFound
	}*/

	response, err := http.Get("http://120.26.166.64:6060/v1/sucai/video/geturl?videoid=" + strconv.FormatUint(videoId, 10))
	defer response.Body.Close()

	if err != nil {
		return nil, MaterialVideoUrlGetError
	}

	rbody, _ := ioutil.ReadAll(response.Body)

	var videoUrlBody *domain.VideoUrlBody

	if err := json.Unmarshal(rbody, &videoUrlBody); err != nil {
		return nil, MaterialVideoUrlGetError
	} else {
		return videoUrlBody, nil
	}
}

func (muc *MaterialUsecase) GetPromotions(ctx context.Context, pageNum, pageSize, promotionId uint64, ptype string) (*domain.Promotion, error) {
	materials, err := muc.repo.ListByPromotionId(ctx, int(pageNum), int(pageSize), promotionId, ptype)

	if err != nil {
		return nil, MaterialMaterialListError
	}

	total, err := muc.repo.CountByPromotionId(ctx, promotionId, ptype)

	if err != nil {
		return nil, MaterialMaterialListError
	}

	promotion := &domain.Promotion{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
	}

	if ptype == "product" {
		if material, err := muc.repo.GetByProductId(ctx, promotionId); err == nil {
			material.SetPlatformName(ctx)

			promotion.PromotionId = material.ProductId
			promotion.PromotionName = material.ProductName
			promotion.PromotionType = "product"
			promotion.PromotionImg = material.ProductImg
			promotion.PromotionLandingPage = material.ProductLandingPage
			promotion.PromotionPrice = material.ProductPrice
			promotion.IndustryName = material.IndustryName
			promotion.ShopName = material.ShopName
			promotion.ShopLogo = material.ShopLogo
			promotion.PromotionPlatformName = material.PlatformName
		}
	} else if ptype == "aweme" {
		if material, err := muc.repo.GetByAwemeId(ctx, promotionId); err == nil {
			material.SetPlatformName(ctx)
			material.SetAwemeFollowersShow(ctx)

			promotion.PromotionId = material.AwemeId
			promotion.PromotionName = material.AwemeName
			promotion.PromotionType = "aweme"
			promotion.PromotionAccount = material.AwemeAccount
			promotion.PromotionImg = material.AwemeImg
			promotion.PromotionLandingPage = material.AwemeLandingPage
			promotion.PromotionFollowers = material.AwemeFollowers
			promotion.PromotionFollowersShow = material.AwemeFollowersShow
			promotion.PromotionPlatformName = material.PlatformName

			promotion.Industry = make([]*domain.MaterialAwemeIndustry, 0)

			if statisticsAwemeIndustries, err := muc.repo.StatisticsAwemeIndustry(ctx, material.AwemeId); err == nil {
				var totalItemNum uint64 = 0

				for _, statisticsAwemeIndustry := range statisticsAwemeIndustries {
					totalItemNum += statisticsAwemeIndustry.TotalItemNum
				}

				for _, statisticsAwemeIndustry := range statisticsAwemeIndustries {
					statisticsAwemeIndustry.SetIndustryRatio(ctx, totalItemNum)

					promotion.Industry = append(promotion.Industry, &domain.MaterialAwemeIndustry{
						IndustryId:    statisticsAwemeIndustry.IndustryId,
						IndustryName:  statisticsAwemeIndustry.IndustryName,
						IndustryRatio: statisticsAwemeIndustry.IndustryRatio,
						TotalItemNum:  statisticsAwemeIndustry.TotalItemNum,
					})
				}
			}
		}
	}

	for _, material := range materials {
		material.SetAwemeFollowersShow(ctx)
		material.SetVideoLikeShowA(ctx)
		material.SetVideoLikeShowB(ctx)
		material.SetPlatformName(ctx)

		promotion.List = append(promotion.List, material)
	}

	return promotion, nil
}

func (muc *MaterialUsecase) GetUploadIdMaterials(ctx context.Context, suffix string) (string, error) {
	objectKey := tool.GetRandCode(time.Now().String())

	createMultipartOutput, err := muc.repo.CreateMultipartUpload(ctx, objectKey+"."+suffix)

	if err != nil {
		return "", MaterialUploadInitializationError
	}

	cacheData := make(map[string]string)
	cacheData["fileName"] = objectKey + "." + suffix
	cacheData["companyId"] = strconv.FormatUint(muc.cconf.DefaultCompanyId, 10)

	if err := muc.repo.SaveCacheHash(ctx, "material:material:"+createMultipartOutput.UploadID, cacheData, muc.conf.Redis.MaterialTokenTimeout.AsDuration()); err != nil {
		return "", MaterialUploadCreateError
	}

	return createMultipartOutput.UploadID, nil
}

func (muc *MaterialUsecase) GetFileSizeMaterials(ctx context.Context, materialUrl string) (uint64, error) {
	object, err := muc.repo.GetObject(ctx, strings.Replace(materialUrl, muc.vconf.Tos.Material.Url+"/", "", -1))

	if err != nil {
		return 0, MaterialFileSizeGetError
	}

	return uint64(object.ContentLength), nil
}

func (muc *MaterialUsecase) GetIsTopMaterials(ctx context.Context, productId uint64) (*domain.Material, error) {
	material, err := muc.repo.GetIsTopByProductId(ctx, productId)

	if err != nil {
		return nil, MaterialMaterialNotFound
	}

	return material, nil
}

func (muc *MaterialUsecase) GetMaterials(ctx context.Context, videoId uint64) (*domain.Material, *v1.GetExternalCompanyProductsReply, error) {
	material, err := muc.repo.GetByVideoId(ctx, videoId)

	if err != nil {
		return nil, nil, MaterialMaterialNotFound
	}

	material.SetPlatformName(ctx)
	material.SetAwemeFollowersShow(ctx)
	material.SetVideoLikeShowA(ctx)
	material.SetVideoLikeShowB(ctx)

	companyProduct, _ := muc.cprepo.GetExternal(ctx, material.ProductId)

	return material, companyProduct, nil
}

func (muc *MaterialUsecase) StatisticsMaterials(ctx context.Context) (*domain.StatisticsMaterials, error) {
	statistics := make([]*domain.StatisticsMaterial, 0)

	statistics = append(statistics, &domain.StatisticsMaterial{
		Key:   "date",
		Value: time.Now().Format("2006/01/02"),
	})

	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		count, _ := muc.repo.Statistics(ctx, "")

		statistics = append(statistics, &domain.StatisticsMaterial{
			Key:   "all",
			Value: strconv.FormatInt(count, 10),
		})

		wg.Done()
	}()

	go func() {
		count, _ := muc.repo.Statistics(ctx, "product")

		statistics = append(statistics, &domain.StatisticsMaterial{
			Key:   "product",
			Value: strconv.FormatInt(count, 10),
		})

		wg.Done()
	}()

	go func() {
		count, _ := muc.repo.Statistics(ctx, "aweme")

		statistics = append(statistics, &domain.StatisticsMaterial{
			Key:   "aweme",
			Value: strconv.FormatInt(count, 10),
		})

		wg.Done()
	}()

	wg.Wait()

	return &domain.StatisticsMaterials{
		Statistics: statistics,
	}, nil
}

func (muc *MaterialUsecase) CreateMaterials(ctx context.Context, productId, videoId uint64, materialUrl, fileTrueName, materialType, fileType string) error {
	if materialType == "product" || (materialType == "material" && fileType != "image") {
		object, err := muc.repo.GetObject(ctx, strings.Replace(materialUrl, muc.vconf.Tos.Material.Url+"/", "", -1))

		if err != nil {
			return MaterialUploadError
		}

		var parentId uint64
		var cparentId uint64

		if companyMaterial, err := muc.cmrepo.GetByLibraryName(ctx, muc.cconf.DefaultCompanyId, 0, muc.cconf.Material.DefaultCompanyMaterialLibraryName); err != nil {
			if pcompanyMaterial, err := muc.cmrepo.Save(ctx, muc.cconf.DefaultCompanyId, 0, 0, 0, uint64(object.ContentLength), 1, muc.cconf.Material.DefaultCompanyMaterialLibraryName, "", ""); err != nil {
				return MaterialCompanyMaterialCreateError
			} else {
				parentId = pcompanyMaterial.Data.CompanyMaterialId
			}
		} else {
			if _, err := muc.cmrepo.UpdateFileSize(ctx, muc.cconf.DefaultCompanyId, companyMaterial.Data.CompanyMaterialId, uint64(object.ContentLength)); err != nil {
				return MaterialCompanyMaterialUpdateError
			}

			parentId = companyMaterial.Data.CompanyMaterialId
		}

		if companyMaterial, err := muc.cmrepo.GetByProductId(ctx, muc.cconf.DefaultCompanyId, parentId, productId); err != nil {
			companyProduct, err := muc.cprepo.GetByProductOutId(ctx, productId)

			if err != nil {
				return MaterialCompanyProductNotFound
			}

			if l := utf8.RuneCountInString(companyProduct.Data.ProductName); l == 0 {
				return MaterialCompanyProductNotFound
			}

			if pcompanyMaterial, err := muc.cmrepo.Save(ctx, muc.cconf.DefaultCompanyId, parentId, productId, 0, uint64(object.ContentLength), 1, companyProduct.Data.ProductName, "", ""); err != nil {
				return MaterialCompanyMaterialCreateError
			} else {
				cparentId = pcompanyMaterial.Data.CompanyMaterialId
			}
		} else {
			if _, err := muc.cmrepo.UpdateFileSize(ctx, muc.cconf.DefaultCompanyId, companyMaterial.Data.CompanyMaterialId, uint64(object.ContentLength)); err != nil {
				return MaterialCompanyMaterialUpdateError
			}

			cparentId = companyMaterial.Data.CompanyMaterialId
		}

		if _, err := muc.cmrepo.Save(ctx, muc.cconf.DefaultCompanyId, cparentId, productId, videoId, uint64(object.ContentLength), 2, fileTrueName, materialUrl, path.Ext(materialUrl)[1:]); err != nil {
			return MaterialCompanyMaterialCreateError
		}
	}

	return nil
}

func (muc *MaterialUsecase) UploadMaterials(ctx context.Context, videoId uint64, videoUrl string) error {
	material, err := muc.repo.GetByVideoId(ctx, videoId)

	if err != nil {
		return MaterialMaterialNotFound
	}

	if err := muc.repo.UpdateVideoUrl(ctx, material.VideoId, videoUrl); err != nil {
		return err
	}

	return nil
}

func (muc *MaterialUsecase) UploadPartMaterials(ctx context.Context, partNumber, totalPart, contentLength uint64, uploadId, content string) error {
	if scompanyId, err := muc.repo.GetCacheHash(ctx, "material:material:"+uploadId, "companyId"); err != nil {
		return MaterialUploadError
	} else {
		if scompanyId != strconv.FormatUint(muc.cconf.DefaultCompanyId, 10) {
			return MaterialUploadError
		}
	}

	if partNumberVal, err := muc.repo.GetCacheHash(ctx, "material:material:"+uploadId, "partNumber:"+strconv.FormatUint(partNumber, 10)); err == nil {
		if len(partNumberVal) > 0 {
			return nil
		}
	}

	fileName, err := muc.repo.GetCacheHash(ctx, "material:material:"+uploadId, "fileName")

	if err != nil {
		return MaterialUploadError
	}

	scontent := strings.Split(content, ",")

	if len(scontent) != 2 {
		return MaterialUploadError
	}

	bcontent, err := base64.StdEncoding.DecodeString(scontent[1])

	if err != nil {
		return MaterialUploadError
	}

	if uint64(len(bcontent)) != contentLength {
		return MaterialUploadError
	}

	partOutput, err := muc.repo.UploadPart(ctx, int(partNumber), int64(contentLength), fileName, uploadId, strings.NewReader(string(bcontent)))

	if err != nil {
		return MaterialUploadError
	}

	cacheData := make(map[string]string)
	cacheData["partNumber:"+strconv.FormatUint(partNumber, 10)] = partOutput.ETag
	cacheData["totalPart"] = strconv.FormatUint(totalPart, 10)

	muc.repo.UpdateCacheHash(ctx, "material:material:"+uploadId, cacheData)

	return nil
}

func (muc *MaterialUsecase) CompleteUploadMaterials(ctx context.Context, uploadId string) (string, error) {
	if scompanyId, err := muc.repo.GetCacheHash(ctx, "material:material:"+uploadId, "companyId"); err != nil {
		return "", MaterialUploadError
	} else {
		if scompanyId != strconv.FormatUint(muc.cconf.DefaultCompanyId, 10) {
			return "", MaterialUploadError
		}
	}

	fileName, err := muc.repo.GetCacheHash(ctx, "material:material:"+uploadId, "fileName")

	if err != nil {
		return "", MaterialUploadError
	}

	totalPart, err := muc.repo.GetCacheHash(ctx, "material:material:"+uploadId, "totalPart")

	if err != nil {
		return "", MaterialUploadError
	}

	itotalPart, _ := strconv.ParseUint(totalPart, 10, 64)
	var index uint64

	parts := make([]tos.UploadedPartV2, 0)

	for index = 0; index < itotalPart; index++ {
		eTag, err := muc.repo.GetCacheHash(ctx, "material:material:"+uploadId, "partNumber:"+strconv.FormatUint((index+1), 10))

		if err != nil {
			return "", MaterialUploadError
		}

		parts = append(parts, tos.UploadedPartV2{
			PartNumber: int(index + 1),
			ETag:       eTag,
		})
	}

	completeOutput, err := muc.repo.CompleteMultipartUpload(ctx, fileName, uploadId, parts)

	if err != nil {
		return "", MaterialUploadError
	}

	muc.repo.DeleteCache(ctx, "material:material:"+uploadId)

	return muc.vconf.Tos.Material.Url + "/" + completeOutput.Key, nil
}

func (muc *MaterialUsecase) AbortUploadMaterials(ctx context.Context, uploadId string) error {
	if scompanyId, err := muc.repo.GetCacheHash(ctx, "material:material:"+uploadId, "companyId"); err != nil {
		return MaterialUploadAbortError
	} else {
		if scompanyId != strconv.FormatUint(muc.cconf.DefaultCompanyId, 10) {
			return MaterialUploadError
		}
	}

	fileName, err := muc.repo.GetCacheHash(ctx, "material:material:"+uploadId, "fileName")

	if err != nil {
		return MaterialUploadAbortError
	}

	if _, err := muc.repo.AbortMultipartUpload(ctx, fileName, uploadId); err != nil {
		return MaterialUploadAbortError
	}

	muc.repo.DeleteCache(ctx, "material:material:"+uploadId)

	return nil
}

func (muc *MaterialUsecase) DownMaterials(ctx context.Context, companyId, videoId, companyMaterialId uint64, downType string) error {
	material, err := muc.repo.GetByVideoId(ctx, videoId)

	if err != nil {
		return MaterialMaterialNotFound
	}

	if downType == "local" {
		messageAd := domain.MaterialAd{
			Type: "material_material_data_sync",
		}

		messageAd.Message.Name = "material"
		messageAd.Message.CompanyId = companyId
		messageAd.Message.VideoId = material.VideoId
		messageAd.Message.Content = downType
		messageAd.Message.SendTime = tool.TimeToString("2006-01-02 15:04:05", time.Now())

		bmessageAd, _ := json.Marshal(messageAd)

		muc.repo.Send(ctx, kafka.NewMessage(strconv.FormatUint(1, 10), bmessageAd))
	} else if downType == "cloud" {
		companyMaterial, err := muc.cmrepo.Get(ctx, companyId, companyMaterialId)

		if err != nil {
			return MaterialCompanyMaterialNotFound
		}

		if companyMaterial.Data.CompanyMaterialLibraryType != 1 {
			return MaterialCompanyMaterialNotFound
		}

		messageAd := domain.MaterialAd{
			Type: "material_material_data_sync",
		}

		messageAd.Message.Name = "material"
		messageAd.Message.CompanyId = companyId
		messageAd.Message.VideoId = material.VideoId
		messageAd.Message.CompanyMaterialId = companyMaterial.Data.CompanyMaterialId
		messageAd.Message.Content = downType
		messageAd.Message.SendTime = tool.TimeToString("2006-01-02 15:04:05", time.Now())

		bmessageAd, _ := json.Marshal(messageAd)

		muc.repo.Send(ctx, kafka.NewMessage(strconv.FormatUint(1, 10), bmessageAd))
	}

	return nil
}

package biz

import (
	v1 "company/api/service/douyin/v1"
	"company/internal/conf"
	"company/internal/domain"
	"company/internal/pkg/event/event"
	"company/internal/pkg/event/kafka"
	"company/internal/pkg/tool"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	ctos "github.com/volcengine/ve-tos-golang-sdk/v2/tos"
)

var (
	CompanyProductNotFound                  = errors.NotFound("COMPANY_PRODUCT_NOT_FOUND", "企业商品不存在")
	CompanyProductExistError                = errors.InternalServer("COMPANY_PRODUCT_EXIST_ERROR", "企业商品已经存在")
	CompanyProductCreateError               = errors.InternalServer("COMPANY_PRODUCT_CREATE_ERROR", "企业商品创建失败")
	CompanyProductCommissionOutUrlError     = errors.InternalServer("COMPANY_PRODUCT_COMMISSION_OUT_URL_ERROR", "企业商品链接错误")
	CompanyProductUpdateError               = errors.InternalServer("COMPANY_PRODUCT_UPDATE_ERROR", "企业商品更新失败")
	CompanyProductUploadInitializationError = errors.InternalServer("COMPANY_PRODUCT_UPLOAD_INITIALIZATION_ERROR", "上传组件初始化失败")
	CompanyProductUploadCreateError         = errors.InternalServer("COMPANY_PRODUCT_UPLOAD_CREATE_ERROR", "上传任务创建失败")
	CompanyProductUploadError               = errors.InternalServer("COMPANY_PRODUCT_UPLOAD_ERROR", "上传失败")
	CompanyProductUploadAbortError          = errors.InternalServer("COMPANY_PRODUCT_UPLOAD_ABORT_ERROR", "上传任务取消失败")
	CompanyProductListError                 = errors.InternalServer("COMPANY_PRODUCT_LIST_ERROR", "企业商品列表获取失败")
	CompanyProductCountError                = errors.InternalServer("COMPANY_PRODUCT_COUNT_ERROR", "企业商品列表总数获取失败")
	CompanyProductReportProductListError    = errors.InternalServer("COMPANY_PRODUCT_REPORT_PRODUCT_LIST_ERROR", "企业商品商品报表列表获取失败")
	CompanyProductMaterialOutUrlError       = errors.InternalServer("COMPANY_PRODUCT_MATERIAL_OUT_URL_ERROR", "企业商品素材链接错误")
	CompanyProductMaterialUpdateError       = errors.InternalServer("COMPANY_PRODUCT_MATERIAL_UPDATE_ERROR", "企业商品素材链接更新失败")
	CompanyProductDeleteError               = errors.InternalServer("COMPANY_PRODUCT_DELETE_ERROR", "企业商品删除失败")
	CompanyDoukeProductParseError           = errors.InternalServer("COMPANY_DOUKE_PRODUCT_PARSE_ERROR", "抖客抖口令转解析失败")
)

type CompanyProductRepo interface {
	GetById(context.Context, uint64, string, string) (*domain.CompanyProduct, error)
	GetByProductOutId(context.Context, uint64, string, string) (*domain.CompanyProduct, error)
	List(context.Context, int, int, uint64, uint64, uint64, string, string, string) ([]*domain.CompanyProduct, error)
	ListByProductOutIds(context.Context, string, []uint64) ([]*domain.CompanyProduct, error)
	ListExternal(context.Context, int, int, uint64, uint64, uint64, uint8, string, string) ([]*domain.CompanyProduct, error)
	ListByProductOutIdOrName(context.Context, int, int, string) ([]*domain.CompanyProduct, error)
	Count(context.Context, uint64, uint64, uint64, uint8, string, string, string) (int64, error)
	CountByProductOutIdOrName(context.Context, string) (int64, error)
	Statistics(context.Context, uint64, uint64, uint64, string, string, string) (int64, error)
	Save(context.Context, *domain.CompanyProduct) (*domain.CompanyProduct, error)
	Update(context.Context, *domain.CompanyProduct) (*domain.CompanyProduct, error)

	GetCacheHash(context.Context, string, string) (string, error)
	SaveCacheHash(context.Context, string, map[string]string, time.Duration) error
	SaveCacheString(context.Context, string, string, time.Duration) (bool, error)
	UpdateCacheHash(context.Context, string, map[string]string) error
	DeleteCache(context.Context, string) error

	Send(context.Context, event.Event) error

	CreateMultipartUpload(context.Context, string) (*ctos.CreateMultipartUploadV2Output, error)
	UploadPart(context.Context, int, int64, string, string, io.Reader) (*ctos.UploadPartV2Output, error)
	CompleteMultipartUpload(context.Context, string, string, []ctos.UploadedPartV2) (*ctos.CompleteMultipartUploadV2Output, error)
	AbortMultipartUpload(context.Context, string, string) (*ctos.AbortMultipartUploadOutput, error)
}

type CompanyProductUsecase struct {
	repo     CompanyProductRepo
	cpcrepo  CompanyProductCategoryRepo
	cmlrepo  CompanyMaterialLibraryRepo
	crepo    CompanyRepo
	csrepo   CompanySetRepo
	ctrepo   CompanyTaskRepo
	oduirepo OpenDouyinUserInfoRepo
	jsrepo   JinritemaiStoreRepo
	jorepo   JinritemaiOrderRepo
	wusrrepo WeixinUserScanRecordRepo
	mmrepo   MaterialMaterialRepo
	dprepo   DoukeProductRepo
	tm       Transaction
	conf     *conf.Data
	cconf    *conf.Company
	vconf    *conf.Volcengine
	econf    *conf.Event
	log      *log.Helper
}

func NewCompanyProductUsecase(repo CompanyProductRepo, cpcrepo CompanyProductCategoryRepo, cmlrepo CompanyMaterialLibraryRepo, crepo CompanyRepo, csrepo CompanySetRepo, ctrepo CompanyTaskRepo, oduirepo OpenDouyinUserInfoRepo, jsrepo JinritemaiStoreRepo, jorepo JinritemaiOrderRepo, wusrrepo WeixinUserScanRecordRepo, mmrepo MaterialMaterialRepo, dprepo DoukeProductRepo, tm Transaction, conf *conf.Data, cconf *conf.Company, vconf *conf.Volcengine, econf *conf.Event, logger log.Logger) *CompanyProductUsecase {
	return &CompanyProductUsecase{repo: repo, cpcrepo: cpcrepo, cmlrepo: cmlrepo, crepo: crepo, csrepo: csrepo, ctrepo: ctrepo, oduirepo: oduirepo, jsrepo: jsrepo, jorepo: jorepo, wusrrepo: wusrrepo, mmrepo: mmrepo, dprepo: dprepo, tm: tm, conf: conf, cconf: cconf, vconf: vconf, econf: econf, log: log.NewHelper(logger)}
}

func (cpuc *CompanyProductUsecase) GetCompanyProducts(ctx context.Context, productId uint64, productStatus string) (*domain.CompanyProduct, error) {
	companyProduct, err := cpuc.repo.GetById(ctx, productId, productStatus, "1")

	if err != nil {
		return nil, CompanyProductNotFound
	}

	companyProduct.SetProductUrl(ctx)
	companyProduct.SetProductImgs(ctx)
	companyProduct.SetCommissions(ctx)
	companyProduct.SetMaterialOutUrls(ctx)

	return companyProduct, nil
}

func (cpuc *CompanyProductUsecase) GetCompanyProductByProductOutIds(ctx context.Context, productOutId uint64, productStatus string) (*domain.CompanyProduct, error) {
	companyProduct, err := cpuc.repo.GetByProductOutId(ctx, productOutId, productStatus, "1")

	if err != nil {
		return nil, CompanyProductNotFound
	}

	companyProduct.SetProductUrl(ctx)
	companyProduct.SetProductImgs(ctx)
	companyProduct.SetCommissions(ctx)
	companyProduct.SetMaterialOutUrls(ctx)

	return companyProduct, nil
}

func (cpuc *CompanyProductUsecase) GetExternalCompanyProducts(ctx context.Context, productId uint64) (*domain.CompanyProduct, error) {
	companyProduct, err := cpuc.repo.GetByProductOutId(ctx, productId, "2", "")

	if err != nil {
		return nil, CompanyProductNotFound
	}

	awemes := make([]*domain.Aweme, 0)

	if openDouyinUserInfos, err := cpuc.oduirepo.ListByProductId(ctx, strconv.FormatUint(companyProduct.ProductOutId, 10)); err == nil {
		for _, openDouyinUserInfo := range openDouyinUserInfos.Data.List {
			awemes = append(awemes, &domain.Aweme{
				Nickname:     openDouyinUserInfo.Nickname,
				AccountId:    openDouyinUserInfo.AccountId,
				Avatar:       openDouyinUserInfo.Avatar,
				AvatarLarger: openDouyinUserInfo.AvatarLarger,
			})
		}
	}

	uday, _ := strconv.ParseUint(time.Now().Format("20060102"), 10, 64)

	sampleThreshold, err := cpuc.csrepo.GetByCompanyIdAndDayAndSetKey(ctx, cpuc.cconf.DefaultCompanyId, uint32(uday), "sampleThreshold")

	if err == nil {
		sampleThreshold.GetSetValue(ctx)
	}

	if companyProduct.IsExist == 1 {
		companyProduct.SetCommissions(ctx)

		if companyProduct.SampleThresholdType == 0 && sampleThreshold != nil {
			companyProduct.SetSampleThresholdType(ctx, sampleThreshold.SetValueSampleThreshold.Type)
			companyProduct.SetSampleThresholdValue(ctx, sampleThreshold.SetValueSampleThreshold.Value)
		}

		pureCommission, pureServiceCommission, _ := companyProduct.GetCommission(ctx)

		companyProduct.SetPureCommission(ctx, pureCommission)
		companyProduct.SetPureServiceCommission(ctx, pureServiceCommission)
	} else {
		var sampleThresholdType uint8 = 0
		var sampleThresholdValue uint64 = 0

		if sampleThreshold != nil {
			sampleThresholdType = sampleThreshold.SetValueSampleThreshold.Type
			sampleThresholdValue = sampleThreshold.SetValueSampleThreshold.Value
		}

		companyProduct.SetSampleThresholdType(ctx, sampleThresholdType)
		companyProduct.SetSampleThresholdValue(ctx, sampleThresholdValue)
	}

	companyProduct.SetProductUrl(ctx)
	companyProduct.SetProductImgs(ctx)
	companyProduct.SetProductDetailImgs(ctx)
	companyProduct.SetMaterialOutUrls(ctx)
	companyProduct.SetAwemes(ctx, awemes)

	return companyProduct, nil
}

func (cpuc *CompanyProductUsecase) GetExternalProductShareCompanyProducts(ctx context.Context, productId, userId uint64) (*v1.CreateShareDoukeProductsReply, error) {
	productUrl := ""

	companyProduct, err := cpuc.repo.GetByProductOutId(ctx, productId, "2", "")

	if err != nil {
		return nil, CompanyProductNotFound
	}

	if companyProduct.IsExist == 1 {
		companyProduct.SetCommissions(ctx)

		_, _, productUrl = companyProduct.GetCommission(ctx)

		if len(productUrl) == 0 {
			productUrl = fmt.Sprintf("http://haohuo.jinritemai.com/views/product/item2?id=%d", companyProduct.ProductOutId)
		}
	} else {
		productUrl = fmt.Sprintf("http://haohuo.jinritemai.com/views/product/item2?id=%d", companyProduct.ProductOutId)
	}

	productShare, err := cpuc.dprepo.Save(ctx, productUrl, strconv.FormatUint(userId, 10))

	if err != nil {
		return nil, CompanyDoukeProductShareCreateError
	}

	return productShare, nil
}

func (cpuc *CompanyProductUsecase) GetUploadIdCompanyProducts(ctx context.Context, suffix string) (string, error) {
	objectKey := tool.GetRandCode(time.Now().String())

	createMultipartOutput, err := cpuc.repo.CreateMultipartUpload(ctx, cpuc.vconf.Tos.Product.SubFolder+"/"+objectKey+"."+suffix)

	if err != nil {
		return "", CompanyProductUploadCreateError
	}

	cacheData := make(map[string]string)
	cacheData["fileName"] = cpuc.vconf.Tos.Product.SubFolder + "/" + objectKey + "." + suffix

	if err := cpuc.repo.SaveCacheHash(ctx, "company:product:"+createMultipartOutput.UploadID, cacheData, cpuc.conf.Redis.ProductTokenTimeout.AsDuration()); err != nil {
		return "", CompanyProductUploadCreateError
	}

	return createMultipartOutput.UploadID, nil
}

func (cpuc *CompanyProductUsecase) ListCompanyProducts(ctx context.Context, pageNum, pageSize, industryId, categoryId, subCategoryId uint64, productStatus, keyword string) (*domain.CompanyProductList, error) {
	list := make([]*domain.CompanyProduct, 0)

	companyProducts, err := cpuc.repo.List(ctx, int(pageNum), int(pageSize), industryId, categoryId, subCategoryId, productStatus, "1", keyword)

	if err != nil {
		return nil, CompanyProductListError
	}

	productOutIds := make([]string, 0)

	for _, companyProduct := range companyProducts {
		companyProduct.SetProductUrl(ctx)
		companyProduct.SetProductImgs(ctx)
		companyProduct.SetCommissions(ctx)
		companyProduct.SetMaterialOutUrls(ctx)

		productOutIds = append(productOutIds, strconv.FormatUint(companyProduct.ProductOutId, 10))

		list = append(list, companyProduct)
	}

	companyTasks, _ := cpuc.ctrepo.ListByProductOutId(ctx, productOutIds)

	for _, l := range list {
		for _, companyTask := range companyTasks {
			if l.ProductOutId == companyTask.ProductOutId {
				l.IsTask = 1

				break
			}
		}
	}

	total, err := cpuc.repo.Count(ctx, industryId, categoryId, subCategoryId, 0, productStatus, "1", keyword)

	if err != nil {
		return nil, CompanyProductListError
	}

	return &domain.CompanyProductList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (cpuc *CompanyProductUsecase) ListExternalCompanyProducts(ctx context.Context, pageNum, pageSize, industryId, categoryId, subCategoryId uint64, isInvestment uint8, keyword string) (*domain.CompanyProductList, error) {
	list := make([]*domain.CompanyProduct, 0)

	companyProducts, err := cpuc.repo.ListExternal(ctx, int(pageNum), int(pageSize), industryId, categoryId, subCategoryId, isInvestment, "2", keyword)

	if err != nil {
		return nil, CompanyProductListError
	}

	total, err := cpuc.repo.Count(ctx, industryId, categoryId, subCategoryId, isInvestment, "2", "", keyword)

	if err != nil {
		return nil, CompanyProductListError
	}

	productOutIds := make([]string, 0)

	for _, companyProduct := range companyProducts {
		if companyProduct.IsExist == 1 {
			companyProduct.SetCommissions(ctx)

			pureCommission, pureServiceCommission, _ := companyProduct.GetCommission(ctx)

			companyProduct.SetPureCommission(ctx, pureCommission)
			companyProduct.SetPureServiceCommission(ctx, pureServiceCommission)

			companyProduct.SetMaterialOutUrls(ctx)
		}

		companyProduct.SetProductUrl(ctx)
		companyProduct.SetProductImgs(ctx)

		productOutIds = append(productOutIds, strconv.FormatUint(companyProduct.ProductOutId, 10))

		list = append(list, companyProduct)
	}

	companyTasks, _ := cpuc.ctrepo.ListByProductOutId(ctx, productOutIds)

	for _, l := range list {
		for _, companyTask := range companyTasks {
			if l.ProductOutId == companyTask.ProductOutId {
				l.IsTask = 1

				break
			}
		}
	}

	return &domain.CompanyProductList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (cpuc *CompanyProductUsecase) ListCompanyProductCategorys(ctx context.Context) ([]*domain.CompanyProductCategory, error) {
	companyProductCategories, err := cpuc.cpcrepo.List(ctx)

	if err != nil {
		return nil, CompanyProductCategoryNotFound
	}

	list := make([]*domain.CompanyProductCategory, 0)

	for _, companyProductCategory := range companyProductCategories {
		if companyProductCategory.ParentId == 0 {
			childList := make([]*domain.CompanyProductCategory, 0)

			list = append(list, &domain.CompanyProductCategory{
				CategoryId:   companyProductCategory.CategoryId,
				ParentId:     companyProductCategory.ParentId,
				CategoryName: companyProductCategory.CategoryName,
				Sort:         companyProductCategory.Sort,
				ChildList:    childList,
			})
		}
	}

	for _, l := range list {
		for _, companyProductCategory := range companyProductCategories {
			if l.CategoryId == companyProductCategory.ParentId {
				childList := make([]*domain.CompanyProductCategory, 0)

				l.ChildList = append(l.ChildList, &domain.CompanyProductCategory{
					CategoryId:   companyProductCategory.CategoryId,
					ParentId:     companyProductCategory.ParentId,
					CategoryName: companyProductCategory.CategoryName,
					Sort:         companyProductCategory.Sort,
					ChildList:    childList,
				})
			}
		}
	}

	for _, l := range list {
		for _, ll := range l.ChildList {
			for _, companyProductCategory := range companyProductCategories {
				if ll.CategoryId == companyProductCategory.ParentId {
					childList := make([]*domain.CompanyProductCategory, 0)

					ll.ChildList = append(ll.ChildList, &domain.CompanyProductCategory{
						CategoryId:   companyProductCategory.CategoryId,
						ParentId:     companyProductCategory.ParentId,
						CategoryName: companyProductCategory.CategoryName,
						Sort:         companyProductCategory.Sort,
						ChildList:    childList,
					})
				}
			}
		}
	}

	return list, nil
}

func (cpuc *CompanyProductUsecase) ListCompanyTaskProducts(ctx context.Context, pageNum, pageSize uint64, keyword string) (*domain.CompanyProductList, error) {
	list := make([]*domain.CompanyProduct, 0)
	productMap := make(map[uint64]bool)

	companyProducts, err := cpuc.repo.ListByProductOutIdOrName(ctx, int(pageNum), int(pageSize), keyword)

	if err != nil {
		return nil, CompanyProductListError
	}

	productOutIds := []string{}

	for _, companyProduct := range companyProducts {
		companyProduct.SetProductUrl(ctx)
		companyProduct.SetProductImgs(ctx)
		companyProduct.SetCommissions(ctx)
		companyProduct.SetMaterialOutUrls(ctx)

		productOutIds = append(productOutIds, strconv.FormatUint(companyProduct.ProductOutId, 10))

		list = append(list, companyProduct)
	}

	companyTasks, err := cpuc.ctrepo.ListByProductOutId(ctx, productOutIds)

	if err != nil {
		return nil, CompanyProductListError
	}

	for _, companyTask := range companyTasks {
		productMap[companyTask.ProductOutId] = true
	}

	for _, l := range list {
		if productMap[l.ProductOutId] {
			l.SetIsTask(ctx, 1)
		}
	}

	total, err := cpuc.repo.CountByProductOutIdOrName(ctx, keyword)

	if err != nil {
		return nil, CompanyProductCountError
	}

	return &domain.CompanyProductList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (cpuc *CompanyProductUsecase) StatisticsCompanyProducts(ctx context.Context, industryId, categoryId, subCategoryId uint64, productStatus, keyword string) ([]*domain.StatisticsCompanyProduct, error) {
	statistics := make([]*domain.StatisticsCompanyProduct, 0)

	var count int64

	if productStatus == "2" {
		count, _ = cpuc.repo.Statistics(ctx, industryId, categoryId, subCategoryId, productStatus, "", keyword)
	} else {
		count, _ = cpuc.repo.Statistics(ctx, industryId, categoryId, subCategoryId, productStatus, "1", keyword)
	}

	statistics = append(statistics, &domain.StatisticsCompanyProduct{
		Key:   "所有",
		Value: strconv.FormatInt(count, 10),
	})

	return statistics, nil
}

func (cpuc *CompanyProductUsecase) CreateCompanyProducts(ctx context.Context, commission string) (*domain.CompanyProduct, error) {
	inCompanyProduct := domain.NewCompanyProduct(ctx, 0, 0, 0, 0, 1, 0, 0, 1, 0.00, "", "", "", "", "", "")
	inCompanyProduct.SetCommission(ctx, commission)
	inCompanyProduct.SetCreateTime(ctx)
	inCompanyProduct.SetUpdateTime(ctx, time.Now())

	if ok := inCompanyProduct.VerifyCommission(ctx); !ok {
		return nil, CompanyProductCommissionOutUrlError
	}

	commissionOutUrl, _ := url.Parse(inCompanyProduct.Commissions[0].CommissionOutUrl)
	productOutId := commissionOutUrl.Query().Get("id")

	iproductOutId, _ := strconv.ParseUint(productOutId, 10, 64)

	isNotExist := true

	if tmpCompanyProduct, err := cpuc.repo.GetByProductOutId(ctx, iproductOutId, "", ""); err == nil {
		if tmpCompanyProduct.IsExist == 1 {
			return nil, CompanyProductExistError
		}

		tmpCompanyProduct.SetCommission(ctx, commission)
		tmpCompanyProduct.SetIsExist(ctx, 1)
		tmpCompanyProduct.SetUpdateTime(ctx, time.Now())

		companyProduct, err := cpuc.repo.Update(ctx, tmpCompanyProduct)

		if err != nil {
			return nil, CompanyProductCreateError
		}

		return companyProduct, nil
	}

	inCompanyProduct.SetProductOutId(ctx, iproductOutId)

	companyProduct, err := cpuc.repo.Save(ctx, inCompanyProduct)

	if err != nil {
		return nil, CompanyProductCreateError
	}

	companyProduct.SetProductUrl(ctx)
	companyProduct.SetProductImgs(ctx)
	companyProduct.SetCommissions(ctx)
	companyProduct.SetMaterialOutUrls(ctx)

	if isNotExist {
		messageAd := domain.ProductAd{
			Type: "company_product_data_sync",
		}

		messageAd.Message.Name = "company"
		messageAd.Message.ProductId = companyProduct.Id
		messageAd.Message.Content = productOutId
		messageAd.Message.SendTime = tool.TimeToString("2006-01-02 15:04:05", time.Now())

		bmessageAd, _ := json.Marshal(messageAd)

		cpuc.repo.Send(ctx, kafka.NewMessage(strconv.FormatUint(1, 10), bmessageAd))
	}

	return companyProduct, nil
}

func (cpuc *CompanyProductUsecase) CreateJinritemaiStoreCompanyProducts(ctx context.Context, userId, productId uint64, openDouyinUserIds string) (*domain.CompanyProductCreateJinritemaiStore, error) {
	companyProduct, err := cpuc.repo.GetByProductOutId(ctx, productId, "", "")

	if err != nil {
		return nil, CompanyProductNotFound
	}

	if companyProduct.IsExist == 1 {
		companyProduct.SetCommissions(ctx)

		_, _, pureCommissionUrl := companyProduct.GetCommission(ctx)

		if len(pureCommissionUrl) > 0 {
			if stores, err := cpuc.jsrepo.Save(ctx, userId, cpuc.cconf.DefaultCompanyId, companyProduct.ProductOutId, openDouyinUserIds, pureCommissionUrl); err != nil {
				return nil, CompanyJinritemaiStoreCreateError
			} else {
				message := &domain.CompanyProductCreateJinritemaiStore{
					Content: "已完成橱窗设置，现在可以发视频带货了~",
				}

				if len(stores.Data.List) > 0 {
					message.Messages = make([]*domain.CompanyProductCreateJinritemaiStoreMessage, 0)

					for _, store := range stores.Data.List {
						message.Messages = append(message.Messages, &domain.CompanyProductCreateJinritemaiStoreMessage{
							ProductName: store.ProductName,
							AwemeName:   store.AwemeName,
							Content:     store.Content,
						})
					}
				}

				return message, nil
			}
		} else {
			return nil, CompanyJinritemaiStoreCreateError
		}
	} else {
		if stores, err := cpuc.jsrepo.Save(ctx, userId, cpuc.cconf.DefaultCompanyId, companyProduct.ProductOutId, openDouyinUserIds, ""); err != nil {
			return nil, CompanyJinritemaiStoreCreateError
		} else {
			message := &domain.CompanyProductCreateJinritemaiStore{
				Content: "已完成橱窗设置，现在可以发视频带货了~",
			}

			if len(stores.Data.List) > 0 {
				message.Messages = make([]*domain.CompanyProductCreateJinritemaiStoreMessage, 0)

				for _, store := range stores.Data.List {
					message.Messages = append(message.Messages, &domain.CompanyProductCreateJinritemaiStoreMessage{
						ProductName: store.ProductName,
						AwemeName:   store.AwemeName,
						Content:     store.Content,
					})
				}
			}

			return message, nil
		}
	}

	return nil, nil
}

func (cpuc *CompanyProductUsecase) UpdateStatusCompanyProducts(ctx context.Context, productId uint64, status uint8) (*domain.CompanyProduct, error) {
	inCompanyProduct, err := cpuc.repo.GetById(ctx, productId, "", "1")

	if err != nil {
		return nil, CompanyProductNotFound
	}

	inCompanyProduct.SetProductStatus(ctx, status)
	inCompanyProduct.SetUpdateTime(ctx, time.Now())

	if status == 0 {
		inCompanyProduct.SetForbidReason(ctx, "下架")
	} else if status == 2 {
		inCompanyProduct.SetForbidReason(ctx, "上架")
	}

	companyProduct, err := cpuc.repo.Update(ctx, inCompanyProduct)

	if err != nil {
		return nil, CompanyProductUpdateError
	}

	companyProduct.SetProductUrl(ctx)
	companyProduct.SetProductImgs(ctx)
	companyProduct.SetCommissions(ctx)
	companyProduct.SetMaterialOutUrls(ctx)

	return companyProduct, nil
}

func (cpuc *CompanyProductUsecase) UpdateIsTopCompanyProducts(ctx context.Context, productId uint64, isTop uint8) (*domain.CompanyProduct, error) {
	inCompanyProduct, err := cpuc.repo.GetById(ctx, productId, "", "1")

	if err != nil {
		return nil, CompanyProductNotFound
	}

	inCompanyProduct.SetIsTop(ctx, isTop)
	inCompanyProduct.SetUpdateTime(ctx, time.Now())

	companyProduct, err := cpuc.repo.Update(ctx, inCompanyProduct)

	if err != nil {
		return nil, CompanyProductUpdateError
	}

	companyProduct.SetProductUrl(ctx)
	companyProduct.SetProductImgs(ctx)
	companyProduct.SetCommissions(ctx)
	companyProduct.SetMaterialOutUrls(ctx)

	return companyProduct, nil
}

func (cpuc *CompanyProductUsecase) UpdateSampleThresholdCompanyProducts(ctx context.Context, productId, sampleThresholdValue uint64, sampleThresholdType uint8) (*domain.CompanyProduct, error) {
	inCompanyProduct, err := cpuc.repo.GetById(ctx, productId, "", "1")

	if err != nil {
		return nil, CompanyProductNotFound
	}

	if sampleThresholdValue == 3 {
		inCompanyProduct.SetSampleThresholdType(ctx, sampleThresholdType)
		inCompanyProduct.SetSampleThresholdValue(ctx, 0)
	} else {
		inCompanyProduct.SetSampleThresholdType(ctx, sampleThresholdType)
		inCompanyProduct.SetSampleThresholdValue(ctx, sampleThresholdValue)
	}

	inCompanyProduct.SetUpdateTime(ctx, time.Now())

	companyProduct, err := cpuc.repo.Update(ctx, inCompanyProduct)

	if err != nil {
		return nil, CompanyProductUpdateError
	}

	companyProduct.SetProductUrl(ctx)
	companyProduct.SetProductImgs(ctx)
	companyProduct.SetCommissions(ctx)
	companyProduct.SetMaterialOutUrls(ctx)

	return companyProduct, nil
}

func (cpuc *CompanyProductUsecase) UpdateCommissionCompanyProducts(ctx context.Context, productId uint64, commission string) (*domain.CompanyProduct, error) {
	inCompanyProduct, err := cpuc.repo.GetById(ctx, productId, "", "1")

	if err != nil {
		return nil, CompanyProductNotFound
	}

	inCompanyProduct.SetCommission(ctx, commission)

	if ok := inCompanyProduct.VerifyCommission(ctx); !ok {
		return nil, CompanyProductCommissionOutUrlError
	}

	isNotExist := true

	commissionOutUrl, _ := url.Parse(inCompanyProduct.Commissions[0].CommissionOutUrl)
	productOutId := commissionOutUrl.Query().Get("id")

	iproductOutId, _ := strconv.ParseUint(productOutId, 10, 64)

	var companyProduct *domain.CompanyProduct

	if inCompanyProduct.ProductOutId != iproductOutId {
		if inTmpCompanyProduct, err := cpuc.repo.GetByProductOutId(ctx, iproductOutId, "", ""); err == nil {
			if inTmpCompanyProduct.IsExist == 1 {
				return nil, CompanyProductExistError
			}

			err = cpuc.tm.InTx(ctx, func(ctx context.Context) error {
				inTmpCompanyProduct.SetIsExist(ctx, 1)
				inTmpCompanyProduct.SetCommission(ctx, commission)
				inTmpCompanyProduct.SetUpdateTime(ctx, time.Now())

				companyProduct, err = cpuc.repo.Update(ctx, inTmpCompanyProduct)

				if err != nil {
					return err
				}

				inCompanyProduct.SetProductStatus(ctx, 1)
				inCompanyProduct.SetIsTop(ctx, 0)
				inCompanyProduct.SetIsExist(ctx, 0)
				inCompanyProduct.SetSampleThresholdType(ctx, 0)
				inCompanyProduct.SetSampleThresholdValue(ctx, 0)
				inCompanyProduct.SetMaterialOutUrl(ctx, "")
				inCompanyProduct.SetCommission(ctx, "")
				inCompanyProduct.SetInvestmentRatio(ctx, 0.00)
				inCompanyProduct.SetForbidReason(ctx, "")
				inCompanyProduct.SetUpdateTime(ctx, time.Now())

				if _, err = cpuc.repo.Update(ctx, inCompanyProduct); err != nil {
					return err
				}

				return nil
			})

			if err != nil {
				return nil, CompanyProductUpdateError
			}

			isNotExist = false
		} else {
			err = cpuc.tm.InTx(ctx, func(ctx context.Context) error {
				inTmpCompanyProduct = domain.NewCompanyProduct(ctx, iproductOutId, 0, 0, 0, 1, 0, 0, 1, 0.00, "", "", "", "", "", "")
				inTmpCompanyProduct.SetCommission(ctx, commission)
				inTmpCompanyProduct.SetCreateTime(ctx)
				inTmpCompanyProduct.SetUpdateTime(ctx, time.Now())

				companyProduct, err = cpuc.repo.Save(ctx, inTmpCompanyProduct)

				if err != nil {
					return err
				}

				inCompanyProduct.SetProductStatus(ctx, 1)
				inCompanyProduct.SetIsTop(ctx, 0)
				inCompanyProduct.SetIsExist(ctx, 0)
				inCompanyProduct.SetSampleThresholdType(ctx, 0)
				inCompanyProduct.SetSampleThresholdValue(ctx, 0)
				inCompanyProduct.SetMaterialOutUrl(ctx, "")
				inCompanyProduct.SetCommission(ctx, "")
				inCompanyProduct.SetInvestmentRatio(ctx, 0.00)
				inCompanyProduct.SetForbidReason(ctx, "")
				inCompanyProduct.SetUpdateTime(ctx, time.Now())

				if _, err = cpuc.repo.Update(ctx, inCompanyProduct); err != nil {
					return err
				}

				return nil
			})

			if err != nil {
				return nil, CompanyProductUpdateError
			}
		}

	} else {
		inCompanyProduct.SetUpdateTime(ctx, time.Now())

		companyProduct, err = cpuc.repo.Update(ctx, inCompanyProduct)

		if err != nil {
			return nil, CompanyProductUpdateError
		}

		isNotExist = false
	}

	if isNotExist {
		messageAd := domain.ProductAd{
			Type: "company_product_data_sync",
		}

		messageAd.Message.Name = "company"
		messageAd.Message.Content = productOutId
		messageAd.Message.SendTime = tool.TimeToString("2006-01-02 15:04:05", time.Now())

		bmessageAd, _ := json.Marshal(messageAd)

		cpuc.repo.Send(ctx, kafka.NewMessage(strconv.FormatUint(1, 10), bmessageAd))
	}

	companyProduct.SetProductUrl(ctx)
	companyProduct.SetProductImgs(ctx)
	companyProduct.SetCommissions(ctx)
	companyProduct.SetMaterialOutUrls(ctx)

	return companyProduct, nil
}

func (cpuc *CompanyProductUsecase) UpdateMaterialCompanyProducts(ctx context.Context, productId uint64, productMaterial string) error {
	inCompanyProduct, err := cpuc.repo.GetById(ctx, productId, "", "1")

	if err != nil {
		return CompanyProductNotFound
	}

	if len(productMaterial) == 0 {
		inCompanyProduct.SetMaterialOutUrl(ctx, "")
		inCompanyProduct.SetUpdateTime(ctx, time.Now())

		if _, err := cpuc.repo.Update(ctx, inCompanyProduct); err != nil {
			return CompanyProductUpdateError
		}
	} else {
		materialUrls := make([]string, 0)

		materialOutUrls := inCompanyProduct.GetMaterialOutUrls(ctx)

		inCompanyProduct.SetMaterialOutUrl(ctx, productMaterial)

		if ok := inCompanyProduct.VerifyMaterialOutUrl(ctx); !ok {
			return CompanyProductMaterialOutUrlError
		}

		newMaterialOutUrls := inCompanyProduct.GetMaterialOutUrls(ctx)

		for _, newMaterialOutUrl := range newMaterialOutUrls {
			isNotExist := true

			for _, lmaterialOutUrl := range materialOutUrls {
				if lmaterialOutUrl == newMaterialOutUrl {
					isNotExist = false

					break
				}
			}

			if isNotExist {
				materialOutUrls = append(materialOutUrls, newMaterialOutUrl)
				materialUrls = append(materialUrls, newMaterialOutUrl)
			}
		}

		inCompanyProduct.SetMaterialOutUrl(ctx, strings.Join(materialOutUrls, ","))
		inCompanyProduct.SetUpdateTime(ctx, time.Now())

		if _, err := cpuc.repo.Update(ctx, inCompanyProduct); err != nil {
			return CompanyProductMaterialUpdateError
		}

		if len(materialUrls) > 0 {
			messageAd := domain.ProductAd{
				Type: "company_product_material_data_sync",
			}

			messageAd.Message.Name = "company"
			messageAd.Message.ProductId = inCompanyProduct.ProductOutId
			messageAd.Message.Content = strings.Join(materialUrls, ",")
			messageAd.Message.SendTime = tool.TimeToString("2006-01-02 15:04:05", time.Now())

			bmessageAd, _ := json.Marshal(messageAd)

			cpuc.repo.Send(ctx, kafka.NewMessage(strconv.FormatUint(1, 10), bmessageAd))
		}
	}

	return nil
}

func (cpuc *CompanyProductUsecase) UpdateInvestmentRatioCompanyProducts(ctx context.Context, productId uint64, investmentRatio float64) (*domain.CompanyProduct, error) {
	inCompanyProduct, err := cpuc.repo.GetById(ctx, productId, "", "1")

	if err != nil {
		return nil, CompanyProductNotFound
	}

	inCompanyProduct.SetInvestmentRatio(ctx, float32(tool.Decimal(investmentRatio, 2)))
	inCompanyProduct.SetUpdateTime(ctx, time.Now())

	companyProduct, err := cpuc.repo.Update(ctx, inCompanyProduct)

	if err != nil {
		return nil, CompanyProductUpdateError
	}

	companyProduct.SetProductUrl(ctx)
	companyProduct.SetProductImgs(ctx)
	companyProduct.SetCommissions(ctx)
	companyProduct.SetMaterialOutUrls(ctx)

	return companyProduct, nil
}

func (cpuc *CompanyProductUsecase) UpdateOutProductCompanyProducts(ctx context.Context, productOutId, industryId, categoryId, subCategoryId, totalSale uint64, productStatus uint8, shopScore, commissionRatio float64, productName, productImg, productDetailImg, productPrice, industryName, categoryName, subCategoryName, shopName, shopLogo, forbidReason string) error {
	for num := 0; num <= 1; num++ {
		result, err := cpuc.repo.SaveCacheString(ctx, "company:product:product:"+strconv.FormatUint(productOutId, 10), "go", cpuc.conf.Redis.ProductLockTimeout.AsDuration())

		if err != nil {
			return CompanyProductUpdateError
		}

		if result {
			inCompanyProduct, err := cpuc.repo.GetByProductOutId(ctx, productOutId, "", "")

			if err != nil {
				cpuc.repo.DeleteCache(ctx, "company:product:product:"+strconv.FormatUint(productOutId, 10))

				return CompanyProductUpdateError
			}

			updateTime := time.Now()

			if materials, err := cpuc.mmrepo.List(ctx, 1, 1, productOutId); err == nil {
				for _, material := range materials.Data.List {
					if tupdateTime, err := tool.StringToTime("2006-01-02 15:04:05", material.UpdateTimeF); err == nil {
						updateTime = tupdateTime
					}
				}
			}

			var isHot uint8 = 0

			if inCompanyProduct.IsExist == 1 {
				if jinritemaiOrder, err := cpuc.jorepo.GetIsTop(ctx, productOutId); err == nil {
					isHot = uint8(jinritemaiOrder.Data.IsTop)
				}
			} else {
				if material, err := cpuc.mmrepo.GetIsTop(ctx, productOutId); err == nil {
					isHot = uint8(material.Data.IsTop)
				}
			}

			if productStatus == 1 {
				inCompanyProduct.SetProductStatus(ctx, productStatus)
				inCompanyProduct.SetForbidReason(ctx, forbidReason)
			} else {
				if inCompanyProduct.ProductStatus == 1 {
					inCompanyProduct.SetProductStatus(ctx, 0)
					inCompanyProduct.SetForbidReason(ctx, "下架")
				}
			}

			inCompanyProduct.SetProductName(ctx, productName)
			inCompanyProduct.SetProductImg(ctx, productImg)
			inCompanyProduct.SetProductDetailImg(ctx, productDetailImg)
			inCompanyProduct.SetProductPrice(ctx, productPrice)
			inCompanyProduct.SetIndustryId(ctx, industryId)
			inCompanyProduct.SetIndustryName(ctx, industryName)
			inCompanyProduct.SetCategoryId(ctx, categoryId)
			inCompanyProduct.SetCategoryName(ctx, categoryName)
			inCompanyProduct.SetCategoryId(ctx, subCategoryId)
			inCompanyProduct.SetCategoryName(ctx, subCategoryName)
			inCompanyProduct.SetShopName(ctx, shopName)
			inCompanyProduct.SetShopScore(ctx, shopScore)
			inCompanyProduct.SetShopLogo(ctx, shopLogo)
			inCompanyProduct.SetTotalSale(ctx, totalSale)
			inCompanyProduct.SetCommissionRatio(ctx, float32(commissionRatio))
			inCompanyProduct.SetIsHot(ctx, isHot)
			inCompanyProduct.SetUpdateTime(ctx, updateTime)

			if _, err := cpuc.repo.Update(ctx, inCompanyProduct); err != nil {
				cpuc.repo.DeleteCache(ctx, "company:product:product:"+strconv.FormatUint(productOutId, 10))

				return CompanyProductUpdateError
			}

			cpuc.repo.DeleteCache(ctx, "company:product:product:"+strconv.FormatUint(productOutId, 10))

			break
		} else {
			time.Sleep(200 * time.Millisecond)

			continue
		}
	}

	return nil
}

func (cpuc *CompanyProductUsecase) ParseCompanyProducts(ctx context.Context, content string) (*domain.CompanyProduct, error) {
	product, err := cpuc.dprepo.Parse(ctx, content)

	if err != nil {
		return nil, errors.InternalServer("COMPANY_DOUKE_PRODUCT_PARSE_ERROR", tool.GetGRPCErrorInfo(err))
	}

	productDetail, err := cpuc.dprepo.Get(ctx, product.Data.ProductId)

	if err != nil {
		return nil, CompanyDoukeProductParseError
	}

	var companyProduct *domain.CompanyProduct

	if inCompanyProduct, err := cpuc.repo.GetByProductOutId(ctx, productDetail.Data.ProductOutId, "", ""); err != nil {
		inCompanyProduct = &domain.CompanyProduct{}
		inCompanyProduct.SetProductOutId(ctx, productDetail.Data.ProductOutId)
		inCompanyProduct.SetProductStatus(ctx, 2)
		inCompanyProduct.SetProductName(ctx, productDetail.Data.ProductName)
		inCompanyProduct.SetProductImg(ctx, productDetail.Data.ProductImg)
		inCompanyProduct.SetProductPrice(ctx, productDetail.Data.ProductPrice)
		inCompanyProduct.SetIndustryId(ctx, productDetail.Data.IndustryId)
		inCompanyProduct.SetCategoryId(ctx, productDetail.Data.CategoryId)
		inCompanyProduct.SetSubCategoryId(ctx, productDetail.Data.SubCategoryId)
		inCompanyProduct.SetShopName(ctx, productDetail.Data.ShopName)
		inCompanyProduct.SetShopScore(ctx, productDetail.Data.ShopScore)
		inCompanyProduct.SetTotalSale(ctx, productDetail.Data.TotalSale)
		inCompanyProduct.SetCommissionRatio(ctx, float32(productDetail.Data.CommissionRatio))
		inCompanyProduct.SetCreateTime(ctx)
		inCompanyProduct.SetUpdateTime(ctx, time.Now())

		companyProduct, err = cpuc.repo.Save(ctx, inCompanyProduct)

		if err != nil {
			return nil, CompanyProductCreateError
		}
	} else {
		inCompanyProduct.SetProductOutId(ctx, productDetail.Data.ProductOutId)
		inCompanyProduct.SetProductStatus(ctx, 2)
		inCompanyProduct.SetProductName(ctx, productDetail.Data.ProductName)
		inCompanyProduct.SetProductImg(ctx, productDetail.Data.ProductImg)
		inCompanyProduct.SetProductPrice(ctx, productDetail.Data.ProductPrice)
		inCompanyProduct.SetIndustryId(ctx, productDetail.Data.IndustryId)
		inCompanyProduct.SetCategoryId(ctx, productDetail.Data.CategoryId)
		inCompanyProduct.SetSubCategoryId(ctx, productDetail.Data.SubCategoryId)
		inCompanyProduct.SetShopName(ctx, productDetail.Data.ShopName)
		inCompanyProduct.SetShopScore(ctx, productDetail.Data.ShopScore)
		inCompanyProduct.SetTotalSale(ctx, productDetail.Data.TotalSale)
		inCompanyProduct.SetCommissionRatio(ctx, float32(productDetail.Data.CommissionRatio))

		companyProduct, err = cpuc.repo.Update(ctx, inCompanyProduct)

		if err != nil {
			return nil, CompanyProductUpdateError
		}
	}

	if companyProduct.IsExist == 1 {
		companyProduct.SetCommissions(ctx)

		pureCommission, pureServiceCommission, _ := companyProduct.GetCommission(ctx)

		companyProduct.SetPureCommission(ctx, pureCommission)
		companyProduct.SetPureServiceCommission(ctx, pureServiceCommission)

		companyProduct.SetMaterialOutUrls(ctx)
	}

	companyProduct.SetProductUrl(ctx)
	companyProduct.SetProductImgs(ctx)

	companyTasks, _ := cpuc.ctrepo.ListByProductOutId(ctx, []string{strconv.FormatUint(companyProduct.ProductOutId, 10)})

	for _, companyTask := range companyTasks {
		if companyProduct.ProductOutId == companyTask.ProductOutId {
			companyProduct.IsTask = 1

			break
		}
	}

	return companyProduct, nil
}

func (cpuc *CompanyProductUsecase) UploadPartCompanyProducts(ctx context.Context, partNumber, totalPart, contentLength uint64, uploadId, content string) error {
	if partNumberVal, err := cpuc.repo.GetCacheHash(ctx, "company:product:"+uploadId, "partNumber:"+strconv.FormatUint(partNumber, 10)); err == nil {
		if len(partNumberVal) > 0 {
			return nil
		}
	}

	fileName, err := cpuc.repo.GetCacheHash(ctx, "company:product:"+uploadId, "fileName")

	if err != nil {
		return CompanyProductUploadError
	}

	scontent := strings.Split(content, ",")

	if len(scontent) != 2 {
		return CompanyProductUploadError
	}

	bcontent, err := base64.StdEncoding.DecodeString(scontent[1])

	if err != nil {
		return CompanyProductUploadError
	}

	if uint64(len(bcontent)) != contentLength {
		return CompanyProductUploadError
	}

	partOutput, err := cpuc.repo.UploadPart(ctx, int(partNumber), int64(contentLength), fileName, uploadId, strings.NewReader(string(bcontent)))

	if err != nil {
		return CompanyProductUploadError
	}

	cacheData := make(map[string]string)
	cacheData["partNumber:"+strconv.FormatUint(partNumber, 10)] = partOutput.ETag
	cacheData["totalPart"] = strconv.FormatUint(totalPart, 10)

	cpuc.repo.UpdateCacheHash(ctx, "company:product:"+uploadId, cacheData)

	return nil
}

func (cpuc *CompanyProductUsecase) CompleteUploadCompanyProducts(ctx context.Context, uploadId string) (string, error) {
	fileName, err := cpuc.repo.GetCacheHash(ctx, "company:product:"+uploadId, "fileName")

	if err != nil {
		return "", CompanyProductUploadError
	}

	totalPart, err := cpuc.repo.GetCacheHash(ctx, "company:product:"+uploadId, "totalPart")

	if err != nil {
		return "", CompanyProductUploadError
	}

	itotalPart, _ := strconv.ParseUint(totalPart, 10, 64)
	var index uint64

	parts := make([]ctos.UploadedPartV2, 0)

	for index = 0; index < itotalPart; index++ {
		eTag, err := cpuc.repo.GetCacheHash(ctx, "company:product:"+uploadId, "partNumber:"+strconv.FormatUint((index+1), 10))

		if err != nil {
			return "", CompanyProductUploadError
		}

		parts = append(parts, ctos.UploadedPartV2{
			PartNumber: int(index + 1),
			ETag:       eTag,
		})
	}

	completeOutput, err := cpuc.repo.CompleteMultipartUpload(ctx, fileName, uploadId, parts)

	if err != nil {
		return "", CompanyProductUploadError
	}

	cpuc.repo.DeleteCache(ctx, "company:product:"+uploadId)

	return cpuc.vconf.Tos.Product.Url + "/" + completeOutput.Key, nil
}

func (cpuc *CompanyProductUsecase) AbortUploadCompanyProducts(ctx context.Context, uploadId string) error {
	fileName, err := cpuc.repo.GetCacheHash(ctx, "company:product:"+uploadId, "fileName")

	if err != nil {
		return CompanyProductUploadAbortError
	}

	if _, err := cpuc.repo.AbortMultipartUpload(ctx, fileName, uploadId); err != nil {
		return CompanyProductUploadAbortError
	}

	cpuc.repo.DeleteCache(ctx, "company:product:"+uploadId)

	return nil
}

func (cpuc *CompanyProductUsecase) DeleteCompanyProducts(ctx context.Context, productId uint64) error {
	inCompanyProduct, err := cpuc.repo.GetById(ctx, productId, "", "1")

	if err != nil {
		return CompanyProductNotFound
	}

	inCompanyProduct.SetProductStatus(ctx, 0)
	inCompanyProduct.SetIsTop(ctx, 0)
	inCompanyProduct.SetIsExist(ctx, 0)
	inCompanyProduct.SetSampleThresholdType(ctx, 0)
	inCompanyProduct.SetSampleThresholdValue(ctx, 0)
	inCompanyProduct.SetMaterialOutUrl(ctx, "")
	inCompanyProduct.SetCommission(ctx, "")
	inCompanyProduct.SetInvestmentRatio(ctx, 0.00)
	inCompanyProduct.SetForbidReason(ctx, "")
	inCompanyProduct.SetUpdateTime(ctx, time.Now())

	if _, err = cpuc.repo.Update(ctx, inCompanyProduct); err != nil {
		return CompanyProductDeleteError
	}

	return nil
}

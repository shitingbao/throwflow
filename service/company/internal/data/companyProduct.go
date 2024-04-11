package data

import (
	"company/internal/biz"
	"company/internal/domain"
	"company/internal/pkg/event/event"
	"context"
	"errors"
	"io"
	"time"
	"unicode/utf8"

	"github.com/go-kratos/kratos/v2/log"
	ctos "github.com/volcengine/ve-tos-golang-sdk/v2/tos"
)

// 企业爆品库表
type CompanyProduct struct {
	Id                   uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	ProductOutId         uint64    `gorm:"column:product_out_id;type:bigint(20) UNSIGNED;uniqueIndex:product_out_id;not null;comment:外部商品ID"`
	ProductType          uint8     `gorm:"column:product_type;type:tinyint(3) UNSIGNED;not null;default:1;comment:公司类型，1：抖音平台商品"`
	ProductStatus        uint8     `gorm:"column:product_status;type:tinyint(3) UNSIGNED;not null;default:0;comment:上架：1：否(爬虫下架)，0：否，2：是"`
	ProductName          string    `gorm:"column:product_name;type:varchar(250);not null;comment:商品名称"`
	ProductImg           string    `gorm:"column:product_img;type:longtext;not null;comment:商品图片地址"`
	ProductDetailImg     string    `gorm:"column:product_detail_img;type:longtext;not null;comment:商品详情图片地址"`
	ProductPrice         string    `gorm:"column:product_price;type:varchar(150);not null;comment:商品售价"`
	IndustryId           uint64    `gorm:"column:industry_id;type:bigint(20) UNSIGNED;not null;default:0;comment:行业ID"`
	IndustryName         string    `gorm:"column:industry_name;type:varchar(50);not null;comment:行业名称"`
	CategoryId           uint64    `gorm:"column:category_id;type:bigint(20) UNSIGNED;not null;default:0;comment:分类ID"`
	CategoryName         string    `gorm:"column:category_name;type:varchar(50);not null;comment:分类名称"`
	SubCategoryId        uint64    `gorm:"column:sub_category_id;type:bigint(20) UNSIGNED;not null;default:0;comment:子分类ID"`
	SubCategoryName      string    `gorm:"column:sub_category_name;type:varchar(50);not null;comment:子分类名称"`
	ShopName             string    `gorm:"column:shop_name;type:varchar(250);not null;default:'';comment:店铺名称"`
	ShopLogo             string    `gorm:"column:shop_logo;type:varchar(250);not null;default:'';comment:店铺LOGO"`
	ShopScore            float64   `gorm:"column:shop_score;type:decimal(10, 2) UNSIGNED;not null;default:0.00;comment:店铺评分"`
	IsTop                uint8     `gorm:"column:is_top;type:tinyint(3) UNSIGNED;not null;default:0;comment:置顶：1：是，0：否"`
	IsHot                uint8     `gorm:"column:is_hot;type:tinyint(3) UNSIGNED;not null;default:0;comment:是否爆单，0 没有爆单，1 爆单"`
	IsExist              uint8     `gorm:"column:is_exist;type:tinyint(3) UNSIGNED;not null;default:0;comment:是否存在盟码爆品库，0 不存在，1 存在"`
	TotalSale            uint64    `gorm:"column:total_sale;type:bigint(20) UNSIGNED;not null;default:0;comment:总销量"`
	CommissionRatio      float32   `gorm:"column:commission_ratio;type:decimal(10, 2) UNSIGNED;not null;default:0;comment:佣金比例，单位：%"`
	SampleThresholdType  uint8     `gorm:"column:sample_threshold_type;type:tinyint(3) UNSIGNED;not null;default:0;comment:免费申请样品门槛：1：达人最近30天销量，2：达人最近30天销额，3: 不接受寄样申请"`
	SampleThresholdValue uint64    `gorm:"column:sample_threshold_value;type:bigint(20) UNSIGNED;not null;default:0;comment:免费申请样品门槛数值"`
	MaterialOutUrl       string    `gorm:"column:material_out_url;type:text;not null;comment:外部素材URL"`
	Commission           string    `gorm:"column:commission;type:text;not null;comment:达人佣金机制"`
	InvestmentRatio      float32   `gorm:"column:investment_ratio;type:decimal(10, 2) UNSIGNED;not null;default:0;comment:投流费比，单位：%"`
	ForbidReason         string    `gorm:"column:forbid_reason;type:varchar(250);not null;comment:下架原因"`
	CreateTime           time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime           time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (CompanyProduct) TableName() string {
	return "company_product"
}

type companyProductRepo struct {
	data *Data
	log  *log.Helper
}

func (cp *CompanyProduct) ToDomain() *domain.CompanyProduct {
	companyProduct := &domain.CompanyProduct{
		Id:                   cp.Id,
		ProductOutId:         cp.ProductOutId,
		ProductType:          cp.ProductType,
		ProductStatus:        cp.ProductStatus,
		ProductName:          cp.ProductName,
		ProductImg:           cp.ProductImg,
		ProductDetailImg:     cp.ProductDetailImg,
		ProductPrice:         cp.ProductPrice,
		IndustryId:           cp.IndustryId,
		IndustryName:         cp.IndustryName,
		CategoryId:           cp.CategoryId,
		CategoryName:         cp.CategoryName,
		SubCategoryId:        cp.SubCategoryId,
		SubCategoryName:      cp.SubCategoryName,
		ShopName:             cp.ShopName,
		ShopLogo:             cp.ShopLogo,
		ShopScore:            cp.ShopScore,
		IsTop:                cp.IsTop,
		IsHot:                cp.IsHot,
		IsExist:              cp.IsExist,
		TotalSale:            cp.TotalSale,
		CommissionRatio:      cp.CommissionRatio,
		SampleThresholdType:  cp.SampleThresholdType,
		SampleThresholdValue: cp.SampleThresholdValue,
		MaterialOutUrl:       cp.MaterialOutUrl,
		Commission:           cp.Commission,
		InvestmentRatio:      cp.InvestmentRatio,
		ForbidReason:         cp.ForbidReason,
		CreateTime:           cp.CreateTime,
		UpdateTime:           cp.UpdateTime,
	}

	return companyProduct
}

func NewCompanyProductRepo(data *Data, logger log.Logger) biz.CompanyProductRepo {
	return &companyProductRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cpr *companyProductRepo) GetById(ctx context.Context, productId uint64, productStatus, isExist string) (*domain.CompanyProduct, error) {
	companyProduct := &CompanyProduct{}

	db := cpr.data.db.WithContext(ctx).
		Where("id = ?", productId)

	if len(productStatus) > 0 {
		db = db.Where("product_status = ?", productStatus)
	}

	if len(isExist) > 0 {
		db = db.Where("is_exist = ?", isExist)
	}

	if result := db.First(companyProduct); result.Error != nil {
		return nil, result.Error
	}

	return companyProduct.ToDomain(), nil
}

func (cpr *companyProductRepo) GetByProductOutId(ctx context.Context, productOutId uint64, productStatus, isExist string) (*domain.CompanyProduct, error) {
	companyProduct := &CompanyProduct{}

	db := cpr.data.db.WithContext(ctx).
		Where("product_out_id = ?", productOutId)

	if len(productStatus) > 0 {
		db = db.Where("product_status = ?", productStatus)
	}

	if len(isExist) > 0 {
		db = db.Where("is_exist = ?", isExist)
	}

	if result := db.First(companyProduct); result.Error != nil {
		return nil, result.Error
	}

	return companyProduct.ToDomain(), nil
}

func (cpr *companyProductRepo) List(ctx context.Context, pageNum, pageSize int, industryId, categoryId, subCategoryId uint64, productStatus, isExist, keyword string) ([]*domain.CompanyProduct, error) {
	var companyProducts []CompanyProduct
	list := make([]*domain.CompanyProduct, 0)

	db := cpr.data.db.WithContext(ctx)

	if len(productStatus) > 0 {
		db = db.Where("product_status = ?", productStatus)
	}

	if len(isExist) > 0 {
		db = db.Where("is_exist = ?", isExist)
	}

	if industryId > 0 {
		db = db.Where("industry_id = ?", industryId)
	}

	if categoryId > 0 {
		db = db.Where("category_id = ?", categoryId)
	}

	if subCategoryId > 0 {
		db = db.Where("sub_category_id = ?", subCategoryId)
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		db = db.Where("product_name like ?", "%"+keyword+"%")
	}

	if pageNum == 0 {
		if result := db.Order("is_top DESC,id DESC").
			Find(&companyProducts); result.Error != nil {
			return nil, result.Error
		}
	} else {
		if result := db.Order("is_top DESC,id DESC").
			Limit(pageSize).Offset((pageNum - 1) * pageSize).
			Find(&companyProducts); result.Error != nil {
			return nil, result.Error
		}
	}

	for _, companyProduct := range companyProducts {
		list = append(list, companyProduct.ToDomain())
	}

	return list, nil
}

func (cpr *companyProductRepo) ListByProductOutIds(ctx context.Context, isExist string, productOutIds []uint64) ([]*domain.CompanyProduct, error) {
	var companyProducts []CompanyProduct
	list := make([]*domain.CompanyProduct, 0)

	db := cpr.data.db.WithContext(ctx).Where("product_out_id in (?)", productOutIds)

	if len(isExist) > 0 {
		db = db.Where("is_exist = ?", isExist)
	}

	if result := db.Find(&companyProducts); result.Error != nil {
		return nil, result.Error
	}

	for _, companyProduct := range companyProducts {
		list = append(list, companyProduct.ToDomain())
	}

	return list, nil
}

func (cpr *companyProductRepo) ListExternal(ctx context.Context, pageNum, pageSize int, industryId, categoryId, subCategoryId uint64, isInvestment uint8, productStatus, keyword string) ([]*domain.CompanyProduct, error) {
	var companyProducts []CompanyProduct
	list := make([]*domain.CompanyProduct, 0)

	db := cpr.data.db.WithContext(ctx)

	if len(productStatus) > 0 {
		db = db.Where("product_status = ?", productStatus)
	}

	if industryId > 0 {
		db = db.Where("industry_id = ?", industryId)
	}

	if categoryId > 0 {
		db = db.Where("category_id = ?", categoryId)
	}

	if subCategoryId > 0 {
		db = db.Where("sub_category_id = ?", subCategoryId)
	}

	if isInvestment == 1 {
		db = db.Where("investment_ratio > 0")
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		db = db.Where("product_name like ?", "%"+keyword+"%")
	}

	if result := db.Order("is_exist DESC,is_top DESC,is_hot DESC,update_time DESC,id DESC").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&companyProducts); result.Error != nil {
		return nil, result.Error
	}

	for _, companyProduct := range companyProducts {
		list = append(list, companyProduct.ToDomain())
	}

	return list, nil
}

func (cpr *companyProductRepo) ListByProductOutIdOrName(ctx context.Context, pageNum, pageSize int, keyword string) ([]*domain.CompanyProduct, error) {
	var companyProducts []CompanyProduct
	list := make([]*domain.CompanyProduct, 0)

	db := cpr.data.db.WithContext(ctx).Where("is_exist = 1")

	if len(keyword) > 0 {
		db = db.Where("product_out_id = ? or product_name like ?", keyword, "%"+keyword+"%")
	}

	if err := db.Order("is_top DESC,id DESC").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&companyProducts).Error; err != nil {
		return nil, err
	}

	for _, companyProduct := range companyProducts {
		list = append(list, companyProduct.ToDomain())
	}

	return list, nil
}

func (cpr *companyProductRepo) Count(ctx context.Context, industryId, categoryId, subCategoryId uint64, isInvestment uint8, productStatus, isExist, keyword string) (int64, error) {
	var count int64

	db := cpr.data.db.WithContext(ctx)

	if len(productStatus) > 0 {
		db = db.Where("product_status = ?", productStatus)
	}

	if len(isExist) > 0 {
		db = db.Where("is_exist = ?", isExist)
	}

	if industryId > 0 {
		db = db.Where("industry_id = ?", industryId)
	}

	if categoryId > 0 {
		db = db.Where("category_id = ?", categoryId)
	}

	if subCategoryId > 0 {
		db = db.Where("sub_category_id = ?", subCategoryId)
	}

	if isInvestment == 1 {
		db = db.Where("investment_ratio > 0")
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		db = db.Where("product_name like ?", "%"+keyword+"%")
	}

	db = db.Select("count(id)")

	if result := db.Model(&CompanyProduct{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (cpr *companyProductRepo) CountByProductOutIdOrName(ctx context.Context, keyword string) (int64, error) {
	var count int64

	db := cpr.data.db.WithContext(ctx).Model(&CompanyProduct{}).Where("is_exist = 1")

	if len(keyword) > 0 {
		db = db.Where("product_out_id = ? or product_name like ?", keyword, "%"+keyword+"%")
	}

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (cpr *companyProductRepo) Statistics(ctx context.Context, industryId, categoryId, subCategoryId uint64, productStatus, isExist, keyword string) (int64, error) {
	var count int64

	db := cpr.data.db.WithContext(ctx)

	if len(productStatus) > 0 {
		db = db.Where("product_status = ?", productStatus)
	}

	if len(isExist) > 0 {
		db = db.Where("is_exist = ?", isExist)
	}

	if industryId > 0 {
		db = db.Where("industry_id = ?", industryId)
	}

	if categoryId > 0 {
		db = db.Where("category_id = ?", categoryId)
	}

	if subCategoryId > 0 {
		db = db.Where("sub_category_id = ?", subCategoryId)
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		db = db.Where("product_name like ?", "%"+keyword+"%")
	}

	db = db.Select("count(id)")

	if result := db.Model(&CompanyProduct{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (cpr *companyProductRepo) Save(ctx context.Context, in *domain.CompanyProduct) (*domain.CompanyProduct, error) {
	companyProduct := &CompanyProduct{
		ProductOutId:         in.ProductOutId,
		ProductType:          in.ProductType,
		ProductStatus:        in.ProductStatus,
		ProductName:          in.ProductName,
		ProductImg:           in.ProductImg,
		ProductDetailImg:     in.ProductDetailImg,
		ProductPrice:         in.ProductPrice,
		IndustryId:           in.IndustryId,
		IndustryName:         in.IndustryName,
		CategoryId:           in.CategoryId,
		CategoryName:         in.CategoryName,
		SubCategoryId:        in.SubCategoryId,
		SubCategoryName:      in.SubCategoryName,
		ShopName:             in.ShopName,
		ShopLogo:             in.ShopLogo,
		ShopScore:            in.ShopScore,
		IsTop:                in.IsTop,
		IsHot:                in.IsHot,
		IsExist:              in.IsExist,
		TotalSale:            in.TotalSale,
		CommissionRatio:      in.CommissionRatio,
		SampleThresholdType:  in.SampleThresholdType,
		SampleThresholdValue: in.SampleThresholdValue,
		MaterialOutUrl:       in.MaterialOutUrl,
		Commission:           in.Commission,
		InvestmentRatio:      in.InvestmentRatio,
		ForbidReason:         in.ForbidReason,
		CreateTime:           in.CreateTime,
		UpdateTime:           in.UpdateTime,
	}

	if result := cpr.data.DB(ctx).Create(companyProduct); result.Error != nil {
		return nil, result.Error
	}

	return companyProduct.ToDomain(), nil
}

func (cpr *companyProductRepo) Update(ctx context.Context, in *domain.CompanyProduct) (*domain.CompanyProduct, error) {
	companyProduct := &CompanyProduct{
		Id:                   in.Id,
		ProductOutId:         in.ProductOutId,
		ProductType:          in.ProductType,
		ProductStatus:        in.ProductStatus,
		ProductName:          in.ProductName,
		ProductImg:           in.ProductImg,
		ProductDetailImg:     in.ProductDetailImg,
		ProductPrice:         in.ProductPrice,
		IndustryId:           in.IndustryId,
		IndustryName:         in.IndustryName,
		CategoryId:           in.CategoryId,
		CategoryName:         in.CategoryName,
		SubCategoryId:        in.SubCategoryId,
		SubCategoryName:      in.SubCategoryName,
		ShopName:             in.ShopName,
		ShopLogo:             in.ShopLogo,
		ShopScore:            in.ShopScore,
		IsTop:                in.IsTop,
		IsHot:                in.IsHot,
		IsExist:              in.IsExist,
		TotalSale:            in.TotalSale,
		CommissionRatio:      in.CommissionRatio,
		SampleThresholdType:  in.SampleThresholdType,
		SampleThresholdValue: in.SampleThresholdValue,
		MaterialOutUrl:       in.MaterialOutUrl,
		Commission:           in.Commission,
		InvestmentRatio:      in.InvestmentRatio,
		ForbidReason:         in.ForbidReason,
		CreateTime:           in.CreateTime,
		UpdateTime:           in.UpdateTime,
	}

	if result := cpr.data.DB(ctx).Save(companyProduct); result.Error != nil {
		return nil, result.Error
	}

	return companyProduct.ToDomain(), nil
}

func (cpr *companyProductRepo) GetCacheHash(ctx context.Context, key string, field string) (string, error) {
	val, err := cpr.data.rdb.HGet(ctx, key, field).Result()

	if err != nil {
		return "", err
	}

	return val, nil
}

func (cpr *companyProductRepo) SaveCacheHash(ctx context.Context, key string, val map[string]string, timeout time.Duration) error {
	_, err := cpr.data.rdb.HMSet(ctx, key, val).Result()

	if err != nil {
		return err
	}

	_, err = cpr.data.rdb.Expire(ctx, key, timeout).Result()

	if err != nil {
		return err
	}

	return nil
}

func (cmr *companyProductRepo) SaveCacheString(ctx context.Context, key string, val string, timeout time.Duration) (bool, error) {
	result, err := cmr.data.rdb.SetNX(ctx, key, val, timeout).Result()

	if err != nil {
		return false, err
	}

	return result, nil
}

func (cpr *companyProductRepo) UpdateCacheHash(ctx context.Context, key string, val map[string]string) error {
	_, err := cpr.data.rdb.HMSet(ctx, key, val).Result()

	if err != nil {
		return err
	}

	return nil
}

func (cpr *companyProductRepo) DeleteCache(ctx context.Context, key string) error {
	if _, err := cpr.data.rdb.Del(ctx, key).Result(); err != nil {
		return err
	}

	return nil
}

func (cpr *companyProductRepo) Send(ctx context.Context, message event.Event) error {
	if err := cpr.data.kafka.Send(ctx, message); err != nil {
		return err
	}

	return nil
}

func (cpr *companyProductRepo) CreateMultipartUpload(ctx context.Context, objectKey string) (*ctos.CreateMultipartUploadV2Output, error) {
	for _, ltos := range cpr.data.toses {
		if ltos.name == "product" {
			createMultipartOutput, err := ltos.tos.CreateMultipartUpload(ctx, objectKey)

			if err != nil {
				return nil, err
			}

			return createMultipartOutput, nil
		}
	}

	return nil, errors.New("tos is not exist")
}

func (cpr *companyProductRepo) UploadPart(ctx context.Context, partNumber int, contentLength int64, objectKey, uploadId string, content io.Reader) (*ctos.UploadPartV2Output, error) {
	for _, ltos := range cpr.data.toses {
		if ltos.name == "product" {
			partOutput, err := ltos.tos.UploadPart(ctx, partNumber, contentLength, objectKey, uploadId, content)

			if err != nil {
				return nil, err
			}

			return partOutput, nil
		}
	}

	return nil, errors.New("tos is not exist")
}

func (cpr *companyProductRepo) CompleteMultipartUpload(ctx context.Context, fileName, uploadId string, parts []ctos.UploadedPartV2) (*ctos.CompleteMultipartUploadV2Output, error) {
	for _, ltos := range cpr.data.toses {
		if ltos.name == "product" {
			completeOutput, err := ltos.tos.CompleteMultipartUpload(ctx, fileName, uploadId, parts)

			if err != nil {
				return nil, err
			}

			return completeOutput, nil
		}
	}

	return nil, errors.New("tos is not exist")
}

func (cpr *companyProductRepo) AbortMultipartUpload(ctx context.Context, objectKey, uploadId string) (*ctos.AbortMultipartUploadOutput, error) {
	for _, ltos := range cpr.data.toses {
		if ltos.name == "product" {
			output, err := ltos.tos.AbortMultipartUpload(ctx, objectKey, uploadId)

			if err != nil {
				return nil, err
			}

			return output, nil
		}
	}

	return nil, errors.New("tos is not exist")
}

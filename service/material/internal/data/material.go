package data

import (
	"context"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"io"
	"material/internal/biz"
	"material/internal/domain"
	"material/internal/pkg/event/event"
	"material/internal/pkg/tool"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-kratos/kratos/v2/log"
)

// 市场素材表
type Material struct {
	Id                 uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	VideoId            uint64    `gorm:"column:video_id;type:bigint(20) UNSIGNED;not null;comment:素材ID"`
	VideoName          string    `gorm:"column:video_name;type:text;not null;comment:素材名称"`
	VideoUrl           string    `gorm:"column:video_url;type:text;not null;comment:素材视频地址"`
	VideoCover         string    `gorm:"column:video_cover;type:text;not null;comment:素材视频封面图片地址"`
	VideoLike          uint64    `gorm:"column:video_like;type:bigint(20) UNSIGNED;not null;default:0;comment:素材点赞数"`
	IndustryId         uint64    `gorm:"column:industry_id;type:bigint(20) UNSIGNED;not null;default:0;comment:行业ID"`
	IndustryName       string    `gorm:"column:industry_name;type:varchar(50);not null;comment:行业名称"`
	CategoryId         uint64    `gorm:"column:category_id;type:bigint(20) UNSIGNED;not null;default:0;comment:分类ID"`
	CategoryName       string    `gorm:"column:category_name;type:varchar(50);not null;comment:分类名称"`
	Source             string    `gorm:"column:source;type:varchar(50);not null;comment:来源"`
	AwemeId            uint64    `gorm:"column:aweme_id;type:bigint(20) UNSIGNED;not null;default:0;comment:达人ID"`
	AwemeName          string    `gorm:"column:aweme_name;type:varchar(250);not null;comment:达人名称"`
	AwemeAccount       string    `gorm:"column:aweme_account;type:varchar(50);not null;comment:达人账号"`
	AwemeFollowers     string    `gorm:"column:aweme_followers;type:varchar(50);not null;comment:粉丝数"`
	AwemeImg           string    `gorm:"column:aweme_img;type:text;not null;comment:达人图片地址"`
	AwemeLandingPage   string    `gorm:"column:aweme_landing_page;type:varchar(255);not null;comment:达人落地页"`
	ProductId          uint64    `gorm:"column:product_id;type:bigint(20) UNSIGNED;not null;default:0;comment:商品ID"`
	ProductName        string    `gorm:"column:product_name;type:text;not null;comment:商品名称"`
	ProductImg         string    `gorm:"column:product_img;type:text;not null;comment:商品图片地址"`
	ProductLandingPage string    `gorm:"column:product_landing_page;type:varchar(255);not null;comment:商品落地页"`
	ProductPrice       string    `gorm:"column:product_price;type:varchar(50);not null;comment:商品单价"`
	ShopId             string    `gorm:"column:shop_id;type:varchar(250);not null;comment:店铺ID"`
	ShopName           string    `gorm:"column:shop_name;type:varchar(250);not null;comment:店铺名称"`
	ShopLogo           string    `gorm:"column:shop_logo;type:text;not null;comment:店铺LOGO地址"`
	ShopScore          string    `gorm:"column:shop_score;type:varchar(50);not null;comment:店铺评分"`
	IsHot              uint8     `gorm:"column:is_hot;type:tinyint(3) UNSIGNED;not null;default:0;comment:是否爆单，0：不是，1：是"`
	TotalItemNum       uint64    `gorm:"column:total_item_num;type:bigint(20) UNSIGNED;not null;default:0;comment:销量"`
	Platform           string    `gorm:"column:platform;type:varchar(10);not null;comment:平台：dy：抖音，ks：快手"`
	UpdateDay          time.Time `gorm:"column:update_day;type:datetime;not null;comment:更新时间"`
	CreateTime         time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime         time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (Material) TableName() string {
	return "material_material"
}

type materialRepo struct {
	data *Data
	log  *log.Helper
}

func (m *Material) ToDomain() *domain.Material {
	return &domain.Material{
		Id:                 m.Id,
		VideoId:            m.VideoId,
		VideoName:          m.VideoName,
		VideoUrl:           m.VideoUrl,
		VideoCover:         m.VideoCover,
		VideoLike:          m.VideoLike,
		IndustryId:         m.IndustryId,
		IndustryName:       m.IndustryName,
		CategoryId:         m.CategoryId,
		CategoryName:       m.CategoryName,
		Source:             m.Source,
		AwemeId:            m.AwemeId,
		AwemeName:          m.AwemeName,
		AwemeAccount:       m.AwemeAccount,
		AwemeFollowers:     m.AwemeFollowers,
		AwemeImg:           m.AwemeImg,
		AwemeLandingPage:   m.AwemeLandingPage,
		ProductId:          m.ProductId,
		ProductName:        m.ProductName,
		ProductImg:         m.ProductImg,
		ProductLandingPage: m.ProductLandingPage,
		ProductPrice:       m.ProductPrice,
		ShopId:             m.ShopId,
		ShopName:           m.ShopName,
		ShopLogo:           m.ShopLogo,
		ShopScore:          m.ShopScore,
		IsHot:              m.IsHot,
		TotalItemNum:       m.TotalItemNum,
		Platform:           m.Platform,
		UpdateDay:          m.UpdateDay,
		CreateTime:         m.CreateTime,
		UpdateTime:         m.UpdateTime,
	}
}

func NewMaterialRepo(data *Data, logger log.Logger) biz.MaterialRepo {
	return &materialRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (mr *materialRepo) GetByVideoId(ctx context.Context, videoId uint64) (*domain.Material, error) {
	material := &Material{}

	if result := mr.data.db.WithContext(ctx).Where("video_id = ?", videoId).First(material); result.Error != nil {
		return nil, result.Error
	}

	return material.ToDomain(), nil
}

func (mr *materialRepo) GetByProductId(ctx context.Context, productId uint64) (*domain.Material, error) {
	material := &Material{}

	if result := mr.data.db.WithContext(ctx).Where("product_id = ?", productId).Order("update_day DESC").First(material); result.Error != nil {
		return nil, result.Error
	}

	return material.ToDomain(), nil
}

func (mr *materialRepo) GetIsTopByProductId(ctx context.Context, productId uint64) (*domain.Material, error) {
	material := &Material{}

	if result := mr.data.db.WithContext(ctx).
		Where("product_id = ?", productId).
		Where("total_item_num > 1000").
		Where("update_day >= ?", tool.TimeToString("2006-01-02", time.Now().AddDate(0, 0, -7))+" 00:00:00").
		Order("update_day DESC").First(material); result.Error != nil {
		return nil, result.Error
	}

	return material.ToDomain(), nil
}

func (mr *materialRepo) GetByAwemeId(ctx context.Context, awemeId uint64) (*domain.Material, error) {
	material := &Material{}

	if result := mr.data.db.WithContext(ctx).Where("aweme_id = ?", awemeId).Order("update_day DESC").First(material); result.Error != nil {
		return nil, result.Error
	}

	return material.ToDomain(), nil
}

func (mr *materialRepo) List(ctx context.Context, pageNum, pageSize int, productId uint64, keyword, search, sortBy, mplatform string, categories []domain.CompanyProductCategory) ([]*domain.Material, error) {
	var materials []Material
	list := make([]*domain.Material, 0)

	db := mr.data.db.WithContext(ctx)

	if categories != nil {
		categorySqls := make([]string, 0)

		for _, category := range categories {
			if category.IndustryId > 0 {
				categoryIds := make([]string, 0)

				for _, categoryId := range category.CategoryId {
					categoryIds = append(categoryIds, strconv.FormatUint(categoryId, 10))
				}

				if len(categoryIds) > 0 {
					categorySqls = append(categorySqls, "(industry_id = '"+strconv.FormatUint(category.IndustryId, 10)+"' and category_id in ("+strings.Join(categoryIds, ",")+"))")
				} else {
					categorySqls = append(categorySqls, "(industry_id = '"+strconv.FormatUint(category.IndustryId, 10)+"')")
				}
			}
		}

		if len(categorySqls) > 0 {
			db = db.Where(strings.Join(categorySqls, " or "))
		}
	}

	if l := utf8.RuneCountInString(mplatform); l > 0 {
		db = db.Where("platform = ?", mplatform)
	}

	if productId > 0 {
		db = db.Where("product_id = ?", productId)
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		if search == "name" {
			db = db.Where("video_name like ?", "%"+keyword+"%")
		} else if search == "product" {
			db = db.Where("product_name like ?", "%"+keyword+"%")
		} else if search == "aweme" {
			db = db.Where("(aweme_name like ? or aweme_account like ?)", "%"+keyword+"%", "%"+keyword+"%")
		}
	}

	if result := db.Order(sortBy + " DESC").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&materials); result.Error != nil {
		return nil, result.Error
	}

	for _, material := range materials {
		list = append(list, material.ToDomain())
	}

	return list, nil
}

func (mr *materialRepo) ListAwemeByProductId(ctx context.Context, productId uint64) ([]*domain.Material, error) {
	var materials []Material
	list := make([]*domain.Material, 0)

	if result := mr.data.db.WithContext(ctx).
		Where("product_id = ?", productId).
		Select("aweme_id,max(aweme_name) as aweme_name,max(aweme_account) as aweme_account,max(aweme_followers) as aweme_followers,max(aweme_img) as aweme_img,max(aweme_landing_page) as aweme_landing_page,max(update_day) as update_day").
		Group("aweme_id").
		Order("update_day DESC").
		Find(&materials); result.Error != nil {
		return nil, result.Error
	}

	for _, material := range materials {
		list = append(list, material.ToDomain())
	}

	return list, nil
}

func (mr *materialRepo) ListByPromotionId(ctx context.Context, pageNum, pageSize int, promotionId uint64, ptype string) ([]*domain.Material, error) {
	var materials []Material
	list := make([]*domain.Material, 0)

	db := mr.data.db.WithContext(ctx).Order("update_day DESC")

	if ptype == "product" {
		db.Where("product_id = ?", promotionId)
	} else if ptype == "aweme" {
		db.Where("aweme_id = ?", promotionId)
	}

	if result := db.
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&materials); result.Error != nil {
		return nil, result.Error
	}

	for _, material := range materials {
		list = append(list, material.ToDomain())
	}

	return list, nil
}

func (mr *materialRepo) Count(ctx context.Context, productId uint64, keyword, search, mplatform string, categories []domain.CompanyProductCategory) (int64, error) {
	var count int64

	db := mr.data.db.WithContext(ctx).Table("material_material")

	if categories != nil {
		categorySqls := make([]string, 0)

		for _, category := range categories {
			if category.IndustryId > 0 {
				categoryIds := make([]string, 0)

				for _, categoryId := range category.CategoryId {
					categoryIds = append(categoryIds, strconv.FormatUint(categoryId, 10))
				}

				if len(categoryIds) > 0 {
					categorySqls = append(categorySqls, "(industry_id = '"+strconv.FormatUint(category.IndustryId, 10)+"' and category_id in ("+strings.Join(categoryIds, ",")+"))")
				} else {
					categorySqls = append(categorySqls, "(industry_id = '"+strconv.FormatUint(category.IndustryId, 10)+"')")
				}
			}
		}

		if len(categorySqls) > 0 {
			db = db.Where(strings.Join(categorySqls, " or "))
		}
	}

	if l := utf8.RuneCountInString(mplatform); l > 0 {
		db = db.Where("platform = ?", mplatform)
	}

	if productId > 0 {
		db = db.Where("product_id = ?", productId)
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		if search == "name" {
			db = db.Where("video_name like ?", "%"+keyword+"%")
		} else if search == "product" {
			db = db.Where("product_name like ?", "%"+keyword+"%")
		} else if search == "aweme" {
			db = db.Where("(aweme_name like ? or aweme_account like ?)", "%"+keyword+"%", "%"+keyword+"%")
		}
	}

	if result := db.Model(&Material{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (mr *materialRepo) CountByPromotionId(ctx context.Context, promotionId uint64, ptype string) (int64, error) {
	var count int64

	db := mr.data.db.WithContext(ctx)

	if ptype == "product" {
		db = db.Where("product_id = ?", promotionId)
	} else if ptype == "aweme" {
		db = db.Where("aweme_id = ?", promotionId)
	}

	if result := db.Model(&Material{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (mr *materialRepo) Statistics(ctx context.Context, mtype string) (int64, error) {
	var count int64

	db := mr.data.db.WithContext(ctx).Model(&Material{})

	if mtype == "product" {
		db = db.Group("product_id")
	} else if mtype == "aweme" {
		db = db.Group("aweme_id")
	}

	if result := db.Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (mr *materialRepo) StatisticsAwemeIndustry(ctx context.Context, awemeId uint64) ([]*domain.MaterialAwemeIndustry, error) {
	var materials []Material
	list := make([]*domain.MaterialAwemeIndustry, 0)

	db := mr.data.db.WithContext(ctx).Where("aweme_id = ?", awemeId)

	if result := db.Select("industry_id,max(industry_name) industry_name,sum(total_item_num) as total_item_num").
		Group("industry_id").Order("total_item_num DESC").
		Find(&materials); result.Error != nil {
		return nil, result.Error
	}

	for _, material := range materials {
		list = append(list, &domain.MaterialAwemeIndustry{
			IndustryId:   material.IndustryId,
			IndustryName: material.IndustryName,
			TotalItemNum: material.TotalItemNum,
		})
	}

	return list, nil
}

func (mr *materialRepo) UpdateVideoUrl(ctx context.Context, videoId uint64, videoUrl string) error {
	if result := mr.data.db.WithContext(ctx).Model(Material{}).
		Where("video_id", videoId).
		Updates(map[string]interface{}{"video_url": videoUrl}); result.Error != nil {
		return result.Error
	}

	return nil
}

func (mr *materialRepo) SaveCacheHash(ctx context.Context, key string, val map[string]string, timeout time.Duration) error {
	_, err := mr.data.rdb.HMSet(ctx, key, val).Result()

	if err != nil {
		return err
	}

	_, err = mr.data.rdb.Expire(ctx, key, timeout).Result()

	if err != nil {
		return err
	}

	return nil
}

func (mr *materialRepo) UpdateCacheHash(ctx context.Context, key string, val map[string]string) error {
	_, err := mr.data.rdb.HMSet(ctx, key, val).Result()

	if err != nil {
		return err
	}

	return nil
}

func (mr *materialRepo) GetCacheHash(ctx context.Context, key string, field string) (string, error) {
	val, err := mr.data.rdb.HGet(ctx, key, field).Result()

	if err != nil {
		return "", err
	}

	return val, nil
}

func (mr *materialRepo) DeleteCache(ctx context.Context, key string) error {
	if _, err := mr.data.rdb.Del(ctx, key).Result(); err != nil {
		return err
	}

	return nil
}

func (mr *materialRepo) GetObject(ctx context.Context, key string) (*tos.GetObjectV2Output, error) {
	object, err := mr.data.tos.GetObject(ctx, key)

	if err != nil {
		return nil, err
	}

	return object, nil
}

func (mr *materialRepo) CreateMultipartUpload(ctx context.Context, objectKey string) (*tos.CreateMultipartUploadV2Output, error) {
	createMultipartOutput, err := mr.data.tos.CreateMultipartUpload(ctx, objectKey)

	if err != nil {
		return nil, err
	}

	return createMultipartOutput, nil
}

func (mr *materialRepo) UploadPart(ctx context.Context, partNumber int, contentLength int64, objectKey, uploadId string, content io.Reader) (*tos.UploadPartV2Output, error) {
	partOutput, err := mr.data.tos.UploadPart(ctx, partNumber, contentLength, objectKey, uploadId, content)

	if err != nil {
		return nil, err
	}

	return partOutput, nil
}

func (mr *materialRepo) CompleteMultipartUpload(ctx context.Context, fileName, uploadId string, parts []tos.UploadedPartV2) (*tos.CompleteMultipartUploadV2Output, error) {
	completeOutput, err := mr.data.tos.CompleteMultipartUpload(ctx, fileName, uploadId, parts)

	if err != nil {
		return nil, err
	}

	return completeOutput, nil
}

func (mr *materialRepo) AbortMultipartUpload(ctx context.Context, objectKey, uploadId string) (*tos.AbortMultipartUploadOutput, error) {
	output, err := mr.data.tos.AbortMultipartUpload(ctx, objectKey, uploadId)

	if err != nil {
		return nil, err
	}

	return output, nil
}

func (mr *materialRepo) Send(ctx context.Context, message event.Event) error {
	if err := mr.data.kafka.Send(ctx, message); err != nil {
		return err
	}

	return nil
}

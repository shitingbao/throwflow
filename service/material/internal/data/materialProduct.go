package data

import (
	"context"
	"material/internal/biz"
	"material/internal/domain"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-kratos/kratos/v2/log"
)

// 市场素材爆品表
type MaterialProduct struct {
	Id           uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	ProductId    uint64    `gorm:"column:product_id;type:bigint(20) UNSIGNED;not null;default:0;comment:商品ID"`
	ProductName  string    `gorm:"column:product_name;type:text;not null;comment:商品名称"`
	VideoLike    uint64    `gorm:"column:video_like;type:bigint(20) UNSIGNED;not null;default:0;comment:素材点赞数"`
	IndustryId   uint64    `gorm:"column:industry_id;type:bigint(20) UNSIGNED;not null;default:0;comment:行业ID"`
	IndustryName string    `gorm:"column:industry_name;type:varchar(50);not null;comment:行业名称"`
	CategoryId   uint64    `gorm:"column:category_id;type:bigint(20) UNSIGNED;not null;default:0;comment:分类ID"`
	CategoryName string    `gorm:"column:category_name;type:varchar(50);not null;comment:分类名称"`
	IsHot        uint8     `gorm:"column:is_hot;type:tinyint(3) UNSIGNED;not null;default:0;comment:是否爆单，0：不是，1：是"`
	Videos       uint64    `gorm:"column:videos;type:bigint(20) UNSIGNED;not null;default:0;comment:视频数"`
	Awemes       uint64    `gorm:"column:awemes;type:bigint(20) UNSIGNED;not null;default:0;comment:达人数"`
	Platform     string    `gorm:"column:platform;type:varchar(10);not null;comment:平台：dy：抖音，ks：快手"`
	UpdateDay    time.Time `gorm:"column:update_day;type:datetime;not null;comment:更新时间"`
	CreateTime   time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime   time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (MaterialProduct) TableName() string {
	return "material_material_product"
}

type materialProductRepo struct {
	data *Data
	log  *log.Helper
}

func (mp *MaterialProduct) ToDomain() *domain.MaterialProduct {
	return &domain.MaterialProduct{
		Id:           mp.Id,
		ProductId:    mp.ProductId,
		ProductName:  mp.ProductName,
		VideoLike:    mp.VideoLike,
		IndustryId:   mp.IndustryId,
		IndustryName: mp.IndustryName,
		CategoryId:   mp.CategoryId,
		CategoryName: mp.CategoryName,
		IsHot:        mp.IsHot,
		Videos:       mp.Videos,
		Awemes:       mp.Awemes,
		Platform:     mp.Platform,
		UpdateDay:    mp.UpdateDay,
		CreateTime:   mp.CreateTime,
		UpdateTime:   mp.UpdateTime,
	}
}

func NewMaterialProductRepo(data *Data, logger log.Logger) biz.MaterialProductRepo {
	return &materialProductRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (mpr *materialProductRepo) List(ctx context.Context, pageNum, pageSize int, keyword, search, sortBy, mplatform string, categories []domain.CompanyProductCategory) ([]*domain.MaterialProduct, error) {
	var materialProducts []MaterialProduct
	list := make([]*domain.MaterialProduct, 0)

	db := mpr.data.db.WithContext(ctx)

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
	
	if l := utf8.RuneCountInString(keyword); l > 0 {
		if search == "product" {
			db = db.Where("product_name like ?", "%"+keyword+"%")
		}
	}

	if result := db.Order(sortBy + " DESC").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&materialProducts); result.Error != nil {
		return nil, result.Error
	}

	for _, materialProduct := range materialProducts {
		list = append(list, materialProduct.ToDomain())
	}

	return list, nil
}

func (mpr *materialProductRepo) Count(ctx context.Context, keyword, search, mplatform string, categories []domain.CompanyProductCategory) (int64, error) {
	var count int64

	db := mpr.data.db.WithContext(ctx)

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

	if l := utf8.RuneCountInString(keyword); l > 0 {
		if search == "product" {
			db = db.Where("product_name like ?", "%"+keyword+"%")
		}
	}

	if result := db.Model(&MaterialProduct{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

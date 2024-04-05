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

// 我的收藏表
type Collect struct {
	CompanyId  uint64    `gorm:"column:company_id;type:bigint(20) UNSIGNED;not null;comment:企业库ID"`
	Phone      string    `gorm:"column:phone;type:varchar(20);not null;comment:手机号"`
	VideoId    uint64    `gorm:"column:video_id;type:bigint(20) UNSIGNED;not null;comment:素材ID"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (Collect) TableName() string {
	return "material_collect"
}

type collectrepo struct {
	data *Data
	log  *log.Helper
}

func (c *Collect) ToDomain() *domain.Collect {
	return &domain.Collect{
		CompanyId:  c.CompanyId,
		Phone:      c.Phone,
		VideoId:    c.VideoId,
		CreateTime: c.CreateTime,
		UpdateTime: c.UpdateTime,
	}
}

func NewCollectRepo(data *Data, logger log.Logger) biz.CollectRepo {
	return &collectrepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cr *collectrepo) Get(ctx context.Context, companyId, videoId uint64, phone string) (*domain.Collect, error) {
	collect := &Collect{}

	if result := cr.data.db.WithContext(ctx).
		Where("company_id = ?", companyId).
		Where("phone = ?", phone).
		Where("video_id = ?", videoId).
		First(collect); result.Error != nil {
		return nil, result.Error
	}

	return collect.ToDomain(), nil
}

func (cr *collectrepo) ListByVideoIds(ctx context.Context, companyId uint64, phone string, videoIds []string) ([]*domain.Collect, error) {
	var collects []Collect
	list := make([]*domain.Collect, 0)

	if result := cr.data.db.WithContext(ctx).
		Where("company_id = ?", companyId).
		Where("phone = ?", phone).
		Where("video_id in ? ", videoIds).
		Find(&collects); result.Error != nil {
		return nil, result.Error
	}

	for _, collect := range collects {
		list = append(list, collect.ToDomain())
	}

	return list, nil
}

func (cr *collectrepo) List(ctx context.Context, pageNum, pageSize int, companyId uint64, phone, keyword, search, sortBy, mplatform string, categories []domain.CompanyProductCategory) ([]*domain.Material, error) {
	var materials []Material
	list := make([]*domain.Material, 0)

	db := cr.data.db.WithContext(ctx).
		Table("material_collect").
		Select("material_material.*").
		Joins("left join material_material on material_collect.video_id = material_material.video_id").
		Where("material_collect.phone = ?", phone).
		Where("material_collect.company_id = ?", companyId).
		Where("material_material.video_id != ''")

	if categories != nil {
		categorySqls := make([]string, 0)

		for _, category := range categories {
			if category.IndustryId > 0 {
				categoryIds := make([]string, 0)

				for _, categoryId := range category.CategoryId {
					categoryIds = append(categoryIds, strconv.FormatUint(categoryId, 10))
				}

				if len(categoryIds) > 0 {
					categorySqls = append(categorySqls, "(material_material.industry_id = '"+strconv.FormatUint(category.IndustryId, 10)+"' and material_material.category_id in ("+strings.Join(categoryIds, ",")+"))")
				} else {
					categorySqls = append(categorySqls, "(material_material.industry_id = '"+strconv.FormatUint(category.IndustryId, 10)+"')")
				}
			}
		}

		if len(categorySqls) > 0 {
			db = db.Where(strings.Join(categorySqls, " or "))
		}
	}

	if l := utf8.RuneCountInString(mplatform); l > 0 {
		db = db.Where("material_material.platform = ?", mplatform)
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		if search == "name" {
			db = db.Where("material_material.video_name like ?", "%"+keyword+"%")
		} else if search == "product" {
			db = db.Where("material_material.product_name like ?", "%"+keyword+"%")
		} else if search == "aweme" {
			db = db.Where("(material_material.aweme_name like ? or material_material.aweme_account like ?)", "%"+keyword+"%", "%"+keyword+"%")
		}
	}

	if result := db.Order("material_material." + sortBy + " DESC").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Scan(&materials); result.Error != nil {
		return nil, result.Error
	}

	for _, material := range materials {
		list = append(list, material.ToDomain())
	}

	return list, nil
}

func (cr *collectrepo) Count(ctx context.Context, companyId uint64, phone, keyword, search, mplatform string, categories []domain.CompanyProductCategory) (int64, error) {
	var count int64

	db := cr.data.db.WithContext(ctx).
		Joins("left join material_material on material_collect.video_id = material_material.video_id").
		Where("material_collect.phone = ?", phone).
		Where("material_collect.company_id = ?", companyId).
		Where("material_material.video_id != ''")

	if categories != nil {
		categorySqls := make([]string, 0)

		for _, category := range categories {
			if category.IndustryId > 0 {
				categoryIds := make([]string, 0)

				for _, categoryId := range category.CategoryId {
					categoryIds = append(categoryIds, strconv.FormatUint(categoryId, 10))
				}

				if len(categoryIds) > 0 {
					categorySqls = append(categorySqls, "(material_material.industry_id = '"+strconv.FormatUint(category.IndustryId, 10)+"' and material_material.category_id in ("+strings.Join(categoryIds, ",")+"))")
				} else {
					categorySqls = append(categorySqls, "(material_material.industry_id = '"+strconv.FormatUint(category.IndustryId, 10)+"')")
				}
			}
		}

		if len(categorySqls) > 0 {
			db = db.Where(strings.Join(categorySqls, " or "))
		}
	}

	if l := utf8.RuneCountInString(mplatform); l > 0 {
		db = db.Where("material_material.platform = ?", mplatform)
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		if search == "name" {
			db = db.Where("material_material.video_name like ?", "%"+keyword+"%")
		} else if search == "product" {
			db = db.Where("material_material.product_name like ?", "%"+keyword+"%")
		} else if search == "aweme" {
			db = db.Where("(material_material.aweme_name like ? or material_material.aweme_account like ?)", "%"+keyword+"%", "%"+keyword+"%")
		}
	}

	if result := db.Model(&Collect{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (cr *collectrepo) Save(ctx context.Context, in *domain.Collect) (*domain.Collect, error) {
	collect := &Collect{
		CompanyId:  in.CompanyId,
		Phone:      in.Phone,
		VideoId:    in.VideoId,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
	}

	if result := cr.data.db.WithContext(ctx).Create(collect); result.Error != nil {
		return nil, result.Error
	}

	return collect.ToDomain(), nil
}

func (cr *collectrepo) Delete(ctx context.Context, in *domain.Collect) error {
	if result := cr.data.db.WithContext(ctx).Where("company_id = ?", in.CompanyId).Where("phone = ?", in.Phone).Where("video_id = ?", in.VideoId).Delete(&Collect{}); result.Error != nil {
		return result.Error
	}

	return nil
}

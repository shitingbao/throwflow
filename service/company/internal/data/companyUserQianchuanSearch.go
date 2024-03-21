package data

import (
	"company/internal/biz"
	"company/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// 企业库用户广告管理搜索表
type CompanyUserQianchuanSearch struct {
	Id          uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	CompanyId   uint64    `gorm:"column:company_id;type:bigint(20) UNSIGNED;not null;comment:企业库ID"`
	SearchType  string    `gorm:"column:search_type;type:enum('advertiser','campaign','ad');not null;default:'advertiser';comment:搜索类型，advertiser：千川账户 campaign：千川组 ad：千川计划"`
	SearchValue string    `gorm:"column:search_value;type:varchar(250);not null;comment:广告管理搜索值"`
	SearchCount uint32    `gorm:"column:search_count;type:int(11) UNSIGNED;not null;comment:查询次数"`
	Day         uint32    `gorm:"column:day;type:int(11) UNSIGNED;not null;comment:查询日期"`
	CreateTime  time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime  time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (CompanyUserQianchuanSearch) TableName() string {
	return "company_company_user_qianchuan_search"
}

type companyUserQianchuanSearchRepo struct {
	data *Data
	log  *log.Helper
}

func (cuqs *CompanyUserQianchuanSearch) ToDomain() *domain.CompanyUserQianchuanSearch {
	companyUserQianchuanSearch := &domain.CompanyUserQianchuanSearch{
		Id:          cuqs.Id,
		CompanyId:   cuqs.CompanyId,
		SearchType:  cuqs.SearchType,
		SearchValue: cuqs.SearchValue,
		SearchCount: cuqs.SearchCount,
		Day:         cuqs.Day,
		CreateTime:  cuqs.CreateTime,
		UpdateTime:  cuqs.UpdateTime,
	}

	return companyUserQianchuanSearch
}

func NewCompanyUserQianchuanSearchRepo(data *Data, logger log.Logger) biz.CompanyUserQianchuanSearchRepo {
	return &companyUserQianchuanSearchRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cuqsr *companyUserQianchuanSearchRepo) Get(ctx context.Context, companyId uint64, searchType, searchValue string) (*domain.CompanyUserQianchuanSearch, error) {
	companyUserQianchuanSearch := &CompanyUserQianchuanSearch{}

	if result := cuqsr.data.db.WithContext(ctx).Where("company_id = ?", companyId).Where("search_type = ?", searchType).Where("search_value = ?", searchValue).First(companyUserQianchuanSearch); result.Error != nil {
		return nil, result.Error
	}

	return companyUserQianchuanSearch.ToDomain(), nil
}

func (cuqsr *companyUserQianchuanSearchRepo) List(ctx context.Context, companyId uint64, startDay, endDay uint32, searchType string) ([]*domain.CompanyUserQianchuanSearch, error) {
	var companyUserQianchuanSearches []CompanyUserQianchuanSearch
	list := make([]*domain.CompanyUserQianchuanSearch, 0)

	db := cuqsr.data.db.WithContext(ctx).Where("company_id = ?", companyId).Where("search_type = ?", searchType)

	if endDay > 0 {
		db = db.Where("day >= ?", startDay)
		db = db.Where("day <= ?", endDay)

		if result := db.Select("search_value,sum(search_count) as search_count").Group("search_value").Order("search_count DESC").
			Find(&companyUserQianchuanSearches); result.Error != nil {
			return nil, result.Error
		}
	} else {
		db = db.Where("day = ?", startDay)

		if result := db.Order("search_count DESC").
			Find(&companyUserQianchuanSearches); result.Error != nil {
			return nil, result.Error
		}
	}

	for _, companyUserQianchuanSearch := range companyUserQianchuanSearches {
		list = append(list, companyUserQianchuanSearch.ToDomain())
	}

	return list, nil
}

func (cuqsr *companyUserQianchuanSearchRepo) Save(ctx context.Context, in *domain.CompanyUserQianchuanSearch) (*domain.CompanyUserQianchuanSearch, error) {
	companyUserQianchuanSearch := &CompanyUserQianchuanSearch{
		Id:          in.Id,
		CompanyId:   in.CompanyId,
		SearchType:  in.SearchType,
		SearchValue: in.SearchValue,
		SearchCount: in.SearchCount,
		Day:         in.Day,
		CreateTime:  in.CreateTime,
		UpdateTime:  in.UpdateTime,
	}

	if result := cuqsr.data.DB(ctx).Create(companyUserQianchuanSearch); result.Error != nil {
		return nil, result.Error
	}

	return companyUserQianchuanSearch.ToDomain(), nil
}

func (cuqsr *companyUserQianchuanSearchRepo) Update(ctx context.Context, in *domain.CompanyUserQianchuanSearch) (*domain.CompanyUserQianchuanSearch, error) {
	companyUserQianchuanSearch := &CompanyUserQianchuanSearch{
		Id:          in.Id,
		CompanyId:   in.CompanyId,
		SearchType:  in.SearchType,
		SearchValue: in.SearchValue,
		SearchCount: in.SearchCount,
		Day:         in.Day,
		CreateTime:  in.CreateTime,
		UpdateTime:  in.UpdateTime,
	}

	if result := cuqsr.data.DB(ctx).Save(companyUserQianchuanSearch); result.Error != nil {
		return nil, result.Error
	}

	return companyUserQianchuanSearch.ToDomain(), nil
}

package domain

import (
	"context"
	"time"
)

type CompanyUserQianchuanSearch struct {
	Id          uint64
	CompanyId   uint64
	SearchType  string
	SearchValue string
	SearchCount uint32
	Day         uint32
	CreateTime  time.Time
	UpdateTime  time.Time
}

type CompanyUserQianchuanDataSearch struct {
	Id   uint64
	Name string
}

func NewCompanyUserQianchuanSearch(ctx context.Context, companyId uint64, day, searchCount uint32, searchType, searchValue string) *CompanyUserQianchuanSearch {
	return &CompanyUserQianchuanSearch{
		CompanyId:   companyId,
		SearchType:  searchType,
		SearchValue: searchValue,
		SearchCount: searchCount,
		Day:         day,
	}
}

func (cuqs *CompanyUserQianchuanSearch) SetCompanyId(ctx context.Context, companyId uint64) {
	cuqs.CompanyId = companyId
}

func (cuqs *CompanyUserQianchuanSearch) SetSearchType(ctx context.Context, searchType string) {
	cuqs.SearchType = searchType
}

func (cuqs *CompanyUserQianchuanSearch) SetSearchValue(ctx context.Context, searchValue string) {
	cuqs.SearchValue = searchValue
}

func (cuqs *CompanyUserQianchuanSearch) SetSearchCount(ctx context.Context, searchCount uint32) {
	cuqs.SearchCount = searchCount
}

func (cuqs *CompanyUserQianchuanSearch) SetDay(ctx context.Context, day uint32) {
	cuqs.Day = day
}

func (cuqs *CompanyUserQianchuanSearch) SetUpdateTime(ctx context.Context) {
	cuqs.UpdateTime = time.Now()
}

func (cuqs *CompanyUserQianchuanSearch) SetCreateTime(ctx context.Context) {
	cuqs.CreateTime = time.Now()
}

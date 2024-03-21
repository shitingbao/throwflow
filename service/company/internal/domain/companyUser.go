package domain

import (
	"context"
	"time"
)

type CompanyUser struct {
	Id               uint64
	Username         string
	Job              string
	Phone            string
	CompanyId        uint64
	CurrentCompanyId uint64
	Role             uint8
	RoleName         string
	Status           uint8
	OrganizationId   uint64
	IsWhite          uint8
	CreateTime       time.Time
	UpdateTime       time.Time
	Company          *Company
}

type LoginCompanyUserCompany struct {
	Id          uint64
	CompanyName string
}

type LoginCompanyUser struct {
	Id                   uint64
	CompanyId            uint64
	Username             string
	Job                  string
	Phone                string
	Role                 uint8
	RoleName             string
	IsWhite              uint8
	CompanyType          uint8
	CompanyTypeName      string
	CompanyName          string
	CompanyStartTime     time.Time
	CompanyEndTime       time.Time
	Accounts             uint32
	QianchuanAdvertisers uint32
	IsTermwork           uint8
	Token                string
	CurrentCompanyId     uint64
	Reason               string
	UserCompany          []*LoginCompanyUserCompany
}

type CompanyUserList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*CompanyUser
}

type CompanyUserQianchuanAdvertiser struct {
	AdvertiserId   uint64
	CompanyId      uint64
	AccountId      uint64
	AdvertiserName string
	CompanyName    string
	UserId         uint64
	Username       string
}

type CompanyUserQianchuanAdvertiserList struct {
	PageNum   uint64
	PageSize  uint64
	Total     uint64
	TotalPage uint64
	List      []*CompanyUserQianchuanAdvertiser
}

type CompanyUserQianchuanReportAdvertiserData struct {
	Time           string
	Value          float64
	YesterdayValue float64
}

type CompanyUserQianchuanReportAdvertiserList struct {
	Key  string
	List []*CompanyUserQianchuanReportAdvertiserData
}

type CompanyUserQianchuanCampaign struct {
	CampaignId   uint64
	CampaignName string
	IsSelect     uint8
}

type ExternalSelectQianchuanAdvertiser struct {
	AdvertiserId   uint64
	AdvertiserName string
}

type ExternalSelectQianchuanCampaign struct {
	CampaignId     uint64
	AdvertiserId   uint64
	CampaignName   string
	Budget         string
	BudgetMode     string
	MarketingGoal  string
	MarketingScene string
}

type AccountOpend struct {
	Phone    string
	RoleName string
}

type Role struct {
	Key   string
	Value string
}

type SelectCompanyUsers struct {
	Role []*Role
}

type QianchuanAdvertiserList struct {
	AdvertiserId   uint64
	Status         uint32
	IsSelect       uint32
	AdvertiserName string
	CompanyName    string
}

type StatisticsCompanyUser struct {
	Key   string
	Value string
}

type StatisticsCompanyUsers struct {
	Statistics []*StatisticsCompanyUser
}

func NewSelectCompanyUsers() *SelectCompanyUsers {
	role := make([]*Role, 0)

	role = append(role, &Role{Key: "1", Value: "主管理员"})
	role = append(role, &Role{Key: "2", Value: "副管理员"})
	role = append(role, &Role{Key: "3", Value: "普通成员"})

	return &SelectCompanyUsers{
		Role: role,
	}
}

func NewCompanyUser(ctx context.Context, companyId uint64, username, job, phone string, role, status uint8) *CompanyUser {
	return &CompanyUser{
		Username:  username,
		Job:       job,
		Phone:     phone,
		CompanyId: companyId,
		Role:      role,
		Status:    status,
	}
}

func (cu *CompanyUser) SetUsername(ctx context.Context, username string) {
	cu.Username = username
}

func (cu *CompanyUser) SetJob(ctx context.Context, job string) {
	cu.Job = job
}

func (cu *CompanyUser) SetPhone(ctx context.Context, phone string) {
	cu.Phone = phone
}

func (cu *CompanyUser) SetRole(ctx context.Context, role uint8) {
	cu.Role = role
}

func (cu *CompanyUser) SetStatus(ctx context.Context, status uint8) {
	cu.Status = status
}

func (cu *CompanyUser) SetOrganizationId(ctx context.Context, organizationId uint64) {
	cu.OrganizationId = organizationId
}

func (cu *CompanyUser) SetIsWhite(ctx context.Context, isWhite uint8) {
	cu.IsWhite = isWhite
}

func (cu *CompanyUser) SetUpdateTime(ctx context.Context) {
	cu.UpdateTime = time.Now()
}

func (cu *CompanyUser) SetCreateTime(ctx context.Context) {
	cu.CreateTime = time.Now()
}

package domain

import (
	"context"
	"time"
)

type Company struct {
	Id                   uint64
	ClueId               uint64
	CompanyType          uint8
	CompanyTypeName      string
	Status               uint8
	StartTime            time.Time
	EndTime              time.Time
	MenuId               string
	Accounts             uint32
	QianchuanAdvertisers uint32
	MiniQrCodeUrl        string
	IsTermwork           uint8
	IsDel                uint8
	CreateTime           time.Time
	UpdateTime           time.Time
	Clue                 *Clue
	CompanyUser          []*CompanyUser
}

type CompanyList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*Company
}

type SelectCompanys struct {
	Status      []*Status
	CompanyType []*CompanyType
}

type StatisticsCompany struct {
	Key   string
	Value string
}

type StatisticsCompanys struct {
	Statistics []*StatisticsCompany
}

type PermissionCodes struct {
	Name string
	Code string
}

func NewSelectCompanys() *SelectCompanys {
	status := make([]*Status, 0)

	status = append(status, &Status{Key: "0", Value: "禁用"})
	status = append(status, &Status{Key: "1", Value: "启用"})
	status = append(status, &Status{Key: "2", Value: "过期"})

	companyType := make([]*CompanyType, 0)

	companyType = append(companyType, &CompanyType{Key: "1", Value: "试用版"})
	companyType = append(companyType, &CompanyType{Key: "2", Value: "基础版"})
	companyType = append(companyType, &CompanyType{Key: "3", Value: "专业版"})
	companyType = append(companyType, &CompanyType{Key: "4", Value: "旗舰版"})
	companyType = append(companyType, &CompanyType{Key: "5", Value: "尊享版"})

	return &SelectCompanys{
		Status:      status,
		CompanyType: companyType,
	}
}

func NewPermissionCodes() []*PermissionCodes {
	permissionCodes := make([]*PermissionCodes, 0)

	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "大盘播报", Code: "broadcast"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "广告管理", Code: "ad"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "策略管理", Code: "strategy"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "执行日志", Code: "log"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "智能创建", Code: "create"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "账户授权", Code: "account"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "智能诊断", Code: "diagnosis"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "团队绩效", Code: "team"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "素材仓库", Code: "warehouse"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "素材报表", Code: "mreports"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "素材参谋", Code: "material"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "在线剪辑", Code: "editing"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "外包创作", Code: "outsourcing"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "操作指南", Code: "operate"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "更新日志", Code: "updatelog"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "意见反馈", Code: "feedback"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "代投服务", Code: "pservices"})
	permissionCodes = append(permissionCodes, &PermissionCodes{Name: "在线续费", Code: "online"})

	return permissionCodes
}

func NewCompany(ctx context.Context, clueId uint64) *Company {
	return &Company{
		ClueId: clueId,
	}
}

func (c *Company) SetStatus(ctx context.Context, status uint8) {
	c.Status = status
}

func (c *Company) SetIsDel(ctx context.Context, isDel uint8) {
	c.IsDel = isDel
}

func (c *Company) SetMenuId(ctx context.Context, menuId string) {
	c.MenuId = menuId
}

func (c *Company) SetAccounts(ctx context.Context, accounts uint32) {
	c.Accounts = accounts
}

func (c *Company) SetQianchuanAdvertisers(ctx context.Context, qianchuanAdvertisers uint32) {
	c.QianchuanAdvertisers = qianchuanAdvertisers
}

func (c *Company) SetCompanyType(ctx context.Context, companyType uint8) {
	c.CompanyType = companyType
}

func (c *Company) SetMiniQrCodeUrl(ctx context.Context, miniQrCodeUrl string) {
	c.MiniQrCodeUrl = miniQrCodeUrl
}

func (c *Company) SetIsTermwork(ctx context.Context, isTermwork uint8) {
	c.IsTermwork = isTermwork
}

func (c *Company) SetStartTime(ctx context.Context, time time.Time) {
	c.StartTime = time
}

func (c *Company) SetEndTime(ctx context.Context, time time.Time) {
	c.EndTime = time
}

func (c *Company) SetUpdateTime(ctx context.Context) {
	c.UpdateTime = time.Now()
}

func (c *Company) SetCreateTime(ctx context.Context) {
	c.CreateTime = time.Now()
}

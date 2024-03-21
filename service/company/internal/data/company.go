package data

import (
	"company/internal/biz"
	"company/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"time"
	"unicode/utf8"
)

// 企业库表
type Company struct {
	Id                   uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	ClueId               uint64    `gorm:"column:clue_id;type:bigint(20) UNSIGNED;not null;comment:线索库ID"`
	CompanyType          uint8     `gorm:"column:company_type;type:tinyint(3) UNSIGNED;not null;default:1;comment:公司类型，1：试用版，2：基础版，3：专业版，4：旗舰版，5：尊享版"`
	Status               uint8     `gorm:"column:status;type:tinyint(3) UNSIGNED;not null;default:0;comment:状态：1：启用，0：禁用, 2: 过期"`
	StartTime            time.Time `gorm:"column:start_time;type:date;not null;comment:开始时间"`
	EndTime              time.Time `gorm:"column:end_time;type:date;not null;comment:到期时间"`
	MenuId               string    `gorm:"column:menu_id;type:varchar(50);not null;comment:菜单ID"`
	Accounts             uint32    `gorm:"column:accounts;type:int(10) UNSIGNED;not null;default:1;comment:账户数"`
	QianchuanAdvertisers uint32    `gorm:"column:qianchuan_advertisers;type:int(10) UNSIGNED;not null;default:0;comment:千川账户数"`
	MiniQrCodeUrl        string    `gorm:"column:mini_qr_code_url;type:varchar(250);not null;comment:公司小程序码 URL"`
	IsTermwork           uint8     `gorm:"column:is_termwork;type:tinyint(3) UNSIGNED;not null;default:1;comment:是否开启团队协同，1：启用，0：禁用"`
	IsDel                uint8     `gorm:"column:is_del;type:tinyint(3) UNSIGNED;not null;default:0;comment:是否移除，1：已移除，0：未移除"`
	CreateTime           time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime           time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (Company) TableName() string {
	return "company_company"
}

type companyRepo struct {
	data *Data
	log  *log.Helper
}

func (c *Company) ToDomain(ctx context.Context) *domain.Company {
	company := &domain.Company{
		Id:                   c.Id,
		ClueId:               c.ClueId,
		CompanyType:          c.CompanyType,
		Status:               c.Status,
		StartTime:            c.StartTime,
		EndTime:              c.EndTime,
		MenuId:               c.MenuId,
		Accounts:             c.Accounts,
		QianchuanAdvertisers: c.QianchuanAdvertisers,
		MiniQrCodeUrl:        c.MiniQrCodeUrl,
		IsTermwork:           c.IsTermwork,
		CreateTime:           c.CreateTime,
		UpdateTime:           c.UpdateTime,
		IsDel:                c.IsDel,
	}

	selectCompanys := domain.NewSelectCompanys()

	for _, companyType := range selectCompanys.CompanyType {
		icompanyType, _ := strconv.Atoi(companyType.Key)

		if uint8(icompanyType) == company.CompanyType {
			company.CompanyTypeName = companyType.Value
			break
		}
	}

	return company
}

func NewCompanyRepo(data *Data, logger log.Logger) biz.CompanyRepo {
	return &companyRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cr *companyRepo) GetById(ctx context.Context, id uint64) (*domain.Company, error) {
	company := &Company{}

	if result := cr.data.db.WithContext(ctx).Where("is_del = ?", 0).First(company, id); result.Error != nil {
		return nil, result.Error
	}

	return company.ToDomain(ctx), nil
}

func (cr *companyRepo) GetByClueId(ctx context.Context, clueId uint64) (*domain.Company, error) {
	company := &Company{}

	if result := cr.data.db.WithContext(ctx).Where("clue_id = ?", clueId).First(company); result.Error != nil {
		return nil, result.Error
	}

	return company.ToDomain(ctx), nil
}

func (cr *companyRepo) List(ctx context.Context, pageNum, pageSize int, industryId uint64, keyword, status string, companyType uint8) ([]*domain.Company, error) {
	var companys []Company
	list := make([]*domain.Company, 0)

	db := cr.data.db.WithContext(ctx).
		Joins("left join company_clue on company_clue.id = company_company.clue_id").
		Where("company_company.is_del = ?", 0)

	if industryId > 0 {
		db = db.Where("find_in_set(?, company_clue.`industry_id`)", industryId)
	}

	if len(status) > 0 {
		db = db.Where("company_company.status = ?", status)
	}

	if companyType > 0 {
		db = db.Where("company_company.company_type = ?", companyType)
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		db = db.Where("(company_clue.company_name like ? or company_clue.follower like ?)", "%"+keyword+"%", "%"+keyword+"%")
	}

	if pageNum == 0 {
		if result := db.Order("company_company.status ASC,company_company.id DESC").
			Find(&companys); result.Error != nil {
			return nil, result.Error
		}
	} else {
		if result := db.Order("company_company.status ASC,company_company.id DESC").
			Limit(pageSize).Offset((pageNum - 1) * pageSize).
			Find(&companys); result.Error != nil {
			return nil, result.Error
		}
	}

	for _, company := range companys {
		list = append(list, company.ToDomain(ctx))
	}

	return list, nil
}

func (cr *companyRepo) Count(ctx context.Context, industryId uint64, keyword, status string, companyType uint8) (int64, error) {
	var count int64

	db := cr.data.db.WithContext(ctx).
		Joins("left join company_clue on company_clue.id = company_company.clue_id").
		Where("company_company.is_del = ?", 0)

	if len(status) > 0 {
		db = db.Where("company_company.status = ?", status)
	}

	if companyType > 0 {
		db = db.Where("company_company.company_type = ?", companyType)
	}

	if industryId > 0 {
		db = db.Where("find_in_set(?, company_clue.`industry_id`)", industryId)
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		db = db.Where("(company_clue.company_name like ? or company_clue.follower like ?)", "%"+keyword+"%", "%"+keyword+"%")
	}

	if result := db.Model(&Company{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (cr *companyRepo) Statistics(ctx context.Context, status uint8) (int64, error) {
	var count int64

	if result := cr.data.db.WithContext(ctx).Model(&Company{}).Where("company_type = ?", status).Where("is_del = ?", 0).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (cr *companyRepo) Save(ctx context.Context, in *domain.Company) (*domain.Company, error) {
	company := &Company{
		ClueId:               in.ClueId,
		CompanyType:          in.CompanyType,
		Status:               in.Status,
		StartTime:            in.StartTime,
		EndTime:              in.EndTime,
		MenuId:               in.MenuId,
		Accounts:             in.Accounts,
		QianchuanAdvertisers: in.QianchuanAdvertisers,
		MiniQrCodeUrl:        in.MiniQrCodeUrl,
		IsTermwork:           in.IsTermwork,
		IsDel:                in.IsDel,
		CreateTime:           in.CreateTime,
		UpdateTime:           in.UpdateTime,
	}

	if result := cr.data.DB(ctx).Create(company); result.Error != nil {
		return nil, result.Error
	}

	return company.ToDomain(ctx), nil
}

func (cr *companyRepo) Update(ctx context.Context, in *domain.Company) (*domain.Company, error) {
	company := &Company{
		Id:                   in.Id,
		ClueId:               in.ClueId,
		CompanyType:          in.CompanyType,
		Status:               in.Status,
		StartTime:            in.StartTime,
		EndTime:              in.EndTime,
		MenuId:               in.MenuId,
		Accounts:             in.Accounts,
		QianchuanAdvertisers: in.QianchuanAdvertisers,
		MiniQrCodeUrl:        in.MiniQrCodeUrl,
		IsTermwork:           in.IsTermwork,
		IsDel:                in.IsDel,
		CreateTime:           in.CreateTime,
		UpdateTime:           in.UpdateTime,
	}

	if result := cr.data.DB(ctx).Save(company); result.Error != nil {
		return nil, result.Error
	}

	return company.ToDomain(ctx), nil
}

func (cr *companyRepo) Delete(ctx context.Context, in *domain.Company) error {
	company := &Company{
		Id:                   in.Id,
		ClueId:               in.ClueId,
		CompanyType:          in.CompanyType,
		Status:               in.Status,
		StartTime:            in.StartTime,
		EndTime:              in.EndTime,
		MenuId:               in.MenuId,
		Accounts:             in.Accounts,
		QianchuanAdvertisers: in.QianchuanAdvertisers,
		MiniQrCodeUrl:        in.MiniQrCodeUrl,
		IsTermwork:           in.IsTermwork,
		IsDel:                in.IsDel,
		CreateTime:           in.CreateTime,
		UpdateTime:           in.UpdateTime,
	}

	if result := cr.data.DB(ctx).Delete(company); result.Error != nil {
		return result.Error
	}

	return nil
}

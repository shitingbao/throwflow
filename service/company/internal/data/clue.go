package data

import (
	"company/internal/biz"
	"company/internal/domain"
	"context"
	"time"
	"unicode/utf8"

	"github.com/go-kratos/kratos/v2/log"
)

// 线索表
type Clue struct {
	Id                 uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	CompanyName        string    `gorm:"column:company_name;type:varchar(250);not null;comment:企业名称"`
	IndustryId         string    `gorm:"column:industry_id;type:varchar(50);not null;default:'';comment:行业ID"`
	ContactInformation string    `gorm:"column:contact_information;type:text;not null;comment:联系人信息"`
	CompanyType        uint8     `gorm:"column:company_type;type:tinyint(3) UNSIGNED;not null;default:1;comment:商家类型，1：服务商，2：品牌商，3：团长"`
	QianchuanUse       uint8     `gorm:"column:qianchuan_use;type:tinyint(3) UNSIGNED;not null;default:1;comment:千川使用情况，1：未使用，2：0<消耗≤10万/月，3：10万/月<消耗≤30万/月，4：30万/月<消耗≤100万/月，5：消耗>100万/月"`
	Sale               uint64    `gorm:"column:sale;type:bigint(20) UNSIGNED;not null;default:0;comment:参考销量"`
	Seller             string    `gorm:"column:seller;type:varchar(20);not null;default:'';comment:销售"`
	Facilitator        string    `gorm:"column:facilitator;type:varchar(20);not null;default:'';comment:辅导"`
	Source             string    `gorm:"column:source;type:varchar(20);not null;comment:信息来源"`
	Status             uint8     `gorm:"column:status;type:tinyint(3) UNSIGNED;not null;default:1;comment:状态，1：公海，2：洽谈，3：过期，4：正式，5：测试"`
	OperationLog       string    `gorm:"column:operation_log;type:text;not null;comment:操作日志"`
	AreaCode           uint64    `gorm:"column:area_code;type:bigint(20) UNSIGNED;not null;default:0;comment:区划代码"`
	Address            string    `gorm:"column:address;type:varchar(250);not null;comment:企业地址"`
	IsDel              uint8     `gorm:"column:is_del;type:tinyint(3) UNSIGNED;not null;default:0;comment:是否删除，1：已删除，0：未删除"`
	CreateTime         time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime         time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (Clue) TableName() string {
	return "company_clue"
}

type clueRepo struct {
	data *Data
	log  *log.Helper
}

func (c *Clue) ToDomain(ctx context.Context) *domain.Clue {
	clue := &domain.Clue{
		Id:                 c.Id,
		CompanyName:        c.CompanyName,
		IndustryId:         c.IndustryId,
		ContactInformation: c.ContactInformation,
		CompanyType:        c.CompanyType,
		QianchuanUse:       c.QianchuanUse,
		Sale:               c.Sale,
		Seller:             c.Seller,
		Facilitator:        c.Facilitator,
		Source:             c.Source,
		Status:             c.Status,
		OperationLog:       c.OperationLog,
		AreaCode:           c.AreaCode,
		Address:            c.Address,
		IsDel:              c.IsDel,
		CreateTime:         c.CreateTime,
		UpdateTime:         c.UpdateTime,
	}

	clue.GetOperationLog(ctx)
	clue.GetContactInformation(ctx)
	clue.GetStatusName(ctx)

	return clue
}

func NewClueRepo(data *Data, logger log.Logger) biz.ClueRepo {
	return &clueRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cr *clueRepo) GetById(ctx context.Context, id uint64) (*domain.Clue, error) {
	clue := &Clue{}

	if result := cr.data.DB(ctx).Where("is_del = ?", 0).First(clue, id); result.Error != nil {
		return nil, result.Error
	}

	return clue.ToDomain(ctx), nil
}

func (cr *clueRepo) List(ctx context.Context, pageNum, pageSize int, industryId uint64, keyword string, status uint8) ([]*domain.Clue, error) {
	var clues []Clue
	list := make([]*domain.Clue, 0)

	db := cr.data.db.WithContext(ctx).Where("is_del = ?", 0)

	if industryId > 0 {
		db = db.Where("find_in_set(?, `industry_id`)", industryId)
	}

	if status > 0 {
		db = db.Where("status = ?", status)
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		db = db.Where("(company_name like ? or seller like ? or facilitator like ?)", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if result := db.Order("create_time DESC,id DESC").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&clues); result.Error != nil {
		return nil, result.Error
	}

	for _, clue := range clues {
		list = append(list, clue.ToDomain(ctx))
	}

	return list, nil
}

func (cr *clueRepo) Count(ctx context.Context, industryId uint64, keyword string, status uint8) (int64, error) {
	var count int64

	db := cr.data.db.WithContext(ctx).Where("is_del = ?", 0)

	if industryId > 0 {
		db = db.Where("find_in_set(?, `industry_id`)", industryId)
	}

	if status > 0 {
		db = db.Where("status = ?", status)
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		db = db.Where("(company_name like ? or seller like ? or facilitator like ?)", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if result := db.Model(&Clue{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (cr *clueRepo) Statistics(ctx context.Context, status uint8) (int64, error) {
	var count int64

	if result := cr.data.db.WithContext(ctx).Model(&Clue{}).Where("is_del = ?", 0).Where("status = ?", status).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (cr *clueRepo) Save(ctx context.Context, in *domain.Clue) (*domain.Clue, error) {
	clue := &Clue{
		CompanyName:        in.CompanyName,
		IndustryId:         in.IndustryId,
		ContactInformation: in.ContactInformation,
		CompanyType:        in.CompanyType,
		QianchuanUse:       in.QianchuanUse,
		Sale:               in.Sale,
		Seller:             in.Seller,
		Facilitator:        in.Facilitator,
		Source:             in.Source,
		Status:             in.Status,
		OperationLog:       in.OperationLog,
		Address:            in.Address,
		AreaCode:           in.AreaCode,
		IsDel:              in.IsDel,
		CreateTime:         in.CreateTime,
		UpdateTime:         in.UpdateTime,
	}

	if result := cr.data.DB(ctx).Create(clue); result.Error != nil {
		return nil, result.Error
	}

	return clue.ToDomain(ctx), nil
}

func (cr *clueRepo) Update(ctx context.Context, in *domain.Clue) (*domain.Clue, error) {
	clue := &Clue{
		Id:                 in.Id,
		CompanyName:        in.CompanyName,
		IndustryId:         in.IndustryId,
		ContactInformation: in.ContactInformation,
		CompanyType:        in.CompanyType,
		QianchuanUse:       in.QianchuanUse,
		Sale:               in.Sale,
		Seller:             in.Seller,
		Facilitator:        in.Facilitator,
		Source:             in.Source,
		Status:             in.Status,
		OperationLog:       in.OperationLog,
		Address:            in.Address,
		AreaCode:           in.AreaCode,
		IsDel:              in.IsDel,
		CreateTime:         in.CreateTime,
		UpdateTime:         in.UpdateTime,
	}

	if result := cr.data.DB(ctx).Save(clue); result.Error != nil {
		return nil, result.Error
	}

	return clue.ToDomain(ctx), nil
}

func (cr *clueRepo) Delete(ctx context.Context, in *domain.Clue) error {
	clue := &Clue{
		Id:                 in.Id,
		CompanyName:        in.CompanyName,
		IndustryId:         in.IndustryId,
		ContactInformation: in.ContactInformation,
		CompanyType:        in.CompanyType,
		QianchuanUse:       in.QianchuanUse,
		Sale:               in.Sale,
		Seller:             in.Seller,
		Facilitator:        in.Facilitator,
		Source:             in.Source,
		Status:             in.Status,
		OperationLog:       in.OperationLog,
		Address:            in.Address,
		AreaCode:           in.AreaCode,
		IsDel:              in.IsDel,
		CreateTime:         in.CreateTime,
		UpdateTime:         in.UpdateTime,
	}

	if result := cr.data.DB(ctx).Delete(clue); result.Error != nil {
		return result.Error
	}

	return nil
}

package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"weixin/internal/biz"
	"weixin/internal/domain"
)

// 微信用户员工用工合同表
type UserContract struct {
	Id               uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	OrganizationId   uint64    `gorm:"column:organization_id;type:bigint(20) UNSIGNED;not null;uniqueIndex:contract_type_organization_id_identity_card_mark;comment:机构ID"`
	Name             string    `gorm:"column:name;type:text;not null;comment:姓名"`
	IdentityCard     string    `gorm:"column:identity_card;type:text;not null;comment:身份证号码"`
	IdentityCardMark string    `gorm:"column:identity_card_mark;type:varchar(250);not null;uniqueIndex:contract_type_organization_id_identity_card_mark;comment:身份证号码标记"`
	ServiceId        uint64    `gorm:"column:service_id;type:bigint(20) UNSIGNED;not null;comment:工猫服务主体ID"`
	TemplateId       uint64    `gorm:"column:template_id;type:bigint(20) UNSIGNED;not null;comment:工猫合同模板编号"`
	ContractId       uint64    `gorm:"column:contract_id;type:bigint(20) UNSIGNED;not null;comment:工猫签署合同ID"`
	ContractStatus   uint8     `gorm:"column:contract_status;type:tinyint(3) UNSIGNED;not null;default:0;comment:工猫签署合同状态，0：未签约（员工未同步），1：待签约（员工已同步，但未确认签约），2：签约文件生成中（用户已确认签署，合同文件正在生成，此状态可发起提现），3：签约完成 （用户已签约，签约合同生成完成）"`
	ContractType     uint8     `gorm:"column:contract_type;type:tinyint(3) UNSIGNED;not null;default:0;uniqueIndex:contract_type_organization_id_identity_card_mark;comment:工猫签署合同类型，1：成本购，2：电商，课程"`
	CreateTime       time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime       time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (UserContract) TableName() string {
	return "weixin_user_contract"
}

type userContractRepo struct {
	data *Data
	log  *log.Helper
}

func (uc *UserContract) ToDomain(ctx context.Context) *domain.UserContract {
	userContract := &domain.UserContract{
		Id:               uc.Id,
		OrganizationId:   uc.OrganizationId,
		Name:             uc.Name,
		IdentityCard:     uc.IdentityCard,
		IdentityCardMark: uc.IdentityCardMark,
		ServiceId:        uc.ServiceId,
		TemplateId:       uc.TemplateId,
		ContractId:       uc.ContractId,
		ContractStatus:   uc.ContractStatus,
		ContractType:     uc.ContractType,
		CreateTime:       uc.CreateTime,
		UpdateTime:       uc.UpdateTime,
	}

	return userContract
}

func NewUserContractRepo(data *Data, logger log.Logger) biz.UserContractRepo {
	return &userContractRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ucr *userContractRepo) Get(ctx context.Context, organizationId uint64, identityCardMark string) (*domain.UserContract, error) {
	userContract := &UserContract{}

	if result := ucr.data.db.WithContext(ctx).
		Where("organization_id = ?", organizationId).
		Where("identity_card_mark = ?", identityCardMark).
		First(userContract); result.Error != nil {
		return nil, result.Error
	}

	return userContract.ToDomain(ctx), nil
}

func (ucr *userContractRepo) GetByContractId(ctx context.Context, contractId uint64) (*domain.UserContract, error) {
	userContract := &UserContract{}

	if result := ucr.data.db.WithContext(ctx).
		Where("contract_id = ?", contractId).
		First(userContract); result.Error != nil {
		return nil, result.Error
	}

	return userContract.ToDomain(ctx), nil
}

func (ucr *userContractRepo) GetByIdentityCardMark(ctx context.Context, contractType uint8, identityCardMark string) (*domain.UserContract, error) {
	userContract := &UserContract{}

	db := ucr.data.db.WithContext(ctx).Where("identity_card_mark = ?", identityCardMark)

	if contractType > 0 {
		db = db.Where("contract_type = ?", contractType)
	}

	if result := db.First(userContract); result.Error != nil {
		return nil, result.Error
	}

	return userContract.ToDomain(ctx), nil
}

func (ucr *userContractRepo) Save(ctx context.Context, in *domain.UserContract) (*domain.UserContract, error) {
	userContract := &UserContract{
		OrganizationId:   in.OrganizationId,
		Name:             in.Name,
		IdentityCard:     in.IdentityCard,
		IdentityCardMark: in.IdentityCardMark,
		ServiceId:        in.ServiceId,
		TemplateId:       in.TemplateId,
		ContractId:       in.ContractId,
		ContractStatus:   in.ContractStatus,
		ContractType:     in.ContractType,
		CreateTime:       in.CreateTime,
		UpdateTime:       in.UpdateTime,
	}

	if result := ucr.data.DB(ctx).Create(userContract); result.Error != nil {
		return nil, result.Error
	}

	return userContract.ToDomain(ctx), nil
}

func (ucr *userContractRepo) Update(ctx context.Context, in *domain.UserContract) (*domain.UserContract, error) {
	userContract := &UserContract{
		Id:               in.Id,
		OrganizationId:   in.OrganizationId,
		Name:             in.Name,
		IdentityCard:     in.IdentityCard,
		IdentityCardMark: in.IdentityCardMark,
		ServiceId:        in.ServiceId,
		TemplateId:       in.TemplateId,
		ContractId:       in.ContractId,
		ContractStatus:   in.ContractStatus,
		ContractType:     in.ContractType,
		CreateTime:       in.CreateTime,
		UpdateTime:       in.UpdateTime,
	}

	if result := ucr.data.DB(ctx).Save(userContract); result.Error != nil {
		return nil, result.Error
	}

	return userContract.ToDomain(ctx), nil
}

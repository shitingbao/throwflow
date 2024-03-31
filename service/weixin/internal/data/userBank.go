package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"weixin/internal/biz"
	"weixin/internal/domain"
)

// 微信用户银行卡表
type UserBank struct {
	Id               uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	UserId           uint64    `gorm:"column:user_id;type:bigint(20) UNSIGNED;not null;uniqueIndex:user_id_identity_card_mark_bank_code;comment:微信小程序用户ID"`
	IdentityCardMark string    `gorm:"column:identity_card_mark;type:varchar(250);not null;uniqueIndex:user_id_identity_card_mark_bank_code;comment:身份证号码标记"`
	BankCode         string    `gorm:"column:bank_code;type:varchar(250);not null;uniqueIndex:user_id_identity_card_mark_bank_code;comment:银行卡卡号"`
	BankName         string    `gorm:"column:bank_name;type:varchar(250);not null;comment:银行名称"`
	CreateTime       time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime       time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (UserBank) TableName() string {
	return "weixin_user_bank"
}

type userBankRepo struct {
	data *Data
	log  *log.Helper
}

func (ub *UserBank) ToDomain(ctx context.Context) *domain.UserBank {
	userBank := &domain.UserBank{
		Id:               ub.Id,
		UserId:           ub.UserId,
		IdentityCardMark: ub.IdentityCardMark,
		BankCode:         ub.BankCode,
		BankName:         ub.BankName,
		CreateTime:       ub.CreateTime,
		UpdateTime:       ub.UpdateTime,
	}

	return userBank
}

func NewUserBankRepo(data *Data, logger log.Logger) biz.UserBankRepo {
	return &userBankRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ubr *userBankRepo) GetByBankCode(ctx context.Context, userId uint64, identityCardMark, bankCode string) (*domain.UserBank, error) {
	userBank := &UserBank{}

	if result := ubr.data.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Where("identity_card_mark = ?", identityCardMark).
		Where("bank_code = ?", bankCode).
		First(userBank); result.Error != nil {
		return nil, result.Error
	}

	return userBank.ToDomain(ctx), nil
}

func (ubr *userBankRepo) List(ctx context.Context, pageNum, pageSize int, userId uint64, identityCardMark string) ([]*domain.UserBank, error) {
	var userBanks []UserBank
	list := make([]*domain.UserBank, 0)

	if result := ubr.data.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Where("identity_card_mark = ?", identityCardMark).
		Order("create_time DESC,id DESC").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&userBanks); result.Error != nil {
		return nil, result.Error
	}

	for _, userBank := range userBanks {
		list = append(list, userBank.ToDomain(ctx))
	}

	return list, nil
}

func (ubr *userBankRepo) Count(ctx context.Context, userId uint64, identityCardMark string) (int64, error) {
	var count int64

	if result := ubr.data.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Where("identity_card_mark = ?", identityCardMark).
		Model(&UserBank{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (ubr *userBankRepo) Save(ctx context.Context, in *domain.UserBank) (*domain.UserBank, error) {
	userBank := &UserBank{
		UserId:           in.UserId,
		IdentityCardMark: in.IdentityCardMark,
		BankCode:         in.BankCode,
		BankName:         in.BankName,
		CreateTime:       in.CreateTime,
		UpdateTime:       in.UpdateTime,
	}

	if result := ubr.data.DB(ctx).Create(userBank); result.Error != nil {
		return nil, result.Error
	}

	return userBank.ToDomain(ctx), nil
}

func (ubr *userBankRepo) Update(ctx context.Context, in *domain.UserBank) (*domain.UserBank, error) {
	userBank := &UserBank{
		Id:               in.Id,
		UserId:           in.UserId,
		IdentityCardMark: in.IdentityCardMark,
		BankCode:         in.BankCode,
		BankName:         in.BankName,
		CreateTime:       in.CreateTime,
		UpdateTime:       in.UpdateTime,
	}

	if result := ubr.data.DB(ctx).Save(userBank); result.Error != nil {
		return nil, result.Error
	}

	return userBank.ToDomain(ctx), nil
}

func (ubr *userBankRepo) Delete(ctx context.Context, in *domain.UserBank) error {
	userBank := &UserBank{
		Id:               in.Id,
		UserId:           in.UserId,
		IdentityCardMark: in.IdentityCardMark,
		BankCode:         in.BankCode,
		BankName:         in.BankName,
		CreateTime:       in.CreateTime,
		UpdateTime:       in.UpdateTime,
	}

	if result := ubr.data.db.WithContext(ctx).Delete(userBank); result.Error != nil {
		return result.Error
	}

	return nil
}

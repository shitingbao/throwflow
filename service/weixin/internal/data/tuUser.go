package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"weixin/internal/biz"
	"weixin/internal/domain"
)

// 微信小程序用户表
type TuUser struct {
	Phone    string  `gorm:"column:phone;type:varchar(50);not null;comment:手机号"`
	Num      uint64  `gorm:"column:num;type:bigint(20) UNSIGNED;not null;default:0;comment:券码数"`
	Id       string  `gorm:"column:id;type:varchar(50);not null;comment:标签ID"`
	ParentId string  `gorm:"column:parent_id;type:varchar(50);not null;comment:父级标签ID"`
	Level    string  `gorm:"column:level;type:varchar(50);not null;comment:级别"`
	Balance  float64 `gorm:"column:balance;type:decimal(10, 2) UNSIGNED;not null;comment:账户余额，单位元"`
}

func (TuUser) TableName() string {
	return "tu_user"
}

type tuUserRepo struct {
	data *Data
	log  *log.Helper
}

func (tu *TuUser) ToDomain(ctx context.Context) *domain.TuUser {
	tuUser := &domain.TuUser{
		Phone:    tu.Phone,
		Num:      tu.Num,
		Id:       tu.Id,
		ParentId: tu.ParentId,
		Level:    tu.Level,
		Balance:  tu.Balance,
	}

	return tuUser
}

func NewTuUserRepo(data *Data, logger log.Logger) biz.TuUserRepo {
	return &tuUserRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (tur *tuUserRepo) List(ctx context.Context) ([]*domain.TuUser, error) {
	var tuUsers []TuUser
	list := make([]*domain.TuUser, 0)

	if result := tur.data.db.WithContext(ctx).
		Find(&tuUsers); result.Error != nil {
		return nil, result.Error
	}

	for _, tuUser := range tuUsers {
		list = append(list, tuUser.ToDomain(ctx))
	}

	return list, nil
}

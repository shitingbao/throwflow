package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"weixin/internal/biz"
	"weixin/internal/domain"
)

// 微信小程序用户券码新增日志表
type UserCouponCreateLog struct {
	Id             uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	UserId         uint64    `gorm:"column:user_id;type:bigint(20) UNSIGNED;not null;comment:微信小程序用户ID"`
	OrganizationId uint64    `gorm:"column:organization_id;type:bigint(20) UNSIGNED;not null;comment:机构ID"`
	Level          uint8     `gorm:"column:level;type:tinyint(3) UNSIGNED;not null;default:0;comment:等级"`
	Num            uint32    `gorm:"column:num;type:int(10) UNSIGNED;not null;default:1;comment:劵码个数"`
	IsHandle       uint8     `gorm:"column:is_handle;type:tinyint(3) UNSIGNED;not null;default:0;comment:是否处理，1：正在处理，0：未处理"`
	CreateTime     time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime     time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (UserCouponCreateLog) TableName() string {
	return "weixin_user_coupon_create_log"
}

type userCouponCreateLogRepo struct {
	data *Data
	log  *log.Helper
}

func (uccl *UserCouponCreateLog) ToDomain() *domain.UserCouponCreateLog {
	return &domain.UserCouponCreateLog{
		Id:             uccl.Id,
		UserId:         uccl.UserId,
		OrganizationId: uccl.OrganizationId,
		Level:          uccl.Level,
		Num:            uccl.Num,
		IsHandle:       uccl.IsHandle,
		CreateTime:     uccl.CreateTime,
		UpdateTime:     uccl.UpdateTime,
	}
}

func NewUserCouponCreateLogRepo(data *Data, logger log.Logger) biz.UserCouponCreateLogRepo {
	return &userCouponCreateLogRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ucclr *userCouponCreateLogRepo) Get(ctx context.Context, isHandle string) (*domain.UserCouponCreateLog, error) {
	userCouponCreateLog := &UserCouponCreateLog{}

	if result := ucclr.data.db.WithContext(ctx).Where("is_handle = ?", isHandle).Order("create_time asc").First(userCouponCreateLog); result.Error != nil {
		return nil, result.Error
	}

	return userCouponCreateLog.ToDomain(), nil
}

func (ucclr *userCouponCreateLogRepo) List(ctx context.Context, isHandle string) ([]*domain.UserCouponCreateLog, error) {
	var userCouponCreateLogs []UserCouponCreateLog
	list := make([]*domain.UserCouponCreateLog, 0)

	if result := ucclr.data.db.WithContext(ctx).Where("is_handle = ?", isHandle).
		Find(&userCouponCreateLogs); result.Error != nil {
		return nil, result.Error
	}

	for _, userCouponCreateLog := range userCouponCreateLogs {
		list = append(list, userCouponCreateLog.ToDomain())
	}

	return list, nil
}

func (ucclr *userCouponCreateLogRepo) Save(ctx context.Context, in *domain.UserCouponCreateLog) (*domain.UserCouponCreateLog, error) {
	userCouponCreateLog := &UserCouponCreateLog{
		UserId:         in.UserId,
		OrganizationId: in.OrganizationId,
		Level:          in.Level,
		Num:            in.Num,
		IsHandle:       in.IsHandle,
		CreateTime:     in.CreateTime,
		UpdateTime:     in.UpdateTime,
	}

	if result := ucclr.data.DB(ctx).Create(userCouponCreateLog); result.Error != nil {
		return nil, result.Error
	}

	return userCouponCreateLog.ToDomain(), nil
}

func (ucclr *userCouponCreateLogRepo) Update(ctx context.Context, in *domain.UserCouponCreateLog) (*domain.UserCouponCreateLog, error) {
	userCouponCreateLog := &UserCouponCreateLog{
		Id:             in.Id,
		UserId:         in.UserId,
		OrganizationId: in.OrganizationId,
		Level:          in.Level,
		Num:            in.Num,
		IsHandle:       in.IsHandle,
		CreateTime:     in.CreateTime,
		UpdateTime:     in.UpdateTime,
	}

	if result := ucclr.data.db.WithContext(ctx).Save(userCouponCreateLog); result.Error != nil {
		return nil, result.Error
	}

	return userCouponCreateLog.ToDomain(), nil
}

func (ucclr *userCouponCreateLogRepo) Delete(ctx context.Context, in *domain.UserCouponCreateLog) error {
	userCouponCreateLog := &UserCouponCreateLog{
		Id:             in.Id,
		UserId:         in.UserId,
		OrganizationId: in.OrganizationId,
		Level:          in.Level,
		Num:            in.Num,
		IsHandle:       in.IsHandle,
		CreateTime:     in.CreateTime,
		UpdateTime:     in.UpdateTime,
	}

	if result := ucclr.data.DB(ctx).Delete(userCouponCreateLog); result.Error != nil {
		return result.Error
	}

	return nil
}

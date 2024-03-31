package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"weixin/internal/biz"
	"weixin/internal/domain"
)

// 微信小程序用户扫码记录
type UserScanRecord struct {
	Id                 uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	UserId             uint64    `gorm:"column:user_id;type:bigint(20) UNSIGNED;not null;comment:微信小程序用户ID"`
	OrganizationId     uint64    `gorm:"column:organization_id;type:bigint(20) UNSIGNED;not null;comment:机构ID"`
	OrganizationUserId uint64    `gorm:"column:organization_user_id;type:bigint(20) UNSIGNED;not null;default:0;comment:企业机构团长ID"`
	AwemeUserId        uint64    `gorm:"column:aweme_user_id;type:bigint(20) UNSIGNED;not null;default:0;comment:达人商务推广用户ID"`
	CreateTime         time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime         time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (UserScanRecord) TableName() string {
	return "weixin_user_scan_record"
}

type userScanRecordRepo struct {
	data *Data
	log  *log.Helper
}

func (usr *UserScanRecord) ToDomain(ctx context.Context) *domain.UserScanRecord {
	userScanRecord := &domain.UserScanRecord{
		Id:                 usr.Id,
		UserId:             usr.UserId,
		OrganizationId:     usr.OrganizationId,
		OrganizationUserId: usr.OrganizationUserId,
		AwemeUserId:        usr.AwemeUserId,
		CreateTime:         usr.CreateTime,
		UpdateTime:         usr.UpdateTime,
	}

	return userScanRecord
}

func NewUserScanRecordRepo(data *Data, logger log.Logger) biz.UserScanRecordRepo {
	return &userScanRecordRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (usrr *userScanRecordRepo) Get(ctx context.Context, userId, organizationId uint64, isOrganization uint8) (*domain.UserScanRecord, error) {
	userScanRecord := &UserScanRecord{}

	db := usrr.data.db.WithContext(ctx).
		Where("user_id = ?", userId)

	if isOrganization == 1 {
		if organizationId > 0 {
			db = db.Where("organization_id = ?", organizationId)
		} else {
			db = db.Where("organization_id > ?", 0)
		}

		db = db.Where("aweme_user_id = ?", 0)
	} else {
		db = db.Where("organization_id = ?", 0)
		db = db.Where("aweme_user_id > ?", 0)
	}

	if result := db.Order("create_time desc").First(userScanRecord); result.Error != nil {
		return nil, result.Error
	}

	return userScanRecord.ToDomain(ctx), nil
}

func (usrr *userScanRecordRepo) Save(ctx context.Context, in *domain.UserScanRecord) (*domain.UserScanRecord, error) {
	userScanRecord := &UserScanRecord{
		UserId:             in.UserId,
		OrganizationId:     in.OrganizationId,
		OrganizationUserId: in.OrganizationUserId,
		AwemeUserId:        in.AwemeUserId,
		CreateTime:         in.CreateTime,
		UpdateTime:         in.UpdateTime,
	}

	if result := usrr.data.DB(ctx).Create(userScanRecord); result.Error != nil {
		return nil, result.Error
	}

	return userScanRecord.ToDomain(ctx), nil
}

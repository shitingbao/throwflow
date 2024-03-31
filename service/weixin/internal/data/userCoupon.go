package data

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	ctos "github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"io"
	"time"
	"weixin/internal/biz"
	"weixin/internal/domain"
)

// 微信小程序用户券码
type UserCoupon struct {
	Id               uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	UserId           uint64    `gorm:"column:user_id;type:bigint(20) UNSIGNED;not null;comment:微信小程序用户ID"`
	CouponCode       string    `gorm:"column:coupon_code;type:varchar(10);not null;comment:机构用户简码"`
	Level            uint8     `gorm:"column:level;type:tinyint(3) UNSIGNED;not null;default:0;comment:等级"`
	Phone            string    `gorm:"column:phone;type:varchar(20);not null;uniqueIndex:phone;comment:手机号"`
	UserCouponStatus uint8     `gorm:"column:user_coupon_status;type:tinyint(3) UNSIGNED;not null;default:0;comment:1：未激活，2：已指定，3：已激活"`
	OrganizationId   uint64    `gorm:"column:organization_id;type:bigint(20) UNSIGNED;not null;comment:机构ID"`
	CreateTime       time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime       time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (UserCoupon) TableName() string {
	return "weixin_user_coupon"
}

type userCouponRepo struct {
	data *Data
	log  *log.Helper
}

func (uc *UserCoupon) ToDomain(ctx context.Context) *domain.UserCoupon {
	userCoupon := &domain.UserCoupon{
		Id:               uc.Id,
		UserId:           uc.UserId,
		CouponCode:       uc.CouponCode,
		Level:            uc.Level,
		Phone:            uc.Phone,
		UserCouponStatus: uc.UserCouponStatus,
		OrganizationId:   uc.OrganizationId,
		CreateTime:       uc.CreateTime,
		UpdateTime:       uc.UpdateTime,
	}

	return userCoupon
}

func NewUserCouponRepo(data *Data, logger log.Logger) biz.UserCouponRepo {
	return &userCouponRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ucr *userCouponRepo) Get(ctx context.Context, userId uint64, level, userCouponStatus string) (*domain.UserCoupon, error) {
	userCoupon := &UserCoupon{}

	db := ucr.data.db.WithContext(ctx).
		Where("user_id = ?", userId)

	if len(level) > 0 {
		db = db.Where("level = ?", level)
	}

	if len(userCouponStatus) > 0 {
		db = db.Where("user_coupon_status = ?", userCouponStatus)
	}

	if result := db.First(userCoupon); result.Error != nil {
		return nil, result.Error
	}

	return userCoupon.ToDomain(ctx), nil
}

func (ucr *userCouponRepo) GetByPhone(ctx context.Context, organizationId uint64, phone, level, userCouponStatus string) (*domain.UserCoupon, error) {
	userCoupon := &UserCoupon{}

	db := ucr.data.db.WithContext(ctx).
		Where("phone = ?", phone)

	if organizationId > 0 {
		db = db.Where("organization_id = ?", organizationId)
	}

	if len(level) > 0 {
		db = db.Where("level = ?", level)
	}

	if len(userCouponStatus) > 0 {
		db = db.Where("user_coupon_status = ?", userCouponStatus)
	}

	if result := db.Order("update_time ASC").First(userCoupon); result.Error != nil {
		return nil, result.Error
	}

	return userCoupon.ToDomain(ctx), nil
}

func (ucr *userCouponRepo) List(ctx context.Context, pageNum, pageSize int, userId, organizationId uint64, userCouponStatus string) ([]*domain.UserCoupon, error) {
	var userCoupons []UserCoupon
	list := make([]*domain.UserCoupon, 0)

	db := ucr.data.db.WithContext(ctx)

	if userId > 0 {
		db = db.Where("user_id = ?", userId)
	}

	if organizationId > 0 {
		db = db.Where("organization_id = ?", organizationId)
	}

	if len(userCouponStatus) > 0 {
		db = db.Where("user_coupon_status = ?", userCouponStatus)
	}

	if pageNum == 0 {
		if result := db.Order("user_coupon_status ASC,id DESC").
			Find(&userCoupons); result.Error != nil {
			return nil, result.Error
		}
	} else {
		if result := db.Order("user_coupon_status ASC,id DESC").
			Limit(pageSize).Offset((pageNum - 1) * pageSize).
			Find(&userCoupons); result.Error != nil {
			return nil, result.Error
		}
	}

	for _, userCoupon := range userCoupons {
		list = append(list, userCoupon.ToDomain(ctx))
	}

	return list, nil
}

func (ucr *userCouponRepo) Count(ctx context.Context, userId, organizationId uint64, userCouponStatus string) (int64, error) {
	var count int64

	db := ucr.data.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Where("organization_id = ?", organizationId)

	if len(userCouponStatus) > 0 {
		db = db.Where("user_coupon_status = ?", userCouponStatus)
	}

	if result := db.Model(&UserCoupon{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (ucr *userCouponRepo) Save(ctx context.Context, in *domain.UserCoupon) (*domain.UserCoupon, error) {
	userCoupon := &UserCoupon{
		UserId:           in.UserId,
		CouponCode:       in.CouponCode,
		Level:            in.Level,
		Phone:            in.Phone,
		UserCouponStatus: in.UserCouponStatus,
		OrganizationId:   in.OrganizationId,
		CreateTime:       in.CreateTime,
		UpdateTime:       in.UpdateTime,
	}

	if result := ucr.data.DB(ctx).Create(userCoupon); result.Error != nil {
		return nil, result.Error
	}

	return userCoupon.ToDomain(ctx), nil
}

func (ucr *userCouponRepo) Update(ctx context.Context, in *domain.UserCoupon) (*domain.UserCoupon, error) {
	userCoupon := &UserCoupon{
		Id:               in.Id,
		UserId:           in.UserId,
		CouponCode:       in.CouponCode,
		Level:            in.Level,
		Phone:            in.Phone,
		UserCouponStatus: in.UserCouponStatus,
		OrganizationId:   in.OrganizationId,
		CreateTime:       in.CreateTime,
		UpdateTime:       in.UpdateTime,
	}

	if result := ucr.data.DB(ctx).Save(userCoupon); result.Error != nil {
		return nil, result.Error
	}

	return userCoupon.ToDomain(ctx), nil
}

func (ucr *userCouponRepo) SaveCacheString(ctx context.Context, key string, val string, timeout time.Duration) (bool, error) {
	result, err := ucr.data.rdb.SetNX(ctx, key, val, timeout).Result()

	if err != nil {
		return false, err
	}

	return result, nil
}

func (ucr *userCouponRepo) DeleteCache(ctx context.Context, key string) error {
	if _, err := ucr.data.rdb.Del(ctx, key).Result(); err != nil {
		return err
	}

	return nil
}

func (ucr *userCouponRepo) PutContent(ctx context.Context, fileName string, content io.Reader) (*ctos.PutObjectV2Output, error) {
	for _, ltos := range ucr.data.toses {
		if ltos.name == "company" {
			output, err := ltos.tos.PutContent(ctx, fileName, content)

			if err != nil {
				return nil, err
			}

			return output, nil
		}
	}

	return nil, errors.New("tos is not exist")
}

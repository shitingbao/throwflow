package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"weixin/internal/biz"
	"weixin/internal/domain"
)

// 微信小程序用户机构关系表
type UserOrganizationRelation struct {
	Id                        uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	UserId                    uint64    `gorm:"column:user_id;type:bigint(20) UNSIGNED;not null;uniqueIndex:user_id_organization_id_is_order_relation;comment:微信小程序用户ID"`
	OrganizationId            uint64    `gorm:"column:organization_id;type:bigint(20) UNSIGNED;not null;uniqueIndex:user_id_organization_id_is_order_relation;comment:机构ID"`
	OrganizationUserId        uint64    `gorm:"column:organization_user_id;type:bigint(20) UNSIGNED;not null;default:0;comment:企业机构团长ID"`
	OrganizationTutorId       uint64    `gorm:"column:organization_tutor_id;type:bigint(20) UNSIGNED;not null;default:0;comment:企业机构导师ID"`
	OrganizationUserQrCodeUrl string    `gorm:"column:organization_user_qr_code_url;type:varchar(250);not null;comment:机构用户小程序码 URL"`
	Level                     uint8     `gorm:"column:level;type:tinyint(3) UNSIGNED;not null;default:0;comment:等级"`
	IsOrderRelation           uint8     `gorm:"column:is_order_relation;type:tinyint(3) UNSIGNED;not null;default:0;uniqueIndex:user_id_organization_id_is_order_relation;comment:状态：，1：账单上下级关系"`
	CreateTime                time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime                time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (UserOrganizationRelation) TableName() string {
	return "weixin_user_organization_relation"
}

type userOrganizationRelationRepo struct {
	data *Data
	log  *log.Helper
}

func (uo *UserOrganizationRelation) ToDomain(ctx context.Context) *domain.UserOrganizationRelation {
	userOrganizationRelation := &domain.UserOrganizationRelation{
		Id:                        uo.Id,
		UserId:                    uo.UserId,
		OrganizationId:            uo.OrganizationId,
		OrganizationUserId:        uo.OrganizationUserId,
		OrganizationTutorId:       uo.OrganizationTutorId,
		OrganizationUserQrCodeUrl: uo.OrganizationUserQrCodeUrl,
		Level:                     uo.Level,
		IsOrderRelation:           uo.IsOrderRelation,
		CreateTime:                uo.CreateTime,
		UpdateTime:                uo.UpdateTime,
	}

	return userOrganizationRelation
}

func NewUserOrganizationRelationRepo(data *Data, logger log.Logger) biz.UserOrganizationRelationRepo {
	return &userOrganizationRelationRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (uorr *userOrganizationRelationRepo) Get(ctx context.Context, organizationId, organizationUserId uint64, isOrderRelation string) (*domain.UserOrganizationRelation, error) {
	userOrganizationRelation := &UserOrganizationRelation{}

	db := uorr.data.db.WithContext(ctx).
		Where("organization_id = ?", organizationId)

	if organizationUserId > 0 {
		db = db.Where("organization_user_id = ?", organizationUserId)
	}

	if len(isOrderRelation) > 0 {
		db = db.Where("is_order_relation = ?", isOrderRelation)
	}

	if result := db.First(userOrganizationRelation); result.Error != nil {
		return nil, result.Error
	}

	return userOrganizationRelation.ToDomain(ctx), nil
}

func (uorr *userOrganizationRelationRepo) GetByUserId(ctx context.Context, userId, organizationId, organizationUserId uint64, isOrderRelation string) (*domain.UserOrganizationRelation, error) {
	userOrganizationRelation := &UserOrganizationRelation{}

	db := uorr.data.db.WithContext(ctx).Where("user_id = ?", userId)

	if organizationId > 0 {
		db = db.Where("organization_id = ?", organizationId)
	}

	if organizationUserId > 0 {
		db = db.Where("organization_user_id = ?", organizationUserId)
	}

	if len(isOrderRelation) > 0 {
		db = db.Where("is_order_relation = ?", isOrderRelation)
	}

	if result := db.First(userOrganizationRelation); result.Error != nil {
		return nil, result.Error
	}

	return userOrganizationRelation.ToDomain(ctx), nil
}

func (uorr *userOrganizationRelationRepo) List(ctx context.Context, organizationId uint64) ([]*domain.UserOrganizationRelation, error) {
	var userOrganizationRelations []UserOrganizationRelation
	list := make([]*domain.UserOrganizationRelation, 0)

	db := uorr.data.db.WithContext(ctx)

	if organizationId > 0 {
		db = db.Where("organization_id = ?", organizationId)
	}

	if result := db.Find(&userOrganizationRelations); result.Error != nil {
		return nil, result.Error
	}

	for _, userOrganizationRelation := range userOrganizationRelations {
		list = append(list, userOrganizationRelation.ToDomain(ctx))
	}

	return list, nil
}

func (uorr *userOrganizationRelationRepo) ListDirectChild(ctx context.Context, userId, organizationId uint64) ([]*domain.UserOrganizationRelation, error) {
	var userOrganizationRelations []UserOrganizationRelation
	list := make([]*domain.UserOrganizationRelation, 0)

	db := uorr.data.db.WithContext(ctx).
		Where("organization_user_id = ?", userId).
		Where("organization_id = ?", organizationId)

	if result := db.Find(&userOrganizationRelations); result.Error != nil {
		return nil, result.Error
	}

	for _, userOrganizationRelation := range userOrganizationRelations {
		list = append(list, userOrganizationRelation.ToDomain(ctx))
	}

	return list, nil
}

func (uorr *userOrganizationRelationRepo) Count(ctx context.Context, userId, organizationId uint64, day string) (int64, error) {
	var count int64

	db := uorr.data.db.WithContext(ctx).
		Where("organization_user_id = ?", userId).
		Where("organization_id = ?", organizationId)

	if len(day) > 0 {
		db = db.Where("create_time >= ?", day+" 00:00:00")
		db = db.Where("create_time <= ?", day+" 23:59:59")
	}

	if result := db.Model(&UserOrganizationRelation{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (uorr *userOrganizationRelationRepo) Save(ctx context.Context, in *domain.UserOrganizationRelation) (*domain.UserOrganizationRelation, error) {
	userOrganizationRelation := &UserOrganizationRelation{
		UserId:                    in.UserId,
		OrganizationId:            in.OrganizationId,
		OrganizationUserId:        in.OrganizationUserId,
		OrganizationTutorId:       in.OrganizationTutorId,
		OrganizationUserQrCodeUrl: in.OrganizationUserQrCodeUrl,
		Level:                     in.Level,
		IsOrderRelation:           in.IsOrderRelation,
		CreateTime:                in.CreateTime,
		UpdateTime:                in.UpdateTime,
	}

	if result := uorr.data.DB(ctx).Create(userOrganizationRelation); result.Error != nil {
		return nil, result.Error
	}

	return userOrganizationRelation.ToDomain(ctx), nil
}

func (uorr *userOrganizationRelationRepo) Update(ctx context.Context, in *domain.UserOrganizationRelation) (*domain.UserOrganizationRelation, error) {
	userOrganizationRelation := &UserOrganizationRelation{
		Id:                        in.Id,
		UserId:                    in.UserId,
		OrganizationId:            in.OrganizationId,
		OrganizationUserId:        in.OrganizationUserId,
		OrganizationTutorId:       in.OrganizationTutorId,
		OrganizationUserQrCodeUrl: in.OrganizationUserQrCodeUrl,
		Level:                     in.Level,
		IsOrderRelation:           in.IsOrderRelation,
		CreateTime:                in.CreateTime,
		UpdateTime:                in.UpdateTime,
	}

	if result := uorr.data.DB(ctx).Save(userOrganizationRelation); result.Error != nil {
		return nil, result.Error
	}

	return userOrganizationRelation.ToDomain(ctx), nil
}

func (uorr *userOrganizationRelationRepo) UpdateSuperior(ctx context.Context, organizationTutorId uint64, childIds []uint64) error {
	if result := uorr.data.DB(ctx).Model(UserOrganizationRelation{}).
		Where("user_id in ?", childIds).
		Updates(map[string]interface{}{"organization_tutor_id": organizationTutorId}); result.Error != nil {
		return result.Error
	}

	return nil
}

func (uorr *userOrganizationRelationRepo) DeleteByUserId(ctx context.Context, userId uint64, isOrderRelation string) error {
	if result := uorr.data.DB(ctx).
		Where("user_id = ?", userId).
		Where("is_order_relation = ?", isOrderRelation).
		Delete(&UserOrganizationRelation{}); result.Error != nil {
		return result.Error
	}

	return nil
}

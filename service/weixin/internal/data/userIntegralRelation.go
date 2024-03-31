package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"weixin/internal/biz"
	"weixin/internal/domain"
)

// 微信小程序用户积分关系表
type UserIntegralRelation struct {
	Id                 uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	UserId             uint64    `gorm:"column:user_id;type:bigint(20) UNSIGNED;not null;uniqueIndex:user_id_organization_id;comment:微信小程序用户ID"`
	OrganizationId     uint64    `gorm:"column:organization_id;type:bigint(20) UNSIGNED;not null;uniqueIndex:user_id_organization_id;comment:机构ID"`
	OrganizationUserId uint64    `gorm:"column:organization_user_id;type:bigint(20) UNSIGNED;not null;default:0;comment:企业机构团长ID"`
	Level              uint8     `gorm:"column:level;type:tinyint(3) UNSIGNED;not null;default:0;comment:等级"`
	CreateTime         time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime         time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (UserIntegralRelation) TableName() string {
	return "weixin_user_integral_relation"
}

type userIntegralRelationRepo struct {
	data *Data
	log  *log.Helper
}

func (uir *UserIntegralRelation) ToDomain(ctx context.Context) *domain.UserIntegralRelation {
	userIntegralRelation := &domain.UserIntegralRelation{
		Id:                 uir.Id,
		UserId:             uir.UserId,
		OrganizationId:     uir.OrganizationId,
		OrganizationUserId: uir.OrganizationUserId,
		Level:              uir.Level,
		CreateTime:         uir.CreateTime,
		UpdateTime:         uir.UpdateTime,
	}

	return userIntegralRelation
}

func NewUserIntegralRelationRepo(data *Data, logger log.Logger) biz.UserIntegralRelationRepo {
	return &userIntegralRelationRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (uirr *userIntegralRelationRepo) Get(ctx context.Context, organizationId, organizationUserId uint64) (*domain.UserIntegralRelation, error) {
	userIntegralRelation := &UserIntegralRelation{}

	db := uirr.data.db.WithContext(ctx).
		Where("organization_id = ?", organizationId)

	if organizationUserId > 0 {
		db = db.Where("organization_user_id = ?", organizationUserId)
	}

	if result := db.First(userIntegralRelation); result.Error != nil {
		return nil, result.Error
	}

	return userIntegralRelation.ToDomain(ctx), nil
}

func (uirr *userIntegralRelationRepo) GetByUserId(ctx context.Context, userId, organizationId, organizationUserId uint64) (*domain.UserIntegralRelation, error) {
	userIntegralRelation := &UserIntegralRelation{}

	db := uirr.data.db.WithContext(ctx).Where("user_id = ?", userId)

	if organizationId > 0 {
		db = db.Where("organization_id = ?", organizationId)
	}

	if organizationUserId > 0 {
		db = db.Where("organization_user_id = ?", organizationUserId)
	}

	if result := db.First(userIntegralRelation); result.Error != nil {
		return nil, result.Error
	}

	return userIntegralRelation.ToDomain(ctx), nil
}

func (uirr *userIntegralRelationRepo) GetSuperior(ctx context.Context, userId uint64, level uint8, userIntegralRelations []*domain.UserIntegralRelation) *domain.UserIntegralRelation {
	for _, userIntegralRelation := range userIntegralRelations {
		if userIntegralRelation.UserId == userId {
			if userIntegralRelation.Level == level {
				return userIntegralRelation
			} else {
				if userIntegralRelation.OrganizationUserId > 0 {
					return uirr.GetSuperior(ctx, userIntegralRelation.OrganizationUserId, level, userIntegralRelations)
				} else {
					return nil
				}
			}
		}
	}

	return nil
}

func (uirr *userIntegralRelationRepo) GetChildNum(ctx context.Context, userId uint64, childNum *uint64, userIntegralRelations []*domain.UserIntegralRelation) {
	for _, userIntegralRelation := range userIntegralRelations {
		if userIntegralRelation.OrganizationUserId == userId {
			*childNum += 1

			uirr.GetChildNum(ctx, userIntegralRelation.UserId, childNum, userIntegralRelations)
		}
	}
}

func (uirr *userIntegralRelationRepo) List(ctx context.Context, organizationId uint64) ([]*domain.UserIntegralRelation, error) {
	var userIntegralRelations []UserIntegralRelation
	list := make([]*domain.UserIntegralRelation, 0)

	db := uirr.data.db.WithContext(ctx)

	if organizationId > 0 {
		db = db.Where("organization_id = ?", organizationId)
	}

	if result := db.Find(&userIntegralRelations); result.Error != nil {
		return nil, result.Error
	}

	for _, userIntegralRelation := range userIntegralRelations {
		list = append(list, userIntegralRelation.ToDomain(ctx))
	}

	return list, nil
}

func (uirr *userIntegralRelationRepo) ListChildId(ctx context.Context, userId uint64, childIds *[]uint64, userIntegralRelations []*domain.UserIntegralRelation) {
	for _, userIntegralRelation := range userIntegralRelations {
		if userIntegralRelation.OrganizationUserId == userId {
			*childIds = append(*childIds, userIntegralRelation.UserId)

			uirr.ListChildId(ctx, userIntegralRelation.UserId, childIds, userIntegralRelations)
		}
	}
}

func (uirr *userIntegralRelationRepo) Save(ctx context.Context, in *domain.UserIntegralRelation) (*domain.UserIntegralRelation, error) {
	userIntegralRelation := &UserIntegralRelation{
		UserId:             in.UserId,
		OrganizationId:     in.OrganizationId,
		OrganizationUserId: in.OrganizationUserId,
		Level:              in.Level,
		CreateTime:         in.CreateTime,
		UpdateTime:         in.UpdateTime,
	}

	if result := uirr.data.DB(ctx).Create(userIntegralRelation); result.Error != nil {
		return nil, result.Error
	}

	return userIntegralRelation.ToDomain(ctx), nil
}

func (uirr *userIntegralRelationRepo) Update(ctx context.Context, in *domain.UserIntegralRelation) (*domain.UserIntegralRelation, error) {
	userIntegralRelation := &UserIntegralRelation{
		Id:                 in.Id,
		UserId:             in.UserId,
		OrganizationId:     in.OrganizationId,
		OrganizationUserId: in.OrganizationUserId,
		Level:              in.Level,
		CreateTime:         in.CreateTime,
		UpdateTime:         in.UpdateTime,
	}

	if result := uirr.data.DB(ctx).Save(userIntegralRelation); result.Error != nil {
		return nil, result.Error
	}

	return userIntegralRelation.ToDomain(ctx), nil
}

package data

import (
	"company/internal/biz"
	"company/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// 企业库用户权限表
type CompanyUserRole struct {
	UserId       uint64    `gorm:"column:user_id;type:bigint(20) UNSIGNED;not null;comment:用户ID"`
	AdvertiserId uint64    `gorm:"column:advertiser_id;type:bigint(20) UNSIGNED;not null;comment:广告主ID"`
	CompanyId    uint64    `gorm:"column:company_id;type:bigint(20) UNSIGNED;not null;comment:企业库ID"`
	RoleType     uint8     `gorm:"column:role_type;type:tinyint(3) UNSIGNED;not null;default:1;comment:企业库用户权限类型，1：增加，2：减少"`
	Day          uint32    `gorm:"column:day;type:int(11);not null;comment:修改日期"`
	CreateTime   time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime   time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (CompanyUserRole) TableName() string {
	return "company_company_user_role"
}

type companyUserRoleRepo struct {
	data *Data
	log  *log.Helper
}

func (cur *CompanyUserRole) ToDomain() *domain.CompanyUserRole {
	companyUserRole := &domain.CompanyUserRole{
		UserId:       cur.UserId,
		AdvertiserId: cur.AdvertiserId,
		CompanyId:    cur.CompanyId,
		RoleType:     cur.RoleType,
		Day:          cur.Day,
		CreateTime:   cur.CreateTime,
		UpdateTime:   cur.UpdateTime,
	}

	return companyUserRole
}

func NewCompanyUserRoleRepo(data *Data, logger log.Logger) biz.CompanyUserRoleRepo {
	return &companyUserRoleRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (curr *companyUserRoleRepo) GetByUserIdAndCompanyIdAndAdvertiserIdAndDay(ctx context.Context, userId, companyId, advertiserId uint64, day uint32) (*domain.CompanyUserRole, error) {
	companyUserRole := &CompanyUserRole{}

	if result := curr.data.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Where("company_id = ?", companyId).
		Where("advertiser_id = ?", advertiserId).
		Where("day = ?", day).
		First(companyUserRole); result.Error != nil {
		return nil, result.Error
	}

	return companyUserRole.ToDomain(), nil
}

func (curr *companyUserRoleRepo) ListByCompanyIdAndDay(ctx context.Context, companyId uint64, day uint32) ([]*domain.CompanyUserRole, error) {
	var companyUserRoles []CompanyUserRole
	list := make([]*domain.CompanyUserRole, 0)

	if result := curr.data.db.WithContext(ctx).
		Where("company_id = ?", companyId).
		Where("day <= ?", day).
		Order("day ASC").
		Find(&companyUserRoles); result.Error != nil {
		return nil, result.Error
	}

	for _, lcompanyUserRole := range companyUserRoles {
		companyUserRole := lcompanyUserRole.ToDomain()

		if companyUserRole.RoleType == 1 {
			list = append(list, companyUserRole)
		} else {
			for index, l := range list {
				if l.AdvertiserId == companyUserRole.AdvertiserId && l.UserId == companyUserRole.UserId {
					list = append(list[:index], list[index+1:]...)
				}
			}
		}
	}

	return list, nil
}

func (curr *companyUserRoleRepo) ListByUserIdAndCompanyIdAndDay(ctx context.Context, userId, companyId uint64, day uint32) ([]*domain.CompanyUserRole, error) {
	var companyUserRoles []CompanyUserRole
	list := make([]*domain.CompanyUserRole, 0)

	if result := curr.data.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Where("company_id = ?", companyId).
		Where("day <= ?", day).
		Order("day ASC").
		Find(&companyUserRoles); result.Error != nil {
		return nil, result.Error
	}

	for _, lcompanyUserRole := range companyUserRoles {
		companyUserRole := lcompanyUserRole.ToDomain()

		if companyUserRole.RoleType == 1 {
			list = append(list, companyUserRole)
		} else {
			for index, l := range list {
				if l.AdvertiserId == companyUserRole.AdvertiserId {
					list = append(list[:index], list[index+1:]...)
				}
			}
		}
	}

	return list, nil
}

func (curr *companyUserRoleRepo) ListByCompanyIdAndAdvertiserIdAndDay(ctx context.Context, companyId, advertiserId uint64, day uint32) ([]*domain.CompanyUserRole, error) {
	var companyUserRoles []CompanyUserRole
	list := make([]*domain.CompanyUserRole, 0)

	if result := curr.data.db.WithContext(ctx).
		Where("company_id = ?", companyId).
		Where("advertiser_id = ?", advertiserId).
		Where("day <= ?", day).
		Order("day ASC").
		Find(&companyUserRoles); result.Error != nil {
		return nil, result.Error
	}

	for _, lcompanyUserRole := range companyUserRoles {
		companyUserRole := lcompanyUserRole.ToDomain()

		if companyUserRole.RoleType == 1 {
			list = append(list, companyUserRole)
		} else {
			for index, l := range list {
				if l.AdvertiserId == companyUserRole.AdvertiserId && l.UserId == companyUserRole.UserId {
					list = append(list[:index], list[index+1:]...)
				}
			}
		}
	}

	return list, nil
}

func (curr *companyUserRoleRepo) Save(ctx context.Context, in *domain.CompanyUserRole) (*domain.CompanyUserRole, error) {
	companyUserRole := &CompanyUserRole{
		UserId:       in.UserId,
		AdvertiserId: in.AdvertiserId,
		CompanyId:    in.CompanyId,
		RoleType:     in.RoleType,
		Day:          in.Day,
		CreateTime:   in.CreateTime,
		UpdateTime:   in.UpdateTime,
	}

	if result := curr.data.DB(ctx).Create(companyUserRole); result.Error != nil {
		return nil, result.Error
	}

	return companyUserRole.ToDomain(), nil
}

func (curr *companyUserRoleRepo) Update(ctx context.Context, in *domain.CompanyUserRole) (*domain.CompanyUserRole, error) {
	if result := curr.data.DB(ctx).Model(CompanyUserRole{}).
		Where("user_id = ?", in.UserId).
		Where("advertiser_id = ?", in.AdvertiserId).
		Where("company_id = ?", in.CompanyId).
		Where("day = ?", in.Day).
		Updates(CompanyUserRole{RoleType: in.RoleType, UpdateTime: in.UpdateTime}); result.Error != nil {
		return nil, result.Error
	}

	return in, nil
}

func (curr *companyUserRoleRepo) DeleteByUserIdAndCompanyIdAndAdvertiserIdAndDay(ctx context.Context, userId, companyId, advertiserId uint64, day uint32) error {
	if result := curr.data.DB(ctx).
		Where("user_id = ?", userId).
		Where("company_id = ?", companyId).
		Where("advertiser_id = ?", advertiserId).
		Where("day = ?", day).
		Delete(&CompanyUserRole{}); result.Error != nil {
		return result.Error
	}

	return nil
}

// ////////////////////要删除的////////////////////////////////////
func (curr *companyUserRoleRepo) ListByUserIdAndCompanyId(ctx context.Context, id, companyId uint64) ([]*domain.CompanyUserRole, error) {
	var companyUserRoles []CompanyUserRole
	list := make([]*domain.CompanyUserRole, 0)

	if result := curr.data.db.WithContext(ctx).
		Where("user_id = ?", id).
		Where("company_id = ?", companyId).
		Find(&companyUserRoles); result.Error != nil {
		return nil, result.Error
	}

	for _, companyUserRole := range companyUserRoles {
		list = append(list, companyUserRole.ToDomain())
	}

	return list, nil
}

//////////////////////////////////////////////////////////

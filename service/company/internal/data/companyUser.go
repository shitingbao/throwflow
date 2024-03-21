package data

import (
	"company/internal/biz"
	"company/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"time"
	"unicode/utf8"
)

// 企业库用户表
type CompanyUser struct {
	Id             uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	Username       string    `gorm:"column:username;type:varchar(50);not null;comment:姓名"`
	Job            string    `gorm:"column:job;type:varchar(50);not null;comment:职位"`
	Phone          string    `gorm:"column:phone;type:varchar(20);not null;comment:手机号"`
	CompanyId      uint64    `gorm:"column:company_id;type:bigint(20) UNSIGNED;not null;comment:企业库ID"`
	Role           uint8     `gorm:"column:role;type:tinyint(3) UNSIGNED;not null;default:3;comment:状态：角色，1：管理员，2：副管理员，3：普通成员"`
	Status         uint8     `gorm:"column:status;type:tinyint(3) UNSIGNED;not null;default:1;comment:状态：1：启用，0：禁用"`
	OrganizationId uint64    `gorm:"column:organization_id;type:bigint(20) UNSIGNED;not null;default:0;comment:企业机构ID"`
	CreateTime     time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime     time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (CompanyUser) TableName() string {
	return "company_company_user"
}

type companyUserRepo struct {
	data *Data
	log  *log.Helper
}

func (cu *CompanyUser) ToDomain() *domain.CompanyUser {
	companyUser := &domain.CompanyUser{
		Id:             cu.Id,
		Username:       cu.Username,
		Job:            cu.Job,
		Phone:          cu.Phone,
		CompanyId:      cu.CompanyId,
		Role:           cu.Role,
		Status:         cu.Status,
		OrganizationId: cu.OrganizationId,
		CreateTime:     cu.CreateTime,
		UpdateTime:     cu.UpdateTime,
	}

	selectCompanyUsers := domain.NewSelectCompanyUsers()

	for _, role := range selectCompanyUsers.Role {
		iRole, _ := strconv.Atoi(role.Key)

		if uint8(iRole) == companyUser.Role {
			companyUser.RoleName = role.Value
			break
		}
	}

	return companyUser
}

func NewCompanyUserRepo(data *Data, logger log.Logger) biz.CompanyUserRepo {
	return &companyUserRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cur *companyUserRepo) GetById(ctx context.Context, id, companyId uint64) (*domain.CompanyUser, error) {
	companyUser := &CompanyUser{}

	if result := cur.data.db.WithContext(ctx).Where("company_id = ?", companyId).Where("id = ?", id).First(companyUser); result.Error != nil {
		return nil, result.Error
	}

	return companyUser.ToDomain(), nil
}

func (cur *companyUserRepo) GetByPhone(ctx context.Context, companyId uint64, phone string) (*domain.CompanyUser, error) {
	companyUser := &CompanyUser{}

	if companyId > 0 {
		if result := cur.data.db.WithContext(ctx).Where("phone = ?", phone).Where("company_id = ?", companyId).First(companyUser); result.Error != nil {
			return nil, result.Error
		}
	} else {
		if result := cur.data.db.WithContext(ctx).Where("phone = ?", phone).First(companyUser); result.Error != nil {
			return nil, result.Error
		}
	}

	return companyUser.ToDomain(), nil
}

func (cur *companyUserRepo) GetByPhoneAndNotInUserId(ctx context.Context, companyId uint64, phone string, userIds []uint64) (*domain.CompanyUser, error) {
	companyUser := &CompanyUser{}

	if result := cur.data.db.WithContext(ctx).Where("id not in ?", userIds).Where("phone = ?", phone).Where("company_id = ?", companyId).First(companyUser); result.Error != nil {
		return nil, result.Error
	}

	return companyUser.ToDomain(), nil
}

func (cur *companyUserRepo) GetByOrganizationIdAndPhone(ctx context.Context, companyId, organizationId uint64, phone string) (*domain.CompanyUser, error) {
	companyUser := &CompanyUser{}

	if result := cur.data.db.WithContext(ctx).Where("phone = ?", phone).Where("company_id = ?", companyId).Where("organization_id = ?", organizationId).First(companyUser); result.Error != nil {
		return nil, result.Error
	}

	return companyUser.ToDomain(), nil
}

func (cur *companyUserRepo) GetByRole(ctx context.Context, companyId uint64, role uint8) (*domain.CompanyUser, error) {
	companyUser := &CompanyUser{}

	if result := cur.data.db.WithContext(ctx).Where("role = ?", role).Where("company_id = ?", companyId).First(companyUser); result.Error != nil {
		return nil, result.Error
	}

	return companyUser.ToDomain(), nil
}

func (cur *companyUserRepo) ListByCompanyId(ctx context.Context, companyId uint64) ([]*domain.CompanyUser, error) {
	var companyUsers []CompanyUser
	list := make([]*domain.CompanyUser, 0)

	if result := cur.data.db.WithContext(ctx).
		Where("company_id = ?", companyId).
		Order("role ASC").
		Find(&companyUsers); result.Error != nil {
		return nil, result.Error
	}

	for _, companyUser := range companyUsers {
		list = append(list, companyUser.ToDomain())
	}

	return list, nil
}

func (cur *companyUserRepo) ListByCompanyIdAndOrganizationId(ctx context.Context, companyId, organizationId uint64) ([]*domain.CompanyUser, error) {
	var companyUsers []CompanyUser
	list := make([]*domain.CompanyUser, 0)

	if result := cur.data.db.WithContext(ctx).
		Where("company_id = ?", companyId).
		Where("organization_id = ?", organizationId).
		Order("role ASC").
		Find(&companyUsers); result.Error != nil {
		return nil, result.Error
	}

	for _, companyUser := range companyUsers {
		list = append(list, companyUser.ToDomain())
	}

	return list, nil
}

func (cur *companyUserRepo) ListByPhone(ctx context.Context, phone string) ([]*domain.CompanyUser, error) {
	var companyUsers []CompanyUser
	list := make([]*domain.CompanyUser, 0)

	if result := cur.data.db.WithContext(ctx).
		Where("phone = ?", phone).
		Find(&companyUsers); result.Error != nil {
		return nil, result.Error
	}

	for _, companyUser := range companyUsers {
		list = append(list, companyUser.ToDomain())
	}

	return list, nil
}

func (cur *companyUserRepo) List(ctx context.Context, companyId uint64, pageNum, pageSize int, keyword string) ([]*domain.CompanyUser, error) {
	var companyUsers []CompanyUser
	list := make([]*domain.CompanyUser, 0)

	db := cur.data.db.WithContext(ctx).Where("company_id = ?", companyId)

	if l := utf8.RuneCountInString(keyword); l > 0 {
		db = db.Where("username like ?", "%"+keyword+"%")
	}

	if pageNum == 0 {
		if result := db.Order("status DESC,role ASC,id DESC").
			Find(&companyUsers); result.Error != nil {
			return nil, result.Error
		}
	} else {
		if result := db.Order("status DESC,role ASC,id DESC").
			Limit(pageSize).Offset((pageNum - 1) * pageSize).
			Find(&companyUsers); result.Error != nil {
			return nil, result.Error
		}
	}

	for _, companyUser := range companyUsers {
		list = append(list, companyUser.ToDomain())
	}

	return list, nil
}

func (cur *companyUserRepo) Count(ctx context.Context, companyId uint64, keyword string) (int64, error) {
	var count int64

	db := cur.data.db.WithContext(ctx).Where("company_id = ?", companyId)

	if l := utf8.RuneCountInString(keyword); l > 0 {
		db = db.Where("username like ?", "%"+keyword+"%")
	}

	if result := db.Model(&CompanyUser{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (cur *companyUserRepo) Save(ctx context.Context, in *domain.CompanyUser) (*domain.CompanyUser, error) {
	companyUser := &CompanyUser{
		Username:       in.Username,
		Job:            in.Job,
		Phone:          in.Phone,
		CompanyId:      in.CompanyId,
		Role:           in.Role,
		Status:         in.Status,
		OrganizationId: in.OrganizationId,
		CreateTime:     in.CreateTime,
		UpdateTime:     in.UpdateTime,
	}

	if result := cur.data.DB(ctx).Create(companyUser); result.Error != nil {
		return nil, result.Error
	}

	return companyUser.ToDomain(), nil
}

func (cur *companyUserRepo) Update(ctx context.Context, in *domain.CompanyUser) (*domain.CompanyUser, error) {
	companyUser := &CompanyUser{
		Id:             in.Id,
		Username:       in.Username,
		Job:            in.Job,
		Phone:          in.Phone,
		CompanyId:      in.CompanyId,
		Role:           in.Role,
		Status:         in.Status,
		OrganizationId: in.OrganizationId,
		CreateTime:     in.CreateTime,
		UpdateTime:     in.UpdateTime,
	}

	if result := cur.data.DB(ctx).Save(companyUser); result.Error != nil {
		return nil, result.Error
	}

	return companyUser.ToDomain(), nil
}

func (cur *companyUserRepo) DeleteByCompanyId(ctx context.Context, companyId uint64) error {
	if result := cur.data.DB(ctx).Where("company_id = ?", companyId).Delete(&CompanyUser{}); result.Error != nil {
		return result.Error
	}

	return nil
}

func (cur *companyUserRepo) DeleteByCompanyIdAndUserId(ctx context.Context, companyId uint64, userIds []uint64) error {
	if result := cur.data.DB(ctx).Where("id in ?", userIds).Where("company_id = ?", companyId).Delete(&CompanyUser{}); result.Error != nil {
		return result.Error
	}

	return nil
}

func (cur *companyUserRepo) DeleteByCompanyIdAndOrganizationId(ctx context.Context, companyId, organizationId uint64) error {
	if result := cur.data.DB(ctx).Where("company_id = ?", companyId).Where("organization_id = ?", organizationId).Delete(&CompanyUser{}); result.Error != nil {
		return result.Error
	}

	return nil
}

func (c *companyUserRepo) Delete(ctx context.Context, in *domain.CompanyUser) error {
	companyUser := &CompanyUser{
		Id:             in.Id,
		Username:       in.Username,
		Job:            in.Job,
		Phone:          in.Phone,
		CompanyId:      in.CompanyId,
		Role:           in.Role,
		Status:         in.Status,
		OrganizationId: in.OrganizationId,
		CreateTime:     in.CreateTime,
		UpdateTime:     in.UpdateTime,
	}

	if result := c.data.db.WithContext(ctx).Delete(companyUser); result.Error != nil {
		return result.Error
	}

	return nil
}

func (c *companyUserRepo) SaveCacheHash(ctx context.Context, key string, val map[string]string, timeout time.Duration) error {
	_, err := c.data.rdb.HMSet(ctx, key, val).Result()

	if err != nil {
		return err
	}

	_, err = c.data.rdb.Expire(ctx, key, timeout).Result()

	if err != nil {
		return err
	}

	return nil
}

func (c *companyUserRepo) UpdateCacheHash(ctx context.Context, key string, val map[string]string) error {
	_, err := c.data.rdb.HMSet(ctx, key, val).Result()

	if err != nil {
		return err
	}

	return nil
}

func (u *companyUserRepo) GetCacheHash(ctx context.Context, key string, field string) (string, error) {
	val, err := u.data.rdb.HGet(ctx, key, field).Result()

	if err != nil {
		return "", err
	}

	return val, nil
}

func (c *companyUserRepo) SaveCacheString(ctx context.Context, key string, val string, timeout time.Duration) error {
	_, err := c.data.rdb.Set(ctx, key, val, timeout).Result()

	if err != nil {
		return err
	}

	return nil
}

func (u *companyUserRepo) GetCacheString(ctx context.Context, key string) (string, error) {
	val, err := u.data.rdb.Get(ctx, key).Result()

	if err != nil {
		return "", err
	}

	return val, nil
}

func (c *companyUserRepo) ExpireCache(ctx context.Context, key string, timeout time.Duration) error {
	_, err := c.data.rdb.Expire(ctx, key, timeout).Result()

	if err != nil {
		return err
	}

	return nil
}

func (u *companyUserRepo) DeleteCache(ctx context.Context, key string) error {
	if _, err := u.data.rdb.Del(ctx, key).Result(); err != nil {
		return err
	}

	return nil
}

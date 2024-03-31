package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"unicode/utf8"
	"weixin/internal/biz"
	"weixin/internal/domain"
)

// 微信小程序用户佣金表
type UserCommission struct {
	Id                      uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	UserId                  uint64    `gorm:"column:user_id;type:bigint(20) UNSIGNED;not null;comment:微信小程序用户ID"`
	OrganizationId          uint64    `gorm:"column:organization_id;type:bigint(20) UNSIGNED;not null;comment:机构ID"`
	ChildUserId             uint64    `gorm:"column:child_user_id;type:bigint(20) UNSIGNED;not null;comment:下级用户微信小程序用户ID"`
	ChildLevel              uint8     `gorm:"column:child_level;type:tinyint(3) UNSIGNED;not null;default:0;comment:下级用户微信小程序等级"`
	Level                   uint8     `gorm:"column:level;type:tinyint(3) UNSIGNED;not null;default:0;comment:等级"`
	Relation                uint8     `gorm:"column:relation;type:tinyint(3) UNSIGNED;not null;default:0;comment:1：直接，2：间接"`
	OrderNum                uint64    `gorm:"column:order_num;type:bigint(20) UNSIGNED;not null;comment:订单数"`
	OrderRefundNum          uint64    `gorm:"column:order_refund_num;type:bigint(20) UNSIGNED;not null;comment:退货订单数"`
	TotalPayAmount          float32   `gorm:"column:total_pay_amount;type:decimal(10, 2) UNSIGNED;not null;comment:订单支付金额，单位元"`
	CommissionPool          float32   `gorm:"column:commission_pool;type:decimal(10, 2) UNSIGNED;not null;comment:分佣池，单位元"`
	EstimatedCommission     float32   `gorm:"column:estimated_commission;type:decimal(10, 2) UNSIGNED;not null;comment:预估佣金收入，单位元"`
	RealCommission          float32   `gorm:"column:real_commission;type:decimal(10, 2) UNSIGNED;not null;comment:实际佣金收入，单位元"`
	EstimatedUserCommission float32   `gorm:"column:estimated_user_commission;type:decimal(10, 2) UNSIGNED;not null;comment:机构返佣预估佣金收入，单位元"`
	RealUserCommission      float32   `gorm:"column:real_user_commission;type:decimal(10, 2) UNSIGNED;not null;comment:机构返佣实际佣金收入，单位元"`
	Day                     uint32    `gorm:"column:day;type:int(11);not null;comment:分佣时间"`
	UserCommissionType      uint8     `gorm:"column:user_commission_type;type:tinyint(3) UNSIGNED;not null;default:0;comment:1：课程，2：电商, 3：成本购"`
	CreateTime              time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime              time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

type UserCommissionData struct {
	ChildUserId             uint64
	ChildNickName           string
	ChildAvatarUrl          string
	ChildPhone              string
	Relation                uint8
	UserCommissionType      uint8
	TotalPayAmount          float32
	CommissionPool          float32
	EstimatedUserCommission float32
	RealUserCommission      float32
}

func (UserCommission) TableName() string {
	return "weixin_user_commission"
}

type userCommissionRepo struct {
	data *Data
	log  *log.Helper
}

func (uc *UserCommission) ToDomain(ctx context.Context) *domain.UserCommission {
	userCommission := &domain.UserCommission{
		Id:                      uc.Id,
		UserId:                  uc.UserId,
		OrganizationId:          uc.OrganizationId,
		ChildUserId:             uc.ChildUserId,
		ChildLevel:              uc.ChildLevel,
		Level:                   uc.Level,
		Relation:                uc.Relation,
		OrderNum:                uc.OrderNum,
		OrderRefundNum:          uc.OrderRefundNum,
		TotalPayAmount:          uc.TotalPayAmount,
		CommissionPool:          uc.CommissionPool,
		EstimatedCommission:     uc.EstimatedCommission,
		RealCommission:          uc.RealCommission,
		EstimatedUserCommission: uc.EstimatedUserCommission,
		RealUserCommission:      uc.RealUserCommission,
		Day:                     uc.Day,
		UserCommissionType:      uc.UserCommissionType,
		CreateTime:              uc.CreateTime,
		UpdateTime:              uc.UpdateTime,
	}

	return userCommission
}

func NewUserCommissionRepo(data *Data, logger log.Logger) biz.UserCommissionRepo {
	return &userCommissionRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ucr *userCommissionRepo) GetByChildUserId(ctx context.Context, userId, childUserId, organizationId uint64, userCommissionType uint8, startDay, endDay string) (*domain.UserCommission, error) {
	userCommission := &UserCommission{}

	db := ucr.data.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Where("child_user_id = ?", childUserId).
		Where("organization_id = ?", organizationId).
		Where("user_commission_type = ?", userCommissionType)

	if len(startDay) > 0 {
		db = db.Where("day >= ?", startDay)
	}

	if len(endDay) > 0 {
		db = db.Where("day <= ?", endDay)
	}

	if result := db.Order("day desc").
		First(userCommission); result.Error != nil {
		return nil, result.Error
	}

	return userCommission.ToDomain(ctx), nil
}

func (ucr *userCommissionRepo) List(ctx context.Context, pageNum, pageSize int, userId, organizationId uint64, isDirect uint8, startDay, endDay, keyword string) ([]*domain.UserCommission, error) {
	var userCommissionDatas []UserCommissionData
	list := make([]*domain.UserCommission, 0)

	db := ucr.data.db.WithContext(ctx).Table("weixin_user_commission").
		Joins("left join weixin_user on weixin_user_commission.child_user_id = weixin_user.id").
		Where("weixin_user_commission.user_id = ?", userId).
		Where("weixin_user_commission.organization_id = ?", organizationId).
		Select("weixin_user.nick_name as child_nick_name, weixin_user.avatar_url as child_avatar_url, weixin_user.phone as child_phone, weixin_user_commission.child_user_id, weixin_user_commission.relation, weixin_user_commission.user_commission_type, sum(total_pay_amount) as total_pay_amount, sum(commission_pool) as commission_pool, sum(estimated_user_commission) as estimated_user_commission, sum(real_user_commission) as real_user_commission")

	if len(startDay) > 0 {
		db = db.Where("weixin_user_commission.day >= ?", startDay)
	}

	if len(endDay) > 0 {
		db = db.Where("weixin_user_commission.day <= ?", endDay)
	}

	if isDirect == 1 {
		db = db.Where("weixin_user_commission.relation = 1")
	} else if isDirect == 2 {
		db = db.Where("weixin_user_commission.relation = 2")
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		db = db.Where("(weixin_user.nick_name like ? or weixin_user.phone like ?)", "%"+keyword+"%", "%"+keyword+"%")
	}

	if result := db.Group("weixin_user_commission.child_user_id,weixin_user_commission.relation,weixin_user_commission.user_commission_type").Order("real_user_commission, child_user_id DESC").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&userCommissionDatas); result.Error != nil {
		return nil, result.Error
	}

	for _, userCommissionData := range userCommissionDatas {
		list = append(list, &domain.UserCommission{
			ChildUserId:             userCommissionData.ChildUserId,
			ChildNickName:           userCommissionData.ChildNickName,
			ChildAvatarUrl:          userCommissionData.ChildAvatarUrl,
			ChildPhone:              userCommissionData.ChildPhone,
			Relation:                userCommissionData.Relation,
			UserCommissionType:      userCommissionData.UserCommissionType,
			TotalPayAmount:          userCommissionData.TotalPayAmount,
			CommissionPool:          userCommissionData.CommissionPool,
			EstimatedUserCommission: userCommissionData.EstimatedUserCommission,
			RealUserCommission:      userCommissionData.RealUserCommission,
		})
	}

	return list, nil
}

func (ucr *userCommissionRepo) ListByDay(ctx context.Context, day uint32, userCommissionType []string) ([]*domain.UserCommission, error) {
	var userCommissions []UserCommission
	list := make([]*domain.UserCommission, 0)

	db := ucr.data.DB(ctx).Where("day = ?", day)

	if len(userCommissionType) > 0 {
		db = db.Where("user_commission_type = ?", userCommissionType)
	}

	if result := db.Find(&userCommissions); result.Error != nil {
		return nil, result.Error
	}

	for _, userCommission := range userCommissions {
		list = append(list, userCommission.ToDomain(ctx))
	}

	return list, nil
}

func (ucr *userCommissionRepo) Count(ctx context.Context, userId, organizationId uint64, isDirect uint8, startDay, endDay, keyword string) (int64, error) {
	var count int64

	db := ucr.data.db.WithContext(ctx).Table("weixin_user_commission").
		Joins("left join weixin_user on weixin_user_commission.child_user_id = weixin_user.id").
		Where("weixin_user_commission.user_id = ?", userId).
		Where("weixin_user_commission.organization_id = ?", organizationId)

	if len(startDay) > 0 {
		db = db.Where("weixin_user_commission.day >= ?", startDay)
	}

	if len(endDay) > 0 {
		db = db.Where("weixin_user_commission.day <= ?", endDay)
	}

	if isDirect == 1 {
		db = db.Where("weixin_user_commission.relation = 1")
	} else if isDirect == 2 {
		db = db.Where("weixin_user_commission.relation = 2")
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		db = db.Where("(weixin_user.nick_name like ? or weixin_user.phone like ?)", "%"+keyword+"%", "%"+keyword+"%")
	}

	db = db.Group("weixin_user_commission.child_user_id,weixin_user_commission.relation,weixin_user_commission.user_commission_type")

	if result := db.Model(&UserCommission{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (ucr *userCommissionRepo) Statistics(ctx context.Context, userId uint64, day uint32, userCommissionType uint8) (*domain.UserCommission, error) {
	userCommission := &UserCommission{}

	if userCommissionType == 1 {
		db := ucr.data.db.WithContext(ctx).
			Where("user_id = ?", userId).
			Where("day <= ?", day).Where("user_commission_type = ?", userCommissionType).
			Select("sum(total_pay_amount) as total_pay_amount, sum(commission_pool) as commission_pool, sum(estimated_user_commission) as estimated_user_commission")

		if result := db.First(userCommission); result.Error != nil {
			return nil, result.Error
		}
	} else {
		db := ucr.data.db.WithContext(ctx).
			Where("id in (?)", ucr.data.db.Table("weixin_user_commission").Select("max(id) as id").
				Where("user_id = ?", userId).
				Where("day <= ?", day).
				Where("user_commission_type = ?", userCommissionType).Group("child_user_id").Order("")).
			Select("sum(total_pay_amount) as total_pay_amount, sum(commission_pool) as commission_pool, sum(estimated_user_commission) as estimated_user_commission")

		if result := db.First(userCommission); result.Error != nil {
			return nil, result.Error
		}
	}

	return userCommission.ToDomain(ctx), nil
}

func (ucr *userCommissionRepo) StatisticsReal(ctx context.Context, userId uint64, day uint32, userCommissionType uint8) (*domain.UserCommission, error) {
	userCommission := &UserCommission{}

	db := ucr.data.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Where("day <= ?", day)

	if userCommissionType > 0 {
		db = db.Where("user_commission_type = ?", userCommissionType)
	}

	db = db.Select("sum(real_user_commission) as real_user_commission")

	if result := db.First(userCommission); result.Error != nil {
		return nil, result.Error
	}

	return userCommission.ToDomain(ctx), nil
}

func (ucr *userCommissionRepo) Update(ctx context.Context, in *domain.UserCommission) (*domain.UserCommission, error) {
	userCommission := &UserCommission{
		Id:                      in.Id,
		UserId:                  in.UserId,
		OrganizationId:          in.OrganizationId,
		ChildUserId:             in.ChildUserId,
		ChildLevel:              in.ChildLevel,
		Level:                   in.Level,
		Relation:                in.Relation,
		OrderNum:                in.OrderNum,
		OrderRefundNum:          in.OrderRefundNum,
		TotalPayAmount:          in.TotalPayAmount,
		CommissionPool:          in.CommissionPool,
		EstimatedCommission:     in.EstimatedCommission,
		RealCommission:          in.RealCommission,
		EstimatedUserCommission: in.EstimatedUserCommission,
		RealUserCommission:      in.RealUserCommission,
		Day:                     in.Day,
		UserCommissionType:      in.UserCommissionType,
		CreateTime:              in.CreateTime,
		UpdateTime:              in.UpdateTime,
	}

	if result := ucr.data.DB(ctx).Save(userCommission); result.Error != nil {
		return nil, result.Error
	}

	return userCommission.ToDomain(ctx), nil
}

func (ucr *userCommissionRepo) Save(ctx context.Context, in *domain.UserCommission) (*domain.UserCommission, error) {
	userCommission := &UserCommission{
		UserId:                  in.UserId,
		OrganizationId:          in.OrganizationId,
		ChildUserId:             in.ChildUserId,
		ChildLevel:              in.ChildLevel,
		Level:                   in.Level,
		Relation:                in.Relation,
		OrderNum:                in.OrderNum,
		OrderRefundNum:          in.OrderRefundNum,
		TotalPayAmount:          in.TotalPayAmount,
		CommissionPool:          in.CommissionPool,
		EstimatedCommission:     in.EstimatedCommission,
		RealCommission:          in.RealCommission,
		EstimatedUserCommission: in.EstimatedUserCommission,
		RealUserCommission:      in.RealUserCommission,
		Day:                     in.Day,
		UserCommissionType:      in.UserCommissionType,
		CreateTime:              in.CreateTime,
		UpdateTime:              in.UpdateTime,
	}

	if result := ucr.data.DB(ctx).Create(userCommission); result.Error != nil {
		return nil, result.Error
	}

	return userCommission.ToDomain(ctx), nil
}

func (ucr *userCommissionRepo) DeleteByDay(ctx context.Context, day uint32, userCommissionTypes []string) error {
	db := ucr.data.DB(ctx).Where("day = ?", day)

	if len(userCommissionTypes) > 0 {
		db = db.Where("user_commission_type in ?", userCommissionTypes)
	}

	if result := db.Delete(&UserCommission{}); result.Error != nil {
		return result.Error
	}

	return nil
}

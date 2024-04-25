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
	Id                 uint64     `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	UserId             uint64     `gorm:"column:user_id;type:bigint(20) UNSIGNED;not null;comment:微信小程序用户ID"`
	OrganizationId     uint64     `gorm:"column:organization_id;type:bigint(20) UNSIGNED;not null;comment:机构ID"`
	ChildUserId        uint64     `gorm:"column:child_user_id;type:bigint(20) UNSIGNED;not null;index:child_user_id_relevance_id_commission_type,priority:1;comment:下级用户微信小程序用户ID"`
	ChildLevel         uint8      `gorm:"column:child_level;type:tinyint(3) UNSIGNED;not null;default:0;comment:下级用户微信小程序等级"`
	Level              uint8      `gorm:"column:level;type:tinyint(3) UNSIGNED;not null;default:0;comment:等级"`
	Relation           uint8      `gorm:"column:relation;type:tinyint(3) UNSIGNED;not null;default:0;comment:1：直接，2：间接"`
	CommissionStatus   uint8      `gorm:"column:commission_status;type:tinyint(3) UNSIGNED;not null;default:0;comment:佣金状态，1：待结算，2：已结算, 3:待体现"`
	CommissionType     uint8      `gorm:"column:commission_type;type:tinyint(3) UNSIGNED;not null;default:0;index:child_user_id_relevance_id_commission_type,priority:3;comment:1：课程，2：电商, 3：成本购，4：成本购返佣，5：任务，6：其他，7：分佣提现，8：成本购提现，前端传递参数，1：成本购，2：电商，课程"`
	RelevanceId        uint64     `gorm:"column:relevance_id;type:bigint(20) UNSIGNED;not null;default:0;index:child_user_id_relevance_id_commission_type,priority:2;comment:金额加入关联的佣金ID"`
	OperationType      uint8      `gorm:"column:operation_type;type:tinyint(3) UNSIGNED;not null;default:0;comment:1：加入，2：减少"`
	TotalPayAmount     float32    `gorm:"column:total_pay_amount;type:decimal(10, 2) UNSIGNED;not null;comment:订单支付金额，单位元"`
	CommissionPool     float32    `gorm:"column:commission_pool;type:decimal(10, 2) UNSIGNED;not null;comment:分佣池，单位元"`
	CommissionMcnRatio float32    `gorm:"column:commission_mcn_ratio;type:decimal(10, 2) UNSIGNED;not null;comment:mcn机构返佣比例，单位元"`
	CommissionAmount   float32    `gorm:"column:commission_amount;type:decimal(10, 2) UNSIGNED;not null;comment:机构返佣实际佣金收入，单位元"`
	Name               string     `gorm:"column:name;type:text;not null;;comment:姓名"`
	IdentityCard       string     `gorm:"column:identity_card;type:text;not null;comment:身份证号码"`
	BankCode           string     `gorm:"column:bank_code;type:text;not null;comment:银行卡号"`
	BalanceStatus      uint8      `gorm:"column:balance_status;type:tinyint(3) UNSIGNED;not null;default:0;comment:余额状态，0：未处理，1：已处理，2：正在处理, 3:处理失败"`
	OutTradeNo         string     `gorm:"column:out_trade_no;type:varchar(32);not null;default:'';comment:提现订单号"`
	InnerTradeNo       string     `gorm:"column:inner_trade_no;type:varchar(100);not null;default:'';comment:工猫提现订单号"`
	RealAmount         float32    `gorm:"column:real_amount;type:decimal(10, 2) UNSIGNED;not null;default:0.00;comment:实发金额，单位元"`
	ApplyTime          *time.Time `gorm:"column:apply_time;type:datetime;null;default:null;comment:申请时间"`
	GongmallCreateTime *time.Time `gorm:"column:gongmall_create_time;type:datetime;null;default:null;comment:工猫进单时间"`
	PayTime            *time.Time `gorm:"column:pay_time;type:datetime;null;default:null;comment:工猫实际付款时间"`
	OperationContent   string     `gorm:"column:operation_content;type:text;not null;comment:说明"`
	CreateTime         time.Time  `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime         time.Time  `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

type UserCommissionData struct {
	ChildUserId      uint64
	ChildLevel       uint8
	ChildNickName    string
	ChildAvatarUrl   string
	ChildPhone       string
	CommissionType   uint8
	OperationType    uint8
	TotalPayAmount   float32
	CommissionPool   float32
	CommissionAmount float32
	BalanceStatus    uint8
	OperationContent string
	CreateTime       time.Time
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
		Id:                 uc.Id,
		UserId:             uc.UserId,
		OrganizationId:     uc.OrganizationId,
		ChildUserId:        uc.ChildUserId,
		ChildLevel:         uc.ChildLevel,
		Level:              uc.Level,
		Relation:           uc.Relation,
		CommissionStatus:   uc.CommissionStatus,
		CommissionType:     uc.CommissionType,
		RelevanceId:        uc.RelevanceId,
		OperationType:      uc.OperationType,
		TotalPayAmount:     uc.TotalPayAmount,
		CommissionPool:     uc.CommissionPool,
		CommissionMcnRatio: uc.CommissionMcnRatio,
		CommissionAmount:   uc.CommissionAmount,
		Name:               uc.Name,
		IdentityCard:       uc.IdentityCard,
		BankCode:           uc.BankCode,
		BalanceStatus:      uc.BalanceStatus,
		OutTradeNo:         uc.OutTradeNo,
		InnerTradeNo:       uc.InnerTradeNo,
		RealAmount:         uc.RealAmount,
		ApplyTime:          uc.ApplyTime,
		GongmallCreateTime: uc.GongmallCreateTime,
		PayTime:            uc.PayTime,
		OperationContent:   uc.OperationContent,
		CreateTime:         uc.CreateTime,
		UpdateTime:         uc.UpdateTime,
	}

	return userCommission
}

func NewUserCommissionRepo(data *Data, logger log.Logger) biz.UserCommissionRepo {
	return &userCommissionRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ucr *userCommissionRepo) GetByOutTradeNo(ctx context.Context, outTradeNo string) (*domain.UserCommission, error) {
	userCommission := &UserCommission{}

	if result := ucr.data.db.WithContext(ctx).
		Where("out_trade_no = ?", outTradeNo).
		First(userCommission); result.Error != nil {
		return nil, result.Error
	}

	return userCommission.ToDomain(ctx), nil
}

func (ucr *userCommissionRepo) GetByRelevanceId(ctx context.Context, relevanceId uint64, commissionType uint8) (*domain.UserCommission, error) {
	userCommission := &UserCommission{}

	if result := ucr.data.db.WithContext(ctx).
		Where("relevance_id = ?", relevanceId).
		Where("commission_type = ?", commissionType).
		First(userCommission); result.Error != nil {
		return nil, result.Error
	}

	return userCommission.ToDomain(ctx), nil
}

func (ucr *userCommissionRepo) List(ctx context.Context, pageNum, pageSize int, userId, organizationId uint64, childIds []uint64, commissionType uint8, startDay, endDay, keyword string) ([]*domain.UserCommission, error) {
	var userCommissionDatas []UserCommissionData
	list := make([]*domain.UserCommission, 0)

	db := ucr.data.db.WithContext(ctx).Table("weixin_user")

	if commissionType == 0 {
		if len(startDay) > 0 {
			db = db.Joins("left join weixin_user_commission on weixin_user_commission.child_user_id = weixin_user.id and weixin_user_commission.user_id = ? and weixin_user_commission.organization_id = ? and weixin_user_commission.commission_type in (1,2,4) and weixin_user_commission.create_time >= ? and weixin_user_commission.create_time <= ?", userId, organizationId, startDay+" 00:00:00", endDay+" 23:59:59")
		} else {
			db = db.Joins("left join weixin_user_commission on weixin_user_commission.child_user_id = weixin_user.id and weixin_user_commission.user_id = ? and weixin_user_commission.organization_id = ? and weixin_user_commission.commission_type in (1,2,4)", userId, organizationId)
		}
	} else {
		if len(startDay) > 0 {
			db = db.Joins("left join weixin_user_commission on weixin_user_commission.child_user_id = weixin_user.id and weixin_user_commission.user_id = ? and weixin_user_commission.organization_id = ? and weixin_user_commission.commission_type = ? and weixin_user_commission.create_time >= ? and weixin_user_commission.create_time <= ?", userId, organizationId, commissionType, startDay+" 00:00:00", endDay+" 23:59:59")
		} else {
			db = db.Joins("left join weixin_user_commission on weixin_user_commission.child_user_id = weixin_user.id and weixin_user_commission.user_id = ? and weixin_user_commission.organization_id = ? and weixin_user_commission.commission_type = ?", userId, organizationId, commissionType)
		}
	}

	db = db.Where("weixin_user.id in (?)", childIds).
		Select("weixin_user.nick_name as child_nick_name, weixin_user.avatar_url as child_avatar_url, weixin_user.phone as child_phone, weixin_user.id as child_user_id, weixin_user_commission.commission_type, sum(weixin_user_commission.total_pay_amount) as total_pay_amount, sum(weixin_user_commission.commission_pool) as commission_pool, sum(weixin_user_commission.commission_amount) as commission_amount")

	if l := utf8.RuneCountInString(keyword); l > 0 {
		db = db.Where("(weixin_user.nick_name like ? or weixin_user.phone like ?)", "%"+keyword+"%", "%"+keyword+"%")
	}

	if result := db.Group("weixin_user.id,weixin_user_commission.commission_type").Order("commission_amount DESC").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&userCommissionDatas); result.Error != nil {
		return nil, result.Error
	}

	for _, userCommissionData := range userCommissionDatas {
		list = append(list, &domain.UserCommission{
			ChildUserId:      userCommissionData.ChildUserId,
			ChildNickName:    userCommissionData.ChildNickName,
			ChildAvatarUrl:   userCommissionData.ChildAvatarUrl,
			ChildPhone:       userCommissionData.ChildPhone,
			CommissionType:   userCommissionData.CommissionType,
			TotalPayAmount:   userCommissionData.TotalPayAmount,
			CommissionPool:   userCommissionData.CommissionPool,
			CommissionAmount: userCommissionData.CommissionAmount,
		})
	}

	return list, nil
}

func (ucr *userCommissionRepo) ListByRelevanceId(ctx context.Context, childUserId, relevanceId uint64, commissionTypes []string) ([]*domain.UserCommission, error) {
	var userCommissions []UserCommission
	list := make([]*domain.UserCommission, 0)

	db := ucr.data.db.WithContext(ctx)

	if childUserId > 0 {
		db = db.Where("child_user_id = ?", childUserId)
	}

	if result := db.Where("relevance_id = ?", relevanceId).Where("commission_type in (?)", commissionTypes).Find(&userCommissions); result.Error != nil {
		return nil, result.Error
	}

	for _, userCommission := range userCommissions {
		list = append(list, userCommission.ToDomain(ctx))
	}

	return list, nil
}

func (ucr *userCommissionRepo) ListTask(ctx context.Context, commissionStatus, commissionType string) ([]*domain.UserCommission, error) {
	var userCommissions []UserCommission
	list := make([]*domain.UserCommission, 0)

	db := ucr.data.db.WithContext(ctx)

	if len(commissionStatus) > 0 {
		db = db.Where("commission_status = ?", commissionStatus)
	}

	if len(commissionType) > 0 {
		db = db.Where("commission_type = ?", commissionType)
	}

	if result := db.Find(&userCommissions); result.Error != nil {
		return nil, result.Error
	}

	for _, userCommission := range userCommissions {
		list = append(list, userCommission.ToDomain(ctx))
	}

	return list, nil
}

func (ucr *userCommissionRepo) ListBalance(ctx context.Context, pageNum, pageSize int, userId uint64, operationType uint8, keyword string) ([]*domain.UserCommission, error) {
	var userCommissionDatas []UserCommissionData
	list := make([]*domain.UserCommission, 0)

	db := ucr.data.db.WithContext(ctx).Table("weixin_user_commission").
		Where("user_id = ?", userId).
		Joins("left join weixin_user on weixin_user.id = weixin_user_commission.child_user_id").
		Select("weixin_user.nick_name as child_nick_name, weixin_user.phone as child_phone, weixin_user_commission.*")

	if operationType > 0 {
		db.Where("operation_type = ?", operationType)
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		db = db.Where("(weixin_user.nick_name like ? or weixin_user.phone like ?)", "%"+keyword+"%", "%"+keyword+"%")
	}

	if result := db.Order("create_time DESC").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&userCommissionDatas); result.Error != nil {
		return nil, result.Error
	}

	for _, userCommissionData := range userCommissionDatas {
		list = append(list, &domain.UserCommission{
			ChildUserId:      userCommissionData.ChildUserId,
			ChildNickName:    userCommissionData.ChildNickName,
			ChildAvatarUrl:   userCommissionData.ChildAvatarUrl,
			ChildLevel:       userCommissionData.ChildLevel,
			ChildPhone:       userCommissionData.ChildPhone,
			CommissionType:   userCommissionData.CommissionType,
			OperationType:    userCommissionData.OperationType,
			TotalPayAmount:   userCommissionData.TotalPayAmount,
			CommissionPool:   userCommissionData.CommissionPool,
			CommissionAmount: userCommissionData.CommissionAmount,
			BalanceStatus:    userCommissionData.BalanceStatus,
			OperationContent: userCommissionData.OperationContent,
			CreateTime:       userCommissionData.CreateTime,
		})
	}

	return list, nil
}

func (ucr *userCommissionRepo) ListCashable(ctx context.Context) ([]*domain.UserCommission, error) {
	var userCommissions []UserCommission
	list := make([]*domain.UserCommission, 0)

	if result := ucr.data.db.WithContext(ctx).
		Where("operation_type = 2").
		Where("commission_type in (7,8)").
		Where("balance_status = 0").Order("create_time ASC").
		Find(&userCommissions); result.Error != nil {
		return nil, result.Error
	}

	for _, userCommission := range userCommissions {
		list = append(list, userCommission.ToDomain(ctx))
	}

	return list, nil
}

func (ucr *userCommissionRepo) ListOperation(ctx context.Context) ([]*domain.UserCommission, error) {
	var userCommissions []UserCommission
	list := make([]*domain.UserCommission, 0)

	if result := ucr.data.db.WithContext(ctx).
		Where("commission_type=2 and commission_mcn_ratio=70 and operation_type=1").
		Where("create_time < ?", "2024-04-01 00:00:00").
		Find(&userCommissions); result.Error != nil {
		return nil, result.Error
	}

	for _, userCommission := range userCommissions {
		list = append(list, userCommission.ToDomain(ctx))
	}

	return list, nil
}

func (ucr *userCommissionRepo) Count(ctx context.Context, userId, organizationId uint64, childIds []uint64, commissionType uint8, startTime, endTime, keyword string) (int64, error) {
	var count int64

	db := ucr.data.db.WithContext(ctx).Table("weixin_user")

	if commissionType == 0 {
		if len(startTime) > 0 {
			db = db.Joins("left join weixin_user_commission on weixin_user_commission.child_user_id = weixin_user.id and weixin_user_commission.user_id = ? and weixin_user_commission.organization_id = ? and weixin_user_commission.commission_type in (1,2,4) and weixin_user_commission.create_time >= ? and weixin_user_commission.create_time <= ?", userId, organizationId, startTime+" 00:00:00", endTime+" 23:59:59")
		} else {
			db = db.Joins("left join weixin_user_commission on weixin_user_commission.child_user_id = weixin_user.id and weixin_user_commission.user_id = ? and weixin_user_commission.organization_id = ? and weixin_user_commission.commission_type in (1,2,4)", userId, organizationId)
		}
	} else {
		if len(startTime) > 0 {
			db = db.Joins("left join weixin_user_commission on weixin_user_commission.child_user_id = weixin_user.id and weixin_user_commission.user_id = ? and weixin_user_commission.organization_id = ? and weixin_user_commission.commission_type = ? and weixin_user_commission.create_time >= ? and weixin_user_commission.create_time <= ?", userId, organizationId, commissionType, startTime+" 00:00:00", endTime+" 23:59:59")
		} else {
			db = db.Joins("left join weixin_user_commission on weixin_user_commission.child_user_id = weixin_user.id and weixin_user_commission.user_id = ? and weixin_user_commission.organization_id = ? and weixin_user_commission.commission_type = ?", userId, organizationId, commissionType)
		}
	}

	db = db.Where("weixin_user.id in (?)", childIds).
		Select("weixin_user.nick_name as child_nick_name, weixin_user.avatar_url as child_avatar_url, weixin_user.phone as child_phone, weixin_user_commission.child_user_id, weixin_user_commission.commission_type, sum(weixin_user_commission.total_pay_amount) as total_pay_amount, sum(weixin_user_commission.commission_pool) as commission_pool, sum(weixin_user_commission.commission_amount) as commission_amount")

	if l := utf8.RuneCountInString(keyword); l > 0 {
		db = db.Where("(weixin_user.nick_name like ? or weixin_user.phone like ?)", "%"+keyword+"%", "%"+keyword+"%")
	}

	if result := db.Group("weixin_user.id,weixin_user_commission.commission_type").Model(&UserCommission{}).
		Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (ucr *userCommissionRepo) CountBalance(ctx context.Context, userId uint64, operationType uint8, keyword string) (int64, error) {
	var count int64

	db := ucr.data.db.WithContext(ctx).Table("weixin_user_commission").
		Where("user_id = ?", userId).
		Joins("left join weixin_user on weixin_user.id = weixin_user_commission.child_user_id")

	if operationType > 0 {
		db.Where("operation_type = ?", operationType)
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		db = db.Where("(weixin_user.nick_name like ? or weixin_user.phone like ?)", "%"+keyword+"%", "%"+keyword+"%")
	}

	if result := db.Model(&UserCommission{}).
		Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (ucr *userCommissionRepo) Statistics(ctx context.Context, userId uint64, commissionStatus, operationType uint8, commissionType, balanceStatus []uint8) (*domain.UserCommission, error) {
	userCommission := &UserCommission{}

	db := ucr.data.db.WithContext(ctx).
		Where("user_id = ?", userId)

	if commissionStatus > 0 {
		db = db.Where("commission_status = ?", commissionStatus)
	}

	if result := db.Where("commission_type in (?)", commissionType).
		Where("operation_type = ?", operationType).
		Where("balance_status in (?)", balanceStatus).
		Select("sum(total_pay_amount) as total_pay_amount, sum(commission_pool) as commission_pool, sum(commission_amount) as commission_amount").
		First(userCommission); result.Error != nil {
		return nil, result.Error
	}

	return userCommission.ToDomain(ctx), nil
}

func (ucr *userCommissionRepo) StatisticsDetail(ctx context.Context, userId, organizationId uint64, childIds []uint64, commissionType uint8, startDay, endDay, keyword string) (*domain.UserCommission, error) {
	userCommission := &UserCommission{}

	db := ucr.data.db.WithContext(ctx).Table("weixin_user").Where("weixin_user.id in (?)", childIds)

	if commissionType == 0 {
		if len(startDay) > 0 {
			db = db.Joins("left join weixin_user_commission on weixin_user_commission.child_user_id = weixin_user.id and weixin_user_commission.user_id = ? and weixin_user_commission.organization_id = ? and weixin_user_commission.commission_type in (1,2,4) and weixin_user_commission.create_time >= ? and weixin_user_commission.create_time <= ?", userId, organizationId, startDay+" 00:00:00", endDay+" 23:59:59")
		} else {
			db = db.Joins("left join weixin_user_commission on weixin_user_commission.child_user_id = weixin_user.id and weixin_user_commission.user_id = ? and weixin_user_commission.organization_id = ? and weixin_user_commission.commission_type in (1,2,4)", userId, organizationId)
		}
	} else {
		if len(startDay) > 0 {
			db = db.Joins("left join weixin_user_commission on weixin_user_commission.child_user_id = weixin_user.id and weixin_user_commission.user_id = ? and weixin_user_commission.organization_id = ? and weixin_user_commission.commission_type = ? and weixin_user_commission.create_time >= ? and weixin_user_commission.create_time <= ?", userId, organizationId, commissionType, startDay+" 00:00:00", endDay+" 23:59:59")
		} else {
			db = db.Joins("left join weixin_user_commission on weixin_user_commission.child_user_id = weixin_user.id and weixin_user_commission.user_id = ? and weixin_user_commission.organization_id = ? and weixin_user_commission.commission_type = ?", userId, organizationId, commissionType)
		}
	}

	db = db.Select("sum(weixin_user_commission.total_pay_amount) as total_pay_amount, sum(weixin_user_commission.commission_pool) as commission_pool, sum(weixin_user_commission.commission_amount) as commission_amount")

	if l := utf8.RuneCountInString(keyword); l > 0 {
		db = db.Where("(weixin_user.nick_name like ? or weixin_user.phone like ?)", "%"+keyword+"%", "%"+keyword+"%")
	}

	if result := db.First(userCommission); result.Error != nil {
		return nil, result.Error
	}

	return userCommission.ToDomain(ctx), nil
}

func (ucr *userCommissionRepo) Update(ctx context.Context, in *domain.UserCommission) (*domain.UserCommission, error) {
	userCommission := &UserCommission{
		Id:                 in.Id,
		UserId:             in.UserId,
		OrganizationId:     in.OrganizationId,
		ChildUserId:        in.ChildUserId,
		ChildLevel:         in.ChildLevel,
		Level:              in.Level,
		Relation:           in.Relation,
		CommissionStatus:   in.CommissionStatus,
		CommissionType:     in.CommissionType,
		RelevanceId:        in.RelevanceId,
		OperationType:      in.OperationType,
		TotalPayAmount:     in.TotalPayAmount,
		CommissionPool:     in.CommissionPool,
		CommissionMcnRatio: in.CommissionMcnRatio,
		CommissionAmount:   in.CommissionAmount,
		Name:               in.Name,
		IdentityCard:       in.IdentityCard,
		BankCode:           in.BankCode,
		BalanceStatus:      in.BalanceStatus,
		OutTradeNo:         in.OutTradeNo,
		InnerTradeNo:       in.InnerTradeNo,
		RealAmount:         in.RealAmount,
		ApplyTime:          in.ApplyTime,
		GongmallCreateTime: in.GongmallCreateTime,
		PayTime:            in.PayTime,
		OperationContent:   in.OperationContent,
		CreateTime:         in.CreateTime,
		UpdateTime:         in.UpdateTime,
	}

	if result := ucr.data.DB(ctx).Save(userCommission); result.Error != nil {
		return nil, result.Error
	}

	return userCommission.ToDomain(ctx), nil
}

func (ucr *userCommissionRepo) Save(ctx context.Context, in *domain.UserCommission) (*domain.UserCommission, error) {
	userCommission := &UserCommission{
		UserId:             in.UserId,
		OrganizationId:     in.OrganizationId,
		ChildUserId:        in.ChildUserId,
		ChildLevel:         in.ChildLevel,
		Level:              in.Level,
		Relation:           in.Relation,
		CommissionStatus:   in.CommissionStatus,
		CommissionType:     in.CommissionType,
		RelevanceId:        in.RelevanceId,
		OperationType:      in.OperationType,
		TotalPayAmount:     in.TotalPayAmount,
		CommissionPool:     in.CommissionPool,
		CommissionMcnRatio: in.CommissionMcnRatio,
		CommissionAmount:   in.CommissionAmount,
		Name:               in.Name,
		IdentityCard:       in.IdentityCard,
		BankCode:           in.BankCode,
		BalanceStatus:      in.BalanceStatus,
		OutTradeNo:         in.OutTradeNo,
		InnerTradeNo:       in.InnerTradeNo,
		RealAmount:         in.RealAmount,
		ApplyTime:          in.ApplyTime,
		GongmallCreateTime: in.GongmallCreateTime,
		PayTime:            in.PayTime,
		OperationContent:   in.OperationContent,
		CreateTime:         in.CreateTime,
		UpdateTime:         in.UpdateTime,
	}

	if result := ucr.data.DB(ctx).Create(userCommission); result.Error != nil {
		return nil, result.Error
	}

	return userCommission.ToDomain(ctx), nil
}

func (ucr *userCommissionRepo) DeleteByRelevanceId(ctx context.Context, relevanceId uint64, commissionType uint8) error {
	if result := ucr.data.DB(ctx).
		Where("relevance_id = ?", relevanceId).
		Where("commission_type = ?", commissionType).Delete(&UserCommission{}); result.Error != nil {
		return result.Error
	}

	return nil
}

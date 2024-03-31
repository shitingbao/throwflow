package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"weixin/internal/biz"
	"weixin/internal/domain"
)

// 微信小程序用户余额操作日志表
type UserBalanceLog struct {
	Id                 uint64     `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	UserId             uint64     `gorm:"column:user_id;type:bigint(20) UNSIGNED;not null;comment:微信小程序用户ID"`
	OrganizationId     uint64     `gorm:"column:organization_id;type:bigint(20) UNSIGNED;not null;comment:机构ID"`
	Name               string     `gorm:"column:name;type:text;not null;;comment:姓名"`
	IdentityCard       string     `gorm:"column:identity_card;type:text;not null;comment:身份证号码"`
	BankCode           string     `gorm:"column:bank_code;type:text;not null;comment:银行卡号"`
	Amount             float32    `gorm:"column:amount;type:decimal(10, 2) UNSIGNED;not null;comment:提现金额，单位元"`
	RelevanceId        uint64     `gorm:"column:relevance_id;type:bigint(20) UNSIGNED;not null;default:0;comment:金额加入关联的佣金ID"`
	BalanceType        uint8      `gorm:"column:balance_type;type:tinyint(3) UNSIGNED;not null;default:0;comment:余额说明,1：课程，2：电商, 3：成本购，4：其他，5，分佣提现，6成本购提现，在数据库里面细分，前端传递参数，1：成本购，2：电商，课程"`
	OperationType      uint8      `gorm:"column:operation_type;type:tinyint(3) UNSIGNED;not null;default:0;comment:1：加入，2：减少"`
	OperationContent   string     `gorm:"column:operation_content;type:text;not null;comment:说明"`
	BalanceStatus      uint8      `gorm:"column:balance_status;type:tinyint(3) UNSIGNED;not null;default:0;comment:余额状态，0：未处理，1：已处理，2：正在处理, 3:处理失败"`
	OutTradeNo         string     `gorm:"column:out_trade_no;type:varchar(32);not null;default:'';comment:提现订单号"`
	InnerTradeNo       string     `gorm:"column:inner_trade_no;type:varchar(100);not null;default:'';comment:工猫提现订单号"`
	RealAmount         float32    `gorm:"column:real_amount;type:decimal(10, 2) UNSIGNED;not null;default:0.00;comment:实发金额，单位元"`
	ApplyTime          *time.Time `gorm:"column:apply_time;type:datetime;null;default:null;comment:申请时间"`
	GongmallCreateTime *time.Time `gorm:"column:gongmall_create_time;type:datetime;null;default:null;comment:工猫进单时间"`
	PayTime            *time.Time `gorm:"column:pay_time;type:datetime;null;default:null;comment:工猫实际付款时间"`
	Day                uint32     `gorm:"column:day;type:int(11);not null;default:0;;comment:佣金产生时间"`
	CreateTime         time.Time  `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime         time.Time  `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (UserBalanceLog) TableName() string {
	return "weixin_user_balance_log"
}

type userBalanceLogRepo struct {
	data *Data
	log  *log.Helper
}

func (ubl *UserBalanceLog) ToDomain(ctx context.Context) *domain.UserBalanceLog {
	userBalanceLog := &domain.UserBalanceLog{
		Id:                 ubl.Id,
		UserId:             ubl.UserId,
		OrganizationId:     ubl.OrganizationId,
		Name:               ubl.Name,
		IdentityCard:       ubl.IdentityCard,
		BankCode:           ubl.BankCode,
		Amount:             ubl.Amount,
		RelevanceId:        ubl.RelevanceId,
		BalanceType:        ubl.BalanceType,
		OperationType:      ubl.OperationType,
		OperationContent:   ubl.OperationContent,
		BalanceStatus:      ubl.BalanceStatus,
		OutTradeNo:         ubl.OutTradeNo,
		InnerTradeNo:       ubl.InnerTradeNo,
		RealAmount:         ubl.RealAmount,
		ApplyTime:          ubl.ApplyTime,
		GongmallCreateTime: ubl.GongmallCreateTime,
		PayTime:            ubl.PayTime,
		Day:                ubl.Day,
		CreateTime:         ubl.CreateTime,
		UpdateTime:         ubl.UpdateTime,
	}

	return userBalanceLog
}

func NewUserBalanceLogRepo(data *Data, logger log.Logger) biz.UserBalanceLogRepo {
	return &userBalanceLogRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ublr *userBalanceLogRepo) NextId(ctx context.Context) (uint64, error) {
	return ublr.data.sonyflake.NextID()
}

func (ublr *userBalanceLogRepo) GetByOutTradeNo(ctx context.Context, outTradeNo string) (*domain.UserBalanceLog, error) {
	userBalanceLog := &UserBalanceLog{}

	if result := ublr.data.db.WithContext(ctx).
		Where("out_trade_no = ?", outTradeNo).
		First(userBalanceLog); result.Error != nil {
		return nil, result.Error
	}

	return userBalanceLog.ToDomain(ctx), nil
}

func (ublr *userBalanceLogRepo) List(ctx context.Context, pageNum, pageSize int, userId uint64, operationType uint8, orderBy string, balanceType, balanceStatus []string) ([]*domain.UserBalanceLog, error) {
	var userBalanceLogs []UserBalanceLog
	list := make([]*domain.UserBalanceLog, 0)

	db := ublr.data.db.WithContext(ctx)

	if userId > 0 {
		db = db.Where("user_id = ?", userId)
	}

	if len(balanceType) > 0 {
		db = db.Where("balance_type in ?", balanceType)
	}

	if operationType > 0 {
		db = db.Where("operation_type = ?", operationType)
	}

	if len(balanceStatus) > 0 {
		db = db.Where("balance_status in ?", balanceStatus)
	}

	if pageNum == 0 {
		if result := db.Order("create_time " + orderBy + ",id " + orderBy).
			Find(&userBalanceLogs); result.Error != nil {
			return nil, result.Error
		}
	} else {
		if result := db.Order("create_time " + orderBy + ",id " + orderBy).
			Limit(pageSize).Offset((pageNum - 1) * pageSize).
			Find(&userBalanceLogs); result.Error != nil {
			return nil, result.Error
		}
	}

	for _, userBalanceLog := range userBalanceLogs {
		list = append(list, userBalanceLog.ToDomain(ctx))
	}

	return list, nil
}

func (ublr *userBalanceLogRepo) Count(ctx context.Context, userId uint64, operationType uint8, balanceType, balanceStatus []string) (int64, error) {
	var count int64

	db := ublr.data.db.WithContext(ctx).Where("user_id = ?", userId)

	if len(balanceType) > 0 {
		db = db.Where("balance_type in ?", balanceType)
	}

	if operationType > 0 {
		db = db.Where("operation_type = ?", operationType)
	}

	if len(balanceStatus) > 0 {
		db = db.Where("balance_status in ?", balanceStatus)
	}

	if result := db.Model(&UserBalanceLog{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (ublr *userBalanceLogRepo) Statistics(ctx context.Context, userId uint64, operationType uint8, balanceType, balanceStatus []string) (*domain.UserBalanceLog, error) {
	userBalanceLog := &UserBalanceLog{}

	db := ublr.data.db.WithContext(ctx).
		Select("sum(amount) as amount").
		Where("user_id = ?", userId)

	if len(balanceType) > 0 {
		db = db.Where("balance_type in ?", balanceType)
	}

	if operationType > 0 {
		db = db.Where("operation_type = ?", operationType)
	}

	if len(balanceStatus) > 0 {
		db = db.Where("balance_status in ?", balanceStatus)
	}

	if result := db.First(userBalanceLog); result.Error != nil {
		return nil, result.Error
	}

	return userBalanceLog.ToDomain(ctx), nil
}

func (ublr *userBalanceLogRepo) Save(ctx context.Context, in *domain.UserBalanceLog) (*domain.UserBalanceLog, error) {
	userBalanceLog := &UserBalanceLog{
		UserId:             in.UserId,
		OrganizationId:     in.OrganizationId,
		Name:               in.Name,
		IdentityCard:       in.IdentityCard,
		BankCode:           in.BankCode,
		Amount:             in.Amount,
		RelevanceId:        in.RelevanceId,
		BalanceType:        in.BalanceType,
		OperationType:      in.OperationType,
		OperationContent:   in.OperationContent,
		BalanceStatus:      in.BalanceStatus,
		OutTradeNo:         in.OutTradeNo,
		InnerTradeNo:       in.InnerTradeNo,
		RealAmount:         in.RealAmount,
		ApplyTime:          in.ApplyTime,
		GongmallCreateTime: in.GongmallCreateTime,
		PayTime:            in.PayTime,
		Day:                in.Day,
		CreateTime:         in.CreateTime,
		UpdateTime:         in.UpdateTime,
	}

	if result := ublr.data.DB(ctx).Create(userBalanceLog); result.Error != nil {
		return nil, result.Error
	}

	return userBalanceLog.ToDomain(ctx), nil
}

func (ublr *userBalanceLogRepo) Update(ctx context.Context, in *domain.UserBalanceLog) (*domain.UserBalanceLog, error) {
	userBalanceLog := &UserBalanceLog{
		Id:                 in.Id,
		UserId:             in.UserId,
		OrganizationId:     in.OrganizationId,
		Name:               in.Name,
		IdentityCard:       in.IdentityCard,
		BankCode:           in.BankCode,
		Amount:             in.Amount,
		RelevanceId:        in.RelevanceId,
		BalanceType:        in.BalanceType,
		OperationType:      in.OperationType,
		OperationContent:   in.OperationContent,
		BalanceStatus:      in.BalanceStatus,
		OutTradeNo:         in.OutTradeNo,
		InnerTradeNo:       in.InnerTradeNo,
		RealAmount:         in.RealAmount,
		ApplyTime:          in.ApplyTime,
		GongmallCreateTime: in.GongmallCreateTime,
		PayTime:            in.PayTime,
		Day:                in.Day,
		CreateTime:         in.CreateTime,
		UpdateTime:         in.UpdateTime,
	}

	if result := ublr.data.DB(ctx).Save(userBalanceLog); result.Error != nil {
		return nil, result.Error
	}

	return userBalanceLog.ToDomain(ctx), nil
}

func (ublr *userBalanceLogRepo) DeleteByDay(ctx context.Context, operationType uint8, day uint32, balanceTypes []string) error {
	db := ublr.data.DB(ctx).Where("operation_type = ?", operationType)

	if day > 0 {
		db = db.Where("day = ?", day)
	}

	if len(balanceTypes) > 0 {
		db = db.Where("balance_type in ?", balanceTypes)
	}

	if result := db.Delete(&UserBalanceLog{}); result.Error != nil {
		return result.Error
	}

	return nil
}

func (ublr *userBalanceLogRepo) SaveCacheString(ctx context.Context, key string, val string, timeout time.Duration) (bool, error) {
	result, err := ublr.data.rdb.SetNX(ctx, key, val, timeout).Result()

	if err != nil {
		return false, err
	}

	return result, nil
}

func (ublr *userBalanceLogRepo) DeleteCache(ctx context.Context, key string) error {
	if _, err := ublr.data.rdb.Del(ctx, key).Result(); err != nil {
		return err
	}

	return nil
}

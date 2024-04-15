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

// 微信小程序用户机构课程订单表
type UserOrder struct {
	Id                  uint64     `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	UserId              uint64     `gorm:"column:user_id;type:bigint(20) UNSIGNED;not null;index:user_id_organization_id;comment:微信小程序用户ID"`
	OrganizationId      uint64     `gorm:"column:organization_id;type:bigint(20) UNSIGNED;not null;index:user_id_organization_id;comment:机构ID"`
	OrganizationUserId  uint64     `gorm:"column:organization_user_id;type:bigint(20) UNSIGNED;not null;default:0;comment:企业机构团长ID"`
	OrganizationTutorId uint64     `gorm:"column:organization_tutor_id;type:bigint(20) UNSIGNED;not null;default:0;comment:企业机构导师ID"`
	Level               uint8      `gorm:"column:level;type:tinyint(3) UNSIGNED;not null;default:0;comment:等级"`
	OutTradeNo          string     `gorm:"column:out_trade_no;type:varchar(32);not null;comment:订单号"`
	TransactionId       string     `gorm:"column:transaction_id;type:varchar(32);not null;comment:百变鱼系统交易号"`
	OutTransactionId    string     `gorm:"column:out_transaction_id;type:varchar(100);not null;comment:第三方订单号"`
	Amount              float32    `gorm:"column:amount;type:decimal(10, 2) UNSIGNED;not null;comment:订单金额，单位元"`
	PayAmount           float32    `gorm:"column:pay_amount;type:decimal(10, 2) UNSIGNED;not null;comment:订单支付金额，单位元"`
	PayTime             *time.Time `gorm:"column:pay_time;type:datetime;null;comment:支付时间"`
	PayStatus           uint8      `gorm:"column:pay_status;type:tinyint(3) UNSIGNED;not null;default:0;comment:状态：0：待付款，1：已付款"`
	OrderType           uint8      `gorm:"column:order_type;type:tinyint(3) UNSIGNED;not null;default:0;comment:状态：1：开通，2：升级"`
	CreateTime          time.Time  `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime          time.Time  `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (UserOrder) TableName() string {
	return "weixin_user_order"
}

type userOrderRepo struct {
	data *Data
	log  *log.Helper
}

func (uo *UserOrder) ToDomain(ctx context.Context) *domain.UserOrder {
	userOrder := &domain.UserOrder{
		Id:                  uo.Id,
		UserId:              uo.UserId,
		OrganizationId:      uo.OrganizationId,
		OrganizationUserId:  uo.OrganizationUserId,
		OrganizationTutorId: uo.OrganizationTutorId,
		Level:               uo.Level,
		OutTradeNo:          uo.OutTradeNo,
		TransactionId:       uo.TransactionId,
		OutTransactionId:    uo.OutTransactionId,
		Amount:              uo.Amount,
		PayAmount:           uo.PayAmount,
		PayTime:             uo.PayTime,
		PayStatus:           uo.PayStatus,
		OrderType:           uo.OrderType,
		CreateTime:          uo.CreateTime,
		UpdateTime:          uo.UpdateTime,
	}

	return userOrder
}

func NewUserOrderRepo(data *Data, logger log.Logger) biz.UserOrderRepo {
	return &userOrderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (uor *userOrderRepo) NextId(ctx context.Context) (uint64, error) {
	return uor.data.sonyflake.NextID()
}

func (uor *userOrderRepo) GetByOutTradeNo(ctx context.Context, outTradeNo string) (*domain.UserOrder, error) {
	userOrder := &UserOrder{}

	if result := uor.data.db.WithContext(ctx).Where("out_trade_no = ?", outTradeNo).First(userOrder); result.Error != nil {
		return nil, result.Error
	}

	return userOrder.ToDomain(ctx), nil
}

func (uor *userOrderRepo) GetByUserId(ctx context.Context, userId, organizationId uint64, payStatus string) (*domain.UserOrder, error) {
	userOrder := &UserOrder{}

	db := uor.data.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Where("organization_id = ?", organizationId)

	if len(payStatus) > 0 {
		db = db.Where("pay_status = ?", payStatus)
	}

	if result := db.Order("pay_time desc").First(userOrder); result.Error != nil {
		return nil, result.Error
	}

	return userOrder.ToDomain(ctx), nil
}

func (uor *userOrderRepo) List(ctx context.Context, pageNum, pageSize int) ([]*domain.UserOrder, error) {
	type userOrder struct {
		Phone     string     `gorm:"column:phone"`
		NickName  string     `gorm:"column:nick_name"`
		AvatarUrl string     `gorm:"column:avatar_url"`
		PayTime   *time.Time `gorm:"column:pay_time"`
	}

	var userOrders []userOrder
	list := make([]*domain.UserOrder, 0)

	db := uor.data.db.WithContext(ctx).Table("weixin_user_order").
		Joins("left join weixin_user on weixin_user_order.user_id=weixin_user.id ").
		Select("weixin_user.phone,weixin_user.nick_name,weixin_user.avatar_url,weixin_user_order.pay_time").
		Where("weixin_user_order.pay_status = ?", 1)

	if result := db.Order("weixin_user_order.pay_time DESC").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&userOrders); result.Error != nil {
		return nil, result.Error
	}

	for _, luserOrder := range userOrders {
		list = append(list, &domain.UserOrder{
			Phone:     luserOrder.Phone,
			NickName:  luserOrder.NickName,
			AvatarUrl: luserOrder.AvatarUrl,
			PayTime:   luserOrder.PayTime,
		})
	}

	return list, nil
}

func (uor *userOrderRepo) ListOperation(ctx context.Context) ([]*domain.UserOrder, error) {
	var userOrders []UserOrder
	list := make([]*domain.UserOrder, 0)

	if result := uor.data.db.WithContext(ctx).Where("pay_status = 1").
		Where("pay_time >= '2024-01-01'").
		Where("organization_id in (5,6,7)").
		Order("pay_time ASC").
		Find(&userOrders); result.Error != nil {
		return nil, result.Error
	}

	for _, userOrder := range userOrders {
		list = append(list, userOrder.ToDomain(ctx))
	}

	return list, nil
}

func (uor *userOrderRepo) Count(ctx context.Context) (int64, error) {
	var count int64

	db := uor.data.db.WithContext(ctx).Table("weixin_user_order").
		Joins("left join weixin_user on weixin_user_order.user_id=weixin_user.id ").
		Select("weixin_user.nick_name,weixin_user.avatar_url,weixin_user_order.pay_time").
		Where("weixin_user_order.pay_status = ?", 1)

	if result := db.Model(&UserOrder{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (uor *userOrderRepo) StatisticsPayAmount(ctx context.Context, userId, organizationId uint64, day string) (*domain.UserOrder, error) {
	userOrder := &UserOrder{}

	db := uor.data.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Where("organization_id = ?", organizationId).
		Where("pay_time >= ?", day+" 00:00:00").
		Where("pay_time <= ?", day+" 23:59:59").
		Where("pay_status = '1'").
		Select("sum(pay_amount) as pay_amount")

	if result := db.First(userOrder); result.Error != nil {
		return nil, result.Error
	}

	return userOrder.ToDomain(ctx), nil
}

func (uor *userOrderRepo) Update(ctx context.Context, in *domain.UserOrder) (*domain.UserOrder, error) {
	userOrder := &UserOrder{
		Id:                  in.Id,
		UserId:              in.UserId,
		OrganizationId:      in.OrganizationId,
		OrganizationUserId:  in.OrganizationUserId,
		OrganizationTutorId: in.OrganizationTutorId,
		Level:               in.Level,
		OutTradeNo:          in.OutTradeNo,
		TransactionId:       in.TransactionId,
		OutTransactionId:    in.OutTransactionId,
		Amount:              in.Amount,
		PayAmount:           in.PayAmount,
		PayTime:             in.PayTime,
		PayStatus:           in.PayStatus,
		OrderType:           in.OrderType,
		CreateTime:          in.CreateTime,
		UpdateTime:          in.UpdateTime,
	}

	if result := uor.data.DB(ctx).Save(userOrder); result.Error != nil {
		return nil, result.Error
	}

	return userOrder.ToDomain(ctx), nil
}

func (uor *userOrderRepo) Save(ctx context.Context, in *domain.UserOrder) (*domain.UserOrder, error) {
	userOrder := &UserOrder{
		UserId:              in.UserId,
		OrganizationId:      in.OrganizationId,
		OrganizationUserId:  in.OrganizationUserId,
		OrganizationTutorId: in.OrganizationTutorId,
		Level:               in.Level,
		OutTradeNo:          in.OutTradeNo,
		TransactionId:       in.TransactionId,
		OutTransactionId:    in.OutTransactionId,
		Amount:              in.Amount,
		PayAmount:           in.PayAmount,
		PayTime:             in.PayTime,
		PayStatus:           in.PayStatus,
		OrderType:           in.OrderType,
		CreateTime:          in.CreateTime,
		UpdateTime:          in.UpdateTime,
	}

	if result := uor.data.DB(ctx).Create(userOrder); result.Error != nil {
		return nil, result.Error
	}

	return userOrder.ToDomain(ctx), nil
}

func (uor *userOrderRepo) PutContent(ctx context.Context, fileName string, content io.Reader) (*ctos.PutObjectV2Output, error) {
	for _, ltos := range uor.data.toses {
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

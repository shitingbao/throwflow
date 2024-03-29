package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm/clause"
)

// 精选联盟抖客订单详情表
type DoukeOrderInfo struct {
	Id                  uint64     `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	UserId              uint64     `gorm:"column:user_id;type:bigint(20) UNSIGNED;not null;comment:微信小程序用户ID"`
	OrderId             string     `gorm:"column:order_id;type:varchar(100);not null;uniqueIndex:order_id;comment:订单ID"`
	ProductId           string     `gorm:"column:product_id;type:varchar(100);not null;comment:商品ID"`
	ProductName         string     `gorm:"column:product_name;type:varchar(100);not null;comment:商品名称"`
	ProductImg          string     `gorm:"column:product_img;type:varchar(250);not null;comment:商品图片URL"`
	PaySuccessTime      time.Time  `gorm:"column:pay_success_time;type:datetime;not null;comment:付款时间"`
	SettleTime          *time.Time `gorm:"column:settle_time;type:datetime;null;default:null;comment:结算时间，结算前为空字符串"`
	TotalPayAmount      float32    `gorm:"column:total_pay_amount;type:decimal(10, 2) UNSIGNED;not null;comment:订单支付金额，单位元"`
	PayGoodsAmount      float32    `gorm:"column:pay_goods_amount;type:decimal(10, 2) UNSIGNED;not null;comment:预估参与结算金额，单位元"`
	FlowPoint           string     `gorm:"column:flow_point;type:varchar(50);not null;comment:订单状态(PAY_SUCC:支付完成 REFUND:退款 SETTLE:结算 CONFIRM: 确认收货)"`
	EstimatedCommission float32    `gorm:"column:estimated_commission;type:decimal(10, 2) UNSIGNED;not null;comment:渠道预估推广费，单位元"`
	RealCommission      float32    `gorm:"column:real_commission;type:decimal(10, 2) UNSIGNED;not null;comment:渠道实际推广费，单位元"`
	ItemNum             uint64     `gorm:"column:item_num;type:bigint(20) UNSIGNED;not null;comment:商品数目"`
	CreateTime          time.Time  `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime          time.Time  `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (DoukeOrderInfo) TableName() string {
	return "douyin_douke_order_info"
}

type doukeOrderInfoRepo struct {
	data *Data
	log  *log.Helper
}

func (doi *DoukeOrderInfo) ToDomain() *domain.DoukeOrderInfo {
	return &domain.DoukeOrderInfo{
		Id:                  doi.Id,
		UserId:              doi.UserId,
		OrderId:             doi.OrderId,
		ProductId:           doi.ProductId,
		ProductName:         doi.ProductName,
		ProductImg:          doi.ProductImg,
		PaySuccessTime:      doi.PaySuccessTime,
		SettleTime:          doi.SettleTime,
		TotalPayAmount:      doi.TotalPayAmount,
		PayGoodsAmount:      doi.PayGoodsAmount,
		FlowPoint:           doi.FlowPoint,
		EstimatedCommission: doi.EstimatedCommission,
		RealCommission:      doi.RealCommission,
		ItemNum:             doi.ItemNum,
		CreateTime:          doi.CreateTime,
		UpdateTime:          doi.UpdateTime,
	}
}

func NewDoukeOrderInfoRepo(data *Data, logger log.Logger) biz.DoukeOrderInfoRepo {
	return &doukeOrderInfoRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (doir *doukeOrderInfoRepo) List(ctx context.Context) ([]*domain.DoukeOrderInfo, error) {
	var doukeOrderInfos []DoukeOrderInfo
	list := make([]*domain.DoukeOrderInfo, 0)

	db := doir.data.db.WithContext(ctx).Where("flow_point = 'SETTLE'")

	if result := db.Find(&doukeOrderInfos); result.Error != nil {
		return nil, result.Error
	}

	for _, doukeOrderInfo := range doukeOrderInfos {
		list = append(list, doukeOrderInfo.ToDomain())
	}

	return list, nil
}

func (doir *doukeOrderInfoRepo) Statistics(ctx context.Context, userId uint64, startDay, endDay, flowPoint string) (*domain.DoukeOrderInfo, error) {
	doukeOrderInfo := &DoukeOrderInfo{}

	db := doir.data.db.WithContext(ctx).
		Select("sum(item_num) as item_num, sum(total_pay_amount) as total_pay_amount, sum(pay_goods_amount) as pay_goods_amount, sum(estimated_commission) as estimated_commission, sum(real_commission) as real_commission")

	if len(startDay) > 0 {
		db = db.Where("pay_success_time >= ?", startDay+" 00:00:00")
	}

	if len(endDay) > 0 {
		db = db.Where("pay_success_time <= ?", endDay+" 23:59:59")
	}

	if flowPoint == "refund" {
		db = db.Where("flow_point = 'REFUND'")
	} else {
		db = db.Where("flow_point != 'REFUND'")
	}

	if userId > 0 {
		db = db.Where("user_id = ?", userId)
	}

	if result := db.First(doukeOrderInfo); result.Error != nil {
		return nil, result.Error
	}

	return doukeOrderInfo.ToDomain(), nil
}

func (doir *doukeOrderInfoRepo) StatisticsRealcommission(ctx context.Context, userId uint64, startDay, endDay string) (*domain.DoukeOrderInfo, error) {
	doukeOrderInfo := &DoukeOrderInfo{}

	db := doir.data.db.WithContext(ctx).
		Where("flow_point = 'SETTLE'").
		Select("sum(real_commission) as real_commission")

	if len(startDay) > 0 {
		db = db.Where("settle_time >= ?", startDay+" 00:00:00")
	}

	if len(endDay) > 0 {
		db = db.Where("settle_time <= ?", endDay+" 23:59:59")
	} else {
		db = db.Where("settle_time <= ?", startDay+" 23:59:59")
	}

	if userId > 0 {
		db = db.Where("user_id = ?", userId)
	}

	if result := db.First(doukeOrderInfo); result.Error != nil {
		return nil, result.Error
	}

	return doukeOrderInfo.ToDomain(), nil
}

func (doir *doukeOrderInfoRepo) Upsert(ctx context.Context, in *domain.DoukeOrderInfo) error {
	if result := doir.data.DB(ctx).Table("douyin_douke_order_info").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "order_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"user_id", "product_id", "product_name", "product_img", "pay_success_time", "settle_time", "total_pay_amount", "pay_goods_amount", "flow_point", "estimated_commission", "real_commission", "item_num", "update_time"}),
	}).Create(&in); result.Error != nil {
		return result.Error
	}

	return nil
}

func (doir *doukeOrderInfoRepo) GetByUserIdAndProductId(ctx context.Context, userId uint64, productId, flowPoint, createTime string) (*domain.DoukeOrderInfo, error) {
	doukeOrderInfo := &DoukeOrderInfo{}

	db := doir.data.DB(ctx).Table("douyin_douke_order_info").
		Where("product_id = ? and user_id = ? and create_time > ?", productId, userId, createTime)

	if len(flowPoint) > 0 {
		db = db.Where("flow_point = ?", flowPoint)
	}

	if err := db.First(&doukeOrderInfo).Error; err != nil {
		return nil, err
	}

	return doukeOrderInfo.ToDomain(), nil
}
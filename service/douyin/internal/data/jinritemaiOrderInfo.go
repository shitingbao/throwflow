package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"douyin/internal/pkg/tool"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm/clause"
	"strconv"
	"strings"
	"time"
)

type JinritemaiOrderInfo struct {
	Id                  uint64     `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	ClientKey           string     `gorm:"column:client_key;type:varchar(50);index:client_key_open_id;not null;comment:应用Client Key"`
	OpenId              string     `gorm:"column:open_id;type:varchar(100);index:client_key_open_id;not null;comment:授权用户唯一标识"`
	BuyinId             string     `gorm:"column:buyin_id;type:varchar(100);not null;comment:百应ID"`
	OrderId             string     `gorm:"column:order_id;type:varchar(100);not null;uniqueIndex:order_id;comment:订单ID"`
	ProductId           string     `gorm:"column:product_id;type:varchar(100);not null;comment:商品ID"`
	ProductName         string     `gorm:"column:product_name;type:varchar(100);not null;comment:商品名称"`
	ProductImg          string     `gorm:"column:product_img;type:varchar(250);not null;comment:商品图片URL"`
	CommissionRate      uint8      `gorm:"column:commission_rate;type:tinyint(3) UNSIGNED;not null;default:0;comment:佣金比例"`
	PaySuccessTime      time.Time  `gorm:"column:pay_success_time;type:datetime;not null;comment:付款时间"`
	SettleTime          *time.Time `gorm:"column:settle_time;type:datetime;null;default:null;comment:结算时间，结算前为空字符串"`
	TotalPayAmount      float32    `gorm:"column:total_pay_amount;type:decimal(10, 2) UNSIGNED;not null;comment:订单支付金额，单位元"`
	PayGoodsAmount      float32    `gorm:"column:pay_goods_amount;type:decimal(10, 2) UNSIGNED;not null;comment:预估参与结算金额，单位元"`
	FlowPoint           string     `gorm:"column:flow_point;type:varchar(50);not null;comment:订单状态(PAY_SUCC:支付完成 REFUND:退款 SETTLE:结算 CONFIRM: 确认收货)"`
	EstimatedCommission float32    `gorm:"column:estimated_commission;type:decimal(10, 2) UNSIGNED;not null;comment:达人预估佣金收入，单位元"`
	RealCommission      float32    `gorm:"column:real_commission;type:decimal(10, 2) UNSIGNED;not null;comment:达人实际佣金收入，单位元"`
	ItemNum             uint64     `gorm:"column:item_num;type:bigint(20) UNSIGNED;not null;comment:商品数目"`
	PickExtra           string     `gorm:"column:pick_extra;type:varchar(150);not null;comment:选品来源自定义参数"`
	MediaType           string     `gorm:"column:media_type;type:varchar(50);not null;comment:带货体裁。shop_list：橱窗；video：视频；live：直播；others：其他(如图文、微头条、问答、西瓜长视频等)"`
	MediaId             uint64     `gorm:"column:media_id;type:bigint(20) UNSIGNED;not null;index:media_id;comment:带货体裁ID"`
	CreateTime          time.Time  `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime          time.Time  `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (JinritemaiOrderInfo) TableName() string {
	return "douyin_jinritemai_order_info"
}

type jinritemaiOrderInfoRepo struct {
	data *Data
	log  *log.Helper
}

func (joi *JinritemaiOrderInfo) ToDomain() *domain.JinritemaiOrderInfo {
	return &domain.JinritemaiOrderInfo{
		Id:                  joi.Id,
		ClientKey:           joi.ClientKey,
		OpenId:              joi.OpenId,
		BuyinId:             joi.BuyinId,
		OrderId:             joi.OrderId,
		ProductId:           joi.ProductId,
		ProductName:         joi.ProductName,
		ProductImg:          joi.ProductImg,
		CommissionRate:      joi.CommissionRate,
		PaySuccessTime:      joi.PaySuccessTime,
		SettleTime:          joi.SettleTime,
		TotalPayAmount:      joi.TotalPayAmount,
		PayGoodsAmount:      joi.PayGoodsAmount,
		FlowPoint:           joi.FlowPoint,
		EstimatedCommission: joi.EstimatedCommission,
		RealCommission:      joi.RealCommission,
		ItemNum:             joi.ItemNum,
		PickExtra:           joi.PickExtra,
		MediaType:           joi.MediaType,
		MediaId:             joi.MediaId,
		CreateTime:          joi.CreateTime,
		UpdateTime:          joi.UpdateTime,
	}
}

func NewJinritemaiOrderInfoRepo(data *Data, logger log.Logger) biz.JinritemaiOrderInfoRepo {
	return &jinritemaiOrderInfoRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (joir *jinritemaiOrderInfoRepo) GetByClientKeyAndOpenId(ctx context.Context, clientKey, openId string) (*domain.JinritemaiOrderInfo, error) {
	jinritemaiOrderInfo := &JinritemaiOrderInfo{}

	if result := joir.data.db.WithContext(ctx).
		Where("client_key = ?", clientKey).
		Where("open_id = ?", openId).
		Order("pay_success_time desc").
		First(jinritemaiOrderInfo); result.Error != nil {
		return nil, result.Error
	}

	return jinritemaiOrderInfo.ToDomain(), nil
}

func (joir *jinritemaiOrderInfoRepo) GetByClientKeyAndOpenIdAndMediaTypeAndMediaId(ctx context.Context, clientKey, openId, mediaType, mediaId string) (*domain.JinritemaiOrderInfo, error) {
	jinritemaiOrderInfo := &JinritemaiOrderInfo{}

	if result := joir.data.db.WithContext(ctx).
		Where("client_key = ?", clientKey).
		Where("open_id = ?", openId).
		Where("media_type = ?", mediaType).
		Where("media_id = ?", mediaId).
		Order("pay_success_time desc").
		First(jinritemaiOrderInfo); result.Error != nil {
		return nil, result.Error
	}

	return jinritemaiOrderInfo.ToDomain(), nil
}

func (joir *jinritemaiOrderInfoRepo) GetIsTopByProductId(ctx context.Context, productId uint64) (*domain.JinritemaiOrderInfo, error) {
	jinritemaiOrderInfo := &JinritemaiOrderInfo{}

	if result := joir.data.db.WithContext(ctx).
		Where("product_id = ?", productId).
		Where("media_type = 'video'").
		Where("flow_point != 'REFUND'").
		Where("update_time >= ?", tool.TimeToString("2006-01-02", time.Now().AddDate(0, 0, -7))+" 00:00:00").
		Select("sum(item_num) as item_num").
		Group("media_id").
		Order("item_num DESC").
		Having("item_num > 100").
		Take(jinritemaiOrderInfo); result.Error != nil {
		return nil, result.Error
	}

	return jinritemaiOrderInfo.ToDomain(), nil
}

func (joir *jinritemaiOrderInfoRepo) List(ctx context.Context, pageNum, pageSize int, openDouyinTokens []*domain.OpenDouyinToken, startDay, endDay string) ([]*domain.JinritemaiOrderInfo, error) {
	var jinritemaiOrderInfos []JinritemaiOrderInfo
	list := make([]*domain.JinritemaiOrderInfo, 0)

	db := joir.data.db.WithContext(ctx).
		Where("pay_success_time >= ?", startDay+" 00:00:00").
		Where("pay_success_time <= ?", endDay+" 23:59:59").
		Where("flow_point != 'REFUND'")

	openDouyinTokenSqls := make([]string, 0)

	for _, openDouyinToken := range openDouyinTokens {
		openDouyinTokenSqls = append(openDouyinTokenSqls, "(client_key = '"+openDouyinToken.ClientKey+"' and open_id = '"+openDouyinToken.OpenId+"')")
	}

	db = db.Where(strings.Join(openDouyinTokenSqls, " or ")).
		Select("client_key, open_id, media_id, media_type, product_id, max(product_name) as product_name, max(product_img) as product_img, sum(total_pay_amount) as total_pay_amount, sum(pay_goods_amount) as pay_goods_amount, sum(estimated_commission) as estimated_commission, sum(real_commission) as real_commission, sum(item_num) as item_num, max(pay_success_time) as pay_success_time")

	if result := db.Group("client_key,open_id,media_id,media_type,product_id").
		Order("pay_success_time desc,product_id desc").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&jinritemaiOrderInfos); result.Error != nil {
		return nil, result.Error
	}

	for _, jinritemaiOrderInfo := range jinritemaiOrderInfos {
		list = append(list, jinritemaiOrderInfo.ToDomain())
	}

	return list, nil
}

func (joir *jinritemaiOrderInfoRepo) ListOperation(ctx context.Context, pageNum, pageSize int) ([]*domain.JinritemaiOrderInfo, error) {
	var jinritemaiOrderInfos []JinritemaiOrderInfo
	list := make([]*domain.JinritemaiOrderInfo, 0)

	if result := joir.data.db.WithContext(ctx).
		Order("id asc").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&jinritemaiOrderInfos); result.Error != nil {
		return nil, result.Error
	}

	for _, jinritemaiOrderInfo := range jinritemaiOrderInfos {
		list = append(list, jinritemaiOrderInfo.ToDomain())
	}

	return list, nil
}

func (joir *jinritemaiOrderInfoRepo) ListProductByClientKeyAndOpenIdAndMediaType(ctx context.Context, clientKey, openId, mediaType string) ([]*domain.JinritemaiOrderInfo, error) {
	var jinritemaiOrderInfos []JinritemaiOrderInfo
	list := make([]*domain.JinritemaiOrderInfo, 0)

	if result := joir.data.db.WithContext(ctx).
		Where("client_key = ?", clientKey).
		Where("open_id = ?", openId).
		Where("media_type = ?", mediaType).
		Select("media_id,max(product_id) as product_id,max(product_name) as product_name,max(product_img) as product_img").
		Group("media_id").
		Find(&jinritemaiOrderInfos); result.Error != nil {
		return nil, result.Error
	}

	for _, jinritemaiOrderInfo := range jinritemaiOrderInfos {
		list = append(list, jinritemaiOrderInfo.ToDomain())
	}

	return list, nil
}

func (joir *jinritemaiOrderInfoRepo) ListByProductIds(ctx context.Context, openDouyinTokens []*domain.OpenDouyinToken, productIds []string) ([]*domain.JinritemaiOrderInfo, error) {
	var jinritemaiOrderInfos []JinritemaiOrderInfo
	list := make([]*domain.JinritemaiOrderInfo, 0)

	db := joir.data.db.WithContext(ctx).Where("flow_point != 'REFUND'").
		Where("product_id in ?", productIds).
		Where("media_type = 'video'")

	openDouyinTokenSqls := make([]string, 0)

	for _, openDouyinToken := range openDouyinTokens {
		openDouyinTokenSqls = append(openDouyinTokenSqls, "(client_key = '"+openDouyinToken.ClientKey+"' and open_id = '"+openDouyinToken.OpenId+"')")
	}

	db = db.Where(strings.Join(openDouyinTokenSqls, " or ")).
		Select("media_id, sum(total_pay_amount) as total_pay_amount")

	if result := db.Group("product_id,media_id").
		Find(&jinritemaiOrderInfos); result.Error != nil {
		return nil, result.Error
	}

	for _, jinritemaiOrderInfo := range jinritemaiOrderInfos {
		list = append(list, jinritemaiOrderInfo.ToDomain())
	}

	return list, nil
}

func (joir *jinritemaiOrderInfoRepo) ListNotRefund(ctx context.Context) ([]*domain.JinritemaiOrderInfo, error) {
	var jinritemaiOrderInfos []JinritemaiOrderInfo
	list := make([]*domain.JinritemaiOrderInfo, 0)

	db := joir.data.db.WithContext(ctx).
		Where("flow_point != 'REFUND'").
		Where("length(pick_extra) > 0").
		Select("pick_extra, client_key, open_id, media_id, media_type, max(media_title) media_title, min(media_create_time) media_create_time, product_id, max(product_name) product_name, max(product_img) product_img, max(product_price) product_price, max(industry_id) industry_id, max(industry_name) industry_name, max(category_id) category_id, max(category_name) category_name, max(product_create_time) product_create_time, sum(total_pay_amount) as total_pay_amount, sum(pay_goods_amount) as pay_goods_amount, sum(estimated_commission) as estimated_commission, sum(real_commission) as real_commission,sum(estimated_service) as estimated_service, sum(real_service) as real_service, sum(item_num) as item_num, min(pay_success_time) as pay_success_time")

	if result := db.Group("pick_extra,client_key,open_id,media_id,media_type,product_id").
		Find(&jinritemaiOrderInfos); result.Error != nil {
		return nil, result.Error
	}

	for _, jinritemaiOrderInfo := range jinritemaiOrderInfos {
		list = append(list, jinritemaiOrderInfo.ToDomain())
	}

	return list, nil
}

func (joir *jinritemaiOrderInfoRepo) ListRefund(ctx context.Context) ([]*domain.JinritemaiOrderInfo, error) {
	var jinritemaiOrderInfos []JinritemaiOrderInfo
	list := make([]*domain.JinritemaiOrderInfo, 0)

	db := joir.data.db.WithContext(ctx).
		Where("flow_point = 'REFUND'").
		Where("length(pick_extra) > 0").
		Select("pick_extra, client_key, open_id, media_id, media_type, max(media_title) media_title, min(media_create_time) media_create_time, product_id, max(product_name) product_name, max(product_img) product_img, max(product_price) product_price, max(industry_id) industry_id, max(industry_name) industry_name, max(category_id) category_id, max(category_name) category_name, max(product_create_time) product_create_time, sum(total_pay_amount) as total_pay_amount, sum(pay_goods_amount) as pay_goods_amount, sum(estimated_commission) as estimated_commission, sum(real_commission) as real_commission,sum(estimated_service) as estimated_service, sum(real_service) as real_service, sum(item_num) as item_num, min(pay_success_time) as pay_success_time")

	if result := db.Group("pick_extra,client_key,open_id,media_id,media_type,product_id").
		Find(&jinritemaiOrderInfos); result.Error != nil {
		return nil, result.Error
	}

	for _, jinritemaiOrderInfo := range jinritemaiOrderInfos {
		list = append(list, jinritemaiOrderInfo.ToDomain())
	}

	return list, nil
}

func (joir *jinritemaiOrderInfoRepo) ListMediaId(ctx context.Context) ([]*domain.JinritemaiOrderInfo, error) {
	var jinritemaiOrderInfos []JinritemaiOrderInfo
	list := make([]*domain.JinritemaiOrderInfo, 0)

	db := joir.data.db.WithContext(ctx).
		Where("media_type = 'video'").
		Where("media_title = ''").
		Select("distinct(media_id) as media_id")

	if result := db.Find(&jinritemaiOrderInfos); result.Error != nil {
		return nil, result.Error
	}

	for _, jinritemaiOrderInfo := range jinritemaiOrderInfos {
		list = append(list, jinritemaiOrderInfo.ToDomain())
	}

	return list, nil
}

func (joir *jinritemaiOrderInfoRepo) ListProductId(ctx context.Context) ([]*domain.JinritemaiOrderInfo, error) {
	var jinritemaiOrderInfos []JinritemaiOrderInfo
	list := make([]*domain.JinritemaiOrderInfo, 0)

	db := joir.data.db.WithContext(ctx).
		Where("industry_id = 0").
		Select("product_id").
		Group("product_id")

	if result := db.Find(&jinritemaiOrderInfos); result.Error != nil {
		return nil, result.Error
	}

	for _, jinritemaiOrderInfo := range jinritemaiOrderInfos {
		list = append(list, jinritemaiOrderInfo.ToDomain())
	}

	return list, nil
}

func (joir *jinritemaiOrderInfoRepo) ListByPickExtra(ctx context.Context, isServiceRate uint8) ([]*domain.JinritemaiOrderInfo, error) {
	var jinritemaiOrderInfos []JinritemaiOrderInfo
	list := make([]*domain.JinritemaiOrderInfo, 0)

	db := joir.data.db.WithContext(ctx).Where("pick_extra != ''")

	if isServiceRate == 1 {
		db = db.Where("service_rate > 0").
			Where("flow_point = 'SETTLE'").
			Where("real_service = 0")
	} else {
		db = db.Where("service_rate = 0")
	}

	if result := db.Find(&jinritemaiOrderInfos); result.Error != nil {
		return nil, result.Error
	}

	for _, jinritemaiOrderInfo := range jinritemaiOrderInfos {
		list = append(list, jinritemaiOrderInfo.ToDomain())
	}

	return list, nil
}

func (joir *jinritemaiOrderInfoRepo) ListByProductIdAndMediaIds(ctx context.Context, commissionRateJinritemaiOrders []*domain.CommissionRateJinritemaiOrder) ([]*domain.JinritemaiOrderInfo, error) {
	var jinritemaiOrderInfos []JinritemaiOrderInfo
	list := make([]*domain.JinritemaiOrderInfo, 0)

	db := joir.data.db.WithContext(ctx).Where("media_type = 'video'")

	commissionRateJinritemaiOrderSqls := make([]string, 0)

	for _, commissionRateJinritemaiOrder := range commissionRateJinritemaiOrders {
		commissionRateJinritemaiOrderSqls = append(commissionRateJinritemaiOrderSqls, "(product_id = '"+strconv.FormatUint(commissionRateJinritemaiOrder.ProductId, 10)+"' and media_id = '"+strconv.FormatUint(commissionRateJinritemaiOrder.VideoId, 10)+"')")
	}

	db = db.Where(strings.Join(commissionRateJinritemaiOrderSqls, " or ")).
		Select("product_id, media_id, commission_rate")

	if result := db.Group("product_id,media_id,commission_rate").
		Find(&jinritemaiOrderInfos); result.Error != nil {
		return nil, result.Error
	}

	for _, jinritemaiOrderInfo := range jinritemaiOrderInfos {
		list = append(list, jinritemaiOrderInfo.ToDomain())
	}

	return list, nil
}

func (joir *jinritemaiOrderInfoRepo) ListProductIdAndMediaIds(ctx context.Context, pageNum, pageSize int) ([]*domain.JinritemaiOrderInfo, error) {
	var jinritemaiOrderInfos []JinritemaiOrderInfo
	list := make([]*domain.JinritemaiOrderInfo, 0)

	db := joir.data.db.WithContext(ctx).
		Where("media_type = 'video'").
		Select("media_id,product_id,sum(item_num) item_num")

	if result := db.Group("media_id,product_id").
		Having("item_num >= ?", 10).
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&jinritemaiOrderInfos); result.Error != nil {
		return nil, result.Error
	}

	for _, jinritemaiOrderInfo := range jinritemaiOrderInfos {
		list = append(list, jinritemaiOrderInfo.ToDomain())
	}

	return list, nil
}

func (joir *jinritemaiOrderInfoRepo) Count(ctx context.Context, openDouyinTokens []*domain.OpenDouyinToken, startDay, endDay string) (int64, error) {
	var count int64

	db := joir.data.db.WithContext(ctx).
		Where("pay_success_time >= ?", startDay+" 00:00:00").
		Where("pay_success_time <= ?", endDay+" 23:59:59").
		Where("flow_point != 'REFUND'")

	openDouyinTokenSqls := make([]string, 0)

	for _, openDouyinToken := range openDouyinTokens {
		openDouyinTokenSqls = append(openDouyinTokenSqls, "(client_key = '"+openDouyinToken.ClientKey+"' and open_id = '"+openDouyinToken.OpenId+"')")
	}

	db = db.Where(strings.Join(openDouyinTokenSqls, " or "))

	if result := joir.data.db.WithContext(ctx).Model(&JinritemaiOrderInfo{}).
		Table("(?) as u", db.Model(&JinritemaiOrderInfo{}).Select("client_key").Group("client_key,open_id,media_id,media_type,product_id")).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (joir *jinritemaiOrderInfoRepo) CountOperation(ctx context.Context) (int64, error) {
	var count int64

	if result := joir.data.db.WithContext(ctx).Model(&JinritemaiOrderInfo{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (joir *jinritemaiOrderInfoRepo) CountProductIdAndMediaIds(ctx context.Context) (int64, error) {
	var count int64

	if result := joir.data.db.WithContext(ctx).Table("(?) as joi", joir.data.db.WithContext(ctx).
		Model(&JinritemaiOrderInfo{}).
		Where("media_type = 'video'").
		Select("media_id,product_id,sum(item_num) item_num").
		Group("media_id,product_id").
		Having("item_num >= ?", 10)).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	/*db := joir.data.db.WithContext(ctx).
		Where("media_type = 'video'").
		Select("media_id,product_id,sum(item_num) item_num").
		Group("media_id,product_id").
		Having("item_num >= ?", 10)

	if result := db.Model(&JinritemaiOrderInfo{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}*/

	return count, nil
}

func (joir *jinritemaiOrderInfoRepo) Statistics(ctx context.Context, openDouyinTokens []*domain.OpenDouyinToken, startDay, endDay, flowPoint, pickExtra string) (*domain.JinritemaiOrderInfo, error) {
	jinritemaiOrderInfo := &JinritemaiOrderInfo{}

	db := joir.data.db.WithContext(ctx).
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

	if len(pickExtra) > 0 {
		db = db.Where("pick_extra = ?", pickExtra)
	}

	openDouyinTokenSqls := make([]string, 0)

	for _, openDouyinToken := range openDouyinTokens {
		openDouyinTokenSqls = append(openDouyinTokenSqls, "(client_key = '"+openDouyinToken.ClientKey+"' and open_id = '"+openDouyinToken.OpenId+"')")
	}

	db = db.Where(strings.Join(openDouyinTokenSqls, " or "))

	if result := db.First(jinritemaiOrderInfo); result.Error != nil {
		return nil, result.Error
	}

	return jinritemaiOrderInfo.ToDomain(), nil
}

func (joir *jinritemaiOrderInfoRepo) StatisticsRealcommission(ctx context.Context, openDouyinTokens []*domain.OpenDouyinToken, startDay, endDay, pickExtra string) (*domain.JinritemaiOrderInfo, error) {
	jinritemaiOrderInfo := &JinritemaiOrderInfo{}

	db := joir.data.db.WithContext(ctx).
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

	if len(pickExtra) > 0 {
		db = db.Where("pick_extra = ?", pickExtra)
	}

	openDouyinTokenSqls := make([]string, 0)

	for _, openDouyinToken := range openDouyinTokens {
		openDouyinTokenSqls = append(openDouyinTokenSqls, "(client_key = '"+openDouyinToken.ClientKey+"' and open_id = '"+openDouyinToken.OpenId+"')")
	}

	db = db.Where(strings.Join(openDouyinTokenSqls, " or "))

	if result := db.First(jinritemaiOrderInfo); result.Error != nil {
		return nil, result.Error
	}

	return jinritemaiOrderInfo.ToDomain(), nil
}

func (joir *jinritemaiOrderInfoRepo) StatisticsRealcommissionPayTime(ctx context.Context, openDouyinTokens []*domain.OpenDouyinToken, payTime, startDay, endDay, pickExtra string) (*domain.JinritemaiOrderInfo, error) {
	jinritemaiOrderInfo := &JinritemaiOrderInfo{}

	db := joir.data.db.WithContext(ctx).
		Where("flow_point = 'SETTLE'").
		Select("sum(real_commission) as real_commission")

	if len(payTime) > 0 {
		db = db.Where("pay_success_time > ?", payTime+" 23:59:59")
	}

	if len(startDay) > 0 {
		db = db.Where("settle_time >= ?", startDay+" 00:00:00")
	}

	if len(endDay) > 0 {
		db = db.Where("settle_time <= ?", endDay+" 23:59:59")
	} else {
		db = db.Where("settle_time <= ?", startDay+" 23:59:59")
	}

	if len(pickExtra) > 0 {
		db = db.Where("pick_extra = ?", pickExtra)
	}

	openDouyinTokenSqls := make([]string, 0)

	for _, openDouyinToken := range openDouyinTokens {
		openDouyinTokenSqls = append(openDouyinTokenSqls, "(client_key = '"+openDouyinToken.ClientKey+"' and open_id = '"+openDouyinToken.OpenId+"')")
	}

	db = db.Where(strings.Join(openDouyinTokenSqls, " or "))

	if result := db.First(jinritemaiOrderInfo); result.Error != nil {
		return nil, result.Error
	}

	return jinritemaiOrderInfo.ToDomain(), nil
}

func (joir *jinritemaiOrderInfoRepo) StatisticsAwemeIndustry(ctx context.Context, companyId uint64, startDay, endDay string, openDouyinUserInfos []*domain.OpenDouyinUserInfo) ([]*domain.JinritemaiOrderInfoStatisticsAwemeIndustry, error) {
	list := make([]*domain.JinritemaiOrderInfoStatisticsAwemeIndustry, 0)

	where := make([]string, 0)

	if companyId > 0 {
		where = append(where, "company_id='"+strconv.FormatUint(companyId, 10)+"'")
	}

	if len(startDay) > 0 {
		where = append(where, "pay_success_time>='"+startDay+" 00:00:00'")
	}

	if len(endDay) > 0 {
		where = append(where, "pay_success_time<='"+endDay+" 23:59:59'")
	}

	openDouyinUserInfoSqls := make([]string, 0)

	for _, openDouyinUserInfo := range openDouyinUserInfos {
		openDouyinUserInfoSqls = append(openDouyinUserInfoSqls, "(client_key = '"+openDouyinUserInfo.ClientKey+"' and open_id = '"+openDouyinUserInfo.OpenId+"')")
	}

	where = append(where, strings.Join(openDouyinUserInfoSqls, " or "))

	group := " group by client_key,open_id,industry_id"

	sql := "SELECT client_key,open_id,industry_id,max(industry_name) industry_name,sum(item_num + refund_item_num) as item_num FROM jinritemai_report_order WHERE " + strings.Join(where, " AND ") + group
	rows, err := joir.data.cdb.Query(ctx, sql)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			clientKey    string
			openId       string
			industryId   uint64
			industryName string
			itemNum      uint64
		)

		if err := rows.Scan(&clientKey, &openId, &industryId, &industryName, &itemNum); err != nil {
			return nil, err
		}

		list = append(list, &domain.JinritemaiOrderInfoStatisticsAwemeIndustry{
			ClientKey:    clientKey,
			OpenId:       openId,
			IndustryId:   industryId,
			IndustryName: industryName,
			ItemNum:      itemNum,
		})
	}

	rows.Close()

	return list, nil
}

func (joir *jinritemaiOrderInfoRepo) Upsert(ctx context.Context, in *domain.JinritemaiOrderInfo) error {
	if result := joir.data.DB(ctx).Table("douyin_jinritemai_order_info").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "order_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"client_key", "open_id", "buyin_id", "order_id", "product_id", "product_name", "product_img", "commission_rate", "pay_success_time", "settle_time", "total_pay_amount", "pay_goods_amount", "flow_point", "estimated_commission", "real_commission", "item_num", "pick_extra", "media_type", "media_id", "update_time"}),
	}).Create(&domain.JinritemaiOrderInfoGorm{
		ClientKey:           in.ClientKey,
		OpenId:              in.OpenId,
		BuyinId:             in.BuyinId,
		OrderId:             in.OrderId,
		ProductId:           in.ProductId,
		ProductName:         in.ProductName,
		ProductImg:          in.ProductImg,
		CommissionRate:      in.CommissionRate,
		PaySuccessTime:      in.PaySuccessTime,
		SettleTime:          in.SettleTime,
		TotalPayAmount:      in.TotalPayAmount,
		PayGoodsAmount:      in.PayGoodsAmount,
		FlowPoint:           in.FlowPoint,
		EstimatedCommission: in.EstimatedCommission,
		RealCommission:      in.RealCommission,
		ItemNum:             in.ItemNum,
		PickExtra:           in.PickExtra,
		MediaType:           in.MediaType,
		MediaId:             in.MediaId,
		CreateTime:          in.CreateTime,
		UpdateTime:          in.UpdateTime,
	}); result.Error != nil {
		return result.Error
	}

	return nil
}

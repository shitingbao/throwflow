package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"unicode/utf8"
	"weixin/internal/biz"
	"weixin/internal/domain"
)

// 微信小程序用户申请寄样表
type UserSampleOrder struct {
	Id               uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	UserId           uint64    `gorm:"column:user_id;type:bigint(20) UNSIGNED;not null;comment:微信小程序用户ID"`
	OpenDouyinUserId uint64    `gorm:"column:open_douyin_user_id;type:bigint(20) UNSIGNED;not null;comment:微信小程序用户关联抖音开放平台用户ID"`
	ProductOutId     uint64    `gorm:"column:product_out_id;type:bigint(20) UNSIGNED;index:product_out_id;not null;comment:外部商品ID"`
	ProductName      string    `gorm:"column:product_name;type:varchar(250);not null;comment:商品名称"`
	ProductImg       string    `gorm:"column:product_img;type:text;not null;comment:商品图片地址"`
	OrderSn          string    `gorm:"column:order_sn;type:varchar(20);not null;comment:订单编号"`
	Name             string    `gorm:"column:name;type:varchar(100);not null;comment:收货人"`
	Phone            string    `gorm:"column:phone;type:varchar(20);not null;comment:手机号码"`
	ProvinceAreaName string    `gorm:"column:province_area_name;type:varchar(50);not null;comment:省区划名称"`
	CityAreaName     string    `gorm:"column:city_area_name;type:varchar(50);not null;comment:市区划名称"`
	AreaAreaName     string    `gorm:"column:area_area_name;type:varchar(50);not null;comment:区区划名称"`
	AddressInfo      string    `gorm:"column:address_info;type:text;not null;comment:详细地址"`
	Note             string    `gorm:"column:note;type:text;not null;comment:备注"`
	IsCancel         uint8     `gorm:"column:is_cancel;type:tinyint(3) UNSIGNED;not null;default:0;comment:取消：1：已取消，0：未取消"`
	CancelNote       string    `gorm:"column:cancel_note;type:text;not null;comment:取消备注"`
	KuaidiCompany    string    `gorm:"column:kuaidi_company;type:char(100);not null;comment:快递公司"`
	KuaidiCode       string    `gorm:"column:kuaidi_code;type:char(100);not null;comment:快递公司简码"`
	KuaidiNum        string    `gorm:"column:kuaidi_num;type:char(100);not null;comment:运单号"`
	KuaidiStateName  string    `gorm:"column:kuaidi_state_name;type:char(100);not null;comment:运单号状态"`
	CreateTime       time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime       time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (UserSampleOrder) TableName() string {
	return "weixin_user_sample_order"
}

type userSampleOrderRepo struct {
	data *Data
	log  *log.Helper
}

func (uso *UserSampleOrder) ToDomain(ctx context.Context) *domain.UserSampleOrder {
	userSampleOrder := &domain.UserSampleOrder{
		Id:               uso.Id,
		UserId:           uso.UserId,
		OpenDouyinUserId: uso.OpenDouyinUserId,
		ProductOutId:     uso.ProductOutId,
		ProductName:      uso.ProductName,
		ProductImg:       uso.ProductImg,
		OrderSn:          uso.OrderSn,
		Name:             uso.Name,
		Phone:            uso.Phone,
		ProvinceAreaName: uso.ProvinceAreaName,
		CityAreaName:     uso.CityAreaName,
		AreaAreaName:     uso.AreaAreaName,
		AddressInfo:      uso.AddressInfo,
		Note:             uso.Note,
		IsCancel:         uso.IsCancel,
		CancelNote:       uso.CancelNote,
		KuaidiCompany:    uso.KuaidiCompany,
		KuaidiCode:       uso.KuaidiCode,
		KuaidiNum:        uso.KuaidiNum,
		KuaidiStateName:  uso.KuaidiStateName,
		CreateTime:       uso.CreateTime,
		UpdateTime:       uso.UpdateTime,
	}

	return userSampleOrder
}

func NewUserSampleOrderRepo(data *Data, logger log.Logger) biz.UserSampleOrderRepo {
	return &userSampleOrderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (usor *userSampleOrderRepo) NextId(ctx context.Context) (uint64, error) {
	return usor.data.sonyflake.NextID()
}

func (usor *userSampleOrderRepo) Get(ctx context.Context, userSampleOrderId uint64) (*domain.UserSampleOrder, error) {
	userSampleOrder := &UserSampleOrder{}

	if result := usor.data.db.WithContext(ctx).
		Where("id = ?", userSampleOrderId).
		First(userSampleOrder); result.Error != nil {
		return nil, result.Error
	}

	return userSampleOrder.ToDomain(ctx), nil
}

func (usor *userSampleOrderRepo) List(ctx context.Context, pageNum, pageSize int, userId uint64) ([]*domain.UserSampleOrder, error) {
	var userSampleOrders []UserSampleOrder
	list := make([]*domain.UserSampleOrder, 0)

	if result := usor.data.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Order("create_time DESC,id DESC").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&userSampleOrders); result.Error != nil {
		return nil, result.Error
	}

	for _, userSampleOrder := range userSampleOrders {
		list = append(list, userSampleOrder.ToDomain(ctx))
	}

	return list, nil
}

func (usor *userSampleOrderRepo) ListProduct(ctx context.Context, keyword string) ([]*domain.UserSampleOrder, error) {
	var userSampleOrders []UserSampleOrder
	list := make([]*domain.UserSampleOrder, 0)

	db := usor.data.db.WithContext(ctx)

	if l := utf8.RuneCountInString(keyword); l > 0 {
		db = db.Where("product_name like ?", "%"+keyword+"%")
	}

	if result := db.Select("product_out_id,max(product_name) as product_name").
		Group("product_out_id").
		Order("product_out_id DESC").
		Find(&userSampleOrders); result.Error != nil {
		return nil, result.Error
	}

	for _, userSampleOrder := range userSampleOrders {
		list = append(list, userSampleOrder.ToDomain(ctx))
	}

	return list, nil
}

func (usor *userSampleOrderRepo) ListOpenDouyinUser(ctx context.Context, pageNum, pageSize int, day, keyword, searchType string) ([]*domain.UserOpenDouyin, error) {
	var userOpenDouyins []UserOpenDouyin
	list := make([]*domain.UserOpenDouyin, 0)

	db := usor.data.db.WithContext(ctx).Table("weixin_user_sample_order").
		Joins("left join weixin_user_open_douyin on weixin_user_sample_order.open_douyin_user_id = weixin_user_open_douyin.id").
		Select("weixin_user_open_douyin.id,weixin_user_open_douyin.nickname,weixin_user_open_douyin.account_id,weixin_user_open_douyin.avatar,weixin_user_open_douyin.avatar_larger,weixin_user_open_douyin.fans")

	if len(day) > 0 {
		db = db.Where("weixin_user_sample_order.update_time >= ?", day+" 00:00:00")
		db = db.Where("weixin_user_sample_order.update_time <= ?", day+" 23:59:59")
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		if searchType == "product" {
			db = db.Where("weixin_user_sample_order.product_name like ?", "%"+keyword+"%")
		} else if searchType == "aweme" {
			db = db.Where("weixin_user_open_douyin.nickname like ?", "%"+keyword+"%")
		} else if searchType == "address" {
			db = db.Where("weixin_user_sample_order.address_info like ?", "%"+keyword+"%")
		}
	}

	if result := db.Group("weixin_user_sample_order.open_douyin_user_id").
		Order("weixin_user_sample_order.open_douyin_user_id ASC").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&userOpenDouyins); result.Error != nil {
		return nil, result.Error
	}

	for _, userOpenDouyin := range userOpenDouyins {
		list = append(list, userOpenDouyin.ToDomain(ctx))
	}

	return list, nil
}

func (usor *userSampleOrderRepo) ListByOpenDouyinUserIds(ctx context.Context, openDouyinUserIds []uint64, day, keyword, searchType string) ([]*domain.UserSampleOrder, error) {
	var userSampleOrders []UserSampleOrder
	list := make([]*domain.UserSampleOrder, 0)

	db := usor.data.db.WithContext(ctx).
		Joins("left join weixin_user_open_douyin on weixin_user_sample_order.open_douyin_user_id = weixin_user_open_douyin.id").
		Select("weixin_user_sample_order.*,weixin_user_open_douyin.nickname,weixin_user_open_douyin.avatar,weixin_user_open_douyin.avatar_larger,weixin_user_open_douyin.fans").
		Where("weixin_user_sample_order.open_douyin_user_id in ?", openDouyinUserIds)

	if len(day) > 0 {
		db = db.Where("weixin_user_sample_order.create_time >= ?", day+" 00:00:00")
		db = db.Where("weixin_user_sample_order.create_time <= ?", day+" 23:59:59")
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		if searchType == "product" {
			db = db.Where("weixin_user_sample_order.product_name like ?", "%"+keyword+"%")
		} else if searchType == "aweme" {
			db = db.Where("weixin_user_open_douyin.nickname like ?", "%"+keyword+"%")
		} else if searchType == "address" {
			db = db.Where("weixin_user_sample_order.address_info like ?", "%"+keyword+"%")
		}
	}

	if result := db.Order("weixin_user_sample_order.id DESC").
		Find(&userSampleOrders); result.Error != nil {
		return nil, result.Error
	}

	for _, userSampleOrder := range userSampleOrders {
		list = append(list, userSampleOrder.ToDomain(ctx))
	}

	return list, nil
}

func (usor *userSampleOrderRepo) Count(ctx context.Context, userId uint64) (int64, error) {
	var count int64

	if result := usor.data.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Model(&UserSampleOrder{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (usor *userSampleOrderRepo) CountOpenDouyinUser(ctx context.Context, day, keyword, searchType string) (int64, error) {
	var count int64

	db := usor.data.db.WithContext(ctx).
		Joins("left join weixin_user_open_douyin on weixin_user_sample_order.open_douyin_user_id = weixin_user_open_douyin.id").
		Select("weixin_user_sample_order.open_douyin_user_id")

	if len(day) > 0 {
		db = db.Where("weixin_user_sample_order.update_time >= ?", day+" 00:00:00")
		db = db.Where("weixin_user_sample_order.update_time <= ?", day+" 23:59:59")
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		if searchType == "product" {
			db = db.Where("weixin_user_sample_order.product_name like ?", "%"+keyword+"%")
		} else if searchType == "aweme" {
			db = db.Where("weixin_user_open_douyin.nickname like ?", "%"+keyword+"%")
		} else if searchType == "address" {
			db = db.Where("weixin_user_sample_order.address_info like ?", "%"+keyword+"%")
		}
	}

	db = db.Group("weixin_user_sample_order.open_douyin_user_id")

	if result := db.Model(&UserSampleOrder{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (usor *userSampleOrderRepo) Statistics(ctx context.Context, day, keyword, searchType string) (int64, error) {
	var count int64

	db := usor.data.db.WithContext(ctx).
		Joins("left join weixin_user_open_douyin on weixin_user_sample_order.open_douyin_user_id = weixin_user_open_douyin.id").
		Select("weixin_user_sample_order.id")

	if len(day) > 0 {
		db = db.Where("weixin_user_sample_order.update_time >= ?", day+" 00:00:00")
		db = db.Where("weixin_user_sample_order.update_time <= ?", day+" 23:59:59")
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		if searchType == "product" {
			db = db.Where("weixin_user_sample_order.product_name like ?", "%"+keyword+"%")
		} else if searchType == "aweme" {
			db = db.Where("weixin_user_open_douyin.nickname like ?", "%"+keyword+"%")
		} else if searchType == "address" {
			db = db.Where("weixin_user_sample_order.address_info like ?", "%"+keyword+"%")
		}
	}

	if result := db.Model(&UserSampleOrder{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (usor *userSampleOrderRepo) Save(ctx context.Context, in *domain.UserSampleOrder) (*domain.UserSampleOrder, error) {
	userSampleOrder := &UserSampleOrder{
		UserId:           in.UserId,
		OpenDouyinUserId: in.OpenDouyinUserId,
		ProductOutId:     in.ProductOutId,
		ProductName:      in.ProductName,
		ProductImg:       in.ProductImg,
		OrderSn:          in.OrderSn,
		Name:             in.Name,
		Phone:            in.Phone,
		ProvinceAreaName: in.ProvinceAreaName,
		CityAreaName:     in.CityAreaName,
		AreaAreaName:     in.AreaAreaName,
		AddressInfo:      in.AddressInfo,
		Note:             in.Note,
		IsCancel:         in.IsCancel,
		CancelNote:       in.CancelNote,
		KuaidiCompany:    in.KuaidiCompany,
		KuaidiCode:       in.KuaidiCode,
		KuaidiNum:        in.KuaidiNum,
		KuaidiStateName:  in.KuaidiStateName,
		CreateTime:       in.CreateTime,
		UpdateTime:       in.UpdateTime,
	}

	if result := usor.data.DB(ctx).Create(userSampleOrder); result.Error != nil {
		return nil, result.Error
	}

	return userSampleOrder.ToDomain(ctx), nil
}

func (usor *userSampleOrderRepo) Update(ctx context.Context, in *domain.UserSampleOrder) (*domain.UserSampleOrder, error) {
	userSampleOrder := &UserSampleOrder{
		Id:               in.Id,
		UserId:           in.UserId,
		OpenDouyinUserId: in.OpenDouyinUserId,
		ProductOutId:     in.ProductOutId,
		ProductName:      in.ProductName,
		ProductImg:       in.ProductImg,
		OrderSn:          in.OrderSn,
		Name:             in.Name,
		Phone:            in.Phone,
		ProvinceAreaName: in.ProvinceAreaName,
		CityAreaName:     in.CityAreaName,
		AreaAreaName:     in.AreaAreaName,
		AddressInfo:      in.AddressInfo,
		Note:             in.Note,
		IsCancel:         in.IsCancel,
		CancelNote:       in.CancelNote,
		KuaidiCompany:    in.KuaidiCompany,
		KuaidiCode:       in.KuaidiCode,
		KuaidiNum:        in.KuaidiNum,
		KuaidiStateName:  in.KuaidiStateName,
		CreateTime:       in.CreateTime,
		UpdateTime:       in.UpdateTime,
	}

	if result := usor.data.DB(ctx).Save(userSampleOrder); result.Error != nil {
		return nil, result.Error
	}

	return userSampleOrder.ToDomain(ctx), nil
}

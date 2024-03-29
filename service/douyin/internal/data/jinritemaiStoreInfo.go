package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
	"time"
)

// 精选联盟达人橱窗商品信息表
type JinritemaiStoreInfo struct {
	Id                uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	ClientKey         string    `gorm:"column:client_key;type:varchar(50);not null;comment:应用Client Key"`
	OpenId            string    `gorm:"column:open_id;type:varchar(100);not null;comment:授权用户唯一标识"`
	ProductId         string    `gorm:"column:product_id;type:varchar(100);not null;comment:商品ID"`
	ProductName       string    `gorm:"column:product_name;type:varchar(100);not null;comment:商品名称"`
	ProductImg        string    `gorm:"column:product_img;type:varchar(250);not null;comment:商品图片URL"`
	ProductPrice      float32   `gorm:"column:product_price;type:decimal(10, 2) UNSIGNED;not null;comment:商品售价（单位为元）"`
	CommissionType    uint8     `gorm:"column:commission_type;type:tinyint(3) UNSIGNED;not null;default:0;comment:佣金类型，0 未定义（异常） 1 专属团长 2 普通佣金 3 定向佣金，4 提报活动 5 招募佣金"`
	CommissionRatio   uint8     `gorm:"column:commission_ratio;type:tinyint(3) UNSIGNED;not null;comment:佣金率（10表示佣金率为10百分比）"`
	PromotionId       uint64    `gorm:"column:promotion_id;type:bigint(20) UNSIGNED;not null;comment:推广ID"`
	PromotionType     uint8     `gorm:"column:promotion_type;type:tinyint(3) UNSIGNED;not null;comment:推广类型，0 非团长，1 团长"`
	ColonelActivityId uint64    `gorm:"column:colonel_activity_id;type:bigint(20) UNSIGNED;not null;comment:团长活动id"`
	CreateTime        time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime        time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (JinritemaiStoreInfo) TableName() string {
	return "douyin_jinritemai_store_info"
}

type jinritemaiStoreInfoRepo struct {
	data *Data
	log  *log.Helper
}

func (jsi *JinritemaiStoreInfo) ToDomain() *domain.JinritemaiStoreInfo {
	return &domain.JinritemaiStoreInfo{
		Id:                jsi.Id,
		ClientKey:         jsi.ClientKey,
		OpenId:            jsi.OpenId,
		ProductId:         jsi.ProductId,
		ProductName:       jsi.ProductName,
		ProductImg:        jsi.ProductImg,
		ProductPrice:      jsi.ProductPrice,
		CommissionType:    jsi.CommissionType,
		CommissionRatio:   jsi.CommissionRatio,
		PromotionId:       jsi.PromotionId,
		PromotionType:     jsi.PromotionType,
		ColonelActivityId: jsi.ColonelActivityId,
		CreateTime:        jsi.CreateTime,
		UpdateTime:        jsi.UpdateTime,
	}
}

func NewJinritemaiStoreInfoRepo(data *Data, logger log.Logger) biz.JinritemaiStoreInfoRepo {
	return &jinritemaiStoreInfoRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (jsir *jinritemaiStoreInfoRepo) List(ctx context.Context, pageNum, pageSize int, openDouyinTokens []*domain.OpenDouyinToken) ([]*domain.JinritemaiStoreInfo, error) {
	var jinritemaiStoreInfos []JinritemaiStoreInfo
	list := make([]*domain.JinritemaiStoreInfo, 0)

	db := jsir.data.db.WithContext(ctx)

	openDouyinTokenSqls := make([]string, 0)

	for _, openDouyinToken := range openDouyinTokens {
		openDouyinTokenSqls = append(openDouyinTokenSqls, "(client_key = '"+openDouyinToken.ClientKey+"' and open_id = '"+openDouyinToken.OpenId+"')")
	}

	db = db.Where(strings.Join(openDouyinTokenSqls, " or "))

	if pageNum == 0 {
		if result := db.Order("product_id desc, id desc").
			Find(&jinritemaiStoreInfos); result.Error != nil {
			return nil, result.Error
		}
	} else {
		if result := db.Order("product_id desc, id desc").
			Limit(pageSize).Offset((pageNum - 1) * pageSize).
			Find(&jinritemaiStoreInfos); result.Error != nil {
			return nil, result.Error
		}
	}

	for _, jinritemaiStoreInfo := range jinritemaiStoreInfos {
		list = append(list, jinritemaiStoreInfo.ToDomain())
	}

	return list, nil
}

func (jsir *jinritemaiStoreInfoRepo) ListByIds(ctx context.Context, storeIds []uint64) ([]*domain.JinritemaiStoreInfo, error) {
	var jinritemaiStoreInfos []JinritemaiStoreInfo
	list := make([]*domain.JinritemaiStoreInfo, 0)

	if result := jsir.data.db.WithContext(ctx).
		Where("id IN (?)", storeIds).
		Order("open_id desc, id desc").
		Find(&jinritemaiStoreInfos); result.Error != nil {
		return nil, result.Error
	}

	for _, jinritemaiStoreInfo := range jinritemaiStoreInfos {
		list = append(list, jinritemaiStoreInfo.ToDomain())
	}

	return list, nil
}

func (jsir *jinritemaiStoreInfoRepo) ListByProductId(ctx context.Context, productId string) ([]*domain.OpenDouyinUserInfo, error) {
	var openDouyinUserInfos []OpenDouyinUserInfo
	list := make([]*domain.OpenDouyinUserInfo, 0)

	db := jsir.data.db.WithContext(ctx).
		Model(JinritemaiStoreInfo{}).
		Select("douyin_open_douyin_user_info.*").
		Joins("left join douyin_open_douyin_user_info on douyin_open_douyin_user_info.client_key = douyin_jinritemai_store_info.client_key and douyin_open_douyin_user_info.open_id = douyin_jinritemai_store_info.open_id").
		Where("douyin_jinritemai_store_info.product_id = ?", productId)

	if result := db.Order("douyin_jinritemai_store_info.id DESC").
		Find(&openDouyinUserInfos); result.Error != nil {
		return nil, result.Error
	}

	for _, openDouyinUserInfo := range openDouyinUserInfos {
		list = append(list, openDouyinUserInfo.ToDomain())
	}

	return list, nil
}

func (jsir *jinritemaiStoreInfoRepo) ListByClientKeyAndOpenIdAndProductIds(ctx context.Context, clientKey, openId string, productIds []string) ([]*domain.JinritemaiStoreInfo, error) {
	var jinritemaiStoreInfos []JinritemaiStoreInfo
	list := make([]*domain.JinritemaiStoreInfo, 0)

	if result := jsir.data.db.WithContext(ctx).
		Where("client_key = ?", clientKey).
		Where("open_id = ?", openId).
		Where("product_id IN (?)", productIds).
		Order("id desc").
		Find(&jinritemaiStoreInfos); result.Error != nil {
		return nil, result.Error
	}

	for _, jinritemaiStoreInfo := range jinritemaiStoreInfos {
		list = append(list, jinritemaiStoreInfo.ToDomain())
	}

	return list, nil
}

func (jsir *jinritemaiStoreInfoRepo) ListProductId(ctx context.Context, pageNum, pageSize int) ([]*domain.JinritemaiStoreInfo, error) {
	var jinritemaiStoreInfos []JinritemaiStoreInfo
	list := make([]*domain.JinritemaiStoreInfo, 0)

	db := jsir.data.db.WithContext(ctx)

	if result := db.Select("product_id").Group("product_id").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&jinritemaiStoreInfos); result.Error != nil {
		return nil, result.Error
	}

	for _, jinritemaiStoreInfo := range jinritemaiStoreInfos {
		list = append(list, jinritemaiStoreInfo.ToDomain())
	}

	return list, nil
}

func (jsir *jinritemaiStoreInfoRepo) Count(ctx context.Context, openDouyinTokens []*domain.OpenDouyinToken) (int64, error) {
	var count int64

	db := jsir.data.db.WithContext(ctx)

	openDouyinTokenSqls := make([]string, 0)

	for _, openDouyinToken := range openDouyinTokens {
		openDouyinTokenSqls = append(openDouyinTokenSqls, "(client_key = '"+openDouyinToken.ClientKey+"' and open_id = '"+openDouyinToken.OpenId+"')")
	}

	db = db.Where(strings.Join(openDouyinTokenSqls, " or "))

	if result := db.Model(&JinritemaiStoreInfo{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (jsir *jinritemaiStoreInfoRepo) CountProductId(ctx context.Context) (int64, error) {
	var count int64

	if result := jsir.data.db.WithContext(ctx).Group("product_id").Model(&JinritemaiStoreInfo{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (jsir *jinritemaiStoreInfoRepo) DeleteByDayAndClientKeyAndOpenIdAndProductIds(ctx context.Context, clientKey, openId string, productIds []string) error {
	if result := jsir.data.DB(ctx).
		Where("client_key = ?", clientKey).
		Where("open_id = ?", openId).
		Where("product_id IN (?)", productIds).
		Delete(&JinritemaiStoreInfo{}); result.Error != nil {
		return result.Error
	}

	return nil
}

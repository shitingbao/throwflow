package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 抖客订单表
type DoukeOrder struct {
	AuthorAccount          string             `json:"author_account" bson:"author_account"`
	MediaType              string             `json:"media_type" bson:"media_type"`
	AdsActivityId          int64              `json:"ads_activity_id" bson:"ads_activity_id"`
	ProductImg             string             `json:"product_img" bson:"product_img"`
	UpdateTime             string             `json:"update_time" bson:"update_time"`
	PaySuccessTime         string             `json:"pay_success_time" bson:"pay_success_time"`
	AdsRealCommission      int64              `json:"ads_real_commission" bson:"ads_real_commission"`
	AdsEstimatedCommission int64              `json:"ads_estimated_commission" bson:"ads_estimated_commission"`
	ProductId              string             `json:"product_id" bson:"product_id"`
	TotalPayAmount         int64              `json:"total_pay_amount" bson:"total_pay_amount"`
	FlowPoint              string             `json:"flow_point" bson:"flow_point"`
	SettleTime             string             `json:"settle_time" bson:"settle_time"`
	SettledGoodsAmount     int64              `json:"settled_goods_amount" bson:"settled_goods_amount"`
	PidInfo                domain.PidInfo     `json:"pid_info" bson:"pid_info"`
	ItemNum                int64              `json:"item_num" bson:"item_num"`
	AuthorBuyinId          string             `json:"author_buyin_id" bson:"author_buyin_id"`
	ShopId                 int64              `json:"shop_id" bson:"shop_id"`
	PayGoodsAmount         int64              `json:"pay_goods_amount" bson:"pay_goods_amount"`
	AdsDistributorId       int64              `json:"ads_distributor_id" bson:"ads_distributor_id"`
	AdsPromotionTate       int64              `json:"ads_promotion_rate" bson:"ads_promotion_rate"`
	OrderId                string             `json:"order_id" bson:"order_id"`
	ProductName            string             `json:"product_name" bson:"product_name"`
	DistributionType       string             `json:"distribution_type" bson:"distribution_type"`
	AuthorUid              int64              `json:"author_uid" bson:"author_uid"`
	ShopName               string             `json:"shop_name" bson:"shop_name"`
	ProductActivityId      string             `json:"product_activity_id" bson:"product_activity_id"`
	MaterialId             string             `json:"material_id" bson:"material_id"`
	RefundTime             string             `json:"refund_time" bson:"refund_time"`
	ConfirmTime            string             `json:"confirm_time" bson:"confirm_time"`
	ProductTags            domain.ProductTags `json:"product_tags" bson:"product_tags"`
	BuyerAppId             string             `json:"buyer_app_id" bson:"buyer_app_id"`
	DistributorRightType   string             `json:"distributor_right_type" bson:"distributor_right_type"`
}

type doukeOrderRepo struct {
	data *Data
	log  *log.Helper
}

func (do *DoukeOrder) ToDomain() *domain.DoukeOrder {
	return &domain.DoukeOrder{
		AuthorAccount:          do.AuthorAccount,
		MediaType:              do.MediaType,
		AdsActivityId:          do.AdsActivityId,
		ProductImg:             do.ProductImg,
		UpdateTime:             do.UpdateTime,
		PaySuccessTime:         do.PaySuccessTime,
		AdsRealCommission:      do.AdsRealCommission,
		AdsEstimatedCommission: do.AdsEstimatedCommission,
		ProductId:              do.ProductId,
		TotalPayAmount:         do.TotalPayAmount,
		FlowPoint:              do.FlowPoint,
		SettleTime:             do.SettleTime,
		SettledGoodsAmount:     do.SettledGoodsAmount,
		PidInfo:                do.PidInfo,
		ItemNum:                do.ItemNum,
		AuthorBuyinId:          do.AuthorBuyinId,
		ShopId:                 do.ShopId,
		PayGoodsAmount:         do.PayGoodsAmount,
		AdsDistributorId:       do.AdsDistributorId,
		AdsPromotionTate:       do.AdsPromotionTate,
		OrderId:                do.OrderId,
		ProductName:            do.ProductName,
		DistributionType:       do.DistributionType,
		AuthorUid:              do.AuthorUid,
		ShopName:               do.ShopName,
		ProductActivityId:      do.ProductActivityId,
		MaterialId:             do.MaterialId,
		RefundTime:             do.RefundTime,
		ConfirmTime:            do.ConfirmTime,
		ProductTags:            do.ProductTags,
		BuyerAppId:             do.BuyerAppId,
		DistributorRightType:   do.DistributorRightType,
	}
}

func NewDoukeOrderRepo(data *Data, logger log.Logger) biz.DoukeOrderRepo {
	return &doukeOrderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (dor *doukeOrderRepo) SaveIndex(ctx context.Context) {
	collection := dor.data.mdb.Database(dor.data.conf.Mongo.Dbname).Collection("douke_order")

	isNotExistIndex := true

	if indexSpecifications, err := collection.Indexes().ListSpecifications(ctx); err == nil {
		for _, indexSpecification := range indexSpecifications {
			if indexSpecification.Name == "order_id_-1" {
				isNotExistIndex = false
			}
		}

		if isNotExistIndex {
			collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{Key: "order_id", Value: -1},
				},
			})
		}
	}
}

func (dor *doukeOrderRepo) Upsert(ctx context.Context, in *domain.DoukeOrder) error {
	collection := dor.data.mdb.Database(dor.data.conf.Mongo.Dbname).Collection("douke_order")

	if _, err := collection.UpdateOne(ctx, bson.D{
		{"order_id", in.OrderId},
	}, bson.D{
		{"$set", bson.D{
			{"author_account", in.AuthorAccount},
			{"media_type", in.MediaType},
			{"ads_activity_id", in.AdsActivityId},
			{"product_img", in.ProductImg},
			{"update_time", in.UpdateTime},
			{"pay_success_time", in.PaySuccessTime},
			{"ads_real_commission", in.AdsRealCommission},
			{"ads_estimated_commission", in.AdsEstimatedCommission},
			{"product_id", in.ProductId},
			{"total_pay_amount", in.TotalPayAmount},
			{"flow_point", in.FlowPoint},
			{"settle_time", in.SettleTime},
			{"settled_goods_amount", in.SettledGoodsAmount},
			{"pid_info", in.PidInfo},
			{"item_num", in.ItemNum},
			{"shop_id", in.ShopId},
			{"pay_goods_amount", in.PayGoodsAmount},
			{"ads_distributor_id", in.AdsDistributorId},
			{"ads_promotion_rate", in.AdsPromotionTate},
			{"product_name", in.ProductName},
			{"distribution_type", in.DistributionType},
			{"author_uid", in.AuthorUid},
			{"shop_name", in.ShopName},
			{"product_activity_id", in.ProductActivityId},
			{"material_id", in.MaterialId},
			{"refund_time", in.RefundTime},
			{"confirm_time", in.ConfirmTime},
			{"product_tags", in.ProductTags},
			{"buyer_app_id", in.BuyerAppId},
			{"distributor_right_type", in.DistributorRightType},
		}},
	}, options.Update().SetUpsert(true)); err != nil {
		return err
	}

	return nil
}

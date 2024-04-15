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
	OrderId                 string  `json:"order_id" bson:"order_id"`
	AppId                   string  `json:"app_id" bson:"app_id"`
	ProductId               string  `json:"product_id" bson:"product_id"`
	ProductName             string  `json:"product_name" bson:"product_name"`
	AuthorAccount           string  `json:"author_account" bson:"author_account"`
	AdsAttribution          string  `json:"ads_attribution" bson:"ads_attribution"`
	ProductImg              string  `json:"product_img" bson:"product_img"`
	TotalPayAmount          int     `json:"total_pay_amount" bson:"total_pay_amount"`
	PaySuccessTime          string  `json:"pay_success_time" bson:"pay_success_time"`
	RefundTime              string  `json:"refund_time" bson:"refund_time"`
	PayGoodsAmount          int     `json:"pay_goods_amount" bson:"pay_goods_amount"`
	EstimatedCommission     float32 `json:"estimated_commission" bson:"estimated_commission"`
	AdsRealCommission       float32 `json:"ads_real_commission" bson:"ads_real_commission"`
	SplitRate               float32 `json:"split_rate" bson:"split_rate"`
	AfterSalesStatus        int     `json:"after_sales_status" bson:"after_sales_status"`
	FlowPoint               string  `json:"flow_point" bson:"flow_point"`
	ExternalInfo            string  `json:"external_info" bson:"external_info"`
	SettleTime              string  `json:"settle_time" bson:"settle_time"`
	ConfirmTime             string  `json:"confirm_time" bson:"confirm_time"`
	MediaTypeName           string  `json:"media_type_name" bson:"media_type_name"`
	UpdateTime              string  `json:"update_time" bson:"update_time"`
	EstimatedTechServiceFee int     `json:"estimated_tech_service_fee" bson:"estimated_tech_service_fee"`
}

type doukeOrderRepo struct {
	data *Data
	log  *log.Helper
}

func (do *DoukeOrder) ToDomain() *domain.DoukeOrder {
	return &domain.DoukeOrder{
		OrderId:                 do.OrderId,
		AppId:                   do.AppId,
		ProductId:               do.ProductId,
		ProductName:             do.ProductName,
		AuthorAccount:           do.AuthorAccount,
		AdsAttribution:          do.AdsAttribution,
		ProductImg:              do.ProductImg,
		TotalPayAmount:          do.TotalPayAmount,
		PaySuccessTime:          do.PaySuccessTime,
		RefundTime:              do.RefundTime,
		PayGoodsAmount:          do.PayGoodsAmount,
		EstimatedCommission:     do.EstimatedCommission,
		AdsRealCommission:       do.AdsRealCommission,
		SplitRate:               do.SplitRate,
		AfterSalesStatus:        do.AfterSalesStatus,
		FlowPoint:               do.FlowPoint,
		ExternalInfo:            do.ExternalInfo,
		SettleTime:              do.SettleTime,
		ConfirmTime:             do.ConfirmTime,
		MediaTypeName:           do.MediaTypeName,
		UpdateTime:              do.UpdateTime,
		EstimatedTechServiceFee: do.EstimatedTechServiceFee,
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
			{"app_id", in.AppId},
			{"product_id", in.ProductId},
			{"product_name", in.ProductName},
			{"author_account", in.AuthorAccount},
			{"ads_attribution", in.AdsAttribution},
			{"product_img", in.ProductImg},
			{"total_pay_amount", in.TotalPayAmount},
			{"pay_success_time", in.PaySuccessTime},
			{"refund_time", in.RefundTime},
			{"pay_goods_amount", in.PayGoodsAmount},
			{"estimated_commission", in.EstimatedCommission},
			{"ads_real_commission", in.AdsRealCommission},
			{"split_rate", in.SplitRate},
			{"after_sales_status", in.AfterSalesStatus},
			{"flow_point", in.FlowPoint},
			{"external_info", in.ExternalInfo},
			{"settle_time", in.SettleTime},
			{"confirm_time", in.ConfirmTime},
			{"media_type_name", in.MediaTypeName},
			{"update_time", in.UpdateTime},
			{"estimated_tech_service_fee", in.EstimatedTechServiceFee},
		}},
	}, options.Update().SetUpsert(true)); err != nil {
		return err
	}

	return nil
}

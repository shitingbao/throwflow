package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// 连山云RDS-千川计划数据表
type QianchuanReportAd struct {
	AdId              uint64    `json:"ad_id" bson:"ad_id"`
	AdvertiserId      uint64    `json:"advertiser_id" bson:"advertiser_id"`
	AwemeId           uint64    `json:"aweme_id" bson:"aweme_id"`
	MarketingGoal     int64     `json:"marketing_goal" bson:"marketing_goal"`
	StatCost          float64   `json:"stat_cost" bson:"stat_cost"`
	ShowCnt           int64     `json:"show_cnt" bson:"show_cnt"`
	ClickCnt          int64     `json:"click_cnt" bson:"click_cnt"`
	PayOrderCount     int64     `json:"pay_order_count" bson:"pay_order_count"`
	CreateOrderAmount float64   `json:"create_order_amount" bson:"create_order_amount"`
	CreateOrderCount  int64     `json:"create_order_count" bson:"create_order_count"`
	PayOrderAmount    float64   `json:"pay_order_amount" bson:"pay_order_amount"`
	DyFollow          int64     `json:"dy_follow" bson:"dy_follow"`
	ConvertCnt        int64     `json:"convert_cnt" bson:"convert_cnt"`
	CreateTime        time.Time `json:"create_time" bson:"create_time"`
	UpdateTime        time.Time `json:"update_time" bson:"update_time"`
}

type qianchuanReportAdRepo struct {
	data *Data
	log  *log.Helper
}

func (qra *QianchuanReportAd) ToDomain() *domain.QianchuanReportAd {
	return &domain.QianchuanReportAd{
		AdId:              qra.AdId,
		AdvertiserId:      qra.AdvertiserId,
		AwemeId:           qra.AwemeId,
		MarketingGoal:     qra.MarketingGoal,
		StatCost:          qra.StatCost,
		ShowCnt:           qra.ShowCnt,
		ClickCnt:          qra.ClickCnt,
		PayOrderCount:     qra.PayOrderCount,
		CreateOrderAmount: qra.CreateOrderAmount,
		CreateOrderCount:  qra.CreateOrderCount,
		PayOrderAmount:    qra.PayOrderAmount,
		DyFollow:          qra.DyFollow,
		ConvertCnt:        qra.ConvertCnt,
		CreateTime:        qra.CreateTime,
		UpdateTime:        qra.UpdateTime,
	}
}

func NewQianchuanReportAdRepo(data *Data, logger log.Logger) biz.QianchuanReportAdRepo {
	return &qianchuanReportAdRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (qrar *qianchuanReportAdRepo) Get(ctx context.Context, advertiserId, adId uint64, day string) (*domain.QianchuanReportAd, error) {
	var qianchuanReportAd QianchuanReportAd

	collection := qrar.data.mdb.Database(qrar.data.conf.Mongo.Dbname).Collection("qianchuan_report_ad_" + day)

	if err := collection.FindOne(ctx, bson.M{
		"$and": []bson.M{
			bson.M{"advertiser_id": advertiserId},
			bson.M{"ad_id": adId},
		},
	}).Decode(&qianchuanReportAd); err != nil {
		return nil, err
	}

	return qianchuanReportAd.ToDomain(), nil
}

func (qrar *qianchuanReportAdRepo) ListByMarketingGoal(ctx context.Context, marketingGoal string, day string) ([]*domain.QianchuanReportAd, error) {
	list := make([]*domain.QianchuanReportAd, 0)

	collection := qrar.data.mdb.Database(qrar.data.conf.Mongo.Dbname).Collection("qianchuan_report_ad_" + day)

	var matchStage bson.D

	if marketingGoal == "VIDEO_PROM_GOODS" {
		matchStage = bson.D{
			{"$match", bson.M{"marketing_goal": 1}},
		}
	} else if marketingGoal == "LIVE_PROM_GOODS" {
		matchStage = bson.D{
			{"$match", bson.M{"marketing_goal": 2}},
		}
	} else {
		matchStage = bson.D{
			{"$match", bson.M{"marketing_goal": -1}},
		}
	}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var qianchuanReportAds []*QianchuanReportAd

	err = cursor.All(ctx, &qianchuanReportAds)

	if err != nil {
		return nil, err
	}

	for _, qianchuanReportAd := range qianchuanReportAds {
		list = append(list, &domain.QianchuanReportAd{
			AdId:              qianchuanReportAd.AdId,
			AdvertiserId:      qianchuanReportAd.AdvertiserId,
			AwemeId:           qianchuanReportAd.AwemeId,
			MarketingGoal:     qianchuanReportAd.MarketingGoal,
			StatCost:          qianchuanReportAd.StatCost,
			ShowCnt:           qianchuanReportAd.ShowCnt,
			ClickCnt:          qianchuanReportAd.ClickCnt,
			PayOrderCount:     qianchuanReportAd.PayOrderCount,
			CreateOrderAmount: qianchuanReportAd.CreateOrderAmount,
			CreateOrderCount:  qianchuanReportAd.CreateOrderCount,
			PayOrderAmount:    qianchuanReportAd.PayOrderAmount,
			DyFollow:          qianchuanReportAd.DyFollow,
			ConvertCnt:        qianchuanReportAd.ConvertCnt,
			CreateTime:        qianchuanReportAd.CreateTime,
			UpdateTime:        qianchuanReportAd.UpdateTime,
		})
	}

	return list, nil
}

func (qrar *qianchuanReportAdRepo) SaveIndex(ctx context.Context, day string) {
	collection := qrar.data.mdb.Database(qrar.data.conf.Mongo.Dbname).Collection("qianchuan_report_ad_" + day)

	isNotExistIndex := true

	if indexSpecifications, err := collection.Indexes().ListSpecifications(ctx); err == nil {
		for _, indexSpecification := range indexSpecifications {
			if indexSpecification.Name == "advertiser_id_-1_ad_id_-1" {
				isNotExistIndex = false
			}
		}

		if isNotExistIndex {
			collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{Key: "advertiser_id", Value: -1},
					{Key: "ad_id", Value: -1},
				},
			})
		}
	}
}

func (qrar *qianchuanReportAdRepo) Upsert(ctx context.Context, day string, in *domain.QianchuanReportAd) error {
	collection := qrar.data.mdb.Database(qrar.data.conf.Mongo.Dbname).Collection("qianchuan_report_ad_" + day)

	if _, err := collection.UpdateOne(ctx, bson.D{
		{"ad_id", in.AdId},
		{"advertiser_id", in.AdvertiserId},
	}, bson.D{
		{"$set", bson.D{
			{"aweme_id", in.AwemeId},
			{"marketing_goal", in.MarketingGoal},
			{"stat_cost", in.StatCost},
			{"show_cnt", in.ShowCnt},
			{"click_cnt", in.ClickCnt},
			{"pay_order_count", in.PayOrderCount},
			{"create_order_amount", in.CreateOrderAmount},
			{"create_order_count", in.CreateOrderCount},
			{"pay_order_amount", in.PayOrderAmount},
			{"dy_follow", in.DyFollow},
			{"convert_cnt", in.ConvertCnt},
			{"create_time", in.CreateTime},
			{"update_time", in.UpdateTime},
		}},
	}, options.Update().SetUpsert(true)); err != nil {
		return err
	}

	return nil
}

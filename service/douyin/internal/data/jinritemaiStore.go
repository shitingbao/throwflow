package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"douyin/internal/pkg/event/event"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 精选联盟达人橱窗商品表
type JinritemaiStore struct {
	ClientKey         string  `json:"client_key" bson:"client_key"`
	OpenId            string  `json:"open_id" bson:"open_id"`
	ProductId         int64   `json:"product_id" bson:"product_id"`
	PromotionId       int64   `json:"promotion_id" bson:"promotion_id"`
	Title             string  `json:"title" bson:"title"`
	Cover             string  `json:"cover" bson:"cover"`
	PromotionType     int64   `json:"promotion_type" bson:"promotion_type"`
	Price             int64   `json:"price" bson:"price"`
	CosType           int64   `json:"cos_type" bson:"cos_type"`
	CosRatio          float64 `json:"cos_ratio" bson:"cos_ratio"`
	ColonelActivityId int64   `json:"colonel_activity_id" bson:"colonel_activity_id"`
	HideStatus        bool    `json:"hide_status" bson:"hide_status"`
}

type jinritemaiStoreRepo struct {
	data *Data
	log  *log.Helper
}

func (js *JinritemaiStore) ToDomain() *domain.JinritemaiStore {
	return &domain.JinritemaiStore{
		ClientKey:         js.ClientKey,
		OpenId:            js.OpenId,
		ProductId:         js.ProductId,
		PromotionId:       js.PromotionId,
		Title:             js.Title,
		Cover:             js.Cover,
		PromotionType:     js.PromotionType,
		Price:             js.Price,
		CosType:           js.CosType,
		CosRatio:          js.CosRatio,
		ColonelActivityId: js.ColonelActivityId,
		HideStatus:        js.HideStatus,
	}
}

func NewJinritemaiStoreRepo(data *Data, logger log.Logger) biz.JinritemaiStoreRepo {
	return &jinritemaiStoreRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (jsr *jinritemaiStoreRepo) SaveIndex(ctx context.Context, day string) {
	collection := jsr.data.mdb.Database(jsr.data.conf.Mongo.Dbname).Collection("jinritemai_store_" + day)

	isNotExistIndex := true

	if indexSpecifications, err := collection.Indexes().ListSpecifications(ctx); err == nil {
		for _, indexSpecification := range indexSpecifications {
			if indexSpecification.Name == "client_key_-1_open_id_-1_product_id_-1" {
				isNotExistIndex = false
			}
		}

		if isNotExistIndex {
			collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{Key: "client_key", Value: -1},
					{Key: "open_id", Value: -1},
					{Key: "product_id", Value: -1},
				},
			})
		}
	}
}

func (jsr *jinritemaiStoreRepo) Upsert(ctx context.Context, day string, in *domain.JinritemaiStore) error {
	collection := jsr.data.mdb.Database(jsr.data.conf.Mongo.Dbname).Collection("jinritemai_store_" + day)

	if _, err := collection.UpdateOne(ctx, bson.D{
		{"client_key", in.ClientKey},
		{"open_id", in.OpenId},
		{"product_id", in.ProductId},
	}, bson.D{
		{"$set", bson.D{
			{"promotion_id", in.PromotionId},
			{"title", in.Title},
			{"cover", in.Cover},
			{"promotion_type", in.PromotionType},
			{"price", in.Price},
			{"cos_type", in.CosType},
			{"cos_ratio", in.CosRatio},
			{"colonel_activity_id", in.ColonelActivityId},
			{"hide_status", in.HideStatus},
		}},
	}, options.Update().SetUpsert(true)); err != nil {
		return err
	}

	return nil
}

func (jsr *jinritemaiStoreRepo) DeleteByDayAndClientKeyAndOpenId(ctx context.Context, day, clientKey, openId string) error {
	collection := jsr.data.mdb.Database(jsr.data.conf.Mongo.Dbname).Collection("jinritemai_store_" + day)

	if _, err := collection.DeleteMany(ctx, bson.D{
		{"client_key", clientKey},
		{"open_id", openId},
	}); err != nil {
		return err
	}

	return nil
}

func (jsr *jinritemaiStoreRepo) Send(ctx context.Context, message event.Event) error {
	if err := jsr.data.kafka.Send(ctx, message); err != nil {
		return err
	}

	return nil
}

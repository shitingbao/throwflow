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

// 千川账户下已授权抖音号表
type QianchuanAweme struct {
	AwemeId                 uint64    `json:"aweme_id" bson:"aweme_id"`
	AdvertiserId            uint64    `json:"advertiser_id" bson:"advertiser_id"`
	AwemeAvatar             string    `json:"aweme_avatar" bson:"aweme_avatar"`
	AwemeShowId             string    `json:"aweme_show_id" bson:"aweme_show_id"`
	AwemeName               string    `json:"aweme_name" bson:"aweme_name"`
	AwemeStatus             string    `json:"aweme_status" bson:"aweme_status"`
	BindType                []*string `json:"bind_type" bson:"bind_type"`
	AwemeHasVideoPermission bool      `json:"aweme_has_video_permission" bson:"aweme_has_video_permission"`
	AwemeHasLivePermission  bool      `json:"aweme_has_live_permission" bson:"aweme_has_live_permission"`
	AwemeHasUniProm         bool      `json:"aweme_has_uni_prom" bson:"aweme_has_uni_prom"`
	CreateTime              time.Time `json:"create_time" bson:"create_time"`
	UpdateTime              time.Time `json:"update_time" bson:"update_time"`
}

type qianchuanAwemeRepo struct {
	data *Data
	log  *log.Helper
}

func NewQianchuanAwemeRepo(data *Data, logger log.Logger) biz.QianchuanAwemeRepo {
	return &qianchuanAwemeRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (qar *qianchuanAwemeRepo) SaveIndex(ctx context.Context, day string) {
	collection := qar.data.mdb.Database(qar.data.conf.Mongo.Dbname).Collection("qianchuan_aweme_" + day)

	isNotExistIndex := true

	if indexSpecifications, err := collection.Indexes().ListSpecifications(ctx); err == nil {
		for _, indexSpecification := range indexSpecifications {
			if indexSpecification.Name == "aweme_id_-1_advertiser_id_-1" {
				isNotExistIndex = false
			}
		}

		if isNotExistIndex {
			collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{Key: "aweme_id", Value: -1},
					{Key: "advertiser_id", Value: -1},
				},
			})
		}
	}
}

func (qar *qianchuanAwemeRepo) Upsert(ctx context.Context, day string, in *domain.QianchuanAweme) error {
	collection := qar.data.mdb.Database(qar.data.conf.Mongo.Dbname).Collection("qianchuan_aweme_" + day)

	if _, err := collection.UpdateOne(ctx, bson.D{
		{"aweme_id", in.AwemeId},
		{"advertiser_id", in.AdvertiserId},
	}, bson.D{
		{"$set", bson.D{
			{"aweme_avatar", in.AwemeAvatar},
			{"aweme_show_id", in.AwemeShowId},
			{"aweme_name", in.AwemeName},
			{"aweme_status", in.AwemeStatus},
			{"bind_type", in.BindType},
			{"aweme_has_video_permission", in.AwemeHasVideoPermission},
			{"aweme_has_live_permission", in.AwemeHasLivePermission},
			{"aweme_has_uni_prom", in.AwemeHasUniProm},
			{"create_time", in.CreateTime},
			{"update_time", in.UpdateTime},
		}},
	}, options.Update().SetUpsert(true)); err != nil {
		return err
	}

	return nil
}

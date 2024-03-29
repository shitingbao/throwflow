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

// 商家可投商品表
type QianchuanProduct struct {
	Id                  uint64           `json:"id" bson:"id"`
	AdvertiserId        uint64           `json:"advertiser_id" bson:"advertiser_id"`
	Name                string           `json:"name" bson:"name"`
	Img                 string           `json:"img" bson:"img"`
	CategoryName        string           `json:"category_name" bson:"category_name"`
	DiscountPrice       float64          `json:"discount_price" bson:"discount_price"`
	DiscountLowerPrice  float64          `json:"discount_lower_price" bson:"discount_lower_price"`
	DiscountHigherPrice float64          `json:"discount_higher_price" bson:"discount_higher_price"`
	ImgList             []*domain.ImgUrl `json:"img_list" bson:"img_list"`
	Inventory           uint64           `json:"inventory" bson:"inventory"`
	MarketPrice         float64          `json:"market_price" bson:"market_price"`
	ProductRate         float64          `json:"product_rate" bson:"product_rate"`
	SaleTime            string           `json:"sale_time" bson:"sale_time"`
	Tags                string           `json:"tags" bson:"tags"`
	CreateTime          time.Time        `json:"create_time" bson:"create_time"`
	UpdateTime          time.Time        `json:"update_time" bson:"update_time"`
}

type qianchuanProductRepo struct {
	data *Data
	log  *log.Helper
}

func NewQianchuanProductRepo(data *Data, logger log.Logger) biz.QianchuanProductRepo {
	return &qianchuanProductRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (qpr *qianchuanProductRepo) SaveIndex(ctx context.Context, day string) {
	collection := qpr.data.mdb.Database(qpr.data.conf.Mongo.Dbname).Collection("qianchuan_product_" + day)

	isNotExistIndex := true

	if indexSpecifications, err := collection.Indexes().ListSpecifications(ctx); err == nil {
		for _, indexSpecification := range indexSpecifications {
			if indexSpecification.Name == "id_-1_advertiser_id_-1" {
				isNotExistIndex = false
			}
		}

		if isNotExistIndex {
			collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{Key: "id", Value: -1},
					{Key: "advertiser_id", Value: -1},
				},
			})
		}
	}
}

func (qpr *qianchuanProductRepo) Upsert(ctx context.Context, day string, in *domain.QianchuanProduct) error {
	collection := qpr.data.mdb.Database(qpr.data.conf.Mongo.Dbname).Collection("qianchuan_product_" + day)

	if _, err := collection.UpdateOne(ctx, bson.D{
		{"id", in.Id},
		{"advertiser_id", in.AdvertiserId},
	}, bson.D{
		{"$set", bson.D{
			{"name", in.Name},
			{"img", in.Img},
			{"category_name", in.CategoryName},
			{"discount_price", in.DiscountPrice},
			{"discount_lower_price", in.DiscountLowerPrice},
			{"discount_higher_price", in.DiscountHigherPrice},
			{"img_list", in.ImgList},
			{"inventory", in.Inventory},
			{"market_price", in.MarketPrice},
			{"product_rate", in.ProductRate},
			{"sale_time", in.SaleTime},
			{"tags", in.Tags},
			{"create_time", in.CreateTime},
			{"update_time", in.UpdateTime},
		}},
	}, options.Update().SetUpsert(true)); err != nil {
		return err
	}

	return nil
}

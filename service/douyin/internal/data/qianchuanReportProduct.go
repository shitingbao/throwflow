package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// 千川账户商品数据表
type QianchuanReportProduct struct {
	AdvertiserId   uint64    `json:"advertiser_id" bson:"advertiser_id"`
	AdvertiserName string    `json:"advertiser_name" bson:"advertiser_name"`
	ProductId      uint64    `json:"product_id" bson:"product_id"`
	DiscountPrice  float64   `json:"discount_price" bson:"discount_price"`
	ProductName    string    `json:"product_name" bson:"product_name"`
	ProductImg     string    `json:"product_img" bson:"product_img"`
	StatCost       float64   `json:"stat_cost" bson:"stat_cost"`
	PayOrderCount  int64     `json:"pay_order_count" bson:"pay_order_count"`
	PayOrderAmount float64   `json:"pay_order_amount" bson:"pay_order_amount"`
	ShowCnt        int64     `json:"show_cnt" bson:"show_cnt"`
	ClickCnt       int64     `json:"click_cnt" bson:"click_cnt"`
	ConvertCnt     int64     `json:"convert_cnt" bson:"convert_cnt"`
	DyFollow       int64     `json:"dy_follow" bson:"dy_follow"`
	CreateTime     time.Time `json:"create_time" bson:"create_time"`
	UpdateTime     time.Time `json:"update_time" bson:"update_time"`
}

type TotalQianchuanReportProduct struct {
	Total int64 `json:"total" bson:"total"`
}

type StatisticsQianchuanReportProductPayOrderCount struct {
	Total int64 `json:"total" bson:"total"`
}

type StatisticsQianchuanReportProductPayOrderAmount struct {
	Total float64 `json:"total" bson:"total"`
}

type StatisticsQianchuanReportProductStatCost struct {
	Total float64 `json:"total" bson:"total"`
}

type qianchuanReportProductRepo struct {
	data *Data
	log  *log.Helper
}

func (qrp *QianchuanReportProduct) ToDomain() *domain.QianchuanReportProduct {
	return &domain.QianchuanReportProduct{
		AdvertiserId:   qrp.AdvertiserId,
		AdvertiserName: qrp.AdvertiserName,
		ProductId:      qrp.ProductId,
		DiscountPrice:  qrp.DiscountPrice,
		ProductName:    qrp.ProductName,
		ProductImg:     qrp.ProductImg,
		StatCost:       qrp.StatCost,
		PayOrderCount:  qrp.PayOrderCount,
		PayOrderAmount: qrp.PayOrderAmount,
		ShowCnt:        qrp.ShowCnt,
		ClickCnt:       qrp.ClickCnt,
		ConvertCnt:     qrp.ConvertCnt,
		DyFollow:       qrp.DyFollow,
		CreateTime:     qrp.CreateTime,
		UpdateTime:     qrp.UpdateTime,
	}
}

func NewQianchuanReportProductRepo(data *Data, logger log.Logger) biz.QianchuanReportProductRepo {
	return &qianchuanReportProductRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (qrpr *qianchuanReportProductRepo) List(ctx context.Context, advertiserIds, day, keyword string, isDistinction uint8, pageNum, pageSize uint64) ([]*domain.QianchuanReportProduct, error) {
	list := make([]*domain.QianchuanReportProduct, 0)

	collection := qrpr.data.mdb.Database(qrpr.data.conf.Mongo.Dbname).Collection("qianchuan_report_product_" + day)

	var matchStage bson.D
	var aadvertiserIds bson.A

	sadvertiserIds := strings.Split(advertiserIds, ",")

	for _, sadvertiserId := range sadvertiserIds {
		if uisadvertiserId, err := strconv.ParseUint(sadvertiserId, 10, 64); err == nil {
			aadvertiserIds = append(aadvertiserIds, uisadvertiserId)
		}
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		if productId, err := strconv.ParseUint(keyword, 10, 64); err == nil {
			matchStage = bson.D{
				{"$match", bson.M{
					"$and": []bson.M{
						bson.M{"advertiser_id": bson.D{{"$in", aadvertiserIds}}},
						bson.M{"$or": []bson.M{
							bson.M{"product_id": productId},
							bson.M{"product_name": bsonx.Regex(keyword, "")},
						}},
					},
				}},
			}
		} else {
			matchStage = bson.D{
				{"$match", bson.M{
					"$and": []bson.M{
						bson.M{"advertiser_id": bson.D{{"$in", aadvertiserIds}}},
						bson.M{"product_name": bsonx.Regex(keyword, "")},
					},
				}},
			}
		}
	} else {
		matchStage = bson.D{
			{"$match", bson.M{"advertiser_id": bson.D{{"$in", aadvertiserIds}}}},
		}
	}

	var groupStage bson.D

	if isDistinction == 1 {
		groupStage = bson.D{
			{"$group", bson.M{
				"_id":              bson.D{{"advertiser_id", "$advertiser_id"}, {"product_id", "$product_id"}},
				"advertiser_id":    bson.D{{"$first", "$advertiser_id"}},
				"advertiser_name":  bson.D{{"$first", "$advertiser_name"}},
				"product_id":       bson.D{{"$first", "$product_id"}},
				"discount_price":   bson.D{{"$first", "$discount_price"}},
				"product_name":     bson.D{{"$first", "$product_name"}},
				"product_img":      bson.D{{"$first", "$product_img"}},
				"pay_order_amount": bson.D{{"$sum", "$pay_order_amount"}},
				"pay_order_count":  bson.D{{"$sum", "$pay_order_count"}},
				"stat_cost":        bson.D{{"$sum", "$stat_cost"}},
				"click_cnt":        bson.D{{"$sum", "$click_cnt"}},
				"convert_cnt":      bson.D{{"$sum", "$convert_cnt"}},
				"show_cnt":         bson.D{{"$sum", "$show_cnt"}},
				"dy_follow":        bson.D{{"$sum", "$dy_follow"}},
			}},
		}
	} else {
		groupStage = bson.D{
			{"$group", bson.M{
				"_id":              bson.D{{"product_id", "$product_id"}},
				"advertiser_id":    bson.D{{"$first", "$advertiser_id"}},
				"advertiser_name":  bson.D{{"$first", "$advertiser_name"}},
				"product_id":       bson.D{{"$first", "$product_id"}},
				"discount_price":   bson.D{{"$first", "$discount_price"}},
				"product_name":     bson.D{{"$first", "$product_name"}},
				"product_img":      bson.D{{"$first", "$product_img"}},
				"pay_order_amount": bson.D{{"$sum", "$pay_order_amount"}},
				"pay_order_count":  bson.D{{"$sum", "$pay_order_count"}},
				"stat_cost":        bson.D{{"$sum", "$stat_cost"}},
				"click_cnt":        bson.D{{"$sum", "$click_cnt"}},
				"convert_cnt":      bson.D{{"$sum", "$convert_cnt"}},
				"show_cnt":         bson.D{{"$sum", "$show_cnt"}},
				"dy_follow":        bson.D{{"$sum", "$dy_follow"}},
			}},
		}
	}

	sortStage := bson.D{
		{"$sort", bson.M{"pay_order_amount": -1}},
	}

	var cursor *mongo.Cursor
	var err error

	if pageNum > 0 {
		skipStage := bson.D{
			{"$skip", (pageNum - 1) * pageSize},
		}

		limitStage := bson.D{
			{"$limit", pageSize},
		}

		cursor, err = collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, sortStage, skipStage, limitStage})
	} else {
		cursor, err = collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, sortStage})
	}

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var qianchuanReportProducts []*QianchuanReportProduct

	err = cursor.All(ctx, &qianchuanReportProducts)

	if err != nil {
		return nil, err
	}

	for _, qianchuanReportProduct := range qianchuanReportProducts {
		list = append(list, qianchuanReportProduct.ToDomain())
	}

	return list, nil
}

func (qrpr *qianchuanReportProductRepo) Count(ctx context.Context, advertiserIds, day, keyword string, isDistinction uint8) (int64, error) {
	collection := qrpr.data.mdb.Database(qrpr.data.conf.Mongo.Dbname).Collection("qianchuan_report_product_" + day)

	var matchStage bson.D
	var aadvertiserIds bson.A

	sadvertiserIds := strings.Split(advertiserIds, ",")

	for _, sadvertiserId := range sadvertiserIds {
		if uisadvertiserId, err := strconv.ParseUint(sadvertiserId, 10, 64); err == nil {
			aadvertiserIds = append(aadvertiserIds, uisadvertiserId)
		}
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		if productId, err := strconv.ParseUint(keyword, 10, 64); err == nil {
			matchStage = bson.D{
				{"$match", bson.M{
					"$and": []bson.M{
						bson.M{"advertiser_id": bson.D{{"$in", aadvertiserIds}}},
						bson.M{"$or": []bson.M{
							bson.M{"product_id": productId},
							bson.M{"product_name": bsonx.Regex(keyword, "")},
						}},
					},
				}},
			}
		} else {
			matchStage = bson.D{
				{"$match", bson.M{"advertiser_id": bson.D{{"$in", aadvertiserIds}}, "product_name": bson.D{{"$regex", keyword}}}},
			}
		}
	} else {
		matchStage = bson.D{
			{"$match", bson.M{"advertiser_id": bson.D{{"$in", aadvertiserIds}}}},
		}
	}

	var groupStage bson.D

	if isDistinction == 1 {
		groupStage = bson.D{
			{"$group", bson.M{
				"_id": bson.D{{"advertiser_id", "$advertiser_id"}, {"product_id", "$product_id"}},
			}},
		}
	} else {
		groupStage = bson.D{
			{"$group", bson.M{
				"_id": bson.D{{"product_id", "$product_id"}},
			}},
		}
	}

	countStage := bson.D{
		{"$count", "total"},
	}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, countStage})

	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)

	var totalQianchuanReportProducts []*TotalQianchuanReportProduct

	err = cursor.All(ctx, &totalQianchuanReportProducts)

	if err != nil {
		return 0, err
	}

	var total int64

	for _, totalQianchuanReportProduct := range totalQianchuanReportProducts {
		total = totalQianchuanReportProduct.Total
	}

	return total, nil
}

func (qrpr *qianchuanReportProductRepo) StatisticsPayOrderCount(ctx context.Context, advertiserIds, day string) (int64, error) {
	collection := qrpr.data.mdb.Database(qrpr.data.conf.Mongo.Dbname).Collection("qianchuan_report_product_" + day)

	var matchStage bson.D
	var aadvertiserIds bson.A

	sadvertiserIds := strings.Split(advertiserIds, ",")

	for _, sadvertiserId := range sadvertiserIds {
		if uisadvertiserId, err := strconv.ParseUint(sadvertiserId, 10, 64); err == nil {
			aadvertiserIds = append(aadvertiserIds, uisadvertiserId)
		}
	}

	matchStage = bson.D{
		{"$match", bson.M{"advertiser_id": bson.D{{"$in", aadvertiserIds}}}},
	}

	groupStage := bson.D{
		{"$group", bson.M{
			"_id":   "",
			"total": bson.D{{"$sum", "$pay_order_count"}},
		}},
	}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})

	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)

	var statisticsQianchuanReportProductPayOrderCounts []*StatisticsQianchuanReportProductPayOrderCount

	err = cursor.All(ctx, &statisticsQianchuanReportProductPayOrderCounts)

	if err != nil {
		return 0, err
	}

	var total int64

	for _, statisticsQianchuanReportProductPayOrderCount := range statisticsQianchuanReportProductPayOrderCounts {
		total = statisticsQianchuanReportProductPayOrderCount.Total

		break
	}

	return total, nil
}

func (qrpr *qianchuanReportProductRepo) StatisticsPayOrderAmount(ctx context.Context, advertiserIds, day string) (float64, error) {
	collection := qrpr.data.mdb.Database(qrpr.data.conf.Mongo.Dbname).Collection("qianchuan_report_product_" + day)

	var matchStage bson.D
	var aadvertiserIds bson.A

	sadvertiserIds := strings.Split(advertiserIds, ",")

	for _, sadvertiserId := range sadvertiserIds {
		if uisadvertiserId, err := strconv.ParseUint(sadvertiserId, 10, 64); err == nil {
			aadvertiserIds = append(aadvertiserIds, uisadvertiserId)
		}
	}

	matchStage = bson.D{
		{"$match", bson.M{"advertiser_id": bson.D{{"$in", aadvertiserIds}}}},
	}

	groupStage := bson.D{
		{"$group", bson.M{
			"_id":   "",
			"total": bson.D{{"$sum", "$pay_order_amount"}},
		}},
	}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})

	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)

	var statisticsQianchuanReportProductPayOrderAmounts []*StatisticsQianchuanReportProductPayOrderAmount

	err = cursor.All(ctx, &statisticsQianchuanReportProductPayOrderAmounts)

	if err != nil {
		return 0, err
	}

	var total float64

	for _, statisticsQianchuanReportProductPayOrderAmount := range statisticsQianchuanReportProductPayOrderAmounts {
		total = statisticsQianchuanReportProductPayOrderAmount.Total

		break
	}

	return total, nil
}

func (qrpr *qianchuanReportProductRepo) StatisticsStatCost(ctx context.Context, advertiserIds, day string) (float64, error) {
	collection := qrpr.data.mdb.Database(qrpr.data.conf.Mongo.Dbname).Collection("qianchuan_report_product_" + day)

	var matchStage bson.D
	var aadvertiserIds bson.A

	sadvertiserIds := strings.Split(advertiserIds, ",")

	for _, sadvertiserId := range sadvertiserIds {
		if uisadvertiserId, err := strconv.ParseUint(sadvertiserId, 10, 64); err == nil {
			aadvertiserIds = append(aadvertiserIds, uisadvertiserId)
		}
	}

	matchStage = bson.D{
		{"$match", bson.M{"advertiser_id": bson.D{{"$in", aadvertiserIds}}}},
	}

	groupStage := bson.D{
		{"$group", bson.M{
			"_id":   "",
			"total": bson.D{{"$sum", "$stat_cost"}},
		}},
	}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})

	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)

	var statisticsQianchuanReportProductStatCosts []*StatisticsQianchuanReportProductStatCost

	err = cursor.All(ctx, &statisticsQianchuanReportProductStatCosts)

	if err != nil {
		return 0, err
	}

	var total float64

	for _, statisticsQianchuanReportProductStatCost := range statisticsQianchuanReportProductStatCosts {
		total = statisticsQianchuanReportProductStatCost.Total

		break
	}

	return total, nil
}

func (qrpr *qianchuanReportProductRepo) SaveIndex(ctx context.Context, day string) {
	collection := qrpr.data.mdb.Database(qrpr.data.conf.Mongo.Dbname).Collection("qianchuan_report_product_" + day)

	isNotExistIndex := true

	if indexSpecifications, err := collection.Indexes().ListSpecifications(ctx); err == nil {
		for _, indexSpecification := range indexSpecifications {
			if indexSpecification.Name == "advertiser_id_-1_product_id_-1" {
				isNotExistIndex = false
			}
		}

		if isNotExistIndex {
			collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{Key: "advertiser_id", Value: -1},
					{Key: "product_id", Value: -1},
				},
			})
		}
	}
}

func (qrpr *qianchuanReportProductRepo) UpsertQianchuanReportProductInfo(ctx context.Context, day string, in *domain.QianchuanReportProduct) error {
	collection := qrpr.data.mdb.Database(qrpr.data.conf.Mongo.Dbname).Collection("qianchuan_report_product_" + day)

	if _, err := collection.UpdateOne(ctx, bson.D{
		{"advertiser_id", in.AdvertiserId},
		{"product_id", in.ProductId},
	}, bson.D{
		{"$setOnInsert", bson.D{
			{"stat_cost", in.StatCost},
			{"pay_order_count", in.PayOrderCount},
			{"pay_order_amount", in.PayOrderAmount},
			{"show_cnt", in.ShowCnt},
			{"click_cnt", in.ClickCnt},
			{"convert_cnt", in.ConvertCnt},
			{"dy_follow", in.DyFollow},
			{"create_time", in.CreateTime},
		}},
		{"$set", bson.D{
			{"advertiser_name", in.AdvertiserName},
			{"discount_price", in.DiscountPrice},
			{"product_name", in.ProductName},
			{"product_img", in.ProductImg},
			{"update_time", in.UpdateTime},
		}},
	}, options.Update().SetUpsert(true)); err != nil {
		return err
	}

	return nil
}

func (qrpr *qianchuanReportProductRepo) UpsertQianchuanReportProduct(ctx context.Context, day string, in *domain.QianchuanReportProduct) error {
	collection := qrpr.data.mdb.Database(qrpr.data.conf.Mongo.Dbname).Collection("qianchuan_report_product_" + day)

	if _, err := collection.UpdateOne(ctx, bson.D{
		{"advertiser_id", in.AdvertiserId},
		{"product_id", in.ProductId},
	}, bson.D{
		{"$setOnInsert", bson.D{
			{"create_time", in.CreateTime},
		}},
		{"$set", bson.D{
			{"discount_price", in.DiscountPrice},
			{"product_name", in.ProductName},
			{"product_img", in.ProductImg},
			{"stat_cost", in.StatCost},
			{"pay_order_count", in.PayOrderCount},
			{"pay_order_amount", in.PayOrderAmount},
			{"show_cnt", in.ShowCnt},
			{"click_cnt", in.ClickCnt},
			{"convert_cnt", in.ConvertCnt},
			{"dy_follow", in.DyFollow},
			{"update_time", in.UpdateTime},
		}},
	}, options.Update().SetUpsert(true)); err != nil {
		return err
	}

	return nil
}

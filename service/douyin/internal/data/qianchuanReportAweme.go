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

// 千川账户达人数据表
type QianchuanReportAweme struct {
	AdvertiserId   uint64    `json:"advertiser_id" bson:"advertiser_id"`
	AdvertiserName string    `json:"advertiser_name" bson:"advertiser_name"`
	AwemeId        uint64    `json:"aweme_id" bson:"aweme_id"`
	AwemeName      string    `json:"aweme_name" bson:"aweme_name"`
	AwemeShowId    string    `json:"aweme_show_id" bson:"aweme_show_id"`
	AwemeAvatar    string    `json:"aweme_avatar" bson:"aweme_avatar"`
	DyFollow       int64     `json:"dy_follow" bson:"dy_follow"`
	StatCost       float64   `json:"stat_cost" bson:"stat_cost"`
	PayOrderCount  int64     `json:"pay_order_count" bson:"pay_order_count"`
	PayOrderAmount float64   `json:"pay_order_amount" bson:"pay_order_amount"`
	ShowCnt        int64     `json:"show_cnt" bson:"show_cnt"`
	ClickCnt       int64     `json:"click_cnt" bson:"click_cnt"`
	ConvertCnt     int64     `json:"convert_cnt" bson:"convert_cnt"`
	CreateTime     time.Time `json:"create_time" bson:"create_time"`
	UpdateTime     time.Time `json:"update_time" bson:"update_time"`
}

type TotalQianchuanReportAweme struct {
	Total int64 `json:"total" bson:"total"`
}

type StatisticsQianchuanReportAwemePayOrderCount struct {
	Total int64 `json:"total" bson:"total"`
}

type StatisticsQianchuanReportAwemePayOrderAmount struct {
	Total float64 `json:"total" bson:"total"`
}

type StatisticsQianchuanReportAwemeStatCost struct {
	Total float64 `json:"total" bson:"total"`
}

type qianchuanReportAwemeRepo struct {
	data *Data
	log  *log.Helper
}

func (qra *QianchuanReportAweme) ToDomain() *domain.QianchuanReportAweme {
	return &domain.QianchuanReportAweme{
		AdvertiserId:   qra.AdvertiserId,
		AdvertiserName: qra.AdvertiserName,
		AwemeId:        qra.AwemeId,
		AwemeName:      qra.AwemeName,
		AwemeShowId:    qra.AwemeShowId,
		AwemeAvatar:    qra.AwemeAvatar,
		DyFollow:       qra.DyFollow,
		StatCost:       qra.StatCost,
		PayOrderCount:  qra.PayOrderCount,
		PayOrderAmount: qra.PayOrderAmount,
		ShowCnt:        qra.ShowCnt,
		ClickCnt:       qra.ClickCnt,
		ConvertCnt:     qra.ConvertCnt,
		CreateTime:     qra.CreateTime,
		UpdateTime:     qra.UpdateTime,
	}
}

func NewQianchuanReportAwemeRepo(data *Data, logger log.Logger) biz.QianchuanReportAwemeRepo {
	return &qianchuanReportAwemeRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (qrar *qianchuanReportAwemeRepo) List(ctx context.Context, advertiserIds, day, keyword string, isDistinction uint8, pageNum, pageSize uint64) ([]*domain.QianchuanReportAweme, error) {
	list := make([]*domain.QianchuanReportAweme, 0)

	collection := qrar.data.mdb.Database(qrar.data.conf.Mongo.Dbname).Collection("qianchuan_report_aweme_" + day)

	var matchStage bson.D
	var aadvertiserIds bson.A

	sadvertiserIds := strings.Split(advertiserIds, ",")

	for _, sadvertiserId := range sadvertiserIds {
		if uisadvertiserId, err := strconv.ParseUint(sadvertiserId, 10, 64); err == nil {
			aadvertiserIds = append(aadvertiserIds, uisadvertiserId)
		}
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		if awemeId, err := strconv.ParseUint(keyword, 10, 64); err == nil {
			matchStage = bson.D{
				{"$match", bson.M{
					"$and": []bson.M{
						bson.M{"advertiser_id": bson.D{{"$in", aadvertiserIds}}},
						bson.M{"$or": []bson.M{
							bson.M{"aweme_id": awemeId},
							bson.M{"aweme_name": bsonx.Regex(keyword, "")},
						}},
					},
				}},
			}
		} else {
			matchStage = bson.D{
				{"$match", bson.M{
					"$and": []bson.M{
						bson.M{"advertiser_id": bson.D{{"$in", aadvertiserIds}}},
						bson.M{"aweme_name": bsonx.Regex(keyword, "")},
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
				"_id":              bson.D{{"advertiser_id", "$advertiser_id"}, {"aweme_id", "$aweme_id"}},
				"advertiser_id":    bson.D{{"$first", "$advertiser_id"}},
				"advertiser_name":  bson.D{{"$first", "$advertiser_name"}},
				"aweme_id":         bson.D{{"$first", "$aweme_id"}},
				"aweme_avatar":     bson.D{{"$first", "$aweme_avatar"}},
				"aweme_name":       bson.D{{"$first", "$aweme_name"}},
				"aweme_show_id":    bson.D{{"$first", "$aweme_show_id"}},
				"dy_follow":        bson.D{{"$sum", "$dy_follow"}},
				"pay_order_amount": bson.D{{"$sum", "$pay_order_amount"}},
				"pay_order_count":  bson.D{{"$sum", "$pay_order_count"}},
				"stat_cost":        bson.D{{"$sum", "$stat_cost"}},
				"click_cnt":        bson.D{{"$sum", "$click_cnt"}},
				"convert_cnt":      bson.D{{"$sum", "$convert_cnt"}},
				"show_cnt":         bson.D{{"$sum", "$show_cnt"}},
			}},
		}
	} else {
		groupStage = bson.D{
			{"$group", bson.M{
				"_id":              bson.D{{"aweme_id", "$aweme_id"}},
				"advertiser_id":    bson.D{{"$first", "$advertiser_id"}},
				"advertiser_name":  bson.D{{"$first", "$advertiser_name"}},
				"aweme_id":         bson.D{{"$first", "$aweme_id"}},
				"aweme_avatar":     bson.D{{"$first", "$aweme_avatar"}},
				"aweme_name":       bson.D{{"$first", "$aweme_name"}},
				"aweme_show_id":    bson.D{{"$first", "$aweme_show_id"}},
				"dy_follow":        bson.D{{"$sum", "$dy_follow"}},
				"pay_order_amount": bson.D{{"$sum", "$pay_order_amount"}},
				"pay_order_count":  bson.D{{"$sum", "$pay_order_count"}},
				"stat_cost":        bson.D{{"$sum", "$stat_cost"}},
				"click_cnt":        bson.D{{"$sum", "$click_cnt"}},
				"convert_cnt":      bson.D{{"$sum", "$convert_cnt"}},
				"show_cnt":         bson.D{{"$sum", "$show_cnt"}},
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

	var qianchuanReportAwemes []*QianchuanReportAweme

	err = cursor.All(ctx, &qianchuanReportAwemes)

	if err != nil {
		return nil, err
	}

	for _, qianchuanReportAweme := range qianchuanReportAwemes {
		list = append(list, qianchuanReportAweme.ToDomain())
	}

	return list, nil
}

func (qrar *qianchuanReportAwemeRepo) Count(ctx context.Context, advertiserIds, day, keyword string, isDistinction uint8) (int64, error) {
	collection := qrar.data.mdb.Database(qrar.data.conf.Mongo.Dbname).Collection("qianchuan_report_aweme_" + day)

	var matchStage bson.D
	var aadvertiserIds bson.A

	sadvertiserIds := strings.Split(advertiserIds, ",")

	for _, sadvertiserId := range sadvertiserIds {
		if uisadvertiserId, err := strconv.ParseUint(sadvertiserId, 10, 64); err == nil {
			aadvertiserIds = append(aadvertiserIds, uisadvertiserId)
		}
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		if awemeId, err := strconv.ParseUint(keyword, 10, 64); err == nil {
			matchStage = bson.D{
				{"$match", bson.M{
					"$and": []bson.M{
						bson.M{"advertiser_id": bson.D{{"$in", aadvertiserIds}}},
						bson.M{"$or": []bson.M{
							bson.M{"aweme_id": awemeId},
							bson.M{"aweme_name": bsonx.Regex(keyword, "")},
						}},
					},
				}},
			}
		} else {
			matchStage = bson.D{
				{"$match", bson.M{"advertiser_id": bson.D{{"$in", aadvertiserIds}}, "aweme_name": bson.D{{"$regex", keyword}}}},
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
				"_id": bson.D{{"advertiser_id", "$advertiser_id"}, {"aweme_id", "$aweme_id"}},
			}},
		}
	} else {
		groupStage = bson.D{
			{"$group", bson.M{
				"_id": bson.D{{"aweme_id", "$aweme_id"}},
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

	var totalQianchuanReportAwemes []*TotalQianchuanReportAweme

	err = cursor.All(ctx, &totalQianchuanReportAwemes)

	if err != nil {
		return 0, err
	}

	var total int64

	for _, totalQianchuanReportAweme := range totalQianchuanReportAwemes {
		total = totalQianchuanReportAweme.Total
	}

	return total, nil
}

func (qrar *qianchuanReportAwemeRepo) StatisticsPayOrderCount(ctx context.Context, advertiserIds, day string) (int64, error) {
	collection := qrar.data.mdb.Database(qrar.data.conf.Mongo.Dbname).Collection("qianchuan_report_aweme_" + day)

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

	var statisticsQianchuanReportAwemePayOrderCounts []*StatisticsQianchuanReportAwemePayOrderCount

	err = cursor.All(ctx, &statisticsQianchuanReportAwemePayOrderCounts)

	if err != nil {
		return 0, err
	}

	var total int64

	for _, statisticsQianchuanReportAwemePayOrderCount := range statisticsQianchuanReportAwemePayOrderCounts {
		total = statisticsQianchuanReportAwemePayOrderCount.Total

		break
	}

	return total, nil
}

func (qrar *qianchuanReportAwemeRepo) StatisticsPayOrderAmount(ctx context.Context, advertiserIds, day string) (float64, error) {
	collection := qrar.data.mdb.Database(qrar.data.conf.Mongo.Dbname).Collection("qianchuan_report_aweme_" + day)

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

	var statisticsQianchuanReportAwemePayOrderAmounts []*StatisticsQianchuanReportAwemePayOrderAmount

	err = cursor.All(ctx, &statisticsQianchuanReportAwemePayOrderAmounts)

	if err != nil {
		return 0, err
	}

	var total float64

	for _, statisticsQianchuanReportAwemePayOrderAmount := range statisticsQianchuanReportAwemePayOrderAmounts {
		total = statisticsQianchuanReportAwemePayOrderAmount.Total

		break
	}

	return total, nil
}

func (qrar *qianchuanReportAwemeRepo) StatisticsStatCost(ctx context.Context, advertiserIds, day string) (float64, error) {
	collection := qrar.data.mdb.Database(qrar.data.conf.Mongo.Dbname).Collection("qianchuan_report_aweme_" + day)

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

	var statisticsQianchuanReportAwemeStatCosts []*StatisticsQianchuanReportAwemeStatCost

	err = cursor.All(ctx, &statisticsQianchuanReportAwemeStatCosts)

	if err != nil {
		return 0, err
	}

	var total float64

	for _, statisticsQianchuanReportAwemeStatCost := range statisticsQianchuanReportAwemeStatCosts {
		total = statisticsQianchuanReportAwemeStatCost.Total

		break
	}

	return total, nil
}

func (qrar *qianchuanReportAwemeRepo) SaveIndex(ctx context.Context, day string) {
	collection := qrar.data.mdb.Database(qrar.data.conf.Mongo.Dbname).Collection("qianchuan_report_aweme_" + day)

	isNotExistIndex := true

	if indexSpecifications, err := collection.Indexes().ListSpecifications(ctx); err == nil {
		for _, indexSpecification := range indexSpecifications {
			if indexSpecification.Name == "advertiser_id_-1_aweme_id_-1" {
				isNotExistIndex = false
			}
		}

		if isNotExistIndex {
			collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{Key: "advertiser_id", Value: -1},
					{Key: "aweme_id", Value: -1},
				},
			})
		}
	}
}

func (qrar *qianchuanReportAwemeRepo) UpsertQianchuanReportAwemeInfo(ctx context.Context, day string, in *domain.QianchuanReportAweme) error {
	collection := qrar.data.mdb.Database(qrar.data.conf.Mongo.Dbname).Collection("qianchuan_report_aweme_" + day)

	if _, err := collection.UpdateOne(ctx, bson.D{
		{"advertiser_id", in.AdvertiserId},
		{"aweme_id", in.AwemeId},
	}, bson.D{
		{"$setOnInsert", bson.D{
			{"dy_follow", in.DyFollow},
			{"stat_cost", in.StatCost},
			{"pay_order_count", in.PayOrderCount},
			{"pay_order_amount", in.PayOrderAmount},
			{"show_cnt", in.ShowCnt},
			{"click_cnt", in.ClickCnt},
			{"convert_cnt", in.ConvertCnt},
			{"create_time", in.CreateTime},
		}},
		{"$set", bson.D{
			{"advertiser_name", in.AdvertiserName},
			{"aweme_name", in.AwemeName},
			{"aweme_show_id", in.AwemeShowId},
			{"aweme_avatar", in.AwemeAvatar},
			{"update_time", in.UpdateTime},
		}},
	}, options.Update().SetUpsert(true)); err != nil {
		return err
	}

	return nil
}

func (qrar *qianchuanReportAwemeRepo) UpsertQianchuanReportAweme(ctx context.Context, day string, in *domain.QianchuanReportAweme) error {
	collection := qrar.data.mdb.Database(qrar.data.conf.Mongo.Dbname).Collection("qianchuan_report_aweme_" + day)

	if _, err := collection.UpdateOne(ctx, bson.D{
		{"advertiser_id", in.AdvertiserId},
		{"aweme_id", in.AwemeId},
	}, bson.D{
		{"$setOnInsert", bson.D{
			{"create_time", in.CreateTime},
		}},
		{"$set", bson.D{
			{"aweme_name", in.AwemeName},
			{"aweme_show_id", in.AwemeShowId},
			{"aweme_avatar", in.AwemeAvatar},
			{"dy_follow", in.DyFollow},
			{"stat_cost", in.StatCost},
			{"pay_order_count", in.PayOrderCount},
			{"pay_order_amount", in.PayOrderAmount},
			{"show_cnt", in.ShowCnt},
			{"click_cnt", in.ClickCnt},
			{"convert_cnt", in.ConvertCnt},
			{"update_time", in.UpdateTime},
		}},
	}, options.Update().SetUpsert(true)); err != nil {
		return err
	}

	return nil
}

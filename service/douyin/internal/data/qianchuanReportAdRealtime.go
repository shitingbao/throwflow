package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
	"time"
)

// 连山云RDS-千川计划数据实时表
type QianchuanReportAdRealtime struct {
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
	Time              int64     `json:"time" bson:"time"`
	ConvertCnt        int64     `json:"convert_cnt" bson:"convert_cnt"`
	CreateTime        time.Time `json:"create_time" bson:"create_time"`
	UpdateTime        time.Time `json:"update_time" bson:"update_time"`
}

type StatisticsAdvertiserCount struct {
	Total int64 `json:"total" bson:"total"`
}

type qianchuanReportAdRealtimeRepo struct {
	data *Data
	log  *log.Helper
}

func (qrar *QianchuanReportAdRealtime) ToDomain() *domain.QianchuanReportAdRealtime {
	return &domain.QianchuanReportAdRealtime{
		AdId:              qrar.AdId,
		AdvertiserId:      qrar.AdvertiserId,
		AwemeId:           qrar.AwemeId,
		MarketingGoal:     qrar.MarketingGoal,
		StatCost:          qrar.StatCost,
		ShowCnt:           qrar.ShowCnt,
		ClickCnt:          qrar.ClickCnt,
		PayOrderCount:     qrar.PayOrderCount,
		CreateOrderAmount: qrar.CreateOrderAmount,
		CreateOrderCount:  qrar.CreateOrderCount,
		PayOrderAmount:    qrar.PayOrderAmount,
		DyFollow:          qrar.DyFollow,
		ConvertCnt:        qrar.ConvertCnt,
		Time:              qrar.Time,
		CreateTime:        qrar.CreateTime,
		UpdateTime:        qrar.UpdateTime,
	}
}

func NewQianchuanReportAdRealtimeRepo(data *Data, logger log.Logger) biz.QianchuanReportAdRealtimeRepo {
	return &qianchuanReportAdRealtimeRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (qrar *qianchuanReportAdRealtimeRepo) Get(ctx context.Context, adId uint64, day string) (*domain.QianchuanReportAdRealtime, error) {
	var qianchuanReportAdRealtime QianchuanReportAdRealtime

	collection := qrar.data.mdb.Database(qrar.data.conf.Mongo.Dbname).Collection("qianchuan_report_ad_realtime_" + day)

	if err := collection.FindOne(ctx, bson.M{
		"$and": []bson.M{
			bson.M{"ad_id": adId},
		},
	}, &options.FindOneOptions{
		Sort: bson.D{
			bson.E{"time", -1},
		},
	}).Decode(&qianchuanReportAdRealtime); err != nil {
		return nil, err
	}

	return qianchuanReportAdRealtime.ToDomain(), nil
}

func (qrar *qianchuanReportAdRealtimeRepo) GetByTime(ctx context.Context, adId uint64, time int64, day string) (*domain.QianchuanReportAdRealtime, error) {
	var qianchuanReportAdRealtime QianchuanReportAdRealtime

	collection := qrar.data.mdb.Database(qrar.data.conf.Mongo.Dbname).Collection("qianchuan_report_ad_realtime_" + day)

	if err := collection.FindOne(ctx, bson.M{
		"$and": []bson.M{
			bson.M{"ad_id": adId},
			bson.M{"time": time},
		},
	}, &options.FindOneOptions{
		Sort: bson.D{
			bson.E{"time", -1},
		},
	}).Decode(&qianchuanReportAdRealtime); err != nil {
		return nil, err
	}

	return qianchuanReportAdRealtime.ToDomain(), nil
}

func (qrar *qianchuanReportAdRealtimeRepo) List(ctx context.Context, adId uint64, day string) ([]*domain.QianchuanReportAdRealtime, error) {
	list := make([]*domain.QianchuanReportAdRealtime, 0)
	var qianchuanReportAdRealtimes []*QianchuanReportAdRealtime

	collection := qrar.data.mdb.Database(qrar.data.conf.Mongo.Dbname).Collection("qianchuan_report_ad_realtime_" + day)

	cursor, err := collection.Find(ctx, bson.M{
		"$and": []bson.M{
			bson.M{"ad_id": adId},
		},
	}, &options.FindOptions{
		Sort: bson.D{
			bson.E{"time", 1},
		},
	})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &qianchuanReportAdRealtimes)

	if err != nil {
		return nil, err
	}

	for _, qianchuanReportAdRealtime := range qianchuanReportAdRealtimes {
		list = append(list, qianchuanReportAdRealtime.ToDomain())
	}

	return list, nil
}

func (qrar *qianchuanReportAdRealtimeRepo) ListAdvertisers(ctx context.Context, advertiserIds, day string) ([]*domain.QianchuanReportAdRealtime, error) {
	list := make([]*domain.QianchuanReportAdRealtime, 0)

	where := make([]string, 0)

	where = append(where, "day='"+day+"'")
	where = append(where, "advertiser_id in ("+advertiserIds+")")

	sql := "SELECT time,sum(stat_cost) as stat_cost,sum(pay_order_amount) pay_order_amount FROM qianchuan_report_ad_realtime WHERE " + strings.Join(where, " AND ") + " group by time order by time asc"

	rows, err := qrar.data.cdb.Query(ctx, sql)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			time           int64
			statCost       float64
			payOrderAmount float64
		)

		if err := rows.Scan(&time, &statCost, &payOrderAmount); err != nil {
			return nil, err
		}

		list = append(list, &domain.QianchuanReportAdRealtime{
			Time:           time,
			StatCost:       statCost,
			PayOrderAmount: payOrderAmount,
		})
	}

	rows.Close()

	return list, nil
}

func (qrar *qianchuanReportAdRealtimeRepo) ListByAdIds(ctx context.Context, adIds []uint64, time int64, day string) ([]*domain.QianchuanReportAdRealtime, error) {
	list := make([]*domain.QianchuanReportAdRealtime, 0)
	var qianchuanReportAdRealtimes []*QianchuanReportAdRealtime

	collection := qrar.data.mdb.Database(qrar.data.conf.Mongo.Dbname).Collection("qianchuan_report_ad_realtime_" + day)

	var aadIds bson.A

	for _, adId := range adIds {
		aadIds = append(aadIds, adId)
	}

	cursor, err := collection.Find(ctx, bson.M{
		"$and": []bson.M{
			bson.M{"ad_id": bson.D{{"$in", aadIds}}},
			bson.M{"time": time},
		},
	})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &qianchuanReportAdRealtimes)

	if err != nil {
		return nil, err
	}

	for _, qianchuanReportAdRealtime := range qianchuanReportAdRealtimes {
		list = append(list, qianchuanReportAdRealtime.ToDomain())
	}

	return list, nil
}

func (qrar *qianchuanReportAdRealtimeRepo) StatisticsAdvertisers(ctx context.Context, advertiserIds, day string) (int64, error) {
	collection := qrar.data.mdb.Database(qrar.data.conf.Mongo.Dbname).Collection("qianchuan_report_ad_realtime_" + day)

	var aadvertiserIds bson.A

	sadvertiserIds := strings.Split(advertiserIds, ",")

	for _, sadvertiserId := range sadvertiserIds {
		if uisadvertiserId, err := strconv.ParseUint(sadvertiserId, 10, 64); err == nil {
			aadvertiserIds = append(aadvertiserIds, uisadvertiserId)
		}
	}

	matchStage := bson.D{
		{"$match", bson.M{"advertiser_id": bson.D{{"$in", aadvertiserIds}}, "stat_cost": bson.D{{"$gt", 0}}}},
	}

	groupAStage := bson.D{
		{"$group", bson.M{
			"_id": "$advertiser_id",
		}},
	}

	projectAStage := bson.D{
		{"$project", bson.M{
			"_id":           0,
			"advertiser_id": "$_id",
		}},
	}

	groupBStage := bson.D{
		{"$group", bson.M{
			"_id":   "$advertiser_id",
			"total": bson.D{{"$sum", 1}},
		}},
	}

	projectBStage := bson.D{
		{"$project", bson.M{
			"_id":   0,
			"total": 1,
		}},
	}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupAStage, projectAStage, groupBStage, projectBStage})

	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)

	var statisticsAdvertiserCounts []*StatisticsAdvertiserCount

	err = cursor.All(ctx, &statisticsAdvertiserCounts)

	if err != nil {
		return 0, err
	}

	var total int64

	for _, statisticsAdvertiserCount := range statisticsAdvertiserCounts {
		total = statisticsAdvertiserCount.Total

		break
	}

	return total, nil
}

func (qrar *qianchuanReportAdRealtimeRepo) SaveIndex(ctx context.Context, day string) {
	collection := qrar.data.mdb.Database(qrar.data.conf.Mongo.Dbname).Collection("qianchuan_report_ad_realtime_" + day)

	isNotExistAIndex := true
	isNotExistBIndex := true

	if indexSpecifications, err := collection.Indexes().ListSpecifications(ctx); err == nil {
		for _, indexSpecification := range indexSpecifications {
			if indexSpecification.Name == "advertiser_id_-1_ad_id_-1_time_-1" {
				isNotExistAIndex = false
			}

			if indexSpecification.Name == "ad_id_-1_time_-1" {
				isNotExistBIndex = false
			}
		}

		if isNotExistAIndex {
			collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{Key: "advertiser_id", Value: -1},
					{Key: "ad_id", Value: -1},
					{Key: "time", Value: -1},
				},
			})
		}

		if isNotExistBIndex {
			collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{Key: "ad_id", Value: -1},
					{Key: "time", Value: -1},
				},
			})
		}
	}
}

func (qrar *qianchuanReportAdRealtimeRepo) Upsert(ctx context.Context, day string, in *domain.QianchuanReportAdRealtime) error {
	collection := qrar.data.mdb.Database(qrar.data.conf.Mongo.Dbname).Collection("qianchuan_report_ad_realtime_" + day)

	if _, err := collection.UpdateOne(ctx, bson.D{
		{"ad_id", in.AdId},
		{"advertiser_id", in.AdvertiserId},
		{"time", in.Time},
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

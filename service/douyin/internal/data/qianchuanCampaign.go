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

// 千川广告组表
type QianchuanCampaign struct {
	Id             uint64    `json:"id" bson:"id"`
	AdvertiserId   uint64    `json:"advertiser_id" bson:"advertiser_id"`
	Name           string    `json:"name" bson:"name"`
	Budget         float64   `json:"budget" bson:"budget"`
	BudgetMode     string    `json:"budget_mode" bson:"budget_mode"`
	MarketingGoal  string    `json:"marketing_goal" bson:"marketing_goal"`
	MarketingScene string    `json:"marketing_scene" bson:"marketing_scene"`
	Status         string    `json:"status" bson:"status"`
	CreateDate     string    `json:"create_date" bson:"create_date"`
	CreateTime     time.Time `json:"create_time" bson:"create_time"`
	UpdateTime     time.Time `json:"update_time" bson:"update_time"`
}

type TotalQianchuanCampaign struct {
	Total int64 `json:"total" bson:"total"`
}

type CountQianchuanCampaign struct {
	AdvertiserId uint64 `json:"advertiser_id" bson:"advertiser_id"`
	Total        int64  `json:"total" bson:"total"`
}

type ExternalQianchuanCampaign struct {
	Id                uint64               `json:"id" bson:"id"`
	Name              string               `json:"name" bson:"name"`
	Budget            float64              `json:"budget" bson:"budget"`
	BudgetMode        string               `json:"budget_mode" bson:"budget_mode"`
	MarketingGoal     string               `json:"marketing_goal" bson:"marketing_goal"`
	MarketingScene    string               `json:"marketing_scene" bson:"marketing_scene"`
	Status            string               `json:"status" bson:"status"`
	CreateDate        string               `json:"create_date" bson:"create_date"`
	QianchuanReportAd []*QianchuanReportAd `json:"qianchuan_report_ad" bson:"qianchuan_report_ad"`
	QianchuanAd       []*QianchuanAd       `json:"qianchuan_ad" bson:"qianchuan_ad"`
}

type qianchuanCampaignRepo struct {
	data *Data
	log  *log.Helper
}

func (qc *QianchuanCampaign) ToDomain() *domain.QianchuanCampaign {
	return &domain.QianchuanCampaign{
		Id:             qc.Id,
		AdvertiserId:   qc.AdvertiserId,
		Name:           qc.Name,
		Budget:         qc.Budget,
		BudgetMode:     qc.BudgetMode,
		MarketingGoal:  qc.MarketingGoal,
		MarketingScene: qc.MarketingScene,
		Status:         qc.Status,
		CreateDate:     qc.CreateDate,
		CreateTime:     qc.CreateTime,
		UpdateTime:     qc.UpdateTime,
	}
}

func NewQianchuanCampaignRepo(data *Data, logger log.Logger) biz.QianchuanCampaignRepo {
	return &qianchuanCampaignRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (qcr *qianchuanCampaignRepo) GetByAdvertiserId(ctx context.Context, advertiserId, campaignId uint64, day string) (*domain.QianchuanCampaign, error) {
	var qianchuanCampaign QianchuanCampaign

	collection := qcr.data.mdb.Database(qcr.data.conf.Mongo.Dbname).Collection("qianchuan_campaign_" + day)

	if err := collection.FindOne(ctx, bson.M{
		"$and": []bson.M{
			bson.M{"status": bson.D{{"$ne", "DELETE"}}},
			bson.M{"advertiser_id": advertiserId},
			bson.M{"id": campaignId},
		},
	}).Decode(&qianchuanCampaign); err != nil {
		return nil, err
	}

	return qianchuanCampaign.ToDomain(), nil
}

func (qcr *qianchuanCampaignRepo) ListByAdvertiserId(ctx context.Context, advertiserId uint64, day, keyword string) ([]*domain.QianchuanCampaign, error) {
	list := make([]*domain.QianchuanCampaign, 0)
	var qianchuanCampaigns []*QianchuanCampaign

	collection := qcr.data.mdb.Database(qcr.data.conf.Mongo.Dbname).Collection("qianchuan_campaign_" + day)

	var and []bson.M

	if l := utf8.RuneCountInString(keyword); l > 0 {
		keywords := strings.Split(keyword, " ")

		var or []bson.M

		for _, lkeyword := range keywords {
			if campaignId, err := strconv.ParseUint(lkeyword, 10, 64); err == nil {
				or = append(or, bson.M{"id": campaignId})
				or = append(or, bson.M{"name": bsonx.Regex(lkeyword, "")})
			} else {
				or = append(or, bson.M{"name": bsonx.Regex(lkeyword, "")})
			}
		}

		and = append(and, bson.M{"$or": or})
	}

	and = append(and, bson.M{"status": bson.D{{"$ne", "DELETE"}}})
	and = append(and, bson.M{"advertiser_id": advertiserId})

	cursor, err := collection.Find(ctx, bson.M{
		"$and": and,
	})
	
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &qianchuanCampaigns)

	if err != nil {
		return nil, err
	}

	for _, qianchuanCampaign := range qianchuanCampaigns {
		list = append(list, &domain.QianchuanCampaign{
			Id:             qianchuanCampaign.Id,
			AdvertiserId:   qianchuanCampaign.AdvertiserId,
			Name:           qianchuanCampaign.Name,
			Budget:         qianchuanCampaign.Budget,
			BudgetMode:     qianchuanCampaign.BudgetMode,
			MarketingGoal:  qianchuanCampaign.MarketingGoal,
			MarketingScene: qianchuanCampaign.MarketingScene,
			Status:         qianchuanCampaign.Status,
			CreateDate:     qianchuanCampaign.CreateDate,
			CreateTime:     qianchuanCampaign.CreateTime,
			UpdateTime:     qianchuanCampaign.UpdateTime,
		})
	}

	return list, nil
}

func (qcr *qianchuanCampaignRepo) List(ctx context.Context, advertiserIds, day, keyword, campaignstatus, marketingGoal string, pageNum, pageSize int64) ([]*domain.QianchuanCampaign, error) {
	list := make([]*domain.QianchuanCampaign, 0)
	var qianchuanCampaigns []*QianchuanCampaign
	var aadvertiserIds bson.A

	sadvertiserIds := strings.Split(advertiserIds, ",")

	for _, sadvertiserId := range sadvertiserIds {
		if uisadvertiserId, err := strconv.ParseUint(sadvertiserId, 10, 64); err == nil {
			aadvertiserIds = append(aadvertiserIds, uisadvertiserId)
		}
	}

	collection := qcr.data.mdb.Database(qcr.data.conf.Mongo.Dbname).Collection("qianchuan_campaign_" + day)

	var cursor *mongo.Cursor
	var fOpts *options.FindOptions
	var and []bson.M
	var err error

	if pageNum > 0 {
		skip := (pageNum - 1) * pageSize

		fOpts = &options.FindOptions{
			Skip:  &skip,
			Limit: &pageSize,
		}
	}

	and = append(and, bson.M{"advertiser_id": bson.D{{"$in", aadvertiserIds}}})

	if campaignstatus == "" {
		and = append(and, bson.M{"status": bson.D{{"$ne", "DELETE"}}})
	} else {
		if campaignstatus != "all" {
			and = append(and, bson.M{"status": bson.D{{"$eq", strings.ToUpper(campaignstatus)}}})
		}
	}

	if marketingGoal != "" {
		and = append(and, bson.M{"marketing_goal": bson.D{{"$eq", marketingGoal}}})
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		keywords := strings.Split(keyword, " ")

		var or []bson.M

		for _, lkeyword := range keywords {
			if campaignId, err := strconv.ParseUint(lkeyword, 10, 64); err == nil {
				or = append(or, bson.M{"id": campaignId})
				or = append(or, bson.M{"name": bsonx.Regex(lkeyword, "")})
			} else {
				or = append(or, bson.M{"name": bsonx.Regex(lkeyword, "")})
			}
		}

		and = append(and, bson.M{"$or": or})
	}

	cursor, err = collection.Find(ctx, bson.M{
		"$and": and,
	}, fOpts)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &qianchuanCampaigns)

	if err != nil {
		return nil, err
	}

	for _, qianchuanCampaign := range qianchuanCampaigns {
		list = append(list, &domain.QianchuanCampaign{
			Id:             qianchuanCampaign.Id,
			AdvertiserId:   qianchuanCampaign.AdvertiserId,
			Name:           qianchuanCampaign.Name,
			Budget:         qianchuanCampaign.Budget,
			BudgetMode:     qianchuanCampaign.BudgetMode,
			MarketingGoal:  qianchuanCampaign.MarketingGoal,
			MarketingScene: qianchuanCampaign.MarketingScene,
			Status:         qianchuanCampaign.Status,
			CreateDate:     qianchuanCampaign.CreateDate,
			CreateTime:     qianchuanCampaign.CreateTime,
			UpdateTime:     qianchuanCampaign.UpdateTime,
		})
	}

	return list, nil
}

func (qcr *qianchuanCampaignRepo) All(ctx context.Context, day string) ([]*domain.QianchuanCampaign, error) {
	list := make([]*domain.QianchuanCampaign, 0)
	var qianchuanCampaigns []*QianchuanCampaign

	collection := qcr.data.mdb.Database(qcr.data.conf.Mongo.Dbname).Collection("qianchuan_campaign_" + day)

	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &qianchuanCampaigns)

	if err != nil {
		return nil, err
	}

	for _, qianchuanCampaign := range qianchuanCampaigns {
		list = append(list, &domain.QianchuanCampaign{
			Id:             qianchuanCampaign.Id,
			AdvertiserId:   qianchuanCampaign.AdvertiserId,
			Name:           qianchuanCampaign.Name,
			Budget:         qianchuanCampaign.Budget,
			BudgetMode:     qianchuanCampaign.BudgetMode,
			MarketingGoal:  qianchuanCampaign.MarketingGoal,
			MarketingScene: qianchuanCampaign.MarketingScene,
			Status:         qianchuanCampaign.Status,
			CreateDate:     qianchuanCampaign.CreateDate,
			CreateTime:     qianchuanCampaign.CreateTime,
			UpdateTime:     qianchuanCampaign.UpdateTime,
		})
	}

	return list, nil
}

func (qcr *qianchuanCampaignRepo) Count(ctx context.Context, advertiserIds, day, keyword, campaignstatus, marketingGoal string) (int64, error) {
	collection := qcr.data.mdb.Database(qcr.data.conf.Mongo.Dbname).Collection("qianchuan_campaign_" + day)

	var aadvertiserIds bson.A
	var and []bson.M

	sadvertiserIds := strings.Split(advertiserIds, ",")

	for _, sadvertiserId := range sadvertiserIds {
		if uisadvertiserId, err := strconv.ParseUint(sadvertiserId, 10, 64); err == nil {
			aadvertiserIds = append(aadvertiserIds, uisadvertiserId)
		}
	}

	and = append(and, bson.M{"advertiser_id": bson.D{{"$in", aadvertiserIds}}})

	if campaignstatus == "" {
		and = append(and, bson.M{"status": bson.D{{"$ne", "DELETE"}}})
	} else {
		if campaignstatus != "all" {
			and = append(and, bson.M{"status": bson.D{{"$eq", strings.ToUpper(campaignstatus)}}})
		}
	}

	if marketingGoal != "" {
		and = append(and, bson.M{"marketing_goal": bson.D{{"$eq", marketingGoal}}})
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		keywords := strings.Split(keyword, " ")

		var or []bson.M

		for _, lkeyword := range keywords {
			if campaignId, err := strconv.ParseUint(lkeyword, 10, 64); err == nil {
				or = append(or, bson.M{"id": campaignId})
				or = append(or, bson.M{"name": bsonx.Regex(lkeyword, "")})
			} else {
				or = append(or, bson.M{"name": bsonx.Regex(lkeyword, "")})
			}
		}

		and = append(and, bson.M{"$or": or})
	}

	matchStage := bson.D{
		{"$match", bson.M{
			"$and": and,
		}},
	}

	countStage := bson.D{
		{"$count", "total"},
	}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, countStage})

	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)

	var totalQianchuanCampaigns []*TotalQianchuanCampaign

	err = cursor.All(ctx, &totalQianchuanCampaigns)

	if err != nil {
		return 0, err
	}

	var total int64

	for _, totalQianchuanCampaign := range totalQianchuanCampaigns {
		total = totalQianchuanCampaign.Total
	}

	return total, nil
}

func (qcr *qianchuanCampaignRepo) CountByAdvertiserIds(ctx context.Context, advertiserIds, day string) (map[uint64]int64, error) {
	list := make(map[uint64]int64)

	collection := qcr.data.mdb.Database(qcr.data.conf.Mongo.Dbname).Collection("qianchuan_campaign_" + day)

	var aadvertiserIds bson.A

	sadvertiserIds := strings.Split(advertiserIds, ",")

	for _, sadvertiserId := range sadvertiserIds {
		if uisadvertiserId, err := strconv.ParseUint(sadvertiserId, 10, 64); err == nil {
			aadvertiserIds = append(aadvertiserIds, uisadvertiserId)
		}
	}

	matchStage := bson.D{
		{"$match", bson.M{"status": bson.D{{"$ne", "DELETE"}}, "advertiser_id": bson.D{{"$in", aadvertiserIds}}}},
	}

	groupStage := bson.D{
		{"$group", bson.M{
			"_id":   bson.D{{"advertiser_id", "$advertiser_id"}},
			"total": bson.D{{"$sum", 1}},
		}},
	}

	projectStage := bson.D{
		{"$project", bson.M{
			"_id":           0,
			"advertiser_id": "$_id.advertiser_id",
			"total":         1,
		}},
	}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})

	if err != nil {
		return list, err
	}

	defer cursor.Close(ctx)

	var countQianchuanCampaigns []*CountQianchuanCampaign

	err = cursor.All(ctx, &countQianchuanCampaigns)

	if err != nil {
		return list, err
	}

	for _, countQianchuanCampaign := range countQianchuanCampaigns {
		list[countQianchuanCampaign.AdvertiserId] = countQianchuanCampaign.Total
	}

	return list, nil
}

func (qcr *qianchuanCampaignRepo) SaveIndex(ctx context.Context, day string) {
	collection := qcr.data.mdb.Database(qcr.data.conf.Mongo.Dbname).Collection("qianchuan_campaign_" + day)

	isNotExistAIndex := true
	isNotExistBIndex := true

	if indexSpecifications, err := collection.Indexes().ListSpecifications(ctx); err == nil {
		for _, indexSpecification := range indexSpecifications {
			if indexSpecification.Name == "id_-1_advertiser_id_-1" {
				isNotExistAIndex = false
			}

			if indexSpecification.Name == "advertiser_id_-1_status_-1_marketing_goal_-1" {
				isNotExistBIndex = false
			}
		}

		if isNotExistAIndex {
			collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{Key: "id", Value: -1},
					{Key: "advertiser_id", Value: -1},
				},
			})
		}

		if isNotExistBIndex {
			collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{Key: "advertiser_id", Value: -1},
					{Key: "status", Value: -1},
					{Key: "marketing_goal", Value: -1},
				},
			})
		}
	}
}

func (qcr *qianchuanCampaignRepo) Upsert(ctx context.Context, day string, in *domain.QianchuanCampaign) error {
	collection := qcr.data.mdb.Database(qcr.data.conf.Mongo.Dbname).Collection("qianchuan_campaign_" + day)

	if _, err := collection.UpdateOne(ctx, bson.D{
		{"id", in.Id},
		{"advertiser_id", in.AdvertiserId},
	}, bson.D{
		{"$set", bson.D{
			{"name", in.Name},
			{"budget", in.Budget},
			{"budget_mode", in.BudgetMode},
			{"marketing_goal", in.MarketingGoal},
			{"marketing_scene", in.MarketingScene},
			{"status", in.Status},
			{"create_date", in.CreateDate},
			{"create_time", in.CreateTime},
			{"update_time", in.UpdateTime},
		}},
	}, options.Update().SetUpsert(true)); err != nil {
		return err
	}

	return nil
}

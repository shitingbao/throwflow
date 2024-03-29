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
	"go.mongodb.org/mongo-driver/x/bsonx"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// 千川广告计划表
type QianchuanAd struct {
	AdId            uint64                 `json:"ad_id" bson:"ad_id"`
	AdvertiserId    uint64                 `json:"advertiser_id" bson:"advertiser_id"`
	CampaignId      uint64                 `json:"campaign_id" bson:"campaign_id"`
	PromotionWay    string                 `json:"promotion_way" bson:"promotion_way"`
	MarketingGoal   string                 `json:"marketing_goal" bson:"marketing_goal"`
	MarketingScene  string                 `json:"marketing_scene" bson:"marketing_scene"`
	Name            string                 `json:"name" bson:"name"`
	Status          string                 `json:"status" bson:"status"`
	OptStatus       string                 `json:"opt_status" bson:"opt_status"`
	AdCreateTime    string                 `json:"ad_create_time" bson:"ad_create_time"`
	AdModifyTime    string                 `json:"ad_modify_time" bson:"ad_modify_time"`
	LabAdType       string                 `json:"lab_ad_type" bson:"lab_ad_type"`
	ProductInfo     []*domain.ProductInfo  `json:"product_info" bson:"product_info"`
	AwemeInfo       []*domain.AwemeInfo    `json:"aweme_info" bson:"aweme_info"`
	DeliverySetting domain.DeliverySetting `json:"delivery_setting" bson:"delivery_setting"`
	CreateTime      time.Time              `json:"create_time" bson:"create_time"`
	UpdateTime      time.Time              `json:"update_time" bson:"update_time"`
}

type QianchuanNotLadAd struct {
	Id           uint64 `json:"id" bson:"id"`
	AdvertiserId uint64 `json:"advertiser_id" bson:"advertiser_id"`
	CampaignId   uint64 `json:"campaign_id" bson:"campaign_id"`
}

type ProductQianchuanAd struct {
	AdvertiserId      uint64             `json:"advertiser_id" bson:"advertiser_id"`
	AdId              uint64             `json:"ad_id" bson:"ad_id"`
	ProductId         uint64             `json:"product_id" bson:"product_id"`
	DiscountPrice     float64            `json:"discount_price" bson:"discount_price"`
	ProductName       string             `json:"product_name" bson:"product_name"`
	ProductImg        string             `json:"product_img" bson:"product_img"`
	QianchuanReportAd *QianchuanReportAd `json:"qianchuan_report_ad" bson:"qianchuan_report_ad"`
}

type AwemeQianchuanAd struct {
	AdvertiserId      uint64             `json:"advertiser_id" bson:"advertiser_id"`
	AdId              uint64             `json:"ad_id" bson:"ad_id"`
	AwemeId           uint64             `json:"aweme_id" bson:"aweme_id"`
	AwemeName         string             `json:"aweme_name" bson:"aweme_name"`
	AwemeShowId       string             `json:"aweme_show_id" bson:"aweme_show_id"`
	AwemeAvatar       string             `json:"aweme_avatar" bson:"aweme_avatar"`
	QianchuanReportAd *QianchuanReportAd `json:"qianchuan_report_ad" bson:"qianchuan_report_ad"`
}

type QianchuanProductAd struct {
	AdId            uint64                 `json:"ad_id" bson:"ad_id"`
	AdvertiserId    uint64                 `json:"advertiser_id" bson:"advertiser_id"`
	CampaignId      uint64                 `json:"campaign_id" bson:"campaign_id"`
	PromotionWay    string                 `json:"promotion_way" bson:"promotion_way"`
	MarketingGoal   string                 `json:"marketing_goal" bson:"marketing_goal"`
	MarketingScene  string                 `json:"marketing_scene" bson:"marketing_scene"`
	Name            string                 `json:"name" bson:"name"`
	Status          string                 `json:"status" bson:"status"`
	OptStatus       string                 `json:"opt_status" bson:"opt_status"`
	AdCreateTime    string                 `json:"ad_create_time" bson:"ad_create_time"`
	AdModifyTime    string                 `json:"ad_modify_time" bson:"ad_modify_time"`
	LabAdType       string                 `json:"lab_ad_type" bson:"lab_ad_type"`
	ProductInfo     *domain.ProductInfo    `json:"product_info" bson:"product_info"`
	AwemeInfo       []*domain.AwemeInfo    `json:"aweme_info" bson:"aweme_info"`
	DeliverySetting domain.DeliverySetting `json:"delivery_setting" bson:"delivery_setting"`
	CreateTime      time.Time              `json:"create_time" bson:"create_time"`
	UpdateTime      time.Time              `json:"update_time" bson:"update_time"`
}

type QianchuanAwemeAd struct {
	AdId            uint64                 `json:"ad_id" bson:"ad_id"`
	AdvertiserId    uint64                 `json:"advertiser_id" bson:"advertiser_id"`
	CampaignId      uint64                 `json:"campaign_id" bson:"campaign_id"`
	PromotionWay    string                 `json:"promotion_way" bson:"promotion_way"`
	MarketingGoal   string                 `json:"marketing_goal" bson:"marketing_goal"`
	MarketingScene  string                 `json:"marketing_scene" bson:"marketing_scene"`
	Name            string                 `json:"name" bson:"name"`
	Status          string                 `json:"status" bson:"status"`
	OptStatus       string                 `json:"opt_status" bson:"opt_status"`
	AdCreateTime    string                 `json:"ad_create_time" bson:"ad_create_time"`
	AdModifyTime    string                 `json:"ad_modify_time" bson:"ad_modify_time"`
	LabAdType       string                 `json:"lab_ad_type" bson:"lab_ad_type"`
	ProductInfo     []*domain.ProductInfo  `json:"product_info" bson:"product_info"`
	AwemeInfo       *domain.AwemeInfo      `json:"aweme_info" bson:"aweme_info"`
	DeliverySetting domain.DeliverySetting `json:"delivery_setting" bson:"delivery_setting"`
	CreateTime      time.Time              `json:"create_time" bson:"create_time"`
	UpdateTime      time.Time              `json:"update_time" bson:"update_time"`
}

type ListQianchuanAdByPromotionId struct {
	AdId              uint64               `json:"ad_id" bson:"ad_id"`
	QianchuanReportAd []*QianchuanReportAd `json:"qianchuan_report_ad" bson:"qianchuan_report_ad"`
}

type QianchuanAdQianchuanCampaign struct {
	Id             uint64  `json:"id" bson:"id"`
	AdvertiserId   uint64  `json:"advertiser_id" bson:"advertiser_id"`
	Name           string  `json:"name" bson:"name"`
	Budget         float64 `json:"budget" bson:"budget"`
	BudgetMode     string  `json:"budget_mode" bson:"budget_mode"`
	MarketingGoal  string  `json:"marketing_goal" bson:"marketing_goal"`
	MarketingScene string  `json:"marketing_scene" bson:"marketing_scene"`
	Status         string  `json:"status" bson:"status"`
	CreateDate     string  `json:"create_date" bson:"create_date"`
	Ads            int64   `json:"ads" bson:"ads"`
}

type qianchuanAdRepo struct {
	data *Data
	log  *log.Helper
}

func (qa *QianchuanAd) ToDomain() *domain.QianchuanAd {
	return &domain.QianchuanAd{
		AdId:            qa.AdId,
		AdvertiserId:    qa.AdvertiserId,
		CampaignId:      qa.CampaignId,
		PromotionWay:    qa.PromotionWay,
		MarketingGoal:   qa.MarketingGoal,
		MarketingScene:  qa.MarketingScene,
		Name:            qa.Name,
		Status:          qa.Status,
		OptStatus:       qa.OptStatus,
		AdCreateTime:    qa.AdCreateTime,
		AdModifyTime:    qa.AdModifyTime,
		LabAdType:       qa.LabAdType,
		ProductInfo:     qa.ProductInfo,
		AwemeInfo:       qa.AwemeInfo,
		DeliverySetting: qa.DeliverySetting,
		CreateTime:      qa.CreateTime,
		UpdateTime:      qa.UpdateTime,
	}
}

func NewQianchuanAdRepo(data *Data, logger log.Logger) biz.QianchuanAdRepo {
	return &qianchuanAdRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (qar *qianchuanAdRepo) GetByAdIdAndDay(ctx context.Context, adId uint64, day string) (*domain.QianchuanAd, error) {
	var qianchuanAd QianchuanAd

	collection := qar.data.mdb.Database(qar.data.conf.Mongo.Dbname).Collection("qianchuan_ad_" + day)

	if err := collection.FindOne(ctx, bson.D{{"ad_id", adId}}).Decode(&qianchuanAd); err != nil {
		return nil, err
	}

	return qianchuanAd.ToDomain(), nil
}

func (qar *qianchuanAdRepo) ListNotLadAd(ctx context.Context, marketingGoal, day string) (map[uint64]map[uint64][]*domain.QianchuanCampaign, error) {
	list := make(map[uint64]map[uint64][]*domain.QianchuanCampaign)

	collection := qar.data.mdb.Database(qar.data.conf.Mongo.Dbname).Collection("qianchuan_ad_" + day)

	matchStage := bson.D{
		{"$match", bson.M{"lab_ad_type": "NOT_LAB_AD", "marketing_goal": marketingGoal}},
	}

	var unwindStage bson.D
	var groupStage bson.D
	var projectStage bson.D

	if marketingGoal == "VIDEO_PROM_GOODS" {
		unwindStage = bson.D{
			{"$unwind", "$product_info"},
		}

		groupStage = bson.D{
			{"$group", bson.M{"_id": bson.D{{"product_id", "$product_info.id"}, {"aweme_id", "$advertiser_id"}, {"advertiser_id", "$advertiser_id"}, {"campaign_id", "$campaign_id"}}}},
		}

		projectStage = bson.D{
			{"$project", bson.M{
				"_id":           0,
				"id":            "$_id.product_id",
				"advertiser_id": "$_id.advertiser_id",
				"campaign_id":   "$_id.campaign_id",
			}},
		}
	} else {
		unwindStage = bson.D{
			{"$unwind", "$aweme_info"},
		}

		groupStage = bson.D{
			{"$group", bson.M{"_id": bson.D{{"aweme_id", "$aweme_info.aweme_id"}, {"advertiser_id", "$advertiser_id"}, {"campaign_id", "$campaign_id"}}}},
		}

		projectStage = bson.D{
			{"$project", bson.M{
				"_id":           0,
				"id":            "$_id.aweme_id",
				"advertiser_id": "$_id.advertiser_id",
				"campaign_id":   "$_id.campaign_id",
			}},
		}
	}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, unwindStage, groupStage, projectStage})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var qianchuanNotLadAds []QianchuanNotLadAd

	err = cursor.All(ctx, &qianchuanNotLadAds)

	if err != nil {
		return nil, err
	}

	for _, qianchuanNotLadAd := range qianchuanNotLadAds {
		if _, ok := list[qianchuanNotLadAd.Id]; !ok {
			list[qianchuanNotLadAd.Id] = make(map[uint64][]*domain.QianchuanCampaign)
		}

		if _, ok := list[qianchuanNotLadAd.Id][qianchuanNotLadAd.AdvertiserId]; !ok {
			list[qianchuanNotLadAd.Id][qianchuanNotLadAd.AdvertiserId] = make([]*domain.QianchuanCampaign, 0)
		}

		list[qianchuanNotLadAd.Id][qianchuanNotLadAd.AdvertiserId] = append(list[qianchuanNotLadAd.Id][qianchuanNotLadAd.AdvertiserId], &domain.QianchuanCampaign{
			Id: qianchuanNotLadAd.CampaignId,
		})
	}

	return list, nil
}

func (qar *qianchuanAdRepo) ListProductAd(ctx context.Context, day string) ([]*domain.QianchuanReportProduct, error) {
	list := make([]*domain.QianchuanReportProduct, 0)

	collection := qar.data.mdb.Database(qar.data.conf.Mongo.Dbname).Collection("qianchuan_ad_" + day)

	matchStage := bson.D{
		{"$match", bson.M{"marketing_goal": "VIDEO_PROM_GOODS"}},
	}

	unwindAStage := bson.D{
		{"$unwind", "$product_info"},
	}

	projectStage := bson.D{
		{"$project", bson.M{
			"advertiser_id":  "$advertiser_id",
			"ad_id":          "$ad_id",
			"product_id":     "$product_info.id",
			"discount_price": "$product_info.discount_price",
			"product_name":   "$product_info.name",
			"product_img":    "$product_info.img",
		}},
	}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, unwindAStage, projectStage})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var productQianchuanAds []*ProductQianchuanAd

	err = cursor.All(ctx, &productQianchuanAds)

	if err != nil {
		return nil, err
	}

	for _, productQianchuanAd := range productQianchuanAds {
		isNotExist := true

		for _, l := range list {
			if l.ProductId == productQianchuanAd.ProductId && l.AdvertiserId == productQianchuanAd.AdvertiserId {
				l.AdIds = append(l.AdIds, productQianchuanAd.AdId)

				isNotExist = false

				break
			}
		}

		if isNotExist {
			adIds := make([]uint64, 0)
			adIds = append(adIds, productQianchuanAd.AdId)

			qianchuanReportProduct := &domain.QianchuanReportProduct{
				AdvertiserId:  productQianchuanAd.AdvertiserId,
				ProductId:     productQianchuanAd.ProductId,
				DiscountPrice: productQianchuanAd.DiscountPrice,
				ProductName:   productQianchuanAd.ProductName,
				ProductImg:    productQianchuanAd.ProductImg,
				AdIds:         adIds,
			}

			list = append(list, qianchuanReportProduct)
		}
	}

	return list, nil
}

func (qar *qianchuanAdRepo) ListAwemeAd(ctx context.Context, day string) ([]*domain.QianchuanReportAweme, error) {
	list := make([]*domain.QianchuanReportAweme, 0)

	collection := qar.data.mdb.Database(qar.data.conf.Mongo.Dbname).Collection("qianchuan_ad_" + day)

	matchStage := bson.D{
		{"$match", bson.M{"marketing_goal": "LIVE_PROM_GOODS"}},
	}

	unwindAStage := bson.D{
		{"$unwind", "$aweme_info"},
	}

	projectStage := bson.D{
		{"$project", bson.M{
			"advertiser_id": "$advertiser_id",
			"ad_id":         "$ad_id",
			"aweme_id":      "$aweme_info.aweme_id",
			"aweme_name":    "$aweme_info.aweme_name",
			"aweme_show_id": "$aweme_info.aweme_show_id",
			"aweme_avatar":  "$aweme_info.aweme_avatar",
		}},
	}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, unwindAStage, projectStage})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var awemeQianchuanAds []*AwemeQianchuanAd

	err = cursor.All(ctx, &awemeQianchuanAds)

	if err != nil {
		return nil, err
	}

	for _, awemeQianchuanAd := range awemeQianchuanAds {
		isNotExist := true

		for _, l := range list {
			if l.AwemeId == awemeQianchuanAd.AwemeId && l.AdvertiserId == awemeQianchuanAd.AdvertiserId {
				l.AdIds = append(l.AdIds, awemeQianchuanAd.AdId)

				isNotExist = false

				break
			}
		}

		if isNotExist {
			adIds := make([]uint64, 0)
			adIds = append(adIds, awemeQianchuanAd.AdId)

			qianchuanReportAweme := &domain.QianchuanReportAweme{
				AdvertiserId: awemeQianchuanAd.AdvertiserId,
				AwemeId:      awemeQianchuanAd.AwemeId,
				AwemeName:    awemeQianchuanAd.AwemeName,
				AwemeShowId:  awemeQianchuanAd.AwemeShowId,
				AwemeAvatar:  awemeQianchuanAd.AwemeAvatar,
				AdIds:        adIds,
			}

			list = append(list, qianchuanReportAweme)
		}
	}

	return list, nil
}

func (qar *qianchuanAdRepo) ListAwemeByAdvertiserId(ctx context.Context, advertiserId uint64, day string) ([]*domain.QianchuanReportAweme, error) {
	list := make([]*domain.QianchuanReportAweme, 0)

	collection := qar.data.mdb.Database(qar.data.conf.Mongo.Dbname).Collection("qianchuan_ad_" + day)

	matchStage := bson.D{
		{"$match", bson.M{"marketing_goal": "LIVE_PROM_GOODS", "advertiser_id": advertiserId}},
	}

	unwindAStage := bson.D{
		{"$unwind", "$aweme_info"},
	}

	groupStage := bson.D{
		{"$group", bson.M{
			"_id":           "$aweme_info.aweme_id",
			"advertiser_id": bson.D{{"$first", "$advertiser_id"}},
		}},
	}

	projectStage := bson.D{
		{"$project", bson.M{
			"_id":           0,
			"aweme_id":      "$_id",
			"advertiser_id": 1,
		}},
	}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, unwindAStage, groupStage, projectStage})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var awemeQianchuanAds []*AwemeQianchuanAd

	err = cursor.All(ctx, &awemeQianchuanAds)

	if err != nil {
		return nil, err
	}

	for _, awemeQianchuanAd := range awemeQianchuanAds {
		list = append(list, &domain.QianchuanReportAweme{
			AdvertiserId: awemeQianchuanAd.AdvertiserId,
			AwemeId:      awemeQianchuanAd.AwemeId,
		})
	}

	return list, nil
}

func (qar *qianchuanAdRepo) ListByPromotionId(ctx context.Context, advertiserId, promotionId uint64, promotionType, day string) ([]*domain.QianchuanReportAd, error) {
	list := make([]*domain.QianchuanReportAd, 0)

	collection := qar.data.mdb.Database(qar.data.conf.Mongo.Dbname).Collection("qianchuan_ad_" + day)

	var matchStage bson.D

	if promotionType == "aweme" {
		matchStage = bson.D{
			{"$match", bson.M{
				"marketing_goal":      "LIVE_PROM_GOODS",
				"advertiser_id":       advertiserId,
				"lab_ad_type":         "LAB_AD",
				"aweme_info.aweme_id": promotionId,
			}},
		}
	} else {
		matchStage = bson.D{
			{"$match", bson.M{
				"marketing_goal":  "VIDEO_PROM_GOODS",
				"advertiser_id":   advertiserId,
				"lab_ad_type":     "LAB_AD",
				"product_info.id": int64(promotionId),
			}},
		}
	}

	projectStage := bson.D{
		{"$project", bson.M{
			"ad_id": 1,
		}},
	}

	lookupStage := bson.D{
		{"$lookup", bson.M{
			"from":         "qianchuan_report_ad_" + day,
			"localField":   "ad_id",
			"foreignField": "ad_id",
			"as":           "qianchuan_report_ad",
		}},
	}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, projectStage, lookupStage})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var listQianchuanAdByPromotionIds []*ListQianchuanAdByPromotionId

	err = cursor.All(ctx, &listQianchuanAdByPromotionIds)

	if err != nil {
		return nil, err
	}

	for _, listQianchuanAdByPromotionId := range listQianchuanAdByPromotionIds {
		qianchuanReportAd := &domain.QianchuanReportAd{
			AdId: listQianchuanAdByPromotionId.AdId,
		}

		for _, lqianchuanReportAd := range listQianchuanAdByPromotionId.QianchuanReportAd {
			qianchuanReportAd.StatCost = lqianchuanReportAd.StatCost
			qianchuanReportAd.ShowCnt = lqianchuanReportAd.ShowCnt
			qianchuanReportAd.ClickCnt = lqianchuanReportAd.ClickCnt
			qianchuanReportAd.PayOrderCount = lqianchuanReportAd.PayOrderCount
			qianchuanReportAd.CreateOrderAmount = lqianchuanReportAd.CreateOrderAmount
			qianchuanReportAd.CreateOrderCount = lqianchuanReportAd.CreateOrderCount
			qianchuanReportAd.PayOrderAmount = lqianchuanReportAd.PayOrderAmount
			qianchuanReportAd.DyFollow = lqianchuanReportAd.DyFollow
			qianchuanReportAd.ConvertCnt = lqianchuanReportAd.ConvertCnt
		}

		list = append(list, qianchuanReportAd)
	}

	return list, nil
}

func (qar *qianchuanAdRepo) ListByCampaignId(ctx context.Context, campaignId uint64, day string) ([]*domain.QianchuanAd, error) {
	list := make([]*domain.QianchuanAd, 0)

	collection := qar.data.mdb.Database(qar.data.conf.Mongo.Dbname).Collection("qianchuan_ad_" + day)

	matchStage := bson.D{
		{"$match", bson.M{"campaign_id": campaignId}},
	}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var qianchuanAds []*QianchuanAd

	err = cursor.All(ctx, &qianchuanAds)

	if err != nil {
		return nil, err
	}

	for _, qianchuanAd := range qianchuanAds {
		list = append(list, qianchuanAd.ToDomain())
	}

	return list, nil
}

func (qar *qianchuanAdRepo) List(ctx context.Context, advertiserIds, day, keyword string, pageNum, pageSize int64) ([]*domain.QianchuanAd, error) {
	list := make([]*domain.QianchuanAd, 0)
	var qianchuanAds []*QianchuanAd
	var aadvertiserIds bson.A

	sadvertiserIds := strings.Split(advertiserIds, ",")

	for _, sadvertiserId := range sadvertiserIds {
		if uisadvertiserId, err := strconv.ParseUint(sadvertiserId, 10, 64); err == nil {
			aadvertiserIds = append(aadvertiserIds, uisadvertiserId)
		}
	}

	collection := qar.data.mdb.Database(qar.data.conf.Mongo.Dbname).Collection("qianchuan_ad_" + day)

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

	if l := utf8.RuneCountInString(keyword); l > 0 {
		keywords := strings.Split(keyword, " ")

		var or []bson.M

		for _, lkeyword := range keywords {
			if campaignId, err := strconv.ParseUint(lkeyword, 10, 64); err == nil {
				or = append(or, bson.M{"ad_id": campaignId})
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

	err = cursor.All(ctx, &qianchuanAds)

	if err != nil {
		return nil, err
	}

	for _, qianchuanAd := range qianchuanAds {
		list = append(list, qianchuanAd.ToDomain())
	}

	return list, nil
}

func (qar *qianchuanAdRepo) AllNotLadAd(ctx context.Context, day string) ([]*domain.QianchuanAd, error) {
	list := make([]*domain.QianchuanAd, 0)

	collection := qar.data.mdb.Database(qar.data.conf.Mongo.Dbname).Collection("qianchuan_ad_" + day)

	matchStage := bson.D{
		{"$match", bson.M{"lab_ad_type": "NOT_LAB_AD"}},
	}

	projectStage := bson.D{
		{"$project", bson.M{
			"ad_id":       1,
			"campaign_id": 1,
		}},
	}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, projectStage})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var qianchuanAds []*QianchuanAd

	err = cursor.All(ctx, &qianchuanAds)

	if err != nil {
		return nil, err
	}

	for _, qianchuanAd := range qianchuanAds {
		list = append(list, qianchuanAd.ToDomain())
	}

	return list, nil
}

func (qar *qianchuanAdRepo) CountByAdIdAndAdvertiserIdAndDay(ctx context.Context, advertiserId uint64, day string) (int64, error) {
	collection := qar.data.mdb.Database(qar.data.conf.Mongo.Dbname).Collection("qianchuan_ad_" + day)

	return collection.CountDocuments(ctx, bson.M{"advertiser_id": advertiserId})
}

func (qar *qianchuanAdRepo) Count(ctx context.Context, advertiserIds, day, keyword, adStatus, marketingGoal string) (int64, error) {
	collection := qar.data.mdb.Database(qar.data.conf.Mongo.Dbname).Collection("qianchuan_ad_" + day)

	var and []bson.M
	var aadvertiserIds bson.A

	sadvertiserIds := strings.Split(advertiserIds, ",")

	for _, sadvertiserId := range sadvertiserIds {
		if uisadvertiserId, err := strconv.ParseUint(sadvertiserId, 10, 64); err == nil {
			aadvertiserIds = append(aadvertiserIds, uisadvertiserId)
		}
	}

	and = append(and, bson.M{"advertiser_id": bson.D{{"$in", aadvertiserIds}}})

	if l := utf8.RuneCountInString(adStatus); l > 0 {
		if adStatus == "timeDone" {
			and = append(and, bson.M{"campaign_status": bson.D{{"$eq", "TIME_DONE"}}})
		} else if adStatus == "timeNoReach" {
			and = append(and, bson.M{"campaign_status": bson.D{{"$eq", "TIME_NO_REACH"}}})
		} else if adStatus == "offlineBalance" {
			and = append(and, bson.M{"campaign_status": bson.D{{"$eq", "OFFLINE_BALANCE"}}})
		} else if adStatus == "offlineBudget" {
			and = append(and, bson.M{"campaign_status": bson.D{{"$eq", "OFFLINE_BUDGET"}}})
		} else if adStatus == "deliveryOk" {
			and = append(and, bson.M{"campaign_status": bson.D{{"$eq", "DELIVERY_OK"}}})
		} else if adStatus == "noSchedule" {
			and = append(and, bson.M{"campaign_status": bson.D{{"$eq", "NO_SCHEDULE"}}})
		} else if adStatus == "offlineAudit" {
			and = append(and, bson.M{"campaign_status": bson.D{{"$eq", "OFFLINE_AUDIT"}}})
		} else if adStatus == "externalUrlDisable" {
			and = append(and, bson.M{"campaign_status": bson.D{{"$eq", "EXTERNAL_URL_DISABLE"}}})
		} else if adStatus == "liveRoomOff" {
			and = append(and, bson.M{"campaign_status": bson.D{{"$eq", "LIVE_ROOM_OFF"}}})
		} else if adStatus == "systemDisable" {
			and = append(and, bson.M{"campaign_status": bson.D{{"$eq", "SYSTEM_DISABLE"}}})
		} else if adStatus == "allIncludeDeleted" {
			and = append(and, bson.M{"campaign_status": bson.D{{"$eq", "ALL_INCLUDE_DELETED"}}})
		} else if adStatus == "quotaDisable" {
			and = append(and, bson.M{"campaign_status": bson.D{{"$eq", "QUOTA_DISABLE"}}})
		} else {
			and = append(and, bson.M{"campaign_status": bson.D{{"$eq", strings.ToUpper(adStatus)}}})
		}
	}

	if marketingGoal == "videoPromGoods" {
		and = append(and, bson.M{"marketing_goal": bson.D{{"$eq", "VIDEO_PROM_GOODS"}}})
	} else if marketingGoal == "livePromGoods" {
		and = append(and, bson.M{"marketing_goal": bson.D{{"$eq", "LIVE_PROM_GOODS"}}})
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		keywords := strings.Split(keyword, " ")

		var or []bson.M

		for _, lkeyword := range keywords {
			if campaignId, err := strconv.ParseUint(lkeyword, 10, 64); err == nil {
				or = append(or, bson.M{"ad_id": campaignId})
				or = append(or, bson.M{"name": bsonx.Regex(lkeyword, "")})
			} else {
				or = append(or, bson.M{"name": bsonx.Regex(lkeyword, "")})
			}
		}

		and = append(and, bson.M{"$or": or})
	}

	return collection.CountDocuments(ctx, bson.M{"$and": and})
}

func (qar *qianchuanAdRepo) SaveIndex(ctx context.Context, day string) {
	collection := qar.data.mdb.Database(qar.data.conf.Mongo.Dbname).Collection("qianchuan_ad_" + day)

	isNotExistAIndex := true
	isNotExistBIndex := true
	isNotExistCIndex := true
	isNotExistDIndex := true

	if indexSpecifications, err := collection.Indexes().ListSpecifications(ctx); err == nil {
		for _, indexSpecification := range indexSpecifications {
			if indexSpecification.Name == "ad_id_-1_advertiser_id_-1" {
				isNotExistAIndex = false
			}

			if indexSpecification.Name == "advertiser_id_-1" {
				isNotExistBIndex = false
			}

			if indexSpecification.Name == "marketing_goal_-1" {
				isNotExistCIndex = false
			}

			if indexSpecification.Name == "ad_id_-1" {
				isNotExistDIndex = false
			}
		}

		if isNotExistAIndex {
			collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{Key: "ad_id", Value: -1},
					{Key: "advertiser_id", Value: -1},
				},
			})
		}

		if isNotExistBIndex {
			collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{Key: "advertiser_id", Value: -1},
				},
			})
		}

		if isNotExistCIndex {
			collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{Key: "marketing_goal", Value: -1},
				},
			})
		}

		if isNotExistDIndex {
			collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{Key: "ad_id", Value: -1},
				},
			})
		}
	}
}

func (qar *qianchuanAdRepo) Upsert(ctx context.Context, day string, in *domain.QianchuanAd) error {
	collection := qar.data.mdb.Database(qar.data.conf.Mongo.Dbname).Collection("qianchuan_ad_" + day)

	if _, err := collection.UpdateOne(ctx, bson.D{
		{"ad_id", in.AdId},
		{"advertiser_id", in.AdvertiserId},
	}, bson.D{
		{"$set", bson.D{
			{"campaign_id", in.CampaignId},
			{"promotion_way", in.PromotionWay},
			{"marketing_goal", in.MarketingGoal},
			{"marketing_scene", in.MarketingScene},
			{"name", in.Name},
			{"status", in.Status},
			{"opt_status", in.OptStatus},
			{"ad_create_time", in.AdCreateTime},
			{"ad_modify_time", in.AdModifyTime},
			{"lab_ad_type", in.LabAdType},
			{"product_info", in.ProductInfo},
			{"aweme_info", in.AwemeInfo},
			{"delivery_setting", in.DeliverySetting},
			{"create_time", in.CreateTime},
			{"update_time", in.UpdateTime},
		}},
	}, options.Update().SetUpsert(true)); err != nil {
		return err
	}

	return nil
}

func (qar *qianchuanAdRepo) Send(ctx context.Context, message event.Event) error {
	if err := qar.data.kafka.Send(ctx, message); err != nil {
		return err
	}

	return nil
}

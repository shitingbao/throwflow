package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// 千川广告计划完整信息表
type QianchuanAdInfo struct {
	AdId                    uint64                 `json:"ad_id" bson:"ad_id"`
	AdvertiserId            uint64                 `json:"advertiser_id" bson:"advertiser_id"`
	AdvertiserName          string                 `json:"advertiser_name" bson:"advertiser_name"`
	CampaignId              uint64                 `json:"campaign_id" bson:"campaign_id"`
	CampaignName            string                 `json:"campaign_name" bson:"campaign_name"`
	CampaignBudget          float64                `json:"campaign_budget" bson:"campaign_budget"`
	CampaignBudgetMode      string                 `json:"campaign_budget_mode" bson:"campaign_budget_mode"`
	CampaignStatus          string                 `json:"campaign_status" bson:"campaign_status"`
	CampaignCreateDate      string                 `json:"campaign_create_date" bson:"campaign_create_date"`
	PromotionWay            string                 `json:"promotion_way" bson:"promotion_way"`
	MarketingGoal           string                 `json:"marketing_goal" bson:"marketing_goal"`
	MarketingScene          string                 `json:"marketing_scene" bson:"marketing_scene"`
	Name                    string                 `json:"name" bson:"name"`
	Status                  string                 `json:"status" bson:"status"`
	OptStatus               string                 `json:"opt_status" bson:"opt_status"`
	AdCreateTime            string                 `json:"ad_create_time" bson:"ad_create_time"`
	AdModifyTime            string                 `json:"ad_modify_time" bson:"ad_modify_time"`
	LabAdType               string                 `json:"lab_ad_type" bson:"lab_ad_type"`
	StatCost                float64                `json:"stat_cost" bson:"stat_cost"`
	Roi                     float64                `json:"roi" bson:"roi"`
	ShowCnt                 int64                  `json:"show_cnt" bson:"show_cnt"`
	ClickCnt                int64                  `json:"click_cnt" bson:"click_cnt"`
	PayOrderCount           int64                  `json:"pay_order_count" bson:"pay_order_count"`
	CreateOrderAmount       float64                `json:"create_order_amount" bson:"create_order_amount"`
	CreateOrderCount        int64                  `json:"create_order_count" bson:"create_order_count"`
	PayOrderAmount          float64                `json:"pay_order_amount" bson:"pay_order_amount"`
	DyFollow                int64                  `json:"dy_follow" bson:"dy_follow"`
	ConvertCnt              int64                  `json:"convert_cnt" bson:"convert_cnt"`
	ClickRate               float64                `json:"click_rate" bson:"click_rate"`
	CpmPlatform             float64                `json:"cpm_platform" bson:"cpm_platform"`
	PayConvertRate          float64                `json:"pay_convert_rate" bson:"pay_convert_rate"`
	ConvertCost             float64                `json:"convert_cost" bson:"convert_cost"`
	ConvertRate             float64                `json:"convert_rate" bson:"convert_rate"`
	AveragePayOrderStatCost float64                `json:"average_pay_order_stat_cost" bson:"average_pay_order_stat_cost"`
	PayOrderAveragePrice    float64                `json:"pay_order_average_price" bson:"pay_order_average_price"`
	ProductInfo             []*domain.ProductInfo  `json:"product_info" bson:"product_info"`
	AwemeInfo               []*domain.AwemeInfo    `json:"aweme_info" bson:"aweme_info"`
	DeliverySetting         domain.DeliverySetting `json:"delivery_setting" bson:"delivery_setting"`
	CreateTime              time.Time              `json:"create_time" bson:"create_time"`
	UpdateTime              time.Time              `json:"update_time" bson:"update_time"`
}

type TotalQianchuanAdInfo struct {
	Total int64 `json:"total" bson:"total"`
}

type qianchuanAdInfoRepo struct {
	data *Data
	log  *log.Helper
}

func (qai *QianchuanAdInfo) ToDomain() *domain.QianchuanAdInfo {
	return &domain.QianchuanAdInfo{
		AdId:                    qai.AdId,
		AdvertiserId:            qai.AdvertiserId,
		AdvertiserName:          qai.AdvertiserName,
		CampaignId:              qai.CampaignId,
		CampaignName:            qai.CampaignName,
		CampaignBudget:          qai.CampaignBudget,
		CampaignBudgetMode:      qai.CampaignBudgetMode,
		CampaignStatus:          qai.CampaignStatus,
		CampaignCreateDate:      qai.CampaignCreateDate,
		PromotionWay:            qai.PromotionWay,
		MarketingGoal:           qai.MarketingGoal,
		MarketingScene:          qai.MarketingScene,
		Name:                    qai.Name,
		Status:                  qai.Status,
		OptStatus:               qai.OptStatus,
		AdCreateTime:            qai.AdCreateTime,
		AdModifyTime:            qai.AdModifyTime,
		LabAdType:               qai.LabAdType,
		StatCost:                qai.StatCost,
		Roi:                     qai.Roi,
		ShowCnt:                 qai.ShowCnt,
		ClickCnt:                qai.ClickCnt,
		PayOrderCount:           qai.PayOrderCount,
		CreateOrderAmount:       qai.CreateOrderAmount,
		CreateOrderCount:        qai.CreateOrderCount,
		PayOrderAmount:          qai.PayOrderAmount,
		DyFollow:                qai.DyFollow,
		ConvertCnt:              qai.ConvertCnt,
		ClickRate:               qai.ClickRate,
		CpmPlatform:             qai.CpmPlatform,
		PayConvertRate:          qai.PayConvertRate,
		ConvertCost:             qai.ConvertCost,
		ConvertRate:             qai.ConvertRate,
		AveragePayOrderStatCost: qai.AveragePayOrderStatCost,
		PayOrderAveragePrice:    qai.PayOrderAveragePrice,
		ProductInfo:             qai.ProductInfo,
		AwemeInfo:               qai.AwemeInfo,
		DeliverySetting:         qai.DeliverySetting,
		CreateTime:              qai.CreateTime,
		UpdateTime:              qai.UpdateTime,
	}
}

func NewQianchuanAdInfoRepo(data *Data, logger log.Logger) biz.QianchuanAdInfoRepo {
	return &qianchuanAdInfoRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (qair *qianchuanAdInfoRepo) Get(ctx context.Context, adId uint64, day string) (*domain.QianchuanAdInfo, error) {
	where := make([]string, 0)

	where = append(where, "ad_id="+strconv.FormatUint(adId, 10))
	where = append(where, "day='"+day+"'")

	sql := "SELECT * FROM qianchuan_ad_info WHERE " + strings.Join(where, " AND ")

	row := qair.data.cdb.QueryRow(ctx, sql)

	var (
		qadId                   uint64
		qadvertiserId           uint64
		advertiserName          string
		campaignId              uint64
		campaignName            string
		campaignBudget          float64
		campaignBudgetMode      string
		campaignStatus          string
		campaignCreateDate      string
		promotionWay            string
		marketingGoal           string
		marketingScene          string
		name                    string
		status                  string
		optStatus               string
		adCreateTime            string
		adModifyTime            string
		labAdType               string
		statCost                float64
		roi                     float64
		showCnt                 int64
		clickCnt                int64
		payOrderCount           int64
		createOrderAmount       float64
		createOrderCount        int64
		payOrderAmount          float64
		dyFollow                int64
		convertCnt              int64
		clickRate               float64
		cpmPlatform             float64
		payConvertRate          float64
		convertCost             float64
		convertRate             float64
		averagePayOrderStatCost float64
		payOrderAveragePrice    float64
		productInfo             string
		awemeInfo               string
		deliverySetting         string
		qday                    time.Time
		createTime              time.Time
		updateTime              time.Time
	)

	if err := row.Scan(&qadId, &qadvertiserId, &advertiserName, &campaignId, &campaignName, &campaignBudget, &campaignBudgetMode, &campaignStatus, &campaignCreateDate, &promotionWay, &marketingGoal, &marketingScene, &name, &status, &optStatus, &adCreateTime, &adModifyTime, &labAdType, &statCost, &roi, &showCnt, &clickCnt, &payOrderCount, &createOrderAmount, &createOrderCount, &payOrderAmount, &dyFollow, &convertCnt, &clickRate, &cpmPlatform, &payConvertRate, &convertCost, &convertRate, &averagePayOrderStatCost, &payOrderAveragePrice, &productInfo, &awemeInfo, &deliverySetting, &qday, &createTime, &updateTime); err != nil {
		return nil, err
	}

	var sproductInfo []*domain.ProductInfo
	var sawemeInfo []*domain.AwemeInfo
	var sdeliverySetting domain.DeliverySetting

	json.Unmarshal([]byte(productInfo), &sproductInfo)
	json.Unmarshal([]byte(awemeInfo), &sawemeInfo)
	json.Unmarshal([]byte(deliverySetting), &sdeliverySetting)

	return &domain.QianchuanAdInfo{
		AdId:                    adId,
		AdvertiserId:            qadvertiserId,
		AdvertiserName:          advertiserName,
		CampaignId:              campaignId,
		CampaignName:            campaignName,
		CampaignBudget:          campaignBudget,
		CampaignBudgetMode:      campaignBudgetMode,
		CampaignStatus:          campaignStatus,
		CampaignCreateDate:      campaignCreateDate,
		PromotionWay:            promotionWay,
		MarketingGoal:           marketingGoal,
		MarketingScene:          marketingScene,
		Name:                    name,
		Status:                  status,
		OptStatus:               optStatus,
		AdCreateTime:            adCreateTime,
		AdModifyTime:            adModifyTime,
		LabAdType:               labAdType,
		StatCost:                statCost,
		Roi:                     roi,
		PayOrderCount:           payOrderCount,
		PayOrderAmount:          payOrderAmount,
		CreateOrderAmount:       createOrderAmount,
		CreateOrderCount:        createOrderCount,
		ClickCnt:                clickCnt,
		ShowCnt:                 showCnt,
		ConvertCnt:              convertCnt,
		ClickRate:               clickRate,
		CpmPlatform:             cpmPlatform,
		DyFollow:                dyFollow,
		PayConvertRate:          payConvertRate,
		ConvertCost:             convertCost,
		ConvertRate:             convertRate,
		AveragePayOrderStatCost: averagePayOrderStatCost,
		PayOrderAveragePrice:    payOrderAveragePrice,
		ProductInfo:             sproductInfo,
		AwemeInfo:               sawemeInfo,
		DeliverySetting:         sdeliverySetting,
		CreateTime:              createTime,
		UpdateTime:              updateTime,
	}, nil
}

func (qair *qianchuanAdInfoRepo) GetByDay(ctx context.Context, adId uint64, startDay, endDay string) (*domain.QianchuanAdInfo, error) {
	where := make([]string, 0)

	where = append(where, "ad_id="+strconv.FormatUint(adId, 10))
	where = append(where, "day>='"+startDay+"'")
	where = append(where, "day<='"+endDay+"'")

	sql := "SELECT * FROM qianchuan_ad_info WHERE " + strings.Join(where, " AND ") + " order by day desc limit 1"

	row := qair.data.cdb.QueryRow(ctx, sql)

	var (
		qadId                   uint64
		qadvertiserId           uint64
		advertiserName          string
		campaignId              uint64
		campaignName            string
		campaignBudget          float64
		campaignBudgetMode      string
		campaignStatus          string
		campaignCreateDate      string
		promotionWay            string
		marketingGoal           string
		marketingScene          string
		name                    string
		status                  string
		optStatus               string
		adCreateTime            string
		adModifyTime            string
		labAdType               string
		statCost                float64
		roi                     float64
		showCnt                 int64
		clickCnt                int64
		payOrderCount           int64
		createOrderAmount       float64
		createOrderCount        int64
		payOrderAmount          float64
		dyFollow                int64
		convertCnt              int64
		clickRate               float64
		cpmPlatform             float64
		payConvertRate          float64
		convertCost             float64
		convertRate             float64
		averagePayOrderStatCost float64
		payOrderAveragePrice    float64
		productInfo             string
		awemeInfo               string
		deliverySetting         string
		qday                    time.Time
		createTime              time.Time
		updateTime              time.Time
	)

	if err := row.Scan(&qadId, &qadvertiserId, &advertiserName, &campaignId, &campaignName, &campaignBudget, &campaignBudgetMode, &campaignStatus, &campaignCreateDate, &promotionWay, &marketingGoal, &marketingScene, &name, &status, &optStatus, &adCreateTime, &adModifyTime, &labAdType, &statCost, &roi, &showCnt, &clickCnt, &payOrderCount, &createOrderAmount, &createOrderCount, &payOrderAmount, &dyFollow, &convertCnt, &clickRate, &cpmPlatform, &payConvertRate, &convertCost, &convertRate, &averagePayOrderStatCost, &payOrderAveragePrice, &productInfo, &awemeInfo, &deliverySetting, &qday, &createTime, &updateTime); err != nil {
		return nil, err
	}

	var sproductInfo []*domain.ProductInfo
	var sawemeInfo []*domain.AwemeInfo
	var sdeliverySetting domain.DeliverySetting

	json.Unmarshal([]byte(productInfo), &sproductInfo)
	json.Unmarshal([]byte(awemeInfo), &sawemeInfo)
	json.Unmarshal([]byte(deliverySetting), &sdeliverySetting)

	return &domain.QianchuanAdInfo{
		AdId:                    adId,
		AdvertiserId:            qadvertiserId,
		AdvertiserName:          advertiserName,
		CampaignId:              campaignId,
		CampaignName:            campaignName,
		CampaignBudget:          campaignBudget,
		CampaignBudgetMode:      campaignBudgetMode,
		CampaignStatus:          campaignStatus,
		CampaignCreateDate:      campaignCreateDate,
		PromotionWay:            promotionWay,
		MarketingGoal:           marketingGoal,
		MarketingScene:          marketingScene,
		Name:                    name,
		Status:                  status,
		OptStatus:               optStatus,
		AdCreateTime:            adCreateTime,
		AdModifyTime:            adModifyTime,
		LabAdType:               labAdType,
		StatCost:                statCost,
		Roi:                     roi,
		PayOrderCount:           payOrderCount,
		PayOrderAmount:          payOrderAmount,
		CreateOrderAmount:       createOrderAmount,
		CreateOrderCount:        createOrderCount,
		ClickCnt:                clickCnt,
		ShowCnt:                 showCnt,
		ConvertCnt:              convertCnt,
		ClickRate:               clickRate,
		CpmPlatform:             cpmPlatform,
		DyFollow:                dyFollow,
		PayConvertRate:          payConvertRate,
		ConvertCost:             convertCost,
		ConvertRate:             convertRate,
		AveragePayOrderStatCost: averagePayOrderStatCost,
		PayOrderAveragePrice:    payOrderAveragePrice,
		ProductInfo:             sproductInfo,
		AwemeInfo:               sawemeInfo,
		DeliverySetting:         sdeliverySetting,
		CreateTime:              createTime,
		UpdateTime:              updateTime,
	}, nil
}

func (qair *qianchuanAdInfoRepo) List(ctx context.Context, advertiserIds, startDay, endDay, keyword, filter, orderName, orderType string, pageNum, pageSize uint64) ([]*domain.QianchuanAdInfo, error) {
	list := make([]*domain.QianchuanAdInfo, 0)

	where := make([]string, 0)

	where = append(where, "advertiser_id in ("+advertiserIds+")")
	where = append(where, "day>='"+startDay+"'")
	where = append(where, "day<='"+endDay+"'")
	where = append(where, "status!='DELETE'")

	if filter == "videoPromGoods" {
		where = append(where, "marketing_goal='VIDEO_PROM_GOODS'")
	} else if filter == "livePromGoods" {
		where = append(where, "marketing_goal='LIVE_PROM_GOODS'")
	} else if filter == "notLabAd" {
		where = append(where, "lab_ad_type='NOT_LAB_AD'")
	} else if filter == "labAd" {
		where = append(where, "lab_ad_type='LAB_AD'")
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		if adId, err := strconv.ParseUint(keyword, 10, 64); err == nil {
			where = append(where, "(ad_id="+strconv.FormatUint(adId, 10)+" or name like '%"+keyword+"%')")
		} else {
			where = append(where, "name like '%"+keyword+"%'")
		}
	}

	var skip string

	if pageNum > 0 {
		skip = " LIMIT " + strconv.FormatUint((pageNum-1)*pageSize, 10) + "," + strconv.FormatUint(pageSize, 10)
	}

	sort := ""

	if orderName == "statCost" {
		if orderType == "asc" {
			sort = " order by stat_cost asc "
		} else {
			sort = " order by stat_cost desc "
		}
	} else if orderName == "adCreateTime" {
		if orderType == "asc" {
			sort = " order by ad_create_time asc "
		} else {
			sort = " order by ad_create_time desc "
		}
	}

	group := " group by ad_id "

	sql := "SELECT ad_id,sum(stat_cost) as stat_cost,sum(show_cnt) as show_cnt,sum(click_cnt) as click_cnt,sum(pay_order_count) as pay_order_count,sum(create_order_amount) as create_order_amount,sum(create_order_count) as create_order_count,sum(pay_order_amount) as pay_order_amount,sum(dy_follow) as dy_follow,sum(convert_cnt) as convert_cnt,max(ad_create_time) as ad_create_time FROM qianchuan_ad_info WHERE " + strings.Join(where, " AND ") + group + sort + skip

	rows, err := qair.data.cdb.Query(ctx, sql)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			adId              uint64
			statCost          float64
			showCnt           int64
			clickCnt          int64
			payOrderCount     int64
			createOrderAmount float64
			createOrderCount  int64
			payOrderAmount    float64
			dyFollow          int64
			convertCnt        int64
			adCreateTime      string
		)

		if err := rows.Scan(&adId, &statCost, &showCnt, &clickCnt, &payOrderCount, &createOrderAmount, &createOrderCount, &payOrderAmount, &dyFollow, &convertCnt, &adCreateTime); err != nil {
			return nil, err
		}

		list = append(list, &domain.QianchuanAdInfo{
			AdId:              adId,
			StatCost:          statCost,
			PayOrderCount:     payOrderCount,
			PayOrderAmount:    payOrderAmount,
			CreateOrderAmount: createOrderAmount,
			CreateOrderCount:  createOrderCount,
			ClickCnt:          clickCnt,
			ShowCnt:           showCnt,
			ConvertCnt:        convertCnt,
			DyFollow:          dyFollow,
		})
	}

	rows.Close()

	return list, nil
}

func (qair *qianchuanAdInfoRepo) ListNotLabAd(ctx context.Context, advertiserIds, day string) ([]*domain.QianchuanAdInfo, error) {
	list := make([]*domain.QianchuanAdInfo, 0)
	var qianchuanAdInfos []*QianchuanAdInfo
	var aadvertiserIds bson.A

	collection := qair.data.mdb.Database(qair.data.conf.Mongo.Dbname).Collection("qianchuan_ad_info_" + day)

	sadvertiserIds := strings.Split(advertiserIds, ",")

	for _, sadvertiserId := range sadvertiserIds {
		if uisadvertiserId, err := strconv.ParseUint(sadvertiserId, 10, 64); err == nil {
			aadvertiserIds = append(aadvertiserIds, uisadvertiserId)
		}
	}

	cursor, err := collection.Find(ctx, bson.M{
		"$and": []bson.M{
			bson.M{"advertiser_id": bson.D{{"$in", aadvertiserIds}}},
			bson.M{"lab_ad_type": bson.D{{"$eq", "NOT_LAB_AD"}}},
		},
	})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &qianchuanAdInfos)

	if err != nil {
		return nil, err
	}

	for _, qianchuanAdInfo := range qianchuanAdInfos {
		list = append(list, qianchuanAdInfo.ToDomain())
	}

	return list, nil
}

func (qair *qianchuanAdInfoRepo) ListAdvertiser(ctx context.Context, day string) ([]*domain.QianchuanAdvertiserInfo, error) {
	list := make([]*domain.QianchuanAdvertiserInfo, 0)

	collection := qair.data.mdb.Database(qair.data.conf.Mongo.Dbname).Collection("qianchuan_ad_info_" + day)

	groupStage := bson.D{
		{"$group", bson.M{
			"_id":                 "$advertiser_id",
			"id":                  bson.D{{"$first", "$advertiser_id"}},
			"stat_cost":           bson.D{{"$sum", "$stat_cost"}},
			"show_cnt":            bson.D{{"$sum", "$show_cnt"}},
			"click_cnt":           bson.D{{"$sum", "$click_cnt"}},
			"pay_order_count":     bson.D{{"$sum", "$pay_order_count"}},
			"create_order_amount": bson.D{{"$sum", "$create_order_amount"}},
			"create_order_count":  bson.D{{"$sum", "$create_order_count"}},
			"pay_order_amount":    bson.D{{"$sum", "$pay_order_amount"}},
			"dy_follow":           bson.D{{"$sum", "$dy_follow"}},
			"convert_cnt":         bson.D{{"$sum", "$convert_cnt"}},
		}},
	}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{groupStage})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var qianchuanAdvertiserInfos []QianchuanAdvertiserInfo

	err = cursor.All(ctx, &qianchuanAdvertiserInfos)

	if err != nil {
		return nil, err
	}

	for _, qianchuanAdvertiserInfo := range qianchuanAdvertiserInfos {
		list = append(list, qianchuanAdvertiserInfo.ToDomain())
	}

	return list, nil
}

func (qair *qianchuanAdInfoRepo) ListAdvertiserCampaigns(ctx context.Context, day string) ([]*domain.QianchuanAdvertiserInfo, error) {
	list := make([]*domain.QianchuanAdvertiserInfo, 0)

	collection := qair.data.mdb.Database(qair.data.conf.Mongo.Dbname).Collection("qianchuan_ad_info_" + day)

	matchStage := bson.D{
		{"$match", bson.M{"campaign_id": bson.D{{"$gt", 0}}}},
	}

	groupAStage := bson.D{
		{"$group", bson.M{
			"_id": bson.D{{"advertiser_id", "$advertiser_id"}, {"campaign_id", "$campaign_id"}},
		}},
	}

	projectStage := bson.D{
		{"$project", bson.M{
			"_id":           0,
			"advertiser_id": "$_id.advertiser_id",
			"campaign_id":   "$_id.campaign_id",
		}},
	}

	groupBStage := bson.D{
		{"$group", bson.M{
			"_id":       "$advertiser_id",
			"id":        bson.D{{"$first", "$advertiser_id"}},
			"campaigns": bson.D{{"$sum", 1}},
		}},
	}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupAStage, projectStage, groupBStage})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var qianchuanAdvertiserInfos []QianchuanAdvertiserInfo

	err = cursor.All(ctx, &qianchuanAdvertiserInfos)

	if err != nil {
		return nil, err
	}

	for _, qianchuanAdvertiserInfo := range qianchuanAdvertiserInfos {
		list = append(list, qianchuanAdvertiserInfo.ToDomain())
	}

	return list, nil
}

func (qair *qianchuanAdInfoRepo) Count(ctx context.Context, advertiserIds, startDay, endDay, keyword, filter string) (uint64, error) {
	where := make([]string, 0)

	where = append(where, "advertiser_id in ("+advertiserIds+")")
	where = append(where, "day>='"+startDay+"'")
	where = append(where, "day<='"+endDay+"'")
	where = append(where, "status!='DELETE'")

	if filter == "videoPromGoods" {
		where = append(where, "marketing_goal='VIDEO_PROM_GOODS'")
	} else if filter == "livePromGoods" {
		where = append(where, "marketing_goal='LIVE_PROM_GOODS'")
	} else if filter == "notLabAd" {
		where = append(where, "lab_ad_type='NOT_LAB_AD'")
	} else if filter == "labAd" {
		where = append(where, "lab_ad_type='LAB_AD'")
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		if adId, err := strconv.ParseUint(keyword, 10, 64); err == nil {
			where = append(where, "(ad_id="+strconv.FormatUint(adId, 10)+" or name like '%"+keyword+"%')")
		} else {
			where = append(where, "name like '%"+keyword+"%'")
		}
	}

	sql := "SELECT count(distinct(ad_id)) from qianchuan_ad_info WHERE " + strings.Join(where, " AND ")

	row := qair.data.cdb.QueryRow(ctx, sql)

	var total uint64

	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (qair *qianchuanAdInfoRepo) Statistics(ctx context.Context, advertiserIds, startDay, endDay, keyword, filter string) (*domain.QianchuanAdInfo, error) {
	where := make([]string, 0)

	where = append(where, "advertiser_id in ("+advertiserIds+")")

	if filter == "videoPromGoods" {
		where = append(where, "marketing_goal='VIDEO_PROM_GOODS'")
	} else if filter == "livePromGoods" {
		where = append(where, "marketing_goal='LIVE_PROM_GOODS'")
	} else if filter == "notLabAd" {
		where = append(where, "lab_ad_type='NOT_LAB_AD'")
	} else if filter == "labAd" {
		where = append(where, "lab_ad_type='LAB_AD'")
	}

	if l := utf8.RuneCountInString(keyword); l > 0 {
		if adId, err := strconv.ParseUint(keyword, 10, 64); err == nil {
			where = append(where, "(ad_id="+strconv.FormatUint(adId, 10)+" or name like '%"+keyword+"%')")
		} else {
			where = append(where, "name like '%"+keyword+"%'")
		}
	}

	where = append(where, "day>='"+startDay+"'")
	where = append(where, "day<='"+endDay+"'")

	sql := "SELECT sum(stat_cost),sum(show_cnt),sum(click_cnt),sum(pay_order_count),sum(create_order_amount),sum(create_order_count),sum(pay_order_amount),sum(dy_follow),sum(convert_cnt) from qianchuan_ad_info WHERE " + strings.Join(where, " AND ")

	row := qair.data.cdb.QueryRow(ctx, sql)

	var (
		statCost          float64
		showCnt           int64
		clickCnt          int64
		payOrderCount     int64
		createOrderAmount float64
		createOrderCount  int64
		payOrderAmount    float64
		dyFollow          int64
		convertCnt        int64
	)

	if err := row.Scan(&statCost, &showCnt, &clickCnt, &payOrderCount, &createOrderAmount, &createOrderCount, &payOrderAmount, &dyFollow, &convertCnt); err != nil {
		return nil, err
	}

	return &domain.QianchuanAdInfo{
		StatCost:          statCost,
		ShowCnt:           showCnt,
		ClickCnt:          clickCnt,
		PayOrderCount:     payOrderCount,
		CreateOrderAmount: createOrderAmount,
		CreateOrderCount:  createOrderCount,
		PayOrderAmount:    payOrderAmount,
		DyFollow:          dyFollow,
		ConvertCnt:        convertCnt,
	}, nil
}

func (qair *qianchuanAdInfoRepo) SaveIndex(ctx context.Context, day string) {
	collection := qair.data.mdb.Database(qair.data.conf.Mongo.Dbname).Collection("qianchuan_ad_info_" + day)

	isNotExistAIndex := true
	isNotExistBIndex := true
	isNotExistCIndex := true

	if indexSpecifications, err := collection.Indexes().ListSpecifications(ctx); err == nil {
		for _, indexSpecification := range indexSpecifications {
			if indexSpecification.Name == "ad_id_-1_advertiser_id_-1" {
				isNotExistAIndex = false
			}

			if indexSpecification.Name == "advertiser_id_-1" {
				isNotExistBIndex = false
			}

			if indexSpecification.Name == "campaign_id__-1" {
				isNotExistCIndex = false
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
					{Key: "campaign_id", Value: -1},
				},
			})
		}
	}
}

func (qlacr *qianchuanAdInfoRepo) UpsertQianchuanReportAd(ctx context.Context, day string, in *domain.QianchuanAdInfo) error {
	collection := qlacr.data.mdb.Database(qlacr.data.conf.Mongo.Dbname).Collection("qianchuan_ad_info_" + day)

	if _, err := collection.UpdateOne(ctx, bson.D{
		{"ad_id", in.AdId},
		{"advertiser_id", in.AdvertiserId},
	}, bson.D{
		{"$set", bson.D{
			{"stat_cost", in.StatCost},
			{"roi", in.Roi},
			{"show_cnt", in.ShowCnt},
			{"click_cnt", in.ClickCnt},
			{"pay_order_count", in.PayOrderCount},
			{"create_order_amount", in.CreateOrderAmount},
			{"create_order_count", in.CreateOrderCount},
			{"pay_order_amount", in.PayOrderAmount},
			{"dy_follow", in.DyFollow},
			{"convert_cnt", in.ConvertCnt},
			{"click_rate", in.ClickRate},
			{"cpm_platform", in.CpmPlatform},
			{"pay_convert_rate", in.PayConvertRate},
			{"convert_cost", in.ConvertCost},
			{"convert_rate", in.ConvertRate},
			{"average_pay_order_stat_cost", in.AveragePayOrderStatCost},
			{"pay_order_average_price", in.PayOrderAveragePrice},
			{"create_time", in.CreateTime},
			{"update_time", in.UpdateTime},
		}},
	}, options.Update().SetUpsert(true)); err != nil {
		return err
	}

	return nil
}

func (qlacr *qianchuanAdInfoRepo) UpsertQianchuanAd(ctx context.Context, day string, in *domain.QianchuanAdInfo) error {
	collection := qlacr.data.mdb.Database(qlacr.data.conf.Mongo.Dbname).Collection("qianchuan_ad_info_" + day)

	if _, err := collection.UpdateOne(ctx, bson.D{
		{"ad_id", in.AdId},
		{"advertiser_id", in.AdvertiserId},
	}, bson.D{
		{"$set", bson.D{
			{"campaign_id", in.CampaignId},
			{"advertiser_name", in.AdvertiserName},
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

func (qlacr *qianchuanAdInfoRepo) UpsertQianchuanCampaign(ctx context.Context, day string, in *domain.QianchuanAdInfo) error {
	collection := qlacr.data.mdb.Database(qlacr.data.conf.Mongo.Dbname).Collection("qianchuan_ad_info_" + day)

	if _, err := collection.UpdateMany(ctx, bson.D{
		{"advertiser_id", in.AdvertiserId},
		{"campaign_id", in.CampaignId},
	}, bson.D{
		{"$set", bson.D{
			{"campaign_name", in.CampaignName},
			{"campaign_budget", in.CampaignBudget},
			{"campaign_budget_mode", in.CampaignBudgetMode},
			{"campaign_status", in.CampaignStatus},
			{"campaign_create_date", in.CampaignCreateDate},
		}},
	}); err != nil {
		return err
	}

	return nil
}

package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
	"strings"
	"time"
)

// 千川随心推订单完整信息表
type QianchuanAwemeOrderInfo struct {
	OrderId              uint64                  `json:"order_id" bson:"order_id"`
	AdvertiserId         uint64                  `json:"advertiser_id" bson:"advertiser_id"`
	AdId                 uint64                  `json:"ad_id" bson:"ad_id"`
	MarketingGoal        string                  `json:"marketing_goal" bson:"marketing_goal"`
	Status               string                  `json:"status" bson:"status"`
	OrderCreateTime      time.Time               `json:"order_create_time" bson:"order_create_time"`
	AwemeInfo            *domain.AwemeInfo       `json:"aweme_info" bson:"aweme_info"`
	VideoInfo            *domain.VideoInfo       `json:"video_info" bson:"video_info"`
	RoomInfo             *domain.RoomInfo        `json:"room_info" bson:"room_info"`
	DeliverySetting      *domain.DeliverySetting `json:"delivery_setting" bson:"delivery_setting"`
	FailList             []uint64                `json:"fail_list" bson:"fail_list"`
	PayOrderAmount       float64                 `json:"pay_order_amount" bson:"pay_order_amount"`
	StatCost             float64                 `json:"stat_cost" bson:"stat_cost"`
	PrepayAndPayOrderRoi float64                 `json:"prepay_and_pay_order_roi" bson:"prepay_and_pay_order_roi"`
	TotalPlay            float64                 `json:"total_play" bson:"total_play"`
	ShowCnt              int64                   `json:"show_cnt" bson:"show_cnt"`
	Ctr                  int64                   `json:"ctr" bson:"ctr"`
	ClickCnt             int64                   `json:"click_cnt" bson:"click_cnt"`
	PayOrderCount        int64                   `json:"pay_order_count" bson:"pay_order_count"`
	PrepayOrderCount     int64                   `json:"prepay_order_count" bson:"prepay_order_count"`
	PrepayOrderAmount    int64                   `json:"prepay_order_amount" bson:"prepay_order_amount"`
	DyFollow             int64                   `json:"dy_follow" bson:"dy_follow"`
	DyShare              int64                   `json:"dy_share" bson:"dy_share"`
}

type qianchuanAwemeOrderInfoRepo struct {
	data *Data
	log  *log.Helper
}

func (qaoi *QianchuanAwemeOrderInfo) ToDomain() *domain.QianchuanAwemeOrderInfo {
	return &domain.QianchuanAwemeOrderInfo{
		OrderId:              qaoi.OrderId,
		AdvertiserId:         qaoi.AdvertiserId,
		AdId:                 qaoi.AdId,
		MarketingGoal:        qaoi.MarketingGoal,
		Status:               qaoi.Status,
		OrderCreateTime:      qaoi.OrderCreateTime,
		AwemeInfo:            qaoi.AwemeInfo,
		VideoInfo:            qaoi.VideoInfo,
		RoomInfo:             qaoi.RoomInfo,
		DeliverySetting:      qaoi.DeliverySetting,
		FailList:             qaoi.FailList,
		PayOrderAmount:       qaoi.PayOrderAmount,
		StatCost:             qaoi.StatCost,
		PrepayAndPayOrderRoi: qaoi.PrepayAndPayOrderRoi,
		TotalPlay:            qaoi.TotalPlay,
		ShowCnt:              qaoi.ShowCnt,
		Ctr:                  qaoi.Ctr,
		ClickCnt:             qaoi.ClickCnt,
		PayOrderCount:        qaoi.PayOrderCount,
		PrepayOrderCount:     qaoi.PrepayOrderCount,
		PrepayOrderAmount:    qaoi.PrepayOrderAmount,
		DyFollow:             qaoi.DyFollow,
		DyShare:              qaoi.DyShare,
	}
}

func NewQianchuanAwemeOrderInfoRepo(data *Data, logger log.Logger) biz.QianchuanAwemeOrderInfoRepo {
	return &qianchuanAwemeOrderInfoRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (qaoir *qianchuanAwemeOrderInfoRepo) List(ctx context.Context, advertiserIds, day, marketingGoal string) ([]*domain.AwemeVideoProductQianchuanAwemeOrderInfo, error) {
	list := make([]*domain.AwemeVideoProductQianchuanAwemeOrderInfo, 0)

	var aadvertiserIds bson.A
	var and []bson.M
	var cursor *mongo.Cursor
	var err error

	sadvertiserIds := strings.Split(advertiserIds, ",")

	for _, sadvertiserId := range sadvertiserIds {
		if uisadvertiserId, err := strconv.ParseUint(sadvertiserId, 10, 64); err == nil {
			aadvertiserIds = append(aadvertiserIds, uisadvertiserId)
		}
	}

	collection := qaoir.data.mdb.Database(qaoir.data.conf.Mongo.Dbname).Collection("qianchuan_aweme_order_info_" + day)

	and = append(and, bson.M{"advertiser_id": bson.D{{"$in", aadvertiserIds}}})

	if marketingGoal == "VIDEO_PROM_GOODS" {
		and = append(and, bson.M{"marketing_goal": bson.D{{"$eq", "VIDEO_PROM_GOODS"}}})

		matchStage := bson.D{
			{"$match", bson.M{
				"$and": and,
			}},
		}

		projectAStage := bson.D{
			{"$project", bson.M{
				"aweme_id":         "$aweme_info.aweme_id",
				"video_id":         "$video_info.aweme_item_id",
				"product_id":       "$product_info.id",
				"pay_order_amount": 1,
				"stat_cost":        1,
			}},
		}

		groupStage := bson.D{
			{"$group", bson.M{
				"_id":              bson.D{{"aweme_id", "$aweme_id"}, {"video_id", "$video_id"}, {"product_id", "$product_id"}},
				"pay_order_amount": bson.D{{"$sum", "$pay_order_amount"}},
				"stat_cost":        bson.D{{"$sum", "$stat_cost"}},
			}},
		}

		projectBStage := bson.D{
			{"$project", bson.M{
				"aweme_id":         "$_id.aweme_id",
				"video_id":         "$_id.video_id",
				"product_id":       "$_id.product_id",
				"pay_order_amount": 1,
				"stat_cost":        1,
			}},
		}

		cursor, err = collection.Aggregate(ctx, mongo.Pipeline{matchStage, projectAStage, groupStage, projectBStage})

		if err != nil {
			return nil, err
		}

		defer cursor.Close(ctx)

		err = cursor.All(ctx, &list)

		if err != nil {
			return nil, err
		}
	} else if marketingGoal == "LIVE_PROM_GOODS" {
		and = append(and, bson.M{"marketing_goal": bson.D{{"$eq", "LIVE_PROM_GOODS"}}})

		matchStage := bson.D{
			{"$match", bson.M{
				"$and": and,
			}},
		}

		projectAStage := bson.D{
			{"$project", bson.M{
				"aweme_id":         "$aweme_info.aweme_id",
				"pay_order_amount": 1,
				"stat_cost":        1,
			}},
		}

		groupStage := bson.D{
			{"$group", bson.M{
				"_id":              bson.D{{"aweme_id", "$aweme_id"}},
				"pay_order_amount": bson.D{{"$sum", "$pay_order_amount"}},
				"stat_cost":        bson.D{{"$sum", "$stat_cost"}},
			}},
		}

		projectBStage := bson.D{
			{"$project", bson.M{
				"aweme_id":         "$_id.aweme_id",
				"pay_order_amount": 1,
				"stat_cost":        1,
			}},
		}

		cursor, err = collection.Aggregate(ctx, mongo.Pipeline{matchStage, projectAStage, groupStage, projectBStage})

		if err != nil {
			return nil, err
		}

		defer cursor.Close(ctx)

		err = cursor.All(ctx, &list)

		if err != nil {
			return nil, err
		}
	}

	return list, nil
}

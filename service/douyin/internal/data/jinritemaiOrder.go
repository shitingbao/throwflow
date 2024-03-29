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

// 精选联盟达人订单表
type JinritemaiOrder struct {
	OrderId                        string         `json:"order_id" bson:"order_id"`
	ProductId                      string         `json:"product_id" bson:"product_id"`
	ProductName                    string         `json:"product_name" bson:"product_name"`
	ProductImg                     string         `json:"product_img" bson:"product_img"`
	AuthorAccount                  string         `json:"author_account" bson:"author_account"`
	AuthorClientKey                string         `json:"author_client_key" bson:"author_client_key"`
	AuthorOpenId                   string         `json:"author_open_id" bson:"author_open_id"`
	ShopName                       string         `json:"shop_name" bson:"shop_name"`
	TotalPayAmount                 float64        `json:"total_pay_amount" bson:"total_pay_amount"`
	CommissionRate                 float64        `json:"commission_rate" bson:"commission_rate"`
	FlowPoint                      string         `json:"flow_point" bson:"flow_point"`
	App                            string         `json:"app" bson:"app"`
	UpdateTime                     string         `json:"update_time" bson:"update_time"`
	PaySuccessTime                 string         `json:"pay_success_time" bson:"pay_success_time"`
	SettleTime                     string         `json:"settle_time" bson:"settle_time"`
	PayGoodsAmount                 int64          `json:"pay_goods_amount" bson:"pay_goods_amount"`
	SettledGoodsAmount             int64          `json:"settled_goods_amount" bson:"settled_goods_amount"`
	EstimatedCommission            int64          `json:"estimated_commission" bson:"estimated_commission"`
	RealCommission                 int64          `json:"real_commission" bson:"real_commission"`
	Extra                          string         `json:"extra" bson:"extra"`
	ItemNum                        int64          `json:"item_num" bson:"item_num"`
	ShopId                         int64          `json:"shop_id" bson:"shop_id"`
	RefundTime                     string         `json:"refund_time" bson:"refund_time"`
	PidInfo                        domain.PidInfo `json:"pid_info" bson:"pid_info"`
	EstimatedTotalCommission       int64          `json:"estimated_total_commission" bson:"estimated_total_commission"`
	EstimatedTechServiceFee        int64          `json:"estimated_tech_service_fee" bson:"estimated_tech_service_fee"`
	PickSourceClientKey            string         `json:"pick_source_client_key" bson:"pick_source_client_key"`
	PickExtra                      string         `json:"pick_extra" bson:"pick_extra"`
	AuthorShortId                  string         `json:"author_short_id" bson:"author_short_id"`
	MediaType                      string         `json:"media_type" bson:"media_type"`
	IsSteppedPlan                  bool           `json:"is_stepped_plan" bson:"is_stepped_plan"`
	PlatformSubsidy                int64          `json:"platform_subsidy" bson:"platform_subsidy"`
	AuthorSubsidy                  int64          `json:"author_subsidy" bson:"author_subsidy"`
	ProductActivityId              string         `json:"product_activity_id" bson:"product_activity_id"`
	AppId                          int64          `json:"app_id" bson:"app_id"`
	SettleUserSteppedCommission    int64          `json:"settle_user_stepped_commission" bson:"settle_user_stepped_commission"`
	SettleInstSteppedCommission    int64          `json:"settle_inst_stepped_commission" bson:"settle_inst_stepped_commission"`
	PaySubsidy                     int64          `json:"pay_subsidy" bson:"pay_subsidy"`
	MediaId                        int64          `json:"media_id" bson:"media_id"`
	AuthorBuyinId                  string         `json:"author_buyin_id" bson:"author_buyin_id"`
	ConfirmTime                    string         `json:"confirm_time" bson:"confirm_time"`
	EstimatedInstSteppedCommission int64          `json:"estimated_inst_stepped_commission" bson:"estimated_inst_stepped_commission"`
	EstimatedUserSteppedCommission int64          `json:"estimated_user_stepped_commission" bson:"estimated_user_stepped_commission"`
}

type jinritemaiOrderRepo struct {
	data *Data
	log  *log.Helper
}

func (jo *JinritemaiOrder) ToDomain() *domain.JinritemaiOrder {
	return &domain.JinritemaiOrder{
		OrderId:                        jo.OrderId,
		ProductId:                      jo.ProductId,
		ProductName:                    jo.ProductName,
		ProductImg:                     jo.ProductImg,
		AuthorAccount:                  jo.AuthorAccount,
		AuthorClientKey:                jo.AuthorClientKey,
		AuthorOpenId:                   jo.AuthorOpenId,
		ShopName:                       jo.ShopName,
		TotalPayAmount:                 jo.TotalPayAmount,
		CommissionRate:                 jo.CommissionRate,
		FlowPoint:                      jo.FlowPoint,
		App:                            jo.App,
		UpdateTime:                     jo.UpdateTime,
		PaySuccessTime:                 jo.PaySuccessTime,
		SettleTime:                     jo.SettleTime,
		PayGoodsAmount:                 jo.PayGoodsAmount,
		SettledGoodsAmount:             jo.SettledGoodsAmount,
		EstimatedCommission:            jo.EstimatedCommission,
		RealCommission:                 jo.RealCommission,
		Extra:                          jo.Extra,
		ItemNum:                        jo.ItemNum,
		ShopId:                         jo.ShopId,
		RefundTime:                     jo.RefundTime,
		PidInfo:                        jo.PidInfo,
		EstimatedTotalCommission:       jo.EstimatedTotalCommission,
		EstimatedTechServiceFee:        jo.EstimatedTechServiceFee,
		PickSourceClientKey:            jo.PickSourceClientKey,
		PickExtra:                      jo.PickExtra,
		AuthorShortId:                  jo.AuthorShortId,
		MediaType:                      jo.MediaType,
		IsSteppedPlan:                  jo.IsSteppedPlan,
		PlatformSubsidy:                jo.PlatformSubsidy,
		AuthorSubsidy:                  jo.AuthorSubsidy,
		ProductActivityId:              jo.ProductActivityId,
		AppId:                          jo.AppId,
		SettleUserSteppedCommission:    jo.SettleUserSteppedCommission,
		SettleInstSteppedCommission:    jo.SettleInstSteppedCommission,
		PaySubsidy:                     jo.PaySubsidy,
		MediaId:                        jo.MediaId,
		AuthorBuyinId:                  jo.AuthorBuyinId,
		ConfirmTime:                    jo.ConfirmTime,
		EstimatedInstSteppedCommission: jo.EstimatedInstSteppedCommission,
		EstimatedUserSteppedCommission: jo.EstimatedUserSteppedCommission,
	}
}

func NewJinritemaiOrderRepo(data *Data, logger log.Logger) biz.JinritemaiOrderRepo {
	return &jinritemaiOrderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (jor *jinritemaiOrderRepo) ListByPickExtra(ctx context.Context) ([]*domain.JinritemaiOrder, error) {
	list := make([]*domain.JinritemaiOrder, 0)

	collection := jor.data.mdb.Database(jor.data.conf.Mongo.Dbname).Collection("jinritemai_order")

	matchStage := bson.D{
		{"$match", bson.M{"pick_extra": bson.D{{"$exists", true}, {"$ne", ""}}}},
	}

	groupStage := bson.D{
		{"$group", bson.M{"_id": bson.D{{"author_client_key", "$author_client_key"}, {"author_open_id", "$author_open_id"}}}},
	}

	projectStage := bson.D{
		{"$project", bson.M{
			"_id":               0,
			"author_client_key": "$_id.author_client_key",
			"author_open_id":    "$_id.author_open_id",
		}},
	}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var jinritemaiOrders []JinritemaiOrder

	err = cursor.All(ctx, &jinritemaiOrders)

	if err != nil {
		return nil, err
	}

	for _, jinritemaiOrder := range jinritemaiOrders {
		list = append(list, &domain.JinritemaiOrder{
			AuthorClientKey: jinritemaiOrder.AuthorClientKey,
			AuthorOpenId:    jinritemaiOrder.AuthorOpenId,
		})
	}

	return list, nil
}

func (jor *jinritemaiOrderRepo) SaveIndex(ctx context.Context) {
	collection := jor.data.mdb.Database(jor.data.conf.Mongo.Dbname).Collection("jinritemai_order")

	isNotExistIndex := true

	if indexSpecifications, err := collection.Indexes().ListSpecifications(ctx); err == nil {
		for _, indexSpecification := range indexSpecifications {
			if indexSpecification.Name == "author_client_key_-1_author_open_id_-1_order_id_-1" {
				isNotExistIndex = false
			}
		}

		if isNotExistIndex {
			collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{Key: "author_client_key", Value: -1},
					{Key: "author_open_id", Value: -1},
					{Key: "order_id", Value: -1},
				},
			})
		}
	}
}

func (jor *jinritemaiOrderRepo) Upsert(ctx context.Context, in *domain.JinritemaiOrder) error {
	collection := jor.data.mdb.Database(jor.data.conf.Mongo.Dbname).Collection("jinritemai_order")

	if _, err := collection.UpdateOne(ctx, bson.D{
		{"author_client_key", in.AuthorClientKey},
		{"author_open_id", in.AuthorOpenId},
		{"order_id", in.OrderId},
	}, bson.D{
		{"$set", bson.D{
			{"product_id", in.ProductId},
			{"product_name", in.ProductName},
			{"product_img", in.ProductImg},
			{"author_account", in.AuthorAccount},
			{"shop_name", in.ShopName},
			{"total_pay_amount", in.TotalPayAmount},
			{"commission_rate", in.CommissionRate},
			{"flow_point", in.FlowPoint},
			{"app", in.App},
			{"update_time", in.UpdateTime},
			{"pay_success_time", in.PaySuccessTime},
			{"settle_time", in.SettleTime},
			{"pay_goods_amount", in.PayGoodsAmount},
			{"settled_goods_amount", in.SettledGoodsAmount},
			{"estimated_commission", in.EstimatedCommission},
			{"real_commission", in.RealCommission},
			{"extra", in.Extra},
			{"item_num", in.ItemNum},
			{"shop_id", in.ShopId},
			{"refund_time", in.RefundTime},
			{"pid_info", in.PidInfo},
			{"estimated_total_commission", in.EstimatedTotalCommission},
			{"estimated_tech_service_fee", in.EstimatedTechServiceFee},
			{"pick_source_client_key", in.PickSourceClientKey},
			{"pick_extra", in.PickExtra},
			{"author_short_id", in.AuthorShortId},
			{"media_type", in.MediaType},
			{"is_stepped_plan", in.IsSteppedPlan},
			{"platform_subsidy", in.PlatformSubsidy},
			{"author_subsidy", in.AuthorSubsidy},
			{"product_activity_id", in.ProductActivityId},
			{"app_id", in.AppId},
			{"settle_user_stepped_commission", in.SettleUserSteppedCommission},
			{"settle_inst_stepped_commission", in.SettleInstSteppedCommission},
			{"pay_subsidy", in.PaySubsidy},
			{"media_id", in.MediaId},
			{"author_buyin_id", in.AuthorBuyinId},
			{"confirm_time", in.ConfirmTime},
			{"estimated_inst_stepped_commission", in.EstimatedInstSteppedCommission},
			{"estimated_user_stepped_commission", in.EstimatedUserSteppedCommission},
		}},
	}, options.Update().SetUpsert(true)); err != nil {
		return err
	}

	return nil
}

func (jor *jinritemaiOrderRepo) Send(ctx context.Context, message event.Event) error {
	if err := jor.data.kafka.Send(ctx, message); err != nil {
		return err
	}

	return nil
}

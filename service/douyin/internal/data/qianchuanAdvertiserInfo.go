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

// 千川账户完整信息表
type QianchuanAdvertiserInfo struct {
	Id                      uint64    `json:"id" bson:"id"`
	Name                    string    `json:"name" bson:"name"`
	GeneralTotalBalance     float64   `json:"general_total_balance" bson:"general_total_balance"`
	Campaigns               uint64    `json:"campaigns" bson:"campaigns"`
	StatCost                float64   `json:"stat_cost" bson:"stat_cost"`
	Roi                     float64   `json:"roi" bson:"roi"`
	PayOrderCount           int64     `json:"pay_order_count" bson:"pay_order_count"`
	PayOrderAmount          float64   `json:"pay_order_amount" bson:"pay_order_amount"`
	CreateOrderAmount       float64   `json:"create_order_amount" bson:"create_order_amount"`
	CreateOrderCount        int64     `json:"create_order_count" bson:"create_order_count"`
	ClickCnt                int64     `json:"click_cnt" bson:"click_cnt"`
	ShowCnt                 int64     `json:"show_cnt" bson:"show_cnt"`
	ConvertCnt              int64     `json:"convert_cnt" bson:"convert_cnt"`
	ClickRate               float64   `json:"click_rate" bson:"click_rate"`
	CpmPlatform             float64   `json:"cpm_platform" bson:"cpm_platform"`
	DyFollow                int64     `json:"dy_follow" bson:"dy_follow"`
	PayConvertRate          float64   `json:"pay_convert_rate" bson:"pay_convert_rate"`
	ConvertCost             float64   `json:"convert_cost" bson:"convert_cost"`
	ConvertRate             float64   `json:"convert_rate" bson:"convert_rate"`
	AveragePayOrderStatCost float64   `json:"average_pay_order_stat_cost" bson:"average_pay_order_stat_cost"`
	PayOrderAveragePrice    float64   `json:"pay_order_average_price" bson:"pay_order_average_price"`
	CreateTime              time.Time `json:"create_time" bson:"create_time"`
	UpdateTime              time.Time `json:"update_time" bson:"update_time"`
}

type qianchuanAdvertiserInfoRepo struct {
	data *Data
	log  *log.Helper
}

func (qai *QianchuanAdvertiserInfo) ToDomain() *domain.QianchuanAdvertiserInfo {
	return &domain.QianchuanAdvertiserInfo{
		Id:                      qai.Id,
		Name:                    qai.Name,
		GeneralTotalBalance:     qai.GeneralTotalBalance,
		Campaigns:               qai.Campaigns,
		StatCost:                qai.StatCost,
		Roi:                     qai.Roi,
		PayOrderCount:           qai.PayOrderCount,
		PayOrderAmount:          qai.PayOrderAmount,
		CreateOrderAmount:       qai.CreateOrderAmount,
		CreateOrderCount:        qai.CreateOrderCount,
		ClickCnt:                qai.ClickCnt,
		ShowCnt:                 qai.ShowCnt,
		ConvertCnt:              qai.ConvertCnt,
		ClickRate:               qai.ClickRate,
		CpmPlatform:             qai.CpmPlatform,
		DyFollow:                qai.DyFollow,
		PayConvertRate:          qai.PayConvertRate,
		ConvertCost:             qai.ConvertCost,
		ConvertRate:             qai.ConvertRate,
		AveragePayOrderStatCost: qai.AveragePayOrderStatCost,
		PayOrderAveragePrice:    qai.PayOrderAveragePrice,
		CreateTime:              qai.CreateTime,
		UpdateTime:              qai.UpdateTime,
	}
}

func NewQianchuanAdvertiserInfoRepo(data *Data, logger log.Logger) biz.QianchuanAdvertiserInfoRepo {
	return &qianchuanAdvertiserInfoRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (qair *qianchuanAdvertiserInfoRepo) Get(ctx context.Context, advertiserId uint64, day string) (*domain.QianchuanAdvertiserInfo, error) {
	where := make([]string, 0)

	where = append(where, "id="+strconv.FormatUint(advertiserId, 10))
	where = append(where, "day='"+day+"'")

	sql := "SELECT * FROM qianchuan_advertiser_info WHERE " + strings.Join(where, " AND ")

	row := qair.data.cdb.QueryRow(ctx, sql)

	var (
		id                      uint64
		name                    string
		generalTotalBalance     float64
		campaigns               uint64
		statCost                float64
		roi                     float64
		payOrderCount           int64
		payOrderAmount          float64
		createOrderAmount       float64
		createOrderCount        int64
		clickCnt                int64
		showCnt                 int64
		convertCnt              int64
		clickRate               float64
		cpmPlatform             float64
		dyFollow                int64
		payConvertRate          float64
		convertCost             float64
		convertRate             float64
		averagePayOrderStatCost float64
		payOrderAveragePrice    float64
		qday                    time.Time
		createTime              time.Time
		updateTime              time.Time
	)

	if err := row.Scan(&id, &name, &generalTotalBalance, &campaigns, &statCost, &roi, &showCnt, &clickCnt, &payOrderCount, &createOrderAmount, &createOrderCount, &payOrderAmount, &dyFollow, &convertCnt, &clickRate, &cpmPlatform, &payConvertRate, &convertCost, &convertRate, &averagePayOrderStatCost, &payOrderAveragePrice, &qday, &createTime, &updateTime); err != nil {
		return nil, err
	}

	return &domain.QianchuanAdvertiserInfo{
		Id:                      id,
		Name:                    name,
		GeneralTotalBalance:     generalTotalBalance,
		Campaigns:               campaigns,
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
		CreateTime:              createTime,
		UpdateTime:              updateTime,
	}, nil
}

func (qair *qianchuanAdvertiserInfoRepo) GetByDay(ctx context.Context, advertiserId uint64, startDay, endDay string) (*domain.QianchuanAdvertiserInfo, error) {
	where := make([]string, 0)

	where = append(where, "id="+strconv.FormatUint(advertiserId, 10))
	where = append(where, "day>='"+startDay+"'")
	where = append(where, "day<='"+endDay+"'")

	sql := "SELECT * FROM qianchuan_advertiser_info WHERE " + strings.Join(where, " AND ") + " order by day desc limit 1"

	row := qair.data.cdb.QueryRow(ctx, sql)

	var (
		id                      uint64
		name                    string
		generalTotalBalance     float64
		campaigns               uint64
		statCost                float64
		roi                     float64
		payOrderCount           int64
		payOrderAmount          float64
		createOrderAmount       float64
		createOrderCount        int64
		clickCnt                int64
		showCnt                 int64
		convertCnt              int64
		clickRate               float64
		cpmPlatform             float64
		dyFollow                int64
		payConvertRate          float64
		convertCost             float64
		convertRate             float64
		averagePayOrderStatCost float64
		payOrderAveragePrice    float64
		qday                    time.Time
		createTime              time.Time
		updateTime              time.Time
	)

	if err := row.Scan(&id, &name, &generalTotalBalance, &campaigns, &statCost, &roi, &showCnt, &clickCnt, &payOrderCount, &createOrderAmount, &createOrderCount, &payOrderAmount, &dyFollow, &convertCnt, &clickRate, &cpmPlatform, &payConvertRate, &convertCost, &convertRate, &averagePayOrderStatCost, &payOrderAveragePrice, &qday, &createTime, &updateTime); err != nil {
		return nil, err
	}

	return &domain.QianchuanAdvertiserInfo{
		Id:                      id,
		Name:                    name,
		GeneralTotalBalance:     generalTotalBalance,
		Campaigns:               campaigns,
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
		CreateTime:              createTime,
		UpdateTime:              updateTime,
	}, nil
}

func (qair *qianchuanAdvertiserInfoRepo) List(ctx context.Context, advertiserIds, startDay, endDay string, pageNum, pageSize uint64) ([]*domain.QianchuanAdvertiserInfo, error) {
	list := make([]*domain.QianchuanAdvertiserInfo, 0)

	where := make([]string, 0)

	where = append(where, "id in ("+advertiserIds+")")
	where = append(where, "day>='"+startDay+"'")
	where = append(where, "day<='"+endDay+"'")

	var skip string

	if pageNum > 0 {
		skip = " LIMIT " + strconv.FormatUint((pageNum-1)*pageSize, 10) + "," + strconv.FormatUint(pageSize, 10)
	}

	sort := " order by stat_cost desc "
	group := " group by id "

	sql := "SELECT id,sum(stat_cost) as stat_cost,sum(show_cnt) as show_cnt,sum(click_cnt) as click_cnt,sum(pay_order_count) as pay_order_count,sum(create_order_amount) as create_order_amount,sum(create_order_count) as create_order_count,sum(pay_order_amount) as pay_order_amount,sum(dy_follow) as dy_follow,sum(convert_cnt) as convert_cnt FROM qianchuan_advertiser_info WHERE " + strings.Join(where, " AND ") + group + sort + skip

	rows, err := qair.data.cdb.Query(ctx, sql)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id                uint64
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

		if err := rows.Scan(&id, &statCost, &showCnt, &clickCnt, &payOrderCount, &createOrderAmount, &createOrderCount, &payOrderAmount, &dyFollow, &convertCnt); err != nil {
			return nil, err
		}

		list = append(list, &domain.QianchuanAdvertiserInfo{
			Id:                id,
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

func (qair *qianchuanAdvertiserInfoRepo) Count(ctx context.Context, advertiserIds, startDay, endDay string) (uint64, error) {
	where := make([]string, 0)

	where = append(where, "id in ("+advertiserIds+")")
	where = append(where, "day>='"+startDay+"'")
	where = append(where, "day<='"+endDay+"'")

	sql := "SELECT count(distinct(id)) from qianchuan_advertiser_info WHERE " + strings.Join(where, " AND ")

	row := qair.data.cdb.QueryRow(ctx, sql)

	var total uint64

	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (qair *qianchuanAdvertiserInfoRepo) Statistics(ctx context.Context, advertiserIds, startDay, endDay string) (*domain.QianchuanAdvertiserInfo, error) {
	where := make([]string, 0)

	where = append(where, "id in ("+advertiserIds+")")
	where = append(where, "day>='"+startDay+"'")
	where = append(where, "day<='"+endDay+"'")

	sql := "SELECT sum(campaigns),sum(stat_cost),sum(pay_order_count),sum(pay_order_amount),sum(create_order_amount),sum(create_order_count),sum(click_cnt),sum(show_cnt),sum(convert_cnt),sum(dy_follow) from qianchuan_advertiser_info WHERE " + strings.Join(where, " AND ")
	row := qair.data.cdb.QueryRow(ctx, sql)

	var (
		campaigns         uint64
		statCost          float64
		payOrderCount     int64
		payOrderAmount    float64
		createOrderAmount float64
		createOrderCount  int64
		clickCnt          int64
		showCnt           int64
		convertCnt        int64
		dyFollow          int64
	)

	if err := row.Scan(&campaigns, &statCost, &payOrderCount, &payOrderAmount, &createOrderAmount, &createOrderCount, &clickCnt, &showCnt, &convertCnt, &dyFollow); err != nil {
		return nil, err
	}

	where = make([]string, 0)

	where = append(where, "id in ("+advertiserIds+")")
	where = append(where, "day='"+endDay+"'")

	sql = "SELECT sum(general_total_balance) from qianchuan_advertiser_info WHERE " + strings.Join(where, " AND ")
	row = qair.data.cdb.QueryRow(ctx, sql)

	var (
		generalTotalBalance float64
	)

	if err := row.Scan(&generalTotalBalance); err != nil {
		return nil, err
	}

	return &domain.QianchuanAdvertiserInfo{
		GeneralTotalBalance: generalTotalBalance,
		Campaigns:           campaigns,
		StatCost:            statCost,
		PayOrderCount:       payOrderCount,
		PayOrderAmount:      payOrderAmount,
		CreateOrderAmount:   createOrderAmount,
		CreateOrderCount:    createOrderCount,
		ClickCnt:            clickCnt,
		ShowCnt:             showCnt,
		ConvertCnt:          convertCnt,
		DyFollow:            dyFollow,
	}, nil
}

func (qair *qianchuanAdvertiserInfoRepo) SaveIndex(ctx context.Context, day string) {
	collection := qair.data.mdb.Database(qair.data.conf.Mongo.Dbname).Collection("qianchuan_advertiser_info_" + day)

	isNotExistIndex := true

	if indexSpecifications, err := collection.Indexes().ListSpecifications(ctx); err == nil {
		for _, indexSpecification := range indexSpecifications {
			if indexSpecification.Name == "id_-1" {
				isNotExistIndex = false
			}
		}

		if isNotExistIndex {
			collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{Key: "id", Value: -1},
				},
			})
		}
	}
}

func (qair *qianchuanAdvertiserInfoRepo) Upsert(ctx context.Context, day string, in *domain.QianchuanAdvertiserInfo) error {
	collection := qair.data.mdb.Database(qair.data.conf.Mongo.Dbname).Collection("qianchuan_advertiser_info_" + day)

	if _, err := collection.UpdateOne(ctx, bson.D{
		{"id", in.Id},
	}, bson.D{
		{"$set", bson.D{
			{"name", in.Name},
			{"general_total_balance", in.GeneralTotalBalance},
			{"campaigns", in.Campaigns},
			{"stat_cost", in.StatCost},
			{"roi", in.Roi},
			{"pay_order_count", in.PayOrderCount},
			{"pay_order_amount", in.PayOrderAmount},
			{"create_order_amount", in.CreateOrderAmount},
			{"create_order_count", in.CreateOrderCount},
			{"click_cnt", in.ClickCnt},
			{"show_cnt", in.ShowCnt},
			{"convert_cnt", in.ConvertCnt},
			{"click_rate", in.ClickRate},
			{"cpm_platform", in.CpmPlatform},
			{"dy_follow", in.DyFollow},
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

func (qair *qianchuanAdvertiserInfoRepo) Delete(ctx context.Context, day string, advertiserIds []uint64) error {
	collection := qair.data.mdb.Database(qair.data.conf.Mongo.Dbname).Collection("qianchuan_advertiser_info_" + day)

	var aadvertiserIds bson.A

	for _, advertiserId := range advertiserIds {
		aadvertiserIds = append(aadvertiserIds, advertiserId)
	}

	if _, err := collection.DeleteMany(ctx, bson.D{{"advertiser_id", bson.D{{"$nin", aadvertiserIds}}}}); err != nil {
		return err
	}

	return nil
}

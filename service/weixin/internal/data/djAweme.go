package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"weixin/internal/biz"
	"weixin/internal/domain"
)

// 星达当家mcn机构达人表
type DjAweme struct {
	UserName      string `json:"user_name" bson:"user_name"`
	Avatar        string `json:"avatar" bson:"avatar"`
	AwemeId       string `json:"aweme_id" bson:"aweme_id"`
	HotsoonId     string `json:"hotsoon_id" bson:"hotsoon_id"`
	FansCount     uint64 `json:"fans_count" bson:"fans_count"`
	Ratio         string `json:"ratio" bson:"ratio"`
	Account       string `json:"account" bson:"account"`
	BindStartTime string `json:"bind_start_time" bson:"bind_start_time"`
	BindEndTime   string `json:"bind_end_time" bson:"bind_end_time"`
}

type djAwemeRepo struct {
	data *Data
	log  *log.Helper
}

func (dj *DjAweme) ToDomain() *domain.DjAweme {
	return &domain.DjAweme{
		UserName:      dj.UserName,
		Avatar:        dj.Avatar,
		AwemeId:       dj.AwemeId,
		HotsoonId:     dj.HotsoonId,
		FansCount:     dj.FansCount,
		Ratio:         dj.Ratio,
		Account:       dj.Account,
		BindStartTime: dj.BindStartTime,
		BindEndTime:   dj.BindEndTime,
	}
}

func NewDjAwemeRepo(data *Data, logger log.Logger) biz.DjAwemeRepo {
	return &djAwemeRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (dar *djAwemeRepo) Get(ctx context.Context, accountId, accounts string) (*domain.DjAweme, error) {
	var djAweme DjAweme
	var aaccounts bson.A

	collection := dar.data.mdb.Database(dar.data.conf.Mongo.Dbname).Collection("dj_aweme")

	saccounts := strings.Split(accounts, ",")

	for _, saccount := range saccounts {
		aaccounts = append(aaccounts, saccount)
	}

	if err := collection.FindOne(ctx, bson.M{
		"$and": []bson.M{
			bson.M{"account": bson.D{{"$in", aaccounts}}},
			bson.M{"aweme_id": accountId},
		},
	}).Decode(&djAweme); err != nil {
		return nil, err
	}

	return djAweme.ToDomain(), nil
}

func (dar *djAwemeRepo) List(ctx context.Context, ratio, accounts string) ([]*domain.DjAweme, error) {
	list := make([]*domain.DjAweme, 0)

	var aaccounts bson.A

	collection := dar.data.mdb.Database(dar.data.conf.Mongo.Dbname).Collection("dj_aweme")

	saccounts := strings.Split(accounts, ",")

	for _, saccount := range saccounts {
		aaccounts = append(aaccounts, saccount)
	}

	matchStage := bson.D{
		{"$match", bson.M{"ratio": bson.D{{"$eq", ratio}}, "account": bson.D{{"$in", aaccounts}}}},
	}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var djAwemes []DjAweme

	err = cursor.All(ctx, &djAwemes)

	if err != nil {
		return nil, err
	}

	for _, djAweme := range djAwemes {
		list = append(list, djAweme.ToDomain())
	}

	return list, nil
}

func (dar *djAwemeRepo) SaveIndex(ctx context.Context) {
	collection := dar.data.mdb.Database(dar.data.conf.Mongo.Dbname).Collection("dj_aweme")

	isNotExistIndex := true

	if indexSpecifications, err := collection.Indexes().ListSpecifications(ctx); err == nil {
		for _, indexSpecification := range indexSpecifications {
			if indexSpecification.Name == "ratio_-1" {
				isNotExistIndex = false
			}
		}

		if isNotExistIndex {
			collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{Key: "ratio", Value: -1},
				},
			})
		}
	}
}

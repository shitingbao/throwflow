package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// 精选联盟平台调用日志表
type JinritemaiApiLog struct {
	ClientKey   string    `json:"client_key" bson:"client_key"`
	OpenId      string    `json:"open_id" bson:"open_id"`
	AccessToken string    `json:"access_token" bson:"access_token"`
	Content     string    `json:"content" bson:"content"`
	CreateTime  time.Time `json:"create_time" bson:"create_time"`
}

type jinritemaiApiLogRepo struct {
	data *Data
	log  *log.Helper
}

func NewJinritemaiApiLogRepo(data *Data, logger log.Logger) biz.JinritemaiApiLogRepo {
	return &jinritemaiApiLogRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
func (jalr *jinritemaiApiLogRepo) Save(ctx context.Context, in *domain.JinritemaiApiLog) error {
	collection := jalr.data.mdb.Database(jalr.data.conf.Mongo.Dbname).Collection("jinritemai_api_log")

	if _, err := collection.InsertOne(ctx, &JinritemaiApiLog{
		ClientKey:   in.ClientKey,
		OpenId:      in.OpenId,
		AccessToken: in.AccessToken,
		Content:     in.Content,
		CreateTime:  in.CreateTime,
	}); err != nil {
		return err
	}

	return nil
}

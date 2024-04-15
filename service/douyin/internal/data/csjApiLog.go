package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// 抖客穿山甲平台调用日志表
type CsjApiLog struct {
	Content    string    `json:"content" bson:"content"`
	CreateTime time.Time `json:"create_time" bson:"create_time"`
}

type csjApiLogRepo struct {
	data *Data
	log  *log.Helper
}

func NewCsjApiLogRepo(data *Data, logger log.Logger) biz.CsjApiLogRepo {
	return &csjApiLogRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
func (calr *csjApiLogRepo) Save(ctx context.Context, in *domain.CsjApiLog) error {
	collection := calr.data.mdb.Database(calr.data.conf.Mongo.Dbname).Collection("csj_api_log")

	if _, err := collection.InsertOne(ctx, &CsjApiLog{
		Content:    in.Content,
		CreateTime: in.CreateTime,
	}); err != nil {
		return err
	}

	return nil
}

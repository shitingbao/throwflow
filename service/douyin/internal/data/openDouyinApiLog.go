package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// 抖音开放平台调用日志表
type OpenDouyinApiLog struct {
	ClientKey   string    `json:"client_key" bson:"client_key"`
	OpenId      string    `json:"open_id" bson:"open_id"`
	AccessToken string    `json:"access_token" bson:"access_token"`
	Content     string    `json:"content" bson:"content"`
	CreateTime  time.Time `json:"create_time" bson:"create_time"`
}

type openDouyinApiLogRepo struct {
	data *Data
	log  *log.Helper
}

func NewOpenDouyinApiLogRepo(data *Data, logger log.Logger) biz.OpenDouyinApiLogRepo {
	return &openDouyinApiLogRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
func (odalr *openDouyinApiLogRepo) Save(ctx context.Context, in *domain.OpenDouyinApiLog) error {
	collection := odalr.data.mdb.Database(odalr.data.conf.Mongo.Dbname).Collection("open_douyin_api_log")

	if _, err := collection.InsertOne(ctx, &OpenDouyinApiLog{
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

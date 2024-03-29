package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// 巨量引擎调用日志表
type OceanengineApiLog struct {
	CompanyId    uint64    `json:"company_id" bson:"company_id"`
	AccountId    uint64    `json:"account_id" bson:"account_id"`
	AppId        string    `json:"app_id" bson:"app_id"`
	AdvertiserId uint64    `json:"advertiser_id" bson:"advertiser_id"`
	CampaignId   uint64    `json:"campaign_id" bson:"campaign_id"`
	AdId         uint64    `json:"ad_id" bson:"ad_id"`
	AccessToken  string    `json:"access_token" bson:"access_token"`
	Content      string    `json:"content" bson:"content"`
	CreateTime   time.Time `json:"create_time" bson:"create_time"`
}

type oceanengineApiLogRepo struct {
	data *Data
	log  *log.Helper
}

func NewOceanengineApiLogRepo(data *Data, logger log.Logger) biz.OceanengineApiLogRepo {
	return &oceanengineApiLogRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
func (oalr *oceanengineApiLogRepo) Save(ctx context.Context, in *domain.OceanengineApiLog) error {
	collection := oalr.data.mdb.Database(oalr.data.conf.Mongo.Dbname).Collection("oceanengine_api_log")

	if _, err := collection.InsertOne(ctx, &OceanengineApiLog{
		CompanyId:    in.CompanyId,
		AccountId:    in.AccountId,
		AppId:        in.AppId,
		AdvertiserId: in.AdvertiserId,
		CampaignId:   in.CampaignId,
		AdId:         in.AdId,
		AccessToken:  in.AccessToken,
		Content:      in.Content,
		CreateTime:   in.CreateTime,
	}); err != nil {
		return err
	}

	return nil
}

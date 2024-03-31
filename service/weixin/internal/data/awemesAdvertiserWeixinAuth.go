package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"weixin/internal/biz"
	"weixin/internal/domain"
)

type AwemesAdvertiserWeixinAuth struct {
	ClientKey       string `json:"client_key" bson:"client_key"`
	OpenId          string `json:"open_id" bson:"open_id"`
	CooperativeCode string `json:"cooperative_code" bson:"cooperative_code"`
	AuthStatus      int32  `json:"auth_status" bson:"auth_status"`
}

type awemesAdvertiserWeixinAuthRepo struct {
	data *Data
	log  *log.Helper
}

func (aawa *AwemesAdvertiserWeixinAuth) ToDomain() *domain.AwemesAdvertiserWeixinAuth {
	return &domain.AwemesAdvertiserWeixinAuth{
		ClientKey:       aawa.ClientKey,
		OpenId:          aawa.OpenId,
		CooperativeCode: aawa.CooperativeCode,
		AuthStatus:      aawa.AuthStatus,
	}
}

func NewAwemesAdvertiserWeixinAuthRepo(data *Data, logger log.Logger) biz.AwemesAdvertiserWeixinAuthRepo {
	return &awemesAdvertiserWeixinAuthRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (aawar *awemesAdvertiserWeixinAuthRepo) Upsert(ctx context.Context, in *domain.AwemesAdvertiserWeixinAuth) error {
	collection := aawar.data.mdb.Database(aawar.data.conf.Mongo.Dbname).Collection("awemes_advertiser_weixin_auth")

	if _, err := collection.UpdateOne(ctx, bson.D{
		{"client_key", in.ClientKey},
		{"open_id", in.OpenId},
	}, bson.D{
		{"$set", bson.D{
			{"cooperative_code", in.CooperativeCode},
			{"auth_status", in.AuthStatus},
		}},
	}, options.Update().SetUpsert(true)); err != nil {
		return err
	}

	return nil
}

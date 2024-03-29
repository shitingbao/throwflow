package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// 账户钱包信息表
type QianchuanWallet struct {
	AdvertiserId                    uint64                            `json:"advertiser_id" bson:"advertiser_id"`
	TotalBalanceAbs                 float64                           `json:"total_balance_abs" bson:"total_balance_abs"`
	GrantBalance                    float64                           `json:"grant_balance" bson:"grant_balance"`
	UnionValidGrantBalance          float64                           `json:"union_valid_grant_balance" bson:"union_valid_grant_balance"`
	SearchValidGrantBalance         float64                           `json:"search_valid_grant_balance" bson:"search_valid_grant_balance"`
	CommonValidGrantBalance         float64                           `json:"common_valid_grant_balance" bson:"common_valid_grant_balance"`
	DefaultValidGrantBalance        float64                           `json:"default_valid_grant_balance" bson:"default_valid_grant_balance"`
	GeneralTotalBalance             float64                           `json:"general_total_balance" bson:"general_total_balance"`
	GeneralBalanceValid             float64                           `json:"general_balance_valid" bson:"general_balance_valid"`
	GeneralBalanceValidNonGrant     float64                           `json:"general_balance_valid_non_grant" bson:"general_balance_valid_non_grant"`
	GeneralBalanceValidGrantUnion   float64                           `json:"general_balance_valid_grant_union" bson:"general_balance_valid_grant_union"`
	GeneralBalanceValidGrantSearch  float64                           `json:"general_balance_valid_grant_search" bson:"general_balance_valid_grant_search"`
	GeneralBalanceValidGrantCommon  float64                           `json:"general_balance_valid_grant_common" bson:"general_balance_valid_grant_common"`
	GeneralBalanceValidGrantDefault float64                           `json:"general_balance_valid_grant_default" bson:"general_balance_valid_grant_default"`
	GeneralBalanceInvalid           float64                           `json:"general_balance_invalid" bson:"general_balance_invalid"`
	GeneralBalanceInvalidOrder      float64                           `json:"general_balance_invalid_order" bson:"general_balance_invalid_order"`
	GeneralBalanceInvalidFrozen     float64                           `json:"general_balance_invalid_frozen" bson:"general_balance_invalid_frozen"`
	BrandBalance                    float64                           `json:"brand_balance" bson:"brand_balance"`
	BrandBalanceValid               float64                           `json:"brand_balance_valid" bson:"brand_balance_valid"`
	BrandBalanceValidNonGrant       float64                           `json:"brand_balance_valid_non_grant" bson:"brand_balance_valid_non_grant"`
	BrandBalanceValidGrant          float64                           `json:"brand_balance_valid_grant" bson:"brand_balance_valid_grant"`
	BrandBalanceInvalid             float64                           `json:"brand_balance_invalid" bson:"brand_balance_invalid"`
	BrandBalanceInvalidFrozen       float64                           `json:"brand_balance_invalid_frozen" bson:"brand_balance_invalid_frozen"`
	DeductionCouponBalance          float64                           `json:"deduction_coupon_balance" bson:"deduction_coupon_balance"`
	DeductionCouponBalanceAll       float64                           `json:"deduction_coupon_balance_all" bson:"deduction_coupon_balance_all"`
	DeductionCouponBalanceOther     float64                           `json:"deduction_coupon_balance_other" bson:"deduction_coupon_balance_other"`
	DeductionCouponBalanceSelf      float64                           `json:"deduction_coupon_balance_self" bson:"deduction_coupon_balance_self"`
	GrantExpiring                   float64                           `json:"grant_expiring" bson:"grant_expiring"`
	ShareBalance                    float64                           `json:"share_balance" bson:"share_balance"`
	ShareBalanceValidGrantUnion     float64                           `json:"share_balance_valid_grant_union" bson:"share_balance_valid_grant_union"`
	ShareBalanceValidGrantSearch    float64                           `json:"share_balance_valid_grant_search" bson:"share_balance_valid_grant_search"`
	ShareBalanceValidGrantCommon    float64                           `json:"share_balance_valid_grant_common" bson:"share_balance_valid_grant_common"`
	ShareBalanceValidGrantDefault   float64                           `json:"share_balance_valid_grant_default" bson:"share_balance_valid_grant_default"`
	ShareBalanceValid               float64                           `json:"share_balance_valid" bson:"share_balance_valid"`
	ShareBalanceExpiring            float64                           `json:"share_balance_expiring" bson:"share_balance_expiring"`
	ShareExpiringDetailList         []*domain.ShareExpiringDetailList `json:"share_expiring_detail_list" bson:"share_expiring_detail_list"`
	CreateTime                      time.Time                         `json:"create_time" bson:"create_time"`
	UpdateTime                      time.Time                         `json:"update_time" bson:"update_time"`
}

type qianchuanWalletRepo struct {
	data *Data
	log  *log.Helper
}

func (qw *QianchuanWallet) ToDomain() *domain.QianchuanWallet {
	return &domain.QianchuanWallet{
		AdvertiserId:                    qw.AdvertiserId,
		TotalBalanceAbs:                 qw.TotalBalanceAbs,
		GrantBalance:                    qw.GrantBalance,
		UnionValidGrantBalance:          qw.UnionValidGrantBalance,
		SearchValidGrantBalance:         qw.SearchValidGrantBalance,
		CommonValidGrantBalance:         qw.CommonValidGrantBalance,
		DefaultValidGrantBalance:        qw.DefaultValidGrantBalance,
		GeneralTotalBalance:             qw.GeneralTotalBalance,
		GeneralBalanceValid:             qw.GeneralBalanceValid,
		GeneralBalanceValidNonGrant:     qw.GeneralBalanceValidNonGrant,
		GeneralBalanceValidGrantUnion:   qw.GeneralBalanceValidGrantUnion,
		GeneralBalanceValidGrantSearch:  qw.GeneralBalanceValidGrantSearch,
		GeneralBalanceValidGrantCommon:  qw.GeneralBalanceValidGrantCommon,
		GeneralBalanceValidGrantDefault: qw.GeneralBalanceValidGrantDefault,
		GeneralBalanceInvalid:           qw.GeneralBalanceInvalid,
		GeneralBalanceInvalidOrder:      qw.GeneralBalanceInvalidOrder,
		GeneralBalanceInvalidFrozen:     qw.GeneralBalanceInvalidFrozen,
		BrandBalance:                    qw.BrandBalance,
		BrandBalanceValid:               qw.BrandBalanceValid,
		BrandBalanceValidNonGrant:       qw.BrandBalanceValidNonGrant,
		BrandBalanceValidGrant:          qw.BrandBalanceValidGrant,
		BrandBalanceInvalid:             qw.BrandBalanceInvalid,
		BrandBalanceInvalidFrozen:       qw.BrandBalanceInvalidFrozen,
		DeductionCouponBalance:          qw.DeductionCouponBalance,
		DeductionCouponBalanceAll:       qw.DeductionCouponBalanceAll,
		DeductionCouponBalanceOther:     qw.DeductionCouponBalanceOther,
		DeductionCouponBalanceSelf:      qw.DeductionCouponBalanceSelf,
		GrantExpiring:                   qw.GrantExpiring,
		ShareBalance:                    qw.ShareBalance,
		ShareBalanceValidGrantUnion:     qw.ShareBalanceValidGrantUnion,
		ShareBalanceValidGrantSearch:    qw.ShareBalanceValidGrantSearch,
		ShareBalanceValidGrantCommon:    qw.ShareBalanceValidGrantCommon,
		ShareBalanceValidGrantDefault:   qw.ShareBalanceValidGrantDefault,
		ShareBalanceValid:               qw.ShareBalanceValid,
		ShareBalanceExpiring:            qw.ShareBalanceExpiring,
		ShareExpiringDetailList:         qw.ShareExpiringDetailList,
		CreateTime:                      qw.CreateTime,
		UpdateTime:                      qw.UpdateTime,
	}
}

func NewQianchuanWalletRepo(data *Data, logger log.Logger) biz.QianchuanWalletRepo {
	return &qianchuanWalletRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (qwr *qianchuanWalletRepo) List(ctx context.Context, day string) ([]*domain.QianchuanWallet, error) {
	list := make([]*domain.QianchuanWallet, 0)

	collection := qwr.data.mdb.Database(qwr.data.conf.Mongo.Dbname).Collection("qianchuan_wallet_" + day)

	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var qianchuanWallets []*QianchuanWallet

	err = cursor.All(ctx, &qianchuanWallets)

	if err != nil {
		return nil, err
	}

	for _, qianchuanWallet := range qianchuanWallets {
		list = append(list, qianchuanWallet.ToDomain())
	}

	return list, nil
}

func (qwr *qianchuanWalletRepo) SaveIndex(ctx context.Context, day string) {
	collection := qwr.data.mdb.Database(qwr.data.conf.Mongo.Dbname).Collection("qianchuan_wallet_" + day)

	isNotExistIndex := true

	if indexSpecifications, err := collection.Indexes().ListSpecifications(ctx); err == nil {
		for _, indexSpecification := range indexSpecifications {
			if indexSpecification.Name == "advertiser_id_-1" {
				isNotExistIndex = false
			}
		}

		if isNotExistIndex {
			collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{Key: "advertiser_id", Value: -1},
				},
			})
		}
	}
}

func (qwr *qianchuanWalletRepo) Upsert(ctx context.Context, day string, in *domain.QianchuanWallet) error {
	collection := qwr.data.mdb.Database(qwr.data.conf.Mongo.Dbname).Collection("qianchuan_wallet_" + day)

	if _, err := collection.UpdateOne(ctx, bson.D{
		{"advertiser_id", in.AdvertiserId},
	}, bson.D{
		{"$set", bson.D{
			{"total_balance_abs", in.TotalBalanceAbs},
			{"grant_balance", in.GrantBalance},
			{"union_valid_grant_balance", in.UnionValidGrantBalance},
			{"search_valid_grant_balance", in.SearchValidGrantBalance},
			{"common_valid_grant_balance", in.CommonValidGrantBalance},
			{"default_valid_grant_balance", in.DefaultValidGrantBalance},
			{"general_total_balance", in.GeneralTotalBalance},
			{"general_balance_valid", in.GeneralBalanceValid},
			{"general_balance_valid_non_grant", in.GeneralBalanceValidNonGrant},
			{"general_balance_valid_grant_union", in.GeneralBalanceValidGrantUnion},
			{"general_balance_valid_grant_search", in.GeneralBalanceValidGrantSearch},
			{"general_balance_valid_grant_common", in.GeneralBalanceValidGrantCommon},
			{"general_balance_valid_grant_default", in.GeneralBalanceValidGrantDefault},
			{"general_balance_invalid", in.GeneralBalanceInvalid},
			{"general_balance_invalid_order", in.GeneralBalanceInvalidOrder},
			{"general_balance_invalid_frozen", in.GeneralBalanceInvalidFrozen},
			{"brand_balance_valid_grant", in.BrandBalanceValidGrant},
			{"brand_balance_invalid", in.BrandBalanceInvalid},
			{"brand_balance_invalid_frozen", in.BrandBalanceInvalidFrozen},
			{"deduction_coupon_balance", in.DeductionCouponBalance},
			{"deduction_coupon_balance_all", in.DeductionCouponBalanceAll},
			{"deduction_coupon_balance_other", in.DeductionCouponBalanceOther},
			{"deduction_coupon_balance_self", in.DeductionCouponBalanceSelf},
			{"grant_expiring", in.GrantExpiring},
			{"share_balance", in.ShareBalance},
			{"share_balance_valid_grant_union", in.ShareBalanceValidGrantUnion},
			{"share_balance_valid_grant_search", in.ShareBalanceValidGrantSearch},
			{"share_balance_valid_grant_common", in.ShareBalanceValidGrantCommon},
			{"share_balance_valid_grant_default", in.ShareBalanceValidGrantDefault},
			{"share_balance_valid", in.ShareBalanceValid},
			{"share_balance_expiring", in.ShareBalanceExpiring},
			{"share_expiring_detail_list", in.ShareExpiringDetailList},
			{"create_time", in.CreateTime},
			{"update_time", in.UpdateTime},
		}},
	}, options.Update().SetUpsert(true)); err != nil {
		return err
	}

	return nil
}

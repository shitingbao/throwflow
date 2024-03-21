package data

import (
	"company/internal/biz"
	"company/internal/domain"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
)

type qianchuanAdAdvertiserRepo struct {
	data *Data
	log  *log.Helper
}

func NewQianchuanAdAdvertiserRepo(data *Data, logger log.Logger) biz.QianchuanAdAdvertiserRepo {
	return &qianchuanAdAdvertiserRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (qar *qianchuanAdAdvertiserRepo) List(ctx context.Context, companyId uint64, updateDay string) ([]*domain.QianchuanAdAdvertiser, error) {
	var qianchuanAdAdvertisers []struct {
		AdvertiserId   int64
		StatCost       float64
		PayOrderAmount float64
		MarketingGoal  string
	}
	var list []*domain.QianchuanAdAdvertiser

	sql := fmt.Sprintf("SELECT advertiser_id as AdvertiserId, marketing_goal as MarketingGoal, sum(stat_cost) as StatCost, sum(pay_order_amount) as PayOrderAmount FROM qc_consume_advert_daily WHERE day_id='%s' AND operator_uid = %d GROUP BY advertiser_id, marketing_goal", updateDay, companyId)

	if err := qar.data.cdb.Select(ctx, &qianchuanAdAdvertisers, sql); err != nil {
		return nil, err
	}

	for _, qianchuanAdAdvertiser := range qianchuanAdAdvertisers {
		list = append(list, &domain.QianchuanAdAdvertiser{
			AdvertiserId:   uint64(qianchuanAdAdvertiser.AdvertiserId),
			StatCost:       float32(qianchuanAdAdvertiser.StatCost),
			PayOrderAmount: float32(qianchuanAdAdvertiser.PayOrderAmount),
			MarketingGoal:  qianchuanAdAdvertiser.MarketingGoal,
		})
	}

	return list, nil
}

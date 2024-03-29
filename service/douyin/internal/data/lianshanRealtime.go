package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// 连山云RDS数据表
type LianshanRealtime struct {
	Id                         uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:主键id"`
	AdvertiserId               int64     `gorm:"column:advertiser_id;type:bigint(20);default:null;comment:广告主id"`
	AwemeId                    int64     `gorm:"column:aweme_id;type:bigint(20);default:null;comment:关联的抖音号"`
	AdId                       int64     `gorm:"column:ad_id;type:bigint(20);default:null;comment:广告计划id"`
	CreativeId                 int64     `gorm:"column:creative_id;type:bigint(20);default:null;comment:广告创意id"`
	MarketingGoal              int64     `gorm:"column:marketing_goal;type:bigint(20);default:null;comment:营销目标"`
	OrderPlatform              int64     `gorm:"column:order_platform;type:bigint(20);default:null;comment:下单平台"`
	MarketingScene             int64     `gorm:"column:marketing_scene;type:bigint(20);default:null;comment:营销场景"`
	PromotionWay               int64     `gorm:"column:promotion_way;type:bigint(20);default:null;comment:推广方式"`
	CreativeMaterialMode       int64     `gorm:"column:creative_material_mode;type:bigint(20);default:null;comment:创意类型"`
	ImageMode                  int64     `gorm:"column:image_mode;type:bigint(20);default:null;comment:素材样式"`
	SmartBidType               int64     `gorm:"column:smart_bid_type;type:bigint(20);default:null;comment:投放场景"`
	PricingCategory            int64     `gorm:"column:pricing_category;type:bigint(20);default:null;comment:广告类型"`
	StatCost                   float64   `gorm:"column:stat_cost;type:double;default:null;comment:热度"`
	ShowCnt                    int64     `gorm:"column:show_cnt;type:bigint(20);default:null;comment:展示次数"`
	ClickCnt                   int64     `gorm:"column:click_cnt;type:bigint(20);default:null;comment:点击次数"`
	ConvertCnt                 int64     `gorm:"column:convert_cnt;type:bigint(20);default:null;comment:转化数"`
	PayOrderCount              int64     `gorm:"column:pay_order_count;type:bigint(20);default:null;comment:直接成交订单数"`
	PayOrderAmount             float64   `gorm:"column:pay_order_amount;type:double;default:null;comment:直接成交金额"`
	CreateOrderCount           int64     `gorm:"column:create_order_count;type:bigint(20);default:null;comment:直接下单订单数"`
	CreateOrderAmount          float64   `gorm:"column:create_order_amount;type:double;default:null;comment:直接下单金额"`
	PrepayOrderCount           int64     `gorm:"column:prepay_order_count;type:bigint(20);default:null;comment:直接预售订单数"`
	PrepayOrderAmount          float64   `gorm:"column:prepay_order_amount;type:double;default:null;comment:直接预售金额"`
	DyFollow                   int64     `gorm:"column:dy_follow;type:bigint(20);default:null;comment:新增粉丝数"`
	LubanLiveEnterCnt          int64     `gorm:"column:luban_live_enter_cnt;type:bigint(20);default:null;comment:直播间观看人数"`
	LiveWatchOneMinuteCount    int64     `gorm:"column:live_watch_one_minute_count;type:bigint(20);default:null;comment:直播间超过1分钟观看人数"`
	LiveFansClubJoinCnt        int64     `gorm:"column:live_fans_club_join_cnt;type:bigint(20);default:null;comment:直播间新加团人次"`
	LubanLiveSlidecartClickCnt int64     `gorm:"column:luban_live_slidecart_click_cnt;type:bigint(20);default:null;comment:直播间查看购物车次数"`
	LubanLiveClickProductCnt   int64     `gorm:"column:luban_live_click_product_cnt;type:bigint(20);default:null;comment:直播间商品点击次数"`
	LubanLiveCommentCnt        int64     `gorm:"column:luban_live_comment_cnt;type:bigint(20);default:null;comment:直播间评论次数"`
	LubanLiveShareCnt          int64     `gorm:"column:luban_live_share_cnt;type:bigint(20);default:null;comment:直播间分享次数"`
	LubanLiveGiftCnt           int64     `gorm:"column:luban_live_gift_cnt;type:bigint(20);default:null;comment:直播间打赏次数"`
	LubanLiveGiftAmount        int64     `gorm:"column:luban_live_gift_amount;type:bigint(20);default:null;comment:直播间音浪收入"`
	DyShare                    int64     `gorm:"column:dy_share;type:bigint(20);default:null;comment:分享次数"`
	DyComment                  int64     `gorm:"column:dy_comment;type:bigint(20);default:null;comment:评论次数"`
	DyLike                     int64     `gorm:"column:dy_like;type:bigint(20);default:null;comment:点赞次数"`
	TotalPlay                  int64     `gorm:"column:total_play;type:bigint(20);default:null;comment:播放数"`
	PlayDuration3s             int64     `gorm:"column:play_duration_3s;type:bigint(20);default:null;comment:3S播放数"`
	Play25FeedBreak            int64     `gorm:"column:play_25_feed_break;type:bigint(20);default:null;comment:25%进度播放数"`
	Play50FeedBreak            int64     `gorm:"column:play_50_feed_break;type:bigint(20);default:null;comment:50%进度播放数"`
	Play75FeedBreak            int64     `gorm:"column:play_75_feed_break;type:bigint(20);default:null;comment:75%进度播放数"`
	PlayOver                   int64     `gorm:"column:play_over;type:bigint(20);default:null;comment:播放完成数"`
	StatTime                   time.Time `gorm:"column:stat_time;type:datetime;default:null;comment:数据产生时间"`
	Uid                        string    `gorm:"column:uid;type:varchar(128);default:null;comment:对应docId"`
	CreateTime                 time.Time `gorm:"column:create_time;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:记录创建时间"`
	UpdateTime                 time.Time `gorm:"column:update_time;type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:记录修改时间"`
}

type lianshanRealtimeRepo struct {
	data *Data
	log  *log.Helper
}

func (lr *LianshanRealtime) ToDomain() *domain.LianshanRealtime {
	return &domain.LianshanRealtime{
		Id:                         lr.Id,
		AdvertiserId:               lr.AdvertiserId,
		AwemeId:                    lr.AwemeId,
		AdId:                       lr.AdId,
		CreativeId:                 lr.CreativeId,
		MarketingGoal:              lr.MarketingGoal,
		OrderPlatform:              lr.OrderPlatform,
		MarketingScene:             lr.MarketingScene,
		PromotionWay:               lr.PromotionWay,
		CreativeMaterialMode:       lr.CreativeMaterialMode,
		ImageMode:                  lr.ImageMode,
		SmartBidType:               lr.SmartBidType,
		PricingCategory:            lr.PricingCategory,
		StatCost:                   lr.StatCost,
		ShowCnt:                    lr.ShowCnt,
		ClickCnt:                   lr.ClickCnt,
		ConvertCnt:                 lr.ConvertCnt,
		PayOrderCount:              lr.PayOrderCount,
		PayOrderAmount:             lr.PayOrderAmount,
		CreateOrderCount:           lr.CreateOrderCount,
		CreateOrderAmount:          lr.CreateOrderAmount,
		PrepayOrderCount:           lr.PrepayOrderCount,
		PrepayOrderAmount:          lr.PrepayOrderAmount,
		DyFollow:                   lr.DyFollow,
		LubanLiveEnterCnt:          lr.LubanLiveEnterCnt,
		LiveWatchOneMinuteCount:    lr.LiveWatchOneMinuteCount,
		LiveFansClubJoinCnt:        lr.LiveFansClubJoinCnt,
		LubanLiveSlidecartClickCnt: lr.LubanLiveSlidecartClickCnt,
		LubanLiveClickProductCnt:   lr.LubanLiveClickProductCnt,
		LubanLiveCommentCnt:        lr.LubanLiveCommentCnt,
		LubanLiveShareCnt:          lr.LubanLiveShareCnt,
		LubanLiveGiftCnt:           lr.LubanLiveGiftCnt,
		LubanLiveGiftAmount:        lr.LubanLiveGiftAmount,
		DyShare:                    lr.DyShare,
		DyComment:                  lr.DyComment,
		DyLike:                     lr.DyLike,
		TotalPlay:                  lr.TotalPlay,
		PlayDuration3s:             lr.PlayDuration3s,
		Play25FeedBreak:            lr.Play25FeedBreak,
		Play50FeedBreak:            lr.Play50FeedBreak,
		Play75FeedBreak:            lr.Play75FeedBreak,
		PlayOver:                   lr.PlayOver,
		StatTime:                   lr.StatTime,
		Uid:                        lr.Uid,
		CreateTime:                 lr.CreateTime,
		UpdateTime:                 lr.UpdateTime,
	}
}

func NewLianshanRealtimeRepo(data *Data, logger log.Logger) biz.LianshanRealtimeRepo {
	return &lianshanRealtimeRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (lrr *lianshanRealtimeRepo) List(ctx context.Context, appId, monday, day, groupName string) ([]*domain.LianshanRealtime, error) {
	var lianshanRealtimes []LianshanRealtime
	list := make([]*domain.LianshanRealtime, 0)

	for _, lianshandb := range lrr.data.lianshandbs {
		if lianshandb.appId == appId {
			if result := lianshandb.db.Table("lianshan_realtime_"+monday).
				WithContext(ctx).
				Select("ad_id, advertiser_id, aweme_id, marketing_goal, SUM(stat_cost) as stat_cost, SUM(show_cnt) show_cnt, SUM(click_cnt) click_cnt, SUM(pay_order_count) pay_order_count, SUM(create_order_amount) create_order_amount, SUM(create_order_count) create_order_count, SUM(pay_order_amount) pay_order_amount, SUM(dy_follow)dy_follow, SUM(convert_cnt) convert_cnt").
				Where("order_platform = 1").
				Where("marketing_scene = 1").
				Where("stat_time >= ?", day+" 00:00:00").
				Where("stat_time <= ?", day+" 23:59:59").
				Group(groupName).
				Find(&lianshanRealtimes); result.Error != nil {
				return nil, result.Error
			}

			break
		}
	}

	for _, lianshanRealtime := range lianshanRealtimes {
		list = append(list, lianshanRealtime.ToDomain())
	}

	return list, nil
}

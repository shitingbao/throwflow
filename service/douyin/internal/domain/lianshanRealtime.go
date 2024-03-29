package domain

import (
	"time"
)

type LianshanRealtime struct {
	Id                         uint64
	AdvertiserId               int64
	AwemeId                    int64
	AdId                       int64
	CreativeId                 int64
	MarketingGoal              int64
	OrderPlatform              int64
	MarketingScene             int64
	PromotionWay               int64
	CreativeMaterialMode       int64
	ImageMode                  int64
	SmartBidType               int64
	PricingCategory            int64
	StatCost                   float64
	ShowCnt                    int64
	ClickCnt                   int64
	ConvertCnt                 int64
	PayOrderCount              int64
	PayOrderAmount             float64
	CreateOrderCount           int64
	CreateOrderAmount          float64
	PrepayOrderCount           int64
	PrepayOrderAmount          float64
	DyFollow                   int64
	LubanLiveEnterCnt          int64
	LiveWatchOneMinuteCount    int64
	LiveFansClubJoinCnt        int64
	LubanLiveSlidecartClickCnt int64
	LubanLiveClickProductCnt   int64
	LubanLiveCommentCnt        int64
	LubanLiveShareCnt          int64
	LubanLiveGiftCnt           int64
	LubanLiveGiftAmount        int64
	DyShare                    int64
	DyComment                  int64
	DyLike                     int64
	TotalPlay                  int64
	PlayDuration3s             int64
	Play25FeedBreak            int64
	Play50FeedBreak            int64
	Play75FeedBreak            int64
	PlayOver                   int64
	StatTime                   time.Time
	Uid                        string
	CreateTime                 time.Time
	UpdateTime                 time.Time
}

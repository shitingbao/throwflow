package domain

import (
	"context"
	"time"
)

type VideoInfo struct {
	AwemeItemId    uint64 `json:"aweme_item_id" bson:"aweme_item_id"`
	AwemeItemTitle string `json:"aweme_item_title" bson:"aweme_item_title"`
	AwemeItemCover string `json:"aweme_item_cover" bson:"aweme_item_cover"`
	ItemType       uint8  `json:"item_type" bson:"item_type"`
}

type RoomInfo struct {
	RoomId     uint64 `json:"room_id" bson:"room_id"`
	RoomTitle  string `json:"room_title" bson:"room_title"`
	RoomCover  string `json:"room_cover" bson:"room_cover"`
	RoomStatus string `json:"room_status" bson:"room_status"`
}

type QianchuanAwemeOrderInfo struct {
	OrderId              uint64
	AdvertiserId         uint64
	AdId                 uint64
	MarketingGoal        string
	Status               string
	OrderCreateTime      time.Time
	AwemeInfo            *AwemeInfo
	VideoInfo            *VideoInfo
	RoomInfo             *RoomInfo
	DeliverySetting      *DeliverySetting
	FailList             []uint64
	PayOrderAmount       float64
	StatCost             float64
	PrepayAndPayOrderRoi float64
	TotalPlay            float64
	ShowCnt              int64
	Ctr                  int64
	ClickCnt             int64
	PayOrderCount        int64
	PrepayOrderCount     int64
	PrepayOrderAmount    int64
	DyFollow             int64
	DyShare              int64
}

func NewQianchuanAwemeOrderInfo(ctx context.Context, orderId, advertiserId, adId uint64, showCnt, ctr, clickCnt, payOrderCount, prepayOrderCount, prepayOrderAmount, dyFollow, dyShare int64, failList []uint64, payOrderAmount, statCost, prepayAndPayOrderRoi, totalPlay float64, marketingGoal, status string, awemeInfo *AwemeInfo, videoInfo *VideoInfo, roomInfo *RoomInfo, deliverySetting *DeliverySetting, orderCreateTime time.Time) *QianchuanAwemeOrderInfo {
	return &QianchuanAwemeOrderInfo{
		OrderId:              orderId,
		AdvertiserId:         advertiserId,
		AdId:                 adId,
		MarketingGoal:        marketingGoal,
		Status:               status,
		OrderCreateTime:      orderCreateTime,
		AwemeInfo:            awemeInfo,
		VideoInfo:            videoInfo,
		RoomInfo:             roomInfo,
		DeliverySetting:      deliverySetting,
		FailList:             failList,
		PayOrderAmount:       payOrderAmount,
		StatCost:             statCost,
		PrepayAndPayOrderRoi: prepayAndPayOrderRoi,
		TotalPlay:            totalPlay,
		ShowCnt:              showCnt,
		Ctr:                  ctr,
		ClickCnt:             clickCnt,
		PayOrderCount:        payOrderCount,
		PrepayOrderCount:     prepayOrderCount,
		PrepayOrderAmount:    prepayOrderAmount,
		DyFollow:             dyFollow,
		DyShare:              dyShare,
	}
}

type AwemeVideoProductQianchuanAwemeOrderInfo struct {
	AwemeId        uint64  `json:"aweme_id" bson:"aweme_id"`
	VideoId        uint64  `json:"video_id" bson:"video_id"`
	ProductId      uint64  `json:"product_id" bson:"product_id"`
	PayOrderAmount float64 `json:"pay_order_amount" bson:"pay_order_amount"`
	StatCost       float64 `json:"stat_cost" bson:"stat_cost"`
}

func (qaoi *QianchuanAwemeOrderInfo) SetOrderId(ctx context.Context, orderId uint64) {
	qaoi.OrderId = orderId
}

func (qaoi *QianchuanAwemeOrderInfo) SetAdvertiserId(ctx context.Context, advertiserId uint64) {
	qaoi.AdvertiserId = advertiserId
}

func (qaoi *QianchuanAwemeOrderInfo) SetAdId(ctx context.Context, adId uint64) {
	qaoi.AdId = adId
}

func (qaoi *QianchuanAwemeOrderInfo) SetMarketingGoal(ctx context.Context, marketingGoal string) {
	qaoi.MarketingGoal = marketingGoal
}

func (qaoi *QianchuanAwemeOrderInfo) SetStatus(ctx context.Context, status string) {
	qaoi.Status = status
}

func (qaoi *QianchuanAwemeOrderInfo) SetOrderCreateTime(ctx context.Context, orderCreateTime time.Time) {
	qaoi.OrderCreateTime = orderCreateTime
}

func (qaoi *QianchuanAwemeOrderInfo) SetAwemeInfo(ctx context.Context, awemeInfo *AwemeInfo) {
	qaoi.AwemeInfo = awemeInfo
}

func (qaoi *QianchuanAwemeOrderInfo) SetVideoInfo(ctx context.Context, videoInfo *VideoInfo) {
	qaoi.VideoInfo = videoInfo
}

func (qaoi *QianchuanAwemeOrderInfo) SetRoomInfo(ctx context.Context, roomInfo *RoomInfo) {
	qaoi.RoomInfo = roomInfo
}

func (qaoi *QianchuanAwemeOrderInfo) SetDeliverySetting(ctx context.Context, deliverySetting *DeliverySetting) {
	qaoi.DeliverySetting = deliverySetting
}

func (qaoi *QianchuanAwemeOrderInfo) SetFailList(ctx context.Context, failList []uint64) {
	qaoi.FailList = failList
}

func (qaoi *QianchuanAwemeOrderInfo) SetPayOrderAmount(ctx context.Context, payOrderAmount float64) {
	qaoi.PayOrderAmount = payOrderAmount
}

func (qaoi *QianchuanAwemeOrderInfo) SetStatCost(ctx context.Context, statCost float64) {
	qaoi.StatCost = statCost
}

func (qaoi *QianchuanAwemeOrderInfo) SetPrepayAndPayOrderRoi(ctx context.Context, prepayAndPayOrderRoi float64) {
	qaoi.PrepayAndPayOrderRoi = prepayAndPayOrderRoi
}

func (qaoi *QianchuanAwemeOrderInfo) SetTotalPlay(ctx context.Context, totalPlay float64) {
	qaoi.TotalPlay = totalPlay
}

func (qaoi *QianchuanAwemeOrderInfo) SetShowCnt(ctx context.Context, showCnt int64) {
	qaoi.ShowCnt = showCnt
}

func (qaoi *QianchuanAwemeOrderInfo) SetCtr(ctx context.Context, ctr int64) {
	qaoi.Ctr = ctr
}

func (qaoi *QianchuanAwemeOrderInfo) SetClickCnt(ctx context.Context, clickCnt int64) {
	qaoi.ClickCnt = clickCnt
}

func (qaoi *QianchuanAwemeOrderInfo) SetPayOrderCount(ctx context.Context, payOrderCount int64) {
	qaoi.PayOrderCount = payOrderCount
}

func (qaoi *QianchuanAwemeOrderInfo) SetPrepayOrderCount(ctx context.Context, prepayOrderCount int64) {
	qaoi.PrepayOrderCount = prepayOrderCount
}

func (qaoi *QianchuanAwemeOrderInfo) SetPrepayOrderAmount(ctx context.Context, prepayOrderAmount int64) {
	qaoi.PrepayOrderAmount = prepayOrderAmount
}

func (qaoi *QianchuanAwemeOrderInfo) SetDyFollow(ctx context.Context, dyFollow int64) {
	qaoi.DyFollow = dyFollow
}

func (qaoi *QianchuanAwemeOrderInfo) SetDyShare(ctx context.Context, dyShare int64) {
	qaoi.DyShare = dyShare
}

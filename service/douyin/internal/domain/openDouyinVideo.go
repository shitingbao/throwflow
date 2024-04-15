package domain

import (
	"context"
)

type VideoStatistics struct {
	CommentCount  int32 `json:"comment_count" bson:"comment_count"`
	DiggCount     int32 `json:"digg_count" bson:"digg_count"`
	DownloadCount int32 `json:"download_count" bson:"download_count"`
	ForwardCount  int32 `json:"forward_count" bson:"forward_count"`
	PlayCount     int32 `json:"play_count" bson:"play_count"`
	ShareCount    int32 `json:"share_count" bson:"share_count"`
}

type SucaiVideoRequest struct {
	ProductId  string   `json:"product_id"`
	OpenId     string   `json:"open_id"`
	CreateTime int64    `json:"create_time"`
	VideoIds   []string `json:"video_list"`
}

func NewSucaiVideoRequest(ctx context.Context, productId, openId string, createTime int64, videoIds []string) *SucaiVideoRequest {
	return &SucaiVideoRequest{
		ProductId:  productId,
		OpenId:     openId,
		CreateTime: createTime,
		VideoIds:   videoIds,
	}
}

type SucaiVideo struct {
	VideoId    string `json:"video_id"`
	ProductId  string `json:"product_id"`
	CreateTime int32  `json:"create_time"`
	OpenId     string `json:"open_id"`
}

type SucaiVideoResult struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data []*SucaiVideo `json:"data"`
}

type OpenDouyinVideo struct {
	ClientKey     string
	OpenId        string
	AwemeId       uint64
	AccountId     string
	Nickname      string
	Avatar        string
	Title         string
	Cover         string
	CreateTime    int64
	IsReviewed    bool
	ItemId        string
	Statistics    VideoStatistics
	IsTop         bool
	MediaType     int
	ShareUrl      string
	VideoId       string
	VideoStatus   int32
	ProductId     string
	ProductName   string
	ProductImg    string
	IsUpdateCover uint8
}

type OpenDouyinVideoList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*OpenDouyinVideo
}

func NewOpenDouyinVideo(ctx context.Context, clientKey, openId, accountId, nickname, avatar, title, cover, itemId, shareUrl, videoId string, awemeId uint64, createTime int64, mediaType int, videoStatus int32, isReviewed, isTop bool, statistics VideoStatistics) *OpenDouyinVideo {
	return &OpenDouyinVideo{
		ClientKey:   clientKey,
		OpenId:      openId,
		AwemeId:     awemeId,
		AccountId:   accountId,
		Nickname:    nickname,
		Avatar:      avatar,
		Title:       title,
		Cover:       cover,
		CreateTime:  createTime,
		IsReviewed:  isReviewed,
		ItemId:      itemId,
		Statistics:  statistics,
		IsTop:       isTop,
		MediaType:   mediaType,
		ShareUrl:    shareUrl,
		VideoId:     videoId,
		VideoStatus: videoStatus,
	}
}

func (odv *OpenDouyinVideo) SetClientKey(ctx context.Context, clientKey string) {
	odv.ClientKey = clientKey
}

func (odv *OpenDouyinVideo) SetOpenId(ctx context.Context, openId string) {
	odv.OpenId = openId
}

func (odv *OpenDouyinVideo) SetAwemeId(ctx context.Context, awemeId uint64) {
	odv.AwemeId = awemeId
}

func (odv *OpenDouyinVideo) SetAccountId(ctx context.Context, accountId string) {
	odv.AccountId = accountId
}

func (odv *OpenDouyinVideo) SetNickname(ctx context.Context, nickname string) {
	odv.Nickname = nickname
}

func (odv *OpenDouyinVideo) SetAvatar(ctx context.Context, avatar string) {
	odv.Avatar = avatar
}

func (odv *OpenDouyinVideo) SetTitle(ctx context.Context, title string) {
	odv.Title = title
}

func (odv *OpenDouyinVideo) SetCover(ctx context.Context, cover string) {
	odv.Cover = cover
}

func (odv *OpenDouyinVideo) SetCreateTime(ctx context.Context, createTime int64) {
	odv.CreateTime = createTime
}

func (odv *OpenDouyinVideo) SetIsReviewed(ctx context.Context, isReviewed bool) {
	odv.IsReviewed = isReviewed
}

func (odv *OpenDouyinVideo) SetItemId(ctx context.Context, itemId string) {
	odv.ItemId = itemId
}

func (odv *OpenDouyinVideo) SetStatistics(ctx context.Context, statistics VideoStatistics) {
	odv.Statistics = statistics
}

func (odv *OpenDouyinVideo) SetIsTop(ctx context.Context, isTop bool) {
	odv.IsTop = isTop
}

func (odv *OpenDouyinVideo) SetMediaType(ctx context.Context, mediaType int) {
	odv.MediaType = mediaType
}

func (odv *OpenDouyinVideo) SetShareUrl(ctx context.Context, shareUrl string) {
	odv.ShareUrl = shareUrl
}

func (odv *OpenDouyinVideo) SetVideoId(ctx context.Context, videoId string) {
	odv.VideoId = videoId
}

func (odv *OpenDouyinVideo) SetVideoStatus(ctx context.Context, videoStatus int32) {
	odv.VideoStatus = videoStatus
}

func (odv *OpenDouyinVideo) SetProductId(ctx context.Context, productId string) {
	odv.ProductId = productId
}

func (odv *OpenDouyinVideo) SetProductName(ctx context.Context, productName string) {
	odv.ProductName = productName
}

func (odv *OpenDouyinVideo) SetProductImg(ctx context.Context, productImg string) {
	odv.ProductImg = productImg
}

func (odv *OpenDouyinVideo) SetIsUpdateCover(ctx context.Context, isUpdateCover uint8) {
	odv.IsUpdateCover = isUpdateCover
}

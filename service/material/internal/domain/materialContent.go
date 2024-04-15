package domain

import (
	"context"
	"time"
)

type MaterialContentResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		VideoId    uint64 `json:"video_id"`
		OcrContent string `json:"ocr_content"`
	} `json:"data"`
}

type MaterialContent struct {
	Id         uint64
	ProductId  uint64
	UserId     uint64
	VideoId    uint64
	Content    string
	VideoName  string
	VideoUrl   string
	VideoCover string
	CreateTime time.Time
	UpdateTime time.Time
}

type MaterialContentList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*MaterialContent
}

func NewMaterialContent(ctx context.Context, productId, userId, videoId uint64, content string) *MaterialContent {
	return &MaterialContent{
		ProductId: productId,
		UserId:    userId,
		VideoId:   videoId,
		Content:   content,
	}
}

func (mc *MaterialContent) SetProductId(ctx context.Context, productId, userId uint64) {
	mc.ProductId = productId
}

func (mc *MaterialContent) SetUserId(ctx context.Context, userId uint64) {
	mc.UserId = userId
}

func (mc *MaterialContent) SetVideoId(ctx context.Context, videoId uint64) {
	mc.VideoId = videoId
}

func (mc *MaterialContent) SetContent(ctx context.Context, content string) {
	mc.Content = content
}

func (mc *MaterialContent) SetVideoName(ctx context.Context, videoName string) {
	mc.VideoName = videoName
}

func (mc *MaterialContent) SetVideoUrl(ctx context.Context, videoUrl string) {
	mc.VideoUrl = videoUrl
}

func (mc *MaterialContent) SetVideoCover(ctx context.Context, videoCover string) {
	mc.VideoCover = videoCover
}

func (mc *MaterialContent) SetCreateTime(ctx context.Context) {
	mc.CreateTime = time.Now()
}

func (mc *MaterialContent) SetUpdateTime(ctx context.Context) {
	mc.UpdateTime = time.Now()
}

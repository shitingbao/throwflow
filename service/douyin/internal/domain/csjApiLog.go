package domain

import (
	"context"
	"time"
)

type CsjApiLog struct {
	Content    string
	CreateTime time.Time
}

func NewCsjApiLog(ctx context.Context, content string) *CsjApiLog {
	return &CsjApiLog{
		Content: content,
	}
}

func (cal *CsjApiLog) SetContent(ctx context.Context, content string) {
	cal.Content = content
}

func (cal *CsjApiLog) SetCreateTime(ctx context.Context) {
	cal.CreateTime = time.Now()
}

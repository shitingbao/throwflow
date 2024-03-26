package domain

import (
	"context"
	"time"
)

type PayAsyncNotificationLog struct {
	Id         uint64
	Content    string
	CreateTime time.Time
	UpdateTime time.Time
}

func NewPayAsyncNotificationLog(ctx context.Context, content string) *PayAsyncNotificationLog {
	return &PayAsyncNotificationLog{
		Content: content,
	}
}

func (panl *PayAsyncNotificationLog) SetContent(ctx context.Context, content string) {
	panl.Content = content
}

func (panl *PayAsyncNotificationLog) SetUpdateTime(ctx context.Context) {
	panl.UpdateTime = time.Now()
}

func (panl *PayAsyncNotificationLog) SetCreateTime(ctx context.Context) {
	panl.CreateTime = time.Now()
}

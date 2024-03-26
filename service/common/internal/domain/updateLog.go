package domain

import (
	"context"
	"time"
)

type UpdateLog struct {
	Id         uint64
	Name       string
	Content    string
	CreateTime time.Time
	UpdateTime time.Time
}

func NewUpdateLog(ctx context.Context, name, content string) *UpdateLog {
	return &UpdateLog{
		Name:    name,
		Content: content,
	}
}

func (ul *UpdateLog) SetName(ctx context.Context, name string) {
	ul.Name = name
}

func (ul *UpdateLog) SetContent(ctx context.Context, content string) {
	ul.Content = content
}

func (ul *UpdateLog) SetUpdateTime(ctx context.Context) {
	ul.UpdateTime = time.Now()
}

func (ul *UpdateLog) SetCreateTime(ctx context.Context) {
	ul.CreateTime = time.Now()
}

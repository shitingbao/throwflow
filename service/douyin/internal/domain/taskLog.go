package domain

import (
	"context"
	"time"
)

type TaskLog struct {
	TaskType   string
	Content    string
	CreateTime time.Time
}

func NewTaskLog(ctx context.Context, taskType, content string) *TaskLog {
	return &TaskLog{
		TaskType: taskType,
		Content:  content,
	}
}

func (tl *TaskLog) SetCreateTime(ctx context.Context) {
	tl.CreateTime = time.Now()
}

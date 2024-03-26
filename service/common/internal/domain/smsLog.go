package domain

import (
	"context"
	"time"
)

type SmsLog struct {
	Id           uint64
	SendPhone    string
	Code         string
	ReplyCode    string
	SendContent  string
	ReplyContent string
	SendType     string
	SendIp       string
	CreateTime   time.Time
	UpdateTime   time.Time
}

type SmsLogList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*SmsLog
}

func NewSmsLog(ctx context.Context, sendPhone, code, sendType, sendContent, sendIp string) *SmsLog {
	return &SmsLog{
		SendPhone:   sendPhone,
		Code:        code,
		SendContent: sendContent,
		SendType:    sendType,
		SendIp:      sendIp,
	}
}

func (sl *SmsLog) SetReplyCode(ctx context.Context, replyCode string) {
	sl.ReplyCode = replyCode
}

func (sl *SmsLog) SetReplyContent(ctx context.Context, replyContent string) {
	sl.ReplyContent = replyContent
}

func (sl *SmsLog) SetUpdateTime(ctx context.Context) {
	sl.UpdateTime = time.Now()
}

func (sl *SmsLog) SetCreateTime(ctx context.Context) {
	sl.CreateTime = time.Now()
}

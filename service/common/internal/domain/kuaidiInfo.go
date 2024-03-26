package domain

import (
	"common/internal/pkg/kuaidi/kuaidi"
	"context"
	"time"
)

type KuaidiInfo struct {
	Id         uint64
	Code       string
	Num        string
	Phone      string
	State      uint8
	Content    string
	CreateTime time.Time
	UpdateTime time.Time
}

type KuaidiInfoData struct {
	Code      string
	Name      string
	Num       string
	State     uint8
	StateName string
	Content   []*kuaidi.Data
}

func NewKuaidiInfo(ctx context.Context, state uint8, code, num, phone, content string) *KuaidiInfo {
	return &KuaidiInfo{
		Code:    code,
		Num:     num,
		Phone:   phone,
		State:   state,
		Content: content,
	}
}

func (ki *KuaidiInfo) SetCode(ctx context.Context, code string) {
	ki.Code = code
}

func (ki *KuaidiInfo) SetNum(ctx context.Context, num string) {
	ki.Num = num
}

func (ki *KuaidiInfo) SetPhone(ctx context.Context, phone string) {
	ki.Phone = phone
}

func (ki *KuaidiInfo) SetState(ctx context.Context, state uint8) {
	ki.State = state
}

func (ki *KuaidiInfo) GetStateName(ctx context.Context) (stateName string) {
	if ki.State == 0 {
		stateName = "在途"
	} else if ki.State == 1 {
		stateName = "揽收"
	} else if ki.State == 2 {
		stateName = "疑难"
	} else if ki.State == 3 {
		stateName = "签收"
	} else if ki.State == 4 {
		stateName = "退签"
	} else if ki.State == 5 {
		stateName = "派件"
	} else if ki.State == 6 {
		stateName = "退回"
	} else if ki.State == 7 {
		stateName = "转投"
	} else if ki.State == 8 {
		stateName = "清关"
	} else if ki.State == 14 {
		stateName = "拒签"
	}

	return
}

func (ki *KuaidiInfo) SetContent(ctx context.Context, content string) {
	ki.Content = content
}

func (ki *KuaidiInfo) SetUpdateTime(ctx context.Context) {
	ki.UpdateTime = time.Now()
}

func (ki *KuaidiInfo) SetCreateTime(ctx context.Context) {
	ki.CreateTime = time.Now()
}

package domain

import (
	"context"
)

type DjAweme struct {
	UserName      string
	Avatar        string
	AwemeId       string
	HotsoonId     string
	FansCount     uint64
	Ratio         string
	Account       string
	BindStartTime string
	BindEndTime   string
}

func NewDjAweme(ctx context.Context, fansCount uint64, userName, avatar, awemeId, hotsoonId, ratio, account, bindStartTime, bindEndTime string) *DjAweme {
	return &DjAweme{
		UserName:      userName,
		Avatar:        avatar,
		AwemeId:       awemeId,
		HotsoonId:     hotsoonId,
		FansCount:     fansCount,
		Ratio:         ratio,
		Account:       account,
		BindStartTime: bindStartTime,
		BindEndTime:   bindEndTime,
	}
}

func (da *DjAweme) SetUserName(ctx context.Context, userName string) {
	da.UserName = userName
}

func (da *DjAweme) SetAvatar(ctx context.Context, avatar string) {
	da.Avatar = avatar
}

func (da *DjAweme) SetAwemeId(ctx context.Context, awemeId string) {
	da.AwemeId = awemeId
}

func (da *DjAweme) SetHotsoonId(ctx context.Context, hotsoonId string) {
	da.HotsoonId = hotsoonId
}

func (da *DjAweme) SetAuthorAccount(ctx context.Context, fansCount uint64) {
	da.FansCount = fansCount
}

func (da *DjAweme) SetRatio(ctx context.Context, ratio string) {
	da.Ratio = ratio
}

func (da *DjAweme) SetAccount(ctx context.Context, account string) {
	da.Account = account
}

func (da *DjAweme) SetBindStartTime(ctx context.Context, bindStartTime string) {
	da.BindStartTime = bindStartTime
}

func (da *DjAweme) SetBindEndTime(ctx context.Context, bindEndTime string) {
	da.BindEndTime = bindEndTime
}

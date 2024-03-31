package domain

import (
	"context"
	"time"
)

type User struct {
	Id                uint64
	Phone             string
	CountryCode       string
	NickName          string
	AvatarUrl         string
	Balance           float64
	Integral          uint64
	IntegralLevelName string
	QrCodeUrl         string
	IdentityCardMark  string
	Ranking           uint64
	Total             uint64
	TotalRanking      uint64
	Token             string
	OpenIds           []*UserOpenId
	CreateTime        time.Time
	UpdateTime        time.Time
}

type FollowData struct {
	FollowType    string
	FollowName    string
	FollowLogoUrl string
	QrCodeUrl     string
	TotalNum      uint64
}

func NewUser(ctx context.Context, totalRanking uint64, phone, countryCode string, balance float64) *User {
	return &User{
		CountryCode:  countryCode,
		Phone:        phone,
		Balance:      balance,
		TotalRanking: totalRanking,
	}
}

func (u *User) SetPhone(ctx context.Context, phone string) {
	u.Phone = phone
}

func (u *User) SetCountryCode(ctx context.Context, countryCode string) {
	u.CountryCode = countryCode
}

func (u *User) SetNickName(ctx context.Context, nickName string) {
	u.NickName = nickName
}

func (u *User) SetAvatarUrl(ctx context.Context, avatarUrl string) {
	u.AvatarUrl = avatarUrl
}

func (u *User) SetDefaultAvatarUrl(ctx context.Context, defaultAvatarUrl string) {
	if len(u.AvatarUrl) == 0 {
		u.AvatarUrl = defaultAvatarUrl
	}
}

func (u *User) SetBalance(ctx context.Context, balance float64) {
	u.Balance = balance
}

func (u *User) SetIntegral(ctx context.Context, integral uint64) {
	u.Integral = integral
}

func (u *User) SetIntegralLevelName(ctx context.Context) {
	if u.Ranking < 3000 {
		u.IntegralLevelName = "铁器"
	} else if u.Ranking >= 3000 && u.Ranking < 10000 {
		u.IntegralLevelName = "青铜"
	} else if u.Ranking >= 10000 && u.Ranking < 100000 {
		u.IntegralLevelName = "白银"
	} else if u.Ranking >= 100000 && u.Ranking < 500000 {
		u.IntegralLevelName = "黄金"
	} else if u.Ranking >= 500000 && u.Ranking < 1000000 {
		u.IntegralLevelName = "铂金"
	} else if u.Ranking >= 1000000 && u.Ranking < 3000000 {
		u.IntegralLevelName = "钻石"
	} else if u.Ranking >= 3000000 {
		u.IntegralLevelName = "王者"
	}
}

func (u *User) SetQrCodeUrl(ctx context.Context, qrCodeUrl string) {
	u.QrCodeUrl = qrCodeUrl
}

func (u *User) SetIdentityCardMark(ctx context.Context, identityCardMark string) {
	u.IdentityCardMark = identityCardMark
}

func (u *User) SetRanking(ctx context.Context, ranking uint64) {
	u.Ranking = ranking
}

func (u *User) SetTotal(ctx context.Context, total uint64) {
	u.Total = total
}

func (u *User) SetTotalRanking(ctx context.Context, totalRanking uint64) {
	u.TotalRanking = totalRanking
}

func (u *User) SetToken(ctx context.Context, token string) {
	u.Token = token
}

func (u *User) SetUpdateTime(ctx context.Context) {
	u.UpdateTime = time.Now()
}

func (u *User) SetCreateTime(ctx context.Context) {
	u.CreateTime = time.Now()
}

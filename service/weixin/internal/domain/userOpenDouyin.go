package domain

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

type UserOpenDouyin struct {
	Id              uint64
	UserId          uint64
	ClientKey       string
	OpenId          string
	AwemeId         uint64
	AccountId       string
	Nickname        string
	Avatar          string
	AvatarLarger    string
	CooperativeCode string
	Fans            uint64
	FansShow        string
	Area            string
	CreateTime      time.Time
	UpdateTime      time.Time
}

type ExternalUserOpenDouyin struct {
	Id               uint64
	UserId           uint64
	ClientKey        string
	OpenId           string
	AccountId        string
	Nickname         string
	Avatar           string
	AvatarLarger     string
	CooperativeCode  string
	Fans             uint64
	FansShow         string
	CreateTime       time.Time
	UpdateTime       time.Time
	UserSampleOrders []*UserSampleOrder
}

type UserOpenDouyinList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*UserOpenDouyin
}

type UserOpenDouyinClientKeyAndOpenId struct {
	ClientKey string
	OpenId    string
}

func NewUserOpenDouyin(ctx context.Context, userId uint64, clientKey, openId, accountId, nickname, avatar, avatarLarger, cooperativeCode string) *UserOpenDouyin {
	return &UserOpenDouyin{
		UserId:          userId,
		ClientKey:       clientKey,
		OpenId:          openId,
		AccountId:       accountId,
		Nickname:        nickname,
		Avatar:          avatar,
		AvatarLarger:    avatarLarger,
		CooperativeCode: cooperativeCode,
	}
}

func (uod *UserOpenDouyin) SetUserId(ctx context.Context, userId uint64) {
	uod.UserId = userId
}

func (uod *UserOpenDouyin) SetClientKey(ctx context.Context, clientKey string) {
	uod.ClientKey = clientKey
}

func (uod *UserOpenDouyin) SetOpenId(ctx context.Context, openId string) {
	uod.OpenId = openId
}

func (uod *UserOpenDouyin) SetAwemeId(ctx context.Context, awemeId uint64) {
	uod.AwemeId = awemeId
}

func (uod *UserOpenDouyin) SetAccountId(ctx context.Context, accountId string) {
	uod.AccountId = accountId
}

func (uod *UserOpenDouyin) SetNickname(ctx context.Context, nickname string) {
	uod.Nickname = nickname
}

func (uod *UserOpenDouyin) SetAvatar(ctx context.Context, avatar string) {
	uod.Avatar = avatar
}

func (uod *UserOpenDouyin) SetAvatarLarger(ctx context.Context, avatarLarger string) {
	uod.AvatarLarger = avatarLarger
}

func (uod *UserOpenDouyin) SetCooperativeCode(ctx context.Context, cooperativeCode string) {
	uod.CooperativeCode = cooperativeCode
}

func (uod *UserOpenDouyin) SetFans(ctx context.Context, fans uint64) {
	uod.Fans = fans
}

func (uod *UserOpenDouyin) SetArea(ctx context.Context, area string) {
	uod.Area = area
}

func (uod *UserOpenDouyin) SetUpdateTime(ctx context.Context) {
	uod.UpdateTime = time.Now()
}

func (uod *UserOpenDouyin) SetCreateTime(ctx context.Context) {
	uod.CreateTime = time.Now()
}

func (uod *UserOpenDouyin) SetFansShow(ctx context.Context) {
	if uod.Fans > 10000 {
		uod.FansShow = fmt.Sprintf("%.1f", float64(uod.Fans)/float64(10000)) + "w"
	} else {
		uod.FansShow = strconv.FormatUint(uod.Fans, 10)
	}
}

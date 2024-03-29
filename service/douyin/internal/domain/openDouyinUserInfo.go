package domain

import (
	"context"
	"time"
)

type OpenDouyinUserInfo struct {
	Id              uint64
	ClientKey       string
	OpenId          string
	UnionId         string
	AwemeId         uint64
	AccountId       string
	BuyinId         string
	Nickname        string
	Avatar          string
	AvatarLarger    string
	Phone           string
	Gender          uint8
	Country         string
	Province        string
	City            string
	District        string
	EAccountRole    string
	CooperativeCode string
	Fans            uint64
	Area            string
	VideoId         string
	CreateTime      time.Time
	UpdateTime      time.Time
}

type OpenDouyinUserInfoList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*OpenDouyinUserInfo
}

func NewOpenDouyinUserInfo(ctx context.Context, clientKey, openId, unionId, accountId, buyinId, nickname, avatar, avatarLarger, phone, country, province, city, district, eAccountRole, cooperativeCode string, gender uint8) *OpenDouyinUserInfo {
	return &OpenDouyinUserInfo{
		ClientKey:       clientKey,
		OpenId:          openId,
		UnionId:         unionId,
		AccountId:       accountId,
		BuyinId:         buyinId,
		Nickname:        nickname,
		Avatar:          avatar,
		AvatarLarger:    avatarLarger,
		Phone:           phone,
		Gender:          gender,
		Country:         country,
		Province:        province,
		City:            city,
		District:        district,
		EAccountRole:    eAccountRole,
		CooperativeCode: cooperativeCode,
	}
}

func (odui *OpenDouyinUserInfo) SetClientKey(ctx context.Context, clientKey string) {
	odui.ClientKey = clientKey
}

func (odui *OpenDouyinUserInfo) SetOpenId(ctx context.Context, openId string) {
	odui.OpenId = openId
}

func (odui *OpenDouyinUserInfo) SetUnionId(ctx context.Context, unionId string) {
	odui.UnionId = unionId
}

func (odui *OpenDouyinUserInfo) SetAwemeId(ctx context.Context, awemeId uint64) {
	odui.AwemeId = awemeId
}

func (odui *OpenDouyinUserInfo) SetAccountId(ctx context.Context, accountId string) {
	odui.AccountId = accountId
}

func (odui *OpenDouyinUserInfo) SetBuyinId(ctx context.Context, buyinId string) {
	odui.BuyinId = buyinId
}

func (odui *OpenDouyinUserInfo) SetNickname(ctx context.Context, nickname string) {
	odui.Nickname = nickname
}

func (odui *OpenDouyinUserInfo) SetAvatar(ctx context.Context, avatar string) {
	odui.Avatar = avatar
}

func (odui *OpenDouyinUserInfo) SetAvatarLarger(ctx context.Context, avatarLarger string) {
	odui.AvatarLarger = avatarLarger
}

func (odui *OpenDouyinUserInfo) SetPhone(ctx context.Context, phone string) {
	odui.Phone = phone
}

func (odui *OpenDouyinUserInfo) SetGender(ctx context.Context, gender uint8) {
	odui.Gender = gender
}

func (odui *OpenDouyinUserInfo) SetCountry(ctx context.Context, country string) {
	odui.Country = country
}

func (odui *OpenDouyinUserInfo) SetProvince(ctx context.Context, province string) {
	odui.Province = province
}

func (odui *OpenDouyinUserInfo) SetCity(ctx context.Context, city string) {
	odui.City = city
}

func (odui *OpenDouyinUserInfo) SetDistrict(ctx context.Context, district string) {
	odui.District = district
}

func (odui *OpenDouyinUserInfo) SetEAccountRole(ctx context.Context, eAccountRole string) {
	odui.EAccountRole = eAccountRole
}

func (odui *OpenDouyinUserInfo) SetCooperativeCode(ctx context.Context, cooperativeCode string) {
	odui.CooperativeCode = cooperativeCode
}

func (odui *OpenDouyinUserInfo) SetFans(ctx context.Context, fans uint64) {
	odui.Fans = fans
}

func (odui *OpenDouyinUserInfo) SetArea(ctx context.Context, area string) {
	odui.Area = area
}

func (odui *OpenDouyinUserInfo) SetVideoId(ctx context.Context, videoId string) {
	odui.VideoId = videoId
}

func (odui *OpenDouyinUserInfo) SetUpdateTime(ctx context.Context) {
	odui.UpdateTime = time.Now()
}

func (odui *OpenDouyinUserInfo) SetCreateTime(ctx context.Context) {
	odui.CreateTime = time.Now()
}

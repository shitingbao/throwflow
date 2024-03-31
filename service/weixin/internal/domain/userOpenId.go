package domain

import (
	"context"
	"time"
)

type UserOpenId struct {
	Id         uint64
	UserId     uint64
	Appid      string
	OpenId     string
	CreateTime time.Time
	UpdateTime time.Time
}

func NewUserOpenId(ctx context.Context, userId uint64, appid, openId string) *UserOpenId {
	return &UserOpenId{
		UserId: userId,
		Appid:  appid,
		OpenId: openId,
	}
}

func (uoi *UserOpenId) SetUserId(ctx context.Context, userId uint64) {
	uoi.UserId = userId
}

func (uoi *UserOpenId) SetAppid(ctx context.Context, appid string) {
	uoi.Appid = appid
}

func (uoi *UserOpenId) SetOpenId(ctx context.Context, openId string) {
	uoi.OpenId = openId
}

func (uoi *UserOpenId) SetUpdateTime(ctx context.Context) {
	uoi.UpdateTime = time.Now()
}

func (uoi *UserOpenId) SetCreateTime(ctx context.Context) {
	uoi.CreateTime = time.Now()
}

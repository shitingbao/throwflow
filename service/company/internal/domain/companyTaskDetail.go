package domain

import (
	v1 "company/api/service/weixin/v1"
	"context"
	"time"
)

type CompanyTaskDetail struct {
	Id                           uint64
	CompanyTaskId                uint64
	CompanyTaskAccountRelationId uint64
	ProductName                  string
	UserId                       uint64
	ClientKey                    string
	OpenId                       string
	VideoId                      string
	ItemId                       string
	PlayCount                    uint64
	Cover                        string
	ReleaseTime                  time.Time
	IsPlaySuccess                uint8
	CreateTime                   time.Time
	UpdateTime                   time.Time
	Nickname                     string
	Avatar                       string
	ClaimTime                    time.Time
	IsReleaseVideo               uint8
}

type CompanyTaskDetailList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*CompanyTaskDetail
}

func NewCompanyTaskDetail(ctx context.Context, companyTaskId, relationId, userId, playCount uint64, productName, clientKey, openId, itemId, cover, nickname, avatar string, releaseTime time.Time) *CompanyTaskDetail {
	return &CompanyTaskDetail{
		CompanyTaskId:                companyTaskId,
		CompanyTaskAccountRelationId: relationId,
		ProductName:                  productName,
		UserId:                       userId,
		ClientKey:                    clientKey,
		OpenId:                       openId,
		ItemId:                       itemId,
		PlayCount:                    playCount,
		Cover:                        cover,
		Nickname:                     nickname,
		Avatar:                       avatar,
		ReleaseTime:                  releaseTime,
	}
}

func (c *CompanyTaskDetail) SetUpdateTime(ctx context.Context) {
	c.UpdateTime = time.Now()
}

func (c *CompanyTaskDetail) SetCreateTime(ctx context.Context) {
	c.CreateTime = time.Now()
}

type CompanyTaskClientKeyAndOpenId struct {
	ClientKey string
	OpenId    string
	VideoId   string
}

func (c *CompanyTaskDetail) SetPlayCount(ctx context.Context, count uint64) {
	c.PlayCount = count
}

func (c *CompanyTaskDetail) SetVideoId(ctx context.Context, videoId string) {
	c.VideoId = videoId
}

func (c *CompanyTaskDetail) SetNicknameAndAvatar(ctx context.Context, list []*v1.ListByClientKeyAndOpenIdsReply_OpenDouyinUser) {
	for _, v := range list {
		if v.ClientKey == c.ClientKey && v.OpenId == c.OpenId {
			c.Nickname = v.Nickname
			c.Avatar = v.Avatar
		}
	}
}

func (c *CompanyTaskDetail) SetNicknameAndAvatarByCompanyIds(ctx context.Context, list []*v1.ListByClientKeyAndOpenIdsReply_OpenDouyinUser) {
	for _, v := range list {
		if v.ClientKey == c.ClientKey && v.OpenId == c.OpenId {
			c.Nickname = v.Nickname
			c.Avatar = v.Avatar
		}
	}
}

func (c *CompanyTaskDetail) SetIsReleaseVideo(ctx context.Context) {
	c.IsReleaseVideo = 1
}

func (c *CompanyTaskDetail) SetIsPlaySuccess(ctx context.Context, isPlaySuccess uint8) {
	c.IsPlaySuccess = isPlaySuccess
}

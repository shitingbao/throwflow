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
	CompanyTaskAccountRelation   CompanyTaskAccountRelation
}

type CompanyTaskDetailList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*CompanyTaskDetail
}

func NewCompanyTaskDetail(ctx context.Context, isPlaySuccess uint8, companyTaskId, companyTaskAccountRelationId, playCount uint64, itemId, cover string, releaseTime time.Time) *CompanyTaskDetail {
	return &CompanyTaskDetail{
		CompanyTaskId:                companyTaskId,
		CompanyTaskAccountRelationId: companyTaskAccountRelationId,
		ItemId:                       itemId,
		PlayCount:                    playCount,
		Cover:                        cover,
		ReleaseTime:                  releaseTime,
		IsPlaySuccess:                isPlaySuccess,
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
}

func (c *CompanyTaskDetail) SetPlayCount(ctx context.Context, count uint64) {
	c.PlayCount = count
}

// func (c *CompanyTaskDetail) SetIsScreenshotAvailable(ctx context.Context, isScreenshotAvailable uint8) {
// 	c.IsScreenshotAvailable = isScreenshotAvailable
// }

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

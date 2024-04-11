package domain

import (
	"context"
	"time"
)

const (
	GoingStatus = iota
	SuccessStatus
	ExpireStatus  // 过期
	SettledStatus // 已经结算
)

const (
	ScreenshotInvalid = iota
	ScreenshotAvailable
)

const (
	DoukeOrderSuccess = "PAY_SUCC"
	DoukeOrderREFUND  = "REFUND"  // 退款
	DoukeOrderSETTLE  = "SETTLE"  // 结算
	DoukeOrderCONFIRM = "CONFIRM" // 确认收货
)

type CompanyTaskAccountRelation struct {
	Id                    uint64
	NickName              string // 微信昵称
	AvatarUrl             string // 微信头像
	CompanyTaskId         uint64
	ProductOutId          uint64
	ProductName           string
	UserId                uint64
	ClaimTime             time.Time
	ExpireTime            time.Time
	Status                uint8
	IsDel                 uint8
	CreateTime            time.Time
	UpdateTime            time.Time
	IsCostBuy             uint8
	ScreenshotAddress     string
	IsScreenshotAvailable uint8
	IsPlaySuccess         uint8
	CompanyTaskDetails    []*CompanyTaskDetail
	CompanyTask           CompanyTask
	CompanyProduct        *CompanyProduct
}

type CompanyTaskAccountRelationList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*CompanyTaskAccountRelation
}

func NewCompanyTaskAccountRelation(ctx context.Context, companyTaskId, userId, productOutId uint64) *CompanyTaskAccountRelation {
	return &CompanyTaskAccountRelation{
		CompanyTaskId: companyTaskId,
		ProductOutId:  productOutId,
		UserId:        userId,
	}
}

func (c *CompanyTaskAccountRelation) SetNickName(ctx context.Context, nickName string) {
	c.NickName = nickName
}

func (c *CompanyTaskAccountRelation) SetAvatarUrl(ctx context.Context, avatarUrl string) {
	c.AvatarUrl = avatarUrl
}

func (c *CompanyTaskAccountRelation) SetClaimTime(ctx context.Context) {
	c.ClaimTime = time.Now()
}

func (c *CompanyTaskAccountRelation) SetExpireTime(ctx context.Context, tm time.Time) {
	// c.ExpireTime = time.Now().AddDate(0, 0, day)
	c.ExpireTime = tm
}

func (c *CompanyTaskAccountRelation) SetUpdateTime(ctx context.Context) {
	c.UpdateTime = time.Now()
}

func (c *CompanyTaskAccountRelation) SetCreateTime(ctx context.Context) {
	c.CreateTime = time.Now()
}

func (c *CompanyTaskAccountRelation) SetIsScreenshotAvailable(ctx context.Context, isScreenshotAvailable uint8) {
	c.IsScreenshotAvailable = isScreenshotAvailable
}

func (c *CompanyTaskAccountRelation) SetScreenshotAddress(ctx context.Context, address string) {
	c.ScreenshotAddress = address
}

func (c *CompanyTaskAccountRelation) SetStatus(ctx context.Context, status uint8) {
	c.Status = status
}

func (c *CompanyTaskAccountRelation) SetIsCostBuy(ctx context.Context, isCostBuy uint8) {
	c.IsCostBuy = isCostBuy
}

func (c *CompanyTaskAccountRelation) SetIsPlaySuccess(ctx context.Context, isPlaySuccess uint8) {
	c.IsPlaySuccess = isPlaySuccess
}

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

package domain

import (
	"context"
	"time"
)

const (
	GoingStatus = iota
	SuccessStatus
	ExpireStatus
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

func NewCompanyTaskAccountRelation(ctx context.Context,
	companyTaskId, userId, productOutId uint64, productName string) *CompanyTaskAccountRelation {
	return &CompanyTaskAccountRelation{
		CompanyTaskId: companyTaskId,
		ProductOutId:  productOutId,
		ProductName:   productName,
		UserId:        userId,
	}
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

type VideoStatistics struct {
	CommentCount  int32 `json:"comment_count" bson:"comment_count"`
	DiggCount     int32 `json:"digg_count" bson:"digg_count"`
	DownloadCount int32 `json:"download_count" bson:"download_count"`
	ForwardCount  int32 `json:"forward_count" bson:"forward_count"`
	PlayCount     int32 `json:"play_count" bson:"play_count"`
	ShareCount    int32 `json:"share_count" bson:"share_count"`
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

type OpenDouyinVideo struct {
	ClientKey     string
	OpenId        string
	AwemeId       uint64
	AccountId     string
	Nickname      string
	Avatar        string
	Title         string
	Cover         string
	CreateTime    uint64
	IsReviewed    bool
	ItemId        string
	Statistics    VideoStatistics
	IsTop         bool
	MediaType     uint64
	ShareUrl      string
	VideoId       string
	VideoStatus   uint32
	ProductId     string
	ProductName   string
	ProductImg    string
	IsUpdateCover uint8
}

type DoukeOrderInfo struct {
	Id                  uint64
	UserId              uint64
	OrderId             string
	ProductId           string
	ProductName         string
	ProductImg          string
	PaySuccessTime      time.Time
	SettleTime          *time.Time
	TotalPayAmount      float32
	PayGoodsAmount      float32
	FlowPoint           string
	EstimatedCommission float32
	RealCommission      float32
	ItemNum             uint64
	CreateTime          time.Time
	UpdateTime          time.Time
}

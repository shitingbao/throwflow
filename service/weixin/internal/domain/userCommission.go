package domain

import (
	"context"
	"time"
)

type UserCommission struct {
	Id                      uint64
	UserId                  uint64
	OrganizationId          uint64
	ChildUserId             uint64
	ChildNickName           string
	ChildAvatarUrl          string
	ChildPhone              string
	ChildLevel              uint8
	ChildLevelName          string
	Level                   uint8
	Relation                uint8
	RelationName            string
	OrderNum                uint64
	OrderRefundNum          uint64
	TotalPayAmount          float32
	CommissionPool          float32
	EstimatedCommission     float32
	RealCommission          float32
	EstimatedUserCommission float32
	RealUserCommission      float32
	CommissionRatio         float32
	Day                     uint32
	UserCommissionType      uint8
	UserCommissionTypeName  string
	ActivationTime          string
	CreateTime              time.Time
	UpdateTime              time.Time
}

type UserCommissionList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*UserCommission
}

type UserOrganizationCommission struct {
	UserId                       uint64
	OrganizationUserId           uint64
	OrganizationTutorId          uint64
	Level                        uint8
	IsOrderRelation              uint8
	OrderNum                     uint64
	OrderRefundNum               uint64
	TotalPayAmount               float64
	EstimatedCommission          float64
	RealCommission               float64
	CoursePayAmount              float64
	CostOrderTotalPayAmount      float64
	CostOrderEstimatedCommission float64
	CostOrderRealCommission      float64
}

type UserCommissionOpenDouyin struct {
	ClientKey string `json:"clientKey"`
	OpenId    string `json:"openId"`
}

func NewUserCommission(ctx context.Context, userId, organizationId, childUserId, orderNum, orderRefundNum uint64, day uint32, childLevel, level, relation, userCommissionType uint8, totalPayAmount, commissionPool, estimatedCommission, realCommission, estimatedUserCommission, realUserCommission float32) *UserCommission {
	return &UserCommission{
		UserId:                  userId,
		OrganizationId:          organizationId,
		ChildUserId:             childUserId,
		ChildLevel:              childLevel,
		Level:                   level,
		Relation:                relation,
		OrderNum:                orderNum,
		OrderRefundNum:          orderRefundNum,
		TotalPayAmount:          totalPayAmount,
		CommissionPool:          commissionPool,
		EstimatedCommission:     estimatedCommission,
		RealCommission:          realCommission,
		EstimatedUserCommission: estimatedUserCommission,
		RealUserCommission:      realUserCommission,
		Day:                     day,
		UserCommissionType:      userCommissionType,
	}
}

func (uc *UserCommission) SetUserId(ctx context.Context, userId uint64) {
	uc.UserId = userId
}

func (uc *UserCommission) SetOrganizationId(ctx context.Context, organizationId uint64) {
	uc.OrganizationId = organizationId
}

func (uc *UserCommission) SetChildUserId(ctx context.Context, childUserId uint64) {
	uc.ChildUserId = childUserId
}

func (uc *UserCommission) SetChildLevel(ctx context.Context, childLevel uint8) {
	uc.ChildLevel = childLevel
}

func (uc *UserCommission) SetLevel(ctx context.Context, level uint8) {
	uc.Level = level
}

func (uc *UserCommission) SetRelation(ctx context.Context, relation uint8) {
	uc.Relation = relation
}

func (uc *UserCommission) GetRelationName(ctx context.Context) (relationName string) {
	if uc.Relation == 1 {
		relationName = "直接"
	} else if uc.Relation == 2 {
		relationName = "间接"
	}

	return
}

func (uc *UserCommission) SetOrderNum(ctx context.Context, orderNum uint64) {
	uc.OrderNum = orderNum
}

func (uc *UserCommission) SetOrderRefundNum(ctx context.Context, orderRefundNum uint64) {
	uc.OrderRefundNum = orderRefundNum
}

func (uc *UserCommission) SetTotalPayAmount(ctx context.Context, totalPayAmount float32) {
	uc.TotalPayAmount = totalPayAmount
}

func (uc *UserCommission) SetCommissionPool(ctx context.Context, commissionPool float32) {
	uc.CommissionPool = commissionPool
}

func (uc *UserCommission) SetEstimatedCommission(ctx context.Context, estimatedCommission float32) {
	uc.EstimatedCommission = estimatedCommission
}

func (uc *UserCommission) SetRealCommission(ctx context.Context, realCommission float32) {
	uc.RealCommission = realCommission
}

func (uc *UserCommission) SetEstimatedUserCommission(ctx context.Context, estimatedUserCommission float32) {
	uc.EstimatedUserCommission = estimatedUserCommission
}

func (uc *UserCommission) SetRealUserCommission(ctx context.Context, realUserCommission float32) {
	uc.RealUserCommission = realUserCommission
}

func (uc *UserCommission) SetDay(ctx context.Context, day uint32) {
	uc.Day = day
}

func (uc *UserCommission) SetUserCommissionType(ctx context.Context, userCommissionType uint8) {
	uc.UserCommissionType = userCommissionType
}

func (uc *UserCommission) GetUserCommissionTypeName(ctx context.Context) (userCommissionTypeName string) {
	if uc.UserCommissionType == 1 {
		userCommissionTypeName = "课程"
	} else if uc.UserCommissionType == 2 {
		userCommissionTypeName = "电商"
	} else if uc.UserCommissionType == 3 {
		userCommissionTypeName = "成本购"
	}

	return
}

func (uc *UserCommission) SetUpdateTime(ctx context.Context) {
	uc.UpdateTime = time.Now()
}

func (uc *UserCommission) SetCreateTime(ctx context.Context) {
	uc.CreateTime = time.Now()
}

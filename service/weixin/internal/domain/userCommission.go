package domain

import (
	"context"
	"time"
)

type UserCommission struct {
	Id                 uint64
	UserId             uint64
	OrganizationId     uint64
	ChildUserId        uint64
	ChildNickName      string
	ChildAvatarUrl     string
	ChildPhone         string
	ChildLevel         uint8
	Level              uint8
	Relation           uint8
	RelationName       string
	CommissionStatus   uint8
	CommissionType     uint8
	CommissionTypeName string
	RelevanceId        uint64
	OperationType      uint8
	TotalPayAmount     float32
	CommissionPool     float32
	CommissionMcnRatio float32
	CommissionAmount   float32
	Name               string
	IdentityCard       string
	BankCode           string
	BalanceStatus      uint8
	OutTradeNo         string
	InnerTradeNo       string
	RealAmount         float32
	CommissionRatio    float32
	ApplyTime          *time.Time
	GongmallCreateTime *time.Time
	PayTime            *time.Time
	ActivationTime     string
	OperationContent   string
	CreateTime         time.Time
	UpdateTime         time.Time
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

func NewUserCommission(ctx context.Context, userId, organizationId, childUserId, relevanceId uint64, childLevel, level, relation, commissionStatus, commissionType, operationType, balanceStatus uint8, totalPayAmount, commissionPool, commissionAmount float32) *UserCommission {
	return &UserCommission{
		UserId:           userId,
		OrganizationId:   organizationId,
		ChildUserId:      childUserId,
		ChildLevel:       childLevel,
		Level:            level,
		Relation:         relation,
		CommissionStatus: commissionStatus,
		CommissionType:   commissionType,
		RelevanceId:      relevanceId,
		OperationType:    operationType,
		TotalPayAmount:   totalPayAmount,
		CommissionPool:   commissionPool,
		CommissionAmount: commissionAmount,
		BalanceStatus:    balanceStatus,
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

func (uc *UserCommission) SetCommissionStatus(ctx context.Context, commissionStatus uint8) {
	uc.CommissionStatus = commissionStatus
}

func (uc *UserCommission) SetCommissionType(ctx context.Context, commissionType uint8) {
	uc.CommissionType = commissionType
}

func (uc *UserCommission) SetBalanceCommissionTypeName(ctx context.Context) {
	if uc.CommissionType == 2 {
		uc.CommissionTypeName = "带货团队分佣"
	} else if uc.CommissionType == 3 {
		uc.CommissionTypeName = "成本购返佣"
	} else if uc.CommissionType == 4 {
		uc.CommissionTypeName = "成本购推荐"
	} else if uc.CommissionType == 1 {
		if uc.ChildLevel == 1 {
			uc.CommissionTypeName = "零级推荐"
		} else if uc.ChildLevel == 2 {
			uc.CommissionTypeName = "初级推荐"
		} else if uc.ChildLevel == 3 {
			uc.CommissionTypeName = "中级推荐"
		} else if uc.ChildLevel == 4 {
			uc.CommissionTypeName = "高级推荐"
		}
	} else if uc.CommissionType == 6 {
		uc.CommissionTypeName = "其他"
	} else if uc.CommissionType == 5 {
		uc.CommissionTypeName = "返现任务"
	} else if uc.CommissionType == 7 || uc.CommissionType == 8 {
		if uc.BalanceStatus == 0 || uc.BalanceStatus == 2 {
			uc.CommissionTypeName = "提现处理中"
		} else if uc.BalanceStatus == 1 {
			uc.CommissionTypeName = "提现成功"
		} else if uc.BalanceStatus == 3 {
			uc.CommissionTypeName = "提现失败"
		}
	}
}

func (uc *UserCommission) GetCommissionTypeName(ctx context.Context) string {
	if uc.CommissionType == 1 {
		return "会员"
	} else if uc.CommissionType == 2 {
		return "带货"
	} else if uc.CommissionType == 4 {
		return "成本购"
	}

	return ""
}

func (uc *UserCommission) SetRelevanceId(ctx context.Context, relevanceId uint64) {
	uc.RelevanceId = relevanceId
}

func (uc *UserCommission) SetOperationType(ctx context.Context, operationType uint8) {
	uc.OperationType = operationType
}

func (uc *UserCommission) SetTotalPayAmount(ctx context.Context, totalPayAmount float32) {
	uc.TotalPayAmount = totalPayAmount
}

func (uc *UserCommission) SetCommissionPool(ctx context.Context, commissionPool float32) {
	uc.CommissionPool = commissionPool
}

func (uc *UserCommission) SetCommissionMcnRatio(ctx context.Context, commissionMcnRatio float32) {
	uc.CommissionMcnRatio = commissionMcnRatio
}

func (uc *UserCommission) SetCommissionAmount(ctx context.Context, commissionAmount float32) {
	uc.CommissionAmount = commissionAmount
}

func (uc *UserCommission) SetName(ctx context.Context, name string) {
	uc.Name = name
}

func (uc *UserCommission) SetIdentityCard(ctx context.Context, identityCard string) {
	uc.IdentityCard = identityCard
}

func (uc *UserCommission) SetBankCode(ctx context.Context, bankCode string) {
	uc.BankCode = bankCode
}

func (uc *UserCommission) SetBalanceStatus(ctx context.Context, balanceStatus uint8) {
	uc.BalanceStatus = balanceStatus
}

func (uc *UserCommission) SetOutTradeNo(ctx context.Context, outTradeNo string) {
	uc.OutTradeNo = outTradeNo
}

func (uc *UserCommission) SetInnerTradeNo(ctx context.Context, innerTradeNo string) {
	uc.InnerTradeNo = innerTradeNo
}

func (uc *UserCommission) SetRealAmount(ctx context.Context, realAmount float32) {
	uc.RealAmount = realAmount
}

func (uc *UserCommission) SetApplyTime(ctx context.Context, applyTime *time.Time) {
	uc.ApplyTime = applyTime
}

func (uc *UserCommission) SetGongmallCreateTime(ctx context.Context, gongmallCreateTime *time.Time) {
	uc.GongmallCreateTime = gongmallCreateTime
}

func (uc *UserCommission) SetPayTime(ctx context.Context, payTime *time.Time) {
	uc.PayTime = payTime
}

func (uc *UserCommission) SetOperationContent(ctx context.Context, operationContent string) {
	uc.OperationContent = operationContent
}

func (uc *UserCommission) SetUpdateTime(ctx context.Context) {
	uc.UpdateTime = time.Now()
}

func (uc *UserCommission) SetCreateTime(ctx context.Context, createTime time.Time) {
	uc.CreateTime = createTime
}

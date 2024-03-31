package domain

import (
	"context"
	"time"
)

type UserBalanceLog struct {
	Id                 uint64
	UserId             uint64
	OrganizationId     uint64
	Name               string
	IdentityCard       string
	BankCode           string
	Amount             float32
	RelevanceId        uint64
	BalanceType        uint8
	OperationType      uint8
	OperationContent   string
	BalanceStatus      uint8
	OutTradeNo         string
	InnerTradeNo       string
	RealAmount         float32
	ApplyTime          *time.Time
	GongmallCreateTime *time.Time
	PayTime            *time.Time
	Day                uint32
	CreateTime         time.Time
	UpdateTime         time.Time
}

type UserBalanceLogList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*UserBalanceLog
}

type UserBalance struct {
	EstimatedCommissionBalance float32
	RealCommissionBalance      float32
	EstimatedCostBalance       float32
	RealCostBalance            float32
}

func NewUserBalanceLog(ctx context.Context, userId, relevanceId uint64, balanceType, operationType, balanceStatus uint8, amount float32, operationContent string) *UserBalanceLog {
	return &UserBalanceLog{
		UserId:           userId,
		Amount:           amount,
		RelevanceId:      relevanceId,
		BalanceType:      balanceType,
		OperationType:    operationType,
		OperationContent: operationContent,
		BalanceStatus:    balanceStatus,
	}
}

func (ubl *UserBalanceLog) SetUserId(ctx context.Context, userId uint64) {
	ubl.UserId = userId
}

func (ubl *UserBalanceLog) SetOrganizationId(ctx context.Context, organizationId uint64) {
	ubl.OrganizationId = organizationId
}

func (ubl *UserBalanceLog) SetName(ctx context.Context, name string) {
	ubl.Name = name
}

func (ubl *UserBalanceLog) SetIdentityCard(ctx context.Context, identityCard string) {
	ubl.IdentityCard = identityCard
}

func (ubl *UserBalanceLog) SetBankCode(ctx context.Context, bankCode string) {
	ubl.BankCode = bankCode
}

func (ubl *UserBalanceLog) SetAmount(ctx context.Context, amount float32) {
	ubl.Amount = amount
}

func (ubl *UserBalanceLog) SetRelevanceId(ctx context.Context, relevanceId uint64) {
	ubl.RelevanceId = relevanceId
}

func (ubl *UserBalanceLog) SetBalanceType(ctx context.Context, balanceType uint8) {
	ubl.BalanceType = balanceType
}

func (ubl *UserBalanceLog) SetOperationType(ctx context.Context, operationType uint8) {
	ubl.OperationType = operationType
}

func (ubl *UserBalanceLog) SetOperationContent(ctx context.Context, operationContent string) {
	ubl.OperationContent = operationContent
}

func (ubl *UserBalanceLog) SetBalanceStatus(ctx context.Context, balanceStatus uint8) {
	ubl.BalanceStatus = balanceStatus
}

func (ubl *UserBalanceLog) SetOutTradeNo(ctx context.Context, outTradeNo string) {
	ubl.OutTradeNo = outTradeNo
}

func (ubl *UserBalanceLog) SetInnerTradeNo(ctx context.Context, innerTradeNo string) {
	ubl.InnerTradeNo = innerTradeNo
}

func (ubl *UserBalanceLog) SetRealAmount(ctx context.Context, realAmount float32) {
	ubl.RealAmount = realAmount
}

func (ubl *UserBalanceLog) SetApplyTime(ctx context.Context, applyTime *time.Time) {
	ubl.ApplyTime = applyTime
}

func (ubl *UserBalanceLog) SetGongmallCreateTime(ctx context.Context, gongmallCreateTime *time.Time) {
	ubl.GongmallCreateTime = gongmallCreateTime
}

func (ubl *UserBalanceLog) SetPayTime(ctx context.Context, payTime *time.Time) {
	ubl.PayTime = payTime
}

func (ubl *UserBalanceLog) SetDay(ctx context.Context, day uint32) {
	ubl.Day = day
}

func (ubl *UserBalanceLog) SetUpdateTime(ctx context.Context) {
	ubl.UpdateTime = time.Now()
}

func (ubl *UserBalanceLog) SetCreateTime(ctx context.Context) {
	ubl.CreateTime = time.Now()
}

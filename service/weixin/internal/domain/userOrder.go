package domain

import (
	"context"
	"time"
)

type UserOrder struct {
	Id                  uint64
	UserId              uint64
	Phone               string
	NickName            string
	AvatarUrl           string
	OrganizationId      uint64
	OrganizationUserId  uint64
	OrganizationTutorId uint64
	Level               uint8
	OutTradeNo          string
	TransactionId       string
	OutTransactionId    string
	Amount              float32
	PayAmount           float32
	PayTime             *time.Time
	PayStatus           uint8
	OrderType           uint8
	CreateTime          time.Time
	UpdateTime          time.Time
}

type UserOrderList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*UserOrder
}

type UserOrderRequestPayment struct {
	TimeStamp  string
	NonceStr   string
	Package    string
	SignType   string
	PaySign    string
	OutTradeNo string
	PayAmount  string
	LevelName  string
}

func NewUserOrder(ctx context.Context, userId, organizationId, organizationUserId, organizationTutorId uint64, outTradeNo, transactionId, outTransactionId string, amount, payAmount float32, payTime *time.Time, level, payStatus, orderType uint8) *UserOrder {
	return &UserOrder{
		UserId:              userId,
		OrganizationId:      organizationId,
		OrganizationUserId:  organizationUserId,
		OrganizationTutorId: organizationTutorId,
		Level:               level,
		OutTradeNo:          outTradeNo,
		TransactionId:       transactionId,
		OutTransactionId:    outTransactionId,
		Amount:              amount,
		PayAmount:           payAmount,
		PayTime:             payTime,
		PayStatus:           payStatus,
		OrderType:           orderType,
	}
}

func (uo *UserOrder) SetUserId(ctx context.Context, userId uint64) {
	uo.UserId = userId
}

func (uo *UserOrder) SetOrganizationId(ctx context.Context, organizationId uint64) {
	uo.OrganizationId = organizationId
}

func (uo *UserOrder) SetOrganizationUserId(ctx context.Context, organizationUserId uint64) {
	uo.OrganizationUserId = organizationUserId
}

func (uo *UserOrder) SetOrganizationTutorId(ctx context.Context, organizationTutorId uint64) {
	uo.OrganizationTutorId = organizationTutorId
}

func (uo *UserOrder) SetLevel(ctx context.Context, level uint8) {
	uo.Level = level
}

func (uo *UserOrder) SetOutTradeNo(ctx context.Context, outTradeNo string) {
	uo.OutTradeNo = outTradeNo
}

func (uo *UserOrder) SetTransactionId(ctx context.Context, transactionId string) {
	uo.TransactionId = transactionId
}

func (uo *UserOrder) SetOutTransactionId(ctx context.Context, outTransactionId string) {
	uo.OutTransactionId = outTransactionId
}

func (uo *UserOrder) SetAmount(ctx context.Context, amount float32) {
	uo.Amount = amount
}

func (uo *UserOrder) SetPayAmount(ctx context.Context, payAmount float32) {
	uo.PayAmount = payAmount
}

func (uo *UserOrder) SetPayTime(ctx context.Context, payTime *time.Time) {
	uo.PayTime = payTime
}

func (uo *UserOrder) SetPayStatus(ctx context.Context, payStatus uint8) {
	uo.PayStatus = payStatus
}

func (uo *UserOrder) SetOrderType(ctx context.Context, orderType uint8) {
	uo.OrderType = orderType
}

func (uo *UserOrder) SetUpdateTime(ctx context.Context) {
	uo.UpdateTime = time.Now()
}

func (uo *UserOrder) SetCreateTime(ctx context.Context) {
	uo.CreateTime = time.Now()
}

package domain

import (
	"context"
	"time"
)

type UserBank struct {
	Id               uint64
	UserId           uint64
	IdentityCardMark string
	BankCode         string
	BankName         string
	CreateTime       time.Time
	UpdateTime       time.Time
}

type UserBankList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*UserBank
}

func NewUserBank(ctx context.Context, userId uint64, identityCardMark, bankCode, bankName string) *UserBank {
	return &UserBank{
		UserId:           userId,
		IdentityCardMark: identityCardMark,
		BankCode:         bankCode,
		BankName:         bankName,
	}
}

func (ub *UserBank) SetUserId(ctx context.Context, userId uint64) {
	ub.UserId = userId
}

func (ub *UserBank) SetIdentityCardMark(ctx context.Context, identityCardMark string) {
	ub.IdentityCardMark = identityCardMark
}

func (ub *UserBank) SetBankCode(ctx context.Context, bankCode string) {
	ub.BankCode = bankCode
}

func (ub *UserBank) SetBankName(ctx context.Context, bankName string) {
	ub.BankName = bankName
}

func (ub *UserBank) SetUpdateTime(ctx context.Context) {
	ub.UpdateTime = time.Now()
}

func (ub *UserBank) SetCreateTime(ctx context.Context) {
	ub.CreateTime = time.Now()
}

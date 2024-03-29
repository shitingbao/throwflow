package domain

import (
	"context"
	"time"
)

type OceanengineAccount struct {
	Id          uint64
	AppId       string
	CompanyId   uint64
	AccountId   uint64
	AccountName string
	AccountRole string
	IsValid     uint8
	CreateTime  time.Time
	UpdateTime  time.Time
}

func NewOceanengineAccount(ctx context.Context, appId, accountName, accountRole string, companyId, accountId uint64, isValid uint8) *OceanengineAccount {
	return &OceanengineAccount{
		AppId:       appId,
		CompanyId:   companyId,
		AccountId:   accountId,
		AccountName: accountName,
		AccountRole: accountRole,
		IsValid:     isValid,
	}
}

func (oa *OceanengineAccount) SetAppId(ctx context.Context, appId string) {
	oa.AppId = appId
}

func (oa *OceanengineAccount) SetAccountName(ctx context.Context, accountName string) {
	oa.AccountName = accountName
}

func (oa *OceanengineAccount) SetAccountRole(ctx context.Context, accountRole string) {
	oa.AccountRole = accountRole
}

func (oa *OceanengineAccount) SetCompanyId(ctx context.Context, companyId uint64) {
	oa.CompanyId = companyId
}

func (oa *OceanengineAccount) SetAccountId(ctx context.Context, accountId uint64) {
	oa.AccountId = accountId
}

func (oa *OceanengineAccount) SetIsValid(ctx context.Context, isValid uint8) {
	oa.IsValid = isValid
}

func (oa *OceanengineAccount) SetUpdateTime(ctx context.Context) {
	oa.UpdateTime = time.Now()
}

func (oa *OceanengineAccount) SetCreateTime(ctx context.Context) {
	oa.CreateTime = time.Now()
}

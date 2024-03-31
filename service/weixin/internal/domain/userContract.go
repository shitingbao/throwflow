package domain

import (
	"context"
	"time"
)

type UserContract struct {
	Id               uint64
	OrganizationId   uint64
	Name             string
	IdentityCard     string
	IdentityCardMark string
	ServiceId        uint64
	TemplateId       uint64
	ContractId       uint64
	ContractStatus   uint8
	ContractType     uint8
	CreateTime       time.Time
	UpdateTime       time.Time
}

func NewUserContract(ctx context.Context, organizationId, serviceId, templateId, contractId uint64, contractStatus, contractType uint8, name, identityCard, identityCardMark string) *UserContract {
	return &UserContract{
		OrganizationId:   organizationId,
		Name:             name,
		IdentityCard:     identityCard,
		IdentityCardMark: identityCardMark,
		ServiceId:        serviceId,
		TemplateId:       templateId,
		ContractId:       contractId,
		ContractStatus:   contractStatus,
		ContractType:     contractType,
	}
}

func (uc *UserContract) SetOrganizationId(ctx context.Context, organizationId uint64) {
	uc.OrganizationId = organizationId
}

func (uc *UserContract) SetName(ctx context.Context, name string) {
	uc.Name = name
}

func (uc *UserContract) SetIdentityCard(ctx context.Context, identityCard string) {
	uc.IdentityCard = identityCard
}

func (uc *UserContract) SetIdentityCardMark(ctx context.Context, identityCardMark string) {
	uc.IdentityCardMark = identityCardMark
}

func (uc *UserContract) SetServiceId(ctx context.Context, serviceId uint64) {
	uc.ServiceId = serviceId
}

func (uc *UserContract) SetTemplateId(ctx context.Context, templateId uint64) {
	uc.TemplateId = templateId
}

func (uc *UserContract) SetContractId(ctx context.Context, contractId uint64) {
	uc.ContractId = contractId
}

func (uc *UserContract) SetContractStatus(ctx context.Context, contractStatus uint8) {
	uc.ContractStatus = contractStatus
}

func (uc *UserContract) SetContractType(ctx context.Context, contractType uint8) {
	uc.ContractType = contractType
}

func (uc *UserContract) SetUpdateTime(ctx context.Context) {
	uc.UpdateTime = time.Now()
}

func (uc *UserContract) SetCreateTime(ctx context.Context) {
	uc.CreateTime = time.Now()
}

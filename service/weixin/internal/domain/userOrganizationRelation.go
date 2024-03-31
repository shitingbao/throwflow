package domain

import (
	"context"
	"time"
)

type UserOrganizationRelation struct {
	Id                        uint64
	UserId                    uint64
	OrganizationId            uint64
	OrganizationUserId        uint64
	OrganizationTutorId       uint64
	OrganizationUserQrCodeUrl string
	Level                     uint8
	IsOrderRelation           uint8
	CreateTime                time.Time
	UpdateTime                time.Time
}

type UserOrganizationRelationInfo struct {
	OrganizationId            uint64
	OrganizationName          string
	OrganizationLogoUrl       string
	OrganizationCourses       OrganizationCourses
	CompanyName               string
	BankCode                  string
	BankDeposit               string
	ActivationTime            time.Time
	LevelName                 string
	Level                     uint8
	OrganizationUserQrCodeUrl string
	ParentUserId              uint64
	ParentNickName            string
	ParentAvatarUrl           string
	Total                     uint64
}

type ParentUserOrganizationRelation struct {
	ParentUserId    uint64
	ParentNickName  string
	ParentAvatarUrl string
	ParentUserType  string
}

type Mcn struct {
	Name          string
	BindStartTime string
	BindEndTime   string
}

type BindUserOrganizationRelationInfo struct {
	OrganizationId uint64
	ParentNickName string
	TutorId        uint64
	TutorNickName  string
	CreateTime     time.Time
	Mcn            []*Mcn
}

type StatisticsUserOrganizationRelation struct {
	Key   string
	Value string
}

type StatisticsUserOrganizationRelations struct {
	Statistics []*StatisticsUserOrganizationRelation
}

func NewUserOrganizationRelation(ctx context.Context, userId, organizationId, organizationUserId, organizationTutorId uint64, level, isOrderRelation uint8, organizationUserQrCodeUrl string) *UserOrganizationRelation {
	return &UserOrganizationRelation{
		UserId:                    userId,
		OrganizationId:            organizationId,
		OrganizationUserId:        organizationUserId,
		OrganizationTutorId:       organizationTutorId,
		OrganizationUserQrCodeUrl: organizationUserQrCodeUrl,
		Level:                     level,
		IsOrderRelation:           isOrderRelation,
	}
}

func (uor *UserOrganizationRelation) SetUserId(ctx context.Context, userId uint64) {
	uor.UserId = userId
}

func (uor *UserOrganizationRelation) SetOrganizationId(ctx context.Context, organizationId uint64) {
	uor.OrganizationId = organizationId
}

func (uor *UserOrganizationRelation) SetOrganizationUserId(ctx context.Context, organizationUserId uint64) {
	uor.OrganizationUserId = organizationUserId
}

func (uor *UserOrganizationRelation) SetOrganizationTutorId(ctx context.Context, organizationTutorId uint64) {
	uor.OrganizationTutorId = organizationTutorId
}

func (uor *UserOrganizationRelation) SetOrganizationUserQrCodeUrl(ctx context.Context, organizationUserQrCodeUrl string) {
	uor.OrganizationUserQrCodeUrl = organizationUserQrCodeUrl
}

func (uor *UserOrganizationRelation) SetLevel(ctx context.Context, level uint8) {
	uor.Level = level
}

func (uor *UserOrganizationRelation) SetIsOrderRelation(ctx context.Context, isOrderRelation uint8) {
	uor.IsOrderRelation = isOrderRelation
}

func (uor *UserOrganizationRelation) SetUpdateTime(ctx context.Context) {
	uor.UpdateTime = time.Now()
}

func (uor *UserOrganizationRelation) SetCreateTime(ctx context.Context) {
	uor.CreateTime = time.Now()
}

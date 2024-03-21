package domain

import (
	"context"
	"encoding/json"
	"regexp"
	"time"
)

type ContactInformation struct {
	ContactUsername string `json:"contactUsername"`
	ContactPosition string `json:"contactPosition"`
	ContactPhone    string `json:"contactPhone"`
	ContactWeixin   string `json:"contactWeixin"`
}

type ContactInformations []ContactInformation

type OperationLog struct {
	UserId     uint64 `json:"userId"`
	UserName   string `json:"userName"`
	Content    string `json:"content"`
	CreateTime string `json:"createTime"`
}

type OperationLogs []*OperationLog

type Clue struct {
	Id                  uint64
	CompanyName         string
	IndustryId          string
	IndustryName        string
	ContactInformation  string
	ContactInformations ContactInformations
	CompanyType         uint8
	CompanyTypeName     string
	QianchuanUse        uint8
	QianchuanUseName    string
	Sale                uint64
	Seller              string
	Facilitator         string
	Source              string
	Status              uint8
	StatusName          string
	OperationLog        string
	AreaCode            uint64
	AreaName            string
	Address             string
	IsAffiliates        uint8
	AdminName           string
	AdminPhone          string
	IsDel               uint8
	OperationLogs       OperationLogs
	CreateTime          time.Time
	UpdateTime          time.Time
}

type ClueList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*Clue
}

type Status struct {
	Key   string
	Value string
}

type CompanyType struct {
	Key   string
	Value string
}

type QianchuanUse struct {
	Key   string
	Value string
}

type SelectClues struct {
	Status       []*Status
	CompanyType  []*CompanyType
	QianchuanUse []*QianchuanUse
}

type StatisticsClue struct {
	Key   string
	Value string
}

type StatisticsClues struct {
	Statistics []*StatisticsClue
}

func NewSelectClues() *SelectClues {
	status := make([]*Status, 0)
	companyType := make([]*CompanyType, 0)
	qianchuanUse := make([]*QianchuanUse, 0)

	status = append(status, &Status{Key: "1", Value: "公海"})
	status = append(status, &Status{Key: "2", Value: "洽谈"})
	status = append(status, &Status{Key: "3", Value: "过期"})
	status = append(status, &Status{Key: "4", Value: "正式"})
	status = append(status, &Status{Key: "5", Value: "测试"})

	companyType = append(companyType, &CompanyType{Key: "1", Value: "服务商"})
	companyType = append(companyType, &CompanyType{Key: "2", Value: "品牌商"})
	companyType = append(companyType, &CompanyType{Key: "3", Value: "团长"})

	qianchuanUse = append(qianchuanUse, &QianchuanUse{Key: "1", Value: "未使用"})
	qianchuanUse = append(qianchuanUse, &QianchuanUse{Key: "2", Value: "0<消耗≤10万/月"})
	qianchuanUse = append(qianchuanUse, &QianchuanUse{Key: "3", Value: "10万/月<消耗≤30万/月"})
	qianchuanUse = append(qianchuanUse, &QianchuanUse{Key: "4", Value: "30万/月<消耗≤100万/月"})
	qianchuanUse = append(qianchuanUse, &QianchuanUse{Key: "5", Value: "消耗>100万/月"})

	return &SelectClues{
		Status:       status,
		CompanyType:  companyType,
		QianchuanUse: qianchuanUse,
	}
}

func NewClue(ctx context.Context, companyName, contactInformation, source, seller, facilitator, address, industryId string, areaCode uint64, companyType, qianchuanUse, status uint8) *Clue {
	return &Clue{
		CompanyName:        companyName,
		IndustryId:         industryId,
		AreaCode:           areaCode,
		ContactInformation: contactInformation,
		CompanyType:        companyType,
		QianchuanUse:       qianchuanUse,
		Source:             source,
		Status:             status,
		Seller:             seller,
		Facilitator:        facilitator,
		Address:            address,
	}
}

func (c *Clue) SetCompanyName(ctx context.Context, companyName string) {
	c.CompanyName = companyName
}

func (c *Clue) SetCompanyType(ctx context.Context, companyType uint8) {
	c.CompanyType = companyType
}

func (c *Clue) SetCompanyTypeName(ctx context.Context) {
	switch c.CompanyType {
	case 1:
		c.CompanyTypeName = "服务商"
	case 2:
		c.CompanyTypeName = "品牌商"
	case 3:
		c.CompanyTypeName = "团长"
	}

	return
}

func (c *Clue) SetQianchuanUse(ctx context.Context, qianchuanUse uint8) {
	c.QianchuanUse = qianchuanUse
}

func (c *Clue) SetSeller(ctx context.Context, seller string) {
	c.Seller = seller
}

func (c *Clue) SetFacilitator(ctx context.Context, facilitator string) {
	c.Facilitator = facilitator
}

func (c *Clue) SetAreaCode(ctx context.Context, areaCode uint64) {
	c.AreaCode = areaCode
}

func (c *Clue) SetAreaName(ctx context.Context, areaName string) {
	c.AreaName = areaName
}

func (c *Clue) SetIndustryId(ctx context.Context, industryId string) {
	c.IndustryId = industryId
}

func (c *Clue) SetIndustryName(ctx context.Context, industryName string) {
	c.IndustryName = industryName
}

func (c *Clue) SetContactInformation(ctx context.Context, contactInformation string) {
	c.ContactInformation = contactInformation
}

func (c *Clue) SetAddress(ctx context.Context, address string) {
	c.Address = address
}

func (c *Clue) SetStatus(ctx context.Context, status uint8) {
	c.Status = status
}

func (c *Clue) SetIsDel(ctx context.Context, isDel uint8) {
	c.IsDel = isDel
}

func (c *Clue) SetUpdateTime(ctx context.Context) {
	c.UpdateTime = time.Now()
}

func (c *Clue) SetCreateTime(ctx context.Context) {
	c.CreateTime = time.Now()
}

func (c *Clue) SetOperationLog(ctx context.Context, userId uint64, userName, content, operationTime string) {
	var operationLogs OperationLogs

	if len(c.OperationLog) == 0 {
		operationLogs = make(OperationLogs, 0)
	} else {
		if err := json.Unmarshal([]byte(c.OperationLog), &operationLogs); err != nil {
			return
		}
	}

	operationLogs = append(operationLogs, &OperationLog{
		UserId:     userId,
		UserName:   userName,
		Content:    content,
		CreateTime: operationTime,
	})

	if ols, err := json.Marshal(operationLogs); err == nil {
		c.OperationLog = string(ols)
	}
}

func (c *Clue) GetStatusName(ctx context.Context) {
	switch c.Status {
	case 1:
		c.StatusName = "公海"
	case 2:
		c.StatusName = "洽谈"
	case 3:
		c.StatusName = "过期"
	case 4:
		c.StatusName = "正式"
	case 5:
		c.StatusName = "测试"
	}
}

func (c *Clue) GetOperationLog(ctx context.Context) {
	var operationLogs OperationLogs

	if err := json.Unmarshal([]byte(c.OperationLog), &operationLogs); err == nil {
		c.OperationLogs = operationLogs
	}
}

func (c *Clue) GetContactInformation(ctx context.Context) {
	var contactInformations ContactInformations

	if err := json.Unmarshal([]byte(c.ContactInformation), &contactInformations); err == nil {
		c.ContactInformations = contactInformations
	}
}

func (c *Clue) VerifyContactInformation(ctx context.Context) bool {
	var contactInformations ContactInformations

	if err := json.Unmarshal([]byte(c.ContactInformation), &contactInformations); err != nil {
		return false
	}

	for _, contactInformation := range contactInformations {
		if !regexp.MustCompile("^((\\+|00)86)?(1[3-9]|9[28])\\d{9}$").MatchString(contactInformation.ContactPhone) {
			return false
		}
	}

	return true
}

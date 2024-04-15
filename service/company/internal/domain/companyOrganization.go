package domain

import (
	"company/internal/pkg/tool"
	"context"
	"encoding/json"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

type OrganizationCommission struct {
	CostOrderCommissionRatio float32 `json:"costOrderCommissionRatio"`
	OrderCommissionRatio     float32 `json:"orderCommissionRatio"`
}

type OrganizationColonelCommission struct {
	ZeroCourseRatio                                               float32 `json:"zeroCourseRatio"`
	ZeroAdvancedPresenterZeroCourseCommissionRule                 float32 `json:"zeroAdvancedPresenterZeroCourseCommissionRule"`
	ZeroAdvancedTutorZeroCourseCommissionRule                     float32 `json:"zeroAdvancedTutorZeroCourseCommissionRule"`
	PrimaryAdvancedPresenterZeroCourseCommissionRule              float32 `json:"primaryAdvancedPresenterZeroCourseCommissionRule"`
	PrimaryAdvancedTutorZeroCourseCommissionRule                  float32 `json:"primaryAdvancedTutorZeroCourseCommissionRule"`
	IntermediateAdvancedPresenterZeroCourseCommissionRule         float32 `json:"intermediateAdvancedPresenterZeroCourseCommissionRule"`
	IntermediateAdvancedTutorZeroCourseCommissionRule             float32 `json:"intermediateAdvancedTutorZeroCourseCommissionRule"`
	AdvancedPresenterZeroCourseCommissionRule                     float32 `json:"advancedPresenterZeroCourseCommissionRule"`
	PrimaryCourseRatio                                            float32 `json:"primaryCourseRatio"`
	PrimaryAdvancedPresenterPrimaryCourseCommissionRule           float32 `json:"primaryAdvancedPresenterPrimaryCourseCommissionRule"`
	PrimaryAdvancedTutorPrimaryCourseCommissionRule               float32 `json:"primaryAdvancedTutorPrimaryCourseCommissionRule"`
	IntermediateAdvancedPresenterPrimaryCourseCommissionRule      float32 `json:"intermediateAdvancedPresenterPrimaryCourseCommissionRule"`
	IntermediateAdvancedTutorPrimaryCourseCommissionRule          float32 `json:"intermediateAdvancedTutorPrimaryCourseCommissionRule"`
	AdvancedPresenterPrimaryCourseCommissionRule                  float32 `json:"advancedPresenterPrimaryCourseCommissionRule"`
	IntermediateCourseRatio                                       float32 `json:"intermediateCourseRatio"`
	PrimaryAdvancedPresenterIntermediateCourseCommissionRule      float32 `json:"primaryAdvancedPresenterIntermediateCourseCommissionRule"`
	PrimaryAdvancedTutorIntermediateCourseCommissionRule          float32 `json:"primaryAdvancedTutorIntermediateCourseCommissionRule"`
	IntermediateAdvancedPresenterIntermediateCourseCommissionRule float32 `json:"intermediateAdvancedPresenterIntermediateCourseCommissionRule"`
	IntermediateAdvancedTutorIntermediateCourseCommissionRule     float32 `json:"intermediateAdvancedTutorIntermediateCourseCommissionRule"`
	AdvancedPresenterIntermediateCourseCommissionRule             float32 `json:"advancedPresenterIntermediateCourseCommissionRule"`
	AdvancedCourseRatio                                           float32 `json:"advancedCourseRatio"`
	PrimaryAdvancedPresenterAdvancedCourseCommissionRule          float32 `json:"primaryAdvancedPresenterAdvancedCourseCommissionRule"`
	PrimaryAdvancedTutorAdvancedCourseCommissionRule              float32 `json:"primaryAdvancedTutorAdvancedCourseCommissionRule"`
	IntermediateAdvancedPresenterAdvancedCourseCommissionRule     float32 `json:"intermediateAdvancedPresenterAdvancedCourseCommissionRule"`
	IntermediateAdvancedTutorAdvancedCourseCommissionRule         float32 `json:"intermediateAdvancedTutorAdvancedCourseCommissionRule"`
	AdvancedPresenterAdvancedCourseCommissionRule                 float32 `json:"advancedPresenterAdvancedCourseCommissionRule"`
	OrderRatio                                                    float32 `json:"orderRatio"`
	PrimaryAdvancedPresenterOrderCommissionRule                   float32 `json:"primaryAdvancedPresenterOrderCommissionRule"`
	PrimaryAdvancedTutorOrderCommissionRule                       float32 `json:"primaryAdvancedTutorOrderCommissionRule"`
	IntermediateAdvancedPresenterOrderCommissionRule              float32 `json:"intermediateAdvancedPresenterOrderCommissionRule"`
	IntermediateAdvancedTutorOrderCommissionRule                  float32 `json:"intermediateAdvancedTutorOrderCommissionRule"`
	AdvancedPresenterOrderCommissionRule                          float32 `json:"advancedPresenterOrderCommissionRule"`
	CostOrderRatio                                                float32 `json:"costOrderRatio"`
	ZeroAdvancedPresenterCostOrderCommissionRule                  float32 `json:"zeroAdvancedPresenterCostOrderCommissionRule"`
	ZeroAdvancedTutorCostOrderCommissionRule                      float32 `json:"zeroAdvancedTutorCostOrderCommissionRule"`
	PrimaryAdvancedPresenterCostOrderCommissionRule               float32 `json:"primaryAdvancedPresenterCostOrderCommissionRule"`
	PrimaryAdvancedTutorCostOrderCommissionRule                   float32 `json:"primaryAdvancedTutorCostOrderCommissionRule"`
	IntermediateAdvancedPresenterCostOrderCommissionRule          float32 `json:"intermediateAdvancedPresenterCostOrderCommissionRule"`
	IntermediateAdvancedTutorCostOrderCommissionRule              float32 `json:"intermediateAdvancedTutorCostOrderCommissionRule"`
	AdvancedPresenterCostOrderCommissionRule                      float32 `json:"advancedPresenterCostOrderCommissionRule"`
}

type OrganizationCourses []OrganizationCourse

type OrganizationCourseModule struct {
	CourseModuleName    string `json:"courseModuleName"`
	CourseModuleContent string `json:"courseModuleContent"`
}

type OrganizationCourse struct {
	CourseName          string                     `json:"courseName"`
	CourseSubName       string                     `json:"courseSubName"`
	CoursePrice         float32                    `json:"coursePrice"`
	CourseDuration      uint64                     `json:"courseDuration"`
	CourseOriginalPrice float32                    `json:"courseOriginalPrice"`
	CourseLevel         uint8                      `json:"courseLevel"`
	CourseModules       []OrganizationCourseModule `json:"courseModule"`
}

type OrganizationUser struct {
	UserId   uint64 `json:"userId"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
}

type CompanyOrganization struct {
	Id                             uint64
	OrganizationName               string
	OrganizationMcn                string
	OrganizationMcns               []string
	CompanyName                    string
	BankCode                       string
	BankDeposit                    string
	OrganizationLogoUrl            string
	OrganizationCode               string
	OrganizationQrCodeUrl          string
	OrganizationShortUrl           string
	OrganizationCommission         string
	OrganizationCommissions        OrganizationCommission
	OrganizationColonelCommission  string
	OrganizationColonelCommissions OrganizationColonelCommission
	OrganizationCourse             string
	OrganizationCourses            []*OrganizationCourse
	OrganizationUser               string
	OrganizationUsers              []*OrganizationUser
	CreateTime                     time.Time
	UpdateTime                     time.Time
}

type CompanyOrganizationList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*CompanyOrganization
}

type CourseLevel struct {
	Key   string
	Value string
}

type OrganizationMcn struct {
	Key   string
	Value string
}

type SelectCompanyOrganizations struct {
	CourseLevel     []*CourseLevel
	OrganizationMcn []*OrganizationMcn
}

func NewSelectCompanyOrganizations() *SelectCompanyOrganizations {
	courseLevel := make([]*CourseLevel, 0)
	organizationMcn := make([]*OrganizationMcn, 0)

	return &SelectCompanyOrganizations{
		CourseLevel:     courseLevel,
		OrganizationMcn: organizationMcn,
	}
}

func NewCompanyOrganization(ctx context.Context, organizationName, organizationMcn, companyName, bankCode, bankDeposit, organizationLogoUrl, organizationCode, organizationQrCodeUrl, organizationShortUrl, organizationCommission, organizationColonelCommission, organizationCourse string) *CompanyOrganization {
	return &CompanyOrganization{
		OrganizationName:              organizationName,
		OrganizationMcn:               organizationMcn,
		CompanyName:                   companyName,
		BankCode:                      bankCode,
		BankDeposit:                   bankDeposit,
		OrganizationLogoUrl:           organizationLogoUrl,
		OrganizationCode:              organizationCode,
		OrganizationQrCodeUrl:         organizationQrCodeUrl,
		OrganizationShortUrl:          organizationShortUrl,
		OrganizationCommission:        organizationCommission,
		OrganizationColonelCommission: organizationColonelCommission,
		OrganizationCourse:            organizationCourse,
	}
}

func (co *CompanyOrganization) SetOrganizationName(ctx context.Context, organizationName string) {
	co.OrganizationName = organizationName
}

func (co *CompanyOrganization) SetOrganizationMcn(ctx context.Context, organizationMcn string) {
	co.OrganizationMcn = organizationMcn
}

func (co *CompanyOrganization) SetOrganizationMcns(ctx context.Context) {
	co.OrganizationMcns = make([]string, 0)

	if len(co.OrganizationMcn) > 0 {
		co.OrganizationMcns = strings.Split(co.OrganizationMcn, ",")
	}
}

func (co *CompanyOrganization) SetCompanyName(ctx context.Context, companyName string) {
	co.CompanyName = companyName
}

func (co *CompanyOrganization) SetBankCode(ctx context.Context, bankCode string) {
	co.BankCode = bankCode
}

func (co *CompanyOrganization) SetBankDeposit(ctx context.Context, bankDeposit string) {
	co.BankDeposit = bankDeposit
}

func (co *CompanyOrganization) SetOrganizationLogoUrl(ctx context.Context, organizationLogoUrl string) {
	co.OrganizationLogoUrl = organizationLogoUrl
}

func (co *CompanyOrganization) SetOrganizationCode(ctx context.Context, organizationCode string) {
	co.OrganizationCode = organizationCode
}

func (co *CompanyOrganization) SetOrganizationQrCodeUrl(ctx context.Context, organizationQrCodeUrl string) {
	co.OrganizationQrCodeUrl = organizationQrCodeUrl
}

func (co *CompanyOrganization) SetOrganizationShortUrl(ctx context.Context, organizationShortUrl string) {
	co.OrganizationShortUrl = organizationShortUrl
}

func (co *CompanyOrganization) SetOrganizationCommission(ctx context.Context, organizationCommission string) {
	co.OrganizationCommission = organizationCommission
}

func (co *CompanyOrganization) SetOrganizationCommissions(ctx context.Context) {
	var organizationCommission OrganizationCommission

	if err := json.Unmarshal([]byte(co.OrganizationCommission), &organizationCommission); err == nil {
		co.OrganizationCommissions = organizationCommission
	}
}

func (co *CompanyOrganization) SetOrganizationColonelCommission(ctx context.Context, organizationColonelCommission string) {
	co.OrganizationColonelCommission = organizationColonelCommission
}

func (co *CompanyOrganization) SetOrganizationColonelCommissions(ctx context.Context) {
	var organizationColonelCommission OrganizationColonelCommission

	if err := json.Unmarshal([]byte(co.OrganizationColonelCommission), &organizationColonelCommission); err == nil {
		co.OrganizationColonelCommissions = organizationColonelCommission
	}
}

func (co *CompanyOrganization) SetOrganizationCourse(ctx context.Context, organizationCourse string) {
	co.OrganizationCourse = organizationCourse
}

func (co *CompanyOrganization) SetOrganizationCourses(ctx context.Context) {
	var organizationCourses []*OrganizationCourse

	if err := json.Unmarshal([]byte(co.OrganizationCourse), &organizationCourses); err == nil {
		co.OrganizationCourses = organizationCourses
	}
}

func (co *CompanyOrganization) SetOrganizationUser(ctx context.Context, organizationUser string) {
	co.OrganizationUser = organizationUser
}

func (co *CompanyOrganization) SetUpdateTime(ctx context.Context) {
	co.UpdateTime = time.Now()
}

func (co *CompanyOrganization) SetCreateTime(ctx context.Context) {
	co.CreateTime = time.Now()
}

func (co *CompanyOrganization) VerifyOrganizationMcn(ctx context.Context) bool {
	mcns := [4]string{"星达小当家", "小当家", "莲宝兔旺", "壹玖传媒"}

	organizationMcns := make([]string, 0)

	tmpOrganizationMcns := tool.RemoveEmptyString(strings.Split(co.OrganizationMcn, ","))

	for _, tmpOrganizationMcn := range tmpOrganizationMcns {
		for _, mcn := range mcns {
			if mcn == tmpOrganizationMcn {
				isNotExist := true

				for _, organizationMcn := range organizationMcns {
					if organizationMcn == mcn {
						isNotExist = false

						break
					}
				}

				if isNotExist {
					organizationMcns = append(organizationMcns, mcn)
				}

				break
			}
		}
	}

	co.OrganizationMcns = organizationMcns
	co.OrganizationMcn = strings.Join(organizationMcns, ",")

	return true
}

func (co *CompanyOrganization) VerifyOrganizationUser(ctx context.Context) bool {
	var organizationUsers []*OrganizationUser

	if err := json.Unmarshal([]byte(co.OrganizationUser), &organizationUsers); err != nil {
		return false
	}

	phones := make([]string, 0)

	for _, organizationUser := range organizationUsers {
		if !regexp.MustCompile("^((\\+|00)86)?(1[3-9]|9[28])\\d{9}$").MatchString(organizationUser.Phone) {
			return false
		}

		if l := utf8.RuneCountInString(organizationUser.Username); l == 0 {
			return false
		}

		isExist := false

		for _, phone := range phones {
			if phone == organizationUser.Phone {
				isExist = true

				break
			}
		}

		if isExist {
			return false
		} else {
			phones = append(phones, organizationUser.Phone)
		}
	}

	co.OrganizationUsers = organizationUsers

	return true
}

func (co *CompanyOrganization) VerifyOrganizationCommission(ctx context.Context) bool {
	var organizationCommission OrganizationCommission

	if err := json.Unmarshal([]byte(co.OrganizationCommission), &organizationCommission); err != nil {
		return false
	}

	if organizationCommission.CostOrderCommissionRatio <= 0 {
		return false
	}

	if organizationCommission.CostOrderCommissionRatio > 100 {
		return false
	}

	if organizationCommission.OrderCommissionRatio <= 0 {
		return false
	}

	if organizationCommission.OrderCommissionRatio > 100 {
		return false
	}

	return true
}

func (co *CompanyOrganization) VerifyOrganizationColonelCommission(ctx context.Context) bool {
	var organizationColonelCommission OrganizationColonelCommission

	if err := json.Unmarshal([]byte(co.OrganizationColonelCommission), &organizationColonelCommission); err != nil {
		return false
	}

	if organizationColonelCommission.ZeroCourseRatio < 0 {
		return false
	}

	if organizationColonelCommission.ZeroAdvancedPresenterZeroCourseCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.ZeroAdvancedTutorZeroCourseCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.PrimaryAdvancedPresenterZeroCourseCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.PrimaryAdvancedTutorZeroCourseCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.IntermediateAdvancedPresenterZeroCourseCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.IntermediateAdvancedTutorZeroCourseCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.AdvancedPresenterZeroCourseCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.PrimaryCourseRatio < 0 {
		return false
	}

	if organizationColonelCommission.PrimaryAdvancedPresenterPrimaryCourseCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.PrimaryAdvancedTutorPrimaryCourseCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.IntermediateAdvancedPresenterPrimaryCourseCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.IntermediateAdvancedTutorPrimaryCourseCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.AdvancedPresenterPrimaryCourseCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.IntermediateCourseRatio < 0 {
		return false
	}

	if organizationColonelCommission.PrimaryAdvancedPresenterIntermediateCourseCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.PrimaryAdvancedTutorIntermediateCourseCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.IntermediateAdvancedPresenterIntermediateCourseCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.IntermediateAdvancedTutorIntermediateCourseCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.AdvancedPresenterIntermediateCourseCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.AdvancedCourseRatio < 0 {
		return false
	}

	if organizationColonelCommission.PrimaryAdvancedPresenterAdvancedCourseCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.PrimaryAdvancedTutorAdvancedCourseCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.IntermediateAdvancedPresenterAdvancedCourseCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.IntermediateAdvancedTutorAdvancedCourseCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.AdvancedPresenterAdvancedCourseCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.OrderRatio < 0 {
		return false
	}

	if organizationColonelCommission.PrimaryAdvancedPresenterOrderCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.PrimaryAdvancedTutorOrderCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.IntermediateAdvancedPresenterOrderCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.IntermediateAdvancedTutorOrderCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.AdvancedPresenterOrderCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.CostOrderRatio < 0 {
		return false
	}

	if organizationColonelCommission.ZeroAdvancedPresenterCostOrderCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.ZeroAdvancedTutorCostOrderCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.PrimaryAdvancedPresenterCostOrderCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.PrimaryAdvancedTutorCostOrderCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.IntermediateAdvancedPresenterCostOrderCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.IntermediateAdvancedTutorCostOrderCommissionRule < 0 {
		return false
	}

	if organizationColonelCommission.AdvancedPresenterCostOrderCommissionRule < 0 {
		return false
	}

	return true
}

func (co *CompanyOrganization) VerifyOrganizationCourse(ctx context.Context) bool {
	var organizationCourses []*OrganizationCourse

	if err := json.Unmarshal([]byte(co.OrganizationCourse), &organizationCourses); err != nil {
		return false
	}

	if len(organizationCourses) < 1 {
		return false
	}

	tmpOrganizationCourses := make([]OrganizationCourse, 0)

	for _, organizationCourse := range organizationCourses {
		if organizationCourse.CoursePrice <= 0 {
			return false
		}

		if organizationCourse.CourseDuration < 0 {
			return false
		}

		if organizationCourse.CourseOriginalPrice < 0 {
			return false
		}

		if organizationCourse.CourseLevel <= 0 {
			return false
		}

		if l := utf8.RuneCountInString(organizationCourse.CourseName); l == 0 {
			return false
		}

		if l := utf8.RuneCountInString(organizationCourse.CourseSubName); l == 0 {
			return false
		}

		for _, courseModule := range organizationCourse.CourseModules {
			if l := utf8.RuneCountInString(courseModule.CourseModuleName); l == 0 {
				return false
			}

			if l := utf8.RuneCountInString(courseModule.CourseModuleContent); l == 0 {
				return false
			}
		}

		tmpOrganizationCourses = append(tmpOrganizationCourses, OrganizationCourse{
			CourseName:          organizationCourse.CourseName,
			CourseSubName:       organizationCourse.CourseSubName,
			CoursePrice:         organizationCourse.CoursePrice,
			CourseDuration:      organizationCourse.CourseDuration,
			CourseOriginalPrice: organizationCourse.CourseOriginalPrice,
			CourseLevel:         organizationCourse.CourseLevel,
			CourseModules:       organizationCourse.CourseModules,
		})
	}

	sort.Sort(OrganizationCourses(tmpOrganizationCourses))

	var level uint8 = 0

	for _, tmpOrganizationCourse := range tmpOrganizationCourses {
		if tmpOrganizationCourse.CourseLevel <= level {
			return false
		} else {
			level = tmpOrganizationCourse.CourseLevel
		}
	}

	return true
}

func (o OrganizationCourses) Len() int {
	return len(o)
}

func (o OrganizationCourses) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (o OrganizationCourses) Less(i, j int) bool {
	return o[i].CoursePrice < o[j].CoursePrice
}

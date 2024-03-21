package domain

import (
	v1 "company/api/service/douyin/v1"
	"company/internal/pkg/tool"
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type CompanyPerformanceRule struct {
	Id                uint64
	CompanyId         uint64
	PerformanceName   string
	AdvertiserIds     string
	Advertisers       []*QianchuanAdvertiserList
	PerformanceRules  string
	PerformanceRulest []*PerformanceRule
	CreateTime        time.Time
	UpdateTime        time.Time
}

type Condition struct {
	Min   string `json:"min"`
	Max   string `json:"max"`
	Ctype string `json:"ctype"`
}

type Commissions struct {
	Ctype      string `json:"ctype"`
	Percentage string `json:"percentage"`
}

type PRules struct {
	Condition  Condition   `json:"condition"`
	Commission Commissions `json:"commission"`
}

type PerformanceRule struct {
	Ptype string    `json:"ptype"`
	Rules []*PRules `json:"rules"`
}

type PerformanceRules []*PerformanceRule

func NewCompanyPerformanceRule(ctx context.Context, companyId uint64, performanceName, performanceRules string) *CompanyPerformanceRule {
	return &CompanyPerformanceRule{
		CompanyId:        companyId,
		PerformanceName:  performanceName,
		PerformanceRules: performanceRules,
	}
}

func (cpr *CompanyPerformanceRule) SetCompanyId(ctx context.Context, companyId uint64) {
	cpr.CompanyId = companyId
}

func (cpr *CompanyPerformanceRule) SetPerformanceName(ctx context.Context, performanceName string) {
	cpr.PerformanceName = performanceName
}

func (cpr *CompanyPerformanceRule) SetAdvertiserIds(ctx context.Context, advertiserIds string, qianchuanAdvertisers *v1.ListQianchuanAdvertisersReply, companyPerformanceRules []*CompanyPerformanceRule) {
	if l := utf8.RuneCountInString(advertiserIds); l > 0 {
		radvertiserIds := make([]string, 0)
		qAdvertiserIds := make([]string, 0)
		eAdvertiserIds := make([]string, 0)

		sadvertiserIds := tool.RemoveEmptyString(strings.Split(advertiserIds, ","))

		for _, acompanyPerformanceRule := range companyPerformanceRules {
			if acompanyPerformanceRule.CompanyId != cpr.CompanyId {
				tAdvertiserIds := tool.RemoveEmptyString(strings.Split(acompanyPerformanceRule.AdvertiserIds, ","))

				radvertiserIds = append(radvertiserIds, tAdvertiserIds...)
			}
		}

		for _, qianchuanAdvertiser := range qianchuanAdvertisers.Data.List {
			qAdvertiserIds = append(qAdvertiserIds, strconv.FormatUint(qianchuanAdvertiser.AdvertiserId, 10))
		}

		for _, sadvertiserId := range sadvertiserIds {
			isNotExit := true

			for _, radvertiserId := range radvertiserIds {
				if sadvertiserId == radvertiserId {
					isNotExit = false
					break
				}
			}

			if isNotExit {
				isNotExit = false

				for _, qAdvertiserId := range qAdvertiserIds {
					if sadvertiserId == qAdvertiserId {
						isNotExit = true
						break
					}
				}
			}

			if isNotExit {
				eAdvertiserIds = append(eAdvertiserIds, sadvertiserId)
			}
		}

		cpr.AdvertiserIds = strings.Join(eAdvertiserIds, ",")
	} else {
		cpr.AdvertiserIds = ""
	}
}

func (cpr *CompanyPerformanceRule) GetAdvertiserIds(ctx context.Context, advertiserIds string, qianchuanAdvertisers *v1.ListQianchuanAdvertisersReply) {
	sadvertiserIds := tool.RemoveEmptyString(strings.Split(advertiserIds, ","))

	for _, sadvertiserId := range sadvertiserIds {
		for _, qianchuanAdvertiser := range qianchuanAdvertisers.Data.List {
			if strconv.FormatUint(qianchuanAdvertiser.AdvertiserId, 10) == sadvertiserId {
				cpr.Advertisers = append(cpr.Advertisers, &QianchuanAdvertiserList{
					AdvertiserId:   qianchuanAdvertiser.AdvertiserId,
					Status:         qianchuanAdvertiser.Status,
					AdvertiserName: qianchuanAdvertiser.AdvertiserName,
					CompanyName:    qianchuanAdvertiser.CompanyName,
				})

				break
			}
		}
	}
}

func (cpr *CompanyPerformanceRule) SetPerformanceRules(ctx context.Context, performanceRules string) {
	cpr.PerformanceRules = performanceRules
}

func (cpr *CompanyPerformanceRule) GetPerformanceRules(ctx context.Context) {
	var performanceRules PerformanceRules

	if err := json.Unmarshal([]byte(cpr.PerformanceRules), &performanceRules); err == nil {
		for _, performanceRule := range performanceRules {
			sort.SliceStable(performanceRule.Rules, func(i, j int) bool {
				return performanceRule.Rules[i].Condition.Min < performanceRule.Rules[j].Condition.Min
			})
		}

		cpr.PerformanceRulest = performanceRules
	}
}

func (cpr *CompanyPerformanceRule) SetUpdateTime(ctx context.Context) {
	cpr.UpdateTime = time.Now()
}

func (cpr *CompanyPerformanceRule) SetCreateTime(ctx context.Context) {
	cpr.CreateTime = time.Now()
}

func (cpr *CompanyPerformanceRule) VerifyPerformanceRules(ctx context.Context) bool {
	ptype := make(map[string]string)
	ptype["live"] = ""
	ptype["video"] = ""

	if l := utf8.RuneCountInString(cpr.PerformanceRules); l == 0 {
		return true
	}

	var performanceRules PerformanceRules

	if err := json.Unmarshal([]byte(cpr.PerformanceRules), &performanceRules); err != nil {
		return false
	}

	if len(performanceRules) > 2 {
		return false
	}

	for _, performanceRule := range performanceRules {
		if performanceRule.Ptype != "live" && performanceRule.Ptype != "video" {
			return false
		}

		if _, ok := ptype[performanceRule.Ptype]; !ok {
			return false
		}

		delete(ptype, performanceRule.Ptype)

		sort.SliceStable(performanceRule.Rules, func(i, j int) bool {
			return performanceRule.Rules[i].Condition.Min < performanceRule.Rules[j].Condition.Min
		})

		var pMin float64
		var pMax float64

		for _, rule := range performanceRule.Rules {
			if rule.Condition.Ctype != "1" {
				return false
			}

			if rule.Commission.Ctype != "1" {
				return false
			}

			cMin, err := strconv.ParseFloat(rule.Condition.Min, 10)

			if err != nil {
				return false
			}

			cMax, err := strconv.ParseFloat(rule.Condition.Max, 10)

			if err != nil {
				return false
			}

			if 0 >= cMin || cMin > 100 {
				return false
			}

			if 0 >= cMax || cMax > 100 {
				return false
			}

			if cMin >= cMax {
				return false
			}

			cPercentage, err := strconv.ParseFloat(rule.Commission.Percentage, 10)

			if err != nil {
				return false
			}

			if 0 >= cPercentage || cPercentage > 100 {
				return false
			}

			if cMin <= pMin {
				return false
			}

			if pMax > cMin {
				return false
			}

			pMin = cMin
			pMax = cMax
		}
	}

	return true
}

package domain

import (
	"context"
	"encoding/json"
	"time"
)

type SetValueSampleThreshold struct {
	Type  uint8  `json:"type"`
	Value uint64 `json:"value"`
}

type CompanySet struct {
	CompanyId               uint64
	Day                     uint32
	SetKey                  string
	SetValue                string
	SetValueSampleThreshold *SetValueSampleThreshold
	CreateTime              time.Time
	UpdateTime              time.Time
}

func NewCompanySet(ctx context.Context, companyId uint64, day uint32, setKey, setValue string) *CompanySet {
	return &CompanySet{
		CompanyId: companyId,
		Day:       day,
		SetKey:    setKey,
		SetValue:  setValue,
	}
}

func (cs *CompanySet) SetCompanyId(ctx context.Context, companyId uint64) {
	cs.CompanyId = companyId
}

func (cs *CompanySet) SetDay(ctx context.Context, day uint32) {
	cs.Day = day
}

func (cs *CompanySet) SetSetValue(ctx context.Context, setValue string) {
	cs.SetValue = setValue
}

func (cs *CompanySet) GetSetValue(ctx context.Context) {
	if cs.SetKey == "sampleThreshold" {
		var setValueSampleThreshold SetValueSampleThreshold

		if err := json.Unmarshal([]byte(cs.SetValue), &setValueSampleThreshold); err == nil {
			cs.SetValueSampleThreshold = &setValueSampleThreshold
		}
	}
}

func (cs *CompanySet) SetSetKey(ctx context.Context, setKey string) {
	cs.SetKey = setKey
}

func (cs *CompanySet) SetUpdateTime(ctx context.Context) {
	cs.UpdateTime = time.Now()
}

func (cs *CompanySet) SetCreateTime(ctx context.Context) {
	cs.CreateTime = time.Now()
}

func (cs *CompanySet) VerifySetValue(ctx context.Context) bool {
	if cs.SetKey == "sampleThreshold" {
		var setValueSampleThreshold SetValueSampleThreshold

		if err := json.Unmarshal([]byte(cs.SetValue), &setValueSampleThreshold); err != nil {
			return false
		}

		if setValueSampleThreshold.Type != 1 && setValueSampleThreshold.Type != 2 {
			return false
		}

		if setValueSampleThreshold.Value <= 0 {
			return false
		}

		return true
	}

	return false
}

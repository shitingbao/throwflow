package domain

import (
	"context"
	"time"
)

type KuaidiCompany struct {
	Id         uint64
	Name       string
	Code       string
	CreateTime time.Time
	UpdateTime time.Time
}

type KuaidiCompanyList struct {
	PageNum  uint64
	PageSize uint64
	Total    uint64
	List     []*KuaidiCompany
}

func NewKuaidiCompany(ctx context.Context, name, code string) *KuaidiCompany {
	return &KuaidiCompany{
		Name: name,
		Code: code,
	}
}

func (kc *KuaidiCompany) SetName(ctx context.Context, name string) {
	kc.Name = name
}

func (kc *KuaidiCompany) SetCode(ctx context.Context, code string) {
	kc.Code = code
}

func (kc *KuaidiCompany) SetUpdateTime(ctx context.Context) {
	kc.UpdateTime = time.Now()
}

func (kc *KuaidiCompany) SetCreateTime(ctx context.Context) {
	kc.CreateTime = time.Now()
}

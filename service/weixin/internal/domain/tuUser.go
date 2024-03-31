package domain

import (
	"context"
)

type TuUser struct {
	UserId   uint64
	Phone    string
	Num      uint64
	Id       string
	ParentId string
	Level    string
	Balance  float64
}

func NewTuUser(ctx context.Context, num uint64, phone, id, parentId, level string) *TuUser {
	return &TuUser{
		Phone:    phone,
		Num:      num,
		Id:       id,
		ParentId: parentId,
		Level:    level,
	}
}

func (tu *TuUser) SetPhone(ctx context.Context, phone string) {
	tu.Phone = phone
}

func (tu *TuUser) SetNum(ctx context.Context, num uint64) {
	tu.Num = num
}

func (tu *TuUser) SetCountryCode(ctx context.Context, id string) {
	tu.Id = id
}

func (tu *TuUser) SetNickName(ctx context.Context, parentId string) {
	tu.ParentId = parentId
}

func (tu *TuUser) SetAvatarUrl(ctx context.Context, level string) {
	tu.Level = level
}

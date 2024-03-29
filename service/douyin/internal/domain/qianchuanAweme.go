package domain

import (
	"context"
	"time"
)

type QianchuanAweme struct {
	AwemeId                 uint64
	AdvertiserId            uint64
	AwemeAvatar             string
	AwemeShowId             string
	AwemeName               string
	AwemeStatus             string
	BindType                []*string
	AwemeHasVideoPermission bool
	AwemeHasLivePermission  bool
	AwemeHasUniProm         bool
	CreateTime              time.Time
	UpdateTime              time.Time
}

func NewQianchuanAweme(ctx context.Context, awemeId, advertiserId uint64, awemeAvatar, awemeShowId, awemeName, awemeStatus string, awemeHasVideoPermission, awemeHasLivePermission, awemeHasUniProm bool, bindType []*string) *QianchuanAweme {
	return &QianchuanAweme{
		AwemeId:                 awemeId,
		AdvertiserId:            advertiserId,
		AwemeAvatar:             awemeAvatar,
		AwemeShowId:             awemeShowId,
		AwemeName:               awemeName,
		AwemeStatus:             awemeStatus,
		BindType:                bindType,
		AwemeHasVideoPermission: awemeHasVideoPermission,
		AwemeHasLivePermission:  awemeHasLivePermission,
		AwemeHasUniProm:         awemeHasUniProm,
	}
}

func (qw *QianchuanAweme) SetAwemeId(ctx context.Context, awemeId uint64) {
	qw.AwemeId = awemeId
}

func (qw *QianchuanAweme) SetAdvertiserId(ctx context.Context, advertiserId uint64) {
	qw.AdvertiserId = advertiserId
}

func (qw *QianchuanAweme) SetAwemeAvatar(ctx context.Context, awemeAvatar string) {
	qw.AwemeAvatar = awemeAvatar
}

func (qw *QianchuanAweme) SetAwemeShowId(ctx context.Context, awemeShowId string) {
	qw.AwemeShowId = awemeShowId
}

func (qw *QianchuanAweme) SetAwemeName(ctx context.Context, awemeName string) {
	qw.AwemeName = awemeName
}

func (qw *QianchuanAweme) SetAwemeStatus(ctx context.Context, awemeStatus string) {
	qw.AwemeStatus = awemeStatus
}

func (qw *QianchuanAweme) SetBindType(ctx context.Context, bindType []*string) {
	qw.BindType = bindType
}

func (qw *QianchuanAweme) SetAwemeHasVideoPermission(ctx context.Context, awemeHasVideoPermission bool) {
	qw.AwemeHasVideoPermission = awemeHasVideoPermission
}

func (qw *QianchuanAweme) SetAwemeHasLivePermission(ctx context.Context, awemeHasLivePermission bool) {
	qw.AwemeHasLivePermission = awemeHasLivePermission
}

func (qw *QianchuanAweme) SetAwemeHasUniProm(ctx context.Context, awemeHasUniProm bool) {
	qw.AwemeHasUniProm = awemeHasUniProm
}

func (qw *QianchuanAweme) SetUpdateTime(ctx context.Context) {
	qw.UpdateTime = time.Now()
}

func (qw *QianchuanAweme) SetCreateTime(ctx context.Context) {
	qw.CreateTime = time.Now()
}

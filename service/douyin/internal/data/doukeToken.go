package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// 抖音开放平台达人授权token表
type DoukeToken struct {
	Id              uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	AuthorityId     string    `gorm:"column:authority_id;type:varchar(50);not null;uniqueIndex:authority_id;comment:授权ID"`
	AuthSubjectType string    `gorm:"column:auth_subject_type;type:varchar(50);not null;uniqueIndex:authority_id;comment:授权主体类型"`
	AccessToken     string    `gorm:"column:access_token;type:varchar(250);not null;comment:token值"`
	ExpiresIn       uint64    `gorm:"column:expires_in;type:bigint(20) UNSIGNED;not null;comment:过期时间(秒级时间戳)"`
	RefreshToken    string    `gorm:"column:refresh_token;type:varchar(250);not null;comment:刷新token值。用于刷新access_token的刷新令牌（有效期：14 天）"`
	CreateTime      time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime      time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (DoukeToken) TableName() string {
	return "douyin_douke_token"
}

type doukeTokenRepo struct {
	data *Data
	log  *log.Helper
}

func (dt *DoukeToken) ToDomain() *domain.DoukeToken {
	return &domain.DoukeToken{
		Id:              dt.Id,
		AuthorityId:     dt.AuthorityId,
		AuthSubjectType: dt.AuthSubjectType,
		AccessToken:     dt.AccessToken,
		ExpiresIn:       dt.ExpiresIn,
		RefreshToken:    dt.RefreshToken,
		CreateTime:      dt.CreateTime,
		UpdateTime:      dt.UpdateTime,
	}
}

func NewDoukeTokenRepo(data *Data, logger log.Logger) biz.DoukeTokenRepo {
	return &doukeTokenRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (dtr *doukeTokenRepo) Get(ctx context.Context, authorityId, authSubjectType string) (*domain.DoukeToken, error) {
	doukeToken := &DoukeToken{}

	if result := dtr.data.db.WithContext(ctx).Where("authority_id = ?", authorityId).Where("auth_subject_type = ?", authSubjectType).First(doukeToken); result.Error != nil {
		return nil, result.Error
	}

	return doukeToken.ToDomain(), nil
}

func (dtr *doukeTokenRepo) List(ctx context.Context) ([]*domain.DoukeToken, error) {
	var doukeTokens []DoukeToken
	list := make([]*domain.DoukeToken, 0)

	if result := dtr.data.db.WithContext(ctx).
		Find(&doukeTokens); result.Error != nil {
		return nil, result.Error
	}

	for _, doukeToken := range doukeTokens {
		list = append(list, doukeToken.ToDomain())
	}

	return list, nil
}

func (dtr *doukeTokenRepo) Save(ctx context.Context, in *domain.DoukeToken) (*domain.DoukeToken, error) {
	doukeToken := &DoukeToken{
		AuthorityId:     in.AuthorityId,
		AuthSubjectType: in.AuthSubjectType,
		AccessToken:     in.AccessToken,
		ExpiresIn:       in.ExpiresIn,
		RefreshToken:    in.RefreshToken,
		CreateTime:      in.CreateTime,
		UpdateTime:      in.UpdateTime,
	}

	if result := dtr.data.DB(ctx).Create(doukeToken); result.Error != nil {
		return nil, result.Error
	}

	return doukeToken.ToDomain(), nil
}

func (dtr *doukeTokenRepo) Update(ctx context.Context, in *domain.DoukeToken) (*domain.DoukeToken, error) {
	doukeToken := &DoukeToken{
		Id:              in.Id,
		AuthorityId:     in.AuthorityId,
		AuthSubjectType: in.AuthSubjectType,
		AccessToken:     in.AccessToken,
		ExpiresIn:       in.ExpiresIn,
		RefreshToken:    in.RefreshToken,
		CreateTime:      in.CreateTime,
		UpdateTime:      in.UpdateTime,
	}

	if result := dtr.data.db.WithContext(ctx).Save(doukeToken); result.Error != nil {
		return nil, result.Error
	}

	return doukeToken.ToDomain(), nil
}

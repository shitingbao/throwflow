package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"weixin/internal/biz"
	"weixin/internal/domain"
)

// 微信小程序用户openId表
type UserOpenId struct {
	Id         uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	UserId     uint64    `gorm:"column:user_id;type:bigint(20) UNSIGNED;not null;index:user_id;comment:微信小程序用户ID"`
	Appid      string    `gorm:"column:appid;type:varchar(100);not null;comment:小程序appid"`
	OpenId     string    `gorm:"column:open_id;type:varchar(100);not null;comment:用户唯一标识"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (UserOpenId) TableName() string {
	return "weixin_user_open_id"
}

type userOpenIdRepo struct {
	data *Data
	log  *log.Helper
}

func (uoi *UserOpenId) ToDomain(ctx context.Context) *domain.UserOpenId {
	userOpenId := &domain.UserOpenId{
		Id:         uoi.Id,
		UserId:     uoi.UserId,
		Appid:      uoi.Appid,
		OpenId:     uoi.OpenId,
		CreateTime: uoi.CreateTime,
		UpdateTime: uoi.UpdateTime,
	}

	return userOpenId
}

func NewUserOpenIdRepo(data *Data, logger log.Logger) biz.UserOpenIdRepo {
	return &userOpenIdRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (uoir *userOpenIdRepo) Get(ctx context.Context, userId uint64, appid, openId string) (*domain.UserOpenId, error) {
	userOpenId := &UserOpenId{}

	db := uoir.data.db.WithContext(ctx).Where("user_id = ?", userId).Where("appid = ?", appid)

	if len(openId) > 0 {
		db = db.Where("open_id = ?", openId)
	}

	if result := db.First(userOpenId); result.Error != nil {
		return nil, result.Error
	}

	return userOpenId.ToDomain(ctx), nil
}

func (uoir *userOpenIdRepo) List(ctx context.Context, userId uint64) ([]*domain.UserOpenId, error) {
	var userOpenIds []UserOpenId
	list := make([]*domain.UserOpenId, 0)

	db := uoir.data.db.WithContext(ctx).Where("user_id = ?", userId)

	if result := db.Order("create_time DESC,id DESC").
		Find(&userOpenIds); result.Error != nil {
		return nil, result.Error
	}

	for _, userOpenId := range userOpenIds {
		list = append(list, userOpenId.ToDomain(ctx))
	}

	return list, nil
}

func (uoir *userOpenIdRepo) Update(ctx context.Context, in *domain.UserOpenId) (*domain.UserOpenId, error) {
	userOpenId := &UserOpenId{
		Id:         in.Id,
		UserId:     in.UserId,
		Appid:      in.Appid,
		OpenId:     in.OpenId,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
	}

	if result := uoir.data.DB(ctx).Save(userOpenId); result.Error != nil {
		return nil, result.Error
	}

	return userOpenId.ToDomain(ctx), nil
}

func (uoir *userOpenIdRepo) Save(ctx context.Context, in *domain.UserOpenId) (*domain.UserOpenId, error) {
	userOpenId := &UserOpenId{
		UserId:     in.UserId,
		Appid:      in.Appid,
		OpenId:     in.OpenId,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
	}

	if result := uoir.data.DB(ctx).Create(userOpenId); result.Error != nil {
		return nil, result.Error
	}

	return userOpenId.ToDomain(ctx), nil
}

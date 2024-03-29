package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// 抖音开放平台达人授权token表
type OpenDouyinToken struct {
	Id               uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	ClientKey        string    `gorm:"column:client_key;type:varchar(50);not null;uniqueIndex:client_key_open_id;comment:应用Client Key"`
	OpenId           string    `gorm:"column:open_id;type:varchar(100);not null;uniqueIndex:client_key_open_id;comment:授权用户唯一标识"`
	AccessToken      string    `gorm:"column:access_token;type:varchar(250);not null;comment:接口调用凭证"`
	ExpiresIn        uint64    `gorm:"column:expires_in;type:bigint(20) UNSIGNED;not null;comment:access_token接口调用凭证超时时间，单位（秒)"`
	RefreshToken     string    `gorm:"column:refresh_token;type:varchar(250);not null;comment:用户刷新access_token"`
	RefreshExpiresIn uint64    `gorm:"column:refresh_expires_in;type:bigint(20) UNSIGNED;not null;comment:refresh_token凭证超时时间，单位（秒)"`
	CreateTime       time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime       time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (OpenDouyinToken) TableName() string {
	return "douyin_open_douyin_token"
}

type openDouyinTokenRepo struct {
	data *Data
	log  *log.Helper
}

func (odt *OpenDouyinToken) ToDomain() *domain.OpenDouyinToken {
	return &domain.OpenDouyinToken{
		Id:               odt.Id,
		ClientKey:        odt.ClientKey,
		OpenId:           odt.OpenId,
		AccessToken:      odt.AccessToken,
		ExpiresIn:        odt.ExpiresIn,
		RefreshToken:     odt.RefreshToken,
		RefreshExpiresIn: odt.RefreshExpiresIn,
		CreateTime:       odt.CreateTime,
		UpdateTime:       odt.UpdateTime,
	}
}

func NewOpenDouyinTokenRepo(data *Data, logger log.Logger) biz.OpenDouyinTokenRepo {
	return &openDouyinTokenRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (odtr *openDouyinTokenRepo) GetByClientKeyAndOpenId(ctx context.Context, clientKey, openId string) (*domain.OpenDouyinToken, error) {
	openDouyinToken := &OpenDouyinToken{}

	if result := odtr.data.db.WithContext(ctx).Where("client_key = ?", clientKey).Where("open_id = ?", openId).First(openDouyinToken); result.Error != nil {
		return nil, result.Error
	}

	return openDouyinToken.ToDomain(), nil
}

func (odtr *openDouyinTokenRepo) List(ctx context.Context) ([]*domain.OpenDouyinToken, error) {
	var openDouyinTokens []OpenDouyinToken
	list := make([]*domain.OpenDouyinToken, 0)

	if result := odtr.data.db.WithContext(ctx).
		Find(&openDouyinTokens); result.Error != nil {
		return nil, result.Error
	}

	for _, openDouyinToken := range openDouyinTokens {
		list = append(list, openDouyinToken.ToDomain())
	}

	return list, nil
}

func (odtr *openDouyinTokenRepo) ListByClientKeyAndOpenId(ctx context.Context, wopenDouyinTokens []*domain.OpenDouyinToken) ([]*domain.OpenDouyinToken, error) {
	var openDouyinTokens []OpenDouyinToken
	list := make([]*domain.OpenDouyinToken, 0)

	db := odtr.data.db.WithContext(ctx)

	openDouyinTokenSqls := make([]string, 0)

	for _, wopenDouyinToken := range wopenDouyinTokens {
		openDouyinTokenSqls = append(openDouyinTokenSqls, "(client_key = '"+wopenDouyinToken.ClientKey+"' and open_id = '"+wopenDouyinToken.OpenId+"')")
	}

	db = db.Where(strings.Join(openDouyinTokenSqls, " or "))

	if result := db.Find(&openDouyinTokens); result.Error != nil {
		return nil, result.Error
	}

	for _, openDouyinToken := range openDouyinTokens {
		list = append(list, openDouyinToken.ToDomain())
	}

	return list, nil
}

func (odtr *openDouyinTokenRepo) ListByCreateTime(ctx context.Context, day string) ([]*domain.OpenDouyinToken, error) {
	var openDouyinTokens []OpenDouyinToken
	list := make([]*domain.OpenDouyinToken, 0)

	if result := odtr.data.db.WithContext(ctx).
		Where("create_time >= ?", day+" 00:00:00").
		Where("create_time <= ?", day+" 23:59:59").
		Find(&openDouyinTokens); result.Error != nil {
		return nil, result.Error
	}

	for _, openDouyinToken := range openDouyinTokens {
		list = append(list, openDouyinToken.ToDomain())
	}

	return list, nil
}

func (odtr *openDouyinTokenRepo) Save(ctx context.Context, in *domain.OpenDouyinToken) (*domain.OpenDouyinToken, error) {
	openDouyinToken := &OpenDouyinToken{
		ClientKey:        in.ClientKey,
		OpenId:           in.OpenId,
		AccessToken:      in.AccessToken,
		ExpiresIn:        in.ExpiresIn,
		RefreshToken:     in.RefreshToken,
		RefreshExpiresIn: in.RefreshExpiresIn,
		CreateTime:       in.CreateTime,
		UpdateTime:       in.UpdateTime,
	}

	if result := odtr.data.DB(ctx).Create(openDouyinToken); result.Error != nil {
		return nil, result.Error
	}

	return openDouyinToken.ToDomain(), nil
}

func (odtr *openDouyinTokenRepo) Update(ctx context.Context, in *domain.OpenDouyinToken) (*domain.OpenDouyinToken, error) {
	openDouyinToken := &OpenDouyinToken{
		Id:               in.Id,
		ClientKey:        in.ClientKey,
		OpenId:           in.OpenId,
		AccessToken:      in.AccessToken,
		ExpiresIn:        in.ExpiresIn,
		RefreshToken:     in.RefreshToken,
		RefreshExpiresIn: in.RefreshExpiresIn,
		CreateTime:       in.CreateTime,
		UpdateTime:       in.UpdateTime,
	}

	if result := odtr.data.db.WithContext(ctx).Save(openDouyinToken); result.Error != nil {
		return nil, result.Error
	}

	return openDouyinToken.ToDomain(), nil
}

func (odtr *openDouyinTokenRepo) SaveCacheString(ctx context.Context, key string, val string, timeout time.Duration) (bool, error) {
	result, err := odtr.data.rdb.SetNX(ctx, key, val, timeout).Result()

	if err != nil {
		return false, err
	}

	return result, nil
}

func (odtr *openDouyinTokenRepo) GetCacheString(ctx context.Context, key string) (string, error) {
	val, err := odtr.data.rdb.Get(ctx, key).Result()

	if err != nil {
		return "", err
	}

	return val, nil
}

func (odtr *openDouyinTokenRepo) DeleteCache(ctx context.Context, key string) error {
	if _, err := odtr.data.rdb.Del(ctx, key).Result(); err != nil {
		return err
	}

	return nil
}

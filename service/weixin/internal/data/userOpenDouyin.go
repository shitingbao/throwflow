package data

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"
	"weixin/internal/biz"
	"weixin/internal/domain"

	"github.com/go-kratos/kratos/v2/log"
)

// 微信小程序用户关联抖音开放平台用户表
type UserOpenDouyin struct {
	Id              uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	UserId          uint64    `gorm:"column:user_id;type:bigint(20) UNSIGNED;not null;uniqueIndex:user_id_client_key_open_id;comment:微信小程序用户ID"`
	ClientKey       string    `gorm:"column:client_key;type:varchar(50);not null;uniqueIndex:user_id_client_key_open_id;comment:抖音开放平台应用Client Key"`
	OpenId          string    `gorm:"column:open_id;type:varchar(100);not null;uniqueIndex:user_id_client_key_open_id;comment:抖音开放平台授权用户唯一标识"`
	AwemeId         uint64    `gorm:"column:aweme_id;type:bigint(20) UNSIGNED;not null;comment:抖音号ID"`
	AccountId       string    `gorm:"column:account_id;type:varchar(100);not null;comment:抖音账户"`
	Nickname        string    `gorm:"column:nickname;type:varchar(250);not null;comment:达人昵称"`
	Avatar          string    `gorm:"column:avatar;type:varchar(250);not null;comment:达人头像"`
	AvatarLarger    string    `gorm:"column:avatar_larger;type:varchar(250);not null;comment:达人头像(大图)"`
	CooperativeCode string    `gorm:"column:cooperative_code;type:varchar(150);not null;comment:合作码"`
	Fans            uint64    `gorm:"column:fans;type:bigint(20) UNSIGNED;not null;comment:粉丝数"`
	Area            string    `gorm:"column:area;type:varchar(100);not null;comment:地区"`
	CreateTime      time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime      time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (UserOpenDouyin) TableName() string {
	return "weixin_user_open_douyin"
}

type userOpenDouyinRepo struct {
	data *Data
	log  *log.Helper
}

func (uod *UserOpenDouyin) ToDomain(ctx context.Context) *domain.UserOpenDouyin {
	userOpenDouyin := &domain.UserOpenDouyin{
		Id:              uod.Id,
		UserId:          uod.UserId,
		ClientKey:       uod.ClientKey,
		OpenId:          uod.OpenId,
		AwemeId:         uod.AwemeId,
		AccountId:       uod.AccountId,
		Nickname:        uod.Nickname,
		Avatar:          uod.Avatar,
		AvatarLarger:    uod.AvatarLarger,
		CooperativeCode: uod.CooperativeCode,
		Fans:            uod.Fans,
		Area:            uod.Area,
		CreateTime:      uod.CreateTime,
		UpdateTime:      uod.UpdateTime,
	}

	return userOpenDouyin
}

func NewUserOpenDouyinRepo(data *Data, logger log.Logger) biz.UserOpenDouyinRepo {
	return &userOpenDouyinRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (uodr *userOpenDouyinRepo) Get(ctx context.Context, userId uint64, clientKey, openId string) (*domain.UserOpenDouyin, error) {
	userOpenDouyin := &UserOpenDouyin{}

	if result := uodr.data.db.WithContext(ctx).Where("user_id = ?", userId).Where("client_key = ?", clientKey).Where("open_id = ?", openId).First(userOpenDouyin); result.Error != nil {
		return nil, result.Error
	}

	return userOpenDouyin.ToDomain(ctx), nil
}

func (uodr *userOpenDouyinRepo) GetById(ctx context.Context, userId, openDouyinUserId uint64) (*domain.UserOpenDouyin, error) {
	userOpenDouyin := &UserOpenDouyin{}

	if result := uodr.data.db.WithContext(ctx).Where("user_id = ?", userId).Where("id = ?", openDouyinUserId).First(userOpenDouyin); result.Error != nil {
		return nil, result.Error
	}

	return userOpenDouyin.ToDomain(ctx), nil
}

func (uodr *userOpenDouyinRepo) GetByClientKeyAndOpenId(ctx context.Context, clientKey, openId string) (*domain.UserOpenDouyin, error) {
	userOpenDouyin := &UserOpenDouyin{}

	if result := uodr.data.db.WithContext(ctx).Where("client_key = ?", clientKey).Where("open_id = ?", openId).First(userOpenDouyin); result.Error != nil {
		return nil, result.Error
	}

	return userOpenDouyin.ToDomain(ctx), nil
}

func (uodr *userOpenDouyinRepo) List(ctx context.Context, pageNum, pageSize int, userId uint64, keyword string) ([]*domain.UserOpenDouyin, error) {
	var userOpenDouyins []UserOpenDouyin
	list := make([]*domain.UserOpenDouyin, 0)

	db := uodr.data.db.WithContext(ctx).Where("user_id = ?", userId)

	if l := utf8.RuneCountInString(keyword); l > 0 {
		db = db.Where("nickname like ?", "%"+keyword+"%")
	}

	if pageNum == 0 {
		if result := db.Order("create_time DESC,id DESC").
			Find(&userOpenDouyins); result.Error != nil {
			return nil, result.Error
		}
	} else {
		if result := db.Order("create_time DESC,id DESC").
			Limit(pageSize).Offset((pageNum - 1) * pageSize).
			Find(&userOpenDouyins); result.Error != nil {
			return nil, result.Error
		}
	}

	for _, userOpenDouyin := range userOpenDouyins {
		list = append(list, userOpenDouyin.ToDomain(ctx))
	}

	return list, nil
}

func (uodr *userOpenDouyinRepo) ListByClientKeyAndOpenId(ctx context.Context, pageNum, pageSize int, clientKeyAndOpenIds []*domain.UserOpenDouyin, keyword string) ([]*domain.UserOpenDouyin, error) {
	var userOpenDouyins []UserOpenDouyin
	list := make([]*domain.UserOpenDouyin, 0)

	db := uodr.data.db.WithContext(ctx)

	if l := utf8.RuneCountInString(keyword); l > 0 {
		db = db.Where("nickname like ?", "%"+keyword+"%")
	}

	clientKeyAndOpenIdSqls := make([]string, 0)

	for _, clientKeyAndOpenId := range clientKeyAndOpenIds {
		clientKeyAndOpenIdSqls = append(clientKeyAndOpenIdSqls, "(client_key = '"+clientKeyAndOpenId.ClientKey+"' and open_id = '"+clientKeyAndOpenId.OpenId+"')")
	}

	if len(clientKeyAndOpenIdSqls) > 0 {
		db = db.Where(strings.Join(clientKeyAndOpenIdSqls, " or "))
	}

	if pageNum == 0 {
		if result := db.Order("create_time DESC,id DESC").
			Find(&userOpenDouyins); result.Error != nil {
			return nil, result.Error
		}
	} else {
		if result := db.Order("create_time DESC,id DESC").
			Limit(pageSize).Offset((pageNum - 1) * pageSize).
			Find(&userOpenDouyins); result.Error != nil {
			return nil, result.Error
		}
	}

	for _, userOpenDouyin := range userOpenDouyins {
		list = append(list, userOpenDouyin.ToDomain(ctx))
	}

	return list, nil
}

func (uodr *userOpenDouyinRepo) Count(ctx context.Context, userId uint64, keyword string) (int64, error) {
	var count int64

	db := uodr.data.db.WithContext(ctx).Where("user_id = ?", userId)

	if l := utf8.RuneCountInString(keyword); l > 0 {
		db = db.Where("nickname like ?", "%"+keyword+"%")
	}

	if result := db.Model(&UserOpenDouyin{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (uodr *userOpenDouyinRepo) CountByClientKeyAndOpenId(ctx context.Context, clientKeyAndOpenIds []*domain.UserOpenDouyin, keyword string) (int64, error) {
	var count int64

	db := uodr.data.db.WithContext(ctx)

	if l := utf8.RuneCountInString(keyword); l > 0 {
		db = db.Where("nickname like ?", "%"+keyword+"%")
	}

	clientKeyAndOpenIdSqls := make([]string, 0)

	for _, clientKeyAndOpenId := range clientKeyAndOpenIds {
		clientKeyAndOpenIdSqls = append(clientKeyAndOpenIdSqls, "(client_key = '"+clientKeyAndOpenId.ClientKey+"' and open_id = '"+clientKeyAndOpenId.OpenId+"')")
	}

	if len(clientKeyAndOpenIdSqls) > 0 {
		db = db.Where(strings.Join(clientKeyAndOpenIdSqls, " or "))
	}

	if result := db.Model(&UserOpenDouyin{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (uodr *userOpenDouyinRepo) Save(ctx context.Context, in *domain.UserOpenDouyin) (*domain.UserOpenDouyin, error) {
	userOpenDouyin := &UserOpenDouyin{
		UserId:          in.UserId,
		ClientKey:       in.ClientKey,
		OpenId:          in.OpenId,
		AwemeId:         in.AwemeId,
		AccountId:       in.AccountId,
		Nickname:        in.Nickname,
		Avatar:          in.Avatar,
		AvatarLarger:    in.AvatarLarger,
		CooperativeCode: in.CooperativeCode,
		Fans:            in.Fans,
		Area:            in.Area,
		CreateTime:      in.CreateTime,
		UpdateTime:      in.UpdateTime,
	}

	if result := uodr.data.DB(ctx).Create(userOpenDouyin); result.Error != nil {
		return nil, result.Error
	}

	return userOpenDouyin.ToDomain(ctx), nil
}

func (uodr *userOpenDouyinRepo) Update(ctx context.Context, in *domain.UserOpenDouyin) (*domain.UserOpenDouyin, error) {
	userOpenDouyin := &UserOpenDouyin{
		Id:              in.Id,
		UserId:          in.UserId,
		ClientKey:       in.ClientKey,
		OpenId:          in.OpenId,
		AwemeId:         in.AwemeId,
		AccountId:       in.AccountId,
		Nickname:        in.Nickname,
		Avatar:          in.Avatar,
		AvatarLarger:    in.AvatarLarger,
		CooperativeCode: in.CooperativeCode,
		Fans:            in.Fans,
		Area:            in.Area,
		CreateTime:      in.CreateTime,
		UpdateTime:      in.UpdateTime,
	}

	if result := uodr.data.DB(ctx).Save(userOpenDouyin); result.Error != nil {
		return nil, result.Error
	}

	return userOpenDouyin.ToDomain(ctx), nil
}

func (uodr *userOpenDouyinRepo) UpdateUserInfos(ctx context.Context, awemeId uint64, clientKey, openId, accountId, nickname, avatar, avatarLarger, area string) error {
	if result := uodr.data.DB(ctx).Model(UserOpenDouyin{}).
		Where("client_key", clientKey).
		Where("open_id", openId).
		Updates(map[string]interface{}{"aweme_id": awemeId, "account_id": accountId, "nickname": nickname, "avatar": avatar, "avatar_larger": avatarLarger, "area": area}); result.Error != nil {
		return result.Error
	}

	return nil
}

func (uodr *userOpenDouyinRepo) UpdateCooperativeCodes(ctx context.Context, clientKey, openId, cooperativeCode string) error {
	if result := uodr.data.DB(ctx).Model(UserOpenDouyin{}).
		Where("client_key", clientKey).
		Where("open_id", openId).
		Updates(map[string]interface{}{"cooperative_code": cooperativeCode}); result.Error != nil {
		return result.Error
	}

	return nil
}

func (uodr *userOpenDouyinRepo) UpdateFans(ctx context.Context, clientKey, openId string, fans uint64) error {
	if result := uodr.data.DB(ctx).Model(UserOpenDouyin{}).
		Where("client_key", clientKey).
		Where("open_id", openId).
		Updates(map[string]interface{}{"fans": fans}); result.Error != nil {
		return result.Error
	}

	return nil
}

func (uodr *userOpenDouyinRepo) Delete(ctx context.Context, in *domain.UserOpenDouyin) error {
	userOpenDouyin := &UserOpenDouyin{
		Id:              in.Id,
		UserId:          in.UserId,
		ClientKey:       in.ClientKey,
		OpenId:          in.OpenId,
		AwemeId:         in.AwemeId,
		AccountId:       in.AccountId,
		Nickname:        in.Nickname,
		Avatar:          in.Avatar,
		AvatarLarger:    in.AvatarLarger,
		CooperativeCode: in.CooperativeCode,
		Fans:            in.Fans,
		Area:            in.Area,
		CreateTime:      in.CreateTime,
		UpdateTime:      in.UpdateTime,
	}

	if result := uodr.data.DB(ctx).Delete(userOpenDouyin); result.Error != nil {
		return result.Error
	}

	return nil
}

func (uodr *userOpenDouyinRepo) DeleteByUserId(ctx context.Context, userId uint64, clientkey, openId string) error {
	if result := uodr.data.DB(ctx).
		Where("user_id != ?", userId).
		Where("client_key = ?", clientkey).
		Where("open_id = ?", openId).
		Delete(&UserOpenDouyin{}); result.Error != nil {
		return result.Error
	}

	return nil
}

package data

import (
	"context"
	"errors"
	"io"
	"time"
	"weixin/internal/biz"
	"weixin/internal/domain"

	"github.com/go-kratos/kratos/v2/log"
	ctos "github.com/volcengine/ve-tos-golang-sdk/v2/tos"
)

// 微信小程序用户表
type User struct {
	Id               uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	Phone            string    `gorm:"column:phone;type:varchar(20);not null;uniqueIndex:phone;comment:手机号"`
	CountryCode      string    `gorm:"column:country_code;type:varchar(50);not null;comment:区号"`
	NickName         string    `gorm:"column:nick_name;type:varchar(200);not null;comment:用户昵称"`
	AvatarUrl        string    `gorm:"column:avatar_url;type:varchar(250);not null;comment:用户头像图片的 URL"`
	Balance          float64   `gorm:"column:balance;type:decimal(10, 2) UNSIGNED;not null;comment:账户余额，单位元"`
	Integral         uint64    `gorm:"column:integral;type:bigint(20) UNSIGNED;not null;default:0;comment:积分"`
	QrCodeUrl        string    `gorm:"column:qr_code_url;type:varchar(250);not null;comment:微信用户小程序码 URL"`
	IdentityCardMark string    `gorm:"column:identity_card_mark;type:varchar(250);not null;comment:身份证号码标记"`
	Ranking          uint64    `gorm:"column:ranking;type:bigint(20) UNSIGNED;not null;default:0;comment:排名"`
	TotalRanking     uint64    `gorm:"column:total_ranking;type:bigint(20) UNSIGNED;not null;default:0;comment:第几个加入"`
	CreateTime       time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime       time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

type Ranking struct {
	Ranking uint64
}

func (User) TableName() string {
	return "weixin_user"
}

type userRepo struct {
	data *Data
	log  *log.Helper
}

func (u *User) ToDomain(ctx context.Context) *domain.User {
	user := &domain.User{
		Id:               u.Id,
		Phone:            u.Phone,
		CountryCode:      u.CountryCode,
		NickName:         u.NickName,
		AvatarUrl:        u.AvatarUrl,
		Balance:          u.Balance,
		Integral:         u.Integral,
		QrCodeUrl:        u.QrCodeUrl,
		IdentityCardMark: u.IdentityCardMark,
		Ranking:          u.Ranking,
		TotalRanking:     u.TotalRanking,
		CreateTime:       u.CreateTime,
		UpdateTime:       u.UpdateTime,
	}

	return user
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ur *userRepo) Get(ctx context.Context, userId uint64) (*domain.User, error) {
	user := &User{}

	if result := ur.data.db.WithContext(ctx).Where("id = ?", userId).First(user); result.Error != nil {
		return nil, result.Error
	}

	return user.ToDomain(ctx), nil
}

func (ur *userRepo) GetByPhoneAndCountryCode(ctx context.Context, phone, countryCode string) (*domain.User, error) {
	user := &User{}

	if result := ur.data.db.WithContext(ctx).Where("phone = ?", phone).Where("country_code = ?", countryCode).First(user); result.Error != nil {
		return nil, result.Error
	}

	return user.ToDomain(ctx), nil
}

func (ur *userRepo) List(ctx context.Context) ([]*domain.User, error) {
	var users []User
	list := make([]*domain.User, 0)

	if result := ur.data.db.WithContext(ctx).
		Order("create_time DESC,id DESC").
		Find(&users); result.Error != nil {
		return nil, result.Error
	}

	for _, user := range users {
		list = append(list, user.ToDomain(ctx))
	}

	return list, nil
}

func (ur *userRepo) ListRanking(ctx context.Context) ([]*domain.User, error) {
	var users []User
	list := make([]*domain.User, 0)

	ur.data.db.WithContext(ctx).Raw("SELECT wu.id,wu.integral, @cur_count := @cur_count + 1, if(@pre_score = wu.integral,@cur_rank,@cur_rank := @cur_count) ranking, @pre_score := wu.integral FROM weixin_user wu,(SELECT @cur_count := 0,@cur_rank:=0,@pre_score := NULL) tmp ORDER BY wu.integral DESC").Scan(&users)

	for _, user := range users {
		list = append(list, user.ToDomain(ctx))
	}

	return list, nil
}

func (ur *userRepo) ListByIds(ctx context.Context, phone, keyword string, ids []uint64) ([]*domain.User, error) {
	var users []User
	list := make([]*domain.User, 0)

	db := ur.data.db.WithContext(ctx)

	if len(ids) > 0 {
		db = db.Where("id in (?)", ids)
	}

	if len(keyword) > 0 {
		db = db.Where("nick_name like ?", "%"+keyword+"%")
	}

	if len(phone) > 0 {
		db = db.Where("phone = ?", phone)
	}

	if result := db.
		Order("create_time DESC,id DESC").
		Find(&users); result.Error != nil {
		return nil, result.Error
	}

	for _, user := range users {
		list = append(list, user.ToDomain(ctx))
	}

	return list, nil
}

func (ur *userRepo) Count(ctx context.Context) (int64, error) {
	var count int64

	if result := ur.data.db.WithContext(ctx).
		Model(&User{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (ur *userRepo) CountByUserId(ctx context.Context, userId uint64) (int64, error) {
	var count int64

	if result := ur.data.db.WithContext(ctx).Where("id <= ?", userId).
		Model(&User{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (ur *userRepo) CountByNickNameOrPhone(ctx context.Context, phone, keyword string) (int64, error) {
	var count int64

	db := ur.data.db.Model(&User{})

	if len(phone) > 0 {
		db = db.Where("phone = ?", phone)
	}

	if len(keyword) > 0 {
		db = db.Where("locate(?, nick_name)>0", keyword)
	}

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (ur *userRepo) Update(ctx context.Context, in *domain.User) (*domain.User, error) {
	user := &User{
		Id:               in.Id,
		Phone:            in.Phone,
		CountryCode:      in.CountryCode,
		NickName:         in.NickName,
		AvatarUrl:        in.AvatarUrl,
		Balance:          in.Balance,
		Integral:         in.Integral,
		QrCodeUrl:        in.QrCodeUrl,
		IdentityCardMark: in.IdentityCardMark,
		Ranking:          in.Ranking,
		TotalRanking:     in.TotalRanking,
		CreateTime:       in.CreateTime,
		UpdateTime:       in.UpdateTime,
	}

	if result := ur.data.DB(ctx).Save(user); result.Error != nil {
		return nil, result.Error
	}

	return user.ToDomain(ctx), nil
}

func (ur *userRepo) UpdateRanking(ctx context.Context, userId, ranking uint64) error {
	if result := ur.data.DB(ctx).Model(User{}).
		Where("id", userId).
		Updates(map[string]interface{}{"ranking": ranking}); result.Error != nil {
		return result.Error
	}

	return nil
}

func (ur *userRepo) Save(ctx context.Context, in *domain.User) (*domain.User, error) {
	user := &User{
		Phone:            in.Phone,
		CountryCode:      in.CountryCode,
		NickName:         in.NickName,
		AvatarUrl:        in.AvatarUrl,
		Balance:          in.Balance,
		Integral:         in.Integral,
		QrCodeUrl:        in.QrCodeUrl,
		IdentityCardMark: in.IdentityCardMark,
		Ranking:          in.Ranking,
		TotalRanking:     in.TotalRanking,
		CreateTime:       in.CreateTime,
		UpdateTime:       in.UpdateTime,
	}

	if result := ur.data.DB(ctx).Create(user); result.Error != nil {
		return nil, result.Error
	}

	return user.ToDomain(ctx), nil
}

func (ur *userRepo) GetCacheHash(ctx context.Context, key string, field string) (string, error) {
	val, err := ur.data.rdb.HGet(ctx, key, field).Result()

	if err != nil {
		return "", err
	}

	return val, nil
}

func (ur *userRepo) SaveCacheHash(ctx context.Context, key string, val map[string]string, timeout time.Duration) error {
	_, err := ur.data.rdb.HMSet(ctx, key, val).Result()

	if err != nil {
		return err
	}

	_, err = ur.data.rdb.Expire(ctx, key, timeout).Result()

	if err != nil {
		return err
	}

	return nil
}

func (ur *userRepo) PutContent(ctx context.Context, fileName string, content io.Reader) (*ctos.PutObjectV2Output, error) {
	for _, ltos := range ur.data.toses {
		if ltos.name == "avatar" {
			output, err := ltos.tos.PutContent(ctx, fileName, content)

			if err != nil {
				return nil, err
			}

			return output, nil
		}
	}

	return nil, errors.New("tos is not exist")
}

func (ur *userRepo) PutContentCode(ctx context.Context, fileName string, content io.Reader) (*ctos.PutObjectV2Output, error) {
	for _, ltos := range ur.data.toses {
		if ltos.name == "company" {
			output, err := ltos.tos.PutContent(ctx, fileName, content)

			if err != nil {
				return nil, err
			}

			return output, nil
		}
	}

	return nil, errors.New("tos is not exist")
}

func (ur *userRepo) NextId(ctx context.Context) (uint64, error) {
	return ur.data.sonyflake.NextID()
}

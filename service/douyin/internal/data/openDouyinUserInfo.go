package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
	"time"
)

// 抖音开放平台达人信息表
type OpenDouyinUserInfo struct {
	Id              uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	ClientKey       string    `gorm:"column:client_key;type:varchar(50);not null;uniqueIndex:client_key_open_id;comment:应用Client Key"`
	OpenId          string    `gorm:"column:open_id;type:varchar(100);not null;uniqueIndex:client_key_open_id;comment:授权用户唯一标识"`
	UnionId         string    `gorm:"column:union_id;type:varchar(100);not null;comment:用户在当前开发者账号下的唯一标识"`
	AwemeId         uint64    `gorm:"column:aweme_id;type:bigint(20) UNSIGNED;not null;comment:抖音号ID"`
	AccountId       string    `gorm:"column:account_id;type:varchar(100);not null;comment:抖音账户"`
	BuyinId         string    `gorm:"column:buyin_id;type:varchar(100);not null;comment:百应ID"`
	Nickname        string    `gorm:"column:nickname;type:varchar(250);not null;comment:达人昵称"`
	Avatar          string    `gorm:"column:avatar;type:varchar(250);not null;comment:达人头像"`
	AvatarLarger    string    `gorm:"column:avatar_larger;type:varchar(250);not null;comment:达人头像(大图)"`
	Phone           string    `gorm:"column:avatar_larger;type:varchar(20);not null;comment:手机号"`
	Gender          uint8     `gorm:"column:gender;type:tinyint(3) UNSIGNED;not null;default:0;comment:性别"`
	Country         string    `gorm:"column:country;type:varchar(150);not null;comment:国家"`
	Province        string    `gorm:"column:province;type:varchar(150);not null;comment:省"`
	City            string    `gorm:"column:city;type:varchar(150);not null;comment:市"`
	District        string    `gorm:"column:district;type:varchar(250);not null;comment:街道"`
	EAccountRole    string    `gorm:"column:e_account_role;type:varchar(100);not null;comment:类型：EAccountM：普通企业号，EAccountS：认证企业号，EAccountK：品牌企业号"`
	CooperativeCode string    `gorm:"column:cooperative_code;type:varchar(150);not null;comment:合作码"`
	Fans            uint64    `gorm:"column:fans;type:bigint(20) UNSIGNED;not null;comment:粉丝数"`
	Area            string    `gorm:"column:area;type:varchar(100);not null;comment:地区"`
	CreateTime      time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime      time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (OpenDouyinUserInfo) TableName() string {
	return "douyin_open_douyin_user_info"
}

type openDouyinUserInfoRepo struct {
	data *Data
	log  *log.Helper
}

func (odui *OpenDouyinUserInfo) ToDomain() *domain.OpenDouyinUserInfo {
	return &domain.OpenDouyinUserInfo{
		Id:              odui.Id,
		ClientKey:       odui.ClientKey,
		OpenId:          odui.OpenId,
		UnionId:         odui.UnionId,
		AwemeId:         odui.AwemeId,
		AccountId:       odui.AccountId,
		BuyinId:         odui.BuyinId,
		Nickname:        odui.Nickname,
		Avatar:          odui.Avatar,
		AvatarLarger:    odui.AvatarLarger,
		Phone:           odui.Phone,
		Gender:          odui.Gender,
		Country:         odui.Country,
		Province:        odui.Province,
		City:            odui.City,
		District:        odui.District,
		EAccountRole:    odui.EAccountRole,
		CooperativeCode: odui.CooperativeCode,
		Fans:            odui.Fans,
		Area:            odui.Area,
		CreateTime:      odui.CreateTime,
		UpdateTime:      odui.UpdateTime,
	}
}

func NewOpenDouyinUserInfoRepo(data *Data, logger log.Logger) biz.OpenDouyinUserInfoRepo {
	return &openDouyinUserInfoRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (oduir *openDouyinUserInfoRepo) GetByClientKeyAndOpenId(ctx context.Context, clientKey, openId string) (*domain.OpenDouyinUserInfo, error) {
	openDouyinUserInfo := &OpenDouyinUserInfo{}

	if result := oduir.data.db.WithContext(ctx).Where("client_key = ?", clientKey).Where("open_id = ?", openId).First(openDouyinUserInfo); result.Error != nil {
		return nil, result.Error
	}

	return openDouyinUserInfo.ToDomain(), nil
}

func (oduir *openDouyinUserInfoRepo) List(ctx context.Context, pageNum, pageSize int) ([]*domain.OpenDouyinUserInfo, error) {
	var openDouyinUserInfos []OpenDouyinUserInfo
	list := make([]*domain.OpenDouyinUserInfo, 0)

	if pageNum == 0 {
		if result := oduir.data.db.WithContext(ctx).
			Order("create_time DESC").
			Find(&openDouyinUserInfos); result.Error != nil {
			return nil, result.Error
		}
	} else {
		if result := oduir.data.db.WithContext(ctx).
			Order("create_time DESC").
			Limit(pageSize).Offset((pageNum - 1) * pageSize).
			Find(&openDouyinUserInfos); result.Error != nil {
			return nil, result.Error
		}
	}

	for _, openDouyinUserInfo := range openDouyinUserInfos {
		list = append(list, openDouyinUserInfo.ToDomain())
	}

	return list, nil
}

func (oduir *openDouyinUserInfoRepo) ListByClientKeyAndOpenId(ctx context.Context, wopenDouyinUserInfos []*domain.OpenDouyinUserInfo) ([]*domain.OpenDouyinUserInfo, error) {
	var openDouyinUserInfos []OpenDouyinUserInfo
	list := make([]*domain.OpenDouyinUserInfo, 0)

	db := oduir.data.db.WithContext(ctx)

	openDouyinUserInfoSqls := make([]string, 0)

	for _, wopenDouyinUserInfo := range wopenDouyinUserInfos {
		openDouyinUserInfoSqls = append(openDouyinUserInfoSqls, "(client_key = '"+wopenDouyinUserInfo.ClientKey+"' and open_id = '"+wopenDouyinUserInfo.OpenId+"')")
	}

	db = db.Where(strings.Join(openDouyinUserInfoSqls, " or "))

	if result := db.Find(&openDouyinUserInfos); result.Error != nil {
		return nil, result.Error
	}

	for _, openDouyinUserInfo := range openDouyinUserInfos {
		list = append(list, openDouyinUserInfo.ToDomain())
	}

	return list, nil
}

func (oduir *openDouyinUserInfoRepo) Count(ctx context.Context) (int64, error) {
	var count int64

	if result := oduir.data.db.WithContext(ctx).Model(&OpenDouyinUserInfo{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (oduir *openDouyinUserInfoRepo) Save(ctx context.Context, in *domain.OpenDouyinUserInfo) (*domain.OpenDouyinUserInfo, error) {
	openDouyinUserInfo := &OpenDouyinUserInfo{
		ClientKey:       in.ClientKey,
		OpenId:          in.OpenId,
		UnionId:         in.UnionId,
		AwemeId:         in.AwemeId,
		AccountId:       in.AccountId,
		BuyinId:         in.BuyinId,
		Nickname:        in.Nickname,
		Avatar:          in.Avatar,
		AvatarLarger:    in.AvatarLarger,
		Phone:           in.Phone,
		Gender:          in.Gender,
		Country:         in.Country,
		Province:        in.Province,
		City:            in.City,
		District:        in.District,
		EAccountRole:    in.EAccountRole,
		CooperativeCode: in.CooperativeCode,
		Fans:            in.Fans,
		Area:            in.Area,
		CreateTime:      in.CreateTime,
		UpdateTime:      in.UpdateTime,
	}

	if result := oduir.data.DB(ctx).Create(openDouyinUserInfo); result.Error != nil {
		return nil, result.Error
	}

	return openDouyinUserInfo.ToDomain(), nil
}

func (oduir *openDouyinUserInfoRepo) Update(ctx context.Context, in *domain.OpenDouyinUserInfo) (*domain.OpenDouyinUserInfo, error) {
	openDouyinUserInfo := &OpenDouyinUserInfo{
		Id:              in.Id,
		ClientKey:       in.ClientKey,
		OpenId:          in.OpenId,
		UnionId:         in.UnionId,
		AwemeId:         in.AwemeId,
		AccountId:       in.AccountId,
		BuyinId:         in.BuyinId,
		Nickname:        in.Nickname,
		Avatar:          in.Avatar,
		AvatarLarger:    in.AvatarLarger,
		Phone:           in.Phone,
		Gender:          in.Gender,
		Country:         in.Country,
		Province:        in.Province,
		City:            in.City,
		District:        in.District,
		EAccountRole:    in.EAccountRole,
		CooperativeCode: in.CooperativeCode,
		Fans:            in.Fans,
		Area:            in.Area,
		CreateTime:      in.CreateTime,
		UpdateTime:      in.UpdateTime,
	}

	if result := oduir.data.db.WithContext(ctx).Save(openDouyinUserInfo); result.Error != nil {
		return nil, result.Error
	}

	return openDouyinUserInfo.ToDomain(), nil
}

func (oduir *openDouyinUserInfoRepo) UpdateCooperativeCodes(ctx context.Context, clientKey, openId, cooperativeCode string) error {
	if result := oduir.data.DB(ctx).Model(OpenDouyinUserInfo{}).
		Where("client_key", clientKey).
		Where("open_id", openId).
		Updates(map[string]interface{}{"cooperative_code": cooperativeCode}); result.Error != nil {
		return result.Error
	}

	return nil
}

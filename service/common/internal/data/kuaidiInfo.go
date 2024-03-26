package data

import (
	"common/internal/biz"
	"common/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// 快递运单号信息表
type KuaidiInfo struct {
	Id         uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	Code       string    `gorm:"column:code;type:char(100);not null;comment:快递公司简码"`
	Num        string    `gorm:"column:num;type:char(100);not null;comment:运单号"`
	Phone      string    `gorm:"column:phone;type:char(100);not null;comment:收、寄件人的电话号码（手机和固定电话均可，只能填写一个，顺丰速运、顺丰快运必填，其他快递公司选填。如座机号码有分机号，分机号无需传入。）"`
	State      uint8     `gorm:"column:state;type:tinyint(3) UNSIGNED;not null;default:1;comment:快递单当前状态，默认为0在途，1揽收，2疑难，3签收，4退签，5派件，8清关，14拒签等10个基础物流状态，如需要返回高级物流状态，请参考 resultv2 传值"`
	Content    string    `gorm:"column:video_name;type:text;not null;comment:运单号信息"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (KuaidiInfo) TableName() string {
	return "common_kuaidi_info"
}

type kuaidiInfoRepo struct {
	data *Data
	log  *log.Helper
}

func (ki *KuaidiInfo) ToDomain() *domain.KuaidiInfo {
	return &domain.KuaidiInfo{
		Id:         ki.Id,
		Code:       ki.Code,
		Num:        ki.Num,
		Phone:      ki.Phone,
		State:      ki.State,
		Content:    ki.Content,
		CreateTime: ki.CreateTime,
		UpdateTime: ki.UpdateTime,
	}
}

func NewKuaidiInfoRepo(data *Data, logger log.Logger) biz.KuaidiInfoRepo {
	return &kuaidiInfoRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (kir *kuaidiInfoRepo) Get(ctx context.Context, code, num string) (*domain.KuaidiInfo, error) {
	kuaidiInfo := &KuaidiInfo{}

	if result := kir.data.db.WithContext(ctx).Where("code = ?", code).Where("num = ?", num).First(kuaidiInfo); result.Error != nil {
		return nil, result.Error
	}

	return kuaidiInfo.ToDomain(), nil
}

func (kir *kuaidiInfoRepo) Save(ctx context.Context, in *domain.KuaidiInfo) (*domain.KuaidiInfo, error) {
	kuaidiInfo := &KuaidiInfo{
		Code:       in.Code,
		Num:        in.Num,
		Phone:      in.Phone,
		State:      in.State,
		Content:    in.Content,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
	}

	if result := kir.data.db.WithContext(ctx).Create(kuaidiInfo); result.Error != nil {
		return nil, result.Error
	}

	return kuaidiInfo.ToDomain(), nil
}

func (kir *kuaidiInfoRepo) Update(ctx context.Context, in *domain.KuaidiInfo) (*domain.KuaidiInfo, error) {
	kuaidiInfo := &KuaidiInfo{
		Id:         in.Id,
		Code:       in.Code,
		Num:        in.Num,
		Phone:      in.Phone,
		State:      in.State,
		Content:    in.Content,
		CreateTime: in.CreateTime,
		UpdateTime: in.UpdateTime,
	}

	if result := kir.data.db.WithContext(ctx).Save(kuaidiInfo); result.Error != nil {
		return nil, result.Error
	}

	return kuaidiInfo.ToDomain(), nil
}

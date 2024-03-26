package data

import (
	"common/internal/biz"
	"common/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// 短信日志表
type SmsLog struct {
	Id           uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	SendPhone    string    `gorm:"column:send_phone;type:char(11);not null;index:phone;comment:手机号"`
	Code         string    `gorm:"column:code;type:char(6);not null;comment:验证码"`
	ReplyCode    string    `gorm:"column:reply_code;type:char(10);not null;comment:返回码"`
	SendContent  string    `gorm:"column:send_content;type:text;not null;comment:短信内容"`
	ReplyContent string    `gorm:"column:reply_content;type:text;not null;comment:短信平台响应数据"`
	SendType     string    `gorm:"column:type;type:enum('login','apply','accountOpend');not null;default:'login';comment:类型，login：登录验证码 apply：申请试用验证码 accountOpend：账户开通后推送短信"`
	SendIp       string    `gorm:"column:send_ip;type:char(20);not null;comment:客户端IP"`
	CreateTime   time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime   time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (SmsLog) TableName() string {
	return "common_sms_log"
}

type smsLogRepo struct {
	data *Data
	log  *log.Helper
}

func (sl *SmsLog) ToDomain() *domain.SmsLog {
	return &domain.SmsLog{
		Id:           sl.Id,
		SendPhone:    sl.SendPhone,
		Code:         sl.Code,
		ReplyCode:    sl.ReplyCode,
		SendContent:  sl.SendContent,
		ReplyContent: sl.ReplyContent,
		SendType:     sl.SendType,
		SendIp:       sl.SendIp,
		CreateTime:   sl.CreateTime,
		UpdateTime:   sl.UpdateTime,
	}
}

func NewSmsLogRepo(data *Data, logger log.Logger) biz.SmsRepo {
	return &smsLogRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (slr *smsLogRepo) GetByPhone(ctx context.Context, phone, types, code string) (*domain.SmsLog, error) {
	smsLog := &SmsLog{}

	if result := slr.data.db.WithContext(ctx).
		Where("send_phone = ?", phone).
		Where("type = ?", types).
		Where("reply_code = ?", code).
		Order("update_time DESC").
		First(smsLog); result.Error != nil {
		return nil, result.Error
	}

	return smsLog.ToDomain(), nil
}

func (slr *smsLogRepo) ListByPhone(ctx context.Context, phone, types, replyCode, time string) ([]*domain.SmsLog, error) {
	var smsLogs []SmsLog
	list := make([]*domain.SmsLog, 0)

	db := slr.data.db.WithContext(ctx).Where("send_phone = ?", phone).Where("reply_code = ?", replyCode)

	if result := db.
		Where("type = ?", types).
		Where("update_time between ? and ?", time+" 00:00:00", time+" 23:59:59").
		Order("update_time DESC").
		Find(&smsLogs); result.Error != nil {
		return nil, result.Error
	}

	for _, smsLog := range smsLogs {
		list = append(list, smsLog.ToDomain())
	}

	return list, nil
}

func (slr *smsLogRepo) List(ctx context.Context, pageNum, pageSize int) ([]*domain.SmsLog, error) {
	var smsLogs []SmsLog
	list := make([]*domain.SmsLog, 0)

	if result := slr.data.db.WithContext(ctx).
		Order("id DESC").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).
		Find(&smsLogs); result.Error != nil {
		return nil, result.Error
	}

	for _, smsLog := range smsLogs {
		list = append(list, smsLog.ToDomain())
	}

	return list, nil
}

func (slr *smsLogRepo) Count(ctx context.Context) (int64, error) {
	var count int64

	if result := slr.data.db.Model(&SmsLog{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (slr *smsLogRepo) Save(ctx context.Context, in *domain.SmsLog) (*domain.SmsLog, error) {
	smsLog := &SmsLog{
		SendPhone:    in.SendPhone,
		Code:         in.Code,
		ReplyCode:    in.ReplyCode,
		SendContent:  in.SendContent,
		ReplyContent: in.ReplyContent,
		SendType:     in.SendType,
		SendIp:       in.SendIp,
		CreateTime:   in.CreateTime,
		UpdateTime:   in.UpdateTime,
	}

	if result := slr.data.db.WithContext(ctx).Create(smsLog); result.Error != nil {
		return nil, result.Error
	}

	return smsLog.ToDomain(), nil
}

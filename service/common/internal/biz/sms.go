package biz

import (
	"common/internal/conf"
	"common/internal/domain"
	"common/internal/pkg/sms"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"math/rand"
	"time"
)

var (
	CommonSmsLimitError     = errors.InternalServer("COMMON_SMS_LIMIT_ERROR", "超出验证码发送限制")
	CommonSmsSendError      = errors.InternalServer("COMMON_SMS_SEND_ERROR", "发送信息失败")
	CommonSmsLogCreateError = errors.InternalServer("COMMON_SMS_LOG_CREATE_ERROR", "短信日志创建失败")
	CommonSmsLogNotFound    = errors.InternalServer("COMMON_SMS_LOG_NOT_FOUND", "短信发送记录不存在")
	CommonSmsCodeEffective  = errors.InternalServer("COMMON_SMS_CODE_EFFECTIVE", "短信验证码已过期")
	CommonSmsVerifyError    = errors.InternalServer("COMMON_SMS_VERIFY_ERROR", "短信验证码验证失败")
)

type SmsRepo interface {
	GetByPhone(context.Context, string, string, string) (*domain.SmsLog, error)
	ListByPhone(context.Context, string, string, string, string) ([]*domain.SmsLog, error)
	List(context.Context, int, int) ([]*domain.SmsLog, error)
	Count(context.Context) (int64, error)
	Save(context.Context, *domain.SmsLog) (*domain.SmsLog, error)
}

type SmsUsecase struct {
	repo  SmsRepo
	dconf *conf.Data
	sconf *conf.Sms
	log   *log.Helper
}

type AccountOpend struct {
	Phone    string `json:"phone"`
	RoleName string `json:"roleName"`
}

func NewSmsUsecase(repo SmsRepo, dconf *conf.Data, sconf *conf.Sms, logger log.Logger) *SmsUsecase {
	return &SmsUsecase{repo: repo, dconf: dconf, sconf: sconf, log: log.NewHelper(logger)}
}

func (suc *SmsUsecase) ListSmsLogs(ctx context.Context, pageNum, pageSize uint64) (*domain.SmsLogList, error) {
	list, err := suc.repo.List(ctx, int(pageNum), int(pageSize))

	if err != nil {
		return nil, CommonDataError
	}

	total, err := suc.repo.Count(ctx)

	if err != nil {
		return nil, CommonDataError
	}

	return &domain.SmsLogList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (suc *SmsUsecase) SendSms(ctx context.Context, phone, types, content, ip string) (*domain.SmsLog, error) {
	if types == "login" || types == "apply" {
		code := randData()
		sendContent := fmt.Sprintf(suc.sconf.Chuanglan.Code.Content, code)
		day := time.Now().Format("2006-01-02")

		inSmsLog := domain.NewSmsLog(ctx, phone, code, types, sendContent, ip)

		logs, err := suc.repo.ListByPhone(ctx, inSmsLog.SendPhone, types, "0", day)

		if err != nil {
			return nil, CommonDataError
		}

		logsCount := len(logs)
		now := time.Now()

		if suc.sconf.Chuanglan.Code.DayLimit < int64(logsCount) {
			return nil, CommonSmsLimitError
		}

		if logsCount != 0 {
			limitTime, _ := time.ParseDuration(suc.sconf.Chuanglan.Code.SecondLimit)
			nowTime := now.Add(limitTime)

			if nowTime.Before(logs[0].UpdateTime) {
				return nil, CommonSmsLimitError
			}
		}

		replyContent, err := sms.Send(suc.sconf, inSmsLog.SendPhone, inSmsLog.SendContent)

		if err != nil {
			return nil, CommonSmsSendError
		}

		replyData, err := json.Marshal(replyContent)

		if err != nil {
			return nil, CommonDataError
		}

		inSmsLog.SetReplyCode(ctx, replyContent.Code)
		inSmsLog.SetReplyContent(ctx, string(replyData))
		inSmsLog.SetUpdateTime(ctx)
		inSmsLog.SetCreateTime(ctx)

		smsLog, err := suc.repo.Save(ctx, inSmsLog)

		if err != nil {
			return nil, CommonSmsLogCreateError
		}

		if replyContent.Code != "0" {
			return nil, errors.InternalServer("SMS_SEND_ERROR", replyContent.ErrorMsg)
		}

		return smsLog, nil
	} else if types == "accountOpend" {
		var accountOpend AccountOpend

		if err := json.Unmarshal([]byte(content), &accountOpend); err != nil {
			return nil, CommonDataError
		}

		sendContent := fmt.Sprintf(suc.sconf.Chuanglan.AccountOpened.Content, accountOpend.Phone, accountOpend.RoleName)
		inSmsLog := domain.NewSmsLog(ctx, phone, "", types, sendContent, ip)

		replyContent, err := sms.Send(suc.sconf, inSmsLog.SendPhone, inSmsLog.SendContent)

		if err != nil {
			return nil, CommonSmsSendError
		}

		replyData, err := json.Marshal(replyContent)

		if err != nil {
			return nil, CommonDataError
		}

		inSmsLog.SetReplyCode(ctx, replyContent.Code)
		inSmsLog.SetReplyContent(ctx, string(replyData))
		inSmsLog.SetUpdateTime(ctx)
		inSmsLog.SetCreateTime(ctx)

		smsLog, err := suc.repo.Save(ctx, inSmsLog)

		if err != nil {
			return nil, CommonSmsLogCreateError
		}

		if replyContent.Code != "0" {
			return nil, errors.InternalServer("SMS_SEND_ERROR", replyContent.ErrorMsg)
		}

		return smsLog, nil
	}

	return nil, nil
}

func (suc *SmsUsecase) VerifySms(ctx context.Context, phone, types, code string) error {
	smsLog, err := suc.repo.GetByPhone(ctx, phone, types, "0")

	if err != nil {
		return CommonSmsLogNotFound
	}

	effectiveTime, _ := time.ParseDuration(suc.sconf.Chuanglan.Code.EffectiveDate)

	if smsLog.UpdateTime.Before(time.Now().Add(effectiveTime)) {
		return CommonSmsCodeEffective
	}

	if smsLog.Code != code {
		return CommonSmsVerifyError
	}

	return nil
}

func randData() string {
	randData := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06v", randData.Int31n(1000000))
}

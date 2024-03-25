package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/common/v1"
	"interface/internal/pkg/tool"
)

var (
	InterfaceSmsSendError   = errors.InternalServer("INTERFACE_SMS_SEND_ERROR", "短信发送失败")
	InterfaceSmsVerifyError = errors.InternalServer("INTERFACE_SMS_VERIFY_ERROR", "短信验证码验证失败")
)

type SmsRepo interface {
	Send(context.Context, string, string, string, string) (*v1.SendSmsReply, error)
	Verify(context.Context, string, string, string) (*v1.VerifySmsReply, error)
}

type SmsUsecase struct {
	repo   SmsRepo
	curepo CompanyUserRepo
	log    *log.Helper
}

func NewSmsUsecase(repo SmsRepo, curepo CompanyUserRepo, logger log.Logger) *SmsUsecase {
	return &SmsUsecase{repo: repo, curepo: curepo, log: log.NewHelper(logger)}
}

func (suc *SmsUsecase) SendSms(ctx context.Context, phone, types string) (*v1.SendSmsReply, error) {
	content := ""
	ip := tool.GetClientIp(ctx)

	if types == "login" {
		companyUsers, err := suc.curepo.ListByPhone(ctx, phone)

		if err != nil {
			return nil, InterfaceDataError
		}

		if len(companyUsers.Data.List) == 0 {
			return nil, errors.NotFound("INTERFACE_COMPANY_USER_NOT_FOUND", "当前用户不存在")
		}
	}

	sms, err := suc.repo.Send(ctx, phone, content, types, ip)

	if err != nil {
		return nil, InterfaceSmsSendError
	}

	return sms, nil
}

func (suc *SmsUsecase) VerifyCode(ctx context.Context, phone, types, code string) bool {
	if _, err := suc.repo.Verify(ctx, phone, types, code); err != nil {
		return false
	}

	return true
}

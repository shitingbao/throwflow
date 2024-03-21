package biz

import (
	v1 "company/api/service/weixin/v1"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	CompanyQrCodeGetError = errors.InternalServer("COMPANY_QR_CODE_GET_ERROR", "公司小程序码获取失败")
)

type QrCodeRepo interface {
	Get(context.Context, uint64, string) (*v1.GetQrCodesReply, error)
}

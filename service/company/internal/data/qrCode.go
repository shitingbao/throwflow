package data

import (
	v1 "company/api/service/weixin/v1"
	"company/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type qrCodeRepo struct {
	data *Data
	log  *log.Helper
}

func NewQrCodeRepo(data *Data, logger log.Logger) biz.QrCodeRepo {
	return &qrCodeRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (qcr *qrCodeRepo) Get(ctx context.Context, organizationId uint64, scene string) (*v1.GetQrCodesReply, error) {
	qrCode, err := qcr.data.weixinuc.GetQrCodes(ctx, &v1.GetQrCodesRequest{
		OrganizationId: organizationId,
		Scene:          scene,
	})

	if err != nil {
		return nil, err
	}

	return qrCode, err
}

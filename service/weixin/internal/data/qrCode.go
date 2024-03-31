package data

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	ctos "github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"io"
	"weixin/internal/biz"
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

func (qcr *qrCodeRepo) PutContent(ctx context.Context, fileName string, content io.Reader) (*ctos.PutObjectV2Output, error) {
	for _, ltos := range qcr.data.toses {
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

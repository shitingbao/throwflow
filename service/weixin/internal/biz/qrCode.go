package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	ctos "github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"io"
	"strings"
	"time"
	"weixin/internal/conf"
	"weixin/internal/pkg/mini/oauth2"
	"weixin/internal/pkg/mini/wxa"
	"weixin/internal/pkg/tool"
)

var (
	WeixinQrCodeGetError = errors.InternalServer("WEIXIN_QR_CODE_GET_ERROR", "小程序码获取失败")
)

type QrCodeRepo interface {
	PutContent(context.Context, string, io.Reader) (*ctos.PutObjectV2Output, error)
}

type QrCodeUsecase struct {
	repo  QrCodeRepo
	conf  *conf.Data
	wconf *conf.Weixin
	oconf *conf.Organization
	vconf *conf.Volcengine
	gconf *conf.Gongmall
	log   *log.Helper
}

func NewQrCodeUsecase(repo QrCodeRepo, conf *conf.Data, wconf *conf.Weixin, oconf *conf.Organization, vconf *conf.Volcengine, gconf *conf.Gongmall, logger log.Logger) *QrCodeUsecase {
	return &QrCodeUsecase{repo: repo, conf: conf, wconf: wconf, oconf: oconf, vconf: vconf, gconf: gconf, log: log.NewHelper(logger)}
}

func (qcuc *QrCodeUsecase) GetQrCodes(ctx context.Context, organizationId uint64, scene string) (string, error) {
	var content *wxa.GetUnlimitedQRCodeResponse

	if organizationId == qcuc.oconf.DjOrganizationId {
		accessToken, err := oauth2.GetAccessToken(qcuc.wconf.DjMini.Appid, qcuc.wconf.DjMini.Secret)

		if err != nil {
			return "", errors.InternalServer("WEIXIN_QR_CODE_GET_ERROR", err.Error())
		}

		content, err = wxa.GetUnlimitedQRCode(accessToken.AccessToken, scene, qcuc.wconf.DjMini.QrCodeEnvVersion)

		if err != nil {
			return "", errors.InternalServer("WEIXIN_QR_CODE_GET_ERROR", err.Error())
		}
	} else {
		accessToken, err := oauth2.GetAccessToken(qcuc.wconf.Mini.Appid, qcuc.wconf.Mini.Secret)

		if err != nil {
			return "", errors.InternalServer("WEIXIN_QR_CODE_GET_ERROR", err.Error())
		}

		content, err = wxa.GetUnlimitedQRCode(accessToken.AccessToken, scene, qcuc.wconf.Mini.QrCodeEnvVersion)

		if err != nil {
			return "", errors.InternalServer("WEIXIN_QR_CODE_GET_ERROR", err.Error())
		}
	}

	objectKey := tool.GetRandCode(time.Now().String())

	if _, err := qcuc.repo.PutContent(ctx, qcuc.vconf.Tos.Company.SubFolder+"/"+objectKey+".png", strings.NewReader(content.Buffer)); err != nil {
		return "", WeixinQrCodeGetError
	}

	return qcuc.vconf.Tos.Company.Url + "/" + qcuc.vconf.Tos.Company.SubFolder + "/" + objectKey + ".png", nil
}

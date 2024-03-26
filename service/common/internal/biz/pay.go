package biz

import (
	"common/internal/conf"
	"common/internal/domain"
	"common/internal/pkg/pay/bby"
	"common/internal/pkg/tool"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"time"
)

var (
	CommonPayCreateError = errors.InternalServer("COMMON_PAY_CREATE_ERROR", "支付创建失败")
)

type PayUsecase struct {
	panlrepo PayAsyncNotificationLogRepo
	dconf    *conf.Data
	pconf    *conf.Pay
	oconf    *conf.Organization
	log      *log.Helper
}

func NewPayUsecase(panlrepo PayAsyncNotificationLogRepo, dconf *conf.Data, pconf *conf.Pay, oconf *conf.Organization, logger log.Logger) *PayUsecase {
	return &PayUsecase{panlrepo: panlrepo, dconf: dconf, pconf: pconf, oconf: oconf, log: log.NewHelper(logger)}
}

func (puc *PayUsecase) Pay(ctx context.Context, organizationId uint64, totalFee float64, outTradeNo, content, nonceStr, openId, clientIp string) (*bby.DataPackage, error) {
	var conf *conf.Pay_BbyAccount

	if organizationId == puc.oconf.DjOrganizationId {
		conf = puc.pconf.Bby.DjAccount
	} else {
		conf = puc.pconf.Bby.DefaultAccount
	}

	utotalFee := uint64(tool.Decimal(totalFee*100, 2))

	data, err := bby.Pay(outTradeNo, content, clientIp, nonceStr, openId, utotalFee, conf)

	if err != nil {
		return nil, CommonPayCreateError
	}

	if data.ResultCode == "0" && data.Status == "0" {
		var dataPackage bby.DataPackage

		if err := json.Unmarshal([]byte(data.DataPackage), &dataPackage); err != nil {
			return nil, errors.InternalServer("COMMON_PAY_CREATE_ERROR", fmt.Sprintf("failed json Unmarshal : %v", err.Error()))
		}

		return &dataPackage, nil
	} else {
		if len(data.ErrMsg) > 0 {
			return nil, errors.InternalServer("COMMON_PAY_CREATE_ERROR", data.ErrMsg)
		} else {
			return nil, CommonPayCreateError
		}
	}
}

func (puc *PayUsecase) PayAsyncNotification(ctx context.Context, content string) (*bby.PayAsyncNotificationData, error) {
	inPayAsyncNotificationLog := domain.NewPayAsyncNotificationLog(ctx, content)
	inPayAsyncNotificationLog.SetCreateTime(ctx)
	inPayAsyncNotificationLog.SetUpdateTime(ctx)

	if _, err := puc.panlrepo.Save(ctx, inPayAsyncNotificationLog); err != nil {
		return nil, CommonPayAsyncNotificationLogCreateError
	}

	data, err := bby.AsyncNotification(content, puc.pconf.Bby)

	if err != nil {
		return nil, errors.InternalServer("COMMON_PAY_ASYNC_NOTIFICATION_ERROR", err.Error())
	}

	if data.Status == "0" && data.PayResult == "0" && data.ResultCode == "0" {
		return data, nil
	} else {
		if len(data.ErrMsg) > 0 {
			return nil, errors.InternalServer("COMMON_PAY_ASYNC_NOTIFICATION_ERROR", data.ErrMsg)
		} else {
			return nil, errors.InternalServer("COMMON_PAY_ASYNC_NOTIFICATION_ERROR", "支付失败")
		}
	}
}

func (puc *PayUsecase) Divide(ctx context.Context, outTradeNo, transactionNo string) (*bby.DivideReplyData, error) {
	conf := puc.pconf.Bby.DjAccount

	outDivideNo := strconv.FormatUint(uint64(time.Now().UnixMilli()), 10)
	nonceStr := strconv.FormatUint(uint64(time.Now().UnixMilli()), 10)

	data, err := bby.Divide(outTradeNo, transactionNo, outDivideNo, nonceStr, "20002381", 197248, conf)
	fmt.Println("#########################")
	fmt.Println(err)
	fmt.Println(outDivideNo)
	fmt.Println(data)
	fmt.Println(data.Status)
	fmt.Println(data.ResultCode)
	fmt.Println(data.DivideState)
	fmt.Println("#########################")
	return nil, nil
}

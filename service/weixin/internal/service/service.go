package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	v1 "weixin/api/weixin/v1"
	"weixin/internal/biz"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewWeixinService)

type WeixinService struct {
	v1.UnimplementedWeixinServer

	uuc    *biz.UserUsecase
	uoduc  *biz.UserOpenDouyinUsecase
	qcuc   *biz.QrCodeUsecase
	uauc   *biz.UserAddressUsecase
	usouc  *biz.UserSampleOrderUsecase
	usruc  *biz.UserScanRecordUsecase
	uouc   *biz.UserOrderUsecase
	uoruc  *biz.UserOrganizationRelationUsecase
	ucuc   *biz.UserCommissionUsecase
	ucouc  *biz.UserCouponUsecase
	ubuc   *biz.UserBalanceUsecase
	uconuc *biz.UserContractUsecase
	ubauc  *biz.UserBankUsecase

	log *log.Helper
}

func NewWeixinService(uuc *biz.UserUsecase, uoduc *biz.UserOpenDouyinUsecase, qcuc *biz.QrCodeUsecase, uauc *biz.UserAddressUsecase, usouc *biz.UserSampleOrderUsecase, usruc *biz.UserScanRecordUsecase, uouc *biz.UserOrderUsecase, uoruc *biz.UserOrganizationRelationUsecase, ucuc *biz.UserCommissionUsecase, ucouc *biz.UserCouponUsecase, ubuc *biz.UserBalanceUsecase, uconuc *biz.UserContractUsecase, ubauc *biz.UserBankUsecase, logger log.Logger) *WeixinService {
	log := log.NewHelper(log.With(logger, "module", "service/douyin"))

	return &WeixinService{uuc: uuc, uoduc: uoduc, qcuc: qcuc, uauc: uauc, usouc: usouc, usruc: usruc, uouc: uouc, uoruc: uoruc, ucuc: ucuc, ucouc: ucouc, ubuc: ubuc, uconuc: uconuc, ubauc: ubauc, log: log}
}

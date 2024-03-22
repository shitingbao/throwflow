package service

import (
	"task/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewTaskService)

type TaskService struct {
	oatuc  *biz.OceanengineAccountTokenUsecase
	odtuc  *biz.OpenDouyinTokenUsecase
	odvuc  *biz.OpenDouyinVideoUsecase
	qauc   *biz.QianchuanAdvertiserUsecase
	qaduc  *biz.QianchuanAdUsecase
	cuc    *biz.CompanyUsecase
	couc   *biz.CompanyOrganizationUsecase
	ctuc   *biz.CompanyTaskUsecase
	jsuc   *biz.JinritemaiStoreUsecase
	jouc   *biz.JinritemaiOrderUsecase
	wuuc   *biz.WeixinUserUsecase
	wuouc  *biz.WeixinUserOrganizationUsecase
	wucuc  *biz.WeixinUserCommissionUsecase
	wucouc *biz.WeixinUserCouponUsecase
	wubuc  *biz.WeixinUserBalanceUsecase
	dtuc   *biz.DoukeTokenUsecase
	log    *log.Helper
}

func NewTaskService(logger log.Logger, oatuc *biz.OceanengineAccountTokenUsecase, odtuc *biz.OpenDouyinTokenUsecase, odvuc *biz.OpenDouyinVideoUsecase, qauc *biz.QianchuanAdvertiserUsecase, qaduc *biz.QianchuanAdUsecase, cuc *biz.CompanyUsecase, couc *biz.CompanyOrganizationUsecase, ctuc *biz.CompanyTaskUsecase, jsuc *biz.JinritemaiStoreUsecase, jouc *biz.JinritemaiOrderUsecase, wuuc *biz.WeixinUserUsecase, wuouc *biz.WeixinUserOrganizationUsecase, wucuc *biz.WeixinUserCommissionUsecase, wucouc *biz.WeixinUserCouponUsecase, wubuc *biz.WeixinUserBalanceUsecase, dtuc *biz.DoukeTokenUsecase) *TaskService {
	log := log.NewHelper(log.With(logger, "module", "service/task"))

	task := &TaskService{
		oatuc:  oatuc,
		odtuc:  odtuc,
		odvuc:  odvuc,
		qauc:   qauc,
		qaduc:  qaduc,
		cuc:    cuc,
		couc:   couc,
		ctuc:   ctuc,
		jsuc:   jsuc,
		jouc:   jouc,
		wuuc:   wuuc,
		wuouc:  wuouc,
		wucuc:  wucuc,
		wucouc: wucouc,
		wubuc:  wubuc,
		dtuc:   dtuc,

		log: log,
	}

	return task
}

var DefaultJobs map[string]JobFunc

type JobFunc func()

func (ts *TaskService) Init() {
	DefaultJobs = map[string]JobFunc{
		// "RefreshOceanengineAccountTokens":      ts.RefreshOceanengineAccountTokens,
		// "RefreshOpenDouyinTokens":              ts.RefreshOpenDouyinTokens,
		// "RenewRefreshTokensOpenDouyinTokens":   ts.RenewRefreshTokensOpenDouyinTokens,
		// "SyncQianchuanDatas":                   ts.SyncQianchuanDatas,
		// "SyncRdsDatas":                         ts.SyncRdsDatas,
		// "SyncQianchuanAds":                     ts.SyncQianchuanAds,
		// "SyncUpdateStatusCompanys":             ts.SyncUpdateStatusCompanys,
		// "SyncUpdateQrCodeCompanyOrganizations": ts.SyncUpdateQrCodeCompanyOrganizations,
		// "SyncJinritemaiStores":                 ts.SyncJinritemaiStores,
		// "SyncOpenDouyinVideos":                 ts.SyncOpenDouyinVideos,
		// "SyncUserFansOpenDouyinTokens":         ts.SyncUserFansOpenDouyinTokens,
		// "Sync90DayJinritemaiOrders":            ts.Sync90DayJinritemaiOrders,
		// "SyncIntegralUsers":                    ts.SyncIntegralUsers,
		// "SyncUserOrganizationRelations":        ts.SyncUserOrganizationRelations,
		// "SyncOrderUserCommissions":             ts.SyncOrderUserCommissions,
		// "SyncCostOrderUserCommissions":         ts.SyncCostOrderUserCommissions,
		// "SyncUserCoupons":                      ts.SyncUserCoupons,
		// "RefreshDoukeTokens":                   ts.RefreshDoukeTokens,
		// "SyncUserBalances":                     ts.SyncUserBalances,
		"SyncCompanyTasks": ts.SyncCompanyTasks,
	}
}

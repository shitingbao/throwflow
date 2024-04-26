package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"task/internal/biz"
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
	cpuc   *biz.CompanyProductUsecase
	jsuc   *biz.JinritemaiStoreUsecase
	jouc   *biz.JinritemaiOrderUsecase
	wuuc   *biz.WeixinUserUsecase
	wucuc  *biz.WeixinUserCommissionUsecase
	wucouc *biz.WeixinUserCouponUsecase
	wubuc  *biz.WeixinUserBalanceUsecase
	douc   *biz.DoukeOrderUsecase
	log    *log.Helper
}

func NewTaskService(logger log.Logger, oatuc *biz.OceanengineAccountTokenUsecase, odtuc *biz.OpenDouyinTokenUsecase, odvuc *biz.OpenDouyinVideoUsecase, qauc *biz.QianchuanAdvertiserUsecase, qaduc *biz.QianchuanAdUsecase, cuc *biz.CompanyUsecase, couc *biz.CompanyOrganizationUsecase, ctuc *biz.CompanyTaskUsecase, cpuc *biz.CompanyProductUsecase, jsuc *biz.JinritemaiStoreUsecase, jouc *biz.JinritemaiOrderUsecase, wuuc *biz.WeixinUserUsecase, wucuc *biz.WeixinUserCommissionUsecase, wucouc *biz.WeixinUserCouponUsecase, wubuc *biz.WeixinUserBalanceUsecase, douc *biz.DoukeOrderUsecase) *TaskService {
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
		cpuc:   cpuc,
		jsuc:   jsuc,
		jouc:   jouc,
		wuuc:   wuuc,
		wucuc:  wucuc,
		wucouc: wucouc,
		wubuc:  wubuc,
		douc:   douc,

		log: log,
	}

	return task
}

var DefaultJobs map[string]JobFunc

type JobFunc func()

func (ts *TaskService) Init() {
	DefaultJobs = map[string]JobFunc{
		"RefreshOceanengineAccountTokens":      ts.RefreshOceanengineAccountTokens,
		"RefreshOpenDouyinTokens":              ts.RefreshOpenDouyinTokens,
		"RenewRefreshTokensOpenDouyinTokens":   ts.RenewRefreshTokensOpenDouyinTokens,
		"SyncQianchuanDatas":                   ts.SyncQianchuanDatas,
		"SyncRdsDatas":                         ts.SyncRdsDatas,
		"SyncQianchuanAds":                     ts.SyncQianchuanAds,
		"SyncUpdateStatusCompanys":             ts.SyncUpdateStatusCompanys,
		"SyncUpdateQrCodeCompanyOrganizations": ts.SyncUpdateQrCodeCompanyOrganizations,
		"SyncJinritemaiStores":                 ts.SyncJinritemaiStores,
		"SyncOpenDouyinVideos":                 ts.SyncOpenDouyinVideos,
		"SyncUserFansOpenDouyinTokens":         ts.SyncUserFansOpenDouyinTokens,
		"Sync90DayJinritemaiOrders":            ts.Sync90DayJinritemaiOrders,
		"SyncIntegralUsers":                    ts.SyncIntegralUsers,
		"SyncTaskUserCommissions":              ts.SyncTaskUserCommissions,
		"SyncUserCoupons":                      ts.SyncUserCoupons,
		"SyncDoukeOrders":                      ts.SyncDoukeOrders,
		"SyncUserBalances":                     ts.SyncUserBalances,
		"SyncCompanyTasks":                     ts.SyncCompanyTasks,
		"SyncCompanyProducts":                  ts.SyncCompanyProducts,
	}
}

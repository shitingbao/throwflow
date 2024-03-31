package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"weixin/internal/conf"
	"weixin/internal/domain"
)

var (
	WeixinUserScanRecordNotFound    = errors.NotFound("WEIXIN_USER_SCAN_RECORD_NOT_FOUND", "微信用户扫码记录不存在")
	WeixinUserScanRecordCreateError = errors.InternalServer("WEIXIN_USER_SCAN_RECORD_CREATE_ERROR", "微信用户扫码记录创建失败")
)

type UserScanRecordRepo interface {
	Get(context.Context, uint64, uint64, uint8) (*domain.UserScanRecord, error)
	Save(context.Context, *domain.UserScanRecord) (*domain.UserScanRecord, error)
}

type UserScanRecordUsecase struct {
	repo   UserScanRecordRepo
	urepo  UserRepo
	corepo CompanyOrganizationRepo
	tm     Transaction
	conf   *conf.Data
	log    *log.Helper
}

func NewUserScanRecordUsecase(repo UserScanRecordRepo, urepo UserRepo, corepo CompanyOrganizationRepo, tm Transaction, conf *conf.Data, logger log.Logger) *UserScanRecordUsecase {
	return &UserScanRecordUsecase{repo: repo, urepo: urepo, corepo: corepo, tm: tm, conf: conf, log: log.NewHelper(logger)}
}

func (usruc *UserScanRecordUsecase) CreateUserScanRecords(ctx context.Context, userId, organizationId, parentUserId uint64) error {
	user, err := usruc.urepo.Get(ctx, userId)

	if err != nil {
		return WeixinLoginError
	}

	if userId == parentUserId {
		return WeixinUserScanRecordCreateError
	}

	companyOrganization, err := usruc.corepo.Get(ctx, organizationId)

	if err != nil {
		return WeixinCompanyOrganizationNotFound
	}

	var inUserScanRecord *domain.UserScanRecord

	if parentUserId > 0 {
		weixinUser, err := usruc.urepo.Get(ctx, parentUserId)

		if err != nil {
			return WeixinUserNotFound

		}
		inUserScanRecord = domain.NewUserScanRecord(ctx, user.Id, companyOrganization.Data.OrganizationId, weixinUser.Id, 0)
	} else {
		inUserScanRecord = domain.NewUserScanRecord(ctx, user.Id, companyOrganization.Data.OrganizationId, 0, 0)
	}

	inUserScanRecord.SetCreateTime(ctx)
	inUserScanRecord.SetUpdateTime(ctx)

	if _, err := usruc.repo.Save(ctx, inUserScanRecord); err != nil {
		return WeixinUserScanRecordCreateError
	}

	return nil
}

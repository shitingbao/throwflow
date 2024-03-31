package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"weixin/internal/conf"
	"weixin/internal/domain"
	"weixin/internal/pkg/gongmall"
	"weixin/internal/pkg/gongmall/employee"
	"weixin/internal/pkg/tool"
)

var (
	WeixinUserContractNotFound        = errors.NotFound("WEIXIN_USER_CONTRACT_NOT_FOUND", "微信用户电签合同签署不存在")
	WeixinUserContractGetError        = errors.InternalServer("WEIXIN_USER_CONTRACT_GET_ERROR", "微信用户电签合同签署合同获取失败")
	WeixinUserContractRsaEncryptError = errors.InternalServer("WEIXIN_USER_CONTRACT_RSA_ENCRYPT_ERROR", "微信用户敏感信息加密失败")
	WeixinUserContractRsaDecryptError = errors.InternalServer("WEIXIN_USER_CONTRACT_RSA_DECRYPT_ERROR", "微信用户敏感信息解密失败")
	WeixinUserContractSignError       = errors.InternalServer("WEIXIN_USER_CONTRACT_SIGN_ERROR", "微信用户工猫回调数据验签失败")
	WeixinUserContractCreateError     = errors.InternalServer("WEIXIN_USER_CONTRACT_CREATE_ERROR", "微信用户电签合同签署创建失败")
	WeixinUserContractConfirmError    = errors.InternalServer("WEIXIN_USER_CONTRACT_CONFIRM_ERROR", "微信用户电签合同签署确认失败")
	WeixinUserContractVerifyError     = errors.InternalServer("WEIXIN_USER_CONTRACT_VERIFY_ERROR", "微信用户电签合同签署用户信息不一致")
	WeixinUserContractExist           = errors.InternalServer("WEIXIN_USER_CONTRACT_EXIST", "微信用户电签合同已签署")
	WeixinUserContractNotExist        = errors.InternalServer("WEIXIN_USER_CONTRACT_NOT_EXIST", "微信用户电签合同不存在")
	WeixinUserContractCompleted       = errors.InternalServer("WEIXIN_USER_CONTRACT_COMPLETED", "微信用户电签合同已签署完成")
)

type UserContractRepo interface {
	Get(context.Context, uint64, string) (*domain.UserContract, error)
	GetByContractId(context.Context, uint64) (*domain.UserContract, error)
	GetByIdentityCardMark(context.Context, uint8, string) (*domain.UserContract, error)
	Save(context.Context, *domain.UserContract) (*domain.UserContract, error)
	Update(context.Context, *domain.UserContract) (*domain.UserContract, error)
}

type UserContractUsecase struct {
	repo    UserContractRepo
	urepo   UserRepo
	ubrepo  UserBankRepo
	uorrepo UserOrganizationRelationRepo
	crepo   CompanyRepo
	corepo  CompanyOrganizationRepo
	tm      Transaction
	conf    *conf.Data
	oconf   *conf.Organization
	gconf   *conf.Gongmall
	log     *log.Helper
}

func NewUserContractUsecase(repo UserContractRepo, urepo UserRepo, ubrepo UserBankRepo, uorrepo UserOrganizationRelationRepo, crepo CompanyRepo, corepo CompanyOrganizationRepo, tm Transaction, conf *conf.Data, oconf *conf.Organization, gconf *conf.Gongmall, logger log.Logger) *UserContractUsecase {
	return &UserContractUsecase{repo: repo, urepo: urepo, ubrepo: ubrepo, uorrepo: uorrepo, crepo: crepo, corepo: corepo, tm: tm, conf: conf, oconf: oconf, gconf: gconf, log: log.NewHelper(logger)}
}

func (ucuc *UserContractUsecase) GetContractUserContracts(ctx context.Context, userId uint64, contractType uint8) (*employee.ListContractDataResponse, error) {
	user, err := ucuc.urepo.Get(ctx, userId)

	if err != nil {
		return nil, WeixinLoginError
	}

	var gongmallConf *conf.Gongmall_Gongmall

	if contractType == 1 {
		gongmallConf = ucuc.gconf.Default
	} else if contractType == 2 {
		userOrganizationRelation, err := ucuc.uorrepo.GetByUserId(ctx, user.Id, 0, 0, "0")

		if err != nil {
			return nil, WeixinUserOrganizationRelationNotFound
		}

		if ucuc.oconf.DjOrganizationId == userOrganizationRelation.OrganizationId {
			gongmallConf = ucuc.gconf.Dj
		} else if ucuc.oconf.DefaultOrganizationId == userOrganizationRelation.OrganizationId {
			gongmallConf = ucuc.gconf.Default
		} else {
			return nil, WeixinCompanyNotFound
		}
	}

	return &employee.ListContractDataResponse{
		ContractAddr: gongmallConf.ContraceUrl,
	}, nil

	/*list, err := employee.ListContract(gongmallConf.ServiceId, gongmallConf.AppKey, gongmallConf.AppSecret)

	if err != nil {
		return nil, WeixinUserContractGetError
	}

	var contract *employee.ListContractDataResponse

	for _, l := range list.Data {
		if l.TemplateId == strconv.FormatUint(gongmallConf.TemplateId, 10) {
			contract = &l

			break
		}
	}

	if contract == nil {
		return nil, WeixinUserContractGetError
	}

	return contract, nil*/
}

func (ucuc *UserContractUsecase) GetUserContracts(ctx context.Context, userId uint64, contractType uint8) (*domain.UserContract, error) {
	user, err := ucuc.urepo.Get(ctx, userId)

	if err != nil {
		return nil, WeixinLoginError
	}

	if len(user.IdentityCardMark) == 0 {
		return nil, WeixinUserContractNotExist
	}

	userContract, err := ucuc.repo.GetByIdentityCardMark(ctx, contractType, user.IdentityCardMark)

	if err != nil {
		return nil, WeixinUserContractNotFound
	}

	return userContract, nil
}

func (ucuc *UserContractUsecase) CreateUserContracts(ctx context.Context, userId uint64, contractType uint8, name, phone, identityCard string) (*domain.UserContract, error) {
	inUser, err := ucuc.urepo.Get(ctx, userId)

	if err != nil {
		return nil, WeixinLoginError
	}

	identityCardMark := identityCard[0:6] + identityCard[:4] + tool.GetMd5(identityCard)

	var gongmallConf *conf.Gongmall_Gongmall
	var organizationId uint64

	if contractType == 1 {
		organizationId = 0
		gongmallConf = ucuc.gconf.Default
	} else if contractType == 2 {
		userOrganizationRelation, err := ucuc.uorrepo.GetByUserId(ctx, inUser.Id, 0, 0, "0")

		if err != nil {
			return nil, WeixinUserOrganizationRelationNotFound
		}

		if _, err := ucuc.corepo.Get(ctx, userOrganizationRelation.OrganizationId); err != nil {
			return nil, WeixinCompanyOrganizationNotFound
		}

		organizationId = userOrganizationRelation.OrganizationId

		if ucuc.oconf.DjOrganizationId == userOrganizationRelation.OrganizationId {
			gongmallConf = ucuc.gconf.Dj
		} else if ucuc.oconf.DefaultOrganizationId == userOrganizationRelation.OrganizationId {
			gongmallConf = ucuc.gconf.Default
		} else {
			return nil, WeixinCompanyNotFound
		}
	}

	fmt.Println("####################################")
	fmt.Println(contractType)
	fmt.Println(organizationId)
	fmt.Println(gongmallConf)
	fmt.Println("####################################")

	enName, err := gongmall.RsaEncrypt(gongmallConf.GongmallPublicKey, name)

	if err != nil {
		return nil, WeixinUserContractRsaEncryptError
	}

	enMengMaName, err := gongmall.RsaEncrypt(gongmallConf.PublicKey, name)

	if err != nil {
		return nil, WeixinUserContractRsaEncryptError
	}

	enPhone, err := gongmall.RsaEncrypt(gongmallConf.GongmallPublicKey, phone)

	if err != nil {
		return nil, WeixinUserContractRsaEncryptError
	}

	enIdentityCard, err := gongmall.RsaEncrypt(gongmallConf.GongmallPublicKey, identityCard)

	if err != nil {
		return nil, WeixinUserContractRsaEncryptError
	}

	enMengMaIdentityCard, err := gongmall.RsaEncrypt(gongmallConf.PublicKey, identityCard)

	if err != nil {
		return nil, WeixinUserContractRsaEncryptError
	}

	if len(inUser.IdentityCardMark) > 0 {
		if identityCardMark != inUser.IdentityCardMark {
			return nil, WeixinUserContractVerifyError
		}
	}

	var userContract *domain.UserContract

	err = ucuc.tm.InTx(ctx, func(ctx context.Context) error {
		if userContract, err = ucuc.repo.Get(ctx, organizationId, identityCardMark); err != nil {
			syncInfo, err := employee.SyncInfo(gongmallConf.ServiceId, gongmallConf.TemplateId, gongmallConf.AppKey, gongmallConf.AppSecret, enName, enPhone, "1", enIdentityCard)

			if err != nil {
				return err
			}

			contractId, err := strconv.ParseUint(syncInfo.Data.ContractId, 10, 64)

			if err != nil {
				return err
			}

			inUserContract := domain.NewUserContract(ctx, organizationId, gongmallConf.ServiceId, gongmallConf.TemplateId, contractId, syncInfo.Data.ProcessStatus, contractType, enMengMaName, enMengMaIdentityCard, identityCardMark)
			inUserContract.SetCreateTime(ctx)
			inUserContract.SetUpdateTime(ctx)

			userContract, err = ucuc.repo.Save(ctx, inUserContract)

			if err != nil {
				return err
			}
		}

		if len(inUser.IdentityCardMark) == 0 {
			inUser.SetIdentityCardMark(ctx, identityCardMark)
			inUser.SetUpdateTime(ctx)

			if _, err := ucuc.urepo.Update(ctx, inUser); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, WeixinUserContractCreateError
	}

	return userContract, nil
}

func (ucuc *UserContractUsecase) ConfirmUserContracts(ctx context.Context, userId, contractId uint64, phone, code string) error {
	if _, err := ucuc.urepo.Get(ctx, userId); err != nil {
		return WeixinLoginError
	}

	inUserContract, err := ucuc.repo.GetByContractId(ctx, contractId)

	if err != nil {
		return WeixinUserContractNotExist
	}

	if inUserContract.ContractStatus == 3 || inUserContract.ContractStatus == 2 {
		return WeixinUserContractCompleted
	}

	var gongmallConf *conf.Gongmall_Gongmall

	if inUserContract.ContractType == 1 {
		gongmallConf = ucuc.gconf.Default
	} else if inUserContract.ContractType == 2 {
		if ucuc.oconf.DjOrganizationId == inUserContract.OrganizationId {
			gongmallConf = ucuc.gconf.Dj
		} else if ucuc.oconf.DefaultOrganizationId == inUserContract.OrganizationId {
			gongmallConf = ucuc.gconf.Default
		} else {
			return WeixinCompanyNotFound
		}
	}

	enPhone, err := gongmall.RsaEncrypt(gongmallConf.GongmallPublicKey, phone)

	if err != nil {
		return WeixinUserContractRsaEncryptError
	}

	if _, err := employee.SignContract(gongmallConf.ServiceId, contractId, gongmallConf.AppKey, gongmallConf.AppSecret, enPhone, code); err != nil {
		return WeixinUserContractConfirmError
	} else {
		inUserContract.SetContractStatus(ctx, 2)
		inUserContract.SetUpdateTime(ctx)

		if _, err := ucuc.repo.Update(ctx, inUserContract); err != nil {
			return WeixinUserContractConfirmError
		}
	}

	return nil
}

func (ucuc *UserContractUsecase) AsyncNotificationUserContracts(ctx context.Context, content string) error {
	contractAsyncNotificationData, err := gongmall.ContractAsyncNotification(content)

	if err != nil {
		return err
	}

	inUserContract, err := ucuc.repo.GetByContractId(ctx, contractAsyncNotificationData.ContractId)

	if err != nil {
		return WeixinUserContractNotFound
	}

	var gongmallConf *conf.Gongmall_Gongmall

	if ucuc.oconf.DjOrganizationId == inUserContract.OrganizationId {
		gongmallConf = ucuc.gconf.Dj
	} else if ucuc.oconf.DefaultOrganizationId == inUserContract.OrganizationId {
		gongmallConf = ucuc.gconf.Default
	} else {
		return WeixinCompanyNotFound
	}

	if ok := gongmall.VerifySign("appKey="+contractAsyncNotificationData.AppKey+"&contractFileUrl="+contractAsyncNotificationData.ContractFileUrl+"&contractId="+strconv.FormatUint(contractAsyncNotificationData.ContractId, 10)+"&identity="+contractAsyncNotificationData.Identity+"&mobile="+contractAsyncNotificationData.Mobile+"&name="+contractAsyncNotificationData.Name+"&nonce="+contractAsyncNotificationData.Nonce+"&serviceId="+strconv.FormatUint(contractAsyncNotificationData.ServiceId, 10)+"&timestamp="+strconv.FormatUint(contractAsyncNotificationData.Timestamp, 10), gongmallConf.AppSecret, contractAsyncNotificationData.Sign); !ok {
		return WeixinUserContractSignError
	}

	inUserContract.SetName(ctx, contractAsyncNotificationData.Name)
	inUserContract.SetIdentityCard(ctx, contractAsyncNotificationData.Identity)
	inUserContract.SetContractStatus(ctx, 3)
	inUserContract.SetUpdateTime(ctx)

	if _, err := ucuc.repo.Update(ctx, inUserContract); err != nil {
		return WeixinUserContractConfirmError
	}

	return nil
}

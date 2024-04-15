package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"weixin/internal/conf"
	"weixin/internal/domain"
	"weixin/internal/pkg/gongmall"
	"weixin/internal/pkg/gongmall/employee"
)

var (
	WeixinUserBankExist       = errors.InternalServer("WEIXIN_USER_BANK_EXIST", "微信用户银行卡已存在")
	WeixinUserBankNotExist    = errors.InternalServer("WEIXIN_USER_BANK_NOT_EXIST", "微信用户银行卡不存在")
	WeixinUserBankListError   = errors.InternalServer("WEIXIN_USER_BANK_LIST_ERROR", "微信用户银行卡列表获取失败")
	WeixinUserBankCreateError = errors.InternalServer("WEIXIN_USER_BANK_CREATE_ERROR", "微信用户银行卡创建失败")
	WeixinUserBankDeleteError = errors.InternalServer("WEIXIN_USER_BANK_DELETE_ERROR", "微信用户银行卡删除失败")
)

type UserBankRepo interface {
	GetByBankCode(context.Context, uint64, string, string) (*domain.UserBank, error)
	List(context.Context, int, int, uint64, string) ([]*domain.UserBank, error)
	Count(context.Context, uint64, string) (int64, error)
	Save(context.Context, *domain.UserBank) (*domain.UserBank, error)
	Update(context.Context, *domain.UserBank) (*domain.UserBank, error)
	Delete(context.Context, *domain.UserBank) error
}

type UserBankUsecase struct {
	repo   UserBankRepo
	urepo  UserRepo
	ucrepo UserContractRepo
	tm     Transaction
	conf   *conf.Data
	oconf  *conf.Organization
	gconf  *conf.Gongmall
	log    *log.Helper
}

func NewUserBankUsecase(repo UserBankRepo, urepo UserRepo, ucrepo UserContractRepo, tm Transaction, conf *conf.Data, oconf *conf.Organization, gconf *conf.Gongmall, logger log.Logger) *UserBankUsecase {
	return &UserBankUsecase{repo: repo, urepo: urepo, ucrepo: ucrepo, tm: tm, conf: conf, oconf: oconf, gconf: gconf, log: log.NewHelper(logger)}
}

func (ubuc *UserBankUsecase) ListUserBanks(ctx context.Context, pageNum, pageSize, userId uint64) (*domain.UserBankList, error) {
	user, err := ubuc.urepo.Get(ctx, userId)

	if err != nil {
		return nil, WeixinLoginError
	}

	list, err := ubuc.repo.List(ctx, int(pageNum), int(pageSize), user.Id, user.IdentityCardMark)

	if err != nil {
		return nil, WeixinUserBankListError
	}

	total, err := ubuc.repo.Count(ctx, user.Id, user.IdentityCardMark)

	if err != nil {
		return nil, WeixinUserBankListError
	}

	return &domain.UserBankList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (ubuc *UserBankUsecase) CreateUserBanks(ctx context.Context, userId uint64, bankCode string) (*domain.UserBank, error) {
	user, err := ubuc.urepo.Get(ctx, userId)

	if err != nil {
		return nil, WeixinLoginError
	}

	if len(user.IdentityCardMark) == 0 {
		return nil, WeixinUserContractNotExist
	}

	if _, err := ubuc.repo.GetByBankCode(ctx, user.Id, user.IdentityCardMark, bankCode); err == nil {
		return nil, WeixinUserBankExist
	}

	userContract, err := ubuc.ucrepo.GetByIdentityCardMark(ctx, 0, user.IdentityCardMark)

	if err != nil {
		return nil, WeixinUserContractNotExist
	}

	var gongmallConf *conf.Gongmall_Gongmall

	if userContract.OrganizationId == 0 {
		gongmallConf = ubuc.gconf.Default
	} else {
		if ubuc.oconf.DjOrganizationId == userContract.OrganizationId {
			gongmallConf = ubuc.gconf.Dj
		} else if ubuc.oconf.DefaultOrganizationId == userContract.OrganizationId {
			gongmallConf = ubuc.gconf.Default
		} else if ubuc.oconf.LbOrganizationId == userContract.OrganizationId {
			gongmallConf = ubuc.gconf.Lb
		} else {
			return nil, WeixinCompanyNotFound
		}
	}
	fmt.Println("######################################")
	fmt.Println(userContract)
	fmt.Println(gongmallConf)
	fmt.Println("######################################")
	enBankCode, err := gongmall.RsaEncrypt(gongmallConf.GongmallPublicKey, bankCode)

	if err != nil {
		return nil, WeixinUserContractRsaEncryptError
	}

	name, err := gongmall.RsaDecrypt(gongmallConf.PrivateKey, userContract.Name)

	if err != nil {
		fmt.Println("####################")
		fmt.Println(userContract.Name)
		fmt.Println(1)
		fmt.Println("####################")
		return nil, WeixinUserContractRsaDecryptError
	}

	enName, err := gongmall.RsaEncrypt(gongmallConf.GongmallPublicKey, name)

	if err != nil {
		return nil, WeixinUserContractRsaEncryptError
	}

	identityCard, err := gongmall.RsaDecrypt(gongmallConf.PrivateKey, userContract.IdentityCard)

	if err != nil {
		fmt.Println("####################")
		fmt.Println(userContract.IdentityCard)
		fmt.Println(2)
		fmt.Println("####################")
		return nil, WeixinUserContractRsaDecryptError
	}

	enIdentityCard, err := gongmall.RsaEncrypt(gongmallConf.GongmallPublicKey, identityCard)

	if err != nil {
		return nil, WeixinUserContractRsaEncryptError
	}

	bank, err := employee.AddBankAccount(gongmallConf.AppKey, gongmallConf.AppSecret, enName, enIdentityCard, enBankCode)

	if err != nil {
		return nil, errors.InternalServer("WEIXIN_USER_BANK_CREATE_ERROR", err.Error())
	}

	deBankCode, err := gongmall.RsaDecrypt(gongmallConf.PrivateKey, bank.Data.BankAccountNo)

	if err != nil {
		fmt.Println("####################")
		fmt.Println(bank.Data.BankAccountNo)
		fmt.Println(3)
		fmt.Println("####################")
		return nil, WeixinUserContractRsaDecryptError
	}

	inUserBank := domain.NewUserBank(ctx, user.Id, user.IdentityCardMark, deBankCode, bank.Data.BankName)
	inUserBank.SetCreateTime(ctx)
	inUserBank.SetUpdateTime(ctx)

	userBank, err := ubuc.repo.Save(ctx, inUserBank)

	if err != nil {
		return nil, WeixinUserBankCreateError
	}

	return userBank, nil
}

func (ubuc *UserBankUsecase) DeleteUserBanks(ctx context.Context, userId uint64, bankCode string) error {
	user, err := ubuc.urepo.Get(ctx, userId)

	if err != nil {
		return WeixinLoginError
	}

	if len(user.IdentityCardMark) == 0 {
		return WeixinUserContractNotExist
	}

	inUserBank, err := ubuc.repo.GetByBankCode(ctx, user.Id, user.IdentityCardMark, bankCode)

	if err != nil {
		return WeixinUserBankNotExist
	}

	if err := ubuc.repo.Delete(ctx, inUserBank); err != nil {
		return WeixinUserAddressDeleteError
	}

	return nil
}

func (ubuc *UserBankUsecase) DecryptDatas(ctx context.Context, organizationId uint64, content string) (string, error) {
	var gongmallConf *conf.Gongmall_Gongmall

	if ubuc.oconf.DjOrganizationId == organizationId {
		gongmallConf = ubuc.gconf.Dj
	} else if ubuc.oconf.DefaultOrganizationId == organizationId {
		gongmallConf = ubuc.gconf.Default
	} else if ubuc.oconf.LbOrganizationId == organizationId {
		gongmallConf = ubuc.gconf.Lb
	} else {
		return "", WeixinCompanyNotFound
	}

	content, err := gongmall.RsaDecrypt(gongmallConf.PrivateKey, content)

	fmt.Println("#################################")
	fmt.Println(err)
	fmt.Println("#################################")

	if err != nil {
		return "", WeixinUserContractRsaDecryptError
	}

	return content, nil
}

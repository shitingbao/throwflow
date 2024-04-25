package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"strings"
	"sync"
	"time"
	"weixin/internal/conf"
	"weixin/internal/domain"
	"weixin/internal/pkg/gongmall"
	"weixin/internal/pkg/gongmall/merchant"
	"weixin/internal/pkg/tool"
)

var (
	WeixinUserBalanceGetError          = errors.NotFound("WEIXIN_USER_BALANCE_GET_ERROR", "微信用户余额获取失败")
	WeixinUserBalanceInsufficientError = errors.NotFound("WEIXIN_USER_BALANCE_INSUFFICIENT_ERROR", "微信用户余额不足")
	WeixinUserBalanceCreateError       = errors.NotFound("WEIXIN_USER_BALANCE_CREATE_ERROR", "微信用户余额提现失败")
)

type UserBalanceUsecase struct {
	urepo   UserRepo
	ubrepo  UserBankRepo
	ublrepo UserBalanceLogRepo
	ucrepo  UserContractRepo
	ucorepo UserCommissionRepo
	tlrepo  TaskLogRepo
	tm      Transaction
	conf    *conf.Data
	oconf   *conf.Organization
	gconf   *conf.Gongmall
	log     *log.Helper
}

func NewUserBalanceUsecase(urepo UserRepo, ubrepo UserBankRepo, ublrepo UserBalanceLogRepo, ucrepo UserContractRepo, ucorepo UserCommissionRepo, tlrepo TaskLogRepo, tm Transaction, conf *conf.Data, oconf *conf.Organization, gconf *conf.Gongmall, logger log.Logger) *UserBalanceUsecase {
	return &UserBalanceUsecase{urepo: urepo, ubrepo: ubrepo, ublrepo: ublrepo, ucrepo: ucrepo, ucorepo: ucorepo, tlrepo: tlrepo, tm: tm, conf: conf, oconf: oconf, gconf: gconf, log: log.NewHelper(logger)}
}

func (ubuc *UserBalanceUsecase) GetUserBalances(ctx context.Context, userId uint64) (*domain.UserBalance, error) {
	if _, err := ubuc.urepo.Get(ctx, userId); err != nil {
		return nil, WeixinLoginError
	}

	var wg sync.WaitGroup

	var userEstimatedCommissionIncomeBalance, userSettleCommissionIncomeBalance, userCashableCommissionIncomeBalance, userCommissionExpenseBalance, userEstimatedCommissionCostIncomeBalance, userSettleCommissionCostIncomeBalance, userCashableCommissionCostIncomeBalance, userCommissionCostExpenseBalance *domain.UserCommission

	wg.Add(8)

	// 会员/带货-待结算收入
	go func() {
		defer wg.Done()

		userEstimatedCommissionIncomeBalance, _ = ubuc.ucorepo.Statistics(ctx, userId, 1, 1, []uint8{1, 2}, []uint8{1})
	}()

	// 会员/带货-已结算收入
	go func() {
		defer wg.Done()

		userSettleCommissionIncomeBalance, _ = ubuc.ucorepo.Statistics(ctx, userId, 2, 1, []uint8{1, 2}, []uint8{1})
	}()

	// 会员/带货-可提现收入
	go func() {
		defer wg.Done()

		userCashableCommissionIncomeBalance, _ = ubuc.ucorepo.Statistics(ctx, userId, 3, 1, []uint8{1, 2}, []uint8{1})
	}()

	// 会员/带货-提现
	go func() {
		defer wg.Done()

		userCommissionExpenseBalance, _ = ubuc.ucorepo.Statistics(ctx, userId, 0, 2, []uint8{7}, []uint8{0, 1, 2})
	}()

	// 成本购/任务-待结算收入
	go func() {
		defer wg.Done()

		userEstimatedCommissionCostIncomeBalance, _ = ubuc.ucorepo.Statistics(ctx, userId, 1, 1, []uint8{3, 4, 5, 6}, []uint8{1})
	}()

	// 成本购/任务-已结算收入
	go func() {
		defer wg.Done()

		userSettleCommissionCostIncomeBalance, _ = ubuc.ucorepo.Statistics(ctx, userId, 2, 1, []uint8{3, 4, 5, 6}, []uint8{1})
	}()

	// 成本购/任务-可提现收入
	go func() {
		defer wg.Done()

		userCashableCommissionCostIncomeBalance, _ = ubuc.ucorepo.Statistics(ctx, userId, 3, 1, []uint8{3, 4, 5, 6}, []uint8{1})
	}()

	// 成本购/任务-提现
	go func() {
		defer wg.Done()

		userCommissionCostExpenseBalance, _ = ubuc.ucorepo.Statistics(ctx, userId, 0, 2, []uint8{8}, []uint8{0, 1, 2})
	}()

	wg.Wait()

	return &domain.UserBalance{
		EstimatedCommissionBalance:     userEstimatedCommissionIncomeBalance.CommissionAmount,
		SettleCommissionBalance:        userSettleCommissionIncomeBalance.CommissionAmount,
		CashableCommissionBalance:      userCashableCommissionIncomeBalance.CommissionAmount - userCommissionExpenseBalance.CommissionAmount,
		EstimatedCommissionCostBalance: userEstimatedCommissionCostIncomeBalance.CommissionAmount,
		SettleCommissionCostBalance:    userSettleCommissionCostIncomeBalance.CommissionAmount,
		CashableCommissionCostBalance:  userCashableCommissionCostIncomeBalance.CommissionAmount - userCommissionCostExpenseBalance.CommissionAmount,
	}, nil
}

func (ubuc *UserBalanceUsecase) ListUserBalances(ctx context.Context, pageNum, pageSize, userId uint64, operationType uint8, keyword string) (*domain.UserCommissionList, error) {
	user, err := ubuc.urepo.Get(ctx, userId)

	if err != nil {
		return nil, WeixinLoginError
	}

	list, err := ubuc.ucorepo.ListBalance(ctx, int(pageNum), int(pageSize), user.Id, operationType, keyword)

	if err != nil {
		return nil, WeixinUserBalanceLogListError
	}

	for _, l := range list {
		l.SetBalanceCommissionTypeName(ctx)
	}

	total, err := ubuc.ucorepo.CountBalance(ctx, user.Id, operationType, keyword)

	if err != nil {
		return nil, WeixinUserBalanceLogListError
	}

	return &domain.UserCommissionList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (ubuc *UserBalanceUsecase) CreateUserBalances(ctx context.Context, userId uint64, balanceType uint8, bankCode string, amount float64) error {
	user, err := ubuc.urepo.Get(ctx, userId)

	if err != nil {
		return WeixinLoginError
	}

	if _, err := ubuc.ubrepo.GetByBankCode(ctx, user.Id, user.IdentityCardMark, bankCode); err != nil {
		return WeixinUserBankNotExist
	}

	incomeCommissionTypes := make([]uint8, 0)
	expenseCommissionTypes := make([]uint8, 0)
	var scommissionType uint8
	var userContract *domain.UserContract

	if balanceType == 1 {
		userContract, err = ubuc.ucrepo.GetByIdentityCardMark(ctx, 1, user.IdentityCardMark)

		if err != nil {
			return WeixinUserContractNotExist
		}

		scommissionType = 8

		incomeCommissionTypes = append(incomeCommissionTypes, 3)
		incomeCommissionTypes = append(incomeCommissionTypes, 4)
		incomeCommissionTypes = append(incomeCommissionTypes, 5)
		incomeCommissionTypes = append(incomeCommissionTypes, 6)
		expenseCommissionTypes = append(expenseCommissionTypes, 8)
	} else if balanceType == 2 {
		userContract, err = ubuc.ucrepo.GetByIdentityCardMark(ctx, 2, user.IdentityCardMark)

		if err != nil {
			return WeixinUserContractNotExist
		}

		scommissionType = 7

		incomeCommissionTypes = append(incomeCommissionTypes, 1)
		incomeCommissionTypes = append(incomeCommissionTypes, 2)
		expenseCommissionTypes = append(expenseCommissionTypes, 7)
	} else {
		return WeixinValidatorError
	}

	isFail := true

	for num := 0; num <= 1; num++ {
		result, err := ubuc.ublrepo.SaveCacheString(ctx, "weixin:balance:log:"+user.IdentityCardMark+strconv.FormatUint(uint64(balanceType), 10), "go", ubuc.conf.Redis.UserBalanceLogTimeout.AsDuration())

		if err != nil {
			return WeixinUserBalanceCreateError
		}

		if result {
			userIncomeBalance, err := ubuc.ucorepo.Statistics(ctx, user.Id, 3, 1, incomeCommissionTypes, []uint8{1})

			if err != nil {
				return WeixinUserBalanceGetError
			}

			userExpenseBalance, err := ubuc.ucorepo.Statistics(ctx, user.Id, 0, 2, expenseCommissionTypes, []uint8{0, 1, 2})

			if err != nil {
				return WeixinUserBalanceGetError
			}

			userBalance := userIncomeBalance.CommissionAmount - userExpenseBalance.CommissionAmount

			if float32(tool.Decimal(float64(userBalance), 2)) < float32(amount) {
				return WeixinUserBalanceInsufficientError
			}

			inUserCommission := domain.NewUserCommission(ctx, user.Id, userContract.OrganizationId, 0, 0, 0, 0, 0, 0, scommissionType, 2, 0, 0.00, 0.00, float32(tool.Decimal(amount, 2)))
			inUserCommission.SetName(ctx, userContract.Name)
			inUserCommission.SetIdentityCard(ctx, userContract.IdentityCard)
			inUserCommission.SetBankCode(ctx, bankCode)
			inUserCommission.SetCreateTime(ctx, time.Now())
			inUserCommission.SetUpdateTime(ctx)

			_, err = ubuc.ucorepo.Save(ctx, inUserCommission)

			ubuc.ublrepo.DeleteCache(ctx, "weixin:balance:log:"+user.IdentityCardMark+strconv.FormatUint(uint64(balanceType), 10))

			if err != nil {
				return WeixinUserBalanceCreateError
			}

			isFail = false

			break
		} else {
			time.Sleep(200 * time.Millisecond)

			continue
		}
	}

	if isFail {
		return WeixinUserBalanceCreateError
	}

	return nil
}

func (ubuc *UserBalanceUsecase) AsyncNotificationUserBalances(ctx context.Context, content string) error {
	balanceAsyncNotificationData, err := gongmall.BalanceAsyncNotificationSuccess(content)

	if err != nil {
		return err
	}

	inUserCommission, err := ubuc.ucorepo.GetByOutTradeNo(ctx, balanceAsyncNotificationData.RequestId)

	if err != nil {
		return WeixinUserBalanceLogNotFound
	}

	var gongmallConf *conf.Gongmall_Gongmall

	if inUserCommission.OrganizationId == 0 {
		gongmallConf = ubuc.gconf.Default
	} else if inUserCommission.OrganizationId > 0 {
		if ubuc.oconf.DjOrganizationId == inUserCommission.OrganizationId {
			gongmallConf = ubuc.gconf.Dj
		} else if ubuc.oconf.DefaultOrganizationId == inUserCommission.OrganizationId {
			gongmallConf = ubuc.gconf.Default
		} else if ubuc.oconf.LbOrganizationId == inUserCommission.OrganizationId {
			gongmallConf = ubuc.gconf.Lb
		} else {
			return WeixinCompanyNotFound
		}
	}

	if balanceAsyncNotificationData.Status == 30 {
		balanceAsyncNotificationFailData, err := gongmall.BalanceAsyncNotificationFail(content)

		if err != nil {
			return err
		}

		balanceAsyncNotificationDataMap := tool.StructToMap(balanceAsyncNotificationFailData)
		balanceAsyncNotificationDataKeys := tool.SortMapKeys(balanceAsyncNotificationDataMap)

		var paramPattern []string

		for _, balanceAsyncNotificationDataKey := range balanceAsyncNotificationDataKeys {
			if balanceAsyncNotificationDataKey != "sign" {
				if balanceAsyncNotificationDataKey == "amount" ||
					balanceAsyncNotificationDataKey == "currentRealWage" ||
					balanceAsyncNotificationDataKey == "currentTax" ||
					balanceAsyncNotificationDataKey == "currentManageFee" ||
					balanceAsyncNotificationDataKey == "currentAddTax" ||
					balanceAsyncNotificationDataKey == "currentAddValueTax" {

					balanceAsyncNotificationDataValue, _ := balanceAsyncNotificationDataMap[balanceAsyncNotificationDataKey].(float32)

					paramPattern = append(paramPattern, fmt.Sprintf("%s=%v", balanceAsyncNotificationDataKey, fmt.Sprintf("%.2f", tool.Decimal(float64(balanceAsyncNotificationDataValue), 2))))
				} else {
					paramPattern = append(paramPattern, fmt.Sprintf("%s=%v", balanceAsyncNotificationDataKey, balanceAsyncNotificationDataMap[balanceAsyncNotificationDataKey]))
				}
			}
		}

		if ok := gongmall.VerifySign(strings.Join(paramPattern, "&"), gongmallConf.AppSecret, balanceAsyncNotificationFailData.Sign); !ok {
			return WeixinUserContractRsaDecryptError
		}

		inUserCommission.SetBalanceStatus(ctx, 3)
		inUserCommission.SetOperationContent(ctx, balanceAsyncNotificationFailData.FailReason)
		inUserCommission.SetUpdateTime(ctx)
	} else if balanceAsyncNotificationData.Status == 20 {
		balanceAsyncNotificationDataMap := tool.StructToMap(balanceAsyncNotificationData)
		balanceAsyncNotificationDataKeys := tool.SortMapKeys(balanceAsyncNotificationDataMap)

		var paramPattern []string

		for _, balanceAsyncNotificationDataKey := range balanceAsyncNotificationDataKeys {
			if balanceAsyncNotificationDataKey != "sign" {
				if balanceAsyncNotificationDataKey == "amount" ||
					balanceAsyncNotificationDataKey == "currentRealWage" ||
					balanceAsyncNotificationDataKey == "currentTax" ||
					balanceAsyncNotificationDataKey == "currentManageFee" ||
					balanceAsyncNotificationDataKey == "currentAddTax" ||
					balanceAsyncNotificationDataKey == "currentAddValueTax" {

					balanceAsyncNotificationDataValue, _ := balanceAsyncNotificationDataMap[balanceAsyncNotificationDataKey].(float32)

					paramPattern = append(paramPattern, fmt.Sprintf("%s=%v", balanceAsyncNotificationDataKey, fmt.Sprintf("%.2f", tool.Decimal(float64(balanceAsyncNotificationDataValue), 2))))
				} else {
					paramPattern = append(paramPattern, fmt.Sprintf("%s=%v", balanceAsyncNotificationDataKey, balanceAsyncNotificationDataMap[balanceAsyncNotificationDataKey]))
				}
			}
		}
		fmt.Println(strings.Join(paramPattern, "&"))
		if ok := gongmall.VerifySign(strings.Join(paramPattern, "&"), gongmallConf.AppSecret, balanceAsyncNotificationData.Sign); !ok {
			return WeixinUserContractRsaDecryptError
		}

		dateTime, err := tool.StringToTime("20060102150405", balanceAsyncNotificationData.DateTime)

		if err != nil {
			return WeixinUserBalanceLogAsyncNotificationError
		}

		createTime, err := tool.StringToTime("20060102150405", balanceAsyncNotificationData.CreateTime)

		if err != nil {
			return WeixinUserBalanceLogAsyncNotificationError
		}

		payTime, err := tool.StringToTime("20060102150405", balanceAsyncNotificationData.PayTime)

		if err != nil {
			return WeixinUserBalanceLogAsyncNotificationError
		}

		inUserCommission.SetBalanceStatus(ctx, 1)
		inUserCommission.SetInnerTradeNo(ctx, balanceAsyncNotificationData.InnerTradeNo)
		inUserCommission.SetRealAmount(ctx, balanceAsyncNotificationData.Amount)
		inUserCommission.SetApplyTime(ctx, &dateTime)
		inUserCommission.SetGongmallCreateTime(ctx, &createTime)
		inUserCommission.SetPayTime(ctx, &payTime)
		inUserCommission.SetUpdateTime(ctx)
	} else {
		return WeixinUserBalanceLogAsyncNotificationError
	}

	if _, err := ubuc.ucorepo.Update(ctx, inUserCommission); err != nil {
		return WeixinUserBalanceLogAsyncNotificationError
	}

	return nil
}

func (ubuc *UserBalanceUsecase) SyncUserBalances(ctx context.Context) error {
	userCommissions, err := ubuc.ucorepo.ListCashable(ctx)

	if err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "SyncUserBalances", fmt.Sprintf("[SyncUserBalancesError] Description=%s", "获取微信用户提现列表失败"))
		inTaskLog.SetCreateTime(ctx)

		ubuc.tlrepo.Save(ctx, inTaskLog)

		return err
	}

	for _, inUserCommission := range userCommissions {
		user, err := ubuc.urepo.Get(ctx, inUserCommission.UserId)

		if err != nil {
			continue
		}

		var gongmallConf *conf.Gongmall_Gongmall
		incomeCommissionTypes := make([]uint8, 0)
		expenseCommissionTypes := make([]uint8, 0)

		if inUserCommission.OrganizationId == 0 {
			gongmallConf = ubuc.gconf.Default

			incomeCommissionTypes = append(incomeCommissionTypes, 3)
			incomeCommissionTypes = append(incomeCommissionTypes, 4)
			incomeCommissionTypes = append(incomeCommissionTypes, 5)
			incomeCommissionTypes = append(incomeCommissionTypes, 6)
			expenseCommissionTypes = append(expenseCommissionTypes, 8)
		} else if inUserCommission.OrganizationId > 0 {
			if ubuc.oconf.DjOrganizationId == inUserCommission.OrganizationId {
				gongmallConf = ubuc.gconf.Dj
			} else if ubuc.oconf.DefaultOrganizationId == inUserCommission.OrganizationId {
				gongmallConf = ubuc.gconf.Default
			} else if ubuc.oconf.LbOrganizationId == inUserCommission.OrganizationId {
				gongmallConf = ubuc.gconf.Lb
			} else {
				continue
			}

			incomeCommissionTypes = append(incomeCommissionTypes, 1)
			incomeCommissionTypes = append(incomeCommissionTypes, 2)
			expenseCommissionTypes = append(expenseCommissionTypes, 7)
		}

		userIncomeBalance, err := ubuc.ucorepo.Statistics(ctx, user.Id, 3, 1, incomeCommissionTypes, []uint8{1})

		if err != nil {
			continue
		}

		userExpenseBalance, err := ubuc.ucorepo.Statistics(ctx, user.Id, 0, 2, expenseCommissionTypes, []uint8{1, 2})

		if err != nil {
			continue
		}

		payAmount := userIncomeBalance.CommissionAmount - userExpenseBalance.CommissionAmount

		if float32(tool.Decimal(float64(payAmount), 2)) < float32(inUserCommission.CommissionAmount) {
			inUserCommission.SetBalanceStatus(ctx, 3)
			inUserCommission.SetOperationContent(ctx, "账户余额不足")
			inUserCommission.SetUpdateTime(ctx)

			ubuc.ucorepo.Update(ctx, inUserCommission)
		} else {
			enPhone, err := gongmall.RsaEncrypt(gongmallConf.GongmallPublicKey, user.Phone)

			if err != nil {
				continue
			}

			enBankCode, err := gongmall.RsaEncrypt(gongmallConf.GongmallPublicKey, inUserCommission.BankCode)

			if err != nil {
				continue
			}

			name, err := gongmall.RsaDecrypt(gongmallConf.PrivateKey, inUserCommission.Name)

			if err != nil {
				continue
			}

			enName, err := gongmall.RsaEncrypt(gongmallConf.GongmallPublicKey, name)

			if err != nil {
				continue
			}

			identityCard, err := gongmall.RsaDecrypt(gongmallConf.PrivateKey, inUserCommission.IdentityCard)

			if err != nil {
				continue
			}

			enIdentityCard, err := gongmall.RsaEncrypt(gongmallConf.GongmallPublicKey, identityCard)

			if err != nil {
				continue
			}

			requestId, err := ubuc.ublrepo.NextId(ctx)

			if err != nil {
				continue
			}

			if _, err := merchant.DoSinglePayment(gongmallConf.ServiceId, gongmallConf.AppKey, gongmallConf.AppSecret, enName, enPhone, enIdentityCard, enBankCode, fmt.Sprintf("%.2f", tool.Decimal(float64(inUserCommission.CommissionAmount*0.935-3), 2)), time.Now().Format("20060102150405"), strconv.FormatUint(requestId, 10)); err != nil {
				inUserCommission.SetBalanceStatus(ctx, 3)
				inUserCommission.SetOperationContent(ctx, err.Error())
				inUserCommission.SetUpdateTime(ctx)

				ubuc.ucorepo.Update(ctx, inUserCommission)
			} else {
				inUserCommission.SetOutTradeNo(ctx, strconv.FormatUint(requestId, 10))
				inUserCommission.SetBalanceStatus(ctx, 2)
				inUserCommission.SetUpdateTime(ctx)

				ubuc.ucorepo.Update(ctx, inUserCommission)
			}
		}
	}

	return nil
}

package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"strings"
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
	user, err := ubuc.urepo.Get(ctx, userId)

	if err != nil {
		return nil, WeixinLoginError
	}

	// 成本购-收入
	userCostIncomeBalance, err := ubuc.ublrepo.Statistics(ctx, user.Id, 1, []string{"3", "4"}, []string{"1"})

	if err != nil {
		return nil, WeixinUserBalanceGetError
	}

	// 成本购-提现
	userCostExpenseBalance, err := ubuc.ublrepo.Statistics(ctx, user.Id, 2, []string{"6"}, []string{"0", "1", "2"})

	if err != nil {
		return nil, WeixinUserBalanceGetError
	}

	// 分佣-收入
	userCommissionIncomeBalance, err := ubuc.ublrepo.Statistics(ctx, user.Id, 1, []string{"1", "2"}, []string{"1"})

	if err != nil {
		return nil, WeixinUserBalanceGetError
	}

	// 分佣-提现
	userCommissionExpenseBalance, err := ubuc.ublrepo.Statistics(ctx, user.Id, 2, []string{"5"}, []string{"0", "1", "2"})

	if err != nil {
		return nil, WeixinUserBalanceGetError
	}

	uiday, _ := strconv.ParseUint(time.Now().Format("20060102"), 10, 64)

	// 电商-预估佣金
	statisticEstimatedOrder, err := ubuc.ucorepo.Statistics(ctx, user.Id, uint32(uiday), 2)

	if err != nil {
		return nil, WeixinUserBalanceGetError
	}

	// 电商-结算佣金
	/*statisticRealOrder, err := ubuc.ucorepo.StatisticsReal(ctx, user.Id, uint32(uiday), 2)

	if err != nil {
		return nil, WeixinUserBalanceGetError
	}*/

	// 成本购-预估佣金
	statisticEstimatedCostOrder, err := ubuc.ucorepo.Statistics(ctx, user.Id, uint32(uiday), 3)

	if err != nil {
		return nil, WeixinUserBalanceGetError
	}

	// 成本购-结算佣金
	statisticRealCostOrder, err := ubuc.ucorepo.StatisticsReal(ctx, user.Id, uint32(uiday), 3)

	if err != nil {
		return nil, WeixinUserBalanceGetError
	}

	return &domain.UserBalance{
		EstimatedCommissionBalance: statisticEstimatedOrder.EstimatedUserCommission,
		RealCommissionBalance:      userCommissionIncomeBalance.Amount - userCommissionExpenseBalance.Amount,
		EstimatedCostBalance:       statisticEstimatedCostOrder.EstimatedUserCommission - statisticRealCostOrder.RealUserCommission,
		RealCostBalance:            userCostIncomeBalance.Amount - userCostExpenseBalance.Amount,
	}, nil
}

func (ubuc *UserBalanceUsecase) ListUserBalances(ctx context.Context, pageNum, pageSize, userId uint64, operationType uint8) (*domain.UserBalanceLogList, error) {
	user, err := ubuc.urepo.Get(ctx, userId)

	if err != nil {
		return nil, WeixinLoginError
	}

	balanceTypes := make([]string, 0)
	balanceStatuses := make([]string, 0)

	list, err := ubuc.ublrepo.List(ctx, int(pageNum), int(pageSize), user.Id, operationType, "DESC", balanceTypes, balanceStatuses)

	if err != nil {
		return nil, WeixinUserBalanceLogListError
	}

	total, err := ubuc.ublrepo.Count(ctx, user.Id, operationType, balanceTypes, balanceStatuses)

	if err != nil {
		return nil, WeixinUserBalanceLogListError
	}

	return &domain.UserBalanceLogList{
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

	incomeBalanceTypes := make([]string, 0)
	expenseBalanceTypes := make([]string, 0)
	var sbalanceType uint8
	var userContract *domain.UserContract

	if balanceType == 1 {
		userContract, err = ubuc.ucrepo.GetByIdentityCardMark(ctx, 1, user.IdentityCardMark)

		if err != nil {
			return WeixinUserContractNotExist
		}

		sbalanceType = 6

		incomeBalanceTypes = append(incomeBalanceTypes, "3")
		incomeBalanceTypes = append(incomeBalanceTypes, "4")
		expenseBalanceTypes = append(incomeBalanceTypes, "6")
	} else if balanceType == 2 {
		userContract, err = ubuc.ucrepo.GetByIdentityCardMark(ctx, 2, user.IdentityCardMark)

		if err != nil {
			return WeixinUserContractNotExist
		}

		sbalanceType = 5

		incomeBalanceTypes = append(incomeBalanceTypes, "1")
		incomeBalanceTypes = append(incomeBalanceTypes, "2")
		expenseBalanceTypes = append(incomeBalanceTypes, "5")
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
			userIncomeBalance, err := ubuc.ublrepo.Statistics(ctx, user.Id, 1, incomeBalanceTypes, []string{"1"})

			if err != nil {
				return WeixinUserBalanceGetError
			}

			userExpenseBalance, err := ubuc.ublrepo.Statistics(ctx, user.Id, 2, expenseBalanceTypes, []string{"0", "1", "2"})

			if err != nil {
				return WeixinUserBalanceGetError
			}

			userBalance := userIncomeBalance.Amount - userExpenseBalance.Amount

			if userBalance < float32(amount) {
				return WeixinUserBalanceInsufficientError
			}

			inUserBalanceLog := domain.NewUserBalanceLog(ctx, user.Id, 0, sbalanceType, 2, 0, float32(tool.Decimal(amount, 2)), "提现")
			inUserBalanceLog.SetOrganizationId(ctx, userContract.OrganizationId)
			inUserBalanceLog.SetName(ctx, userContract.Name)
			inUserBalanceLog.SetIdentityCard(ctx, userContract.IdentityCard)
			inUserBalanceLog.SetBankCode(ctx, bankCode)
			inUserBalanceLog.SetCreateTime(ctx)
			inUserBalanceLog.SetUpdateTime(ctx)

			_, err = ubuc.ublrepo.Save(ctx, inUserBalanceLog)

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

	inUserBalanceLog, err := ubuc.ublrepo.GetByOutTradeNo(ctx, balanceAsyncNotificationData.RequestId)

	if err != nil {
		return WeixinUserBalanceLogNotFound
	}

	var gongmallConf *conf.Gongmall_Gongmall

	if inUserBalanceLog.OrganizationId == 0 {
		gongmallConf = ubuc.gconf.Default
	} else if inUserBalanceLog.OrganizationId > 0 {
		if ubuc.oconf.DjOrganizationId == inUserBalanceLog.OrganizationId {
			gongmallConf = ubuc.gconf.Dj
		} else if ubuc.oconf.DefaultOrganizationId == inUserBalanceLog.OrganizationId {
			gongmallConf = ubuc.gconf.Default
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

		inUserBalanceLog.SetBalanceStatus(ctx, 3)
		inUserBalanceLog.SetOperationContent(ctx, balanceAsyncNotificationFailData.FailReason)
		inUserBalanceLog.SetUpdateTime(ctx)
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

		inUserBalanceLog.SetBalanceStatus(ctx, 1)
		inUserBalanceLog.SetInnerTradeNo(ctx, balanceAsyncNotificationData.InnerTradeNo)
		inUserBalanceLog.SetRealAmount(ctx, balanceAsyncNotificationData.Amount)
		inUserBalanceLog.SetApplyTime(ctx, &dateTime)
		inUserBalanceLog.SetGongmallCreateTime(ctx, &createTime)
		inUserBalanceLog.SetPayTime(ctx, &payTime)
		inUserBalanceLog.SetUpdateTime(ctx)
	} else {
		return WeixinUserBalanceLogAsyncNotificationError
	}

	if _, err := ubuc.ublrepo.Update(ctx, inUserBalanceLog); err != nil {
		return WeixinUserBalanceLogAsyncNotificationError
	}

	return nil
}

func (ubuc *UserBalanceUsecase) SyncUserBalances(ctx context.Context) error {
	userBalances, err := ubuc.ublrepo.List(ctx, 0, 40, 0, 2, "ASC", []string{"5", "6"}, []string{"0"})

	if err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "SyncUserBalances", fmt.Sprintf("[SyncUserBalancesError] Description=%s", "获取微信用户提现列表失败"))
		inTaskLog.SetCreateTime(ctx)

		ubuc.tlrepo.Save(ctx, inTaskLog)

		return err
	}

	for _, inUserBalance := range userBalances {
		user, err := ubuc.urepo.Get(ctx, inUserBalance.UserId)

		if err != nil {
			continue
		}

		var gongmallConf *conf.Gongmall_Gongmall
		incomeBalanceTypes := make([]string, 0)
		expenseBalanceTypes := make([]string, 0)

		if inUserBalance.OrganizationId == 0 {
			gongmallConf = ubuc.gconf.Default

			incomeBalanceTypes = append(incomeBalanceTypes, "3")
			incomeBalanceTypes = append(incomeBalanceTypes, "4")
			expenseBalanceTypes = append(expenseBalanceTypes, "6")
		} else if inUserBalance.OrganizationId > 0 {
			if ubuc.oconf.DjOrganizationId == inUserBalance.OrganizationId {
				gongmallConf = ubuc.gconf.Dj
			} else if ubuc.oconf.DefaultOrganizationId == inUserBalance.OrganizationId {
				gongmallConf = ubuc.gconf.Default
			} else {
				continue
			}

			incomeBalanceTypes = append(incomeBalanceTypes, "1")
			incomeBalanceTypes = append(incomeBalanceTypes, "2")
			expenseBalanceTypes = append(expenseBalanceTypes, "5")
		}

		userIncomeBalance, err := ubuc.ublrepo.Statistics(ctx, user.Id, 1, incomeBalanceTypes, []string{"1"})

		if err != nil {
			continue
		}

		userExpenseBalance, err := ubuc.ublrepo.Statistics(ctx, user.Id, 2, expenseBalanceTypes, []string{"1", "2"})

		if err != nil {
			continue
		}

		payAmount := userIncomeBalance.Amount - userExpenseBalance.Amount

		if payAmount < float32(inUserBalance.Amount) {
			inUserBalance.SetBalanceStatus(ctx, 3)
			inUserBalance.SetOperationContent(ctx, "账户余额不足")
			inUserBalance.SetUpdateTime(ctx)

			ubuc.ublrepo.Update(ctx, inUserBalance)
		} else {
			enPhone, err := gongmall.RsaEncrypt(gongmallConf.GongmallPublicKey, user.Phone)

			if err != nil {
				continue
			}

			enBankCode, err := gongmall.RsaEncrypt(gongmallConf.GongmallPublicKey, inUserBalance.BankCode)

			if err != nil {
				continue
			}

			name, err := gongmall.RsaDecrypt(gongmallConf.PrivateKey, inUserBalance.Name)

			if err != nil {
				continue
			}

			enName, err := gongmall.RsaEncrypt(gongmallConf.GongmallPublicKey, name)

			if err != nil {
				continue
			}

			identityCard, err := gongmall.RsaDecrypt(gongmallConf.PrivateKey, inUserBalance.IdentityCard)

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

			if _, err := merchant.DoSinglePayment(gongmallConf.ServiceId, gongmallConf.AppKey, gongmallConf.AppSecret, enName, enPhone, enIdentityCard, enBankCode, fmt.Sprintf("%.2f", tool.Decimal(float64(inUserBalance.Amount*0.935-3), 2)), time.Now().Format("20060102150405"), strconv.FormatUint(requestId, 10)); err != nil {
				inUserBalance.SetBalanceStatus(ctx, 3)
				inUserBalance.SetOperationContent(ctx, err.Error())
				inUserBalance.SetUpdateTime(ctx)

				ubuc.ublrepo.Update(ctx, inUserBalance)
			} else {
				inUserBalance.SetOutTradeNo(ctx, strconv.FormatUint(requestId, 10))
				inUserBalance.SetBalanceStatus(ctx, 2)
				inUserBalance.SetUpdateTime(ctx)

				ubuc.ublrepo.Update(ctx, inUserBalance)
			}
		}
	}

	return nil
}

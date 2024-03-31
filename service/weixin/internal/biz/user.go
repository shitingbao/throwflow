package biz

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	ctos "github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"io"
	"strconv"
	"strings"
	"time"
	"weixin/internal/conf"
	"weixin/internal/domain"
	"weixin/internal/pkg/mini/account"
	"weixin/internal/pkg/mini/oauth2"
	"weixin/internal/pkg/mini/wxa"
	"weixin/internal/pkg/tool"
)

var (
	WeixinUserNotFound            = errors.NotFound("WEIXIN_USER_NOT_FOUND", "微信用户不存在")
	WeixinUserToOrganizationError = errors.NotFound("WEIXIN_USER_TO_ORGANIZATION_ERROR", "微信用户和机构不匹配")
	WeixinUserCreateError         = errors.InternalServer("WEIXIN_USER_CREATE_ERROR", "微信用户创建失败")
	WeixinUserUpdateError         = errors.InternalServer("WEIXIN_USER_UPDATE_ERROR", "微信用户更新失败")
	WeixinUserFollowCreateError   = errors.InternalServer("WEIXIN_USER_FOLLOW_CREATE_ERROR", "微信用户绑定关系创建失败")
	WeixinLoginError              = errors.InternalServer("WEIXIN_LOGIN_ERROR", "微信用户异常错误")

	Mime = map[string]string{
		"image/jpeg": ".jpeg",
		"image/png":  ".png",
		"image/gif":  ".gif",
	}
)

type UserRepo interface {
	Get(context.Context, uint64) (*domain.User, error)
	GetByPhoneAndCountryCode(context.Context, string, string) (*domain.User, error)
	List(context.Context) ([]*domain.User, error)
	ListRanking(context.Context) ([]*domain.User, error)
	Count(context.Context) (int64, error)
	CountByUserId(context.Context, uint64) (int64, error)
	Update(context.Context, *domain.User) (*domain.User, error)
	UpdateRanking(context.Context, uint64, uint64) error
	Save(context.Context, *domain.User) (*domain.User, error)

	GetCacheHash(context.Context, string, string) (string, error)
	SaveCacheHash(context.Context, string, map[string]string, time.Duration) error

	PutContent(context.Context, string, io.Reader) (*ctos.PutObjectV2Output, error)
	PutContentCode(context.Context, string, io.Reader) (*ctos.PutObjectV2Output, error)

	NextId(context.Context) (uint64, error)
}

type UserUsecase struct {
	repo     UserRepo
	uoirepo  UserOpenIdRepo
	uorerepo UserOrganizationRelationRepo
	uirrepo  UserIntegralRelationRepo
	usrrepo  UserScanRecordRepo
	uorrepo  UserOrderRepo
	crepo    CompanyRepo
	cprepo   CompanyProductRepo
	corepo   CompanyOrganizationRepo
	uodrepo  UserOpenDouyinRepo
	jorepo   JinritemaiOrderRepo
	tlrepo   TaskLogRepo
	turepo   TuUserRepo
	curepo   CouponUserRepo
	surepo   ShortUrlRepo
	screpo   ShortCodeRepo
	ucclrepo UserCouponCreateLogRepo
	tm       Transaction
	conf     *conf.Data
	wconf    *conf.Weixin
	vconf    *conf.Volcengine
	oconf    *conf.Organization
	log      *log.Helper
}

func NewUserUsecase(repo UserRepo, uoirepo UserOpenIdRepo, uorerepo UserOrganizationRelationRepo, uirrepo UserIntegralRelationRepo, usrrepo UserScanRecordRepo, uorrepo UserOrderRepo, crepo CompanyRepo, cprepo CompanyProductRepo, corepo CompanyOrganizationRepo, uodrepo UserOpenDouyinRepo, jorepo JinritemaiOrderRepo, tlrepo TaskLogRepo, turepo TuUserRepo, curepo CouponUserRepo, surepo ShortUrlRepo, screpo ShortCodeRepo, ucclrepo UserCouponCreateLogRepo, tm Transaction, conf *conf.Data, wconf *conf.Weixin, vconf *conf.Volcengine, oconf *conf.Organization, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, uoirepo: uoirepo, uorerepo: uorerepo, uirrepo: uirrepo, usrrepo: usrrepo, uorrepo: uorrepo, crepo: crepo, cprepo: cprepo, corepo: corepo, uodrepo: uodrepo, jorepo: jorepo, tlrepo: tlrepo, turepo: turepo, curepo: curepo, surepo: surepo, screpo: screpo, ucclrepo: ucclrepo, tm: tm, conf: conf, wconf: wconf, vconf: vconf, oconf: oconf, log: log.NewHelper(logger)}
}

func (uuc *UserUsecase) GetUsers(ctx context.Context, token string) (*domain.User, error) {
	suserId, err := uuc.repo.GetCacheHash(ctx, "weixin:user:"+token, "userId")

	if err != nil {
		return nil, WeixinLoginError
	}

	userId, err := strconv.ParseUint(suserId, 10, 64)

	if err != nil {
		return nil, WeixinLoginError
	}

	user, err := uuc.getUserById(ctx, userId)

	if err != nil {
		return nil, WeixinLoginError
	}

	return user, nil
}

func (uuc *UserUsecase) GetByIdUsers(ctx context.Context, userId uint64) (*domain.User, error) {
	user, err := uuc.getUserById(ctx, userId)

	if err != nil {
		return nil, WeixinUserNotFound
	}

	return user, nil
}

func (uuc *UserUsecase) GetFollowUsers(ctx context.Context, organizationId, parentUserId uint64) (*domain.FollowData, error) {
	followData := &domain.FollowData{}

	companyOrganization, err := uuc.corepo.Get(ctx, organizationId)

	if err != nil {
		return nil, WeixinCompanyOrganizationNotFound
	}

	if parentUserId > 0 {
		weixinUser, err := uuc.repo.Get(ctx, parentUserId)

		if err != nil {
			return nil, WeixinUserNotFound
		}

		userOrganization, err := uuc.uorerepo.GetByUserId(ctx, weixinUser.Id, organizationId, 0, "0")

		if err != nil {
			return nil, WeixinUserToOrganizationError
		}

		if userOrganization.OrganizationId != companyOrganization.Data.OrganizationId {
			return nil, WeixinUserToOrganizationError
		}

		followData.FollowType = "organizationColonel"
		followData.FollowName = weixinUser.NickName
		followData.FollowLogoUrl = weixinUser.AvatarUrl
		followData.QrCodeUrl = companyOrganization.Data.OrganizationQrCodeUrl
	} else {
		followData.FollowType = "organization"
		followData.FollowName = companyOrganization.Data.OrganizationName
		followData.FollowLogoUrl = companyOrganization.Data.OrganizationLogoUrl
		followData.QrCodeUrl = companyOrganization.Data.OrganizationQrCodeUrl
	}

	totalNum, _ := uuc.repo.Count(ctx)

	followData.TotalNum = uint64(totalNum)

	return followData, nil
}

func (uuc *UserUsecase) CreateUsers(ctx context.Context, organizationId uint64, loginCode, phoneCode string) (*domain.User, error) {
	var code2session *account.Code2SessionResponse
	var phoneNumber *account.GetPhoneNumberResponse
	var err error
	appid := ""

	if organizationId == uuc.oconf.DjOrganizationId {
		appid = uuc.wconf.DjMini.Appid

		code2session, err = account.Code2Session(uuc.wconf.DjMini.Appid, uuc.wconf.DjMini.Secret, loginCode)

		if err != nil {
			return nil, errors.InternalServer("WEIXIN_CREATE_USER_ERROR", err.Error())
		}

		accessToken, err := oauth2.GetAccessToken(uuc.wconf.DjMini.Appid, uuc.wconf.DjMini.Secret)

		if err != nil {
			return nil, errors.InternalServer("WEIXIN_CREATE_USER_ERROR", err.Error())
		}

		phoneNumber, err = account.GetPhoneNumber(accessToken.AccessToken, phoneCode)

		if err != nil {
			return nil, errors.InternalServer("WEIXIN_CREATE_USER_ERROR", err.Error())
		}
	} else {
		appid = uuc.wconf.Mini.Appid

		code2session, err = account.Code2Session(uuc.wconf.Mini.Appid, uuc.wconf.Mini.Secret, loginCode)

		if err != nil {
			return nil, errors.InternalServer("WEIXIN_CREATE_USER_ERROR", err.Error())
		}

		accessToken, err := oauth2.GetAccessToken(uuc.wconf.Mini.Appid, uuc.wconf.Mini.Secret)

		if err != nil {
			return nil, errors.InternalServer("WEIXIN_CREATE_USER_ERROR", err.Error())
		}

		phoneNumber, err = account.GetPhoneNumber(accessToken.AccessToken, phoneCode)

		if err != nil {
			return nil, errors.InternalServer("WEIXIN_CREATE_USER_ERROR", err.Error())
		}
	}

	var tmpUser *domain.User

	if tmpUser, err = uuc.repo.GetByPhoneAndCountryCode(ctx, phoneNumber.PhoneInfo.PurePhoneNumber, phoneNumber.PhoneInfo.CountryCode); err != nil {
		err = uuc.tm.InTx(ctx, func(ctx context.Context) error {
			count, err := uuc.repo.Count(ctx)

			if err != nil {
				return err
			}

			inUser := domain.NewUser(ctx, uint64(count+1), phoneNumber.PhoneInfo.PurePhoneNumber, phoneNumber.PhoneInfo.CountryCode, 0.00)
			inUser.SetCreateTime(ctx)
			inUser.SetUpdateTime(ctx)

			tmpUser, err = uuc.repo.Save(ctx, inUser)

			if err != nil {
				return err
			}

			if inUserOpenId, err := uuc.uoirepo.Get(ctx, tmpUser.Id, appid, code2session.Openid); err != nil {
				inUserOpenId = domain.NewUserOpenId(ctx, tmpUser.Id, appid, code2session.Openid)
				inUserOpenId.SetCreateTime(ctx)
				inUserOpenId.SetUpdateTime(ctx)

				if _, err = uuc.uoirepo.Save(ctx, inUserOpenId); err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			return nil, WeixinUserCreateError
		}
	} else {
		if _, err := uuc.uoirepo.Get(ctx, tmpUser.Id, appid, code2session.Openid); err != nil {
			inUserOpenId := domain.NewUserOpenId(ctx, tmpUser.Id, appid, code2session.Openid)
			inUserOpenId.SetCreateTime(ctx)
			inUserOpenId.SetUpdateTime(ctx)

			if _, err = uuc.uoirepo.Save(ctx, inUserOpenId); err != nil {
				return nil, WeixinUserCreateError
			}
		}
	}

	token := tool.GetToken()

	cacheData := make(map[string]string)
	cacheData["userId"] = strconv.FormatUint(tmpUser.Id, 10)

	if err := uuc.repo.SaveCacheHash(ctx, "weixin:user:"+token, cacheData, uuc.conf.Redis.LoginTokenTimeout.AsDuration()); err != nil {
		return nil, WeixinUserCreateError
	}

	user, err := uuc.getUserById(ctx, tmpUser.Id)

	if err != nil {
		return nil, WeixinUserCreateError
	}

	user.SetToken(ctx, token)

	return user, nil
}

func (uuc *UserUsecase) UpdateUsers(ctx context.Context, userId uint64, nickName, avatar string) (*domain.User, error) {
	inUser, err := uuc.getUserById(ctx, userId)

	if err != nil {
		return nil, WeixinLoginError
	}

	savatar := strings.Split(avatar, ",")

	if len(savatar) != 2 {
		return nil, WeixinUserUpdateError
	}

	if _, ok := Mime[savatar[0][5:len(savatar[0])-7]]; !ok {
		return nil, WeixinUserUpdateError
	}

	imagePath := uuc.vconf.Tos.Avatar.SubFolder + "/" + tool.GetRandCode(time.Now().String()) + Mime[savatar[0][5:len(savatar[0])-7]]
	imageContent, err := base64.StdEncoding.DecodeString(savatar[1])

	if err != nil {
		return nil, WeixinUserUpdateError
	}

	if _, err = uuc.repo.PutContent(ctx, imagePath, strings.NewReader(string(imageContent))); err != nil {
		return nil, WeixinUserUpdateError
	}

	inUser.SetNickName(ctx, nickName)
	inUser.SetAvatarUrl(ctx, uuc.vconf.Tos.Avatar.Url+"/"+imagePath)
	inUser.SetUpdateTime(ctx)

	if _, err := uuc.repo.Update(ctx, inUser); err != nil {
		return nil, WeixinUserUpdateError
	}

	user, err := uuc.getUserById(ctx, inUser.Id)

	if err != nil {
		return nil, WeixinUserUpdateError
	}

	return user, nil
}

func (uuc *UserUsecase) SyncIntegralUsers(ctx context.Context) error {
	users, err := uuc.repo.List(ctx)

	if err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "SyncIntegralUsers", fmt.Sprintf("[SyncIntegralUsersError] Description=%s", "获取微信用户列表失败"))
		inTaskLog.SetCreateTime(ctx)

		uuc.tlrepo.Save(ctx, inTaskLog)

		return err
	}

	for _, inUser := range users {
		startDay := inUser.CreateTime.AddDate(0, 0, -30).Format("2006-01-02")
		endDay := time.Now().Format("2006-01-02")

		if statistics, err := uuc.jorepo.Statistics(ctx, inUser.Id, startDay, endDay); err == nil {
			for _, lstatistics := range statistics.Data.Statistics {
				if lstatistics.Key == "销额" {
					if integral, err := strconv.ParseFloat(lstatistics.Value, 64); err == nil {
						inUser.SetIntegral(ctx, uint64(integral))
						inUser.SetUpdateTime(ctx)

						uuc.repo.Update(ctx, inUser)
					}
				}
			}
		}
	}

	users, err = uuc.repo.ListRanking(ctx)

	if err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "SyncIntegralUsers", fmt.Sprintf("[SyncIntegralUsersError] Description=%s", "获取微信用户排序列表失败"))
		inTaskLog.SetCreateTime(ctx)

		uuc.tlrepo.Save(ctx, inTaskLog)

		return err
	}

	for _, user := range users {
		uuc.repo.UpdateRanking(ctx, user.Id, user.Ranking)
	}

	return nil
}

func (uuc *UserUsecase) ImportDatas(ctx context.Context) error {
	var organizationId uint64 = 6

	tuUsers, err := uuc.turepo.List(ctx)

	if err != nil {
		return err
	}

	accessToken, err := oauth2.GetAccessToken(uuc.wconf.Mini.Appid, uuc.wconf.Mini.Secret)

	if err != nil {
		return err
	}

	err = uuc.tm.InTx(ctx, func(ctx context.Context) error {
		payTime, err := tool.StringToTime("2006-01-02 15:04:05", "2023-09-01 18:18:18")

		if err != nil {
			return err
		}

		for _, tuUser := range tuUsers {
			var user *domain.User

			user, err = uuc.repo.GetByPhoneAndCountryCode(ctx, tuUser.Phone, "86")

			if err != nil {
				count, err := uuc.repo.Count(ctx)

				if err != nil {
					return err
				}

				inUser := domain.NewUser(ctx, uint64(count+1), tuUser.Phone, "86", 0.00)
				inUser.SetCreateTime(ctx)
				inUser.SetUpdateTime(ctx)

				user, err = uuc.repo.Save(ctx, inUser)

				if err != nil {
					return err
				}
			}

			tuUser.UserId = user.Id
		}

		for _, tuUser := range tuUsers {
			var organizationUserId uint64 = 0

			if len(tuUser.ParentId) == 0 {
				organizationUserId = 0
			} else {
				for _, ltuUser := range tuUsers {
					if tuUser.ParentId == ltuUser.Id {
						organizationUserId = ltuUser.UserId

						break
					}
				}
			}

			if userScanRecord, err := uuc.usrrepo.Get(ctx, tuUser.UserId, organizationId, 1); err == nil {
				if userScanRecord.OrganizationUserId != organizationUserId {
					inUserScanRecord := domain.NewUserScanRecord(ctx, tuUser.UserId, organizationId, organizationUserId, 0)
					inUserScanRecord.SetCreateTime(ctx)
					inUserScanRecord.SetUpdateTime(ctx)

					if _, err := uuc.usrrepo.Save(ctx, inUserScanRecord); err != nil {
						return err
					}
				}
			} else {
				inUserScanRecord := domain.NewUserScanRecord(ctx, tuUser.UserId, organizationId, organizationUserId, 0)
				inUserScanRecord.SetCreateTime(ctx)
				inUserScanRecord.SetUpdateTime(ctx)

				if _, err := uuc.usrrepo.Save(ctx, inUserScanRecord); err != nil {
					return err
				}
			}

			var level uint8

			if tuUser.Level == "2" {
				level = 2
			} else if tuUser.Level == "3" {
				level = 3
			} else if tuUser.Level == "4" {
				level = 4
			}

			if inUserOrganizationRelation, err := uuc.uorerepo.GetByUserId(ctx, tuUser.UserId, organizationId, 0, "0"); err == nil {
				content, err := wxa.GetUnlimitedQRCode(accessToken.AccessToken, "oId="+strconv.FormatUint(organizationId, 10)+"&uId="+strconv.FormatUint(tuUser.UserId, 10), uuc.wconf.Mini.QrCodeEnvVersion)

				if err != nil {
					return err
				}

				objectKey := tool.GetRandCode(time.Now().String())

				if _, err := uuc.repo.PutContentCode(ctx, uuc.vconf.Tos.Company.SubFolder+"/"+objectKey+".png", strings.NewReader(content.Buffer)); err != nil {
					return err
				}

				inUserOrganizationRelation.SetOrganizationUserId(ctx, organizationUserId)
				inUserOrganizationRelation.SetOrganizationUserQrCodeUrl(ctx, uuc.vconf.Tos.Company.Url+"/"+uuc.vconf.Tos.Company.SubFolder+"/"+objectKey+".png")
				inUserOrganizationRelation.SetLevel(ctx, level)
				inUserOrganizationRelation.SetUpdateTime(ctx)

				if _, err := uuc.uorerepo.Update(ctx, inUserOrganizationRelation); err != nil {
					return err
				}

				if inUserIntegralRelation, err := uuc.uirrepo.GetByUserId(ctx, tuUser.UserId, organizationId, 0); err == nil {
					if inUserIntegralRelation.OrganizationUserId != organizationUserId {
						inUserIntegralRelation.SetOrganizationUserId(ctx, organizationUserId)
						inUserIntegralRelation.SetLevel(ctx, level)
						inUserIntegralRelation.SetUpdateTime(ctx)

						if _, err := uuc.uirrepo.Update(ctx, inUserIntegralRelation); err != nil {
							return err
						}
					}
				} else {
					inUserIntegralRelation = domain.NewUserIntegralRelation(ctx, tuUser.UserId, organizationId, organizationUserId, level)
					inUserIntegralRelation.SetCreateTime(ctx)
					inUserIntegralRelation.SetUpdateTime(ctx)

					if _, err := uuc.uirrepo.Save(ctx, inUserIntegralRelation); err != nil {
						return err
					}
				}
			} else {
				outTradeNo, err := uuc.repo.NextId(ctx)

				if err != nil {
					return err
				}

				inUserOrder := domain.NewUserOrder(ctx, tuUser.UserId, organizationId, organizationUserId, 0, strconv.FormatUint(outTradeNo, 10), strconv.FormatUint(outTradeNo, 10), strconv.FormatUint(outTradeNo, 10), 1888.00, 1888.00, &payTime, level, 1, 1)
				inUserOrder.SetCreateTime(ctx)
				inUserOrder.SetUpdateTime(ctx)

				if _, err := uuc.uorrepo.Save(ctx, inUserOrder); err != nil {
					return err
				}

				content, err := wxa.GetUnlimitedQRCode(accessToken.AccessToken, "oId="+strconv.FormatUint(organizationId, 10)+"&uId="+strconv.FormatUint(inUserOrder.UserId, 10), uuc.wconf.Mini.QrCodeEnvVersion)

				if err != nil {
					return err
				}

				objectKey := tool.GetRandCode(time.Now().String())

				if _, err := uuc.repo.PutContentCode(ctx, uuc.vconf.Tos.Company.SubFolder+"/"+objectKey+".png", strings.NewReader(content.Buffer)); err != nil {
					return err
				}

				inUserOrganizationRelation = domain.NewUserOrganizationRelation(ctx, tuUser.UserId, organizationId, organizationUserId, 0, level, 0, uuc.vconf.Tos.Company.Url+"/"+uuc.vconf.Tos.Company.SubFolder+"/"+objectKey+".png")
				inUserOrganizationRelation.SetCreateTime(ctx)
				inUserOrganizationRelation.SetUpdateTime(ctx)

				if _, err := uuc.uorerepo.Save(ctx, inUserOrganizationRelation); err != nil {
					return err
				}

				inUserIntegralRelation := domain.NewUserIntegralRelation(ctx, tuUser.UserId, organizationId, organizationUserId, level)
				inUserIntegralRelation.SetCreateTime(ctx)
				inUserIntegralRelation.SetUpdateTime(ctx)

				if _, err := uuc.uirrepo.Save(ctx, inUserIntegralRelation); err != nil {
					return err
				}
			}
		}

		return nil
	})

	return nil
}

func (uuc *UserUsecase) ParentUserDatas(ctx context.Context) error {
	var organizationId uint64 = 6

	userIntegralRelations, err := uuc.uirrepo.List(ctx, organizationId)

	if err != nil {
		fmt.Println(err)
	}

	for _, userIntegralRelation := range userIntegralRelations {
		var childNum uint64 = 0
		fmt.Println(userIntegralRelation.UserId)
		uuc.uirrepo.GetChildNum(ctx, userIntegralRelation.UserId, &childNum, userIntegralRelations)

		fmt.Println("用户ID:" + strconv.FormatUint(userIntegralRelation.UserId, 10) + "，孩子节点的个数：" + strconv.FormatUint(childNum, 10))
	}

	return nil
}

func (uuc *UserUsecase) getUserById(ctx context.Context, userId uint64) (*domain.User, error) {
	user, err := uuc.repo.Get(ctx, userId)

	if err != nil {
		return nil, err
	}

	user.SetIntegralLevelName(ctx)

	user.OpenIds, _ = uuc.uoirepo.List(ctx, userId)

	if total, err := uuc.repo.Count(ctx); err == nil {
		user.SetTotal(ctx, uint64(total))
	}

	return user, nil
}

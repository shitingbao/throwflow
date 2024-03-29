package biz

import (
	"context"
	"douyin/internal/conf"
	"douyin/internal/domain"
	"douyin/internal/pkg/jinritemai/kol"
	"douyin/internal/pkg/openDouyin/js"
	"douyin/internal/pkg/openDouyin/oauth2"
	"douyin/internal/pkg/openDouyin/user"
	"douyin/internal/pkg/tool"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	DouyinOpenDouyinTokenNotFound       = errors.InternalServer("DOUYIN_OPEN_DOUYIN_TOKEN_NOT_FOUND", "授权账户不存在")
	DouyinOpenDouyinTokenExist          = errors.InternalServer("DOUYIN_OPEN_DOUYIN_TOKEN_EXIST", "授权账户已被授权给微信用户")
	DouyinOpenTicketGetError            = errors.InternalServer("DOUYIN_OPEN_TICKET_GET_ERROR", "jsb_ticket获取失败")
	DouyinOpenDouyinTokenCreateError    = errors.InternalServer("DOUYIN_OPEN_DOUYIN_TOKEN_CREATE_ERROR", "授权账户信息创建失败")
	DouyinOpenDouyinTokenConfigNotFound = errors.NotFound("DOUYIN_OPEN_DOUYIN_TOKEN_CONFIG_NOT_FOUND", "应用配置不存在")
	DouyinOpenDouyinTokenListError      = errors.InternalServer("DOUYIN_OPEN_DOUYIN_TOKEN_LIST_ERROR", "获取授权账户列表失败")
	DouyinOpenQrCodeGetError            = errors.InternalServer("DOUYIN_OPEN_QR_CODE_GET_ERROR", "授权二维码获取失败")
	DouyinOpenStateGetError             = errors.InternalServer("DOUYIN_OPEN_QR_CODE_GET_ERROR", "授权状态获取失败")
	DouyinOpenUrlGetError               = errors.InternalServer("DOUYIN_OPEN_URL_GET_ERROR", "授权链接获取失败")
)

type OpenDouyinTokenRepo interface {
	GetByClientKeyAndOpenId(context.Context, string, string) (*domain.OpenDouyinToken, error)
	List(context.Context) ([]*domain.OpenDouyinToken, error)
	ListByClientKeyAndOpenId(context.Context, []*domain.OpenDouyinToken) ([]*domain.OpenDouyinToken, error)
	ListByCreateTime(context.Context, string) ([]*domain.OpenDouyinToken, error)
	Update(context.Context, *domain.OpenDouyinToken) (*domain.OpenDouyinToken, error)
	Save(context.Context, *domain.OpenDouyinToken) (*domain.OpenDouyinToken, error)

	SaveCacheString(context.Context, string, string, time.Duration) (bool, error)
	GetCacheString(context.Context, string) (string, error)

	DeleteCache(context.Context, string) error
}

type OpenDouyinTokenUsecase struct {
	repo       OpenDouyinTokenRepo
	oduirepo   OpenDouyinUserInfoRepo
	oduiclrepo OpenDouyinUserInfoCreateLogRepo
	odalrepo   OpenDouyinApiLogRepo
	jalrepo    JinritemaiApiLogRepo
	wurepo     WeixinUserRepo
	wuodrepo   WeixinUserOpenDouyinRepo
	joirepo    JinritemaiOrderInfoRepo
	tlrepo     TaskLogRepo
	tm         Transaction
	conf       *conf.Data
	dconf      *conf.Developer
	econf      *conf.Event
	log        *log.Helper
}

func NewOpenDouyinTokenUsecase(repo OpenDouyinTokenRepo, oduirepo OpenDouyinUserInfoRepo, oduiclrepo OpenDouyinUserInfoCreateLogRepo, odalrepo OpenDouyinApiLogRepo, jalrepo JinritemaiApiLogRepo, wurepo WeixinUserRepo, wuodrepo WeixinUserOpenDouyinRepo, joirepo JinritemaiOrderInfoRepo, tlrepo TaskLogRepo, tm Transaction, conf *conf.Data, dconf *conf.Developer, econf *conf.Event, logger log.Logger) *OpenDouyinTokenUsecase {
	return &OpenDouyinTokenUsecase{repo: repo, oduirepo: oduirepo, oduiclrepo: oduiclrepo, odalrepo: odalrepo, jalrepo: jalrepo, wurepo: wurepo, wuodrepo: wuodrepo, joirepo: joirepo, tlrepo: tlrepo, tm: tm, conf: conf, dconf: dconf, econf: econf, log: log.NewHelper(logger)}
}

func (odtuc *OpenDouyinTokenUsecase) GetUrlOpenDouyinTokens(ctx context.Context, userId uint64) (string, error) {
	redirectUrl := "https://open.douyin.com/platform/oauth/connect?client_key=CLIENT_KEY&response_type=code&scope=SCOPE&redirect_uri=REDIRECT_URI&state=STATE"

	scope := make([]string, 0)
	scope = append(scope, "user_info")
	scope = append(scope, "renew_refresh_token")
	scope = append(scope, "h5.share")
	scope = append(scope, "alliance.kol.store_manage")
	scope = append(scope, "alliance.picksource.convert")
	scope = append(scope, "alliance.kol.reputation")
	scope = append(scope, "alliance.kol.orders")
	scope = append(scope, "video.list.bind")
	scope = append(scope, "video.data.bind")
	scope = append(scope, "alliance.kol.buyin_id")
	scope = append(scope, "data.external.user")

	state := tool.GetRandCode(time.Now().String())

	redirectUrl = strings.Replace(redirectUrl, "CLIENT_KEY", odtuc.dconf.Aweme.OpenDouyin.ClientKey, -1)
	redirectUrl = strings.Replace(redirectUrl, "REDIRECT_URI", odtuc.dconf.Aweme.OpenDouyin.CallBackUrl, -1)
	redirectUrl = strings.Replace(redirectUrl, "SCOPE", strings.Join(scope, ","), -1)
	redirectUrl = strings.Replace(redirectUrl, "STATE", state, -1)

	stateParam := make(map[string]string)
	stateParam["userId"] = strconv.FormatUint(userId, 10)
	stateParam["clientKey"] = odtuc.dconf.Aweme.OpenDouyin.ClientKey

	sstateParam, _ := json.Marshal(stateParam)

	if _, err := odtuc.repo.SaveCacheString(ctx, "douyin:openDouyin:auth:"+state, string(sstateParam), odtuc.conf.Redis.OpenDouyinUrlTokenTimeout.AsDuration()); err != nil {
		return "", DouyinOpenUrlGetError
	}

	return redirectUrl, nil
}

func (odtuc *OpenDouyinTokenUsecase) GetTicketOpenDouyinTokens(ctx context.Context) (*js.GetTicketResponse, error) {
	clientToken, err := oauth2.ClientToken(odtuc.dconf.Aweme.OpenDouyin.ClientKey, odtuc.dconf.Aweme.OpenDouyin.ClientSecret)

	if err != nil {
		return nil, DouyinOpenTicketGetError
	}

	ticket, err := js.GetTicket(clientToken.Data.AccessToken)

	if err != nil {
		return nil, DouyinOpenTicketGetError
	}

	return ticket, nil
}

func (odtuc *OpenDouyinTokenUsecase) GetQrCodeOpenDouyinTokens(ctx context.Context, userId uint64) (*domain.QrCode, error) {
	qrCodeUrl := "https://open.douyin.com/oauth/get_qrcode/?&client_key=CLIENT_KEY&scope=SCOPE&next=NEXT&state=STATE&jump_type=native&optional_scope_check=&optional_scope_uncheck=&customize_params={\"comment_id\":\"\",\"source\":\"pc_auth\",\"not_skip_confirm\":\"true\"}"

	scope := make([]string, 0)
	scope = append(scope, "user_info")
	scope = append(scope, "renew_refresh_token")
	scope = append(scope, "h5.share")
	scope = append(scope, "alliance.kol.store_manage")
	scope = append(scope, "alliance.picksource.convert")
	scope = append(scope, "alliance.kol.reputation")
	scope = append(scope, "alliance.kol.orders")
	scope = append(scope, "video.list.bind")
	scope = append(scope, "video.data.bind")
	scope = append(scope, "alliance.kol.buyin_id")
	scope = append(scope, "data.external.user")

	state := tool.GetRandCode(time.Now().String())

	qrCodeUrl = strings.Replace(qrCodeUrl, "CLIENT_KEY", odtuc.dconf.Aweme.OpenDouyin.ClientKey, -1)
	qrCodeUrl = strings.Replace(qrCodeUrl, "NEXT", odtuc.dconf.Aweme.OpenDouyin.CallBackUrl, -1)
	qrCodeUrl = strings.Replace(qrCodeUrl, "SCOPE", strings.Join(scope, ","), -1)
	qrCodeUrl = strings.Replace(qrCodeUrl, "STATE", state, -1)

	response, err := http.Get(qrCodeUrl)
	defer response.Body.Close()

	if err != nil {
		return nil, DouyinOpenQrCodeGetError
	}

	rbody, _ := ioutil.ReadAll(response.Body)

	var qrCodeBody *domain.QrCodeOpenDouyinToken

	if err := json.Unmarshal(rbody, &qrCodeBody); err != nil {
		return nil, DouyinOpenQrCodeGetError
	}

	if qrCodeBody.Message != "success" {
		return nil, DouyinOpenQrCodeGetError
	}

	stateParam := make(map[string]string)
	stateParam["userId"] = strconv.FormatUint(userId, 10)
	stateParam["clientKey"] = odtuc.dconf.Aweme.OpenDouyin.ClientKey

	sstateParam, _ := json.Marshal(stateParam)

	if _, err := odtuc.repo.SaveCacheString(ctx, "douyin:openDouyin:auth:"+state, string(sstateParam), odtuc.conf.Redis.QrCodeTokenTimeout.AsDuration()); err != nil {
		return nil, DouyinOpenQrCodeGetError
	}

	return &domain.QrCode{
		QrCode: "data:image/png;base64," + qrCodeBody.Data.Qrcode,
		State:  state,
		Token:  qrCodeBody.Data.Token,
	}, nil
}

func (odtuc *OpenDouyinTokenUsecase) GetStatusOpenDouyinTokens(ctx context.Context, token, state, timestamp string) (*domain.QrCodeStatus, error) {
	checkQrCodeUrl := "https://open.douyin.com/oauth/check_qrcode/?client_key=CLIENT_KEY&scope=SCOPE&next=NEXT&state=STATE&token=TOKEN&timestamp=TIMESTAMP&jump_type=native&optional_scope_check=&optional_scope_uncheck=&customize_params={\"comment_id\":\"\",\"source\":\"pc_auth\",\"not_skip_confirm\":\"true\"}"

	scope := make([]string, 0)
	scope = append(scope, "user_info")
	scope = append(scope, "renew_refresh_token")
	scope = append(scope, "h5.share")
	scope = append(scope, "alliance.kol.store_manage")
	scope = append(scope, "alliance.picksource.convert")
	scope = append(scope, "alliance.kol.reputation")
	scope = append(scope, "alliance.kol.orders")
	scope = append(scope, "video.list.bind")
	scope = append(scope, "video.data.bind")
	scope = append(scope, "alliance.kol.buyin_id")
	scope = append(scope, "data.external.user")

	checkQrCodeUrl = strings.Replace(checkQrCodeUrl, "CLIENT_KEY", odtuc.dconf.Aweme.OpenDouyin.ClientKey, -1)
	checkQrCodeUrl = strings.Replace(checkQrCodeUrl, "NEXT", odtuc.dconf.Aweme.OpenDouyin.CallBackUrl, -1)
	checkQrCodeUrl = strings.Replace(checkQrCodeUrl, "SCOPE", strings.Join(scope, ","), -1)
	checkQrCodeUrl = strings.Replace(checkQrCodeUrl, "STATE", state, -1)
	checkQrCodeUrl = strings.Replace(checkQrCodeUrl, "TOKEN", token, -1)
	checkQrCodeUrl = strings.Replace(checkQrCodeUrl, "TIMESTAMP", timestamp, -1)

	response, err := http.Get(checkQrCodeUrl)
	defer response.Body.Close()

	if err != nil {
		return nil, DouyinOpenStateGetError
	}

	rbody, _ := ioutil.ReadAll(response.Body)

	var qrCodeStatusBody *domain.QrCodeOpenDouyinStatus

	if err := json.Unmarshal(rbody, &qrCodeStatusBody); err != nil {
		return nil, DouyinOpenStateGetError
	}

	if qrCodeStatusBody.Message != "success" {
		return nil, DouyinOpenStateGetError
	}

	return &domain.QrCodeStatus{
		Code:   qrCodeStatusBody.Data.Code,
		Status: qrCodeStatusBody.Data.Status,
	}, nil
}

func (odtuc *OpenDouyinTokenUsecase) CreateOpenDouyinTokens(ctx context.Context, state, code string) error {
	sstateParam, err := odtuc.repo.GetCacheString(ctx, "douyin:openDouyin:auth:"+state)

	if err != nil {
		return DouyinOpenDouyinTokenCreateError
	}

	var stateParam domain.StateParam
	err = json.Unmarshal([]byte(sstateParam), &stateParam)

	if err != nil {
		return DouyinOpenDouyinTokenCreateError
	}

	userId, err := strconv.ParseUint(stateParam.UserId, 10, 64)

	if err != nil {
		return DouyinOpenDouyinTokenCreateError
	}

	if stateParam.ClientKey != odtuc.dconf.Aweme.OpenDouyin.ClientKey {
		return DouyinOpenDouyinTokenConfigNotFound
	}

	weixinUser, err := odtuc.wurepo.GetById(ctx, userId)

	if err != nil {
		return DouyinWeixinUserNotFound
	}

	accessToken, err := oauth2.AccessToken(stateParam.ClientKey, odtuc.dconf.Aweme.OpenDouyin.ClientSecret, code)

	if err != nil {
		inOpenDouyinApiLog := domain.NewOpenDouyinApiLog(ctx, stateParam.ClientKey, "", "", err.Error())
		inOpenDouyinApiLog.SetCreateTime(ctx)

		odtuc.odalrepo.Save(ctx, inOpenDouyinApiLog)

		return DouyinOpenDouyinTokenCreateError
	}

	if _, err := odtuc.wuodrepo.Get(ctx, stateParam.ClientKey, accessToken.Data.OpenId); err == nil {
		return DouyinOpenDouyinTokenExist
	}

	userInfo, err := oauth2.GetUserInfo(accessToken.Data.AccessToken, accessToken.Data.OpenId)

	if err != nil {
		inOpenDouyinApiLog := domain.NewOpenDouyinApiLog(ctx, stateParam.ClientKey, "", "", err.Error())
		inOpenDouyinApiLog.SetCreateTime(ctx)

		odtuc.odalrepo.Save(ctx, inOpenDouyinApiLog)

		return DouyinOpenDouyinTokenCreateError
	}

	buyinId, err := kol.GetBuyinId(accessToken.Data.AccessToken, accessToken.Data.OpenId)

	if err != nil {
		inOpenDouyinApiLog := domain.NewOpenDouyinApiLog(ctx, stateParam.ClientKey, "", "", err.Error())
		inOpenDouyinApiLog.SetCreateTime(ctx)

		odtuc.odalrepo.Save(ctx, inOpenDouyinApiLog)

		return DouyinOpenDouyinTokenCreateError
	}

	err = odtuc.tm.InTx(ctx, func(ctx context.Context) error {
		var inOpenDouyinToken *domain.OpenDouyinToken

		if inOpenDouyinToken, err = odtuc.repo.GetByClientKeyAndOpenId(ctx, stateParam.ClientKey, accessToken.Data.OpenId); err != nil {
			inOpenDouyinToken = domain.NewOpenDouyinToken(ctx, stateParam.ClientKey, accessToken.Data.OpenId, accessToken.Data.AccessToken, accessToken.Data.RefreshToken, accessToken.Data.ExpiresIn, accessToken.Data.RefreshExpiresIn)
			inOpenDouyinToken.SetCreateTime(ctx)
			inOpenDouyinToken.SetUpdateTime(ctx)

			if _, err := odtuc.repo.Save(ctx, inOpenDouyinToken); err != nil {
				return err
			}
		} else {
			inOpenDouyinToken.SetAccessToken(ctx, accessToken.Data.AccessToken)
			inOpenDouyinToken.SetExpiresIn(ctx, accessToken.Data.ExpiresIn)
			inOpenDouyinToken.SetRefreshToken(ctx, accessToken.Data.RefreshToken)
			inOpenDouyinToken.SetRefreshExpiresIn(ctx, accessToken.Data.RefreshExpiresIn)
			inOpenDouyinToken.SetUpdateTime(ctx)

			if _, err := odtuc.repo.Update(ctx, inOpenDouyinToken); err != nil {
				return err
			}
		}

		var inOpenDouyinUserInfo *domain.OpenDouyinUserInfo

		if inOpenDouyinUserInfo, err = odtuc.oduirepo.GetByClientKeyAndOpenId(ctx, stateParam.ClientKey, userInfo.Data.OpenId); err != nil {
			inOpenDouyinUserInfo = domain.NewOpenDouyinUserInfo(ctx, stateParam.ClientKey, userInfo.Data.OpenId, userInfo.Data.UnionId, "", buyinId.Data.BuyinId, userInfo.Data.Nickname, userInfo.Data.Avatar, userInfo.Data.AvatarLarger, "", userInfo.Data.Country, userInfo.Data.Province, userInfo.Data.City, userInfo.Data.District, userInfo.Data.EAccountRole, "", userInfo.Data.Gender)
			inOpenDouyinUserInfo.SetCreateTime(ctx)
			inOpenDouyinUserInfo.SetUpdateTime(ctx)

			if _, err := odtuc.oduirepo.Save(ctx, inOpenDouyinUserInfo); err != nil {
				return err
			}
		} else {
			inOpenDouyinUserInfo.SetUnionId(ctx, userInfo.Data.UnionId)
			inOpenDouyinUserInfo.SetBuyinId(ctx, buyinId.Data.BuyinId)
			inOpenDouyinUserInfo.SetNickname(ctx, userInfo.Data.Nickname)
			inOpenDouyinUserInfo.SetAvatar(ctx, userInfo.Data.Avatar)
			inOpenDouyinUserInfo.SetAvatarLarger(ctx, userInfo.Data.AvatarLarger)
			inOpenDouyinUserInfo.SetPhone(ctx, "")
			inOpenDouyinUserInfo.SetCountry(ctx, userInfo.Data.Country)
			inOpenDouyinUserInfo.SetProvince(ctx, userInfo.Data.Province)
			inOpenDouyinUserInfo.SetCity(ctx, userInfo.Data.City)
			inOpenDouyinUserInfo.SetDistrict(ctx, userInfo.Data.District)
			inOpenDouyinUserInfo.SetEAccountRole(ctx, userInfo.Data.EAccountRole)
			inOpenDouyinUserInfo.SetGender(ctx, userInfo.Data.Gender)
			inOpenDouyinUserInfo.SetUpdateTime(ctx)

			if _, err := odtuc.oduirepo.Update(ctx, inOpenDouyinUserInfo); err != nil {
				return err
			}
		}

		if _, err := odtuc.wuodrepo.Update(ctx, weixinUser.Data.UserId, stateParam.ClientKey, userInfo.Data.OpenId, userInfo.Data.Nickname, userInfo.Data.Avatar, userInfo.Data.AvatarLarger); err != nil {
			return err
		}

		inOpenDouyinUserInfoCreateLog := domain.NewOpenDouyinUserInfoCreateLog(ctx, stateParam.ClientKey, userInfo.Data.OpenId)
		inOpenDouyinUserInfoCreateLog.SetCreateTime(ctx)
		inOpenDouyinUserInfoCreateLog.SetUpdateTime(ctx)

		odtuc.oduiclrepo.Save(ctx, inOpenDouyinUserInfoCreateLog)

		return nil
	})

	if err != nil {
		return DouyinOpenDouyinTokenCreateError
	}

	odtuc.repo.DeleteCache(ctx, "douyin:openDouyin:auth:"+state)

	return nil
}

func (odtuc *OpenDouyinTokenUsecase) UpdateCooperativeCodeDouyinTokens(ctx context.Context, clientKey, openId, cooperativeCode string) error {
	if clientKey != odtuc.dconf.Aweme.OpenDouyin.ClientKey {
		return DouyinOpenDouyinTokenConfigNotFound
	}

	if err := odtuc.oduirepo.UpdateCooperativeCodes(ctx, clientKey, openId, cooperativeCode); err != nil {
		return DouyinOpenDouyinUserInfoUpdateError
	}

	return nil
}

func (odtuc *OpenDouyinTokenUsecase) RefreshOpenDouyinTokens(ctx context.Context) error {
	list, err := odtuc.repo.List(ctx)

	if err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "refreshToken", fmt.Sprintf("[Open Douyin RefreshTokenError] Description=%s", "获取达人账户token表失败"))
		inTaskLog.SetCreateTime(ctx)

		odtuc.tlrepo.Save(ctx, inTaskLog)

		return err
	}

	for _, inAwemeToken := range list {
		if time.Now().After(inAwemeToken.UpdateTime.Add((time.Duration(inAwemeToken.ExpiresIn) * time.Second) - 60*60*time.Second)) {
			if accessToken, err := oauth2.RefreshToken(inAwemeToken.ClientKey, inAwemeToken.RefreshToken); err == nil {
				inAwemeToken.SetAccessToken(ctx, accessToken.Data.AccessToken)
				inAwemeToken.SetExpiresIn(ctx, accessToken.Data.ExpiresIn)
				inAwemeToken.SetRefreshToken(ctx, accessToken.Data.RefreshToken)
				inAwemeToken.SetRefreshExpiresIn(ctx, accessToken.Data.RefreshExpiresIn)
				inAwemeToken.SetUpdateTime(ctx)

				if _, err := odtuc.repo.Update(ctx, inAwemeToken); err != nil {
					inTaskLog := domain.NewTaskLog(ctx, "refreshToken", fmt.Sprintf("[Open Douyin RefreshTokenError] ClientKey=%s, OpenId=%s, AccessToken=%s, RefreshToken=%s, ExpiresIn=%d, RefreshTokenExpiresIn=%d, Description=%s", inAwemeToken.ClientKey, inAwemeToken.OpenId, inAwemeToken.AccessToken, inAwemeToken.RefreshToken, inAwemeToken.ExpiresIn, inAwemeToken.RefreshExpiresIn, "获取新的token插入数据库失败"))
					inTaskLog.SetCreateTime(ctx)

					odtuc.tlrepo.Save(ctx, inTaskLog)
				}
			} else {
				inTaskLog := domain.NewTaskLog(ctx, "refreshToken", fmt.Sprintf("[Open Douyin RefreshTokenError] ClientKey=%d, OpenId=%d, Description=%s", inAwemeToken.ClientKey, inAwemeToken.OpenId, "调用抖音开放平台刷新token接口失败"))
				inTaskLog.SetCreateTime(ctx)

				odtuc.tlrepo.Save(ctx, inTaskLog)
			}
		}
	}

	return nil
}

func (odtuc *OpenDouyinTokenUsecase) RenewRefreshTokensOpenDouyinTokens(ctx context.Context) error {
	list, err := odtuc.repo.List(ctx)

	if err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "renewRefreshToken", fmt.Sprintf("[Open Douyin RenewRefreshTokenError] Description=%s", "获取达人账户token表失败"))
		inTaskLog.SetCreateTime(ctx)

		odtuc.tlrepo.Save(ctx, inTaskLog)

		return err
	}

	for _, inAwemeToken := range list {
		if time.Now().After(inAwemeToken.UpdateTime.Add((time.Duration(inAwemeToken.RefreshExpiresIn) * time.Second) - 60*60*time.Second)) {
			if refreshToken, err := oauth2.RenewRefreshToken(inAwemeToken.ClientKey, inAwemeToken.RefreshToken); err == nil {
				inAwemeToken.SetRefreshToken(ctx, refreshToken.Data.RefreshToken)
				inAwemeToken.SetRefreshExpiresIn(ctx, refreshToken.Data.ExpiresIn)
				inAwemeToken.SetUpdateTime(ctx)

				if _, err := odtuc.repo.Update(ctx, inAwemeToken); err != nil {
					inTaskLog := domain.NewTaskLog(ctx, "renewRefreshToken", fmt.Sprintf("[Open Douyin RenewRefreshTokenError] ClientKey=%s, OpenId=%s, AccessToken=%s, RefreshToken=%s, ExpiresIn=%d, RefreshTokenExpiresIn=%d, Description=%s", inAwemeToken.ClientKey, inAwemeToken.OpenId, inAwemeToken.AccessToken, inAwemeToken.RefreshToken, inAwemeToken.ExpiresIn, inAwemeToken.RefreshExpiresIn, "获取新的refresh token插入数据库失败"))
					inTaskLog.SetCreateTime(ctx)

					odtuc.tlrepo.Save(ctx, inTaskLog)
				}
			} else {
				inTaskLog := domain.NewTaskLog(ctx, "renewRefreshToken", fmt.Sprintf("[Open Douyin RenewRefreshTokenError] ClientKey=%d, OpenId=%d, Description=%s", inAwemeToken.ClientKey, inAwemeToken.OpenId, "调用抖音开放平台刷新refresh token接口失败"))
				inTaskLog.SetCreateTime(ctx)

				odtuc.tlrepo.Save(ctx, inTaskLog)
			}
		}
	}

	return nil
}

func (odtuc *OpenDouyinTokenUsecase) SyncUserFansOpenDouyinTokens(ctx context.Context) error {
	list, err := odtuc.repo.List(ctx)

	if err != nil {
		inTaskLog := domain.NewTaskLog(ctx, "syncUserFans", fmt.Sprintf("[Open Douyin SyncUserFansError] Description=%s", "获取达人账户token表失败"))
		inTaskLog.SetCreateTime(ctx)

		odtuc.tlrepo.Save(ctx, inTaskLog)

		return err
	}

	for _, l := range list {
		if userFans, err := user.ListUserFans(l.AccessToken, l.OpenId); err == nil {
			date, _ := tool.StringToTime("2006-01-02", "2006-01-02")
			var fans uint64

			for _, userFan := range userFans.Data.ResultList {
				if day, err := tool.StringToTime("2006-01-02", userFan.Date); err == nil {
					if day.After(date) {
						date = day
						fans = uint64(userFan.TotalFans)
					}
				}
			}

			if inOpenDouyinUserInfo, err := odtuc.oduirepo.GetByClientKeyAndOpenId(ctx, l.ClientKey, l.OpenId); err == nil {
				inOpenDouyinUserInfo.SetFans(ctx, fans)
				inOpenDouyinUserInfo.SetUpdateTime(ctx)

				odtuc.oduirepo.Update(ctx, inOpenDouyinUserInfo)

				odtuc.wuodrepo.UpdateUserFans(ctx, l.ClientKey, l.OpenId, fans)
			}
		}
	}

	return nil
}

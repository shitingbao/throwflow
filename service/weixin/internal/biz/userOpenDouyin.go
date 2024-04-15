package biz

import (
	"context"
	"encoding/json"
	"weixin/internal/conf"
	"weixin/internal/domain"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	WeixinUserOpenDouyinNotFound    = errors.NotFound("WEIXIN_USER_OPEN_DOUYIN_NOT_FOUND", "微信用户关联抖音用户不存在")
	WeixinUserOpenDouyinListError   = errors.InternalServer("WEIXIN_USER_OPEN_DOUYIN_LIST_ERROR", "微信用户关联抖音用户列表获取失败")
	WeixinUserOpenDouyinCreateError = errors.InternalServer("WEIXIN_USER_OPEN_DOUYIN_CREATE_ERROR", "微信用户关联抖音用户创建失败")
	WeixinUserOpenDouyinUpdateError = errors.InternalServer("WEIXIN_USER_OPEN_DOUYIN_UPDATE_ERROR", "微信用户关联抖音用户更新失败")
	WeixinUserOpenDouyinDeleteError = errors.InternalServer("WEIXIN_USER_OPEN_DOUYIN_DELETE_ERROR", "微信用户关联抖音用户删除失败")
)

type UserOpenDouyinRepo interface {
	Get(context.Context, uint64, string, string) (*domain.UserOpenDouyin, error)
	GetById(context.Context, uint64, uint64) (*domain.UserOpenDouyin, error)
	GetByClientKeyAndOpenId(context.Context, string, string) (*domain.UserOpenDouyin, error)
	List(context.Context, int, int, uint64, string) ([]*domain.UserOpenDouyin, error)
	ListByClientKeyAndOpenId(context.Context, int, int, []*domain.UserOpenDouyin, string) ([]*domain.UserOpenDouyin, error)
	Count(context.Context, uint64, string) (int64, error)
	CountByClientKeyAndOpenId(context.Context, []*domain.UserOpenDouyin, string) (int64, error)
	Save(context.Context, *domain.UserOpenDouyin) (*domain.UserOpenDouyin, error)
	Update(context.Context, *domain.UserOpenDouyin) (*domain.UserOpenDouyin, error)
	UpdateUserInfos(context.Context, uint64, string, string, string, string, string, string, string) error
	UpdateCooperativeCodes(context.Context, string, string, string) error
	UpdateFans(context.Context, string, string, uint64) error
	Delete(context.Context, *domain.UserOpenDouyin) error
	DeleteByUserId(context.Context, uint64, string, string) error
}

type UserOpenDouyinUsecase struct {
	repo     UserOpenDouyinRepo
	oduirepo OpenDouyinUserInfoRepo
	aawarepo AwemesAdvertiserWeixinAuthRepo
	darepo   DjAwemeRepo
	urepo    UserRepo
	tm       Transaction
	conf     *conf.Data
	log      *log.Helper
}

func NewUserOpenDouyinUsecase(repo UserOpenDouyinRepo, oduirepo OpenDouyinUserInfoRepo, aawarepo AwemesAdvertiserWeixinAuthRepo, darepo DjAwemeRepo, urepo UserRepo, tm Transaction, conf *conf.Data, logger log.Logger) *UserOpenDouyinUsecase {
	return &UserOpenDouyinUsecase{repo: repo, oduirepo: oduirepo, aawarepo: aawarepo, darepo: darepo, urepo: urepo, tm: tm, conf: conf, log: log.NewHelper(logger)}
}

func (uoduc *UserOpenDouyinUsecase) GetOpenDouyinUsers(ctx context.Context, clientKey, openId string) (*domain.UserOpenDouyin, error) {
	openDouyinUser, err := uoduc.repo.GetByClientKeyAndOpenId(ctx, clientKey, openId)

	if err != nil {
		return nil, WeixinUserOpenDouyinNotFound
	}

	openDouyinUser.SetFansShow(ctx)

	return &domain.UserOpenDouyin{
		Id:              openDouyinUser.Id,
		UserId:          openDouyinUser.UserId,
		ClientKey:       openDouyinUser.ClientKey,
		OpenId:          openDouyinUser.OpenId,
		AwemeId:         openDouyinUser.AwemeId,
		AccountId:       openDouyinUser.AccountId,
		Nickname:        openDouyinUser.Nickname,
		Avatar:          openDouyinUser.Avatar,
		AvatarLarger:    openDouyinUser.AvatarLarger,
		CooperativeCode: openDouyinUser.CooperativeCode,
		Fans:            openDouyinUser.Fans,
		FansShow:        openDouyinUser.FansShow,
		Area:            openDouyinUser.Area,
		CreateTime:      openDouyinUser.CreateTime,
		UpdateTime:      openDouyinUser.UpdateTime,
	}, nil
}

func (uoduc *UserOpenDouyinUsecase) ListOpenDouyinUsers(ctx context.Context, pageNum, pageSize, userId uint64, keyword string) (*domain.UserOpenDouyinList, error) {
	openDouyinUsers, err := uoduc.repo.List(ctx, int(pageNum), int(pageSize), userId, keyword)

	if err != nil {
		return nil, WeixinUserOpenDouyinListError
	}

	total, err := uoduc.repo.Count(ctx, userId, keyword)

	if err != nil {
		return nil, WeixinUserOpenDouyinListError
	}

	list := make([]*domain.UserOpenDouyin, 0)

	for _, openDouyinUser := range openDouyinUsers {
		if len(openDouyinUser.AccountId) > 0 {
			if djAweme, err := uoduc.darepo.Get(ctx, openDouyinUser.AccountId, ""); err == nil {
				openDouyinUser.SetLevel(ctx, djAweme.Level)
			}
		}

		openDouyinUser.SetFansShow(ctx)

		list = append(list, &domain.UserOpenDouyin{
			Id:              openDouyinUser.Id,
			UserId:          openDouyinUser.UserId,
			ClientKey:       openDouyinUser.ClientKey,
			OpenId:          openDouyinUser.OpenId,
			AwemeId:         openDouyinUser.AwemeId,
			AccountId:       openDouyinUser.AccountId,
			Nickname:        openDouyinUser.Nickname,
			Avatar:          openDouyinUser.Avatar,
			AvatarLarger:    openDouyinUser.AvatarLarger,
			CooperativeCode: openDouyinUser.CooperativeCode,
			Fans:            openDouyinUser.Fans,
			FansShow:        openDouyinUser.FansShow,
			Area:            openDouyinUser.Area,
			Level:           openDouyinUser.Level,
			CreateTime:      openDouyinUser.CreateTime,
			UpdateTime:      openDouyinUser.UpdateTime,
		})
	}

	return &domain.UserOpenDouyinList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (uoduc *UserOpenDouyinUsecase) ListByClientKeyAndOpenIds(ctx context.Context, pageNum, pageSize uint64, clientKeyAndOpenIds, keyword string) (*domain.UserOpenDouyinList, error) {
	userOpenDouyins := make([]*domain.UserOpenDouyin, 0)

	if len(clientKeyAndOpenIds) > 0 {
		json.Unmarshal([]byte(clientKeyAndOpenIds), &userOpenDouyins)
	}

	openDouyinUsers, err := uoduc.repo.ListByClientKeyAndOpenId(ctx, int(pageNum), int(pageSize), userOpenDouyins, keyword)

	if err != nil {
		return nil, WeixinUserOpenDouyinListError
	}

	total, err := uoduc.repo.CountByClientKeyAndOpenId(ctx, userOpenDouyins, keyword)

	if err != nil {
		return nil, WeixinUserOpenDouyinListError
	}

	list := make([]*domain.UserOpenDouyin, 0)

	for _, openDouyinUser := range openDouyinUsers {
		openDouyinUser.SetFansShow(ctx)

		list = append(list, &domain.UserOpenDouyin{
			Id:              openDouyinUser.Id,
			UserId:          openDouyinUser.UserId,
			ClientKey:       openDouyinUser.ClientKey,
			OpenId:          openDouyinUser.OpenId,
			AwemeId:         openDouyinUser.AwemeId,
			AccountId:       openDouyinUser.AccountId,
			Nickname:        openDouyinUser.Nickname,
			Avatar:          openDouyinUser.Avatar,
			AvatarLarger:    openDouyinUser.AvatarLarger,
			CooperativeCode: openDouyinUser.CooperativeCode,
			Fans:            openDouyinUser.Fans,
			FansShow:        openDouyinUser.FansShow,
			Area:            openDouyinUser.Area,
			CreateTime:      openDouyinUser.CreateTime,
			UpdateTime:      openDouyinUser.UpdateTime,
		})
	}

	return &domain.UserOpenDouyinList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (uoduc *UserOpenDouyinUsecase) UpdateOpenDouyinUsers(ctx context.Context, userId, awemeId uint64, clientKey, openId, accountId, nickname, avatar, avatarLarger, area string) error {
	user, err := uoduc.urepo.Get(ctx, userId)

	if err != nil {
		return WeixinUserNotFound
	}

	if inUserOpenDouyin, err := uoduc.repo.Get(ctx, user.Id, clientKey, openId); err != nil {
		inUserOpenDouyin = domain.NewUserOpenDouyin(ctx, userId, clientKey, openId, accountId, nickname, avatar, avatarLarger, "")
		inUserOpenDouyin.SetAwemeId(ctx, awemeId)
		inUserOpenDouyin.SetArea(ctx, area)
		inUserOpenDouyin.SetCreateTime(ctx)
		inUserOpenDouyin.SetUpdateTime(ctx)

		if _, err := uoduc.repo.Save(ctx, inUserOpenDouyin); err != nil {
			return WeixinUserOpenDouyinCreateError
		}
	} else {
		inUserOpenDouyin.SetAwemeId(ctx, awemeId)
		inUserOpenDouyin.SetAccountId(ctx, accountId)
		inUserOpenDouyin.SetNickname(ctx, nickname)
		inUserOpenDouyin.SetAvatar(ctx, avatar)
		inUserOpenDouyin.SetAvatarLarger(ctx, avatarLarger)
		inUserOpenDouyin.SetArea(ctx, area)
		inUserOpenDouyin.SetUpdateTime(ctx)

		if _, err := uoduc.repo.Update(ctx, inUserOpenDouyin); err != nil {
			return WeixinUserOpenDouyinUpdateError
		}
	}

	if err := uoduc.repo.DeleteByUserId(ctx, user.Id, clientKey, openId); err != nil {
		return WeixinUserOpenDouyinUpdateError
	}

	return nil
}

func (uoduc *UserOpenDouyinUsecase) UpdateUserInfoOpenDouyinUsers(ctx context.Context, awemeId uint64, clientKey, openId, accountId, nickname, avatar, avatarLarger, area string) error {
	if err := uoduc.repo.UpdateUserInfos(ctx, awemeId, clientKey, openId, accountId, nickname, avatar, avatarLarger, area); err != nil {
		return WeixinUserOpenDouyinUpdateError
	}

	return nil
}

func (uoduc *UserOpenDouyinUsecase) UpdateCooperativeCodeOpenDouyinUsers(ctx context.Context, userId, openDouyinUserId uint64, cooperativeCode string) ([]*domain.UserOpenDouyin, error) {
	openDouyinUser, err := uoduc.repo.GetById(ctx, userId, openDouyinUserId)

	if err != nil {
		return nil, WeixinUserOpenDouyinNotFound
	}

	if cooperativeCode != openDouyinUser.CooperativeCode {
		err = uoduc.tm.InTx(ctx, func(ctx context.Context) error {
			if err := uoduc.repo.UpdateCooperativeCodes(ctx, openDouyinUser.ClientKey, openDouyinUser.OpenId, cooperativeCode); err != nil {
				return err
			}

			if _, err := uoduc.oduirepo.UpdateCooperativeCodes(ctx, openDouyinUser.ClientKey, openDouyinUser.OpenId, cooperativeCode); err != nil {
				return err
			}

			inAwemesAdvertiserWeixinAuth := domain.NewAwemesAdvertiserWeixinAuth(ctx, 1, openDouyinUser.ClientKey, openDouyinUser.OpenId, cooperativeCode)

			if err := uoduc.aawarepo.Upsert(ctx, inAwemesAdvertiserWeixinAuth); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return nil, WeixinUserOpenDouyinUpdateError
		}
	}

	list, err := uoduc.repo.List(ctx, 0, 40, userId, "")

	if err != nil {
		return nil, WeixinUserOpenDouyinListError
	}

	return list, nil
}

func (uoduc *UserOpenDouyinUsecase) UpdateUserFansOpenDouyinUsers(ctx context.Context, clientKey, openId string, fans uint64) error {
	if err := uoduc.repo.UpdateFans(ctx, clientKey, openId, fans); err != nil {
		return WeixinUserOpenDouyinUpdateError
	}

	return nil
}

func (uoduc *UserOpenDouyinUsecase) DeleteOpenDouyinUsers(ctx context.Context, userId, openDouyinUserId uint64) error {
	inOpenDouyinUser, err := uoduc.repo.GetById(ctx, userId, openDouyinUserId)

	if err != nil {
		return WeixinUserOpenDouyinNotFound
	}

	if err := uoduc.repo.Delete(ctx, inOpenDouyinUser); err != nil {
		return WeixinUserOpenDouyinDeleteError
	}

	return nil
}

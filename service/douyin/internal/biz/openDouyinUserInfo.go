package biz

import (
	"context"
	"douyin/internal/conf"
	"douyin/internal/domain"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	DouyinOpenDouyinUserInfoNotFound    = errors.InternalServer("DOUYIN_OPEN_DOUYIN_USER_INFO_NOT_FOUND", "抖音开放平台用户信息不存在")
	DouyinOpenDouyinUserInfoListError   = errors.InternalServer("DOUYIN_OPEN_DOUYIN_USER_INFO_LIST_ERROR", "抖音开放平台用户信息列表获取失败")
	DouyinOpenDouyinUserInfoUpdateError = errors.InternalServer("DOUYIN_OPEN_DOUYIN_USER_INFO_UPDATE_ERROR", "抖音开放平台用户信息更新失败")
)

type OpenDouyinUserInfoRepo interface {
	GetByClientKeyAndOpenId(ctx context.Context, clientKey, openId string) (*domain.OpenDouyinUserInfo, error)
	List(context.Context, int, int) ([]*domain.OpenDouyinUserInfo, error)
	ListByClientKeyAndOpenId(context.Context, []*domain.OpenDouyinUserInfo) ([]*domain.OpenDouyinUserInfo, error)
	Count(context.Context) (int64, error)
	Save(ctx context.Context, in *domain.OpenDouyinUserInfo) (*domain.OpenDouyinUserInfo, error)
	Update(ctx context.Context, in *domain.OpenDouyinUserInfo) (*domain.OpenDouyinUserInfo, error)
	UpdateCooperativeCodes(context.Context, string, string, string) error
}

type OpenDouyinUserInfoUsecase struct {
	repo     OpenDouyinUserInfoRepo
	odvrepo  OpenDouyinVideoRepo
	jsirepo  JinritemaiStoreInfoRepo
	joirepo  JinritemaiOrderInfoRepo
	wuodrepo WeixinUserOpenDouyinRepo
	tm       Transaction
	conf     *conf.Data
	log      *log.Helper
}

func NewOpenDouyinUserInfoUsecase(repo OpenDouyinUserInfoRepo, odvrepo OpenDouyinVideoRepo, jsirepo JinritemaiStoreInfoRepo, joirepo JinritemaiOrderInfoRepo, wuodrepo WeixinUserOpenDouyinRepo, tm Transaction, conf *conf.Data, logger log.Logger) *OpenDouyinUserInfoUsecase {
	return &OpenDouyinUserInfoUsecase{repo: repo, odvrepo: odvrepo, jsirepo: jsirepo, joirepo: joirepo, wuodrepo: wuodrepo, tm: tm, conf: conf, log: log.NewHelper(logger)}
}

func (oduiuc *OpenDouyinUserInfoUsecase) ListOpenDouyinUserInfos(ctx context.Context, pageNum, pageSize uint64) (*domain.OpenDouyinUserInfoList, error) {
	openDouyinUserInfos, err := oduiuc.repo.List(ctx, int(pageNum), int(pageSize))

	if err != nil {
		return nil, DouyinOpenDouyinUserInfoListError
	}

	list := make([]*domain.OpenDouyinUserInfo, 0)

	for _, openDouyinUserInfo := range openDouyinUserInfos {
		if openDouyinVideo, err := oduiuc.odvrepo.GetByClientKeyAndOpenId(ctx, openDouyinUserInfo.ClientKey, openDouyinUserInfo.OpenId, "2,4", 1); err == nil {
			openDouyinUserInfo.SetVideoId(ctx, openDouyinVideo.VideoId)
		}

		list = append(list, openDouyinUserInfo)
	}

	total, err := oduiuc.repo.Count(ctx)

	if err != nil {
		return nil, DouyinOpenDouyinUserInfoListError
	}

	return &domain.OpenDouyinUserInfoList{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (oduiuc *OpenDouyinUserInfoUsecase) ListOpenDouyinUserInfosByProductId(ctx context.Context, productId string) ([]*domain.OpenDouyinUserInfo, error) {
	list, err := oduiuc.jsirepo.ListByProductId(ctx, productId)

	if err != nil {
		return nil, DouyinOpenDouyinUserInfoListError
	}

	return list, nil
}

func (oduiuc *OpenDouyinUserInfoUsecase) UpdateOpenDouyinUserInfos(ctx context.Context, awemeId uint64, clientKey, openId, accountId, nickname, avatar, avatarLarger, area string) error {
	inOpenDouyinUserInfo, err := oduiuc.repo.GetByClientKeyAndOpenId(ctx, clientKey, openId)

	if err != nil {
		return DouyinOpenDouyinUserInfoNotFound
	}

	if inOpenDouyinUserInfo.AwemeId == awemeId &&
		inOpenDouyinUserInfo.AccountId == accountId &&
		inOpenDouyinUserInfo.Nickname == nickname &&
		inOpenDouyinUserInfo.Avatar == avatar &&
		inOpenDouyinUserInfo.AvatarLarger == avatarLarger &&
		inOpenDouyinUserInfo.Area == area {
		return nil
	}

	inOpenDouyinUserInfo.SetAwemeId(ctx, awemeId)
	inOpenDouyinUserInfo.SetAccountId(ctx, accountId)
	inOpenDouyinUserInfo.SetNickname(ctx, nickname)
	inOpenDouyinUserInfo.SetAvatar(ctx, avatar)
	inOpenDouyinUserInfo.SetAvatarLarger(ctx, avatarLarger)
	inOpenDouyinUserInfo.SetArea(ctx, area)
	inOpenDouyinUserInfo.SetUpdateTime(ctx)

	if _, err := oduiuc.repo.Update(ctx, inOpenDouyinUserInfo); err != nil {
		return DouyinOpenDouyinUserInfoUpdateError
	}

	if _, err := oduiuc.wuodrepo.UpdateUserInfos(ctx, inOpenDouyinUserInfo.AwemeId, inOpenDouyinUserInfo.ClientKey, inOpenDouyinUserInfo.OpenId, inOpenDouyinUserInfo.AccountId, inOpenDouyinUserInfo.Nickname, inOpenDouyinUserInfo.Avatar, inOpenDouyinUserInfo.AvatarLarger, inOpenDouyinUserInfo.Area); err != nil {
		return DouyinOpenDouyinUserInfoUpdateError
	}
	
	return nil
}

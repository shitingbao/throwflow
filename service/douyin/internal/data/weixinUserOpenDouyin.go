package data

import (
	"context"
	v1 "douyin/api/service/weixin/v1"
	"douyin/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type weixinUserOpenDouyinRepo struct {
	data *Data
	log  *log.Helper
}

func NewWeixinUserOpenDouyinRepo(data *Data, logger log.Logger) biz.WeixinUserOpenDouyinRepo {
	return &weixinUserOpenDouyinRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (wuodr *weixinUserOpenDouyinRepo) Get(ctx context.Context, clientKey, openId string) (*v1.GetOpenDouyinUsersReply, error) {
	openDouyinUser, err := wuodr.data.weixinuc.GetOpenDouyinUsers(ctx, &v1.GetOpenDouyinUsersRequest{
		ClientKey: clientKey,
		OpenId:    openId,
	})

	if err != nil {
		return nil, err
	}

	return openDouyinUser, err
}

func (wuodr *weixinUserOpenDouyinRepo) List(ctx context.Context, userId uint64) (*v1.ListOpenDouyinUsersReply, error) {
	list, err := wuodr.data.weixinuc.ListOpenDouyinUsers(ctx, &v1.ListOpenDouyinUsersRequest{
		PageSize: 40,
		UserId:   userId,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (wuodr *weixinUserOpenDouyinRepo) Update(ctx context.Context, userId uint64, clientKey, openId, nickname, avatar, avatarLarger string) (*v1.UpdateOpenDouyinUsersReply, error) {
	openDouyinUser, err := wuodr.data.weixinuc.UpdateOpenDouyinUsers(ctx, &v1.UpdateOpenDouyinUsersRequest{
		UserId:       userId,
		ClientKey:    clientKey,
		OpenId:       openId,
		Nickname:     nickname,
		Avatar:       avatar,
		AvatarLarger: avatarLarger,
	})

	if err != nil {
		return nil, err
	}

	return openDouyinUser, err
}

func (wuodr *weixinUserOpenDouyinRepo) UpdateUserInfos(ctx context.Context, awemeId uint64, clientKey, openId, accountId, nickname, avatar, avatarLarger, area string) (*v1.UpdateUserInfoOpenDouyinUsersReply, error) {
	openDouyinUser, err := wuodr.data.weixinuc.UpdateUserInfoOpenDouyinUsers(ctx, &v1.UpdateUserInfoOpenDouyinUsersRequest{
		ClientKey:    clientKey,
		OpenId:       openId,
		AwemeId:      awemeId,
		AccountId:    accountId,
		Nickname:     nickname,
		Avatar:       avatar,
		AvatarLarger: avatarLarger,
		Area:         area,
	})

	if err != nil {
		return nil, err
	}

	return openDouyinUser, err
}

func (wuodr *weixinUserOpenDouyinRepo) UpdateUserFans(ctx context.Context, clientKey, openId string, fans uint64) (*v1.UpdateUserFansOpenDouyinUsersReply, error) {
	openDouyinUser, err := wuodr.data.weixinuc.UpdateUserFansOpenDouyinUsers(ctx, &v1.UpdateUserFansOpenDouyinUsersRequest{
		ClientKey: clientKey,
		OpenId:    openId,
		Fans:      fans,
	})

	if err != nil {
		return nil, err
	}

	return openDouyinUser, err
}

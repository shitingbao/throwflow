package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/weixin/v1"
	"interface/internal/biz"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ur *userRepo) Get(ctx context.Context, token string) (*v1.GetUsersReply, error) {
	user, err := ur.data.weixinuc.GetUsers(ctx, &v1.GetUsersRequest{
		Token: token,
	})

	if err != nil {
		return nil, err
	}

	return user, err
}

func (ur *userRepo) GetById(ctx context.Context, userId uint64) (*v1.GetByIdUsersReply, error) {
	user, err := ur.data.weixinuc.GetByIdUsers(ctx, &v1.GetByIdUsersRequest{
		UserId: userId,
	})

	if err != nil {
		return nil, err
	}

	return user, err
}

func (ur *userRepo) GetFollow(ctx context.Context, organizationId, parentUserId uint64) (*v1.GetFollowUsersReply, error) {
	user, err := ur.data.weixinuc.GetFollowUsers(ctx, &v1.GetFollowUsersRequest{
		OrganizationId: organizationId,
		ParentUserId:   parentUserId,
	})

	if err != nil {
		return nil, err
	}

	return user, err
}

func (ur *userRepo) Create(ctx context.Context, organizationId uint64, loginCode, phoneCode string) (*v1.CreateUsersReply, error) {
	user, err := ur.data.weixinuc.CreateUsers(ctx, &v1.CreateUsersRequest{
		OrganizationId: organizationId,
		LoginCode:      loginCode,
		PhoneCode:      phoneCode,
	})

	if err != nil {
		return nil, err
	}

	return user, err
}

func (ur *userRepo) Update(ctx context.Context, userId uint64, nickName, avatar string) (*v1.UpdateUsersReply, error) {
	user, err := ur.data.weixinuc.UpdateUsers(ctx, &v1.UpdateUsersRequest{
		UserId:   userId,
		NickName: nickName,
		Avatar:   avatar,
	})

	if err != nil {
		return nil, err
	}

	return user, err
}

package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/weixin/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"
)

type UserRepo interface {
	Get(context.Context, string) (*v1.GetUsersReply, error)
	GetById(context.Context, uint64) (*v1.GetByIdUsersReply, error)
	GetFollow(context.Context, uint64, uint64) (*v1.GetFollowUsersReply, error)
	Create(context.Context, uint64, string, string) (*v1.CreateUsersReply, error)
	Update(context.Context, uint64, string, string) (*v1.UpdateUsersReply, error)
}

type UserUsecase struct {
	repo    UserRepo
	usrrepo UserScanRecordRepo
	conf    *conf.Data
	log     *log.Helper
}

func NewUserUsecase(repo UserRepo, usrrepo UserScanRecordRepo, conf *conf.Data, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, usrrepo: usrrepo, conf: conf, log: log.NewHelper(logger)}
}

func (uuc *UserUsecase) GetUsers(ctx context.Context, userId uint64) (*v1.GetByIdUsersReply, error) {
	user, err := uuc.repo.GetById(ctx, userId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_USER_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return user, nil
}

func (uuc *UserUsecase) GetFollows(ctx context.Context, organizationId, parentUserId uint64) (*v1.GetFollowUsersReply, error) {
	follow, err := uuc.repo.GetFollow(ctx, organizationId, parentUserId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_FOLLOW_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return follow, nil
}

func (uuc *UserUsecase) GetUser(ctx context.Context) (*v1.GetUsersReply, error) {
	token := ctx.Value("token")

	user, err := uuc.repo.Get(ctx, token.(string))

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LOGIN_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return user, nil
}

func (uuc *UserUsecase) CreateUsers(ctx context.Context, organizationId uint64, loginCode, phoneCode string) (*v1.CreateUsersReply, error) {
	user, err := uuc.repo.Create(ctx, organizationId, loginCode, phoneCode)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_CREATE_USER_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return user, nil
}

func (uuc *UserUsecase) UpdateUsers(ctx context.Context, userId uint64, nickName, avatar string) (*v1.UpdateUsersReply, error) {
	user, err := uuc.repo.Update(ctx, userId, nickName, avatar)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_UPDATE_USER_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return user, nil
}

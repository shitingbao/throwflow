package biz

import (
	"admin/internal/conf"
	"admin/internal/domain"
	"admin/internal/pkg/tool"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	AdminLoginError                   = errors.InternalServer("ADMIN_LOGIN_ERROR", "登录异常错误")
	AdminLoginUsernameOrPasswordError = errors.InternalServer("ADMIN_LOGIN_USERNAME_OR_PASSWORD_ERROR", "账号或密码错误")
	AdminLoginPermissionError         = errors.InternalServer("ADMIN_LOGIN_PERMISSION_ERROR", "权限不够")
	AdminLoginTokenError              = errors.InternalServer("ADMIN_LOGIN_TOKEN_ERROR", "token验证失败")

	AdminUserCreateError = errors.InternalServer("ADMIN_USER_CREATE_ERROR", "管理员创建失败")
	AdminUserUpdateError = errors.InternalServer("ADMIN_USER_UPDATE_ERROR", "管理员更新失败")
	AdminUserNotDelete   = errors.NotFound("ADMIN_USER_NOT_DELETE", "管理员不能删除自己")
	AdminUserDeleteError = errors.InternalServer("ADMIN_USER_DELETE_ERROR", "管理员删除失败")
	AdminUserNotFound    = errors.NotFound("ADMIN_USER_NOT_FOUND", "管理员不存在")
	AdminUserNotClaim    = errors.InternalServer("ADMIN_USER_NOT_CLAIM", "超级管理员无需认领线索")
)

type UserRepo interface {
	GetByUsername(context.Context, string) (*domain.User, error)
	GetById(context.Context, uint64) (*domain.User, error)
	List(context.Context, int) ([]*domain.User, error)
	Count(context.Context) (int64, error)
	Save(context.Context, *domain.User) (*domain.User, error)
	Update(context.Context, *domain.User) (*domain.User, error)
	Delete(context.Context, *domain.User) error

	SaveCacheHash(context.Context, string, map[string]string, time.Duration) error
	GetCacheHash(context.Context, string, string) (uint64, error)

	DeleteCache(context.Context, string) error
}

type UserUsecase struct {
	repo  UserRepo
	mrepo MenuRepo
	conf  *conf.Data
	log   *log.Helper
}

func NewUserUsecase(repo UserRepo, mrepo MenuRepo, conf *conf.Data, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, mrepo: mrepo, conf: conf, log: log.NewHelper(logger)}
}

func (uuc *UserUsecase) GetUserById(ctx context.Context, id uint64) (*domain.User, error) {
	user, err := uuc.getUserById(ctx, id)

	if err != nil {
		return nil, AdminUserNotFound
	}

	return user, nil
}

func (uuc *UserUsecase) GetUserByToken(ctx context.Context) (*domain.User, error) {
	token := ctx.Value("token")

	id, err := uuc.repo.GetCacheHash(ctx, "admin:admin:"+token.(string), "id")

	if err != nil {
		return nil, AdminLoginError
	}

	user, err := uuc.getUserById(ctx, id)

	if err != nil {
		return nil, AdminLoginError
	}

	return user, nil
}

func (uuc *UserUsecase) ListUsers(ctx context.Context, pageNum uint64) (*domain.UserList, error) {
	list, err := uuc.repo.List(ctx, int(pageNum))

	if err != nil {
		return nil, AdminDataError
	}

	total, err := uuc.repo.Count(ctx)

	if err != nil {
		return nil, AdminDataError
	}

	return &domain.UserList{
		PageNum:  pageNum,
		PageSize: uint64(uuc.conf.Database.PageSize),
		Total:    uint64(total),
		List:     list,
	}, nil
}

func (uuc *UserUsecase) CreateUsers(ctx context.Context, username, nickname, email, password string, roleId uint64, status uint8) (*domain.User, error) {
	inUser := domain.NewUser(ctx, username, nickname, email, password, roleId, status)
	inUser.SetCreateTime(ctx)
	inUser.SetUpdateTime(ctx)
	inUser.SetLastLoginIp(ctx)
	inUser.SetLastLoginTime(ctx)
	inUser.SetPassword(ctx, password)

	user, err := uuc.repo.Save(ctx, inUser)

	if err != nil {
		return nil, AdminUserCreateError
	}

	return user, nil
}

func (uuc *UserUsecase) UpdateUsers(ctx context.Context, id uint64, username, nickname, email, password string, roleId uint64, status uint8) (*domain.User, error) {
	inUser, err := uuc.getUserById(ctx, id)

	if err != nil {
		return nil, AdminUserNotFound
	}

	if len := utf8.RuneCountInString(password); len > 0 {
		inUser.SetPassword(ctx, password)
	}

	inUser.SetUsername(ctx, username)
	inUser.SetNickname(ctx, nickname)
	inUser.SetEmail(ctx, email)
	inUser.SetRoleId(ctx, roleId)
	inUser.SetStatus(ctx, status)
	inUser.SetUpdateTime(ctx)

	user, err := uuc.repo.Update(ctx, inUser)

	if err != nil {
		return nil, AdminUserUpdateError
	}

	return user, nil
}

func (uuc *UserUsecase) UpdateStatusUsers(ctx context.Context, id uint64, status uint8) (*domain.User, error) {
	inUser, err := uuc.getUserById(ctx, id)

	if err != nil {
		return nil, AdminUserNotFound
	}

	inUser.SetStatus(ctx, status)
	inUser.SetUpdateTime(ctx)

	user, err := uuc.repo.Update(ctx, inUser)

	if err != nil {
		return nil, AdminUserUpdateError
	}

	return user, nil
}

func (uuc *UserUsecase) DeleteUsers(ctx context.Context, loginUserId, id uint64) error {
	inUser, err := uuc.getUserById(ctx, id)

	if err != nil {
		return AdminUserNotFound
	}

	if loginUserId == inUser.Id {
		return AdminUserNotDelete
	}

	if err := uuc.repo.Delete(ctx, inUser); err != nil {
		return AdminUserDeleteError
	}

	return nil
}

func (uuc *UserUsecase) Login(ctx context.Context, username, password string) (*domain.Login, error) {
	user, err := uuc.getUserByUsername(ctx, username)

	if err != nil {
		return nil, AdminLoginUsernameOrPasswordError
	}

	if ok := user.VerifyPassword(ctx, password); !ok {
		return nil, AdminLoginUsernameOrPasswordError
	}

	if ok := user.VerifyStatus(ctx); !ok {
		return nil, AdminLoginError
	}

	if ok := user.Role.VerifyStatus(ctx); !ok {
		return nil, AdminLoginError
	}

	user.SetLastLoginIp(ctx)
	user.SetLastLoginTime(ctx)
	user.SetUpdateTime(ctx)

	if _, err := uuc.repo.Update(ctx, user); err != nil {
		return nil, AdminLoginError
	}

	cacheData := make(map[string]string)
	cacheData["id"] = strconv.FormatUint(user.Id, 10)

	token := tool.GetToken()

	if err := uuc.repo.SaveCacheHash(ctx, "admin:admin:"+token, cacheData, uuc.conf.Redis.LoginTokenTimeout.AsDuration()); err != nil {
		return nil, AdminLoginError
	}

	return &domain.Login{
		Username:      user.Username,
		Nickname:      user.Nickname,
		Email:         user.Email,
		LastLoginTime: user.LastLoginTime,
		LastLoginIp:   user.LastLoginIp,
		Token:         token,
		RoleName:      user.Role.RoleName,
		Menus:         user.Menus,
	}, nil
}

func (uuc *UserUsecase) Logout(ctx context.Context) error {
	token := ctx.Value("token")

	if err := uuc.repo.DeleteCache(ctx, "admin:admin:"+token.(string)); err != nil {
		return err
	}

	return nil
}

func (uuc *UserUsecase) getUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := uuc.repo.GetByUsername(ctx, username)

	if err != nil {
		return nil, err
	}

	user.Menus, _ = uuc.mrepo.ListByIds(ctx, strings.Split(user.Role.MenuIds, ","))

	return user, nil
}

func (uuc *UserUsecase) getUserById(ctx context.Context, id uint64) (*domain.User, error) {
	user, err := uuc.repo.GetById(ctx, id)

	if err != nil {
		return nil, err
	}

	user.Menus, _ = uuc.mrepo.ListByIds(ctx, strings.Split(user.Role.MenuIds, ","))

	return user, nil
}

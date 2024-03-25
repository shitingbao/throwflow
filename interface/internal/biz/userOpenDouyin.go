package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	dv1 "interface/api/service/douyin/v1"
	v1 "interface/api/service/weixin/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"
)

type UserOpenDouyinRepo interface {
	GetUrl(context.Context, uint64) (*dv1.GetUrlOpenDouyinTokensReply, error)
	GetTicket(context.Context) (*dv1.GetTicketOpenDouyinTokensReply, error)
	GetQrCode(context.Context, uint64) (*dv1.GetQrCodeOpenDouyinTokensReply, error)
	GetStatusQrCode(context.Context, string, string, string) (*dv1.GetStatusOpenDouyinTokensReply, error)
	List(context.Context, uint64, uint64, uint64, string) (*v1.ListOpenDouyinUsersReply, error)
	Save(context.Context, string, string) (*dv1.CreateOpenDouyinTokensReply, error)
	UpdateCooperativeCode(context.Context, uint64, uint64, string) (*v1.UpdateCooperativeCodeOpenDouyinUsersReply, error)
	Delete(context.Context, uint64, uint64) (*v1.DeleteOpenDouyinUsersReply, error)
}

type UserOpenDouyinUsecase struct {
	repo UserOpenDouyinRepo
	conf *conf.Data
	log  *log.Helper
}

func NewUserOpenDouyinUsecase(repo UserOpenDouyinRepo, conf *conf.Data, logger log.Logger) *UserOpenDouyinUsecase {
	return &UserOpenDouyinUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (uoduc *UserOpenDouyinUsecase) GetUrlUserOpenDouyins(ctx context.Context, userId uint64) (*dv1.GetUrlOpenDouyinTokensReply, error) {
	redirectUrl, err := uoduc.repo.GetUrl(ctx, userId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_URL_USER_OPEN_DOUYIN_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return redirectUrl, nil
}

func (uoduc *UserOpenDouyinUsecase) GetTicketUserOpenDouyins(ctx context.Context) (*dv1.GetTicketOpenDouyinTokensReply, error) {
	ticket, err := uoduc.repo.GetTicket(ctx)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_TICKET_USER_OPEN_DOUYIN_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return ticket, nil
}

func (uoduc *UserOpenDouyinUsecase) GetQrCodeUserOpenDouyins(ctx context.Context, userId uint64) (*dv1.GetQrCodeOpenDouyinTokensReply, error) {
	qrCode, err := uoduc.repo.GetQrCode(ctx, userId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_QR_CODE_USER_OPEN_DOUYIN_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return qrCode, nil
}

func (uoduc *UserOpenDouyinUsecase) GetStatusQrCodeUserOpenDouyins(ctx context.Context, token, state, timestamp string) (*dv1.GetStatusOpenDouyinTokensReply, error) {
	qrCodeStatus, err := uoduc.repo.GetStatusQrCode(ctx, token, state, timestamp)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_STATUS_QR_CODE_USER_OPEN_DOUYIN_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return qrCodeStatus, nil
}

func (uoduc *UserOpenDouyinUsecase) ListUserOpenDouyins(ctx context.Context, pageNum, pageSize, userId uint64, keyword string) (*v1.ListOpenDouyinUsersReply, error) {
	if pageSize == 0 {
		pageSize = uint64(uoduc.conf.Database.PageSize)
	}

	list, err := uoduc.repo.List(ctx, userId, pageNum, pageSize, keyword)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_USER_OPEN_DOUYIN_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (uoduc *UserOpenDouyinUsecase) CreateUserOpenDouyins(ctx context.Context, state, code string) (*dv1.CreateOpenDouyinTokensReply, error) {
	userOpenDouyin, err := uoduc.repo.Save(ctx, state, code)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_CREATE_USER_OPEN_DOUYIN_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userOpenDouyin, nil
}

func (uoduc *UserOpenDouyinUsecase) UpdateCooperativeCodeUserOpenDouyins(ctx context.Context, userId, openDouyinUserId uint64, cooperativeCode string) (*v1.UpdateCooperativeCodeOpenDouyinUsersReply, error) {
	list, err := uoduc.repo.UpdateCooperativeCode(ctx, userId, openDouyinUserId, cooperativeCode)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_UPDATE_COOPERATIVE_CODE_USER_OPEN_DOUYIN_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (uoduc *UserOpenDouyinUsecase) DeleteUserOpenDouyins(ctx context.Context, userId, openDouyinUserId uint64) (*v1.DeleteOpenDouyinUsersReply, error) {
	userOpenDouyin, err := uoduc.repo.Delete(ctx, userId, openDouyinUserId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_DELETE_USER_OPEN_DOUYIN_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userOpenDouyin, nil
}

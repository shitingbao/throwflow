package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/weixin/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"
)

type UserSampleOrderRepo interface {
	GetKuaidi(context.Context, uint64) (*v1.GetKuaidiInfoUserSampleOrdersReply, error)
	List(context.Context, uint64, uint64, string, string, string) (*v1.ListUserSampleOrdersReply, error)
	Statistics(context.Context, string, string, string) (*v1.StatisticsUserSampleOrdersReply, error)
	Verify(context.Context, uint64, uint64, uint64) (*v1.VerifyUserSampleOrdersReply, error)
	Cancel(context.Context, uint64, string) (*v1.CancelUserSampleOrdersReply, error)
	Save(context.Context, uint64, uint64, uint64, uint64, string) (*v1.CreateUserSampleOrdersReply, error)
	UpdateKuaidi(context.Context, uint64, string, string) (*v1.UpdateKuaidiUserSampleOrdersReply, error)
}

type UserSampleOrderUsecase struct {
	repo UserSampleOrderRepo
	conf *conf.Data
	log  *log.Helper
}

func NewUserSampleOrderUsecase(repo UserSampleOrderRepo, conf *conf.Data, logger log.Logger) *UserSampleOrderUsecase {
	return &UserSampleOrderUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (usouc *UserSampleOrderUsecase) ListSampleOrders(ctx context.Context, pageNum, pageSize uint64, day, keyword, searchType string) (*v1.ListUserSampleOrdersReply, error) {
	if pageSize == 0 {
		pageSize = uint64(usouc.conf.Database.PageSize)
	}

	list, err := usouc.repo.List(ctx, pageNum, pageSize, day, keyword, searchType)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_SAMPLE_ORDER_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (usouc *UserSampleOrderUsecase) GetKuaidiInfoSampleOrders(ctx context.Context, UserSampleOrderId uint64) (*v1.GetKuaidiInfoUserSampleOrdersReply, error) {
	list, err := usouc.repo.GetKuaidi(ctx, UserSampleOrderId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_KUAIDI_INFO_SAMPLE_ORDER_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (usouc *UserSampleOrderUsecase) StatisticsSampleOrders(ctx context.Context, day, keyword, searchType string) (*v1.StatisticsUserSampleOrdersReply, error) {
	list, err := usouc.repo.Statistics(ctx, day, keyword, searchType)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_statistics_SAMPLE_ORDER_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (usouc *UserSampleOrderUsecase) VerifyUserSampleOrders(ctx context.Context, userId, openDouyinUserId, productId uint64) (*v1.VerifyUserSampleOrdersReply, error) {
	userSampleOrder, err := usouc.repo.Verify(ctx, userId, openDouyinUserId, productId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_VERIFY_USER_SAMPLE_ORDER_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userSampleOrder, nil
}

func (usouc *UserSampleOrderUsecase) CancelSampleOrders(ctx context.Context, UserSampleOrderId uint64, cancelNote string) (*v1.CancelUserSampleOrdersReply, error) {
	userSampleOrder, err := usouc.repo.Cancel(ctx, UserSampleOrderId, cancelNote)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_CANCEL_USER_SAMPLE_ORDER_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userSampleOrder, nil
}

func (usouc *UserSampleOrderUsecase) CreateUserSampleOrders(ctx context.Context, userId, openDouyinUserId, productId, userAddressId uint64, note string) (*v1.CreateUserSampleOrdersReply, error) {
	userSampleOrder, err := usouc.repo.Save(ctx, userId, openDouyinUserId, productId, userAddressId, note)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_CREATE_USER_SAMPLE_ORDER_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userSampleOrder, nil
}

func (usouc *UserSampleOrderUsecase) UpdateKuaidiSampleOrders(ctx context.Context, userSampleOrderId uint64, kuaidiCode, kuaidiNum string) (*v1.UpdateKuaidiUserSampleOrdersReply, error) {
	userSampleOrder, err := usouc.repo.UpdateKuaidi(ctx, userSampleOrderId, kuaidiCode, kuaidiNum)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_UPDATE_KUAIDI_USER_SAMPLE_ORDER_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userSampleOrder, nil
}

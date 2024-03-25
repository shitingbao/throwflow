package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/weixin/v1"
	"interface/internal/conf"
	"interface/internal/pkg/tool"
)

type UserOrganizationRepo interface {
	Get(context.Context, uint64) (*v1.GetUserOrganizationRelationsReply, error)
	GetBind(context.Context, uint64, uint64) (*v1.GetBindUserOrganizationRelationsReply, error)
	List(context.Context, uint64, uint64) (*v1.ListUserOrdersReply, error)
	ListParent(context.Context, uint64, uint64, string) (*v1.ListParentUserOrganizationRelationsReply, error)
	ListCommission(context.Context, uint64, uint64, uint64, uint64, uint32, string, string) (*v1.ListUserCommissionsReply, error)
	Statistics(context.Context, uint64, uint64) (*v1.StatisticsUserCommissionsReply, error)
}

type UserOrganizationUsecase struct {
	repo   UserOrganizationRepo
	uorepo UserOrderRepo
	conf   *conf.Data
	log    *log.Helper
}

func NewUserOrganizationUsecase(repo UserOrganizationRepo, uorepo UserOrderRepo, conf *conf.Data, logger log.Logger) *UserOrganizationUsecase {
	return &UserOrganizationUsecase{repo: repo, uorepo: uorepo, conf: conf, log: log.NewHelper(logger)}
}

func (uouc *UserOrganizationUsecase) GetMiniUserOrganizationRelations(ctx context.Context, userId uint64) (*v1.GetUserOrganizationRelationsReply, error) {
	userOrganization, err := uouc.repo.Get(ctx, userId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_USER_ORGANIZATION_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userOrganization, nil
}

func (uouc *UserOrganizationUsecase) GetBindMiniUserOrganizations(ctx context.Context, userId, organizationId uint64) (*v1.GetBindUserOrganizationRelationsReply, error) {
	bindUserOrganization, err := uouc.repo.GetBind(ctx, userId, organizationId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_GET_BIND_USER_ORGANIZATION_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return bindUserOrganization, nil
}

func (uouc *UserOrganizationUsecase) ListMinUserOrders(ctx context.Context, pageNum, pageSize uint64) (*v1.ListUserOrdersReply, error) {
	if pageSize == 0 {
		pageSize = uint64(uouc.conf.Database.PageSize)
	}

	list, err := uouc.repo.List(ctx, pageNum, pageSize)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_USER_ORDER_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (uouc *UserOrganizationUsecase) ListParentMiniUserOrganizationRelations(ctx context.Context, userId, organizationId uint64, relationType string) (*v1.ListParentUserOrganizationRelationsReply, error) {
	list, err := uouc.repo.ListParent(ctx, userId, organizationId, relationType)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_PARENT_USER_ORGANIZATION_RELATION_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (uouc *UserOrganizationUsecase) ListMiniUserOrganizations(ctx context.Context, pageNum, pageSize, userId, organizationId uint64, isDirect uint32, month, keyword string) (*v1.ListUserCommissionsReply, error) {
	if pageSize == 0 {
		pageSize = uint64(uouc.conf.Database.PageSize)
	}

	list, err := uouc.repo.ListCommission(ctx, pageNum, pageSize, userId, organizationId, isDirect, month, keyword)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_LIST_USER_COMMISSION_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return list, nil
}

func (uouc *UserOrganizationUsecase) CreateMiniUserOrders(ctx context.Context, userId, parentUserId, organizationId uint64, payAmount float64) (*v1.CreateUserOrdersReply, error) {
	clientIp := tool.GetClientIp(ctx)

	userOrder, err := uouc.uorepo.Create(ctx, userId, parentUserId, organizationId, payAmount, clientIp)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_CREATE_USER_ORDER_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userOrder, nil
}

func (uouc *UserOrganizationUsecase) UpgradeMiniUserOrders(ctx context.Context, userId, organizationId uint64, payAmount float64) (*v1.UpgradeUserOrdersReply, error) {
	clientIp := tool.GetClientIp(ctx)

	userOrder, err := uouc.uorepo.Upgrade(ctx, userId, organizationId, payAmount, clientIp)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_UPGRADE_USER_ORDER_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return userOrder, nil
}

func (uouc *UserOrganizationUsecase) StatisticsMiniUserOrganizations(ctx context.Context, userId, organizationId uint64) (*v1.StatisticsUserCommissionsReply, error) {
	statistics, err := uouc.repo.Statistics(ctx, userId, organizationId)

	if err != nil {
		return nil, errors.InternalServer("INTERFACE_STATISTIC_USER_ORGANIZATION_FAILED", tool.GetGRPCErrorInfo(err))
	}

	return statistics, nil
}

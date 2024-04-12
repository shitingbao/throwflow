package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/weixin/v1"
	"interface/internal/biz"
)

type userOrganizationRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserOrganizationRepo(data *Data, logger log.Logger) biz.UserOrganizationRepo {
	return &userOrganizationRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (uor *userOrganizationRepo) Get(ctx context.Context, userId uint64) (*v1.GetUserOrganizationRelationsReply, error) {
	userOrganization, err := uor.data.weixinuc.GetUserOrganizationRelations(ctx, &v1.GetUserOrganizationRelationsRequest{
		UserId: userId,
	})

	if err != nil {
		return nil, err
	}

	return userOrganization, err
}

func (uor *userOrganizationRepo) GetBind(ctx context.Context, userId, organizationId uint64) (*v1.GetBindUserOrganizationRelationsReply, error) {
	userOrganization, err := uor.data.weixinuc.GetBindUserOrganizationRelations(ctx, &v1.GetBindUserOrganizationRelationsRequest{
		UserId:         userId,
		OrganizationId: organizationId,
	})

	if err != nil {
		return nil, err
	}

	return userOrganization, err
}

func (uor *userOrganizationRepo) List(ctx context.Context, pageNum, pageSize uint64) (*v1.ListUserOrdersReply, error) {
	list, err := uor.data.weixinuc.ListUserOrders(ctx, &v1.ListUserOrdersRequest{
		PageNum:  pageNum,
		PageSize: pageSize,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (uor *userOrganizationRepo) ListParent(ctx context.Context, userId, organizationId uint64, relationType string) (*v1.ListParentUserOrganizationRelationsReply, error) {
	list, err := uor.data.weixinuc.ListParentUserOrganizationRelations(ctx, &v1.ListParentUserOrganizationRelationsRequest{
		UserId:         userId,
		OrganizationId: organizationId,
		RelationType:   relationType,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (uor *userOrganizationRepo) ListCommission(ctx context.Context, pageNum, pageSize, userId, organizationId uint64, commissionType uint32, month, keyword string) (*v1.ListUserCommissionsReply, error) {
	list, err := uor.data.weixinuc.ListUserCommissions(ctx, &v1.ListUserCommissionsRequest{
		PageNum:        pageNum,
		PageSize:       pageSize,
		UserId:         userId,
		OrganizationId: organizationId,
		Month:          month,
		CommissionType: commissionType,
		Keyword:        keyword,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (uor *userOrganizationRepo) Statistics(ctx context.Context, userId, organizationId uint64) (*v1.StatisticsUserCommissionsReply, error) {
	statistics, err := uor.data.weixinuc.StatisticsUserCommissions(ctx, &v1.StatisticsUserCommissionsRequest{
		UserId:         userId,
		OrganizationId: organizationId,
	})

	if err != nil {
		return nil, err
	}

	return statistics, err
}

func (uor *userOrganizationRepo) StatisticsDetail(ctx context.Context, userId, organizationId uint64, commissionType uint32, month, keyword string) (*v1.StatisticsDetailUserCommissionsReply, error) {
	statistics, err := uor.data.weixinuc.StatisticsDetailUserCommissions(ctx, &v1.StatisticsDetailUserCommissionsRequest{
		UserId:         userId,
		OrganizationId: organizationId,
		Month:          month,
		CommissionType: commissionType,
		Keyword:        keyword,
	})

	if err != nil {
		return nil, err
	}

	return statistics, err
}

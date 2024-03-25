package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "interface/api/service/weixin/v1"
	"interface/internal/biz"
)

type userSampleOrderRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserSampleOrderRepo(data *Data, logger log.Logger) biz.UserSampleOrderRepo {
	return &userSampleOrderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (usor *userSampleOrderRepo) GetKuaidi(ctx context.Context, userSampleOrderId uint64) (*v1.GetKuaidiInfoUserSampleOrdersReply, error) {
	list, err := usor.data.weixinuc.GetKuaidiInfoUserSampleOrders(ctx, &v1.GetKuaidiInfoUserSampleOrdersRequest{
		UserSampleOrderId: userSampleOrderId,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (usor *userSampleOrderRepo) List(ctx context.Context, pageNum, pageSize uint64, day, keyword, searchType string) (*v1.ListUserSampleOrdersReply, error) {
	list, err := usor.data.weixinuc.ListUserSampleOrders(ctx, &v1.ListUserSampleOrdersRequest{
		PageNum:    pageNum,
		PageSize:   pageSize,
		Day:        day,
		Keyword:    keyword,
		SearchType: searchType,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (usor *userSampleOrderRepo) Statistics(ctx context.Context, day, keyword, searchType string) (*v1.StatisticsUserSampleOrdersReply, error) {
	list, err := usor.data.weixinuc.StatisticsUserSampleOrders(ctx, &v1.StatisticsUserSampleOrdersRequest{
		Day:        day,
		Keyword:    keyword,
		SearchType: searchType,
	})

	if err != nil {
		return nil, err
	}

	return list, err
}

func (usor *userSampleOrderRepo) Verify(ctx context.Context, userId, openDouyinUserId, productId uint64) (*v1.VerifyUserSampleOrdersReply, error) {
	sampleOrder, err := usor.data.weixinuc.VerifyUserSampleOrders(ctx, &v1.VerifyUserSampleOrdersRequest{
		UserId:           userId,
		OpenDouyinUserId: openDouyinUserId,
		ProductId:        productId,
	})

	if err != nil {
		return nil, err
	}

	return sampleOrder, err
}

func (usor *userSampleOrderRepo) Cancel(ctx context.Context, userSampleOrderId uint64, cancelNote string) (*v1.CancelUserSampleOrdersReply, error) {
	sampleOrder, err := usor.data.weixinuc.CancelUserSampleOrders(ctx, &v1.CancelUserSampleOrdersRequest{
		UserSampleOrderId: userSampleOrderId,
		CancelNote:        cancelNote,
	})

	if err != nil {
		return nil, err
	}

	return sampleOrder, err
}

func (usor *userSampleOrderRepo) Save(ctx context.Context, userId, openDouyinUserId, productId, userAddressId uint64, note string) (*v1.CreateUserSampleOrdersReply, error) {
	sampleOrder, err := usor.data.weixinuc.CreateUserSampleOrders(ctx, &v1.CreateUserSampleOrdersRequest{
		UserId:           userId,
		OpenDouyinUserId: openDouyinUserId,
		ProductId:        productId,
		UserAddressId:    userAddressId,
		Note:             note,
	})

	if err != nil {
		return nil, err
	}

	return sampleOrder, err
}

func (usor *userSampleOrderRepo) UpdateKuaidi(ctx context.Context, userSampleOrderId uint64, kuaidiCode, kuaidiNum string) (*v1.UpdateKuaidiUserSampleOrdersReply, error) {
	sampleOrder, err := usor.data.weixinuc.UpdateKuaidiUserSampleOrders(ctx, &v1.UpdateKuaidiUserSampleOrdersRequest{
		UserSampleOrderId: userSampleOrderId,
		KuaidiCode:        kuaidiCode,
		KuaidiNum:         kuaidiNum,
	})

	if err != nil {
		return nil, err
	}

	return sampleOrder, err
}

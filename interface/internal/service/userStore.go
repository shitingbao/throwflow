package service

import (
	"context"
	v1 "interface/api/interface/v1"
	"interface/internal/biz"
)

func (is *InterfaceService) ListUserStores(ctx context.Context, in *v1.ListUserStoresRequest) (*v1.ListUserStoresReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	userStores, err := is.usuc.ListUserStores(ctx, in.PageNum, in.PageSize, userInfo.Data.UserId)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListUserStoresReply_Store, 0)

	for _, userStore := range userStores.Data.List {
		list = append(list, &v1.ListUserStoresReply_Store{
			ClientKey:       userStore.ClientKey,
			OpenId:          userStore.OpenId,
			ProductId:       userStore.ProductId,
			ProductName:     userStore.ProductName,
			ProductImg:      userStore.ProductImg,
			ProductPrice:    userStore.ProductPrice,
			CommissionRatio: userStore.CommissionRatio,
		})
	}

	return &v1.ListUserStoresReply{
		Code: 200,
		Data: &v1.ListUserStoresReply_Data{
			PageNum:   userStores.Data.PageNum,
			PageSize:  userStores.Data.PageSize,
			Total:     userStores.Data.Total,
			TotalPage: userStores.Data.TotalPage,
			List:      list,
		},
	}, nil
}

func (is *InterfaceService) CreateUserStores(ctx context.Context, in *v1.CreateUserStoresRequest) (*v1.CreateUserStoresReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	userStore, err := is.usuc.CreateUserStores(ctx, userInfo.Data.UserId, in.ProductId, in.OpenDouyinUserIds)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.CreateUserStoresReply_Message, 0)

	for _, message := range userStore.Data.List {
		list = append(list, &v1.CreateUserStoresReply_Message{
			ProductName: message.ProductName,
			AwemeName:   message.AwemeName,
			Content:     message.Content,
		})
	}

	return &v1.CreateUserStoresReply{
		Code: 200,
		Data: &v1.CreateUserStoresReply_Data{
			List:    list,
			Content: userStore.Data.Content,
		},
	}, nil
}

func (is *InterfaceService) UpdateUserStores(ctx context.Context, in *v1.UpdateUserStoresRequest) (*v1.UpdateUserStoresReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	userStore, err := is.usuc.UpdateUserStores(ctx, userInfo.Data.UserId, in.Stores)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateUserStoresReply{
		Code: 200,
		Data: &v1.UpdateUserStoresReply_Data{
			Content: userStore.Data.Content,
		},
	}, nil
}

func (is *InterfaceService) DeleteUserStores(ctx context.Context, in *v1.DeleteUserStoresRequest) (*v1.DeleteUserStoresReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	userStore, err := is.usuc.DeleteUserStores(ctx, userInfo.Data.UserId, in.Stores)

	if err != nil {
		return nil, err
	}

	return &v1.DeleteUserStoresReply{
		Code: 200,
		Data: &v1.DeleteUserStoresReply_Data{
			Content: userStore.Data.Content,
		},
	}, nil
}

package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "interface/api/interface/v1"
	"interface/internal/biz"
)

func (is *InterfaceService) GetUrlUserOpenDouyins(ctx context.Context, in *emptypb.Empty) (*v1.GetUrlUserOpenDouyinsReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	url, err := is.uoduc.GetUrlUserOpenDouyins(ctx, userInfo.Data.UserId)

	if err != nil {
		return nil, err
	}

	return &v1.GetUrlUserOpenDouyinsReply{
		Code: 200,
		Data: &v1.GetUrlUserOpenDouyinsReply_Data{
			RedirectUrl: url.Data.RedirectUrl,
		},
	}, nil
}

func (is *InterfaceService) GetTicketUserOpenDouyins(ctx context.Context, in *emptypb.Empty) (*v1.GetTicketUserOpenDouyinsReply, error) {
	ticket, err := is.uoduc.GetTicketUserOpenDouyins(ctx)

	if err != nil {
		return nil, err
	}

	return &v1.GetTicketUserOpenDouyinsReply{
		Code: 200,
		Data: &v1.GetTicketUserOpenDouyinsReply_Data{
			Ticket: ticket.Data.Ticket,
		},
	}, nil
}

func (is *InterfaceService) GetQrCodeUserOpenDouyins(ctx context.Context, in *emptypb.Empty) (*v1.GetQrCodeUserOpenDouyinsReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	qrCode, err := is.uoduc.GetQrCodeUserOpenDouyins(ctx, userInfo.Data.UserId)

	if err != nil {
		return nil, err
	}

	return &v1.GetQrCodeUserOpenDouyinsReply{
		Code: 200,
		Data: &v1.GetQrCodeUserOpenDouyinsReply_Data{
			QrCode: qrCode.Data.QrCode,
			State:  qrCode.Data.State,
			Token:  qrCode.Data.Token,
		},
	}, nil
}

func (is *InterfaceService) GetStatusQrCodeUserOpenDouyins(ctx context.Context, in *v1.GetStatusQrCodeUserOpenDouyinsRequest) (*v1.GetStatusQrCodeUserOpenDouyinsReply, error) {
	if _, err := is.verifyMiniUserLogin(ctx); err != nil {
		return nil, err
	}

	qrCodeStatus, err := is.uoduc.GetStatusQrCodeUserOpenDouyins(ctx, in.Token, in.State, in.Timestamp)

	if err != nil {
		return nil, err
	}

	return &v1.GetStatusQrCodeUserOpenDouyinsReply{
		Code: 200,
		Data: &v1.GetStatusQrCodeUserOpenDouyinsReply_Data{
			Code:   qrCodeStatus.Data.Code,
			Status: qrCodeStatus.Data.Status,
		},
	}, nil
}

func (is *InterfaceService) ListUserOpenDouyins(ctx context.Context, in *v1.ListUserOpenDouyinsRequest) (*v1.ListUserOpenDouyinsReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	userOpenDouyins, err := is.uoduc.ListUserOpenDouyins(ctx, in.PageNum, in.PageSize, userInfo.Data.UserId, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListUserOpenDouyinsReply_OpenDouyinUser, 0)

	for _, userOpenDouyin := range userOpenDouyins.Data.List {
		list = append(list, &v1.ListUserOpenDouyinsReply_OpenDouyinUser{
			OpenDouyinUserId: userOpenDouyin.OpenDouyinUserId,
			AccountId:        userOpenDouyin.AccountId,
			Nickname:         userOpenDouyin.Nickname,
			Avatar:           userOpenDouyin.Avatar,
			AvatarLarger:     userOpenDouyin.AvatarLarger,
			CooperativeCode:  userOpenDouyin.CooperativeCode,
			Fans:             userOpenDouyin.FansShow,
		})
	}

	return &v1.ListUserOpenDouyinsReply{
		Code: 200,
		Data: &v1.ListUserOpenDouyinsReply_Data{
			PageNum:   userOpenDouyins.Data.PageNum,
			PageSize:  userOpenDouyins.Data.PageSize,
			Total:     userOpenDouyins.Data.Total,
			TotalPage: userOpenDouyins.Data.TotalPage,
			List:      list,
		},
	}, nil
}

func (is *InterfaceService) CreateUserOpenDouyins(ctx context.Context, in *v1.CreateUserOpenDouyinsRequest) (*v1.CreateUserOpenDouyinsReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	if _, err := is.uoduc.CreateUserOpenDouyins(ctx, in.State, in.Code); err != nil {
		return nil, err
	}

	return &v1.CreateUserOpenDouyinsReply{
		Code: 200,
		Data: &v1.CreateUserOpenDouyinsReply_Data{},
	}, nil
}

func (is *InterfaceService) UpdateCooperativeCodeUserOpenDouyins(ctx context.Context, in *v1.UpdateCooperativeCodeUserOpenDouyinsRequest) (*v1.UpdateCooperativeCodeUserOpenDouyinsReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	userOpenDouyins, err := is.uoduc.UpdateCooperativeCodeUserOpenDouyins(ctx, userInfo.Data.UserId, in.OpenDouyinUserId, in.CooperativeCode)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.UpdateCooperativeCodeUserOpenDouyinsReply_OpenDouyinUser, 0)

	for _, userOpenDouyin := range userOpenDouyins.Data.List {
		list = append(list, &v1.UpdateCooperativeCodeUserOpenDouyinsReply_OpenDouyinUser{
			OpenDouyinUserId: userOpenDouyin.OpenDouyinUserId,
			AccountId:        userOpenDouyin.AccountId,
			Nickname:         userOpenDouyin.Nickname,
			Avatar:           userOpenDouyin.Avatar,
			AvatarLarger:     userOpenDouyin.AvatarLarger,
			CooperativeCode:  userOpenDouyin.CooperativeCode,
		})
	}

	return &v1.UpdateCooperativeCodeUserOpenDouyinsReply{
		Code: 200,
		Data: &v1.UpdateCooperativeCodeUserOpenDouyinsReply_Data{
			List: list,
		},
	}, nil
}

func (is *InterfaceService) DeleteUserOpenDouyins(ctx context.Context, in *v1.DeleteUserOpenDouyinsRequest) (*v1.DeleteUserOpenDouyinsReply, error) {
	userInfo, err := is.verifyMiniUserLogin(ctx)

	if err != nil {
		return nil, err
	}

	if _, err := is.uoduc.DeleteUserOpenDouyins(ctx, userInfo.Data.UserId, in.OpenDouyinUserId); err != nil {
		return nil, err
	}

	return &v1.DeleteUserOpenDouyinsReply{
		Code: 200,
		Data: &v1.DeleteUserOpenDouyinsReply_Data{},
	}, nil
}

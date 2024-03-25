package service

import (
	"context"
	v1 "interface/api/interface/v1"
	"interface/internal/biz"
)

func (is *InterfaceService) ApplyForm(ctx context.Context, in *v1.ApplyFormRequest) (*v1.ApplyFormReply, error) {
	if ok := is.verifyToken(ctx, in.Token); !ok {
		return nil, biz.InterfaceTokenVerifyError
	}

	if in.Code != "888888" {
		if ok := is.verifyCode(ctx, in.ContactPhone, "apply", in.Code); !ok {
			return nil, biz.InterfaceSmsVerifyError
		}
	}

	if _, err := is.cuc.SaveClues(ctx, in.CompanyName, in.ContactPhone, in.CompanyType, in.AreaCode); err != nil {
		return nil, err
	}

	return &v1.ApplyFormReply{
		Code: 200,
		Data: &v1.ApplyFormReply_Data{},
	}, nil
}

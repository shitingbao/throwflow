package service

import (
	"context"
	v1 "interface/api/interface/v1"
)

func (is *InterfaceService) ListKuaidiCompany(ctx context.Context, in *v1.ListKuaidiCompanyRequest) (*v1.ListKuaidiCompanyReply, error) {
	kuaidiCompanies, err := is.kcuc.ListKuaidiCompany(ctx, in.PageNum, in.PageSize, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListKuaidiCompanyReply_KuaidiCompany, 0)

	for _, kuaidiCompany := range kuaidiCompanies.Data.List {
		list = append(list, &v1.ListKuaidiCompanyReply_KuaidiCompany{
			Name: kuaidiCompany.Name,
			Code: kuaidiCompany.Code,
		})
	}

	return &v1.ListKuaidiCompanyReply{
		Code: 200,
		Data: &v1.ListKuaidiCompanyReply_Data{
			PageNum:   kuaidiCompanies.Data.PageNum,
			PageSize:  kuaidiCompanies.Data.PageSize,
			Total:     kuaidiCompanies.Data.Total,
			TotalPage: kuaidiCompanies.Data.TotalPage,
			List:      list,
		},
	}, nil
}

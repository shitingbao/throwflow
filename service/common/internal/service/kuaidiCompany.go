package service

import (
	v1 "common/api/common/v1"
	"context"
	"math"
)

func (cs *CommonService) ListKuaidiCompanys(ctx context.Context, in *v1.ListKuaidiCompanysRequest) (*v1.ListKuaidiCompanysReply, error) {
	kuaidiCompanies, err := cs.kcuc.ListKuaidiCompanys(ctx, in.PageNum, in.PageSize, in.Keyword)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListKuaidiCompanysReply_KuaidiCompany, 0)

	for _, kuaidiCompany := range kuaidiCompanies.List {
		list = append(list, &v1.ListKuaidiCompanysReply_KuaidiCompany{
			Name: kuaidiCompany.Name,
			Code: kuaidiCompany.Code,
		})
	}

	totalPage := uint64(math.Ceil(float64(kuaidiCompanies.Total) / float64(kuaidiCompanies.PageSize)))

	return &v1.ListKuaidiCompanysReply{
		Code: 200,
		Data: &v1.ListKuaidiCompanysReply_Data{
			PageNum:   kuaidiCompanies.PageNum,
			PageSize:  kuaidiCompanies.PageSize,
			Total:     kuaidiCompanies.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

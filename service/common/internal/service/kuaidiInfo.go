package service

import (
	v1 "common/api/common/v1"
	"context"
)

func (cs *CommonService) GetKuaidiInfos(ctx context.Context, in *v1.GetKuaidiInfosRequest) (*v1.GetKuaidiInfosReply, error) {
	kuaidiInfoData, err := cs.kiuc.GetKuaidiInfos(ctx, in.Code, in.Num, in.Phone)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.GetKuaidiInfosReply_KuaidiInfo, 0)

	for _, kuaidiInfo := range kuaidiInfoData.Content {
		list = append(list, &v1.GetKuaidiInfosReply_KuaidiInfo{
			Time:    kuaidiInfo.Time,
			Content: kuaidiInfo.Context,
		})
	}

	return &v1.GetKuaidiInfosReply{
		Code: 200,
		Data: &v1.GetKuaidiInfosReply_Data{
			Code:      kuaidiInfoData.Code,
			Name:      kuaidiInfoData.Name,
			Num:       kuaidiInfoData.Num,
			State:     uint32(kuaidiInfoData.State),
			StateName: kuaidiInfoData.StateName,
			List:      list,
		},
	}, nil
}

package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "interface/api/interface/v1"
	"interface/internal/pkg/tool"
)

func (is *InterfaceService) ListUpdateLogs(ctx context.Context, in *emptypb.Empty) (*v1.ListUpdateLogsReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, "updatelog"); err != nil {
		return nil, err
	}

	updateLogs, err := is.uluc.ListUpdateLogs(ctx)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListUpdateLogsReply_UpdateLogs, 0)

	for _, updateLog := range updateLogs.Data.List {
		createTime, _ := tool.StringToTime("2006-01-02 15:04:05", updateLog.CreateTime)

		list = append(list, &v1.ListUpdateLogsReply_UpdateLogs{
			Name:       updateLog.Name,
			Content:    updateLog.Content,
			CreateTime: tool.TimeToString("2006-01-02 15:04", createTime),
		})
	}

	return &v1.ListUpdateLogsReply{
		Code: 200,
		Data: &v1.ListUpdateLogsReply_Data{
			List: list,
		},
	}, nil
}

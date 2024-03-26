package service

import (
	v1 "common/api/common/v1"
	"common/internal/pkg/tool"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (cs *CommonService) ListUpdateLogs(ctx context.Context, in *emptypb.Empty) (*v1.ListUpdateLogsReply, error) {
	updateLogs, err := cs.uluc.ListUpdateLogs(ctx)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListUpdateLogsReply_UpdateLogs, 0)

	for _, updateLog := range updateLogs {
		list = append(list, &v1.ListUpdateLogsReply_UpdateLogs{
			Id:         updateLog.Id,
			Name:       updateLog.Name,
			Content:    updateLog.Content,
			CreateTime: tool.TimeToString("2006-01-02 15:04:05", updateLog.CreateTime),
			UpdateTime: tool.TimeToString("2006-01-02 15:04:05", updateLog.UpdateTime),
		})
	}

	return &v1.ListUpdateLogsReply{
		Code: 200,
		Data: &v1.ListUpdateLogsReply_Data{
			List: list,
		},
	}, nil
}

func (cs *CommonService) CreateUpdateLogs(ctx context.Context, in *v1.CreateUpdateLogsRequest) (*v1.CreateUpdateLogsReply, error) {
	updateLog, err := cs.uluc.CreateUpdateLogs(ctx, in.Name, in.Content)

	if err != nil {
		return nil, err
	}

	return &v1.CreateUpdateLogsReply{
		Code: 200,
		Data: &v1.CreateUpdateLogsReply_Data{
			Id:         updateLog.Id,
			Name:       updateLog.Name,
			Content:    updateLog.Content,
			CreateTime: tool.TimeToString("2006-01-02 15:04:05", updateLog.CreateTime),
			UpdateTime: tool.TimeToString("2006-01-02 15:04:05", updateLog.UpdateTime),
		},
	}, nil
}

func (cs *CommonService) UpdateUpdateLogs(ctx context.Context, in *v1.UpdateUpdateLogsRequest) (*v1.UpdateUpdateLogsReply, error) {
	updateLog, err := cs.uluc.UpdateUpdateLogs(ctx, in.Id, in.Name, in.Content)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateUpdateLogsReply{
		Code: 200,
		Data: &v1.UpdateUpdateLogsReply_Data{
			Id:         updateLog.Id,
			Name:       updateLog.Name,
			Content:    updateLog.Content,
			CreateTime: tool.TimeToString("2006-01-02 15:04:05", updateLog.CreateTime),
			UpdateTime: tool.TimeToString("2006-01-02 15:04:05", updateLog.UpdateTime),
		},
	}, nil
}

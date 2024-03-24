package service

import (
	v1 "admin/api/admin/v1"
	"admin/internal/biz"
	"admin/internal/pkg/tool"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (as *AdminService) ListUpdateLogs(ctx context.Context, in *emptypb.Empty) (*v1.ListUpdateLogsReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:updateLog:list"); err != nil {
		return nil, err
	}

	updateLogs, err := as.uluc.ListUpdateLogs(ctx)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListUpdateLogsReply_UpdateLogs, 0)

	for _, updateLog := range updateLogs.Data.List {
		createTime, _ := tool.StringToTime("2006-01-02 15:04:05", updateLog.CreateTime)

		list = append(list, &v1.ListUpdateLogsReply_UpdateLogs{
			Id:         updateLog.Id,
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

func (as *AdminService) CreateUpdateLogs(ctx context.Context, in *v1.CreateUpdateLogsRequest) (*v1.CreateUpdateLogsReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:updateLog:create"); err != nil {
		return nil, err
	}

	if ok := as.verifyToken(ctx, in.Token); !ok {
		return nil, biz.AdminTokenVerifyError
	}

	updateLog, err := as.uluc.CreateUpdateLogs(ctx, in.Name, in.Content)

	if err != nil {
		return nil, err
	}

	createTime, _ := tool.StringToTime("2006-01-02 15:04:05", updateLog.Data.CreateTime)

	return &v1.CreateUpdateLogsReply{
		Code: 200,
		Data: &v1.CreateUpdateLogsReply_Data{
			Id:         updateLog.Data.Id,
			Name:       updateLog.Data.Name,
			Content:    updateLog.Data.Content,
			CreateTime: tool.TimeToString("2006-01-02 15:04", createTime),
		},
	}, nil
}

func (as *AdminService) UpdateUpdateLogs(ctx context.Context, in *v1.UpdateUpdateLogsRequest) (*v1.UpdateUpdateLogsReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:updateLog:update"); err != nil {
		return nil, err
	}

	updateLog, err := as.uluc.UpdateUpdateLogs(ctx, in.Id, in.Name, in.Content)

	if err != nil {
		return nil, err
	}

	createTime, _ := tool.StringToTime("2006-01-02 15:04:05", updateLog.Data.CreateTime)

	return &v1.UpdateUpdateLogsReply{
		Code: 200,
		Data: &v1.UpdateUpdateLogsReply_Data{
			Id:         updateLog.Data.Id,
			Name:       updateLog.Data.Name,
			Content:    updateLog.Data.Content,
			CreateTime: tool.TimeToString("2006-01-02 15:04", createTime),
		},
	}, nil
}

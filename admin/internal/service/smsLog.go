package service

import (
	v1 "admin/api/admin/v1"
	"context"
	"math"
)

func (as *AdminService) ListSmsLogs(ctx context.Context, in *v1.ListSmsLogsRequest) (*v1.ListSmsLogsReply, error) {
	if _, err := as.verifyPermission(ctx, "admin:smslog:list"); err != nil {
		return nil, err
	}

	smsLogs, err := as.sluc.ListSmsLogs(ctx, in.PageNum)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListSmsLogsReply_SmsLogs, 0)

	for _, smsLog := range smsLogs.Data.List {
		list = append(list, &v1.ListSmsLogsReply_SmsLogs{
			Id:         smsLog.Id,
			Phone:      smsLog.Phone,
			Content:    smsLog.Content,
			Reply:      smsLog.Reply,
			Type:       formatType(smsLog.Type),
			Ip:         smsLog.Ip,
			CreateTime: smsLog.CreateTime,
			UpdateTime: smsLog.UpdateTime,
		})
	}

	totalPage := uint64(math.Ceil(float64(smsLogs.Data.Total) / float64(smsLogs.Data.PageSize)))

	return &v1.ListSmsLogsReply{
		Code: 200,
		Data: &v1.ListSmsLogsReply_Data{
			PageNum:   smsLogs.Data.PageNum,
			PageSize:  smsLogs.Data.PageSize,
			Total:     smsLogs.Data.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func formatType(types string) string {
	if types == "login" {
		return "登录验证码短信"
	} else if types == "apply" {
		return "申请试用验证码短信"
	} else if types == "marketing" {
		return "营销短信"
	}
	return ""
}

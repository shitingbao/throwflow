package service

import (
	"common/internal/pkg/tool"
	"context"
	"math"

	v1 "common/api/common/v1"
)

func (cs *CommonService) SendSms(ctx context.Context, in *v1.SendSmsRequest) (*v1.SendSmsReply, error) {
	_, err := cs.suc.SendSms(ctx, in.Phone, in.Types, in.Content, in.Ip)

	if err != nil {
		return nil, err
	}

	return &v1.SendSmsReply{
		Code: 200,
		Data: &v1.SendSmsReply_Data{},
	}, nil
}

func (cs *CommonService) VerifySms(ctx context.Context, in *v1.VerifySmsRequest) (*v1.VerifySmsReply, error) {
	err := cs.suc.VerifySms(ctx, in.Phone, in.Types, in.Code)

	if err != nil {
		return nil, err
	}

	return &v1.VerifySmsReply{
		Code: 200,
		Data: &v1.VerifySmsReply_Data{},
	}, nil
}

func (cs *CommonService) ListSms(ctx context.Context, in *v1.ListSmsRequest) (*v1.ListSmsReply, error) {
	smsLogs, err := cs.suc.ListSmsLogs(ctx, in.PageNum, in.PageSize)

	if err != nil {
		return nil, err
	}

	list := make([]*v1.ListSmsReply_Smss, 0)

	for _, smsLog := range smsLogs.List {
		list = append(list, &v1.ListSmsReply_Smss{
			Id:         smsLog.Id,
			Phone:      smsLog.SendPhone,
			Content:    smsLog.SendContent,
			Reply:      smsLog.ReplyContent,
			Type:       smsLog.SendType,
			Ip:         smsLog.SendIp,
			CreateTime: tool.TimeToString("2006-01-02 15:04:05", smsLog.CreateTime),
			UpdateTime: tool.TimeToString("2006-01-02 15:04:05", smsLog.UpdateTime),
		})
	}

	totalPage := uint64(math.Ceil(float64(smsLogs.Total) / float64(smsLogs.PageSize)))

	return &v1.ListSmsReply{
		Code: 200,
		Data: &v1.ListSmsReply_Data{
			PageNum:   smsLogs.PageNum,
			PageSize:  smsLogs.PageSize,
			Total:     smsLogs.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

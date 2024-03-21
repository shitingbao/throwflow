package service

import (
	v1 "company/api/company/v1"
	"company/internal/pkg/tool"
	"context"
	"math"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (cs *CompanyService) CreateCompanyTask(ctx context.Context, in *v1.CreateCompanyTaskRequest) (*v1.CreateCompanyTaskReply, error) {
	task, err := cs.ctuc.CreateCompanyTask(ctx, in.ProductOutId, in.ExpireTime, in.PlayNum, in.Price, in.Quota, uint8(in.IsGoodReviews))

	if err != nil {
		return nil, err
	}

	return &v1.CreateCompanyTaskReply{
		Code: 200,
		Data: &v1.CreateCompanyTaskReply_Data{
			Id:            task.Id,
			ProductOutId:  task.ProductOutId,
			ExpireTime:    task.ExpireTime,
			PlayNum:       task.PlayNum,
			Price:         task.Price,
			Quota:         task.Quota,
			ClaimQuota:    task.ClaimQuota,
			SuccessQuota:  task.SuccessQuota,
			IsTop:         uint32(task.IsTop),
			IsDel:         uint32(task.IsDel),
			IsGoodReviews: uint32(task.IsGoodReviews),
			CreateTime:    tool.TimeToString("2006-01-02 15:04:05", task.CreateTime),
		},
	}, nil
}

func (cs *CompanyService) GetByProductOutId(ctx context.Context, in *v1.GetByProductOutIdRequest) (*v1.GetByProductOutIdReply, error) {
	existCount, task, err := cs.ctuc.GetByProductOutId(ctx, in.ProductOutId, in.UserId)

	if err != nil {
		return nil, err
	}

	return &v1.GetByProductOutIdReply{
		Code: 200,
		Data: &v1.GetByProductOutIdReply_Data{
			Id:            task.Id,
			ProductOutId:  task.ProductOutId,
			ExpireTime:    task.ExpireTime,
			PlayNum:       task.PlayNum,
			Price:         task.Price,
			Quota:         task.Quota,
			ClaimQuota:    task.ClaimQuota,
			SuccessQuota:  task.SuccessQuota,
			IsTop:         uint32(task.IsTop),
			IsDel:         uint32(task.IsDel),
			IsGoodReviews: uint32(task.IsGoodReviews),
			CreateTime:    tool.TimeToString("2006-01-02 15:04:05", task.CreateTime),
			IsExist:       uint32(existCount),
		},
	}, nil
}

func (cs *CompanyService) UpdateCompanyTaskQuota(ctx context.Context, in *v1.UpdateCompanyTaskQuotaRequest) (*v1.UpdateCompanyTaskReply, error) {
	task, err := cs.ctuc.UpdateCompanyTaskQuota(ctx, in.CompanyTaskId, in.Quota)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateCompanyTaskReply{Code: 200,
		Data: &v1.UpdateCompanyTaskReply_Data{
			Id:            task.Id,
			ProductOutId:  task.ProductOutId,
			ExpireTime:    task.ExpireTime,
			PlayNum:       task.PlayNum,
			Price:         task.Price,
			Quota:         task.Quota,
			ClaimQuota:    task.ClaimQuota,
			SuccessQuota:  task.SuccessQuota,
			IsTop:         uint32(task.IsTop),
			IsDel:         uint32(task.IsDel),
			IsGoodReviews: uint32(task.IsGoodReviews),
			CreateTime:    tool.TimeToString("2006-01-02 15:04:05", task.CreateTime),
		}}, nil
}

func (cs *CompanyService) UpdateCompanyTaskIsTop(ctx context.Context, in *v1.UpdateCompanyTaskIsTopRequest) (*v1.UpdateCompanyTaskReply, error) {
	task, err := cs.ctuc.UpdateCompanyTaskIsTop(ctx, in.CompanyTaskId, uint8(in.IsTop))

	if err != nil {
		return nil, err
	}

	return &v1.UpdateCompanyTaskReply{Code: 200,
		Data: &v1.UpdateCompanyTaskReply_Data{
			Id:            task.Id,
			ProductOutId:  task.ProductOutId,
			ExpireTime:    task.ExpireTime,
			PlayNum:       task.PlayNum,
			Price:         task.Price,
			Quota:         task.Quota,
			ClaimQuota:    task.ClaimQuota,
			SuccessQuota:  task.SuccessQuota,
			IsTop:         uint32(task.IsTop),
			IsDel:         uint32(task.IsDel),
			IsGoodReviews: uint32(task.IsGoodReviews),
			CreateTime:    tool.TimeToString("2006-01-02 15:04:05", task.CreateTime),
		}}, nil
}

func (cs *CompanyService) ListCompanyTask(ctx context.Context, in *v1.ListCompanyTaskRequest) (*v1.ListCompanyTaskReply, error) {
	tasks, err := cs.ctuc.ListCompanyTask(ctx, int(in.IsDel), in.IsTop, in.PageNum, in.PageSize, in.Keyword)

	if err != nil {
		return nil, err
	}

	totalPage := uint64(math.Ceil(float64(tasks.Total) / float64(tasks.PageSize)))
	list := []*v1.ListCompanyTaskReply_CompanyTask{}

	for _, v := range tasks.List {
		c := &v1.ListCompanyTaskReply_CompanyTask_CompanyProduct{
			ProductOutId:          v.CompanyProduct.ProductOutId,
			ProductName:           v.CompanyProduct.ProductName,
			ProductPrice:          v.CompanyProduct.ProductPrice,
			PureCommission:        v.CompanyProduct.PureCommission,
			PureServiceCommission: v.CompanyProduct.PureServiceCommission,
		}

		if len(v.CompanyProduct.ProductImgs) > 0 {
			c.ProductImg = v.CompanyProduct.ProductImgs[0]
		}

		task := &v1.ListCompanyTaskReply_CompanyTask{
			Id:             v.Id,
			ProductOutId:   v.ProductOutId,
			ExpireTime:     v.ExpireTime,
			PlayNum:        v.PlayNum,
			Price:          v.Price,
			Quota:          v.Quota,
			ClaimQuota:     v.ClaimQuota,
			SuccessQuota:   v.SuccessQuota,
			IsDel:          uint32(v.IsDel),
			IsTop:          uint32(v.IsTop),
			IsGoodReviews:  uint32(v.IsGoodReviews),
			CreateTime:     tool.TimeToString("2006-01-02 15:04:05", v.CreateTime),
			CompanyProduct: c,
		}

		list = append(list, task)
	}

	return &v1.ListCompanyTaskReply{
		Code: 200,
		Data: &v1.ListCompanyTaskReply_Data{
			PageNum:   tasks.PageNum,
			PageSize:  tasks.PageSize,
			Total:     tasks.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

func (cs *CompanyService) UpdateCompanyTaskIsDel(ctx context.Context, in *v1.UpdateCompanyTaskIsDelRequest) (*v1.UpdateCompanyTaskIsDelReply, error) {
	if err := cs.ctuc.UpdateCompanyTaskIsDel(ctx, in.CompanyTaskId); err != nil {
		return nil, err
	}

	return &v1.UpdateCompanyTaskIsDelReply{
		Code: 200,
	}, nil
}

func (cs *CompanyService) CreateCompanyTaskAccountRelation(ctx context.Context, in *v1.CreateCompanyTaskAccountRelationRequest) (*v1.CreateCompanyTaskAccountRelationReply, error) {
	relation, err := cs.ctaruc.CreateCompanyTaskAccountRelation(ctx, in.CompanyTaskId, in.ProductOutId, in.UserId)

	if err != nil {
		return nil, err
	}

	return &v1.CreateCompanyTaskAccountRelationReply{
		Code: 200,
		Data: &v1.CreateCompanyTaskAccountRelationReply_Data{
			Id:                    relation.Id,
			CompanyTaskId:         relation.CompanyTaskId,
			ProductName:           relation.ProductName,
			ProductOutId:          relation.ProductOutId,
			UserId:                relation.UserId,
			ClaimTime:             tool.TimeToString("2006-01-02 15:04:05", relation.ClaimTime),
			ExpireTime:            tool.TimeToString("2006-01-02 15:04:05", relation.ExpireTime),
			Status:                uint32(relation.Status),
			IsDel:                 uint32(relation.IsDel),
			IsCostBuy:             uint32(relation.IsCostBuy),
			ScreenshotAddress:     relation.ScreenshotAddress,
			IsScreenshotAvailable: uint32(relation.IsScreenshotAvailable),
			CreateTime:            tool.TimeToString("2006-01-02 15:04:05", relation.CreateTime),
			UpdateTime:            tool.TimeToString("2006-01-02 15:04:05", relation.UpdateTime),
		},
	}, nil
}

func (cs *CompanyService) ListCompanyTaskAccountRelation(ctx context.Context, in *v1.ListCompanyTaskAccountRelationRequest) (*v1.ListCompanyTaskAccountRelationReply, error) {
	relation, err := cs.ctaruc.ListCompanyTaskAccountRelation(ctx, in.Status, in.CompanyTaskId, in.UserId, in.PageNum, in.PageSize, in.ExpireTime, in.ExpiredTime, in.ProductName)

	if err != nil {
		return nil, err
	}

	list := []*v1.ListCompanyTaskAccountRelationReply_CompanyTaskAccountRelation{}

	for _, v := range relation.List {
		dts := []*v1.ListCompanyTaskAccountRelationReply_CompanyTaskAccountRelation_CompanyTaskDetail{}
		var isPlayCount uint32 = 0

		for _, det := range v.CompanyTaskDetails {

			if det.PlayCount >= v.CompanyTask.PlayNum {
				isPlayCount = 1
			}

			dts = append(dts, &v1.ListCompanyTaskAccountRelationReply_CompanyTaskAccountRelation_CompanyTaskDetail{
				Id:                           det.Id,
				CompanyTaskId:                det.CompanyTaskId,
				UserId:                       det.UserId,
				ClientKey:                    det.ClientKey,
				OpenId:                       det.OpenId,
				CompanyTaskAccountRelationId: det.CompanyTaskAccountRelationId,
				ProductName:                  det.ProductName,
				ItemId:                       det.ItemId,
				PlayCount:                    det.PlayCount,
				Cover:                        det.Cover,
				ReleaseTime:                  tool.TimeToString("2006-01-02 15:04:05", det.ReleaseTime),
				IsReleaseVideo:               uint32(det.IsReleaseVideo),
				IsPlaySuccess:                uint32(det.IsPlaySuccess),
				CreateTime:                   tool.TimeToString("2006-01-02 15:04:05", det.CreateTime),
				Nickname:                     det.Nickname,
				Avatar:                       det.Avatar,
			})
		}

		re := &v1.ListCompanyTaskAccountRelationReply_CompanyTaskAccountRelation{
			Id:                    v.Id,
			CompanyTaskId:         v.CompanyTaskId,
			ProductName:           v.ProductName,
			ProductOutId:          v.ProductOutId,
			UserId:                v.UserId,
			ClaimTime:             tool.TimeToString("2006-01-02 15:04:05", v.ClaimTime),
			ExpireTime:            tool.TimeToString("2006-01-02 15:04:05", v.ExpireTime),
			Status:                uint32(v.Status),
			IsDel:                 uint32(v.IsDel),
			IsCostBuy:             uint32(v.IsCostBuy),
			ScreenshotAddress:     v.ScreenshotAddress,
			IsScreenshotAvailable: uint32(v.IsScreenshotAvailable),
			CreateTime:            tool.TimeToString("2006-01-02 15:04:05", v.CreateTime),
			UpdateTime:            tool.TimeToString("2006-01-02 15:04:05", v.UpdateTime),
			CompanyTaskDetails:    dts,
			CompanyTask: &v1.ListCompanyTaskAccountRelationReply_CompanyTaskAccountRelation_CompanyTask{
				Id:            v.CompanyTask.Id,
				ProductOutId:  v.CompanyTask.ProductOutId,
				ExpireTime:    v.CompanyTask.ExpireTime,
				PlayNum:       v.CompanyTask.PlayNum,
				Price:         v.CompanyTask.Price,
				Quota:         v.CompanyTask.Quota,
				IsTop:         uint32(v.CompanyTask.IsTop),
				IsDel:         uint32(v.CompanyTask.IsDel),
				CreateTime:    tool.TimeToString("2006-01-02 15:04:05", v.CompanyTask.CreateTime),
				IsGoodReviews: uint32(v.CompanyTask.IsGoodReviews),
				ClaimQuota:    v.CompanyTask.ClaimQuota,
				SuccessQuota:  v.CompanyTask.SuccessQuota,
			},
			CompanyProduct: &v1.ListCompanyTaskAccountRelationReply_CompanyTaskAccountRelation_CompanyProduct{
				ProductOutId: v.CompanyProduct.ProductOutId,
				ProductName:  v.CompanyProduct.ProductName,
				ProductPrice: v.CompanyProduct.ProductPrice,
				ProductImg:   v.CompanyProduct.ProductImg,
			},
			IsPlayCount: isPlayCount,
		}

		list = append(list, re)
	}

	totalPage := uint64(math.Ceil(float64(relation.Total) / float64(relation.PageSize)))

	return &v1.ListCompanyTaskAccountRelationReply{
		Code: 200,
		Data: &v1.ListCompanyTaskAccountRelationReply_Data{
			PageNum:   relation.PageNum,
			PageSize:  relation.PageSize,
			Total:     relation.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

// SyncUpdateCompanyTaskDetail
func (cs *CompanyService) SyncUpdateCompanyTaskDetail(ctx context.Context, in *emptypb.Empty) (*v1.SyncUpdateCompanyTaskDetailReply, error) {
	err := cs.ctaruc.SyncUpdateCompanyTaskDetail(ctx)

	if err != nil {
		return nil, err
	}

	return &v1.SyncUpdateCompanyTaskDetailReply{
		Code: 200,
		Data: &v1.SyncUpdateCompanyTaskDetailReply_Data{},
	}, nil
}

func (cs *CompanyService) UpdateCompanyTaskDetailScreenshot(ctx context.Context, in *v1.UpdateCompanyTaskDetailScreenshotRequest) (*v1.UpdateCompanyTaskDetailScreenshotReply, error) {
	_, err := cs.ctaruc.UpdateScreenshotById(ctx, in.Id, in.Screenshot)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateCompanyTaskDetailScreenshotReply{
		Code: 200,
		Data: &v1.UpdateCompanyTaskDetailScreenshotReply_Data{},
	}, nil
}

func (cs *CompanyService) UpdateCompanyTaskDetailScreenshotAvailable(ctx context.Context, in *v1.UpdateCompanyTaskDetailScreenshotAvailableRequest) (*v1.UpdateCompanyTaskDetailScreenshotAvailableReply, error) {
	_, err := cs.ctaruc.UpdateScreenshotAvailableById(ctx, uint8(in.IsScreenshotAvailable), in.Id)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateCompanyTaskDetailScreenshotAvailableReply{
		Code: 200,
		Data: &v1.UpdateCompanyTaskDetailScreenshotAvailableReply_Data{},
	}, nil
}

func (cs *CompanyService) ListCompanyTaskDetail(ctx context.Context, in *v1.ListCompanyTaskDetailRequest) (*v1.ListCompanyTaskDetailReply, error) {
	tasks, err := cs.ctduc.ListCompanyTaskDetail(ctx, in.PageNum, in.PageSize, in.CompanyTaskId, in.Nickname)

	if err != nil {
		return nil, err
	}

	totalPage := uint64(math.Ceil(float64(tasks.Total) / float64(tasks.PageSize)))
	list := []*v1.ListCompanyTaskDetailReply_CompanyTaskDetail{}

	for _, v := range tasks.List {
		task := &v1.ListCompanyTaskDetailReply_CompanyTaskDetail{
			Id:                           v.Id,
			CompanyTaskId:                v.CompanyTaskId,
			UserId:                       v.UserId,
			ClientKey:                    v.ClientKey,
			OpenId:                       v.OpenId,
			CompanyTaskAccountRelationId: v.CompanyTaskAccountRelationId,
			ProductName:                  v.ProductName,
			ItemId:                       v.ItemId,
			PlayCount:                    v.PlayCount,
			Cover:                        v.Cover,
			ReleaseTime:                  tool.TimeToString("2006-01-02 15:04:05", v.ReleaseTime),
			IsPlaySuccess:                uint32(v.IsPlaySuccess),
			IsReleaseVideo:               uint32(v.IsReleaseVideo),
			CreateTime:                   tool.TimeToString("2006-01-02 15:04:05", v.CreateTime),
			Nickname:                     v.Nickname,
			Avatar:                       v.Avatar,
			CompanyTaskAccountRelation: &v1.ListCompanyTaskDetailReply_CompanyTaskAccountRelation{
				Id:                    v.CompanyTaskAccountRelation.Id,
				CompanyTaskId:         v.CompanyTaskAccountRelation.CompanyTaskId,
				ProductName:           v.CompanyTaskAccountRelation.ProductName,
				ProductOutId:          v.CompanyTaskAccountRelation.ProductOutId,
				UserId:                v.CompanyTaskAccountRelation.UserId,
				ClaimTime:             tool.TimeToString("2006-01-02 15:04:05", v.CompanyTaskAccountRelation.ClaimTime),
				ExpireTime:            tool.TimeToString("2006-01-02 15:04:05", v.CompanyTaskAccountRelation.ExpireTime),
				Status:                uint32(v.CompanyTaskAccountRelation.Status),
				IsDel:                 uint32(v.CompanyTaskAccountRelation.IsDel),
				IsCostBuy:             uint32(v.CompanyTaskAccountRelation.IsCostBuy),
				ScreenshotAddress:     v.CompanyTaskAccountRelation.ScreenshotAddress,
				IsScreenshotAvailable: uint32(v.CompanyTaskAccountRelation.IsScreenshotAvailable),
				CreateTime:            tool.TimeToString("2006-01-02 15:04:05", v.CompanyTaskAccountRelation.CreateTime),
				UpdateTime:            tool.TimeToString("2006-01-02 15:04:05", v.CompanyTaskAccountRelation.UpdateTime),
			},
		}

		list = append(list, task)
	}

	return &v1.ListCompanyTaskDetailReply{
		Code: 200,
		Data: &v1.ListCompanyTaskDetailReply_Data{
			PageNum:   tasks.PageNum,
			PageSize:  tasks.PageSize,
			Total:     tasks.Total,
			TotalPage: totalPage,
			List:      list,
		},
	}, nil
}

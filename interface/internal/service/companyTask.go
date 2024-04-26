package service

import (
	"context"
	v1 "interface/api/interface/v1"
	"interface/internal/pkg/tool"
	"time"
)

func (is *InterfaceService) GetCompanyTaskByProductOutId(ctx context.Context, in *v1.GetCompanyTaskByProductOutIdRequest) (*v1.GetCompanyTaskByProductOutIdReply, error) {
	if _, err := is.verifyMiniUserLogin(ctx); err != nil {
		return nil, err
	}

	task, err := is.ctuc.GetCompanyTaskByProductOutId(ctx, in.ProductOutId, in.UserId)

	if err != nil {
		return nil, err
	}

	return &v1.GetCompanyTaskByProductOutIdReply{
		Code: 200,
		Data: &v1.GetCompanyTaskByProductOutIdReply_Data{
			Id:            task.Data.Id,
			ProductOutId:  task.Data.ProductOutId,
			ExpireTime:    task.Data.ExpireTime,
			PlayNum:       task.Data.PlayNum,
			Price:         task.Data.Price,
			Quota:         task.Data.Quota,
			ClaimQuota:    task.Data.ClaimQuota,
			SuccessQuota:  task.Data.SuccessQuota,
			IsTop:         task.Data.IsTop,
			IsDel:         task.Data.IsDel,
			CreateTime:    task.Data.CreateTime,
			IsGoodReviews: task.Data.IsGoodReviews,
			IsExist:       task.Data.IsExist,
		},
	}, nil
}

func (is *InterfaceService) ListCompanyTask(ctx context.Context, in *v1.ListCompanyTaskRequest) (*v1.ListCompanyTaskReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, ""); err != nil {
		return nil, err
	}

	res, err := is.ctuc.ListCompanyTask(ctx, 1, -1, in.PageNum, in.PageSize, in.Keyword, "")

	if err != nil {
		return nil, err
	}

	tasks := &v1.ListCompanyTaskReply{
		Code: 200,
		Data: &v1.ListCompanyTaskReply_Data{
			PageNum:   res.Data.PageNum,
			PageSize:  res.Data.PageSize,
			Total:     res.Data.Total,
			TotalPage: res.Data.TotalPage,
		},
	}

	for _, v := range res.Data.List {
		companyProduct := &v1.ListCompanyTaskReply_CompanyTask_CompanyProduct{
			ProductOutId:          v.CompanyProduct.ProductOutId,
			ProductName:           v.CompanyProduct.ProductName,
			ProductPrice:          v.CompanyProduct.ProductPrice,
			ProductImg:            v.CompanyProduct.ProductImg,
			PureCommission:        v.CompanyProduct.PureCommission,
			PureServiceCommission: v.CompanyProduct.PureServiceCommission,
			CommissionRatio:       v.CompanyProduct.CommissionRatio,
		}

		tasks.Data.List = append(tasks.Data.List, &v1.ListCompanyTaskReply_CompanyTask{
			Id:             v.Id,
			ProductOutId:   v.ProductOutId,
			ExpireTime:     v.ExpireTime,
			PlayNum:        v.PlayNum,
			Price:          v.Price,
			Quota:          v.Quota,
			ClaimQuota:     v.ClaimQuota,
			SuccessQuota:   v.SuccessQuota,
			IsTop:          v.IsTop,
			IsDel:          v.IsDel,
			CreateTime:     v.CreateTime,
			CompanyProduct: companyProduct,
			IsGoodReviews:  v.IsGoodReviews,
			ReleaseTime:    v.ReleaseTime,
		})
	}

	return tasks, nil
}

func (is *InterfaceService) ListCompanyTaskUsable(ctx context.Context, in *v1.ListCompanyTaskUsableRequest) (*v1.ListCompanyTaskReply, error) {
	res, err := is.ctuc.ListCompanyTask(ctx, 1, 0, in.PageNum, in.PageSize, "", tool.TimeToString("2006-01-02 15:04:05", time.Now()))

	if err != nil {
		return nil, err
	}

	tasks := &v1.ListCompanyTaskReply{
		Code: 200,
		Data: &v1.ListCompanyTaskReply_Data{
			PageNum:   res.Data.PageNum,
			PageSize:  res.Data.PageSize,
			Total:     res.Data.Total,
			TotalPage: res.Data.TotalPage,
		},
	}

	for _, v := range res.Data.List {
		companyProduct := &v1.ListCompanyTaskReply_CompanyTask_CompanyProduct{
			ProductOutId:          v.CompanyProduct.ProductOutId,
			ProductName:           v.CompanyProduct.ProductName,
			ProductPrice:          v.CompanyProduct.ProductPrice,
			ProductImg:            v.CompanyProduct.ProductImg,
			PureCommission:        v.CompanyProduct.PureCommission,
			PureServiceCommission: v.CompanyProduct.PureServiceCommission,
			CommissionRatio:       v.CompanyProduct.CommissionRatio,
		}

		tasks.Data.List = append(tasks.Data.List, &v1.ListCompanyTaskReply_CompanyTask{
			Id:             v.Id,
			ProductOutId:   v.ProductOutId,
			ExpireTime:     v.ExpireTime,
			PlayNum:        v.PlayNum,
			Price:          v.Price,
			Quota:          v.Quota,
			ClaimQuota:     v.ClaimQuota,
			SuccessQuota:   v.SuccessQuota,
			IsTop:          v.IsTop,
			IsDel:          v.IsDel,
			CreateTime:     v.CreateTime,
			CompanyProduct: companyProduct,
			IsGoodReviews:  v.IsGoodReviews,
			ReleaseTime:    v.ReleaseTime,
		})
	}

	return tasks, nil
}

// 达人个人任务列表
func (is *InterfaceService) ListCompanyTaskAccountRelation(ctx context.Context, in *v1.ListCompanyTaskAccountRelationRequest) (*v1.ListCompanyTaskAccountRelationReply, error) {
	if _, err := is.verifyMiniUserLogin(ctx); err != nil {
		return nil, err
	}

	task, err := is.ctuc.ListCompanyTaskAccountRelation(ctx, in.PageNum, in.PageSize, in.CompanyTaskId, in.UserId, in.Status, in.ExpireTime, in.ExpiredTime, in.ProductName)

	if err != nil {
		return nil, err
	}

	list := []*v1.ListCompanyTaskAccountRelationReply_CompanyTaskAccountRelation{}

	for _, v := range task.Data.List {
		detail := []*v1.ListCompanyTaskAccountRelationReply_CompanyTaskAccountRelation_CompanyTaskDetail{}

		for _, des := range v.CompanyTaskDetails {
			detail = append(detail, &v1.ListCompanyTaskAccountRelationReply_CompanyTaskAccountRelation_CompanyTaskDetail{
				Id:                           des.Id,
				CompanyTaskId:                des.CompanyTaskId,
				UserId:                       des.UserId,
				ClientKey:                    des.ClientKey,
				OpenId:                       des.OpenId,
				CompanyTaskAccountRelationId: des.CompanyTaskAccountRelationId,
				ProductName:                  des.ProductName,
				ItemId:                       des.ItemId,
				PlayCount:                    des.PlayCount,
				Cover:                        des.Cover,
				ReleaseTime:                  des.ReleaseTime,
				IsReleaseVideo:               des.IsReleaseVideo,
				IsPlaySuccess:                des.IsPlaySuccess,
				CreateTime:                   des.CreateTime,
				Nickname:                     des.Nickname,
				Avatar:                       des.Avatar,
			})
		}

		companyTask := &v1.ListCompanyTaskAccountRelationReply_CompanyTaskAccountRelation_CompanyTask{}

		if v.CompanyTask != nil {
			companyTask = &v1.ListCompanyTaskAccountRelationReply_CompanyTaskAccountRelation_CompanyTask{
				Id:            v.CompanyTask.Id,
				ProductOutId:  v.CompanyTask.ProductOutId,
				ExpireTime:    v.CompanyTask.ExpireTime,
				PlayNum:       v.CompanyTask.PlayNum,
				Price:         v.CompanyTask.Price,
				Quota:         v.CompanyTask.Quota,
				IsTop:         v.CompanyTask.IsTop,
				IsDel:         v.CompanyTask.IsDel,
				CreateTime:    v.CompanyTask.CreateTime,
				IsGoodReviews: v.CompanyTask.IsGoodReviews,
				ClaimQuota:    v.CompanyTask.ClaimQuota,
				SuccessQuota:  v.CompanyTask.SuccessQuota,
			}
		}

		companyProduct := &v1.ListCompanyTaskAccountRelationReply_CompanyTaskAccountRelation_CompanyProduct{}

		if v.CompanyProduct != nil {
			companyProduct = &v1.ListCompanyTaskAccountRelationReply_CompanyTaskAccountRelation_CompanyProduct{
				ProductOutId: v.CompanyProduct.ProductOutId,
				ProductName:  v.CompanyProduct.ProductName,
				ProductPrice: v.CompanyProduct.ProductPrice,
				ProductImg:   v.CompanyProduct.ProductImg,
			}
		}

		list = append(list, &v1.ListCompanyTaskAccountRelationReply_CompanyTaskAccountRelation{
			Id:                    v.Id,
			CompanyTaskId:         v.CompanyTaskId,
			ProductName:           v.ProductName,
			ProductOutId:          v.ProductOutId,
			UserId:                v.UserId,
			ClaimTime:             v.ClaimTime,
			ExpireTime:            v.ExpireTime,
			Status:                v.Status,
			IsDel:                 v.IsDel,
			IsCostBuy:             v.IsCostBuy,
			ScreenshotAddress:     v.ScreenshotAddress,
			IsScreenshotAvailable: v.IsScreenshotAvailable,
			CreateTime:            v.CreateTime,
			UpdateTime:            v.UpdateTime,
			CompanyTaskDetails:    detail,
			CompanyTask:           companyTask,
			CompanyProduct:        companyProduct,
			IsPlayCount:           v.IsPlayCount,
		})
	}

	return &v1.ListCompanyTaskAccountRelationReply{
		Code: 200,
		Data: &v1.ListCompanyTaskAccountRelationReply_Data{
			PageNum:   task.Data.PageNum,
			PageSize:  task.Data.PageSize,
			Total:     task.Data.Total,
			TotalPage: task.Data.TotalPage,
			List:      list,
		},
	}, nil
}

func (is *InterfaceService) ListCompanyTaskDetail(ctx context.Context, in *v1.ListCompanyTaskDetailRequest) (*v1.ListCompanyTaskDetailReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, ""); err != nil {
		return nil, err
	}

	taskDetails, err := is.ctuc.ListCompanyTaskDetail(ctx, in.PageNum, in.PageSize, in.CompanyTaskId, in.Nickname)

	if err != nil {
		return nil, err
	}

	list := []*v1.ListCompanyTaskDetailReply_CompanyTaskAccountRelation{}

	for _, relation := range taskDetails.Data.List {
		details := []*v1.ListCompanyTaskDetailReply_CompanyTaskDetail{}

		for _, detail := range relation.CompanyTaskDetails {
			details = append(details, &v1.ListCompanyTaskDetailReply_CompanyTaskDetail{
				Id:                           detail.Id,
				CompanyTaskId:                detail.CompanyTaskId,
				UserId:                       detail.UserId,
				ClientKey:                    detail.ClientKey,
				OpenId:                       detail.OpenId,
				CompanyTaskAccountRelationId: detail.CompanyTaskAccountRelationId,
				ProductName:                  detail.ProductName,
				ItemId:                       detail.ItemId,
				PlayCount:                    detail.PlayCount,
				Cover:                        detail.Cover,
				ReleaseTime:                  detail.ReleaseTime,
				IsReleaseVideo:               detail.IsReleaseVideo,
				IsPlaySuccess:                detail.IsPlaySuccess,
				CreateTime:                   detail.CreateTime,
				Nickname:                     detail.Nickname,
				Avatar:                       detail.Avatar,
			})
		}

		list = append(list, &v1.ListCompanyTaskDetailReply_CompanyTaskAccountRelation{
			Id:                    relation.Id,
			NickName:              relation.NickName,
			AvatarUrl:             relation.AvatarUrl,
			CompanyTaskId:         relation.CompanyTaskId,
			ProductName:           relation.ProductName,
			ProductOutId:          relation.ProductOutId,
			UserId:                relation.UserId,
			ClaimTime:             relation.ClaimTime,
			ExpireTime:            relation.ExpireTime,
			Status:                relation.Status,
			IsDel:                 relation.IsDel,
			IsCostBuy:             relation.IsCostBuy,
			ScreenshotAddress:     relation.ScreenshotAddress,
			IsScreenshotAvailable: relation.IsScreenshotAvailable,
			IsPlaySuccess:         uint32(relation.IsPlaySuccess),
			CreateTime:            relation.CreateTime,
			UpdateTime:            relation.UpdateTime,
			CompanyTaskDetails:    details,
		})
	}

	return &v1.ListCompanyTaskDetailReply{
		Code: 200,
		Data: &v1.ListCompanyTaskDetailReply_Data{
			PageNum:   taskDetails.Data.PageNum,
			PageSize:  taskDetails.Data.PageSize,
			Total:     taskDetails.Data.Total,
			TotalPage: taskDetails.Data.TotalPage,
			List:      list,
		},
	}, nil
}

func (is *InterfaceService) CreateCompanyTask(ctx context.Context, in *v1.CreateCompanyTaskRequest) (*v1.CreateCompanyTaskReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, ""); err != nil {
		return nil, err
	}

	task, err := is.ctuc.CreateCompanyTask(ctx, in.ProductOutId, in.ExpireTime, in.PlayNum, in.Quota, in.IsGoodReviews, in.Price, in.ReleaseTime)

	if err != nil {
		return nil, err
	}

	return &v1.CreateCompanyTaskReply{
		Code: 200,
		Data: &v1.CreateCompanyTaskReply_CompanyTask{
			Id:            task.Data.Id,
			ProductOutId:  task.Data.ProductOutId,
			ExpireTime:    task.Data.ExpireTime,
			PlayNum:       task.Data.PlayNum,
			Price:         task.Data.Price,
			Quota:         task.Data.Quota,
			ClaimQuota:    task.Data.ClaimQuota,
			SuccessQuota:  task.Data.SuccessQuota,
			IsTop:         task.Data.IsTop,
			IsDel:         task.Data.IsDel,
			CreateTime:    task.Data.CreateTime,
			IsGoodReviews: task.Data.IsGoodReviews,
			ReleaseTime:   task.Data.ReleaseTime,
		},
	}, nil
}

func (is *InterfaceService) CreateCompanyTaskAccountRelation(ctx context.Context, in *v1.CreateCompanyTaskAccountRelationRequest) (*v1.CreateCompanyTaskAccountRelationReply, error) {
	if _, err := is.verifyMiniUserLogin(ctx); err != nil {
		return nil, err
	}

	task, err := is.ctuc.CreateCompanyTaskAccountRelation(ctx, in.CompanyTaskId, in.ProductOutId, in.UserId)

	if err != nil {
		return nil, err
	}

	return &v1.CreateCompanyTaskAccountRelationReply{
		Code: 200,
		Data: &v1.CreateCompanyTaskAccountRelationReply_Data{
			Id:                    task.Data.Id,
			CompanyTaskId:         task.Data.CompanyTaskId,
			ProductName:           task.Data.ProductName,
			ProductOutId:          task.Data.ProductOutId,
			UserId:                task.Data.UserId,
			ClaimTime:             task.Data.ClaimTime,
			ExpireTime:            task.Data.ExpireTime,
			Status:                task.Data.Status,
			IsDel:                 task.Data.IsDel,
			IsCostBuy:             task.Data.IsCostBuy,
			ScreenshotAddress:     task.Data.ScreenshotAddress,
			IsScreenshotAvailable: task.Data.IsScreenshotAvailable,
			CreateTime:            task.Data.CreateTime,
			UpdateTime:            task.Data.UpdateTime,
		},
	}, nil
}

func (is *InterfaceService) UpdateCompanyTaskQuota(ctx context.Context, in *v1.UpdateCompanyTaskQuotaRequest) (*v1.UpdateCompanyTaskReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, ""); err != nil {
		return nil, err
	}

	task, err := is.ctuc.UpdateCompanyTaskQuota(ctx, in.CompanyTaskId, in.Quota)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateCompanyTaskReply{
		Code: 200,
		Data: &v1.UpdateCompanyTaskReply_Data{
			Id:            task.Data.Id,
			ProductOutId:  task.Data.ProductOutId,
			ExpireTime:    task.Data.ExpireTime,
			PlayNum:       task.Data.PlayNum,
			Price:         task.Data.Price,
			Quota:         task.Data.Quota,
			ClaimQuota:    task.Data.ClaimQuota,
			SuccessQuota:  task.Data.SuccessQuota,
			IsTop:         task.Data.IsTop,
			IsDel:         task.Data.IsDel,
			CreateTime:    task.Data.CreateTime,
			IsGoodReviews: task.Data.IsGoodReviews,
		},
	}, nil
}

func (is *InterfaceService) UpdateCompanyTaskIsTop(ctx context.Context, in *v1.UpdateCompanyTaskIsTopRequest) (*v1.UpdateCompanyTaskReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, ""); err != nil {
		return nil, err
	}

	task, err := is.ctuc.UpdateCompanyTaskIsTop(ctx, in.CompanyTaskId, in.IsTop)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateCompanyTaskReply{
		Code: 200,
		Data: &v1.UpdateCompanyTaskReply_Data{
			Id:            task.Data.Id,
			ProductOutId:  task.Data.ProductOutId,
			ExpireTime:    task.Data.ExpireTime,
			PlayNum:       task.Data.PlayNum,
			Price:         task.Data.Price,
			Quota:         task.Data.Quota,
			ClaimQuota:    task.Data.ClaimQuota,
			SuccessQuota:  task.Data.SuccessQuota,
			IsTop:         task.Data.IsTop,
			IsDel:         task.Data.IsDel,
			CreateTime:    task.Data.CreateTime,
			IsGoodReviews: task.Data.IsGoodReviews,
		},
	}, nil
}

func (is *InterfaceService) UpdateCompanyTaskIsDel(ctx context.Context, in *v1.UpdateCompanyTaskIsDelRequest) (*v1.UpdateCompanyTaskIsDelReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, ""); err != nil {
		return nil, err
	}

	_, err := is.ctuc.UpdateCompanyTaskIsDel(ctx, in.CompanyTaskId)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateCompanyTaskIsDelReply{
		Code: 200,
		Data: &v1.UpdateCompanyTaskIsDelReply_Data{},
	}, nil
}

func (is *InterfaceService) UpdateCompanyTaskDetailScreenshot(ctx context.Context, in *v1.UpdateCompanyTaskDetailScreenshotRequest) (*v1.UpdateCompanyTaskDetailScreenshotReply, error) {
	if _, err := is.verifyMiniUserLogin(ctx); err != nil {
		return nil, err
	}

	_, err := is.ctuc.UpdateCompanyTaskDetailScreenshot(ctx, in.Screenshot, in.Id)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateCompanyTaskDetailScreenshotReply{
		Code: 200,
		Data: &v1.UpdateCompanyTaskDetailScreenshotReply_Data{},
	}, nil
}

func (is *InterfaceService) UpdateCompanyTaskDetailScreenshotAvailable(ctx context.Context, in *v1.UpdateCompanyTaskDetailScreenshotAvailableRequest) (*v1.UpdateCompanyTaskDetailScreenshotAvailableReply, error) {
	if _, err := is.verifyLogin(ctx, true, false, ""); err != nil {
		return nil, err
	}

	_, err := is.ctuc.UpdateCompanyTaskDetailScreenshotAvailable(ctx, in.IsScreenshotAvailable, in.Id)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateCompanyTaskDetailScreenshotAvailableReply{
		Code: 200,
		Data: &v1.UpdateCompanyTaskDetailScreenshotAvailableReply_Data{},
	}, nil
}

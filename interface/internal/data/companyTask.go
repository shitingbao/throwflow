package data

import (
	"context"
	v1 "interface/api/service/company/v1"
	"interface/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type companyTaskRepo struct {
	data *Data
	log  *log.Helper
}

func NewCompanyTaskRepo(data *Data, logger log.Logger) biz.CompanyTaskRepo {
	return &companyTaskRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ctr *companyTaskRepo) CreateCompanyTask(ctx context.Context, productOutId, expireTime, playNum, price, quota uint64, isGoodReviews uint32) (*v1.CreateCompanyTaskReply, error) {
	return ctr.data.companyuc.CreateCompanyTask(ctx, &v1.CreateCompanyTaskRequest{
		ProductOutId:  productOutId,
		ExpireTime:    expireTime,
		PlayNum:       playNum,
		Price:         price,
		Quota:         quota,
		IsGoodReviews: isGoodReviews,
	})
}

func (ctr *companyTaskRepo) GetByProductOutId(ctx context.Context, productOutId, userId uint64) (*v1.GetByProductOutIdReply, error) {
	return ctr.data.companyuc.GetByProductOutId(ctx, &v1.GetByProductOutIdRequest{
		ProductOutId: productOutId,
		UserId:       userId,
	})
}

func (ctr *companyTaskRepo) UpdateCompanyTaskQuota(ctx context.Context, taskId, quota uint64) (*v1.UpdateCompanyTaskReply, error) {
	return ctr.data.companyuc.UpdateCompanyTaskQuota(ctx, &v1.UpdateCompanyTaskQuotaRequest{
		CompanyTaskId: taskId,
		Quota:         quota,
	})
}

func (ctr *companyTaskRepo) UpdateCompanyTaskIsTop(ctx context.Context, taskId uint64, isTop uint32) (*v1.UpdateCompanyTaskReply, error) {
	return ctr.data.companyuc.UpdateCompanyTaskIsTop(ctx, &v1.UpdateCompanyTaskIsTopRequest{
		CompanyTaskId: taskId,
		IsTop:         isTop,
	})
}

func (ctr *companyTaskRepo) ListCompanyTask(ctx context.Context, isTop uint32, isDel int32, pageNum, pageSize uint64, keyword string) (*v1.ListCompanyTaskReply, error) {
	return ctr.data.companyuc.ListCompanyTask(ctx, &v1.ListCompanyTaskRequest{
		PageNum:  pageNum,
		PageSize: pageSize,
		Keyword:  keyword,
		IsTop:    isTop,
		IsDel:    isDel,
	})
}

func (ctr *companyTaskRepo) UpdateCompanyTaskIsDel(ctx context.Context, taskId uint64) (*v1.UpdateCompanyTaskIsDelReply, error) {
	return ctr.data.companyuc.UpdateCompanyTaskIsDel(ctx, &v1.UpdateCompanyTaskIsDelRequest{
		CompanyTaskId: taskId,
	})
}

func (ctr *companyTaskRepo) CreateCompanyTaskAccountRelation(ctx context.Context,
	companyTaskId, productOutId, userId uint64) (*v1.CreateCompanyTaskAccountRelationReply, error) {
	return ctr.data.companyuc.CreateCompanyTaskAccountRelation(ctx, &v1.CreateCompanyTaskAccountRelationRequest{
		CompanyTaskId: companyTaskId,
		ProductOutId:  productOutId,
		UserId:        userId,
	})
}

func (ctr *companyTaskRepo) ListCompanyTaskAccountRelation(ctx context.Context,
	pageNum, pageSize, companyTaskId, userId uint64, status int32, expireTime, expiredTime, productName string) (*v1.ListCompanyTaskAccountRelationReply, error) {
	return ctr.data.companyuc.ListCompanyTaskAccountRelation(ctx, &v1.ListCompanyTaskAccountRelationRequest{
		PageNum:       pageNum,
		PageSize:      pageSize,
		CompanyTaskId: companyTaskId,
		UserId:        userId,
		Status:        status,
		ExpireTime:    expireTime,
		ExpiredTime:   expiredTime,
		ProductName:   productName,
	})
}

func (ctr *companyTaskRepo) UpdateCompanyTaskDetailScreenshotAvailable(ctx context.Context,
	isScreenshotAvailable uint32, id uint64) (*v1.UpdateCompanyTaskDetailScreenshotAvailableReply, error) {
	return ctr.data.companyuc.UpdateCompanyTaskDetailScreenshotAvailable(ctx, &v1.UpdateCompanyTaskDetailScreenshotAvailableRequest{
		Id:                    id,
		IsScreenshotAvailable: isScreenshotAvailable,
	})
}

func (ctr *companyTaskRepo) UpdateCompanyTaskDetailScreenshot(ctx context.Context,
	Screenshot string, id uint64) (*v1.UpdateCompanyTaskDetailScreenshotReply, error) {
	return ctr.data.companyuc.UpdateCompanyTaskDetailScreenshot(ctx, &v1.UpdateCompanyTaskDetailScreenshotRequest{
		Id:         id,
		Screenshot: Screenshot,
	})
}

func (ctr *companyTaskRepo) ListCompanyTaskDetail(ctx context.Context,
	pageNum, pageSize, taskId uint64, nickname string) (*v1.ListCompanyTaskDetailReply, error) {
	return ctr.data.companyuc.ListCompanyTaskDetail(ctx, &v1.ListCompanyTaskDetailRequest{
		PageNum:       pageNum,
		PageSize:      pageSize,
		CompanyTaskId: taskId,
		Nickname:      nickname,
	})
}

package biz

import (
	"context"
	v1 "interface/api/service/company/v1"
	"interface/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

type CompanyTaskRepo interface {
	CreateCompanyTask(context.Context, uint64, uint64, uint64, uint64, uint64, uint32) (*v1.CreateCompanyTaskReply, error)
	GetByProductOutId(context.Context, uint64, uint64) (*v1.GetByProductOutIdReply, error)
	UpdateCompanyTaskQuota(context.Context, uint64, uint64) (*v1.UpdateCompanyTaskReply, error)
	UpdateCompanyTaskIsTop(context.Context, uint64, uint32) (*v1.UpdateCompanyTaskReply, error)
	UpdateCompanyTaskIsDel(context.Context, uint64) (*v1.UpdateCompanyTaskIsDelReply, error)
	ListCompanyTask(context.Context, uint32, int32, uint64, uint64, string) (*v1.ListCompanyTaskReply, error)
	CreateCompanyTaskAccountRelation(context.Context, uint64, uint64, uint64) (*v1.CreateCompanyTaskAccountRelationReply, error)
	ListCompanyTaskAccountRelation(context.Context, uint64, uint64, uint64, uint64, int32, string, string, string) (*v1.ListCompanyTaskAccountRelationReply, error)
	UpdateCompanyTaskDetailScreenshotAvailable(context.Context, uint32, uint64) (*v1.UpdateCompanyTaskDetailScreenshotAvailableReply, error)
	UpdateCompanyTaskDetailScreenshot(context.Context, string, uint64) (*v1.UpdateCompanyTaskDetailScreenshotReply, error)
	ListCompanyTaskDetail(context.Context, uint64, uint64, uint64, string) (*v1.ListCompanyTaskDetailReply, error)
}

type CompanyTaskUsecase struct {
	repo CompanyTaskRepo
	conf *conf.Data
	log  *log.Helper
}

func NewCompanyTaskUsecase(repo CompanyTaskRepo, conf *conf.Data, logger log.Logger) *CompanyTaskUsecase {
	return &CompanyTaskUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (cuc *CompanyTaskUsecase) CreateCompanyTask(ctx context.Context, productId, expireTime, playNum, price, quota uint64, isGoodReviews uint32) (*v1.CreateCompanyTaskReply, error) {
	return cuc.repo.CreateCompanyTask(ctx, productId, expireTime, playNum, price, quota, isGoodReviews)
}

func (cuc *CompanyTaskUsecase) GetByProductOutId(ctx context.Context, productOutId, userId uint64) (*v1.GetByProductOutIdReply, error) {
	return cuc.repo.GetByProductOutId(ctx, productOutId, userId)
}

func (cuc *CompanyTaskUsecase) UpdateCompanyTaskQuota(ctx context.Context, taskId, quota uint64) (*v1.UpdateCompanyTaskReply, error) {
	return cuc.repo.UpdateCompanyTaskQuota(ctx, taskId, quota)
}

func (cuc *CompanyTaskUsecase) UpdateCompanyTaskIsTop(ctx context.Context, taskId uint64, isTop uint32) (*v1.UpdateCompanyTaskReply, error) {
	return cuc.repo.UpdateCompanyTaskIsTop(ctx, taskId, isTop)
}

func (cuc *CompanyTaskUsecase) ListCompanyTask(ctx context.Context, isTop uint32, isDel int32, pageNum, pageSize uint64, keyword string) (*v1.ListCompanyTaskReply, error) {
	return cuc.repo.ListCompanyTask(ctx, isTop, isDel, pageNum, pageSize, keyword)
}

func (cuc *CompanyTaskUsecase) UpdateCompanyTaskIsDel(ctx context.Context, taskId uint64) (*v1.UpdateCompanyTaskIsDelReply, error) {
	return cuc.repo.UpdateCompanyTaskIsDel(ctx, taskId)
}

func (cuc *CompanyTaskUsecase) CreateCompanyTaskAccountRelation(ctx context.Context,
	companyTaskId, productId, userId uint64) (*v1.CreateCompanyTaskAccountRelationReply, error) {
	return cuc.repo.CreateCompanyTaskAccountRelation(ctx, companyTaskId, productId, userId)
}

func (cuc *CompanyTaskUsecase) ListCompanyTaskAccountRelation(ctx context.Context, pageNum uint64, pageSize uint64, companyTaskId uint64, userId uint64, status int32, expireTime string, expiredTime, productName string) (*v1.ListCompanyTaskAccountRelationReply, error) {
	return cuc.repo.ListCompanyTaskAccountRelation(ctx, pageNum, pageSize, companyTaskId, userId, status, expireTime, expiredTime, productName)
}

func (cuc *CompanyTaskUsecase) UpdateCompanyTaskDetailScreenshotAvailable(ctx context.Context,
	isScreenshotAvailable uint32, id uint64) (*v1.UpdateCompanyTaskDetailScreenshotAvailableReply, error) {
	return cuc.repo.UpdateCompanyTaskDetailScreenshotAvailable(ctx, isScreenshotAvailable, id)
}

func (cuc *CompanyTaskUsecase) UpdateCompanyTaskDetailScreenshot(ctx context.Context, screenshot string, id uint64) (*v1.UpdateCompanyTaskDetailScreenshotReply, error) {
	return cuc.repo.UpdateCompanyTaskDetailScreenshot(ctx, screenshot, id)
}

func (cuc *CompanyTaskUsecase) ListCompanyTaskDetail(ctx context.Context,
	pageNum, pageSize, taskId uint64, nickname string) (*v1.ListCompanyTaskDetailReply, error) {
	return cuc.repo.ListCompanyTaskDetail(ctx, pageNum, pageSize, taskId, nickname)
}

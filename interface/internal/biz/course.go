package biz

import (
	"context"
	v1 "interface/api/service/weixin/v1"
	"interface/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

type CourseRepo interface {
	ListCourseRoot(context.Context, uint64, uint64) (*v1.ListCourseRootReply, error)
	ListCourseVideo(context.Context, uint64, uint64, uint64) (*v1.ListCourseVideoReply, error)
	UpdateCourseUser(context.Context, uint64, uint64, uint64) (*v1.UpdateCourseUserReply, error)
}

type CourseUsecase struct {
	repo CourseRepo
	conf *conf.Data
	log  *log.Helper
}

func NewCourseUsecase(repo CourseRepo, conf *conf.Data, logger log.Logger) *CourseUsecase {
	return &CourseUsecase{repo: repo, conf: conf, log: log.NewHelper(logger)}
}

func (cuc *CourseUsecase) ListCourseRoot(ctx context.Context, userId uint64, organizationId uint64) (*v1.ListCourseRootReply, error) {
	return cuc.repo.ListCourseRoot(ctx, userId, organizationId)
}

func (cuc *CourseUsecase) ListCourseVideo(ctx context.Context, userId, organizationId, courseId uint64) (*v1.ListCourseVideoReply, error) {
	return cuc.repo.ListCourseVideo(ctx, userId, organizationId, courseId)
}

func (cuc *CourseUsecase) UpdateCourseUser(ctx context.Context, userId, courseId, duration uint64) (*v1.UpdateCourseUserReply, error) {
	return cuc.repo.UpdateCourseUser(ctx, userId, courseId, duration)
}

package data

import (
	"context"
	v1 "interface/api/service/weixin/v1"
	"interface/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type courseRepo struct {
	data *Data
	log  *log.Helper
}

func NewCourseRepo(data *Data, logger log.Logger) biz.CourseRepo {
	return &courseRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ctr *courseRepo) ListCourseRoot(ctx context.Context, userId, organizationId uint64) (*v1.ListCourseRootReply, error) {
	return ctr.data.weixinuc.ListCourseRoot(ctx, &v1.ListCourseRootRequest{
		UserId:         userId,
		OrganizationId: organizationId,
	})
}

func (ctr *courseRepo) ListCourseVideo(ctx context.Context, userId, organizationId, courseId uint64) (*v1.ListCourseVideoReply, error) {
	return ctr.data.weixinuc.ListCourseVideo(ctx, &v1.ListCourseVideoRequest{
		UserId:         userId,
		OrganizationId: organizationId,
		CourseId:       courseId,
	})
}

func (ctr *courseRepo) UpdateCourseUser(ctx context.Context, userId, courseId, duration uint64) (*v1.UpdateCourseUserReply, error) {
	return ctr.data.weixinuc.UpdateCourseUser(ctx, &v1.UpdateCourseUserRequest{
		UserId:   userId,
		CourseId: courseId,
		Duration: duration,
	})
}

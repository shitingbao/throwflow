package biz

import (
	"context"
	"weixin/internal/domain"

	"github.com/go-kratos/kratos/v2/errors"
)

var (
	WeixinUserCourseError = errors.NotFound("WEIXIN_USER_COURSE_ERROR", "保存用户课程关系出错")
)

type CourseUserRepo interface {
	Get(context.Context, uint8, uint64, uint64) (*domain.CourseUser, error)
	Count(context.Context, uint8, []uint64) ([]*domain.CourseUserRelation, error)
	Statistics(context.Context, uint64, []uint64) (*domain.CourseUserRelation, error)
	Save(context.Context, *domain.CourseUser) (*domain.CourseUser, error)
}

// ListCourseVideo 课程内容，并记录用户和课程关系
func (cu *CourseUsecase) UpdateCourseUser(ctx context.Context, userId, courseId, duration uint64) (*domain.CourseUser, error) {
	course, err := cu.repo.Get(ctx, courseId)

	if err != nil {
		return nil, WeixinListCourseError
	}

	if duration > course.Duration {
		return nil, WeixinUserCourseError
	}

	courseUser, err := cu.curepo.Get(ctx, domain.CourseUserTypeVideo, courseId, userId)

	if err != nil {
		// create
		courseUser = domain.NewCourseUser(ctx, courseId, userId, duration, domain.CourseUserTypeVideo)
		courseUser.SetCreateTime(ctx)
	} else {
		// udpate,传递的是当前观看时长，存储的是历史最长播放时长
		if duration <= courseUser.Duration {
			return courseUser, nil
		}

		courseUser.SetDuration(ctx, duration)
	}

	courseUser.SetUpdateTime(ctx)

	courseUser, err = cu.curepo.Save(ctx, courseUser)

	if err != nil {
		return nil, WeixinUserCourseError
	}

	return courseUser, nil
}

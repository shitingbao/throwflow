package biz

import (
	"context"
	"strconv"
	"strings"
	"weixin/internal/conf"
	"weixin/internal/domain"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	WeixinListCourseError      = errors.NotFound("WEIXIN_LIST_COURSE_ERROR", "课程列表获取出错")
	WeixinListCourseCountError = errors.NotFound("WEIXIN_LIST_COURSE_COUNT_ERROR", "课程列表人数获取出错")
	WeixinListCourseSumError   = errors.NotFound("WEIXIN_LIST_COURSE_SUM_ERROR", "课程列表学习进度获取出错")
)

type CourseRepo interface {
	Get(context.Context, uint64) (*domain.Course, error)
	GetIds(context.Context, uint64, uint64) ([]uint64, error)
	List(context.Context, uint8, uint64, uint64, uint64) ([]*domain.Course, error)
	ListRoot(context.Context, uint8) ([]*domain.Course, error)
}

type CourseUsecase struct {
	repo   CourseRepo
	curepo CourseUserRepo
	uorepo UserOrganizationRelationRepo
	conf   *conf.Data
	log    *log.Helper
}

func NewCourseUsecase(repo CourseRepo, curepo CourseUserRepo, uorepo UserOrganizationRelationRepo, conf *conf.Data, logger log.Logger) *CourseUsecase {
	return &CourseUsecase{repo: repo, curepo: curepo, uorepo: uorepo, conf: conf, log: log.NewHelper(logger)}
}

// ListCourseRoot 课程大类封面
func (cu *CourseUsecase) ListCourseRoot(ctx context.Context, userId, organizationId uint64) ([]*domain.Course, error) {
	courses, err := cu.repo.ListRoot(ctx, domain.AdvancedLevel)

	if err != nil {
		return nil, WeixinListCourseError
	}

	ids := []uint64{}

	for _, course := range courses {
		ids = append(ids, course.Id)
	}

	relationCounts, err := cu.curepo.Count(ctx, domain.CourseUserType, ids)

	if err != nil {
		return nil, WeixinListCourseCountError
	}

	relationMap := make(map[uint64]int64)

	for _, relation := range relationCounts {
		relationMap[relation.CourseId] = relation.Count
	}

	for _, course := range courses {
		ids, err := cu.repo.GetIds(ctx, course.CourseSort, course.CourseSort+99999)

		if err != nil {
			return nil, WeixinListCourseSumError
		}

		relationDuration, err := cu.curepo.Statistics(ctx, userId, ids)

		if err != nil {
			return nil, WeixinListCourseSumError
		}

		course.SetTotalDuration(ctx, course.Duration)
		course.SetStudyCount(ctx, relationMap[course.Id])
		course.SetStudyDuration(ctx, relationDuration.Duration)
	}

	return courses, nil
}

// ListCourseVideo 课程内容，并记录用户和课程关系
func (cu *CourseUsecase) ListCourseVideo(ctx context.Context, userId, organizationId, courseId uint64) ([]*domain.Course, error) {
	course, err := cu.repo.Get(ctx, courseId)

	if err != nil {
		return nil, WeixinListCourseError
	}

	list, err := cu.repo.List(ctx, domain.AdvancedLevel, userId, course.CourseSort, course.CourseSort+99999)

	if err != nil {
		return nil, WeixinListCourseError
	}

	courses := []*domain.Course{}
	courseVideos := []*domain.Course{}

	for _, v := range list {
		if v.CourseType == domain.CouseTypeContent {
			courses = append(courses, v)
		} else {
			courseVideos = append(courseVideos, v)
		}
	}

	for _, course := range courseVideos {
		course.SetUserDuration(ctx, course.CourseUser.Duration)
		// 找出下级，并设置时长该视频用户观看的历史时长
		for _, cour := range courses {
			if strings.HasPrefix(strconv.Itoa(int(course.CourseSort)), strconv.Itoa(int(cour.CourseSort/100))) {
				cour.Children = append(cour.Children, *course)
			}
		}
	}

	// 考虑没有目录的情况，直接反馈视频
	if len(courses) == 0 {
		courses = courseVideos
	}

	_, err = cu.curepo.Get(ctx, domain.CourseUserType, courseId, userId)

	if err == nil {
		return courses, nil
	}

	// 保存该用户和该课程的关系
	courseUser := domain.NewCourseUser(ctx, courseId, userId, 0, domain.CourseUserType)
	courseUser.SetCreateTime(ctx)
	courseUser.SetUpdateTime(ctx)

	if _, err := cu.curepo.Save(ctx, courseUser); err != nil {
		return nil, WeixinUserCourseError
	}

	return courses, nil
}

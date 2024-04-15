package service

import (
	"context"
	v1 "weixin/api/weixin/v1"
	"weixin/internal/pkg/tool"
)

func (ws *WeixinService) ListCourseRoot(ctx context.Context, in *v1.ListCourseRootRequest) (*v1.ListCourseRootReply, error) {
	courses, err := ws.cuc.ListCourseRoot(ctx, in.UserId, in.OrganizationId)

	if err != nil {
		return nil, err
	}

	list := []*v1.ListCourseRootReply_Course{}

	for _, course := range courses {
		list = append(list, &v1.ListCourseRootReply_Course{
			Id:            course.Id,
			CourseName:    course.CourseName,
			Imgurl:        course.Imgurl,
			Level:         uint32(course.Level),
			CourseSort:    course.CourseSort,
			StudyCount:    uint64(course.StudyCount),
			StudyDuration: course.StudyDuration,
			TotalDuration: course.TotalDuration,
			CreateTime:    tool.TimeToString("2006-01-02 15:04:05", course.CreateTime),
			UpdateTime:    tool.TimeToString("2006-01-02 15:04:05", course.UpdateTime),
		})
	}

	return &v1.ListCourseRootReply{
		Code: 200,
		Data: &v1.ListCourseRootReply_Data{
			List: list,
		},
	}, nil
}

func (ws *WeixinService) ListCourseVideo(ctx context.Context, in *v1.ListCourseVideoRequest) (*v1.ListCourseVideoReply, error) {
	courseVideos, err := ws.cuc.ListCourseVideo(ctx, in.UserId, in.OrganizationId, in.CourseId)

	if err != nil {
		return nil, err
	}

	list := []*v1.ListCourseVideoReply_Course{}

	for _, courseVideo := range courseVideos {
		children := []*v1.ListCourseVideoReply_Children{}

		for _, child := range courseVideo.Children {
			children = append(children, &v1.ListCourseVideoReply_Children{
				Id:           child.Id,
				CourseName:   child.CourseName,
				Imgurl:       child.Imgurl,
				Level:        uint32(child.Level),
				VideoUrl:     child.VideoUrl,
				CourseSort:   child.CourseSort,
				Duration:     child.Duration,
				UserDuration: child.UserDuration,
				CreateTime:   tool.TimeToString("2006-01-02 15:04:05", child.CreateTime),
				UpdateTime:   tool.TimeToString("2006-01-02 15:04:05", child.UpdateTime),
			})
		}

		list = append(list, &v1.ListCourseVideoReply_Course{
			Id:           courseVideo.Id,
			CourseName:   courseVideo.CourseName,
			Imgurl:       courseVideo.Imgurl,
			Level:        uint32(courseVideo.Level),
			VideoUrl:     courseVideo.VideoUrl,
			CourseSort:   courseVideo.CourseSort,
			Duration:     courseVideo.Duration,
			UserDuration: courseVideo.UserDuration,
			CreateTime:   tool.TimeToString("2006-01-02 15:04:05", courseVideo.CreateTime),
			UpdateTime:   tool.TimeToString("2006-01-02 15:04:05", courseVideo.UpdateTime),
			Children:     children,
		})
	}

	return &v1.ListCourseVideoReply{
		Code: 200,
		Data: &v1.ListCourseVideoReply_Data{
			List: list,
		},
	}, nil
}

func (ws *WeixinService) UpdateCourseUser(ctx context.Context, in *v1.UpdateCourseUserRequest) (*v1.UpdateCourseUserReply, error) {
	_, err := ws.cuc.UpdateCourseUser(ctx, in.UserId, in.CourseId, in.Duration)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateCourseUserReply{
		Code: 200,
		Data: &v1.UpdateCourseUserReply_Data{},
	}, nil
}

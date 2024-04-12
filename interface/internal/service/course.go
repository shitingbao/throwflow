package service

import (
	"context"
	v1 "interface/api/interface/v1"
)

func (is *InterfaceService) ListCourseRoot(ctx context.Context, in *v1.ListCourseRootRequest) (*v1.ListCourseRootReply, error) {
	res, err := is.couuc.ListCourseRoot(ctx, in.UserId, in.OrganizationId)

	if err != nil {
		return nil, err
	}

	list := []*v1.ListCourseRootReply_Course{}

	for _, v := range res.Data.List {
		list = append(list, &v1.ListCourseRootReply_Course{
			Id:            v.Id,
			CourseName:    v.CourseName,
			Imgurl:        v.Imgurl,
			Level:         v.Level,
			CourseSort:    v.CourseSort,
			StudyCount:    v.StudyCount,
			StudyDuration: v.StudyDuration,
			TotalDuration: v.TotalDuration,
			CreateTime:    v.CreateTime,
			UpdateTime:    v.UpdateTime,
		})
	}

	return &v1.ListCourseRootReply{
		Code: 200,
		Data: &v1.ListCourseRootReply_Data{
			List: list,
		},
	}, nil
}

func (is *InterfaceService) ListCourseVideo(ctx context.Context, in *v1.ListCourseVideoRequest) (*v1.ListCourseVideoReply, error) {
	res, err := is.couuc.ListCourseVideo(ctx, in.UserId, in.OrganizationId, in.CourseId)

	if err != nil {
		return nil, err
	}

	list := []*v1.ListCourseVideoReply_Course{}

	for _, v := range res.Data.List {

		children := []*v1.ListCourseVideoReply_Children{}

		for _, child := range v.Children {
			children = append(children, &v1.ListCourseVideoReply_Children{
				Id:           child.Id,
				CourseName:   child.CourseName,
				Imgurl:       child.Imgurl,
				Level:        uint32(child.Level),
				VideoUrl:     child.VideoUrl,
				CourseSort:   child.CourseSort,
				Duration:     child.Duration,
				UserDuration: child.UserDuration,
				CreateTime:   child.CreateTime,
				UpdateTime:   child.UpdateTime,
			})
		}

		list = append(list, &v1.ListCourseVideoReply_Course{
			Id:           v.Id,
			CourseName:   v.CourseName,
			Imgurl:       v.Imgurl,
			Level:        v.Level,
			VideoUrl:     v.VideoUrl,
			CourseSort:   v.CourseSort,
			Duration:     v.Duration,
			UserDuration: v.UserDuration,
			CreateTime:   v.CreateTime,
			UpdateTime:   v.UpdateTime,
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

func (is *InterfaceService) UpdateCourseUser(ctx context.Context, in *v1.UpdateCourseUserRequest) (*v1.UpdateCourseUserReply, error) {
	_, err := is.couuc.UpdateCourseUser(ctx, in.UserId, in.CourseId, in.Duration)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateCourseUserReply{
		Code: 200,
		Data: &v1.UpdateCourseUserReply_Data{},
	}, nil
}

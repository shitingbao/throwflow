package domain

import (
	"context"
	"time"
)

const (
	CouseTypeContent = iota + 1 // 课程目录
	CouseTypeVideo              // 课程视频

)

type Course struct {
	Id            uint64
	ParentId      uint64
	CourseType    uint8
	CourseName    string
	Imgurl        string
	Level         uint8
	VideoUrl      string
	CourseSort    uint64
	Duration      uint64
	StudyCount    int64
	StudyDuration uint64
	TotalDuration uint64
	CreateTime    time.Time
	UpdateTime    time.Time
	UserDuration  uint64
	CourseUser    CourseUser
	Children      []Course
}

func (c *Course) SetParentId(ctx context.Context, parentId uint64) {
	c.ParentId = parentId
}

func (c *Course) SetCourseType(ctx context.Context, courseType uint8) {
	c.CourseType = courseType
}

func (c *Course) SetCourseName(ctx context.Context, courseName string) {
	c.CourseName = courseName
}

func (c *Course) SetImgurl(ctx context.Context, imgurl string) {
	c.Imgurl = imgurl
}

func (c *Course) SetLevel(ctx context.Context, level uint8) {
	c.Level = level
}

func (c *Course) SetVideoUrl(ctx context.Context, videoUrl string) {
	c.VideoUrl = videoUrl
}

func (c *Course) SetCourseSort(ctx context.Context, courseSort uint64) {
	c.CourseSort = courseSort
}

func (c *Course) SetDuration(ctx context.Context, duration uint64) {
	c.Duration = duration
}

func (c *Course) SetStudyCount(ctx context.Context, studyCount int64) {
	c.StudyCount = studyCount
}

func (c *Course) SetStudyDuration(ctx context.Context, studyDuration uint64) {
	c.StudyDuration = studyDuration
}

func (c *Course) SetTotalDuration(ctx context.Context, totalDuration uint64) {
	c.TotalDuration = totalDuration
}

func (c *Course) SetCreateTime(ctx context.Context) {
	c.CreateTime = time.Now()
}

func (c *Course) SetUpdateTime(ctx context.Context) {
	c.UpdateTime = time.Now()
}

func (c *Course) SetUserDuration(ctx context.Context, userDuration uint64) {
	c.UserDuration = userDuration
}

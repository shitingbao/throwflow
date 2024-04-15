package domain

import (
	"context"
	"time"
)

const (
	CourseUserType      = iota // 用户课程关系,进入过课程
	CourseUserTypeVideo        // 用户课程视频观看时长关系
)

type CourseUser struct {
	Id             uint64
	CourseId       uint64
	UserId         uint64
	Duration       uint64
	CourseUserType uint8
	CreateTime     time.Time
	UpdateTime     time.Time
}

func NewCourseUser(ctx context.Context, courseId, userId, duration uint64, courseUserType uint8) *CourseUser {
	return &CourseUser{
		CourseId:       courseId,
		UserId:         userId,
		Duration:       duration,
		CourseUserType: courseUserType,
	}
}

func (c *CourseUser) SetCourseId(ctx context.Context, courseId uint64) {
	c.CourseId = courseId
}

func (c *CourseUser) SetUserId(ctx context.Context, userId uint64) {
	c.UserId = userId
}

func (c *CourseUser) SetDuration(ctx context.Context, duration uint64) {
	c.Duration = duration
}

func (c *CourseUser) SetCourseUserType(ctx context.Context, courseUserType uint8) {
	c.CourseUserType = courseUserType
}

func (c *CourseUser) SetCreateTime(ctx context.Context) {
	c.CreateTime = time.Now()
}

func (c *CourseUser) SetUpdateTime(ctx context.Context) {
	c.UpdateTime = time.Now()
}

type CourseUserRelation struct {
	UserId   uint64
	CourseId uint64
	Count    int64
	Duration uint64
}

func NewCourseUserRelation(ctx context.Context, userId, courseId uint64) *CourseUserRelation {
	return &CourseUserRelation{
		UserId:   userId,
		CourseId: courseId,
	}
}

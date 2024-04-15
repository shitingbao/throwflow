package data

import (
	"context"
	"time"
	"weixin/internal/biz"
	"weixin/internal/domain"

	"github.com/go-kratos/kratos/v2/log"
)

// 课程视频表
type Course struct {
	Id         uint64     `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	ParentId   uint64     `gorm:"column:parent_id;type:bigint(20) UNSIGNED;comment:上级ID"`
	CourseType uint8      `gorm:"column:course_type;type:tinyint(3) UNSIGNED;not null;default:1;comment:课程类别:1目录,2视频"`
	CourseName string     `gorm:"column:course_name;type:varchar(255);not null;comment:视频大类名称"`
	Imgurl     string     `gorm:"column:imgurl;type:varchar(255);comment:视频大类封面"`
	Level      uint8      `gorm:"column:level;type:tinyint(3) UNSIGNED;not null;default:0;comment:等级"`
	VideoUrl   string     `gorm:"column:video_url;type:varchar(255);comment:url"`
	CourseSort uint64     `gorm:"column:course_sort;type:int(10);comment:三级目录,100000,101000,101010"` // 对应目录排序
	Duration   uint64     `gorm:"column:duration;type:bigint(20) UNSIGNED;not null;default:0;comment:对应视频总时长(秒)"`
	CreateTime time.Time  `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime time.Time  `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
	CourseUser CourseUser `gorm:"foreignKey:CourseId;references:Id"`
}

func (Course) TableName() string {
	return "weixin_course"
}

type courseRepo struct {
	data *Data
	log  *log.Helper
}

func (c *Course) ToDomain(ctx context.Context) *domain.Course {
	course := &domain.Course{
		Id:         c.Id,
		ParentId:   c.ParentId,
		CourseType: c.CourseType,
		CourseName: c.CourseName,
		Imgurl:     c.Imgurl,
		Level:      c.Level,
		VideoUrl:   c.VideoUrl,
		CourseSort: c.CourseSort,
		Duration:   c.Duration,
		CreateTime: c.CreateTime,
		UpdateTime: c.UpdateTime,
		CourseUser: domain.CourseUser{
			Id:             c.CourseUser.Id,
			CourseId:       c.CourseUser.CourseId,
			UserId:         c.CourseUser.UserId,
			Duration:       c.CourseUser.Duration,
			CourseUserType: c.CourseUser.CourseUserType,
			CreateTime:     c.CourseUser.CreateTime,
			UpdateTime:     c.CourseUser.UpdateTime,
		},
	}

	return course
}

func NewCourseRepo(data *Data, logger log.Logger) biz.CourseRepo {
	return &courseRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cr *courseRepo) Get(ctx context.Context, id uint64) (*domain.Course, error) {
	course := &Course{}

	if err := cr.data.db.WithContext(ctx).Model(&Course{}).First(course, id).Error; err != nil {
		return nil, err
	}

	return course.ToDomain(ctx), nil
}

func (cr *courseRepo) GetIds(ctx context.Context, minSort, maxSort uint64) ([]uint64, error) {
	ids := []uint64{}

	if err := cr.data.db.WithContext(ctx).Model(&Course{}).Select("id").
		Where("course_type = 2 and course_sort >= ? and course_sort <= ?", minSort, maxSort).
		Find(&ids).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

func (cr *courseRepo) List(ctx context.Context, level uint8, userId, minSort, maxSort uint64) ([]*domain.Course, error) {
	var courses []Course
	list := make([]*domain.Course, 0)

	db := cr.data.db.WithContext(ctx).Model(&Course{}).
		Where("parent_id is not null and level <= ? and course_sort >= ? and course_sort <= ?", level, minSort, maxSort).
		Preload("CourseUser", "course_user_type = 1 and user_id = ?", userId)

	if result := db.Order("course_sort").Find(&courses); result.Error != nil {
		return nil, result.Error
	}

	for _, course := range courses {
		list = append(list, course.ToDomain(ctx))
	}

	return list, nil
}

func (cr *courseRepo) ListRoot(ctx context.Context, level uint8) ([]*domain.Course, error) {
	var courses []Course
	list := make([]*domain.Course, 0)

	db := cr.data.db.WithContext(ctx).
		Model(&Course{}).Where("parent_id is null and level <= ?", level)

	if result := db.Order("course_sort").Find(&courses); result.Error != nil {
		return nil, result.Error
	}

	for _, course := range courses {
		list = append(list, course.ToDomain(ctx))
	}

	return list, nil
}

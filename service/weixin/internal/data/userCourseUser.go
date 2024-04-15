package data

import (
	"context"
	"time"
	"weixin/internal/biz"
	"weixin/internal/domain"

	"github.com/go-kratos/kratos/v2/log"
)

// 微信用户查看课程表
// 保存该用户是否看过该类课程,或者对应观看时间
type CourseUser struct {
	Id             uint64    `gorm:"column:id;primarykey;type:bigint(20) UNSIGNED;autoIncrement;not null;comment:自增ID"`
	CourseId       uint64    `gorm:"column:course_id;type:bigint(20) UNSIGNED;uniqueIndex:idx_unique_course_id_user_id;not null;comment:课程ID"`
	UserId         uint64    `gorm:"column:user_id;type:bigint(20) UNSIGNED;uniqueIndex:idx_unique_course_id_user_id;not null;comment:用户ID"`
	Duration       uint64    `gorm:"column:duration;type:bigint(20) UNSIGNED;not null;default:0;comment:观看对应视频总时长(秒)"`
	CourseUserType uint8     `gorm:"column:course_user_type;type:tinyint(3) UNSIGNED;uniqueIndex:idx_unique_course_id_user_id;not null;default:0;comment:课程用户关系:0为进入课程,1为观看有时间"`
	CreateTime     time.Time `gorm:"column:create_time;type:datetime;not null;comment:新增时间"`
	UpdateTime     time.Time `gorm:"column:update_time;type:datetime;not null;comment:修改时间"`
}

func (CourseUser) TableName() string {
	return "weixin_user_course"
}

type courseUserRepo struct {
	data *Data
	log  *log.Helper
}

func (cu *CourseUser) ToDomain(ctx context.Context) *domain.CourseUser {
	courseUser := &domain.CourseUser{
		Id:             cu.Id,
		CourseId:       cu.CourseId,
		UserId:         cu.UserId,
		Duration:       cu.Duration,
		CourseUserType: cu.CourseUserType,
		CreateTime:     cu.CreateTime,
		UpdateTime:     cu.UpdateTime,
	}

	return courseUser
}

func NewCourseUserRepo(data *Data, logger log.Logger) biz.CourseUserRepo {
	return &courseUserRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (cur *courseUserRepo) Get(ctx context.Context, courseUserType uint8, courseId, userId uint64) (*domain.CourseUser, error) {
	courseUser := &CourseUser{}
	db := cur.data.db.WithContext(ctx).Model(&CourseUser{}).
		Where("course_id = ? and user_id = ? and course_user_type = ?", courseId, userId, courseUserType)

	if err := db.First(courseUser).Error; err != nil {
		return nil, err
	}

	return courseUser.ToDomain(ctx), nil
}

func (cur *courseUserRepo) Count(ctx context.Context, courseUserType uint8, courseIds []uint64) ([]*domain.CourseUserRelation, error) {
	relations := []*domain.CourseUserRelation{}

	db := cur.data.db.WithContext(ctx).Model(&CourseUser{}).
		Where("course_user_type = ? and course_id in (?)", courseUserType, courseIds).
		Group("course_id").
		Select("course_id,count(*) as count")

	if err := db.
		Find(&relations).Error; err != nil {
		return nil, err
	}

	return relations, nil
}

// 注意是每个视频的时间相加，统计的是整个类目的观看时间
func (cur *courseUserRepo) Statistics(ctx context.Context, userId uint64, courseIds []uint64) (*domain.CourseUserRelation, error) {
	course := &domain.CourseUserRelation{}

	db := cur.data.db.WithContext(ctx).Model(&CourseUser{}).
		Select("user_id,sum(duration) as duration").
		Where("course_user_type = 1 and user_id = ? and course_id in (?)", userId, courseIds).
		Group("user_id")

	if result := db.Find(course); result.Error != nil {
		return nil, result.Error
	}

	return course, nil
}

func (cur *courseUserRepo) Save(ctx context.Context, in *domain.CourseUser) (*domain.CourseUser, error) {
	course := &CourseUser{
		Id:             in.Id,
		CourseId:       in.CourseId,
		UserId:         in.UserId,
		Duration:       in.Duration,
		CourseUserType: in.CourseUserType,
		CreateTime:     in.CreateTime,
		UpdateTime:     in.UpdateTime,
	}

	if err := cur.data.DB(ctx).Save(course).Error; err != nil {
		return nil, err
	}

	return course.ToDomain(ctx), nil
}

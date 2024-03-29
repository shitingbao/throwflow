package data

import (
	"context"
	"douyin/internal/biz"
	"douyin/internal/domain"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// 计划任务日志表
type TaskLog struct {
	TaskType   string    `json:"task_type" bson:"task_type"`
	Content    string    `json:"content" bson:"content"`
	CreateTime time.Time `json:"create_time" bson:"create_time"`
}

type taskLogRepo struct {
	data *Data
	log  *log.Helper
}

func NewTaskLogRepo(data *Data, logger log.Logger) biz.TaskLogRepo {
	return &taskLogRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
func (tlr *taskLogRepo) Save(ctx context.Context, in *domain.TaskLog) error {
	collection := tlr.data.mdb.Database(tlr.data.conf.Mongo.Dbname).Collection("task_log")

	if _, err := collection.InsertOne(ctx, &TaskLog{
		TaskType:   in.TaskType,
		Content:    in.Content,
		CreateTime: in.CreateTime,
	}); err != nil {
		return err
	}

	return nil
}

package server

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron/v3"
	"task/internal/conf"
	"task/internal/service"
)

type Task struct {
	c    *conf.Task
	cron *cron.Cron
}

func NewTaskServer(c *conf.Task, taskService *service.TaskService) (t *Task) {
	taskService.Init()
	
	t = &Task{
		c:    c,
		cron: cron.New(),
	}

	for _, ljob := range c.Jobs {
		job, ok := service.DefaultJobs[ljob.Name]

		if !ok {
			log.Warnf("can not find job: %s", ljob.Name)

			continue
		}

		t.cron.AddFunc(ljob.Schedule, job)
	}

	return t
}

func (t *Task) Start(ctx context.Context) error {
	fmt.Println("计划任务开始了")
	t.cron.Start()

	return nil
}

func (t *Task) Stop(ctx context.Context) error {
	fmt.Println("计划任务结束了")
	t.cron.Stop()

	return nil
}

func (t *Task) RunSrv(name string) {
	log.Info("run job{%s}", name)
}

// HeartBeat .
func (t *Task) HeartBeat() {
	log.Info("alive...")
}

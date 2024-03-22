package service

import (
	"context"
	"log"
	"time"
)

func (ts *TaskService) SyncCompanyTasks() {
	ctx := context.Background()
	log.Println("start********:", time.Now())
	ts.ctuc.SyncCompanyTasks(ctx)
	log.Println("stop!!!!!!!!!:", time.Now())
}

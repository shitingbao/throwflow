package service

import (
	"context"
)

func (ts *TaskService) SyncCompanyTasks() {
	ctx := context.Background()

	ts.ctuc.SyncCompanyTasks(ctx)
}

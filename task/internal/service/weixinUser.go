package service

import (
	"context"
)

func (ts *TaskService) SyncIntegralUsers() {
	ctx := context.Background()

	ts.wuuc.SyncIntegralUsers(ctx)
}

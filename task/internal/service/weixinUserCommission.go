package service

import (
	"context"
)

func (ts *TaskService) SyncTaskUserCommissions() {
	ctx := context.Background()

	ts.wucuc.SyncTaskUserCommissions(ctx)
}

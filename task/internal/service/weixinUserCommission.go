package service

import (
	"context"
)

func (ts *TaskService) SyncOrderUserCommissions() {
	ctx := context.Background()

	ts.wucuc.SyncOrderUserCommissions(ctx)
}

func (ts *TaskService) SyncCostOrderUserCommissions() {
	ctx := context.Background()

	ts.wucuc.SyncCostOrderUserCommissions(ctx)
}

package service

import (
	"context"
)

func (ts *TaskService) SyncUserBalances() {
	ctx := context.Background()

	ts.wubuc.SyncUserBalances(ctx)
}

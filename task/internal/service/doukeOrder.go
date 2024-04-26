package service

import (
	"context"
)

func (ts *TaskService) SyncDoukeOrders() {
	ctx := context.Background()

	ts.douc.SyncDoukeOrders(ctx)
}

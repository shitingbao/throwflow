package service

import (
	"context"
)

func (ts *TaskService) SyncJinritemaiStores() {
	ctx := context.Background()

	ts.jsuc.SyncJinritemaiStores(ctx)
}

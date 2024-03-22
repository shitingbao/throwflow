package service

import (
	"context"
)

func (ts *TaskService) SyncUpdateStatusCompanys() {
	ctx := context.Background()

	ts.cuc.SyncUpdateStatusCompanys(ctx)
}

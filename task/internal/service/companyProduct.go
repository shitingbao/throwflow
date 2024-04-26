package service

import (
	"context"
)

func (ts *TaskService) SyncCompanyProducts() {
	ctx := context.Background()

	ts.cpuc.SyncCompanyProducts(ctx)
}

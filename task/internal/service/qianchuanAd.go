package service

import (
	"context"
)

func (ts *TaskService) SyncQianchuanAds() {
	ctx := context.Background()

	ts.qaduc.SyncQianchuanAds(ctx)
}

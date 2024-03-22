package service

import (
	"context"
)

func (ts *TaskService) SyncQianchuanDatas() {
	ctx := context.Background()

	ts.qauc.SyncQianchuanDatas(ctx)
}

func (ts *TaskService) SyncRdsDatas() {
	ctx := context.Background()

	ts.qauc.SyncRdsDatas(ctx)
}

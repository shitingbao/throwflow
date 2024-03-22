package service

import (
	"context"
)

func (ts *TaskService) SyncOpenDouyinVideos() {
	ctx := context.Background()

	ts.odvuc.SyncOpenDouyinVideos(ctx)
}

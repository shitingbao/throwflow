package service

import (
	"context"
)

func (ts *TaskService) RefreshDoukeTokens() {
	ctx := context.Background()

	ts.dtuc.RefreshDoukeTokens(ctx)
}

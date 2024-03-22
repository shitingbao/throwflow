package service

import (
	"context"
)

func (ts *TaskService) RefreshOpenDouyinTokens() {
	ctx := context.Background()

	ts.odtuc.RefreshOpenDouyinTokens(ctx)
}

func (ts *TaskService) RenewRefreshTokensOpenDouyinTokens() {
	ctx := context.Background()

	ts.odtuc.RenewRefreshTokensOpenDouyinTokens(ctx)
}

func (ts *TaskService) SyncUserFansOpenDouyinTokens() {
	ctx := context.Background()

	ts.odtuc.SyncUserFansOpenDouyinTokens(ctx)
}

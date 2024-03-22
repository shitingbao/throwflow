package service

import (
	"context"
)

func (ts *TaskService) SyncUserCoupons() {
	ctx := context.Background()

	ts.wucouc.SyncUserCoupons(ctx)
}

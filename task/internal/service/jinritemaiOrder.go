package service

import (
	"context"
)

func (ts *TaskService) Sync90DayJinritemaiOrders() {
	ctx := context.Background()

	ts.jouc.Sync90DayJinritemaiOrders(ctx)
}

package service

import (
	"context"
)

func (ts *TaskService) SyncUserOrganizationRelations() {
	ctx := context.Background()

	ts.wuouc.SyncUserOrganizationRelations(ctx)
}

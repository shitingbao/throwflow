package service

import (
	"context"
)

func (ts *TaskService) SyncUpdateQrCodeCompanyOrganizations() {
	ctx := context.Background()

	ts.couc.SyncUpdateQrCodeCompanyOrganizations(ctx)
}

package service

import (
	"context"
)

func (ts *TaskService) RefreshOceanengineAccountTokens() {
	ctx := context.Background()

	ts.oatuc.RefreshOceanengineAccountTokens(ctx)
}

package service

import (
	"context"
	"main/internal/model"
	"main/pkg/apperror"
)

type Interface interface {
	GetBillSplitsByFilter(ctx context.Context, filter map[string]any) ([]model.BillSplit, apperror.Error)
	CalculateAndSaveBillSplits(ctx context.Context, userID, groupID uint64) (model.BillSplits, apperror.Error)
	RecalculateBillSplits(ctx context.Context, userID, groupID uint64) (model.BillSplits, apperror.Error)
	ClearBillSplitsForGroup(ctx context.Context, groupID uint64) apperror.Error
}

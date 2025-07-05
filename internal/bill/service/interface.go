package service

import (
	"context"
	"main/internal/model"
	"main/pkg/apperror"
)

type Interface interface {
	GetBills(ctx context.Context, filter map[string]any) (model.Bills, apperror.Error)
	CreateBill(ctx context.Context, bill model.Bill) apperror.Error
	UpdateBill(ctx context.Context, billID uint64, updates any) apperror.Error
	DeleteBill(ctx context.Context, billID uint64) apperror.Error
}

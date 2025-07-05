package service

import (
	"context"
	"log"
	"main/constants"
	"main/internal/bill/repository"
	"main/internal/model"
	"main/pkg/apperror"
	"main/util"
	"net/http"
	"sync"
	"time"
)

type Service struct {
	repository.Interface
}

var (
	syncOnce sync.Once
	svc      *Service
)

func NewService(r repository.Interface) *Service {
	syncOnce.Do(func() {
		svc = &Service{r}
	})

	return svc
}

func (s *Service) GetBills(ctx context.Context, filter map[string]any) (model.Bills, apperror.Error) {
	return s.GetAll(ctx, filter)
}

func (s *Service) CreateBill(ctx context.Context, bill model.Bill) apperror.Error {
	logTag := util.LogPrefix(ctx, "CreateBillForGroup")

	err := s.Create(ctx, &bill)
	if err.Exists() {
		log.Printf("%s failed to create bill for bill %v: %v", logTag, bill, err)
		return apperror.NewWithMessage("Failed to create bill", http.StatusBadRequest)
	}

	return apperror.Error{}
}

func (s *Service) UpdateBill(ctx context.Context, billID uint64, updates any) apperror.Error {
	logTag := util.LogPrefix(ctx, "UpdateUserBill")

	bill, err := s.Get(ctx, map[string]any{
		constants.ID: billID,
	})
	if err.Exists() || bill.ID == 0 {
		log.Printf("%s attempted to update invalid or non-owned bill %d: %v", logTag, billID, err)

		return apperror.NewWithMessage("Bill not found or unauthorized", http.StatusForbidden)
	}

	err = s.Update(ctx, map[string]any{
		constants.ID: billID,
	}, updates)
	if err.Exists() {
		log.Printf("%s failed to update bill %d %v", logTag, billID, err)

		return apperror.NewWithMessage("Failed to update bill", http.StatusBadRequest)
	}

	return apperror.Error{}
}

func (s *Service) DeleteBill(ctx context.Context, billID uint64) apperror.Error {
	logTag := util.LogPrefix(ctx, "DeleteBillByID")

	bill, err := s.Get(ctx, map[string]any{
		constants.ID: billID,
	})
	if err.Exists() || bill.ID == 0 {
		log.Printf("%s bill with ID %d not found or already deleted: %v", logTag, billID, err)

		return apperror.NewWithMessage("Bill not found", http.StatusNotFound)
	}

	err = s.Update(ctx, map[string]any{
		constants.ID: billID,
	}, map[string]any{
		constants.DeletedAt: time.Now(),
	})
	if err.Exists() {
		log.Printf("%s failed to soft delete bill ID %d: %v", logTag, billID, err)

		return apperror.NewWithMessage("Failed to delete bill", http.StatusBadRequest)
	}

	return apperror.Error{}
}

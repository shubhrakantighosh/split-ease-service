package service

import (
	"context"
	"log"
	"main/constants"
	billSvc "main/internal/bill/service"
	billSplitRepo "main/internal/bill_split/repository"
	groupSvc "main/internal/group/service"
	"main/internal/model"
	"main/pkg/apperror"
	"main/util"
	"math"
	"net/http"
	"sync"
)

type Service struct {
	billSplitRepo billSplitRepo.Interface
	billSvc       billSvc.Interface
	groupSvc      groupSvc.Interface
}

var (
	syncOnce sync.Once
	svc      *Service
)

func NewService(
	billSplitRepo billSplitRepo.Interface,
	billSvc billSvc.Interface,
	groupSvc groupSvc.Interface,
) *Service {
	syncOnce.Do(func() {
		svc = &Service{billSplitRepo: billSplitRepo, billSvc: billSvc, groupSvc: groupSvc}
	})

	return svc
}

func (s *Service) CalculateAndSaveBillSplits(ctx context.Context, userID, groupID uint64) (model.BillSplits, apperror.Error) {
	logTag := util.LogPrefix(ctx, "CalculateAndSaveBillSplits")

	isValid, err := s.validateUserGroupBillSplitsAccess(ctx, userID, groupID, model.Create)
	if err.Exists() {
		return nil, err
	}

	if !isValid {
		return nil, apperror.NewWithMessage("User is not authorized or no bill splits exist for this group", http.StatusForbidden)
	}

	bills, err := s.billSvc.GetBills(ctx, map[string]any{
		constants.GroupID: groupID,
	})
	if err.Exists() {
		log.Printf("%s failed to retrieve bills for group %d: %v", logTag, groupID, err)

		return nil, apperror.NewWithMessage("Failed to fetch bills", http.StatusBadRequest)
	}

	if len(bills) == 0 {
		return nil, apperror.NewWithMessage("No bills found for group", http.StatusNotFound)
	}

	memberSpend := make(map[uint64]float64)
	var total float64
	for _, bill := range bills {
		memberSpend[bill.UserID] += bill.PaidAmount
		total += bill.PaidAmount
	}

	numMembers := len(memberSpend)
	if numMembers == 0 {
		return nil, apperror.NewWithMessage("No members in group", http.StatusBadRequest)
	}

	perHead := total / float64(numMembers)

	debtors := make(map[uint64]float64)
	creditors := make(map[uint64]float64)

	for uid, paid := range memberSpend {
		diff := paid - perHead
		if diff < 0 {
			debtors[uid] = -diff
		} else if diff > 0 {
			creditors[uid] = diff
		}
	}

	billSplits := make(model.BillSplits, 0)

	for debtorID, due := range debtors {
		for creditorID, credit := range creditors {
			if due == 0 {
				break
			}
			if credit == 0 {
				continue
			}

			amount := math.Min(due, credit)

			billSplits = append(billSplits, &model.BillSplit{
				GroupID:     groupID,
				UserID:      debtorID,
				ToPayUserID: creditorID,
				AmountDue:   amount,
			})

			debtors[debtorID] -= amount
			creditors[creditorID] -= amount
			due -= amount
		}
	}

	err = s.billSplitRepo.CreateMany(ctx, billSplits)
	if err.Exists() {
		log.Printf("%s failed to save bill splits: %v", logTag, err)

		return nil, apperror.NewWithMessage("Failed to store bill splits", http.StatusBadRequest)
	}

	return billSplits, apperror.Error{}
}

func (s *Service) RecalculateBillSplits(ctx context.Context, userID, groupID uint64) (model.BillSplits, apperror.Error) {
	logTag := util.LogPrefix(ctx, "RecalculateBillSplits")

	err := s.ClearBillSplitsForGroup(ctx, groupID)
	if err.Exists() {
		log.Printf("%s failed to clear old bill splits for group %d: %v", logTag, groupID, err)
		return nil, apperror.NewWithMessage("Failed to clear old bill splits", http.StatusBadRequest)
	}

	return s.CalculateAndSaveBillSplits(ctx, userID, groupID)
}

func (s *Service) ClearBillSplitsForGroup(ctx context.Context, groupID uint64) apperror.Error {
	logTag := util.LogPrefix(ctx, "ClearBillSplitsForGroup")

	err := s.billSplitRepo.Delete(ctx,
		map[string]any{
			constants.GroupID: groupID,
		},
	)
	if err.Exists() {
		log.Printf("%s failed to soft delete bill splits for group %d : %v", logTag, groupID, err)

		return apperror.NewWithMessage("Failed to clear bill splits", http.StatusBadRequest)
	}

	return apperror.Error{}
}

func (s *Service) validateUserGroupBillSplitsAccess(
	ctx context.Context,
	userID,
	groupID uint64,
	permissionType model.PermissionType,
) (bool, apperror.Error) {
	logTag := util.LogPrefix(ctx, "validateUserGroupBillSplitsAccess")

	hasPermission, err := s.groupSvc.ValidateUserGroupPermission(ctx, userID, groupID, permissionType)
	if err.Exists() || !hasPermission {
		log.Printf("%s user %d does not have '%s' permission on group %d: %v", logTag, userID, permissionType, groupID, err)

		return false, apperror.NewWithMessage("Permission denied for accessing bill splits", http.StatusForbidden)
	}

	bills, err := s.billSplitRepo.GetAll(ctx, map[string]interface{}{constants.GroupID: groupID})
	if err.Exists() {
		log.Printf("%s failed to retrieve bill splits for group %d: %v", logTag, groupID, err)

		return false, apperror.NewWithMessage("Failed to retrieve bill splits for group", http.StatusBadRequest)
	}

	if len(bills) > 0 {
		log.Printf("%s failed to retrieve bill splits for group %d: %v", logTag, groupID, err)

		return false, apperror.NewWithMessage("Bill alreday splited", http.StatusBadRequest)
	}

	return len(bills) == 0, apperror.Error{}
}

package service

import (
	"context"
	"log"
	"main/constants"
	"main/internal/controller/request"
	"main/internal/model"
	"main/pkg/apperror"
	"main/util"
	"net/http"
)

func (s *Service) CreateGroupBill(
	ctx context.Context,
	currentUserID, userID, groupID uint64,
	req request.CreateBillRequest,
) apperror.Error {
	logTag := util.LogPrefix(ctx, "CreateGroupBill")
	hasPermission, err := s.ValidateUserGroupPermission(ctx, currentUserID, groupID, model.Create)
	if err.Exists() {
		log.Printf("%s failed to validate permission for user %d on group %d: %v", logTag, userID, groupID, err)
		return err
	}
	if !hasPermission {
		return apperror.NewWithMessage("Permission denied", http.StatusForbidden)
	}

	if currentUserID != userID && s.userSvc.IsUserValid(ctx, userID) {
		log.Printf("%s: invalid user ID %d", logTag, userID)

		return apperror.NewWithMessage("Please provide a valid user", http.StatusBadRequest)
	}

	bill := model.Bill{
		UserID:      userID,
		GroupID:     groupID,
		PaidAmount:  req.PaidAmount,
		Description: req.Description,
	}
	err = s.billSvc.CreateBill(ctx, bill)
	if err.Exists() {
		log.Printf("%s failed to create bill for user %d in group %d: %v", logTag, userID, groupID, err)
		return apperror.NewWithMessage("Failed to create bill", http.StatusBadRequest)
	}

	return apperror.Error{}
}

func (s *Service) UpdateGroupBill(
	ctx context.Context,
	userID, groupID, billID uint64,
	req request.UpdateBillRequest,
) apperror.Error {
	logTag := util.LogPrefix(ctx, "UpdateGroupBill")

	hasPermission, err := s.ValidateUserGroupPermission(ctx, userID, groupID, model.Edit)
	if err.Exists() {
		log.Printf("%s permission validation failed for user %d: %v", logTag, userID, err)
		return err
	}
	if !hasPermission {
		log.Printf("%s user %d lacks permission to update bills in group %d", logTag, userID, groupID)
		return apperror.NewWithMessage("Permission denied", http.StatusForbidden)
	}

	bill := model.Bill{
		PaidAmount:  req.PaidAmount,
		Description: req.Description,
	}
	err = s.billSvc.UpdateBill(ctx, billID, bill)
	if err.Exists() {
		log.Printf("%s failed to update bill for user %d: %v", logTag, userID, err)

		return apperror.NewWithMessage("Failed to update bill", http.StatusBadRequest)
	}

	return apperror.Error{}
}

func (s *Service) DeleteGroupBill(
	ctx context.Context,
	userID, groupID, billID uint64,
) apperror.Error {
	logTag := util.LogPrefix(ctx, "DeleteGroupBill")

	hasPermission, err := s.ValidateUserGroupPermission(ctx, userID, groupID, model.Delete)
	if err.Exists() {
		log.Printf("%s permission validation failed for user %d: %v", logTag, userID, err)
		return err
	}
	if !hasPermission {
		log.Printf("%s user %d lacks permission to delete bills in group %d", logTag, userID, groupID)

		return apperror.NewWithMessage("Permission denied", http.StatusForbidden)
	}

	err = s.billSvc.DeleteBill(ctx, billID)
	if err.Exists() {
		log.Printf("%s failed to delete bill %d for user %d: %v", logTag, billID, userID, err)
		return err
	}
	return apperror.Error{}
}

func (s *Service) ValidateUserGroupPermission(
	ctx context.Context,
	userID,
	groupID uint64,
	permissionType model.PermissionType,
) (bool, apperror.Error) {
	logTag := util.LogPrefix(ctx, "ValidateUserGroupPermission")

	group, err := s.groupRepo.Get(ctx, map[string]any{
		constants.ID: groupID,
	})
	if err.Exists() || group.ID == 0 {
		log.Printf("%s group not found with ID %d: %v", logTag, groupID, err)
		return false, apperror.NewWithMessage("Group not found or unauthorized access", http.StatusForbidden)
	}

	hasPermission, err := s.groupPermissionSvc.HasUserPermissionInGroup(ctx, userID, groupID, permissionType)
	if err.Exists() || !hasPermission {
		log.Printf("%s user %d does not have '%s' permission for group %d: %v", logTag, userID, permissionType, groupID, err)

		return false, apperror.NewWithMessage("Permission denied", http.StatusForbidden)
	}

	return true, apperror.Error{}
}

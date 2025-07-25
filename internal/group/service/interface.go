package service

import (
	"context"
	"main/internal/controller/request"
	"main/internal/controller/response"
	"main/internal/model"
	"main/pkg/apperror"
)

type Interface interface {
	CreateGroup(ctx context.Context, userID uint64, req request.CreateGroupRequest) apperror.Error
	UpdateGroup(ctx context.Context, userID, groupID uint64, req request.UpdateGroupRequest) apperror.Error
	RemoveGroup(ctx context.Context, userID, groupID uint64) apperror.Error

	GetUserGroupsWithPermissions(
		ctx context.Context,
		userID uint64,
		filters map[string]any,
	) (model.Groups, model.GroupUserPermissions, apperror.Error)

	FetchGroupDetailsByUserAccess(ctx context.Context, userID, groupID uint64) (*response.GroupDetails, apperror.Error)

	AssignUserToGroup(
		ctx context.Context,
		currentUserID, userID, groupID uint64,
	) apperror.Error

	CreateGroupBill(
		ctx context.Context,
		currentUserID, userID, groupID uint64,
		req request.CreateBillRequest,
	) apperror.Error

	UpdateGroupBill(
		ctx context.Context,
		userID, groupID, billID uint64,
		req request.UpdateBillRequest,
	) apperror.Error

	DeleteGroupBill(
		ctx context.Context,
		userID, groupID, billID uint64,
	) apperror.Error

	ValidateUserGroupPermission(
		ctx context.Context,
		userID,
		groupID uint64,
		permissionType model.PermissionType,
	) (bool, apperror.Error)
}

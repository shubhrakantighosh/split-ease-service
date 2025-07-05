package service

import (
	"context"
	"main/internal/controller/request"
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

	GetGroupDetails(ctx context.Context, userID, groupID uint64) (model.Group, apperror.Error)

	CreateGroupBill(
		ctx context.Context,
		userID, groupID uint64,
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
}

package service

import (
	"context"
	"main/internal/model"
	"main/pkg/apperror"
)

type Interface interface {
	AssignGroupPermissionsToUser(
		ctx context.Context,
		userID, groupID uint64,
		permissions []model.PermissionType,
	) apperror.Error

	FetchUserGroup(
		ctx context.Context,
		userID uint64,
	) (model.GroupUserPermissions, apperror.Error)

	HasUserPermissionInGroup(
		ctx context.Context,
		userID uint64,
		groupID uint64,
		permission model.PermissionType,
	) (bool, apperror.Error)

	DeleteGroupPermissions(
		ctx context.Context,
		groupID uint64,
	) apperror.Error
}

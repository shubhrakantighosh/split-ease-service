package service

import (
	"context"
	"log"
	"main/constants"
	"main/internal/group_permission/repository"
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

func (s *Service) AssignGroupPermissionsToUser(
	ctx context.Context,
	userID, groupID uint64,
	permissions []model.PermissionType,
) apperror.Error {
	logTag := util.LogPrefix(ctx, "AssignGroupPermissionsToUser")

	groupUserPermission := make(model.GroupUserPermissions, 0)
	for _, permission := range permissions {
		groupUserPermission = append(groupUserPermission, model.GroupUserPermission{
			UserID:         userID,
			GroupID:        groupID,
			PermissionType: permission,
			IsActive:       true,
		})
	}

	err := s.UpdateMany(ctx, groupUserPermission)
	if err.Exists() {
		log.Printf("%s failed to assign permissions %v to user %d in group %d: %v",
			logTag, permissions, userID, groupID, err)

		return apperror.NewWithMessage("Failed to assign permissions", http.StatusBadRequest)
	}

	return apperror.Error{}
}

func (s *Service) FetchUserGroup(
	ctx context.Context,
	userID uint64,
) (model.GroupUserPermissions, apperror.Error) {
	logTag := util.LogPrefix(ctx, "FetchUserGroupPermissions")

	filters := map[string]any{
		constants.UserID:   userID,
		constants.IsActive: true,
	}
	records, err := s.GetAll(ctx, filters)
	if err.Exists() {
		log.Printf("%s failed to fetch permissions for userID: %d. Error: %v", logTag, userID, err)

		return nil, apperror.NewWithMessage("Failed to fetch group permissions", http.StatusBadRequest)
	}

	return records, apperror.Error{}
}

func (s *Service) HasUserPermissionInGroup(
	ctx context.Context,
	userID uint64,
	groupID uint64,
	permission model.PermissionType,
) (bool, apperror.Error) {
	logTag := util.LogPrefix(ctx, "HasUserPermissionInGroup")

	filters := map[string]any{
		constants.UserID:         userID,
		constants.GroupID:        groupID,
		constants.PermissionType: permission,
		constants.IsActive:       true,
	}

	record, err := s.Get(ctx, filters)
	if err.Exists() {
		log.Printf("%s failed to check permission [%s] for user %d in group %d: %v",
			logTag, permission, userID, groupID, err)

		return false, apperror.NewWithMessage("Failed to check user permission", http.StatusBadRequest)
	}

	if record.ID == 0 {
		log.Printf("%s no permission [%s] found for user %d in group %d",
			logTag, permission, userID, groupID)

		return false, apperror.Error{}
	}

	return true, apperror.Error{}
}

func (s *Service) DeleteGroupPermissions(
	ctx context.Context,
	groupID uint64,
) apperror.Error {
	logTag := util.LogPrefix(ctx, "DeleteGroupPermissions")

	err := s.Update(ctx, map[string]interface{}{
		constants.GroupID: groupID,
	}, map[string]interface{}{
		constants.DeletedAt: time.Now(),
	})
	if err.Exists() {
		log.Printf("%s failed to soft-delete permissions for group %d: %v", logTag, groupID, err)

		return apperror.NewWithMessage("Failed to remove permissions from group", http.StatusBadRequest)
	}

	return apperror.Error{}
}

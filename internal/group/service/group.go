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
	"strconv"
	"sync"
	"time"
)

func (s *Service) CreateGroup(ctx context.Context, userID uint64, req request.CreateGroupRequest) apperror.Error {
	logTag := util.LogPrefix(ctx, "CreateGroup")

	group := model.Group{
		OwnerID:     userID,
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   strconv.FormatUint(userID, 10),
		UpdatedBy:   strconv.FormatUint(userID, 10),
	}
	err := s.groupRepo.Create(ctx, &group)
	if err.Exists() {
		log.Printf("%s failed to create group: %v", logTag, err)

		return apperror.NewWithMessage("Failed to create group", http.StatusBadRequest)
	}

	return s.groupPermissionSvc.AssignGroupPermissionsToUser(
		ctx,
		userID,
		group.ID,
		[]model.PermissionType{model.View, model.Edit, model.Delete},
	)
}

func (s *Service) UpdateGroup(ctx context.Context, userID, groupID uint64, req request.UpdateGroupRequest) apperror.Error {
	logTag := util.LogPrefix(ctx, "UpdateGroup")

	group, err := s.groupRepo.Get(ctx, map[string]any{constants.ID: groupID})
	if err.Exists() {
		return apperror.NewWithMessage("Permission denied", http.StatusForbidden)
	}

	if group.ID == 0 {
		return apperror.NewWithMessage("Permission denied", http.StatusForbidden)
	}

	hasPermission, err := s.groupPermissionSvc.HasUserPermissionInGroup(ctx, userID, groupID, model.Edit)
	if err.Exists() || !hasPermission {
		log.Printf("%s user %d does not have edit permission for group %d. Error: %v", logTag, userID, groupID, err)

		return apperror.NewWithMessage("Permission denied", http.StatusForbidden)
	}

	update := map[string]any{
		constants.Name:        req.Name,
		constants.Description: req.Description,
	}
	err = s.groupRepo.Update(ctx, map[string]any{constants.ID: groupID}, update)
	if err.Exists() {
		log.Printf("%s failed to update group %d: %v", logTag, groupID, err)

		return apperror.NewWithMessage("Failed to update group", http.StatusBadRequest)
	}

	return apperror.Error{}
}

func (s *Service) RemoveGroup(ctx context.Context, userID, groupID uint64) apperror.Error {
	logTag := util.LogPrefix(ctx, "RemoveGroup")

	group, err := s.groupRepo.Get(ctx, map[string]any{constants.ID: groupID})
	if err.Exists() || group.ID == 0 {
		log.Printf("%s failed to find group %d: %v", logTag, groupID, err)
		return apperror.NewWithMessage("Permission denied", http.StatusForbidden)
	}

	hasPermission, err := s.groupPermissionSvc.HasUserPermissionInGroup(ctx, userID, groupID, model.Delete)
	if err.Exists() || !hasPermission {
		log.Printf("%s user %d does not have delete permission for group %d. Error: %v", logTag, userID, groupID, err)
		return apperror.NewWithMessage("Permission denied", http.StatusForbidden)
	}

	var (
		wg      sync.WaitGroup
		errChan = make(chan apperror.Error, 3)
	)

	// 1. Soft-delete the group
	wg.Add(1)
	go func() {
		defer wg.Done()
		updateErr := s.groupRepo.Update(ctx, map[string]any{
			constants.ID: groupID,
		}, map[string]any{
			constants.DeletedAt: time.Now(),
		})
		if updateErr.Exists() {
			log.Printf("%s failed to mark group %d as deleted: %v", logTag, groupID, updateErr)
			errChan <- apperror.NewWithMessage("Failed to delete group", http.StatusInternalServerError)
		}
	}()

	// 2. Mark all group permissions as deleted
	wg.Add(1)
	go func() {
		defer wg.Done()
		updateErr := s.groupPermissionSvc.DeleteGroupPermissions(ctx, groupID)
		if updateErr.Exists() {
			log.Printf("%s failed to mark group permissions for group %d as deleted: %v", logTag, groupID, updateErr)
			errChan <- apperror.NewWithMessage("Failed to update group permissions", http.StatusInternalServerError)
		}
	}()

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	// Return first error if any occurred
	for e := range errChan {
		if e.Exists() {
			return e
		}
	}

	return apperror.Error{}
}

func (s *Service) GetUserGroupsWithPermissions(
	ctx context.Context,
	userID uint64,
	filters map[string]any,
) (model.Groups, model.GroupUserPermissions, apperror.Error) {
	logTag := util.LogPrefix(ctx, "FetchUserAccessibleGroups")

	groupPermissions, err := s.groupPermissionSvc.FetchUserGroup(ctx, userID)
	if err.Exists() {
		log.Printf("%s failed to fetch group permissions for user %d: %v", logTag, userID, err)

		return nil, nil, apperror.NewWithMessage("Failed to fetch user group permissions", http.StatusBadRequest)
	}

	filters[constants.ID] = groupPermissions.GetUniqueGroupIDs()
	groups, err := s.groupRepo.GetAll(ctx, filters)
	if err.Exists() {
		log.Printf("%s failed to fetch groups for user %d: %v", logTag, userID, err)

		return nil, nil, apperror.NewWithMessage("Unable to fetch user groups", http.StatusBadRequest)
	}

	return groups, groupPermissions, apperror.Error{}
}

func (s *Service) GetGroupDetails(ctx context.Context, userID, groupID uint64) (model.Group, apperror.Error) {
	logTag := util.LogPrefix(ctx, "GetUserGroupDetails")

	group, err := s.groupRepo.Get(ctx, map[string]any{
		constants.ID: groupID,
	})
	if err.Exists() {
		log.Printf("%s failed to fetch group %d for user %d: %v", logTag, groupID, userID, err)

		return model.Group{}, apperror.NewWithMessage("Group not found or access denied", http.StatusNotFound)
	}

	hasPermission, err := s.groupPermissionSvc.HasUserPermissionInGroup(ctx, userID, groupID, model.View)
	if err.Exists() || !hasPermission {
		log.Printf("%s user %d does not have delete permission for group %d. Error: %v", logTag, userID, groupID, err)

		// fix
		return model.Group{}, apperror.NewWithMessage("Permission denied", http.StatusForbidden)
	}

	return group, apperror.Error{}
}

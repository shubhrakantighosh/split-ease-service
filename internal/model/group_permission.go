package model

import (
	"gorm.io/gorm"
	"time"
)

type PermissionType string

const (
	View   PermissionType = "View"
	Edit   PermissionType = "Edit"
	Create PermissionType = "Create"
	Delete PermissionType = "Delete"
)

var PermissionTypeToString = map[PermissionType]string{
	View:   "view",
	Edit:   "edit",
	Create: "create",
	Delete: "delete",
}

func (p PermissionTypes) ToStringSlice() []string {
	strings := make([]string, 0, len(p))
	for _, perm := range p {
		if str, ok := PermissionTypeToString[perm]; ok {
			strings = append(strings, str)
		}
	}

	return strings
}

type PermissionTypes []PermissionType

type GroupUserPermission struct {
	ID             uint64         `json:"id"`
	GroupID        uint64         `json:"group_id"`
	UserID         uint64         `json:"user_id"`
	PermissionType PermissionType `json:"permission_type"`
	IsActive       bool           `json:"is_active"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at"`
}

type GroupUserPermissions []GroupUserPermission

func (g GroupUserPermissions) GetUniqueGroupIDs() []uint64 {
	uniqueGroupIDs := make([]uint64, 0)
	uniqueGroupID := make(map[uint64]struct{})
	if g == nil || len(g) == 0 {
		return uniqueGroupIDs
	}

	for _, permission := range g {
		if _, ok := uniqueGroupID[permission.GroupID]; !ok {
			uniqueGroupIDs = append(uniqueGroupIDs, permission.GroupID)
		}

		uniqueGroupID[permission.GroupID] = struct{}{}
	}

	return uniqueGroupIDs
}

func (g GroupUserPermissions) MapGroupIDToPermissions() map[uint64]PermissionTypes {
	groupIDMapPermissionTypes := make(map[uint64]PermissionTypes)
	if g == nil || len(g) == 0 {
		return groupIDMapPermissionTypes
	}

	for _, permission := range g {
		if _, ok := groupIDMapPermissionTypes[permission.GroupID]; !ok {
			groupIDMapPermissionTypes[permission.GroupID] = make(PermissionTypes, 0)
		}

		groupIDMapPermissionTypes[permission.GroupID] = append(groupIDMapPermissionTypes[permission.GroupID], permission.PermissionType)
	}

	return groupIDMapPermissionTypes
}

package adapter

import (
	"main/internal/controller/response"
	"main/internal/model"
)

func BuildAuthTokenResponse(req model.AuthToken) response.AuthTokenResponse {
	return response.AuthTokenResponse{
		AccessToken:      req.AccessToken,
		RefreshToken:     req.RefreshToken,
		AccessExpiresAt:  req.AccessExpiresAt,
		RefreshExpiresAt: req.RefreshExpiresAt,
	}
}

func BuildGroupPermissionsResponse(
	groups model.Groups,
	groupUserPermissions model.GroupUserPermissions,
) []response.GroupPermissionResponse {
	groupIDToPermissions := groupUserPermissions.MapGroupIDToPermissions()

	result := make([]response.GroupPermissionResponse, 0, len(groups))
	for _, group := range groups {
		result = append(result, response.GroupPermissionResponse{
			ID:          group.ID,
			Name:        group.Name,
			Description: group.Description,
			Permissions: groupIDToPermissions[group.ID].ToStringSlice(),
		})
	}

	return result
}

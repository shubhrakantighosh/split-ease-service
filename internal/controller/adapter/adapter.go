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

func BuildUsersResponse(
	users model.Users,
) response.Users {
	result := make(response.Users, 0, len(users))
	for _, user := range users {
		result = append(result, response.User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	return result
}

func BuildGroupDetailsResponse(
	group model.Group,
	users model.Users,
	bills model.Bills,
) *response.GroupDetails {
	idMap := users.MapByID()
	responseBills := make([]response.Bill, 0, len(bills))

	for _, bill := range bills {
		payer := idMap[bill.UserID]
		responseBills = append(responseBills, response.Bill{
			User: response.User{
				ID:    payer.ID,
				Name:  payer.Name,
				Email: payer.Email,
			},
			PaidAmount:  bill.PaidAmount,
			Description: bill.Description,
		})
	}

	return &response.GroupDetails{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description,
		Bills:       responseBills,
	}
}

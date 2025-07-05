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

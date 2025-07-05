package service

import (
	"context"
	ctrlReq "main/internal/controller/request"
	"main/internal/model"
	"main/pkg/apperror"
)

type Interface interface {
	FetchFilteredUsers(ctx context.Context, filter map[string]any) (model.Users, apperror.Error)
	GetUsers(ctx context.Context, currentUserID uint64) (model.Users, apperror.Error)
	CreateUserAccount(ctx context.Context, req ctrlReq.RegisterRequest) apperror.Error
	AuthenticateUser(ctx context.Context, email, password string) (model.AuthToken, apperror.Error)
	SendActivationEmail(ctx context.Context, email string) apperror.Error
	ActivateUserAccount(ctx context.Context, email, password, otp string) apperror.Error
	UpdateUserProfile(ctx context.Context, user model.User) apperror.Error
	IsUserValid(ctx context.Context, userID uint64) bool
}

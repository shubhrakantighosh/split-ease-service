package service

import (
	"context"
	"main/internal/model"
	"main/pkg/apperror"
)

type Interface interface {
	AuthenticateUser(ctx context.Context, email, password string) (model.AuthToken, apperror.Error)
	SendActivationEmail(ctx context.Context, email string) apperror.Error
	ActivateUserAccount(ctx context.Context, email, password, otp string) apperror.Error
	UpdateUserProfile(ctx context.Context, user model.User) apperror.Error
	IsUserValid(ctx context.Context, userID uint64) bool
}

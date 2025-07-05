package service

import (
	"context"
	"main/internal/model"
	"main/pkg/apperror"
)

type Interface interface {
	GenerateOTP(ctx context.Context, userID uint64, purpose model.Purpose) (string, apperror.Error)

	ValidateOTP(
		ctx context.Context,
		userID uint64,
		purpose model.Purpose,
	) (bool, apperror.Error)

	MarkOTPAsUsed(ctx context.Context, userID uint64, otpCode string) apperror.Error
}

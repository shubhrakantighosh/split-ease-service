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
		otp string,
	) (bool, apperror.Error)

	MarkOTPUsed(ctx context.Context, userID uint64, otpCode string) apperror.Error
}

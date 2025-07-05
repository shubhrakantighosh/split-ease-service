package service

import (
	"context"
	"gorm.io/gorm"
	"log"
	"main/constants"
	"main/internal/model"
	"main/internal/otp/repository"
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

func (s *Service) GenerateOTP(ctx context.Context, userID uint64, purpose model.Purpose) (string, apperror.Error) {
	logTag := util.LogPrefix(ctx, "GenerateOTP")

	code, err := util.GenerateRandomNumericCode(6)
	if err != nil {
		log.Println(logTag, "Failed to generate OTP code:", err)
		return "", apperror.NewWithMessage("Failed to generate OTP", http.StatusBadRequest)
	}

	otp := model.OTP{
		UserID:    userID,
		Code:      code,
		Purpose:   purpose,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	createErr := s.Create(ctx, &otp)
	if createErr.Exists() {
		log.Println(logTag, "Failed to store OTP in DB:", createErr)

		return "", apperror.NewWithMessage("Failed to create OTP", http.StatusBadRequest)
	}

	return code, apperror.Error{}
}

func (s *Service) ValidateOTP(
	ctx context.Context,
	userID uint64,
	purpose model.Purpose,
	otp string,
) (bool, apperror.Error) {
	logTag := util.LogPrefix(ctx, "ValidateOTP")

	scopes := []func(db *gorm.DB) *gorm.DB{
		func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		},
		func(db *gorm.DB) *gorm.DB {
			return db.Limit(1)
		},
		func(db *gorm.DB) *gorm.DB {
			return db.Where("used = ?", false)
		},
	}

	otps, err := s.GetAll(ctx, map[string]any{
		constants.UserID:  userID,
		constants.Purpose: purpose,
	}, scopes...)

	if err.Exists() {
		log.Println(logTag, "Failed to fetch OTPs:", err)

		return false, apperror.NewWithMessage("Unable to validate OTP", http.StatusBadRequest)
	}

	if len(otps) == 0 {
		log.Println(logTag, "No valid OTP found for user:", userID)

		return false, apperror.NewWithMessage("Invalid or expired OTP", http.StatusBadRequest)
	}

	latestOTP := otps[0]
	if latestOTP.ExpiresAt.Before(time.Now()) {
		log.Printf("%s OTP expired for user ID: %d", logTag, userID)

		return false, apperror.NewWithMessage("OTP has expired", http.StatusBadRequest)
	}

	if latestOTP.Code != otp {
		log.Printf("%s OTP code mismatch for user ID: %d", logTag, userID)

		return false, apperror.NewWithMessage("Invalid OTP", http.StatusBadRequest)
	}

	return true, apperror.Error{}
}

func (s *Service) MarkOTPUsed(ctx context.Context, userID uint64, otpCode string) apperror.Error {
	logTag := util.LogPrefix(ctx, "MarkOTPUsed")

	err := s.Update(ctx, map[string]any{
		constants.UserID: userID,
		constants.Code:   otpCode,
	}, map[string]any{
		constants.Used: true,
	})

	if err.Exists() {
		log.Printf("%s: failed to mark OTP as used for user %d and code %s: %v", logTag, userID, otpCode, err)

		return apperror.NewWithMessage("Unable to mark OTP as used", http.StatusBadRequest)
	}

	return apperror.Error{}
}

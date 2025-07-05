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
	syncOnce *sync.Once
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
		return "", apperror.NewWithMessage("Failed to generate OTP", http.StatusInternalServerError)
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

		return "", apperror.NewWithMessage("Failed to create OTP", http.StatusInternalServerError)
	}

	return code, apperror.Error{}
}

func (s *Service) ValidateOTP(
	ctx context.Context,
	userID uint64,
	purpose model.Purpose,
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

		return false, apperror.NewWithMessage("Unable to validate OTP", http.StatusInternalServerError)
	}

	if len(otps) == 0 {
		log.Println(logTag, "No valid OTP found for user:", userID)

		return false, apperror.NewWithMessage("Invalid or expired OTP", http.StatusBadRequest)
	}

	otp := otps[0]
	if otp.ExpiresAt.Before(time.Now()) {
		log.Println(logTag, "OTP expired for user:", userID)

		return false, apperror.NewWithMessage("OTP has expired", http.StatusBadRequest)
	}

	return true, apperror.Error{}
}

func (s *Service) MarkOTPAsUsed(ctx context.Context, userID uint64, otpCode string) apperror.Error {
	logTag := util.LogPrefix(ctx, "MarkOTPAsUsed")

	err := s.Update(ctx, map[string]any{
		constants.UserID: userID,
		constants.Code:   otpCode,
	}, map[string]any{
		constants.Used: true,
	})

	if err.Exists() {
		log.Println(logTag, "Failed to mark OTP as used:", err)
		return apperror.NewWithMessage("Failed to update OTP status", http.StatusInternalServerError)
	}

	return apperror.Error{}
}

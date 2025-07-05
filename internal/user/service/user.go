package service

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"main/constants"
	ctrlReq "main/internal/controller/request"
	"main/internal/model"
	"main/pkg/apperror"
	"main/util"
	"net/http"
)

func (s *Service) CreateUserAccount(ctx context.Context, req ctrlReq.RegisterRequest) apperror.Error {
	logTag := util.LogPrefix(ctx, "CreateUserAccount")

	users, err := s.repo.GetAll(ctx, map[string]any{constants.Email: req.Email})
	if err.Exists() {
		log.Printf("%s: failed to check existing user for email %s: %v", logTag, req.Email, err)

		return apperror.NewWithMessage("Failed to validate user", http.StatusBadRequest)
	}

	if len(users) > 0 {
		log.Printf("%s: user already exists with email %s", logTag, req.Email)
		return apperror.NewWithMessage("User already exists", http.StatusConflict)
	}

	hashedPass, hashErr := hashPassword(req.Password)
	if hashErr.Exists() {
		log.Printf("%s: failed to hash password for email %s: %v", logTag, req.Email, hashErr)

		return apperror.NewWithMessage("Failed to process password", http.StatusBadRequest)
	}

	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPass,
		IsActive: true,
	}
	createErr := s.repo.Create(ctx, &user)
	if createErr.Exists() {
		log.Printf("%s: failed to create user for email %s: %v", logTag, req.Email, createErr)

		return apperror.NewWithMessage("Failed to create user", http.StatusBadRequest)
	}

	return apperror.Error{}
}

func (s *Service) AuthenticateUser(ctx context.Context, email, password string) (model.AuthToken, apperror.Error) {
	logTag := util.LogPrefix(ctx, "AuthenticateUser")

	user, err := s.repo.Get(ctx, map[string]any{constants.Email: email})
	if err.Exists() {
		log.Printf("%s failed to get user by email %s: %v", logTag, email, err)

		return model.AuthToken{}, apperror.NewWithMessage("User not found", http.StatusNotFound)
	}

	if !user.IsActive {
		log.Printf("%s user %s is not active", logTag, email)

		return model.AuthToken{}, apperror.NewWithMessage("User account is not active", http.StatusUnauthorized)
	}

	if !checkPasswordHash(user.Password, password) {
		log.Printf("%s invalid password for user %s", logTag, email)

		return model.AuthToken{}, apperror.NewWithMessage("Invalid credentials", http.StatusUnauthorized)
	}

	return s.authSvc.GenerateOrUpdateAuthToken(ctx, user.ID)
}

func (s *Service) SendActivationEmail(ctx context.Context, email string) apperror.Error {
	logTag := util.LogPrefix(ctx, "SendActivationEmail")

	user, err := s.repo.Get(ctx, map[string]any{constants.Email: email})
	if err.Exists() {
		log.Printf("%s failed to fetch user %s: %v", logTag, email, err)

		return apperror.NewWithMessage("User lookup failed", http.StatusBadRequest)
	}

	if user.IsActive {
		log.Printf("%s user %s is already active", logTag, email)

		return apperror.NewWithMessage("Account is already activated", http.StatusBadRequest)
	}

	otp, err := s.otpSvc.GenerateOTP(ctx, user.ID, model.Activation)
	if err.Exists() {
		log.Printf("%s failed to generate OTP for user %d: %v", logTag, user.ID, err)

		return apperror.NewWithMessage("Failed to generate OTP", http.StatusBadRequest)
	}

	// TODO: Send OTP via email service
	fmt.Println("OTP:", otp)
	return apperror.Error{}
}

func (s *Service) ActivateUserAccount(ctx context.Context, email, password, otp string) apperror.Error {
	logTag := util.LogPrefix(ctx, "ActivateUserAccount")

	user, err := s.repo.Get(ctx, map[string]any{constants.Email: email})
	if err.Exists() {
		log.Printf("%s failed to get user by email %s: %v", logTag, email, err)

		return apperror.NewWithMessage("Failed to fetch user", http.StatusBadRequest)
	}

	if user.IsActive {
		log.Printf("%s user %s is already active", logTag, email)

		return apperror.NewWithMessage("Account already activated", http.StatusBadRequest)
	}

	isValid, err := s.otpSvc.ValidateOTP(ctx, user.ID, model.Activation, otp)
	if err.Exists() {
		log.Printf("%s OTP validation failed for user %d: %v", logTag, user.ID, err)
		return err
	}

	if !isValid {
		log.Printf("%s invalid or expired OTP for user %d", logTag, user.ID)
		return apperror.NewWithMessage("Invalid or expired OTP", http.StatusBadRequest)
	}

	user.IsActive = true
	user.Password = password
	err = s.UpdateUserProfile(ctx, user)
	if err.Exists() {
		log.Printf("%s - Failed to update user profile after successful OTP validation for user ID %d: %v", logTag, user.ID, err)
		return err
	}

	return s.otpSvc.MarkOTPUsed(ctx, user.ID, otp)
}

func (s *Service) UpdateUserProfile(ctx context.Context, user model.User) apperror.Error {
	logTag := util.LogPrefix(ctx, "UpdateUserProfile")

	if len(user.Email) > 0 {
		users, err := s.repo.GetAll(ctx, map[string]any{constants.Email: user.Email})
		if err.Exists() {
			log.Printf("%s failed to check email %s: %v", logTag, user.Email, err)

			return apperror.NewWithMessage("Something went wrong while checking email", http.StatusBadRequest)
		}

		if len(users) > 0 && users[0].ID != user.ID {
			log.Printf("%s email %s is already used by user %d", logTag, user.Email, users[0].ID)

			return apperror.NewWithMessage("This email is already registered", http.StatusBadRequest)
		}
	}

	if len(user.Password) > 0 {
		hashedPassword, err := hashPassword(user.Password)
		if err.Exists() {
			log.Printf("%s failed to hash password: %v", logTag, err)

			return apperror.NewWithMessage("Failed to hash password", http.StatusBadRequest)
		}
		user.Password = hashedPassword
	}

	err := s.repo.Update(ctx, map[string]any{constants.ID: user.ID}, &user)
	if err.Exists() {
		log.Printf("%s failed to update user %d: %v", logTag, user.ID, err)

		return apperror.NewWithMessage("Failed to update user", http.StatusBadRequest)
	}

	return s.authSvc.MarkTokenExpired(ctx, user.ID)
}

func (s *Service) IsUserValid(ctx context.Context, userID uint64) bool {
	logTag := util.LogPrefix(ctx, "IsUserValid")

	user, err := s.repo.Get(ctx, map[string]interface{}{constants.ID: userID})
	if err.Exists() {
		log.Printf("%s failed to fetch user %d: %v", logTag, userID, err)
		return false
	}

	if user.ID == 0 {
		log.Printf("%s user %d not found", logTag, userID)
		return false
	}

	if !user.IsActive {
		log.Printf("%s user %d is inactive", logTag, userID)
		return false
	}

	return true
}

func hashPassword(password string) (string, apperror.Error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Hashing error: %v", err)

		return "", apperror.NewWithMessage("Failed to hash password", http.StatusBadRequest)
	}

	return string(hashed), apperror.Error{}
}

func checkPasswordHash(hashedPassword, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	return err == nil
}

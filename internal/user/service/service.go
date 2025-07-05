package service

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"main/constants"
	authSvc "main/internal/auth/service"
	"main/internal/model"
	otpSvc "main/internal/otp/service"
	"main/internal/user/repository"
	"main/pkg/apperror"
	"main/util"
	"net/http"
	"sync"
)

type Service struct {
	repo    repository.Interface
	authSvc authSvc.Interface
	otpSvc  otpSvc.Interface
}

var (
	syncOnce *sync.Once
	svc      *Service
)

func NewService(repo repository.Interface, authSvc authSvc.Interface, otpSvc otpSvc.Interface) *Service {
	syncOnce.Do(func() {
		svc = &Service{repo: repo, authSvc: authSvc, otpSvc: otpSvc}
	})

	return svc
}

func (s *Service) AuthenticateUser(ctx context.Context, email, password string) (model.AuthToken, apperror.Error) {
	logTag := util.LogPrefix(ctx, "AuthenticateUser")
	filter := map[string]any{
		constants.Email: email,
	}

	user, err := s.repo.Get(ctx, filter)
	if err.Exists() {
		log.Println(logTag, err)

		return model.AuthToken{}, apperror.NewWithMessage("", http.StatusNotFound)
	}

	if !user.IsActive {
		log.Println(logTag, err)

		return model.AuthToken{}, apperror.NewWithMessage("Please acrive", http.StatusNotFound)
	}

	if !checkPasswordHash(user.Password, password) {
		log.Println(logTag, err)

		return model.AuthToken{}, apperror.NewWithMessage("Invalid Password", http.StatusNotFound)
	}

	return s.authSvc.GenerateOrUpdateAuthToken(ctx, user.ID)
}

func (s *Service) SendActivationEmail(ctx context.Context, email string) apperror.Error {
	logTag := util.LogPrefix(ctx, "SendActivationEmail")

	user, err := s.repo.Get(ctx, map[string]any{constants.Email: email})
	if err.Exists() {
		log.Printf("%s failed to check email existence for %s: %v", logTag, email, err)

		return apperror.NewWithMessage("Something went wrong while checking email", http.StatusInternalServerError)
	}

	if user.IsActive {
		log.Printf(fmt.Sprintf("%s is already active", email))

		return apperror.NewWithMessage("This email is already registered", http.StatusBadRequest)
	}

	otp, err := s.otpSvc.GenerateOTP(ctx, user.ID, model.Activation)
	if err.Exists() {
		log.Printf("%s failed to generate otp for user %s: %v", logTag, user.ID, err)

		return apperror.NewWithMessage("Something went wrong while generating otp for user", http.StatusInternalServerError)
	}

	// sent otp to email

	fmt.Println(otp)
	return apperror.Error{}
}

func (s *Service) ActivateUserAccount(ctx context.Context, email, password, otp string) apperror.Error {
	logTag := util.LogPrefix(ctx, "ActivateUserAccount")

	user, err := s.repo.Get(ctx, map[string]any{constants.Email: email})
	if err.Exists() {
		log.Printf("%s failed to check email existence for %s: %v", logTag, email, err)

		return apperror.NewWithMessage("Something went wrong while checking email", http.StatusInternalServerError)
	}

	if user.IsActive {
		log.Printf(fmt.Sprintf("%s is already active", email))

		return apperror.NewWithMessage("This email is already registered", http.StatusBadRequest)
	}

	// check on for this otp valid or not based opn time
	ok, err := s.otpSvc.ValidateOTP(ctx, user.ID, model.Activation)
	if err.Exists() {
		log.Printf("%s failed to validate otp for user %s: %v", logTag, user.ID, err)

		return err
	}

	if !ok {
		return apperror.NewWithMessage("generated new otp", http.StatusBadRequest)
	}

	hassPass, hasErr := hashPassword(password)
	if hasErr.Exists() {
		log.Println(logTag, err)

		return apperror.NewWithMessage("password is incorrect", http.StatusBadRequest)
	}

	user.IsActive = true
	user.Password = hassPass
	return s.UpdateUserProfile(ctx, user)
}

func (s *Service) UpdateUserProfile(ctx context.Context, user model.User) apperror.Error {
	logTag := util.LogPrefix(ctx, "UpdateUserProfile")

	if len(user.Email) > 0 {
		users, err := s.repo.GetAll(ctx, map[string]any{constants.Email: user.Email})
		if err.Exists() {
			log.Println("%s failed to check email existence for %s: %v", logTag, user.Email, err)

			return apperror.NewWithMessage("Something went wrong while checking email", http.StatusInternalServerError)
		}

		if len(users) > 0 && users[0].ID != user.ID {
			log.Printf("%s email %s already in use by user %d", logTag, user.Email, users[0].ID)

			return apperror.NewWithMessage("This email is already registered", http.StatusBadRequest)
		}
	}

	if len(user.Password) > 0 {
		hashPass, err := hashPassword(user.Password)
		if err.Exists() {
			log.Printf("%s failed to hash password: %v", logTag, err)

			return apperror.NewWithMessage("Please try again", http.StatusInternalServerError)
		}

		user.Password = hashPass
	}

	err := s.repo.Update(ctx, map[string]any{constants.ID: user.ID}, &user)
	if err.Exists() {
		log.Printf("%s failed to update user %d: %v", logTag, user.ID, err)

		return apperror.NewWithMessage("Failed to update user profile", http.StatusInternalServerError)
	}

	return s.authSvc.MarkTokenExpired(ctx, user.ID)
}

func (s *Service) IsUserValid(ctx context.Context, userID uint64) bool {
	logTag := util.LogPrefix(ctx, "IsUserValid")

	user, err := s.repo.Get(ctx, map[string]interface{}{
		constants.ID: userID,
	})

	if err.Exists() {
		log.Printf("%s failed to fetch user %d: %v", logTag, userID, err)
		return false
	}

	if user.ID == 0 {
		log.Printf("%s user %d does not exist", logTag, userID)
		return false
	}

	if !user.IsActive {
		log.Printf("%s user %d is not active", logTag, userID)
		return false
	}

	return true
}

func hashPassword(password string) (string, apperror.Error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to hash password: %v", err)

		return "", apperror.NewWithMessage("Please try gain", http.StatusBadRequest)
	}

	return string(hashedBytes), apperror.Error{}
}

func checkPasswordHash(hashedPassword, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	return err == nil
}

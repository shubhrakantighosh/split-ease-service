package service

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"main/constants"
	"main/internal/auth/repository"
	"main/internal/jwt/private"
	"main/internal/model"
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

func (s *Service) GenerateOrUpdateAuthToken(ctx context.Context, userID uint64) (model.AuthToken, apperror.Error) {
	logTag := util.LogPrefix(ctx, "GenerateOrUpdateAuthToken")

	// Generate access + refresh token pair
	token, tokenErr := generateTokenPair(userID)
	if tokenErr != nil {
		log.Println(logTag, "Failed to generate token pair:", tokenErr)

		return model.AuthToken{}, apperror.NewWithMessage("Failed to generate tokens", http.StatusBadRequest)
	}

	// Prepare AuthToken model
	authToken := model.AuthToken{
		UserID:           userID,
		AccessToken:      token.AccessToken,
		RefreshToken:     token.RefreshToken,
		AccessExpiresAt:  token.AccessExpiresAt,
		RefreshExpiresAt: token.RefreshExpiresAt,
	}

	// Check if token already exists for user
	existingTokens, err := s.GetAll(ctx, map[string]any{
		constants.UserID: userID,
	})
	if err.Exists() {
		log.Println(logTag, "Failed to check existing tokens:", err)

		return model.AuthToken{}, apperror.NewWithMessage("Failed to update token", http.StatusBadRequest)
	}

	if len(existingTokens) > 0 {
		updateErr := s.Update(ctx, map[string]any{
			constants.UserID: userID,
		}, &authToken)
		if updateErr.Exists() {
			log.Println(logTag, "Failed to update existing token:", updateErr)

			return model.AuthToken{}, apperror.NewWithMessage("Failed to update token", http.StatusBadRequest)
		}

		return authToken, apperror.Error{}
	}

	createErr := s.Create(ctx, &authToken)
	if createErr.Exists() {
		log.Println(logTag, "Failed to create new token:", createErr)

		return model.AuthToken{}, apperror.NewWithMessage("Failed to create token", http.StatusBadRequest)
	}

	return authToken, apperror.Error{}
}

func (s *Service) MarkTokenExpired(ctx context.Context, userID uint64) apperror.Error {
	logTag := util.LogPrefix(ctx, "MarkTokenExpired")

	existingTokens, err := s.GetAll(ctx, map[string]any{
		constants.UserID: userID,
	})
	if err.Exists() {
		log.Println(logTag, "Failed to check existing tokens:", err)

		return apperror.NewWithMessage("Failed to update token", http.StatusBadRequest)
	}

	if len(existingTokens) > 0 {
		return apperror.Error{}
	}

	token := existingTokens[0]
	token.RefreshExpiresAt = time.Now().Add(-10)

	err = s.Update(ctx, map[string]any{
		constants.ID: token.ID,
	}, &token)
	if err.Exists() {
		log.Println(logTag, "Failed to update existing token:", token, err)

		return apperror.NewWithMessage("Failed to update token", http.StatusBadRequest)
	}

	return apperror.Error{}
}

var jwtKey = []byte("access_secret")      // Use env var
var refreshKey = []byte("refresh_secret") // Use env var

type tokenPair struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
}

func generateTokenPair(userID uint64) (*tokenPair, error) {
	userDetails := private.UserDetails{UserID: userID}

	accessExpiresAt := time.Now().Add(15 * time.Minute)
	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	// Access Token
	accessClaims := &private.Claims{
		UserDetails: userDetails,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	// Refresh Token
	refreshClaims := &private.Claims{
		UserDetails: userDetails,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(refreshKey)
	if err != nil {
		return nil, err
	}

	return &tokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresAt:  accessExpiresAt,
		RefreshExpiresAt: refreshExpiresAt,
	}, nil
}

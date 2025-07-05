package service

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"main/internal/auth/repository"
	"main/internal/model"
	"time"
)

type Service struct {
	repository.Interface
}

func NewService(r repository.Interface) *Service {
	return &Service{r}
}

func (s *Service) Get(ctx context.Context) {
	err := s.Create(ctx, &model.AuthToken{
		UserID:       0,
		AccessToken:  "hbjn",
		RefreshToken: "abhjnk",
		ExpiresAt:    time.Time{},
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
		DeletedAt:    gorm.DeletedAt{},
	})

	fmt.Println(err)
}

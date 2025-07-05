package service

import (
	"main/internal/auth/repository"
)

type Service struct {
	repository.Interface
}

func NewService(r repository.Interface) *Service {
	return &Service{r}
}

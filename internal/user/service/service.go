package service

import (
	authSvc "main/internal/auth/service"
	otpSvc "main/internal/otp/service"
	"main/internal/user/repository"
	"sync"
)

type Service struct {
	repo    repository.Interface
	authSvc authSvc.Interface
	otpSvc  otpSvc.Interface
}

var (
	syncOnce sync.Once
	svc      *Service
)

func NewService(repo repository.Interface, authSvc authSvc.Interface, otpSvc otpSvc.Interface) *Service {
	syncOnce.Do(func() {
		svc = &Service{repo: repo, authSvc: authSvc, otpSvc: otpSvc}
	})

	return svc
}

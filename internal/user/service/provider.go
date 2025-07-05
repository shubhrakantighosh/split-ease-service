package service

import (
	"github.com/google/wire"
	authRepo "main/internal/auth/repository"
	authSvc "main/internal/auth/service"
	otpRepo "main/internal/otp/repository"
	otpSvc "main/internal/otp/service"
	userRepo "main/internal/user/repository"
)

var ProviderSet = wire.NewSet(
	NewService,
	userRepo.NewRepository,
	authRepo.NewRepository,
	authSvc.NewService,
	otpRepo.NewRepository,
	otpSvc.NewService,

	// bind each one of the interfaces
	wire.Bind(new(Interface), new(*Service)),
	wire.Bind(new(userRepo.Interface), new(*userRepo.Repository)),
	wire.Bind(new(authRepo.Interface), new(*authRepo.Repository)),
	wire.Bind(new(authSvc.Interface), new(*authSvc.Service)),
	wire.Bind(new(otpRepo.Interface), new(*otpRepo.Repository)),
	wire.Bind(new(otpSvc.Interface), new(*otpSvc.Service)),
)

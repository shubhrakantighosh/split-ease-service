package controller

import (
	"github.com/google/wire"
	authRepo "main/internal/auth/repository"
	authSvc "main/internal/auth/service"
	otpRepo "main/internal/otp/repository"
	otpSvc "main/internal/otp/service"
	userRepo "main/internal/user/repository"
	userSvc "main/internal/user/service"
)

var ProviderSet = wire.NewSet(
	NewController,
	userSvc.NewService,
	userRepo.NewRepository,
	otpSvc.NewService,
	otpRepo.NewRepository,
	authSvc.NewService,
	authRepo.NewRepository,

	// bind each one of the interfaces
	wire.Bind(new(Interface), new(*Controller)),
	wire.Bind(new(userSvc.Interface), new(*userSvc.Service)),
	wire.Bind(new(userRepo.Interface), new(*userRepo.Repository)),
	wire.Bind(new(otpSvc.Interface), new(*otpSvc.Service)),
	wire.Bind(new(otpRepo.Interface), new(*otpRepo.Repository)),
	wire.Bind(new(authSvc.Interface), new(*authSvc.Service)),
	wire.Bind(new(authRepo.Interface), new(*authRepo.Repository)),
)

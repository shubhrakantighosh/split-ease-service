package service

import (
	"github.com/google/wire"
	authRepo "main/internal/auth/repository"
	authSvc "main/internal/auth/service"
	billRepo "main/internal/bill/repository"
	billSvc "main/internal/bill/service"
	billSplitRepo "main/internal/bill_split/repository"
	groupRepo "main/internal/group/repository"
	groupSvc "main/internal/group/service"
	groupPermissionRepo "main/internal/group_permission/repository"
	groupPermissionSvc "main/internal/group_permission/service"
	otpRepo "main/internal/otp/repository"
	otpSvc "main/internal/otp/service"
	userRepo "main/internal/user/repository"
	userSvc "main/internal/user/service"
)

var ProviderSet = wire.NewSet(
	NewService,
	billSplitRepo.NewRepository,
	billSvc.NewService,
	billRepo.NewRepository,
	groupRepo.NewRepository,
	groupSvc.NewService,
	groupPermissionRepo.NewRepository,
	groupPermissionSvc.NewService,
	userRepo.NewRepository,
	userSvc.NewService,
	authRepo.NewRepository,
	authSvc.NewService,
	otpRepo.NewRepository,
	otpSvc.NewService,

	// bind each one of the interfaces
	wire.Bind(new(Interface), new(*Service)),
	wire.Bind(new(billSplitRepo.Interface), new(*billSplitRepo.Repository)),
	wire.Bind(new(billSvc.Interface), new(*billSvc.Service)),
	wire.Bind(new(billRepo.Interface), new(*billRepo.Repository)),
	wire.Bind(new(groupRepo.Interface), new(*groupRepo.Repository)),
	wire.Bind(new(groupSvc.Interface), new(*groupSvc.Service)),
	wire.Bind(new(groupPermissionRepo.Interface), new(*groupPermissionRepo.Repository)),
	wire.Bind(new(groupPermissionSvc.Interface), new(*groupPermissionSvc.Service)),
	wire.Bind(new(userRepo.Interface), new(*userRepo.Repository)),
	wire.Bind(new(userSvc.Interface), new(*userSvc.Service)),
	wire.Bind(new(authRepo.Interface), new(*authRepo.Repository)),
	wire.Bind(new(authSvc.Interface), new(*authSvc.Service)),
	wire.Bind(new(otpSvc.Interface), new(*otpSvc.Service)),
	wire.Bind(new(otpRepo.Interface), new(*otpRepo.Repository)),
)

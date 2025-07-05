package service

import (
	"github.com/google/wire"
	billRepo "main/internal/bill/repository"
	billSvc "main/internal/bill/service"
	"main/internal/group/repository"
	groupPermissionRepo "main/internal/group_permission/repository"
	groupPermissionSvc "main/internal/group_permission/service"
)

var ProviderSet = wire.NewSet(
	NewService,
	repository.NewRepository,
	groupPermissionSvc.NewService,
	groupPermissionRepo.NewRepository,
	billSvc.NewService,
	billRepo.NewRepository,

	// bind each one of the interfaces
	wire.Bind(new(Interface), new(*Service)),
	wire.Bind(new(repository.Interface), new(*repository.Repository)),
	wire.Bind(new(groupPermissionSvc.Interface), new(*groupPermissionSvc.Service)),
	wire.Bind(new(groupPermissionRepo.Interface), new(*groupPermissionRepo.Repository)),
	wire.Bind(new(billSvc.Interface), new(*billSvc.Service)),
	wire.Bind(new(billRepo.Interface), new(*billRepo.Repository)),
)

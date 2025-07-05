package service

import (
	"github.com/google/wire"
	"main/internal/group/repository"
	groupPermissionRepo "main/internal/group_permission/repository"
	groupPermissionSvc "main/internal/group_permission/service"
)

var ProviderSet = wire.NewSet(
	NewService,
	repository.NewRepository,
	groupPermissionSvc.NewService,
	groupPermissionRepo.NewRepository,

	// bind each one of the interfaces
	wire.Bind(new(Interface), new(*Service)),
	wire.Bind(new(repository.Interface), new(*repository.Repository)),
	wire.Bind(new(groupPermissionSvc.Interface), new(*groupPermissionSvc.Service)),
	wire.Bind(new(groupPermissionRepo.Interface), new(*groupPermissionRepo.Repository)),
)

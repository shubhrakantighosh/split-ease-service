package service

import (
	"github.com/google/wire"
	"main/internal/group_permission/repository"
)

var ProviderSet = wire.NewSet(
	NewService,
	repository.NewRepository,

	// bind each one of the interfaces
	wire.Bind(new(Interface), new(*Service)),
	wire.Bind(new(repository.Interface), new(*repository.Repository)),
)

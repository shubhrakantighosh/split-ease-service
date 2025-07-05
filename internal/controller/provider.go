package controller

import (
	"github.com/google/wire"
	"main/internal/auth/repository"
	"main/internal/auth/service"
)

var ProviderSet = wire.NewSet(
	NewController,
	service.NewService,
	repository.NewRepository,

	// bind each one of the interfaces
	wire.Bind(new(Interface), new(*Controller)),
	wire.Bind(new(service.Interface), new(*service.Service)),
	wire.Bind(new(repository.Interface), new(*repository.Repository)),
)

// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package service

import (
	"context"
	repository2 "main/internal/auth/repository"
	"main/internal/auth/service"
	repository3 "main/internal/otp/repository"
	service2 "main/internal/otp/service"
	"main/internal/user/repository"
	"main/pkg/db/postgres"
)

// Injectors from wire.go:

func Wire(ctx context.Context, db *postgres.DbCluster) *Service {
	repositoryRepository := repository.NewRepository(db)
	repository4 := repository2.NewRepository(db)
	serviceService := service.NewService(repository4)
	repository5 := repository3.NewRepository(db)
	service3 := service2.NewService(repository5)
	service4 := NewService(repositoryRepository, serviceService, service3)
	return service4
}

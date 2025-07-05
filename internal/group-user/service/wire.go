//go:build wireinject
// +build wireinject

package service

import (
	"context"
	"github.com/google/wire"
	"main/pkg/db/postgres"
)

func Wire(ctx context.Context, db *postgres.DbCluster) *Service {
	panic(wire.Build(ProviderSet))
}

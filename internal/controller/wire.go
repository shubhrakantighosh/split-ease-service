//go:build wireinject
// +build wireinject

package controller

import (
	"context"
	"github.com/google/wire"
	"main/pkg/db/postgres"
)

func Wire(ctx context.Context, db *postgres.DbCluster) *Controller {
	panic(wire.Build(ProviderSet))
}

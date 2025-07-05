package service

import (
	"context"
	"main/internal/model"
	"main/pkg/apperror"
)

type Interface interface {
	GenerateOrUpdateAuthToken(ctx context.Context, userID uint64) (model.AuthToken, apperror.Error)
	MarkTokenExpired(ctx context.Context, userID uint64) apperror.Error
}

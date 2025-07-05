package repository

import (
	"context"
	"gorm.io/gorm"
	"main/pkg/apperror"
)

type Interface[T any] interface {
	Create(ctx context.Context, data *T) apperror.Error

	CreateMany(ctx context.Context, data []*T) apperror.Error

	Update(
		ctx context.Context,
		filter map[string]interface{},
		updates map[string]interface{},
	) apperror.Error

	UpdateMany(ctx context.Context, data []T) apperror.Error

	GetAll(
		ctx context.Context,
		filter map[string]interface{},
		scopes ...func(db *gorm.DB) *gorm.DB,
	) (results []T, err apperror.Error)

	GetAllWithPagination(
		ctx context.Context,
		filter map[string]interface{},
		scopes ...func(db *gorm.DB) *gorm.DB,
	) (results []T, count int64, err apperror.Error)

	Get(
		ctx context.Context,
		filter map[string]interface{},
		scopes ...func(db *gorm.DB) *gorm.DB,
	) (result T, err apperror.Error)
}

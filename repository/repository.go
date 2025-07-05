package repository

import (
	"context"
	"gorm.io/gorm"
	"log"
	"main/pkg/apperror"
	"main/pkg/db/postgres"
	"main/util"
	"net/http"
)

type Repository[T any] struct {
	Db *postgres.DbCluster
}

func (r *Repository[T]) GetAll(
	ctx context.Context,
	filter map[string]interface{},
	scopes ...func(db *gorm.DB) *gorm.DB,
) (results []T, err apperror.Error) {
	logTag := util.LogPrefix(ctx, "Repository.GetAll")

	tx := r.Db.GetSlaveDB(ctx).Model(&results).Where(filter).Scopes(scopes...).Find(&results)
	if tx.Error != nil {
		log.Println(logTag, "Error while fetching records:", tx.Error)

		return nil, apperror.New(tx.Error, http.StatusBadRequest)
	}

	return results, apperror.Error{}
}

func (r *Repository[T]) GetAllWithPagination(
	ctx context.Context,
	filter map[string]interface{},
	scopes ...func(db *gorm.DB) *gorm.DB,
) (results []T, count int64, err apperror.Error) {
	logTag := util.LogPrefix(ctx, "Repository.GetAllWithPagination")

	db := r.Db.GetSlaveDB(ctx).Model(&results).Where(filter).Scopes(scopes...)

	if tx := db.Count(&count); tx.Error != nil {
		log.Println(logTag, "Error counting records:", tx.Error)

		return nil, 0, apperror.New(tx.Error, http.StatusBadRequest)
	}
	if count == 0 {
		log.Println(logTag, "No records found")
		return nil, 0, apperror.Error{}
	}

	if tx := db.Find(&results); tx.Error != nil {
		log.Println(logTag, "Error fetching paginated records:", tx.Error)

		return nil, 0, apperror.New(tx.Error, http.StatusBadRequest)
	}

	return results, count, apperror.Error{}
}

func (r *Repository[T]) Get(
	ctx context.Context,
	filter map[string]interface{},
	scopes ...func(db *gorm.DB) *gorm.DB,
) (result T, err apperror.Error) {
	logTag := util.LogPrefix(ctx, "Repository.Get")

	tx := r.Db.GetSlaveDB(ctx).Model(&result).Where(filter).Scopes(scopes...).First(&result)
	if tx.Error != nil {
		log.Println(logTag, "Error fetching record:", tx.Error)

		return result, apperror.New(tx.Error, http.StatusNotFound)
	}

	return result, apperror.Error{}
}

func (r *Repository[T]) Create(ctx context.Context, data *T) apperror.Error {
	logTag := util.LogPrefix(ctx, "Repository.Create")

	tx := r.Db.GetMasterDB(ctx).Model(data).Create(data)
	if tx.Error != nil {
		log.Println(logTag, "Error creating record:", tx.Error)

		return apperror.New(tx.Error, http.StatusBadRequest)
	}

	return apperror.Error{}
}

func (r *Repository[T]) CreateMany(ctx context.Context, data []*T) apperror.Error {
	logTag := util.LogPrefix(ctx, "Repository.CreateMany")

	tx := r.Db.GetMasterDB(ctx).Model(data).CreateInBatches(data, 1500)
	if tx.Error != nil {
		log.Println(logTag, "Error creating records in bulk:", tx.Error)

		return apperror.New(tx.Error, http.StatusBadRequest)
	}

	return apperror.Error{}
}

func (r *Repository[T]) Update(
	ctx context.Context,
	filter map[string]interface{},
	updates map[string]interface{},
) apperror.Error {
	logTag := util.LogPrefix(ctx, "Repository.Update")

	tx := r.Db.GetMasterDB(ctx).Model(new(T)).Where(filter).Updates(updates)
	if tx.Error != nil {
		log.Println(logTag, "Error updating record:", tx.Error)

		return apperror.New(tx.Error, http.StatusBadRequest)
	}

	return apperror.Error{}
}

func (r *Repository[T]) UpdateMany(ctx context.Context, data []T) apperror.Error {
	logTag := util.LogPrefix(ctx, "Repository.UpdateMany")

	for _, item := range data {
		tx := r.Db.GetMasterDB(ctx).Save(&item)
		if tx.Error != nil {
			log.Println(logTag, "Error updating item in bulk:", tx.Error)

			return apperror.New(tx.Error, http.StatusBadRequest)
		}
	}

	return apperror.Error{}
}

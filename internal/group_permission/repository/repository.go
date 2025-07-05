package repository

import (
	"main/internal/model"
	"main/pkg/db/postgres"
	"main/repository"
	"sync"
)

type Repository struct {
	Interface
	db *postgres.DbCluster
}

var (
	syncOnce sync.Once
	repo     *Repository
)

func NewRepository(db *postgres.DbCluster) *Repository {
	syncOnce.Do(func() {
		repo = &Repository{
			Interface: &repository.Repository[model.GroupUserPermission]{Db: db},
			db:        db,
		}
	})

	return repo
}

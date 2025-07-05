package repository

import (
	"main/internal/model"
	"main/pkg/db/postgres"
	"main/repository"
	"sync"
)

type Repository struct {
	Interface
}

var (
	syncOnce sync.Once
	repo     *Repository
)

func NewRepository(db *postgres.DbCluster) *Repository {
	syncOnce.Do(func() {
		repo = &Repository{&repository.Repository[model.User]{Db: db}}
	})

	return repo
}

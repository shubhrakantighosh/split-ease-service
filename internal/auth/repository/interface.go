package repository

import (
	"main/internal/model"
	"main/repository"
)

type Interface interface {
	repository.Interface[model.AuthToken]
}

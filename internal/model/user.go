package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint64
	Name      string
	Email     string
	Password  string
	IsActive  bool
	CreatedBy string
	UpdatedBy string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

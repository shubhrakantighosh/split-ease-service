package model

import (
	"gorm.io/gorm"
	"time"
)

type GroupUser struct {
	ID        uint64
	GroupID   uint64
	UserID    uint64
	IsActive  bool
	CreatedBy string
	UpdatedBy string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

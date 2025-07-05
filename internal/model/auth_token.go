package model

import (
	"gorm.io/gorm"
	"time"
)

type AuthToken struct {
	ID           uint64
	UserID       uint64
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

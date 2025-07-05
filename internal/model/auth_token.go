package model

import (
	"gorm.io/gorm"
	"time"
)

type AuthToken struct {
	ID               uint64
	UserID           uint64
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt
}

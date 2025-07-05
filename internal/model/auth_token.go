package model

import (
	"gorm.io/gorm"
	"time"
)

type AuthToken struct {
	ID               uint64         `json:"id"`
	UserID           uint64         `json:"user_id"`
	AccessToken      string         `json:"access_token"`
	RefreshToken     string         `json:"refresh_token"`
	AccessExpiresAt  time.Time      `json:"access_expires_at"`
	RefreshExpiresAt time.Time      `json:"refresh_expires_at"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at"`
}

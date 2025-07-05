package model

import (
	"gorm.io/gorm"
	"time"
)

type Purpose string

const (
	Activation    Purpose = "activation"
	PasswordReset Purpose = "password_reset"
)

type OTP struct {
	ID        uint64
	UserID    uint64
	Code      string
	Purpose   Purpose
	ExpiresAt time.Time
	Used      bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

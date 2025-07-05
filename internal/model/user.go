package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID          uint64
	Name        string
	PhoneNumber string
	Email       string
	Password    string
	CreatedBy   string
	UpdatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

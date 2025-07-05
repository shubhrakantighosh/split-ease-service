package model

import (
	"gorm.io/gorm"
	"time"
)

type Bill struct {
	ID             uint64
	GroupID        uint64
	TotalAmountDue float64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt
}

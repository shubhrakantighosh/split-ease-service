package model

import (
	"gorm.io/gorm"
	"time"
)

type BillSplit struct {
	ID        uint64
	BillID    uint64
	AmountDue float64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

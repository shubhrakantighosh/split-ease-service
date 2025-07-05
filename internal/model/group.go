package model

import (
	"gorm.io/gorm"
	"time"
)

type Group struct {
	ID          uint64
	Name        string
	Description string
	CreatedBy   string
	UpdatedBy   string
	CreateTime  time.Time
	UpdateTime  time.Time
	DeleteTime  gorm.DeletedAt
}

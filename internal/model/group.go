package model

import (
	"gorm.io/gorm"
	"time"
)

type Group struct {
	ID          uint64         `json:"id"`
	OwnerID     uint64         `json:"owner_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	CreatedBy   string         `json:"created_by"`
	UpdatedBy   string         `json:"updated_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

type Groups []Group

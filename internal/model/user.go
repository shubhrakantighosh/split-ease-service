package model

import (
	"gorm.io/gorm"
	"main/util"
	"time"
)

type User struct {
	ID        uint64         `json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	IsActive  bool           `json:"is_active"`
	CreatedBy string         `json:"created_by"`
	UpdatedBy string         `json:"updated_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type Users []User

func (u Users) MapByID() map[uint64]User {
	idMapUser := make(map[uint64]User)
	u = util.DeduplicateSlice(u)
	if u == nil || len(u) == 0 {
		return idMapUser
	}

	for _, user := range u {
		idMapUser[user.ID] = user
	}

	return idMapUser
}

package model

import (
	"gorm.io/gorm"
	"time"
)

type Bill struct {
	ID          uint64         `json:"id"`
	UserID      uint64         `json:"user_id"`
	GroupID     uint64         `json:"group_id"`
	PaidAmount  float64        `json:"paid_amount"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

type Bills []Bill

func (b Bills) ExtractUniqueUserIDs() []uint64 {
	uniqueUserIDs := make([]uint64, 0)
	uniqueUserIDMap := make(map[uint64]struct{})

	for _, bill := range b {
		if _, exists := uniqueUserIDMap[bill.UserID]; !exists {
			uniqueUserIDs = append(uniqueUserIDs, bill.UserID)
			uniqueUserIDMap[bill.UserID] = struct{}{}
		}
	}

	return uniqueUserIDs
}

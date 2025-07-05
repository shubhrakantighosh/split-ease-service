package model

import (
	"gorm.io/gorm"
	"time"
)

type BillSplit struct {
	ID          uint64         `json:"id"`
	GroupID     uint64         `json:"group_id"`
	ToPayUserID uint64         `json:"to_pay_user_id"`
	UserID      uint64         `json:"user_id"`
	AmountDue   float64        `json:"amount_due"`
	IsPaid      bool           `json:"is_paid"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

type BillSplits []*BillSplit

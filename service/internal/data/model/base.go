package model

import "time"

// Base contains common fields for all models
type Base struct {
	ID uint `gorm:"primarykey;autoIncrement"`
}

type Others struct {
	CreatedAt time.Time `gorm:"column:created_at;type:DATETIME;not null;autoCreateTime;default:CURRENT_TIMESTAMP"`
}

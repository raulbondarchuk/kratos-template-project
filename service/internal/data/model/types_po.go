package model

import "time"

// Types represents the types model
type Types struct {
	Base
	Name string `gorm:"column:name;type:varchar(255);not null;unique"` // Type1/Type2
	Others
	UpdatedAt time.Time `gorm:"column:updated_at;type:DATETIME;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// TableName returns the name of the table for the User model
func (Types) TableName() string {
	return "types"
}

package model

import "time"

// Templates represents the templates model
type Templates struct {
	Base
	TypeID uint   `gorm:"column:type_id;type:uint;not null"` // Type1/Type2
	Type   Types  `gorm:"foreignKey:TypeID"`
	Name   string `gorm:"column:name;type:varchar(255);not null;unique"` // Template 1, Template 2, etc
	Others
	UpdatedAt time.Time `gorm:"column:updated_at;type:DATETIME;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// TableName returns the name of the table for the User model
func (Templates) TableName() string {
	return "templates"
}

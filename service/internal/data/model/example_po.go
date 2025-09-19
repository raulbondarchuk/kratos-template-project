package model

import "time"

// Examples represents the examples model
type Examples struct {
	Base
	TypeExamplesID uint          `gorm:"column:type_examples_id;type:uint;not null"` // Type1/Type2
	TypeExamples   TypesExamples `gorm:"foreignKey:TypeExamplesID"`
	Name           string        `gorm:"column:name;type:varchar(255);not null;unique"` // Template 1, Template 2, etc
	Others
	UpdatedAt time.Time `gorm:"column:updated_at;type:DATETIME;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// TableName returns the name of the table for the User model
func (Examples) TableName() string {
	return "examples"
}

// Types represents the types model
type TypesExamples struct {
	Base
	Name string `gorm:"column:name;type:varchar(255);not null;unique"` // Type1/Type2
	Others
	UpdatedAt time.Time `gorm:"column:updated_at;type:DATETIME;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// TableName returns the name of the table for the User model
func (TypesExamples) TableName() string {
	return "types_examples"
}

package example_biz

import "time"

// Template represents the example business model
type Example struct {
	ID        uint
	Type      TypeExample
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TypeExample struct {
	ID        uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

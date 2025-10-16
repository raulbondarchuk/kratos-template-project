package example_biz

import "time"

type Example struct {
	ID        uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
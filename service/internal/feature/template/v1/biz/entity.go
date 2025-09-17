package template_biz

import "time"

// Template represents the template business model
type Template struct {
	ID        uint
	Type      Type
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Type struct {
	ID        uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

package grm

import "time"

// Model a basic GoLang struct which includes the following fields: ID, CreatedAt, UpdatedAt, DeletedAt
// It may be embedded into your model or you may build your own model without it
//
//	type User struct {
//	  grm.Model
//	}
type Model struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
}

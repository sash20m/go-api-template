package model

import "time"

// Base creates the default model that every other model is based on.
type Base struct {
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

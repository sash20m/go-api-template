package model

import (
	"time"
)


type User struct {
	ID                     string     `json:"id,omitempty" db:"id"`
	FirstName              string     `json:"firstName" db:"first_name"`
	LastName               string     `json:"lastName" db:"last_name"`
	Email                  string     `json:"email" db:"email"`
	Password               string     `json:"-" db:"password"`
	CreatedAt              time.Time  `json:"-" db:"created_at"`
	UpdatedAt              time.Time  `json:"-" db:"updated_at"`
}

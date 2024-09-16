package models

import "time"

type Employee struct {
	ID        int       `json:"id"`
	Username  string    `json:"username" validate:"required"`
	FirstName string    `json:"first_name" validate:"required"`
	LastName  string    `json:"last_name" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

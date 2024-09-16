package models

import "time"

type OrganizationType string

const (
	IE  OrganizationType = "IE"
	LLS OrganizationType = "LLS"
	JSC OrganizationType = "JSC"
)

type Organization struct {
	ID          int              `json:"id"`
	Name        string           `json:"name" validate:"required"`
	Description string           `json:"description" validate:"required"`
	Type        OrganizationType `json:"type" validate:"required"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

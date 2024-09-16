package models

type OrganizationResponsible struct {
	ID              int `json:"id"`
	Organization_ID int `json:"organization_id"`
	User_ID         int `json:"user_id"`
}
